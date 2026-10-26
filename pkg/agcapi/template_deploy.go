package agcapi

import (
	"archive/zip"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/Createitv/agc-cli/pkg/domain"
)

func invokeTemplateDeploy(ctx context.Context, client *http.Client, req InvokeRequest) (InvokeResponse, error) {
	out := InvokeResponse{Method: "POST", URL: "native:templates/deploy", DryRun: req.DryRun}
	var body map[string]any
	dec := json.NewDecoder(strings.NewReader(string(req.Body)))
	dec.UseNumber()
	if dec.Decode(&body) != nil {
		return out, uploadIssue(req.Endpoint, "body", "Provide a valid SAM JSON manifest.")
	}
	name, _ := body["name"].(string)
	bucket, _ := body["storageInstance"].(string)
	components, _ := body["components"].([]any)
	if name == "" || bucket == "" || len(components) == 0 {
		return out, uploadIssue(req.Endpoint, "body", "Provide name, storageInstance and one or more CloudFunctions components.")
	}
	type artifact struct {
		Path     string
		SHA256   string
		Bytes    int64
		ObjectID string
	}
	artifacts := make([]artifact, len(components))
	// Every artifact is validated before loading authorization or starting uploads.
	for i, v := range components {
		c, ok := v.(map[string]any)
		if !ok {
			return out, uploadIssue(req.Endpoint, fmt.Sprintf("body.components[%d]", i), "Component must be a CloudFunctions object.")
		}
		props, _ := c["properties"].(map[string]any)
		path, _ := props["codeUri"].(string)
		file, err := os.Open(path)
		if err != nil {
			return out, uploadIssue(req.Endpoint, fmt.Sprintf("body.components[%d].properties.codeUri", i), "Use a readable local ZIP codeUri, smaller than the source SAM limit of 300 MiB.")
		}
		stat, err := file.Stat()
		if err != nil || !stat.Mode().IsRegular() || stat.Size() >= 300<<20 || !strings.EqualFold(filepath.Ext(path), ".zip") {
			file.Close()
			return out, uploadIssue(req.Endpoint, fmt.Sprintf("body.components[%d].properties.codeUri", i), "Use a regular local ZIP smaller than 300 MiB.")
		}
		z, err := zip.NewReader(file, stat.Size())
		if err != nil || len(z.File) == 0 {
			file.Close()
			return out, uploadIssue(req.Endpoint, fmt.Sprintf("body.components[%d].properties.codeUri", i), "Code ZIP must contain at least one entry.")
		}
		for _, entry := range z.File {
			reader, openErr := entry.Open()
			if openErr != nil {
				file.Close()
				return out, uploadIssue(req.Endpoint, fmt.Sprintf("body.components[%d].properties.codeUri", i), "Code ZIP entry cannot be opened.")
			}
			_, readErr := io.Copy(io.Discard, contextZIPReader{ctx: ctx, reader: reader})
			reader.Close()
			if readErr != nil {
				file.Close()
				return out, uploadIssue(req.Endpoint, fmt.Sprintf("body.components[%d].properties.codeUri", i), "Code ZIP payload failed integrity verification or the request context expired.")
			}
		}
		sum := sha256.New()
		if _, err = io.Copy(sum, file); err != nil {
			file.Close()
			return out, uploadIssue(req.Endpoint, "body.components", "Cannot hash code ZIP.")
		}
		file.Close()
		artifacts[i] = artifact{Path: path, SHA256: hex.EncodeToString(sum.Sum(nil)), Bytes: stat.Size()}
		componentName, _ := c["name"].(string)
		combined := name + "-" + componentName
		if len(componentName) > 16 || !regexp.MustCompile(`^[a-z][a-z0-9-]{0,30}[a-z0-9]$`).MatchString(combined) {
			return out, uploadIssue(req.Endpoint, fmt.Sprintf("body.components[%d].name", i), "The SAM application prefix plus component name must fit the reviewed 32-character function name limit.")
		}
		c["name"] = combined
	}
	if req.DryRun {
		return out, nil
	}
	type uploadedArtifact struct {
		ObjectID string `json:"objectId"`
		SHA256   string `json:"sha256"`
		Bytes    int64  `json:"bytes"`
	}
	uploaded := []uploadedArtifact{}
	fail := func(stage string, err error) (InvokeResponse, error) {
		return out, fmt.Errorf("SAM deployment stopped at %s; %d uploaded artifacts may be retained, no cleanup or retry was performed: %w", stage, len(uploaded), err)
	}
	call := func(id string, b map[string]any, q map[string]string) (InvokeResponse, error) {
		e, ok := domain.EndpointByID("templates", id)
		if !ok {
			return InvokeResponse{}, fmt.Errorf("missing SAM transport contract")
		}
		var data []byte
		if b != nil {
			data, _ = json.Marshal(b)
		}
		return InvokeEndpoint(ctx, client, InvokeRequest{Endpoint: e, Body: data, Query: q, SessionFile: req.SessionFile, ProjectID: req.ProjectID, BaseURL: req.BaseURL})
	}
	search := func() (map[string]any, error) {
		response, err := call("search-applications", nil, map[string]string{"appName": name})
		if err != nil {
			return nil, err
		}
		var matches []map[string]any
		if json.Unmarshal(response.Body, &matches) != nil {
			return nil, fmt.Errorf("SAM search omitted its source-defined array response")
		}
		var found map[string]any
		for _, m := range matches {
			if m["name"] == name {
				if found != nil {
					return nil, fmt.Errorf("multiple exact SAM application matches; refusing ambiguous target")
				}
				found = m
			}
		}
		return found, nil
	}
	prior, err := search()
	if err != nil {
		return fail("application preflight", err)
	}
	if prior != nil {
		return fail("application preflight", fmt.Errorf("application already exists; use explicit create-plan/get-plan/apply-plan instead of overwriting it"))
	}
	for i, a := range artifacts {
		random := make([]byte, 16)
		if _, err = rand.Read(random); err != nil {
			return fail("object naming", err)
		}
		object := "agc-cli-" + hex.EncodeToString(random) + ".zip"
		e := domain.TemplateCodeUploadEndpoint()
		_, err := InvokeEndpoint(ctx, client, InvokeRequest{Endpoint: e, SessionFile: req.SessionFile, ProjectID: req.ProjectID, Params: map[string]string{"bucket": bucket, "object_name": object}, FilePath: a.Path, Headers: map[string]string{"X-Agc-Trace-Id": strings.TrimSuffix(object, ".zip")}})
		if err != nil {
			return fail("artifact upload", err)
		}
		key := bucket + "/" + object
		uploaded = append(uploaded, uploadedArtifact{ObjectID: key, SHA256: a.SHA256, Bytes: a.Bytes})
		components[i].(map[string]any)["properties"].(map[string]any)["codeUri"] = key
	}
	wait, _ := body["waitForReady"].(bool)
	delete(body, "storageInstance")
	delete(body, "waitForReady")
	accepted, err := call("create-application", body, nil)
	if err != nil {
		return fail("application creation", err)
	}
	var created map[string]any
	if json.Unmarshal(accepted.Body, &created) != nil || created["id"] == nil {
		return fail("application creation readback", fmt.Errorf("creation response omitted application id; creation may have succeeded"))
	}
	ready := false
	for i := 0; i < 15; i++ {
		current, err := search()
		if err != nil {
			return fail("readiness readback", err)
		}
		if current != nil {
			if current["id"] != created["id"] {
				return fail("readiness readback", fmt.Errorf("application id differs from created resource"))
			}
			status, _ := current["status"].(map[string]any)
			ready, _ = status["ready"].(bool)
		}
		if ready || !wait {
			break
		}
		if i == 14 {
			break
		}
		timer := time.NewTimer(2400 * time.Millisecond)
		select {
		case <-ctx.Done():
			timer.Stop()
			return fail("readiness wait", ctx.Err())
		case <-timer.C:
		}
	}
	out.StatusCode = accepted.StatusCode
	out.ContentType = "application/json"
	out.Body, _ = json.Marshal(map[string]any{"applicationAccepted": true, "applicationId": created["id"], "ready": ready, "uploadedArtifacts": uploaded, "retention": "Uploaded code objects retained; no automatic cleanup. Existing applications were not modified."})
	out.RawBody = out.Body
	return out, nil
}

type contextZIPReader struct {
	ctx    context.Context
	reader io.Reader
}

func (r contextZIPReader) Read(p []byte) (int, error) {
	select {
	case <-r.ctx.Done():
		return 0, r.ctx.Err()
	default:
		return r.reader.Read(p)
	}
}
