package domain

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestReviewedManagementDiscovery(t *testing.T) {
	for family, operations := range map[string][]string{
		"cache":     {"execute-command", "connection-info", "usage"},
		"clouddb":   {"schema-versions", "get-object-type", "create-object-type", "update-object-type", "import-object-types", "list-zones", "create-zone", "delete-zone"},
		"functions": {"create-alias", "get-alias", "list-aliases", "update-alias", "delete-alias", "create-http-trigger", "create-cron-trigger", "get-http-trigger", "get-trigger", "delete-http-trigger", "delete-trigger", "update-cron-trigger", "list-layers", "list-layer-versions", "get-layer-version", "create-layer-version"},
		"templates": {"search-applications", "create-application", "create-plan", "get-plan", "apply-plan"},
		"prefetch":  {"fetch-install", "fetch-periodic", "fetch-deferred"},
	} {
		for _, id := range operations {
			t.Run(family+"/"+id, func(t *testing.T) {
				e, ok := EndpointByID(family, id)
				if !ok {
					t.Fatalf("missing reviewed operation %s/%s", family, id)
				}
				if e.Method == "" || e.SourceURL == "" {
					t.Fatal("missing method/source contract")
				}
				if _, ok := EndpointReference(family, id); !ok {
					t.Fatal("missing embedded evidence")
				}
			})
		}
	}
}

func TestPrefetchDocsRetainFullLinkedConfigurationGuide(t *testing.T) {
	raw, ok := EndpointReference("prefetch", "fetch-periodic")
	if !ok {
		t.Fatal("docs absent")
	}
	var doc struct {
		RelatedDocuments []struct{ SourceURL, RenderedText string }
	}
	if json.Unmarshal(raw, &doc) != nil {
		t.Fatal("invalid document")
	}
	for _, related := range doc.RelatedDocuments {
		if strings.HasSuffix(related.SourceURL, "cloudfoundation-prefetch-config") && strings.Contains(related.RenderedText, "仅支持填写一个服务器地址") {
			return
		}
	}
	t.Fatal("full source configuration guide was replaced by unrelated supplements")
}

func TestAliasMapOpenAPIContainsDynamicWeightContract(t *testing.T) {
	e, _ := EndpointByID("functions", "create-alias")
	props := invokeRequestSchemaProperties(e)
	body := props["body"].(map[string]any)["properties"].(map[string]any)
	routing := body["routingConfig"].(map[string]any)
	values, ok := routing["additionalProperties"].(map[string]any)
	if !ok || values["type"] != "integer" || values["minimum"] != int64(0) || values["maximum"] != int64(100) || routing["maxProperties"] != 2 {
		t.Fatalf("missing dynamic version-weight shape: %+v", routing)
	}
}
