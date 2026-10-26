package agcapi

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/Createitv/agc-cli/pkg/domain"
)

func TestProductImportFeedbackIdentifiesNestedInvalidInputs(t *testing.T) {
	endpoint, _ := domain.EndpointByID("pms", "batchaddproduct-harmonyosnext")
	base := `{"productNo":"fixture","productName":"fixture","productDesc":"fixture","purchaseType":"consumable","defaultLocale":"zh_CN",`
	cases := []struct{ name, extra, field string }{
		{"language-length", `"languages":[{"locale":"en_US","productName":"` + strings.Repeat("a", 56) + `","productDesc":"fixture"}]}`, "body.products[0].languages[0].productName"},
		{"price-required", `"defaultPriceInfo":{"country":"CN"}}`, "body.products[0].defaultPriceInfo.currency"},
		{"price-object", `"defaultPriceInfo":"invalid"}`, "body.products[0].defaultPriceInfo"},
		{"language-array", `"languages":"invalid"}`, "body.products[0].languages"},
		{"language-item", `"languages":["invalid"]}`, "body.products[0].languages[0]"},
		{"integer", `"autoAddSaleCountry":"1"}`, "body.products[0].autoAddSaleCountry"},
		{"integer-enum", `"autoAddSaleCountry":2}`, "body.products[0].autoAddSaleCountry"},
		{"required-text", `"reviewNote":"fixture","productName":""}`, "body.products[0].productName"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateEndpointRequest(endpoint, nil, map[string]string{"appId": "fixture"}, []byte(`{"requestId":"fixture","products":[`+base+tc.extra+`]}`))
			var feedback *RequestValidationError
			if !errors.As(err, &feedback) {
				t.Fatalf("expected structured invalid-input error, got %v", err)
			}
			for _, issue := range feedback.Violations {
				if issue.Field == tc.field {
					if issue.SourceURL == "" || issue.Suggestion == "" {
						t.Fatal("missing actionable source guidance")
					}
					return
				}
			}
			t.Fatalf("missing field %s: %+v", tc.field, feedback.Violations)
		})
	}
}

func TestProductImportValidNestedModelsAreAccepted(t *testing.T) {
	endpoint, _ := domain.EndpointByID("pms", "batchaddproduct-harmonyosnext")
	body := []byte(`{"requestId":"fixture","products":[{"productNo":"fixture","productName":"fixture","productDesc":"fixture","purchaseType":"auto_subscription","subGroupId":"fixture-group","subPeriod":"1","subPeriodUnit":"M","defaultLocale":"zh_CN","defaultPriceInfo":{"country":"CN","currency":"CNY","price":"199"},"languages":[{"locale":"en_US","productName":"fixture","productDesc":"fixture"}],"autoAddSaleCountry":0}]}`)
	if err := ValidateEndpointRequest(endpoint, nil, map[string]string{"appId": "fixture"}, body); err != nil {
		t.Fatal(err)
	}
}

func TestReportDateAndFilterRulesUseAllQueryValues(t *testing.T) {
	endpoint, _ := domain.EndpointByID("reports", "appdownloadexport")
	for _, tc := range []struct {
		name, start, end   string
		conditions, values []string
		valid              bool
	}{
		{"valid", "20260901", "20261001", []string{"countryId", "businessType"}, []string{"CN", "AppGallery"}, true},
		{"invalid-calendar-date", "20260230", "20261001", nil, nil, false},
		{"reversed-date", "20261001", "20260901", nil, nil, false},
		{"too-long", "20260101", "20261001", nil, nil, false},
		{"misaligned-filters", "20260901", "20261001", []string{"countryId", "businessType"}, []string{"CN"}, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			q := map[string]string{"language": "en-US", "startTime": tc.start, "endTime": tc.end}
			lists := map[string][]string{}
			if tc.conditions != nil {
				lists["filterCondition"] = tc.conditions
			}
			if tc.values != nil {
				lists["filterConditionValue"] = tc.values
			}
			err := ValidateEndpointRequest(endpoint, q, nil, nil, lists)
			if tc.valid {
				if err != nil {
					t.Fatal(err)
				}
				return
			}
			var detail *RequestValidationError
			if !errors.As(err, &detail) {
				t.Fatalf("expected documented feedback, got %v", err)
			}
			for _, issue := range detail.Violations {
				if issue.SourceURL == "" || issue.Suggestion == "" {
					t.Fatal("missing correction guidance")
				}
			}
		})
	}
}

func TestDomainNestedInputsFollowOfficialModels(t *testing.T) {
	for _, tc := range []struct{ id, body, field string }{
		{"domain-api-precheck", `{"domain":{"type":"httpRequest","value":"https://example.com"}}`, "body.domain.type"},
		{"domain-api-post-domain", `{"domains":[{"type":"unknown","value":"https://example.com"}]}`, "body.domains[0].type"},
		{"domain-api-post-domain", `{"domains":[{"type":"httpRequest"}]}`, "body.domains[0].value"},
	} {
		t.Run(tc.id+tc.field, func(t *testing.T) {
			e, _ := domain.EndpointByID("domains", tc.id)
			err := ValidateEndpointRequest(e, map[string]string{"category": "server"}, map[string]string{"appId": "fixture"}, []byte(tc.body))
			var detail *RequestValidationError
			if !errors.As(err, &detail) {
				t.Fatalf("expected model feedback, got %v", err)
			}
			for _, v := range detail.Violations {
				if v.Field == tc.field && v.SourceURL != "" {
					return
				}
			}
			t.Fatalf("missing nested source guidance: %+v", detail.Violations)
		})
	}
}

func TestProjectApplicationSelectionAndOAuthTeamRules(t *testing.T) {
	app, _ := domain.EndpointByID("projects", "appbriefinfo")
	if err := ValidateEndpointRequest(app, nil, map[string]string{"userID": "1", "oauth2Token": "offline-fixture"}, nil); err == nil {
		t.Fatal("application query accepted without either documented identity")
	}
	if err := ValidateEndpointRequest(app, map[string]string{"packageNames": "app.example.fixture"}, map[string]string{"userID": "1", "oauth2Token": "offline-fixture"}, nil); err != nil {
		t.Fatal(err)
	}
	listing, _ := domain.EndpointByID("projects", "queryprojectlist")
	if err := ValidateEndpointRequest(listing, nil, nil, nil); err != nil {
		t.Fatalf("Service Account branch incorrectly requires OAuth team: %v", err)
	}
	err := ValidateEndpointRequest(listing, nil, map[string]string{"oauth2Token": "offline-fixture"}, nil)
	var feedback *RequestValidationError
	if !errors.As(err, &feedback) || feedback.Violations[0].Field != "header.teamId" {
		t.Fatalf("missing OAuth team guidance: %v", err)
	}
	if err := ValidateEndpointRequest(listing, nil, map[string]string{"oauth2Token": "offline-fixture", "teamId": "1"}, nil); err != nil {
		t.Fatal(err)
	}
}

func TestProjectFingerprintRejectsDocumentedForbiddenColon(t *testing.T) {
	e, _ := domain.EndpointByID("projects", "addappfingerprint")
	err := ValidateEndpointRequest(e, nil, map[string]string{"oauth2Token": "offline-fixture", "uid": "1", "teamId": "1"}, []byte(`{"addCertFingerprints":["AA:BB"]}`))
	var feedback *RequestValidationError
	if !errors.As(err, &feedback) || feedback.Violations[0].Field != "body.addCertFingerprints[0]" {
		t.Fatalf("missing fingerprint correction: %v", err)
	}
}

func TestProvisioningNestedModelsReturnActionableFeedback(t *testing.T) {
	for _, tc := range []struct{ id, body, field string }{
		{"provision-api-add-device", `{"deviceList":[{"deviceName":"offline","udid":"offline","deviceType":99}]}`, "body.deviceList[0].deviceType"},
		{"provision-api-applyacl", `{"aclPermissionList":[{"name":"ohos.permission.fixture"}]}`, "body.aclPermissionList[0].reason"},
		{"provision-api-applyacl", `{"aclPermissionList":{"name":"ohos.permission.fixture","reason":"offline"}}`, "body.aclPermissionList"},
	} {
		t.Run(tc.id+tc.field, func(t *testing.T) {
			e, _ := domain.EndpointByID("provisioning", tc.id)
			err := ValidateEndpointRequest(e, nil, map[string]string{"appId": "fixture"}, []byte(tc.body))
			var detail *RequestValidationError
			if !errors.As(err, &detail) {
				t.Fatalf("expected structured model error, got %v", err)
			}
			for _, v := range detail.Violations {
				if v.Field == tc.field && v.SourceURL != "" && v.Suggestion != "" {
					return
				}
			}
			t.Fatalf("missing model field guidance: %+v", detail.Violations)
		})
	}
}

func TestFeedbackDeletionRequiresActualFeedbackIdentifiers(t *testing.T) {
	e, _ := domain.EndpointByID("testing", "test-api-delete-feedback")
	err := ValidateEndpointRequest(e, nil, map[string]string{"appId": "fixture"}, []byte(`{"delFeedBackInfos":[{"problemId":"fixture"}],"deleteMode":1}`))
	var feedback *RequestValidationError
	if !errors.As(err, &feedback) {
		t.Fatalf("missing feedback model validation: %v", err)
	}
	if feedback.Violations[0].Field != "body.delFeedBackInfos[0].testerId" || feedback.Violations[0].SourceURL == "" {
		t.Fatalf("missing indexed model guidance: %+v", feedback.Violations)
	}
	if err := ValidateEndpointRequest(e, nil, map[string]string{"appId": "fixture"}, []byte(`{"delFeedBackInfos":[{"problemId":"fixture","testerId":"fixture"}],"deleteMode":1}`)); err != nil {
		t.Fatal(err)
	}
}

func TestResourceVersionFieldsRemainNested(t *testing.T) {
	e, _ := domain.EndpointByID("resources", "respackapi-addversion")
	err := ValidateEndpointRequest(e, nil, map[string]string{"requestId": "offline-fixture"}, []byte(`{"resourceVersion":{"appId":"fixture","startTime":1,"endTime":2}}`))
	var detail *RequestValidationError
	if !errors.As(err, &detail) {
		t.Fatalf("missing nested version validation: %v", err)
	}
	for _, v := range detail.Violations {
		if v.Field == "body.resourceVersion.version" && v.SourceURL != "" {
			return
		}
	}
	t.Fatalf("missing root-model guidance: %+v", detail.Violations)
}

func TestCommentReplyTargetsAreMutuallyExclusive(t *testing.T) {
	e, _ := domain.EndpointByID("comments", "com-reviewreply")
	err := ValidateEndpointRequest(e, nil, map[string]string{"requestId": "offline-fixture"}, []byte(`{"lang":"zh_CN","appId":"fixture","reviewId":"fixture","devReplyContent":"offline","countryCode":"CN","toReplyId":"user-reply","updateReplyId":"developer-reply"}`))
	var detail *RequestValidationError
	if !errors.As(err, &detail) {
		t.Fatalf("conflicting reply targets accepted: %v", err)
	}
	for _, v := range detail.Violations {
		if v.Code == "reply-target-exclusive" && v.SourceURL != "" && v.Suggestion != "" {
			return
		}
	}
	t.Fatalf("missing documented target feedback: %+v", detail.Violations)
}

func TestGameResourceJSONStringValidatesItsInnerModel(t *testing.T) {
	e, _ := domain.EndpointByID("gameplay", "resources")
	for _, tc := range []struct{ name, resources, field string }{
		{"invalid-json", "not-json", "body.resources"},
		{"missing-upload-url", `[{"index":"offline","name":"offline","type":1,"operType":0,"updateTime":1}]`, "body.resources[0].url"},
		{"invalid-resource-type", `[{"index":"offline","name":"offline","type":99,"operType":1}]`, "body.resources[0].type"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			body, _ := json.Marshal(map[string]any{"method": "external.hms.gs.syncGameResources", "appId": "fixture", "resources": tc.resources, "ts": 1, "sign": "offline-sign"})
			err := ValidateEndpointRequest(e, nil, nil, body)
			var feedback *RequestValidationError
			if !errors.As(err, &feedback) {
				t.Fatalf("invalid encoded resource accepted: %v", err)
			}
			for _, v := range feedback.Violations {
				if v.Field == tc.field && v.SourceURL != "" {
					return
				}
			}
			t.Fatalf("missing resource model field guidance: %+v", feedback.Violations)
		})
	}
}

func TestLegacySubmitPhaseFieldsAreConditional(t *testing.T) {
	e, _ := domain.EndpointByID("publishing", "app-submit")
	if err := ValidateEndpointRequest(e, map[string]string{"appId": "fixture"}, nil, nil); err != nil {
		t.Fatalf("commercial submit incorrectly requires phase Body: %v", err)
	}
	err := ValidateEndpointRequest(e, map[string]string{"appId": "fixture", "releaseType": "3"}, nil, nil)
	var detail *RequestValidationError
	if !errors.As(err, &detail) {
		t.Fatal("phased submit accepted without its required fields")
	}
	found := false
	for _, v := range detail.Violations {
		found = found || v.Field == "body.phasedReleaseStartTime"
	}
	if !found {
		t.Fatal("missing phase start correction")
	}
}

func TestPermissionVideoRequiresReviewedObjectIdentifier(t *testing.T) {
	e, _ := domain.EndpointByID("publishing", "app-file-info-v3")
	err := ValidateEndpointRequest(e, map[string]string{"appId": "fixture"}, nil, []byte(`{"packagePermissionIntroVideoList":[{"lang":"zh_CN","permissionName":"ohos.permission.fixture","deviceType":4}]}`))
	var detail *RequestValidationError
	if !errors.As(err, &detail) {
		t.Fatalf("permission video accepted without objectId: %v", err)
	}
	for _, v := range detail.Violations {
		if v.Field == "body.packagePermissionIntroVideoList[0].objectId" && v.SourceURL != "" {
			return
		}
	}
	t.Fatalf("missing documented model field: %+v", detail.Violations)
}

func TestTestingLanguageRequiresDescriptionAndTestFeatures(t *testing.T) {
	e, _ := domain.EndpointByID("testing", "test-api-modify-test-version")
	err := ValidateEndpointRequest(e, map[string]string{"appId": "fixture"}, nil, []byte(`{"versionId":"fixture","languages":[{"language":"zh_CN"}]}`))
	var detail *RequestValidationError
	if !errors.As(err, &detail) {
		t.Fatalf("incomplete test language accepted: %v", err)
	}
	for _, field := range []string{"body.languages[0].appDesc", "body.languages[0].newFeatures"} {
		found := false
		for _, v := range detail.Violations {
			found = found || v.Field == field
		}
		if !found {
			t.Errorf("missing documented field %s", field)
		}
	}
}

func TestTestingDeviceMaterialReportsDocumentedNestedViolations(t *testing.T) {
	e, _ := domain.EndpointByID("testing", "test-api-modify-test-version")
	err := ValidateEndpointRequest(e, map[string]string{"appId": "fixture"}, nil, []byte(`{"versionId":"fixture","languages":[{"language":"zh_CN","appDesc":"fixture","newFeatures":"fixture","deviceMaterials":[{"deviceType":19,"screenShots":["one","two"],"showType":2}]}]}`))
	var detail *RequestValidationError
	if !errors.As(err, &detail) {
		t.Fatalf("invalid device material accepted: %v", err)
	}
	for _, field := range []string{"body.languages[0].deviceMaterials[0].deviceType", "body.languages[0].deviceMaterials[0].screenShots", "body.languages[0].deviceMaterials[0].showType"} {
		found := false
		for _, v := range detail.Violations {
			found = found || v.Field == field && strings.Contains(v.SourceURL, "devicematerial")
		}
		if !found {
			t.Errorf("missing source-backed violation for %s: %+v", field, detail.Violations)
		}
	}
}

func TestTestingVersionRequiresReviewedDeviceAndTimeFields(t *testing.T) {
	e, _ := domain.EndpointByID("testing", "test-api-modify-test-version")
	for _, tc := range []struct{ body, field string }{
		{`{"versionId":"fixture","deviceTypes":[{}]}`, "body.deviceTypes[0].deviceType"},
		{`{"versionId":"fixture","deviceTypes":[{"deviceType":19}]}`, "body.deviceTypes[0].deviceType"},
		{`{"versionId":"fixture","openTestInfo":{"startTime":1704211200000}}`, "body.openTestInfo.endTime"},
		{`{"versionId":"fixture","openTestInfo":{"startTime":"1704211200000","endTime":1704211200000}}`, "body.openTestInfo.startTime"},
	} {
		err := ValidateEndpointRequest(e, map[string]string{"appId": "fixture"}, nil, []byte(tc.body))
		var detail *RequestValidationError
		if !errors.As(err, &detail) {
			t.Fatalf("missing nested validation for %s: %v", tc.field, err)
		}
		found := false
		for _, v := range detail.Violations {
			found = found || v.Field == tc.field && v.SourceURL != ""
		}
		if !found {
			t.Errorf("missing field correction %s: %+v", tc.field, detail.Violations)
		}
	}
}

func TestTestingPermissionVideoUsesTestingSpecificDeviceEnums(t *testing.T) {
	e, _ := domain.EndpointByID("testing", "test-api-modify-test-version")
	for _, tc := range []struct{ video, field string }{
		{`{"lang":"zh_CN","permissionName":"fixture","deviceType":4}`, "objectId"},
		{`{"lang":"zh_CN","permissionName":"fixture","deviceType":19,"objectId":"fixture"}`, "deviceType"},
	} {
		err := ValidateEndpointRequest(e, map[string]string{"appId": "fixture"}, nil, []byte(`{"versionId":"fixture","packagePermissionIntroVideoList":[`+tc.video+`]}`))
		var detail *RequestValidationError
		if !errors.As(err, &detail) {
			t.Fatalf("invalid testing video accepted: %v", err)
		}
		found := false
		for _, v := range detail.Violations {
			found = found || v.Field == "body.packagePermissionIntroVideoList[0]."+tc.field && strings.Contains(v.SourceURL, "test-api-data-packagepermission")
		}
		if !found {
			t.Errorf("missing testing-specific violation: %+v", detail.Violations)
		}
	}
	if err := ValidateEndpointRequest(e, map[string]string{"appId": "fixture"}, nil, []byte(`{"versionId":"fixture","packagePermissionIntroVideoList":[{"lang":"zh_CN","permissionName":"fixture","deviceType":14,"objectId":"fixture"}]}`)); err != nil {
		t.Fatal(err)
	}
}

func TestTestingVersionExtensionValidatesReviewedFields(t *testing.T) {
	e, _ := domain.EndpointByID("testing", "test-api-modify-test-version")
	if err := ValidateEndpointRequest(e, map[string]string{"appId": "fixture"}, nil, []byte(`{"versionId":"fixture","extInfo":{"isAiGenerate":1,"aiTypes":[-1,1,2,3,4,5]}}`)); err != nil {
		t.Fatalf("documented AI service types rejected: %v", err)
	}
	for _, tc := range []struct{ ext, field string }{
		{`{"isAiGenerate":1,"aiTypes":[99]}`, "body.extInfo.aiTypes[0]"},
		{`{"encrypted":2}`, "body.extInfo.encrypted"},
		{`{"registeredIdType":0}`, "body.extInfo.registeredIdType"},
		{`{"isAiGenerate":1}`, "body.extInfo.aiTypes"},
		{`{"userAgreementIds":["1","2","3","4","5","6"]}`, "body.extInfo.userAgreementIds"},
		{`{"appReviewPhoneInfo":{"account":1}}`, "body.extInfo.appReviewPhoneInfo.account"},
	} {
		err := ValidateEndpointRequest(e, map[string]string{"appId": "fixture"}, nil, []byte(`{"versionId":"fixture","extInfo":`+tc.ext+`}`))
		var detail *RequestValidationError
		if !errors.As(err, &detail) {
			t.Fatalf("invalid extension accepted for %s: %v", tc.field, err)
		}
		found := false
		for _, v := range detail.Violations {
			found = found || v.Field == tc.field && v.SourceURL != ""
		}
		if !found {
			t.Errorf("missing correction for %s: %+v", tc.field, detail.Violations)
		}
	}
}

func TestTestingResolvedIssueRequiresOfficialIdentifiers(t *testing.T) {
	e, _ := domain.EndpointByID("testing", "test-api-modify-test-version")
	for _, field := range []string{"versionId", "packageId", "vocId"} {
		issue := map[string]any{"versionId": "fixture", "packageId": "fixture", "vocId": "fixture"}
		delete(issue, field)
		body, _ := json.Marshal(map[string]any{"versionId": "fixture", "versionVocReq": []any{issue}})
		err := ValidateEndpointRequest(e, map[string]string{"appId": "fixture"}, nil, body)
		var detail *RequestValidationError
		if !errors.As(err, &detail) {
			t.Fatalf("issue without %s accepted: %v", field, err)
		}
		found := false
		for _, v := range detail.Violations {
			found = found || v.Field == "body.versionVocReq[0]."+field && strings.Contains(v.SourceURL, "versionvocreq")
		}
		if !found {
			t.Errorf("missing documented identifier correction: %+v", detail.Violations)
		}
	}
}

func TestTestingTaskNestedGroupBindingUsesDocumentedEnum(t *testing.T) {
	e, _ := domain.EndpointByID("testing", "test-api-modify-test-version")
	if err := ValidateEndpointRequest(e, map[string]string{"appId": "fixture"}, nil, []byte(`{"versionId":"fixture","openTestInfo":{"startTime":1704211200000,"endTime":1704211200000,"testTaskInfo":{"groupInfos":[{"bind":0},{"bind":1}],"needShareLink":1,"needNotify":0}}}`)); err != nil {
		t.Fatalf("documented task values rejected or optional group identifiers made mandatory: %v", err)
	}
	err := ValidateEndpointRequest(e, map[string]string{"appId": "fixture"}, nil, []byte(`{"versionId":"fixture","openTestInfo":{"startTime":1704211200000,"endTime":1704211200000,"testTaskInfo":{"groupInfos":[{"groupId":"fixture","bind":2}],"needShareLink":2}}}`))
	var detail *RequestValidationError
	if !errors.As(err, &detail) {
		t.Fatalf("invalid nested testing task accepted: %v", err)
	}
	for _, field := range []string{"body.openTestInfo.testTaskInfo.groupInfos[0].bind", "body.openTestInfo.testTaskInfo.needShareLink"} {
		found := false
		for _, v := range detail.Violations {
			found = found || v.Field == field && v.SourceURL != ""
		}
		if !found {
			t.Errorf("missing documented correction for %s: %+v", field, detail.Violations)
		}
	}
}

func TestOpenTestLinkOAuthRequiresSelectedTeam(t *testing.T) {
	e, _ := domain.EndpointByID("testing", "test-api-create-open-test-link")
	for _, headers := range []map[string]string{
		{"Authorization": "Bearer offline-test-token"},
		{"oauth2Token": "offline-test-token", "teamId": "selected-team"},
	} {
		if err := ValidateEndpointRequest(e, nil, headers, []byte(`{"groupId":"fixture","limit":1}`)); err != nil {
			t.Fatalf("documented authentication branch rejected: %v", err)
		}
	}
	err := ValidateEndpointRequest(e, nil, map[string]string{"oauth2Token": "offline-test-token"}, []byte(`{"groupId":"fixture","limit":1}`))
	var detail *RequestValidationError
	if !errors.As(err, &detail) {
		t.Fatalf("OAuth public link accepted without team: %v", err)
	}
	for _, v := range detail.Violations {
		if v.Field == "header.teamId" && v.SourceURL != "" {
			return
		}
	}
	t.Fatalf("missing selected-team correction: %+v", detail.Violations)
}

func TestTestingOptionalNameRejectsExplicitlyEmptyText(t *testing.T) {
	e, _ := domain.EndpointByID("testing", "test-api-modify-test-version")
	for _, name := range []string{"", "a"} {
		body, _ := json.Marshal(map[string]any{"versionId": "fixture", "languages": []any{map[string]any{"language": "zh_CN", "appName": name, "appDesc": "fixture", "newFeatures": "fixture"}}})
		err := ValidateEndpointRequest(e, map[string]string{"appId": "fixture"}, nil, body)
		var detail *RequestValidationError
		if !errors.As(err, &detail) {
			t.Fatalf("name below official minimum accepted: %v", err)
		}
		found := false
		for _, v := range detail.Violations {
			found = found || v.Code == "min_length" && v.Field == "body.languages[0].appName" && v.SourceURL != ""
		}
		if !found {
			t.Errorf("missing minimum-length correction: %+v", detail.Violations)
		}
	}
	if err := ValidateEndpointRequest(e, map[string]string{"appId": "fixture"}, nil, []byte(`{"versionId":"fixture","languages":[{"language":"zh_CN","appDesc":"fixture","newFeatures":"fixture"}]}`)); err != nil {
		t.Fatalf("omitted optional name rejected: %v", err)
	}
}

func TestResourceUploadRejectsOfficiallyForbiddenFileNameCharacters(t *testing.T) {
	e, _ := domain.EndpointByID("resources", "respackapi-uploadapplyfor")
	for _, name := range []string{"资源包.apk", "resource-name_1.apk"} {
		body, _ := json.Marshal(map[string]any{"appId": "fixture", "versionId": 1, "fileName": name, "fileSHA256": "fixture", "fileLength": 6000000, "fileType": 0})
		if err := ValidateEndpointRequest(e, nil, map[string]string{"requestId": "fixture"}, body); err != nil {
			t.Fatalf("valid resource name rejected: %v", err)
		}
	}
	for _, character := range []string{"|", ":", "/", "\\", "*", "?", "\"", "<", ">"} {
		body, _ := json.Marshal(map[string]any{"appId": "fixture", "versionId": 1, "fileName": "resource" + character + "file", "fileSHA256": "fixture", "fileLength": 6000000, "fileType": 0})
		err := ValidateEndpointRequest(e, nil, map[string]string{"requestId": "fixture"}, body)
		var detail *RequestValidationError
		if !errors.As(err, &detail) {
			t.Fatalf("forbidden filename character %q accepted: %v", character, err)
		}
		found := false
		for _, v := range detail.Violations {
			found = found || v.Field == "body.fileName" && v.SourceURL != ""
		}
		if !found {
			t.Errorf("missing filename correction: %+v", detail.Violations)
		}
	}
}

func TestResourceFileQueryEnforcesOfficialCSVLimit(t *testing.T) {
	e, _ := domain.EndpointByID("resources", "respackapi-queryfile")
	for _, ids := range []string{strings.TrimSuffix(strings.Repeat("fixture,", 501), ","), "fixture,,other"} {
		err := ValidateEndpointRequest(e, map[string]string{"appId": "fixture", "fileIds": ids}, nil, nil)
		var detail *RequestValidationError
		if !errors.As(err, &detail) {
			t.Fatalf("invalid file ID list accepted: %v", err)
		}
		found := false
		for _, v := range detail.Violations {
			found = found || v.Field == "query.fileIds" && v.SourceURL != ""
		}
		if !found {
			t.Errorf("missing CSV guidance: %+v", detail.Violations)
		}
	}
	if err := ValidateEndpointRequest(e, map[string]string{"appId": "fixture", "fileIds": strings.TrimSuffix(strings.Repeat("fixture,", 500), ",")}, nil, nil); err != nil {
		t.Fatalf("documented 500-ID boundary rejected: %v", err)
	}
}

func TestResourceFileQuerySelectsOneDocumentedIdentity(t *testing.T) {
	e, _ := domain.EndpointByID("resources", "respackapi-queryfile")
	err := ValidateEndpointRequest(e, map[string]string{"appId": "fixture", "versionId": "1", "fileIds": "fixture"}, nil, nil)
	var detail *RequestValidationError
	if !errors.As(err, &detail) {
		t.Fatalf("both file identities accepted: %v", err)
	}
	for _, v := range detail.Violations {
		if v.Code == "resource-file-identity-exclusive" && v.SourceURL != "" {
			return
		}
	}
	t.Fatalf("missing identity selection guidance: %+v", detail.Violations)
}

func TestProjectApplicationQueryBoundsPackageListWithAppIDPrecedence(t *testing.T) {
	e, _ := domain.EndpointByID("projects", "appbriefinfo")
	headers := map[string]string{"oauth2Token": "offline-token", "userID": "1"}
	list := strings.TrimSuffix(strings.Repeat("com.example.fixture,", 101), ",")
	err := ValidateEndpointRequest(e, map[string]string{"packageNames": list}, headers, nil)
	var detail *RequestValidationError
	if !errors.As(err, &detail) {
		t.Fatalf("overlong package list accepted: %v", err)
	}
	found := false
	for _, v := range detail.Violations {
		found = found || v.Field == "query.packageNames" && v.SourceURL != ""
	}
	if !found {
		t.Fatalf("missing package-list guidance: %+v", detail.Violations)
	}
	for _, query := range []map[string]string{{"packageNames": strings.TrimSuffix(strings.Repeat("com.example.fixture,", 100), ",")}, {"appId": "fixture", "packageNames": list}} {
		if err := ValidateEndpointRequest(e, query, headers, nil); err != nil {
			t.Fatalf("documented boundary or appId precedence rejected: %v", err)
		}
	}
}

func TestProjectApplicationQueryUsesOfficialAppTypes(t *testing.T) {
	e, _ := domain.EndpointByID("projects", "appbriefinfo")
	headers := map[string]string{"oauth2Token": "offline-token", "userID": "1"}
	err := ValidateEndpointRequest(e, map[string]string{"appId": "fixture", "appType": "999"}, headers, nil)
	var detail *RequestValidationError
	if !errors.As(err, &detail) {
		t.Fatalf("undefined app type accepted: %v", err)
	}
	found := false
	for _, v := range detail.Violations {
		found = found || v.Field == "query.appType" && v.Code == "enum" && v.SourceURL != ""
	}
	if !found {
		t.Fatalf("missing official enum correction: %+v", detail.Violations)
	}
	for _, typ := range []string{"1", "2", "4", "8", "9", "19", "23", "27", "28", "29", "30", "9998", "9999"} {
		if err := ValidateEndpointRequest(e, map[string]string{"appId": "fixture", "appType": typ}, headers, nil); err != nil {
			t.Fatalf("official appType %s rejected: %v", typ, err)
		}
	}
}

func TestProjectApplicationQueryValidatesEachDistributionStatus(t *testing.T) {
	e,_:=domain.EndpointByID("projects","appbriefinfo")
	headers:=map[string]string{"oauth2Token":"offline-token","userID":"1"}
	err:=ValidateEndpointRequest(e,map[string]string{"appId":"fixture","appDistStatus":"101,999"},headers,nil)
	var detail *RequestValidationError
	if !errors.As(err,&detail) {t.Fatalf("undefined distribution status accepted: %v",err)}
	found:=false;for _,v:=range detail.Violations {found=found || v.Field=="query.appDistStatus" && v.SourceURL!=""}
	if !found {t.Fatalf("missing distribution-list correction: %+v",detail.Violations)}
	if err:=ValidateEndpointRequest(e,map[string]string{"appId":"fixture","appDistStatus":"-999,-995,0,1,62,101,99999,9998,9997"},headers,nil);err!=nil {t.Fatalf("official status list rejected: %v",err)}
}
