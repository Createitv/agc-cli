package agcapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func ServerlessTokenBaseURL(base string) (string, error) {
	u, err := url.Parse(base)
	if err != nil || u.Host == "" {
		return "", fmt.Errorf("invalid serverless base URL")
	}
	if u.Hostname() == "127.0.0.1" || u.Hostname() == "localhost" {
		return u.Scheme + "://" + u.Host, nil
	}
	if u.Scheme != "https" {
		return "", fmt.Errorf("serverless authentication requires HTTPS")
	}
	region := map[string]string{"ops-server-drcn.agcstorage.link": "", "ops-server-dre.agcstorage.link": "-dre", "ops-server-dra.agcstorage.link": "-dra", "ops-server-drru.agcstorage.link": "-drru"}
	if suffix, ok := region[u.Hostname()]; ok {
		return "https://connect-api" + suffix + ".cloud.huawei.com", nil
	}
	switch u.Hostname() {
	case "connect-api.cloud.huawei.com", "connect-api-dre.cloud.huawei.com", "connect-api-dra.cloud.huawei.com", "connect-api-drru.cloud.huawei.com":
		return u.Scheme + "://" + u.Host, nil
	}
	return "", fmt.Errorf("automatic serverless authorization only supports reviewed Huawei service hosts; use an explicit token for an explicitly selected custom host")
}
func ServerlessClientToken(ctx context.Context, client *http.Client, base string, credential Credential) (string, error) {
	if credential.Mode != "api-client" {
		return "", fmt.Errorf("client REST requires an API client profile")
	}
	u, err := url.Parse(base)
	if err != nil || u.Host == "" {
		return "", fmt.Errorf("invalid client token base URL")
	}
	if u.Hostname() != "connect-drcn.dbankcloud.cn" && u.Hostname() != "localhost" && u.Hostname() != "127.0.0.1" {
		return "", fmt.Errorf("client token region is not documented in this reference; use an explicitly acquired token")
	}
	if u.Hostname() == "connect-drcn.dbankcloud.cn" && u.Scheme != "https" {
		return "", fmt.Errorf("client token requires HTTPS")
	}
	b, _ := json.Marshal(map[string]any{"useJwt": 1, "grant_type": "client_credentials", "client_id": credential.ClientID, "client_secret": credential.ClientKey})
	req, err := http.NewRequestWithContext(ctx, "POST", strings.TrimRight(base, "/")+"/agc/apigw/oauth2/v1/token", bytes.NewReader(b))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	if client == nil {
		client = http.DefaultClient
	}
	clone := *client
	clone.CheckRedirect = func(*http.Request, []*http.Request) error { return fmt.Errorf("token redirects refused") }
	resp, err := clone.Do(req)
	if err != nil {
		return "", fmt.Errorf("client token transport failed")
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return "", err
	}
	var token struct {
		AccessToken string `json:"access_token"`
	}
	if resp.StatusCode != 200 || json.Unmarshal(data, &token) != nil || token.AccessToken == "" {
		return "", fmt.Errorf("client token authorization failed (HTTP %d); verify the API client, project region and enabled service", resp.StatusCode)
	}
	return token.AccessToken, nil
}
