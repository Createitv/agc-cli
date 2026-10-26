package domain

import "testing"

func TestCloudFunctionsDiscoveryAndNativeContract(t *testing.T) {
	c, ok := CapabilityByID("functions")
	if !ok || c.EndpointCount != 29 {
		t.Fatal("functions family with 9 native workflows and 20 management transports missing")
	}
	for _, id := range []string{"create-function", "update-function-code", "update-function-config", "publish-function", "update-and-publish-function", "invoke", "list-functions", "query-function-version", "list-function-versions"} {
		e, ok := EndpointByID("functions", id)
		if !ok || e.Authentication != "agcli-user-oauth" || e.Source != "Huawei official AGCLI 1.12.0 source" {
			t.Fatalf("missing source/auth contract for %s: %+v", id, e)
		}
		if _, ok := EndpointReference("functions", id); !ok {
			t.Fatalf("source evidence missing for %s", id)
		}
	}
	spec := OpenAPISpec()
	paths := spec["paths"].(map[string]any)
	op := paths["/api/v1/functions/endpoints/invoke/invoke"].(map[string]any)["post"].(map[string]any)
	schema := op["requestBody"].(map[string]any)["content"].(map[string]any)["application/json"].(map[string]any)["schema"].(map[string]any)
	props := schema["properties"].(map[string]any)
	if props["payload"] == nil || props["sessionFile"] == nil || props["name"] == nil {
		t.Fatal("native request fields missing in OpenAPI")
	}
}
