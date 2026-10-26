package domain

import "strings"

// Current Publishing API operations and parameters from Huawei reference tables.
// Existing V2 commands are retained; V3 updates/submission have explicit V3 IDs.
func currentPublishingEndpoints() []Endpoint {
	return []Endpoint{
		currentPublishingEndpoint("app-create", "创建APP ID", "POST", "/api/publish/v3/app", "https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-provision-api-create-appid-0000002526543939", []Parameter{
			{Name: "projectId", In: "query", Type: "string", Required: false, Description: "项目ID，获取方法请参见查看应用信息。 如果不传projectId，系统会先在内部创建一个与应用同名的项目，随后再创建应用。"},
			{Name: "appName", In: "body", Type: "string", Required: true, Description: "应用名称。"},
			{Name: "packageName", In: "body", Type: "string", Required: false, Description: "应用包名，必须唯一，且不能与其他应用（包括Android应用）包名相同。 说明 installationFree取值为0（HarmonyOS应用）时，此参数必填。 installationFree取值为1（元服务）时，系统将自动生成包名，此参数无需填写。"},
			{Name: "parentType", In: "body", Type: "integer", Required: true, Description: "应用分类。 取值范围： 2：游戏 13：应用"},
			{Name: "installationFree", In: "body", Type: "integer", Required: true, Description: "应用类型。 取值范围： 0：HarmonyOS应用 1：元服务 默认值：0"},
		}),
		currentPublishingEndpoint("appid-list", "查询应用包名对应的appid", "GET", "/api/publish/v2/appid-list", "https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-appid-list-0000002271000617", []Parameter{
			{Name: "packageName", In: "query", Type: "string", Required: true, Description: "需要查询的应用包名。 多个包名以逗号分隔，最多支持50个。"},
			{Name: "packageTypes", In: "query", Type: "string", Required: false, Description: "通过软件包类型过滤应用，只查询指定软件包类型的应用。 取值范围： 7：APP（HarmonyOS应用/元服务）"},
		}),
		currentPublishingEndpoint("members-role-list", "查询团队账号下的成员角色列表", "GET", "/api/publish/v3/members-role-list", "https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-teamid-member-list-0000002679253743", []Parameter{
			{Name: "fromPage", In: "query", Type: "integer", Required: true, Description: "起始页。 最小值：1，表示从第1页开始查询。"},
			{Name: "pageSize", In: "query", Type: "integer", Required: true, Description: "每页查询的条数。 取值范围：[1,50]"},
			{Name: "languageCode", In: "query", Type: "string", Required: true, Description: "语言码。 取值范围： zh_cn：中文 en_us：英文 ru_ru：俄语"},
		}),
		currentPublishingEndpoint("app-info-query-v3", "查询应用信息", "GET", "/api/publish/v3/app-info", "https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-appinfo-query-0000002236041422", []Parameter{
			{Name: "appId", In: "query", Type: "string", Required: true, Description: "需要查询的应用ID。"},
			{Name: "lang", In: "query", Type: "string", Required: false, Description: "需要查询的语言。 说明 取值参考语言类型。 不传该参数，则查询全部语言的应用信息。"},
			{Name: "releaseType", In: "query", Type: "integer", Required: false, Description: "应用发布方式。 取值范围： 1：全网 6：HarmonyOS测试发布方式（API>=10） 默认值：1"},
			{Name: "versionId", In: "query", Type: "string", Required: false, Description: "需要查询的版本ID。 注意 releaseType为6，该字段必填。"},
		}),
		currentPublishingEndpoint("app-info-update-v3", "更新应用基本信息", "PUT", "/api/publish/v3/app-info", "https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-appinfo-update-0000002236201246", []Parameter{
			{Name: "appId", In: "query", Type: "string", Required: true, Description: "需要更新的应用ID。"},
			{Name: "releaseType", In: "query", Type: "integer", Required: false, Description: "应用发布方式。 取值范围： 1：全网 默认值：1"},
			{Name: "releasePhase", In: "query", Type: "integer", Required: false, Description: "分阶段发布标识。 取值范围： 0：全网 3：分阶段 默认值：0"},
			{Name: "defaultLang", In: "body", Type: "string", Required: false, Description: "应用的默认语言。 说明 取值参见语言类型。"},
			{Name: "harmonyChildType", In: "body", Type: "integer", Required: false, Description: "HarmonyOS应用/元服务的二级分类。 固定ID。 说明 取值参见应用游戏分类。"},
			{Name: "kindMainTag", In: "body", Type: "integer", Required: false, Description: "HarmonyOS应用/元服务主标签。 固定ID。 说明 取值参见应用游戏分类的“标签”列。"},
			{Name: "kindSubTags", In: "body", Type: "array", Required: false, Description: "HarmonyOS应用/元服务附属标签。 固定ID。最多可选4个附属标签。 说明 取值参见应用游戏分类的“标签”列。", ItemsType: "integer"},
			{Name: "privacyPolicy", In: "body", Type: "string", Required: false, Description: "隐私政策网址。"},
			{Name: "appNetType", In: "body", Type: "integer", Required: false, Description: "应用联网类型。 取值范围： 1：单机 2：网游 默认值：1"},
			{Name: "publishCountry", In: "body", Type: "string", Required: true, Description: "应用发布国家的国家码。 多个国家以英文逗号分隔，例如：“CN,DE,MC”。 如果取值为“OTHER”，则表示当App Gallery支持的国家或地区列表中新增某个国家或地区时，应用将自动分发到新增的国家和地区。 注意 取值不能只包含OTHER。 如果国家码有设置非CN的地区，将存在个人信息出境的情况，具体请参考个人信息出境授权声明。"},
			{Name: "appTariffType", In: "body", Type: "string", Required: false, Description: "应用内付费道具类型。 取值范围： 1：激活收费 2：道具收费 3：关卡收费 4：购买虚拟币 5：部分章节收费（图书阅读类） 6：其他的 7：课程收费 8：会员收费 支持设置多个，使用英文逗号分隔，例如“1,2,3”。"},
			{Name: "publicationNumber", In: "body", Type: "string", Required: false, Description: "版号信息。 此参数只支持游戏类应用。"},
			{Name: "developerWebsite", In: "body", Type: "string", Required: false, Description: "开发者官网。"},
			{Name: "customInfo", In: "body", Type: "object", Required: false, Description: "客服信息。"},
			{Name: "deviceTypes", In: "body", Type: "array", Required: false, Description: "设备类型信息。", ItemsType: "object"},
			{Name: "privacyRightsUrl", In: "body", Type: "string", Required: false, Description: "隐私权利网址。 即关于您的用户实施其权利的网站地址，如用户删除、修改、导出其个人数据的入口。"},
			{Name: "registeredDclType", In: "body", Type: "integer", Required: false, Description: "APP类型信息。 取值范围： -1：清空APP类型信息 0：APP服务器在中国大陆 1：APP服务器不在中国大陆 2：APP为单机APP"},
			{Name: "registeredIdType", In: "body", Type: "integer", Required: false, Description: "主办单位类型。 取值范围： -1：清空数据（会联动清空registeredIdNumber的内容） 1：企业（社会信用代码） 2：个人（身份证号码） 3：机构（机构代码）"},
			{Name: "registeredIdNumber", In: "body", Type: "string", Required: false, Description: "主办单位证件号。 取值范围： registeredIdType为1时，对应统一社会信用代码 registeredIdType为2时，对应身份证号码 registeredIdType为3时，对应组织机构代码"},
			{Name: "registeredEntity", In: "body", Type: "integer", Required: false, Description: "主办单位。 取值范围： 1：与开发者账号主体一致 2：与开发者账号主体不一致 说明 当registeredDclType值为0（APP服务器在中国大陆）时，该字段必填。若不填写或填写错误，会导致审核不通过。 若在版本提交时，该字段值选择“一致”，则主办单位名称和证件号，将直接复用开发者联盟注册账号信息。 若您提交的主办单位信息与联盟账号信息不一致，该字段值却选择了“一致”，则系统将自动判定为“不一致”，将会影响审核结果。"},
			{Name: "registeredEntityName", In: "body", Type: "string", Required: false, Description: "主办单位名称。 说明 当registeredDclType值为0（APP服务器在中国大陆）时，该字段必填。若不填写或填写错误，会导致审核不通过。"},
			{Name: "privacyAgreementId", In: "body", Type: "string", Required: false, Description: "隐私协议ID。"},
			{Name: "userAgreementUrl", In: "body", Type: "string", Required: false, Description: "用户协议URL。"},
			{Name: "userAgreementIds", In: "body", Type: "array", Required: false, Description: "用户协议ID列表。 当userAgreementUrl和userAgreementIds同时传值时，以userAgreementIds值为准。 数组长度不超过5。", ItemsType: "string"},
			{Name: "encrypted", In: "body", Type: "integer", Required: true, Description: "包是否加密。 取值范围： 0：否，表示不加密。 1：是，表示加密。 默认值：1"},
			{Name: "isAiGenerate", In: "body", Type: "integer", Required: false, Description: "是否涉及AI生成合成服务。 取值范围： 0：不涉及。 1：涉及。"},
			{Name: "aiTypes", In: "body", Type: "array", Required: false, Description: "涉及的AI生成合成服务类型。 取值范围： -1：其他 1：文本 2：图片 3：音频 4：视频 5：虚拟场景 说明 isAiGenerate取值为1（涉及）时，此参数必填。", ItemsType: "integer"},
			{Name: "testUserName", In: "body", Type: "string", Required: false, Description: "应用审核信息的测试账号用户名。 说明 如果存在需要登录后才可使用的功能，需要提供测试账号用户名。"},
			{Name: "testUserPassword", In: "body", Type: "string", Required: false, Description: "应用审核信息的测试账号密码。 说明 如果存在需要登录后才可使用的功能，需要提供测试账号密码。"},
			{Name: "appRemark", In: "body", Type: "string", Required: false, Description: "应用审核信息的备注。 提供有助于审核人员更准确、高效测试您应用的额外信息，比如在审核时需要的特别设置等。"},
			{Name: "selfTestVideos", In: "body", Type: "string", Required: false, Description: "应用审核信息的自测文件。 格式为JSON字符串，示例：{\"objectId1\":\"文件名1\",\"objectId2\":\"文件名2\"...} 说明 objectId的获取流程：先调用获取上传文件地址接口，然后调用上传文件接口上传文件，最后得到获取上传文件地址接口返回的objectId。 最多可上传5个自测视频，总大小不超过500MB。 注意 此字段只支持元服务。 元服务需要与设备联动时，此参数必填。"},
			{Name: "appReviewPhoneInfo", In: "body", Type: "object", Required: false, Description: "应用负责人的手机号信息。用于沟通审核问题。"},
			{Name: "appReviewEmailInfo", In: "body", Type: "object", Required: false, Description: "应用负责人的邮箱账号信息。用于接收上架审核结果，应用整改或下架通知。"},
			{Name: "appReviewName", In: "body", Type: "string", Required: false, Description: "应用负责人的姓名。"},
			{Name: "versionVocs", In: "body", Type: "array", Required: false, Description: "应用版本已解决问题。 数组长度不超过1000。 说明 调用查询应用未处理的问题工单接口可以查询用户反馈的应用问题。如果在应用版本更新时某些问题已经解决，可以将这些问题提交审核，工单平台验证通过后会关闭工单。", ItemsType: "object"},
		}),
		currentPublishingEndpoint("language-info-update-v3", "更新应用语言信息", "PUT", "/api/publish/v3/app-language-info", "https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-language-info-update-0000002271160569", []Parameter{
			{Name: "appId", In: "query", Type: "string", Required: true, Description: "需要更新语言描述信息的应用ID。"},
			{Name: "releaseType", In: "query", Type: "integer", Required: false, Description: "应用发布方式。 取值范围： 1：全网 默认值：1"},
			{Name: "releasePhase", In: "query", Type: "integer", Required: false, Description: "分阶段发布标识。 取值范围： 0：全网 3：分阶段 默认值：0"},
			{Name: "lang", In: "body", Type: "string", Required: true, Description: "语言种类。 说明 取值参见语言类型。"},
			{Name: "appName", In: "body", Type: "string", Required: false, Description: "对应语言的应用名称。 新增一个语言种类时，此字段必选。"},
			{Name: "appDesc", In: "body", Type: "string", Required: false, Description: "对应语言的应用描述信息。"},
			{Name: "briefInfo", In: "body", Type: "string", Required: false, Description: "对应语言的一句话简介。"},
			{Name: "newFeatures", In: "body", Type: "string", Required: false, Description: "对应语言的新版本简介。 说明 升级版本时，此参数必填。"},
		}),
		currentPublishingEndpoint("language-info-delete", "删除应用语言信息", "DELETE", "/api/publish/v2/app-language-info", "https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-language-info-delete-0000002271000621", []Parameter{
			{Name: "appId", In: "query", Type: "string", Required: true, Description: "需要删除语言描述信息的应用ID。"},
			{Name: "lang", In: "query", Type: "string", Required: true, Description: "需要删除的语言。 说明 取值参见语言类型。"},
			{Name: "releaseType", In: "query", Type: "integer", Required: false, Description: "应用发布方式。 取值范围： 1：全网 默认值：1"},
		}),
		currentPublishingEndpoint("app-file-info-v3", "更新应用文件信息", "PUT", "/api/publish/v3/app-file-info", "https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-app-file-info-update-0000002236041430", []Parameter{
			{Name: "appId", In: "query", Type: "string", Required: true, Description: "需要更新的应用ID。"},
			{Name: "releaseType", In: "query", Type: "integer", Required: false, Description: "应用发布方式。 取值范围： 1：全网 默认值：1"},
			{Name: "releasePhase", In: "query", Type: "integer", Required: false, Description: "分阶段发布标识。 取值范围： 0：全网 3：分阶段 默认值：0"},
			{Name: "appIconList", In: "body", Type: "array", Required: false, Description: "应用图标信息。 数组长度不超过20。 说明 应用图标必须与软件包中的应用图标一致，且需符合审核政策与应用文件要求。", ItemsType: "object"},
			{Name: "screenShotList", In: "body", Type: "array", Required: false, Description: "应用介绍截图信息。 数组长度不超过20。 说明 在同一语言、同一设备下，同时更新的截图素材和视频素材方向必须保持一致，否则将报错。", ItemsType: "object"},
			{Name: "introVideoList", In: "body", Type: "array", Required: false, Description: "应用介绍视频信息。 数组长度不超过20。 说明 在同一语言、同一设备下，同时更新的截图素材和视频素材方向必须保持一致，否则将报错。", ItemsType: "object"},
			{Name: "rcmdVideoList", In: "body", Type: "array", Required: false, Description: "推荐视频信息。 数组长度不超过20。", ItemsType: "object"},
			{Name: "rcmdPicList", In: "body", Type: "array", Required: false, Description: "推荐图片信息。 数组长度不超过20。", ItemsType: "object"},
			{Name: "elecCertificate", In: "body", Type: "string", Required: false, Description: "电子版本证书。"},
			{Name: "certificateList", In: "body", Type: "array", Required: false, Description: "应用版权证书或代理证书。 数组长度不超过10。", ItemsType: "string"},
			{Name: "publicationList", In: "body", Type: "array", Required: false, Description: "版号证明。 数组长度不超过10。", ItemsType: "string"},
			{Name: "packagePermissionIntroVideoList", In: "body", Type: "array", Required: false, Description: "包权限使用场景说明视频。 数组长度不超过100。 说明 隐私权限说明的“权限描述”、“权限使用设备范围”和“权限使用理由”从上传的软件包体中读取，无需通过接口更新。", ItemsType: "object"},
			{Name: "otherCertificatesURL", In: "body", Type: "string", Required: false, Description: "授权书及其他材料文件压缩包的objectId。 说明 objectId的获取流程：先调用获取上传文件地址接口，然后调用上传文件接口上传文件，最后得到获取上传文件地址接口返回的objectId。 若文件格式不正确，提交发布时会出错，具体限制如下： 压缩包必须为zip格式，不超过200MB。 解压后的大小不超过1GB。 文件夹层级不超过3层。 只能包含JPG、JPEG、BMP格式图片和PDF格式的文件，累计不超过200个。 说明 此字段只支持HarmonyOS应用，不支持元服务。"},
			{Name: "updateScreenShotList", In: "body", Type: "array", Required: false, Description: "更新截图素材。 数组长度不超过20。", ItemsType: "object"},
			{Name: "updateVideoList", In: "body", Type: "array", Required: false, Description: "更新视频素材。 数组长度不超过20。", ItemsType: "object"},
		}),
		currentPublishingEndpoint("app-package-info-update", "更新应用软件包信息", "PUT", "/api/publish/v3/app-package-info", "https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-app-package-info-update-0000002236201250", []Parameter{
			{Name: "appId", In: "query", Type: "string", Required: true, Description: "需要更新的应用ID。"},
			{Name: "releaseType", In: "query", Type: "integer", Required: false, Description: "应用发布方式。 取值范围： 1：全网 默认值：1"},
			{Name: "releasePhase", In: "query", Type: "integer", Required: false, Description: "分阶段发布标识。 取值范围： 0：全网 3：分阶段 默认值：0"},
			{Name: "fileName", In: "body", Type: "string", Required: true, Description: "文件名称，包括文件的后缀名。"},
			{Name: "objectId", In: "body", Type: "string", Required: true, Description: "文件在文件服务器中的对象ID。 说明 objectId的获取流程：先调用获取上传文件地址接口，然后调用上传文件接口上传文件，最后得到获取上传文件地址接口返回的objectId。"},
		}),
		currentPublishingEndpoint("update-releasetime", "更新版本上架时间", "PUT", "/api/publish/v2/on-shelf-time", "https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-update-release-time-0000002271160573", []Parameter{
			{Name: "appId", In: "query", Type: "string", Required: true, Description: "应用ID。"},
			{Name: "changeType", In: "body", Type: "integer", Required: true, Description: "修改类型。 取值范围： 2：指定时间上架改为审核通过立即上架 3：变更定时上架时间"},
			{Name: "releaseTime", In: "body", Type: "string", Required: true, Description: "指定的版本上架时间。 UTC时间格式：yyyy-MM-ddTHH:mm:ssZZ，例如“2015-01-01T01:01:01+0800”。 changeType为3时有效。"},
			{Name: "releaseType", In: "body", Type: "integer", Required: true, Description: "发布方式。 取值范围： 1：全网 目前支持全网及releasePhase为3的分阶段发布。"},
		}),
		currentPublishingEndpoint("phased-release-update", "更新分阶段发布", "PUT", "/api/publish/v2/version/phased-release", "https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-phased-release-0000002271000625", []Parameter{
			{Name: "appId", In: "header", Type: "string", Required: true, Description: "应用ID，获取方法参考查看应用信息。"},
			{Name: "versionId", In: "body", Type: "string", Required: true, Description: "版本ID。 可从查询应用所有版本类型的版本列表接口返回的versionId参数中获取。"},
			{Name: "releasePhase", In: "body", Type: "integer", Required: false, Description: "分阶段发布标识。 取值范围： 3：分阶段。 4：分阶段转全网。"},
			{Name: "state", In: "body", Type: "string", Required: false, Description: "分阶段发布的状态。 取值范围： SUSPEND：暂停。 RELEASE：恢复为发布中。"},
			{Name: "description", In: "body", Type: "string", Required: false, Description: "分阶段发布说明。"},
			{Name: "phaseDay", In: "body", Type: "integer", Required: false, Description: "更新分阶段天数。 取值范围2、3、4、5、6、7，对应分阶段发布比例加速到2%、5%、10%、20%、50%、100%。"},
		}),
		currentPublishingEndpoint("app-package-info-query", "查询应用软件包信息", "GET", "/api/publish/v2/app-package-info", "https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-query-package-info-0000002328848469", []Parameter{
			{Name: "appId", In: "header", Type: "string", Required: true, Description: "应用ID，获取方法参考查看应用信息。"},
			{Name: "packageId", In: "query", Type: "string", Required: true, Description: "待查询的软件包ID。"},
		}),
		currentPublishingEndpoint("package-compile-status-v3", "查询软件包编译状态", "GET", "/api/publish/v3/package/compile/status", "https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-query-compile-status-0000002236041434", []Parameter{
			{Name: "appId", In: "query", Type: "string", Required: true, Description: "应用ID。"},
			{Name: "pkgIds", In: "query", Type: "string", Required: true, Description: "待查询的软件包ID。 多个ID之间用英文逗号分隔。 可从更新应用软件包信息接口返回的packageId参数中获取。"},
		}),
		currentPublishingEndpoint("app-version-list", "查询应用所有版本类型的版本列表", "POST", "/api/publish/v3/version/brief-info/list", "https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-query-brief-info-list-0000002236201254", []Parameter{
			{Name: "appId", In: "header", Type: "string", Required: true, Description: "应用ID，获取方法参考查看应用信息。"},
			{Name: "packageName", In: "body", Type: "string", Required: false, Description: "应用包名。 说明 appId和packageName至少传一个。"},
			{Name: "state", In: "body", Type: "string", Required: false, Description: "版本状态。 多个状态用逗号分隔。 说明 state字段的取值请参考AppInfo中的releaseState参数。"},
		}),
		currentPublishingEndpoint("app-version-voc", "查询应用未处理的问题工单", "POST", "/api/publish/v1/app-voc/version-voc", "https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-post-app-version-voc-0000002495233986", []Parameter{
			{Name: "appId", In: "header", Type: "string", Required: true, Description: "应用ID，获取方法参考查看应用信息。"},
			{Name: "vocId", In: "body", Type: "string", Required: false, Description: "工单ID。"},
			{Name: "vocDesc", In: "body", Type: "string", Required: false, Description: "工单描述。"},
			{Name: "pageNo", In: "body", Type: "integer", Required: true, Description: "当前页。 最小值：1"},
			{Name: "pageSize", In: "body", Type: "integer", Required: true, Description: "每页数量。 最大值：100"},
		}),
		currentPublishingEndpoint("tag-qual-dependent", "查询需要上传的资质文件信息", "GET", "/api/publish/v2/tag-qual-dependent", "https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-query-tag-qual-dependent-0000002271160577", []Parameter{}),
		currentPublishingEndpoint("qualification-review-create", "提交资质审核任务", "POST", "/api/publish/v2/qualification-review-task", "https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-post-qualification-review-0000002271000629", []Parameter{
			{Name: "appId", In: "header", Type: "string", Required: true, Description: "应用ID，获取方法参考查看应用信息。"},
			{Name: "tagQualInfos", In: "body", Type: "array", Required: true, Description: "标签和标签对应的资质信息列表。 数组长度不超过5。", ItemsType: "object"},
		}),
		currentPublishingEndpoint("qualification-review-query", "查询资质审核任务详情", "GET", "/api/publish/v2/qualification-review-task", "https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-get-qualification-review-0000002236041438", []Parameter{
			{Name: "appId", In: "header", Type: "string", Required: true, Description: "应用ID，获取方法参考查看应用信息。"},
			{Name: "tagIds", In: "query", Type: "string", Required: true, Description: "审核任务包含的标签ID列表，逗号分隔。"},
		}),
		currentPublishingEndpoint("qualification-review-update", "修改资质审核任务", "PUT", "/api/publish/v2/qualification-review-task", "https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-put-qualification-review-0000002422530746", []Parameter{
			{Name: "appId", In: "header", Type: "string", Required: true, Description: "应用ID，获取方法参考查看应用信息。"},
			{Name: "tagQualInfos", In: "body", Type: "array", Required: true, Description: "标签和标签对应的资质信息列表。 数组长度不超过5。", ItemsType: "object"},
		}),
		currentPublishingEndpoint("qualification-review-cancel", "撤销资质审核任务", "PUT", "/api/publish/v2/qualification-review-task/cancel", "https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-cancel-qualification-review-0000002236201258", []Parameter{
			{Name: "appId", In: "header", Type: "string", Required: true, Description: "应用ID，获取方法参考查看应用信息。"},
			{Name: "taskId", In: "body", Type: "string", Required: true, Description: "资质审核任务ID。"},
		}),
		currentPublishingEndpoint("qualification-review-delete", "删除资质审核任务", "DELETE", "/api/publish/v2/qualification-review-task", "https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-delete-qualification-review-0000002526538788", []Parameter{
			{Name: "appId", In: "header", Type: "string", Required: true, Description: "应用ID，获取方法参考查看应用信息。"},
			{Name: "tagIds", In: "body", Type: "array", Required: true, Description: "需要删除的资质审核标签ID列表，逗号分隔。 数组长度不超过5。 说明 具体取值请参见元服务的“标签ID”列。", ItemsType: "string"},
		}),
		currentPublishingEndpoint("user-agreement-create", "新建用户协议", "POST", "/api/publish/v2/agreement", "https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-create-user-agreement-0000002271160581", []Parameter{
			{Name: "appId", In: "header", Type: "string", Required: true, Description: "应用ID，获取方法参考查看应用信息。"},
			{Name: "type", In: "body", Type: "integer", Required: true, Description: "协议类型。 取值范围： 2：用户协议"},
			{Name: "defaultLang", In: "body", Type: "string", Required: false, Description: "用户协议默认语言。 取值范围： zh：简体中文 en：英语 ug：维吾尔语 bo：藏语 zh-Hant-HK：繁体中文（中国香港特别行政区） zh-Hant-TW：繁体中文（中国台湾）"},
			{Name: "agreementInfos", In: "body", Type: "array", ItemsType: "object", MaxItems: 100, Required: false, Description: "用户协议列表；以官方请求样例的数组格式为准，最多100项。"},
			{Name: "channel", In: "body", Type: "integer", Required: false, Description: "官方创建用户协议请求样例中的渠道字段。"},
		}),
		currentPublishingEndpoint("user-agreement-list", "查询用户协议列表", "POST", "/api/publish/v2/agreement/list", "https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-query-user-agreement-0000002236041442", []Parameter{
			{Name: "appId", In: "header", Type: "string", Required: true, Description: "应用ID，获取方法参考查看应用信息。"},
			{Name: "agreementIdList", In: "body", Type: "array", Required: false, Description: "用户协议ID。 数组长度为1~1000。", ItemsType: "string"},
			{Name: "name", In: "body", Type: "string", Required: true, Description: "用户协议名称。 使用模糊查询方式。", AllowEmpty: true},
			{Name: "status", In: "body", Type: "integer", Required: false, Description: "用户协议状态。 取值范围： 2：已完成 3：已删除"},
			{Name: "type", In: "body", Type: "integer", Required: false, Description: "协议类型。 取值范围： 2：用户协议"},
			{Name: "fromRecCount", In: "body", Type: "integer", Required: false, Description: "查询起始位置。 默认值：1"},
			{Name: "maxReqCount", In: "body", Type: "integer", Required: false, Description: "最大请求查询数。 取值范围：[1,100] 默认值：30"},
		}),
		currentPublishingEndpoint("user-agreement-update", "更新用户协议", "PUT", "/api/publish/v2/agreement", "https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-update-user-agreement-0000002294845626", []Parameter{
			{Name: "appId", In: "header", Type: "string", Required: true, Description: "应用ID，获取方法参考查看应用信息。"},
			{Name: "id", In: "body", Type: "string", Required: true, Description: "用户协议ID。"},
			{Name: "defaultLang", In: "body", Type: "string", Required: false, Description: "用户协议默认语言。 取值范围： zh：简体中文 en：英语 ug：维吾尔语 bo：藏语 zh-Hant-HK：繁体中文（中国香港特别行政区） zh-Hant-TW：繁体中文（中国台湾）"},
			{Name: "agreementInfos", In: "body", Type: "array", Required: false, Description: "用户协议列表。 数组长度不超过100。", ItemsType: "object"},
		}),
		currentPublishingEndpoint("privacy-agreement-submit", "提交隐私政策协议", "PUT", "/api/publish/v2/agreement/submit", "https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-put-privacy-agreement-0000002271000633", []Parameter{
			{Name: "appId", In: "header", Type: "string", Required: true, Description: "应用ID，获取方法参考查看应用信息。"},
			{Name: "id", In: "body", Type: "string", Required: true, Description: "隐私政策协议ID。"},
			{Name: "requireConfirm", In: "body", Type: "integer", Required: false, Description: "是否弹窗确认。 取值范围： 1：是 0：否"},
		}),
		currentPublishingEndpoint("privacy-agreement-list", "查询隐私政策协议列表", "POST", "/api/publish/v2/agreement/list", "https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-query-privacy-agreement-0000002328924925", []Parameter{
			{Name: "appId", In: "header", Type: "string", Required: true, Description: "应用ID，获取方法参考查看应用信息。"},
			{Name: "agreementIdList", In: "body", Type: "array", Required: false, Description: "隐私政策协议ID。 数组长度为1~1000。", ItemsType: "string"},
			{Name: "name", In: "body", Type: "string", Required: true, Description: "隐私政策协议名称。 使用模糊查询方式。", AllowEmpty: true},
			{Name: "status", In: "body", Type: "integer", Required: false, Description: "隐私政策协议状态。 取值范围： 1：草稿 2：已完成 3：已删除 不传值：默认查未删除的记录"},
			{Name: "type", In: "body", Type: "integer", Required: false, Description: "协议类型。 取值范围： 1：隐私政策协议"},
			{Name: "fromRecCount", In: "body", Type: "integer", Required: false, Description: "起始页数。 默认值：1"},
			{Name: "maxReqCount", In: "body", Type: "integer", Required: false, Description: "每页返回条数。 取值范围：[1,100] 默认值：30"},
		}),
		currentPublishingEndpoint("privacy-agreement-update", "更新隐私政策协议", "PUT", "/api/publish/v2/agreement", "https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-update-privacy-agreement-0000002328805169", []Parameter{
			{Name: "appId", In: "header", Type: "string", Required: true, Description: "应用ID，获取方法参考查看应用信息。"},
			{Name: "id", In: "body", Type: "string", Required: true, Description: "隐私政策协议ID。"},
			{Name: "defaultLang", In: "body", Type: "string", Required: false, Description: "隐私政策协议默认语言。 取值范围： zh：简体中文 en：英语 ug：维吾尔语 bo：藏语 zh-Hant-HK：繁体中文（中国香港特别行政区） zh-Hant-TW：繁体中文（中国台湾）"},
			{Name: "agreementInfos", In: "body", Type: "array", Required: false, Description: "隐私政策协议列表。 数组长度不超过100。", ItemsType: "object"},
		}),
		currentPublishingEndpoint("app-age-rating", "申请年龄分级", "POST", "/api/publish/v2/app-age-rating", "https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-post-app-age-rating-0000002236201262", []Parameter{
			{Name: "appId", In: "header", Type: "string", Required: true, Description: "应用ID。 获取方法参考查看应用信息。"},
			{Name: "releaseType", In: "query", Type: "integer", Required: false, Description: "应用发布方式。 取值范围： 1：全网 默认值：1"},
			{Name: "releasePhase", In: "query", Type: "integer", Required: false, Description: "分阶段发布标识。 取值范围： 0：全网 3：分阶段 默认值：0"},
			{Name: "feedback", In: "body", Type: "string", Required: true, Description: "问卷内容。 说明 具体请参见feedback获取方式。"},
			{Name: "expectedRating", In: "body", Type: "integer", Required: false, Description: "预期年龄分级。 取值范围： 1：巴西适用所有年龄段，中国区及其他区域年满3周岁 2：中国区年满8周岁，巴西年满10周岁，其他区域年满7周岁 3：巴西年满14周岁，中国区及其他区域年满12周岁 4：年满16周岁 5：年满18周岁"},
			{Name: "ratingRule", In: "body", Type: "integer", Required: false, Description: "年龄分级生效规则。 取值范围： 0：以问卷内容评级结果为准。 1：当预期年龄分级expectedRating设置的值大于问卷评级结果中的年龄值时，以预期年龄分级为准，否则以问卷内容评级结果为准。"},
			{Name: "childFlag", In: "body", Type: "integer", Required: false, Description: "应用是否仅面向儿童。 取值范围： 0：应用不仅面向儿童 1：应用仅面向儿童"},
		}),
		currentPublishingEndpoint("app-submit-v3", "提交发布", "POST", "/api/publish/v3/app-submit", "https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-app-submit-0000002271160585", []Parameter{
			{Name: "appId", In: "query", Type: "string", Required: true, Description: "应用ID。"},
			{Name: "releaseTime", In: "body", Type: "string", Required: false, Description: "指定发布时间。 UTC时间格式：yyyy-MM-ddTHH:mm:ssZZ，例如“2015-01-01T01:01:01+0800”。 可以不用填写，则表示立刻发布。 说明 如果在首次提交发布时已设置过releaseTime参数，撤销审核后再次提交发布时，需要将releaseTime参数重新设置为过去的时间，才能实现立即发布。"},
			{Name: "remark", In: "body", Type: "string", Required: false, Description: "提审发布备注。 可以为空，如果填写则长度范围为10-300。"},
			{Name: "releaseType", In: "body", Type: "integer", Required: false, Description: "应用发布方式。 取值范围： 1：全网 默认值：1"},
			{Name: "releasePhase", In: "body", Type: "integer", Required: false, Description: "分阶段发布标识。 取值范围： 0：全网发布 3：分阶段发布（7天内自动更新） 默认值：0"},
			{Name: "phasedReleaseDescription", In: "body", Type: "string", Required: false, Description: "分阶段发布说明。 说明 releasePhase取值为3（分阶段发布）时，此参数必填。"},
		}),
		currentPublishingEndpoint("on-shelf-cancel", "撤销审核", "PUT", "/api/publish/v3/version/on-shelf/cancel", "https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-on-shelf-cancel-0000002271000637", []Parameter{
			{Name: "appId", In: "header", Type: "string", Required: true, Description: "应用ID，获取方法参考查看应用信息。"},
			{Name: "versionId", In: "body", Type: "string", Required: true, Description: "版本ID。"},
		}),
		currentPublishingEndpoint("app-shelf-off", "下架应用", "POST", "/api/publish/v2/app-shelf-off", "https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-app-shelf-off-0000002236041450", []Parameter{
			{Name: "appId", In: "query", Type: "string", Required: true, Description: "需要下架的应用ID。"},
			{Name: "reason", In: "query", Type: "string", Required: true, Description: "下架的原因描述。 说明 中文字符需要进行URL编码后再传入。"},
			{Name: "releaseType", In: "query", Type: "integer", Required: false, Description: "下架的方式。 取值范围： 1：全网下架。 默认值：1"},
		}),
	}
}

func currentPublishingEndpoint(id, name, method, path, sourceURL string, parameters []Parameter) Endpoint {
	e := officialEndpoint("publishing", id, name, name, method, path, nil, sourceURL[strings.LastIndex(sourceURL, "/")+1:], sourceURL, "developer-to-huawei")
	e.Parameters = parameters
	e.RequiredParams = nil
	for _, p := range parameters {
		if p.Required {
			e.RequiredParams = append(e.RequiredParams, p.Name)
		}
	}
	e.Body = inferBody(method, parameters)
	e = applyReferenceParameters(e)
	e = applyPublishingAgreementModels(e)
	e = applyPublishingMetadataModels(e)
	return applyPublishingQualificationModels(e)
}

func init() {
	for _, current := range currentPublishingEndpoints() {
		replaced := false
		for i, existing := range EndpointFamilies["publishing"] {
			if existing.ID == current.ID {
				EndpointFamilies["publishing"][i] = current
				replaced = true
				break
			}
		}
		if !replaced {
			EndpointFamilies["publishing"] = append(EndpointFamilies["publishing"], current)
		}
	}
}
