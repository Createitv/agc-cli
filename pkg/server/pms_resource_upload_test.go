package server

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRESTPMSResourceForwardsSelectedLocalFileAsMultipart(t *testing.T) {
	file := filepath.Join(t.TempDir(), "fixture.png")
	if err := os.WriteFile(file, []byte("offline-image"), 0600); err != nil {
		t.Fatal(err)
	}
	called := false
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		if err := r.ParseMultipartForm(1024); err != nil {
			t.Error(err)
			return
		}
		defer r.MultipartForm.RemoveAll()
		if r.FormValue("requestId") != "fixture" || r.FormValue("fileType") != "1" {
			t.Error("multipart fields changed")
		}
		part, _, err := r.FormFile("file")
		if err != nil {
			t.Error(err)
			return
		}
		defer part.Close()
		data, _ := io.ReadAll(part)
		if string(data) != "offline-image" {
			t.Error("file bytes changed")
		}
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{"error":{"errorCode":0},"resourceId":"fixture"}`)
	}))
	defer upstream.Close()
	payload, _ := json.Marshal(map[string]any{"baseUrl": upstream.URL, "filePath": file, "headers": map[string]string{"appId": "dedicated"}, "fields": map[string]string{"requestId": "fixture", "fileType": "1", "zoneId": "1", "fileName": "fixture.png"}, "dryRun": false})
	rec := httptest.NewRecorder()
	Handler().ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/api/v1/pms/endpoints/resourcemanage-harmonyosnext/invoke", bytes.NewReader(payload)))
	if rec.Code != http.StatusOK || !called {
		t.Fatalf("REST upload failed: %d %s", rec.Code, rec.Body.String())
	}
}

func TestRESTMissingResourceFileReturnsAgentCorrection(t *testing.T) {
	rec := httptest.NewRecorder()
	Handler().ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/api/v1/pms/endpoints/resourcemanage-harmonyosnext/invoke", strings.NewReader(`{"headers":{"appId":"dedicated"},"fields":{"requestId":"fixture","fileType":"1","zoneId":"1","fileName":"fixture.png"},"dryRun":true}`)))
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("missing file classified incorrectly: %d %s", rec.Code, rec.Body.String())
	}
	var response struct {
		Error struct {
			Code, Endpoint, DocumentationCommand string
			Violations                           []struct{ Field, SourceURL, Suggestion string }
		}
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Fatal(err)
	}
	if response.Error.Code != "invalid_request" || response.Error.Endpoint != "pms/resourcemanage-harmonyosnext" || response.Error.DocumentationCommand == "" || len(response.Error.Violations) != 1 {
		t.Fatalf("incomplete Agent correction: %s", rec.Body.String())
	}
	v := response.Error.Violations[0]
	if v.Field != "filePath" || v.SourceURL == "" || v.Suggestion == "" {
		t.Fatalf("missing file guidance: %+v", v)
	}
}
