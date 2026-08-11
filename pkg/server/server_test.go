package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCapabilitiesRoute(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/capabilities", nil)
	rec := httptest.NewRecorder()
	Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d", rec.Code)
	}
	var body struct {
		Data []map[string]any `json:"data"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if len(body.Data) != 13 {
		t.Fatalf("capabilities = %d, want 13", len(body.Data))
	}
}

func TestEndpointsRoute(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/endpoints", nil)
	rec := httptest.NewRecorder()
	Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d", rec.Code)
	}
	var body struct {
		Data []map[string]any `json:"data"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if len(body.Data) < 150 {
		t.Fatalf("endpoints = %d, want at least 150", len(body.Data))
	}
}

func TestFamilyEndpointsRoute(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/publishing/endpoints", nil)
	rec := httptest.NewRecorder()
	Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d", rec.Code)
	}
	var body struct {
		Data []map[string]any `json:"data"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if len(body.Data) == 0 {
		t.Fatal("publishing endpoints empty")
	}
}

func TestSingleEndpointRoute(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/publishing/endpoints/app-submit", nil)
	rec := httptest.NewRecorder()
	Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "提交发布") {
		t.Fatalf("body = %s", rec.Body.String())
	}
}

func TestEndpointInvokeRouteDryRun(t *testing.T) {
	payload := `{"baseUrl":"https://example.com","query":{"appId":"123","lang":"zh-CN"}}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/publishing/endpoints/app-info-query/invoke", bytes.NewBufferString(payload))
	rec := httptest.NewRecorder()
	Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", rec.Code, rec.Body.String())
	}
	var body struct {
		Data struct {
			Method string `json:"method"`
			URL    string `json:"url"`
			DryRun bool   `json:"dryRun"`
		} `json:"data"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body.Data.Method != http.MethodGet || !body.Data.DryRun {
		t.Fatalf("data = %#v", body.Data)
	}
	if body.Data.URL != "https://example.com/api/publish/v2/app-info?appId=123&lang=zh-CN" {
		t.Fatalf("url = %q", body.Data.URL)
	}
}

func TestEndpointInvokeRouteRequiresPost(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/publishing/endpoints/app-info-query/invoke", nil)
	rec := httptest.NewRecorder()
	Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status = %d, want 405", rec.Code)
	}
}

func TestEndpointInvokeRouteValidatesRequiredParams(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/v1/game-items/endpoints/propapi-order/invoke", bytes.NewBufferString(`{}`))
	rec := httptest.NewRecorder()
	Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "missing params.callbackUrl") {
		t.Fatalf("body = %s", rec.Body.String())
	}
}

func TestOpenAPIRoute(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/openapi.json", nil)
	rec := httptest.NewRecorder()
	Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d", rec.Code)
	}
	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body["openapi"] != "3.1.0" {
		t.Fatalf("openapi = %v", body["openapi"])
	}
	if !strings.Contains(rec.Body.String(), "/api/v1/publishing/endpoints/app-submit/invoke") {
		t.Fatalf("spec missing app-submit route")
	}
}

func TestRootRoute(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1", nil)
	rec := httptest.NewRecorder()
	Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d", rec.Code)
	}
	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body["name"] != "agc-cli" {
		t.Fatalf("name = %v", body["name"])
	}
}

func TestResourceRoute(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/publishing", nil)
	rec := httptest.NewRecorder()
	Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d", rec.Code)
	}
	if got := rec.Body.String(); got == "" || !json.Valid(rec.Body.Bytes()) {
		t.Fatalf("invalid json: %q", got)
	}
}

func TestUnknownRoute(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/nope", nil)
	rec := httptest.NewRecorder()
	Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", rec.Code)
	}
}

func TestRunRouteRequiresPost(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/run", nil)
	rec := httptest.NewRecorder()
	Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status = %d, want 405", rec.Code)
	}
}

func TestRunRouteAcceptsPost(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/run", bytes.NewBufferString(`{"command":"agc capabilities"}`))
	rec := httptest.NewRecorder()
	Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusAccepted {
		t.Fatalf("status = %d, want 202", rec.Code)
	}
}
