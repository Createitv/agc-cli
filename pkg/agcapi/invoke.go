package agcapi

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strings"

	"github.com/Createitv/agc-cli/pkg/domain"
	"github.com/Createitv/agc-cli/pkg/hvigor"
)

type InvokeRequest struct {
	SessionFile string              `json:"sessionFile,omitempty"`
	ProjectID   string              `json:"projectId,omitempty"`
	QueryValues map[string][]string `json:"queryValues,omitempty"`
	Endpoint    domain.Endpoint     `json:"endpoint"`
	BaseURL     string              `json:"baseUrl"`
	Params      map[string]string   `json:"params,omitempty"`
	Query       map[string]string   `json:"query,omitempty"`
	Headers     map[string]string   `json:"headers,omitempty"`
	FilePath    string              `json:"filePath,omitempty"`
	Body        []byte              `json:"-"`
	AccessToken string              `json:"-"`
	DryRun      bool                `json:"dryRun"`
}

type InvokeResponse struct {
	ETag        string          `json:"etag,omitempty"`
	BodyBase64  string          `json:"bodyBase64,omitempty"`
	Method      string          `json:"method"`
	URL         string          `json:"url"`
	StatusCode  int             `json:"statusCode,omitempty"`
	ContentType string          `json:"contentType,omitempty"`
	Body        json.RawMessage `json:"body,omitempty"`
	RawBody     []byte          `json:"-"`
	DryRun      bool            `json:"dryRun"`
}

func InvokeEndpoint(ctx context.Context, client *http.Client, req InvokeRequest) (InvokeResponse, error) {
	if req.Endpoint.Authentication == "agcli-user-oauth" {
		return InvokeResponse{}, fmt.Errorf("native function workflows require the cloudfunctions client and an explicit OAuth session")
	}
	if err := ValidateEndpointPathInput(req.Endpoint, req.Params); err != nil {
		return InvokeResponse{}, err
	}
	if domain.IsServerlessFamily(req.Endpoint.FamilyID) {
		if err := ValidateEndpointRequest(req.Endpoint, req.Query, req.Headers, req.Body, req.QueryValues); err != nil {
			return InvokeResponse{}, err
		}
	}
	if err := ValidateEndpointFileInput(req.Endpoint, req.FilePath, req.Headers); err != nil {
		return InvokeResponse{}, err
	}
	if req.Endpoint.FamilyID == "templates" && req.Endpoint.ID == "deploy" {
		return invokeTemplateDeploy(ctx, client, req)
	}
	if IsManagementEndpoint(req.Endpoint) {
		var err error
		req, err = managementSessionRequest(req)
		if err != nil {
			return InvokeResponse{}, err
		}
	}
	if req.Endpoint.Direction == "local" {
		if req.Endpoint.FamilyID == "prefetch" && req.Endpoint.ID == "device-getcache" {
			return invokePrefetchDevice(ctx, req)
		}
		result := InvokeResponse{Method: req.Endpoint.Method, URL: req.Endpoint.Path, DryRun: req.DryRun}
		if req.Endpoint.FamilyID != "cicd" || req.Endpoint.ID != "build-hvigor" {
			return result, fmt.Errorf("unsupported local endpoint")
		}
		var build hvigor.Request
		if err := json.Unmarshal(req.Body, &build); err != nil {
			return result, fmt.Errorf("local build requires a JSON request")
		}
		if build.Project == "" {
			return result, fmt.Errorf("project directory required")
		}
		if req.DryRun {
			return result, nil
		}
		artifact, err := hvigor.Build(ctx, build)
		if err != nil {
			return result, err
		}
		result.StatusCode, result.ContentType = 200, "application/json"
		result.Body, err = json.Marshal(artifact)
		result.RawBody = result.Body
		return result, err
	}
	if client == nil {
		client = http.DefaultClient
	}
	if req.BaseURL == "" && req.Endpoint.DefaultBaseURL != "" {
		req.BaseURL = req.Endpoint.DefaultBaseURL
	}
	endpointURL, err := BuildEndpointURLWithValues(req.BaseURL, req.Endpoint.Path, req.Params, req.Query, req.QueryValues)
	if err != nil {
		return InvokeResponse{}, err
	}
	if domain.IsServerlessFamily(req.Endpoint.FamilyID) {
		u, parseErr := url.Parse(endpointURL)
		loopback := parseErr == nil && (u.Hostname() == "127.0.0.1" || u.Hostname() == "localhost" || u.Hostname() == "::1")
		if parseErr != nil || u.Host == "" || u.User != nil || (u.Scheme != "https" && !(u.Scheme == "http" && loopback)) {
			field := "baseUrl"
			if req.Endpoint.Direction == "inbound-callback" {
				field = "params.callbackUrl"
			}
			if req.Endpoint.Authentication == "signed-upload" {
				field = "params.uploadUrl"
			}
			return InvokeResponse{}, uploadIssue(req.Endpoint, field, "Serverless and developer callback URLs require HTTPS; only loopback HTTP is allowed for isolated tests, and URL user credentials are forbidden.")
		}
	}
	result := InvokeResponse{Method: req.Endpoint.Method, URL: endpointURL, DryRun: req.DryRun}
	if req.Endpoint.FamilyID == "prefetch" {
		u, _ := url.Parse(result.URL)
		q := u.Query()
		if q.Has("token") {
			q.Set("token", "[redacted]")
		}
		u.RawQuery = q.Encode()
		result.URL = u.String()
	}
	signedUpload := req.Endpoint.Authentication == "signed-upload" || (req.Endpoint.Method == http.MethodPut && req.Endpoint.Path == "{uploadUrl}" && (req.Endpoint.ID == "upload-file-new" || req.Endpoint.ID == "obbfile-upload"))
	storageUpload := req.Endpoint.FamilyID == "storage" && req.Endpoint.ID == "upload-file" || req.Endpoint.ID == "template-code-upload"
	// Publishing JSON operations accept an empty object, not an absent entity.
	if req.Endpoint.FamilyID == "publishing" && req.Endpoint.Body != "" && len(req.Body) == 0 && !signedUpload {
		req.Body = []byte("{}")
	}
	// Vendor schema import passes a JSON file's raw string to Axios, which sets
	// its default form content type but does not form-encode that string.
	rawSchemaImport := IsManagementEndpoint(req.Endpoint) && req.Endpoint.FamilyID == "clouddb" && req.Endpoint.ID == "import-object-types"
	if strings.HasPrefix(req.Endpoint.RequestContentType, "application/x-www-form-urlencoded") && !rawSchemaImport && req.Endpoint.ID != "template-code-upload" {
		decoder := json.NewDecoder(bytes.NewReader(req.Body))
		decoder.UseNumber()
		var fields map[string]any
		if err := decoder.Decode(&fields); err != nil || fields == nil {
			return result, fmt.Errorf("callback form requires a JSON object of documented fields")
		}
		form := url.Values{}
		for name, value := range fields {
			switch typed := value.(type) {
			case string:
				form.Set(name, typed)
			case json.Number:
				form.Set(name, typed.String())
			default:
				return result, fmt.Errorf("callback form field %s must be a string or number", name)
			}
		}
		req.Body = []byte(form.Encode())
		result.ContentType = req.Endpoint.RequestContentType
	}
	var body io.Reader = bytes.NewReader(req.Body)
	var fileSize int64
	multipartResource := req.Endpoint.FamilyID == "pms" && req.Endpoint.ID == "resourcemanage-harmonyosnext"
	if multipartResource {
		if req.FilePath == "" {
			return result, fmt.Errorf("multipart resource upload requires a file")
		}
		file, size, contentType, uploadErr := preparePMSResourceUpload(req.FilePath, req.Body)
		if uploadErr != nil {
			return result, uploadErr
		}
		defer os.Remove(file.Name())
		defer file.Close()
		body, fileSize = file, size
		result.ContentType = contentType
	} else if req.FilePath != "" {
		if !signedUpload && !storageUpload {
			return result, fmt.Errorf("file upload is unsupported for endpoint %s", req.Endpoint.ID)
		}
		if len(req.Body) > 0 {
			return result, fmt.Errorf("file upload cannot also specify a JSON body")
		}
		file, openErr := os.Open(req.FilePath)
		if openErr != nil {
			return result, fmt.Errorf("open upload file: %w", openErr)
		}
		defer file.Close()
		info, statErr := file.Stat()
		if statErr != nil {
			return result, fmt.Errorf("stat upload file: %w", statErr)
		}
		if !info.Mode().IsRegular() {
			return result, fmt.Errorf("upload file must be a regular file")
		}
		body, fileSize = file, info.Size()
	}
	if req.DryRun {
		return result, nil
	}
	if domain.IsServerlessFamily(req.Endpoint.FamilyID) {
		switch req.Endpoint.Authentication {
		case "server-api-client", "client-api-token", "service-account-or-api-client":
			if req.AccessToken == "" && uploadHeader(req.Headers, "Authorization") == "" {
				return result, &AuthenticationError{Message: "Explicit service-compatible authorization is required."}
			}
			if req.Endpoint.Authentication != "service-account-or-api-client" && uploadHeader(req.Headers, "client_id") == "" {
				return result, &AuthenticationError{Message: "API client authorization requires the matching client_id header."}
			}
		}
	}
	httpReq, err := http.NewRequestWithContext(ctx, req.Endpoint.Method, endpointURL, body)
	if err != nil {
		return InvokeResponse{}, err
	}
	if req.AccessToken != "" && !signedUpload && req.Endpoint.Direction != "inbound-callback" && req.Endpoint.Authentication != "body-signature" && req.Endpoint.Authentication != "body-client-credentials" {
		httpReq.Header.Set("Authorization", "Bearer "+req.AccessToken)
	}
	for key, value := range req.Headers {
		if value != "" {
			if strings.EqualFold(key, "Host") {
				httpReq.Host = value
			} else {
				httpReq.Header.Set(key, value)
			}
		}
	}
	if req.FilePath != "" {
		httpReq.ContentLength = fileSize
		if fileSize == 0 {
			httpReq.Body = http.NoBody
		}
	}
	if req.Endpoint.RequestContentType != "" {
		httpReq.Header.Set("Content-Type", req.Endpoint.RequestContentType)
	}
	if multipartResource {
		httpReq.Header.Set("Content-Type", result.ContentType)
	}
	if httpReq.Header.Get("Content-Type") == "" {
		if signedUpload {
			httpReq.Header.Set("Content-Type", "application/octet-stream")
		} else if len(req.Body) > 0 {
			httpReq.Header.Set("Content-Type", "application/json")
		}
	}
	if domain.IsServerlessFamily(req.Endpoint.FamilyID) {
		// Authenticated redirects cannot change origin (including signed uploads).
		clone := *client
		previous := clone.CheckRedirect
		clone.CheckRedirect = func(next *http.Request, via []*http.Request) error {
			if len(via) > 0 && (next.URL.Scheme != via[0].URL.Scheme || next.URL.Host != via[0].URL.Host) {
				return fmt.Errorf("cross-origin redirect refused for serverless credentials")
			}
			if previous != nil {
				return previous(next, via)
			}
			if len(via) >= 10 {
				return fmt.Errorf("too many redirects")
			}
			return nil
		}
		client = &clone
	}
	resp, err := client.Do(httpReq)
	if err != nil {
		return InvokeResponse{}, err
	}
	defer resp.Body.Close()
	var responseReader io.Reader = resp.Body
	responseLimit := int64(0)
	if req.Endpoint.FamilyID == "prefetch" {
		responseLimit = 3 << 20
		if req.Endpoint.ID == "fetch-install" {
			responseLimit = 2 << 20
		}
		responseReader = io.LimitReader(resp.Body, responseLimit+1)
	}
	if IsManagementEndpoint(req.Endpoint) {
		responseLimit = 32 << 20
		responseReader = io.LimitReader(resp.Body, responseLimit+1)
	}
	data, err := io.ReadAll(responseReader)
	if err != nil {
		return InvokeResponse{}, err
	}
	if responseLimit > 0 && int64(len(data)) > responseLimit {
		return InvokeResponse{}, uploadIssue(req.Endpoint, "response.body", "Response exceeds the documented prefetch quota or local management response limit; refusing truncated success.")
	}
	if IsManagementEndpoint(req.Endpoint) {
		data = redactManagementResponse(data, req.Headers, req.Body)
	}
	result.StatusCode = resp.StatusCode
	result.ContentType = resp.Header.Get("Content-Type")
	// 只暴露内容校验值，不把认证头或签名地址复制到响应元数据中。
	result.ETag = resp.Header.Get("ETag")
	result.RawBody = data
	if req.Endpoint.FamilyID == "storage" && req.Endpoint.ID == "download-file" && resp.StatusCode < 400 {
		result.BodyBase64 = base64.StdEncoding.EncodeToString(data)
	} else if len(data) > 0 && json.Valid(data) {
		result.Body = data
	} else if len(data) > 0 {
		encoded, _ := json.Marshal(string(data))
		result.Body = encoded
	}
	if resp.StatusCode >= 400 || (!(req.Endpoint.FamilyID == "storage" && req.Endpoint.ID == "download-file") && req.Endpoint.FamilyID != "prefetch" && businessFailure(data)) {
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
		if key != "uploadUrl" && key != "callbackUrl" {
			if key == "object_name" {
				parts := strings.Split(value, "/")
				for i := range parts {
					parts[i] = url.PathEscape(parts[i])
				}
				value = strings.Join(parts, "/")
			} else {
				value = url.PathEscape(value)
			}
		}
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

// BuildEndpointURLWithValues preserves ordered repeated keys from official
// report examples. Scalar query entries remain compatible with existing clients.
func BuildEndpointURLWithValues(baseURL, path string, params, query map[string]string, values map[string][]string) (string, error) {
	full, err := BuildEndpointURL(baseURL, path, params, query)
	if err != nil || len(values) == 0 {
		return full, err
	}
	parsed, err := url.Parse(full)
	if err != nil {
		return "", err
	}
	encoded := parsed.Query()
	for key, entries := range values {
		encoded.Del(key)
		for _, value := range entries {
			encoded.Add(key, value)
		}
	}
	parsed.RawQuery = encoded.Encode()
	return parsed.String(), nil
}

// ScalarQuery supplies required-field and scalar-rule validation with the
// first value of each ordered query list, without modifying the caller's maps.
func ScalarQuery(query map[string]string, values map[string][]string) map[string]string {
	result := make(map[string]string, len(query)+len(values))
	for key, value := range query {
		result[key] = value
	}
	for key, entries := range values {
		if len(entries) > 0 {
			result[key] = entries[0]
		} else {
			delete(result, key)
		}
	}
	return result
}
