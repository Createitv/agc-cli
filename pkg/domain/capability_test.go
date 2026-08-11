package domain

import "testing"

func TestDecoratedCapabilitiesExposeAllOfficialFamilies(t *testing.T) {
	capabilities := DecoratedCapabilities()
	if got, want := len(capabilities), 13; got != want {
		t.Fatalf("len(capabilities) = %d, want %d", got, want)
	}
	for _, capability := range capabilities {
		if capability.ID == "" || capability.Command == "" || capability.RESTPath == "" {
			t.Fatalf("capability is missing public fields: %#v", capability)
		}
		if capability.Links["self"].Href != capability.RESTPath {
			t.Fatalf("self link = %q, want %q", capability.Links["self"].Href, capability.RESTPath)
		}
		if capability.Affordances["list"] == "" {
			t.Fatalf("%s has no list affordance", capability.ID)
		}
	}
}

func TestCapabilityByID(t *testing.T) {
	capability, ok := CapabilityByID("publishing")
	if !ok {
		t.Fatal("publishing capability not found")
	}
	if capability.Name != "Publishing API" {
		t.Fatalf("name = %q", capability.Name)
	}
	if _, ok := CapabilityByID("missing"); ok {
		t.Fatal("missing capability unexpectedly found")
	}
}
