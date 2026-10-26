package agcapi

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/Createitv/agc-cli/internal/managementtest"
	"github.com/Createitv/agc-cli/pkg/domain"
)

func TestManagementErrorsRedactRequestPasswordsAndSessionValues(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		w.Write([]byte(`{"code":123,"message":"Offline9!fixture offline-auth offline-csrf offline-oauth","password":"Offline9!fixture"}`))
	}))
	defer srv.Close()
	e, _ := domain.EndpointByID("cache", "reset-password")
	out, err := InvokeEndpoint(context.Background(), nil, InvokeRequest{Endpoint: e, SessionFile: managementtest.SessionFile(t, srv.URL), Body: []byte(`{"password":"Offline9!fixture"}`)})
	if err == nil {
		t.Fatal("upstream failure hidden")
	}
	for _, secret := range []string{"Offline9!fixture", "offline-auth", "offline-csrf", "offline-oauth"} {
		if strings.Contains(err.Error()+string(out.Body), secret) {
			t.Fatal("credential exposed in management failure feedback")
		}
	}
}

func TestManagementSessionIsolationRejectsInvalidContextBeforeNetwork(t *testing.T) {
	e, _ := domain.EndpointByID("cache", "connection-info")
	for _, change := range []struct {
		name   string
		update func(map[string]any)
	}{
		{"expired", func(s map[string]any) { s["expiresAt"] = time.Now().Add(-time.Minute) }},
		{"missing-cookie", func(s map[string]any) { delete(s, "authInfo") }},
		{"header-injection", func(s map[string]any) { s["csrfToken"] = "unsafe\nvalue" }},
		{"non-https", func(s map[string]any) { s["projectEdgeBaseUrl"] = "http://example.invalid" }},
	} {
		t.Run(change.name, func(t *testing.T) {
			called := false
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { called = true }))
			defer srv.Close()
			path := managementtest.SessionFile(t, srv.URL)
			b, _ := os.ReadFile(path)
			var s map[string]any
			json.Unmarshal(b, &s)
			change.update(s)
			b, _ = json.Marshal(s)
			os.WriteFile(path, b, 0600)
			_, err := InvokeEndpoint(context.Background(), nil, InvokeRequest{Endpoint: e, SessionFile: path})
			var auth *AuthenticationError
			if called || !errors.As(err, &auth) {
				t.Fatalf("called=%v error=%v", called, err)
			}
		})
	}
}
