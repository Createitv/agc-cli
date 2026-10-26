package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Createitv/agc-cli/internal/managementtest"
	"github.com/Createitv/agc-cli/pkg/domain"
)

func testReviewedManagementREST(t *testing.T, e domain.Endpoint) {
	if e.FamilyID == "templates" && e.ID == "deploy" {
		testTemplateDeployREST(t)
		return
	}
	t.Run(e.FamilyID+"/"+e.ID, func(t *testing.T) {
		f, ok := managementtest.Fixtures()[e.FamilyID+"/"+e.ID]
		if !ok {
			t.Fatal("source fixture missing")
		}
		for _, scenario := range []struct {
			name, body string
			status     int
			fail       bool
		}{{"success", `{"code":0,"result":{"marker":"received"}}`, 200, false}, {"business-error", `{"code":17,"message":"denied"}`, 200, true}, {"http-error", `{"code":19,"message":"forbidden"}`, 403, true}} {
			t.Run(scenario.name, func(t *testing.T) {
				calls := 0
				srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					calls++
					managementtest.AssertRequest(t, r, f)
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(scenario.status)
					w.Write([]byte(scenario.body))
				}))
				defer srv.Close()
				payload := map[string]any{"params": f.Params, "query": f.Query, "body": f.Body, "sessionFile": managementtest.SessionFile(t, srv.URL), "dryRun": false}
				if f.Body == nil {
					delete(payload, "body")
				}
				b, _ := json.Marshal(payload)
				w := httptest.NewRecorder()
				Handler().ServeHTTP(w, httptest.NewRequest("POST", "/api/v1/"+e.FamilyID+"/endpoints/"+e.ID+"/invoke", bytes.NewReader(b)))
				if calls != 1 || (w.Code != 200) != scenario.fail {
					t.Fatalf("calls=%d status=%d %s", calls, w.Code, w.Body.String())
				}
			})
		}
	})
}

func testTemplateDeployREST(t *testing.T) {
	t.Run("templates/deploy", func(t *testing.T) {
		srv, body := managementtest.Deployment(t)
		defer srv.Close()
		p := map[string]any{"body": body, "sessionFile": managementtest.SessionFile(t, srv.URL), "dryRun": false}
		b, _ := json.Marshal(p)
		w := httptest.NewRecorder()
		Handler().ServeHTTP(w, httptest.NewRequest("POST", "/api/v1/templates/endpoints/deploy/invoke", bytes.NewReader(b)))
		if w.Code != 200 {
			t.Fatalf("%d %s", w.Code, w.Body.String())
		}
	})
}

func testPrefetchDeviceREST(t *testing.T, e domain.Endpoint) {
	t.Run(e.FamilyID+"/"+e.ID, func(t *testing.T) {
		payload := map[string]any{"body": map[string]string{"bundleName": "com.example.offline", "deviceId": "offline-device", "hdcPath": managementtest.HDCFile(t)}, "dryRun": false}
		b, _ := json.Marshal(payload)
		w := httptest.NewRecorder()
		Handler().ServeHTTP(w, httptest.NewRequest("POST", "/api/v1/prefetch/endpoints/device-getcache/invoke", bytes.NewReader(b)))
		if w.Code != 200 {
			t.Fatalf("%d %s", w.Code, w.Body.String())
		}
	})
}
