package command

import (
	"encoding/json"
	"github.com/Createitv/agc-cli/pkg/domain"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCLIPreservesBusinessSubcodeOriginWithoutRawBody(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"rtnCode":0,"businessCode":2,"rtnDesc":"denied","invitationCode":"raw-secret"}`))
	}))
	defer server.Close()
	out, err := execute("--project", t.TempDir(), "testing", "test-api-get-test-grouplist", "--invoke", "--dry-run=false", "--header", "appId=fixture", "--token", "offline", "--base-url", server.URL)
	if err == nil {
		t.Fatal("failed business subcode accepted")
	}
	var envelope domain.ErrorEnvelope
	if json.Unmarshal([]byte(out), &envelope) != nil || envelope.Error.CodeField != "businessCode" || envelope.Error.ResultCodes["rtnCode"] != "0" || envelope.Error.ResultCodes["businessCode"] != "2" {
		t.Fatalf("code origin missing: %s", out)
	}
	if strings.Contains(out, "raw-secret") || strings.Contains(out, "rawBody") {
		t.Fatal("raw response secret exposed")
	}
}
