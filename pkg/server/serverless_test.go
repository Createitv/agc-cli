package server

import (
	"archive/zip"
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/Createitv/agc-cli/pkg/agcapi"
	"github.com/Createitv/agc-cli/pkg/domain"
)

func slValue(p domain.Parameter) any {
	if p.Name == "fileSha256" {
		return strings.Repeat("0", 64)
	}
	if len(p.Enum) > 0 {
		if p.Type == "integer" {
			return json.Number(p.Enum[0])
		}
		return p.Enum[0]
	}
	if p.EncodedJSON != nil {
		b, _ := json.Marshal(slValue(*p.EncodedJSON))
		return string(b)
	}
	switch p.Type {
	case "integer", "number":
		return 1
	case "boolean":
		return false
	case "object":
		m := map[string]any{}
		for _, c := range p.Properties {
			if c.Required {
				m[c.Name] = slValue(c)
			}
		}
		return m
	case "array":
		if p.ItemsType == "object" {
			m := map[string]any{}
			for _, c := range p.ItemsProperties {
				if c.Required {
					m[c.Name] = slValue(c)
				}
			}
			return []any{m}
		}
		return []any{"test"}
	}
	return "test"
}
func TestEveryServerlessOperationREST(t *testing.T) {
	for _, e := range domain.AllEndpoints() {
		if !domain.IsServerlessFamily(e.FamilyID) {
			continue
		}
		if e.FamilyID == "prefetch" && e.Direction == "local" {
			testPrefetchDeviceREST(t, e)
			continue
		}
		if agcapi.IsManagementEndpoint(e) {
			testReviewedManagementREST(t, e)
			continue
		}
		t.Run(e.FamilyID+"/"+e.ID, func(t *testing.T) {
			seen := 0
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				seen++
				w.Header().Set("Content-Type", "application/json")
				w.Write([]byte(`{"ret":{"code":0},"function":{"revisionId":"r1","handler":"handler.myHandler"}}`))
			}))
			defer srv.Close()
			params, query, headers := map[string]string{}, map[string]string{}, map[string]string{"client_id": "offline"}
			body := map[string]any{}
			for _, p := range e.Parameters {
				if !p.Required {
					continue
				}
				v := slValue(p)
				switch p.In {
				case "body":
					body[p.Name] = v
				case "header":
					headers[p.Name] = fmt.Sprint(v)
				case "query":
					query[p.Name] = fmt.Sprint(v)
				case "path":
					value := fmt.Sprint(v)
					if p.Name == "uploadUrl" || p.Name == "callbackUrl" {
						value = srv.URL + "/test"
					}
					params[p.Name] = value
				}
			}
			if e.FamilyID == "cloud-auth" {
				switch e.ID {
				case "request-code", "register":
					body["email"] = "test@example.test"
				case "update-account":
					body["newEmail"] = "test@example.test"
				}
			}
			if e.FamilyID == "clouddb" && (e.ID == "update-data" || e.ID == "delete-data") {
				body["schemas"] = []any{}
				body["opData"] = []any{}
			}
			payload := map[string]any{"baseUrl": srv.URL, "params": params, "query": query, "headers": headers, "body": body, "dryRun": false, "token": "offline"}
			if e.ID == "upload-file" || e.ID == "upload-version" {
				p := filepath.Join(t.TempDir(), "test.zip")
				os.WriteFile(p, []byte("offline"), 0600)
				payload["filePath"] = p
			}
			if e.Authentication == "signed-upload" {
				sum := sha256.Sum256([]byte("offline"))
				headers["Authorization"] = "AWS4-HMAC-SHA256 offline"
				headers["x-amz-content-sha256"] = hex.EncodeToString(sum[:])
				headers["x-amz-client-request-id"] = "offline"
				headers["x-amz-date"] = "20261026T000000Z"
				headers["Host"] = strings.TrimPrefix(srv.URL, "http://")
				headers["Content-Type"] = "application/zip"
				delete(payload, "body")
			}
			if e.ID == "upload-file" {
				delete(payload, "body")
			}
			if e.Authentication == "agcli-user-oauth" {
				dir := t.TempDir()
				session := filepath.Join(dir, "session.json")
				b, _ := json.Marshal(map[string]any{"uid": "offline", "teamId": "offline", "productId": "offline", "accessToken": "offline-oauth", "expiresAt": time.Now().Add(time.Hour), "projectBaseUrl": srv.URL})
				os.WriteFile(session, b, 0600)
				payload = map[string]any{"name": "test-func", "sessionFile": session, "dryRun": false, "versionDescription": "test"}
				if e.ID == "create-function" || e.ID == "update-function-code" || e.ID == "update-and-publish-function" {
					p := filepath.Join(dir, "test.zip")
					f, _ := os.Create(p)
					z := zip.NewWriter(f)
					entry, _ := z.Create("handler.js")
					entry.Write([]byte("test"))
					z.Close()
					f.Close()
					payload["codePath"] = p
				}
				if e.ID == "update-function-config" {
					payload["memory"] = 1000
				}
			}
			b, _ := json.Marshal(payload)
			request := httptest.NewRequest("POST", "/api/v1/"+e.FamilyID+"/endpoints/"+e.ID+"/invoke", bytes.NewReader(b))
			response := httptest.NewRecorder()
			Handler().ServeHTTP(response, request)
			if response.Code != 200 || seen == 0 {
				t.Fatalf("REST did not execute operation: status=%d calls=%d %s", response.Code, seen, response.Body.String())
			}
			var parsed any
			if json.Unmarshal(response.Body.Bytes(), &parsed) != nil {
				t.Fatal("REST output is not JSON")
			}
		})
	}
}
func TestServerlessRESTInvalidNestedInputNeverCallsUpstream(t *testing.T) {
	calls := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { calls++; io.WriteString(w, `{}`) }))
	defer srv.Close()
	payload := `{"baseUrl":"` + srv.URL + `","headers":{"productId":"p","service":"hosting","requestid":"test"},"body":{"siteId":"s","version":"v","materialInfos":[{"sceneId":"p","fileName":"x.zip","fileSha256":"s","fileSize":1,"contentType":1}]},"dryRun":false}`
	r := httptest.NewRequest("POST", "/api/v1/hosting/endpoints/upload-url/invoke", strings.NewReader(payload))
	w := httptest.NewRecorder()
	Handler().ServeHTTP(w, r)
	if w.Code != 400 || calls != 0 || !strings.Contains(w.Body.String(), "body.materialInfos[0].contentType") {
		t.Fatalf("invalid request sent upstream: %d %d %s", w.Code, calls, w.Body.String())
	}
}

func TestNativeRESTPreservesUpstreamErrorCode(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(403)
		w.Write([]byte(`{"code":"ROLE_DENIED","message":"permission denied"}`))
	}))
	defer s.Close()
	p := filepath.Join(t.TempDir(), "session.json")
	b, _ := json.Marshal(map[string]any{"uid": "test", "teamId": "test", "productId": "test", "accessToken": "offline", "expiresAt": time.Now().Add(time.Hour), "projectBaseUrl": s.URL})
	os.WriteFile(p, b, 0600)
	b, _ = json.Marshal(map[string]any{"name": "test-func", "sessionFile": p, "dryRun": false})
	r := httptest.NewRequest("POST", "/api/v1/functions/endpoints/invoke/invoke", bytes.NewReader(b))
	w := httptest.NewRecorder()
	Handler().ServeHTTP(w, r)
	var result domain.ErrorEnvelope
	json.Unmarshal(w.Body.Bytes(), &result)
	if result.Error.Code != "ROLE_DENIED" || result.Error.ResultCodes["code"] != "ROLE_DENIED" {
		t.Fatalf("upstream code lost: %s", w.Body.String())
	}
}
