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
	case "projects/getteamlist":
		// 官方仅支持 OAuth；不能将 Service Account JWT 或 API 客户端 token 当作 OAuth 用户授权。
		return []Parameter{typed("oauth2Token", "header", "string", true, "OAuth 用户授权的 Access Token；仅支持 OAuth 客户端方式，以 oauth2Token 请求头传入。")}
	case "publishing/package-compile-status-v3", "publishing/query-aabfile":
		return []Parameter{appQuery, typed("pkgIds", "query", "string", true, "待查询的软件包 ID；多个 ID 以英文逗号分隔")}
	case "pms/batchactproduct-harmonyosnext", "pms/batchactproduct-android", "pms/batchdeactproduct-harmonyosnext", "pms/batchdeactproduct-android":
		return []Parameter{appHeader, typed("requestId", "body", "string", true, "开发者自定义唯一请求标识"), stringsArray("products", true, 0, 0, "待变更状态的商品标识列表；激活限已下线商品，去激活限已上线商品")}
	case "pms/batchaddproduct-harmonyosnext", "pms/batchaddproduct-android", "pms/batchupdateproduct-harmonyosnext", "pms/batchupdateproduct-android":
		products := typed("products", "body", "array", true, "商品列表；各条目遵循对应平台的 ProductImportInfo 模型")
		products.ItemsType = "object"
		return []Parameter{appHeader, typed("requestId", "body", "string", true, "开发者自定义唯一请求标识"), products}
	case "pms/batchactpromotion-harmonyosnext", "pms/batchactpromotion-android", "pms/batchdeactpromotion-harmonyosnext", "pms/batchdeactpromotion-android":
		request := typed("requestId", "body", "string", true, "开发者自定义唯一请求标识")
		request.MaxLength = 64
		parameters := []Parameter{appHeader, request, stringsArray("promotions", true, 0, 100, "促销 ID 列表，每个 ID 最长 32 字符")}
		if id == "batchdeactpromotion-harmonyosnext" || id == "batchdeactpromotion-android" {
			immediate := typed("immediateEnd", "body", "string", false, "1：立即结束促销；不传或其他值不更新促销结束时间")
			immediate.MaxLength = 2
			parameters = append(parameters, immediate)
		}
		return parameters
	case "pms/addproduct-harmonyosnext":
		// ProductInfo 的 status 是预留字段；创建成功不代表审核通过或已生效。
		return []Parameter{appHeader, typed("requestId", "body", "string", true, "唯一请求序列，最长 64 字符"), typed("product", "body", "object", true, "商品定义；创建后需要审核才能上架生效")}
	case "pms/addproductgroup-harmonyosnext", "pms/updateproductgroup-harmonyosnext", "pms/addproductgroup-android", "pms/updateproductgroup-android":
		// 官方 ProductGroupInfo：创建时系统生成 groupId；groupName 必填。
		return []Parameter{appHeader, typed("requestId", "body", "string", true, "开发者生成的唯一请求序列，最长 64 字符"), typed("resource", "body", "object", true, "商品订阅分组信息；groupName 必填，更新时 groupId 必填；状态约束以对应版本官方模型为准")}
	case "pms/getproductgroup-harmonyosnext", "pms/getproductgroup-android":
		return []Parameter{appHeader, typed("requestId", "body", "string", true, "唯一请求序列，最长 64 字符"), typed("pageNum", "body", "integer", false, "起始页，默认 1"), typed("pageSize", "body", "integer", false, "每页 1 至 300 条，默认 20"), typed("orderBy", "body", "string", false, "以 JSON 表示 name、createTime 或 updateTime 的 asc/desc 排序规则")}
	case "cicd/build-hvigor":
		return []Parameter{typed("project", "body", "string", true, "本地 HarmonyOS 项目目录"), typed("hvigorPath", "body", "string", false, "Hvigor 可执行文件路径；默认检测 DevEco Studio 或 PATH"), typed("product", "body", "string", false, "产品名，默认 default"), typed("buildMode", "body", "string", false, "debug 或 release，默认 debug"), typed("task", "body", "string", false, "assembleHap：模块 HAP；assembleApp：项目 APP，默认 assembleHap")}
	case "upload/obbfile-init":
		// 官方初始化接口在 Query 中传参，POST 并不表示这些字段在 Body。
		return []Parameter{appQuery, typed("fileName", "query", "string", true, "带后缀的原始文件名"), typed("contentType", "query", "string", false, "对象 MIME 类型"), {Name: "fileType", In: "query", Type: "integer", Enum: []string{"1", "2", "3"}, Description: "1：HarmonyOS APP；2：OBB（默认）；3：APK/AAB。"}, {Name: "obbType", In: "query", Type: "integer", Enum: []string{"1", "2"}, Description: "OBB文件必填：1主扩展，2补丁扩展；非OBB文件不能传入。"}, typed("releaseType", "query", "integer", false, "1：商用；6：HarmonyOS 测试"), typed("chineseMainlandFlag", "query", "integer", false, "是否分发中国大陆：0 或 1")}
	case "upload/obbfile-uploadurl", "upload/obbfile-compose":
		// Body 是以 additionalProp1、additionalProp2 等编号为键的动态对象，不能将文档占位名 additionalPropX 设为必填字段。
		return []Parameter{typed("objectId", "query", "string", true, "本次分片初始化返回的对象 ID"), typed("nspUploadId", "query", "string", true, "本次分片初始化返回的上传 ID")}
	case "publishing/app-create":
		return []Parameter{typed("projectId", "query", "string", false, "可选项目 ID；省略时创建专用同名项目"), typed("appName", "body", "string", true, "应用名称，最长 64 字符"), typed("packageName", "body", "string", false, "HarmonyOS 应用必填且包名唯一；元服务由系统生成"), typed("parentType", "body", "integer", true, "13：应用；2：游戏"), typed("installationFree", "body", "integer", true, "0：HarmonyOS 应用；1：元服务")}
	case "publishing/app-info-query-v3":
		return []Parameter{appQuery, typed("lang", "query", "string", false, "可选语言；不传则返回全部语言"), typed("releaseType", "query", "integer", false, "1：商用版本；6：HarmonyOS 测试版本"), typed("versionId", "query", "string", false, "releaseType=6 时必填的测试版本 ID")}
	case "publishing/app-version-list":
		return []Parameter{appHeader, typed("packageName", "body", "string", false, "应用包名；CLI 已要求 appId 请求头"), typed("state", "body", "string", false, "可选版本状态，多个状态以逗号分隔")}
	case "publishing/language-info-update":
		return []Parameter{appQuery, typed("lang", "body", "string", true, "目标语言；新建语言与更新语言使用相同入口"), typed("appName", "body", "string", false, "目标语言的应用名称"), typed("appDesc", "body", "string", false, "目标语言的详细描述"), typed("briefInfo", "body", "string", false, "目标语言的简介"), typed("newFeatures", "body", "string", false, "目标语言的新版本特性")}
	case "publishing/language-info-delete":
		return []Parameter{appQuery, typed("lang", "query", "string", true, "待删除的非默认语言"), typed("releaseType", "query", "integer", false, "应用发布方式")}
	case "provisioning/provision-api-get-fingerprints", "provisioning/provision-api-getacl", "testing/test-api-get-test-grouplist":
		return []Parameter{appHeader}
	case "provisioning/provision-api-add-fingerprints", "provisioning/provision-api-delete-fingerprints":
		return []Parameter{appHeader, stringsArray("fingerprintList", true, 1, 50, "SHA-256 证书指纹列表；只修改专用应用的测试指纹")}
	case "provisioning/provision-api-apply-cent":
		return []Parameter{typed("csr", "body", "string", true, "公钥 CSR；私钥不得上传"), typed("certName", "body", "string", true, "专用证书名称"), typed("certType", "body", "integer", true, "1：调试；2：发布；3：In-house；4：受限二进制证书")}
	case "provisioning/provision-api-delete-cent":
		return []Parameter{stringsArray("certIds", true, 1, 100, "待删除的本次运行专用证书 ID 列表")}
	case "provisioning/provision-api-apply-provision":
		devices := stringsArray("deviceIdList", false, 0, 4000, "真实设备 ID 列表，最多 4000 项")
		devices.ItemsMaxLength = 32
		permissions := stringsArray("aclPermissionList", false, 0, 1000, "已获审批的 ACL 列表，最多 1000 项")
		permissions.ItemsMaxLength = 256
		return []Parameter{typed("provisionName", "body", "string", true, "专用 Profile 名称"), typed("provisionType", "body", "integer", true, "1：调试；2：发布；3：In-house；6：指定设备"), typed("certId", "body", "string", true, "本次生成的专用证书 ID"), typed("appId", "body", "string", true, "专用测试应用 ID"), devices, permissions}
	case "provisioning/provision-api-applyacl":
		model := "https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-provision-api-data-applypermissioninfo-0000002469038756"
		return []Parameter{appHeader, typed("attachment", "body", "string", false, "已上传附件的 objectId"), {Name: "aclPermissionList", In: "body", Type: "array", Required: true, ItemsType: "object", Description: "受限 ACL 申请列表；按官方请求示例展开为对象数组。", ItemsProperties: []Parameter{
			{Name: "name", In: "body", Type: "string", Required: true, MaxLength: 256, SourceURL: model},
			{Name: "reason", In: "body", Type: "string", Required: true, MaxLength: 256, SourceURL: model},
		}}}
	case "provisioning/provision-api-add-device":
		model := "https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-provision-api-data-adddeviceinfo-0000002236201314"
		return []Parameter{{Name: "deviceList", In: "body", Type: "array", Required: true, ItemsType: "object", MaxItems: 1000, Description: "真实调试/测试设备信息列表，最多1000项。", ItemsProperties: []Parameter{
			{Name: "deviceName", In: "body", Type: "string", Required: true, MaxLength: 100, SourceURL: model},
			{Name: "udid", In: "body", Type: "string", Required: true, MaxLength: 64, SourceURL: model},
			{Name: "deviceType", In: "body", Type: "integer", Required: true, Enum: []string{"1", "2", "4", "5", "8", "19"}, SourceURL: model},
		}}}
	case "provisioning/provision-api-delete-provision":
		return []Parameter{typed("id", "query", "string", true, "本次创建并登记的专用 Profile ID")}
	case "provisioning/provision-api-update-provision":
		return []Parameter{typed("provisionId", "body", "string", true, "本次创建的专用调试 Profile ID"), stringsArray("deviceIdList", true, 1, 4000, "更新后的真实设备 ID 列表；真实接口不接受空列表")}
	case "provisioning/provision-api-query-provision":
		return []Parameter{appHeader}
	case "domains/domain-api-get-domain":
		return []Parameter{appHeader, p("category", "query", true, "Domain category")}
	case "domains/domain-api-post-domain":
		// 全量替换指定分类的域名；空数组可清空。不得把 category 放在 Body。
		return []Parameter{appHeader, typed("category", "query", "string", true, "server：服务器域名；business：业务域名"), typed("domains", "body", "array", false, "全量域名列表，每项含 type 和 value；最多 3000 项")}
	case "domains/domain-api-precheck":
		return []Parameter{appHeader, typed("domain", "body", "object", false, "待校验的业务域名，type 必须为 webView；需提前部署校验文件")}
	case "domains/domain-api-download-domain-config", "domains/domain-api-get-domain-config":
		return []Parameter{appHeader}
	case "testing/test-api-query-test-user":
		return []Parameter{appHeader, p("groupId", "header", true, "Test group ID")}
	case "testing/test-api-add-test-group-user":
		return []Parameter{appHeader,
			{Name: "groupId", In: "body", Type: "string", Required: true, MaxLength: 32, Description: "测试群组ID。"},
			{Name: "hwAccount", In: "body", Type: "string", Required: true, MaxLength: 256, Description: "真实华为账号、邮箱或手机号；内部群组手机号格式0086手机号，外部群组为86-手机号。"},
			{Name: "nickName", In: "body", Type: "string", MaxLength: 64, Description: "测试成员昵称。"},
			{Name: "uid", In: "body", Type: "string", MaxLength: 128, Description: "联盟账号ID；内部群组添加成员必填，需先确认群组类型。"},
		}
	case "testing/test-api-delete-test-group-user":
		return []Parameter{appHeader,
			{Name: "groupId", In: "body", Type: "string", Required: true, MaxLength: 32, Description: "测试群组ID。"},
			{Name: "deleteTesterIds", In: "body", Type: "array", Required: true, ItemsType: "string", ItemsMaxLength: 100, MaxItems: 100, Description: "真实测试成员ID列表；数组长度不超过100，每个ID最长100字符。"},
		}
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
	case "testing/test-api-add-test-package":
		// file 是已经上传的 .app 文件元数据，不是把二进制文件直接放入此 JSON 请求。
		return []Parameter{appQuery, typed("file", "body", "object", true, "文件信息，包含 fileName 和上传返回的 objectId；需随后轮询异步软件包解析"), typed("distributeMode", "body", "integer", false, "1：测试专区；2：应用市场")}
	case "testing/test-api-modify-test-version":
		// 已核对官方参数表，并在专用应用上读回草稿描述变更。
		return []Parameter{appQuery, typed("versionId", "body", "string", true, "本次创建的测试草稿版本 ID；在架版本不可更新"), typed("pkgId", "body", "string", false, "成功编译的软件包 ID"), typed("languages", "body", "array", false, "版本语言信息，最多 100 项"), typed("deviceTypes", "body", "array", false, "设备类型，最多 6 项"), typed("distCountryList", "body", "string", false, "公开测试国家码，以英文逗号分隔"), typed("extInfo", "body", "object", false, "版本扩展信息"), typed("openTestInfo", "body", "object", false, "测试时间、描述与任务信息"), typed("packagePermissionIntroVideoList", "body", "array", false, "权限说明视频"), typed("versionVocReq", "body", "array", false, "公开测试已解决问题")}
	case "testing/test-api-delete-test-version":
		return []Parameter{appQuery, typed("versionId", "query", "string", true, "本次创建的测试草稿版本 ID")}
	case "upload/upload-url-new":
		return []Parameter{appQuery, typed("fileName", "query", "string", true, "包含后缀的原始文件名"), typed("contentLength", "query", "integer", true, "文件长度，单位字节"), typed("sha256", "query", "string", false, "文件内容 SHA256"), typed("releaseType", "query", "integer", false, "发布方式"), typed("chineseMainlandFlag", "query", "integer", false, "是否分发中国大陆；非中国大陆注册开发者必须设置")}
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
