package agcapi

import (
	"github.com/Createitv/agc-cli/pkg/domain"
	"strings"
	"testing"
)

func TestTestingPackageFileRequiresOfficialModelFields(t *testing.T) {
	e, _ := domain.EndpointByID("testing", "test-api-add-test-package")
	for _, body := range []string{`{"file":{}}`, `{"file":{"fileName":"owned.app"}}`, `{"file":{"objectId":"owned-object"}}`, `{"file":{"fileName":123,"objectId":"owned-object"}}`, `{"file":{"fileName":"` + strings.Repeat("x", 256) + `","objectId":"owned-object"}}`, `{"file":{"fileName":"owned.app","objectId":"` + strings.Repeat("x", 1025) + `"}}`} {
		if err := ValidateEndpointRequest(e, map[string]string{"appId": "dedicated"}, nil, []byte(body)); err == nil {
			t.Error("invalid nested FileInfo accepted")
		}
	}
	if err := ValidateEndpointRequest(e, map[string]string{"appId": "dedicated"}, nil, []byte(`{"file":{"fileName":"owned.app","objectId":"owned-object"}}`)); err != nil {
		t.Fatal(err)
	}
	doc, _ := domain.EndpointReference("testing", "test-api-add-test-package")
	if !strings.Contains(string(doc), "agc-help-test-api-data-fileinfo-0000002272575549") || !strings.Contains(string(doc), "relatedDocuments") {
		t.Fatal("linked FileInfo source not exposed")
	}
}
