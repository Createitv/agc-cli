package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestConfiguredServerRequiresItsOwnToken(t *testing.T) {
	for _, provided := range []string{"", "wrong", "server-secret"} {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/capabilities", nil)
		req.Header.Set("X-AGC-Server-Token", provided)
		rec := httptest.NewRecorder()
		HandlerWithToken("server-secret").ServeHTTP(rec, req)
		want := http.StatusUnauthorized
		if provided == "server-secret" {
			want = http.StatusOK
		}
		if rec.Code != want {
			t.Fatalf("provided token valid=%v: got %d want %d", provided == "server-secret", rec.Code, want)
		}
	}
}
