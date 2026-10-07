package agcapi

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestSaveCredentialReplacesActiveProfile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "credentials.json")
	if err := SaveCredential(path, Credential{Name: "prod", Mode: "api-client", ClientID: "id", ClientKey: "key"}); err != nil {
		t.Fatal(err)
	}
	if err := SaveCredential(path, Credential{Name: "staging", Mode: "service-account", ServiceAccountFile: "service.json"}); err != nil {
		t.Fatal(err)
	}
	store, err := LoadCredentials(path)
	if err != nil {
		t.Fatal(err)
	}
	active, ok := ActiveCredential(store)
	if !ok {
		t.Fatal("no active credential")
	}
	if active.Name != "staging" {
		t.Fatalf("active = %q, want staging", active.Name)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if got := info.Mode().Perm(); runtime.GOOS != "windows" && got != 0600 {
		t.Fatalf("mode = %v, want 0600", got)
	}
}

func TestValidateCredential(t *testing.T) {
	if err := ValidateCredential(Credential{Mode: "service-account"}); err == nil {
		t.Fatal("expected missing service account file error")
	}
	if err := ValidateCredential(Credential{Mode: "api-client", ClientID: "id"}); err == nil {
		t.Fatal("expected missing client key error")
	}
	if err := ValidateCredential(Credential{Mode: "api-client", ClientID: "id", ClientKey: "key"}); err != nil {
		t.Fatal(err)
	}
}

func TestLoadCredentialsMissingFileReturnsEmptyStore(t *testing.T) {
	store, err := LoadCredentials(filepath.Join(t.TempDir(), "missing.json"))
	if err != nil {
		t.Fatal(err)
	}
	if len(store.Accounts) != 0 {
		t.Fatalf("accounts = %d, want 0", len(store.Accounts))
	}
}

func TestActiveCredentialUsesSingleAccountFallback(t *testing.T) {
	active, ok := ActiveCredential(CredentialStore{Accounts: []Credential{{Name: "only", Mode: "api-client"}}})
	if !ok {
		t.Fatal("expected single account fallback")
	}
	if active.Name != "only" {
		t.Fatalf("active = %q", active.Name)
	}
}

func TestCredentialByName(t *testing.T) {
	store := CredentialStore{Accounts: []Credential{
		{Name: "production", Mode: "service-account"},
		{Name: "staging", Mode: "api-client"},
	}}
	credential, ok := CredentialByName(store, "staging")
	if !ok || credential.Name != "staging" {
		t.Fatalf("credential = %#v, ok = %v", credential, ok)
	}
	if _, ok := CredentialByName(store, "missing"); ok {
		t.Fatal("unexpected match for missing profile")
	}
}

func TestClientDefaults(t *testing.T) {
	client := NewClient("")
	if client.BaseURL != "https://connect-api.cloud.huawei.com" {
		t.Fatalf("base url = %q", client.BaseURL)
	}
	if got := client.TokenEndpoint(); got != "https://connect-api.cloud.huawei.com/api/oauth2/v1/token" {
		t.Fatalf("token endpoint = %q", got)
	}
}

func TestSaveCredentialSecuresExistingFileAndReplacesAtomically(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "credentials.json")
	if err := os.WriteFile(path, []byte(`{"accounts":[]}`), 0644); err != nil {
		t.Fatal(err)
	}
	old, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer old.Close()
	oldInfo, err := old.Stat()
	if err != nil {
		t.Fatal(err)
	}
	// Windows does not permit replacing this file while the test holds it open.
	if runtime.GOOS == "windows" {
		if err := old.Close(); err != nil {
			t.Fatal(err)
		}
	}
	if err := SaveCredential(path, Credential{Name: "default", Mode: "api-client", ClientID: "test-id", ClientKey: "test-secret"}); err != nil {
		t.Fatal(err)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if runtime.GOOS != "windows" && info.Mode().Perm() != 0600 {
		t.Fatal("existing credentials file remains readable by others")
	}
	if os.SameFile(info, oldInfo) {
		t.Fatal("credential file was updated in place instead of atomically replaced")
	}
	store, err := LoadCredentials(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(store.Accounts) != 1 || store.Accounts[0].ClientKey != "test-secret" {
		t.Fatal("credential was not preserved")
	}
	files, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(files) != 1 {
		t.Fatal("temporary credential file was not cleaned up")
	}
}

func TestCredentialPathsRespectHomeAndRejectMissingHome(t *testing.T) {
	dir := t.TempDir()
	homeVariable := "HOME"
	if runtime.GOOS == "windows" {
		homeVariable = "USERPROFILE"
	}
	t.Setenv(homeVariable, dir)
	path, err := CredentialsPath()
	if err != nil || path != filepath.Join(dir, ".agc", "credentials.json") {
		t.Fatal("credential path did not use user's home")
	}
	t.Setenv(homeVariable, "")
	if _, err := ConfigDir(); err == nil {
		t.Fatal("missing home accepted")
	}
	if _, err := CredentialsPath(); err == nil {
		t.Fatal("missing home accepted")
	}
}

func TestCredentialImportValidatesAliasesWithoutEchoingSecrets(t *testing.T) {
	for _, content := range []string{`{"client_id":"test-id","client_secret":"test-secret"}`, `{"clientId":"test-id","clientKey":"test-secret"}`, `{"client_id":"test-id","clientId":"test-id","client_secret":"test-secret","clientKey":"test-secret"}`} {
		file := filepath.Join(t.TempDir(), "client.json")
		if err := os.WriteFile(file, []byte(content), 0600); err != nil {
			t.Fatal(err)
		}
		credential, err := LoadAPIClientCredential(file)
		if err != nil || credential.ClientID != "test-id" || credential.ClientKey != "test-secret" {
			t.Fatal("valid import failed")
		}
		data, err := json.Marshal(credential.View())
		if err != nil {
			t.Fatal(err)
		}
		if strings.Contains(string(data), "test-secret") || strings.Contains(string(data), "clientKey") {
			t.Fatal("profile metadata exposes secret")
		}
	}
	for _, content := range []string{`{"client_id":"test-id","client_secret":"test-secret","clientKey":"another-secret"}`, `{"client_id":"test-id","client_secret":false}`, `{"client_id":"test-id","client_secret":null}`, `{"client_id":"test-id","client_secret":""}`, `{"client_id":"test-id"}`, `[]`, `null`, `{"client_secret":"test-secret"`} {
		file := filepath.Join(t.TempDir(), "client.json")
		if err := os.WriteFile(file, []byte(content), 0600); err != nil {
			t.Fatal(err)
		}
		_, err := LoadAPIClientCredential(file)
		if err == nil {
			t.Fatal("invalid credential import accepted")
		}
		if strings.Contains(err.Error(), "test-secret") || strings.Contains(err.Error(), "another-secret") {
			t.Fatal("credential value reflected in error")
		}
	}
	if _, err := LoadAPIClientCredential(filepath.Join(t.TempDir(), "missing.json")); err == nil {
		t.Fatal("missing file accepted")
	}
}

func TestCredentialEnvironmentValidationAndPrecedence(t *testing.T) {
	for _, name := range []string{"AGC_CLIENT_ID", "AGC_CLIENT_KEY", "AGC_CLIENT_SECRET", "AGC_SERVICE_ACCOUNT_FILE"} {
		t.Setenv(name, "")
	}
	if _, ok, err := CredentialFromEnvironment(); err != nil || ok {
		t.Fatal("empty environment creates credential")
	}
	t.Setenv("AGC_SERVICE_ACCOUNT_FILE", "service.json")
	credential, ok, err := CredentialFromEnvironment()
	if err != nil || !ok || credential.Mode != "service-account" {
		t.Fatal("service account environment failed")
	}
	t.Setenv("AGC_CLIENT_ID", "test-id")
	if _, _, err := CredentialFromEnvironment(); err == nil {
		t.Fatal("incomplete API client pair silently fell back to service account")
	}
	t.Setenv("AGC_CLIENT_SECRET", "test-secret")
	credential, ok, err = CredentialFromEnvironment()
	if err != nil || !ok || credential.Mode != "api-client" || credential.ClientKey != "test-secret" {
		t.Fatal("complete API client did not take priority")
	}
	t.Setenv("AGC_CLIENT_KEY", "test-secret")
	if _, ok, err := CredentialFromEnvironment(); err != nil || !ok {
		t.Fatal("matching aliases rejected")
	}
	t.Setenv("AGC_CLIENT_KEY", "different-secret")
	_, _, err = CredentialFromEnvironment()
	if err == nil {
		t.Fatal("conflicting aliases accepted")
	}
	if strings.Contains(err.Error(), "different-secret") || strings.Contains(err.Error(), "test-secret") {
		t.Fatal("secret reflected in environment error")
	}
}

func TestSaveCredentialPreservesProfileActivationAndOriginalStoreOnFailure(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "credentials.json")
	if err := SaveCredential(path, Credential{Mode: "api-client", ClientID: "test-id", ClientKey: "test-secret"}); err != nil {
		t.Fatal(err)
	}
	if err := SaveCredentialWithActivation(path, Credential{Name: "default", Mode: "api-client", ClientID: "replacement-id", ClientKey: "replacement-secret"}, false); err != nil {
		t.Fatal(err)
	}
	store, err := LoadCredentials(path)
	if err != nil {
		t.Fatal(err)
	}
	active, ok := ActiveCredential(store)
	if !ok || active.ClientID != "replacement-id" {
		t.Fatal("updating active profile without activation lost active state")
	}
	if err := SaveCredentialWithActivation(path, Credential{Name: "secondary", Mode: "service-account", ServiceAccountFile: "service.json"}, false); err != nil {
		t.Fatal(err)
	}
	store, err = LoadCredentials(path)
	if err != nil {
		t.Fatal(err)
	}
	active, ok = ActiveCredential(store)
	if !ok || active.Name != "default" {
		t.Fatal("inactive import changed active profile")
	}
	before, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := SaveCredential(path, Credential{Name: "invalid"}); err == nil {
		t.Fatal("missing mode accepted")
	}
	after, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(before) != string(after) {
		t.Fatal("failed save modified original credential store")
	}
	if err := os.WriteFile(path, []byte(`{"accounts":`), 0600); err != nil {
		t.Fatal(err)
	}
	if err := SaveCredential(path, Credential{Mode: "api-client"}); err == nil {
		t.Fatal("malformed existing store overwritten")
	}
	if _, ok := ActiveCredential(CredentialStore{Accounts: []Credential{{Name: "one"}, {Name: "two"}}}); ok {
		t.Fatal("ambiguous inactive store selected a credential")
	}
	if _, err := LoadCredentials(dir); err == nil {
		t.Fatal("directory accepted as credential store")
	}
	if err := ValidateCredential(Credential{Mode: "unsupported"}); err == nil {
		t.Fatal("unsupported mode accepted")
	}
}
