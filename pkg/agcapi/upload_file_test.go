package agcapi

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/Createitv/agc-cli/pkg/domain"
)

func TestInvokeEndpointStreamsUploadFile(t *testing.T) {
	payload := []byte("actual file\x00contents")
	file := filepath.Join(t.TempDir(), "app.apk")
	if err := os.WriteFile(file, payload, 0600); err != nil {
		t.Fatal(err)
	}
	for _, id := range []string{"upload-file-new", "obbfile-upload"} {
		for _, contentType := range []string{"", "application/vnd.android.package-archive"} {
			t.Run(id+contentType, func(t *testing.T) {
				server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					data, err := io.ReadAll(r.Body)
					if err != nil {
						t.Error(err)
					}
					if string(data) != string(payload) {
						t.Errorf("wrong file body: %q", data)
					}
					if r.ContentLength != int64(len(payload)) {
						t.Errorf("wrong length: %d", r.ContentLength)
					}
					expected := contentType
					if expected == "" {
						expected = "application/octet-stream"
					}
					if r.Header.Get("Content-Type") != expected {
						t.Errorf("wrong content type: %s", r.Header.Get("Content-Type"))
					}
					if r.Header.Get("Authorization") != "" {
						t.Error("AGC token forwarded to upload host")
					}
				}))
				defer server.Close()
				_, err := InvokeEndpoint(context.Background(), server.Client(), InvokeRequest{Endpoint: domain.Endpoint{ID: id, Method: "PUT", Path: "{uploadUrl}"}, Params: map[string]string{"uploadUrl": server.URL}, FilePath: file, AccessToken: "test-token", Headers: map[string]string{"Content-Type": contentType}})
				if err != nil {
					t.Fatal(err)
				}
			})
		}
	}
}

func TestInvokeEndpointUploadValidationAndDryRun(t *testing.T) {
	calls := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { calls++ }))
	defer server.Close()
	file := filepath.Join(t.TempDir(), "app.apk")
	req := InvokeRequest{Endpoint: domain.Endpoint{ID: "upload-file-new", Method: "PUT", Path: "{uploadUrl}"}, Params: map[string]string{"uploadUrl": server.URL}, FilePath: file}
	if _, err := InvokeEndpoint(context.Background(), server.Client(), req); err == nil {
		t.Error("missing file accepted")
	}
	if err := os.WriteFile(file, []byte("apk"), 0600); err != nil {
		t.Fatal(err)
	}
	req.DryRun = true
	if _, err := InvokeEndpoint(context.Background(), server.Client(), req); err != nil {
		t.Fatal(err)
	}
	req.FilePath = file + "missing"
	if _, err := InvokeEndpoint(context.Background(), server.Client(), req); err == nil {
		t.Error("dry run missing file accepted")
	}
	if calls != 0 {
		t.Errorf("validation/dryrun made %d requests", calls)
	}
}

func TestInvokeEndpointRejectsUnsupportedFileUploads(t *testing.T) {
	file := filepath.Join(t.TempDir(), "app.apk")
	if err := os.WriteFile(file, []byte("apk"), 0600); err != nil {
		t.Fatal(err)
	}
	for _, req := range []InvokeRequest{
		{Endpoint: domain.Endpoint{ID: "other", Method: "POST", Path: "/test"}, BaseURL: "http://127.0.0.1:1", FilePath: file},
		{Endpoint: domain.Endpoint{ID: "upload-file-new", Method: "PUT", Path: "{uploadUrl}"}, Params: map[string]string{"uploadUrl": "http://127.0.0.1:1"}, FilePath: file, Body: []byte(`{}`)},
		{Endpoint: domain.Endpoint{ID: "upload-file-new", Method: "PUT", Path: "{uploadUrl}"}, Params: map[string]string{"uploadUrl": "http://127.0.0.1:1"}, FilePath: filepath.Dir(file)},
	} {
		req.DryRun = true
		if _, err := InvokeEndpoint(context.Background(), nil, req); err == nil {
			t.Fatal("invalid file request accepted")
		}
	}
}
