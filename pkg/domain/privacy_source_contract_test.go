package domain

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestPrivacyModelsMatchCapturedOfficialLinkContracts(t *testing.T) {
	raw, err := referenceDocuments.ReadFile("reference/supplemental/privacy-agreement-info.contract.json")
	if err != nil {
		t.Fatal(err)
	}
	var model struct{ Parameters []Parameter }
	if err := json.Unmarshal(raw, &model); err != nil {
		t.Fatal(err)
	}
	e, _ := EndpointByID("publishing", "privacy-agreement-update")
	var actual []Parameter
	for _, p := range e.Parameters {
		if p.Name == "agreementInfos" {
			actual = p.ItemsProperties
		}
	}
	var compare func(string, []Parameter, []Parameter)
	compare = func(path string, got, want []Parameter) {
		if len(got) != len(want) {
			t.Errorf("%s field count: got %d want %d", path, len(got), len(want))
		}
		for _, expected := range want {
			found := false
			for _, p := range got {
				if p.Name == expected.Name {
					found = true
					if p.Type != expected.Type || p.Required != expected.Required || p.MaxLength != expected.MaxLength || p.ItemsMaxLength != expected.ItemsMaxLength || p.MinItems != expected.MinItems || p.MaxItems != expected.MaxItems || p.ItemsType != expected.ItemsType || !reflect.DeepEqual(p.Enum, expected.Enum) || p.SourceURL != expected.SourceURL {
						t.Errorf("%s.%s differs from official table: got %+v want %+v", path, p.Name, p, expected)
					}
					compare(path+"."+p.Name, p.Properties, expected.Properties)
					compare(path+"."+p.Name+"[]", p.ItemsProperties, expected.ItemsProperties)
				}
			}
			if !found {
				t.Errorf("%s missing official field %s", path, expected.Name)
			}
		}
	}
	compare("agreementInfos[]", actual, model.Parameters)
}
