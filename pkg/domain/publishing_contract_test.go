package domain

import "testing"

func TestHarmonyCompileStatusContract(t *testing.T) {
	e, ok := EndpointByID("publishing", "package-compile-status-v3")
	if !ok || e.Method != "GET" || e.Path != "/api/publish/v3/package/compile/status" {
		t.Fatal("missing current HarmonyOS compile status endpoint")
	}
	for _, name := range []string{"appId", "pkgIds"} {
		found := false
		for _, p := range e.Parameters {
			if p.Name == name && p.In == "query" && p.Required {
				found = true
			}
		}
		if !found {
			t.Errorf("missing required query %s", name)
		}
	}
}

func TestLanguageCommandsRequireApplicationAndLanguage(t *testing.T) {
	for _, c := range []struct{ id, name, where string }{
		{"language-info-update", "appId", "query"},
		{"language-info-update", "lang", "body"},
		{"language-info-delete", "appId", "query"},
		{"language-info-delete", "lang", "query"},
	} {
		t.Run(c.id+"/"+c.name, func(t *testing.T) {
			e, _ := EndpointByID("publishing", c.id)
			for _, p := range e.Parameters {
				if p.Name == c.name && p.In == c.where && p.Required {
					return
				}
			}
			t.Fatalf("missing %s@%s", c.name, c.where)
		})
	}
}
