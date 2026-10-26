package domain

import (
	"reflect"
	"testing"
)

func TestTestingAITypeOpenAPIUsesOfficialIntegerItemEnums(t *testing.T) {
	op := OpenAPISpec()["paths"].(map[string]any)["/api/v1/testing/endpoints/test-api-modify-test-version/invoke"].(map[string]any)["post"].(map[string]any)
	schema := op["requestBody"].(map[string]any)["content"].(map[string]any)["application/json"].(map[string]any)["schema"].(map[string]any)
	body := schema["properties"].(map[string]any)["body"].(map[string]any)
	ext := body["properties"].(map[string]any)["extInfo"].(map[string]any)
	items := ext["properties"].(map[string]any)["aiTypes"].(map[string]any)["items"].(map[string]any)
	if items["type"] != "integer" || !reflect.DeepEqual(items["enum"], []int64{-1, 1, 2, 3, 4, 5}) {
		t.Fatalf("AI service item schema differs from reviewed official model: %#v", items)
	}
}
