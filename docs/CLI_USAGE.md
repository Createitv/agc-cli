# agc CLI 使用指南

[项目首页](../README.md) · [English](CLI_USAGE.en.md)

## 内置 Agent Skills

```bash
agc skills list --pretty
agc skills show agc-cli --output markdown
```

这些指引内嵌于二进制，可离线读取，无需凭据；不会安装文件或查询远程技能库。默认输出 JSON，也支持 `--output table` 和 `--output markdown`。

旧版 v0.3.0 不包含 `skills` 命令。源码修复需重新构建；已安装的发布版需在包含此修复的新版本发布后升级。

## 快速开始

先准备 AppGallery Connect 中的应用 ID，以及 Service Account JSON 或 API Client 凭据。应用 ID 可从 [AppGallery Connect 控制台](https://developer.huawei.com/consumer/cn/service/josp/agc/index.html) 的应用信息中查找。

### 1. 安装

macOS（需要已安装 Homebrew）：

```bash
brew tap createitv/tap
brew install agc-cli
agc version
```

Windows 和 Linux 用户见[安装与升级](#安装与升级)。使用发布版安装包不需要 Go 或 Node.js。

### 2. 保存凭据

将你的 Service Account JSON 放在本机固定位置，然后运行：

```bash
agc auth login --service-account-file ~/.agc/service-account.json --name production
agc auth check
```

`auth check` 显示当前本地凭据配置；它不验证华为服务器是否接受该凭据。其他登录方式见[鉴权与多账号](#鉴权与多账号)。

### 3. 绑定项目

在你的应用项目目录运行，替换示例 ID：

```bash
agc init --app-id YOUR_APP_ID --default-profile production
```

配置写入 `.agc/project.json`。后续命令自动选择该项目绑定的凭据 profile；**当前仍需在接口参数中显式提供 appId、projectId 或包名**。

### 4. 查询第一个应用

先预览请求，不连接华为服务器：

```bash
agc publishing app-info-query --invoke --query appId=YOUR_APP_ID --query lang=zh-CN --pretty
```

输出中的 `data.dryRun` 为 `true`，并显示 HTTP 方法和目标 URL。确认后发送真实查询：

```bash
agc publishing app-info-query --invoke --query appId=YOUR_APP_ID --query lang=zh-CN --dry-run=false --pretty
```

若所用接口要求 `client_id` 请求头，追加 `--header client_id=YOUR_CLIENT_ID`；CLI 当前不会自动填充这个请求头。

## 功能与当前状态

当前注册 **184 个接口条目**：181 个华为官方接口/回调、2 个上传 URL handoff 操作、1 个本地 Hvigor bridge 条目。注册表示可以发现接口并构建通用请求，不表示每个接口都已通过生产验证。

| 你要做什么 | 命令入口 | 注册条目 |
| --- | --- | ---: |
| 查询/更新应用资料、多语言描述，构建提交审核请求 | `agc publishing` | 42 |
| 查看上传地址、分片上传和确认上传接口 | `agc upload` | 6 |
| 管理鸿蒙证书、Provisioning Profile、设备和指纹 | `agc provisioning` | 17 |
| 查询和更新元服务域名配置 | `agc domains` | 5 |
| 管理测试版本、测试用户、群组和反馈 | `agc testing` | 27 |
| 请求下载、销售等报表 | `agc reports` | 12 |
| 查询团队、项目、应用和 SDK 配置 | `agc projects` | 8 |
| 查询评分、评论并管理回复 | `agc comments` | 8 |
| 管理商品、订阅、价格和促销 | `agc pms` | 40 |
| 查看在玩服务、游戏道具商城回调协议 | `agc gameplay` / `agc game-items` | 8 / 2 |
| 查看资源包预下载接口 | `agc resources` | 8 |
| 查看本地 Hvigor 构建接口定义 | `agc cicd` | 1 |

目前支持接口发现、凭据选择、请求构建、通用 HTTP 调用、本地 REST 与 OpenAPI。签名上传通过 `--file` 发送二进制文件，分片上传需按初始化、获取地址、上传、合并顺序调用。本地构建可执行 `agc --timeout 180s cicd build-hvigor --invoke --dry-run=false --field project=/absolute/project/path`，成功后返回产物路径、大小与 SHA-256；本地构建不需要华为授权。`--field hvigorPath=/absolute/hvigorw/path` 可覆盖构建工具，`product` 默认 default，`buildMode` 默认 debug，也可设为 release。`--field task=assembleApp` 构建用于上传的项目 `.app` 包，默认任务 assembleHap 构建模块 `.hap`。产物是否签名取决于项目签名配置；构建成功不代表远端解析、审核或发布成功。构建输出包含敏感配置的可能性较高，因此 API 不返回原始工具日志。真实接口验证进度见 `docs/verification/goal-progress.json`。

## 按任务查找文档

| 任务 | 从这里开始 |
| --- | --- |
| 安装、登录、第一次查询 | [快速开始](#快速开始) |
| 多账号、项目配置、环境变量 | [鉴权与多账号](#鉴权与多账号) |
| 上传与提交审核 | [应用发布流程](#应用发布流程) |
| 查看所有命令、参数和官方文档地址 | [命令发现与请求参数](#命令发现与请求参数) |
| 用脚本或 AI Agent 操作 | [JSON 与 Agent 支持](#json-与-agent-支持) |
| 使用浏览器界面或本地 API | [Web、REST 与 OpenAPI](#webrest-与-openapi) |
| 解决常见问题 | [常见问题](#常见问题) |
| 开发、测试和工具自身发布 | [参与开发](#参与开发) |

[CLI 使用指南](../docs/CLI_USAGE.md) 包含更多操作示例；[实现计划](../docs/AGC_CLI_FULL_PLAN.md) 介绍架构与规划，实际能力以当前代码和本 README 的状态说明为准。

## 安装与升级

### macOS：Homebrew

```bash
brew tap createitv/tap
brew install agc-cli
# 升级已安装版本
brew update
brew upgrade agc-cli
```

### Windows：Scoop

需要已安装 Scoop，在 PowerShell 中运行：

```powershell
scoop bucket add createitv https://github.com/Createitv/scoop-bucket
scoop install agc-cli
agc version
# 升级已安装版本
scoop update agc-cli
```

Winget 发布记录见 [manifest PR](https://github.com/microsoft/winget-pkgs/pull/415361)。使用前可运行 `winget search --id Createitv.AgcCli -e` 确认可用性。

### macOS / Linux / Windows：Release 安装包

在 [Releases](https://github.com/Createitv/agc-cli/releases) 选择与你的系统和架构匹配的文件，按该版本的 `checksums.txt` 校验后解压，将 `agc`（Windows 为 `agc.exe`）放到 PATH 中，再运行 `agc version`。升级时用新版本替换旧二进制。

发布配置包含 macOS/Linux 的 amd64、arm64 和 Windows amd64。具体资产以发布页为准。

### 使用 Go 安装

需要 Go 1.22 或更高版本：

```bash
go install github.com/Createitv/agc-cli/cmd/agc@latest
```

将 Go 的二进制安装目录加入 PATH（默认通常为 `$(go env GOPATH)/bin`）。这种安装方式的 `agc version` 可能显示 `dev`，发布版安装包通过 GoReleaser 注入版本信息。

## 鉴权与多账号

Service Account 文件应包含 `key_id`、`private_key`、`sub_account`。CLI 使用该文件签署 PS256 JWT；登录时保存文件路径，因此文件需要留在本机。

也可以保存 API Client 凭据（下列变量由你预先设置）：

```bash
agc auth login --client-id "$AGC_CLIENT_ID" --client-key "$AGC_CLIENT_KEY" --name staging
agc --profile staging auth check
```

这两个变量只是 shell 传参示例；CLI 不会自动读取 `AGC_CLIENT_ID` / `AGC_CLIENT_KEY`。

凭据选择顺序：`--profile` → 项目 `.agc/project.json` 中的 `profile` → 全局 active 账号（或唯一账号）。每次登录会将该账号设为 active。临时切换账号或指定其他项目：

```bash
agc --profile staging publishing endpoints
agc --project ../another-app auth check
```

`agc init` 还接受 `--project-id` 和 `--package-name`，用于保存项目上下文。

| 环境变量 | 用途 |
| --- | --- |
| `AGC_ACCESS_TOKEN` | 接口调用的 Bearer token，优先级低于 `--token`，高于凭据 profile |
| `AGC_CREDENTIALS_PATH` | 覆盖默认凭据文件 `~/.agc/credentials.json`；低于 `--credentials-path` |

`agc auth token` 输出令牌。当前 `auth login`、`auth list`、`auth check` 在 API Client 模式下可能输出 `clientKey`；这些输出不应直接上传到 Issue 或公开 CI 日志。凭据文件、Service Account 私钥和访问令牌应保存在本机或 CI secret store 中。

## 命令发现与请求参数

无需登录即可查看接口定义：

```bash
agc --help
agc capabilities --output table
agc publishing endpoints --output table
agc publishing app-info-query --pretty
agc publishing app-info-query --help
agc endpoints --pretty
```

接口定义的 `sourceUrl` 指向对应的华为官方参考文档，例如[查询应用信息](https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-app-info-query-0000001158365045)。完整映射见 [endpoint_catalog.go](../pkg/domain/endpoint_catalog.go)。

| 调用方式 | 行为 |
| --- | --- |
| 不加 `--invoke` | 只显示接口定义 |
| `--invoke` | 构建请求，默认 dry-run；输出方法、URL 和 `dryRun` |
| `--invoke --dry-run=false` | 发送真实请求，包括 GET 查询 |

| 参数 | 用途 |
| --- | --- |
| `--param key=value` | 路径参数，可重复 |
| `--query key=value` | 查询参数，可重复 |
| `--header key=value` | HTTP 请求头，可重复 |
| `--field key=value` | JSON 字段，可重复；已声明的整数、数组和对象按相应 JSON 类型编码 |
| `--body request.json` | 读取完整请求体文件；优先于 `--field` |
| `--token TOKEN` | 显式提供 Bearer token |
| `--out response.json` | 将原始响应体保存到文件，不转换文件格式 |
| `--timeout 120s` | 调整请求超时，默认 60 秒 |

已声明类型的字段可通过 `--field` 传入整数、数组或对象；复杂嵌套请求也可使用 `--body`。dry-run 不展示完整 headers/body，也不校验全部华为业务字段；发送前需要自行核对请求文件、权限和官方协议。

## 现行 Publishing API 命令

各操作的用途、请求位置和协议操作流程见 [Publishing 使用与核验](PUBLISHING_WORKFLOWS.md)。

官网现行 Publishing API 参考的 31 项操作全部有对应命令。Publishing 共 42 个条目，包含保留的旧版命令。HarmonyOS 资料更新与提交发布使用下表的 V3 命令；原有 V2 命令保持兼容。

| 官网操作 | CLI 命令 |
| --- | --- |
| 创建APP ID | `agc publishing app-create` |
| 查询应用包名对应的appid | `agc publishing appid-list` |
| 查询团队账号下的成员角色列表 | `agc publishing members-role-list` |
| 查询应用信息 | `agc publishing app-info-query-v3` |
| 更新应用基本信息 | `agc publishing app-info-update-v3` |
| 更新应用语言信息 | `agc publishing language-info-update-v3` |
| 删除应用语言信息 | `agc publishing language-info-delete` |
| 更新应用文件信息 | `agc publishing app-file-info-v3` |
| 更新应用软件包信息 | `agc publishing app-package-info-update` |
| 更新版本上架时间 | `agc publishing update-releasetime` |
| 更新分阶段发布 | `agc publishing phased-release-update` |
| 查询应用软件包信息 | `agc publishing app-package-info-query` |
| 查询软件包编译状态 | `agc publishing package-compile-status-v3` |
| 查询应用所有版本类型的版本列表 | `agc publishing app-version-list` |
| 查询应用未处理的问题工单 | `agc publishing app-version-voc` |
| 查询需要上传的资质文件信息 | `agc publishing tag-qual-dependent` |
| 提交资质审核任务 | `agc publishing qualification-review-create` |
| 查询资质审核任务详情 | `agc publishing qualification-review-query` |
| 修改资质审核任务 | `agc publishing qualification-review-update` |
| 撤销资质审核任务 | `agc publishing qualification-review-cancel` |
| 删除资质审核任务 | `agc publishing qualification-review-delete` |
| 新建用户协议 | `agc publishing user-agreement-create` |
| 查询用户协议列表 | `agc publishing user-agreement-list` |
| 更新用户协议 | `agc publishing user-agreement-update` |
| 提交隐私政策协议 | `agc publishing privacy-agreement-submit` |
| 查询隐私政策协议列表 | `agc publishing privacy-agreement-list` |
| 更新隐私政策协议 | `agc publishing privacy-agreement-update` |
| 申请年龄分级 | `agc publishing app-age-rating` |
| 提交发布 | `agc publishing app-submit-v3` |
| 撤销审核 | `agc publishing on-shelf-cancel` |
| 下架应用 | `agc publishing app-shelf-off` |

执行 `agc publishing 命令名 --pretty` 可查看 HTTP 方法、路径、参数位置、必填性、类型和官方链接。`--query`、`--header` 接收字符串；`--field` 的整数值直接填写数字，数组与对象填写带引号的 JSON；完整嵌套结构可使用 `--body request.json`。协议列表支持 `--field name=` 空名称筛选。POST 查询类接口可直接调用；删除资质审核任务通过 JSON Body 传 `tagIds`。

```bash
agc --profile default publishing app-info-query-v3 --invoke --dry-run=false --query appId=YOUR_APP_ID
agc --profile default publishing app-version-list --invoke --dry-run=false --header appId=YOUR_APP_ID
agc --profile default publishing user-agreement-list --invoke --dry-run=false --header appId=YOUR_APP_ID --field name= --field type=2
agc publishing app-package-info-update --invoke --query appId=YOUR_APP_ID --field fileName=app-release.app --field objectId=UPLOADED_OBJECT_ID
agc publishing app-info-update-v3 --invoke --query appId=YOUR_APP_ID --body app-info.json
agc publishing app-submit-v3 --invoke --query appId=YOUR_APP_ID --body submission.json
```

最后三条示例默认预览请求，增加 `--dry-run=false` 即执行对应操作。旧版 `app-info-update`、`language-info-update`、`app-file-info`、`app-submit` 继续保留，其现行 V3 命令增加 `-v3` 后缀。

## 应用发布流程

发布通常涉及：**查询应用 → 上传文件 → 更新应用文件信息 → 完善资料 → 提交审核 → 查询远程状态**。

当前 CLI 暴露各步骤的接口，尚未提供自动串联这些步骤的一键发布命令。先查看定义：

```bash
agc upload endpoints --output table
agc publishing app-file-info --pretty
agc publishing app-info-update --pretty
agc publishing language-info-update --pretty
agc publishing app-submit --pretty
```

上传二进制文件或 multipart 内容目前需使用符合华为协议的外部上传工具；`--body` 不能替代完整文件上传流程。准备好符合对应官方文档的 JSON 文件后，可以先构建更新/提交请求：

```bash
agc publishing app-file-info --invoke --body app-file-info.json --pretty
agc publishing app-info-update --invoke --body app-info.json --pretty
agc publishing app-submit --invoke --body submission.json --pretty
```

上述命令默认不发送请求。确认上传结果、文件关联和资料完整后，再为要执行的步骤增加 `--dry-run=false`，并补齐所需请求头。

提交接口返回成功后，仍需查询 AppGallery Connect 远程状态。请求成功、审核通过和正式上架是不同阶段。

## JSON 与 Agent 支持

默认 JSON 输出方便接入 `jq`、脚本和 AI Agent；人工浏览可用 `--output table` 或 `--output markdown`。

```bash
agc capabilities --pretty
agc publishing app-info-query --pretty
```

接口定义中的 `_links` 提供本地 REST 路由，`affordances` 提供后续命令模板。Agent 可先读取定义和 `sourceUrl`，补齐参数后预览请求。模板并不表示业务前置条件已经满足。

## Web、REST 与 OpenAPI

[agccli.app](https://agccli.app/) 可浏览静态接口参考。要读取本地 REST 数据，在源码仓库启动服务和 Web 开发服务器：

```bash
agc web-server --addr :8421
# 另一个终端；需要 Node.js 20+
npm --prefix apps/web ci
npm --prefix apps/web run dev
```

打开 Vite 输出的本地地址。它会将 `/api` 代理到 `127.0.0.1:8421`；没有本地 API 时页面显示参考数据。

```bash
curl http://localhost:8421/api/v1/capabilities
curl http://localhost:8421/api/v1/endpoints
curl http://localhost:8421/api/v1/openapi.json
agc openapi --pretty
```

预览 REST 调用：

```bash
curl -X POST http://localhost:8421/api/v1/publishing/endpoints/app-info-query/invoke   -H 'Content-Type: application/json'   -d '{"query":{"appId":"YOUR_APP_ID","lang":"zh-CN"},"dryRun":true}'
```

REST 调用默认 dry-run；真实调用需设置 `"dryRun": false`，在请求 JSON 的 `token` 字段提供令牌，以及所需的 `headers`。REST 调用器当前不会自动读取 CLI 保存的凭据 profile。

## 常见问题

| 问题 | 处理方式 |
| --- | --- |
| `agc: command not found` | 检查安装目录是否在 PATH 中，重新打开终端，再运行 `agc version` |
| 只输出 URL，没有查询结果 | 检查是否包含 `--invoke --dry-run=false` |
| `credential profile ... not found` | 检查本地 profile 名称，用 `--profile` 覆盖或修改项目配置 |
| `auth check` 成功，真实请求仍失败 | 它只检查本地配置；核对密钥文件、账号权限、项目归属和接口请求头 |
| 运行 init 后仍要填写 appId | 当前只自动选择 profile，应用参数仍需显式传入 |
| 下载文件不是 CSV/Excel | `--out` 保存原始响应；部分报表接口返回下载地址，需要另外下载 |
| Web 显示 reference mode | 启动本地 API 与 Vite 开发服务器，并访问 Vite 地址 |

反馈问题时附上 `agc version`、系统/架构、脱敏后的命令、错误信息和预期结果。请在 [GitHub Issues](https://github.com/Createitv/agc-cli/issues) 提交，不要附上密钥或访问令牌。

## 参与开发

需要 Go 1.22+；Web 开发还需要 Node.js 20+。从源码构建：

```bash
git clone https://github.com/Createitv/agc-cli.git
cd agc-cli
make build
./bin/agc version
npm --prefix apps/web ci
make ci
```

`make ci` 执行 Go vet、Go 测试、80% 覆盖率门槛、Web 测试和 Web 构建。单独运行可用 `make test`、`make coverage-check`、`make web-test` 和 `make web-build`。

主要目录：`cmd/agc/command`（CLI）、`pkg/agcapi`（鉴权与 HTTP）、`pkg/domain`（接口注册表）、`pkg/server`（本地 REST）、`apps/web`（网站）。欢迎通过 [Pull Request](https://github.com/Createitv/agc-cli/pulls) 改进命令、文档和测试。

工具自身的发布配置见 [.goreleaser.yaml](../.goreleaser.yaml) 和 [Release workflow](../.github/workflows/release.yml)：`v*` tag 触发发布，配置包含 Release 安装包与校验文件、Homebrew、Scoop、Winget manifest PR 和 GHCR 镜像。Homebrew/Scoop/Winget 发布需要 `TAP_GITHUB_TOKEN`。实际发布结果以 Actions 和 Release 页面为准。

## 许可证

[MIT](../LICENSE)

## 云开发（Serverless）

新增认证、云数据库、云存储、云托管、云函数和模板回调共 55 项操作。查看 [云开发使用与验证边界](SERVERLESS.md)。默认 dry-run，云端操作需显式选择专用测试资源。
