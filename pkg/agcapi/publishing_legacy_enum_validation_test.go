package agcapi

import (
	"encoding/json"
	"github.com/Createitv/agc-cli/pkg/domain"
	"testing"
)

func TestLegacyAppInfoRejectsValuesOutsideOfficialEnums(t *testing.T) {
	e, _ := domain.EndpointByID("publishing", "app-info-update")
	for _, body := range []string{`{"privacyPolicy":"https://example.com/privacy","registeredDclType":-1}`, `{"privacyPolicy":"https://example.com/privacy","registeredIdType":-1}`, `{"privacyPolicy":"https://example.com/privacy","isFree":1}`} {
		if err := ValidateEndpointRequest(e, map[string]string{"appId": "dedicated"}, nil, []byte(body)); err != nil {
			t.Fatalf("documented value rejected: %v", err)
		}
	}
	for _, name := range []string{"releaseType", "appNetType", "isFree", "isAppForcedUpdate", "hispaceAutoDown", "familyShareTag", "webGameFlag", "registeredDclType", "registeredIdType", "registeredEntity", "encrypted"} {
		t.Run(name, func(t *testing.T) {
			query := map[string]string{"appId": "dedicated"}
			body := map[string]any{"privacyPolicy": "https://example.com/privacy"}
			field := "body." + name
			if name == "releaseType" {
				query[name] = "999"
				field = "query." + name
			} else {
				body[name] = 999
			}
			data, _ := json.Marshal(body)
			err := ValidateEndpointRequest(e, query, nil, data)
			invalid, ok := err.(*RequestValidationError)
			if !ok {
				t.Fatalf("undocumented value accepted: %v", err)
			}
			for _, v := range invalid.Violations {
				if v.Field == field && v.Code == "enum" && v.SourceURL != "" {
					return
				}
			}
			t.Fatalf("missing sourced enum violation for %s: %v", field, invalid.Violations)
		})
	}
}

func TestLegacyAppInfoRequiresDocumentedConditionalFields(t *testing.T) {
	e, _ := domain.EndpointByID("publishing", "app-info-update")
	for _, tc := range []struct{ body, field string }{
		{`{"privacyPolicy":"https://example.com/privacy","isFree":0}`, "body.price"},
		{`{"privacyPolicy":"https://example.com/privacy","isFree":0,"price":"{}"}`, "body.priceDetail"},
		{`{"privacyPolicy":"https://example.com/privacy","registeredDclType":0}`, "body.registeredEntity"},
		{`{"privacyPolicy":"https://example.com/privacy","registeredDclType":0,"registeredEntity":1}`, "body.registeredEntityName"},
	} {
		err := ValidateEndpointRequest(e, map[string]string{"appId": "dedicated"}, nil, []byte(tc.body))
		invalid, ok := err.(*RequestValidationError)
		if !ok {
			t.Fatalf("missing conditional field %s accepted: %v", tc.field, err)
		}
		found := false
		for _, v := range invalid.Violations {
			found = found || v.Field == tc.field && v.SourceURL != ""
		}
		if !found {
			t.Errorf("missing sourced correction for %s: %v", tc.field, invalid.Violations)
		}
	}
}
