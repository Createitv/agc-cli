package agcapi

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestUserTeamList(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/api/ups/user-permission-service/v1/user-team-list" {
			t.Errorf("request = %s %s", r.Method, r.URL.Path)
		}
		if r.Header.Get("oauth2Token") != "oauth-token" {
			t.Error("missing oauth2Token header")
		}
		w.Write([]byte(`{"ret":{"code":0},"teams":[{"teamId":"one","teamName":"Example Studio"}]}`))
	}))
	defer srv.Close()

	result, err := UserTeamList(context.Background(), srv.Client(), srv.URL, "oauth-token")
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Teams) != 1 || result.Teams[0]["teamName"] != "Example Studio" {
		t.Fatalf("result = %#v", result)
	}
}

func TestUserTeamListRejectsMissingTokenAndHTTPFailure(t *testing.T) {
	if _, err := UserTeamList(context.Background(), nil, "", ""); err == nil {
		t.Fatal("expected missing token error")
	}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "denied", http.StatusForbidden)
	}))
	defer srv.Close()
	_, err := UserTeamList(context.Background(), srv.Client(), srv.URL, "oauth-token")
	if err == nil || !strings.Contains(err.Error(), "HTTP 403") {
		t.Fatalf("error = %v", err)
	}
}

func TestUserTeamListRejectsInvalidJSON(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{`))
	}))
	defer srv.Close()
	_, err := UserTeamList(context.Background(), srv.Client(), srv.URL, "oauth-token")
	if err == nil || !strings.Contains(err.Error(), "decode team account response") {
		t.Fatalf("error = %v", err)
	}
}
