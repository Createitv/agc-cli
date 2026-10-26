package domain

import "testing"

func TestPMSUploadOpenAPIRequiresFilePathInsteadOfJSONFile(t *testing.T) {
	op := OpenAPISpec()["paths"].(map[string]any)["/api/v1/pms/endpoints/resourcemanage-harmonyosnext/invoke"].(map[string]any)["post"].(map[string]any)
	schema := op["requestBody"].(map[string]any)["content"].(map[string]any)["application/json"].(map[string]any)["schema"].(map[string]any)
	required, _ := schema["required"].([]string)
	found := false
	for _, name := range required {
		found = found || name == "filePath"
	}
	if !found {
		t.Fatal("multipart filePath not required in REST schema")
	}
	fields := schema["properties"].(map[string]any)["fields"].(map[string]any)
	fieldRequired, _ := fields["required"].([]string)
	for _, name := range fieldRequired {
		if name == "file" {
			t.Fatal("file incorrectly required inside JSON fields")
		}
	}
}
