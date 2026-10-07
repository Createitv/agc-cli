package domain

import "testing"

func TestDraftUpdateRequiresAppQueryAndVersionBody(t *testing.T) {
	e, _ := EndpointByID("testing", "test-api-modify-test-version")
	for _, expected := range []struct{ name, where string }{{"appId", "query"}, {"versionId", "body"}} {
		found := false
		for _, p := range e.Parameters {
			if p.Name == expected.name && p.In == expected.where && p.Required {
				found = true
			}
		}
		if !found {
			t.Errorf("missing required %s in %s", expected.name, expected.where)
		}
	}
}
