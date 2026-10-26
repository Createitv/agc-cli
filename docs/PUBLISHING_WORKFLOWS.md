# Publishing API 使用与核验

Publishing 命令用于管理 HarmonyOS 应用和元服务的上架资料、软件包、资质、协议、年龄分级及发布状态。资料更新成功、软件包上传成功、审核通过和用户可下载是不同结果，必须分别查询确认。

官方入口：[Publishing API参考](https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-reference-0000002271160565)。本次浏览器确认入口更新时间为 2026-09-04 16:04，包含 31 项操作，以及数据模型和错误码入口。命令清单见 [CLI_USAGE.md](CLI_USAGE.md#现行-publishing-api-命令)。

每项用途、角色、参数表及前置条件见 [逐项 CLI 参考](PUBLISHING_API_REFERENCE.md)。

## 发现命令和参数

```sh
agc publishing endpoints --pretty
agc publishing user-agreement-create --docs --pretty
agc publishing privacy-agreement-update --docs --pretty
```

`--docs` 返回仓库保存的官方正文，包括参数表、说明、请求/响应示例和来源 URL。`relatedDocuments` 附带已捕获的关联模型，并沿模型链接继续展开。只附带已读过的文档；没有捕获的链接不代表已验证。`checkedAt` 是捕获日期，`updatedAt` 是官网更新日期。

请求参数按位置传入：`--query key=value` 对应 URL 查询参数，`--header key=value` 对应请求头，`--field key=value` 对应 JSON 字段。复杂对象数组使用 `--body request.json`。同名字段在不同操作中可能位于不同位置；例如基本信息查询使用 `appId` 查询参数，协议操作使用 `appId` 请求头。

实际调用要明确选择 `--profile PROFILE`。下文 `APP_ID`、`PROFILE`、协议 ID 和 URL 都是占位符，需要替换成目标资源的真实值。不要把鉴权令牌写入示例文件。认证方式依官方页面选择一个分支，不需要同时提供所有认证方式的请求头。

不带 `--invoke` 时，命令只展示接口定义，即使提供请求参数也不会执行请求。`--invoke --dry-run` 在本地构建并校验请求；实际发送必须使用 `--invoke --dry-run=false`。下文流程命令使用实际发送的写法；预检查时将 `--dry-run=false` 改成 `--dry-run`。预检查不创建资源，也不证明远端资源状态或权限。

## 用户协议：创建、查询、更新

用户协议描述用户使用应用的约定。它与隐私政策协议使用不同的数据模型。创建接口的 `type` 为 2；`agreementInfos` 是对象数组，而不是名为 AgreementInfo 的 JSON 包装对象。

创建请求 `user-agreement.json`：

```json
{
  "type": 2,
  "defaultLang": "zh",
  "agreementInfos": [
    {
      "lang": "zh",
      "baseInfo": {
        "name": "示例应用用户协议",
        "agreementTitle": "用户协议",
        "userAgreementUrl": "https://example.com/terms"
      }
    }
  ]
}
```

```sh
agc publishing user-agreement-create --invoke --dry-run=false --profile PROFILE \
  --header appId=APP_ID --body user-agreement.json --pretty
agc publishing user-agreement-list --invoke --dry-run=false --profile PROFILE \
  --header appId=APP_ID --field name=示例应用用户协议 --field type=2 --pretty
```

创建后，从响应及列表查询确认协议 ID。更新时在请求体中加入 `id`，使用 `user-agreement-update`；更新后再查询该 ID 的内容，确认目标语言和网址。不要仅以 HTTP 200 判断成功，还应检查华为业务码。

每个 `agreementInfos` 条目必须有 `lang`，支持 `zh`、`en`、`ug`、`bo`、`zh-Hant-HK`、`zh-Hant-TW`。这里的简体中文是 `zh`，不要套用应用语言资料的 `zh-CN`。最多 100 个条目。`baseInfo` 可省略，但提供时 `name` 必填且最多 100 字符；标题最多 100 字符，网址最多 512 字符。

规则来源：[AgreementInfo（用户协议）](https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-data-agreementinfo-user-0000002236041482)、[AgreementBaseInfo（用户协议）](https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-data-agreementbaseinfo-user-0000002236201290)。

## 隐私政策协议：查询、更新、提交

隐私政策说明应用怎样处理个人信息。该目录提供列表查询、更新和提交接口；不能把“新建用户协议”当成“新建隐私政策协议”。先查询现有隐私协议，确认 ID、状态和模板信息，再决定更新方式。

```sh
agc publishing privacy-agreement-list --invoke --dry-run=false --profile PROFILE \
  --header appId=APP_ID --field name= --field type=1 --pretty
```

隐私政策的 `baseInfo` 使用 `name`、`appName`、`companyName`、`briefDesc`、`effectiveTime`。这些字段与用户协议的标题和网址字段不同。`effectiveTime` 是整数毫秒时间戳，不能用字符串日期替代。

更新请求 `privacy-update.json` 的基础信息示例：

```json
{
  "id": "EXISTING_PRIVACY_AGREEMENT_ID",
  "defaultLang": "zh",
  "agreementInfos": [
    {
      "lang": "zh",
      "baseInfo": {
        "name": "示例应用隐私政策",
        "appName": "示例应用",
        "companyName": "实际运营主体名称"
      }
    }
  ]
}
```

这个示例只展示基础信息结构，不能作为完整隐私政策提交材料。`personalInfo`、设备权限、第三方 SDK、共享信息和扩展信息应依据应用实际行为及官方模型填写。部分更新是否保留未传入内容，需要在专用资源中读回验证，不能假定属于 JSON Merge Patch。

```sh
agc publishing privacy-agreement-update --invoke --dry-run=false --profile PROFILE \
  --header appId=APP_ID --body privacy-update.json --pretty
agc publishing privacy-agreement-submit --invoke --dry-run=false --profile PROFILE \
  --header appId=APP_ID --field id=EXISTING_PRIVACY_AGREEMENT_ID \
  --field requireConfirm=0 --pretty
```

每个语言条目的 `lang` 必填；提供 `baseInfo` 时 `name` 必填。名称最多 100 字符，应用名称 255 字符，公司名称 1024 字符，简介 500 字符。这里的协议更新不等于已经完成后台的“隐私标签信息录入”。

若隐私协议被在架版本关联，修改内容需要提交运营审核，审核通过后才能生效。审核中的协议，或被审核中应用关联的协议，不支持编辑。查询响应中的子协议信息和审核结果有助于判断当前状态。[SubAgreement](https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-data-subagreement-0000002328805241)

### 填写 SDK、共享、权限和存储信息

这些字段直接放在 `agreementInfos` 的语言条目中，不增加 `PersonalInfo` 或 `SdkInfo` 等类型名称的包装层。当前 CLI 校验以下 28 个关联模型的字段类型、明确的字符/数组限制，以及敏感数据枚举：

| 语言条目字段 | 用途和结构 | 主要限制 |
| --- | --- | --- |
| `personalInfo` | 通过 `requiredInfo`、`recommendInfo`、`adsInfo` 或相应服务模式描述收集数据的场景 | `sceneItems` 和场景内 `dataItems` 各最多 200 项；`sensitiveData` 为整数 0 或 1 |
| `devicePermissions` | 权限组、用途和组内 `permissions` | 权限组最多 200 项；提供组内权限数组时须为 1–200 项 |
| `thirdSdkInfo.sdkInfos` | SDK 名称、提供方、声明及 `privacyLabels` 收集场景 | SDK 最多 50 项；每个 SDK 的场景最多 200 项 |
| `thirdSharedInfos` | 接收方、共享目的、官网/隐私链接及数据种类 | 最多 200 项；链接最多 1024 字符 |
| `extInfo` | 联系方式、存储、用户权利入口、自定义章节、儿童隐私和版本说明 | 联系邮箱/电话各最多 10 项；电话每项最多 31 字符；额外自定义章节最多 5 项 |

以下片段展示 SDK 和存储字段的类型与嵌套。将它放到前述更新请求的语言条目中，并依据应用真实行为填写；数据项和场景 ID 应使用该协议模板/后台中的真实值。

```json
{
  "lang": "zh",
  "thirdSdkInfo": {
    "sdkInfos": [
      {
        "name": "实际 SDK 名称",
        "providerName": "实际提供方",
        "privacyDeclare": "https://example.com/sdk/privacy",
        "privacyLabels": [
          {
            "sceneReason": "实际使用目的",
            "dataItems": [
              {"itemId": "REAL_DATA_ITEM_ID", "itemDesc": "实际数据种类", "sensitiveData": 0}
            ]
          }
        ]
      }
    ]
  },
  "extInfo": {
    "storageInfo": {"country": "中国境内", "isCrossBorder": false},
    "contactInfo": {"phoneNum": ["0086-13700000000"]}
  }
}
```

`isCrossBorder` 必须是 JSON 布尔值，不能写成字符串 `"false"`。模型表标为可选的对象可以省略；提供对象时才检查其字段，提供有下限的数组时才检查下限。CLI 不会因为存在这些字段而默认要求所有隐私章节都填写，也不根据 `Integer(32)` 编造数值范围。

完整官方请求结构对应两份离线示例：[HarmonyOS 应用](../pkg/agcapi/testdata/privacy-harmonyos-request.json)、[元服务](../pkg/agcapi/testdata/privacy-atomic-request.json)。示例的协议 ID 和联系信息已替换为测试值，不能直接用于真实应用。两份结构均由独立回归测试验证，官网未在参数表定义的示例扩展字段仍保持兼容。

字段错误会返回完整路径，例如 `body.agreementInfos[0].thirdSdkInfo.sdkInfos[0].privacyLabels[0].dataItems[0].sensitiveData`，同时提供 `sourceUrl`、修正建议和 `documentationCommand`。这些检查在鉴权与 HTTP 请求之前执行。

规则来源：[AgreementInfo（隐私政策协议）](https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-data-agreementinfo-privacy-0000002294845630)、[AgreementBaseInfo（隐私政策协议）](https://developer.huawei.com/consumer/cn/doc/doccenter-submission/publish-api-agreementbaseinfo-privacy-0000002328924929)。

## 元服务资质审核：查询要求、上传材料、创建与修改

先查询标签依赖的资质，不猜测标签或资质编码：

```sh
agc publishing tag-qual-dependent --invoke --dry-run=false --profile PROFILE --pretty
```

随后按上传流程取得真实 `objectId`，并保留本地文件名称、字节数和 SHA256。`qualificationFileInfos` 是文件描述对象数组；这个接口不接收本地文件二进制，也不会因一个 objectId 就确认文件内容正确。

创建请求 `qualification.json` 的结构：

```json
{
  "tagQualInfos": [
    {
      "tagId": "REAL_TAG_ID",
      "qualificationInfos": [
        {
          "qualificationCode": "REAL_QUALIFICATION_CODE",
          "qualificationFileInfos": [
            {
              "objectId": "REAL_UPLOADED_OBJECT_ID",
              "fileName": "实际资质文件.pdf",
              "fileSize": 12345,
              "sha256": "REAL_FILE_SHA256"
            }
          ],
          "validityType": 0,
          "startTime": "2030-01-01",
          "expireTime": "2031-01-01"
        }
      ]
    }
  ]
}
```

日期、文件大小和 SHA256 均为格式示意，实际调用必须填写资质文件上的真实日期及真实文件数据。`validityType=0` 表示自定义有效期，日期使用 `yyyy-MM-dd`，提供截止日期时剩余有效期须大于 30 天；永久有效使用 `validityType=1`。日期在模型表中可选，CLI 不会将永久有效分支改成强制提供日期。

```sh
agc publishing qualification-review-create --invoke --dry-run=false --profile PROFILE \
  --header appId=APP_ID --body qualification.json --pretty
agc publishing qualification-review-query --invoke --dry-run=false --profile PROFILE \
  --header appId=APP_ID --query tagIds=REAL_TAG_ID --pretty
```

修改现有任务时，在对应的 `tagQualInfos` 条目中加入真实 `taskId`，使用 `qualification-review-update`，再查询确认结果。创建时 `taskId` 可选，修改时必填。每次最多 5 个标签，每个标签最多 10 项资质，每项资质最多 10 个文件。文件描述中 `objectId`、`fileName`、`fileSize`、`sha256` 都必填。

`qualification-review-cancel` 使用请求体 `taskId` 撤销任务；`qualification-review-delete` 使用请求体 `tagIds` 字符串数组移除旧资质任务。两者用途和输入不同。官方说明：资质已通过但上架因资质被驳回时，需要先移除旧材料再重新提交；不应把这一步用于无关的业务资源。

规则来源：[TagQualInfo](https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-data-tagqualinfo-0000002236041478)、[QualificationInfo](https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-data-qualificationinfo-0000002236201286)、[QualificationFileInfo](https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-data-qualificationfileinfo-0000002271160609)。

## 提交版本已解决的问题工单

先用 `app-version-voc` 查询未处理工单；官方要求应用软件包已上传。将已解决的问题作为 `app-info-update-v3` 请求中的 `versionVocs` 数组提交，每个条目都需要 `versionId`、`packageId`、`vocId`，每个 ID 最多 32 字符。上报已解决不代表工单立即关闭，工单平台验证通过后才会关闭。

```sh
agc publishing app-version-voc --invoke --dry-run=false --profile PROFILE \
  --header appId=APP_ID --field pageNo=1 --field pageSize=20 --pretty
```

调用基本资料更新时，还需按该接口的要求提供 `publishCountry` 和 `encrypted`。不要根据工单 ID 猜测对应的版本或软件包 ID。[VersionVoc 模型](https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-data-versionvoc-0000002495156924)

## 核验边界与复跑

年龄分级页面明确要求使用 API 客户端认证，并先设置应用分类、完成官方问卷取得 `feedback`。`expectedRating` 是评级枚举，不能直接填写年龄。提交发布要求资料完整、软件包上传完成；官方说明软件包异步解析，上传后等待 2 分钟再提交。成功提交仍需查询审核和上架状态。

资质审核系列用于元服务。仅选择智能手表或运动手表的应用不支持提交资质审核任务。下架接口仅向允许清单内用户开放，普通 Publishing API 开发者不能仅凭命令已登记就认为具备调用权限。上述说明的来源与完整参数见逐项参考。

隐私模型 `SdkInfo.privacyLabels` 描述 SDK 收集的信息；这证实隐私政策接口包含相关数据结构，但仍不能证明更新这些字段会完成后台独立的“隐私标签信息录入”步骤。需要以专用应用的实际后台读回结果确认这种关联。

顶层对照、模型捕获、规则实现、离线测试和真实业务验证分别记录。31 项方法/路径及顶层参数对照通过，不意味着所有模型和业务状态均已验证。当前仍在逐页复核模型、错误码和附录；不得将本文件作为全部接口验收完成的证明。

```sh
python3 scripts/audit-publishing-contracts.py
go test ./pkg/domain ./pkg/agcapi -run 'TestPublishingAgreement|TestFileEndpointIncludes' -count=1
go test ./pkg/domain ./pkg/agcapi ./cmd/agc/command ./pkg/server -run 'TestPrivacy|TestPublishingPrivacy' -count=1
python3 scripts/audit-publishing-models.py
go test ./pkg/domain ./pkg/agcapi ./cmd/agc/command -run 'TestPublishingQualification|TestQualificationOpenAPI|TestCurrentPublishingRequestsMatchOfficialReference' -count=1
node --test scripts/test_browser_capture.mjs
make ci
go test -race ./...
```

嵌套错误应提供数组索引、具体子字段、官方来源与修正建议。协议校验独立于 HTTP 和用户凭据；离线回归测试不调用华为。线上写操作需要专用测试应用，执行后查询确认结果；发布、下架和既有业务资源变更不能用作普通冒烟测试。
