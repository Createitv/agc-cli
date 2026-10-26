package command

import (
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
)

func TestAuthWhoAmIShowsActiveIdentityWithoutSecrets(t *testing.T) {
	path := filepath.Join(t.TempDir(), "credentials.json")
	if _, err := execute("auth", "login", "--client-id", "developer-client", "--client-key", "test-secret", "--name", "production", "--credentials-path", path); err != nil {
		t.Fatal(err)
	}
	out, err := execute("auth", "whoami", "--credentials-path", path)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, `"name":"production"`) || !strings.Contains(out, `"clientId":"developer-client"`) || !strings.Contains(out, `"teamLookup":"oauth-required"`) {
		t.Fatalf("output = %s", out)
	}
	if strings.Contains(out, "test-secret") || strings.Contains(out, "clientKey") {
		t.Fatal("whoami output exposes client key")
	}
}

func TestAuthWhoAmIReadsOfficialOAuthTeamList(t *testing.T) {
	t.Setenv("AGC_OAUTH2_TOKEN", "oauth-test-secret")
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/ups/user-permission-service/v1/user-team-list" {
			t.Errorf("path = %s", r.URL.Path)
		}
		if r.Header.Get("oauth2Token") != "oauth-test-secret" {
			t.Error("missing OAuth team-list token")
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"ret":{"code":0,"msg":"success"},"teams":[{"teamId":"team-1","teamName":"Example Studio"}]}`))
	}))
	defer srv.Close()

	out, err := execute("auth", "whoami", "--teams", "--base-url", srv.URL, "--credentials-path", filepath.Join(t.TempDir(), "missing.json"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, `"teamLookup":"verified"`) || !strings.Contains(out, `"teamName":"Example Studio"`) {
		t.Fatalf("output = %s", out)
	}
	if strings.Contains(out, "oauth-test-secret") {
		t.Fatal("whoami output exposes OAuth token")
	}
}

func TestAuthWhoAmITeamsRequiresOAuthToken(t *testing.T) {
	t.Setenv("AGC_OAUTH2_TOKEN", "")
	_, err := execute("auth", "whoami", "--teams", "--credentials-path", filepath.Join(t.TempDir(), "missing.json"))
	if err == nil || !strings.Contains(err.Error(), "AGC_OAUTH2_TOKEN") {
		t.Fatalf("error = %v", err)
	}
}
