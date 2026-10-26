// Package cloudfunctions implements the native Go workflows from the Huawei
// AGCLI source. OAuth session files are explicit; publishing credentials are never read.
package cloudfunctions

import (
	"archive/zip"
	"bytes"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode/utf16"

	"github.com/Createitv/agc-cli/pkg/agcapi"
	"github.com/Createitv/agc-cli/pkg/domain"
)

type Session struct {
	UID            string    `json:"uid"`
	TeamID         string    `json:"teamId"`
	ProductID      string    `json:"productId"`
	AccessToken    string    `json:"accessToken"`
	ExpiresAt      time.Time `json:"expiresAt"`
	ProjectBaseURL string    `json:"projectBaseUrl"`
}
type Client struct {
	HTTP    *http.Client
	Session Session
}
type Response struct {
	Operation       string          `json:"operation"`
	DryRun          bool            `json:"dryRun"`
	StatusCode      int             `json:"statusCode,omitempty"`
	ExecutionStatus string          `json:"executionStatus,omitempty"`
	ExecutedVersion string          `json:"executedVersion,omitempty"`
	TraceID         string          `json:"traceId,omitempty"`
	Updated         bool            `json:"updated,omitempty"`
	Published       bool            `json:"published,omitempty"`
	Body            json.RawMessage `json:"body,omitempty"`
}

func violation(field, msg string) error {
	return &agcapi.RequestValidationError{Violations: []domain.InputViolation{{Code: "function_input", Field: field, Message: msg, SourceURL: "https://www.npmjs.com/package/@hw-agconnect/agc-cli/v/1.12.0", Suggestion: "Read agc functions <operation> --docs --pretty; correct this field before invoking."}}}
}
func Validate(operation string, r domain.FunctionRequest) error {
	known := false
	for _, op := range domain.FunctionOperations {
		known = known || operation == op
	}
	if !known {
		return violation("operation", "Unknown native function operation")
	}
	if operation != "list-functions" && !regexp.MustCompile(`^[a-z][a-z0-9-]{0,30}[a-z0-9]$`).MatchString(r.Name) {
		return violation("name", "Function name must have 2-32 lowercase letters, digits or hyphens, starting with a letter and ending with a letter or digit.")
	}
	if r.Version != "" && r.Version != "$latest" {
		n, err := strconv.Atoi(r.Version)
		if err != nil || n < 1 || n > 1000 {
			return violation("version", "Use $latest or a version integer from 1 to 1000 (AGCLI option contract).")
		}
	}
	if r.Version != "" && r.AliasName != "" {
		return violation("aliasName", "version and aliasName are mutually exclusive.")
	}
	if r.AliasName != "" && !regexp.MustCompile(`^[a-z][a-z0-9-]{0,14}[a-z0-9]$`).MatchString(r.AliasName) {
		return violation("aliasName", "Alias must contain 2-16 lowercase letters, digits or hyphens.")
	}
	if len(r.Payload) > 0 && (!json.Valid(r.Payload) || len(r.Payload) >= 4<<20) {
		return violation("payload", "Supply valid JSON smaller than 4 MiB.")
	}
	if r.Memory != nil {
		switch *r.Memory {
		case 500, 1000, 2000, 4000:
		default:
			return violation("memory", "Use 500, 1000, 2000 or 4000 MiB.")
		}
	}
	if r.Timeout != nil && (*r.Timeout < 1 || *r.Timeout > 1800) {
		return violation("timeout", "Use an integer from 1 to 1800 seconds.")
	}
	if r.InstanceConcurrency != nil && (*r.InstanceConcurrency < 1 || *r.InstanceConcurrency > 1000) {
		return violation("instanceConcurrency", "Use an integer from 1 to 1000.")
	}
	if r.Runtime != "" && r.Runtime != "nodejs14.x" && r.Runtime != "python3" && r.Runtime != "java1.8" && r.Runtime != "custom-runtime" {
		return violation("runtime", "Use a runtime supported by the reviewed AGCLI contract.")
	}
	if agcliWidth(r.Description) > 64 || agcliWidth(r.VersionDescription) > 64 {
		return violation("versionDescription", "Descriptions must contain at most 64 characters.")
	}
	if strings.Contains(operation, "publish") && r.VersionDescription == "" {
		return violation("versionDescription", "A version description is required for publish operations.")
	}
	if r.Handler != nil && (*r.Handler == "" || len(*r.Handler) > 64) {
		return violation("handler", "Supply a nonempty handler of at most 64 characters.")
	}
	if r.Environment != nil {
		raw, _ := json.Marshal(*r.Environment)
		if len(raw) >= 3<<20 {
			return violation("environment", "Environment must be smaller than 3 MiB.")
		}
		for k, v := range *r.Environment {
			if len(k) > 64 || !regexp.MustCompile(`^[-._a-zA-Z][-._a-zA-Z0-9]*$`).MatchString(k) || len(utf16.Encode([]rune(v))) > 2048 {
				return violation("environment", "Environment keys must match the AGCLI name format (max 64); values have max length 2048.")
			}
		}
	}
	if r.Layers != nil {
		if len(*r.Layers) > 5 {
			return violation("layers", "At most five layers may be attached.")
		}
		for _, l := range *r.Layers {
			if !strings.Contains(l, ":") {
				return violation("layers", "Use layerName:version or layerUrn:version.")
			}
		}
	}
	if r.PageSize < 0 || r.PageSize > 1000 || r.PageIndex < 0 {
		return violation("pageSize", "Use pageSize 1-1000 and pageIndex >=1, or omit defaults.")
	}
	if len(r.FuzzyName) > 32 || r.FuzzyName != "" && !regexp.MustCompile(`^[a-z0-9-]+$`).MatchString(r.FuzzyName) {
		return violation("fuzzyName", "Use at most 32 lowercase letters, digits or hyphens.")
	}
	if (operation == "create-function" || operation == "update-function-code") && r.CodePath == "" {
		return violation("codePath", "A ZIP deployment package is required.")
	}
	if r.CodePath != "" {
		info, err := os.Stat(r.CodePath)
		if err != nil || !info.Mode().IsRegular() || info.Size() >= 50<<20 || strings.ToLower(filepath.Ext(r.CodePath)) != ".zip" {
			return violation("codePath", "Use a readable regular ZIP file smaller than 50 MiB.")
		}
		z, err := zip.OpenReader(r.CodePath)
		if err != nil {
			return violation("codePath", "Deployment package must be a valid ZIP archive.")
		}
		z.Close()
	}
	if operation == "update-function-config" && !hasUpdate(r) && r.VersionDescription == "" {
		return violation("body", "Supply at least one configuration change.")
	}
	return nil
}
func hasUpdate(r domain.FunctionRequest) bool {
	return r.CodePath != "" || r.Handler != nil || r.Memory != nil || r.Timeout != nil || r.InstanceConcurrency != nil || r.Environment != nil || r.Layers != nil
}
func (c *Client) Execute(ctx context.Context, op string, r domain.FunctionRequest) (Response, error) {
	out := Response{Operation: op, DryRun: r.DryRun}
	if err := Validate(op, r); err != nil {
		return out, err
	}
	if r.DryRun {
		return out, nil
	}
	session := c.Session
	if r.SessionFile != "" {
		b, err := os.ReadFile(r.SessionFile)
		if err != nil {
			return out, violation("sessionFile", "Cannot read the explicitly selected OAuth session file.")
		}
		if json.Unmarshal(b, &session) != nil {
			return out, violation("sessionFile", "Invalid OAuth session file.")
		}
	}
	if session.UID == "" || session.TeamID == "" || session.ProductID == "" || session.AccessToken == "" || !session.ExpiresAt.After(time.Now()) {
		return out, violation("sessionFile", "An explicit, unexpired user OAuth session with uid, teamId, productId, accessToken and expiresAt is required.")
	}
	if r.ProjectID != "" && r.ProjectID != session.ProductID {
		return out, violation("projectId", "Selected project does not match the OAuth session productId.")
	}
	base := r.BaseURL
	if base == "" {
		base = session.ProjectBaseURL
	}
	u, err := url.Parse(base)
	if err != nil || u.Host == "" || u.User != nil || u.RawQuery != "" || u.Fragment != "" || (u.Scheme != "https" && !(u.Scheme == "http" && (u.Hostname() == "127.0.0.1" || u.Hostname() == "localhost"))) {
		return out, violation("baseUrl", "Supply the project's HTTPS base URL; only loopback HTTP is allowed for isolated tests.")
	}
	client := c.HTTP
	if client == nil {
		client = &http.Client{Timeout: 60 * time.Second}
	}
	clone := *client
	clone.CheckRedirect = func(next *http.Request, via []*http.Request) error {
		if len(via) > 0 && (next.URL.Host != via[0].URL.Host || next.URL.Scheme != via[0].URL.Scheme) {
			return fmt.Errorf("cross-origin OAuth redirect refused")
		}
		if len(via) >= 10 {
			return fmt.Errorf("too many redirects")
		}
		return nil
	}
	client = &clone
	prefix := strings.TrimRight(base, "/") + "/api/wisefunction/serverless/v1/functions"
	named := prefix + "/" + url.PathEscape(r.Name)
	request := func(method, target, ctype string, body io.Reader, headers map[string]string) (json.RawMessage, error) {
		h, err := http.NewRequestWithContext(ctx, method, target, body)
		if err != nil {
			return nil, err
		}
		if file, ok := body.(*os.File); ok {
			if info, statErr := file.Stat(); statErr == nil {
				h.ContentLength = info.Size()
			}
		}
		h.Header.Set("uid", session.UID)
		h.Header.Set("teamId", session.TeamID)
		h.Header.Set("productId", session.ProductID)
		h.Header.Set("oauth2Token", session.AccessToken)
		h.Header.Set("X-HD-DATE", time.Now().UTC().Format("20060102T150405Z"))
		h.Header.Set("Accept", "application/json")
		h.Header.Set("User-Agent", "agc-cli/native-functions")
		trace := make([]byte, 16)
		if _, err = rand.Read(trace); err != nil {
			return nil, err
		}
		h.Header.Set("x-trace-id", hex.EncodeToString(trace))
		if ctype != "" {
			h.Header.Set("Content-Type", ctype)
		}
		for k, v := range headers {
			h.Header.Set(k, v)
		}
		resp, err := client.Do(h)
		if err != nil {
			return nil, fmt.Errorf("native function transport failed: %w", err)
		}
		defer resp.Body.Close()
		out.StatusCode = resp.StatusCode
		out.TraceID = h.Header.Get("x-trace-id")
		out.ExecutionStatus = resp.Header.Get("x-function-execution-status")
		out.ExecutedVersion = resp.Header.Get("x-executed-version")
		data, err := io.ReadAll(io.LimitReader(resp.Body, (16<<20)+1))
		if err != nil {
			return nil, err
		}
		if len(data) > 16<<20 {
			return nil, fmt.Errorf("native function response exceeds the 16 MiB local implementation limit; refusing truncated success")
		}
		if resp.StatusCode >= 400 {
			return nil, agcapi.ParseAGCError(resp.StatusCode, data)
		}
		var obj map[string]json.RawMessage
		if op != "invoke" && json.Unmarshal(data, &obj) == nil {
			for _, key := range []string{"code", "error_code"} {
				var code any
				if json.Unmarshal(obj[key], &code) == nil && fmt.Sprint(code) != "0" && fmt.Sprint(code) != "<nil>" {
					return nil, agcapi.ParseAGCError(resp.StatusCode, data)
				}
			}
		}
		if out.ExecutionStatus != "" && !strings.EqualFold(out.ExecutionStatus, "success") && out.ExecutionStatus != "0" {
			return nil, fmt.Errorf("function execution failed (HTTP %d, execution status %s)", resp.StatusCode, out.ExecutionStatus)
		}
		return data, nil
	}
	getLatest := func() (map[string]json.RawMessage, error) {
		data, err := request("GET", named+"?qualifier=%24latest", "", nil, nil)
		if err != nil {
			return nil, err
		}
		var root struct {
			Function map[string]json.RawMessage `json:"function"`
		}
		if json.Unmarshal(data, &root) != nil || root.Function == nil {
			return nil, fmt.Errorf("function query omitted function model")
		}
		var rev string
		if json.Unmarshal(root.Function["revisionId"], &rev) != nil || rev == "" {
			return nil, fmt.Errorf("function query omitted revisionId; refusing write")
		}
		return root.Function, nil
	}
	publish := func(revision json.RawMessage) (json.RawMessage, error) {
		var rev string
		if json.Unmarshal(revision, &rev) != nil || rev == "" {
			return nil, fmt.Errorf("missing revisionId; refusing publish")
		}
		b, _ := json.Marshal(map[string]string{"revisionId": rev, "versionDesc": r.VersionDescription})
		return request("POST", named+"/versions", "application/json", bytes.NewReader(b), nil)
	}
	update := func(prior map[string]json.RawMessage, create bool) (json.RawMessage, error) {
		// A temporary multipart file bounds memory even for deployment ZIPs.
		f, err := os.CreateTemp("", "agc-function-*.multipart")
		if err != nil {
			return nil, err
		}
		defer os.Remove(f.Name())
		defer f.Close()
		w := multipart.NewWriter(f)
		fields := map[string]string{}
		headers := map[string]string{}
		if create {
			fields["name"] = r.Name
			runtime := r.Runtime
			if runtime == "" {
				runtime = "nodejs14.x"
			}
			fields["runtime"] = runtime
			fields["memory"] = "500"
			fields["cpu"] = "500"
			fields["timeout"] = "55"
			fields["instanceConcurrency"] = "100"
			fields["handler"] = "handler.myHandler"
			if runtime == "java1.8" {
				fields["handler"] = "faas.HelloPojo::handleRequest"
			}
			if runtime == "custom-runtime" {
				fields["handler"] = "/invoke"
			}
			if r.Description != "" {
				fields["description"] = r.Description
			}
		} else {
			for _, k := range []string{"revisionId", "environment", "layers", "handler"} {
				raw := prior[k]
				if len(raw) == 0 {
					continue
				}
				if k == "handler" || k == "revisionId" {
					var text string
					if json.Unmarshal(raw, &text) != nil {
						return nil, fmt.Errorf("invalid prior %s model", k)
					}
					fields[k] = text
				} else {
					fields[k] = string(raw)
				}
			}
		}
		if r.Handler != nil {
			fields["handler"] = *r.Handler
		}
		if r.Memory != nil {
			fields["memory"] = strconv.Itoa(*r.Memory)
			fields["cpu"] = strconv.Itoa(map[int]int{500: 500, 1000: 500, 2000: 1000, 4000: 2000}[*r.Memory])
		}
		if r.Timeout != nil {
			fields["timeout"] = strconv.Itoa(*r.Timeout)
		}
		if r.InstanceConcurrency != nil {
			fields["instanceConcurrency"] = strconv.Itoa(*r.InstanceConcurrency)
		}
		if r.VersionDescription != "" {
			fields["versionDesc"] = r.VersionDescription
		}
		if r.Environment != nil {
			b, _ := json.Marshal(*r.Environment)
			fields["environment"] = string(b)
		}
		if r.Layers != nil {
			layers := []string{}
			for _, l := range *r.Layers {

				pieces := strings.Split(l, ":")
				if len(pieces) != 2 && len(pieces) != 7 {
					return nil, violation("layers", "Use layerName:version or a complete layer URN.")
				}
				target := strings.TrimRight(base, "/") + "/api/wisefunction/serverless/v1/layers/" + url.PathEscape(strings.Join(pieces[:len(pieces)-1], ":")) + "/versions/" + url.PathEscape(pieces[len(pieces)-1])
				raw, err := request("GET", target, "", nil, nil)
				if err != nil {
					return nil, err
				}
				var response struct {
					Result struct {
						LayerVersionUrn string `json:"layerVersionUrn"`
						DeletedFlag     bool   `json:"deletedFlag"`
					} `json:"result"`
				}
				if json.Unmarshal(raw, &response) != nil || response.Result.DeletedFlag {
					return nil, fmt.Errorf("layer unavailable or deleted")
				}
				urn := response.Result.LayerVersionUrn
				if urn == "" {
					return nil, fmt.Errorf("layer query omitted layerVersionUrn")
				}
				layers = append(layers, urn)
			}
			b, _ := json.Marshal(layers)
			fields["layers"] = string(b)
		}
		for k, v := range fields {
			if err = w.WriteField(k, v); err != nil {
				return nil, err
			}
		}
		if r.CodePath != "" {
			source, err := os.Open(r.CodePath)
			if err != nil {
				return nil, err
			}
			defer source.Close()
			part, err := w.CreateFormFile("file", filepath.Base(r.CodePath))
			if err != nil {
				return nil, err
			}
			hash := sha256.New()
			n, err := io.Copy(io.MultiWriter(part, hash), source)
			if err != nil {
				return nil, err
			}
			headers["x-file-digest"] = strconv.FormatInt(n, 10)
			headers["x-file-digest-sha256"] = hex.EncodeToString(hash.Sum(nil))
		}
		if err = w.Close(); err != nil {
			return nil, err
		}
		if _, err = f.Seek(0, 0); err != nil {
			return nil, err
		}
		method, target := "PUT", named
		if create {
			method, target = "POST", prefix
		}
		return request(method, target, w.FormDataContentType(), f, headers)
	}
	var data json.RawMessage
	switch op {
	case "create-function":
		data, err = update(nil, true)
	case "invoke":
		qualifier := r.Version
		if qualifier == "" {
			qualifier = "$latest"
		}
		if r.AliasName != "" {
			qualifier = r.AliasName
		}
		payload := r.Payload
		if len(payload) == 0 {
			payload = json.RawMessage(`{}`)
		}
		data, err = request("POST", named+"/invocations?qualifier="+url.QueryEscape(qualifier), "application/json", bytes.NewReader(payload), nil)
	case "list-functions", "list-function-versions":
		size, index := r.PageSize, r.PageIndex
		if size == 0 {
			size = 10
		}
		if index == 0 {
			index = 1
		}
		query := url.Values{"pageSize": {strconv.Itoa(size)}, "pageIndex": {strconv.Itoa(index)}}
		target := named + "/versions"
		if op == "list-functions" {
			target = prefix
			query.Set("qualifier", "$latest")
			query.Set("status", "1")
			if r.FuzzyName != "" {
				query.Set("name", r.FuzzyName)
			}
		}
		data, err = request("GET", target+"?"+query.Encode(), "", nil, nil)
	case "query-function-version":
		v := r.Version
		if v == "" {
			v = "$latest"
		}
		data, err = request("GET", named+"?qualifier="+url.QueryEscape(v), "", nil, nil)
	default:
		var prior map[string]json.RawMessage
		prior, err = getLatest()
		if err != nil {
			break
		}
		if op == "publish-function" {
			data, err = publish(prior["revisionId"])
			out.Published = err == nil
			break
		}
		if op == "update-and-publish-function" && !hasUpdate(r) {
			data, err = publish(prior["revisionId"])
			out.Published = err == nil
			break
		}
		data, err = update(prior, false)
		if err != nil {
			break
		}
		out.Updated = true
		if op == "update-and-publish-function" {
			var changed struct {
				Function map[string]json.RawMessage `json:"function"`
			}
			json.Unmarshal(data, &changed)
			data, err = publish(changed.Function["revisionId"])
			out.Published = err == nil
			if err != nil {
				err = fmt.Errorf("update succeeded; publish failed; do not blindly repeat update: %w", err)
			}
		}
	}
	if err != nil {
		return out, err
	}
	if op == "invoke" {
		if json.Valid(data) {
			out.Body = append(json.RawMessage(nil), data...)
		} else {
			out.Body, _ = json.Marshal(string(data))
		}
	} else {
		out.Body = redacted(data)
	}
	return out, nil
}
func redacted(data []byte) json.RawMessage {
	var obj any
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.UseNumber()
	if dec.Decode(&obj) != nil {
		b, _ := json.Marshal(string(data))
		return b
	}
	var clean func(any)
	clean = func(x any) {
		switch v := x.(type) {
		case map[string]any:
			for k, c := range v {
				switch strings.ToLower(k) {
				case "environment", "accesstoken", "oauth2token", "uid", "teamid":
					delete(v, k)
				default:
					clean(c)
				}
			}
		case []any:
			for _, c := range v {
				clean(c)
			}
		}
	}
	clean(obj)
	b, _ := json.Marshal(obj)
	return b
}

func agcliWidth(s string) int {
	width := 0
	for _, unit := range utf16.Encode([]rune(s)) {
		if unit <= 128 {
			width++
		} else {
			width += 2
		}
	}
	return width
}
