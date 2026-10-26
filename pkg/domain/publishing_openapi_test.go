package domain

import (
	"encoding/json"
	"os"
	"strings"
	"testing"
)

func TestCurrentPublishingOpenAPIExposesOfficialParameters(t *testing.T) {
	raw, err := os.ReadFile("../../cmd/agc/command/testdata/publishing-current.json")
	if err != nil {
		t.Fatal(err)
	}
	var ref struct {
		Endpoints []struct {
			ID, Method, Path string
			Parameters       []struct {
				Name, In, Type string
				Required       bool
			}
		}
	}
	if err = json.Unmarshal(raw, &ref); err != nil {
		t.Fatal(err)
	}
	paths := OpenAPISpec()["paths"].(map[string]any)
	for _, e := range ref.Endpoints {
		t.Run(e.ID, func(t *testing.T) {
			route := "/api/v1/publishing/endpoints/" + e.ID + "/invoke"
			item, ok := paths[route].(map[string]any)
			if !ok {
				t.Fatalf("missing route %s", route)
			}
			operation := item["post"].(map[string]any)
			if operation["x-huawei-method"] != e.Method || operation["x-huawei-path"] != e.Path {
				t.Fatal("wrong Huawei operation exposed")
			}
			schema := operation["requestBody"].(map[string]any)["content"].(map[string]any)["application/json"].(map[string]any)["schema"].(map[string]any)
			props := schema["properties"].(map[string]any)
			for _, p := range e.Parameters {
				container := map[string]string{"query": "query", "header": "headers", "body": "body"}[p.In]
				object := props[container].(map[string]any)
				fields, ok := object["properties"].(map[string]any)
				if !ok {
					t.Fatalf("missing %s schema", container)
				}
				field, ok := fields[p.Name].(map[string]any)
				if !ok {
					t.Fatalf("missing %s.%s", container, p.Name)
				}
				typ := "object"
				switch {
				case strings.HasPrefix(p.Type, "String"):
					typ = "string"
				case strings.HasPrefix(p.Type, "Integer"):
					typ = "integer"
				case strings.HasPrefix(p.Type, "List<"):
					typ = "array"
				}
				if p.In == "body" {
					if field["type"] != typ {
						t.Errorf("%s wrong type", p.Name)
					}
				} else if field["x-huawei-type"] != typ {
					t.Errorf("%s wrong upstream type", p.Name)
				}
				found := false
				required, _ := object["required"].([]string)
				for _, name := range required {
					if name == p.Name {
						found = true
					}
				}
				if p.In == "query" {
					allOf, _ := schema["allOf"].([]any)
					for _, condition := range allOf {
						branches, _ := condition.(map[string]any)["anyOf"].([]any)
						alternatives := map[string]bool{}
						for _, rawBranch := range branches {
							branch := rawBranch.(map[string]any)
							for _, container := range branch["required"].([]string) {
								nested := branch["properties"].(map[string]any)[container].(map[string]any)
								for _, name := range nested["required"].([]string) {
									if name == p.Name {
										alternatives[container] = true
									}
								}
							}
						}
						if alternatives["query"] && alternatives["queryValues"] {
							found = true
						}
					}
				}
				if found != p.Required {
					t.Errorf("%s wrong requiredness", p.Name)
				}
			}
		})
	}
}
