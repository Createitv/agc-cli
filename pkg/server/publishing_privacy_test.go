package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestPrivacyRESTRejectsDeepInputBeforeNetwork(t *testing.T) {
	called := false
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		t.Error("invalid input reached transport")
	}))
	defer upstream.Close()
	payload := `{"baseUrl":"` + upstream.URL + `","token":"offline-fixture","headers":{"appId":"fixture"},"dryRun":false,"body":{"id":"fixture","agreementInfos":[{"lang":"zh","extInfo":{"storageInfo":{"isCrossBorder":"false"}}}]}}`
	request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(payload))
	response := httptest.NewRecorder()
	invokeEndpoint(response, request, "publishing", "privacy-agreement-update")
	if response.Code != http.StatusBadRequest || called {
		t.Fatalf("invalid input should fail before HTTP: %d %s", response.Code, response.Body.String())
	}
	var result map[string]any
	if json.Unmarshal(response.Body.Bytes(), &result) != nil {
		t.Fatal("REST feedback not JSON")
	}
	for _, want := range []string{"agreementInfos[0].extInfo.storageInfo.isCrossBorder", "data-storageinfo", "invalid_request", "documentationCommand", "suggestion"} {
		if !strings.Contains(response.Body.String(), want) {
			t.Errorf("missing %s: %s", want, response.Body.String())
		}
	}
}
