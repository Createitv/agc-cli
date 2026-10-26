package command

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
)

func TestTeamListRejectsBearerOnlyBeforeNetwork(t *testing.T) {
	var calls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { calls.Add(1); w.Write([]byte(`{"ret":{"code":0}}`)) }))
	defer server.Close()
	out, err := execute("--project", t.TempDir(), "projects", "getteamlist", "--invoke", "--dry-run=false", "--token", "offline-bearer", "--base-url", server.URL)
	if err == nil || calls.Load() != 0 {
		t.Fatalf("OAuth-only API accepted Bearer or contacted network: %v, %d", err, calls.Load())
	}
	if !strings.Contains(out, "oauth2Token") {
		t.Fatalf("missing OAuth correction: %s", out)
	}
}

func TestTeamListUsesDocumentedOAuthHeader(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("oauth2Token") != "offline-oauth" {
			t.Error("OAuth header absent")
		}
		if r.Header.Get("Authorization") != "" {
			t.Error("unexpected Bearer auth")
		}
		if r.URL.Path != "/api/ups/user-permission-service/v1/user-team-list" {
			t.Error("wrong official team path")
		}
		w.Write([]byte(`{"ret":{"code":0},"teams":[]}`))
	}))
	defer server.Close()
	_, err := execute("--project", t.TempDir(), "projects", "getteamlist", "--invoke", "--dry-run=false", "--header", "oauth2Token=offline-oauth", "--base-url", server.URL)
	if err != nil {
		t.Fatal(err)
	}
}
