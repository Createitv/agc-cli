package domain

import (
	"encoding/json"
	"testing"
)

// Agent discovery must include the actual browser-read source and the rule
// that explains why a conditional input is required, rather than a bare URL.
func TestEndpointDiscoveryIncludesBrowserReadReference(t *testing.T) {
	e, _ := EndpointByID("publishing", "app-info-query-v3")
	raw, _ := json.Marshal(e)
	var data map[string]any
	json.Unmarshal(raw, &data)
	doc, ok := data["documentation"].(map[string]any)
	if !ok {
		t.Fatal("Agent discovery is missing embedded official documentation")
	}
	if doc["status"] != "browser-read" {
		t.Fatal("documentation is not identified as browser-read")
	}
	if doc["sourceUrl"] != e.SourceURL {
		t.Fatal("wrong source document attached")
	}
	rules, ok := data["rules"].([]any)
	if !ok || len(rules) == 0 {
		t.Fatal("conditional requirement is missing")
	}
}

func TestFileEndpointIncludesDeviceSpecificRequirements(t *testing.T) {
	raw, ok := EndpointReference("publishing", "app-file-info-v3")
	if !ok {
		t.Fatal("missing file documentation")
	}
	var doc struct {
		RelatedDocuments []struct {
			Title  string
			Tables []json.RawMessage
		} `json:"relatedDocuments"`
	}
	if err := json.Unmarshal(raw, &doc); err != nil {
		t.Fatal(err)
	}
	if len(doc.RelatedDocuments) < 4 || doc.RelatedDocuments[0].Title != "应用文件要求" || len(doc.RelatedDocuments[0].Tables) != 2 {
		t.Fatal("device and application-kind appendix missing from Agent docs")
	}
}

func TestBatchOperationsHaveDocumentedRequestBodies(t *testing.T) {
	for _, platform := range []string{"android", "harmonyosnext"} {
		for _, operation := range []string{"batchactproduct", "batchdeactproduct", "batchaddproduct", "batchupdateproduct", "batchactpromotion", "batchdeactpromotion"} {
			id := operation + "-" + platform
			t.Run(id, func(t *testing.T) {
				e, _ := EndpointByID("pms", id)
				fields := map[string]Parameter{}
				for _, p := range e.Parameters {
					if p.In == "body" {
						fields[p.Name] = p
					}
				}
				if p := fields["requestId"]; !p.Required || p.Type != "string" {
					t.Fatal("missing required request identifier")
				}
				list := "products"
				items := "string"
				if operation == "batchactpromotion" || operation == "batchdeactpromotion" {
					list = "promotions"
				}
				if operation == "batchaddproduct" || operation == "batchupdateproduct" {
					items = "object"
				}
				p := fields[list]
				if !p.Required || p.Type != "array" || p.ItemsType != items {
					t.Fatalf("incorrect %s list contract: %+v", list, p)
				}
				if list == "promotions" && p.MaxItems != 100 {
					t.Fatal("missing documented promotion list maximum")
				}
			})
		}
	}
}

func TestBatchImportDocumentsIncludePlatformSpecificModels(t *testing.T) {
	for _, platform := range []string{"android", "harmonyosnext"} {
		raw, ok := EndpointReference("pms", "batchaddproduct-"+platform)
		if !ok {
			t.Fatal("missing batch documentation")
		}
		var doc struct {
			RelatedDocuments []struct {
				Title     string
				SourceURL string `json:"sourceUrl"`
			} `json:"relatedDocuments"`
		}
		if err := json.Unmarshal(raw, &doc); err != nil {
			t.Fatal(err)
		}
		if len(doc.RelatedDocuments) != 5 {
			t.Fatalf("%s missing linked import/price/language/promotion source documents", platform)
		}
		for _, model := range doc.RelatedDocuments {
			if model.SourceURL == "" || model.Title == "" {
				t.Fatal("model lacks official provenance")
			}
		}
	}
}

func TestProductQueryAndReviewExposeActualTableFields(t *testing.T) {
	for _, platform := range []string{"android", "harmonyosnext"} {
		e, _ := EndpointByID("pms", "bygetproductinfo-"+platform)
		found := map[string]Parameter{}
		for _, p := range e.Parameters {
			found[p.In+"."+p.Name] = p
		}
		for _, name := range []string{"requestId", "productNo", "productName", "purchaseType", "status", "orderBy", "beginTime", "endTime", "pageNum", "pageSize", "country", "groupId"} {
			p, ok := found["body."+name]
			if !ok {
				t.Errorf("%s missing documented query field %s", platform, name)
			}
			if p.SourceURL == "" {
				t.Errorf("%s.%s lacks source paragraph", platform, name)
			}
		}
		if found["body.pageNum"].Type != "integer" {
			t.Error("pagination lost JSON integer type")
		}
	}
	e, _ := EndpointByID("pms", "reviewproduct-harmonyosnext")
	for _, p := range e.Parameters {
		if p.Name == "productIdList" && p.In == "body" && p.Type == "array" && p.ItemsType == "string" && p.Required {
			return
		}
	}
	t.Fatal("missing documented product review list")
}

func TestReportOpenAPIAcceptsScalarOrRepeatedRequiredQueryInputs(t *testing.T) {
	spec := OpenAPISpec()
	path := spec["paths"].(map[string]any)["/api/v1/reports/endpoints/appdownloadexport/invoke"].(map[string]any)
	schema := path["post"].(map[string]any)["requestBody"].(map[string]any)["content"].(map[string]any)["application/json"].(map[string]any)["schema"].(map[string]any)
	for _, name := range schema["required"].([]string) {
		if name == "query" {
			t.Fatal("schema rejects requests supplied exclusively through queryValues")
		}
	}
	alternatives, ok := schema["allOf"].([]any)
	if !ok || len(alternatives) != 3 {
		t.Fatal("missing alternatives for language/startTime/endTime")
	}
}

func TestDomainConfigApplicationIdentityUsesDocumentedLocation(t *testing.T) {
	for _, tc := range []struct{ id, in string }{{"domain-api-get-domain-config", "query"}, {"domain-api-download-domain-config", "header"}} {
		e, _ := EndpointByID("domains", tc.id)
		found := false
		for _, p := range e.Parameters {
			if p.Name == "appId" {
				if p.In != tc.in || !p.Required || p.Type != "string" || p.MaxLength != 32 || p.SourceURL == "" {
					t.Fatalf("%s has wrong application contract: %+v", tc.id, p)
				}
				found = true
			}
		}
		if !found {
			t.Fatal("missing application identity")
		}
	}
}

func TestProjectLookupExposesDocumentedIdentifiersAndPagination(t *testing.T) {
	for _, tc := range []struct {
		id, in, name, kind string
		required           bool
	}{
		{"appbriefinfo", "header", "userID", "integer", true},
		{"appbriefinfo", "query", "packageNames", "string", false},
		{"appbriefinfo", "query", "fromRecCount", "integer", false},
		{"appbriefinfo", "query", "maxReqCount", "integer", false},
		{"queryservice", "query", "appID", "integer", true},
		{"queryappfingerprint", "header", "uid", "string", true},
		{"queryappfingerprint", "header", "teamId", "string", true},
	} {
		e, _ := EndpointByID("projects", tc.id)
		found := false
		for _, p := range e.Parameters {
			if p.Name == tc.name && p.In == tc.in {
				found = true
				if p.Type != tc.kind || p.Required != tc.required || p.SourceURL == "" {
					t.Errorf("%s incorrect documented field: %+v", tc.id, p)
				}
			}
		}
		if !found {
			t.Errorf("%s missing %s.%s", tc.id, tc.in, tc.name)
		}
	}
}

func TestProvisioningContractsUseDocumentedTypesAndBounds(t *testing.T) {
	for _, tc := range []struct {
		id, name, kind string
		max            int
	}{
		{"provision-api-applyacl", "aclPermissionList", "array", 0},
		{"provision-api-apply-provision", "aclPermissionList", "array", 1000},
		{"provision-api-apply-provision", "deviceIdList", "array", 4000},
		{"provision-api-add-device", "deviceList", "array", 1000},
	} {
		e, _ := EndpointByID("provisioning", tc.id)
		found := false
		for _, p := range e.Parameters {
			if p.In == "body" && p.Name == tc.name {
				found = true
				if p.Type != tc.kind || p.MaxItems != tc.max || p.SourceURL == "" {
					t.Errorf("%s.%s wrong source contract: %+v", tc.id, tc.name, p)
				}
			}
		}
		if !found {
			t.Errorf("%s missing %s", tc.id, tc.name)
		}
	}
}

func TestTestingMemberQueryUsesOfficialQueryFields(t *testing.T) {
	e, _ := EndpointByID("testing", "test-api-query-test-user")
	for _, name := range []string{"appId", "groupId", "current", "pageSize"} {
		found := false
		for _, p := range e.Parameters {
			if p.Name == name && p.In == "query" && p.SourceURL != "" {
				found = true
			}
		}
		if !found {
			t.Errorf("missing official query.%s", name)
		}
	}
}

func TestFeedbackContractsExposeActualRootModels(t *testing.T) {
	for _, tc := range []struct{ id, name, kind string }{{"test-api-query-feedback", "models", "array"}, {"test-api-query-feedback", "currentPage", "integer"}, {"test-api-delete-feedback", "delFeedBackInfos", "array"}, {"test-api-delete-feedback", "deleteMode", "integer"}} {
		e, _ := EndpointByID("testing", tc.id)
		found := false
		for _, p := range e.Parameters {
			if p.In == "body" && p.Name == tc.name && p.Type == tc.kind && p.SourceURL != "" {
				found = true
			}
		}
		if !found {
			t.Errorf("%s missing actual root field %s", tc.id, tc.name)
		}
	}
}

func TestResourceRequestsUseActualRootAndNestedFields(t *testing.T) {
	for _, tc := range []struct{ id, name, kind string }{{"respackapi-addversion", "resourceVersion", "object"}, {"respackapi-uploadapplyfor", "filePartApplyInfo", "array"}, {"respackapi-uploadconfirm", "filePartConfirmInfos", "array"}} {
		e, _ := EndpointByID("resources", tc.id)
		found := false
		for _, p := range e.Parameters {
			if p.Name == tc.name && p.In == "body" && p.Type == tc.kind {
				found = true
				if len(p.Properties)+len(p.ItemsProperties) == 0 {
					t.Error("missing reviewed nested fields")
				}
			}
		}
		if !found {
			t.Errorf("%s missing actual root %s", tc.id, tc.name)
		}
	}
}

func TestResourceRootFieldsRetainTheirOwnParagraphs(t *testing.T) {
	e, _ := EndpointByID("resources", "respackapi-uploadapplyfor")
	for _, p := range e.Parameters {
		if p.Name == "fileLength" && p.In == "body" {
			if p.Description != "升级包文件大小（字节数）。" {
				t.Fatalf("nested part table overwrote root description: %s", p.Description)
			}
			return
		}
	}
	t.Fatal("missing root fileLength")
}

func TestCommentReplyContractsFollowPlatformSpecificTypes(t *testing.T) {
	for _, tc := range []struct{ id, name, in, kind string }{{"com-reviewreply-harmonyos", "reviewId", "body", "integer"}, {"com-reviewreply", "reviewId", "body", "string"}, {"com-reviewreply", "requestId", "header", "string"}, {"comapi-getreviews", "apkVersions", "query", "string"}} {
		e, _ := EndpointByID("comments", tc.id)
		found := false
		for _, p := range e.Parameters {
			if p.Name == tc.name && p.In == tc.in && p.Type == tc.kind && p.SourceURL != "" {
				found = true
			}
		}
		if !found {
			t.Errorf("%s missing documented %s.%s %s", tc.id, tc.in, tc.name, tc.kind)
		}
	}
}

func TestUploadAddressExposesDocumentedLengthsAndEnums(t *testing.T) {
	e, _ := EndpointByID("upload", "upload-url-new")
	for _, tc := range []struct {
		name string
		max  int
	}{{"appId", 32}, {"fileName", 256}, {"sha256", 64}} {
		found := false
		for _, p := range e.Parameters {
			if p.Name == tc.name && p.In == "query" {
				found = true
				if p.MaxLength != tc.max {
					t.Errorf("%s missing documented maximum %d", tc.name, tc.max)
				}
			}
		}
		if !found {
			t.Errorf("missing query.%s", tc.name)
		}
	}
	for _, p := range e.Parameters {
		if p.Name == "releaseType" && (len(p.Enum) != 1 || p.Enum[0] != "1") {
			t.Fatal("upload releaseType lost documented allowed value")
		}
	}
}

func TestGameplayContractsSeparateHeadersAndJSONBody(t *testing.T) {
	e, _ := EndpointByID("gameplay", "basic-info")
	found := map[string]Parameter{}
	for _, p := range e.Parameters {
		found[p.In+"."+p.Name] = p
	}
	for _, name := range []string{"appId", "timestamp", "nonce", "sign"} {
		if !found["header."+name].Required {
			t.Errorf("missing mandatory signature header %s", name)
		}
	}
	users := found["body.userIds"]
	if users.Type != "array" || users.ItemsType != "object" || len(users.ItemsProperties) != 2 {
		t.Fatal("missing reviewed player list model")
	}
}

func TestGameResourceOpenAPIPreservesJSONStringEncoding(t *testing.T) {
	e, _ := EndpointByID("gameplay", "resources")
	body := invokeRequestSchemaProperties(e)["body"].(map[string]any)
	resource := body["properties"].(map[string]any)["resources"].(map[string]any)
	if resource["type"] != "string" || resource["contentMediaType"] != "application/json" {
		t.Fatal("resource JSON string wire type was lost")
	}
	inner := resource["contentSchema"].(map[string]any)
	if inner["type"] != "array" || inner["maxItems"] != 50 {
		t.Fatal("missing resource array contract")
	}
}

func TestLegacyPublishingRootAndNestedContracts(t *testing.T) {
	for _, tc := range []struct{ id, name, kind string }{{"add-packageurl", "downloadFileName", "string"}, {"app-info-update", "privacyPolicy", "string"}, {"app-submit-with-file", "phasedReleaseInfo", "object"}} {
		e, _ := EndpointByID("publishing", tc.id)
		found := false
		for _, p := range e.Parameters {
			if p.In == "body" && p.Name == tc.name && p.Type == tc.kind && p.SourceURL != "" {
				found = true
			}
		}
		if !found {
			t.Errorf("%s missing actual field %s", tc.id, tc.name)
		}
	}
}
