package agcapi

import (
	"encoding/json"
	"io"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/Createitv/agc-cli/pkg/domain"
)

// ManagementSession is explicit and never read from the publishing profile or
// the official CLI's encrypted session store. Cookie transport is source-based.
type ManagementSession struct {
	UID                   string    `json:"uid"`
	TeamID                string    `json:"teamId"`
	ProductID             string    `json:"productId"`
	AccessToken           string    `json:"accessToken"`
	AuthInfo              string    `json:"authInfo"`
	CSRFToken             string    `json:"csrfToken"`
	ExpiresAt             time.Time `json:"expiresAt"`
	ProjectBaseURL        string    `json:"projectBaseUrl"`
	ProjectEdgeBaseURL    string    `json:"projectEdgeBaseUrl"`
	TeamBaseURL           string    `json:"teamBaseUrl"`
	ProjectStorageBaseURL string    `json:"projectStorageBaseUrl"`
}

func IsManagementEndpoint(e domain.Endpoint) bool {
	return e.Authentication == "agcli-management-oauth" || e.Authentication == "agcli-management-cookie"
}

func managementSessionRequest(req InvokeRequest) (InvokeRequest, error) {
	if req.DryRun {
		if req.BaseURL == "" {
			req.BaseURL = "https://session-required.invalid"
		}
		req.AccessToken = ""
		return req, nil
	}
	fail := func(message string) (InvokeRequest, error) { return req, &AuthenticationError{Message: message} }
	if req.SessionFile == "" {
		return fail("An explicit --session-file (REST sessionFile) with an unexpired user session is required; publishing profiles are not compatible.")
	}
	f, err := os.Open(req.SessionFile)
	if err != nil {
		return fail("Cannot read explicitly selected management session file.")
	}
	defer f.Close()
	b, err := io.ReadAll(io.LimitReader(f, 65537))
	if err != nil || len(b) > 65536 {
		return fail("Management session file exceeds the 64 KiB local limit or cannot be read.")
	}
	var s ManagementSession
	if json.Unmarshal(b, &s) != nil || s.UID == "" || s.TeamID == "" || s.ProductID == "" || s.AccessToken == "" || !s.ExpiresAt.After(time.Now()) {
		return fail("Explicit user session must include uid, teamId, productId, accessToken and a future RFC3339 expiresAt.")
	}
	if req.ProjectID != "" && req.ProjectID != s.ProductID {
		return fail("Selected projectId does not match management session productId.")
	}
	if req.Query["projectId"] != "" && req.Query["projectId"] != s.ProductID {
		return fail("query.projectId does not match management session productId.")
	}
	base := s.ProjectBaseURL
	if req.Endpoint.Authentication == "agcli-management-cookie" {
		base = s.ProjectEdgeBaseURL
	}
	if req.Endpoint.FamilyID == "templates" {
		base = s.TeamBaseURL
	}
	if req.Endpoint.ID == "template-code-upload" {
		base = s.ProjectStorageBaseURL
	}
	if req.BaseURL != "" {
		base = req.BaseURL
	}
	u, err := url.Parse(base)
	if err != nil || u.Host == "" || u.User != nil || u.RawQuery != "" || u.Fragment != "" || (u.Scheme != "https" && !(u.Scheme == "http" && (u.Hostname() == "127.0.0.1" || u.Hostname() == "localhost" || u.Hostname() == "::1"))) {
		return fail("Management session requires the source-defined HTTPS base URL; only loopback HTTP is supported for offline testing.")
	}
	req.BaseURL = base
	req.AccessToken = ""
	headers := map[string]string{}
	for k, v := range req.Headers {
		headers[k] = v
	}
	for k := range headers {
		switch strings.ToLower(k) {
		case "authorization", "oauth2token", "cookie", "x-hd-csrf", "uid", "teamid", "productid", "projectid", "agcteamid":
			delete(headers, k)
		}
	}
	for _, v := range []string{s.UID, s.TeamID, s.ProductID, s.AccessToken, s.AuthInfo, s.CSRFToken} {
		if strings.ContainsAny(v, "\r\n") {
			return fail("Management session contains invalid header characters.")
		}
	}
	headers["uid"] = s.UID
	headers["teamId"] = s.TeamID
	headers["productId"] = s.ProductID
	headers["oauth2Token"] = s.AccessToken
	headers["X-HD-DATE"] = time.Now().UTC().Format("20060102T150405Z")
	if req.Endpoint.Authentication == "agcli-management-cookie" {
		if s.AuthInfo == "" || s.CSRFToken == "" || strings.ContainsAny(s.AuthInfo+s.CSRFToken, ";") {
			return fail("This vendor transport requires explicit authInfo and csrfToken session fields; they are not inferred from an OAuth token.")
		}
		headers["projectId"] = s.ProductID
		headers["agcTeamId"] = s.TeamID
		headers["Cookie"] = "authInfo=" + s.AuthInfo + ";csrfToken=" + s.CSRFToken + ";"
		headers["X-HD-CSRF"] = s.CSRFToken
	}
	req.Headers = headers
	return req, nil
}

func redactManagementResponse(data []byte, headers map[string]string, requestBody []byte) []byte {
	// Replace credential values even if an upstream error echoes them. Never
	// expose raw HTML/text replies from an authenticated management endpoint.
	var root any
	dec := json.NewDecoder(strings.NewReader(string(data)))
	dec.UseNumber()
	if dec.Decode(&root) != nil {
		return []byte(`{"message":"Non-JSON management response omitted"}`)
	}
	secrets := []string{}
	var input map[string]any
	if json.Unmarshal(requestBody, &input) == nil {
		for _, key := range []string{"password", "accessToken", "authInfo", "csrfToken"} {
			if value, ok := input[key].(string); ok && value != "" {
				secrets = append(secrets, value)
			}
		}
	}
	for k, v := range headers {
		if strings.EqualFold(k, "oauth2Token") || strings.EqualFold(k, "Cookie") || strings.EqualFold(k, "X-HD-CSRF") {
			secrets = append(secrets, v)
			if strings.EqualFold(k, "Cookie") {
				for _, pair := range strings.Split(v, ";") {
					_, value, ok := strings.Cut(pair, "=")
					if ok && value != "" {
						secrets = append(secrets, value)
					}
				}
			}
		}
	}
	var walk func(any) any
	walk = func(v any) any {
		switch x := v.(type) {
		case map[string]any:
			for k, c := range x {
				switch strings.ToLower(k) {
				case "accesstoken", "authinfo", "csrftoken", "password", "authorization", "cookie", "oauth2token", "environment", "env":
					delete(x, k)
				default:
					x[k] = walk(c)
				}
			}
			return x
		case []any:
			for i, c := range x {
				x[i] = walk(c)
			}
			return x
		case string:
			for _, secret := range secrets {
				if secret != "" {
					x = strings.ReplaceAll(x, secret, "[redacted]")
				}
			}
			return x
		}
		return v
	}
	out, _ := json.Marshal(walk(root))
	return out
}
