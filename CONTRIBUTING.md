# 贡献指南

**简体中文** · [English](#english)

欢迎改进命令、API 适配、文档和测试。贡献使用本项目的 [MIT 许可证](LICENSE)。

## 构建与测试

需要 Go 1.22+；Web 开发需要 Node.js 20+。

```bash
git clone https://github.com/Createitv/agc-cli.git
cd agc-cli
make build
./bin/agc version
make test
npm --prefix apps/web ci
make ci
```

`make ci` 运行 Go vet、Go 测试、80% 覆盖率门槛、Web 测试和 Web 构建。发布配置见 [.goreleaser.yaml](.goreleaser.yaml) 和 [Release workflow](.github/workflows/release.yml)。

## 代码组织

| 目录 | 职责 |
| --- | --- |
| `cmd/agc/command` | CLI 命令、参数与输出 |
| `pkg/domain` | 接口定义、JSON 模型、链接与命令模板 |
| `pkg/agcapi` | 鉴权、HTTP 请求与错误处理 |
| `pkg/project` | 项目上下文配置 |
| `pkg/server` | 本地 REST 路由 |
| `apps/web` | 官网与 Web Command Center |

## 提交改动

1. 在 Issue 或 PR 中说明用户遇到的问题、预期行为和范围。
2. 创建分支，一次 PR 聚焦一个可验证的改动。
3. 对新行为或错误修复，先写能复现问题的测试，再实现代码；测试应断言用户可观察的结果。
4. 若改动影响 CLI 与 REST 共用能力，同时检查两种入口，避免参数或默认行为不一致。
5. 更新受影响的中英文文档与 [CHANGELOG](CHANGELOG.md)，运行 `make check-docs` 和相关检查。
6. 提交 PR，写明行为变化、验证结果和尚未验证的部分。

PR 使用仓库的[中英文模板](.github/pull_request_template.md)，填写兼容性和实际验证结果。CI 会在 PR 上检查文档并运行构建与测试。分支保护需由维护者在 GitHub Settings 中配置；模板勾选框不会自动阻止合并。

提交标题可用 `feat:`、`fix:`、`docs:`、`test:` 或 `chore:`。不要添加 `Co-Authored-By` trailer。

本地模拟、华为 API 返回成功、远程状态确认是不同层级的证据。不要把接口注册或 dry-run 成功写成生产功能验证完成。测试中使用临时目录、模拟 HTTP 和示例 ID，避免写入真实账号。

## 文档更新

README 保留介绍、快速开始、Agent 示例、功能表和文档入口。安装细节、完整流程和常见问题写入使用指南。新命令或行为变更需同步中文与英文版本；参数的最终参考是 `--help`。`make check-docs` 检查公开文档的本地链接、目录锚点、图片路径和代码块；不检查外部网站可用性。

错误报告附上版本、系统/架构、脱敏命令、错误信息和预期结果。当前 API Client 的部分 auth 输出包含密钥，分享前需删除密钥、令牌及私钥内容。

## English

Contributions are licensed under [MIT](LICENSE). Use Go 1.22+ and Node.js 20+ for Web work. The build/test commands and directory table above apply to both languages.

1. Describe the user problem, expected behavior, and scope in an issue or PR.
2. Create a branch; keep each PR focused on one verifiable change.
3. For new behavior or bug fixes, write a reproducing test before implementation. Assert observable results.
4. Check both CLI and REST when changing shared behavior.
5. Update affected Chinese/English docs and the [changelog](CHANGELOG.md); run `make check-docs` and relevant checks.
6. Open a PR describing behavior changes, validation, and anything not verified.

Use the [bilingual PR template](.github/pull_request_template.md) to document compatibility and actual validation. PR CI checks docs and runs builds/tests. Maintainers configure branch protection in GitHub Settings; template checkboxes do not enforce merge restrictions.

Use commit prefixes such as `feat:`, `fix:`, `docs:`, `test:`, or `chore:`. Do not add `Co-Authored-By` trailers.

Keep local simulation, successful Huawei responses, and remote state verification distinct. Registration and dry-run are not production verification. Use temporary directories, mocked HTTP, and placeholder IDs in tests.

Keep the README short: introduction, quick start, agent example, feature table, and links. Put detailed workflows in the user guide; synchronize both languages when behavior changes. `--help` is the final flag reference. `make check-docs` checks local links, anchors, images, and code fences in published docs, not external website availability.

Issue reports should include version, OS/architecture, a redacted command, errors, and expected behavior. Some API Client auth output currently includes keys; remove secrets before sharing logs.
