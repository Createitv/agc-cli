package domain

import (
	"strings"
	"testing"
)

func TestTestingResolvedIssueContractIncludesReviewedModel(t *testing.T) {
	e, _ := EndpointByID("testing", "test-api-modify-test-version")
	for _, p := range e.Parameters {
		if p.Name != "versionVocReq" {
			continue
		}
		if p.MaxItems != 1000 || len(p.ItemsProperties) != 3 {
			t.Fatalf("incomplete issue contract: %+v", p)
		}
		for _, child := range p.ItemsProperties {
			if !child.Required || child.Type != "string" || child.MaxLength != 32 || !strings.Contains(child.SourceURL, "versionvocreq") {
				t.Fatalf("wrong issue identifier contract: %+v", child)
			}
		}
		raw, _ := EndpointReference(e.FamilyID, e.ID)
		if !strings.Contains(string(raw), `"title":"VersionVocReq"`) {
			t.Fatal("Agent docs omit the reviewed issue model")
		}
		return
	}
	t.Fatal("missing versionVocReq")
}
