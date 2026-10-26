package command

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestPrivacyCLIRejectsDeepInputBeforeNetwork(t *testing.T) {
	called := false
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		t.Error("invalid input reached transport")
	}))
	defer upstream.Close()
	out, err := execute("--project", t.TempDir(), "publishing", "privacy-agreement-update", "--invoke", "--dry-run=false", "--token", "offline-fixture", "--base-url", upstream.URL, "--header", "appId=fixture", "--field", "id=fixture", "--field", `agreementInfos=[{"lang":"zh","thirdSdkInfo":{"sdkInfos":[{"privacyLabels":[{"dataItems":[{"sensitiveData":2}]}]}]}}]`)
	if err == nil || called {
		t.Fatal("deep input must fail locally")
	}
	var result map[string]any
	if json.Unmarshal([]byte(out), &result) != nil {
		t.Fatalf("JSON stdout not parseable: %s", out)
	}
	for _, want := range []string{"agreementInfos[0].thirdSdkInfo.sdkInfos[0].privacyLabels[0].dataItems[0].sensitiveData", "data-dataitem", "invalid_request", "documentationCommand", "suggestion"} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %s: %s", want, out)
		}
	}
}
