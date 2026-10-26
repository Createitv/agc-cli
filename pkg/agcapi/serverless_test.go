package agcapi

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Createitv/agc-cli/pkg/domain"
)

func slEndpoint(t *testing.T, f, id string) domain.Endpoint {
	t.Helper()
	e, ok := domain.EndpointByID(f, id)
	if !ok {
		t.Fatalf("missing %s/%s", f, id)
	}
	return e
}
func TestServerlessBinaryUploads(t *testing.T) {
	data := []byte{0, 255, 7, 10, 13}
	file := filepath.Join(t.TempDir(), "file.zip")
	os.WriteFile(file, data, 0600)
	for _, f := range []string{"storage", "hosting"} {
		t.Run(f, func(t *testing.T) {
			id := "upload-file"
			if f == "hosting" {
				id = "upload-version"
			}
			e := slEndpoint(t, f, id)
			sum := sha256.Sum256(data)
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				b, _ := io.ReadAll(r.Body)
				if string(b) != string(data) || r.ContentLength != int64(len(data)) {
					t.Errorf("binary body corrupted %v", b)
				}
				if f == "hosting" && r.Header.Get("Authorization") != "AWS4-HMAC-SHA256 signed" {
					t.Error("signed auth replaced")
				}
				if f == "storage" && r.Header.Get("Authorization") != "Bearer offline" {
					t.Error("storage authorization missing")
				}
				w.Write([]byte(`{"uploadStatus":"finalize"}`))
			}))
			defer srv.Close()
			params := map[string]string{"bucket_name": "test", "object_name": "nested/a b.zip", "uploadUrl": srv.URL + "/signed"}
			heads := map[string]string{"productId": "test", "client_id": "test"}
			if f == "hosting" {
				heads["Authorization"] = "AWS4-HMAC-SHA256 signed"
				heads["x-amz-content-sha256"] = hex.EncodeToString(sum[:])
				heads["x-amz-client-request-id"] = "offline"
				heads["x-amz-date"] = "20261026T000000Z"
				heads["Host"] = strings.TrimPrefix(srv.URL, "http://")
				heads["Content-Type"] = "application/zip"
			}
			_, err := InvokeEndpoint(context.Background(), srv.Client(), InvokeRequest{Endpoint: e, BaseURL: srv.URL, Params: params, Headers: heads, FilePath: file, AccessToken: "offline"})
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}
func TestServerlessCallbackNeverForwardsAGCToken(t *testing.T) {
	e := slEndpoint(t, "templates", "player-id")
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "" {
			t.Error("AGC token leaked to developer callback")
		}
		w.Write([]byte(`{"ret":{"code":0},"data":{"playerId":"test"}}`))
	}))
	defer s.Close()
	_, err := InvokeEndpoint(context.Background(), s.Client(), InvokeRequest{Endpoint: e, BaseURL: s.URL, Params: map[string]string{"callbackUrl": s.URL}, Body: []byte(`{"appId":"test","accessToken":"user-token"}`), AccessToken: "must-not-forward"})
	if err != nil {
		t.Fatal(err)
	}
}
func TestServerlessNestedValidation(t *testing.T) {
	tests := []struct{ family, id, body, path string }{
		{"hosting", "upload-url", `{"siteId":"s","version":"v","materialInfos":[{"sceneId":"p","fileName":"x.zip","fileSha256":"abc","fileSize":1,"contentType":1}]}`, "body.materialInfos[0].contentType"},
		{"cloud-auth", "import-users", `{"users":[{"importUid":"u","createTime":"0","providers":[{}]}]}`, "body.users[0].providers[0].provider"},
		{"clouddb", "query-data", `{"msgInfo":{"type":5,"opStore":{"storeName":"z"}},"clientInfo":{"appVer":1},"queryReqMsg":{"queryType":0,"queryTable":"Student","queryCond":"{\"queryConditions\":[{}]}"}}`, "body.queryReqMsg.queryCond.queryConditions[0].conditionType"},
	}
	for _, tt := range tests {
		t.Run(tt.family+tt.id, func(t *testing.T) {
			e := slEndpoint(t, tt.family, tt.id)
			err := ValidateEndpointRequest(e, map[string]string{"_v": "4"}, map[string]string{"productId": "p", "requestid": "test", "service": "hosting"}, []byte(tt.body))
			v, ok := err.(*RequestValidationError)
			if !ok {
				t.Fatalf("expected nested violation, got %v", err)
			}
			found := false
			for _, i := range v.Violations {
				found = found || i.Field == tt.path
			}
			if !found {
				t.Fatalf("missing indexed field %s: %+v", tt.path, v.Violations)
			}
		})
	}
}
func TestServerlessCloudDBBusinessFailure(t *testing.T) {
	data := []byte(`{"resInfo":{"resCode":1001,"resInfo":"permission denied"}}`)
	if !businessFailure(data) {
		t.Fatal("CloudDB error treated as success")
	}
	e := ParseAGCError(200, data)
	if e.Code != "1001" {
		t.Fatalf("lost code: %+v", e)
	}
}
func TestServerlessPathsAreEncodedOnce(t *testing.T) {
	u, err := BuildEndpointURL("https://example.test/v0", "/{bucket_name}/{object_name}", map[string]string{"bucket_name": "test", "object_name": "folder/a b?#.txt"}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if u != "https://example.test/v0/test/folder/a%20b%3F%23.txt" {
		t.Fatalf("unescaped path: %s", u)
	}
}
func TestServerlessEncodedQueryCondRetainsString(t *testing.T) {
	e := slEndpoint(t, "clouddb", "query-data")
	wire, err := MarshalEndpointFields(e, map[string]string{"msgInfo": `{"type":5,"opStore":{"storeName":"z"}}`, "clientInfo": `{"appVer":1}`, "queryReqMsg": `{"queryType":0,"queryTable":"Student","queryCond":"{\"queryConditions\":[]}"}`})
	if err != nil {
		t.Fatal(err)
	}
	var obj map[string]map[string]json.RawMessage
	json.Unmarshal(wire, &obj)
	if string(obj["queryReqMsg"]["queryCond"])[0] != '"' {
		t.Fatal("encoded JSON became object")
	}
}

func TestServerlessDocumentedConditionalRules(t *testing.T) {
	tests := []struct{ family, id, body, field string }{
		{"cloud-auth", "register", `{}`, "body.phone|body.email"},
		{"cloud-auth", "register", `{"phone":"+86-12345678"}`, "body.verifyCode"},
		{"cloud-auth", "request-code", `{"email":"test@example.test","sendInterval":29}`, "body.sendInterval"},
		{"cloud-auth", "import-users", `{"users":[{"importUid":"u","createTime":"0","providers":[{"provider":"12","providerUid":"u"}]}]}`, "body.users[0].email"},
		{"clouddb", "query-data", `{"msgInfo":{"type":5,"opStore":{"storeName":"z"}},"clientInfo":{"appVer":1},"queryReqMsg":{"queryType":2,"queryTable":"Student","queryCond":"{\"queryConditions\":[{\"conditionType\":\"Average\",\"fieldName\":\"age\"}]}"}}`, "body.queryReqMsg.queryCond"},
	}
	for _, tt := range tests {
		t.Run(tt.id+tt.field, func(t *testing.T) {
			e := slEndpoint(t, tt.family, tt.id)
			err := ValidateEndpointRequest(e, map[string]string{"productId": "p", "_v": "4"}, map[string]string{"productId": "p"}, []byte(tt.body))
			v, ok := err.(*RequestValidationError)
			if !ok {
				t.Fatalf("missing documented rule: %v", err)
			}
			found := false
			for _, i := range v.Violations {
				found = found || i.Field == tt.field
			}
			if !found {
				t.Fatalf("missing %s: %+v", tt.field, v.Violations)
			}
		})
	}
}
func TestServerlessMissingFileFailsLocally(t *testing.T) {
	e := slEndpoint(t, "storage", "upload-file")
	if ValidateEndpointFileInput(e, "/nonexistent") == nil {
		t.Fatal("missing local file accepted")
	}
}

func TestStorageDownloadPreservesBinaryInJSON(t *testing.T) {
	data := []byte{0, 255, 128, 10}
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Write(data)
	}))
	defer s.Close()
	result, err := InvokeEndpoint(context.Background(), s.Client(), InvokeRequest{Endpoint: slEndpoint(t, "storage", "download-file"), BaseURL: s.URL, Params: map[string]string{"bucket_name": "b", "object_name": "x.bin"}, AccessToken: "offline", Headers: map[string]string{"productId": "p", "client_id": "c"}})
	if err != nil {
		t.Fatal(err)
	}
	raw, _ := json.Marshal(result)
	var obj map[string]any
	json.Unmarshal(raw, &obj)
	if obj["bodyBase64"] != "AP+ACg==" {
		t.Fatalf("binary response not lossless in JSON: %s", raw)
	}
}

func TestServerlessRegionalTokenRouting(t *testing.T) {
	for base, want := range map[string]string{"https://ops-server-drcn.agcstorage.link/v0": "https://connect-api.cloud.huawei.com", "https://ops-server-dre.agcstorage.link/v0": "https://connect-api-dre.cloud.huawei.com", "https://ops-server-dra.agcstorage.link/v0": "https://connect-api-dra.cloud.huawei.com", "https://ops-server-drru.agcstorage.link/v0": "https://connect-api-drru.cloud.huawei.com"} {
		got, err := ServerlessTokenBaseURL(base)
		if err != nil || got != want {
			t.Fatalf("wrong regional token host %s %v", got, err)
		}
	}
	for _, base := range []string{"http://ops-server-dre.agcstorage.link", "https://unreviewed.example.test", "bad"} {
		if _, err := ServerlessTokenBaseURL(base); err == nil {
			t.Fatalf("unreviewed token destination accepted %s", base)
		}
	}
}
func TestClientTokenUsesClientGatewayContract(t *testing.T) {
	seen := false
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		seen = true
		if r.URL.Path != "/agc/apigw/oauth2/v1/token" || r.Header.Get("Authorization") != "" {
			t.Error("wrong token transport")
		}
		var b map[string]any
		json.NewDecoder(r.Body).Decode(&b)
		if b["useJwt"] != float64(1) || b["grant_type"] != "client_credentials" || b["client_secret"] != "offline-secret" {
			t.Error("wrong client token body")
		}
		w.Write([]byte(`{"access_token":"offline-token"}`))
	}))
	defer s.Close()
	token, err := ServerlessClientToken(context.Background(), s.Client(), s.URL, Credential{Mode: "api-client", ClientID: "offline", ClientKey: "offline-secret"})
	if err != nil || token != "offline-token" || !seen {
		t.Fatalf("client authorization failed %v", err)
	}
}
func TestClientTokenFailureDoesNotEchoCredentialResponse(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(401)
		w.Write([]byte(`{"secret":"must-not-echo"}`))
	}))
	defer s.Close()
	_, err := ServerlessClientToken(context.Background(), s.Client(), s.URL, Credential{Mode: "api-client", ClientID: "offline", ClientKey: "offline-secret"})
	if err == nil || strings.Contains(err.Error(), "must-not-echo") {
		t.Fatalf("unsafe token error %v", err)
	}
}

func TestStorageJSONFilesAreNotBusinessEnvelopes(t *testing.T) {
	data := []byte("{\n  \"code\": 17, \"file\": true\n}\n")
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	}))
	defer s.Close()
	r, err := InvokeEndpoint(context.Background(), s.Client(), InvokeRequest{Endpoint: slEndpoint(t, "storage", "download-file"), BaseURL: s.URL, Params: map[string]string{"bucket_name": "b", "object_name": "data.json"}, AccessToken: "offline", Headers: map[string]string{"productId": "p", "client_id": "c"}})
	if err != nil {
		t.Fatalf("file content interpreted as business failure: %v", err)
	}
	encoded, _ := json.Marshal(r)
	var obj map[string]any
	json.Unmarshal(encoded, &obj)
	b64, _ := obj["bodyBase64"].(string)
	if b64 == "" {
		t.Fatal("JSON file needs lossless Base64 representation too")
	}
}

func TestStorageResumableHeadersAreCheckedLocally(t *testing.T) {
	e := slEndpoint(t, "storage", "upload-file")
	err := ValidateEndpointRequest(e, nil, map[string]string{"productId": "p", "client_id": "c", "X-Agc-File-Offset": "1"}, nil)
	v, ok := err.(*RequestValidationError)
	if !ok {
		t.Fatalf("resumable prerequisites accepted: %v", err)
	}
	seen := map[string]bool{}
	for _, issue := range v.Violations {
		seen[issue.Field] = true
	}
	if !seen["header.X-Agc-File-Size"] || !seen["header.X-Agc-Sha256"] {
		t.Fatalf("missing size/digest corrections: %+v", v.Violations)
	}
}

func TestSignedHostingDigestMismatchFailsBeforeNetwork(t *testing.T) {
	path := filepath.Join(t.TempDir(), "test.zip")
	os.WriteFile(path, []byte("ZIP bytes"), 0600)
	e := slEndpoint(t, "hosting", "upload-version")
	_, err := InvokeEndpoint(context.Background(), nil, InvokeRequest{Endpoint: e, Params: map[string]string{"uploadUrl": "https://example.test/signed"}, FilePath: path, DryRun: true, Headers: map[string]string{"Authorization": "AWS4-HMAC-SHA256 offline", "x-amz-content-sha256": strings.Repeat("0", 64), "x-amz-client-request-id": "offline", "x-amz-date": "20261026T000000Z", "Host": "example.test", "Content-Type": "application/zip"}})
	if err == nil {
		t.Fatal("signed SHA256 mismatch accepted")
	}
}

type serverlessRoundTrip func(*http.Request) (*http.Response, error)

func (f serverlessRoundTrip) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }
func TestTemplateCallbackRequiresHTTPSOutsideLoopback(t *testing.T) {
	called := false
	client := &http.Client{Transport: serverlessRoundTrip(func(r *http.Request) (*http.Response, error) {
		called = true
		return &http.Response{StatusCode: 200, Header: http.Header{}, Body: io.NopCloser(strings.NewReader(`{"ret":{"code":0}}`))}, nil
	})}
	_, err := InvokeEndpoint(context.Background(), client, InvokeRequest{Endpoint: slEndpoint(t, "templates", "user-id"), Params: map[string]string{"callbackUrl": "http://developer.example.test/auth"}, Body: []byte(`{"accessToken":"offline-user-credential"}`)})
	if err == nil || called {
		t.Fatalf("non-HTTPS callback accepted: called=%v err=%v", called, err)
	}
}
