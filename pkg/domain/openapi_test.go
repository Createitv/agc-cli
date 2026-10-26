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

func TestOpenAPIExposesVerifiedTestingMemberContract(t *testing.T) {
	ep, _ := EndpointByID("testing", "test-api-query-test-user")
	props := invokeRequestSchemaProperties(ep)
	headers := props["headers"].(map[string]any)
	required := headers["required"].([]string)
	if len(required) != 1 || required[0] != "appId" {
		t.Fatalf("wrong required headers: %v", required)
	}
	query := props["query"].(map[string]any)["properties"].(map[string]any)
	for _, name := range []string{"appId", "groupId"} {
		field, ok := query[name].(map[string]any)
		if !ok || field["x-huawei-required"] != true {
			t.Fatalf("missing required query field %s", name)
		}
	}
	op := endpointOperation("invoke", ep, "Invoke endpoint", nil)
	schema := op["requestBody"].(map[string]any)["content"].(map[string]any)["application/json"].(map[string]any)["schema"].(map[string]any)
	if len(schema["allOf"].([]any)) != 2 {
		t.Fatal("missing scalar/repeated required query alternatives")
	}
}

func TestOpenAPIDescribesConditionalLocalServerAuthentication(t *testing.T) {
	spec := OpenAPISpec()
	schemes := spec["components"].(map[string]any)["securitySchemes"].(map[string]any)
	if _, ok := schemes["bearerAuth"]; ok {
		t.Fatal("local server must not declare Huawei bearer auth")
	}
	scheme, ok := schemes["localServerToken"].(map[string]string)
	if !ok || scheme["type"] != "apiKey" || scheme["in"] != "header" || scheme["name"] != "X-AGC-Server-Token" {
		t.Fatalf("local scheme = %#v", schemes)
	}
	if _, ok := spec["security"]; !ok {
		t.Fatal("missing conditional global security")
	}
	paths := spec["paths"].(map[string]any)
	for path, raw := range paths {
		for _, operation := range raw.(map[string]any) {
			if _, ok := operation.(map[string]any)["security"]; ok {
				t.Fatalf("%s overrides conditional global security", path)
			}
		}
	}
}
