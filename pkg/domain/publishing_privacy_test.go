package domain

import (
	"encoding/json"
	"path"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

// Compare every deep field to its independently browser-read parameter row.
// This checks more than registry-generated fixtures: missing fields, incorrect
// types, requiredness and explicit string/array bounds all fail the audit.
func TestPrivacyModelsMatchReviewedSourceTables(t *testing.T) {
	seen := map[string]bool{}
	var visit func([]Parameter)
	visit = func(fields []Parameter) {
		if len(fields) == 0 {
			return
		}
		source := fields[0].SourceURL
		if seen[source] {
			return
		}
		seen[source] = true
		raw, err := referenceDocuments.ReadFile("reference/supplemental/" + path.Base(source) + ".json")
		if err != nil {
			t.Fatal(err)
		}
		var doc struct{ Tables []struct{ Rows [][]string } }
		if err = json.Unmarshal(raw, &doc); err != nil {
			t.Fatal(err)
		}
		rows := map[string][]string{}
		for _, table := range doc.Tables {
			for _, row := range table.Rows {
				if len(row) >= 4 && row[0] != "参数名称" {
					rows[row[0]] = row
				}
			}
		}
		if len(fields) != len(rows) {
			t.Errorf("%s: implemented %d/%d fields", source, len(fields), len(rows))
		}
		for _, field := range fields {
			row, ok := rows[field.Name]
			if !ok {
				t.Errorf("undocumented field %s.%s", source, field.Name)
				continue
			}
			if field.Required != (row[1] == "M") {
				t.Errorf("%s wrong requiredness", field.Name)
			}
			kind := "object"
			switch {
			case strings.HasPrefix(row[2], "String"):
				kind = "string"
			case strings.HasPrefix(row[2], "Integer"):
				kind = "integer"
			case row[2] == "boolean":
				kind = "boolean"
			case strings.HasPrefix(row[2], "List<"):
				kind = "array"
			}
			if field.Type != kind {
				t.Errorf("%s wrong type %s, expected %s", field.Name, field.Type, kind)
			}
			if match := regexp.MustCompile(`^String\((\d+)\)$`).FindStringSubmatch(row[2]); len(match) > 0 {
				limit, _ := strconv.Atoi(match[1])
				if field.MaxLength != limit {
					t.Errorf("%s missing string bound %d", field.Name, limit)
				}
			}
			if match := regexp.MustCompile(`^List<String\((\d+)\)>$`).FindStringSubmatch(row[2]); len(match) > 0 {
				limit, _ := strconv.Atoi(match[1])
				if field.ItemsType != "string" || field.ItemsMaxLength != limit {
					t.Errorf("%s missing string-item bound %d", field.Name, limit)
				}
			}
			if match := regexp.MustCompile(`数组长度不超过(\d+)`).FindStringSubmatch(row[3]); len(match) > 0 {
				limit, _ := strconv.Atoi(match[1])
				if field.MaxItems != limit {
					t.Errorf("%s missing array maximum %d", field.Name, limit)
				}
			}
			if strings.Contains(row[3], "数组长度为1~200") && (field.MinItems != 1 || field.MaxItems != 200) {
				t.Errorf("%s missing explicit array interval", field.Name)
			}
			visit(field.Properties)
			visit(field.ItemsProperties)
		}
	}
	for _, model := range []string{"PersonalInfo", "DevicePermission", "ThirdSdkInfo", "ThirdSharedInfo", "ExtendInfo"} {
		visit(publishingPrivacyFields(model))
	}
	if len(seen) != 28 {
		t.Fatalf("reviewed model coverage %d, want 28", len(seen))
	}
}

func TestPrivacyOpenAPIExposesDeepBoundsAndBoolean(t *testing.T) {
	e, _ := EndpointByID("publishing", "privacy-agreement-update")
	schema := invokeRequestSchemaProperties(e)["body"].(map[string]any)
	props := schema["properties"].(map[string]any)
	item := props["agreementInfos"].(map[string]any)["items"].(map[string]any)["properties"].(map[string]any)
	sdks := item["thirdSdkInfo"].(map[string]any)["properties"].(map[string]any)["sdkInfos"].(map[string]any)
	if sdks["maxItems"] != 50 {
		t.Fatal("SDK limit missing")
	}
	labels := sdks["items"].(map[string]any)["properties"].(map[string]any)["privacyLabels"].(map[string]any)
	if labels["maxItems"] != 200 {
		t.Fatal("SDK labels limit missing")
	}
	data := labels["items"].(map[string]any)["properties"].(map[string]any)["dataItems"].(map[string]any)
	sensitive := data["items"].(map[string]any)["properties"].(map[string]any)["sensitiveData"].(map[string]any)
	if sensitive["type"] != "integer" || len(sensitive["enum"].([]int64)) != 2 {
		t.Fatal("sensitive-data enum lost")
	}
	ext := item["extInfo"].(map[string]any)["properties"].(map[string]any)
	storage := ext["storageInfo"].(map[string]any)["properties"].(map[string]any)
	if storage["isCrossBorder"].(map[string]any)["type"] != "boolean" {
		t.Fatal("cross-border flag must remain boolean")
	}
	contact := ext["contactInfo"].(map[string]any)["properties"].(map[string]any)["phoneNum"].(map[string]any)
	if contact["maxItems"] != 10 || contact["items"].(map[string]any)["maxLength"] != 31 {
		t.Fatal("contact array/item bounds missing")
	}
}
