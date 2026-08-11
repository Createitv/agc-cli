package agcapi

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/your-org/agc-cli/pkg/domain"
)

func TestBuildEndpointURLReplacesParams(t *testing.T) {
	url, err := BuildEndpointURL("https://example.com", "/api/reviews/v1/manage/dev/reviews/{reviewId}", map[string]string{"reviewId": "abc"}, map[string]string{"lang": "zh-CN"})
	if err != nil {
		t.Fatal(err)
	}
	if url != "https://example.com/api/reviews/v1/manage/dev/reviews/abc?lang=zh-CN" {
		t.Fatalf("url = %q", url)
	}
}

func TestBuildEndpointURLRejectsMissingParam(t *testing.T) {
	if _, err := BuildEndpointURL("https://example.com", "/apps/{appId}", nil, nil); err == nil {
		t.Fatal("expected missing param error")
	}
}

func TestBuildEndpointURLUsesUploadURL(t *testing.T) {
	url, err := BuildEndpointURL("https://example.com", "{uploadUrl}", map[string]string{"uploadUrl": "https://upload.example.com/file"}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if url != "https://upload.example.com/file" {
		t.Fatalf("url = %q", url)
	}
}

func TestBuildEndpointURLUsesCallbackURL(t *testing.T) {
	url, err := BuildEndpointURL("https://example.com", "{callbackUrl}", map[string]string{"callbackUrl": "https://callback.example.com/order"}, map[string]string{"dryRun": "1"})
	if err != nil {
		t.Fatal(err)
	}
	if url != "https://callback.example.com/order?dryRun=1" {
		t.Fatalf("url = %q", url)
	}
}

func TestBuildEndpointURLEscapesAndSortsQuery(t *testing.T) {
	url, err := BuildEndpointURL("https://example.com", "/apps", nil, map[string]string{"b": "hello world", "a": "zh-CN"})
	if err != nil {
		t.Fatal(err)
	}
	if url != "https://example.com/apps?a=zh-CN&b=hello+world" {
		t.Fatalf("url = %q", url)
	}
}

func TestInvokeEndpointDryRun(t *testing.T) {
	endpoint, _ := domain.EndpointByID("publishing", "app-info-query")
	resp, err := InvokeEndpoint(context.Background(), nil, InvokeRequest{
		Endpoint: endpoint,
		BaseURL:  "https://example.com",
		Query:    map[string]string{"appId": "123"},
		DryRun:   true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if !resp.DryRun || resp.URL != "https://example.com/api/publish/v2/app-info?appId=123" {
		t.Fatalf("response = %#v", resp)
	}
}

func TestInvokeEndpointSendsRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer token" {
			t.Fatalf("authorization = %q", r.Header.Get("Authorization"))
		}
		if r.Header.Get("client_id") != "client" {
			t.Fatalf("client_id = %q", r.Header.Get("client_id"))
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	defer server.Close()
	endpoint := domain.Endpoint{Method: http.MethodGet, Path: "/apps/{appId}"}
	resp, err := InvokeEndpoint(context.Background(), server.Client(), InvokeRequest{
		Endpoint:    endpoint,
		BaseURL:     server.URL,
		Params:      map[string]string{"appId": "123"},
		Headers:     map[string]string{"client_id": "client"},
		AccessToken: "token",
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK || string(resp.Body) != `{"ok":true}` {
		t.Fatalf("response = %#v", resp)
	}
	if string(resp.RawBody) != `{"ok":true}` {
		t.Fatalf("raw body = %q", string(resp.RawBody))
	}
}

func TestInvokeEndpointPreservesNonJSONRawBody(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/csv")
		_, _ = w.Write([]byte("date,value\n2026-08-11,42\n"))
	}))
	defer server.Close()
	endpoint := domain.Endpoint{Method: http.MethodGet, Path: "/report.csv"}
	resp, err := InvokeEndpoint(context.Background(), server.Client(), InvokeRequest{Endpoint: endpoint, BaseURL: server.URL})
	if err != nil {
		t.Fatal(err)
	}
	if resp.ContentType != "text/csv" {
		t.Fatalf("content type = %q", resp.ContentType)
	}
	if string(resp.RawBody) != "date,value\n2026-08-11,42\n" {
		t.Fatalf("raw body = %q", string(resp.RawBody))
	}
}

func TestInvokeEndpointReturnsStructuredAGCError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"code":"PARAM_ERROR","message":"missing appId"}`))
	}))
	defer server.Close()
	endpoint := domain.Endpoint{Method: http.MethodGet, Path: "/apps"}
	_, err := InvokeEndpoint(context.Background(), server.Client(), InvokeRequest{Endpoint: endpoint, BaseURL: server.URL})
	if err == nil {
		t.Fatal("expected AGC error")
	}
	agcErr, ok := err.(AGCError)
	if !ok {
		t.Fatalf("error = %T, want AGCError", err)
	}
	if agcErr.Code != "PARAM_ERROR" || agcErr.Message != "missing appId" {
		t.Fatalf("agc error = %#v", agcErr)
	}
}

func TestParseAGCErrorHandlesEmptyBody(t *testing.T) {
	agcErr := ParseAGCError(http.StatusInternalServerError, nil)
	if agcErr.Message != "empty error response" {
		t.Fatalf("message = %q", agcErr.Message)
	}
}

func TestServiceAccountJWT(t *testing.T) {
	account := serviceAccountForTest(t)
	jwt, err := ServiceAccountJWT(account, time.Unix(1000, 0).UTC())
	if err != nil {
		t.Fatal(err)
	}
	parts := strings.Split(jwt, ".")
	if len(parts) != 3 {
		t.Fatalf("jwt parts = %d, want 3", len(parts))
	}
	var header map[string]any
	if err := decodeJWTPart(parts[0], &header); err != nil {
		t.Fatal(err)
	}
	if header["alg"] != "PS256" || header["kid"] != account.KeyID {
		t.Fatalf("header = %#v", header)
	}
	var payload map[string]any
	if err := decodeJWTPart(parts[1], &payload); err != nil {
		t.Fatal(err)
	}
	if payload["iss"] != account.SubAccount {
		t.Fatalf("payload = %#v", payload)
	}
}

func TestAccessTokenFromServiceAccountFile(t *testing.T) {
	account := serviceAccountForTest(t)
	path := filepath.Join(t.TempDir(), "service-account.json")
	data, err := json.Marshal(account)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, data, 0600); err != nil {
		t.Fatal(err)
	}
	token, err := AccessToken(context.Background(), nil, "", Credential{Mode: "service-account", ServiceAccountFile: path})
	if err != nil {
		t.Fatal(err)
	}
	if token.AccessToken == "" || token.ExpiresIn != 3600 {
		t.Fatalf("token = %#v", token)
	}
}

func TestAPIClientToken(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/oauth2/v1/token" {
			t.Fatalf("path = %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"access_token":"abc","token_type":"Bearer","expires_in":3600}`))
	}))
	defer server.Close()
	token, err := APIClientToken(context.Background(), server.Client(), server.URL, "client", "secret")
	if err != nil {
		t.Fatal(err)
	}
	if token.AccessToken != "abc" {
		t.Fatalf("token = %#v", token)
	}
}

func TestAPIClientTokenRejectsHTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "bad credentials", http.StatusUnauthorized)
	}))
	defer server.Close()
	if _, err := APIClientToken(context.Background(), server.Client(), server.URL, "client", "secret"); err == nil {
		t.Fatal("expected HTTP error")
	}
}

func TestLoadServiceAccountRejectsMissingFields(t *testing.T) {
	path := filepath.Join(t.TempDir(), "service-account.json")
	if err := os.WriteFile(path, []byte(`{"key_id":"kid"}`), 0600); err != nil {
		t.Fatal(err)
	}
	if _, err := LoadServiceAccount(path); err == nil {
		t.Fatal("expected missing fields error")
	}
}

func serviceAccountForTest(t *testing.T) ServiceAccountFile {
	t.Helper()
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	der, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		t.Fatal(err)
	}
	pemBytes := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: der})
	return ServiceAccountFile{KeyID: "kid-1", PrivateKey: string(pemBytes), SubAccount: "sub-1", TokenURI: serviceAccountAudience}
}

func decodeJWTPart(part string, target any) error {
	data, err := base64.RawURLEncoding.DecodeString(part)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, target)
}
