package agcapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type Credential struct {
	Name               string    `json:"name"`
	Mode               string    `json:"mode"`
	ServiceAccountFile string    `json:"serviceAccountFile,omitempty"`
	ClientID           string    `json:"clientId,omitempty"`
	ClientKey          string    `json:"clientKey,omitempty"`
	Active             bool      `json:"active,omitempty"`
	CreatedAt          time.Time `json:"createdAt"`
}

type CredentialStore struct {
	Accounts []Credential `json:"accounts"`
}

func ConfigDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".agc"), nil
}

func CredentialsPath() (string, error) {
	dir, err := ConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "credentials.json"), nil
}

func LoadCredentials(path string) (CredentialStore, error) {
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return CredentialStore{}, nil
	}
	if err != nil {
		return CredentialStore{}, err
	}
	var store CredentialStore
	if err := json.Unmarshal(data, &store); err != nil {
		return CredentialStore{}, err
	}
	return store, nil
}

func SaveCredential(path string, credential Credential) error {
	return SaveCredentialWithActivation(path, credential, true)
}

// SaveCredentialWithActivation can preserve the current active profile when importing.
func SaveCredentialWithActivation(path string, credential Credential, activate bool) error {
	if credential.Name == "" {
		credential.Name = "default"
	}
	if credential.Mode == "" {
		return errors.New("credential mode is required")
	}
	credential.Active = activate
	if credential.CreatedAt.IsZero() {
		credential.CreatedAt = time.Now().UTC()
	}
	store, err := LoadCredentials(path)
	if err != nil {
		return err
	}
	replaced := false
	for i := range store.Accounts {
		if activate {
			store.Accounts[i].Active = false
		}
		if store.Accounts[i].Name == credential.Name {
			if !activate {
				credential.Active = store.Accounts[i].Active
			}
			store.Accounts[i] = credential
			replaced = true
		}
	}
	if !replaced {
		store.Accounts = append(store.Accounts, credential)
	}
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}
	data, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return err
	}
	// Replace the file only after the complete credential store has been written.
	// This prevents partial stores and resets insecure permissions on older files.
	// Concurrent writers still require coordination by the caller.
	temp, err := os.CreateTemp(filepath.Dir(path), ".credentials-*")
	if err != nil {
		return err
	}
	defer os.Remove(temp.Name())
	defer temp.Close()
	if err := temp.Chmod(0600); err != nil {
		return err
	}
	if _, err := temp.Write(data); err != nil {
		return err
	}
	if err := temp.Sync(); err != nil {
		return err
	}
	if err := temp.Close(); err != nil {
		return err
	}
	return os.Rename(temp.Name(), path)
}

func ActiveCredential(store CredentialStore) (Credential, bool) {
	for _, account := range store.Accounts {
		if account.Active {
			return account, true
		}
	}
	if len(store.Accounts) == 1 {
		return store.Accounts[0], true
	}
	return Credential{}, false
}

func CredentialByName(store CredentialStore, name string) (Credential, bool) {
	for _, account := range store.Accounts {
		if account.Name == name {
			return account, true
		}
	}
	return Credential{}, false
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type,omitempty"`
	ExpiresIn   int    `json:"expires_in,omitempty"`
}

type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

func NewClient(baseURL string) Client {
	if baseURL == "" {
		baseURL = "https://connect-api.cloud.huawei.com"
	}
	return Client{BaseURL: baseURL, HTTPClient: http.DefaultClient}
}

func (c Client) TokenEndpoint() string {
	return c.BaseURL + "/api/oauth2/v1/token"
}

func ValidateCredential(credential Credential) error {
	switch credential.Mode {
	case "service-account":
		if credential.ServiceAccountFile == "" {
			return errors.New("service account file is required")
		}
	case "api-client":
		if credential.ClientID == "" || credential.ClientKey == "" {
			return errors.New("client id and client key are required")
		}
	default:
		return fmt.Errorf("unsupported credential mode %q", credential.Mode)
	}
	return nil
}

// CredentialView exposes profile metadata without secret material.
type CredentialView struct {
	Name               string    `json:"name"`
	Mode               string    `json:"mode"`
	ServiceAccountFile string    `json:"serviceAccountFile,omitempty"`
	ClientID           string    `json:"clientId,omitempty"`
	Active             bool      `json:"active,omitempty"`
	CreatedAt          time.Time `json:"createdAt"`
}

func (c Credential) View() CredentialView {
	return CredentialView{c.Name, c.Mode, c.ServiceAccountFile, c.ClientID, c.Active, c.CreatedAt}
}

// LoadAPIClientCredential reads credentials without exposing secret values in errors.
func LoadAPIClientCredential(path string) (Credential, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Credential{}, fmt.Errorf("read API client file: %w", err)
	}
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(data, &fields); err != nil || fields == nil {
		return Credential{}, errors.New("API client file must be a JSON object")
	}
	read := func(names ...string) (string, error) {
		value := ""
		for _, name := range names {
			raw, exists := fields[name]
			if !exists {
				continue
			}
			var candidate string
			if string(raw) == "null" || json.Unmarshal(raw, &candidate) != nil || candidate == "" {
				return "", errors.New("API client credential fields must be nonempty strings")
			}
			if value != "" && value != candidate {
				return "", errors.New("API client file contains conflicting credential aliases")
			}
			value = candidate
		}
		return value, nil
	}
	id, err := read("client_id", "clientId")
	if err != nil {
		return Credential{}, err
	}
	key, err := read("client_secret", "clientKey")
	if err != nil {
		return Credential{}, err
	}
	credential := Credential{Mode: "api-client", ClientID: id, ClientKey: key}
	if err := ValidateCredential(credential); err != nil {
		return Credential{}, err
	}
	return credential, nil
}

// CredentialFromEnvironment resolves ephemeral credentials when no saved profile is selected.
func CredentialFromEnvironment() (Credential, bool, error) {
	id := os.Getenv("AGC_CLIENT_ID")
	key := os.Getenv("AGC_CLIENT_KEY")
	secret := os.Getenv("AGC_CLIENT_SECRET")
	if key != "" && secret != "" && key != secret {
		return Credential{}, false, errors.New("AGC_CLIENT_KEY and AGC_CLIENT_SECRET conflict")
	}
	if key == "" {
		key = secret
	}
	if id != "" || key != "" {
		if id == "" || key == "" {
			return Credential{}, false, errors.New("AGC_CLIENT_ID and AGC_CLIENT_KEY or AGC_CLIENT_SECRET must be set together")
		}
		return Credential{Name: "environment", Mode: "api-client", ClientID: id, ClientKey: key}, true, nil
	}
	if file := os.Getenv("AGC_SERVICE_ACCOUNT_FILE"); file != "" {
		return Credential{Name: "environment", Mode: "service-account", ServiceAccountFile: file}, true, nil
	}
	return Credential{}, false, nil
}
