# Platform copy — agc-cli NPM and Agent Skills

## V2EX

Title: agc-cli 现在支持 NPM 安装和 Codex Agent Skills 了

Body:

我给 agc-cli 补上了两层面向 AI 编程助手的分发能力：

- Windows、macOS、Linux 现在可以统一使用 `npm install --global @createitv/agc-cli`；
- Codex 可以通过 `npx --yes skills add Createitv/agc-cli --global --agent codex` 安装 6 个任务型 Skills。

官网同时增加 JSON-LD Schema、`llms.txt`、`llms-full.txt`、中英文 Agent Guide、Skills 索引和机器可读 API 契约。目标不是让 Agent 猜命令，而是让它沿着“发现能力 → 读取文档 → 检查凭据 → dry-run → 用户授权 → 验证远程状态”的路径工作。

项目：https://github.com/Createitv/agc-cli

官网：https://agccli.app/

NPM：https://www.npmjs.com/package/@createitv/agc-cli

完整中文文章发布后补充链接。

## 掘金 / CSDN / 华为开发者社区 / 51CTO

Title: 现在，Codex 可以直接使用 agc-cli 管理鸿蒙应用了

Summary: agc-cli 已提供跨平台 NPM 安装、6 个任务型 Agent Skills，以及 JSON-LD、llms.txt、Agent Guide 和 API JSON 等机器可读资源。Codex 和其他 Coding Agent 可以更快理解 AppGallery Connect 工作流，同时继续遵守 dry-run、凭据留在本地和用户明确授权等安全边界。

Suggested tags: HarmonyOS, AppGallery Connect, CLI, Codex, AI Agent, NPM, 开源

Source: `agc-cli-npm-agents.zh-CN.md`

## 微信公众号

Title: 现在，Codex 可以直接使用 agc-cli 管理鸿蒙应用了

Digest: 我为 agc-cli 增加了跨平台 NPM 安装、6 个 Agent Skills 和面向 AI 的结构化文档。现在 Codex 可以从发现接口、检查凭据和 dry-run 开始，更安全地参与鸿蒙应用管理流程。

Cover: `assets/agc-cli-npm-agents-cover-zh-CN.png`

Source: `agc-cli-npm-agents.zh-CN.md`

Publication rule: create a draft first. Do not mass-send without explicit confirmation.

## DEV Community

Title: Codex Can Now Use agc-cli for HarmonyOS App Workflows

Description: agc-cli now ships through NPM and provides six task-oriented Agent Skills, JSON-LD Schema, llms.txt, agent guides, and machine-readable API contracts for Codex and other coding agents.

Tags: harmonyos, ai, cli, opensource

Canonical URL: leave empty unless an English canonical article has already been published elsewhere.

Source: `agc-cli-npm-agents.en.md`
