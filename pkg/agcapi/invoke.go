package agcapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"

	"github.com/Createitv/agc-cli/pkg/domain"
)

type InvokeRequest struct {
	Endpoint    domain.Endpoint   `json:"endpoint"`
	BaseURL     string            `json:"baseUrl"`
	Params      map[string]string `json:"params,omitempty"`
	Query       map[string]string `json:"query,omitempty"`
	Headers     map[string]string `json:"headers,omitempty"`
	Body        []byte            `json:"-"`
	AccessToken string            `json:"-"`
	DryRun      bool              `json:"dryRun"`
}

type InvokeResponse struct {
	Method      string          `json:"method"`
	URL         string          `json:"url"`
	StatusCode  int             `json:"statusCode,omitempty"`
	ContentType string          `json:"contentType,omitempty"`
	Body        json.RawMessage `json:"body,omitempty"`
	RawBody     []byte          `json:"-"`
	DryRun      bool            `json:"dryRun"`
}

func InvokeEndpoint(ctx context.Context, client *http.Client, req InvokeRequest) (InvokeResponse, error) {
	if client == nil {
		client = http.DefaultClient
	}
	url, err := BuildEndpointURL(req.BaseURL, req.Endpoint.Path, req.Params, req.Query)
	if err != nil {
		return InvokeResponse{}, err
	}
	result := InvokeResponse{Method: req.Endpoint.Method, URL: url, DryRun: req.DryRun}
	if req.DryRun {
		return result, nil
	}
	httpReq, err := http.NewRequestWithContext(ctx, req.Endpoint.Method, url, bytes.NewReader(req.Body))
	if err != nil {
		return InvokeResponse{}, err
	}
	if req.AccessToken != "" {
		httpReq.Header.Set("Authorization", "Bearer "+req.AccessToken)
	}
	for key, value := range req.Headers {
		if value != "" {
			httpReq.Header.Set(key, value)
		}
	}
	if len(req.Body) > 0 {
		httpReq.Header.Set("Content-Type", "application/json")
	}
	resp, err := client.Do(httpReq)
	if err != nil {
		return InvokeResponse{}, err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return InvokeResponse{}, err
	}
	result.StatusCode = resp.StatusCode
	result.ContentType = resp.Header.Get("Content-Type")
	result.RawBody = data
	if len(data) > 0 && json.Valid(data) {
		result.Body = data
	} else if len(data) > 0 {
		encoded, _ := json.Marshal(string(data))
		result.Body = encoded
	}
	if resp.StatusCode >= 400 {
		return result, ParseAGCError(resp.StatusCode, data)
	}
	return result, nil
}

func BuildEndpointURL(baseURL, path string, params, query map[string]string) (string, error) {
	if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
		baseURL = ""
	}
	if strings.Contains(path, "{uploadUrl}") {
		uploadURL := params["uploadUrl"]
		if uploadURL == "" {
			return "", fmt.Errorf("missing required path param uploadUrl")
		}
		path = uploadURL
		baseURL = ""
	}
	for key, value := range params {
		path = strings.ReplaceAll(path, "{"+key+"}", value)
	}
	if strings.Contains(path, "{") {
		return "", fmt.Errorf("missing required path param in %s", path)
	}
	if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
		if len(query) == 0 {
			return path, nil
		}
		return appendQuery(path, query), nil
	}
	if baseURL == "" {
		baseURL = "https://connect-api.cloud.huawei.com"
	}
	fullURL := strings.TrimRight(baseURL, "/") + "/" + strings.TrimLeft(path, "/")
	if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
		fullURL = path
	}
	if len(query) == 0 {
		return fullURL, nil
	}
	return appendQuery(fullURL, query), nil
}

func appendQuery(fullURL string, query map[string]string) string {
	separator := "?"
	if strings.Contains(fullURL, "?") {
		separator = "&"
	}
	parts := make([]string, 0, len(query))
	for key, value := range query {
		parts = append(parts, url.QueryEscape(key)+"="+url.QueryEscape(value))
	}
	sort.Strings(parts)
	return fullURL + separator + strings.Join(parts, "&")
}
