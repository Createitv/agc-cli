package command

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Createitv/agc-cli/pkg/agcapi"
	"github.com/Createitv/agc-cli/pkg/domain"
)

// 逐接口检查命令是否正确发送请求头、查询参数和 JSON 请求体。
// 同时验证签名上传不泄露 AGC 令牌，HTTP 错误与业务错误不会被当作成功。
// 所有请求都发送到本地测试服务器；通过不代表华为真实业务验收通过。
func TestEveryEndpointHTTPTransport(t *testing.T) {
	for _, endpoint := range domain.AllEndpoints() {
		if domain.IsServerlessFamily(endpoint.FamilyID) {
			testServerlessOperationHTTP(t, endpoint)
			continue
		}
		if endpoint.Direction == "local" {
			continue
		} // 本地构建由独立执行器测试，不应发送 HTTP。
		endpoint := endpoint
		t.Run(endpoint.FamilyID+"/"+endpoint.ID, func(t *testing.T) {
			t.Parallel()
			for _, scenario := range []struct {
				name      string
				status    int
				response  string
				wantError bool
			}{
				{"success", 200, `{"rtnCode":0,"marker":"received"}`, false},
				{"business-error", 200, `{"rtnCode":17,"rtnDesc":"rejected"}`, true},
				{"rate-limit", 429, `{"code":429,"message":"rate limited"}`, true},
			} {
				t.Run(scenario.name, func(t *testing.T) {
					resourceUpload := endpoint.FamilyID == "pms" && endpoint.ID == "resourcemanage-harmonyosnext"
					signed := endpoint.Path == "{uploadUrl}"
					seen := make(chan struct{}, 1)
					server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						seen <- struct{}{}
						if r.Method != endpoint.Method {
							t.Errorf("method = %s, want %s", r.Method, endpoint.Method)
						}
						if r.URL.Query().Get("probe") != "a b&c" {
							t.Error("query value was dropped or encoded incorrectly")
						}
						if r.Header.Get("X-Test-Fixture") != "fixture" {
							t.Error("explicit header missing")
						}
						if (signed || endpoint.Direction == "inbound-callback" || endpoint.Authentication == "body-signature") && r.Header.Get("Authorization") != "" {
							t.Error("AGC bearer leaked into signed upload")
						}
						if !signed && endpoint.Direction != "inbound-callback" && endpoint.Authentication != "body-signature" && r.Header.Get("Authorization") != "Bearer offline-token" {
							t.Error("authorization missing")
						}
						if endpoint.Method == "POST" || endpoint.Method == "PUT" || endpoint.Method == "DELETE" {
							if strings.HasPrefix(endpoint.RequestContentType, "application/x-www-form-urlencoded") {
								if err := r.ParseForm(); err != nil || r.PostForm.Get("probeField") != "payload" {
									t.Error("form body field was lost")
								}
							} else if resourceUpload {
								if err := r.ParseMultipartForm(1024); err != nil || r.FormValue("probeField") != "payload" {
									t.Error("multipart field was lost")
								}
								if r.MultipartForm != nil {
									defer r.MultipartForm.RemoveAll()
								}
							} else {
								body, err := io.ReadAll(r.Body)
								var fields map[string]any
								if err != nil || json.Unmarshal(body, &fields) != nil || (endpoint.DynamicBody == nil && fields["probeField"] != "payload") {
									t.Error("JSON body field was lost")
								}
							}
						}
						w.Header().Set("Content-Type", "application/json")
						w.WriteHeader(scenario.status)
						_, _ = io.WriteString(w, scenario.response)
					}))
					defer server.Close()
					args := []string{endpoint.FamilyID, endpoint.ID, "--invoke", "--dry-run=false", "--base-url", server.URL, "--token", "offline-token", "--query", "probe=a b&c", "--header", "X-Test-Fixture=fixture"}
					for _, p := range endpoint.Parameters {
						if resourceUpload && p.In == "file" {
							continue
						}
						if !p.Required {
							continue
						}
						value := parameterFixture(p)
						if p.Name == "callbackUrl" || p.Name == "uploadUrl" {
							value = server.URL + "/target"
						}
						option := map[string]string{"path": "--param", "query": "--query", "header": "--header", "body": "--field", "file": "--field"}[p.In]
						args = append(args, option, p.Name+"="+value)
					}
					if resourceUpload {
						file := filepath.Join(t.TempDir(), "fixture.png")
						if err := os.WriteFile(file, []byte("offline-image"), 0600); err != nil {
							t.Fatal(err)
						}
						args = append(args, "--file", file)
					}
					if endpoint.Method == "POST" || endpoint.Method == "PUT" || endpoint.Method == "DELETE" {
						if endpoint.DynamicBody == nil {
							args = append(args, "--field", "probeField=payload")
						}
					}
					if endpoint.DynamicBody != nil {
						args = append(args, "--field", "additionalProp1="+parameterFixture(*endpoint.DynamicBody))
					}
					if endpoint.FamilyID == "resources" && endpoint.ID == "respackapi-queryfile" {
						args = append(args, "--query", "versionId=1")
					}
					if endpoint.FamilyID == "projects" && endpoint.ID == "appbriefinfo" {
						args = append(args, "--query", "packageNames=app.example.fixture")
					}
					if endpoint.FamilyID == "publishing" && endpoint.ID == "app-version-list" {
						args = append(args, "--header", "appId=fixture")
					}
					if endpoint.FamilyID == "upload" && endpoint.ID == "obbfile-init" {
						args = append(args, "--query", "fileType=1")
					}
					if endpoint.FamilyID == "publishing" && endpoint.ID == "app-create" {
						args = append(args, "--field", "packageName=app.example.fixture")
					}
					out, err := execute(args...)
					select {
					case <-seen:
					default:
						t.Fatal("CLI did not send its request")
					}
					if (err != nil) != scenario.wantError {
						t.Fatalf("error = %v, wantError %v", err, scenario.wantError)
					}
					if err == nil {
						var envelope domain.Envelope[agcapi.InvokeResponse]
						if json.Unmarshal([]byte(out), &envelope) != nil || envelope.Data.DryRun || envelope.Data.StatusCode != 200 || !strings.Contains(string(envelope.Data.Body), "received") {
							t.Fatalf("CLI response was lost: %s", out)
						}
					}
				})
			}
		})
	}
}

func TestReportFilterQueryPreservesEveryValueAndItsOrder(t *testing.T) {
	called := false
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		for key, want := range map[string][]string{"filterCondition": {"countryId", "businessType"}, "filterConditionValue": {"CN", "AppGallery"}} {
			got := r.URL.Query()[key]
			if strings.Join(got, "|") != strings.Join(want, "|") {
				t.Errorf("%s lost ordered query values: %v", key, got)
			}
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"code":0}`))
	}))
	defer server.Close()
	_, err := execute("--project", t.TempDir(), "reports", "appdownloadexport", "--invoke", "--dry-run=false", "--token", "offline-report-token", "--base-url", server.URL, "--param", "appId=fixture", "--query", "language=en-US", "--query", "startTime=20260901", "--query", "endTime=20261001", "--query", "filterCondition=countryId", "--query", "filterCondition=businessType", "--query", "filterConditionValue=CN", "--query", "filterConditionValue=AppGallery")
	if err != nil {
		t.Fatal(err)
	}
	if !called {
		t.Fatal("transport was not exercised")
	}
}

func TestGameCallbackCLIWorksWithoutHuaweiAuthorization(t *testing.T) {
	received := make(chan struct{}, 1)
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		received <- struct{}{}
		if r.Header.Get("Authorization") != "" {
			t.Error("unexpected Huawei authorization on callback")
		}
		if err := r.ParseForm(); err != nil || r.PostForm.Get("sign") != "a+b/=" {
			t.Error("callback signature was not encoded exactly once")
		}
		_, _ = w.Write([]byte(`{"code":0}`))
	}))
	defer upstream.Close()
	_, err := execute("--project", t.TempDir(), "game-items", "propapi-getroleinfo", "--invoke", "--dry-run=false", "--param", "callbackUrl="+upstream.URL, "--field", "userId=offline-player", "--field", "userIdType=1", "--field", "ts=123", "--field", "sign=a+b/=")
	if err != nil {
		t.Fatal(err)
	}
	select {
	case <-received:
	default:
		t.Fatal("CLI did not invoke callback server")
	}
}
