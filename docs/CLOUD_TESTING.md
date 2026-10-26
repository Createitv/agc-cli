# 云端测试与 Codex 工作流程

推送任意分支、创建或更新 PR、进入合并队列都会运行 GitHub Actions CI。也可以从 Actions 手动运行。检查包括文档、报告生成器、Go vet、Linux/macOS/Windows 的稳定版 Go 测试及 race 检测、Linux 的最低 Go 1.22 兼容检查、80% 总覆盖率门槛，以及 Web 测试与构建。

每个操作系统上传 `test-report-<os>-go-<version>-<commit SHA>`，保存 30 天，包含 JUnit、JSON、Markdown、原始 Go 测试事件和覆盖率。Actions Summary 显示提交、测试结果及覆盖率。报告绑定提交，不提交生成文件到源码仓库。PR 的 Checks 显示结果；失败时查看该次运行的日志和 artifact。

Release 调用相同 CI，全部通过后才执行发布。发布时将三平台报告和提交信息打包为 validation-reports.zip，附加到对应 GitHub Release，便于长期追溯。

Codex 每次修改先阅读根目录 AGENTS.md，修复有回归测试，本地运行 `make ci` 和 `go test -race ./...`，推送工作分支并检查对应 SHA 的云端结果后再合并。仓库管理员可以把 `Required checks` 设置为 main 分支的必需检查，要求 PR 和至少一位审核人，并禁止绕过检查。

离线 CI 不需要华为密钥，也不证明全部真实业务接口已验收。真实 API 验证需要有效授权、测试应用和专用资源，结果单独记录在 docs/verification。不要在外部贡献者 PR 中注入管理员密钥，不使用 pull_request_target 执行 PR 代码。

## Homebrew 与 Scoop 自动更新

发布稳定版时推送 `vMAJOR.MINOR.PATCH` 标签，Release 会先验证 TAP_GITHUB_TOKEN 对 Createitv/homebrew-tap（master）和 Createitv/scoop-bucket（main）的写权限，通过完整 CI 后由 GoReleaser 更新 agc-cli.rb 与 agc-cli.json，并回读版本、下载地址和 SHA256，与该 Release 的 checksums.txt 对照。无需手工编辑安装清单。普通分支提交不会发布安装包。预发布版本不覆盖稳定版安装渠道。

手动触发 Release 时必须选择已有版本标签，例如 `gh workflow run release.yml --ref v0.2.0`。不要对已发布版本重跑完整发布，应使用新版本标签。Check package publishers 工作流只检查授权及最新稳定版清单，不发布或修改资源。

用户安装不会被强制更新；Homebrew 用户执行 `brew update && brew upgrade agc-cli`，Scoop 用户执行 `scoop update && scoop update agc-cli`。跨仓库发布需要 TAP_GITHUB_TOKEN Secret；Actions 默认 GITHUB_TOKEN 只用于本仓库 Release/GHCR。
