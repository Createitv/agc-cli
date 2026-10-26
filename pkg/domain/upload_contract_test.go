package domain

import "testing"

func TestMultipartInitDeclaresPackageKindAndOBBKind(t *testing.T) {
	e, _ := EndpointByID("upload", "obbfile-init")
	for _, name := range []string{"fileType", "obbType"} {
		found := false
		for _, p := range e.Parameters {
			if p.Name == name && p.In == "query" && p.Type == "integer" && !p.Required {
				found = true
			}
		}
		if !found {
			t.Errorf("official init query parameter %s missing", name)
		}
	}
}

func TestMultipartQueryContracts(t *testing.T) {
	for _, c := range []struct{ id, name string }{
		{"obbfile-init", "appId"}, {"obbfile-init", "fileName"},
		{"obbfile-uploadurl", "objectId"}, {"obbfile-uploadurl", "nspUploadId"},
		{"obbfile-compose", "objectId"}, {"obbfile-compose", "nspUploadId"},
	} {
		t.Run(c.id+"/"+c.name, func(t *testing.T) {
			e, _ := EndpointByID("upload", c.id)
			for _, p := range e.Parameters {
				if p.Name == c.name && p.In == "query" && p.Required {
					return
				}
			}
			t.Fatalf("missing required query %s", c.name)
		})
	}
}
