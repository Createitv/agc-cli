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

func TestOpenAPIExposesVerifiedHeaderContract(t *testing.T) {
	ep, _ := EndpointByID("testing", "test-api-query-test-user")
	props := invokeRequestSchemaProperties(ep)
	headers := props["headers"].(map[string]any)
	known, ok := headers["properties"].(map[string]any)
	if !ok {
		t.Fatal("headers has no documented properties")
	}
	if _, ok := known["groupId"]; !ok {
		t.Fatal("missing groupId header")
	}
	required := headers["required"].([]string)
	if len(required) != 2 || required[0] != "appId" || required[1] != "groupId" {
		t.Fatalf("required headers = %v", required)
	}
	op := endpointOperation("invoke", ep, "Invoke endpoint", nil)
	schema := op["requestBody"].(map[string]any)["content"].(map[string]any)["application/json"].(map[string]any)["schema"].(map[string]any)
	containers := schema["required"].([]string)
	if len(containers) != 1 || containers[0] != "headers" {
		t.Fatalf("required containers = %v", containers)
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
