package agcapi

import (
	"github.com/Createitv/agc-cli/pkg/domain"
	"testing"
)

func TestTypedRequiredFieldsRejectNull(t *testing.T) {
	ep, _ := domain.EndpointByID("provisioning", "provision-api-update-provision")
	if _, err := MarshalEndpointFields(ep, map[string]string{"provisionId": "fixture", "deviceIdList": "null"}); err == nil {
		t.Fatal("必填数组不能通过字符串 null 绕过校验")
	}
}

func TestTypedRequiredStringsRejectWhitespace(t *testing.T) {
	ep, _ := domain.EndpointByID("testing", "test-api-add-test-group")
	if err := ValidateEndpointBody(ep, []byte(`{"groupName":"  "}`)); err == nil {
		t.Fatal("必填群组名称不能为空白")
	}
}

func TestRequiredTypedBodyMustBePresent(t *testing.T) {
	ep, _ := domain.EndpointByID("testing", "test-api-add-test-group")
	for _, body := range [][]byte{nil, []byte(`{}`), []byte(`{"groupName":null}`), []byte(`{"groupName":""}`)} {
		if err := ValidateEndpointBody(ep, body); err == nil {
			t.Fatalf("required groupName accepted: %s", body)
		}
	}
}

func TestOptionalNullAndLiteralNullStringRemainValid(t *testing.T) {
	ep, _ := domain.EndpointByID("testing", "test-api-add-test-group")
	for _, body := range []string{`{"groupName":"null","groupType":null}`, `{"groupName":"fixture","groupType":0}`} {
		if err := ValidateEndpointBody(ep, []byte(body)); err != nil {
			t.Fatalf("valid request rejected: %v", err)
		}
	}
}
