# Agent Skills / Agent 技能

仓库按照 ASC Skills 的任务分层模式提供一组可移植 Agent Skills，帮助 Codex、Claude Code 及其他支持 `SKILL.md` 的 Agent 使用 agc-cli。

This repository follows the task-oriented ASC Skills pattern and ships a portable collection for Codex, Claude Code, and other agents that support `SKILL.md`.

| Skill | 用途 / Purpose |
| --- | --- |
| [`agc-cli-usage`](../skills/agc-cli-usage/SKILL.md) | 通用命令发现、dry-run、授权执行和验证 / General discovery, dry-run, authorized execution, and verification |
| [`agc-auth`](../skills/agc-auth/SKILL.md) | 凭据、profile、`whoami` 和授权诊断 / Credentials, profiles, `whoami`, and authorization diagnosis |
| [`agc-publishing`](../skills/agc-publishing/SKILL.md) | 应用资料、软件包、审核提交和发布状态 / Metadata, packages, review submission, and release state |
| [`agc-testing`](../skills/agc-testing/SKILL.md) | 测试版本、群组、用户、邀请码和反馈 / Test releases, groups, users, invitation codes, and feedback |
| [`agc-reports`](../skills/agc-reports/SKILL.md) | 分发、付费、留存和元服务报表 / Distribution, revenue, retention, and HarmonyOS service reports |
| [`agc-cicd`](../skills/agc-cicd/SKILL.md) | 跨平台无人值守安装与发布流水线 / Cross-platform unattended installation and release pipelines |

## 安装 / Install

先安装 CLI：

```bash
npm install --global @createitv/agc-cli
agc version
```

通过 agc-cli 浏览并安装全部技能：

```bash
agc skills list --remote
agc skills install --agent codex --global
```

也可以只安装一个技能：

```bash
agc skills install --name agc-publishing --agent codex --global
```

CLI delegates installation to the open Agent Skills installer, so Node.js and `npx` must be available. The equivalent direct commands are:

```bash
npx --yes skills add Createitv/agc-cli --global --agent codex
npx --yes skills add Createitv/agc-cli --skill agc-publishing --global --agent codex
```

不使用安装器时，可以把所需的完整 Skill 目录复制到 Agent 的技能目录，例如 Codex 的 `~/.codex/skills/` 或 Claude Code 的 `~/.claude/skills/`。保留整个目录，以便 `agents/openai.yaml` 等客户端元数据随技能一起安装。

Without the installer, copy each complete skill directory into the agent's skill directory.

## 使用 / Use

Codex 中可显式调用：

```text
Use $agc-publishing to inspect the app release, preview the request, and execute only the authorized change.
```

支持自动发现的 Agent 也可以根据 `SKILL.md` 的 description 自动选择该技能。技能不会扩大用户授权：查看、预览、提交、删除、发布和审核状态仍是不同操作。

Agents that support automatic discovery can select the skill from its description. The skill does not expand user authorization: inspection, preview, submission, deletion, publication, and review status remain separate operations.

## 设计原则 / Design

这些技能对标 `tddworks/asc-cli-skills` 的分层方式：通用入口保持精简，高频工作流按任务拆分；具体命令优先从当前二进制的 `--help`、endpoint definition、`sourceUrl` 和 OpenAPI 动态发现。技能内容不会扩大用户对远程写入、提交、删除或发布的授权。

The collection follows the layered `tddworks/asc-cli-skills` pattern: a small general entry skill plus task-specific workflow skills. Exact commands come from the installed binary's `--help`, endpoint definitions, `sourceUrl`, and OpenAPI contract. A skill never expands authorization for remote writes, submission, deletion, or publication.
