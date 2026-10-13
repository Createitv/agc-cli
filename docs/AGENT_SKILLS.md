# Agent Skill / Agent 技能

仓库提供可移植的 [`agc-cli-usage`](../skills/agc-cli-usage/SKILL.md) Agent Skill，帮助 Codex、Claude Code 及其他支持 `SKILL.md` 的 Agent 安全地发现命令、选择凭据与项目、预览请求、执行已授权的 AppGallery Connect 操作并验证远程结果。

This repository ships the portable [`agc-cli-usage`](../skills/agc-cli-usage/SKILL.md) skill for Codex, Claude Code, and other agents that support `SKILL.md`. It teaches command discovery, credential and project selection, dry-run review, authorized AppGallery Connect invocation, and remote-state verification.

## 安装 / Install

先安装 CLI：

```bash
npm install --global @createitv/agc-cli
agc version
```

通过 Agent Skills 安装器全局安装仓库中的技能：

```bash
npx --yes skills add Createitv/agc-cli --global --agent codex
```

不使用安装器时，可以将 `skills/agc-cli-usage` 目录复制到 Agent 的技能目录，例如 Codex 的 `~/.codex/skills/` 或 Claude Code 的 `~/.claude/skills/`。保留整个目录，以便 `agents/openai.yaml` 等客户端元数据随技能一起安装。

Without the installer, copy the complete `skills/agc-cli-usage` directory into the agent's skill directory, such as `~/.codex/skills/` for Codex or `~/.claude/skills/` for Claude Code.

## 使用 / Use

Codex 中可显式调用：

```text
Use $agc-cli-usage to inspect the publishing operation, preview the request, and execute only the authorized change.
```

支持自动发现的 Agent 也可以根据 `SKILL.md` 的 description 自动选择该技能。技能不会扩大用户授权：查看、预览、提交、删除、发布和审核状态仍是不同操作。

Agents that support automatic discovery can select the skill from its description. The skill does not expand user authorization: inspection, preview, submission, deletion, publication, and review status remain separate operations.

## 设计原则 / Design

本技能对标 ASC Skills 的分层方式：入口技能只保存会改变 Agent 决策的稳定规则，具体命令优先从当前二进制的 `--help`、endpoint definition、`sourceUrl` 和 OpenAPI 动态发现。随着高频工作流成熟，可以再增加独立技能，例如发布编排、元数据同步或审核诊断，而不是让一个技能变成完整手册。

The skill follows the layered ASC Skills pattern: the entry skill keeps stable decision rules, while exact commands are discovered from the installed binary's `--help`, endpoint definitions, `sourceUrl`, and OpenAPI contract. Mature workflows can later become separate release, metadata, or review-diagnosis skills instead of turning one skill into a copied manual.
