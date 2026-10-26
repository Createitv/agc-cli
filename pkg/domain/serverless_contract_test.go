package domain

import (
	"encoding/json"
	"testing"
)

func TestServerlessOfficialHTTPFamilies(t *testing.T) {
	for family, want := range map[string]int{"cloud-auth": 17, "clouddb": 5, "storage": 6, "hosting": 16, "templates": 2} {
		got := EndpointsByFamily(family)
		publicCount := 0
		for _, e := range got {
			if e.Authentication == "agcli-management-cookie" || e.Authentication == "agcli-management-oauth" {
				continue
			}
			publicCount++
			raw, ok := EndpointReference(family, e.ID)
			if !ok || !json.Valid(raw) {
				t.Errorf("%s missing readable source", e.Command)
			}
			if e.Documentation == nil || e.Documentation.Status != "browser-read" {
				t.Errorf("%s missing reviewed provenance", e.Command)
			}
		}
		if publicCount != want {
			t.Errorf("%s public HTTP operations=%d want %d", family, publicCount, want)
		}
	}
}
func TestServerlessSourceContracts(t *testing.T) {
	// Independent assertions from official parameter tables/examples, not registry-derived fixtures.
	e, ok := EndpointByID("hosting", "upload-url")
	if !ok {
		t.Fatal("hosting upload-url missing")
	}
	p := findServerlessParameter(e.Parameters, "materialInfos")
	if p.Type != "array" || p.ItemsType != "object" {
		t.Fatal("materialInfos must be object array")
	}
	f := findServerlessParameter(p.ItemsProperties, "contentType")
	if !f.Required || len(f.Enum) != 1 || f.Enum[0] != "902" {
		t.Fatal("nested contentType=902 required")
	}
	e, _ = EndpointByID("cloud-auth", "export-users")
	p = findServerlessParameter(e.Parameters, "client_id")
	if p.Required {
		t.Fatal("Service Account alternative must not require client_id")
	}
	e, _ = EndpointByID("clouddb", "query-data")
	p = findServerlessParameter(findServerlessParameter(e.Parameters, "queryReqMsg").Properties, "queryCond")
	if p.Type != "string" || p.EncodedJSON == nil {
		t.Fatal("queryCond must retain encoded JSON string")
	}
	e, _ = EndpointByID("templates", "player-id")
	if e.Direction != "inbound-callback" {
		t.Fatal("player-id is a developer callback")
	}
	e, _ = EndpointByID("storage", "update-metadata")
	if e.Method != "GET" {
		t.Fatal("official metadata-update verb is GET; preserve source")
	}
}
func findServerlessParameter(ps []Parameter, name string) Parameter {
	for _, p := range ps {
		if p.Name == name {
			return p
		}
	}
	return Parameter{}
}

func TestServerlessFileAndNativeOpenAPI(t *testing.T) {
	spec := OpenAPISpec()["paths"].(map[string]any)
	for _, op := range []string{"/api/v1/storage/endpoints/upload-file/invoke", "/api/v1/hosting/endpoints/upload-version/invoke"} {
		schema := spec[op].(map[string]any)["post"].(map[string]any)["requestBody"].(map[string]any)["content"].(map[string]any)["application/json"].(map[string]any)["schema"].(map[string]any)
		required := schema["required"].([]string)
		has := false
		for _, n := range required {
			if n == "fields" {
				t.Fatal("binary upload incorrectly requires fields")
			}
			has = has || n == "filePath"
		}
		if !has {
			t.Fatal("binary upload requires explicit filePath")
		}
	}
	schema := FunctionSchema()
	p := schema["properties"].(map[string]any)
	if p["timeout"].(map[string]any)["maximum"] != 1800 {
		t.Fatal("native timeout bound absent from schema")
	}
}
