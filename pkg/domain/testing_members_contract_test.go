package domain

import "testing"

func TestTestingMemberContractsMatchOfficialBodyModels(t *testing.T) {
	for _, c := range []struct {
		id, name, kind string
		required       bool
	}{
		{"test-api-add-test-group-user", "groupId", "string", true},
		{"test-api-add-test-group-user", "hwAccount", "string", true},
		{"test-api-add-test-group-user", "nickName", "string", false},
		{"test-api-add-test-group-user", "uid", "string", false},
		{"test-api-delete-test-group-user", "groupId", "string", true},
		{"test-api-delete-test-group-user", "deleteTesterIds", "array", true},
	} {
		t.Run(c.id+"/"+c.name, func(t *testing.T) {
			e, _ := EndpointByID("testing", c.id)
			for _, p := range e.Parameters {
				if p.Name == c.name && p.In == "body" && p.Type == c.kind && p.Required == c.required {
					if c.name == "deleteTesterIds" && (p.ItemsType != "string" || p.MaxItems != 100) {
						t.Fatal("missing official string-array maximum")
					}
					return
				}
			}
			t.Fatalf("official body field %s missing or wrong", c.name)
		})
	}
}
