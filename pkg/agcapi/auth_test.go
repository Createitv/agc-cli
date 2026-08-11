package agcapi

import (
	"os"
	"path/filepath"
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
	if got := info.Mode().Perm(); got != 0600 {
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
