# 云开发管理扩展

本轮新增 47 项操作，并保留了 23 页通过 Computer Use 阅读的官方页面、嵌套规则、段落链接及审查日期。管理传输来自公开的 `@hw-agconnect/agc-cli@1.12.0`，每份 `--docs` 都包含对应源码、SHA256 和相关指南。这是源码适配，不是华为承诺稳定的公开 REST 规范。

## 命令与鉴权

```sh
./bin/agc cache endpoints --pretty
./bin/agc functions create-alias --docs --pretty
./bin/agc clouddb create-object-type --docs --pretty
./bin/agc templates deploy --docs --pretty
./bin/agc prefetch fetch-periodic --docs --pretty
```

管理入口使用 `--session-file /PRIVATE_PATH/session.json`，REST 请求根字段为 `sessionFile`，可选 `projectId` 必须匹配会话。`--profile` 的发布凭据不会替代用户会话。默认 dry-run 不读取会话、不访问华为，也不执行设备命令。实际执行需要 `--invoke --dry-run=false`。

会话文件保存在仓库外，建议权限 `0600`。通用字段为 `uid`、`teamId`、`productId`、`accessToken`、RFC3339 `expiresAt`。额外字段取决于入口：

| 入口 | 会话字段及传输 |
| --- | --- |
| 云函数管理 | `projectBaseUrl`；用户 OAuth 头 |
| schema、存储区、云缓存 | `projectEdgeBaseUrl`、`authInfo`、`csrfToken`；同时保留 OAuth、Cookie 和 CSRF 头，与官方源码一致 |
| SAM 应用与计划 | `teamBaseUrl`；用户 OAuth 头 |
| SAM 部署代码上传 | `projectStorageBaseUrl`；用户 OAuth 头，原始 ZIP 文件流 |

不读取官方 CLI 的加密 scratch，不从浏览器提取凭据，不推测 Cookie 与 OAuth 的对应关系。TLS 验证始终保留，认证重定向不能跨源。响应中的凭据与管理环境变量会脱敏。

## 云缓存

六项操作：`execute-command`、`connection-info`、`usage`、`enable-service`、`set-password`、`reset-password`。

```sh
./bin/agc cache execute-command --invoke \
  --session-file /PRIVATE_PATH/session.json --field 'command=GET dedicated:key'
./bin/agc cache usage --invoke --session-file /PRIVATE_PATH/session.json \
  --query views=query_cache_metric_trend \
  --query startTime=1792962000000 --query endTime=1792965600000
```

Web CLI 支持的命令名单、ASCII/空格与 4096 字符限制来自官方页面。Redis 命令参数的完整语法仍由服务端判断。时间戳为毫秒，区间不超过三天。密码须以字母开头、8–32 位，包含大小写字母、数字和文档列出的特殊字符。开通与设置密码是两个请求；密码重置可能影响已有连接，测试不能选用业务缓存。

## 数据库 schema 与存储区

十一项操作：`schema-versions`、`get-object-type`、`create-object-type`、`update-object-type`、`delete-object-type`、`import-object-types`、`export-object-types`、`list-zones`、`create-zone`、`delete-zone`、`process-status`。

```sh
./bin/agc clouddb schema-versions --invoke --session-file /PRIVATE_PATH/session.json
./bin/agc clouddb get-object-type --invoke --session-file /PRIVATE_PATH/session.json \
  --query schemaVersion=1 --query objectTypeName=Book
./bin/agc clouddb create-object-type --invoke --session-file /PRIVATE_PATH/session.json \
  --body ./object.json
```

创建使用 `{"objectType":{...}}`。更新同时需要完整 `objectType` 和查询得到的 `originObjectType`；不能通过改名实现迁移。导入使用包含 `schemaVersion`、`objectTypes`、可选 `permissions` 的原始 JSON 文档，不增加包装字段。官方 CLI 将该 JSON 字符串原样交给 Axios，最终 Content-Type 为 `application/x-www-form-urlencoded`，但字节不是表单编码；适配器保留此行为，并有独立传输断言。

嵌套检查包括名称/保留字、字段类型、1–5 个主键、自增字段、加密/敏感字段总数与互斥关系、默认值、索引字段及 16 个索引/5 个组合字段限制。上线后能否修改已有结构和权限依赖真实应用状态；CLI 不能仅从本地 JSON 证明迁移可行。返回 `processId` 仅表示开始处理，必须用 `process-status` 并重新导出/查询 schema 确认完成。

## 云函数管理

新增二十项传输：别名创建/查询/列表/更新/删除；HTTP 与其他触发器查询/列表/删除；HTTP/定时触发器创建、定时触发器更新；调用策略查询/创建；层列表、版本列表、版本查询和层版本创建。

```sh
./bin/agc functions get-alias --invoke --session-file /PRIVATE_PATH/session.json \
  --param functionName=dedicated-func --param aliasName=dedicated-alias
./bin/agc functions create-alias --invoke --session-file /PRIVATE_PATH/session.json \
  --param functionName=dedicated-func --body ./alias.json
```

别名路由为动态版本键到整数权重的对象，最多两个版本，总权重 100。CLI 要求明确提供完整权重，不自动选择“最新”版本。更新需先读别名并提供 `revisionId`。触发器的 `funcId` 必须来自函数查询并带正确限定符；创建触发器与创建调用策略是独立步骤，不会自动授予调用权限。

层上传的 `body.zipFile` 为 Base64 ZIP。当前指南上限为压缩 150MiB、解压 300MiB、30000 文件；旧 CLI 本地检查仅允许小于 50MiB。此直接 API 适配选择当前指南限制，并保留这一差异。旧 CLI 运行时枚举与当前指南仅允许新建 Node.js 20.x/latest 的说明也有差异，现代运行时的精确 wire 值及云端可用性不能由旧包源码证明。定时触发器开关同样存在“当前指南称创建后不可更改、旧 CLI 提供更新请求”的冲突，不能把离线更新请求通过当作云端支持证明。

## SAM 模板部署

`search-applications`、`create-application`、`create-plan`、`get-plan`、`apply-plan` 暴露官方 SAM 的应用与变更计划传输。`deploy` 提供新应用完整流程：先校验全部本地 ZIP，按应用名前缀转换组件名，查询同名应用，上传 ZIP，创建应用，再独立查询同一应用 ID 的 readiness。

```sh
./bin/agc templates deploy --invoke --session-file /PRIVATE_PATH/session.json \
  --body ./sam-deploy.json
```

部署文件为 SAM JSON 根对象，包含 `name`、`components` 以及本地选项 `storageInstance`、可选 `waitForReady`。组件使用 `Huawei::AGC::CloudFunctions`，`properties.codeUri` 在此工作流中是本地 ZIP；底层 `create-application`/`create-plan` 的 `codeUri` 则是已上传的 `bucket/object.zip`。代码存储实例必须事先存在，不自动开通云存储。

默认创建模式拒绝覆盖同名应用。更新需显式创建并阅读计划后再应用。默认 `waitForReady:false`；实际执行仍进行一次独立读回。`true` 最多查询 15 次，间隔 2.4 秒。结果分开返回 `applicationAccepted` 与 `ready`；上传/创建/读回失败不自动重试，已上传对象可能保留，不自动删除。

这些是 SAM CloudFunctions 应用部署，不是控制台“抽奖活动、排行榜”等模板库的一键部署/升级。浏览器审查到的模板库页面描述控制台流程，尚未提供其管理接口契约。

## 预加载

`fetch-install`、`fetch-periodic`、`fetch-deferred` 模拟华为对开发者服务器的 GET 请求，不加载/转发 AGC 凭据；安装响应最多 2MiB，另外两种最多 3MiB。周期请求的用户 `token` 最多 2048 字符、`params` 最多 1024 字符；响应中任何 `code` 都是应用数据，不能当成管理接口错误。

```sh
./bin/agc prefetch fetch-periodic --invoke --param callbackUrl=https://www.example.com/prefetchData \
  --query appId=DEDICATED_APP --query 'params={"theme":"holiday"}'
./bin/agc prefetch device-getcache --invoke \
  --field bundleName=com.example.dedicated --field deviceId=EXPLICIT_DEVICE
```

`device-getcache` 通过独立 argv 执行 `hdc -t DEVICE shell cf_prefetch getcache -m BUNDLE`，可用 `body.hdcPath` 指定已安装的 hdc。设备须满足官方 ROM、debug HAP 与调试证书条件。进程成功仅表示命令完成，`dataFreshnessVerified:false` 提醒调用方检查 `fetch data timestamp`；旧时间戳表示本次拉取失败。已保留官方错误码与修正指引。云端绑定云函数/下载地址仍只有控制台文档，未猜测管理接口。

## 验证边界

逐项 CLI/REST 本地 HTTP 测试使用独立源码夹具，同时覆盖 dry-run、业务错误、HTTP 错误、会话过期与项目隔离、嵌套输入、原始文件、回调数据和部署阶段。设备入口使用测试可执行文件确认 argv，不会接触用户手机。

真实云端写验证仍需专用项目、显式 profile/用户会话和可清理资源。未执行的项目不会标为云端通过。当前没有取得可复核 wire 契约的控制台模板库部署/升级、预加载绑定配置、层删除和其他服务类型触发器创建保留为未支持，不用猜测 URL 填补。

本次工作区二进制为 `./bin/agc`；系统安装版升级、稳定版发布和 GitHub Actions 验证是独立步骤。
