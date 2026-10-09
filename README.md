# <img src="apps/web/public/agc-app-icon-192.png" width="36" height="36" alt="agc-cli"> agc-cli

**AppGallery Connect 命令中心** — 在终端管理你的华为应用。

[![CI](https://github.com/Createitv/agc-cli/actions/workflows/ci.yml/badge.svg)](https://github.com/Createitv/agc-cli/actions/workflows/ci.yml)
[![Go](https://img.shields.io/badge/Go-1.22%2B-00ADD8)](go.mod)
[![Platform](https://img.shields.io/badge/platform-macOS%20%7C%20Linux%20%7C%20Windows-blue)](https://github.com/Createitv/agc-cli/releases)
[![Release](https://img.shields.io/github/v/release/Createitv/agc-cli)](https://github.com/Createitv/agc-cli/releases)
[![License](https://img.shields.io/badge/license-MIT-green)](LICENSE)

**简体中文** · [English](README.en.md) · [官网](https://agccli.app/)

用 Go 编写的 AppGallery Connect CLI：查询应用、构建发布请求、调用测试、商品和报表 API，在终端或 CI 中完成操作。默认输出结构化 JSON，供脚本和 AI Agent 使用。

## 快速开始

需要应用 ID 和 Service Account JSON，可在 [AppGallery Connect](https://developer.huawei.com/consumer/cn/service/josp/agc/index.html) 管理应用和凭据。macOS 使用 Homebrew：

```bash
brew tap createitv/tap && brew install agc-cli

agc auth login --service-account-file ~/.agc/service-account.json --name production
agc init --app-id YOUR_APP_ID --default-profile production

# 预览查询；添加 --dry-run=false 才会发送请求
agc publishing app-info-query --invoke --query appId=YOUR_APP_ID --query lang=zh-CN
```

多账号、Windows/Linux 安装和 CI 环境变量见[使用指南](docs/CLI_USAGE.md)。从源码构建见[贡献指南](CONTRIBUTING.md)。项目初始化自动选择凭据 profile，应用参数仍需显式填写；接口需要 `client_id` 时用 `--header client_id=YOUR_CLIENT_ID` 提供。

## 为 AI Agent 设计

接口定义带有 `affordances` 命令模板与 `_links` REST 链接，Agent 可以发现下一步操作，无需记住完整命令树。例如 `agc publishing app-info-query --pretty` 的输出节选：

```json
{
  "data": {
    "id": "app-info-query",
    "method": "GET",
    "path": "/api/publish/v2/app-info",
    "affordances": {
      "show": "agc publishing app-info-query --pretty",
      "invoke": "agc publishing app-info-query --invoke"
    }
  }
}
```

默认 JSON；人工浏览可加 `--output table` 或 `--output markdown`。运行 `agc skills add --agent copilot` 安装内置 skill，也可用 `--agent all` 安装到 Copilot、Claude Code 和 Codex。`agc web-server` 提供本地 REST API。命令模板需要补齐参数，当前不会根据远程业务状态判断能否提交。详见 [Agent 与 REST 使用说明](docs/CLI_USAGE.md#json-与-agent-支持)。

## 功能覆盖

| 领域 | 可以做什么 |
| --- | --- |
| **应用与发布** | 查询和更新应用资料、软件包信息、多语言描述，构建提交审核和发布时间请求 → [发布流程](docs/CLI_USAGE.md#应用发布流程) |
| **上传** | 发现上传地址、分片和确认上传接口；完整二进制/multipart 上传尚需外部工具 → [上传说明](docs/CLI_USAGE.md#应用发布流程) |
| **测试** | 调用测试版本、测试用户、群组、邀请码和反馈接口 → `agc testing` |
| **商品与订阅** | 调用商品、价格、促销和多语言展示接口 → `agc pms` |
| **签名与设备** | 调用鸿蒙证书、Provisioning Profile、设备、ACL 和指纹接口 → `agc provisioning` |
| **评论与报表** | 查询评分与评论、管理回复、请求报表和保存原始响应 → `agc comments` / `agc reports` |
| **项目与域名** | 查询团队、项目、SDK 配置，调用元服务域名配置接口 → `agc projects` / `agc domains` |
| **游戏与资源包** | 查看游戏回调协议和资源包预下载接口 → `agc gameplay` / `agc game-items` / `agc resources` |
| **REST 与 Web** | 浏览接口注册表，预览和调用本地 REST，导出 OpenAPI → [本地使用](docs/CLI_USAGE.md#webrest-与-openapi) |

注册表包含 13 个 API 家族、159 个条目，支持通用请求构建与调用；各接口的字段、权限和前置条件以其 `sourceUrl` 华为文档为准。上传编排和本地 Hvigor 构建执行器尚未完成。

所有文档见[文档索引](docs/README.md)。所有命令与参数用 `agc --help`、`agc <家族> <接口> --help` 查看；`agc endpoints --pretty` 列出接口定义和官方参考链接。

## 更多

- [发布流程](docs/CLI_USAGE.md#应用发布流程) — 查询 → 上传 → 完善资料 → 提交审核 → 检查远程状态
- [贡献指南](CONTRIBUTING.md) — 构建、测试、代码组织和提交改动
- [变更日志](CHANGELOG.md) · [Release 下载](https://github.com/Createitv/agc-cli/releases)
- [反馈问题](https://github.com/Createitv/agc-cli/issues) — 附版本、脱敏命令和错误信息

## 许可证

[MIT](LICENSE)。
