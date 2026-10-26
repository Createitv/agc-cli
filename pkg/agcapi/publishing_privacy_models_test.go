package agcapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/Createitv/agc-cli/pkg/domain"
)

// Body shapes are captured from the official HarmonyOS/atomic-service request
// examples. IDs/contact values are replaced with offline fixtures; the model
// nesting, explicit empty optional arrays/objects and boolean types are kept.
func TestPublishingPrivacyOfficialRequestExamples(t *testing.T) {
	e, _ := domain.EndpointByID("publishing", "privacy-agreement-update")
	for _, kind := range []string{"harmonyos", "atomic"} {
		t.Run(kind, func(t *testing.T) {
			body, err := os.ReadFile("testdata/privacy-" + kind + "-request.json")
			if err != nil {
				t.Fatal(err)
			}
			if err = ValidateEndpointRequest(e, nil, map[string]string{"appId": "fixture"}, body); err != nil {
				t.Fatal(err)
			}
		})
	}
}

// Expectations are independent of the registry and come from the reviewed
// SceneItem, DataItem, ThirdSdkInfo, DevicePermission and ExtendInfo pages.
func TestPublishingPrivacyDeepViolations(t *testing.T) {
	e, _ := domain.EndpointByID("publishing", "privacy-agreement-update")
	tooManySDKs, _ := json.Marshal(make([]map[string]any, 51))
	tooManyChapters, _ := json.Marshal(make([]map[string]any, 6))
	for _, tc := range []struct{ name, fields, field, source string }{
		{"personal-object", `"personalInfo":"invalid"`, "personalInfo", "agreementinfo-privacy"},
		{"sensitive-enum", `"personalInfo":{"requiredInfo":{"sceneItems":[{"dataItems":[{"sensitiveData":2}]}]}}`, "personalInfo.requiredInfo.sceneItems[0].dataItems[0].sensitiveData", "dataitem"},
		{"data-id-bound", `"personalInfo":{"adsInfo":{"sceneItems":[{"dataItems":[{"itemId":"` + strings.Repeat("x", 65) + `"}]}]}}`, "personalInfo.adsInfo.sceneItems[0].dataItems[0].itemId", "dataitem"},
		{"sdk-array-limit", `"thirdSdkInfo":{"sdkInfos":` + string(tooManySDKs) + `}`, "thirdSdkInfo.sdkInfos", "thirdsdkinfo"},
		{"sdk-label-shape", `"thirdSdkInfo":{"sdkInfos":[{"privacyLabels":["invalid"]}]}`, "thirdSdkInfo.sdkInfos[0].privacyLabels[0]", "sdkinfo"},
		{"sdk-purpose-bound", `"thirdSdkInfo":{"sdkInfos":[{"purpose":"` + strings.Repeat("x", 501) + `"}]}`, "thirdSdkInfo.sdkInfos[0].purpose", "sdkinfo"},
		{"permission-minimum", `"devicePermissions":[{"permissions":[]}]`, "devicePermissions[0].permissions", "devicepermission"},
		{"permission-desc-bound", `"devicePermissions":[{"permissions":[{"permissionDesc":"` + strings.Repeat("x", 501) + `"}]}]`, "devicePermissions[0].permissions[0].permissionDesc", "data-permission"},
		{"shared-url-bound", `"thirdSharedInfos":[{"officialPrivacyUrl":"` + strings.Repeat("x", 1025) + `"}]`, "thirdSharedInfos[0].officialPrivacyUrl", "thirdsharedinfo"},
		{"cross-border-boolean", `"extInfo":{"storageInfo":{"isCrossBorder":"true"}}`, "extInfo.storageInfo.isCrossBorder", "storageinfo"},
		{"contact-phone-bound", `"extInfo":{"contactInfo":{"phoneNum":["` + strings.Repeat("x", 32) + `"]}}`, "extInfo.contactInfo.phoneNum[0]", "contactinfo"},
		{"custom-chapter-limit", `"extInfo":{"customInfo":{"customInfos":` + string(tooManyChapters) + `}}`, "extInfo.customInfo.customInfos", "data-custominfo"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateEndpointRequest(e, nil, map[string]string{"appId": "fixture"}, []byte(`{"id":"fixture","agreementInfos":[{"lang":"zh",`+tc.fields+`}]}`))
			var detail *RequestValidationError
			if !errors.As(err, &detail) {
				t.Fatalf("expected deep input violation, got %v", err)
			}
			want := "body.agreementInfos[0]." + tc.field
			for _, v := range detail.Violations {
				if v.Field == want && strings.Contains(v.SourceURL, tc.source) && v.Suggestion != "" {
					return
				}
			}
			t.Fatalf("missing source-based violation %s: %+v", want, detail.Violations)
		})
	}
}

func TestPublishingPrivacyDeepValidInputs(t *testing.T) {
	e, _ := domain.EndpointByID("publishing", "privacy-agreement-update")
	for _, fields := range []string{
		`"personalInfo":{},"thirdSdkInfo":{},"extInfo":{}`,
		`"personalInfo":{"requiredInfo":{"sceneItems":[{"dataItems":[{"itemId":"device-id","sensitiveData":1}]}]}},"devicePermissions":[{"permissions":[{"permissionName":"ohos.permission.CAMERA"}]}]`,
		`"thirdSdkInfo":{"sdkInfos":[{"name":"fixture-sdk","privacyLabels":[{"sceneReason":"同步数据","dataItems":[{"sensitiveData":0}]}]}]},"thirdSharedInfos":[{"name":"fixture-provider","dataItems":[{}]}]`,
		`"extInfo":{"storageInfo":{"isCrossBorder":false,"storagePeriod":"-1"},"contactInfo":{"phoneNum":["0086-13700000000"]},"customInfo":{"customInfos":[{"title":"补充条款"}]}}`,
	} {
		body := []byte(`{"id":"fixture","agreementInfos":[{"lang":"zh",` + fields + `}]}`)
		if err := ValidateEndpointRequest(e, nil, map[string]string{"appId": "fixture"}, body); err != nil {
			t.Fatal(err)
		}
	}
}

func TestPublishingPrivacySceneArrayBoundaries(t *testing.T) {
	e, _ := domain.EndpointByID("publishing", "privacy-agreement-update")
	for _, tc := range []struct{ name, pattern, path string }{
		{"request-scenes", `"personalInfo":{"requiredInfo":{"sceneItems":%s}}`, "personalInfo.requiredInfo.sceneItems"},
		{"sdk-scenes", `"thirdSdkInfo":{"sdkInfos":[{"privacyLabels":%s}]}`, "thirdSdkInfo.sdkInfos[0].privacyLabels"},
		{"shared-data", `"thirdSharedInfos":[{"dataItems":%s}]`, "thirdSharedInfos[0].dataItems"},
	} {
		for _, size := range []int{200, 201} {
			t.Run(tc.name+"/"+strconv.Itoa(size), func(t *testing.T) {
				items := make([]map[string]any, size)
				for i := range items {
					items[i] = map[string]any{}
				}
				array, _ := json.Marshal(items)
				body := []byte(`{"id":"fixture","agreementInfos":[{"lang":"zh",` + fmt.Sprintf(tc.pattern, array) + `}]}`)
				err := ValidateEndpointRequest(e, nil, map[string]string{"appId": "fixture"}, body)
				if size == 200 {
					if err != nil {
						t.Fatal(err)
					}
					return
				}
				var detail *RequestValidationError
				if !errors.As(err, &detail) {
					t.Fatalf("201 entries must fail: %v", err)
				}
				for _, v := range detail.Violations {
					if v.Field == "body.agreementInfos[0]."+tc.path && v.Code == "array_length" {
						return
					}
				}
				t.Fatalf("array bound not reported: %+v", detail.Violations)
			})
		}
	}
}
