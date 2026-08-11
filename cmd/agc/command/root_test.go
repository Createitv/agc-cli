package command

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/your-org/agc-cli/pkg/agcapi"
	"github.com/your-org/agc-cli/pkg/domain"
	"github.com/your-org/agc-cli/pkg/project"
)

func execute(args ...string) (string, error) {
	cmd := NewRootCommand()
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)
	cmd.SetArgs(args)
	err := cmd.Execute()
	return buf.String(), err
}

func executeWithContext(ctx context.Context, args ...string) (string, error) {
	cmd := NewRootCommand()
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)
	cmd.SetArgs(args)
	cmd.SetContext(ctx)
	err := cmd.Execute()
	return buf.String(), err
}

func TestCapabilitiesCommandReturnsAllFamilies(t *testing.T) {
	out, err := execute("capabilities")
	if err != nil {
		t.Fatal(err)
	}
	var body domain.Envelope[[]domain.Capability]
	if err := json.Unmarshal([]byte(out), &body); err != nil {
		t.Fatal(err)
	}
	if len(body.Data) != 13 {
		t.Fatalf("capabilities = %d, want 13", len(body.Data))
	}
}

func TestEndpointsCommandReturnsRegisteredInterfaces(t *testing.T) {
	out, err := execute("endpoints")
	if err != nil {
		t.Fatal(err)
	}
	var body domain.Envelope[[]domain.Endpoint]
	if err := json.Unmarshal([]byte(out), &body); err != nil {
		t.Fatal(err)
	}
	if len(body.Data) < 150 {
		t.Fatalf("endpoints = %d, want at least 150", len(body.Data))
	}
}

func TestOpenAPICommand(t *testing.T) {
	out, err := execute("openapi", "--pretty")
	if err != nil {
		t.Fatal(err)
	}
	var body map[string]any
	if err := json.Unmarshal([]byte(out), &body); err != nil {
		t.Fatal(err)
	}
	if body["openapi"] != "3.1.0" {
		t.Fatalf("openapi = %v", body["openapi"])
	}
	if !strings.Contains(out, "/api/v1/publishing/endpoints/app-submit/invoke") {
		t.Fatalf("output missing app-submit invoke route: %s", out)
	}
}

func TestVersionCommand(t *testing.T) {
	out, err := execute("version")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, `"version":"dev"`) {
		t.Fatalf("output = %s", out)
	}
}

func TestEveryRegisteredCommandHasHelp(t *testing.T) {
	root := NewRootCommand()
	for _, cmd := range root.Commands() {
		if cmd.Hidden {
			continue
		}
		out, err := execute(cmd.Name(), "--help")
		if err != nil {
			t.Fatalf("%s --help failed: %v", cmd.Name(), err)
		}
		if !strings.Contains(out, "Usage:") {
			t.Fatalf("%s help missing Usage: %q", cmd.Name(), out)
		}
	}
}

func TestModuleListCommand(t *testing.T) {
	out, err := execute("publishing", "list", "--pretty")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, `"id": "publishing"`) {
		t.Fatalf("output = %s", out)
	}
}

func TestModuleEndpointsCommand(t *testing.T) {
	out, err := execute("publishing", "endpoints")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, `"app-submit"`) {
		t.Fatalf("output = %s", out)
	}
}

func TestEndpointCommandShowsSpec(t *testing.T) {
	out, err := execute("publishing", "app-submit")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, `"method":"POST"`) {
		t.Fatalf("output = %s", out)
	}
}

func TestEveryRegisteredEndpointCanDryRun(t *testing.T) {
	for _, endpoint := range domain.AllEndpoints() {
		args := []string{endpoint.FamilyID, endpoint.ID, "--invoke", "--base-url", "https://example.com"}
		for _, parameter := range endpoint.Parameters {
			if !parameter.Required {
				continue
			}
			switch parameter.In {
			case "path":
				value := "sample"
				if parameter.Name == "callbackUrl" || parameter.Name == "uploadUrl" {
					value = "https://callback.example.com/hook"
				}
				args = append(args, "--param", parameter.Name+"="+value)
			case "query":
				args = append(args, "--query", parameter.Name+"=sample")
			case "header":
				args = append(args, "--header", parameter.Name+"=sample")
			case "body", "file":
				args = append(args, "--field", parameter.Name+"=sample")
			}
		}
		out, err := execute(args...)
		if err != nil {
			t.Fatalf("%s dry-run failed: %v", endpoint.Command, err)
		}
		var body domain.Envelope[agcapi.InvokeResponse]
		if err := json.Unmarshal([]byte(out), &body); err != nil {
			t.Fatalf("%s returned invalid JSON: %v", endpoint.Command, err)
		}
		if body.Data.URL == "" || strings.Contains(body.Data.URL, "{") {
			t.Fatalf("%s built invalid URL %q", endpoint.Command, body.Data.URL)
		}
	}
}

func TestEndpointCommandDryRunBuildsURL(t *testing.T) {
	out, err := execute("publishing", "app-info-query", "--invoke", "--query", "appId=123", "--query", "lang=zh-CN", "--base-url", "https://example.com")
	if err != nil {
		t.Fatal(err)
	}
	var body domain.Envelope[agcapi.InvokeResponse]
	if err := json.Unmarshal([]byte(out), &body); err != nil {
		t.Fatal(err)
	}
	if body.Data.URL != "https://example.com/api/publish/v2/app-info?appId=123&lang=zh-CN" {
		t.Fatalf("url = %s", body.Data.URL)
	}
}

func TestEndpointCommandLoadsActiveCredentialTokenForInvoke(t *testing.T) {
	credentialsPath := filepath.Join(t.TempDir(), "credentials.json")
	serviceAccountPath := writeServiceAccountFile(t)
	if _, err := execute("auth", "login", "--service-account-file", serviceAccountPath, "--name", "prod", "--credentials-path", credentialsPath); err != nil {
		t.Fatal(err)
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") == "" {
			t.Fatal("missing authorization header")
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	defer server.Close()
	out, err := execute(
		"publishing", "app-info-query",
		"--invoke",
		"--dry-run=false",
		"--query", "appId=123",
		"--base-url", server.URL,
		"--credentials-path", credentialsPath,
	)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, `"statusCode":200`) {
		t.Fatalf("output = %s", out)
	}
}

func TestEndpointCommandWritesRawResponseToOutFile(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/csv")
		_, _ = w.Write([]byte("date,value\n2026-08-11,42\n"))
	}))
	defer server.Close()
	outFile := filepath.Join(t.TempDir(), "report.csv")
	out, err := execute(
		"reports", "appdownloadexport",
		"--invoke",
		"--dry-run=false",
		"--param", "appId=123",
		"--query", "from=2026-08-01",
		"--query", "to=2026-08-11",
		"--base-url", server.URL,
		"--token", "token",
		"--out", outFile,
	)
	if err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(outFile)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "date,value\n2026-08-11,42\n" {
		t.Fatalf("file = %q", string(data))
	}
	if strings.Contains(out, "2026-08-11,42") {
		t.Fatalf("raw body should not be embedded when --out is used: %s", out)
	}
}

func TestEndpointCommandRequiresCallbackURL(t *testing.T) {
	if _, err := execute("game-items", "propapi-order", "--invoke"); err == nil {
		t.Fatal("expected missing callbackUrl param error")
	}
}

func TestEndpointCommandAcceptsBodyFieldAndHeader(t *testing.T) {
	out, err := execute("publishing", "add-packageurl", "--invoke", "--field", "appId=123", "--field", "packageUrl=https://example.com/app.app", "--header", "client_id=client", "--base-url", "https://example.com")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, `https://example.com/api/publish/v2/app-package-file/by-url`) {
		t.Fatalf("output = %s", out)
	}
}

func TestEndpointCommandRequiresPathParam(t *testing.T) {
	if _, err := execute("comments", "com-getreviewinfo", "--invoke"); err == nil {
		t.Fatal("expected missing path param error")
	}
}

func TestInitCommandWritesProjectConfig(t *testing.T) {
	dir := t.TempDir()
	out, err := execute("--project", dir, "init", "--app-id", "123", "--project-id", "p1", "--package-name", "com.example.app")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, `"appId":"123"`) {
		t.Fatalf("output = %s", out)
	}
	if _, err := os.Stat(filepath.Join(dir, ".agc", "project.json")); err != nil {
		t.Fatal(err)
	}
}

func TestInitRequiresAppID(t *testing.T) {
	if _, err := execute("init"); err == nil {
		t.Fatal("expected missing app id error")
	}
}

func TestAuthLoginAndCheck(t *testing.T) {
	path := filepath.Join(t.TempDir(), "credentials.json")
	if _, err := execute("auth", "login", "--client-id", "id", "--client-key", "key", "--name", "prod", "--credentials-path", path); err != nil {
		t.Fatal(err)
	}
	out, err := execute("auth", "check", "--credentials-path", path)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, `"name":"prod"`) {
		t.Fatalf("output = %s", out)
	}
}

func TestResolveCredentialUsesProjectProfileAndCLIOverride(t *testing.T) {
	projectDir := t.TempDir()
	if err := project.Save(projectDir, project.Config{AppID: "123", Profile: "production"}); err != nil {
		t.Fatal(err)
	}
	store := agcapi.CredentialStore{Accounts: []agcapi.Credential{
		{Name: "production", Mode: "service-account"},
		{Name: "staging", Mode: "api-client", Active: true},
	}}

	credential, ok, err := resolveCredential(&options{project: projectDir}, store)
	if err != nil || !ok || credential.Name != "production" {
		t.Fatalf("project credential = %#v, ok = %v, err = %v", credential, ok, err)
	}

	credential, ok, err = resolveCredential(&options{project: projectDir, profile: "staging"}, store)
	if err != nil || !ok || credential.Name != "staging" {
		t.Fatalf("CLI credential = %#v, ok = %v, err = %v", credential, ok, err)
	}
}

func TestResolveCredentialFallbackAndErrors(t *testing.T) {
	store := agcapi.CredentialStore{Accounts: []agcapi.Credential{
		{Name: "active", Mode: "api-client", Active: true},
	}}

	credential, ok, err := resolveCredential(&options{project: t.TempDir()}, store)
	if err != nil || !ok || credential.Name != "active" {
		t.Fatalf("active credential = %#v, ok = %v, err = %v", credential, ok, err)
	}

	if _, _, err := resolveCredential(&options{profile: "missing"}, store); err == nil || !strings.Contains(err.Error(), `credential profile "missing" not found`) {
		t.Fatalf("missing profile error = %v", err)
	}

	projectDir := t.TempDir()
	configDir := filepath.Join(projectDir, ".agc")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(configDir, "project.json"), []byte("{"), 0644); err != nil {
		t.Fatal(err)
	}
	if _, _, err := resolveCredential(&options{project: projectDir}, store); err == nil || !strings.Contains(err.Error(), "load project profile") {
		t.Fatalf("corrupt project error = %v", err)
	}
}

func TestAuthCheckUsesProjectBoundProfile(t *testing.T) {
	credentialsPath := filepath.Join(t.TempDir(), "credentials.json")
	if _, err := execute("auth", "login", "--client-id", "prod-id", "--client-key", "prod-key", "--name", "production", "--credentials-path", credentialsPath); err != nil {
		t.Fatal(err)
	}
	if _, err := execute("auth", "login", "--client-id", "stage-id", "--client-key", "stage-key", "--name", "staging", "--credentials-path", credentialsPath); err != nil {
		t.Fatal(err)
	}
	projectDir := t.TempDir()
	if err := project.Save(projectDir, project.Config{AppID: "123", Profile: "production"}); err != nil {
		t.Fatal(err)
	}
	out, err := execute("--project", projectDir, "auth", "check", "--credentials-path", credentialsPath)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, `"name":"production"`) {
		t.Fatalf("output = %s", out)
	}
}

func TestAuthTokenFromServiceAccount(t *testing.T) {
	credentialsPath := filepath.Join(t.TempDir(), "credentials.json")
	serviceAccountPath := writeServiceAccountFile(t)
	if _, err := execute("auth", "login", "--service-account-file", serviceAccountPath, "--name", "prod", "--credentials-path", credentialsPath); err != nil {
		t.Fatal(err)
	}
	out, err := execute("auth", "token", "--credentials-path", credentialsPath)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, `"token_type":"Bearer"`) || !strings.Contains(out, `"access_token"`) {
		t.Fatalf("output = %s", out)
	}
}

func TestAuthList(t *testing.T) {
	path := filepath.Join(t.TempDir(), "credentials.json")
	if _, err := execute("auth", "login", "--client-id", "id", "--client-key", "key", "--name", "prod", "--credentials-path", path); err != nil {
		t.Fatal(err)
	}
	out, err := execute("auth", "list", "--credentials-path", path)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, `"name":"prod"`) {
		t.Fatalf("output = %s", out)
	}
}

func TestAuthLoginRejectsIncompleteCredential(t *testing.T) {
	if _, err := execute("auth", "login", "--client-id", "id", "--credentials-path", filepath.Join(t.TempDir(), "credentials.json")); err == nil {
		t.Fatal("expected validation error")
	}
}

func TestDocsCommand(t *testing.T) {
	out, err := execute("docs", "publishing")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "docs/features/publishing.md") {
		t.Fatalf("output = %s", out)
	}
}

func TestModuleStatusCommand(t *testing.T) {
	out, err := execute("reports", "status")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, `"id":"reports"`) {
		t.Fatalf("output = %s", out)
	}
}

func TestCredentialPathEnvironmentOverride(t *testing.T) {
	t.Setenv("AGC_CREDENTIALS_PATH", "/tmp/agc-test-credentials.json")
	got, err := credentialPath("")
	if err != nil {
		t.Fatal(err)
	}
	if got != "/tmp/agc-test-credentials.json" {
		t.Fatalf("path = %q", got)
	}
}

func TestWebServerCommandCanShutdownFromContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() {
		_, err := executeWithContext(ctx, "web-server", "--addr", ":0")
		done <- err
	}()
	time.Sleep(50 * time.Millisecond)
	cancel()
	select {
	case err := <-done:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("web-server did not shut down")
	}
}

func writeServiceAccountFile(t *testing.T) string {
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
	data, err := json.Marshal(map[string]string{
		"key_id":      "kid-1",
		"private_key": string(pemBytes),
		"sub_account": "sub-1",
		"token_uri":   "https://oauth-login.cloud.huawei.com/oauth2/v3/token",
	})
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), "service-account.json")
	if err := os.WriteFile(path, data, 0600); err != nil {
		t.Fatal(err)
	}
	return path
}
