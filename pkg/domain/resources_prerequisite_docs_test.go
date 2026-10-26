package domain

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestResourceDocsExposeDeveloperIDGuide(t *testing.T) {
	for _, endpoint := range EndpointsByFamily("resources") {
		raw, ok := EndpointReference("resources", endpoint.ID)
		if !ok {
			t.Fatalf("missing reference %s", endpoint.ID)
		}
		var doc struct {
			RelatedDocuments []struct {
				SourceURL    string `json:"sourceUrl"`
				RenderedText string `json:"renderedText"`
			} `json:"relatedDocuments"`
		}
		if err := json.Unmarshal(raw, &doc); err != nil {
			t.Fatal(err)
		}
		found := false
		for _, related := range doc.RelatedDocuments {
			if strings.Contains(related.SourceURL, "agc-get-developerid-projectid-0000001166543063") && strings.Contains(related.RenderedText, "Developer ID and Project ID") {
				found = true
			}
		}
		if !found {
			t.Fatalf("developer ID prerequisite guide absent for %s", endpoint.ID)
		}
	}
}
