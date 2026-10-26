package domain

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestPublishingDiscoveryExplainsPurpose(t *testing.T) {
	for _, tc := range []struct{ id, purpose string }{
		{"app-create", "HarmonyOS应用/元服务"},
		{"qualification-review-create", "元服务"},
		{"app-package-info-update", "软件包"},
	} {
		e, _ := EndpointByID("publishing", tc.id)
		if e.Description == e.Name || !strings.Contains(e.Description, tc.purpose) {
			t.Errorf("%s lacks official purpose: %s", tc.id, e.Description)
		}
	}
}

func TestPublishingAgreementDocsIncludeLinkedModels(t *testing.T) {
	for _, tc := range []struct{ id, info, base string }{
		{"user-agreement-create", "AgreementInfo（用户协议）", "AgreementBaseInfo（用户协议）"},
		{"privacy-agreement-update", "AgreementInfo（隐私政策协议）", "AgreementBaseInfo（隐私政策协议）"},
	} {
		raw, _ := EndpointReference("publishing", tc.id)
		var doc struct {
			RelatedDocuments []struct {
				Title  string
				Tables []json.RawMessage
			} `json:"relatedDocuments"`
		}
		if err := json.Unmarshal(raw, &doc); err != nil {
			t.Fatal(err)
		}
		found := map[string]bool{}
		for _, model := range doc.RelatedDocuments {
			found[model.Title] = len(model.Tables) > 0
		}
		for _, name := range []string{tc.info, tc.base} {
			if !found[name] {
				t.Errorf("%s missing full source model %s", tc.id, name)
			}
		}
	}
}

func TestUserAgreementDocsDoNotFollowAdjacentModelNavigation(t *testing.T) {
	raw, _ := EndpointReference("publishing", "user-agreement-create")
	var doc struct {
		RelatedDocuments []struct{ Title string } `json:"relatedDocuments"`
	}
	if err := json.Unmarshal(raw, &doc); err != nil {
		t.Fatal(err)
	}
	for _, model := range doc.RelatedDocuments {
		if strings.Contains(model.Title, "隐私政策") || model.Title == "AdsInfo" {
			t.Fatalf("unrelated adjacent document included: %s", model.Title)
		}
	}
}

func TestPublishingDocsIncludeOfficialErrorTable(t *testing.T) {
	raw, _ := EndpointReference("publishing", "language-info-delete")
	var doc struct {
		RelatedDocuments []struct {
			Title  string
			Tables []struct{ Rows [][]string }
		} `json:"relatedDocuments"`
	}
	if err := json.Unmarshal(raw, &doc); err != nil {
		t.Fatal(err)
	}
	for _, model := range doc.RelatedDocuments {
		if model.Title != "错误码" {
			continue
		}
		for _, table := range model.Tables {
			for _, row := range table.Rows {
				if len(row) >= 2 && row[0] == "204144651" && strings.Contains(row[1], "默认语言") {
					return
				}
			}
		}
	}
	t.Fatal("missing official default-language business-error guidance")
}

func TestPrivacyDocsIncludeBothApplicationKindExamples(t *testing.T) {
	for _, id := range []string{"privacy-agreement-update", "privacy-agreement-list"} {
		raw, _ := EndpointReference("publishing", id)
		var doc struct {
			Examples []struct{ Section, Text string }
		}
		if err := json.Unmarshal(raw, &doc); err != nil {
			t.Fatal(err)
		}
		section := "请求示例"
		if id == "privacy-agreement-list" {
			section = "响应示例"
		}
		found := false
		for _, example := range doc.Examples {
			if example.Section == section && strings.Contains(example.Text, "HarmonyOS应用") && strings.Contains(example.Text, "元服务") && strings.Contains(example.Text, "agreementInfos") {
				found = true
			}
		}
		if !found {
			t.Fatalf("%s lost the application-kind example bodies", id)
		}
	}
}

func TestPublishingAgreementOpenAPINestedModel(t *testing.T) {
	paths := OpenAPISpec()["paths"].(map[string]any)
	for _, id := range []string{"user-agreement-create", "user-agreement-update", "privacy-agreement-update"} {
		op := paths["/api/v1/publishing/endpoints/"+id+"/invoke"].(map[string]any)["post"].(map[string]any)
		schema := op["requestBody"].(map[string]any)["content"].(map[string]any)["application/json"].(map[string]any)["schema"].(map[string]any)
		body := schema["properties"].(map[string]any)["body"].(map[string]any)
		list := body["properties"].(map[string]any)["agreementInfos"].(map[string]any)
		item := list["items"].(map[string]any)
		props := item["properties"].(map[string]any)
		if props["lang"].(map[string]any)["maxLength"] != 20 {
			t.Fatalf("%s lost language bound", id)
		}
		base := props["baseInfo"].(map[string]any)
		if base["properties"].(map[string]any)["name"].(map[string]any)["maxLength"] != 100 {
			t.Fatalf("%s lost name bound", id)
		}
		if len(item["required"].([]string)) != 1 || item["required"].([]string)[0] != "lang" {
			t.Fatalf("%s must require only lang", id)
		}
		if len(base["required"].([]string)) != 1 || base["required"].([]string)[0] != "name" {
			t.Fatalf("%s missing conditional object requirement", id)
		}
	}
}
