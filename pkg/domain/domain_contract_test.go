package domain

import "testing"

func TestDomainWriteRequiresHeaderAppAndQueryCategory(t *testing.T) {
	e, _ := EndpointByID("domains", "domain-api-post-domain")
	for _, expected := range []struct{ name, where string }{{"appId", "header"}, {"category", "query"}} {
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
