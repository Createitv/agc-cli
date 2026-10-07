package domain

import "testing"

func TestLegacyProductGroupContracts(t *testing.T) {
	for _, id := range []string{"addproductgroup-android", "updateproductgroup-android", "getproductgroup-android"} {
		e, _ := EndpointByID("pms", id)
		for _, expected := range []struct{ name, where string }{{"appId", "header"}, {"requestId", "body"}} {
			found := false
			for _, p := range e.Parameters {
				if p.Name == expected.name && p.In == expected.where && p.Required {
					found = true
				}
			}
			if !found {
				t.Errorf("%s missing %s in %s", id, expected.name, expected.where)
			}
		}
	}
}

func TestHarmonyProductCreateContract(t *testing.T) {
	e, _ := EndpointByID("pms", "addproduct-harmonyosnext")
	for _, expected := range []struct{ name, where, kind string }{{"appId", "header", ""}, {"requestId", "body", "string"}, {"product", "body", "object"}} {
		found := false
		for _, p := range e.Parameters {
			if p.Name == expected.name && p.In == expected.where && p.Required && (expected.kind == "" || p.Type == expected.kind) {
				found = true
			}
		}
		if !found {
			t.Errorf("missing %s in %s of type %s", expected.name, expected.where, expected.kind)
		}
	}
}

func TestHarmonyProductGroupCreateContract(t *testing.T) {
	e, _ := EndpointByID("pms", "addproductgroup-harmonyosnext")
	for _, expected := range []struct{ name, where, kind string }{{"appId", "header", ""}, {"requestId", "body", "string"}, {"resource", "body", "object"}} {
		found := false
		for _, p := range e.Parameters {
			if p.Name == expected.name && p.In == expected.where && p.Required && (expected.kind == "" || p.Type == expected.kind) {
				found = true
			}
		}
		if !found {
			t.Errorf("missing %s in %s of type %s", expected.name, expected.where, expected.kind)
		}
	}
}

func TestHarmonyProductGroupUpdateAndListContract(t *testing.T) {
	for _, id := range []string{"updateproductgroup-harmonyosnext", "getproductgroup-harmonyosnext"} {
		e, _ := EndpointByID("pms", id)
		for _, expected := range []struct{ name, where string }{{"appId", "header"}, {"requestId", "body"}} {
			found := false
			for _, p := range e.Parameters {
				if p.Name == expected.name && p.In == expected.where && p.Required {
					found = true
				}
			}
			if !found {
				t.Errorf("%s missing %s in %s", id, expected.name, expected.where)
			}
		}
	}
}
