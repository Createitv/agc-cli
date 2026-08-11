package agcapi

import (
	"bytes"
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

const serviceAccountAudience = "https://oauth-login.cloud.huawei.com/oauth2/v3/token"

type ServiceAccountFile struct {
	ProjectID  string `json:"project_id"`
	KeyID      string `json:"key_id"`
	PrivateKey string `json:"private_key"`
	SubAccount string `json:"sub_account"`
	TokenURI   string `json:"token_uri"`
}

func AccessToken(ctx context.Context, client *http.Client, baseURL string, credential Credential) (TokenResponse, error) {
	switch credential.Mode {
	case "service-account":
		account, err := LoadServiceAccount(credential.ServiceAccountFile)
		if err != nil {
			return TokenResponse{}, err
		}
		jwt, err := ServiceAccountJWT(account, time.Now().UTC())
		if err != nil {
			return TokenResponse{}, err
		}
		return TokenResponse{AccessToken: jwt, TokenType: "Bearer", ExpiresIn: 3600}, nil
	case "api-client":
		return APIClientToken(ctx, client, baseURL, credential.ClientID, credential.ClientKey)
	default:
		return TokenResponse{}, fmt.Errorf("unsupported credential mode %q", credential.Mode)
	}
}

func LoadServiceAccount(path string) (ServiceAccountFile, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return ServiceAccountFile{}, err
	}
	var account ServiceAccountFile
	if err := json.Unmarshal(data, &account); err != nil {
		return ServiceAccountFile{}, err
	}
	if account.KeyID == "" || account.PrivateKey == "" || account.SubAccount == "" {
		return ServiceAccountFile{}, fmt.Errorf("service account file must include key_id, private_key, and sub_account")
	}
	if account.TokenURI == "" {
		account.TokenURI = serviceAccountAudience
	}
	return account, nil
}

func ServiceAccountJWT(account ServiceAccountFile, now time.Time) (string, error) {
	header := map[string]string{"kid": account.KeyID, "typ": "JWT", "alg": "PS256"}
	payload := map[string]any{
		"aud": serviceAccountAudience,
		"iss": account.SubAccount,
		"iat": now.Unix(),
		"exp": now.Add(time.Hour).Unix(),
	}
	headerJSON, err := json.Marshal(header)
	if err != nil {
		return "", err
	}
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	signingInput := base64.RawURLEncoding.EncodeToString(headerJSON) + "." + base64.RawURLEncoding.EncodeToString(payloadJSON)
	privateKey, err := parseRSAPrivateKey(account.PrivateKey)
	if err != nil {
		return "", err
	}
	digest := sha256.Sum256([]byte(signingInput))
	signature, err := rsa.SignPSS(rand.Reader, privateKey, crypto.SHA256, digest[:], &rsa.PSSOptions{SaltLength: rsa.PSSSaltLengthEqualsHash, Hash: crypto.SHA256})
	if err != nil {
		return "", err
	}
	return signingInput + "." + base64.RawURLEncoding.EncodeToString(signature), nil
}

func parseRSAPrivateKey(privateKey string) (*rsa.PrivateKey, error) {
	clean := strings.TrimSpace(privateKey)
	block, _ := pem.Decode([]byte(clean))
	var der []byte
	if block != nil {
		der = block.Bytes
	} else {
		decoded, err := base64.StdEncoding.DecodeString(clean)
		if err != nil {
			return nil, err
		}
		der = decoded
	}
	key, err := x509.ParsePKCS8PrivateKey(der)
	if err != nil {
		if rsaKey, rsaErr := x509.ParsePKCS1PrivateKey(der); rsaErr == nil {
			return rsaKey, nil
		}
		return nil, err
	}
	rsaKey, ok := key.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("private key is %T, want RSA private key", key)
	}
	return rsaKey, nil
}

func APIClientToken(ctx context.Context, client *http.Client, baseURL, clientID, clientKey string) (TokenResponse, error) {
	if client == nil {
		client = http.DefaultClient
	}
	if baseURL == "" {
		baseURL = "https://connect-api.cloud.huawei.com"
	}
	body, err := json.Marshal(map[string]string{
		"grant_type":    "client_credentials",
		"client_id":     clientID,
		"client_secret": clientKey,
	})
	if err != nil {
		return TokenResponse{}, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, strings.TrimRight(baseURL, "/")+"/api/oauth2/v1/token", bytes.NewReader(body))
	if err != nil {
		return TokenResponse{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return TokenResponse{}, err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return TokenResponse{}, err
	}
	if resp.StatusCode >= 400 {
		return TokenResponse{}, fmt.Errorf("token endpoint returned HTTP %d: %s", resp.StatusCode, string(data))
	}
	var token TokenResponse
	if err := json.Unmarshal(data, &token); err != nil {
		return TokenResponse{}, err
	}
	if token.AccessToken == "" {
		return TokenResponse{}, fmt.Errorf("token endpoint returned no access_token")
	}
	return token, nil
}
