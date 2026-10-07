package domain

var EndpointFamilies = map[string][]Endpoint{
	"publishing": {
		// 创建专用 HarmonyOS 测试应用：省略 projectId 时由华为建立同名项目。
		// 请求：POST /api/publish/v3/app；调用方向：开发者调用华为服务端。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-provision-api-create-appid-0000002526543939
		officialEndpoint("publishing", "app-create", "创建APP ID", "创建 HarmonyOS 应用或元服务，用于建立与现有业务应用隔离的测试资源。", "POST", "/api/publish/v3/app", []string{}, "agc-help-provision-api-create-appid-0000002526543939", "https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-provision-api-create-appid-0000002526543939", "developer-to-huawei"),
		// 查询 HarmonyOS V3 应用或测试版本详情，不替换兼容旧接口的 V2 命令。
		// 请求：GET /api/publish/v3/app-info；调用方向：开发者调用华为服务端。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-appinfo-query-0000002236041422
		officialEndpoint("publishing", "app-info-query-v3", "查询HarmonyOS应用信息", "查询 HarmonyOS 应用和元服务的资料；releaseType=6 时必须提供专用测试版本 ID。", "GET", "/api/publish/v3/app-info", []string{}, "agc-help-publish-api-appinfo-query-0000002236041422", "https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-appinfo-query-0000002236041422", "developer-to-huawei"),
		// 查询所有商用和测试版本：用于确认专用测试草稿创建及清理结果。
		// 请求：POST /api/publish/v3/version/brief-info/list；调用方向：开发者调用华为服务端。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-query-brief-info-list-0000002236201254
		officialEndpoint("publishing", "app-version-list", "查询应用所有版本类型的版本列表", "查询 HarmonyOS 应用或元服务的商用版本和测试版本，读回版本 ID 与状态。", "POST", "/api/publish/v3/version/brief-info/list", []string{}, "agc-help-publish-api-query-brief-info-list-0000002236201254", "https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-query-brief-info-list-0000002236201254", "developer-to-huawei"),
		// 通过下载方式提交软件包：如果您将软件包保存在自己的服务器上，可以直接调用本接口将软件包提交到AppGallery Connect服务器。
		// 请求：POST /api/publish/v2/app-package-file/by-url；调用方向：开发者调用华为服务端。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-add-packageurl-0000001158245065
		officialEndpoint("publishing", "add-packageurl", "通过下载方式提交软件包", "如果您将软件包保存在自己的服务器上，可以直接调用本接口将软件包提交到AppGallery Connect服务器。", "POST", "/api/publish/v2/app-package-file/by-url", []string{}, "agcapi-add-packageurl-0000001158245065", "https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-add-packageurl-0000001158245065", "developer-to-huawei"),
		// 更新应用文件信息：图片、视频、APK、RPK、AAB等文件上传完成后，通过该接口刷新应用文件信息。
		// 请求：PUT /api/publish/v2/app-file-info；调用方向：开发者调用华为服务端。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-app-file-info-0000001111685202
		officialEndpoint("publishing", "app-file-info", "更新应用文件信息", "图片、视频、APK、RPK、AAB等文件上传完成后，通过该接口刷新应用文件信息。", "PUT", "/api/publish/v2/app-file-info", []string{}, "agcapi-app-file-info-0000001111685202", "https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-app-file-info-0000001111685202", "developer-to-huawei"),
		// 查询应用信息：此接口用于根据指定语言和应用ID查询应用的详细信息。
		// 请求：GET /api/publish/v2/app-info；调用方向：开发者调用华为服务端。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-app-info-query-0000001158365045
		officialEndpoint("publishing", "app-info-query", "查询应用信息", "此接口用于根据指定语言和应用ID查询应用的详细信息。", "GET", "/api/publish/v2/app-info", []string{}, "agcapi-app-info-query-0000001158365045", "https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-app-info-query-0000001158365045", "developer-to-huawei"),
		// 更新应用基本信息：此接口用于更新指定应用的应用详情。
		// 请求：PUT /api/publish/v2/app-info；调用方向：开发者调用华为服务端。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-app-info-update-0000001111685198
		officialEndpoint("publishing", "app-info-update", "更新应用基本信息", "此接口用于更新指定应用的应用详情。", "PUT", "/api/publish/v2/app-info", []string{}, "agcapi-app-info-update-0000001111685198", "https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-app-info-update-0000001111685198", "developer-to-huawei"),
		// 提交发布：此接口用于提交应用审核，调用本接口前必须保证应用信息已补充完整，应用软件包已经上传。
		// 请求：POST /api/publish/v2/app-submit；调用方向：开发者调用华为服务端。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-app-submit-0000001158245061
		officialEndpoint("publishing", "app-submit", "提交发布", "此接口用于提交应用审核，调用本接口前必须保证应用信息已补充完整，应用软件包已经上传。", "POST", "/api/publish/v2/app-submit", []string{}, "agcapi-app-submit-0000001158245061", "https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-app-submit-0000001158245061", "developer-to-huawei"),
		// 通过下载方式提交发布：如果您将软件包保存在自己的服务器上，可以直接调用本接口提交应用审核，并传入软件包的下载地址。
		// 请求：POST /api/publish/v2/app-submit-with-file；调用方向：开发者调用华为服务端。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-app-submit-with-file-0000001111845092
		officialEndpoint("publishing", "app-submit-with-file", "通过下载方式提交发布", "如果您将软件包保存在自己的服务器上，可以直接调用本接口提交应用审核，并传入软件包的下载地址。", "POST", "/api/publish/v2/app-submit-with-file", []string{}, "agcapi-app-submit-with-file-0000001111845092", "https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-app-submit-with-file-0000001111845092", "developer-to-huawei"),
		// 查询应用包名对应的appid：此接口用于根据应用包名查询对应的应用ID，支持一次查询多个应用包名。
		// 请求：GET /api/publish/v2/appid-list；调用方向：开发者调用华为服务端。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-appid-list-0000001111845086
		// 现行文档（已核对登记路径）：https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-appid-list-0000002271000617
		officialEndpoint("publishing", "appid-list", "查询应用包名对应的appid", "此接口用于根据应用包名查询对应的应用ID，支持一次查询多个应用包名。", "GET", "/api/publish/v2/appid-list", []string{}, "agcapi-appid-list-0000001111845086", "https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-appid-list-0000001111845086", "developer-to-huawei"),
		// 设置应用GMS依赖属性：此接口用于向AGC上报应用是否依赖GMS。
		// 请求：PUT /api/publish/v2/properties/gms；调用方向：开发者调用华为服务端。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-gms-0000001111845094
		officialEndpoint("publishing", "gms", "设置应用GMS依赖属性", "此接口用于向AGC上报应用是否依赖GMS。", "PUT", "/api/publish/v2/properties/gms", []string{}, "agcapi-gms-0000001111845094", "https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-gms-0000001111845094", "developer-to-huawei"),
		// 删除语言描述信息：此接口用于删除指定应用不需要的语言描述信息，默认语言不支持删除。
		// 请求：DELETE /api/publish/v2/app-language-info；调用方向：开发者调用华为服务端。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-language-info-delete-0000001111845088
		officialEndpoint("publishing", "language-info-delete", "删除语言描述信息", "此接口用于删除指定应用不需要的语言描述信息，默认语言不支持删除。", "DELETE", "/api/publish/v2/app-language-info", []string{}, "agcapi-language-info-delete-0000001111845088", "https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-language-info-delete-0000001111845088", "developer-to-huawei"),
		// 更新语言描述信息：此接口用于更新指定应用的语言描述信息，如果当前没有语言描述信息则新增。
		// 请求：PUT /api/publish/v2/app-language-info；调用方向：开发者调用华为服务端。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-language-info-update-0000001158245057
		officialEndpoint("publishing", "language-info-update", "更新语言描述信息", "此接口用于更新指定应用的语言描述信息，如果当前没有语言描述信息则新增。", "PUT", "/api/publish/v2/app-language-info", []string{}, "agcapi-language-info-update-0000001158245057", "https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-language-info-update-0000001158245057", "developer-to-huawei"),
		// 查询软件包列表：此接口用于查询应用当前关联的软件包列表信息。
		// 请求：GET /api/publish/v2/package-list；调用方向：开发者调用华为服务端。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-obbfile-package-list-0000001111685210
		officialEndpoint("publishing", "obbfile-package-list", "查询软件包列表", "此接口用于查询应用当前关联的软件包列表信息。", "GET", "/api/publish/v2/package-list", []string{}, "agcapi-obbfile-package-list-0000001111685210", "https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-obbfile-package-list-0000001111685210", "developer-to-huawei"),
		// 更新分阶段发布：此接口用于将分阶段发布修改为全网发布，或者更新分阶段发布的设置。
		// 请求：PUT /api/publish/v2/phased-release；调用方向：开发者调用华为服务端。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-phased-release-0000001111685204
		officialEndpoint("publishing", "phased-release", "更新分阶段发布", "此接口用于将分阶段发布修改为全网发布，或者更新分阶段发布的设置。", "PUT", "/api/publish/v2/phased-release", []string{}, "agcapi-phased-release-0000001111685204", "https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-phased-release-0000001111685204", "developer-to-huawei"),
		// 查询软件包编译状态：此接口用于查询软件包的编译状态。
		// 请求：GET /api/publish/v2/package/compile/status；调用方向：开发者调用华为服务端。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-query-aabfile-0000001111685206
		officialEndpoint("publishing", "query-aabfile", "查询软件包编译状态", "此接口用于查询软件包的编译状态。", "GET", "/api/publish/v2/package/compile/status", []string{}, "agcapi-query-aabfile-0000001111685206", "https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-query-aabfile-0000001111685206", "developer-to-huawei"),
		// 更新版本上架时间：此接口用于设置版本的上架时间。
		// 请求：PUT /api/publish/v2/on-shelf-time；调用方向：开发者调用华为服务端。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-update-releasetime-0000001158365053
		// 现行文档（已核对登记路径）：https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-update-release-time-0000002271160573
		officialEndpoint("publishing", "update-releasetime", "更新版本上架时间", "此接口用于设置版本的上架时间。", "PUT", "/api/publish/v2/on-shelf-time", []string{}, "agcapi-update-releasetime-0000001158365053", "https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-update-releasetime-0000001158365053", "developer-to-huawei"),
	},
	"upload": {
		// 合并分片：此接口用于将已经上传成功后的分片合并为一个文件。
		// 请求：POST /api/publish/v2/upload/multipart/compose；调用方向：开发者调用华为服务端。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-obbfile-compose-0000001111845098
		// 现行文档（已核对登记路径）：https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-upload-api-obbfile-compose-0000002271160625
		officialEndpoint("upload", "obbfile-compose", "合并分片", "此接口用于将已经上传成功后的分片合并为一个文件。", "POST", "/api/publish/v2/upload/multipart/compose", []string{}, "agcapi-obbfile-compose-0000001111845098", "https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-obbfile-compose-0000001111845098", "developer-to-huawei"),
		// 分片上传初始化：此接口用于分片信息初始化，在获取分片文件上传地址前调用。
		// 请求：POST /api/publish/v2/upload/multipart/init；调用方向：开发者调用华为服务端。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-obbfile-init-0000001158365055
		// 现行文档（已核对登记路径）：https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-upload-api-obbfile-init-0000002271000677
		officialEndpoint("upload", "obbfile-init", "分片上传初始化", "此接口用于分片信息初始化，在获取分片文件上传地址前调用。", "POST", "/api/publish/v2/upload/multipart/init", []string{}, "agcapi-obbfile-init-0000001158365055", "https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-obbfile-init-0000001158365055", "developer-to-huawei"),
		// 上传分片实体：将单个分片上传到获取分片上传地址接口返回的上传地址。
		// 请求：PUT {uploadUrl}；调用方向：开发者调用华为服务端。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-obbfile-upload-0000001158245067
		officialEndpoint("upload", "obbfile-upload", "上传分片实体", "将单个分片上传到获取分片上传地址接口返回的上传地址。", "PUT", "{uploadUrl}", []string{"uploadUrl"}, "agcapi-obbfile-upload-0000001158245067", "https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-obbfile-upload-0000001158245067", "developer-to-huawei"),
		// 获取分片上传地址：此接口用于获取文件所有分片的上传地址，调用此接口前必须先调用 分片上传初始化 接口完成分片初始化。
		// 请求：POST /api/publish/v2/upload/multipart/parts；调用方向：开发者调用华为服务端。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-obbfile-uploadurl-0000001111685208
		// 现行文档（已核对登记路径）：https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-upload-api-obbfile-uploadurl-0000002236041490
		officialEndpoint("upload", "obbfile-uploadurl", "获取分片上传地址", "此接口用于获取文件所有分片的上传地址，调用此接口前必须先调用 分片上传初始化 接口完成分片初始化。", "POST", "/api/publish/v2/upload/multipart/parts", []string{}, "agcapi-obbfile-uploadurl-0000001111685208", "https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-obbfile-uploadurl-0000001111685208", "developer-to-huawei"),
		// 上传文件：将文件上传到获取文件上传地址接口返回的上传地址。
		// 请求：PUT {uploadUrl}；调用方向：开发者调用华为服务端。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-upload-file-new-0000001111845090
		officialEndpoint("upload", "upload-file-new", "上传文件", "将文件上传到获取文件上传地址接口返回的上传地址。", "PUT", "{uploadUrl}", []string{"uploadUrl"}, "agcapi-upload-file-new-0000001111845090", "https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-upload-file-new-0000001111845090", "developer-to-huawei"),
		// 获取文件上传地址：上传文件至AppGallery Connect前需要先调用此接口获取上传地址，包括图片、视频、PDF等文件。
		// 请求：GET /api/publish/v2/upload-url/for-obs；调用方向：开发者调用华为服务端。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-upload-url-new-0000001111685200
		// 现行文档（已核对登记路径）：https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-upload-api-upload-url-0000002236201294
		officialEndpoint("upload", "upload-url-new", "获取文件上传地址", "上传文件至AppGallery Connect前需要先调用此接口获取上传地址，包括图片、视频、PDF等文件。", "GET", "/api/publish/v2/upload-url/for-obs", []string{}, "agcapi-upload-url-new-0000001111685200", "https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-upload-url-new-0000001111685200", "developer-to-huawei"),
	},
	"provisioning": {
		// 批量添加设备：此接口用于批量添加调试/测试设备。
		// 请求：POST /api/publish/v2/device；调用方向：开发者调用华为服务端。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/app/agc-help-provision-api-add-device-0000002236041498
		// 现行文档（已核对登记路径）：https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-provision-api-add-device-0000002236041498
		officialEndpoint("provisioning", "provision-api-add-device", "批量添加设备", "此接口用于批量添加调试/测试设备。", "POST", "/api/publish/v2/device", []string{}, "agc-help-provision-api-add-device-0000002236041498", "https://developer.huawei.com/consumer/cn/doc/app/agc-help-provision-api-add-device-0000002236041498", "developer-to-huawei"),
		// 新增证书指纹：此接口用于新增证书指纹。
		// 请求：POST /api/provision/v1/fingerprints；调用方向：开发者调用华为服务端。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/app/agc-help-provision-api-add-fingerprints-0000002431201480
		// 现行文档（已核对登记路径）：https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-provision-api-add-fingerprints-0000002431201480
		officialEndpoint("provisioning", "provision-api-add-fingerprints", "新增证书指纹", "此接口用于新增证书指纹。", "POST", "/api/provision/v1/fingerprints", []string{}, "agc-help-provision-api-add-fingerprints-0000002431201480", "https://developer.huawei.com/consumer/cn/doc/app/agc-help-provision-api-add-fingerprints-0000002431201480", "developer-to-huawei"),
		// 申请证书：此接口用于申请二进制证书、In-house发布证书、调试证书或发布证书。
		// 请求：POST /api/publish/v3/cert；调用方向：开发者调用华为服务端。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/app/agc-help-provision-api-apply-cent-0000002236201302
		// 现行文档（已核对登记路径）：https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-provision-api-apply-cent-0000002236201302
		officialEndpoint("provisioning", "provision-api-apply-cent", "申请证书", "此接口用于申请二进制证书、In-house发布证书、调试证书或发布证书。", "POST", "/api/publish/v3/cert", []string{}, "agc-help-provision-api-apply-cent-0000002236201302", "https://developer.huawei.com/consumer/cn/doc/app/agc-help-provision-api-apply-cent-0000002236201302", "developer-to-huawei"),
		// 申请Profile：此接口用于申请指定设备发布Profile、In-house发布Profile、调试Profile或发布Profile。
		// 请求：POST /api/publish/v3/provision；调用方向：开发者调用华为服务端。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/app/agc-help-provision-api-apply-provision-0000002271000689
		// 现行文档（已核对登记路径）：https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-provision-api-apply-provision-0000002271000689
		officialEndpoint("provisioning", "provision-api-apply-provision", "申请Profile", "此接口用于申请指定设备发布Profile、In-house发布Profile、调试Profile或发布Profile。", "POST", "/api/publish/v3/provision", []string{}, "agc-help-provision-api-apply-provision-0000002271000689", "https://developer.huawei.com/consumer/cn/doc/app/agc-help-provision-api-apply-provision-0000002271000689", "developer-to-huawei"),
		// 申请受限ACL权限：此接口用于申请受限ACL权限。
		// 请求：POST /api/provision/v1/user/permission/apply；调用方向：开发者调用华为服务端。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/app/agc-help-provision-api-applyacl-0000002502198721
		// 现行文档（已核对登记路径）：https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-provision-api-applyacl-0000002502198721
		officialEndpoint("provisioning", "provision-api-applyacl", "申请受限ACL权限", "此接口用于申请受限ACL权限。", "POST", "/api/provision/v1/user/permission/apply", []string{}, "agc-help-provision-api-applyacl-0000002502198721", "https://developer.huawei.com/consumer/cn/doc/app/agc-help-provision-api-applyacl-0000002502198721", "developer-to-huawei"),
		// 删除证书：此接口用于删除创建的证书。
		// 请求：POST /api/publish/v2/cert/delete；调用方向：开发者调用华为服务端。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/app/agc-help-provision-api-delete-cent-0000002271000685
		// 现行文档（已核对登记路径）：https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-provision-api-delete-cent-0000002271000685
		officialEndpoint("provisioning", "provision-api-delete-cent", "删除证书", "此接口用于删除创建的证书。", "POST", "/api/publish/v2/cert/delete", []string{}, "agc-help-provision-api-delete-cent-0000002271000685", "https://developer.huawei.com/consumer/cn/doc/app/agc-help-provision-api-delete-cent-0000002271000685", "developer-to-huawei"),
		// 删除设备：此接口用于删除已添加的设备。
		// 请求：POST /api/publish/v2/device/delete；调用方向：开发者调用华为服务端。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/app/agc-help-provision-api-delete-device-0000002271160633
		// 现行文档（已核对登记路径）：https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-provision-api-delete-device-0000002271160633
		officialEndpoint("provisioning", "provision-api-delete-device", "删除设备", "此接口用于删除已添加的设备。", "POST", "/api/publish/v2/device/delete", []string{}, "agc-help-provision-api-delete-device-0000002271160633", "https://developer.huawei.com/consumer/cn/doc/app/agc-help-provision-api-delete-device-0000002271160633", "developer-to-huawei"),
		// 删除证书指纹：此接口用于删除证书指纹。
		// 请求：DELETE /api/provision/v1/fingerprints；调用方向：开发者调用华为服务端。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/app/agc-help-provision-api-delete-fingerprints-0000002464839925
		// 现行文档（已核对登记路径）：https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-provision-api-delete-fingerprints-0000002464839925
		officialEndpoint("provisioning", "provision-api-delete-fingerprints", "删除证书指纹", "此接口用于删除证书指纹。", "DELETE", "/api/provision/v1/fingerprints", []string{}, "agc-help-provision-api-delete-fingerprints-0000002464839925", "https://developer.huawei.com/consumer/cn/doc/app/agc-help-provision-api-delete-fingerprints-0000002464839925", "developer-to-huawei"),
		// 删除Profile：此接口用于删除申请的Profile。
		// 请求：DELETE /api/publish/v2/provision；调用方向：开发者调用华为服务端。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/app/agc-help-provision-api-delete-provision-0000002236201310
		// 现行文档（已核对登记路径）：https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-provision-api-delete-provision-0000002236201310
		officialEndpoint("provisioning", "provision-api-delete-provision", "删除Profile", "此接口用于删除申请的Profile。", "DELETE", "/api/publish/v2/provision", []string{}, "agc-help-provision-api-delete-provision-0000002236201310", "https://developer.huawei.com/consumer/cn/doc/app/agc-help-provision-api-delete-provision-0000002236201310", "developer-to-huawei"),
		// 查询可申请的ACL权限列表：此接口用于查询指定APP ID下所有可申请的ACL权限列表。
		// 请求：POST /api/provision/v1/user/permission/apply/list；调用方向：开发者调用华为服务端。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/app/agc-help-provision-api-eligible-acl-0000002554975044
		// 现行文档（已核对登记路径）：https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-provision-api-eligible-acl-0000002554975044
		officialEndpoint("provisioning", "provision-api-eligible-acl", "查询可申请的ACL权限列表", "此接口用于查询指定APP ID下所有可申请的ACL权限列表。", "POST", "/api/provision/v1/user/permission/apply/list", []string{}, "agc-help-provision-api-eligible-acl-0000002554975044", "https://developer.huawei.com/consumer/cn/doc/app/agc-help-provision-api-eligible-acl-0000002554975044", "developer-to-huawei"),
		// 查询证书指纹：此接口用于查询已创建的证书指纹。
		// 请求：GET /api/provision/v1/fingerprints；调用方向：开发者调用华为服务端。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/app/agc-help-provision-api-get-fingerprints-0000002431361316
		// 现行文档（已核对登记路径）：https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-provision-api-get-fingerprints-0000002431361316
		officialEndpoint("provisioning", "provision-api-get-fingerprints", "查询证书指纹", "此接口用于查询已创建的证书指纹。", "GET", "/api/provision/v1/fingerprints", []string{}, "agc-help-provision-api-get-fingerprints-0000002431361316", "https://developer.huawei.com/consumer/cn/doc/app/agc-help-provision-api-get-fingerprints-0000002431361316", "developer-to-huawei"),
		// 查询已申请的受限ACL权限：此接口用于查询申请的受限ACL权限信息。
		// 请求：GET /api/provision/v1/user/permission；调用方向：开发者调用华为服务端。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/app/agc-help-provision-api-getacl-0000002465770113
		// 现行文档（已核对登记路径）：https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-provision-api-getacl-0000002465770113
		officialEndpoint("provisioning", "provision-api-getacl", "查询已申请的受限ACL权限", "此接口用于查询申请的受限ACL权限信息。", "GET", "/api/provision/v1/user/permission", []string{}, "agc-help-provision-api-getacl-0000002465770113", "https://developer.huawei.com/consumer/cn/doc/app/agc-help-provision-api-getacl-0000002465770113", "developer-to-huawei"),
		// 查询受限ACL权限的申请状态：此接口用于查询受限ACL权限的申请状态，包含“申请中”和“审核通过”。
		// 请求：GET /api/provision/v1/user/permission/apply/status；调用方向：开发者调用华为服务端。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/app/agc-help-provision-api-getaclapplystatus-0000002469246150
		// 现行文档（已核对登记路径）：https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-provision-api-getaclapplystatus-0000002469246150
		officialEndpoint("provisioning", "provision-api-getaclapplystatus", "查询受限ACL权限的申请状态", "此接口用于查询受限ACL权限的申请状态，包含“申请中”和“审核通过”。", "GET", "/api/provision/v1/user/permission/apply/status", []string{}, "agc-help-provision-api-getaclapplystatus-0000002469246150", "https://developer.huawei.com/consumer/cn/doc/app/agc-help-provision-api-getaclapplystatus-0000002469246150", "developer-to-huawei"),
		// 查询证书：此接口用于查询创建的证书。
		// 请求：POST /api/publish/v3/cert/list；调用方向：开发者调用华为服务端。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/app/agc-help-provision-api-query-cent-0000002271160629
		// 现行文档（已核对登记路径）：https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-provision-api-query-cent-0000002271160629
		officialEndpoint("provisioning", "provision-api-query-cent", "查询证书", "此接口用于查询创建的证书。", "POST", "/api/publish/v3/cert/list", []string{}, "agc-help-provision-api-query-cent-0000002271160629", "https://developer.huawei.com/consumer/cn/doc/app/agc-help-provision-api-query-cent-0000002271160629", "developer-to-huawei"),
		// 查询设备列表：此接口用于查询当前设备的信息。
		// 请求：GET /api/publish/v2/device/list；调用方向：开发者调用华为服务端。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/app/agc-help-provision-api-query-device-0000002236201306
		// 现行文档（已核对登记路径）：https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-provision-api-query-device-0000002236201306
		officialEndpoint("provisioning", "provision-api-query-device", "查询设备列表", "此接口用于查询当前设备的信息。", "GET", "/api/publish/v2/device/list", []string{}, "agc-help-provision-api-query-device-0000002236201306", "https://developer.huawei.com/consumer/cn/doc/app/agc-help-provision-api-query-device-0000002236201306", "developer-to-huawei"),
		// 查询Profile列表：此接口用于查询申请的Profile。
		// 请求：GET /api/publish/v3/provision/list；调用方向：开发者调用华为服务端。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/app/agc-help-provision-api-query-provision-0000002236041506
		// 现行文档（已核对登记路径）：https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-provision-api-query-provision-0000002236041506
		officialEndpoint("provisioning", "provision-api-query-provision", "查询Profile列表", "此接口用于查询申请的Profile。", "GET", "/api/publish/v3/provision/list", []string{}, "agc-help-provision-api-query-provision-0000002236041506", "https://developer.huawei.com/consumer/cn/doc/app/agc-help-provision-api-query-provision-0000002236041506", "developer-to-huawei"),
		// 修改Profile绑定的设备：此接口用于修改调试Profile绑定的调试设备，或者修改指定设备发布Profile绑定的测试设备。
		// 请求：PUT /api/publish/v3/provision；调用方向：开发者调用华为服务端。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/app/agc-help-provision-api-update-provision-0000002469198766
		// 现行文档（已核对登记路径）：https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-provision-api-update-provision-0000002469198766
		officialEndpoint("provisioning", "provision-api-update-provision", "修改Profile绑定的设备", "此接口用于修改调试Profile绑定的调试设备，或者修改指定设备发布Profile绑定的测试设备。", "PUT", "/api/publish/v3/provision", []string{}, "agc-help-provision-api-update-provision-0000002469198766", "https://developer.huawei.com/consumer/cn/doc/app/agc-help-provision-api-update-provision-0000002469198766", "developer-to-huawei"),
	},
	"domains": {
		// 下载域名配置文件：此接口用于以文件流的形式进行直接下载域名配置文件。
		// 请求：GET /api/dms/domain-manage/v1/app/domain/verify-file；调用方向：开发者调用华为服务端。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/app/agc-help-domain-api-download-domain-config-0000002271160653
		// 现行文档（已核对登记路径）：https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-domain-api-download-domain-config-0000002271160653
		officialEndpoint("domains", "domain-api-download-domain-config", "下载域名配置文件", "此接口用于以文件流的形式进行直接下载域名配置文件。", "GET", "/api/dms/domain-manage/v1/app/domain/verify-file", []string{}, "agc-help-domain-api-download-domain-config-0000002271160653", "https://developer.huawei.com/consumer/cn/doc/app/agc-help-domain-api-download-domain-config-0000002271160653", "developer-to-huawei"),
		// 查询元服务的域名配置信息：此接口用于查询元服务的域名配置信息。
		// 请求：GET /api/dms/domain-manage/v1/app/domain；调用方向：开发者调用华为服务端。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/app/agc-help-domain-api-get-domain-0000002271000701
		// 现行文档（已核对登记路径）：https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-domain-api-get-domain-0000002271000701
		officialEndpoint("domains", "domain-api-get-domain", "查询元服务的域名配置信息", "此接口用于查询元服务的域名配置信息。", "GET", "/api/dms/domain-manage/v1/app/domain", []string{}, "agc-help-domain-api-get-domain-0000002271000701", "https://developer.huawei.com/consumer/cn/doc/app/agc-help-domain-api-get-domain-0000002271000701", "developer-to-huawei"),
		// 查询域名修改次数/配置上限：此接口用于查询域名修改次数、配置上限等信息。
		// 请求：GET /api/dms/domain-manage/v1/app/domain/config；调用方向：开发者调用华为服务端。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/app/agc-help-domain-api-get-domain-config-0000002236201326
		// 现行文档（已核对登记路径）：https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-domain-api-get-domain-config-0000002236201326
		officialEndpoint("domains", "domain-api-get-domain-config", "查询域名修改次数/配置上限", "此接口用于查询域名修改次数、配置上限等信息。", "GET", "/api/dms/domain-manage/v1/app/domain/config", []string{}, "agc-help-domain-api-get-domain-config-0000002236201326", "https://developer.huawei.com/consumer/cn/doc/app/agc-help-domain-api-get-domain-config-0000002236201326", "developer-to-huawei"),
		// 新增/更新元服务的域名配置信息：此接口用于新增/更新元服务的域名配置信息。
		// 请求：POST /api/dms/domain-manage/v1/app/domain；调用方向：开发者调用华为服务端。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/app/agc-help-domain-api-post-domain-0000002236041522
		// 现行文档（已核对登记路径）：https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-domain-api-post-domain-0000002236041522
		officialEndpoint("domains", "domain-api-post-domain", "新增/更新元服务的域名配置信息", "此接口用于新增/更新元服务的域名配置信息。", "POST", "/api/dms/domain-manage/v1/app/domain", []string{}, "agc-help-domain-api-post-domain-0000002236041522", "https://developer.huawei.com/consumer/cn/doc/app/agc-help-domain-api-post-domain-0000002236041522", "developer-to-huawei"),
		// 预检查业务域名配置：此接口用于对业务域名配置进行预检查。
		// 请求：POST /api/dms/domain-manage/v1/app/domain/pre-check；调用方向：开发者调用华为服务端。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/app/agc-help-domain-api-precheck-0000002295134522
		// 现行文档（已核对登记路径）：https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-domain-api-precheck-0000002295134522
		officialEndpoint("domains", "domain-api-precheck", "预检查业务域名配置", "此接口用于对业务域名配置进行预检查。", "POST", "/api/dms/domain-manage/v1/app/domain/pre-check", []string{}, "agc-help-domain-api-precheck-0000002295134522", "https://developer.huawei.com/consumer/cn/doc/app/agc-help-domain-api-precheck-0000002295134522", "developer-to-huawei"),
	},
	"testing": {
		// 生成邀请码：此接口用于邀请测试生成邀请码。
		// 请求：POST /api/app-test/v1/invitation-code；调用方向：开发者调用华为服务端。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/app/agc-help-test-api-add-invite-code-0000002236201342
		// 现行文档（已核对登记路径）：https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-test-api-add-invite-code-0000002236201342
		officialEndpoint("testing", "test-api-add-invite-code", "生成邀请码", "此接口用于邀请测试生成邀请码。", "POST", "/api/app-test/v1/invitation-code", []string{}, "agc-help-test-api-add-invite-code-0000002236201342", "https://developer.huawei.com/consumer/cn/doc/app/agc-help-test-api-add-invite-code-0000002236201342", "developer-to-huawei"),
		// 新增测试群组：此接口用于邀请测试创建内外部测试群组。
		// 请求：POST /api/app-test/v1/test-group；调用方向：开发者调用华为服务端。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/app/agc-help-test-api-add-test-group-0000002271160661
		// 现行文档（已核对登记路径）：https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-test-api-add-test-group-0000002271160661
		officialEndpoint("testing", "test-api-add-test-group", "新增测试群组", "此接口用于邀请测试创建内外部测试群组。", "POST", "/api/app-test/v1/test-group", []string{}, "agc-help-test-api-add-test-group-0000002271160661", "https://developer.huawei.com/consumer/cn/doc/app/agc-help-test-api-add-test-group-0000002271160661", "developer-to-huawei"),
		// 新增测试群组成员：此接口用于邀请测试新增测试群组成员。
		// 请求：POST /api/publish/v2/test/user；调用方向：开发者调用华为服务端。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/app/agc-help-test-api-add-test-group-user-0000002236201338
		// 现行文档（已核对登记路径）：https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-test-api-add-test-group-user-0000002236201338
		officialEndpoint("testing", "test-api-add-test-group-user", "新增测试群组成员", "此接口用于邀请测试新增测试群组成员。", "POST", "/api/publish/v2/test/user", []string{}, "agc-help-test-api-add-test-group-user-0000002236201338", "https://developer.huawei.com/consumer/cn/doc/app/agc-help-test-api-add-test-group-user-0000002236201338", "developer-to-huawei"),
		// 添加软件包：此接口用于测试版本上传软件包。
		// 请求：POST /api/publish/v2/test/version/pkg；调用方向：开发者调用华为服务端。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/app/agc-help-test-api-add-test-package-0000002236201330
		// 现行文档（已核对登记路径）：https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-test-api-add-test-package-0000002236201330
		officialEndpoint("testing", "test-api-add-test-package", "添加软件包", "此接口用于测试版本上传软件包。", "POST", "/api/publish/v2/test/version/pkg", []string{}, "agc-help-test-api-add-test-package-0000002236201330", "https://developer.huawei.com/consumer/cn/doc/app/agc-help-test-api-add-test-package-0000002236201330", "developer-to-huawei"),
		// 新建测试版本：此接口用于创建HarmonyOS应用/元服务的测试版本。
		// 请求：POST /api/publish/v2/test/app/version；调用方向：开发者调用华为服务端。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/app/agc-help-test-api-add-test-version-0000002236041526
		// 现行文档（已核对登记路径）：https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-test-api-add-test-version-0000002236041526
		officialEndpoint("testing", "test-api-add-test-version", "新建测试版本", "此接口用于创建HarmonyOS应用/元服务的测试版本。", "POST", "/api/publish/v2/test/app/version", []string{}, "agc-help-test-api-add-test-version-0000002236041526", "https://developer.huawei.com/consumer/cn/doc/app/agc-help-test-api-add-test-version-0000002236041526", "developer-to-huawei"),
		// 创建公开链接：此接口用于创建群组级公开测试链接。
		// 请求：POST /api/app-test/v1/test-group/public-link；调用方向：开发者调用华为服务端。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/app/agc-help-test-api-create-open-test-link-0000002665117715
		// 现行文档（已核对登记路径）：https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-test-api-create-open-test-link-0000002665117715
		officialEndpoint("testing", "test-api-create-open-test-link", "创建公开链接", "此接口用于创建群组级公开测试链接。", "POST", "/api/app-test/v1/test-group/public-link", []string{}, "agc-help-test-api-create-open-test-link-0000002665117715", "https://developer.huawei.com/consumer/cn/doc/app/agc-help-test-api-create-open-test-link-0000002665117715", "developer-to-huawei"),
		// 删除测试问题反馈：此接口用于删除测试问题反馈。
		// 请求：DELETE /api/app-test/v1/mgmt/tester-feedback；调用方向：开发者调用华为服务端。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/app/agc-help-test-api-delete-feedback-0000002634718668
		// 现行文档（已核对登记路径）：https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-test-api-delete-feedback-0000002634718668
		officialEndpoint("testing", "test-api-delete-feedback", "删除测试问题反馈", "此接口用于删除测试问题反馈。", "DELETE", "/api/app-test/v1/mgmt/tester-feedback", []string{}, "agc-help-test-api-delete-feedback-0000002634718668", "https://developer.huawei.com/consumer/cn/doc/app/agc-help-test-api-delete-feedback-0000002634718668", "developer-to-huawei"),
		// 删除邀请码：此接口用于邀请测试删除邀请码。
		// 请求：DELETE /api/app-test/v1/invitation-code；调用方向：开发者调用华为服务端。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/app/agc-help-test-api-delete-invite-code-0000002634878570
		// 现行文档（已核对登记路径）：https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-test-api-delete-invite-code-0000002634878570
		officialEndpoint("testing", "test-api-delete-invite-code", "删除邀请码", "此接口用于邀请测试删除邀请码。", "DELETE", "/api/app-test/v1/invitation-code", []string{}, "agc-help-test-api-delete-invite-code-0000002634878570", "https://developer.huawei.com/consumer/cn/doc/app/agc-help-test-api-delete-invite-code-0000002634878570", "developer-to-huawei"),
		// 删除测试群组：此接口用于邀请测试删除测试群组。
		// 请求：DELETE /api/app-test/v1/publish-api/test-group；调用方向：开发者调用华为服务端。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/app/agc-help-test-api-delete-test-group-0000002236041534
		// 现行文档（已核对登记路径）：https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-test-api-delete-test-group-0000002236041534
		officialEndpoint("testing", "test-api-delete-test-group", "删除测试群组", "此接口用于邀请测试删除测试群组。", "DELETE", "/api/app-test/v1/publish-api/test-group", []string{}, "agc-help-test-api-delete-test-group-0000002236041534", "https://developer.huawei.com/consumer/cn/doc/app/agc-help-test-api-delete-test-group-0000002236041534", "developer-to-huawei"),
		// 删除测试群组成员：此接口用于邀请测试删除测试群组中的成员。
		// 请求：DELETE /api/publish/v2/test/user；调用方向：开发者调用华为服务端。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/app/agc-help-test-api-delete-test-group-user-0000002271000721
		// 现行文档（已核对登记路径）：https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-test-api-delete-test-group-user-0000002271000721
		officialEndpoint("testing", "test-api-delete-test-group-user", "删除测试群组成员", "此接口用于邀请测试删除测试群组中的成员。", "DELETE", "/api/publish/v2/test/user", []string{}, "agc-help-test-api-delete-test-group-user-0000002271000721", "https://developer.huawei.com/consumer/cn/doc/app/agc-help-test-api-delete-test-group-user-0000002271000721", "developer-to-huawei"),
		// 删除测试版本：此接口用于删除测试版本。
		// 请求：DELETE /api/publish/v2/test/app/version；调用方向：开发者调用华为服务端。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/app/agc-help-test-api-delete-test-version-0000002236201346
		// 现行文档（已核对登记路径）：https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-test-api-delete-test-version-0000002236201346
		officialEndpoint("testing", "test-api-delete-test-version", "删除测试版本", "此接口用于删除测试版本。", "DELETE", "/api/publish/v2/test/app/version", []string{}, "agc-help-test-api-delete-test-version-0000002236201346", "https://developer.huawei.com/consumer/cn/doc/app/agc-help-test-api-delete-test-version-0000002236201346", "developer-to-huawei"),
		// 修改公开链接：此接口用于编辑群组级公开测试链接。
		// 请求：PUT /api/app-test/v1/test-group/public-link；调用方向：开发者调用华为服务端。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/app/agc-help-test-api-edit-open-test-link-0000002634718666
		// 现行文档（已核对登记路径）：https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-test-api-edit-open-test-link-0000002634718666
		officialEndpoint("testing", "test-api-edit-open-test-link", "修改公开链接", "此接口用于编辑群组级公开测试链接。", "PUT", "/api/app-test/v1/test-group/public-link", []string{}, "agc-help-test-api-edit-open-test-link-0000002634718666", "https://developer.huawei.com/consumer/cn/doc/app/agc-help-test-api-edit-open-test-link-0000002634718666", "developer-to-huawei"),
		// 修改测试群组：此接口用于邀请测试修改测试群组。
		// 请求：PUT /api/app-test/v2/test-group；调用方向：开发者调用华为服务端。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/app/agc-help-test-api-edit-test-group-0000002634718664
		// 现行文档（已核对登记路径）：https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-test-api-edit-test-group-0000002634718664
		officialEndpoint("testing", "test-api-edit-test-group", "修改测试群组", "此接口用于邀请测试修改测试群组。", "PUT", "/api/app-test/v2/test-group", []string{}, "agc-help-test-api-edit-test-group-0000002634718664", "https://developer.huawei.com/consumer/cn/doc/app/agc-help-test-api-edit-test-group-0000002634718664", "developer-to-huawei"),
		// 查询上架自检报告信息：此接口用于查询邀请测试进行上架自检的报告信息。
		// 请求：GET /api/publish/v2/test/app/version/detect-report；调用方向：开发者调用华为服务端。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/app/agc-help-test-api-get-detect-task-0000002525033189
		// 现行文档（已核对登记路径）：https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-test-api-get-detect-task-0000002525033189
		officialEndpoint("testing", "test-api-get-detect-task", "查询上架自检报告信息", "此接口用于查询邀请测试进行上架自检的报告信息。", "GET", "/api/publish/v2/test/app/version/detect-report", []string{}, "agc-help-test-api-get-detect-task-0000002525033189", "https://developer.huawei.com/consumer/cn/doc/app/agc-help-test-api-get-detect-task-0000002525033189", "developer-to-huawei"),
		// 查询测试群组邀请码：此接口用于邀请测试查询测试群组邀请码信息。
		// 请求：GET /api/app-test/v1/test-group/invitation-code；调用方向：开发者调用华为服务端。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/app/agc-help-test-api-get-invite-code-0000002271160669
		// 现行文档（已核对登记路径）：https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-test-api-get-invite-code-0000002271160669
		officialEndpoint("testing", "test-api-get-invite-code", "查询测试群组邀请码", "此接口用于邀请测试查询测试群组邀请码信息。", "GET", "/api/app-test/v1/test-group/invitation-code", []string{}, "agc-help-test-api-get-invite-code-0000002271160669", "https://developer.huawei.com/consumer/cn/doc/app/agc-help-test-api-get-invite-code-0000002271160669", "developer-to-huawei"),
		// 查询测试群组列表：此接口用于邀请测试查询已创建的内外部测试群组列表。
		// 请求：GET /api/app-test/v1/test-group/list；调用方向：开发者调用华为服务端。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/app/agc-help-test-api-get-test-grouplist-0000002271000717
		// 现行文档（已核对登记路径）：https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-test-api-get-test-grouplist-0000002271000717
		officialEndpoint("testing", "test-api-get-test-grouplist", "查询测试群组列表", "此接口用于邀请测试查询已创建的内外部测试群组列表。", "GET", "/api/app-test/v1/test-group/list", []string{}, "agc-help-test-api-get-test-grouplist-0000002271000717", "https://developer.huawei.com/consumer/cn/doc/app/agc-help-test-api-get-test-grouplist-0000002271000717", "developer-to-huawei"),
		// 更新测试版本：此接口用于更新已创建的测试版本。
		// 请求：PUT /api/publish/v2/test/app/version；调用方向：开发者调用华为服务端。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/app/agc-help-test-api-modify-test-version-0000002271160657
		// 现行文档（已核对登记路径）：https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-test-api-modify-test-version-0000002271160657
		officialEndpoint("testing", "test-api-modify-test-version", "更新测试版本", "此接口用于更新已创建的测试版本。", "PUT", "/api/publish/v2/test/app/version", []string{}, "agc-help-test-api-modify-test-version-0000002271160657", "https://developer.huawei.com/consumer/cn/doc/app/agc-help-test-api-modify-test-version-0000002271160657", "developer-to-huawei"),
		// 发送邀请测试通知：此接口用于给内外部群组的测试用户发送邀请测试的通知。
		// 请求：POST /api/publish/v1/test/user/action/notify；调用方向：开发者调用华为服务端。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/app/agc-help-test-api-post-test-notify-0000002236041538
		// 现行文档（已核对登记路径）：https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-test-api-post-test-notify-0000002236041538
		officialEndpoint("testing", "test-api-post-test-notify", "发送邀请测试通知", "此接口用于给内外部群组的测试用户发送邀请测试的通知。", "POST", "/api/publish/v1/test/user/action/notify", []string{}, "agc-help-test-api-post-test-notify-0000002236041538", "https://developer.huawei.com/consumer/cn/doc/app/agc-help-test-api-post-test-notify-0000002236041538", "developer-to-huawei"),
		// 查询测试问题反馈：此接口用于查询测试问题反馈列表。
		// 请求：POST /api/app-test/v1/mgmt/tester-feedback/action/query；调用方向：开发者调用华为服务端。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/app/agc-help-test-api-query-feedback-0000002665117717
		// 现行文档（已核对登记路径）：https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-test-api-query-feedback-0000002665117717
		officialEndpoint("testing", "test-api-query-feedback", "查询测试问题反馈", "此接口用于查询测试问题反馈列表。", "POST", "/api/app-test/v1/mgmt/tester-feedback/action/query", []string{}, "agc-help-test-api-query-feedback-0000002665117717", "https://developer.huawei.com/consumer/cn/doc/app/agc-help-test-api-query-feedback-0000002665117717", "developer-to-huawei"),
		// 查询测试问题反馈维度：此接口用于查询测试反馈维度。
		// 请求：POST /api/app-test/v1/mgmt/tester-feedback-dimension/action/query；调用方向：开发者调用华为服务端。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/app/agc-help-test-api-query-feedback-dimension-0000002634878572
		// 现行文档（已核对登记路径）：https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-test-api-query-feedback-dimension-0000002634878572
		officialEndpoint("testing", "test-api-query-feedback-dimension", "查询测试问题反馈维度", "此接口用于查询测试反馈维度。", "POST", "/api/app-test/v1/mgmt/tester-feedback-dimension/action/query", []string{}, "agc-help-test-api-query-feedback-dimension-0000002634878572", "https://developer.huawei.com/consumer/cn/doc/app/agc-help-test-api-query-feedback-dimension-0000002634878572", "developer-to-huawei"),
		// 查询公开链接：此接口用于查询群组级公开测试链接。
		// 请求：GET /api/app-test/v1/test-group/public-link；调用方向：开发者调用华为服务端。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/app/agc-help-test-api-query-open-test-link-0000002664997775
		// 现行文档（已核对登记路径）：https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-test-api-query-open-test-link-0000002664997775
		officialEndpoint("testing", "test-api-query-open-test-link", "查询公开链接", "此接口用于查询群组级公开测试链接。", "GET", "/api/app-test/v1/test-group/public-link", []string{}, "agc-help-test-api-query-open-test-link-0000002664997775", "https://developer.huawei.com/consumer/cn/doc/app/agc-help-test-api-query-open-test-link-0000002664997775", "developer-to-huawei"),
		// 查询测试群组成员：此接口用于查询测试群组成员。
		// 请求：GET /api/app-test/v1/test-user；调用方向：开发者调用华为服务端。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/app/agc-help-test-api-query-test-user-0000002664997773
		// 现行文档（已核对登记路径）：https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-test-api-query-test-user-0000002664997773
		officialEndpoint("testing", "test-api-query-test-user", "查询测试群组成员", "此接口用于查询测试群组成员。", "GET", "/api/app-test/v1/test-user", []string{}, "agc-help-test-api-query-test-user-0000002664997773", "https://developer.huawei.com/consumer/cn/doc/app/agc-help-test-api-query-test-user-0000002664997773", "developer-to-huawei"),
		// 更新测试生效版本：此接口用于将公开测试版本转为全网版本或者分阶段版本。
		// 请求：PUT /api/publish/v2/test/version/release；调用方向：开发者调用华为服务端。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/app/agc-help-test-api-release-test-version-0000002271000713
		// 现行文档（已核对登记路径）：https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-test-api-release-test-version-0000002271000713
		officialEndpoint("testing", "test-api-release-test-version", "更新测试生效版本", "此接口用于将公开测试版本转为全网版本或者分阶段版本。", "PUT", "/api/publish/v2/test/version/release", []string{}, "agc-help-test-api-release-test-version-0000002271000713", "https://developer.huawei.com/consumer/cn/doc/app/agc-help-test-api-release-test-version-0000002271000713", "developer-to-huawei"),
		// 停止邀请码：此接口用于邀请测试停止邀请码。
		// 请求：POST /api/app-test/v1/invitation-code/action/stop；调用方向：开发者调用华为服务端。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/app/agc-help-test-api-stop-invite-code-0000002271000725
		// 现行文档（已核对登记路径）：https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-test-api-stop-invite-code-0000002271000725
		officialEndpoint("testing", "test-api-stop-invite-code", "停止邀请码", "此接口用于邀请测试停止邀请码。", "POST", "/api/app-test/v1/invitation-code/action/stop", []string{}, "agc-help-test-api-stop-invite-code-0000002271000725", "https://developer.huawei.com/consumer/cn/doc/app/agc-help-test-api-stop-invite-code-0000002271000725", "developer-to-huawei"),
		// 停止测试版本：此接口用于停止测试版本。
		// 请求：POST /api/publish/v2/test/version/stop；调用方向：开发者调用华为服务端。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/app/agc-help-test-api-stop-test-version-0000002236041542
		// 现行文档（已核对登记路径）：https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-test-api-stop-test-version-0000002236041542
		officialEndpoint("testing", "test-api-stop-test-version", "停止测试版本", "此接口用于停止测试版本。", "POST", "/api/publish/v2/test/version/stop", []string{}, "agc-help-test-api-stop-test-version-0000002236041542", "https://developer.huawei.com/consumer/cn/doc/app/agc-help-test-api-stop-test-version-0000002236041542", "developer-to-huawei"),
		// 提交测试版本：此接口用于测试版本提交审核。
		// 请求：POST /api/publish/v2/test/app/version/submit；调用方向：开发者调用华为服务端。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/app/agc-help-test-api-submit-test-version-0000002236201334
		// 现行文档（已核对登记路径）：https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-test-api-submit-test-version-0000002236201334
		officialEndpoint("testing", "test-api-submit-test-version", "提交测试版本", "此接口用于测试版本提交审核。", "POST", "/api/publish/v2/test/app/version/submit", []string{}, "agc-help-test-api-submit-test-version-0000002236201334", "https://developer.huawei.com/consumer/cn/doc/app/agc-help-test-api-submit-test-version-0000002236201334", "developer-to-huawei"),
		// 更新有效期内的测试版本信息：此接口用于更新有效期内测试版本的测试时间和下载安装次数。
		// 请求：PUT /api/publish/v2/test/version/open-test；调用方向：开发者调用华为服务端。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/app/agc-help-test-api-update-test-version-opentest-0000002236041530
		// 现行文档（已核对登记路径）：https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-test-api-update-test-version-opentest-0000002236041530
		officialEndpoint("testing", "test-api-update-test-version-opentest", "更新有效期内的测试版本信息", "此接口用于更新有效期内测试版本的测试时间和下载安装次数。", "PUT", "/api/publish/v2/test/version/open-test", []string{}, "agc-help-test-api-update-test-version-opentest-0000002236041530", "https://developer.huawei.com/consumer/cn/doc/app/agc-help-test-api-update-test-version-opentest-0000002236041530", "developer-to-huawei"),
	},
	"reports": {
		// 获取给指定用户群发奖的报表：获取指定给用户群发奖的报表数据的CSV文件地址。
		// 请求：GET /api/report/distribution-operation-quality/v1/activityAwardExport/{appId}；调用方向：开发者调用华为服务端。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-activityawardexport-0000001158245071
		officialEndpoint("reports", "activityawardexport", "获取给指定用户群发奖的报表", "获取指定给用户群发奖的报表数据的CSV文件地址。", "GET", "/api/report/distribution-operation-quality/v1/activityAwardExport/{appId}", []string{"appId"}, "agcapi-activityawardexport-0000001158245071", "https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-activityawardexport-0000001158245071", "developer-to-huawei"),
		// 获取新增和留存的报表：获取新增和留存的报表数据的CSV或者Excel文件地址。
		// 请求：GET /api/report/distribution-operation-quality/v1/addAdKpExport/{appId}；调用方向：开发者调用华为服务端。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-addadkpexport-0000001111845102
		officialEndpoint("reports", "addadkpexport", "获取新增和留存的报表", "获取新增和留存的报表数据的CSV或者Excel文件地址。", "GET", "/api/report/distribution-operation-quality/v1/addAdKpExport/{appId}", []string{"appId"}, "agcapi-addadkpexport-0000001111845102", "https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-addadkpexport-0000001111845102", "developer-to-huawei"),
		// 获取下载安装的报表：获取下载安装的报表数据的CSV或者Excel文件地址。
		// 请求：GET /api/report/distribution-operation-quality/v1/appDownloadExport/{appId}；调用方向：开发者调用华为服务端。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-appdownloadexport-0000001158365059
		officialEndpoint("reports", "appdownloadexport", "获取下载安装的报表", "获取下载安装的报表数据的CSV或者Excel文件地址。", "GET", "/api/report/distribution-operation-quality/v1/appDownloadExport/{appId}", []string{"appId"}, "agcapi-appdownloadexport-0000001158365059", "https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-appdownloadexport-0000001158365059", "developer-to-huawei"),
		// 获取安装失败数据的报表：获取安装失败数据的报表CSV或EXCEL文件。
		// 请求：GET /api/report/distribution-operation-quality/v1/appDownloadFailExport/{appId}；调用方向：开发者调用华为服务端。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-appdownloadfailexport-0000001158365063
		officialEndpoint("reports", "appdownloadfailexport", "获取安装失败数据的报表", "获取安装失败数据的报表CSV或EXCEL文件。", "GET", "/api/report/distribution-operation-quality/v1/appDownloadFailExport/{appId}", []string{"appId"}, "agcapi-appdownloadfailexport-0000001158365063", "https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-appdownloadfailexport-0000001158365063", "developer-to-huawei"),
		// 获取优惠券活动的报表：获取优惠券活动的报表CSV或EXCEL文件。
		// 请求：GET /api/report/distribution-operation-quality/v1/activityCouponExport/{appId}；调用方向：开发者调用华为服务端。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-couponexport-0000001111845104
		officialEndpoint("reports", "couponexport", "获取优惠券活动的报表", "获取优惠券活动的报表CSV或EXCEL文件。", "GET", "/api/report/distribution-operation-quality/v1/activityCouponExport/{appId}", []string{"appId"}, "agcapi-couponexport-0000001111845104", "https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-couponexport-0000001111845104", "developer-to-huawei"),
		// 获取元服务分发分析的报表文件：获取元服务分发分析的报表数据的CSV或者Excel文件地址。
		// 请求：GET /api/report/distribution-operation-quality/v1/fa/distributeAnalysisExport/；调用方向：开发者调用华为服务端。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-fa-distribute_analysis_export-0000001690307352
		officialEndpoint("reports", "fa-distribute-analysis-export", "获取元服务分发分析的报表文件", "获取元服务分发分析的报表数据的CSV或者Excel文件地址。", "GET", "/api/report/distribution-operation-quality/v1/fa/distributeAnalysisExport/", []string{}, "agcapi-fa-distribute_analysis_export-0000001690307352", "https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-fa-distribute_analysis_export-0000001690307352", "developer-to-huawei"),
		// 获取元服务新增留存的报表文件：获取元服务新增留存的报表数据的CSV或者Excel文件地址。
		// 请求：GET /api/report/distribution-operation-quality/v1/fa/userAnalysisExport/；调用方向：开发者调用华为服务端。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-fa-user_analysis_export-0000001738465949
		officialEndpoint("reports", "fa-user-analysis-export", "获取元服务新增留存的报表文件", "获取元服务新增留存的报表数据的CSV或者Excel文件地址。", "GET", "/api/report/distribution-operation-quality/v1/fa/userAnalysisExport/", []string{}, "agcapi-fa-user_analysis_export-0000001738465949", "https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-fa-user_analysis_export-0000001738465949", "developer-to-huawei"),
		// 获取元服务卡片分析的报表文件：获取元服务卡片分析的报表数据的CSV或者Excel文件地址。
		// 请求：GET /api/report/distribution-operation-quality/v1/fa/widgetAnalysisExport/；调用方向：开发者调用华为服务端。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-fa-widget_analysis_export-0000001690467092
		officialEndpoint("reports", "fa-widget-analysis-export", "获取元服务卡片分析的报表文件", "获取元服务卡片分析的报表数据的CSV或者Excel文件地址。", "GET", "/api/report/distribution-operation-quality/v1/fa/widgetAnalysisExport/", []string{}, "agcapi-fa-widget_analysis_export-0000001690467092", "https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-fa-widget_analysis_export-0000001690467092", "developer-to-huawei"),
		// 获取预约的报表：获取预约报表数据的CSV或EXCEL文件地址。
		// 请求：GET /api/report/distribution-operation-quality/v1/gameReservationExport/{appId}；调用方向：开发者调用华为服务端。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-gamereservationexport-0000001111685212
		officialEndpoint("reports", "gamereservationexport", "获取预约的报表", "获取预约报表数据的CSV或EXCEL文件地址。", "GET", "/api/report/distribution-operation-quality/v1/gameReservationExport/{appId}", []string{"appId"}, "agcapi-gamereservationexport-0000001111685212", "https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-gamereservationexport-0000001111685212", "developer-to-huawei"),
		// 获取应用内付费的报表：获取应用内付费的报表的CSV或者EXCEL文件地址。
		// 请求：GET /api/report/distribution-operation-quality/v1/IAPExport/{appId}；调用方向：开发者调用华为服务端。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-iapexport-0000001158365061
		officialEndpoint("reports", "iapexport", "获取应用内付费的报表", "获取应用内付费的报表的CSV或者EXCEL文件地址。", "GET", "/api/report/distribution-operation-quality/v1/IAPExport/{appId}", []string{"appId"}, "agcapi-iapexport-0000001158365061", "https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-iapexport-0000001158365061", "developer-to-huawei"),
		// 获取付费下载的报表：获取付费下载报表的CSV或者EXCEL文件。
		// 请求：GET /api/report/distribution-operation-quality/v1/orderAnalysisExport/{appId}；调用方向：开发者调用华为服务端。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-orderanalysisexport-0000001111685214
		officialEndpoint("reports", "orderanalysisexport", "获取付费下载的报表", "获取付费下载报表的CSV或者EXCEL文件。", "GET", "/api/report/distribution-operation-quality/v1/orderAnalysisExport/{appId}", []string{"appId"}, "agcapi-orderanalysisexport-0000001111685214", "https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-orderanalysisexport-0000001111685214", "developer-to-huawei"),
		// 获取付费下载明细的报表：获取付费下载明细的报表CSV文件。
		// 请求：GET /api/report/distribution-operation-quality/v1/orderDetailExport/{appId}；调用方向：开发者调用华为服务端。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-orderdetailexport-0000001158245073
		officialEndpoint("reports", "orderdetailexport", "获取付费下载明细的报表", "获取付费下载明细的报表CSV文件。", "GET", "/api/report/distribution-operation-quality/v1/orderDetailExport/{appId}", []string{"appId"}, "agcapi-orderdetailexport-0000001158245073", "https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-orderdetailexport-0000001158245073", "developer-to-huawei"),
	},
	"projects": {
		// 添加证书指纹：此接口用于设置应用的证书指纹信息。
		// 请求：PUT /api/cds/app-distirbution/v1/all/app-extra-info/{appId}；调用方向：开发者调用华为服务端。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-addappfingerprint-0000001111685218
		officialEndpoint("projects", "addappfingerprint", "添加证书指纹", "此接口用于设置应用的证书指纹信息。", "PUT", "/api/cds/app-distirbution/v1/all/app-extra-info/{appId}", []string{"appId"}, "agcapi-addappfingerprint-0000001111685218", "https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-addappfingerprint-0000001111685218", "developer-to-huawei"),
		// 获取应用简略信息：批量查询应用简略信息。
		// 请求：GET /api/cds/app-distirbution/v1/all/app-brief-info/list；调用方向：开发者调用华为服务端。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-appbriefinfo-0000001111845106
		officialEndpoint("projects", "appbriefinfo", "获取应用简略信息", "批量查询应用简略信息。", "GET", "/api/cds/app-distirbution/v1/all/app-brief-info/list", []string{}, "agcapi-appbriefinfo-0000001111845106", "https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-appbriefinfo-0000001111845106", "developer-to-huawei"),
		// 获取配置文件：获取应用的配置文件。
		// 请求：GET /api/cpms/project-management-service/v1/config-file；调用方向：开发者调用华为服务端。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-getconfigfile-0000001158365065
		officialEndpoint("projects", "getconfigfile", "获取配置文件", "获取应用的配置文件。", "GET", "/api/cpms/project-management-service/v1/config-file", []string{}, "agcapi-getconfigfile-0000001158365065", "https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-getconfigfile-0000001158365065", "developer-to-huawei"),
		// 获取团队列表：获取开发者的团队账号列表信息。
		// 请求：GET /api/ups/user-permission-service/v1/user-team-list；调用方向：开发者调用华为服务端。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-getteamlist-0000001158245075
		officialEndpoint("projects", "getteamlist", "获取团队列表", "获取开发者的团队账号列表信息。", "GET", "/api/ups/user-permission-service/v1/user-team-list", []string{}, "agcapi-getteamlist-0000001158245075", "https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-getteamlist-0000001158245075", "developer-to-huawei"),
		// 查询证书指纹和AppSecret：此接口用于根据应用ID查询App Secret、证书指纹信息。
		// 请求：GET /api/cds/app-distirbution/v1/all/app-extra-info/{appId}；调用方向：开发者调用华为服务端。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-queryappfingerprint-0000001158245077
		officialEndpoint("projects", "queryappfingerprint", "查询证书指纹和AppSecret", "此接口用于根据应用ID查询App Secret、证书指纹信息。", "GET", "/api/cds/app-distirbution/v1/all/app-extra-info/{appId}", []string{"appId"}, "agcapi-queryappfingerprint-0000001158245077", "https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-queryappfingerprint-0000001158245077", "developer-to-huawei"),
		// 查询项目详情及项目下的应用：此接口用于查询指定项目的详细项目信息和项目下的应用列表信息。
		// 请求：GET /api/project-service/v1/projects/{projectId}；调用方向：开发者调用华为服务端。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-queryprojectdetail-0000001158365067
		officialEndpoint("projects", "queryprojectdetail", "查询项目详情及项目下的应用", "此接口用于查询指定项目的详细项目信息和项目下的应用列表信息。", "GET", "/api/project-service/v1/projects/{projectId}", []string{"projectId"}, "agcapi-queryprojectdetail-0000001158365067", "https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-queryprojectdetail-0000001158365067", "developer-to-huawei"),
		// 查询项目列表：此接口用于查询指定团队账号下的项目列表信息。
		// 请求：GET /api/project-service/v1/projects；调用方向：开发者调用华为服务端。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-queryprojectlist-0000001111685220
		officialEndpoint("projects", "queryprojectlist", "查询项目列表", "此接口用于查询指定团队账号下的项目列表信息。", "GET", "/api/project-service/v1/projects", []string{}, "agcapi-queryprojectlist-0000001111685220", "https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-queryprojectlist-0000001111685220", "developer-to-huawei"),
		// 查询服务开通状态：此接口用于查询项目、应用开通的服务信息列表以及API的状态。
		// 请求：GET /api/cpms/project-service/v1/services/order-api；调用方向：开发者调用华为服务端。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-queryservice-0000001111845108
		officialEndpoint("projects", "queryservice", "查询服务开通状态", "此接口用于查询项目、应用开通的服务信息列表以及API的状态。", "GET", "/api/cpms/project-service/v1/services/order-api", []string{}, "agcapi-queryservice-0000001111845108", "https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-queryservice-0000001111845108", "developer-to-huawei"),
	},
	"comments": {
		// 查询单个评论详情：获取单个评论详情。
		// 请求：GET /api/reviews/v1/manage/dev/reviews/{reviewId}；调用方向：开发者调用华为服务端。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-com-getreviewinfo-0000001116028248
		officialEndpoint("comments", "com-getreviewinfo", "查询单个评论详情", "获取单个评论详情。", "GET", "/api/reviews/v1/manage/dev/reviews/{reviewId}", []string{"reviewId"}, "agcapi-com-getreviewinfo-0000001116028248", "https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-com-getreviewinfo-0000001116028248", "developer-to-huawei"),
		// 查询单个评论详情：开发者查询对应鸿蒙应用下的单个评论详情及评论下的回复列表。
		// 请求：GET /api/marketing-api/v2/reviews/manage/dev/reviews/{reviewId}；调用方向：开发者调用华为服务端。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-com-getreviewinfo-harmonyos-0000002504053839
		officialEndpoint("comments", "com-getreviewinfo-harmonyos", "查询单个评论详情", "开发者查询对应鸿蒙应用下的单个评论详情及评论下的回复列表。", "GET", "/api/marketing-api/v2/reviews/manage/dev/reviews/{reviewId}", []string{"reviewId"}, "agcapi-com-getreviewinfo-harmonyos-0000002504053839", "https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-com-getreviewinfo-harmonyos-0000002504053839", "developer-to-huawei"),
		// 查询应用评分：开发者查询对应应用下的评分列表及综合评分。
		// 请求：GET /api/reviews/v1/manage/dev/ratings；调用方向：开发者调用华为服务端。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-com-rating-0000001133889426
		officialEndpoint("comments", "com-rating", "查询应用评分", "开发者查询对应应用下的评分列表及综合评分。", "GET", "/api/reviews/v1/manage/dev/ratings", []string{}, "agcapi-com-rating-0000001133889426", "https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-com-rating-0000001133889426", "developer-to-huawei"),
		// 查询应用评分列表：开发者查询对应鸿蒙应用下的评分列表。
		// 请求：GET /api/marketing-api/v2/reviews/manage/dev/ratings；调用方向：开发者调用华为服务端。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-com-rating-harmonyos-0000002503933921
		officialEndpoint("comments", "com-rating-harmonyos", "查询应用评分列表", "开发者查询对应鸿蒙应用下的评分列表。", "GET", "/api/marketing-api/v2/reviews/manage/dev/ratings", []string{}, "agcapi-com-rating-harmonyos-0000002503933921", "https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-com-rating-harmonyos-0000002503933921", "developer-to-huawei"),
		// 回复评论和回复：开发者在自己的应用下回复用户的评论，并继续回复用户的回复。
		// 请求：POST /api/reviews/v1/manage/dev/reviews；调用方向：开发者调用华为服务端。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-com-reviewreply-0000001162548117
		officialEndpoint("comments", "com-reviewreply", "回复评论和回复", "开发者在自己的应用下回复用户的评论，并继续回复用户的回复。", "POST", "/api/reviews/v1/manage/dev/reviews", []string{}, "agcapi-com-reviewreply-0000001162548117", "https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-com-reviewreply-0000001162548117", "developer-to-huawei"),
		// 回复评论和回复：开发者在自己的鸿蒙应用下回复用户的评论，并支持继续回复用户的回复。
		// 请求：POST /api/marketing-api/v2/；调用方向：开发者调用华为服务端。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-com-reviewreply-harmonyos-0000002471053974
		officialEndpoint("comments", "com-reviewreply-harmonyos", "回复评论和回复", "开发者在自己的鸿蒙应用下回复用户的评论，并支持继续回复用户的回复。", "POST", "/api/marketing-api/v2/", []string{}, "agcapi-com-reviewreply-harmonyos-0000002471053974", "https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-com-reviewreply-harmonyos-0000002471053974", "developer-to-huawei"),
		// 查询评论列表：开发者查询对应应用下的评论列表。
		// 请求：GET /api/reviews/v1/manage/dev/reviews；调用方向：开发者调用华为服务端。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-comapi-getreviews-0000001162468139
		officialEndpoint("comments", "comapi-getreviews", "查询评论列表", "开发者查询对应应用下的评论列表。", "GET", "/api/reviews/v1/manage/dev/reviews", []string{}, "agcapi-comapi-getreviews-0000001162468139", "https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-comapi-getreviews-0000001162468139", "developer-to-huawei"),
		// 查询应用评论列表：开发者查询对应鸿蒙应用下的评论列表。
		// 请求：GET /api/marketing-api/v2/reviews/manage/dev/reviews；调用方向：开发者调用华为服务端。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-comapi-getreviews-harmonyos-0000002470893976
		officialEndpoint("comments", "comapi-getreviews-harmonyos", "查询应用评论列表", "开发者查询对应鸿蒙应用下的评论列表。", "GET", "/api/marketing-api/v2/reviews/manage/dev/reviews", []string{}, "agcapi-comapi-getreviews-harmonyos-0000002470893976", "https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-comapi-getreviews-harmonyos-0000002470893976", "developer-to-huawei"),
	},
	"pms": {
		// 创建商品：此接口用于创建商品定义信息，包含商品价格，展示语言等。
		// 请求：POST /api/pms/product-price-service/v1/manage/product；调用方向：开发者调用华为服务端。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-addproduct-android-0000002171407001
		officialEndpoint("pms", "addproduct-android", "创建商品", "此接口用于创建商品定义信息，包含商品价格，展示语言等。", "POST", "/api/pms/product-price-service/v1/manage/product", []string{}, "agcapi-addproduct-android-0000002171407001", "https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-addproduct-android-0000002171407001", "developer-to-huawei"),
		// 创建商品：此接口用于创建商品定义信息，包含商品价格，展示语言等。
		// 请求：POST /api/pms/product-price-service/v2/manage/product；调用方向：开发者调用华为服务端。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-addproduct-harmonyosnext-0000002131508844
		officialEndpoint("pms", "addproduct-harmonyosnext", "创建商品", "此接口用于创建商品定义信息，包含商品价格，展示语言等。", "POST", "/api/pms/product-price-service/v2/manage/product", []string{}, "agcapi-addproduct-harmonyosnext-0000002131508844", "https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-addproduct-harmonyosnext-0000002131508844", "developer-to-huawei"),
		// 创建商品订阅分组信息：此接口用于创建商品订阅分组信息。
		// 请求：POST /api/pms/product-price-service/v1/manage/product/group；调用方向：开发者调用华为服务端。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-addproductgroup-android-0000002135929308
		officialEndpoint("pms", "addproductgroup-android", "创建商品订阅分组信息", "此接口用于创建商品订阅分组信息。", "POST", "/api/pms/product-price-service/v1/manage/product/group", []string{}, "agcapi-addproductgroup-android-0000002135929308", "https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-addproductgroup-android-0000002135929308", "developer-to-huawei"),
		// 创建商品订阅分组信息：此接口用于创建商品订阅分组信息。
		// 请求：POST /api/pms/product-price-service/v2/manage/product/group；调用方向：开发者调用华为服务端。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-addproductgroup-harmonyosnext-0000002166748605
		officialEndpoint("pms", "addproductgroup-harmonyosnext", "创建商品订阅分组信息", "此接口用于创建商品订阅分组信息。", "POST", "/api/pms/product-price-service/v2/manage/product/group", []string{}, "agcapi-addproductgroup-harmonyosnext-0000002166748605", "https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-addproductgroup-harmonyosnext-0000002166748605", "developer-to-huawei"),
		// 创建商品促销信息：此接口用于创建商品促销信息。
		// 请求：POST /api/pms/product-price-service/v1/manage/product/promotion；调用方向：开发者调用华为服务端。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-addpromotion-android-0000002171407021
		officialEndpoint("pms", "addpromotion-android", "创建商品促销信息", "此接口用于创建商品促销信息。", "POST", "/api/pms/product-price-service/v1/manage/product/promotion", []string{}, "agcapi-addpromotion-android-0000002171407021", "https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-addpromotion-android-0000002171407021", "developer-to-huawei"),
		// 创建商品促销信息：此接口用于创建商品促销信息。
		// 请求：POST /api/pms/product-price-service/v2/manage/product/promotion；调用方向：开发者调用华为服务端。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-addpromotion-harmonyosnext-0000002131508856
		officialEndpoint("pms", "addpromotion-harmonyosnext", "创建商品促销信息", "此接口用于创建商品促销信息。", "POST", "/api/pms/product-price-service/v2/manage/product/promotion", []string{}, "agcapi-addpromotion-harmonyosnext-0000002131508856", "https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-addpromotion-harmonyosnext-0000002131508856", "developer-to-huawei"),
		// 查询关联APP已有商品类型：此接口用于查询关联APP已有商品类型。
		// 请求：GET /api/pms/product-price-service/v2/manage/application/basic；调用方向：开发者调用华为服务端。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-basicapplication-harmonyosnext-0000002131519928
		officialEndpoint("pms", "basicapplication-harmonyosnext", "查询关联APP已有商品类型", "此接口用于查询关联APP已有商品类型。", "GET", "/api/pms/product-price-service/v2/manage/application/basic", []string{}, "agcapi-basicapplication-harmonyosnext-0000002131519928", "https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-basicapplication-harmonyosnext-0000002131519928", "developer-to-huawei"),
		// 批量激活商品：此接口用于批量激活商品。
		// 请求：PUT /api/pms/product-price-service/v1/manage/product/batchActive；调用方向：开发者调用华为服务端。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-batchactproduct-android-0000002171288617
		officialEndpoint("pms", "batchactproduct-android", "批量激活商品", "此接口用于批量激活商品。", "PUT", "/api/pms/product-price-service/v1/manage/product/batchActive", []string{}, "agcapi-batchactproduct-android-0000002171288617", "https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-batchactproduct-android-0000002171288617", "developer-to-huawei"),
		// 批量激活商品：此接口用于批量激活商品 ， 生效已下线商品，已结束的促销不会重新生效 。
		// 请求：PUT /api/pms/product-price-service/v2/manage/product/batchActive；调用方向：开发者调用华为服务端。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-batchactproduct-harmonyosnext-0000002131350712
		officialEndpoint("pms", "batchactproduct-harmonyosnext", "批量激活商品", "此接口用于批量激活商品 ， 生效已下线商品，已结束的促销不会重新生效 。", "PUT", "/api/pms/product-price-service/v2/manage/product/batchActive", []string{}, "agcapi-batchactproduct-harmonyosnext-0000002131350712", "https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-batchactproduct-harmonyosnext-0000002131350712", "developer-to-huawei"),
		// 批量上线商品促销活动：此接口用于批量上线商品促销活动。
		// 请求：PUT /api/pms/product-price-service/v1/manage/product/promotion/batchActive；调用方向：开发者调用华为服务端。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-batchactpromotion-android-0000002171407025
		officialEndpoint("pms", "batchactpromotion-android", "批量上线商品促销活动", "此接口用于批量上线商品促销活动。", "PUT", "/api/pms/product-price-service/v1/manage/product/promotion/batchActive", []string{}, "agcapi-batchactpromotion-android-0000002171407025", "https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-batchactpromotion-android-0000002171407025", "developer-to-huawei"),
		// 批量上线商品促销活动：此接口用于批量上线商品促销活动。
		// 请求：PUT /api/pms/product-price-service/v2/manage/product/promotion/batchActive；调用方向：开发者调用华为服务端。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-batchactpromotion-harmonyosnext-0000002131508860
		officialEndpoint("pms", "batchactpromotion-harmonyosnext", "批量上线商品促销活动", "此接口用于批量上线商品促销活动。", "PUT", "/api/pms/product-price-service/v2/manage/product/promotion/batchActive", []string{}, "agcapi-batchactpromotion-harmonyosnext-0000002131508860", "https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-batchactpromotion-harmonyosnext-0000002131508860", "developer-to-huawei"),
		// 批量创建商品：此接口用于批量创建商品定义信息，展示语言，商品价格，促销活动信息。
		// 请求：POST /api/pms/product-price-service/v1/manage/product/batchImportProducts；调用方向：开发者调用华为服务端。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-batchaddproduct-android-0000002171288605
		officialEndpoint("pms", "batchaddproduct-android", "批量创建商品", "此接口用于批量创建商品定义信息，展示语言，商品价格，促销活动信息。", "POST", "/api/pms/product-price-service/v1/manage/product/batchImportProducts", []string{}, "agcapi-batchaddproduct-android-0000002171288605", "https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-batchaddproduct-android-0000002171288605", "developer-to-huawei"),
		// 批量创建商品：此接口用于批量创建商品定义信息，展示语言，商品价格，促销活动信息。
		// 请求：POST /api/pms/product-price-service/v2/manage/product/batchImportProducts；调用方向：开发者调用华为服务端。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-batchaddproduct-harmonyosnext-0000002131350704
		officialEndpoint("pms", "batchaddproduct-harmonyosnext", "批量创建商品", "此接口用于批量创建商品定义信息，展示语言，商品价格，促销活动信息。", "POST", "/api/pms/product-price-service/v2/manage/product/batchImportProducts", []string{}, "agcapi-batchaddproduct-harmonyosnext-0000002131350704", "https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-batchaddproduct-harmonyosnext-0000002131350704", "developer-to-huawei"),
		// 批量去激活商品：此接口用于批量去激活商品。
		// 请求：PUT /api/pms/product-price-service/v1/manage/product/batchDeactive；调用方向：开发者调用华为服务端。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-batchdeactproduct-android-0000002135929320
		officialEndpoint("pms", "batchdeactproduct-android", "批量去激活商品", "此接口用于批量去激活商品。", "PUT", "/api/pms/product-price-service/v1/manage/product/batchDeactive", []string{}, "agcapi-batchdeactproduct-android-0000002135929320", "https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-batchdeactproduct-android-0000002135929320", "developer-to-huawei"),
		// 批量去激活商品：此接口用于批量去激活商品 ，下线已生效商品，结束促销活动 。
		// 请求：PUT /api/pms/product-price-service/v2/manage/product/batchDeactive；调用方向：开发者调用华为服务端。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-batchdeactproduct-harmonyosnext-0000002166748613
		officialEndpoint("pms", "batchdeactproduct-harmonyosnext", "批量去激活商品", "此接口用于批量去激活商品 ，下线已生效商品，结束促销活动 。", "PUT", "/api/pms/product-price-service/v2/manage/product/batchDeactive", []string{}, "agcapi-batchdeactproduct-harmonyosnext-0000002166748613", "https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-batchdeactproduct-harmonyosnext-0000002166748613", "developer-to-huawei"),
		// 批量下线商品促销活动：此接口用于批量下线商品促销活动。
		// 请求：PUT /api/pms/product-price-service/v1/manage/product/promotion/batchDeactive；调用方向：开发者调用华为服务端。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-batchdeactpromotion-android-0000002171288629
		officialEndpoint("pms", "batchdeactpromotion-android", "批量下线商品促销活动", "此接口用于批量下线商品促销活动。", "PUT", "/api/pms/product-price-service/v1/manage/product/promotion/batchDeactive", []string{}, "agcapi-batchdeactpromotion-android-0000002171288629", "https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-batchdeactpromotion-android-0000002171288629", "developer-to-huawei"),
		// 批量下线商品促销活动：此接口用于批量下线商品促销活动。
		// 请求：PUT /api/pms/product-price-service/v2/manage/product/promotion/batchDeactive；调用方向：开发者调用华为服务端。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-batchdeactpromotion-harmonyosnext-0000002131350720
		officialEndpoint("pms", "batchdeactpromotion-harmonyosnext", "批量下线商品促销活动", "此接口用于批量下线商品促销活动。", "PUT", "/api/pms/product-price-service/v2/manage/product/promotion/batchDeactive", []string{}, "agcapi-batchdeactpromotion-harmonyosnext-0000002131350720", "https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-batchdeactpromotion-harmonyosnext-0000002131350720", "developer-to-huawei"),
		// 批量更新商品：此接口用于批量更新商品定义信息，展示语言，商品价格，促销活动信息。
		// 请求：PUT /api/pms/product-price-service/v1/manage/product/batchImportProducts；调用方向：开发者调用华为服务端。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-batchupdateproduct-android-0000002171407013
		officialEndpoint("pms", "batchupdateproduct-android", "批量更新商品", "此接口用于批量更新商品定义信息，展示语言，商品价格，促销活动信息。", "PUT", "/api/pms/product-price-service/v1/manage/product/batchImportProducts", []string{}, "agcapi-batchupdateproduct-android-0000002171407013", "https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-batchupdateproduct-android-0000002171407013", "developer-to-huawei"),
		// 批量更新商品：此接口用于批量更新商品定义信息，展示语言，商品价格，促销活动信息。
		// 请求：PUT /api/pms/product-price-service/v2/manage/product/batchImportProducts；调用方向：开发者调用华为服务端。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-batchupdateproduct-harmonyosnext-0000002131508852
		officialEndpoint("pms", "batchupdateproduct-harmonyosnext", "批量更新商品", "此接口用于批量更新商品定义信息，展示语言，商品价格，促销活动信息。", "PUT", "/api/pms/product-price-service/v2/manage/product/batchImportProducts", []string{}, "agcapi-batchupdateproduct-harmonyosnext-0000002131508852", "https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-batchupdateproduct-harmonyosnext-0000002131508852", "developer-to-huawei"),
		// 按条件查询商品信息：此接口用于商品查询，支持多条件过滤，根据productNo和AppId唯一查询一个商品信息。
		// 请求：POST /api/pms/product-price-service/v1/manage/product/list；调用方向：开发者调用华为服务端。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-bygetproductinfo-android-0000002135929316
		officialEndpoint("pms", "bygetproductinfo-android", "按条件查询商品信息", "此接口用于商品查询，支持多条件过滤，根据productNo和AppId唯一查询一个商品信息。", "POST", "/api/pms/product-price-service/v1/manage/product/list", []string{}, "agcapi-bygetproductinfo-android-0000002135929316", "https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-bygetproductinfo-android-0000002135929316", "developer-to-huawei"),
		// 按条件查询商品信息：此接口用于商品查询，支持多条件过滤，根据productNo和AppId唯一查询一个商品信息。
		// 请求：POST /api/pms/product-price-service/v2/manage/product/list；调用方向：开发者调用华为服务端。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-bygetproductinfo-harmonyosnext-0000002166748609
		officialEndpoint("pms", "bygetproductinfo-harmonyosnext", "按条件查询商品信息", "此接口用于商品查询，支持多条件过滤，根据productNo和AppId唯一查询一个商品信息。", "POST", "/api/pms/product-price-service/v2/manage/product/list", []string{}, "agcapi-bygetproductinfo-harmonyosnext-0000002166748609", "https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-bygetproductinfo-harmonyosnext-0000002166748609", "developer-to-huawei"),
		// 按条件查询商品促销信息：此接口用于商品促销查询，支持多条件过滤。
		// 请求：POST /api/pms/product-price-service/v1/manage/product/promotion/list；调用方向：开发者调用华为服务端。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-bygetpromotioninfo-android-0000002136087416
		officialEndpoint("pms", "bygetpromotioninfo-android", "按条件查询商品促销信息", "此接口用于商品促销查询，支持多条件过滤。", "POST", "/api/pms/product-price-service/v1/manage/product/promotion/list", []string{}, "agcapi-bygetpromotioninfo-android-0000002136087416", "https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-bygetpromotioninfo-android-0000002136087416", "developer-to-huawei"),
		// 按条件查询商品促销信息：此接口用于商品促销查询，支持多条件过滤。
		// 请求：POST /api/pms/product-price-service/v2/manage/product/promotion/list；调用方向：开发者调用华为服务端。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-bygetpromotioninfo-harmonyosnext-0000002166670149
		officialEndpoint("pms", "bygetpromotioninfo-harmonyosnext", "按条件查询商品促销信息", "此接口用于商品促销查询，支持多条件过滤。", "POST", "/api/pms/product-price-service/v2/manage/product/promotion/list", []string{}, "agcapi-bygetpromotioninfo-harmonyosnext-0000002166670149", "https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-bygetpromotioninfo-harmonyosnext-0000002166670149", "developer-to-huawei"),
		// 查询商品订阅分组信息：此接口用于查询商品订阅分组信息。
		// 请求：POST /api/pms/product-price-service/v1/manage/product/group/list；调用方向：开发者调用华为服务端。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-getproductgroup-android-0000002171407009
		officialEndpoint("pms", "getproductgroup-android", "查询商品订阅分组信息", "此接口用于查询商品订阅分组信息。", "POST", "/api/pms/product-price-service/v1/manage/product/group/list", []string{}, "agcapi-getproductgroup-android-0000002171407009", "https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-getproductgroup-android-0000002171407009", "developer-to-huawei"),
		// 查询商品订阅分组信息：此接口用于查询商品订阅分组信息。
		// 请求：POST /api/pms/product-price-service/v2/manage/product/group/list；调用方向：开发者调用华为服务端。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-getproductgroup-harmonyosnext-0000002131508848
		officialEndpoint("pms", "getproductgroup-harmonyosnext", "查询商品订阅分组信息", "此接口用于查询商品订阅分组信息。", "POST", "/api/pms/product-price-service/v2/manage/product/group/list", []string{}, "agcapi-getproductgroup-harmonyosnext-0000002131508848", "https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-getproductgroup-harmonyosnext-0000002131508848", "developer-to-huawei"),
		// 查询商品详情：此接口用于商品详情查询，根据productNo和appId唯一查询一个商品信息。
		// 请求：GET /api/pms/product-price-service/v1/manage/product；调用方向：开发者调用华为服务端。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-getproductinfo-android-0000002171288609
		officialEndpoint("pms", "getproductinfo-android", "查询商品详情", "此接口用于商品详情查询，根据productNo和appId唯一查询一个商品信息。", "GET", "/api/pms/product-price-service/v1/manage/product", []string{}, "agcapi-getproductinfo-android-0000002171288609", "https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-getproductinfo-android-0000002171288609", "developer-to-huawei"),
		// 查询商品详情：此接口用于商品详情查询，根据productNo和appId唯一查询一个商品信息。
		// 请求：GET /api/pms/product-price-service/v2/manage/product；调用方向：开发者调用华为服务端。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-getproductinfo-harmonyosnext-0000002131350708
		officialEndpoint("pms", "getproductinfo-harmonyosnext", "查询商品详情", "此接口用于商品详情查询，根据productNo和appId唯一查询一个商品信息。", "GET", "/api/pms/product-price-service/v2/manage/product", []string{}, "agcapi-getproductinfo-harmonyosnext-0000002131350708", "https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-getproductinfo-harmonyosnext-0000002131350708", "developer-to-huawei"),
		// 查询商品促销详情：此接口用于查询商品促销信息。
		// 请求：GET /api/pms/product-price-service/v1/manage/product/promotion；调用方向：开发者调用华为服务端。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-getpromotioninfo-android-0000002135929328
		officialEndpoint("pms", "getpromotioninfo-android", "查询商品促销详情", "此接口用于查询商品促销信息。", "GET", "/api/pms/product-price-service/v1/manage/product/promotion", []string{}, "agcapi-getpromotioninfo-android-0000002135929328", "https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-getpromotioninfo-android-0000002135929328", "developer-to-huawei"),
		// 查询商品促销详情：此接口用于查询商品促销信息。
		// 请求：GET /api/pms/product-price-service/v2/manage/product/promotion；调用方向：开发者调用华为服务端。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-getpromotioninfo-harmonyosnext-0000002166748617
		officialEndpoint("pms", "getpromotioninfo-harmonyosnext", "查询商品促销详情", "此接口用于查询商品促销信息。", "GET", "/api/pms/product-price-service/v2/manage/product/promotion", []string{}, "agcapi-getpromotioninfo-harmonyosnext-0000002166748617", "https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-getpromotioninfo-harmonyosnext-0000002166748617", "developer-to-huawei"),
		// 查询是否有在线商品：此接口用于根据AppId查询是否有激活的商品。
		// 请求：GET /api/pms/product-price-service/v1/manage/product/check；调用方向：开发者调用华为服务端。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-isproductactive-android-0000002136087412
		officialEndpoint("pms", "isproductactive-android", "查询是否有在线商品", "此接口用于根据AppId查询是否有激活的商品。", "GET", "/api/pms/product-price-service/v1/manage/product/check", []string{}, "agcapi-isproductactive-android-0000002136087412", "https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-isproductactive-android-0000002136087412", "developer-to-huawei"),
		// 查询是否有在线商品：此接口用于根据AppId查询是否有激活的商品。
		// 请求：GET /api/pms/product-price-service/v2/manage/product/check；调用方向：开发者调用华为服务端。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-isproductactive-harmonyosnext-0000002166670145
		officialEndpoint("pms", "isproductactive-harmonyosnext", "查询是否有在线商品", "此接口用于根据AppId查询是否有激活的商品。", "GET", "/api/pms/product-price-service/v2/manage/product/check", []string{}, "agcapi-isproductactive-harmonyosnext-0000002166670145", "https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-isproductactive-harmonyosnext-0000002166670145", "developer-to-huawei"),
		// 更新商品订阅分组等级信息：此接口用于更新商品订阅分组等级信息。
		// 请求：PUT /api/pms/product-price-service/v2/manage/product/group/level；调用方向：开发者调用华为服务端。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-levelgroup-harmonyosnext-0000002175107485
		officialEndpoint("pms", "levelgroup-harmonyosnext", "更新商品订阅分组等级信息", "此接口用于更新商品订阅分组等级信息。", "PUT", "/api/pms/product-price-service/v2/manage/product/group/level", []string{}, "agcapi-levelgroup-harmonyosnext-0000002175107485", "https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-levelgroup-harmonyosnext-0000002175107485", "developer-to-huawei"),
		// 上传资源：此接口用于提交上传资源。
		// 请求：POST /api/pms/product-price-service/v2/manage/resource；调用方向：开发者调用华为服务端。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-resourcemanage-harmonyosnext-0000002131361792
		officialEndpoint("pms", "resourcemanage-harmonyosnext", "上传资源", "此接口用于提交上传资源。", "POST", "/api/pms/product-price-service/v2/manage/resource", []string{}, "agcapi-resourcemanage-harmonyosnext-0000002131361792", "https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-resourcemanage-harmonyosnext-0000002131361792", "developer-to-huawei"),
		// 提交商品审核：此接口用于提交商品审核。
		// 请求：POST /api/pms/product-price-service/v2/manage/product/review；调用方向：开发者调用华为服务端。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-reviewproduct-harmonyosnext-0000002166681241
		officialEndpoint("pms", "reviewproduct-harmonyosnext", "提交商品审核", "此接口用于提交商品审核。", "POST", "/api/pms/product-price-service/v2/manage/product/review", []string{}, "agcapi-reviewproduct-harmonyosnext-0000002166681241", "https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-reviewproduct-harmonyosnext-0000002166681241", "developer-to-huawei"),
		// 更新商品信息：此接口用于更新商品定义信息，包含商品价格，展示语言和商品状态。
		// 请求：PUT /api/pms/product-price-service/v1/manage/product；调用方向：开发者调用华为服务端。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-updateproduct-android-0000002136087404
		officialEndpoint("pms", "updateproduct-android", "更新商品信息", "此接口用于更新商品定义信息，包含商品价格，展示语言和商品状态。", "PUT", "/api/pms/product-price-service/v1/manage/product", []string{}, "agcapi-updateproduct-android-0000002136087404", "https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-updateproduct-android-0000002136087404", "developer-to-huawei"),
		// 更新商品信息：此接口用于更新商品定义信息，包含商品价格，展示语言和商品状态。
		// 请求：PUT /api/pms/product-price-service/v2/manage/product；调用方向：开发者调用华为服务端。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-updateproduct-harmonyosnext-0000002166670141
		officialEndpoint("pms", "updateproduct-harmonyosnext", "更新商品信息", "此接口用于更新商品定义信息，包含商品价格，展示语言和商品状态。", "PUT", "/api/pms/product-price-service/v2/manage/product", []string{}, "agcapi-updateproduct-harmonyosnext-0000002166670141", "https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-updateproduct-harmonyosnext-0000002166670141", "developer-to-huawei"),
		// 更新商品订阅分组信息：此接口用于更新商品订阅分组信息。
		// 请求：PUT /api/pms/product-price-service/v1/manage/product/group；调用方向：开发者调用华为服务端。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-updateproductgroup-android-0000002136087400
		officialEndpoint("pms", "updateproductgroup-android", "更新商品订阅分组信息", "此接口用于更新商品订阅分组信息。", "PUT", "/api/pms/product-price-service/v1/manage/product/group", []string{}, "agcapi-updateproductgroup-android-0000002136087400", "https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-updateproductgroup-android-0000002136087400", "developer-to-huawei"),
		// 更新商品订阅分组信息：此接口用于更新商品订阅分组信息。
		// 请求：PUT /api/pms/product-price-service/v2/manage/product/group；调用方向：开发者调用华为服务端。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-updateproductgroup-harmonyosnext-0000002166670137
		officialEndpoint("pms", "updateproductgroup-harmonyosnext", "更新商品订阅分组信息", "此接口用于更新商品订阅分组信息。", "PUT", "/api/pms/product-price-service/v2/manage/product/group", []string{}, "agcapi-updateproductgroup-harmonyosnext-0000002166670137", "https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-updateproductgroup-harmonyosnext-0000002166670137", "developer-to-huawei"),
		// 更新商品促销信息：此接口用于更新商品促销信息。
		// 请求：PUT /api/pms/product-price-service/v1/manage/product/promotion；调用方向：开发者调用华为服务端。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-updatepromotion-android-0000002171288621
		officialEndpoint("pms", "updatepromotion-android", "更新商品促销信息", "此接口用于更新商品促销信息。", "PUT", "/api/pms/product-price-service/v1/manage/product/promotion", []string{}, "agcapi-updatepromotion-android-0000002171288621", "https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-updatepromotion-android-0000002171288621", "developer-to-huawei"),
		// 更新商品促销信息：此接口用于更新商品促销信息。
		// 请求：PUT /api/pms/product-price-service/v2/manage/product/promotion；调用方向：开发者调用华为服务端。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-updatepromotion-harmonyosnext-0000002131350716
		officialEndpoint("pms", "updatepromotion-harmonyosnext", "更新商品促销信息", "此接口用于更新商品促销信息。", "PUT", "/api/pms/product-price-service/v2/manage/product/promotion", []string{}, "agcapi-updatepromotion-harmonyosnext-0000002131350716", "https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-updatepromotion-harmonyosnext-0000002131350716", "developer-to-huawei"),
	},
	"gameplay": {
		// 查询资产明细：华为游戏服务器调用此接口查询玩家在玩游戏的资产明细数据，并将信息展示在游戏中心的“在玩”页面上。
		// 请求：POST {callbackUrl}；调用方向：华为调用开发者回调服务；本项目的请求构造不等于已实现回调接收与验签。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-asset-details-0000002022225253
		officialEndpoint("gameplay", "asset-details", "查询资产明细", "华为游戏服务器调用此接口查询玩家在玩游戏的资产明细数据，并将信息展示在游戏中心的“在玩”页面上。", "POST", "{callbackUrl}", []string{"callbackUrl", "callbackUrl"}, "agcapi-asset-details-0000002022225253", "https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-asset-details-0000002022225253", "inbound-callback"),
		// 查询资产概要：华为游戏服务器调用此接口查询玩家在玩游戏的资产概要数据，并将信息展示在游戏中心的“在玩”页面上。
		// 请求：POST {callbackUrl}；调用方向：华为调用开发者回调服务；本项目的请求构造不等于已实现回调接收与验签。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-asset-summary-0000002022105729
		officialEndpoint("gameplay", "asset-summary", "查询资产概要", "华为游戏服务器调用此接口查询玩家在玩游戏的资产概要数据，并将信息展示在游戏中心的“在玩”页面上。", "POST", "{callbackUrl}", []string{"callbackUrl", "callbackUrl"}, "agcapi-asset-summary-0000002022105729", "https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-asset-summary-0000002022105729", "inbound-callback"),
		// 批量查询玩家基本信息：华为游戏服务器调用此接口批量查询玩家基本信息，并将信息展示在游戏中心的“在玩”页面上。
		// 请求：POST {callbackUrl}；调用方向：华为调用开发者回调服务；本项目的请求构造不等于已实现回调接收与验签。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-basic-info-0000002038975322
		officialEndpoint("gameplay", "basic-info", "批量查询玩家基本信息", "华为游戏服务器调用此接口批量查询玩家基本信息，并将信息展示在游戏中心的“在玩”页面上。", "POST", "{callbackUrl}", []string{"callbackUrl", "callbackUrl"}, "agcapi-basic-info-0000002038975322", "https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-basic-info-0000002038975322", "inbound-callback"),
		// 查询战绩筛选标签：华为游戏服务器调用此接口查询玩家在玩游戏战绩的筛选标签，并将信息展示在游戏中心的“在玩”页面上。
		// 请求：POST {callbackUrl}；调用方向：华为调用开发者回调服务；本项目的请求构造不等于已实现回调接收与验签。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-label-0000001985625670
		officialEndpoint("gameplay", "label", "查询战绩筛选标签", "华为游戏服务器调用此接口查询玩家在玩游戏战绩的筛选标签，并将信息展示在游戏中心的“在玩”页面上。", "POST", "{callbackUrl}", []string{"callbackUrl", "callbackUrl"}, "agcapi-label-0000001985625670", "https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-label-0000001985625670", "inbound-callback"),
		// 上传/下架游戏资源：开发者服务器调用此接口同步华为服务器的游戏资源，用于上传或下架“在玩”页面上的图片或视频。
		// 请求：POST /gameservice/api/gbClientApi；调用方向：开发者调用华为服务端。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-resources-0000002022105725
		officialEndpoint("gameplay", "resources", "上传/下架游戏资源", "开发者服务器调用此接口同步华为服务器的游戏资源，用于上传或下架“在玩”页面上的图片或视频。", "POST", "/gameservice/api/gbClientApi", []string{}, "agcapi-resources-0000002022105725", "https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-resources-0000002022105725", "developer-to-huawei"),
		// 批量查询玩家角色信息：华为游戏服务器调用此接口批量查询玩家角色信息。
		// 请求：POST {callbackUrl}；调用方向：华为调用开发者回调服务；本项目的请求构造不等于已实现回调接收与验签。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-role-info-0000002075016173
		officialEndpoint("gameplay", "role-info", "批量查询玩家角色信息", "华为游戏服务器调用此接口批量查询玩家角色信息。", "POST", "{callbackUrl}", []string{"callbackUrl", "callbackUrl"}, "agcapi-role-info-0000002075016173", "https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-role-info-0000002075016173", "inbound-callback"),
		// 分页查询玩家战绩：华为游戏服务器调用此接口分页查询玩家在玩游戏的战绩数据，并将信息展示在游戏中心的“在玩”页面上。
		// 请求：POST {callbackUrl}；调用方向：华为调用开发者回调服务；本项目的请求构造不等于已实现回调接收与验签。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-scores-0000001985465938
		officialEndpoint("gameplay", "scores", "分页查询玩家战绩", "华为游戏服务器调用此接口分页查询玩家在玩游戏的战绩数据，并将信息展示在游戏中心的“在玩”页面上。", "POST", "{callbackUrl}", []string{"callbackUrl", "callbackUrl"}, "agcapi-scores-0000001985465938", "https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-scores-0000001985465938", "inbound-callback"),
		// 查询游戏概要信息：华为游戏服务器调用此接口查询玩家在玩游戏的概要数据，并将信息展示在游戏中心的“在玩”页面上。
		// 请求：POST {callbackUrl}；调用方向：华为调用开发者回调服务；本项目的请求构造不等于已实现回调接收与验签。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-summary-0000002022225249
		officialEndpoint("gameplay", "summary", "查询游戏概要信息", "华为游戏服务器调用此接口查询玩家在玩游戏的概要数据，并将信息展示在游戏中心的“在玩”页面上。", "POST", "{callbackUrl}", []string{"callbackUrl", "callbackUrl"}, "agcapi-summary-0000002022225249", "https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-summary-0000002022225249", "inbound-callback"),
	},
	"game-items": {
		// 查询玩家区服角色：此接口用于查询玩家区服角色信息。
		// 请求：POST {callbackUrl}；调用方向：华为调用开发者回调服务；本项目的请求构造不等于已实现回调接收与验签。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-propapi-getroleinfo-0000002355255377
		officialEndpoint("game-items", "propapi-getroleinfo", "查询玩家区服角色", "此接口用于查询玩家区服角色信息。", "POST", "{callbackUrl}", []string{"callbackUrl", "callbackUrl"}, "agcapi-propapi-getroleinfo-0000002355255377", "https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-propapi-getroleinfo-0000002355255377", "inbound-callback"),
		// 下单：下单
		// 请求：POST {callbackUrl}；调用方向：华为调用开发者回调服务；本项目的请求构造不等于已实现回调接收与验签。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-propapi-order-0000002321296666
		officialEndpoint("game-items", "propapi-order", "下单", "下单", "POST", "{callbackUrl}", []string{"callbackUrl", "callbackUrl"}, "agcapi-propapi-order-0000002321296666", "https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-propapi-order-0000002321296666", "inbound-callback"),
	},
	"resources": {
		// 新增资源包版本：创建“使用华为CDN”方式的资源包预下载任务。
		// 请求：POST /api/games-background-assets-service/v1/open-gw/dev/{devId}/package-versions；调用方向：开发者调用华为服务端。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-respackapi-addversion-0000002328396224
		officialEndpoint("resources", "respackapi-addversion", "新增资源包版本", "创建“使用华为CDN”方式的资源包预下载任务。", "POST", "/api/games-background-assets-service/v1/open-gw/dev/{devId}/package-versions", []string{"devId"}, "agcapi-respackapi-addversion-0000002328396224", "https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-respackapi-addversion-0000002328396224", "developer-to-huawei"),
		// 删除资源包文件：删除资源包任务文件，要求该资源包任务为下线状态。
		// 请求：DELETE /api/games-background-assets-service/v1/open-gw/dev/{devId}/package-files；调用方向：开发者调用华为服务端。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-respackapi-deletefile-0000002328396228
		officialEndpoint("resources", "respackapi-deletefile", "删除资源包文件", "删除资源包任务文件，要求该资源包任务为下线状态。", "DELETE", "/api/games-background-assets-service/v1/open-gw/dev/{devId}/package-files", []string{"devId"}, "agcapi-respackapi-deletefile-0000002328396228", "https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-respackapi-deletefile-0000002328396228", "developer-to-huawei"),
		// 查询资源包版本：查询指定游戏下所有资源包版本的详细信息。
		// 请求：GET /api/games-background-assets-service/v1/open-gw/dev/{devId}/package-versions；调用方向：开发者调用华为服务端。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-respackapi-getversion-0000002362274745
		officialEndpoint("resources", "respackapi-getversion", "查询资源包版本", "查询指定游戏下所有资源包版本的详细信息。", "GET", "/api/games-background-assets-service/v1/open-gw/dev/{devId}/package-versions", []string{"devId"}, "agcapi-respackapi-getversion-0000002362274745", "https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-respackapi-getversion-0000002362274745", "developer-to-huawei"),
		// 修改资源包版本：若isDraft传参为0，提交资源包预下载任务。
		// 请求：PUT /api/games-background-assets-service/v1/open-gw/dev/{devId}/package-versions；调用方向：开发者调用华为服务端。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-respackapi-modifyversion-0000002328236372
		officialEndpoint("resources", "respackapi-modifyversion", "修改资源包版本", "若isDraft传参为0，提交资源包预下载任务。", "PUT", "/api/games-background-assets-service/v1/open-gw/dev/{devId}/package-versions", []string{"devId"}, "agcapi-respackapi-modifyversion-0000002328236372", "https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-respackapi-modifyversion-0000002328236372", "developer-to-huawei"),
		// 查询资源包文件：查询资源包文件。
		// 请求：GET /api/games-background-assets-service/v1/open-gw/dev/{devId}/package-files；调用方向：开发者调用华为服务端。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-respackapi-queryfile-0000002362274749
		officialEndpoint("resources", "respackapi-queryfile", "查询资源包文件", "查询资源包文件。", "GET", "/api/games-background-assets-service/v1/open-gw/dev/{devId}/package-files", []string{"devId"}, "agcapi-respackapi-queryfile-0000002362274749", "https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-respackapi-queryfile-0000002362274749", "developer-to-huawei"),
		// 发布资源包版本：发布资源包下载任务。
		// 请求：POST /api/games-background-assets-service/v1/open-gw/dev/{devId}/package-version/publish；调用方向：开发者调用华为服务端。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-respackapi-releaseversion-0000002362394637
		officialEndpoint("resources", "respackapi-releaseversion", "发布资源包版本", "发布资源包下载任务。", "POST", "/api/games-background-assets-service/v1/open-gw/dev/{devId}/package-version/publish", []string{"devId"}, "agcapi-respackapi-releaseversion-0000002362394637", "https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-respackapi-releaseversion-0000002362394637", "developer-to-huawei"),
		// 资源包文件上传申请：获取上传地址和必要参数，根据申请返回的信息完成文件上传。
		// 请求：POST /api/games-background-assets-service/v1/open-gw/dev/{devId}/package-file/apply；调用方向：开发者调用华为服务端。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-respackapi-uploadapplyfor-0000002328236376
		officialEndpoint("resources", "respackapi-uploadapplyfor", "资源包文件上传申请", "获取上传地址和必要参数，根据申请返回的信息完成文件上传。", "POST", "/api/games-background-assets-service/v1/open-gw/dev/{devId}/package-file/apply", []string{"devId"}, "agcapi-respackapi-uploadapplyfor-0000002328236376", "https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-respackapi-uploadapplyfor-0000002328236376", "developer-to-huawei"),
		// 资源包文件上传确认：确认资源包文件已上传。
		// 请求：POST /api/games-background-assets-service/v1/open-gw/dev/{devId}/package-file/confirm；调用方向：开发者调用华为服务端。
		// 官方文档：https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-respackapi-uploadconfirm-0000002362394641
		officialEndpoint("resources", "respackapi-uploadconfirm", "资源包文件上传确认", "确认资源包文件已上传。", "POST", "/api/games-background-assets-service/v1/open-gw/dev/{devId}/package-file/confirm", []string{"devId"}, "agcapi-respackapi-uploadconfirm-0000002362394641", "https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-respackapi-uploadconfirm-0000002362394641", "developer-to-huawei"),
	},
	"cicd": {
		// 本地 Hvigor 构建：Run local HarmonyOS Hvigor build and return artifact metadata.
		// 请求：POST /local/hvigor/build；调用方向：本地构建集成；不是华为远端接口。
		// 项目说明：docs/AGC_CLI_FULL_PLAN.md
		officialEndpoint("cicd", "build-hvigor", "本地 Hvigor 构建", "Run local HarmonyOS Hvigor build and return artifact metadata.", "POST", "/local/hvigor/build", []string{"project"}, "local-hvigor-build", "docs/AGC_CLI_FULL_PLAN.md", "local"),
	},
}
