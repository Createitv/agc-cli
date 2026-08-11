# agc-cli

[![Website](https://img.shields.io/badge/website-createitv.github.io%2Fagc--cli-f04444)](https://createitv.github.io/agc-cli/)
[![Release](https://img.shields.io/github/v/release/Createitv/agc-cli?label=release)](https://github.com/Createitv/agc-cli/releases)
[![Coverage](https://img.shields.io/badge/coverage-80%25%20gate-brightgreen)](#testing--coverage)
[![License](https://img.shields.io/badge/license-MIT-green)](LICENSE)

AppGallery Connect Command Center.

`agc-cli` is a Go CLI for Huawei AppGallery Connect. It lets you discover, dry-run, and invoke AppGallery Connect APIs from the terminal, CI, local REST API, or a web command center. It outputs structured JSON by default so humans, scripts, and AI agents can follow the same release workflow.

- [中文文档](#中文)
- [English](#english)

## 中文

### 快速开始

```bash
brew install agccli

agc auth login \
  --service-account-file ~/.agc/service-account.json \
  --name production

agc auth check

agc init \
  --app-id <app-id> \
  --project-id <project-id> \
  --package-name com.example.app \
  --default-profile production

agc publishing app-info-query \
  --invoke \
  --query appId=<app-id> \
  --query lang=zh-CN
```

`agc init` 会把项目上下文写入 `.agc/project.json`。之后同一个仓库里的命令可以自动解析 app、project、package 和默认 profile，不需要每次重复传入。

### 功能概览

| 分类 | 能做什么 |
| --- | --- |
| Publishing API | 查询和更新应用信息、软件包信息、语言描述、GMS 属性、提交发布、撤回或更新上架时间 |
| Upload Management API | 处理应用包、图标、截图、视频、PDF、OBB 等上传入口和上传回调 |
| Provisioning API | 管理 HarmonyOS 证书、Profile、ACL 权限、测试设备和指纹 |
| Domain Management API | 查询、预校验、下载和更新元服务域名配置 |
| Testing API | 管理测试版本、测试包、测试用户、用户组、邀请码、反馈和公开测试链接 |
| Reports API | 请求和下载 AppGallery Connect 报表，支持 CSV、Excel 等返回形式 |
| Project Management API | 查询团队、项目、应用摘要、SDK 配置文件、服务开关和证书指纹 |
| Comments API | 获取评论、评分、评论详情，并回复或删除回复 |
| PMS API | 管理商品、订阅、促销、价格、语言展示和审核资料 |
| 在玩服务 | 处理游戏联运资源同步和 AppGallery 游戏回调 |
| 游戏道具商城 | 处理角色查询和订单回调 |
| 资源包预下载 | 管理资源包版本、文件上传、确认上传和发布 |
| CI/CD 平台 | 通过同一个命令面桥接本地 Hvigor 构建 |
| AI Agents | JSON envelope、`_links` 和 `affordances` 帮助 agent 发现下一步命令 |
| Web Command Center | 浏览器界面读取同一份接口注册表，可接入本地 REST 数据 |

当前注册 `156` 个接口条目：`153` 个华为官方 AppGallery Connect Reference 接口/回调、`2` 个上传 URL handoff 操作、`1` 个本地 Hvigor bridge。

### 环境要求

- Go 1.22 或更高版本
- Node.js 20 或更高版本，仅在开发 Web Command Center 时需要
- AppGallery Connect Service Account JSON，或 API Client ID/Key
- 对写入类 API 拥有相应 AppGallery Connect 权限

### 安装

#### Homebrew 推荐

```bash
brew install agccli
agc version
```

如果本机 Homebrew 还没有索引 Createitv tap，请先运行一次 `brew tap createitv/tap`，然后重新执行 `brew install agccli`。

#### Scoop 推荐 Windows

```powershell
scoop bucket add createitv https://github.com/Createitv/scoop-bucket
scoop install agc-cli
agc version
```

#### Winget

Winget manifest 已提交给 Microsoft 审核：

https://github.com/microsoft/winget-pkgs/pull/415361

审核通过后可使用：

```powershell
winget install --id Createitv.AgcCli -e
```

#### GitHub Release 包

发布页包含 macOS、Linux、Windows 的二进制包和 `checksums.txt`：

https://github.com/Createitv/agc-cli/releases/tag/v0.1.0

#### Go install

macOS / Linux：

```bash
mkdir -p "$HOME/.local/bin" && \
GOBIN="$HOME/.local/bin" go install github.com/Createitv/agc-cli/cmd/agc@latest && \
export PATH="$HOME/.local/bin:$PATH" && \
agc version
```

Windows PowerShell：

```powershell
$p="$env:LOCALAPPDATA\Programs\agc\bin"; New-Item -ItemType Directory -Force $p; $env:GOBIN=$p; go install github.com/Createitv/agc-cli/cmd/agc@latest; $env:Path="$p;$env:Path"; agc version
```

注意：`go install` 不经过 GoReleaser 注入版本信息，`agc version` 可能显示 `dev`。正式 Release 包和 Homebrew 安装会显示发布版本。

#### 从源码构建

```bash
git clone https://github.com/Createitv/agc-cli.git
cd agc-cli
make build
install -m 0755 bin/agc /usr/local/bin/agc
agc version
```

### 鉴权

#### 持久化登录 推荐

Service Account：

```bash
agc auth login \
  --service-account-file ~/.agc/service-account.json \
  --name production
```

API Client：

```bash
agc auth login \
  --client-id <client-id> \
  --client-key <client-key> \
  --name production
```

常用命令：

```bash
agc auth list
agc auth check
agc auth token
```

凭据保存在 `~/.agc/credentials.json`。项目文件 `.agc/project.json` 只保存 app/project/package/profile 上下文，不保存 client key 或私钥。

#### Profile 解析顺序

```text
--profile 显式参数
  -> .agc/project.json 里的默认 profile
  -> ~/.agc/credentials.json 里的激活或默认账号
```

### 项目初始化

```bash
agc init \
  --app-id <app-id> \
  --project-id <project-id> \
  --package-name com.example.app \
  --default-profile production
```

示例 `.agc/project.json`：

```json
{
  "appId": "123456789",
  "projectId": "987654321",
  "packageName": "com.example.app",
  "profile": "production"
}
```

### 命令参考

#### 全局发现

```bash
agc capabilities --output table
agc endpoints --pretty
agc openapi --pretty
agc docs publishing
```

#### API 家族

| API 家族 | 命令 | 接口数 |
| --- | --- | ---: |
| Publishing API | `agc publishing` | 14 |
| Upload Management API | `agc upload` | 6 |
| Provisioning API | `agc provisioning` | 17 |
| Domain Management API | `agc domains` | 5 |
| Testing API | `agc testing` | 27 |
| Reports API | `agc reports` | 12 |
| Project Management API | `agc projects` | 8 |
| Comments API | `agc comments` | 8 |
| PMS API | `agc pms` | 40 |
| 在玩服务 | `agc gameplay` | 8 |
| 游戏道具商城 | `agc game-items` | 2 |
| 资源包预下载 | `agc resources` | 8 |
| CI/CD 平台 | `agc cicd` | 1 |

#### Publishing 示例

```bash
agc publishing endpoints --output table

agc publishing app-info-query \
  --invoke \
  --query appId=<app-id> \
  --query lang=zh-CN

agc publishing add-packageurl \
  --invoke \
  --field appId=<app-id> \
  --field packageUrl=https://example.com/app.app

agc publishing app-submit \
  --invoke \
  --field appId=<app-id> \
  --field releaseType=1 \
  --dry-run=false
```

#### Reports 示例

```bash
agc reports endpoints --output table

agc reports appdownloadexport \
  --invoke \
  --param appId=<app-id> \
  --query from=2026-08-01 \
  --query to=2026-08-11 \
  --out downloads.csv \
  --dry-run=false
```

#### 回调 URL 示例

部分官方文档是开发者服务接收 AppGallery 回调的接口。`agc` 把它们注册为 `inbound-callback`，调用时需要传入 `callbackUrl`：

```bash
agc game-items propapi-order \
  --invoke \
  --param callbackUrl=https://callback.example.com/order \
  --dry-run
```

### dry-run 与真实调用

接口命令默认 `--dry-run=true`。默认情况下，CLI 只构建请求并输出 method、URL、headers、body，不会发给 AppGallery Connect。

```bash
agc publishing app-info-query \
  --invoke \
  --query appId=123 \
  --query lang=zh-CN \
  --dry-run
```

确认无误后再显式发送：

```bash
agc publishing app-info-query \
  --invoke \
  --query appId=123 \
  --query lang=zh-CN \
  --dry-run=false
```

通用参数：

```bash
--param key=value   # path 参数
--query key=value   # query 参数
--header key=value  # HTTP header
--field key=value   # JSON body 字段
--body body.json    # 原始 JSON body
--token <token>     # Bearer token；默认读取 AGC_ACCESS_TOKEN 或当前 profile
--out file          # 保存原始响应 body
```

### JSON 输出与 Agent 支持

`agc` 默认输出 JSON envelope，并在响应里保留 `_links` 和 `affordances`。脚本或 agent 可以先发现能力，再选择下一条合法命令，而不需要提前硬编码完整命令树。

```bash
agc capabilities --pretty
agc publishing app-info-query --pretty
```

每个接口条目包含：

- `id`
- `familyId`
- `method`
- `path`
- `command`
- `sourceUrl`
- `officialSlug`
- `direction`
- `_links`
- `affordances`

### 本地 REST 与 OpenAPI

启动本地服务：

```bash
agc web-server --addr :8421
```

常用路由：

```bash
curl http://localhost:8421/api/v1/capabilities
curl http://localhost:8421/api/v1/endpoints
curl http://localhost:8421/api/v1/openapi.json
```

每个接口都有 REST 调用路由：

```bash
curl -X POST http://localhost:8421/api/v1/publishing/endpoints/app-info-query/invoke \
  -H 'Content-Type: application/json' \
  -d '{
    "baseUrl": "https://connect-api.cloud.huawei.com",
    "query": {"appId": "123", "lang": "zh-CN"},
    "headers": {"client_id": "<client-id>"},
    "dryRun": true
  }'
```

REST 调用同样默认 dry-run。只有传入 `"dryRun": false` 并提供 token 时才会发送真实请求。

### Web Command Center

线上文档站：

https://createitv.github.io/agc-cli/

本地开发：

```bash
npm --prefix apps/web install
npm --prefix apps/web run dev
```

如果同时启动 `agc web-server --addr :8421`，Web Command Center 会读取本地 REST 数据；否则使用内置参考数据。

### CI/CD 与发布

GitHub Actions：

- `CI`：`go vet`、Go tests、覆盖率检查、Web tests、Web build
- `Pages`：部署 Web Command Center 到 GitHub Pages
- `Release`：tag `v*` 触发 GoReleaser

GoReleaser 发布内容：

- GitHub Release 二进制包和 `checksums.txt`
- Homebrew formula：`createitv/tap/agccli`
- Scoop manifest：`createitv/scoop-bucket/agc-cli`
- Winget manifest PR：`Createitv.AgcCli`
- GHCR 镜像：`ghcr.io/createitv/agc-cli`

本地检查：

```bash
make ci
TAP_GITHUB_TOKEN=dummy make release-snapshot
make release-snapshot-docker # 需要 Docker
```

### 测试与覆盖率

```bash
make test
make coverage
make coverage-check
make ci
```

CI 覆盖率门槛是 `80%`。当前测试覆盖 CLI 命令、接口注册表、dry-run 请求构建、REST 路由、OpenAPI 导出、Web 渲染和发布配置 smoke check。

### 安全边界

- 不保存华为账号密码、浏览器 cookie 或验证码会话
- Service Account 使用 `key_id`、`private_key`、`sub_account` 生成官方 PS256 JWT bearer token
- API Client 使用 `/api/oauth2/v1/token`
- `.agc/project.json` 可以提交到项目仓库；密钥应放在 `~/.agc/credentials.json` 或 CI secret store
- 默认 dry-run，写入类 API 必须显式 `--dry-run=false`
- 接口字段、权限、业务前置条件仍应以华为官方文档为准

## English

### Quick Start

```bash
brew install agccli

agc auth login \
  --service-account-file ~/.agc/service-account.json \
  --name production

agc auth check

agc init \
  --app-id <app-id> \
  --project-id <project-id> \
  --package-name com.example.app \
  --default-profile production

agc publishing app-info-query \
  --invoke \
  --query appId=<app-id> \
  --query lang=en-US
```

`agc init` writes project context to `.agc/project.json`. Later commands in the same repository can resolve the app, project, package name, and default profile automatically.

### Features

| Category | What you can do |
| --- | --- |
| Publishing API | Query and update app information, package state, localized descriptions, GMS settings, submissions, and release timing |
| Upload Management API | Handle package, icon, screenshot, video, PDF, and OBB upload entry points |
| Provisioning API | Manage HarmonyOS certificates, profiles, ACL permissions, test devices, and fingerprints |
| Domain Management API | Query, pre-check, download, and update atomic service domain configuration |
| Testing API | Manage test versions, packages, testers, groups, invitation codes, feedback, and public test links |
| Reports API | Request and download AppGallery Connect reports as CSV or Excel files |
| Project Management API | Query teams, projects, app summaries, SDK configs, services, and certificate fingerprints |
| Comments API | Fetch review lists, ratings, details, and create or delete replies |
| PMS API | Manage products, subscriptions, promotions, prices, display languages, and review assets |
| Game Playing Service | Handle game resource synchronization and AppGallery game callbacks |
| Game Item Mall | Handle role query and order callbacks |
| Resource Package Predownload | Manage resource package versions, file upload, upload confirmation, and release |
| CI/CD Platform | Bridge local Hvigor builds into the same command surface |
| AI Agents | JSON envelopes, `_links`, and `affordances` help agents choose the next legal command |
| Web Command Center | Browser UI backed by the same endpoint registry and optional local REST data |

The project currently registers `156` interface entries: `153` Huawei official AppGallery Connect Reference endpoints/callbacks, `2` upload URL handoff operations, and `1` local Hvigor bridge.

### Requirements

- Go 1.22 or later
- Node.js 20 or later, only for Web Command Center development
- AppGallery Connect Service Account JSON, or API Client ID/Key
- AppGallery Connect permissions for the APIs you invoke

### Installation

#### Homebrew recommended

```bash
brew install agccli
agc version
```

If Homebrew has not indexed the Createitv tap locally yet, run `brew tap createitv/tap` once and then run `brew install agccli` again.

#### Scoop recommended for Windows

```powershell
scoop bucket add createitv https://github.com/Createitv/scoop-bucket
scoop install agc-cli
agc version
```

#### Winget

The Winget manifest has been submitted for Microsoft review:

https://github.com/microsoft/winget-pkgs/pull/415361

After approval, use:

```powershell
winget install --id Createitv.AgcCli -e
```

#### GitHub Release archives

Release assets include macOS, Linux, Windows binaries and `checksums.txt`:

https://github.com/Createitv/agc-cli/releases/tag/v0.1.0

#### Go install

macOS / Linux:

```bash
mkdir -p "$HOME/.local/bin" && \
GOBIN="$HOME/.local/bin" go install github.com/Createitv/agc-cli/cmd/agc@latest && \
export PATH="$HOME/.local/bin:$PATH" && \
agc version
```

Windows PowerShell:

```powershell
$p="$env:LOCALAPPDATA\Programs\agc\bin"; New-Item -ItemType Directory -Force $p; $env:GOBIN=$p; go install github.com/Createitv/agc-cli/cmd/agc@latest; $env:Path="$p;$env:Path"; agc version
```

Note: `go install` does not run through GoReleaser ldflags, so `agc version` may show `dev`. Release archives and Homebrew builds contain the release version.

#### Build from source

```bash
git clone https://github.com/Createitv/agc-cli.git
cd agc-cli
make build
install -m 0755 bin/agc /usr/local/bin/agc
agc version
```

### Authentication

#### Persistent login recommended

Service Account:

```bash
agc auth login \
  --service-account-file ~/.agc/service-account.json \
  --name production
```

API Client:

```bash
agc auth login \
  --client-id <client-id> \
  --client-key <client-key> \
  --name production
```

Common commands:

```bash
agc auth list
agc auth check
agc auth token
```

Credentials are saved to `~/.agc/credentials.json`. Project context in `.agc/project.json` does not store client keys or private keys.

#### Profile resolution

```text
--profile explicit override
  -> default profile in .agc/project.json
  -> active/default account in ~/.agc/credentials.json
```

### Project Init

```bash
agc init \
  --app-id <app-id> \
  --project-id <project-id> \
  --package-name com.example.app \
  --default-profile production
```

Example `.agc/project.json`:

```json
{
  "appId": "123456789",
  "projectId": "987654321",
  "packageName": "com.example.app",
  "profile": "production"
}
```

### Command Reference

#### Discovery

```bash
agc capabilities --output table
agc endpoints --pretty
agc openapi --pretty
agc docs publishing
```

#### API families

| API family | Command | Interfaces |
| --- | --- | ---: |
| Publishing API | `agc publishing` | 14 |
| Upload Management API | `agc upload` | 6 |
| Provisioning API | `agc provisioning` | 17 |
| Domain Management API | `agc domains` | 5 |
| Testing API | `agc testing` | 27 |
| Reports API | `agc reports` | 12 |
| Project Management API | `agc projects` | 8 |
| Comments API | `agc comments` | 8 |
| PMS API | `agc pms` | 40 |
| Game Playing Service | `agc gameplay` | 8 |
| Game Item Mall | `agc game-items` | 2 |
| Resource Package Predownload | `agc resources` | 8 |
| CI/CD Platform | `agc cicd` | 1 |

#### Publishing examples

```bash
agc publishing endpoints --output table

agc publishing app-info-query \
  --invoke \
  --query appId=<app-id> \
  --query lang=en-US

agc publishing add-packageurl \
  --invoke \
  --field appId=<app-id> \
  --field packageUrl=https://example.com/app.app

agc publishing app-submit \
  --invoke \
  --field appId=<app-id> \
  --field releaseType=1 \
  --dry-run=false
```

#### Reports example

```bash
agc reports endpoints --output table

agc reports appdownloadexport \
  --invoke \
  --param appId=<app-id> \
  --query from=2026-08-01 \
  --query to=2026-08-11 \
  --out downloads.csv \
  --dry-run=false
```

#### Callback URL example

Some official reference entries describe callbacks implemented by the developer service. `agc` registers them as `inbound-callback`, and you provide the callback URL explicitly:

```bash
agc game-items propapi-order \
  --invoke \
  --param callbackUrl=https://callback.example.com/order \
  --dry-run
```

### dry-run and real requests

Endpoint commands default to `--dry-run=true`. The CLI builds and prints the request but does not send it to AppGallery Connect.

```bash
agc publishing app-info-query \
  --invoke \
  --query appId=123 \
  --query lang=en-US \
  --dry-run
```

Send the real request only after inspection:

```bash
agc publishing app-info-query \
  --invoke \
  --query appId=123 \
  --query lang=en-US \
  --dry-run=false
```

Common endpoint flags:

```bash
--param key=value   # path parameter
--query key=value   # query parameter
--header key=value  # HTTP header
--field key=value   # JSON body field
--body body.json    # raw JSON body
--token <token>     # bearer token; defaults to AGC_ACCESS_TOKEN or active profile
--out file          # write raw response body
```

### JSON Output and Agent Support

`agc` returns JSON envelopes by default and includes `_links` and `affordances`. Scripts and agents can discover capabilities first, then choose the next legal command without hard-coding the command tree.

```bash
agc capabilities --pretty
agc publishing app-info-query --pretty
```

Each endpoint entry includes:

- `id`
- `familyId`
- `method`
- `path`
- `command`
- `sourceUrl`
- `officialSlug`
- `direction`
- `_links`
- `affordances`

### Local REST and OpenAPI

Start the local server:

```bash
agc web-server --addr :8421
```

Routes:

```bash
curl http://localhost:8421/api/v1/capabilities
curl http://localhost:8421/api/v1/endpoints
curl http://localhost:8421/api/v1/openapi.json
```

Every registered endpoint also has an invocation route:

```bash
curl -X POST http://localhost:8421/api/v1/publishing/endpoints/app-info-query/invoke \
  -H 'Content-Type: application/json' \
  -d '{
    "baseUrl": "https://connect-api.cloud.huawei.com",
    "query": {"appId": "123", "lang": "en-US"},
    "headers": {"client_id": "<client-id>"},
    "dryRun": true
  }'
```

REST invocation defaults to dry-run. Set `"dryRun": false` and provide a token only when you want the local server to send the request.

### Web Command Center

Production docs:

https://createitv.github.io/agc-cli/

Local development:

```bash
npm --prefix apps/web install
npm --prefix apps/web run dev
```

Start `agc web-server --addr :8421` for live local REST data. Without the server, the web app remains useful as a static reference.

### CI/CD and Release

GitHub Actions:

- `CI`: `go vet`, Go tests, coverage gate, Web tests, Web build
- `Pages`: deploys the Web Command Center to GitHub Pages
- `Release`: tags named `v*` trigger GoReleaser

GoReleaser publishes:

- GitHub Release archives and `checksums.txt`
- Homebrew formula: `createitv/tap/agccli`
- Scoop manifest: `createitv/scoop-bucket/agc-cli`
- Winget manifest PR: `Createitv.AgcCli`
- GHCR image: `ghcr.io/createitv/agc-cli`

Local checks:

```bash
make ci
TAP_GITHUB_TOKEN=dummy make release-snapshot
make release-snapshot-docker # requires Docker
```

### Testing & Coverage

```bash
make test
make coverage
make coverage-check
make ci
```

The CI gate requires at least `80%` Go test coverage. Current tests cover CLI commands, endpoint registry behavior, dry-run request building, REST routes, OpenAPI export, Web rendering, and release configuration smoke checks.

### Security

- Do not store Huawei account passwords, browser cookies, or captcha-backed sessions in this CLI
- Service Account credentials generate the official PS256 JWT bearer token from `key_id`, `private_key`, and `sub_account`
- API Client credentials call `/api/oauth2/v1/token`
- `.agc/project.json` is project context only; put secrets in `~/.agc/credentials.json` or CI secret storage
- dry-run is the default; mutating requests require explicit `--dry-run=false`
- Always verify request fields, permissions, and business prerequisites against Huawei official documentation before production writes

## License

MIT
