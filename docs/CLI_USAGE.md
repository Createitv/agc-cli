# agc CLI 使用指南

`agc` 是 AppGallery Connect 的命令行控制面。它把 13 个 API 家族、156 个接口条目、项目上下文、凭据 profile、本地 REST API 和 Web Command Center 放在同一套可发现的命令结构里。

> 当前边界：接口注册表、请求构建、鉴权、dry-run、通用调用器和本地 REST 路由已经可用。正式写入生产数据前，仍应以华为 AppGallery Connect 官方接口文档核对字段、权限和业务前置条件。CLI 默认启用 dry-run，只有显式传入 `--dry-run=false` 才会发送请求。

## 1. 安装到命令行环境

环境要求：Homebrew 或 Go 1.22 及更高版本；Web Command Center 另需 Node.js 20 或更高版本。

macOS 推荐 Homebrew：

```bash
brew tap createitv/tap && brew install agc-cli && agc version
```

这条命令会添加 Createitv tap、安装发布版 formula，并验证 `agc` 已进入 PATH。

Windows 推荐 Scoop：

```powershell
scoop bucket add createitv https://github.com/Createitv/scoop-bucket
scoop install agc-cli
agc version
```

Winget manifest 已提交给 Microsoft 审核。审核通过后可使用：

```powershell
winget install --id Createitv.AgcCli -e
```

没有包管理器时，可以使用 Go install fallback。macOS / Linux：

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

注意：Go install 不经过 GoReleaser 注入版本信息，`agc version` 可能显示 `dev`。Homebrew、Scoop、GitHub Release 包会显示发布版本。

从当前仓库本地构建：

```bash
make build
install -m 0755 bin/agc /usr/local/bin/agc
agc version
```

不希望写入 `/usr/local/bin` 时，可以直接运行 `./bin/agc`，或把仓库的 `bin` 目录加入 `PATH`。

## 2. 创建凭据 profile

推荐使用 AppGallery Connect Service Account：

```bash
agc auth login \
  --service-account-file ~/.agc/service-account.json \
  --name production
```

也可以使用 API Client：

```bash
agc auth login \
  --client-id "$AGC_CLIENT_ID" \
  --client-key "$AGC_CLIENT_KEY" \
  --name legacy
```

查看 profile：

```bash
agc auth list --output table
agc auth check --pretty
agc --profile production auth check --pretty
```

凭据保存在 `~/.agc/credentials.json`，文件权限为 `0600`。不要把该文件、Service Account 私钥或访问令牌提交到 Git。

## 3. 自动绑定项目 profile

在应用仓库根目录运行一次：

```bash
agc init \
  --app-id <app-id> \
  --project-id <project-id> \
  --package-name com.example.app \
  --default-profile production
```

命令会创建 `.agc/project.json`，其中只保存应用上下文和 profile 名称，不保存密钥。

后续命令按以下顺序选择凭据：

1. 当前命令显式传入的 `--profile <name>`。
2. 当前项目 `.agc/project.json` 中的 `profile`。
3. `~/.agc/credentials.json` 中标记为 active 的 profile。

临时切换到测试环境：

```bash
agc --profile staging publishing endpoints --output table
```

指定另一个项目目录：

```bash
agc --project ../another-app auth check --pretty
```

## 4. 发现能力与接口

先列出 API 家族：

```bash
agc capabilities --output table
```

查看所有注册接口：

```bash
agc endpoints --pretty
```

查看某个家族及其接口：

```bash
agc publishing status --pretty
agc publishing endpoints --output table
agc publishing app-info-query
```

不加 `--invoke` 时，具体接口命令只打印方法、路径、参数位置和下一步动作，不发送网络请求。

## 5. 构建与调用请求

接口命令使用统一参数：

- `--param key=value`：路径参数，可重复。
- `--query key=value`：查询参数，可重复。
- `--header key=value`：HTTP Header，可重复，例如 `client_id=<id>`。
- `--field key=value`：简单 JSON body 字段，可重复。
- `--body file.json`：完整 JSON 请求体。
- `--out path`：把 CSV、Excel、PDF 或二进制响应写入文件。
- `--invoke`：进入请求构建/调用路径。
- `--dry-run=false`：真正发送请求；默认值为 `true`。
- `--token`：显式 Bearer token；未提供时依次检查 `AGC_ACCESS_TOKEN` 和选中的凭据 profile。

先 dry-run 检查 URL：

```bash
agc publishing app-info-query \
  --invoke \
  --query appId=<app-id> \
  --query lang=zh-CN \
  --pretty
```

确认后发送真实请求：

```bash
agc publishing app-info-query \
  --invoke \
  --query appId=<app-id> \
  --query lang=zh-CN \
  --dry-run=false \
  --pretty
```

使用 JSON body：

```bash
agc publishing app-info-update \
  --invoke \
  --body app-info.json
```

## 6. 修改 Profile 设备绑定

准备符合华为 Provisioning API 要求的 `devices.json`，先 dry-run：

```bash
agc provisioning provision-api-update-provision \
  --invoke \
  --body devices.json \
  --pretty
```

核对生成的 `PUT /api/publish/v3/provision` 请求后再执行：

```bash
agc provisioning provision-api-update-provision \
  --invoke \
  --body devices.json \
  --dry-run=false \
  --pretty
```

这里的“设备绑定”与“凭据 profile 选择”是两个概念：前者修改 AppGallery Connect Provisioning Profile 的设备集合，后者决定 CLI 使用哪组 Service Account/API Client 凭据。

## 7. 下载报表与原始文件

```bash
agc reports appdownloadexport \
  --invoke \
  --param appId=<app-id> \
  --query from=2026-08-01 \
  --query to=2026-08-11 \
  --out downloads.csv \
  --dry-run=false
```

使用 `--out` 后，原始响应体写入文件，终端 JSON 只保留状态、URL 和响应元数据。

## 8. 输出格式

```bash
agc capabilities --output json --pretty
agc capabilities --output table
agc capabilities --output markdown
```

JSON 是默认格式，适合 CI 和 Agent；table 适合人工浏览；markdown 适合把结果写入任务记录或文档。

## 9. 本地 REST 与 Web Command Center

启动本地服务：

```bash
agc web-server --addr :8421
```

常用路由：

```text
GET  /api/v1/capabilities
GET  /api/v1/endpoints
GET  /api/v1/openapi.json
POST /api/v1/:family/endpoints/:id/invoke
```

启动 Web：

```bash
npm --prefix apps/web install
npm --prefix apps/web run dev
```

Vite 会把 `/api` 代理到 `http://127.0.0.1:8421`。本地服务未启动时，页面自动显示内置的 13 个 API 家族和 156 个接口演示数据。

## 10. 全局参数与环境变量

```text
--output json|table|markdown
--pretty
--timeout 60s
--profile <name>
--project <directory>
```

支持的环境变量：

- `AGC_ACCESS_TOKEN`：直接提供 Bearer token。
- `AGC_CREDENTIALS_PATH`：覆盖默认凭据文件位置。

优先级：接口命令的 `--token` 高于 `AGC_ACCESS_TOKEN`；两者都未设置且命令真实执行时，CLI 才会通过选中的 profile 获取 token。

## 11. 常见问题

### 提示 `credential profile "..." not found`

运行 `agc auth list --output table` 核对名称；修改 `.agc/project.json` 的 `profile`，或显式传入 `--profile`。

### 命令只输出 URL，没有真正请求

这是预期的安全行为。确认命令包含 `--invoke`，并在完成 dry-run 审核后显式增加 `--dry-run=false`。

### Web 页面显示 `reference mode`

静态官网默认使用内置注册表。需要真实本地数据时，同时运行 `agc web-server --addr :8421` 和 Vite 开发服务器。

### 生产请求返回权限或字段错误

先检查 Service Account/API Client 权限、AppGallery Connect 项目归属、接口所需字段和 API 区域；再用 `--pretty` 查看错误 envelope。不要通过保存浏览器 Cookie、华为账号密码或绕过验证码来修复鉴权。

## 12. 开发与验证

```bash
make test
make coverage
make coverage-check
make ci
```

`make ci` 会运行 Go vet、Go 测试、80% 覆盖率门槛、Web 测试和 Web 构建。
