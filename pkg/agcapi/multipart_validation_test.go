package agcapi

import (
	"github.com/Createitv/agc-cli/pkg/domain"
	"testing"
)

func TestMultipartInitEnforcesOfficialOBBBranches(t *testing.T) {
	e, _ := domain.EndpointByID("upload", "obbfile-init")
	for _, c := range []struct {
		name    string
		query   map[string]string
		invalid bool
	}{
		{"default OBB needs subtype", map[string]string{}, true},
		{"explicit OBB needs subtype", map[string]string{"fileType": "2"}, true},
		{"APP forbids OBB subtype", map[string]string{"fileType": "1", "obbType": "1"}, true},
		{"APK forbids OBB subtype", map[string]string{"fileType": "3", "obbType": "2"}, true},
		{"default OBB with subtype", map[string]string{"obbType": "1"}, false},
		{"APP without subtype", map[string]string{"fileType": "1"}, false},
	} {
		t.Run(c.name, func(t *testing.T) {
			c.query["appId"] = "fixture"
			c.query["fileName"] = "fixture.app"
			err := ValidateEndpointRequest(e, c.query, nil, nil)
			if (err != nil) != c.invalid {
				t.Fatalf("invalid=%v, got %v", c.invalid, err)
			}
		})
	}
}
