package agcapi

import (
	"errors"
	"strings"
	"testing"

	"github.com/Createitv/agc-cli/pkg/domain"
)

// These failures follow CustInfo, DeviceTypeInfo and AppReviewContact,
// independently of the top-level Publishing registry.
func TestPublishingMetadataNestedModels(t *testing.T) {
	e, _ := domain.EndpointByID("publishing", "app-info-update-v3")
	for _, tc := range []struct{ name, fields, path string }{
		{"device-required", `"deviceTypes":[{}]`, "body.deviceTypes[0].deviceType"},
		{"device-enum", `"deviceTypes":[{"deviceType":1}]`, "body.deviceTypes[0].deviceType"},
		{"device-integer", `"deviceTypes":[{"deviceType":"4"}]`, "body.deviceTypes[0].deviceType"},
		{"customer-phone", `"customInfo":{"phone":"` + strings.Repeat("1", 51) + `"}`, "body.customInfo.phone"},
		{"review-contact", `"appReviewEmailInfo":{"account":"` + strings.Repeat("a", 501) + `"}`, "body.appReviewEmailInfo.account"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateEndpointRequest(e, map[string]string{"appId": "fixture"}, nil, []byte(`{"publishCountry":"CN","encrypted":0,`+tc.fields+`}`))
			var detail *RequestValidationError
			if !errors.As(err, &detail) {
				t.Fatalf("expected nested violation, got %v", err)
			}
			for _, v := range detail.Violations {
				if v.Field == tc.path && v.SourceURL != "" {
					return
				}
			}
			t.Fatalf("missing source-based violation %s: %+v", tc.path, detail.Violations)
		})
	}
}

func TestPublishingMetadataDocumentedDeviceValues(t *testing.T) {
	e, _ := domain.EndpointByID("publishing", "app-info-update-v3")
	body := []byte(`{"publishCountry":"CN","encrypted":0,"deviceTypes":[{"deviceType":4,"appAdapters":"5"},{"deviceType":8,"remoteControl":"1"},{"deviceType":19}],"customInfo":{},"appReviewEmailInfo":{}}`)
	if err := ValidateEndpointRequest(e, map[string]string{"appId": "fixture"}, nil, body); err != nil {
		t.Fatal(err)
	}
}
