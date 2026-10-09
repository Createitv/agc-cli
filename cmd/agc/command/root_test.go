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

	"github.com/Createitv/agc-cli/pkg/agcapi"
	"github.com/Createitv/agc-cli/pkg/domain"
	"github.com/Createitv/agc-cli/pkg/project"
)

// Tests must never inherit a developer's real local credentials.
func TestMain(m *testing.M) {
	for _, name := range []string{"AGC_PROFILE", "AGC_CLIENT_ID", "AGC_CLIENT_KEY", "AGC_CLIENT_SECRET", "AGC_SERVICE_ACCOUNT_FILE", "AGC_CREDENTIALS_PATH", "AGC_ACCESS_TOKEN", "AGC_SERVER_TOKEN"} {
		if err := os.Unsetenv(name); err != nil {
			os.Exit(1)
		}
	}
	dir, err := os.MkdirTemp("", "agc-command-tests-*")
	if err != nil {
		os.Exit(1)
	}
	if err := os.Setenv("AGC_CREDENTIALS_PATH", filepath.Join(dir, "credentials.json")); err != nil {
		os.RemoveAll(dir)
		os.Exit(1)
	}
	code := m.Run()
	os.RemoveAll(dir)
	os.Exit(code)
}

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

func TestSkillsAddInstallsForAllAgents(t *testing.T) {
	projectDir := t.TempDir()
	out, err := execute("--project", projectDir, "skills", "add", "--agent", "all")
	if err != nil {
		t.Fatal(err)
	}
	var body domain.Envelope[[]struct {
		Agent string `json:"agent"`
		Path  string `json:"path"`
	}]
	if err := json.Unmarshal([]byte(out), &body); err != nil {
		t.Fatal(err)
	}
	want := map[string]string{
		"copilot": ".github/skills/agc-cli/SKILL.md",
		"claude":  ".claude/skills/agc-cli/SKILL.md",
		"codex":   ".agents/skills/agc-cli/SKILL.md",
	}
	if len(body.Data) != len(want) {
		t.Fatalf("installed %d skills, want %d", len(body.Data), len(want))
	}
	for _, skill := range body.Data {
		if want[skill.Agent] != skill.Path {
			t.Errorf("installed agent/path = %q/%q", skill.Agent, skill.Path)
		}
		content, err := os.ReadFile(filepath.Join(projectDir, filepath.FromSlash(skill.Path)))
		if err != nil {
			t.Fatal(err)
		}
		if string(content) != string(agcCLISkill) {
			t.Errorf("%s skill differs from embedded skill", skill.Agent)
		}
	}
}

func TestSkillsAddDoesNotReplaceUnlessForced(t *testing.T) {
	projectDir := t.TempDir()
	destination := filepath.Join(projectDir, ".claude", "skills", "agc-cli", "SKILL.md")
	if err := os.MkdirAll(filepath.Dir(destination), 0755); err != nil {
		t.Fatal(err)
	}
	const existing = "existing skill"
	if err := os.WriteFile(destination, []byte(existing), 0644); err != nil {
		t.Fatal(err)
	}
	if _, err := execute("--project", projectDir, "skills", "add", "--agent", "claude"); err == nil {
		t.Fatal("expected existing skill error")
	}
	content, err := os.ReadFile(destination)
	if err != nil {
		t.Fatal(err)
	}
	if string(content) != existing {
		t.Fatalf("existing skill changed to %q", content)
	}
	if _, err := execute("--project", projectDir, "skills", "add", "--agent", "claude", "--force"); err != nil {
		t.Fatal(err)
	}
	content, err = os.ReadFile(destination)
	if err != nil {
		t.Fatal(err)
	}
	if string(content) != string(agcCLISkill) {
		t.Fatal("forced skill does not match embedded content")
	}
}

func TestSkillsAddRejectsUnsupportedAgent(t *testing.T) {
	if _, err := execute("skills", "add", "--agent", "unknown"); err == nil {
		t.Fatal("expected unsupported agent error")
	}
}

func TestSkillsAddDoesNotReplaceDirectory(t *testing.T) {
	projectDir := t.TempDir()
	destination := filepath.Join(projectDir, ".agents", "skills", "agc-cli", "SKILL.md")
	if err := os.MkdirAll(destination, 0755); err != nil {
		t.Fatal(err)
	}
	if _, err := execute("--project", projectDir, "skills", "add", "--agent", "codex", "--force"); err == nil {
		t.Fatal("expected directory destination error")
	}
	info, err := os.Stat(destination)
	if err != nil {
		t.Fatal(err)
	}
	if !info.IsDir() {
		t.Fatal("destination directory was replaced")
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
				value := "sample"
				if parameter.Type == "integer" || parameter.Type == "number" {
					value = "1"
				}
				if parameter.Type == "array" {
					value = `["fixture"]`
				}
				if parameter.Type == "object" {
					value = `{"key":"fixture"}`
				}
				args = append(args, "--field", parameter.Name+"="+value)
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
		_, err := executeWithContext(ctx, "web-server", "--addr", "127.0.0.1:0")
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

func TestAuthCommandsDoNotExposeClientKey(t *testing.T) {
	path := filepath.Join(t.TempDir(), "credentials.json")
	for _, args := range [][]string{
		{"auth", "login", "--client-id", "test-client", "--client-key", "test-secret", "--credentials-path", path},
		{"auth", "list", "--credentials-path", path},
		{"auth", "check", "--credentials-path", path},
	} {
		out, err := execute(args...)
		if err != nil {
			t.Fatal(err)
		}
		if strings.Contains(out, "test-secret") || strings.Contains(out, "clientKey") {
			t.Fatal("credential output exposes client key")
		}
	}
	store, err := agcapi.LoadCredentials(path)
	if err != nil || store.Accounts[0].ClientKey != "test-secret" {
		t.Fatal("persisted key was lost")
	}
}

func TestEndpointDefaultsContextAndClientID(t *testing.T) {
	t.Setenv("AGC_ACCESS_TOKEN", "")
	dir := t.TempDir()
	path := filepath.Join(dir, "credentials.json")
	if err := project.Save(dir, project.Config{AppID: "app-default"}); err != nil {
		t.Fatal(err)
	}
	if err := agcapi.SaveCredential(path, agcapi.Credential{Name: "default", Mode: "api-client", ClientID: "client-default", ClientKey: "test-key"}); err != nil {
		t.Fatal(err)
	}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/token") {
			w.Write([]byte(`{"access_token":"test-token"}`))
			return
		}
		if r.URL.Query().Get("appId") != "app-default" || r.Header.Get("client_id") != "client-default" {
			t.Error("missing automatic app or client context")
		}
		w.Write([]byte(`{}`))
	}))
	defer srv.Close()
	_, err := execute("--project", dir, "publishing", "app-info-query", "--invoke", "--dry-run=false", "--query", "lang=en-US", "--credentials-path", path, "--base-url", srv.URL)
	if err != nil {
		t.Fatal(err)
	}
}

func TestAuthRemoteCheckValidatesReadRequest(t *testing.T) {
	dir := t.TempDir()
	if err := project.Save(dir, project.Config{AppID: "test-app"}); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(dir, "credentials.json")
	if err := agcapi.SaveCredential(path, agcapi.Credential{Name: "default", Mode: "api-client", ClientID: "test-client", ClientKey: "test-key"}); err != nil {
		t.Fatal(err)
	}
	calls := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		if strings.HasSuffix(r.URL.Path, "/token") {
			w.Write([]byte(`{"access_token":"test-token"}`))
			return
		}
		if r.URL.Path != "/api/publish/v2/app-info" || r.URL.Query().Get("appId") != "test-app" {
			t.Error("API client remote verification did not use app read")
		}
		if r.Header.Get("Authorization") != "Bearer test-token" || r.Header.Get("client_id") != "test-client" {
			t.Error("missing authentication")
		}
		w.Write([]byte(`{"projects":[]}`))
	}))
	defer srv.Close()
	out, err := execute("auth", "check", "--remote", "--base-url", srv.URL, "--credentials-path", path, "--project", dir)
	if err != nil {
		t.Fatal(err)
	}
	if calls != 2 || !strings.Contains(out, `"remoteVerified":true`) || strings.Contains(out, "test-token") || strings.Contains(out, "test-key") {
		t.Fatal("remote verification missing or exposes secret")
	}
}

func TestBodyFileDoesNotBypassRequiredFields(t *testing.T) {
	path := filepath.Join(t.TempDir(), "body.json")
	if err := os.WriteFile(path, []byte(`{}`), 0600); err != nil {
		t.Fatal(err)
	}
	cmd := endpointCommand(&options{project: t.TempDir(), output: "json", timeout: time.Second}, domain.Endpoint{ID: "test-body", Method: "POST", Path: "/test", Parameters: []domain.Parameter{{Name: "required", In: "body", Required: true}}})
	cmd.SetArgs([]string{"--invoke", "--body", path})
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)
	err := cmd.Execute()
	if err == nil {
		t.Fatal("empty body bypassed required fields")
	}
}

func TestEndpointExplicitContextAndHeaderWin(t *testing.T) {
	t.Setenv("AGC_ACCESS_TOKEN", "")
	dir := t.TempDir()
	path := filepath.Join(dir, "credentials.json")
	if err := project.Save(dir, project.Config{AppID: "app-default"}); err != nil {
		t.Fatal(err)
	}
	if err := agcapi.SaveCredential(path, agcapi.Credential{Name: "default", Mode: "api-client", ClientID: "client-default", ClientKey: "test-key"}); err != nil {
		t.Fatal(err)
	}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/token") {
			w.Write([]byte(`{"access_token":"test-token"}`))
			return
		}
		if r.URL.Query().Get("appId") != "app-explicit" || r.Header.Get("client_id") != "client-explicit" {
			t.Error("explicit context was overwritten")
		}
		w.Write([]byte(`{}`))
	}))
	defer srv.Close()
	_, err := execute("--project", dir, "publishing", "app-info-query", "--invoke", "--dry-run=false", "--query", "lang=en-US", "--query", "appId=app-explicit", "--header", "CLIENT_ID=client-explicit", "--credentials-path", path, "--base-url", srv.URL)
	if err != nil {
		t.Fatal(err)
	}
}

func TestProjectContextSupportsAppIDAndHeader(t *testing.T) {
	dir := t.TempDir()
	if err := project.Save(dir, project.Config{AppID: "configured-app"}); err != nil {
		t.Fatal(err)
	}
	for _, parameter := range []domain.Parameter{{Name: "appID", In: "query", Required: true}, {Name: "appId", In: "header", Required: true}} {
		cmd := endpointCommand(&options{project: dir, output: "json", timeout: time.Second}, domain.Endpoint{ID: "context", Method: "GET", Path: "/context", Parameters: []domain.Parameter{parameter}})
		var buf bytes.Buffer
		cmd.SetOut(&buf)
		cmd.SetErr(&buf)
		cmd.SetArgs([]string{"--invoke"})
		if err := cmd.Execute(); err != nil {
			t.Fatal(err)
		}
	}
}

func TestAPIClientRemoteCheckWithoutAppOnlyVerifiesToken(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "credentials.json")
	if err := agcapi.SaveCredential(path, agcapi.Credential{Name: "default", Mode: "api-client", ClientID: "test-client", ClientKey: "test-key"}); err != nil {
		t.Fatal(err)
	}
	calls := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		if !strings.HasSuffix(r.URL.Path, "/token") {
			t.Error("API client attempted unsupported project read")
		}
		w.Write([]byte(`{"access_token":"test-token"}`))
	}))
	defer srv.Close()
	out, err := execute("auth", "check", "--remote", "--base-url", srv.URL, "--credentials-path", path, "--project", dir)
	if err != nil {
		t.Fatal(err)
	}
	if calls != 1 || !strings.Contains(out, `"tokenVerified":true`) || !strings.Contains(out, `"readVerified":false`) {
		t.Fatal("token-only verification state incorrect")
	}
}

func TestProjectHeaderExplicitOverrideIsCaseInsensitive(t *testing.T) {
	dir := t.TempDir()
	if err := project.Save(dir, project.Config{AppID: "configured-app"}); err != nil {
		t.Fatal(err)
	}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("appId") != "explicit-app" {
			t.Error("explicit header overwritten")
		}
		w.Write([]byte(`{}`))
	}))
	defer srv.Close()
	cmd := endpointCommand(&options{project: dir, output: "json", timeout: time.Second}, domain.Endpoint{ID: "context", Method: "GET", Path: "/context", Parameters: []domain.Parameter{{Name: "appId", In: "header", Required: true}}})
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)
	cmd.SetArgs([]string{"--invoke", "--dry-run=false", "--token", "test-token", "--header", "APPID=explicit-app", "--base-url", srv.URL})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
}

func TestWebServerRejectsExternalBindingWithoutToken(t *testing.T) {
	t.Setenv("AGC_SERVER_TOKEN", "")
	for _, addr := range []string{":8421", "0.0.0.0:8421", "[::]:8421", "192.0.2.1:8421"} {
		_, err := execute("web-server", "--addr", addr)
		if err == nil || !strings.Contains(err.Error(), "AGC_SERVER_TOKEN") {
			t.Fatalf("external bind was not rejected: %s", addr)
		}
	}
}

func TestLoopbackListenAddresses(t *testing.T) {
	for _, addr := range []string{"127.0.0.1:8421", "127.0.1.2:8421", "localhost:8421", "[::1]:8421"} {
		if !loopbackListenAddress(addr) {
			t.Fatalf("loopback address rejected: %s", addr)
		}
	}
	for _, addr := range []string{":8421", "0.0.0.0:8421", "[::]:8421", "localhost.example:8421", "invalid"} {
		if loopbackListenAddress(addr) {
			t.Fatalf("external or invalid address accepted: %s", addr)
		}
	}
}

func TestExplicitAuthenticationHeadersDoNotRequireProfile(t *testing.T) {
	for _, envToken := range []string{"", "environment-token"} {
		for _, header := range []string{"AUTHORIZATION=Custom test-auth", "OAUTH2TOKEN=test-oauth"} {
			t.Run(strings.Split(header, "=")[0], func(t *testing.T) {
				t.Setenv("AGC_ACCESS_TOKEN", envToken)
				srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if strings.HasPrefix(header, "AUTHORIZATION") && r.Header.Get("Authorization") != "Custom test-auth" {
						t.Error("explicit authorization overwritten")
					}
					if strings.HasPrefix(header, "OAUTH2TOKEN") && (r.Header.Get("oauth2Token") != "test-oauth" || r.Header.Get("Authorization") != "") {
						t.Error("environment bearer injected into OAuth request")
					}
					w.Write([]byte(`{}`))
				}))
				defer srv.Close()
				dir := t.TempDir()
				cmd := endpointCommand(&options{project: dir, output: "json", timeout: time.Second}, domain.Endpoint{ID: "explicit-auth", Method: "GET", Path: "/test"})
				var buf bytes.Buffer
				cmd.SetOut(&buf)
				cmd.SetErr(&buf)
				cmd.SetArgs([]string{"--invoke", "--dry-run=false", "--header", header, "--credentials-path", filepath.Join(dir, "missing.json"), "--base-url", srv.URL})
				if err := cmd.Execute(); err != nil {
					t.Fatal(err)
				}
			})
		}
	}
}

func TestAuthLoginImportsAPIClientFileWithoutActivating(t *testing.T) {
	for _, content := range []string{`{"client_id":"import-id","client_secret":"import-secret"}`, `{"clientId":"import-id","clientKey":"import-secret"}`} {
		dir := t.TempDir()
		path := filepath.Join(dir, "credentials.json")
		file := filepath.Join(dir, "client.json")
		if err := agcapi.SaveCredential(path, agcapi.Credential{Name: "service", Mode: "service-account", ServiceAccountFile: "service.json"}); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(file, []byte(content), 0600); err != nil {
			t.Fatal(err)
		}
		out, err := execute("auth", "login", "--api-client-file", file, "--credentials-path", path, "--name", "imported", "--activate=false")
		if err != nil {
			t.Fatal(err)
		}
		if strings.Contains(out, "import-secret") {
			t.Fatal("secret exposed")
		}
		store, err := agcapi.LoadCredentials(path)
		if err != nil {
			t.Fatal(err)
		}
		active, ok := agcapi.ActiveCredential(store)
		if !ok || active.Name != "service" {
			t.Fatal("active profile changed")
		}
		imported, ok := agcapi.CredentialByName(store, "imported")
		if !ok || imported.ClientKey != "import-secret" || imported.ClientID != "import-id" {
			t.Fatal("import lost credential")
		}
	}
}

func TestAuthLoginRejectsInvalidOrMixedAPIClientFile(t *testing.T) {
	for _, content := range []string{`{"client_id":"id","client_secret":"secret","clientKey":"different-secret"}`, `{"client_id":"id","client_secret":"secret","clientId":"other"}`, `{"client_id":"id"}`, `{"client_secret":"secret"`, `[]`, `null`} {
		dir := t.TempDir()
		file := filepath.Join(dir, "client.json")
		if err := os.WriteFile(file, []byte(content), 0600); err != nil {
			t.Fatal(err)
		}
		out, err := execute("auth", "login", "--api-client-file", file, "--credentials-path", filepath.Join(dir, "credentials.json"))
		if err == nil {
			t.Fatal("invalid file accepted")
		}
		if strings.Contains(out, "different-secret") {
			t.Fatal("secret reflected in error")
		}
	}
	dir := t.TempDir()
	file := filepath.Join(dir, "client.json")
	if err := os.WriteFile(file, []byte(`{"client_id":"id","client_secret":"secret"}`), 0600); err != nil {
		t.Fatal(err)
	}
	for _, flag := range []string{"--client-id", "--client-key", "--service-account-file"} {
		_, err := execute("auth", "login", "--api-client-file", file, flag, "explicit", "--credentials-path", filepath.Join(dir, "credentials.json"))
		if err == nil {
			t.Fatal("mixed options accepted")
		}
	}
}

func clearCredentialEnvironment(t *testing.T) {
	t.Helper()
	for _, name := range []string{"AGC_PROFILE", "AGC_CLIENT_ID", "AGC_CLIENT_KEY", "AGC_CLIENT_SECRET", "AGC_SERVICE_ACCOUNT_FILE"} {
		t.Setenv(name, "")
	}
}

func TestCredentialEnvironmentFallback(t *testing.T) {
	clearCredentialEnvironment(t)
	opts := &options{project: t.TempDir()}
	t.Setenv("AGC_CLIENT_ID", "environment-id")
	t.Setenv("AGC_CLIENT_SECRET", "environment-secret")
	credential, ok, err := resolveCredential(opts, agcapi.CredentialStore{})
	if err != nil || !ok || credential.Mode != "api-client" || credential.ClientID != "environment-id" || credential.ClientKey != "environment-secret" {
		t.Fatal("API client environment fallback failed")
	}
	t.Setenv("AGC_CLIENT_ID", "")
	t.Setenv("AGC_CLIENT_SECRET", "")
	t.Setenv("AGC_SERVICE_ACCOUNT_FILE", "environment-service.json")
	credential, ok, err = resolveCredential(opts, agcapi.CredentialStore{})
	if err != nil || !ok || credential.Mode != "service-account" || credential.ServiceAccountFile != "environment-service.json" {
		t.Fatal("service account fallback failed")
	}
}

func TestCredentialEnvironmentRejectsIncompleteOrConflictingPair(t *testing.T) {
	for _, values := range []map[string]string{{"AGC_CLIENT_ID": "id"}, {"AGC_CLIENT_KEY": "test-secret"}, {"AGC_CLIENT_ID": "id", "AGC_CLIENT_KEY": "test-secret", "AGC_CLIENT_SECRET": "other-secret"}} {
		clearCredentialEnvironment(t)
		for name, value := range values {
			t.Setenv(name, value)
		}
		_, _, err := resolveCredential(&options{project: t.TempDir()}, agcapi.CredentialStore{})
		if err == nil {
			t.Fatal("invalid credential environment accepted")
		}
		if strings.Contains(err.Error(), "test-secret") || strings.Contains(err.Error(), "other-secret") {
			t.Fatal("secret reflected in error")
		}
	}
}

func TestCredentialProfilePriorityWithEnvironment(t *testing.T) {
	clearCredentialEnvironment(t)
	dir := t.TempDir()
	opts := &options{project: dir}
	store := agcapi.CredentialStore{Accounts: []agcapi.Credential{{Name: "active", Mode: "service-account", Active: true}, {Name: "environment"}, {Name: "project"}, {Name: "explicit"}}}
	t.Setenv("AGC_CLIENT_ID", "id")
	t.Setenv("AGC_CLIENT_KEY", "secret")
	credential, _, err := resolveCredential(opts, store)
	if err != nil || credential.Name != "active" {
		t.Fatal("environment replaced active credential")
	}
	t.Setenv("AGC_PROFILE", "environment")
	credential, _, err = resolveCredential(opts, store)
	if err != nil || credential.Name != "environment" {
		t.Fatal("environment profile not selected")
	}
	if err := project.Save(dir, project.Config{AppID: "app", Profile: "project"}); err != nil {
		t.Fatal(err)
	}
	credential, _, err = resolveCredential(opts, store)
	if err != nil || credential.Name != "project" {
		t.Fatal("environment replaced project profile")
	}
	opts.profile = "explicit"
	credential, _, err = resolveCredential(opts, store)
	if err != nil || credential.Name != "explicit" {
		t.Fatal("explicit profile not prioritized")
	}
	opts.profile = "missing"
	if _, _, err := resolveCredential(opts, store); err == nil {
		t.Fatal("missing explicit profile fell back")
	}
}
