package domain

import (
	"reflect"
	"testing"
)

func TestLegacyAppInfoOpenAPIIncludesDocumentedSignedEnumValues(t *testing.T) {
	op := OpenAPISpec()["paths"].(map[string]any)["/api/v1/publishing/endpoints/app-info-update/invoke"].(map[string]any)["post"].(map[string]any)
	schema := op["requestBody"].(map[string]any)["content"].(map[string]any)["application/json"].(map[string]any)["schema"].(map[string]any)
	body := schema["properties"].(map[string]any)["body"].(map[string]any)["properties"].(map[string]any)
	for name, want := range map[string][]int64{"isFree": {0, 1}, "registeredDclType": {-1, 0, 1, 2}, "registeredIdType": {-1, 1, 2, 3}} {
		p := body[name].(map[string]any)
		if !reflect.DeepEqual(p["enum"], want) {
			t.Fatalf("%s enum=%#v, want %#v", name, p["enum"], want)
		}
	}
}
