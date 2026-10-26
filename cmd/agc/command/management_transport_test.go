package command

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Createitv/agc-cli/internal/managementtest"
	"github.com/Createitv/agc-cli/pkg/domain"
)

func managementArgs(t *testing.T, e domain.Endpoint, f managementtest.Fixture, base string, dry bool) []string {
	t.Helper()
	session := "/must-not-read-management-session"
	if !dry {
		session = managementtest.SessionFile(t, base)
	}
	args := []string{"--project", t.TempDir(), e.FamilyID, e.ID, "--invoke", "--session-file", session, "--base-url", base}
	if !dry {
		args = append(args, "--dry-run=false")
	}
	for k, v := range f.Params {
		args = append(args, "--param", k+"="+v)
	}
	for k, v := range f.Query {
		args = append(args, "--query", k+"="+v)
	}
	if len(f.Body) > 0 {
		p := filepath.Join(t.TempDir(), "body.json")
		b, _ := json.Marshal(f.Body)
		os.WriteFile(p, b, 0600)
		args = append(args, "--body", p)
	}
	return args
}

func testReviewedManagementCLI(t *testing.T, e domain.Endpoint) {
	if e.FamilyID == "templates" && e.ID == "deploy" {
		testTemplateDeployCLI(t, e, false)
		return
	}
	t.Run(e.FamilyID+"/"+e.ID, func(t *testing.T) {
		f, ok := managementtest.Fixtures()[e.FamilyID+"/"+e.ID]
		if !ok {
			t.Fatal("missing independent source fixture")
		}
		for _, scenario := range []struct {
			name   string
			status int
			body   string
			fail   bool
		}{{"success", 200, `{"code":0,"result":{"marker":"received"}}`, false}, {"business-error", 200, `{"code":17,"message":"denied"}`, true}, {"http-error", 403, `{"code":19,"message":"forbidden"}`, true}} {
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
				out, err := execute(managementArgs(t, e, f, srv.URL, false)...)
				if calls != 1 || (err != nil) != scenario.fail {
					t.Fatalf("calls=%d error=%v output=%s", calls, err, out)
				}
				if !json.Valid([]byte(out)) {
					t.Fatal("non-JSON feedback")
				}
			})
		}
	})
}

func testReviewedManagementDryRun(t *testing.T, e domain.Endpoint) {
	if e.FamilyID == "templates" && e.ID == "deploy" {
		testTemplateDeployCLI(t, e, true)
		return
	}
	t.Run(e.FamilyID+"/"+e.ID, func(t *testing.T) {
		f, ok := managementtest.Fixtures()[e.FamilyID+"/"+e.ID]
		if !ok {
			t.Fatal("fixture missing")
		}
		out, err := execute(managementArgs(t, e, f, "https://offline.example", true)...)
		if err != nil || !strings.Contains(out, `"dryRun":true`) {
			t.Fatalf("%v %s", err, out)
		}
	})
}

func testTemplateDeployCLI(t *testing.T, e domain.Endpoint, dry bool) {
	t.Run("templates/deploy", func(t *testing.T) {
		srv, body := managementtest.Deployment(t)
		defer srv.Close()
		f := managementtest.Fixture{Body: body}
		out, err := execute(managementArgs(t, e, f, srv.URL, dry)...)
		if err != nil {
			t.Fatalf("%v %s", err, out)
		}
		if !dry && !strings.Contains(out, `"ready":true`) {
			t.Fatal("deployment readiness not read back")
		}
	})
}

func testPrefetchDeviceCLI(t *testing.T, e domain.Endpoint, dry bool) {
	t.Run(e.FamilyID+"/"+e.ID, func(t *testing.T) {
		args := []string{"prefetch", e.ID, "--invoke", "--field", "bundleName=com.example.offline", "--field", "deviceId=offline-device"}
		hdc := "/must-not-execute"
		if !dry {
			hdc = managementtest.HDCFile(t)
			args = append(args, "--dry-run=false")
		}
		args = append(args, "--field", "hdcPath="+hdc)
		out, err := execute(args...)
		if err != nil {
			t.Fatalf("%v %s", err, out)
		}
		if !dry && !strings.Contains(out, `"dataFreshnessVerified":false`) {
			t.Fatal("device debug misreported")
		}
	})
}
