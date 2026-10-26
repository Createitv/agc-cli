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

func TestPMSResourceSendsDocumentedMultipartFile(t *testing.T) {
	file := filepath.Join(t.TempDir(), "素材.png")
	if err := os.WriteFile(file, []byte("offline-image-bytes"), 0600); err != nil {
		t.Fatal(err)
	}
	called := false
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		if r.Header.Get("appId") != "dedicated" {
			t.Error("missing appId header")
		}
		if err := r.ParseMultipartForm(1024); err != nil {
			t.Errorf("not multipart: %v", err)
			return
		}
		defer r.MultipartForm.RemoveAll()
		for key, want := range map[string]string{"requestId": "fixture", "fileType": "1", "zoneId": "1", "fileName": "素材.png"} {
			if r.FormValue(key) != want {
				t.Errorf("wrong %s", key)
			}
		}
		f, h, err := r.FormFile("file")
		if err != nil {
			t.Error(err)
			return
		}
		defer f.Close()
		data, _ := io.ReadAll(f)
		if string(data) != "offline-image-bytes" || h.Filename != "素材.png" {
			t.Error("file bytes or filename changed")
		}
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{"error":{"errorCode":0},"resourceId":"fixture"}`)
	}))
	defer server.Close()
	e, _ := domain.EndpointByID("pms", "resourcemanage-harmonyosnext")
	_, err := InvokeEndpoint(context.Background(), server.Client(), InvokeRequest{Endpoint: e, BaseURL: server.URL, FilePath: file, Headers: map[string]string{"appId": "dedicated"}, Body: []byte(`{"requestId":"fixture","fileType":"1","zoneId":"1","fileName":"素材.png"}`)})
	if err != nil {
		t.Fatal(err)
	}
	if !called {
		t.Fatal("multipart request not sent")
	}
}
