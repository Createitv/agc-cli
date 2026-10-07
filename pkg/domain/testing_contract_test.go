package domain

import "testing"

// 固定官方业务字段，防止命令仍可 dry-run 却遗漏真实接口必填参数。
func TestTestingLifecycleRequiredContracts(t *testing.T) {
	cases := []struct{ id, name, location string }{
		{"test-api-add-test-group", "appId", "header"},
		{"test-api-add-test-group", "groupName", "body"},
		{"test-api-edit-test-group", "appId", "header"},
		{"test-api-edit-test-group", "groupId", "body"},
		{"test-api-edit-test-group", "groupName", "body"},
		{"test-api-delete-test-group", "appId", "header"},
		{"test-api-delete-test-group", "groupId", "query"},
		{"test-api-add-invite-code", "groupId", "body"},
		{"test-api-add-invite-code", "invitationCodeValidDays", "body"},
		{"test-api-add-invite-code", "invitationCodeInviteLimit", "body"},
		{"test-api-get-invite-code", "appId", "header"},
		{"test-api-get-invite-code", "groupId", "query"},
		{"test-api-stop-invite-code", "appId", "header"},
		{"test-api-stop-invite-code", "invitationCodeId", "body"},
		{"test-api-delete-invite-code", "appId", "header"},
		{"test-api-delete-invite-code", "invitationCodeId", "query"},
		{"test-api-add-test-version", "appId", "query"},
		{"test-api-add-test-version", "releaseType", "body"},
		{"test-api-add-test-version", "testType", "body"},
		{"test-api-add-test-version", "testDesc", "body"},
		{"test-api-delete-test-version", "appId", "query"},
		{"test-api-delete-test-version", "versionId", "query"},
	}
	for _, c := range cases {
		t.Run(c.id+"/"+c.name, func(t *testing.T) {
			ep, _ := EndpointByID("testing", c.id)
			for _, p := range ep.Parameters {
				if p.Name == c.name && p.In == c.location && p.Required {
					return
				}
			}
			t.Fatalf("缺少必填字段 %s@%s", c.name, c.location)
		})
	}
}

func TestOpenAPIExposesTestingBodyContract(t *testing.T) {
	ep, _ := EndpointByID("testing", "test-api-add-invite-code")
	props := invokeRequestSchemaProperties(ep)
	body, ok := props["body"].(map[string]any)
	if !ok {
		t.Fatal("body 没有结构化契约")
	}
	fields, ok := body["properties"].(map[string]any)
	if !ok {
		t.Fatal("body 未声明业务字段")
	}
	if fields["invitationCodeValidDays"].(map[string]any)["type"] != "integer" {
		t.Fatal("有效期必须是整数")
	}
	required := body["required"].([]string)
	if len(required) != 3 {
		t.Fatalf("必填 body 字段 = %v", required)
	}
}

func TestVersionListAllowsTestDraftReadback(t *testing.T) {
	ep, ok := EndpointByID("publishing", "app-version-list")
	if !ok || ep.Method != "POST" || ep.Path != "/api/publish/v3/version/brief-info/list" {
		t.Fatal("缺少官方测试草稿读回接口")
	}
	if len(ep.Parameters) == 0 || ep.Parameters[0].Name != "appId" || ep.Parameters[0].In != "header" {
		t.Fatal("版本查询需要 appId 请求头")
	}
}

func TestDedicatedAppCreationContract(t *testing.T) {
	ep, ok := EndpointByID("publishing", "app-create")
	if !ok || ep.Path != "/api/publish/v3/app" || ep.Method != "POST" {
		t.Fatal("缺少专用测试应用创建接口")
	}
	for _, name := range []string{"appName", "parentType", "installationFree"} {
		found := false
		for _, p := range ep.Parameters {
			if p.Name == name && p.In == "body" && p.Required {
				found = true
			}
		}
		if !found {
			t.Fatalf("缺少创建应用字段 %s", name)
		}
	}
	info, ok := EndpointByID("publishing", "app-info-query-v3")
	if !ok || info.Path != "/api/publish/v3/app-info" {
		t.Fatal("缺少 HarmonyOS V3 应用信息读回")
	}
}
