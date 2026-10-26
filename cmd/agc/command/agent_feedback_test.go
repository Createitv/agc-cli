package command

import (
	"encoding/json"
	"strings"
	"testing"
)

// The missing test-version ID must fail before authentication/network access
// and give the Agent a documented, machine-readable correction.
func TestAgentReceivesConditionalParameterFeedback(t *testing.T) {
	out, err := execute("--project", t.TempDir(), "publishing", "app-info-query-v3", "--invoke", "--query", "appId=fixture", "--query", "releaseType=6")
	if err == nil {
		t.Fatal("test-version request without versionId was accepted")
	}
	var result map[string]any
	if json.Unmarshal([]byte(out), &result) != nil {
		t.Fatalf("missing structured feedback: %q, %v", out, err)
	}
	raw, _ := json.Marshal(result)
	for _, want := range []string{"versionId", "sourceUrl", "suggestion", "invalid_request"} {
		if !strings.Contains(string(raw), want) {
			t.Errorf("feedback missing %s: %s", want, raw)
		}
	}
}

func TestVersionListAcceptsPackageNameWithoutApplicationHeader(t *testing.T) {
	_, err := execute("--project", t.TempDir(), "publishing", "app-version-list", "--invoke", "--field", "packageName=app.example.fixture")
	if err != nil {
		t.Fatalf("official appId-or-packageName branch rejected: %v", err)
	}
}

func TestSourceDocumentArrayLimitIsEnforced(t *testing.T) {
	_, err := execute("--project", t.TempDir(), "publishing", "qualification-review-create", "--invoke", "--header", "appId=fixture", "--field", `tagQualInfos=[{},{},{},{},{},{}]`)
	if err == nil {
		t.Fatal("six qualification groups accepted; official maximum is five")
	}
}

func TestBatchProductFieldsKeepDocumentedArrayShape(t *testing.T) {
	out, err := execute("--project", t.TempDir(), "pms", "batchactproduct-harmonyosnext", "--pretty")
	if err != nil {
		t.Fatal(err)
	}
	var data struct {
		Data struct {
			Parameters []struct{ Name, In, Type string }
		}
	}
	if json.Unmarshal([]byte(out), &data) != nil {
		t.Fatal("invalid discovery output")
	}
	found := false
	for _, p := range data.Data.Parameters {
		if p.Name == "products" && p.In == "body" && p.Type == "array" {
			found = true
		}
	}
	if !found {
		t.Fatal("official products array is missing from PMS contract")
	}
}

func TestUploadURLRequiresActualFileNameAndApplication(t *testing.T) {
	_, err := execute("--project", t.TempDir(), "upload", "upload-url-new", "--invoke", "--query", "contentLength=1")
	if err == nil {
		t.Fatal("upload URL accepted without required appId and fileName")
	}
}

func TestFileGroupRejectsUndocumentedDeviceAndMissingObjects(t *testing.T) {
	out, err := execute("--project", t.TempDir(), "publishing", "app-file-info-v3", "--invoke", "--query", "appId=fixture", "--field", `screenShotList=[{"lang":"zh_CN","fileInfoList":[{"deviceType":999}]}]`)
	if err == nil {
		t.Fatal("invalid nested FileInfo accepted")
	}
	for _, want := range []string{"deviceType", "objectIdList", "sourceUrl"} {
		if !strings.Contains(out, want) {
			t.Errorf("feedback missing %s: %s", want, out)
		}
	}
}

func TestBatchSubscriptionRequiresDocumentedSubscriptionFields(t *testing.T) {
	out, err := execute("--project", t.TempDir(), "pms", "batchaddproduct-harmonyosnext", "--invoke", "--header", "appId=fixture", "--field", "requestId=fixture", "--field", `products=[{"productNo":"fixture","productName":"fixture","productDesc":"fixture","purchaseType":"auto_subscription","defaultLocale":"zh_CN"}]`)
	if err == nil {
		t.Fatal("subscription without subscription group and period accepted")
	}
	for _, want := range []string{"products[0].subGroupId", "products[0].subPeriod", "sourceUrl"} {
		if !strings.Contains(out, want) {
			t.Errorf("feedback missing %s: %s", want, out)
		}
	}
}

func TestBatchProductRejectsMalformedLanguageAndPriceModels(t *testing.T) {
	base := `{"productNo":"fixture","productName":"fixture","productDesc":"fixture","purchaseType":"consumable","defaultLocale":"zh_CN",`
	for _, tc := range []struct{ name, extra, field string }{
		{"long-language-name", `"languages":[{"locale":"en_US","productName":"` + strings.Repeat("a", 56) + `","productDesc":"fixture"}]}`, "products[0].languages[0].productName"},
		{"missing-price-currency", `"defaultPriceInfo":{"country":"CN"}}`, "products[0].defaultPriceInfo.currency"},
		{"wrong-language-type", `"languages":"invalid"}`, "products[0].languages"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			out, err := execute("--project", t.TempDir(), "pms", "batchaddproduct-harmonyosnext", "--invoke", "--header", "appId=fixture", "--field", "requestId=fixture", "--field", "products=["+base+tc.extra+"]")
			if err == nil {
				t.Fatal("malformed nested model accepted")
			}
			if !strings.Contains(out, tc.field) {
				t.Fatalf("missing nested field feedback: %s", out)
			}
		})
	}
}

func TestReportDateWindowAndFilterPairsGiveCorrectiveFeedback(t *testing.T) {
	for _, tc := range []struct {
		name  string
		query []string
		field string
	}{
		{"invalid-date", []string{"startTime=20260230", "endTime=20261001"}, "startTime"},
		{"too-long", []string{"startTime=20260101", "endTime=20261001"}, "endTime"},
		{"mismatched-filter", []string{"startTime=20260901", "endTime=20261001", "filterCondition=countryId", "filterCondition=businessType", "filterConditionValue=CN"}, "filterConditionValue"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			args := []string{"--project", t.TempDir(), "reports", "appdownloadexport", "--invoke", "--param", "appId=fixture", "--query", "language=en-US"}
			for _, q := range tc.query {
				args = append(args, "--query", q)
			}
			out, err := execute(args...)
			if err == nil {
				t.Fatal("invalid report request accepted")
			}
			if !strings.Contains(out, tc.field) || !strings.Contains(out, "sourceUrl") {
				t.Fatalf("unhelpful feedback: %s", out)
			}
		})
	}
}
