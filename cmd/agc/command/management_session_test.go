package command

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestManagementCookieSessionNeverUsesPublishingCredentials(t *testing.T) {
	calls := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		if r.Method != "POST" || r.URL.Path != "/agc/edge/scs/scs-service/v1/eval" {
			t.Errorf("wrong transport %s %s", r.Method, r.URL)
		}
		if r.Header.Get("Authorization") != "" || r.Header.Get("oauth2Token") != "offline-oauth" || r.Header.Get("projectId") != "offline-project" || r.Header.Get("agcTeamId") != "offline-team" || r.Header.Get("Cookie") != "authInfo=offline-auth;csrfToken=offline-csrf;" || r.Header.Get("X-HD-CSRF") != "offline-csrf" {
			t.Error("missing explicit management-session context")
		}
		var body map[string]string
		json.NewDecoder(r.Body).Decode(&body)
		if body["command"] != "GET offline-key" {
			t.Error("command changed")
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"code":0,"result":"ok"}`))
	}))
	defer srv.Close()
	dir := t.TempDir()
	file := filepath.Join(dir, "session.json")
	b, _ := json.Marshal(map[string]any{"uid": "offline-user", "teamId": "offline-team", "productId": "offline-project", "accessToken": "offline-oauth", "authInfo": "offline-auth", "csrfToken": "offline-csrf", "expiresAt": time.Now().Add(time.Hour), "projectEdgeBaseUrl": srv.URL})
	os.WriteFile(file, b, 0600)
	out, err := execute("cache", "execute-command", "--invoke", "--dry-run=false", "--session-file", file, "--field", "command=GET offline-key")
	if err != nil || calls != 1 {
		t.Fatalf("calls=%d err=%v output=%s", calls, err, out)
	}
	for _, forbidden := range []string{"offline-oauth", "offline-auth", "offline-csrf"} {
		if containsString(out, forbidden) {
			t.Error("credential leaked")
		}
	}
}

func containsString(s, v string) bool {
	for i := 0; i+len(v) <= len(s); i++ {
		if s[i:i+len(v)] == v {
			return true
		}
	}
	return false
}

func TestManagementDryRunDoesNotReadSession(t *testing.T) {
	out, err := execute("cache", "execute-command", "--invoke", "--session-file", "/does-not-exist", "--field", "command=GET offline-key")
	if err != nil {
		t.Fatalf("dry-run loaded authorization: %v %s", err, out)
	}
}
