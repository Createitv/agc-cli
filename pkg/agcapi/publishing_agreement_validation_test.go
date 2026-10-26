package agcapi

import (
	"errors"
	"strings"
	"testing"

	"github.com/Createitv/agc-cli/pkg/domain"
)

// Independent expectations from the AgreementInfo and AgreementBaseInfo
// official pages, rather than fixtures generated from the endpoint registry.
func TestPublishingAgreementNestedRequirements(t *testing.T) {
	for _, id := range []string{"user-agreement-create", "user-agreement-update", "privacy-agreement-update"} {
		for _, tc := range []struct{ name, item, field string }{
			{"language-required", `{}`, "body.agreementInfos[0].lang"},
			{"language-enum", `{"lang":"zh-CN"}`, "body.agreementInfos[0].lang"},
			{"base-name-required", `{"lang":"zh","baseInfo":{}}`, "body.agreementInfos[0].baseInfo.name"},
			{"base-name-length", `{"lang":"zh","baseInfo":{"name":"` + strings.Repeat("界", 101) + `"}}`, "body.agreementInfos[0].baseInfo.name"},
		} {
			t.Run(id+"/"+tc.name, func(t *testing.T) {
				e, _ := domain.EndpointByID("publishing", id)
				err := ValidateEndpointRequest(e, nil, map[string]string{"appId": "fixture"}, []byte(`{"id":"fixture","type":2,"agreementInfos":[`+tc.item+`]}`))
				var detail *RequestValidationError
				if !errors.As(err, &detail) {
					t.Fatalf("expected nested violation, got %v", err)
				}
				for _, v := range detail.Violations {
					if v.Field == tc.field && v.SourceURL != "" && v.Suggestion != "" {
						return
					}
				}
				t.Fatalf("missing sourced violation for %s: %+v", tc.field, detail.Violations)
			})
		}
	}
}

func TestPublishingAgreementOptionalBaseInfo(t *testing.T) {
	for _, id := range []string{"user-agreement-create", "user-agreement-update", "privacy-agreement-update"} {
		e, _ := domain.EndpointByID("publishing", id)
		for _, body := range []string{
			`{"id":"fixture","type":2,"agreementInfos":[{"lang":"zh"}]}`,
			`{"id":"fixture","type":2,"agreementInfos":[{"lang":"zh-Hant-TW","baseInfo":{"name":"测试协议"}}]}`,
		} {
			if err := ValidateEndpointRequest(e, nil, map[string]string{"appId": "fixture"}, []byte(body)); err != nil {
				t.Fatalf("%s: %v", id, err)
			}
		}
	}
}
