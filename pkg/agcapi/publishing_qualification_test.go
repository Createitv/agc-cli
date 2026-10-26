package agcapi

import (
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/Createitv/agc-cli/pkg/domain"
)

func TestPublishingQualificationModels(t *testing.T) {
	for _, tc := range []struct{ id, body, path string }{
		{"qualification-review-create", `{"tagQualInfos":[{}]}`, "body.tagQualInfos[0].tagId"},
		{"qualification-review-update", `{"tagQualInfos":[{"tagId":"fixture"}]}`, "body.tagQualInfos[0].taskId"},
		{"qualification-review-create", `{"tagQualInfos":[{"tagId":"fixture","qualificationInfos":[{"validityType":1}]}]}`, "body.tagQualInfos[0].qualificationInfos[0].qualificationCode"},
		{"qualification-review-create", `{"tagQualInfos":[{"tagId":"fixture","qualificationInfos":[{"qualificationCode":"10023","validityType":1}]}]}`, "body.tagQualInfos[0].qualificationInfos[0].qualificationFileInfos"},
		{"qualification-review-create", `{"tagQualInfos":[{"tagId":"fixture","qualificationInfos":[{"qualificationCode":"10023","validityType":1,"qualificationFileInfos":[{}]}]}]}`, "body.tagQualInfos[0].qualificationInfos[0].qualificationFileInfos[0].objectId"},
		{"qualification-review-create", `{"tagQualInfos":[{"tagId":"fixture","qualificationInfos":[{"qualificationCode":"10023","validityType":2,"qualificationFileInfos":[]}]}]}`, "body.tagQualInfos[0].qualificationInfos[0].validityType"},
		{"app-info-update-v3", `{"publishCountry":"CN","encrypted":0,"versionVocs":[{"versionId":"fixture","packageId":"fixture"}]}`, "body.versionVocs[0].vocId"},
	} {
		t.Run(tc.id+"/"+tc.path, func(t *testing.T) {
			e, _ := domain.EndpointByID("publishing", tc.id)
			err := ValidateEndpointRequest(e, map[string]string{"appId": "fixture"}, map[string]string{"appId": "fixture"}, []byte(tc.body))
			var feedback *RequestValidationError
			if !errors.As(err, &feedback) {
				t.Fatalf("expected documented nested error, got %v", err)
			}
			for _, v := range feedback.Violations {
				if v.Field == tc.path && v.SourceURL != "" && v.Suggestion != "" {
					return
				}
			}
			t.Fatalf("missing field %s: %+v", tc.path, feedback.Violations)
		})
	}
}

func TestPublishingQualificationDateRules(t *testing.T) {
	e, _ := domain.EndpointByID("publishing", "qualification-review-create")
	today := time.Now().UTC().Truncate(24 * time.Hour)
	for _, tc := range []struct {
		name, date string
		valid      bool
	}{
		{"empty-date", "", false},
		{"calendar-date", "2026-02-30", false},
		{"wrong-format", "20260510", false},
		{"thirty-days", today.AddDate(0, 0, 30).Format("2006-01-02"), false},
		{"thirty-one-days", today.AddDate(0, 0, 31).Format("2006-01-02"), true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			body := fmt.Sprintf(`{"tagQualInfos":[{"tagId":"fixture","qualificationInfos":[{"qualificationCode":"10023","validityType":0,"expireTime":"%s","qualificationFileInfos":[{"objectId":"fixture","fileName":"fixture.pdf","fileSize":10,"sha256":"fixture"}]}]}]}`, tc.date)
			err := ValidateEndpointRequest(e, nil, map[string]string{"appId": "fixture"}, []byte(body))
			if tc.valid {
				if err != nil {
					t.Fatal(err)
				}
				return
			}
			var feedback *RequestValidationError
			if !errors.As(err, &feedback) {
				t.Fatalf("date should fail: %v", err)
			}
			for _, v := range feedback.Violations {
				if strings.HasSuffix(v.Field, ".expireTime") && strings.Contains(v.SourceURL, "data-qualificationinfo") {
					return
				}
			}
			t.Fatalf("missing expiry violation: %+v", feedback.Violations)
		})
	}
	// Dates are optional in the model table; do not turn a permanent validity
	// branch into an unconditional request for startTime/expireTime.
	if err := ValidateEndpointRequest(e, nil, map[string]string{"appId": "fixture"}, []byte(`{"tagQualInfos":[{"tagId":"fixture","qualificationInfos":[{"qualificationCode":"10023","validityType":1,"qualificationFileInfos":[{"objectId":"fixture","fileName":"fixture.pdf","fileSize":10,"sha256":"fixture"}]}]}]}`)); err != nil {
		t.Fatal(err)
	}
}
