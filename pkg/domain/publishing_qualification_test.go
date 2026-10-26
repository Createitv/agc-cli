package domain

import "testing"

func TestQualificationOpenAPINestedRequiredness(t *testing.T) {
	for _, id := range []string{"qualification-review-create", "qualification-review-update"} {
		e, _ := EndpointByID("publishing", id)
		body := invokeRequestSchemaProperties(e)["body"].(map[string]any)
		tags := body["properties"].(map[string]any)["tagQualInfos"].(map[string]any)["items"].(map[string]any)
		required := map[string]bool{}
		for _, name := range tags["required"].([]string) {
			required[name] = true
		}
		if !required["tagId"] || required["taskId"] != (id == "qualification-review-update") {
			t.Fatalf("%s wrong create/update requirements: %v", id, required)
		}
		qualifications := tags["properties"].(map[string]any)["qualificationInfos"].(map[string]any)
		item := qualifications["items"].(map[string]any)
		props := item["properties"].(map[string]any)
		if props["expireTime"].(map[string]any)["format"] != "date" {
			t.Fatal("expiry format missing")
		}
		files := props["qualificationFileInfos"].(map[string]any)["items"].(map[string]any)
		if len(files["required"].([]string)) != 4 {
			t.Fatal("file model requirements missing")
		}
		size := files["properties"].(map[string]any)["fileSize"].(map[string]any)
		if size["minimum"] != nil || size["maximum"] != nil {
			t.Fatal("integer width must not become an invented size range")
		}
		rules := item["x-agc-rules"].([]ReferenceRule)
		if len(rules) != 1 || rules[0].MinDays != 30 || rules[0].When.Equals != "0" {
			t.Fatal("custom-validity branch rule missing")
		}
	}
}
