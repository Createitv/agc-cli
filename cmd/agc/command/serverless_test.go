package command

import (
	"archive/zip"
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

func serverlessFixture(p domain.Parameter) any {
	if p.Name == "fileSha256" {
		return strings.Repeat("0", 64)
	}
	if len(p.Enum) > 0 {
		if p.Type == "integer" {
			var n json.Number = json.Number(p.Enum[0])
			return n
		}
		return p.Enum[0]
	}
	if p.EncodedJSON != nil {
		b, _ := json.Marshal(serverlessFixture(*p.EncodedJSON))
		return string(b)
	}
	switch p.Type {
	case "boolean":
		return false
	case "integer", "number":
		return 1
	case "object":
		m := map[string]any{}
		for _, c := range p.Properties {
			if c.Required {
				m[c.Name] = serverlessFixture(c)
			}
		}
		return m
	case "array":
		if p.ItemsType == "object" {
			m := map[string]any{}
			for _, c := range p.ItemsProperties {
				if c.Required {
					m[c.Name] = serverlessFixture(c)
				}
			}
			return []any{m}
		}
		if p.ItemsType == "integer" {
			return []any{1}
		}
		return []any{"test"}
	}
	return "test"
}

// Independent source-based transport and nested-contract tests complement these
// registry-wide CLI exchanges. Every operation sends to an isolated local server.
func testServerlessOperationHTTP(t *testing.T, e domain.Endpoint) {
	if e.FamilyID == "prefetch" && e.Direction == "local" {
		testPrefetchDeviceCLI(t, e, false)
		return
	}
	if agcapi.IsManagementEndpoint(e) {
		testReviewedManagementCLI(t, e)
		return
	}
	t.Run(e.FamilyID+"/"+e.ID, func(t *testing.T) {
		if e.Authentication == "agcli-user-oauth" {
			testNativeFunctionCLI(t, e)
			return
		}
		for _, scenario := range []struct {
			name   string
			status int
			body   string
			fail   bool
		}{{"success", 200, `{"ret":{"code":0},"marker":"received"}`, false}, {"business-error", 200, `{"ret":{"code":17,"msg":"rejected"}}`, true}, {"http-error", 429, `{"code":429,"message":"limited"}`, true}} {
			t.Run(scenario.name, func(t *testing.T) {
				seen := false
				srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					seen = true
					if r.Method != e.Method {
						t.Errorf("verb=%s want=%s", r.Method, e.Method)
					}
					if e.Direction == "inbound-callback" || e.Authentication == "body-client-credentials" {
						if r.Header.Get("Authorization") != "" {
							t.Error("unexpected developer auth")
						}
					}
					if e.Authentication == "signed-upload" && r.Header.Get("Authorization") != "AWS4-HMAC-SHA256 offline" {
						t.Error("lost signed header")
					}
					if e.ID == "upload-file" || e.ID == "upload-version" {
						b, _ := io.ReadAll(r.Body)
						if string(b) != "offline bytes" {
							t.Error("not raw bytes")
						}
					}
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(scenario.status)
					io.WriteString(w, scenario.body)
				}))
				defer srv.Close()
				args := []string{"--project", t.TempDir(), e.FamilyID, e.ID, "--invoke", "--dry-run=false", "--base-url", srv.URL}
				if e.Direction != "inbound-callback" && e.Authentication != "body-client-credentials" && e.Authentication != "signed-upload" {
					args = append(args, "--token", "offline")
				}
				body := map[string]any{}
				for _, p := range e.Parameters {
					if !p.Required {
						continue
					}
					if p.In == "file" {
						continue
					}
					v := serverlessFixture(p)
					value := fmt.Sprint(v)
					switch p.In {
					case "body":
						body[p.Name] = v
					case "path":
						if p.Name == "uploadUrl" || p.Name == "callbackUrl" {
							value = srv.URL + "/test"
						}
						args = append(args, "--param", p.Name+"="+value)
					case "header":
						args = append(args, "--header", p.Name+"="+value)
					case "query":
						args = append(args, "--query", p.Name+"="+value)
					}
				}
				if e.FamilyID == "cloud-auth" {
					switch e.ID {
					case "register", "request-code":
						body["email"] = "test@example.test"
					case "update-account":
						body["newEmail"] = "test@example.test"
					}
				}
				if e.FamilyID == "clouddb" && (e.ID == "update-data" || e.ID == "delete-data") {
					body["schemas"] = []any{map[string]any{"n": "StudentInfo", "fs": []any{map[string]any{"n": "id", "t": "TYPE_INT32"}, map[string]any{"n": "naturalbase_version", "t": "TYPE_LONG"}, map[string]any{"n": "naturalbase_deleted", "t": "TYPE_BOOLEAN"}}}}
					body["opData"] = []any{map[string]any{"t": 1, "os": []any{map[string]any{"fs": []any{map[string]any{"i": 1}, map[string]any{"l": nil}, map[string]any{"bl": e.ID == "delete-data"}}}}}}
				}
				if e.Authentication == "server-api-client" || e.Authentication == "client-api-token" {
					args = append(args, "--header", "client_id=offline")
				}
				if e.ID == "upload-file" || e.ID == "upload-version" {
					p := filepath.Join(t.TempDir(), "test.zip")
					os.WriteFile(p, []byte("offline bytes"), 0600)
					args = append(args, "--file", p)
				}
				if e.Authentication == "signed-upload" {
					sum := sha256.Sum256([]byte("offline bytes"))
					args = append(args, "--header", "Authorization=AWS4-HMAC-SHA256 offline", "--header", "x-amz-content-sha256="+hex.EncodeToString(sum[:]), "--header", "x-amz-client-request-id=offline", "--header", "x-amz-date=20261026T000000Z", "--header", "Host="+strings.TrimPrefix(srv.URL, "http://"), "--header", "Content-Type=application/zip")
				}
				if len(body) > 0 {
					p := filepath.Join(t.TempDir(), "body.json")
					b, _ := json.Marshal(body)
					os.WriteFile(p, b, 0600)
					args = append(args, "--body", p)
				}
				out, err := execute(args...)
				if !seen {
					t.Fatalf("CLI did not send request: err=%v output=%s", err, out)
				}
				wantFailure := scenario.fail
				// Downloaded JSON is user file content, including fields that
				// resemble an AGC business error. Only HTTP errors fail here.
				if e.FamilyID == "storage" && e.ID == "download-file" {
					wantFailure = scenario.status >= 400
				}
				if e.FamilyID == "prefetch" {
					wantFailure = scenario.status >= 400
				}
				if (err != nil) != wantFailure {
					t.Fatalf("err=%v want failure=%v", err, wantFailure)
				}
				if err == nil {
					var parsed any
					if json.Unmarshal([]byte(out), &parsed) != nil {
						t.Fatalf("invalid JSON stdout: %s", out)
					}
				}
			})
		}
	})
}
func testNativeFunctionCLI(t *testing.T, e domain.Endpoint) {
	calls := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		if r.Header.Get("oauth2Token") != "offline-token" {
			t.Error("missing selected session")
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"function":{"revisionId":"r1","handler":"handler.myHandler"},"marker":"received"}`))
	}))
	defer srv.Close()
	dir := t.TempDir()
	session := filepath.Join(dir, "session.json")
	b, _ := json.Marshal(map[string]any{"uid": "offline", "teamId": "offline", "productId": "offline", "accessToken": "offline-token", "expiresAt": time.Now().Add(time.Hour), "projectBaseUrl": srv.URL})
	os.WriteFile(session, b, 0600)
	args := []string{"--project", dir, "functions", e.ID, "--invoke", "--dry-run=false", "--name", "test-func", "--session-file", session, "--version-description", "test"}
	if e.ID == "create-function" || e.ID == "update-function-code" || e.ID == "update-and-publish-function" {
		p := filepath.Join(dir, "test.zip")
		f, _ := os.Create(p)
		z := zip.NewWriter(f)
		entry, _ := z.Create("handler.js")
		entry.Write([]byte("test"))
		z.Close()
		f.Close()
		args = append(args, "--code-path", p)
	}
	if e.ID == "update-function-config" {
		args = append(args, "--memory", "1000")
	}
	out, err := execute(args...)
	if err != nil || calls == 0 {
		t.Fatalf("native CLI failed: %v %s", err, out)
	}
}
func TestServerlessDryRunDoesNotLoadAuthorization(t *testing.T) {
	t.Setenv("AGC_CREDENTIALS_PATH", filepath.Join(t.TempDir(), "missing"))
	for _, args := range [][]string{{"functions", "invoke", "--invoke", "--name", "test-func", "--session-file", "/nonexistent"}, {"hosting", "list-sites", "--invoke", "--header", "productId=offline", "--header", "requestid=test", "--header", "service=hosting", "--field", "site={}"}} {
		out, err := execute(args...)
		if err != nil || !strings.Contains(out, `"dryRun":true`) {
			t.Fatalf("dry-run failed: %v %s", err, out)
		}
	}
}
func TestServerlessRejectsServiceAccountForClientOnlyAPI(t *testing.T) {
	dir := t.TempDir()
	creds := filepath.Join(dir, "creds.json")
	agcapi.SaveCredential(creds, agcapi.Credential{Name: "offline", Mode: "service-account", ServiceAccountFile: "/must-not-read"})
	_, err := execute("--project", dir, "--profile", "offline", "storage", "list-files", "--invoke", "--dry-run=false", "--credentials-path", creds, "--param", "bucket_name=test", "--header", "productId=test")
	if err == nil || !strings.Contains(err.Error(), "API client") {
		t.Fatalf("must reject wrong credential before reading key: %v", err)
	}
}

func testServerlessOperationDryRun(t *testing.T, e domain.Endpoint) {
	if e.FamilyID == "prefetch" && e.Direction == "local" {
		testPrefetchDeviceCLI(t, e, true)
		return
	}
	if agcapi.IsManagementEndpoint(e) {
		testReviewedManagementDryRun(t, e)
		return
	}
	t.Run(e.FamilyID+"/"+e.ID, func(t *testing.T) {
		// Capture valid CLI arguments through the real transport fixture, replacing
		// only the execution switch is tested separately below for native workflows.
		if e.Authentication == "agcli-user-oauth" {
			dir := t.TempDir()
			args := []string{"--project", dir, "functions", e.ID, "--invoke", "--name", "test-func", "--session-file", "/must-not-read", "--version-description", "test"}
			if e.ID == "create-function" || e.ID == "update-function-code" {
				p := filepath.Join(dir, "test.zip")
				f, _ := os.Create(p)
				z := zip.NewWriter(f)
				z.Create("handler.js")
				z.Close()
				f.Close()
				args = append(args, "--code-path", p)
			}
			if e.ID == "update-function-config" {
				args = append(args, "--memory", "1000")
			}
			out, err := execute(args...)
			if err != nil || !strings.Contains(out, `"dryRun":true`) {
				t.Fatalf("native dry-run failed: %v %s", err, out)
			}
			return
		}
		args := []string{"--project", t.TempDir(), e.FamilyID, e.ID, "--invoke", "--base-url", "https://example.test"}
		body := map[string]any{}
		for _, p := range e.Parameters {
			if !p.Required {
				continue
			}
			v := serverlessFixture(p)
			switch p.In {
			case "body":
				body[p.Name] = v
			case "header":
				args = append(args, "--header", p.Name+"="+fmt.Sprint(v))
			case "query":
				args = append(args, "--query", p.Name+"="+fmt.Sprint(v))
			case "path":
				value := fmt.Sprint(v)
				if p.Name == "callbackUrl" || p.Name == "uploadUrl" {
					value = "https://example.test/hook"
				}
				args = append(args, "--param", p.Name+"="+value)
			case "file":
				file := filepath.Join(t.TempDir(), "offline.zip")
				os.WriteFile(file, []byte("offline"), 0600)
				args = append(args, "--file", file)
			}
		}
		if e.FamilyID == "cloud-auth" {
			switch e.ID {
			case "register", "request-code":
				body["email"] = "test@example.test"
			case "update-account":
				body["newEmail"] = "test@example.test"
			}
		}
		if e.FamilyID == "clouddb" && (e.ID == "update-data" || e.ID == "delete-data") {
			body["schemas"] = []any{}
			body["opData"] = []any{}
		}
		if e.Authentication == "signed-upload" {
			sum := sha256.Sum256([]byte("offline"))
			args = append(args, "--header", "Authorization=AWS4-HMAC-SHA256 offline", "--header", "x-amz-content-sha256="+hex.EncodeToString(sum[:]), "--header", "Host=example.test", "--header", "Content-Type=application/zip")
		}
		if len(body) > 0 {
			b, _ := json.Marshal(body)
			file := filepath.Join(t.TempDir(), "body.json")
			os.WriteFile(file, b, 0600)
			args = append(args, "--body", file)
		}
		out, err := execute(args...)
		if err != nil || !strings.Contains(out, `"dryRun":true`) {
			t.Fatalf("dry run failed: %v %s", err, out)
		}
	})
}

func TestStorageEmptyDownloadWritesOutputFile(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Length", "0")
		w.WriteHeader(200)
	}))
	defer s.Close()
	out := filepath.Join(t.TempDir(), "empty.bin")
	_, err := execute("--project", t.TempDir(), "storage", "download-file", "--invoke", "--dry-run=false", "--base-url", s.URL, "--token", "offline", "--param", "bucket_name=b", "--param", "object_name=empty.bin", "--header", "productId=p", "--header", "client_id=c", "--out", out)
	if err != nil {
		t.Fatal(err)
	}
	info, err := os.Stat(out)
	if err != nil || info.Size() != 0 {
		t.Fatalf("empty download not written: %v", err)
	}
}

func TestNativeFunctionCLIErrorsAreMachineReadable(t *testing.T) {
	out, err := execute("functions", "invoke", "--invoke", "--name", "BAD")
	if err == nil {
		t.Fatal("invalid name accepted")
	}
	var e domain.ErrorEnvelope
	if json.Unmarshal([]byte(out), &e) != nil || e.Error.Endpoint != "functions/invoke" || len(e.Error.Violations) == 0 || e.Error.DocumentationCommand == "" {
		t.Fatalf("native feedback missing: %s", out)
	}
}
func TestSignedHostingRequiresReturnedAuthorizationBeforeNetwork(t *testing.T) {
	e, _ := domain.EndpointByID("hosting", "upload-version")
	err := agcapi.ValidateEndpointRequest(e, nil, map[string]string{"x-amz-content-sha256": "test", "x-amz-client-request-id": "test", "x-amz-date": "test", "Host": "example.test", "Content-Type": "application/zip"}, nil)
	if err == nil {
		t.Fatal("missing signed Authorization accepted")
	}
}
