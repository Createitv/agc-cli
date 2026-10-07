package domain

import "testing"

func TestTestPackageMetadataContract(t *testing.T) {
	e, _ := EndpointByID("testing", "test-api-add-test-package")
	for _, expected := range []struct{ name, where, kind string }{{"appId", "query", "string"}, {"file", "body", "object"}} {
		found := false
		for _, p := range e.Parameters {
			if p.Name == expected.name && p.In == expected.where && p.Type == expected.kind && p.Required {
				found = true
			}
		}
		if !found {
			t.Errorf("missing required %s in %s of type %s", expected.name, expected.where, expected.kind)
		}
	}
}
