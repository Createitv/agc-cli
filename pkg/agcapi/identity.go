package agcapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// TeamListResponse is returned by Huawei's OAuth-only team account endpoint.
// Team objects remain maps because Huawei does not publish a stable field schema
// for the Teams type in the endpoint overview.
type TeamListResponse struct {
	Ret   json.RawMessage  `json:"ret,omitempty"`
	Teams []map[string]any `json:"teams"`
}

func UserTeamList(ctx context.Context, client *http.Client, baseURL, oauth2Token string) (TeamListResponse, error) {
	if oauth2Token == "" {
		return TeamListResponse{}, fmt.Errorf("OAuth access token is required")
	}
	if client == nil {
		client = http.DefaultClient
	}
	if baseURL == "" {
		baseURL = "https://connect-api.cloud.huawei.com"
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, strings.TrimRight(baseURL, "/")+"/api/ups/user-permission-service/v1/user-team-list", nil)
	if err != nil {
		return TeamListResponse{}, err
	}
	req.Header.Set("oauth2Token", oauth2Token)
	resp, err := client.Do(req)
	if err != nil {
		return TeamListResponse{}, err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return TeamListResponse{}, err
	}
	if resp.StatusCode >= 400 {
		return TeamListResponse{}, fmt.Errorf("team account endpoint returned HTTP %d", resp.StatusCode)
	}
	var result TeamListResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return TeamListResponse{}, fmt.Errorf("decode team account response: %w", err)
	}
	if result.Teams == nil {
		result.Teams = []map[string]any{}
	}
	return result, nil
}
