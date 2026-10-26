package server

import (
	"github.com/Createitv/agc-cli/pkg/domain"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRequiredBodyFieldCannotBeBypassedByUnrelatedJSON(t *testing.T) {
	endpoint := domain.Endpoint{Parameters: []domain.Parameter{{Name: "appId", In: "body", Required: true}}}
	for _, body := range []string{`{}`, `{"other":"x"}`, `{"appId":null}`, `{"appId":""}`, `[]`, `not-json`} {
		if err := validateInvokeParameters(endpoint, nil, nil, nil, nil, []byte(body)); err == nil {
			t.Errorf("accepted missing required field in %s", body)
		}
	}
	for _, body := range []string{`{"appId":"123"}`, `{"appId":123}`} {
		if err := validateInvokeParameters(endpoint, nil, nil, nil, nil, []byte(body)); err != nil {
			t.Errorf("rejected valid body: %v", err)
		}
	}
}

func TestRequiredHeaderNamesAreCaseInsensitive(t *testing.T) {
	endpoint := domain.Endpoint{Parameters: []domain.Parameter{{Name: "appId", In: "header", Required: true}}}
	if err := validateInvokeParameters(endpoint, nil, nil, map[string]string{"AppID": "123"}, nil, nil); err != nil {
		t.Fatal(err)
	}
}

func TestRESTAgentGetsConditionalRuleAndSource(t *testing.T) {
	request := httptest.NewRequest(http.MethodPost, "/api/v1/publishing/endpoints/app-info-query-v3/invoke", strings.NewReader(`{"query":{"appId":"fixture","releaseType":"6"},"dryRun":true}`))
	response := httptest.NewRecorder()
	Handler().ServeHTTP(response, request)
	if response.Code != http.StatusBadRequest {
		t.Fatalf("got %d: %s", response.Code, response.Body.String())
	}
	for _, want := range []string{"versionId", "sourceUrl", "suggestion"} {
		if !strings.Contains(response.Body.String(), want) {
			t.Errorf("feedback missing %s", want)
		}
	}
}

func TestReportRESTPreservesOrderedQueryValues(t *testing.T) {
	called := false
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		for key, want := range map[string][]string{"filterCondition": {"countryId", "businessType"}, "filterConditionValue": {"CN", "AppGallery"}} {
			if got := r.URL.Query()[key]; strings.Join(got, "|") != strings.Join(want, "|") {
				t.Errorf("%s lost ordered values: %v", key, got)
			}
		}
		_, _ = w.Write([]byte(`{"code":0}`))
	}))
	defer upstream.Close()
	payload := `{"baseUrl":"` + upstream.URL + `","params":{"appId":"fixture"},"queryValues":{"language":["en-US"],"startTime":["20260901"],"endTime":["20261001"],"filterCondition":["countryId","businessType"],"filterConditionValue":["CN","AppGallery"]},"dryRun":false}`
	request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(payload))
	response := httptest.NewRecorder()
	invokeEndpoint(response, request, "reports", "appdownloadexport")
	if response.Code != http.StatusOK {
		t.Fatalf("REST invocation failed: %s", response.Body.String())
	}
	if !called {
		t.Fatal("REST did not invoke upstream")
	}
}
