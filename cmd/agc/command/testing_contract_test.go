package command

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGroupCreateRequiresNameBeforeInvocation(t *testing.T) {
	_, err := execute("testing", "test-api-add-test-group", "--invoke", "--header", "appId=fixture")
	if err == nil || !strings.Contains(err.Error(), "groupName") {
		t.Fatalf("缺少名称必须被拒绝，实际错误：%v", err)
	}
}

func TestInvitationFieldsSendNumbers(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, _ := io.ReadAll(r.Body)
		var body map[string]json.RawMessage
		_ = json.Unmarshal(data, &body)
		if string(body["invitationCodeValidDays"]) != "30" {
			t.Errorf("有效期应为 JSON 数字：%s", data)
		}
		if string(body["invitationCodeInviteLimit"]) != "1" {
			t.Errorf("人数应为 JSON 数字：%s", data)
		}
		_, _ = io.WriteString(w, `{"rtnCode":0}`)
	}))
	defer server.Close()
	_, err := execute("testing", "test-api-add-invite-code", "--invoke", "--dry-run=false", "--base-url", server.URL, "--token", "offline", "--header", "appId=fixture", "--field", "groupId=fixture", "--field", "invitationCodeValidDays=30", "--field", "invitationCodeInviteLimit=1")
	if err != nil {
		t.Fatal(err)
	}
}

func TestTestingBodyRejectsIncorrectTypes(t *testing.T) {
	file := filepath.Join(t.TempDir(), "body.json")
	if err := os.WriteFile(file, []byte(`{"groupName":42}`), 0600); err != nil {
		t.Fatal(err)
	}
	_, err := execute("testing", "test-api-add-test-group", "--invoke", "--header", "appId=fixture", "--body", file)
	if err == nil {
		t.Fatal("数字群组名称不能绕过类型校验")
	}
	_, err = execute("testing", "test-api-add-test-group", "--invoke", "--header", "appId=fixture", "--field", "groupName=fixture", "--field", "groupType=abc")
	if err == nil {
		t.Fatal("非整数群组类型必须拒绝")
	}
}

func TestDebugProfileCannotClearAllDevices(t *testing.T) {
	_, err := execute("provisioning", "provision-api-update-provision", "--invoke", "--field", "provisionId=fixture", "--field", "deviceIdList=[]")
	if err == nil {
		t.Fatal("真实接口已经确认调试 Profile 不接受空设备列表")
	}
}
