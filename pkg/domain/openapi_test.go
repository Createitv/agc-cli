package domain

import "testing"

func TestOpenAPISpecIncludesEveryEndpointInvokeRoute(t *testing.T) {
	spec := OpenAPISpec()
	paths, ok := spec["paths"].(map[string]any)
	if !ok {
		t.Fatalf("paths has type %T", spec["paths"])
	}
	for _, endpoint := range AllEndpoints() {
		path := "/api/v1/" + endpoint.FamilyID + "/endpoints/" + endpoint.ID + "/invoke"
		raw, ok := paths[path]
		if !ok {
			t.Fatalf("missing OpenAPI invoke path %s", path)
		}
		operations, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("path item has type %T", raw)
		}
		if _, ok := operations["post"]; !ok {
			t.Fatalf("%s missing POST operation", path)
		}
	}
}

func TestOpenAPISpecCarriesHuaweiEndpointMetadata(t *testing.T) {
	spec := OpenAPISpec()
	paths := spec["paths"].(map[string]any)
	raw := paths["/api/v1/publishing/endpoints/app-submit/invoke"].(map[string]any)
	post := raw["post"].(map[string]any)
	if post["x-huawei-method"] != "POST" {
		t.Fatalf("method = %v", post["x-huawei-method"])
	}
	if post["x-huawei-path"] != "/api/publish/v2/app-submit" {
		t.Fatalf("path = %v", post["x-huawei-path"])
	}
	if post["x-agc-command"] != "agc publishing app-submit" {
		t.Fatalf("command = %v", post["x-agc-command"])
	}
}
