# agc CLI 使用指南

[项目首页](../README.md) · [English](CLI_USAGE.en.md)

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

当前注册 **156 个接口条目**：153 个华为官方接口/回调、2 个上传 URL handoff 操作、1 个本地 Hvigor bridge 条目。注册表示可以发现接口并构建通用请求，不表示每个接口都已通过生产验证。

| 你要做什么 | 命令入口 | 注册条目 |
| --- | --- | ---: |
| 查询/更新应用资料、多语言描述，构建提交审核请求 | `agc publishing` | 14 |
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

目前可用的是接口发现、凭据选择、请求构建、通用 HTTP 调用、本地 REST 与 OpenAPI。文件上传尚未封装为完整的二进制/multipart 上传流程；Hvigor 条目也尚未接入本地构建执行器。

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
| `--field key=value` | JSON 字段，可重复；值按字符串编码 |
| `--body request.json` | 读取完整请求体文件；优先于 `--field` |
| `--token TOKEN` | 显式提供 Bearer token |
| `--out response.json` | 将原始响应体保存到文件，不转换文件格式 |
| `--timeout 120s` | 调整请求超时，默认 60 秒 |

包含数组、对象、数字或布尔值的请求请使用 `--body`。dry-run 不展示完整 headers/body，也不校验全部华为业务字段；发送前需要自行核对请求文件、权限和官方协议。

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
