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
	if credential.Name == "" {
		credential.Name = "default"
	}
	if credential.Mode == "" {
		return errors.New("credential mode is required")
	}
	credential.Active = true
	if credential.CreatedAt.IsZero() {
		credential.CreatedAt = time.Now().UTC()
	}
	store, err := LoadCredentials(path)
	if err != nil {
		return err
	}
	replaced := false
	for i := range store.Accounts {
		store.Accounts[i].Active = false
		if store.Accounts[i].Name == credential.Name {
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
	return os.WriteFile(path, data, 0600)
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
