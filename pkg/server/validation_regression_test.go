package server

import (
	"github.com/Createitv/agc-cli/pkg/domain"
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
