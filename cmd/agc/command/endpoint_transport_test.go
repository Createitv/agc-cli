package command

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
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
						if signed && r.Header.Get("Authorization") != "" {
							t.Error("AGC bearer leaked into signed upload")
						}
						if !signed && r.Header.Get("Authorization") != "Bearer offline-token" {
							t.Error("authorization missing")
						}
						if endpoint.Method == "POST" || endpoint.Method == "PUT" || endpoint.Method == "DELETE" {
							body, err := io.ReadAll(r.Body)
							var fields map[string]any
							if err != nil || json.Unmarshal(body, &fields) != nil || fields["probeField"] != "payload" {
								t.Error("JSON body field was lost")
							}
						}
						w.Header().Set("Content-Type", "application/json")
						w.WriteHeader(scenario.status)
						_, _ = io.WriteString(w, scenario.response)
					}))
					defer server.Close()
					args := []string{endpoint.FamilyID, endpoint.ID, "--invoke", "--dry-run=false", "--base-url", server.URL, "--token", "offline-token", "--query", "probe=a b&c", "--header", "X-Test-Fixture=fixture"}
					for _, p := range endpoint.Parameters {
						if !p.Required {
							continue
						}
						value := "fixture"
						if p.In == "body" && (p.Type == "integer" || p.Type == "number") {
							value = "1"
						}
						if p.In == "body" && p.Type == "array" {
							value = `["fixture"]`
						}
						if p.In == "body" && p.Type == "object" {
							value = `{"key":"fixture"}`
						}
						if p.Name == "callbackUrl" || p.Name == "uploadUrl" {
							value = server.URL + "/target"
						}
						option := map[string]string{"path": "--param", "query": "--query", "header": "--header", "body": "--field", "file": "--field"}[p.In]
						args = append(args, option, p.Name+"="+value)
					}
					if endpoint.Method == "POST" || endpoint.Method == "PUT" || endpoint.Method == "DELETE" {
						args = append(args, "--field", "probeField=payload")
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
