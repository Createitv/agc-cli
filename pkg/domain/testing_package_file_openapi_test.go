package domain

import "testing"

func TestTestingPackageFileOpenAPIMatchesOfficialModel(t *testing.T) {
	op := OpenAPISpec()["paths"].(map[string]any)["/api/v1/testing/endpoints/test-api-add-test-package/invoke"].(map[string]any)["post"].(map[string]any)
	schema := op["requestBody"].(map[string]any)["content"].(map[string]any)["application/json"].(map[string]any)["schema"].(map[string]any)
	body := schema["properties"].(map[string]any)["body"].(map[string]any)
	file := body["properties"].(map[string]any)["file"].(map[string]any)
	props := file["properties"].(map[string]any)
	for name, limit := range map[string]int{"fileName": 255, "objectId": 1024} {
		p := props[name].(map[string]any)
		if p["type"] != "string" || p["maxLength"] != limit {
			t.Fatalf("wrong official %s schema: %v", name, p)
		}
	}
	required := file["required"].([]string)
	if len(required) != 2 {
		t.Fatalf("missing FileInfo requirements: %v", required)
	}
}
