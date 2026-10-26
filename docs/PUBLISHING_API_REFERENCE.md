# Publishing API 逐项 CLI 参考

本文从浏览器保存的官方正文生成，每项说明用途、传输位置、参数和前置条件。生成命令：`python3 scripts/render-publishing-guide.py`。

流程和请求文件示例见 [Publishing 使用与核验](PUBLISHING_WORKFLOWS.md)。表格保留官方必选标记；条件要求和表格/示例冲突需要同时阅读说明。本文件不代表所有模型规则已实现或全部线上操作已验收。

## 参数位置

`query` 使用 `--query key=value`，`header` 使用 `--header key=value`，`body` 使用 JSON 文件 `--body request.json` 或已声明类型的 `--field key=value`。JSON 整数不加引号，对象和数组按实际模型展开。

实际调用选择 `--profile PROFILE`，先使用 `--invoke --dry-run` 查看请求计划；只有准备实际调用时才使用 `--invoke --dry-run=false`。dry-run 不证明资源状态、权限或审核结果。

## 创建APP ID

命令：`agc publishing app-create`。请求：`POST /api/publish/v3/app`。

[官方正文](https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-provision-api-create-appid-0000002526543939)；官网更新：2026-09-04 16:04；浏览器读取：2026-10-07T05:11:11.924Z。

此接口用于为HarmonyOS应用/元服务创建APP ID。

### header（Service Account方式：）

[参数来源](https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-provision-api-create-appid-0000002526543939#section1323714041015)

| 参数名称 | 必选(M)/可选(O) | 类型 | 参数说明 |
| --- | --- | --- | --- |
| Authorization | M | String | 认证信息，格式为“Authorization: Bearer ${JWT}”。JWT为通过Service Account方式获取授权中获取的鉴权令牌。 |

### header（API客户端方式：）

[参数来源](https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-provision-api-create-appid-0000002526543939#section1323714041015)

| 参数名称 | 必选(M)/可选(O) | 类型 | 参数说明 |
| --- | --- | --- | --- |
| client_id | M | String | 客户端ID，获取方法参考创建API客户端。 |
| Authorization | M | String | 认证信息，格式为“Authorization: Bearer ${access_token}”。access_token为获取Token中获取的access_token。 |

### header（OAuth客户端方式：）

[参数来源](https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-provision-api-create-appid-0000002526543939#section1323714041015)

| 参数名称 | 必选(M)/可选(O) | 类型 | 参数说明 |
| --- | --- | --- | --- |
| teamId | M | String(64) | 开发者所在团队的团队ID。 |
| oauth2Token | M | String | 认证信息，传入获取用户授权码中获取的Access Token。 |

### query

[参数来源](https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-provision-api-create-appid-0000002526543939#section16551174881019)

| 参数名称 | 必选(M)/可选(O) | 类型 | 参数说明 |
| --- | --- | --- | --- |
| projectId | O | String | 项目ID，获取方法请参见查看应用信息。 如果不传projectId，系统会先在内部创建一个与应用同名的项目，随后再创建应用。 |

### body

[参数来源](https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-provision-api-create-appid-0000002526543939#section16551174881019)

| 参数名称 | 必选(M)/可选(O) | 类型 | 参数说明 |
| --- | --- | --- | --- |
| appName | M | String(64) | 应用名称。 |
| packageName | O | String(128) | 应用包名，必须唯一，且不能与其他应用（包括Android应用）包名相同。 说明 installationFree取值为0（HarmonyOS应用）时，此参数必填。 installationFree取值为1（元服务）时，系统将自动生成包名，此参数无需填写。 |
| parentType | M | Integer(32) | 应用分类。 取值范围： 2：游戏 13：应用 |
| installationFree | M | Integer(32) | 应用类型。 取值范围： 0：HarmonyOS应用 1：元服务 默认值：0 |

### 请求说明

本接口支持使用Service Account方式、API客户端方式和OAuth客户端方式，区别请参见获取服务端授权。

### 完整模型及示例

```sh
agc publishing app-create --docs --pretty
```

该命令返回官方请求/响应示例及已捕获的关联模型。业务响应中的 `ret.code` 应与 HTTP 状态分别检查；修改后使用相应查询接口确认目标状态。

## 查询应用包名对应的appid

命令：`agc publishing appid-list`。请求：`GET /api/publish/v2/appid-list`。

[官方正文](https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-appid-list-0000002271000617)；官网更新：2026-06-25 19:52；浏览器读取：2026-10-07T05:11:13.650Z。

此接口用于根据HarmonyOS应用/元服务的应用包名查询对应的应用ID，支持一次查询多个应用包名。

接口调用者的角色：账号持有者、管理员、APP管理员、运营。

### header（Service Account方式：）

[参数来源](https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-appid-list-0000002271000617#section119494174419)

| 参数名称 | 必选(M)/可选(O) | 类型 | 参数说明 |
| --- | --- | --- | --- |
| Authorization | M | String | 认证信息，格式为“Authorization: Bearer ${JWT}”。JWT为通过Service Account方式获取授权中获取的鉴权令牌。 |

### header（API客户端方式：）

[参数来源](https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-appid-list-0000002271000617#section119494174419)

| 参数名称 | 必选(M)/可选(O) | 类型 | 参数说明 |
| --- | --- | --- | --- |
| client_id | M | String | 客户端ID，获取方法参考创建API客户端。 |
| Authorization | M | String | 认证信息，格式为“Authorization: Bearer ${access_token}”。access_token为获取Token中获取的access_token。 |

### query

[参数来源](https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-appid-list-0000002271000617#section19247123515514)

| 参数名称 | 必选(M)/可选(O) | 类型 | 参数说明 |
| --- | --- | --- | --- |
| packageName | M | String(4096) | 需要查询的应用包名。 多个包名以逗号分隔，最多支持50个。 |
| packageTypes | O | String(32) | 通过软件包类型过滤应用，只查询指定软件包类型的应用。 取值范围： 7：APP（HarmonyOS应用/元服务） |

### 请求说明

本接口支持使用Service Account方式和API客户端方式，二者区别请参见获取服务端授权。

### 完整模型及示例

```sh
agc publishing appid-list --docs --pretty
```

该命令返回官方请求/响应示例及已捕获的关联模型。业务响应中的 `ret.code` 应与 HTTP 状态分别检查；修改后使用相应查询接口确认目标状态。

## 查询团队账号下的成员角色列表

命令：`agc publishing members-role-list`。请求：`GET /api/publish/v3/members-role-list`。

[官方正文](https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-teamid-member-list-0000002679253743)；官网更新：2026-07-31 17:44；浏览器读取：2026-10-07T05:11:33.888Z。

此接口用于查询团队账号下的成员角色列表和团队主账号的开发者信息。

接口调用者的角色：账号持有者、管理员、APP管理员、运营、开发、客服、财务。

### header（Service Account方式：）

[参数来源](https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-teamid-member-list-0000002679253743#section119494174419)

| 参数名称 | 必选(M)/可选(O) | 类型 | 参数说明 |
| --- | --- | --- | --- |
| Authorization | M | String | 认证信息，格式为“Authorization: Bearer ${JWT}”。JWT为通过Service Account方式获取授权中获取的鉴权令牌。 |

### header（API客户端方式：）

[参数来源](https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-teamid-member-list-0000002679253743#section119494174419)

| 参数名称 | 必选(M)/可选(O) | 类型 | 参数说明 |
| --- | --- | --- | --- |
| client_id | M | String | 客户端ID，获取方法参考创建API客户端。 |
| Authorization | M | String | 认证信息，格式为“Authorization: Bearer ${access_token}”。access_token为获取Token中获取的access_token。 |

### query

[参数来源](https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-teamid-member-list-0000002679253743#section420734618165)

| 参数名称 | 必选(M)/可选(O) | 类型 | 参数说明 |
| --- | --- | --- | --- |
| fromPage | M | Integer | 起始页。 最小值：1，表示从第1页开始查询。 |
| pageSize | M | Integer | 每页查询的条数。 取值范围：[1,50] |
| languageCode | M | String(20) | 语言码。 取值范围： zh_cn：中文 en_us：英文 ru_ru：俄语 |

### 请求说明

本接口支持使用Service Account方式和API客户端方式，二者区别请参见获取服务端授权。

### 完整模型及示例

```sh
agc publishing members-role-list --docs --pretty
```

该命令返回官方请求/响应示例及已捕获的关联模型。业务响应中的 `ret.code` 应与 HTTP 状态分别检查；修改后使用相应查询接口确认目标状态。

## 查询应用信息

命令：`agc publishing app-info-query-v3`。请求：`GET /api/publish/v3/app-info`。

[官方正文](https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-appinfo-query-0000002236041422)；官网更新：2026-07-31 17:44；浏览器读取：2026-10-07T05:11:36.198Z。

此接口用于根据指定语言和应用ID查询HarmonyOS应用/元服务的详细信息。

接口调用者的角色：账号持有者、管理员、APP管理员、运营。

### header（Service Account方式：）

[参数来源](https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-appinfo-query-0000002236041422#section119494174419)

| 参数名称 | 必选(M)/可选(O) | 类型 | 参数说明 |
| --- | --- | --- | --- |
| Authorization | M | String | 认证信息，格式为“Authorization: Bearer ${JWT}”。JWT为通过Service Account方式获取授权中获取的鉴权令牌。 |

### header（API客户端方式：）

[参数来源](https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-appinfo-query-0000002236041422#section119494174419)

| 参数名称 | 必选(M)/可选(O) | 类型 | 参数说明 |
| --- | --- | --- | --- |
| client_id | M | String | 客户端ID，获取方法参考创建API客户端。 |
| Authorization | M | String | 认证信息，格式为“Authorization: Bearer ${access_token}”。access_token为获取Token中获取的access_token。 |

### query

[参数来源](https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-appinfo-query-0000002236041422#section420734618165)

| 参数名称 | 必选(M)/可选(O) | 类型 | 参数说明 |
| --- | --- | --- | --- |
| appId | M | String(32) | 需要查询的应用ID。 |
| lang | O | String(16) | 需要查询的语言。 说明 取值参考语言类型。 不传该参数，则查询全部语言的应用信息。 |
| releaseType | O | Integer | 应用发布方式。 取值范围： 1：全网 6：HarmonyOS测试发布方式（API>=10） 默认值：1 |
| versionId | O | String | 需要查询的版本ID。 注意 releaseType为6，该字段必填。 |

### 请求说明

本接口支持使用Service Account方式和API客户端方式，二者区别请参见获取服务端授权。

### 完整模型及示例

```sh
agc publishing app-info-query-v3 --docs --pretty
```

该命令返回官方请求/响应示例及已捕获的关联模型。业务响应中的 `ret.code` 应与 HTTP 状态分别检查；修改后使用相应查询接口确认目标状态。

## 更新应用基本信息

命令：`agc publishing app-info-update-v3`。请求：`PUT /api/publish/v3/app-info`。

[官方正文](https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-appinfo-update-0000002236201246)；官网更新：2026-06-25 19:52；浏览器读取：2026-10-07T05:11:38.645Z。

此接口用于更新指定HarmonyOS应用/元服务的应用详情。

接口调用者的角色：账号持有者、管理员、APP管理员、运营。

### header（Service Account方式：）

[参数来源](https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-appinfo-update-0000002236201246#section119494174419)

| 参数名称 | 必选(M)/可选(O) | 类型 | 参数说明 |
| --- | --- | --- | --- |
| Authorization | M | String | 认证信息，格式为“Authorization: Bearer ${JWT}”。JWT为通过Service Account方式获取授权中获取的鉴权令牌。 |

### header（API客户端方式：）

[参数来源](https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-appinfo-update-0000002236201246#section119494174419)

| 参数名称 | 必选(M)/可选(O) | 类型 | 参数说明 |
| --- | --- | --- | --- |
| client_id | M | String | 客户端ID，获取方法参考创建API客户端。 |
| Authorization | M | String | 认证信息，格式为“Authorization: Bearer ${access_token}”。access_token为获取Token中获取的access_token。 |

### query

[参数来源](https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-appinfo-update-0000002236201246#section886219295206)

| 参数名称 | 必选(M)/可选(O) | 类型 | 参数说明 |
| --- | --- | --- | --- |
| appId | M | String(32) | 需要更新的应用ID。 |
| releaseType | O | Integer(32) | 应用发布方式。 取值范围： 1：全网 默认值：1 |
| releasePhase | O | Integer(32) | 分阶段发布标识。 取值范围： 0：全网 3：分阶段 默认值：0 |

### body

[参数来源](https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-appinfo-update-0000002236201246#section17512144171520)

| 参数名称 | 必选(M)/可选(O) | 类型 | 参数说明 |
| --- | --- | --- | --- |
| defaultLang | O | String(16) | 应用的默认语言。 说明 取值参见语言类型。 |
| harmonyChildType | O | Integer(4) | HarmonyOS应用/元服务的二级分类。 固定ID。 说明 取值参见应用游戏分类。 |
| kindMainTag | O | Integer(4) | HarmonyOS应用/元服务主标签。 固定ID。 说明 取值参见应用游戏分类的“标签”列。 |
| kindSubTags | O | List<Integer> | HarmonyOS应用/元服务附属标签。 固定ID。最多可选4个附属标签。 说明 取值参见应用游戏分类的“标签”列。 |
| privacyPolicy | O | String(255) | 隐私政策网址。 |
| appNetType | O | Integer | 应用联网类型。 取值范围： 1：单机 2：网游 默认值：1 |
| publishCountry | M | String(2048) | 应用发布国家的国家码。 多个国家以英文逗号分隔，例如：“CN,DE,MC”。 如果取值为“OTHER”，则表示当App Gallery支持的国家或地区列表中新增某个国家或地区时，应用将自动分发到新增的国家和地区。 注意 取值不能只包含OTHER。 如果国家码有设置非CN的地区，将存在个人信息出境的情况，具体请参考个人信息出境授权声明。 |
| appTariffType | O | String(16) | 应用内付费道具类型。 取值范围： 1：激活收费 2：道具收费 3：关卡收费 4：购买虚拟币 5：部分章节收费（图书阅读类） 6：其他的 7：课程收费 8：会员收费 支持设置多个，使用英文逗号分隔，例如“1,2,3”。 |
| publicationNumber | O | String(50) | 版号信息。 此参数只支持游戏类应用。 |
| developerWebsite | O | String(255) | 开发者官网。 |
| customInfo | O | CustInfo | 客服信息。 |
| deviceTypes | O | List<DeviceTypeInfo> | 设备类型信息。 |
| privacyRightsUrl | O | String(256) | 隐私权利网址。 即关于您的用户实施其权利的网站地址，如用户删除、修改、导出其个人数据的入口。 |
| registeredDclType | O | Integer(32) | APP类型信息。 取值范围： -1：清空APP类型信息 0：APP服务器在中国大陆 1：APP服务器不在中国大陆 2：APP为单机APP |
| registeredIdType | O | Integer(32) | 主办单位类型。 取值范围： -1：清空数据（会联动清空registeredIdNumber的内容） 1：企业（社会信用代码） 2：个人（身份证号码） 3：机构（机构代码） |
| registeredIdNumber | O | String(64) | 主办单位证件号。 取值范围： registeredIdType为1时，对应统一社会信用代码 registeredIdType为2时，对应身份证号码 registeredIdType为3时，对应组织机构代码 |
| registeredEntity | O | Integer(32) | 主办单位。 取值范围： 1：与开发者账号主体一致 2：与开发者账号主体不一致 说明 当registeredDclType值为0（APP服务器在中国大陆）时，该字段必填。若不填写或填写错误，会导致审核不通过。 若在版本提交时，该字段值选择“一致”，则主办单位名称和证件号，将直接复用开发者联盟注册账号信息。 若您提交的主办单位信息与联盟账号信息不一致，该字段值却选择了“一致”，则系统将自动判定为“不一致”，将会影响审核结果。 |
| registeredEntityName | O | String(1024) | 主办单位名称。 说明 当registeredDclType值为0（APP服务器在中国大陆）时，该字段必填。若不填写或填写错误，会导致审核不通过。 |
| privacyAgreementId | O | String(32) | 隐私协议ID。 |
| userAgreementUrl | O | String(512) | 用户协议URL。 |
| userAgreementIds | O | List<String(32)> | 用户协议ID列表。 当userAgreementUrl和userAgreementIds同时传值时，以userAgreementIds值为准。 数组长度不超过5。 |
| encrypted | M | Integer(32) | 包是否加密。 取值范围： 0：否，表示不加密。 1：是，表示加密。 默认值：1 |
| isAiGenerate | O | Integer(32) | 是否涉及AI生成合成服务。 取值范围： 0：不涉及。 1：涉及。 |
| aiTypes | O | List<Integer(32)> | 涉及的AI生成合成服务类型。 取值范围： -1：其他 1：文本 2：图片 3：音频 4：视频 5：虚拟场景 说明 isAiGenerate取值为1（涉及）时，此参数必填。 |
| testUserName | O | String(128) | 应用审核信息的测试账号用户名。 说明 如果存在需要登录后才可使用的功能，需要提供测试账号用户名。 |
| testUserPassword | O | String(128) | 应用审核信息的测试账号密码。 说明 如果存在需要登录后才可使用的功能，需要提供测试账号密码。 |
| appRemark | O | String(300) | 应用审核信息的备注。 提供有助于审核人员更准确、高效测试您应用的额外信息，比如在审核时需要的特别设置等。 |
| selfTestVideos | O | String(2000) | 应用审核信息的自测文件。 格式为JSON字符串，示例：{"objectId1":"文件名1","objectId2":"文件名2"...} 说明 objectId的获取流程：先调用获取上传文件地址接口，然后调用上传文件接口上传文件，最后得到获取上传文件地址接口返回的objectId。 最多可上传5个自测视频，总大小不超过500MB。 注意 此字段只支持元服务。 元服务需要与设备联动时，此参数必填。 |
| appReviewPhoneInfo | O | AppReviewContact | 应用负责人的手机号信息。用于沟通审核问题。 |
| appReviewEmailInfo | O | AppReviewContact | 应用负责人的邮箱账号信息。用于接收上架审核结果，应用整改或下架通知。 |
| appReviewName | O | String(256) | 应用负责人的姓名。 |
| versionVocs | O | List<VersionVoc> | 应用版本已解决问题。 数组长度不超过1000。 说明 调用查询应用未处理的问题工单接口可以查询用户反馈的应用问题。如果在应用版本更新时某些问题已经解决，可以将这些问题提交审核，工单平台验证通过后会关闭工单。 |

### 请求说明

本接口支持使用Service Account方式和API客户端方式，二者区别请参见获取服务端授权。

请求Body中使用JSON格式携带更新的应用信息，参数如下表所示。

### 完整模型及示例

```sh
agc publishing app-info-update-v3 --docs --pretty
```

该命令返回官方请求/响应示例及已捕获的关联模型。业务响应中的 `ret.code` 应与 HTTP 状态分别检查；修改后使用相应查询接口确认目标状态。

## 更新应用语言信息

命令：`agc publishing language-info-update-v3`。请求：`PUT /api/publish/v3/app-language-info`。

[官方正文](https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-language-info-update-0000002271160569)；官网更新：2026-06-25 19:52；浏览器读取：2026-10-07T04:55:27.457Z。

此接口用于更新指定的HarmonyOS应用/元服务的语言描述信息，如果当前没有语言描述信息则新增。

接口调用者的角色：账号持有者、管理员、APP管理员、运营。

### header（Service Account方式：）

[参数来源](https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-language-info-update-0000002271160569#section119494174419)

| 参数名称 | 必选(M)/可选(O) | 类型 | 参数说明 |
| --- | --- | --- | --- |
| Authorization | M | String | 认证信息，格式为“Authorization: Bearer ${JWT}”。JWT为通过Service Account方式获取授权中获取的鉴权令牌。 |

### header（API客户端方式：）

[参数来源](https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-language-info-update-0000002271160569#section119494174419)

| 参数名称 | 必选(M)/可选(O) | 类型 | 参数说明 |
| --- | --- | --- | --- |
| client_id | M | String | 客户端ID，获取方法参考创建API客户端。 |
| Authorization | M | String | 认证信息，格式为“Authorization: Bearer ${access_token}”。access_token为获取Token中获取的access_token。 |

### query

[参数来源](https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-language-info-update-0000002271160569#section41341319578)

| 参数名称 | 必选(M)/可选(O) | 类型 | 参数说明 |
| --- | --- | --- | --- |
| appId | M | String | 需要更新语言描述信息的应用ID。 |
| releaseType | O | Integer(32) | 应用发布方式。 取值范围： 1：全网 默认值：1 |
| releasePhase | O | Integer(32) | 分阶段发布标识。 取值范围： 0：全网 3：分阶段 默认值：0 |

### body

[参数来源](https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-language-info-update-0000002271160569#section13576204051611)

| 参数名称 | 必选(M)/可选(O) | 类型 | 参数说明 |
| --- | --- | --- | --- |
| lang | M | String(8) | 语言种类。 说明 取值参见语言类型。 |
| appName | O | String(64) | 对应语言的应用名称。 新增一个语言种类时，此字段必选。 |
| appDesc | O | String(8000) | 对应语言的应用描述信息。 |
| briefInfo | O | String(80) | 对应语言的一句话简介。 |
| newFeatures | O | String(500) | 对应语言的新版本简介。 说明 升级版本时，此参数必填。 |

### 请求说明

本接口支持使用Service Account方式和API客户端方式，二者区别请参见获取服务端授权。

请求Body中使用JSON格式携带更新的语言描述信息，参数如下表所示。

appDesc、briefInfo和newFeatures三个参数的取值不能相同。

### 完整模型及示例

```sh
agc publishing language-info-update-v3 --docs --pretty
```

该命令返回官方请求/响应示例及已捕获的关联模型。业务响应中的 `ret.code` 应与 HTTP 状态分别检查；修改后使用相应查询接口确认目标状态。

## 删除应用语言信息

命令：`agc publishing language-info-delete`。请求：`DELETE /api/publish/v2/app-language-info`。

[官方正文](https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-language-info-delete-0000002271000621)；官网更新：2026-06-25 19:52；浏览器读取：2026-10-07T04:55:29.656Z。

此接口用于删除指定HarmonyOS应用/元服务不需要的语言描述信息，默认语言不支持删除。

接口调用者的角色：账号持有者、管理员、APP管理员、运营。

### header（Service Account方式：）

[参数来源](https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-language-info-delete-0000002271000621#section119494174419)

| 参数名称 | 必选(M)/可选(O) | 类型 | 参数说明 |
| --- | --- | --- | --- |
| Authorization | M | String | 认证信息，格式为“Authorization: Bearer ${JWT}”。JWT为通过Service Account方式获取授权中获取的鉴权令牌。 |

### header（API客户端方式：）

[参数来源](https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-language-info-delete-0000002271000621#section119494174419)

| 参数名称 | 必选(M)/可选(O) | 类型 | 参数说明 |
| --- | --- | --- | --- |
| client_id | M | String | 客户端ID，获取方法参考创建API客户端。 |
| Authorization | M | String | 认证信息，格式为“Authorization: Bearer ${access_token}”。access_token为获取Token中获取的access_token。 |

### query

[参数来源](https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-language-info-delete-0000002271000621#section23242264517)

| 参数名称 | 必选(M)/可选(O) | 类型 | 参数说明 |
| --- | --- | --- | --- |
| appId | M | String | 需要删除语言描述信息的应用ID。 |
| lang | M | String | 需要删除的语言。 说明 取值参见语言类型。 |
| releaseType | O | Integer | 应用发布方式。 取值范围： 1：全网 默认值：1 |

### 请求说明

本接口支持使用Service Account方式和API客户端方式，二者区别请参见获取服务端授权。

### 完整模型及示例

```sh
agc publishing language-info-delete --docs --pretty
```

该命令返回官方请求/响应示例及已捕获的关联模型。业务响应中的 `ret.code` 应与 HTTP 状态分别检查；修改后使用相应查询接口确认目标状态。

## 更新应用文件信息

命令：`agc publishing app-file-info-v3`。请求：`PUT /api/publish/v3/app-file-info`。

[官方正文](https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-app-file-info-update-0000002236041430)；官网更新：2026-06-25 19:52；浏览器读取：2026-10-07T04:55:31.723Z。

图片、视频等文件上传完成后，通过该接口刷新HarmonyOS应用/元服务的文件信息。

接口调用者的角色：账号持有者、管理员、APP管理员、运营。

### header（Service Account方式：）

[参数来源](https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-app-file-info-update-0000002236041430#section119494174419)

| 参数名称 | 必选(M)/可选(O) | 类型 | 参数说明 |
| --- | --- | --- | --- |
| Authorization | M | String | 认证信息，格式为“Authorization: Bearer ${JWT}”。JWT为通过Service Account方式获取授权中获取的鉴权令牌。 |

### header（API客户端方式：）

[参数来源](https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-app-file-info-update-0000002236041430#section119494174419)

| 参数名称 | 必选(M)/可选(O) | 类型 | 参数说明 |
| --- | --- | --- | --- |
| client_id | M | String | 客户端ID，获取方法参考创建API客户端。 |
| Authorization | M | String | 认证信息，格式为“Authorization: Bearer ${access_token}”。access_token为获取Token中获取的access_token。 |

### query

[参数来源](https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-app-file-info-update-0000002236041430#section156021546155912)

| 参数名称 | 必选(M)/可选(O) | 类型 | 参数说明 |
| --- | --- | --- | --- |
| appId | M | String(32) | 需要更新的应用ID。 |
| releaseType | O | Integer(32) | 应用发布方式。 取值范围： 1：全网 默认值：1 |
| releasePhase | O | Integer(32) | 分阶段发布标识。 取值范围： 0：全网 3：分阶段 默认值：0 |

### body

[参数来源](https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-app-file-info-update-0000002236041430#section1131102811712)

| 参数名称 | 必选(M)/可选(O) | 类型 | 参数说明 |
| --- | --- | --- | --- |
| appIconList | O | List<LangFileInfo> | 应用图标信息。 数组长度不超过20。 说明 应用图标必须与软件包中的应用图标一致，且需符合审核政策与应用文件要求。 |
| screenShotList | O | List<LangFileInfo> | 应用介绍截图信息。 数组长度不超过20。 说明 在同一语言、同一设备下，同时更新的截图素材和视频素材方向必须保持一致，否则将报错。 |
| introVideoList | O | List<LangFileInfo> | 应用介绍视频信息。 数组长度不超过20。 说明 在同一语言、同一设备下，同时更新的截图素材和视频素材方向必须保持一致，否则将报错。 |
| rcmdVideoList | O | List<LangFileInfo> | 推荐视频信息。 数组长度不超过20。 |
| rcmdPicList | O | List<LangFileInfo> | 推荐图片信息。 数组长度不超过20。 |
| elecCertificate | O | String(512) | 电子版本证书。 |
| certificateList | O | List<String(512)> | 应用版权证书或代理证书。 数组长度不超过10。 |
| publicationList | O | List<String(512)> | 版号证明。 数组长度不超过10。 |
| packagePermissionIntroVideoList | O | List<PackagePermissionIntroVideo> | 包权限使用场景说明视频。 数组长度不超过100。 说明 隐私权限说明的“权限描述”、“权限使用设备范围”和“权限使用理由”从上传的软件包体中读取，无需通过接口更新。 |
| otherCertificatesURL | O | String(512) | 授权书及其他材料文件压缩包的objectId。 说明 objectId的获取流程：先调用获取上传文件地址接口，然后调用上传文件接口上传文件，最后得到获取上传文件地址接口返回的objectId。 若文件格式不正确，提交发布时会出错，具体限制如下： 压缩包必须为zip格式，不超过200MB。 解压后的大小不超过1GB。 文件夹层级不超过3层。 只能包含JPG、JPEG、BMP格式图片和PDF格式的文件，累计不超过200个。 说明 此字段只支持HarmonyOS应用，不支持元服务。 |
| updateScreenShotList | O | List<LangFileInfo> | 更新截图素材。 数组长度不超过20。 |
| updateVideoList | O | List<LangFileInfo> | 更新视频素材。 数组长度不超过20。 |

### 请求说明

本接口支持使用Service Account方式和API客户端方式，二者区别请参见获取服务端授权。

请求Body中使用JSON格式携带更新的应用文件信息，参数如下表所示。

各文件类型对应的文件上传要求请参见应用文件要求。

### 完整模型及示例

```sh
agc publishing app-file-info-v3 --docs --pretty
```

该命令返回官方请求/响应示例及已捕获的关联模型。业务响应中的 `ret.code` 应与 HTTP 状态分别检查；修改后使用相应查询接口确认目标状态。

## 更新应用软件包信息

命令：`agc publishing app-package-info-update`。请求：`PUT /api/publish/v3/app-package-info`。

[官方正文](https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-app-package-info-update-0000002236201250)；官网更新：2026-06-25 19:52；浏览器读取：2026-10-07T04:55:46.452Z。

APP软件包上传完成后，通过该接口刷新HarmonyOS应用/元服务的软件包信息。

该接口支持解析软件包中的构建版本号（buildVersion），用于区分同一主版本下的不同测试子版本。解析完成后，可以从查询应用软件包信息接口返回的packageInfo参数中获取buildVersion字段。

接口调用者的角色：账号持有者、管理员、APP管理员、运营。

### header（Service Account方式：）

[参数来源](https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-app-package-info-update-0000002236201250#section119494174419)

| 参数名称 | 必选(M)/可选(O) | 类型 | 参数说明 |
| --- | --- | --- | --- |
| Authorization | M | String | 认证信息，格式为“Authorization: Bearer ${JWT}”。JWT为通过Service Account方式获取授权中获取的鉴权令牌。 |

### header（API客户端方式：）

[参数来源](https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-app-package-info-update-0000002236201250#section119494174419)

| 参数名称 | 必选(M)/可选(O) | 类型 | 参数说明 |
| --- | --- | --- | --- |
| client_id | M | String | 客户端ID，获取方法参考创建API客户端。 |
| Authorization | M | String | 认证信息，格式为“Authorization: Bearer ${access_token}”。access_token为获取Token中获取的access_token。 |

### query

[参数来源](https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-app-package-info-update-0000002236201250#section156021546155912)

| 参数名称 | 必选(M)/可选(O) | 类型 | 参数说明 |
| --- | --- | --- | --- |
| appId | M | String | 需要更新的应用ID。 |
| releaseType | O | Integer(32) | 应用发布方式。 取值范围： 1：全网 默认值：1 |
| releasePhase | O | Integer(32) | 分阶段发布标识。 取值范围： 0：全网 3：分阶段 默认值：0 |

### body

[参数来源](https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-app-package-info-update-0000002236201250#section1131102811712)

| 参数名称 | 必选(M)/可选(O) | 类型 | 参数说明 |
| --- | --- | --- | --- |
| fileName | M | String(255) | 文件名称，包括文件的后缀名。 |
| objectId | M | String(512) | 文件在文件服务器中的对象ID。 说明 objectId的获取流程：先调用获取上传文件地址接口，然后调用上传文件接口上传文件，最后得到获取上传文件地址接口返回的objectId。 |

### 请求说明

本接口支持使用Service Account方式和API客户端方式，二者区别请参见获取服务端授权。

请求Body中使用JSON格式携带更新的应用软件包信息，参数如下表所示。

### 完整模型及示例

```sh
agc publishing app-package-info-update --docs --pretty
```

该命令返回官方请求/响应示例及已捕获的关联模型。业务响应中的 `ret.code` 应与 HTTP 状态分别检查；修改后使用相应查询接口确认目标状态。

## 更新版本上架时间

命令：`agc publishing update-releasetime`。请求：`PUT /api/publish/v2/on-shelf-time`。

[官方正文](https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-update-release-time-0000002271160573)；官网更新：2026-06-25 19:52；浏览器读取：2026-10-07T04:55:48.314Z。

此接口用于设置版本的上架时间。

接口调用者的角色：账号持有者、管理员、APP管理员、运营。

### header（Service Account方式：）

[参数来源](https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-update-release-time-0000002271160573#section119494174419)

| 参数名称 | 必选(M)/可选(O) | 类型 | 参数说明 |
| --- | --- | --- | --- |
| Authorization | M | String | 认证信息，格式为“Authorization: Bearer ${JWT}”。JWT为通过Service Account方式获取授权中获取的鉴权令牌。 |

### header（API客户端方式：）

[参数来源](https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-update-release-time-0000002271160573#section119494174419)

| 参数名称 | 必选(M)/可选(O) | 类型 | 参数说明 |
| --- | --- | --- | --- |
| client_id | M | String | 客户端ID，获取方法参考创建API客户端。 |
| Authorization | M | String | 认证信息，格式为“Authorization: Bearer ${access_token}”。access_token为获取Token中获取的access_token。 |

### query

[参数来源](https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-update-release-time-0000002271160573#section165302018915)

| 参数名称 | 必选(M)/可选(O) | 类型 | 参数说明 |
| --- | --- | --- | --- |
| appId | M | String | 应用ID。 |

### body

[参数来源](https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-update-release-time-0000002271160573#section2357705386)

| 参数名称 | 必选(M)/可选(O) | 类型 | 参数说明 |
| --- | --- | --- | --- |
| changeType | M | Integer | 修改类型。 取值范围： 2：指定时间上架改为审核通过立即上架 3：变更定时上架时间 |
| releaseTime | M | String | 指定的版本上架时间。 UTC时间格式：yyyy-MM-ddTHH:mm:ssZZ，例如“2015-01-01T01:01:01+0800”。 changeType为3时有效。 |
| releaseType | M | Integer | 发布方式。 取值范围： 1：全网 目前支持全网及releasePhase为3的分阶段发布。 |

### 请求说明

本接口支持使用Service Account方式和API客户端方式，二者区别请参见获取服务端授权。

请求Body中使用JSON格式携带更新版本上架时间内容，参数如下表所示。

### 完整模型及示例

```sh
agc publishing update-releasetime --docs --pretty
```

该命令返回官方请求/响应示例及已捕获的关联模型。业务响应中的 `ret.code` 应与 HTTP 状态分别检查；修改后使用相应查询接口确认目标状态。

## 更新分阶段发布

命令：`agc publishing phased-release-update`。请求：`PUT /api/publish/v2/version/phased-release`。

[官方正文](https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-phased-release-0000002271000625)；官网更新：2026-06-25 19:52；浏览器读取：2026-10-07T04:55:50.500Z。

此接口用于将分阶段发布修改为全网发布，或者更新分阶段发布的设置。

接口调用者的角色：账号持有者、管理员、APP管理员、运营。

### header（Service Account方式：）

[参数来源](https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-phased-release-0000002271000625#section236517415165)

| 参数名称 | 必选(M)/可选(O) | 类型 | 参数说明 |
| --- | --- | --- | --- |
| Authorization | M | String | 认证信息，格式为“Authorization: Bearer ${JWT}”。JWT为通过Service Account方式获取授权中获取的鉴权令牌。 |
| appId | M | String(32) | 应用ID，获取方法参考查看应用信息。 |

### header（API客户端方式：）

[参数来源](https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-phased-release-0000002271000625#section236517415165)

| 参数名称 | 必选(M)/可选(O) | 类型 | 参数说明 |
| --- | --- | --- | --- |
| client_id | M | String | 客户端ID，获取方法参考创建API客户端。 |
| Authorization | M | String | 认证信息，格式为“Authorization: Bearer ${access_token}”。access_token为获取Token中获取的access_token。 |
| appId | M | String(32) | 应用ID，获取方法参考查看应用信息。 |

### body

[参数来源](https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-phased-release-0000002271000625#section1520919574201)

| 参数名称 | 必选(M)/可选(O) | 类型 | 参数说明 |
| --- | --- | --- | --- |
| versionId | M | String(32) | 版本ID。 可从查询应用所有版本类型的版本列表接口返回的versionId参数中获取。 |
| releasePhase | O | Integer(32) | 分阶段发布标识。 取值范围： 3：分阶段。 4：分阶段转全网。 |
| state | O | String(16) | 分阶段发布的状态。 取值范围： SUSPEND：暂停。 RELEASE：恢复为发布中。 |
| description | O | String(500) | 分阶段发布说明。 |
| phaseDay | O | Integer(32) | 更新分阶段天数。 取值范围2、3、4、5、6、7，对应分阶段发布比例加速到2%、5%、10%、20%、50%、100%。 |

### 请求说明

本接口支持使用Service Account方式和API客户端方式，二者区别请参见获取服务端授权。

请求Body中使用JSON格式携带分阶段发布信息，参数如下表所示。

### 完整模型及示例

```sh
agc publishing phased-release-update --docs --pretty
```

该命令返回官方请求/响应示例及已捕获的关联模型。业务响应中的 `ret.code` 应与 HTTP 状态分别检查；修改后使用相应查询接口确认目标状态。

## 查询应用软件包信息

命令：`agc publishing app-package-info-query`。请求：`GET /api/publish/v2/app-package-info`。

[官方正文](https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-query-package-info-0000002328848469)；官网更新：2026-06-25 19:52；浏览器读取：2026-10-07T04:56:01.749Z。

此接口用于查询HarmonyOS应用/元服务的软件包信息。

接口调用者的角色：账号持有者、管理员、APP管理员、运营。

### header（Service Account方式：）

[参数来源](https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-query-package-info-0000002328848469#section236517415165)

| 参数名称 | 必选(M)/可选(O) | 类型 | 参数说明 |
| --- | --- | --- | --- |
| Authorization | M | String | 认证信息，格式为“Authorization: Bearer ${JWT}”。JWT为通过Service Account方式获取授权中获取的鉴权令牌。 |
| appId | M | String(32) | 应用ID，获取方法参考查看应用信息。 |

### header（API客户端方式：）

[参数来源](https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-query-package-info-0000002328848469#section236517415165)

| 参数名称 | 必选(M)/可选(O) | 类型 | 参数说明 |
| --- | --- | --- | --- |
| client_id | M | String | 客户端ID，获取方法参考创建API客户端。 |
| Authorization | M | String | 认证信息，格式为“Authorization: Bearer ${access_token}”。access_token为获取Token中获取的access_token。 |
| appId | M | String(32) | 应用ID，获取方法参考查看应用信息。 |

### query

[参数来源](https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-query-package-info-0000002328848469#section165302018915)

| 参数名称 | 必选(M)/可选(O) | 类型 | 参数说明 |
| --- | --- | --- | --- |
| packageId | M | String(32) | 待查询的软件包ID。 |

### 请求说明

本接口支持使用Service Account方式和API客户端方式，二者区别请参见获取服务端授权。

### 完整模型及示例

```sh
agc publishing app-package-info-query --docs --pretty
```

该命令返回官方请求/响应示例及已捕获的关联模型。业务响应中的 `ret.code` 应与 HTTP 状态分别检查；修改后使用相应查询接口确认目标状态。

## 查询软件包编译状态

命令：`agc publishing package-compile-status-v3`。请求：`GET /api/publish/v3/package/compile/status`。

[官方正文](https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-query-compile-status-0000002236041434)；官网更新：2026-06-25 19:52；浏览器读取：2026-10-07T04:56:03.411Z。

此接口用于查询HarmonyOS应用/元服务的软件包的编译状态。

接口调用者的角色：账号持有者、管理员、APP管理员、运营。

### header（Service Account方式：）

[参数来源](https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-query-compile-status-0000002236041434#section119494174419)

| 参数名称 | 必选(M)/可选(O) | 类型 | 参数说明 |
| --- | --- | --- | --- |
| Authorization | M | String | 认证信息，格式为“Authorization: Bearer ${JWT}”。JWT为通过Service Account方式获取授权中获取的鉴权令牌。 |

### header（API客户端方式：）

[参数来源](https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-query-compile-status-0000002236041434#section119494174419)

| 参数名称 | 必选(M)/可选(O) | 类型 | 参数说明 |
| --- | --- | --- | --- |
| client_id | M | String | 客户端ID，获取方法参考创建API客户端。 |
| Authorization | M | String | 认证信息，格式为“Authorization: Bearer ${access_token}”。access_token为获取Token中获取的access_token。 |

### query

[参数来源](https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-query-compile-status-0000002236041434#section165302018915)

| 参数名称 | 必选(M)/可选(O) | 类型 | 参数说明 |
| --- | --- | --- | --- |
| appId | M | String(32) | 应用ID。 |
| pkgIds | M | String(500) | 待查询的软件包ID。 多个ID之间用英文逗号分隔。 可从更新应用软件包信息接口返回的packageId参数中获取。 |

### 请求说明

本接口支持使用Service Account方式和API客户端方式，二者区别请参见获取服务端授权。

### 完整模型及示例

```sh
agc publishing package-compile-status-v3 --docs --pretty
```

该命令返回官方请求/响应示例及已捕获的关联模型。业务响应中的 `ret.code` 应与 HTTP 状态分别检查；修改后使用相应查询接口确认目标状态。

## 查询应用所有版本类型的版本列表

命令：`agc publishing app-version-list`。请求：`POST /api/publish/v3/version/brief-info/list`。

[官方正文](https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-query-brief-info-list-0000002236201254)；官网更新：2026-07-31 17:44；浏览器读取：2026-10-07T04:56:06.794Z。

此接口用于获取HarmonyOS应用/元服务的所有版本类型（含商用版本和测试版本）的版本列表，帮助开发者了解版本发布状态。

接口调用者的角色：账号持有者、管理员、APP管理员、运营。

### header（Service Account方式：）

[参数来源](https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-query-brief-info-list-0000002236201254#section236517415165)

| 参数名称 | 必选(M)/可选(O) | 类型 | 参数说明 |
| --- | --- | --- | --- |
| Authorization | M | String | 认证信息，格式为“Authorization: Bearer ${JWT}”。JWT为通过Service Account方式获取授权中获取的鉴权令牌。 |
| appId | M | String(32) | 应用ID，获取方法参考查看应用信息。 |

### header（API客户端方式：）

[参数来源](https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-query-brief-info-list-0000002236201254#section236517415165)

| 参数名称 | 必选(M)/可选(O) | 类型 | 参数说明 |
| --- | --- | --- | --- |
| client_id | M | String | 客户端ID，获取方法参考创建API客户端。 |
| Authorization | M | String | 认证信息，格式为“Authorization: Bearer ${access_token}”。access_token为获取Token中获取的access_token。 |
| appId | M | String(32) | 应用ID，获取方法参考查看应用信息。 |

### body

[参数来源](https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-query-brief-info-list-0000002236201254#section2357705386)

| 参数名称 | 必选(M)/可选(O) | 类型 | 参数说明 |
| --- | --- | --- | --- |
| packageName | O | String(128) | 应用包名。 说明 appId和packageName至少传一个。 |
| state | O | String(64) | 版本状态。 多个状态用逗号分隔。 说明 state字段的取值请参考AppInfo中的releaseState参数。 |

### 请求说明

本接口支持使用Service Account方式和API客户端方式，二者区别请参见获取服务端授权。

请求Body中使用JSON格式携带查询应用所有版本类型的版本列表，参数如下表所示。

### 完整模型及示例

```sh
agc publishing app-version-list --docs --pretty
```

该命令返回官方请求/响应示例及已捕获的关联模型。业务响应中的 `ret.code` 应与 HTTP 状态分别检查；修改后使用相应查询接口确认目标状态。

## 查询应用未处理的问题工单

命令：`agc publishing app-version-voc`。请求：`POST /api/publish/v1/app-voc/version-voc`。

[官方正文](https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-post-app-version-voc-0000002495233986)；官网更新：2026-06-25 19:52；浏览器读取：2026-10-07T04:56:19.563Z。

此接口用于查询HarmonyOS应用/元服务未处理的问题工单。

调用本接口前必须保证应用软件包已经上传。

接口调用者的角色：账号持有者、管理员、APP管理员、开发。

### header（Service Account方式：）

[参数来源](https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-post-app-version-voc-0000002495233986#section236517415165)

| 参数名称 | 必选(M)/可选(O) | 类型 | 参数说明 |
| --- | --- | --- | --- |
| Authorization | M | String | 认证信息，格式为“Authorization: Bearer ${JWT}”。JWT为通过Service Account方式获取授权中获取的鉴权令牌。 |
| appId | M | String(32) | 应用ID，获取方法参考查看应用信息。 |

### header（API客户端方式：）

[参数来源](https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-post-app-version-voc-0000002495233986#section236517415165)

| 参数名称 | 必选(M)/可选(O) | 类型 | 参数说明 |
| --- | --- | --- | --- |
| client_id | M | String | 客户端ID，获取方法参考创建API客户端。 |
| Authorization | M | String | 认证信息，格式为“Authorization: Bearer ${access_token}”。access_token为获取Token中获取的access_token。 |
| appId | M | String(32) | 应用ID，获取方法参考查看应用信息。 |

### body

[参数来源](https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-post-app-version-voc-0000002495233986#section11431111271820)

| 参数名称 | 必选(M)/可选(O) | 类型 | 参数说明 |
| --- | --- | --- | --- |
| vocId | O | String(32) | 工单ID。 |
| vocDesc | O | String(255) | 工单描述。 |
| pageNo | M | Integer(32) | 当前页。 最小值：1 |
| pageSize | M | Integer(32) | 每页数量。 最大值：100 |

### 请求说明

本接口支持使用Service Account方式和API客户端方式，二者区别请参见获取服务端授权。

### 完整模型及示例

```sh
agc publishing app-version-voc --docs --pretty
```

该命令返回官方请求/响应示例及已捕获的关联模型。业务响应中的 `ret.code` 应与 HTTP 状态分别检查；修改后使用相应查询接口确认目标状态。

## 查询需要上传的资质文件信息

命令：`agc publishing tag-qual-dependent`。请求：`GET /api/publish/v2/tag-qual-dependent`。

[官方正文](https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-query-tag-qual-dependent-0000002271160577)；官网更新：2026-06-25 19:52；浏览器读取：2026-10-07T04:56:24.391Z。

此接口用于查询元服务需要上传的所有资质文件信息。

接口调用者的角色：账号持有者、管理员、APP管理员、运营。

### header（Service Account方式：）

[参数来源](https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-query-tag-qual-dependent-0000002271160577#section119494174419)

| 参数名称 | 必选(M)/可选(O) | 类型 | 参数说明 |
| --- | --- | --- | --- |
| Authorization | M | String | 认证信息，格式为“Authorization: Bearer ${JWT}”。JWT为通过Service Account方式获取授权中获取的鉴权令牌。 |

### header（API客户端方式：）

[参数来源](https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-query-tag-qual-dependent-0000002271160577#section119494174419)

| 参数名称 | 必选(M)/可选(O) | 类型 | 参数说明 |
| --- | --- | --- | --- |
| client_id | M | String | 客户端ID，获取方法参考创建API客户端。 |
| Authorization | M | String | 认证信息，格式为“Authorization: Bearer ${access_token}”。access_token为获取Token中获取的access_token。 |

### 请求说明

本接口支持使用Service Account方式和API客户端方式，二者区别请参见获取服务端授权。

### 完整模型及示例

```sh
agc publishing tag-qual-dependent --docs --pretty
```

该命令返回官方请求/响应示例及已捕获的关联模型。业务响应中的 `ret.code` 应与 HTTP 状态分别检查；修改后使用相应查询接口确认目标状态。

## 提交资质审核任务

命令：`agc publishing qualification-review-create`。请求：`POST /api/publish/v2/qualification-review-task`。

[官方正文](https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-post-qualification-review-0000002271000629)；官网更新：2026-06-25 19:52；浏览器读取：2026-10-07T04:56:28.028Z。

此接口用于提交元服务的资质审核任务。

应用设备类型仅选择智能手表或运动手表时，不支持提交资质审核任务。

接口调用者的角色：账号持有者、管理员、APP管理员、运营。

### header（Service Account方式：）

[参数来源](https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-post-qualification-review-0000002271000629#section236517415165)

| 参数名称 | 必选(M)/可选(O) | 类型 | 参数说明 |
| --- | --- | --- | --- |
| Authorization | M | String | 认证信息，格式为“Authorization: Bearer ${JWT}”。JWT为通过Service Account方式获取授权中获取的鉴权令牌。 |
| appId | M | String(32) | 应用ID，获取方法参考查看应用信息。 |

### header（API客户端方式：）

[参数来源](https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-post-qualification-review-0000002271000629#section236517415165)

| 参数名称 | 必选(M)/可选(O) | 类型 | 参数说明 |
| --- | --- | --- | --- |
| client_id | M | String | 客户端ID，获取方法参考创建API客户端。 |
| Authorization | M | String | 认证信息，格式为“Authorization: Bearer ${access_token}”。access_token为获取Token中获取的access_token。 |
| appId | M | String(32) | 应用ID，获取方法参考查看应用信息。 |

### body

[参数来源](https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-post-qualification-review-0000002271000629#section2357705386)

| 参数名称 | 必选(M)/可选(O) | 类型 | 参数说明 |
| --- | --- | --- | --- |
| tagQualInfos | M | List<TagQualInfo> | 标签和标签对应的资质信息列表。 数组长度不超过5。 |

### 请求说明

本接口支持使用Service Account方式和API客户端方式，二者区别请参见获取服务端授权。

### 完整模型及示例

```sh
agc publishing qualification-review-create --docs --pretty
```

该命令返回官方请求/响应示例及已捕获的关联模型。业务响应中的 `ret.code` 应与 HTTP 状态分别检查；修改后使用相应查询接口确认目标状态。

## 查询资质审核任务详情

命令：`agc publishing qualification-review-query`。请求：`GET /api/publish/v2/qualification-review-task`。

[官方正文](https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-get-qualification-review-0000002236041438)；官网更新：2026-06-25 19:52；浏览器读取：2026-10-07T04:56:44.080Z。

此接口用于查询元服务的资质审核任务详情。

接口调用者的角色：账号持有者、管理员、APP管理员、运营。

### header（Service Account方式：）

[参数来源](https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-get-qualification-review-0000002236041438#section236517415165)

| 参数名称 | 必选(M)/可选(O) | 类型 | 参数说明 |
| --- | --- | --- | --- |
| Authorization | M | String | 认证信息，格式为“Authorization: Bearer ${JWT}”。JWT为通过Service Account方式获取授权中获取的鉴权令牌。 |
| appId | M | String(32) | 应用ID，获取方法参考查看应用信息。 |

### header（API客户端方式：）

[参数来源](https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-get-qualification-review-0000002236041438#section236517415165)

| 参数名称 | 必选(M)/可选(O) | 类型 | 参数说明 |
| --- | --- | --- | --- |
| client_id | M | String | 客户端ID，获取方法参考创建API客户端。 |
| Authorization | M | String | 认证信息，格式为“Authorization: Bearer ${access_token}”。access_token为获取Token中获取的access_token。 |
| appId | M | String(32) | 应用ID，获取方法参考查看应用信息。 |

### query

[参数来源](https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-get-qualification-review-0000002236041438#section165302018915)

| 参数名称 | 必选(M)/可选(O) | 类型 | 参数说明 |
| --- | --- | --- | --- |
| tagIds | M | String(256) | 审核任务包含的标签ID列表，逗号分隔。 |

### 请求说明

本接口支持使用Service Account方式和API客户端方式，二者区别请参见获取服务端授权。

### 完整模型及示例

```sh
agc publishing qualification-review-query --docs --pretty
```

该命令返回官方请求/响应示例及已捕获的关联模型。业务响应中的 `ret.code` 应与 HTTP 状态分别检查；修改后使用相应查询接口确认目标状态。

## 修改资质审核任务

命令：`agc publishing qualification-review-update`。请求：`PUT /api/publish/v2/qualification-review-task`。

[官方正文](https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-put-qualification-review-0000002422530746)；官网更新：2026-06-25 19:52；浏览器读取：2026-10-07T04:56:45.618Z。

此接口用于修改元服务的资质审核任务。

接口调用者的角色：账号持有者、管理员、APP管理员、运营。

### header（Service Account方式：）

[参数来源](https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-put-qualification-review-0000002422530746#section236517415165)

| 参数名称 | 必选(M)/可选(O) | 类型 | 参数说明 |
| --- | --- | --- | --- |
| Authorization | M | String | 认证信息，格式为“Authorization: Bearer ${JWT}”。JWT为通过Service Account方式获取授权中获取的鉴权令牌。 |
| appId | M | String(32) | 应用ID，获取方法参考查看应用信息。 |

### header（API客户端方式：）

[参数来源](https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-put-qualification-review-0000002422530746#section236517415165)

| 参数名称 | 必选(M)/可选(O) | 类型 | 参数说明 |
| --- | --- | --- | --- |
| client_id | M | String | 客户端ID，获取方法参考创建API客户端。 |
| Authorization | M | String | 认证信息，格式为“Authorization: Bearer ${access_token}”。access_token为获取Token中获取的access_token。 |
| appId | M | String(32) | 应用ID，获取方法参考查看应用信息。 |

### body

[参数来源](https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-put-qualification-review-0000002422530746#section2357705386)

| 参数名称 | 必选(M)/可选(O) | 类型 | 参数说明 |
| --- | --- | --- | --- |
| tagQualInfos | M | List<TagQualInfo> | 标签和标签对应的资质信息列表。 数组长度不超过5。 |

### 请求说明

本接口支持使用Service Account方式和API客户端方式，二者区别请参见获取服务端授权。

### 完整模型及示例

```sh
agc publishing qualification-review-update --docs --pretty
```

该命令返回官方请求/响应示例及已捕获的关联模型。业务响应中的 `ret.code` 应与 HTTP 状态分别检查；修改后使用相应查询接口确认目标状态。

## 撤销资质审核任务

命令：`agc publishing qualification-review-cancel`。请求：`PUT /api/publish/v2/qualification-review-task/cancel`。

[官方正文](https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-cancel-qualification-review-0000002236201258)；官网更新：2026-06-25 19:52；浏览器读取：2026-10-07T04:56:47.305Z。

此接口用于撤销元服务的资质审核任务。

接口调用者的角色：账号持有者、管理员、APP管理员、运营。

### header（Service Account方式：）

[参数来源](https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-cancel-qualification-review-0000002236201258#section236517415165)

| 参数名称 | 必选(M)/可选(O) | 类型 | 参数说明 |
| --- | --- | --- | --- |
| Authorization | M | String | 认证信息，格式为“Authorization: Bearer ${JWT}”。JWT为通过Service Account方式获取授权中获取的鉴权令牌。 |
| appId | M | String(32) | 应用ID，获取方法参考查看应用信息。 |

### header（API客户端方式：）

[参数来源](https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-cancel-qualification-review-0000002236201258#section236517415165)

| 参数名称 | 必选(M)/可选(O) | 类型 | 参数说明 |
| --- | --- | --- | --- |
| client_id | M | String | 客户端ID，获取方法参考创建API客户端。 |
| Authorization | M | String | 认证信息，格式为“Authorization: Bearer ${access_token}”。access_token为获取Token中获取的access_token。 |
| appId | M | String(32) | 应用ID，获取方法参考查看应用信息。 |

### body

[参数来源](https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-cancel-qualification-review-0000002236201258#section2357705386)

| 参数名称 | 必选(M)/可选(O) | 类型 | 参数说明 |
| --- | --- | --- | --- |
| taskId | M | String(32) | 资质审核任务ID。 |

### 请求说明

本接口支持使用Service Account方式和API客户端方式，二者区别请参见获取服务端授权。

### 完整模型及示例

```sh
agc publishing qualification-review-cancel --docs --pretty
```

该命令返回官方请求/响应示例及已捕获的关联模型。业务响应中的 `ret.code` 应与 HTTP 状态分别检查；修改后使用相应查询接口确认目标状态。

## 删除资质审核任务

命令：`agc publishing qualification-review-delete`。请求：`DELETE /api/publish/v2/qualification-review-task`。

[官方正文](https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-delete-qualification-review-0000002526538788)；官网更新：2026-06-25 19:52；浏览器读取：2026-10-07T04:57:03.034Z。

此接口用于删除元服务的资质审核任务。

如果元服务资质审核已通过，但在提交上架时因资质问题被驳回，需要先调用该接口移除旧的资质材料，然后重新提交资质审核任务。

接口调用者的角色：账号持有者、管理员、APP管理员、运营。

### header（Service Account方式：）

[参数来源](https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-delete-qualification-review-0000002526538788#section236517415165)

| 参数名称 | 必选(M)/可选(O) | 类型 | 参数说明 |
| --- | --- | --- | --- |
| Authorization | M | String | 认证信息，格式为“Authorization: Bearer ${JWT}”。JWT为通过Service Account方式获取授权中获取的鉴权令牌。 |
| appId | M | String(32) | 应用ID，获取方法参考查看应用信息。 |

### header（API客户端方式：）

[参数来源](https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-delete-qualification-review-0000002526538788#section236517415165)

| 参数名称 | 必选(M)/可选(O) | 类型 | 参数说明 |
| --- | --- | --- | --- |
| client_id | M | String | 客户端ID，获取方法参考创建API客户端。 |
| Authorization | M | String | 认证信息，格式为“Authorization: Bearer ${access_token}”。access_token为获取Token中获取的access_token。 |
| appId | M | String(32) | 应用ID，获取方法参考查看应用信息。 |

### body

[参数来源](https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-delete-qualification-review-0000002526538788#section2357705386)

| 参数名称 | 必选(M)/可选(O) | 类型 | 参数说明 |
| --- | --- | --- | --- |
| tagIds | M | List<String(64)> | 需要删除的资质审核标签ID列表，逗号分隔。 数组长度不超过5。 说明 具体取值请参见元服务的“标签ID”列。 |

### 请求说明

本接口支持使用Service Account方式和API客户端方式，二者区别请参见获取服务端授权。

### 完整模型及示例

```sh
agc publishing qualification-review-delete --docs --pretty
```

该命令返回官方请求/响应示例及已捕获的关联模型。业务响应中的 `ret.code` 应与 HTTP 状态分别检查；修改后使用相应查询接口确认目标状态。

## 新建用户协议

命令：`agc publishing user-agreement-create`。请求：`POST /api/publish/v2/agreement`。

[官方正文](https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-create-user-agreement-0000002271160581)；官网更新：2026-06-25 19:52；浏览器读取：2026-10-07T04:57:04.732Z。

此接口用于新建HarmonyOS应用/元服务的用户协议。

接口调用者的角色：账号持有者、管理员、APP管理员、运营。

### header（Service Account方式：）

[参数来源](https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-create-user-agreement-0000002271160581#section236517415165)

| 参数名称 | 必选(M)/可选(O) | 类型 | 参数说明 |
| --- | --- | --- | --- |
| Authorization | M | String | 认证信息，格式为“Authorization: Bearer ${JWT}”。JWT为通过Service Account方式获取授权中获取的鉴权令牌。 |
| appId | M | String(32) | 应用ID，获取方法参考查看应用信息。 |

### header（API客户端方式：）

[参数来源](https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-create-user-agreement-0000002271160581#section236517415165)

| 参数名称 | 必选(M)/可选(O) | 类型 | 参数说明 |
| --- | --- | --- | --- |
| client_id | M | String | 客户端ID，获取方法参考创建API客户端。 |
| Authorization | M | String | 认证信息，格式为“Authorization: Bearer ${access_token}”。access_token为获取Token中获取的access_token。 |
| appId | M | String(32) | 应用ID，获取方法参考查看应用信息。 |

### body

[参数来源](https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-create-user-agreement-0000002271160581#section11431111271820)

| 参数名称 | 必选(M)/可选(O) | 类型 | 参数说明 |
| --- | --- | --- | --- |
| type | M | Integer(32) | 协议类型。 取值范围： 2：用户协议 |
| defaultLang | O | String(20) | 用户协议默认语言。 取值范围： zh：简体中文 en：英语 ug：维吾尔语 bo：藏语 zh-Hant-HK：繁体中文（中国香港特别行政区） zh-Hant-TW：繁体中文（中国台湾） |
| agreementInfos | O | AgreementInfo | 用户协议列表。 数组长度不超过100。 |

### 请求说明

本接口支持使用Service Account方式和API客户端方式，二者区别请参见获取服务端授权。

### 完整模型及示例

```sh
agc publishing user-agreement-create --docs --pretty
```

该命令返回官方请求/响应示例及已捕获的关联模型。业务响应中的 `ret.code` 应与 HTTP 状态分别检查；修改后使用相应查询接口确认目标状态。

## 查询用户协议列表

命令：`agc publishing user-agreement-list`。请求：`POST /api/publish/v2/agreement/list`。

[官方正文](https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-query-user-agreement-0000002236041442)；官网更新：2026-06-25 19:52；浏览器读取：2026-10-07T04:57:06.324Z。

此接口用于查询HarmonyOS应用/元服务的用户协议信息。

接口调用者的角色：账号持有者、管理员、APP管理员、运营。

### header（Service Account方式：）

[参数来源](https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-query-user-agreement-0000002236041442#section236517415165)

| 参数名称 | 必选(M)/可选(O) | 类型 | 参数说明 |
| --- | --- | --- | --- |
| Authorization | M | String | 认证信息，格式为“Authorization: Bearer ${JWT}”。JWT为通过Service Account方式获取授权中获取的鉴权令牌。 |
| appId | M | String(32) | 应用ID，获取方法参考查看应用信息。 |

### header（API客户端方式：）

[参数来源](https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-query-user-agreement-0000002236041442#section236517415165)

| 参数名称 | 必选(M)/可选(O) | 类型 | 参数说明 |
| --- | --- | --- | --- |
| client_id | M | String | 客户端ID，获取方法参考创建API客户端。 |
| Authorization | M | String | 认证信息，格式为“Authorization: Bearer ${access_token}”。access_token为获取Token中获取的access_token。 |
| appId | M | String(32) | 应用ID，获取方法参考查看应用信息。 |

### body

[参数来源](https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-query-user-agreement-0000002236041442#section2357705386)

| 参数名称 | 必选(M)/可选(O) | 类型 | 参数说明 |
| --- | --- | --- | --- |
| agreementIdList | O | List<String(64)> | 用户协议ID。 数组长度为1~1000。 |
| name | M | String(100) | 用户协议名称。 使用模糊查询方式。 |
| status | O | Integer(32) | 用户协议状态。 取值范围： 2：已完成 3：已删除 |
| type | O | Integer(32) | 协议类型。 取值范围： 2：用户协议 |
| fromRecCount | O | Integer(32) | 查询起始位置。 默认值：1 |
| maxReqCount | O | Integer | 最大请求查询数。 取值范围：[1,100] 默认值：30 |

### 请求说明

本接口支持使用Service Account方式和API客户端方式，二者区别请参见获取服务端授权。

### 完整模型及示例

```sh
agc publishing user-agreement-list --docs --pretty
```

该命令返回官方请求/响应示例及已捕获的关联模型。业务响应中的 `ret.code` 应与 HTTP 状态分别检查；修改后使用相应查询接口确认目标状态。

## 更新用户协议

命令：`agc publishing user-agreement-update`。请求：`PUT /api/publish/v2/agreement`。

[官方正文](https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-update-user-agreement-0000002294845626)；官网更新：2026-06-25 19:52；浏览器读取：2026-10-07T04:57:18.109Z。

此接口用于更新HarmonyOS应用/元服务的用户协议信息。

接口调用者的角色：账号持有者、管理员、APP管理员、运营。

### header（Service Account方式：）

[参数来源](https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-update-user-agreement-0000002294845626#section236517415165)

| 参数名称 | 必选(M)/可选(O) | 类型 | 参数说明 |
| --- | --- | --- | --- |
| Authorization | M | String | 认证信息，格式为“Authorization: Bearer ${JWT}”。JWT为通过Service Account方式获取授权中获取的鉴权令牌。 |
| appId | M | String(32) | 应用ID，获取方法参考查看应用信息。 |

### header（API客户端方式：）

[参数来源](https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-update-user-agreement-0000002294845626#section236517415165)

| 参数名称 | 必选(M)/可选(O) | 类型 | 参数说明 |
| --- | --- | --- | --- |
| client_id | M | String | 客户端ID，获取方法参考创建API客户端。 |
| Authorization | M | String | 认证信息，格式为“Authorization: Bearer ${access_token}”。access_token为获取Token中获取的access_token。 |
| appId | M | String(32) | 应用ID，获取方法参考查看应用信息。 |

### body

[参数来源](https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-update-user-agreement-0000002294845626#section11431111271820)

| 参数名称 | 必选(M)/可选(O) | 类型 | 参数说明 |
| --- | --- | --- | --- |
| id | M | String(32) | 用户协议ID。 |
| defaultLang | O | String(20) | 用户协议默认语言。 取值范围： zh：简体中文 en：英语 ug：维吾尔语 bo：藏语 zh-Hant-HK：繁体中文（中国香港特别行政区） zh-Hant-TW：繁体中文（中国台湾） |
| agreementInfos | O | List<AgreementInfo> | 用户协议列表。 数组长度不超过100。 |

### 请求说明

本接口支持使用Service Account方式和API客户端方式，二者区别请参见获取服务端授权。

### 完整模型及示例

```sh
agc publishing user-agreement-update --docs --pretty
```

该命令返回官方请求/响应示例及已捕获的关联模型。业务响应中的 `ret.code` 应与 HTTP 状态分别检查；修改后使用相应查询接口确认目标状态。

## 提交隐私政策协议

命令：`agc publishing privacy-agreement-submit`。请求：`PUT /api/publish/v2/agreement/submit`。

[官方正文](https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-put-privacy-agreement-0000002271000633)；官网更新：2026-06-25 19:52；浏览器读取：2026-10-07T04:57:20.479Z。

此接口用于提交HarmonyOS应用/元服务的隐私政策协议。

接口调用者的角色：账号持有者、管理员、APP管理员、运营。

### header（Service Account方式：）

[参数来源](https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-put-privacy-agreement-0000002271000633#section236517415165)

| 参数名称 | 必选(M)/可选(O) | 类型 | 参数说明 |
| --- | --- | --- | --- |
| Authorization | M | String | 认证信息，格式为“Authorization: Bearer ${JWT}”。JWT为通过Service Account方式获取授权中获取的鉴权令牌。 |
| appId | M | String(32) | 应用ID，获取方法参考查看应用信息。 |

### header（API客户端方式：）

[参数来源](https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-put-privacy-agreement-0000002271000633#section236517415165)

| 参数名称 | 必选(M)/可选(O) | 类型 | 参数说明 |
| --- | --- | --- | --- |
| client_id | M | String | 客户端ID，获取方法参考创建API客户端。 |
| Authorization | M | String | 认证信息，格式为“Authorization: Bearer ${access_token}”。access_token为获取Token中获取的access_token。 |
| appId | M | String(32) | 应用ID，获取方法参考查看应用信息。 |

### body

[参数来源](https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-put-privacy-agreement-0000002271000633#section1131102811712)

| 参数名称 | 必选(M)/可选(O) | 类型 | 参数说明 |
| --- | --- | --- | --- |
| id | M | String(32) | 隐私政策协议ID。 |
| requireConfirm | O | Integer(32) | 是否弹窗确认。 取值范围： 1：是 0：否 |

### 请求说明

本接口支持使用Service Account方式和API客户端方式，二者区别请参见获取服务端授权。

### 完整模型及示例

```sh
agc publishing privacy-agreement-submit --docs --pretty
```

该命令返回官方请求/响应示例及已捕获的关联模型。业务响应中的 `ret.code` 应与 HTTP 状态分别检查；修改后使用相应查询接口确认目标状态。

## 查询隐私政策协议列表

命令：`agc publishing privacy-agreement-list`。请求：`POST /api/publish/v2/agreement/list`。

[官方正文](https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-query-privacy-agreement-0000002328924925)；官网更新：2026-06-25 19:52；浏览器读取：2026-10-07T05:10:26.407Z。

此接口用于查询HarmonyOS应用/元服务的隐私政策协议信息。

接口调用者的角色：账号持有者、管理员、APP管理员、运营。

### header（Service Account方式：）

[参数来源](https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-query-privacy-agreement-0000002328924925#section236517415165)

| 参数名称 | 必选(M)/可选(O) | 类型 | 参数说明 |
| --- | --- | --- | --- |
| Authorization | M | String | 认证信息，格式为“Authorization: Bearer ${JWT}”。JWT为通过Service Account方式获取授权中获取的鉴权令牌。 |
| appId | M | String(32) | 应用ID，获取方法参考查看应用信息。 |

### header（API客户端方式：）

[参数来源](https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-query-privacy-agreement-0000002328924925#section236517415165)

| 参数名称 | 必选(M)/可选(O) | 类型 | 参数说明 |
| --- | --- | --- | --- |
| client_id | M | String | 客户端ID，获取方法参考创建API客户端。 |
| Authorization | M | String | 认证信息，格式为“Authorization: Bearer ${access_token}”。access_token为获取Token中获取的access_token。 |
| appId | M | String(32) | 应用ID，获取方法参考查看应用信息。 |

### body

[参数来源](https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-query-privacy-agreement-0000002328924925#section2357705386)

| 参数名称 | 必选(M)/可选(O) | 类型 | 参数说明 |
| --- | --- | --- | --- |
| agreementIdList | O | List<String(64)> | 隐私政策协议ID。 数组长度为1~1000。 |
| name | M | String(100) | 隐私政策协议名称。 使用模糊查询方式。 |
| status | O | Integer(32) | 隐私政策协议状态。 取值范围： 1：草稿 2：已完成 3：已删除 不传值：默认查未删除的记录 |
| type | O | Integer(32) | 协议类型。 取值范围： 1：隐私政策协议 |
| fromRecCount | O | Integer(32) | 起始页数。 默认值：1 |
| maxReqCount | O | Integer | 每页返回条数。 取值范围：[1,100] 默认值：30 |

### 请求说明

本接口支持使用Service Account方式和API客户端方式，二者区别请参见获取服务端授权。

### 完整模型及示例

```sh
agc publishing privacy-agreement-list --docs --pretty
```

该命令返回官方请求/响应示例及已捕获的关联模型。业务响应中的 `ret.code` 应与 HTTP 状态分别检查；修改后使用相应查询接口确认目标状态。

## 更新隐私政策协议

命令：`agc publishing privacy-agreement-update`。请求：`PUT /api/publish/v2/agreement`。

[官方正文](https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-update-privacy-agreement-0000002328805169)；官网更新：2026-06-25 19:52；浏览器读取：2026-10-07T05:06:42.704Z。

此接口用于更新HarmonyOS应用/元服务的隐私政策协议。接口调用者的角色：账号持有者、管理员、APP管理员、运营。

若隐私政策协议被在架版本关联，编辑后需要提交华为运营审核，待协议审核通过后，更新的内容方可生效。

若隐私政策协议被状态为“审核中”的应用关联，将不支持编辑。

### header（Service Account方式：）

[参数来源](https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-update-privacy-agreement-0000002328805169#section236517415165)

| 参数名称 | 必选(M)/可选(O) | 类型 | 参数说明 |
| --- | --- | --- | --- |
| Authorization | M | String | 认证信息，格式为“Authorization: Bearer ${JWT}”。JWT为通过Service Account方式获取授权中获取的鉴权令牌。 |
| appId | M | String(32) | 应用ID，获取方法参考查看应用信息。 |

### header（API客户端方式：）

[参数来源](https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-update-privacy-agreement-0000002328805169#section236517415165)

| 参数名称 | 必选(M)/可选(O) | 类型 | 参数说明 |
| --- | --- | --- | --- |
| client_id | M | String | 客户端ID，获取方法参考创建API客户端。 |
| Authorization | M | String | 认证信息，格式为“Authorization: Bearer ${access_token}”。access_token为获取Token中获取的access_token。 |
| appId | M | String(32) | 应用ID，获取方法参考查看应用信息。 |

### body

[参数来源](https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-update-privacy-agreement-0000002328805169#section11431111271820)

| 参数名称 | 必选(M)/可选(O) | 类型 | 参数说明 |
| --- | --- | --- | --- |
| id | M | String(32) | 隐私政策协议ID。 |
| defaultLang | O | String(20) | 隐私政策协议默认语言。 取值范围： zh：简体中文 en：英语 ug：维吾尔语 bo：藏语 zh-Hant-HK：繁体中文（中国香港特别行政区） zh-Hant-TW：繁体中文（中国台湾） |
| agreementInfos | O | List<AgreementInfo> | 隐私政策协议列表。 数组长度不超过100。 |

### 请求说明

本接口支持使用Service Account方式和API客户端方式，二者区别请参见获取服务端授权。

### 完整模型及示例

```sh
agc publishing privacy-agreement-update --docs --pretty
```

该命令返回官方请求/响应示例及已捕获的关联模型。业务响应中的 `ret.code` 应与 HTTP 状态分别检查；修改后使用相应查询接口确认目标状态。

## 申请年龄分级

命令：`agc publishing app-age-rating`。请求：`POST /api/publish/v2/app-age-rating`。

[官方正文](https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-post-app-age-rating-0000002236201262)；官网更新：2026-06-25 19:52；浏览器读取：2026-10-07T04:57:45.706Z。

此接口用于申请或修改HarmonyOS应用/元服务的年龄分级。

接口调用者的角色：账号持有者、管理员、APP管理员、运营。

调用本接口前需完成应用分类设置，明确是游戏类应用还是普通应用。

### header（认证方式以正文说明为准）

[参数来源](https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-post-app-age-rating-0000002236201262#section236517415165)

| 参数名称 | 必选(M)/可选(O) | 类型 | 参数说明 |
| --- | --- | --- | --- |
| client_id | M | String | 客户端ID。 |
| Authorization | M | String | 认证信息。 格式：Authorization: Bearer ${access_token} |
| appId | M | String(32) | 应用ID。 获取方法参考查看应用信息。 |

### query

[参数来源](https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-post-app-age-rating-0000002236201262#section23242264517)

| 参数名称 | 必选(M)/可选(O) | 类型 | 参数说明 |
| --- | --- | --- | --- |
| releaseType | O | Integer(32) | 应用发布方式。 取值范围： 1：全网 默认值：1 |
| releasePhase | O | Integer(32) | 分阶段发布标识。 取值范围： 0：全网 3：分阶段 默认值：0 |

### body

[参数来源](https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-post-app-age-rating-0000002236201262#section2357705386)

| 参数名称 | 必选(M)/可选(O) | 类型 | 参数说明 |
| --- | --- | --- | --- |
| feedback | M | String(8192) | 问卷内容。 说明 具体请参见feedback获取方式。 |
| expectedRating | O | Integer(32) | 预期年龄分级。 取值范围： 1：巴西适用所有年龄段，中国区及其他区域年满3周岁 2：中国区年满8周岁，巴西年满10周岁，其他区域年满7周岁 3：巴西年满14周岁，中国区及其他区域年满12周岁 4：年满16周岁 5：年满18周岁 |
| ratingRule | O | Integer(32) | 年龄分级生效规则。 取值范围： 0：以问卷内容评级结果为准。 1：当预期年龄分级expectedRating设置的值大于问卷评级结果中的年龄值时，以预期年龄分级为准，否则以问卷内容评级结果为准。 |
| childFlag | O | Integer(32) | 应用是否仅面向儿童。 取值范围： 0：应用不仅面向儿童 1：应用仅面向儿童 |

### 请求说明

本接口请使用API客户端认证方式，请参见API客户端方式获取服务端授权获取客户端ID和access_token。

feedback获取方式

该问卷为独立页面，需要开发者通过window.open()方法打开，并监听message事件获取问卷填写内容。

调用window.open()方法打开页面，代码如下： window.open("https://developer.huawei.com/consumer/cn/service/josp/agc/amp/ageRatingQuestionnaire.html#/?appType={appType}&lang={lang}") 其中，appType和lang为参数变量。 appType可选值为app/game，app为应用年龄分级问卷，game为游戏年龄分级问卷。 lang可选值为cn/en/ru，cn为中文问卷，en为英文问卷，ru为俄文问卷。 说明 如果无法打开问卷，请先退出AppGallery Connect登录，然后重新访问该页面。

window.open("https://developer.huawei.com/consumer/cn/service/josp/agc/amp/ageRatingQuestionnaire.html#/?appType={appType}&lang={lang}")

其中，appType和lang为参数变量。

appType可选值为app/game，app为应用年龄分级问卷，game为游戏年龄分级问卷。

lang可选值为cn/en/ru，cn为中文问卷，en为英文问卷，ru为俄文问卷。

如果无法打开问卷，请先退出AppGallery Connect登录，然后重新访问该页面。

填写完问卷后点击提交，问卷页面会发起postMessage请求，调用方通过监听message事件获取问卷内容feedback，代码如下： window.addEventListener("message", (event) => { console.log(event.data); //问卷内容feedback })

window.addEventListener("message", (event) => {

console.log(event.data); //问卷内容feedback

})

### 完整模型及示例

```sh
agc publishing app-age-rating --docs --pretty
```

该命令返回官方请求/响应示例及已捕获的关联模型。业务响应中的 `ret.code` 应与 HTTP 状态分别检查；修改后使用相应查询接口确认目标状态。

## 提交发布

命令：`agc publishing app-submit-v3`。请求：`POST /api/publish/v3/app-submit`。

[官方正文](https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-app-submit-0000002271160585)；官网更新：2026-06-25 19:52；浏览器读取：2026-10-07T04:57:47.737Z。

此接口用于提交HarmonyOS应用/元服务审核，调用本接口前必须保证应用信息已补充完整，应用软件包已经上传。

接口调用者的角色：账号持有者、管理员、APP管理员。

软件包采用异步解析方式，请您在传包后等候2分钟再调用提交发布接口。

### header（Service Account方式：）

[参数来源](https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-app-submit-0000002271160585#section119494174419)

| 参数名称 | 必选(M)/可选(O) | 类型 | 参数说明 |
| --- | --- | --- | --- |
| Authorization | M | String | 认证信息，格式为“Authorization: Bearer ${JWT}”。JWT为通过Service Account方式获取授权中获取的鉴权令牌。 |

### header（API客户端方式：）

[参数来源](https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-app-submit-0000002271160585#section119494174419)

| 参数名称 | 必选(M)/可选(O) | 类型 | 参数说明 |
| --- | --- | --- | --- |
| client_id | M | String | 客户端ID，获取方法参考创建API客户端。 |
| Authorization | M | String | 认证信息，格式为“Authorization: Bearer ${access_token}”。access_token为获取Token中获取的access_token。 |

### query

[参数来源](https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-app-submit-0000002271160585#section292816392364)

| 参数名称 | 必选(M)/可选(O) | 类型 | 参数说明 |
| --- | --- | --- | --- |
| appId | M | String(32) | 应用ID。 |

### body

[参数来源](https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-app-submit-0000002271160585#section11431111271820)

| 参数名称 | 必选(M)/可选(O) | 类型 | 参数说明 |
| --- | --- | --- | --- |
| releaseTime | O | String(64) | 指定发布时间。 UTC时间格式：yyyy-MM-ddTHH:mm:ssZZ，例如“2015-01-01T01:01:01+0800”。 可以不用填写，则表示立刻发布。 说明 如果在首次提交发布时已设置过releaseTime参数，撤销审核后再次提交发布时，需要将releaseTime参数重新设置为过去的时间，才能实现立即发布。 |
| remark | O | String(300) | 提审发布备注。 可以为空，如果填写则长度范围为10-300。 |
| releaseType | O | Integer(32) | 应用发布方式。 取值范围： 1：全网 默认值：1 |
| releasePhase | O | Integer(32) | 分阶段发布标识。 取值范围： 0：全网发布 3：分阶段发布（7天内自动更新） 默认值：0 |
| phasedReleaseDescription | O | String(500) | 分阶段发布说明。 说明 releasePhase取值为3（分阶段发布）时，此参数必填。 |

### 请求说明

本接口支持使用Service Account方式和API客户端方式，二者区别请参见获取服务端授权。

Body中使用JSON格式携带发布信息，参数如下表所示。

### 完整模型及示例

```sh
agc publishing app-submit-v3 --docs --pretty
```

该命令返回官方请求/响应示例及已捕获的关联模型。业务响应中的 `ret.code` 应与 HTTP 状态分别检查；修改后使用相应查询接口确认目标状态。

## 撤销审核

命令：`agc publishing on-shelf-cancel`。请求：`PUT /api/publish/v3/version/on-shelf/cancel`。

[官方正文](https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-on-shelf-cancel-0000002271000637)；官网更新：2026-06-25 19:52；浏览器读取：2026-10-07T04:58:00.948Z。

此接口用于HarmonyOS应用/元服务的撤销审核，当应用发布方式为“全网”时，只有处于“审核中”、“预审中”、“升级审核中”、“待上架/预约上架”、“申请下架”状态的应用可以撤销审核；当应用发布方式为“HarmonyOS测试发布方式（API>=10）”时，只有处于“正在审核”状态的应用可以撤销审核。

可通过查询应用信息接口的appInfo参数获取应用状态，具体请参见releaseState应用状态说明。

接口调用者的角色：账号持有者、管理员、APP管理员。

### header（Service Account方式：）

[参数来源](https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-on-shelf-cancel-0000002271000637#section236517415165)

| 参数名称 | 必选(M)/可选(O) | 类型 | 参数说明 |
| --- | --- | --- | --- |
| Authorization | M | String | 认证信息，格式为“Authorization: Bearer ${JWT}”。JWT为通过Service Account方式获取授权中获取的鉴权令牌。 |
| appId | M | String(32) | 应用ID，获取方法参考查看应用信息。 |

### header（API客户端方式：）

[参数来源](https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-on-shelf-cancel-0000002271000637#section236517415165)

| 参数名称 | 必选(M)/可选(O) | 类型 | 参数说明 |
| --- | --- | --- | --- |
| client_id | M | String | 客户端ID，获取方法参考创建API客户端。 |
| Authorization | M | String | 认证信息，格式为“Authorization: Bearer ${access_token}”。access_token为获取Token中获取的access_token。 |
| appId | M | String(32) | 应用ID，获取方法参考查看应用信息。 |

### body

[参数来源](https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-on-shelf-cancel-0000002271000637#section13576204051611)

| 参数名称 | 必选(M)/可选(O) | 类型 | 参数说明 |
| --- | --- | --- | --- |
| versionId | M | String | 版本ID。 |

### 请求说明

本接口支持使用Service Account方式和API客户端方式，二者区别请参见获取服务端授权。

请求Body中使用JSON格式携带撤销审核信息，参数如下表所示。

### 完整模型及示例

```sh
agc publishing on-shelf-cancel --docs --pretty
```

该命令返回官方请求/响应示例及已捕获的关联模型。业务响应中的 `ret.code` 应与 HTTP 状态分别检查；修改后使用相应查询接口确认目标状态。

## 下架应用

命令：`agc publishing app-shelf-off`。请求：`POST /api/publish/v2/app-shelf-off`。

[官方正文](https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-app-shelf-off-0000002236041450)；官网更新：2026-07-31 09:08；浏览器读取：2026-10-07T04:58:02.750Z。

此接口用于对已经上架的应用进行下架处理。

接口调用者的角色：账号持有者、管理员、APP管理员。

本接口仅对允许清单内用户开放，普通Publishing API开发者请勿调用本接口。

### header（Service Account方式：）

[参数来源](https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-app-shelf-off-0000002236041450#section119494174419)

| 参数名称 | 必选(M)/可选(O) | 类型 | 参数说明 |
| --- | --- | --- | --- |
| Authorization | M | String | 认证信息，格式为“Authorization: Bearer ${JWT}”。JWT为通过Service Account方式获取授权中获取的鉴权令牌。 |

### header（API客户端方式：）

[参数来源](https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-app-shelf-off-0000002236041450#section119494174419)

| 参数名称 | 必选(M)/可选(O) | 类型 | 参数说明 |
| --- | --- | --- | --- |
| client_id | M | String | 客户端ID，获取方法参考创建API客户端。 |
| Authorization | M | String | 认证信息，格式为“Authorization: Bearer ${access_token}”。access_token为获取Token中获取的access_token。 |

### query

[参数来源](https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-app-shelf-off-0000002236041450#section198511327144414)

| 参数名称 | 必选(M)/可选(O) | 类型 | 参数说明 |
| --- | --- | --- | --- |
| appId | M | String(32) | 需要下架的应用ID。 |
| reason | M | String(1024) | 下架的原因描述。 说明 中文字符需要进行URL编码后再传入。 |
| releaseType | O | Integer | 下架的方式。 取值范围： 1：全网下架。 默认值：1 |

### 请求说明

本接口支持使用Service Account方式和API客户端方式，二者区别请参见获取服务端授权。

### 完整模型及示例

```sh
agc publishing app-shelf-off --docs --pretty
```

该命令返回官方请求/响应示例及已捕获的关联模型。业务响应中的 `ret.code` 应与 HTTP 状态分别检查；修改后使用相应查询接口确认目标状态。
