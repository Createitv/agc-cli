package agcapi

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Createitv/agc-cli/internal/managementtest"
	"github.com/Createitv/agc-cli/pkg/domain"
)

func TestTemplateDeployRejectsDuplicateComponentsDuringDryRun(t *testing.T) {
	e, _ := domain.EndpointByID("templates", "deploy")
	server, body := managementtest.Deployment(t)
	defer server.Close()
	components := body["components"].([]any)
	body["components"] = append(components, components[0])
	b, _ := json.Marshal(body)
	_, err := InvokeEndpoint(context.Background(), nil, InvokeRequest{Endpoint: e, Body: b, DryRun: true, SessionFile: "/must-not-read"})
	if err == nil || !strings.Contains(err.Error(), "unique") {
		t.Fatalf("duplicate component accepted in local preflight: %v", err)
	}
}

func TestTemplateDeployRejectsCorruptZIPPayloadDuringDryRun(t *testing.T) {
	e, _ := domain.EndpointByID("templates", "deploy")
	server, body := managementtest.Deployment(t)
	defer server.Close()
	path := body["components"].([]any)[0].(map[string]any)["properties"].(map[string]any)["codeUri"].(string)
	// A valid central directory does not establish that the entry bytes match CRC.
	f, _ := os.Create(path)
	w := zip.NewWriter(f)
	entry, _ := w.CreateHeader(&zip.FileHeader{Name: "handler.js", Method: zip.Store})
	entry.Write([]byte("offline-payload"))
	w.Close()
	f.Close()
	data, _ := os.ReadFile(path)
	offset := bytes.Index(data, []byte("offline-payload"))
	if offset < 0 {
		t.Fatal("fixture payload missing")
	}
	data[offset] ^= 1
	os.WriteFile(path, data, 0600)
	b, _ := json.Marshal(body)
	_, err := InvokeEndpoint(context.Background(), nil, InvokeRequest{Endpoint: e, Body: b, DryRun: true, SessionFile: "/must-not-read"})
	if err == nil || !strings.Contains(err.Error(), "ZIP") {
		t.Fatalf("corrupt ZIP accepted: %v", err)
	}
}

func TestTemplateDeployUploadsZIPAndReadsIndependentReadiness(t *testing.T) {
	e, ok := domain.EndpointByID("templates", "deploy")
	if !ok {
		t.Fatal("SAM deployment workflow absent")
	}
	file := filepath.Join(t.TempDir(), "function.zip")
	f, _ := os.Create(file)
	z := zip.NewWriter(f)
	entry, _ := z.Create("handler.js")
	entry.Write([]byte("exports.myHandler=()=>1"))
	z.Close()
	f.Close()
	want, _ := os.ReadFile(file)
	fixture := managementtest.Fixtures()["templates/create-application"].Body
	fixture["storageInstance"] = "offline-bucket"
	fixture["waitForReady"] = true
	fixture["components"].([]any)[0].(map[string]any)["properties"].(map[string]any)["codeUri"] = file
	created, uploaded := false, false
	calls := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		w.Header().Set("Content-Type", "application/json")
		if r.Header.Get("oauth2Token") != "offline-oauth" || r.Header.Get("Authorization") != "" {
			t.Error("wrong deployment auth")
		}
		switch {
		case r.Method == "GET" && r.URL.Path == "/api/adm/v1/applications/search":
			if r.URL.Query().Get("appName") != "offline-app" {
				t.Error("wrong app selection")
			}
			if created {
				io.WriteString(w, `[{"id":"application:offline","name":"offline-app","status":{"ready":true}}]`)
			} else {
				io.WriteString(w, `[]`)
			}
		case r.Method == "PUT" && strings.HasPrefix(r.URL.Path, "/v0/offline-bucket/"):
			b, _ := io.ReadAll(r.Body)
			if string(b) != string(want) || r.ContentLength != int64(len(want)) || r.Header.Get("Content-Type") != "application/x-www-form-urlencoded" {
				t.Error("ZIP transport changed")
			}
			uploaded = true
			io.WriteString(w, `{}`)
		case r.Method == "POST" && r.URL.Path == "/api/adm/v1/applications":
			if !uploaded {
				t.Error("application created before upload")
			}
			var b map[string]any
			json.NewDecoder(r.Body).Decode(&b)
			if _, ok := b["storageInstance"]; ok {
				t.Error("local-only input leaked upstream")
			}
			c := b["components"].([]any)[0].(map[string]any)
			if c["name"] != "offline-app-test-func" || !strings.HasPrefix(c["properties"].(map[string]any)["codeUri"].(string), "offline-bucket/") {
				t.Error("SAM component prefix/codeUri changed")
			}
			created = true
			io.WriteString(w, `{"id":"application:offline"}`)
		default:
			t.Errorf("unexpected %s %s", r.Method, r.URL)
			http.Error(w, "unexpected", 500)
		}
	}))
	defer srv.Close()
	session := managementtest.SessionFile(t, srv.URL)
	var s map[string]any
	b, _ := os.ReadFile(session)
	json.Unmarshal(b, &s)
	s["projectStorageBaseUrl"] = srv.URL
	b, _ = json.Marshal(s)
	os.WriteFile(session, b, 0600)
	b, _ = json.Marshal(fixture)
	result, err := InvokeEndpoint(context.Background(), nil, InvokeRequest{Endpoint: e, Body: b, SessionFile: session})
	if err != nil || calls != 4 || !strings.Contains(string(result.Body), `"ready":true`) {
		t.Fatalf("calls=%d error=%v body=%s", calls, err, result.Body)
	}
}

func TestTemplateDeployValidatesAllArtifactsBeforeReadingSession(t *testing.T) {
	e, ok := domain.EndpointByID("templates", "deploy")
	if !ok {
		t.Fatal("deploy absent")
	}
	fixture := managementtest.Fixtures()["templates/create-application"].Body
	fixture["storageInstance"] = "offline-bucket"
	fixture["components"].([]any)[0].(map[string]any)["properties"].(map[string]any)["codeUri"] = "/missing.zip"
	b, _ := json.Marshal(fixture)
	_, err := InvokeEndpoint(context.Background(), nil, InvokeRequest{Endpoint: e, Body: b, SessionFile: "/must-not-read"})
	if err == nil || !strings.Contains(err.Error(), "ZIP") {
		t.Fatalf("expected artifact error before auth: %v", err)
	}
}
