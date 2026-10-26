package cloudfunctions

import (
	"archive/zip"
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/Createitv/agc-cli/pkg/domain"
)

func testZIP(t *testing.T) string {
	t.Helper()
	p := filepath.Join(t.TempDir(), "func.zip")
	f, _ := os.Create(p)
	w := zip.NewWriter(f)
	e, _ := w.Create("handler.js")
	e.Write([]byte("exports.myHandler = () => 1"))
	w.Close()
	f.Close()
	return p
}
func offlineSession() Session {
	return Session{UID: "offline-user", TeamID: "offline-team", ProductID: "offline-project", AccessToken: "offline-oauth", ExpiresAt: time.Now().Add(time.Hour)}
}
func TestNativeFunctionEveryWorkflow(t *testing.T) {
	for _, op := range domain.FunctionOperations {
		t.Run(op, func(t *testing.T) {
			calls := []string{}
			file := testZIP(t)
			fileData, _ := os.ReadFile(file)
			sum := sha256.Sum256(fileData)
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				calls = append(calls, r.Method+" "+r.URL.Path)
				if r.Header.Get("oauth2Token") != "offline-oauth" || r.Header.Get("Authorization") != "" || r.Header.Get("productId") != "offline-project" {
					t.Error("wrong OAuth headers")
				}
				w.Header().Set("Content-Type", "application/json")
				if r.Method == "PUT" || r.Method == "POST" && strings.HasSuffix(r.URL.Path, "/functions") {
					if err := r.ParseMultipartForm(2 << 20); err != nil {
						t.Fatal(err)
					}
					if r.Method == "PUT" {
						if r.FormValue("revisionId") != "r1" || r.FormValue("handler") != "handler.original" || r.FormValue("environment") != `{"KEEP":"private"}` || r.FormValue("layers") != `["urn:keep:1"]` {
							t.Error("lost prior revision/config")
						}
					}
					if op == "create-function" || op == "update-function-code" || op == "update-and-publish-function" {
						f, _, err := r.FormFile("file")
						if err != nil {
							t.Fatal(err)
						}
						b, _ := io.ReadAll(f)
						if !bytes.Equal(b, fileData) || r.Header.Get("x-file-digest-sha256") != hex.EncodeToString(sum[:]) {
							t.Error("ZIP/digest mismatch")
						}
					}
					w.Write([]byte(`{"function":{"revisionId":"r2","environment":{"SECRET":"private"}}}`))
					return
				}
				if strings.HasSuffix(r.URL.Path, "/invocations") {
					b, _ := io.ReadAll(r.Body)
					if string(b) != `[9007199254740993,"hello"]` {
						t.Errorf("payload corrupted %s", b)
					}
					w.Header().Set("x-function-execution-status", "success")
					w.Write([]byte(`{"result":1}`))
					return
				}
				if r.Method == "POST" && strings.HasSuffix(r.URL.Path, "/versions") {
					var b map[string]string
					json.NewDecoder(r.Body).Decode(&b)
					want := "r1"
					if op == "update-and-publish-function" {
						want = "r2"
					}
					if b["revisionId"] != want {
						t.Errorf("publish revision=%s", b["revisionId"])
					}
					w.Write([]byte(`{"function":{"version":"1"}}`))
					return
				}
				w.Write([]byte(`{"function":{"revisionId":"r1","handler":"handler.original","environment":{"KEEP":"private"},"layers":["urn:keep:1"]}}`))
			}))
			defer srv.Close()
			c := Client{HTTP: srv.Client(), Session: offlineSession()}
			req := domain.FunctionRequest{Name: "test-func", BaseURL: srv.URL, ProjectID: "offline-project", VersionDescription: "test", Payload: json.RawMessage(`[9007199254740993,"hello"]`)}
			if op == "create-function" || op == "update-function-code" || op == "update-and-publish-function" {
				req.CodePath = file
			}
			m := 1000
			if op == "update-function-config" {
				req.Memory = &m
			}
			resp, err := c.Execute(context.Background(), op, req)
			if err != nil {
				t.Fatal(err)
			}
			b, _ := json.Marshal(resp)
			if strings.Contains(string(b), "private") || strings.Contains(string(b), "offline-oauth") {
				t.Error("response leaked authorization/environment")
			}
			if len(calls) == 0 {
				t.Fatal("no actual HTTP exchange")
			}
			if (op == "publish-function" || op == "update-function-code" || op == "update-function-config") && len(calls) != 2 {
				t.Fatalf("expected revision read before write: %v", calls)
			}
			if op == "update-and-publish-function" && len(calls) != 3 {
				t.Fatalf("expected read/update/publish: %v", calls)
			}
		})
	}
}
func TestUpdateFailureNeverPublishes(t *testing.T) {
	calls := 0
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		if r.Method == "GET" {
			w.Write([]byte(`{"function":{"revisionId":"r1"}}`))
			return
		}
		w.WriteHeader(409)
		w.Write([]byte(`{"code":"CONFLICT","message":"revision conflict"}`))
	}))
	defer s.Close()
	n := 1000
	_, err := (&Client{HTTP: s.Client(), Session: offlineSession()}).Execute(context.Background(), "update-and-publish-function", domain.FunctionRequest{Name: "test-func", BaseURL: s.URL, VersionDescription: "test", Memory: &n})
	if err == nil || calls != 2 {
		t.Fatalf("conflict must stop, calls=%d err=%v", calls, err)
	}
}
func TestNativeFunctionDryRunNeverLoadsSession(t *testing.T) {
	_, err := (&Client{}).Execute(context.Background(), "invoke", domain.FunctionRequest{Name: "test-func", SessionFile: "/must/not/read", DryRun: true})
	if err != nil {
		t.Fatal(err)
	}
}
func TestFunctionInvalidBeforeAuthentication(t *testing.T) {
	for _, r := range []domain.FunctionRequest{{Name: "BAD"}, {Name: "test-func", Payload: json.RawMessage(`{`)}, {Name: "test-func", Version: "bad"}, {Name: "test-func", ProjectID: "wrong"}} {
		_, err := (&Client{Session: offlineSession()}).Execute(context.Background(), "invoke", r)
		if err == nil {
			t.Fatal("invalid input accepted")
		}
	}
}

func TestPublishFailureReportsAppliedUpdate(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			w.Write([]byte(`{"function":{"revisionId":"r1","handler":"handler.myHandler"}}`))
		case "PUT":
			w.Write([]byte(`{"function":{"revisionId":"r2"}}`))
		case "POST":
			w.WriteHeader(409)
			w.Write([]byte(`{"code":"CONFLICT"}`))
		}
	}))
	defer s.Close()
	m := 1000
	result, err := (&Client{HTTP: s.Client(), Session: offlineSession()}).Execute(context.Background(), "update-and-publish-function", domain.FunctionRequest{Name: "test-func", BaseURL: s.URL, VersionDescription: "test", Memory: &m})
	if err == nil || !result.Updated || result.Published || !strings.Contains(err.Error(), "update succeeded") {
		t.Fatalf("partial state lost: %+v %v", result, err)
	}
}
func TestLayerReferencesResolveThroughDocumentedResultWrapper(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "/layers/") {
			w.Write([]byte(`{"code":0,"result":{"layerVersionUrn":"urn:resolved:1","deletedFlag":false}}`))
			return
		}
		if r.Method == "GET" {
			w.Write([]byte(`{"function":{"revisionId":"r1","handler":"handler.myHandler"}}`))
			return
		}
		r.ParseMultipartForm(1 << 20)
		if r.FormValue("layers") != `["urn:resolved:1"]` {
			t.Error("resolved layer lost")
		}
		w.Write([]byte(`{"function":{"revisionId":"r2"}}`))
	}))
	defer s.Close()
	layers := []string{"test-layer:1"}
	_, err := (&Client{HTTP: s.Client(), Session: offlineSession()}).Execute(context.Background(), "update-function-config", domain.FunctionRequest{Name: "test-func", BaseURL: s.URL, Layers: &layers})
	if err != nil {
		t.Fatal(err)
	}
}
func TestInvokeFailureHeaderIsIndependentOfHTTP200(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("x-function-execution-status", "failed")
		w.Write([]byte(`{"message":"failed"}`))
	}))
	defer s.Close()
	r, err := (&Client{HTTP: s.Client(), Session: offlineSession()}).Execute(context.Background(), "invoke", domain.FunctionRequest{Name: "test-func", BaseURL: s.URL})
	if err == nil || r.StatusCode != 200 || r.ExecutionStatus != "failed" {
		t.Fatalf("execution failed unnoticed: %+v %v", r, err)
	}
}
func TestNativeOAuthRedirectCannotReachOtherOrigin(t *testing.T) {
	leaked := false
	target := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { leaked = true }))
	defer target.Close()
	source := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { http.Redirect(w, r, target.URL, 307) }))
	defer source.Close()
	_, err := (&Client{Session: offlineSession()}).Execute(context.Background(), "invoke", domain.FunctionRequest{Name: "test-func", BaseURL: source.URL})
	if err == nil || leaked {
		t.Fatal("cross-origin OAuth redirect allowed")
	}
}

func TestSourceDocumentedFunctionInputBounds(t *testing.T) {
	n := func(v int) *int { return &v }
	text := func(v string) *string { return &v }
	env := map[string]string{"INVALID KEY": "value"}
	layers := []string{"bad"}
	tooMany := []string{"a:1", "b:1", "c:1", "d:1", "e:1", "f:1"}
	tests := []struct {
		name, op string
		r        domain.FunctionRequest
	}{
		{"memory", "invoke", domain.FunctionRequest{Name: "test-func", Memory: n(512)}}, {"timeout-low", "invoke", domain.FunctionRequest{Name: "test-func", Timeout: n(0)}}, {"timeout-high", "invoke", domain.FunctionRequest{Name: "test-func", Timeout: n(1801)}}, {"concurrency", "invoke", domain.FunctionRequest{Name: "test-func", InstanceConcurrency: n(1001)}}, {"runtime", "create-function", domain.FunctionRequest{Name: "test-func", Runtime: "unknown"}}, {"description", "invoke", domain.FunctionRequest{Name: "test-func", Description: strings.Repeat("中", 33)}}, {"alias-version", "invoke", domain.FunctionRequest{Name: "test-func", Version: "1", AliasName: "live"}}, {"alias", "invoke", domain.FunctionRequest{Name: "test-func", AliasName: "BAD"}}, {"version", "invoke", domain.FunctionRequest{Name: "test-func", Version: "1001"}}, {"handler", "invoke", domain.FunctionRequest{Name: "test-func", Handler: text("")}}, {"environment", "invoke", domain.FunctionRequest{Name: "test-func", Environment: &env}}, {"layer-format", "invoke", domain.FunctionRequest{Name: "test-func", Layers: &layers}}, {"layer-limit", "invoke", domain.FunctionRequest{Name: "test-func", Layers: &tooMany}}, {"page-size", "list-functions", domain.FunctionRequest{PageSize: 1001}}, {"page-index", "list-functions", domain.FunctionRequest{PageIndex: -1}}, {"filter", "list-functions", domain.FunctionRequest{FuzzyName: "BAD"}}, {"missing-code", "create-function", domain.FunctionRequest{Name: "test-func"}}, {"missing-description", "publish-function", domain.FunctionRequest{Name: "test-func"}}, {"no-update", "update-function-config", domain.FunctionRequest{Name: "test-func"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.r.DryRun = true
			if _, err := (&Client{}).Execute(context.Background(), tt.op, tt.r); err == nil {
				t.Fatal("invalid documented input accepted")
			}
		})
	}
}
func TestExplicitSessionFileExpiryAndProjectIsolation(t *testing.T) {
	for _, kind := range []string{"expired", "project-mismatch", "invalid-json", "missing"} {
		t.Run(kind, func(t *testing.T) {
			p := filepath.Join(t.TempDir(), "session.json")
			s := offlineSession()
			if kind == "expired" {
				s.ExpiresAt = time.Now().Add(-time.Hour)
			}
			b, _ := json.Marshal(s)
			if kind == "invalid-json" {
				b = []byte(`{"accessToken":`)
			}
			if kind != "missing" {
				os.WriteFile(p, b, 0600)
			}
			r := domain.FunctionRequest{Name: "test-func", SessionFile: p}
			if kind == "project-mismatch" {
				r.ProjectID = "wrong"
			}
			_, err := (&Client{}).Execute(context.Background(), "invoke", r)
			if err == nil || strings.Contains(err.Error(), "offline-oauth") {
				t.Fatalf("session boundary failed %v", err)
			}
		})
	}
}

func TestConfigUpdateAllowsOnlyVersionDescription(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			w.Write([]byte(`{"function":{"revisionId":"r1","handler":"handler.myHandler"}}`))
			return
		}
		r.ParseMultipartForm(1 << 20)
		if r.FormValue("versionDesc") != "new description" {
			t.Error("description omitted")
		}
		w.Write([]byte(`{"function":{"revisionId":"r2"}}`))
	}))
	defer s.Close()
	_, err := (&Client{HTTP: s.Client(), Session: offlineSession()}).Execute(context.Background(), "update-function-config", domain.FunctionRequest{Name: "test-func", BaseURL: s.URL, VersionDescription: "new description"})
	if err != nil {
		t.Fatal(err)
	}
}

func TestInvocationPreservesArbitraryApplicationJSON(t *testing.T) {
	want := `{"code":17,"uid":"application-user","environment":{"label":"business-result"}}`
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.Write([]byte(want)) }))
	defer s.Close()
	result, err := (&Client{Session: offlineSession()}).Execute(context.Background(), "invoke", domain.FunctionRequest{Name: "test-func", BaseURL: s.URL})
	if err != nil || string(result.Body) != want {
		t.Fatalf("application result was misinterpreted/redacted: %s %v", result.Body, err)
	}
}
func TestMemoryCPUAllocationMatchesVendorParser(t *testing.T) {
	for memory, cpu := range map[int]string{500: "500", 1000: "500", 2000: "1000", 4000: "2000"} {
		t.Run(cpu+"-"+strconv.Itoa(memory), func(t *testing.T) {
			s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method == "GET" {
					w.Write([]byte(`{"function":{"revisionId":"r1"}}`))
					return
				}
				r.ParseMultipartForm(1 << 20)
				if r.FormValue("cpu") != cpu {
					t.Errorf("memory %d requires cpu %s, got %q", memory, cpu, r.FormValue("cpu"))
				}
				w.Write([]byte(`{"function":{"revisionId":"r2"}}`))
			}))
			defer s.Close()
			_, err := (&Client{Session: offlineSession()}).Execute(context.Background(), "update-function-config", domain.FunctionRequest{Name: "test-func", BaseURL: s.URL, Memory: &memory})
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestOversizedInvocationFailsInsteadOfTruncating(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { io.WriteString(w, strings.Repeat("x", (16<<20)+1)) }))
	defer s.Close()
	_, err := (&Client{Session: offlineSession()}).Execute(context.Background(), "invoke", domain.FunctionRequest{Name: "test-func", BaseURL: s.URL})
	if err == nil || !strings.Contains(err.Error(), "16 MiB") {
		t.Fatalf("oversize success/truncation: %v", err)
	}
}
