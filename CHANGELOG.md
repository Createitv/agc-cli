# 变更日志 / Changelog

记录用户可见的变化；已发布安装包与版本见 [GitHub Releases](https://github.com/Createitv/agc-cli/releases)。

User-visible changes are listed here. Published binaries and versions are available in [GitHub Releases](https://github.com/Createitv/agc-cli/releases).

## Unreleased

## v0.2.0 — 2026-10-07

### Added

- 团队级 API Client JSON 导入、环境变量授权配置和远端授权检查。
- 注册接口的本地并发 HTTP 回归测试、测试资源生命周期脚本及参数契约测试。
- Linux、macOS、Windows 云端测试、最低 Go 1.22 兼容检查，以及与提交绑定的测试与覆盖率报告。
- 版本发布前授权检查、Homebrew/Scoop 清单自动更新和版本、下载地址、SHA256 回读验证；Release 附加验证报告。
- 中英文静态参考页、接口注册表导出及站点 SEO 检查。

### Fixed

- API 返回 HTTP 200 时的业务错误识别、原始响应体处理和文件流式上传。
- 凭证脱敏、原子保存、授权参数冲突检查及 profile 激活状态保留。
- REST 服务授权保护，以及 Provisioning、Testing、HarmonyOS Comments 等接口参数契约。

### Validation scope

- 离线测试验证程序和请求契约，不代表全部华为业务接口已经完成真实验收。


### Changed

- 默认中文 README，提供对应英文版，按快速开始、Agent 支持、功能覆盖和文档入口组织。/ Chinese-first README with an English version, organized around quick start, agents, features, and documentation links.
- 详细安装、鉴权、发布流程和常见问题移入中英文使用指南。/ Detailed installation, authentication, release workflows, and troubleshooting live in bilingual user guides.

### Added

- 文档索引、贡献指南与中英文 PR 模板。/ Documentation index, contributor guide, and bilingual PR template.
- `make check-docs` 检查本地链接、锚点、图片和代码块，并接入 PR CI。/ `make check-docs` validates local links, anchors, images, and code fences in PR CI.
