package domain

// verifiedParameters overrides the method-based legacy inference only for
// contracts verified against the linked Huawei reference or upstream responses.
// Optional entries expose known locations without asserting requiredness that
// has not been verified. This is intentionally not a complete API schema.
func verifiedParameters(family, id string, inferred []Parameter) []Parameter {
	p := func(name, location string, required bool, description string) Parameter {
		return Parameter{Name: name, In: location, Required: required, Description: description}
	}
	appHeader := p("appId", "header", true, "Huawei application ID")
	// 下列群组、邀请码和测试草稿字段已经逐项核对官方参数表。
	// 认证方式的可选分支由授权层处理，不能把三种认证头全部设为必填。
	typed := func(name, location, kind string, required bool, description string) Parameter {
		field := p(name, location, required, description)
		field.Type = kind
		return field
	}
	appQuery := typed("appId", "query", "string", true, "应用 ID")
	stringsArray := func(name string, required bool, minimum, maximum int, description string) Parameter {
		field := typed(name, "body", "array", required, description)
		field.ItemsType, field.MinItems, field.MaxItems = "string", minimum, maximum
		return field
	}
	switch family + "/" + id {
	case "publishing/app-create":
		return []Parameter{typed("projectId", "query", "string", false, "可选项目 ID；省略时创建专用同名项目"), typed("appName", "body", "string", true, "应用名称，最长 64 字符"), typed("packageName", "body", "string", false, "HarmonyOS 应用必填且包名唯一；元服务由系统生成"), typed("parentType", "body", "integer", true, "13：应用；2：游戏"), typed("installationFree", "body", "integer", true, "0：HarmonyOS 应用；1：元服务")}
	case "publishing/app-info-query-v3":
		return []Parameter{appQuery, typed("lang", "query", "string", false, "可选语言；不传则返回全部语言"), typed("releaseType", "query", "integer", false, "1：商用版本；6：HarmonyOS 测试版本"), typed("versionId", "query", "string", false, "releaseType=6 时必填的测试版本 ID")}
	case "publishing/app-version-list":
		return []Parameter{appHeader, typed("packageName", "body", "string", false, "应用包名；CLI 已要求 appId 请求头"), typed("state", "body", "string", false, "可选版本状态，多个状态以逗号分隔")}
	case "provisioning/provision-api-get-fingerprints", "provisioning/provision-api-getacl", "testing/test-api-get-test-grouplist":
		return []Parameter{appHeader}
	case "provisioning/provision-api-add-fingerprints", "provisioning/provision-api-delete-fingerprints":
		return []Parameter{appHeader, stringsArray("fingerprintList", true, 1, 50, "SHA-256 证书指纹列表；只修改专用应用的测试指纹")}
	case "provisioning/provision-api-apply-cent":
		return []Parameter{typed("csr", "body", "string", true, "公钥 CSR；私钥不得上传"), typed("certName", "body", "string", true, "专用证书名称"), typed("certType", "body", "integer", true, "1：调试；2：发布；3：In-house；4：受限二进制证书")}
	case "provisioning/provision-api-delete-cent":
		return []Parameter{stringsArray("certIds", true, 1, 100, "待删除的本次运行专用证书 ID 列表")}
	case "provisioning/provision-api-apply-provision":
		return []Parameter{typed("provisionName", "body", "string", true, "专用 Profile 名称"), typed("provisionType", "body", "integer", true, "1：调试；2：发布；3：In-house；6：指定设备"), typed("certId", "body", "string", true, "本次生成的专用证书 ID"), typed("appId", "body", "string", true, "专用测试应用 ID"), typed("deviceIdList", "body", "array", false, "真实设备 ID 列表，最多 4000 项"), typed("aclPermissionList", "body", "array", false, "已获审批的 ACL 列表，最多 1000 项")}
	case "provisioning/provision-api-delete-provision":
		return []Parameter{typed("id", "query", "string", true, "本次创建并登记的专用 Profile ID")}
	case "provisioning/provision-api-update-provision":
		return []Parameter{typed("provisionId", "body", "string", true, "本次创建的专用调试 Profile ID"), stringsArray("deviceIdList", true, 1, 4000, "更新后的真实设备 ID 列表；真实接口不接受空列表")}
	case "provisioning/provision-api-query-provision":
		return []Parameter{appHeader}
	case "domains/domain-api-get-domain":
		return []Parameter{appHeader, p("category", "query", true, "Domain category")}
	case "testing/test-api-query-test-user":
		return []Parameter{appHeader, p("groupId", "header", true, "Test group ID")}
	case "testing/test-api-add-test-group":
		return []Parameter{appHeader, typed("groupName", "body", "string", true, "群组名称，最长 50 个字符"), typed("groupType", "body", "integer", false, "0：外部群组；1：内部群组（需要 AppTest）")}
	case "testing/test-api-edit-test-group":
		return []Parameter{appHeader, typed("groupId", "body", "string", true, "专用测试群组 ID"), typed("groupName", "body", "string", true, "群组名称，最长 50 个字符")}
	case "testing/test-api-delete-test-group":
		return []Parameter{appHeader, typed("groupId", "query", "string", true, "待删除的专用群组 ID")}
	case "testing/test-api-add-invite-code":
		return []Parameter{appHeader, typed("groupId", "body", "string", true, "测试群组 ID"), typed("invitationCodeValidDays", "body", "integer", true, "有效期天数；官方正文要求 30 天"), typed("invitationCodeInviteLimit", "body", "integer", true, "可邀请人数上限，不超过 10000")}
	case "testing/test-api-get-invite-code":
		return []Parameter{appHeader, typed("groupId", "query", "string", true, "测试群组 ID")}
	case "testing/test-api-stop-invite-code":
		return []Parameter{appHeader, typed("invitationCodeId", "body", "string", true, "要停止的邀请码 ID")}
	case "testing/test-api-delete-invite-code":
		// 官方 Query 表缺项，以接口原型中的 query 和请求示例为依据。
		return []Parameter{appHeader, typed("invitationCodeId", "query", "string", true, "要删除的邀请码 ID；AppTest 邀请测试不支持此接口")}
	case "testing/test-api-add-test-version":
		return []Parameter{appQuery, typed("releaseType", "body", "integer", true, "HarmonyOS 测试发布方式，默认 6"), typed("testType", "body", "integer", true, "3：邀请测试；4：公开测试"), typed("testDesc", "body", "string", true, "测试版本描述，最长 50 字符"), typed("onshelfSelfDetect", "body", "integer", false, "0：不进行正式上架自检；1：进行自检")}
	case "testing/test-api-delete-test-version":
		return []Parameter{appQuery, typed("versionId", "query", "string", true, "本次创建的测试草稿版本 ID")}
	case "upload/upload-url-new":
		return []Parameter{p("contentLength", "query", true, "Upload file length in bytes"), p("appId", "query", false, "Huawei application ID"), p("suffix", "query", false, "Upload file suffix"), p("fileName", "query", false, "Upload file name")}
	case "publishing/app-info-query":
		return []Parameter{p("appId", "query", true, "Huawei application ID"), p("lang", "query", false, "Application information language")}
	case "projects/queryservice":
		return []Parameter{p("projectId", "query", true, "Huawei project ID"), p("appID", "query", true, "Huawei application ID; parameter name is case sensitive")}
	case "comments/comapi-getreviews-harmonyos", "comments/com-rating-harmonyos":
		return harmonyOSCommentsParameters(id)
	case "reports/appdownloadexport":
		return append(inferred, p("language", "query", false, "Report language"), p("startTime", "query", false, "Start date in YYYYMMDD format"), p("endTime", "query", false, "End date in YYYYMMDD format"))
	case "reports/iapexport":
		return append(inferred, p("currency", "query", true, "Report currency"))
	case "reports/orderanalysisexport":
		return append(inferred, p("teamId", "header", true, "Huawei team ID"))
	default:
		return inferred
	}
}

// Huawei query tables for agcapi-comapi-getreviews-harmonyos-0000002470893976
// and agcapi-com-rating-harmonyos-0000002503933921, updated 2026-04-29.
func harmonyOSCommentsParameters(id string) []Parameter {
	p := func(name, typ, format string, required bool, description string) Parameter {
		return Parameter{Name: name, In: "query", Type: typ, Format: format, Required: required, Description: description}
	}
	result := []Parameter{
		p("appId", "string", "", true, "Huawei application ID"),
		p("beginTime", "integer", "int64", true, "Start timestamp in milliseconds since Unix epoch; query interval must not exceed six months"),
		p("endTime", "integer", "int64", true, "End timestamp in milliseconds since Unix epoch; query interval must not exceed six months"),
		p("countries", "string", "", true, "Comma-separated country codes, such as CN; all countries must belong to the same regional site"),
		p("ratings", "string", "", false, "Comma-separated ratings from 1 to 5"),
		p("appVersions", "string", "", false, "Comma-separated application versions"),
		p("sort", "integer", "int32", false, "0: descending time (default); 1: descending score; 2: ascending score"),
		p("page", "integer", "int32", false, "Page number; defaults to 1"),
	}
	if id == "comapi-getreviews-harmonyos" {
		result = append(result, p("limit", "integer", "int32", false, "Page size; defaults to 20, maximum 100"),
			p("content", "string", "", false, "Comment content keyword"),
			p("devReplyStates", "string", "", false, "Comma-separated developer reply states: 0, 1, 6, 3"),
			p("langs", "string", "", false, "Comma-separated ISO 639-1 language codes"))
	} else {
		result = append(result, p("limit", "integer", "int32", false, "Page size; defaults to 50, maximum 100"))
	}
	return result
}
