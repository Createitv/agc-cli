# 现在，Codex 可以直接使用 agc-cli 管理鸿蒙应用了

![agc-cli NPM 与 Agent Skills 中文封面](https://typora-1300715298.cos.ap-shanghai.myqcloud.com/uPic/articles/agc-cli/20261013-npm-agents/agc-cli-npm-agents-cover-zh-CN.png)

AI 编程助手已经能写代码、改测试、整理发布说明，但当工作走到 AppGallery Connect，很多流程仍然需要开发者手动查文档、找接口、拼参数。

我希望 Codex 和其他 Coding Agent 不只知道“有一个命令行工具”，还能够理解：应该先检查什么、如何选择任务、怎样预演请求，以及什么时候必须停下来等待用户确认。

因此，我为 **agc-cli** 补上了两层新的分发能力：

1. 通过 NPM 在 Windows、macOS 和 Linux 上安装 CLI；
2. 通过 6 个任务型 Agent Skills，把 AppGallery Connect 工作流交给 Codex、Claude 等 Agent 使用。

同时，官网增加了面向搜索引擎和 AI 的结构化 Schema、`llms.txt`、完整 Agent 指南与 API 契约，让模型能够直接读取项目能力，而不是从营销页面里猜命令。

- 项目：[github.com/Createitv/agc-cli](https://github.com/Createitv/agc-cli)
- 官网：[agccli.app](https://agccli.app/)
- NPM：[npmjs.com/package/@createitv/agc-cli](https://www.npmjs.com/package/@createitv/agc-cli)

## 一条 NPM 命令，安装到三个系统

过去，不同开发环境通常需要分别准备 Homebrew、Scoop 或手动下载脚本。现在，只要机器上有 Node.js 18 或更高版本，就可以使用同一条命令：

```bash
npm install --global @createitv/agc-cli
agc version
```

这个 NPM 包不是在安装时临时下载未知二进制。它内置 Windows x64、macOS x64/arm64、Linux x64/arm64 的已发布文件，并通过 checksum 进行校验，也不依赖 `postinstall` 下载脚本。

这意味着开发电脑和无人值守的 CI/CD 机器可以使用同一种安装入口。

Homebrew、Scoop 和 GitHub Releases 仍然保留；NPM 只是增加了更统一的跨平台选择。

## 给 Agent 的不只是一份命令列表

![从 NPM 安装到 Coding Agent 使用 Skills 的流程](https://typora-1300715298.cos.ap-shanghai.myqcloud.com/uPic/articles/agc-cli/20261013-npm-agents/agc-cli-agent-workflow.png)

如果只把 CLI 的 `--help` 交给 AI，它仍然需要自己判断任务边界。对于应用提交、测试用户或报表下载，这种猜测并不可靠。

agc-cli 现在提供 6 个任务型 Skills：

| Skill | 负责的任务 |
| --- | --- |
| `agc-cli-usage` | 命令发现、dry-run、授权执行和结果验证 |
| `agc-auth` | 凭据 Profile、身份查询和权限诊断 |
| `agc-publishing` | 应用资料、软件包、提交审核和发布状态 |
| `agc-testing` | 测试版本、群组、测试用户和反馈 |
| `agc-reports` | 分发、安装、收入与留存报表 |
| `agc-cicd` | 跨平台安装、受保护凭据和 CI 门禁 |

当前可以直接从 GitHub 安装全部 Skills：

```bash
npx --yes skills add Createitv/agc-cli --global --agent codex
```

只安装发布相关 Skill：

```bash
npx --yes skills add Createitv/agc-cli \
  --skill agc-publishing \
  --global \
  --agent codex
```

安装后，你可以直接告诉 Codex：

> 使用 agc-cli 检查这个鸿蒙项目的发布准备情况。先发现接口并 dry-run，不要提交审核；把缺少的参数和下一步命令列出来。

Skill 会给 Agent 一条更稳妥的工作路径：先发现能力、读取接口文档、检查凭据和目标，再预演请求。只有用户明确授权时，才进入真实调用。

## Schema 让模型知道“这是什么”，文档让模型知道“怎么用”

![人类页面与 Agent 可读取资源共同描述 agc-cli](https://typora-1300715298.cos.ap-shanghai.myqcloud.com/uPic/articles/agc-cli/20261013-npm-agents/agc-cli-agent-schema.png)

官网现在同时为人和 AI 提供内容：

- JSON-LD `SoftwareApplication` Schema：说明项目名称、作者、许可证、运行平台、安装地址和主要能力；
- `llms.txt`：提供轻量的 AI 导航入口；
- `llms-full.txt`：提供一次读取所需的完整项目上下文；
- 中英文 Agent Guide：解释任务选择、安全边界和验证方式；
- `skills.md`：列出 6 个 Skills 及安装方式；
- `api-reference.json`：提供机器可读的接口方法、路径、参数和参考链接。

这些文件不会自动让项目获得搜索排名，也不会赋予 Agent 华为账号权限。它们解决的是另一个问题：当 Codex、Claude 或其他模型读取 agc-cli 时，可以获得一致、明确、可引用的项目事实。

## 五分钟把 agc-cli 交给 Codex

第一步，安装 CLI：

```bash
npm install --global @createitv/agc-cli
agc version
```

第二步，安装 Skills：

```bash
npx --yes skills add Createitv/agc-cli --global --agent codex
```

第三步，在你的鸿蒙项目目录中打开 Codex，并用自然语言说明目标：

> 使用 agc-cli 查看当前可用的发布接口，检查项目需要哪些凭据。先执行只读检查和 dry-run，不要上传软件包，也不要提交审核。

涉及 Service Account 时，凭据仍然保存在本机。不要把私钥、OAuth 密钥或 `~/.agc/credentials.json` 提交到 Git，也不要把它们粘贴到公共对话中。

## AI 可以执行命令，但授权仍然属于开发者

让 Agent 理解 CLI，不等于让 Agent 无限制操作应用。

agc-cli 的接口命令默认 dry-run。上传软件包、修改资料、提交审核和发布版本，仍然需要明确的目标、正确的账号权限，以及执行后的远程状态确认。

这也是我设计 Skills 和 Agent 文档时保留的边界：

- 先发现，再执行；
- 先预演，再修改；
- 凭据留在本地；
- 提交和发布由用户明确决定；
- 命令成功不等于审核成功，必须读取远程状态。

如果你正在使用 Codex、Claude 或其他 Coding Agent 开发鸿蒙应用，现在可以把 AppGallery Connect 的一部分重复工作也放进同一个开发流程。

项目地址：[github.com/Createitv/agc-cli](https://github.com/Createitv/agc-cli)

安装与 Agent 文档：[agccli.app](https://agccli.app/)

问题和功能建议：[GitHub Issues](https://github.com/Createitv/agc-cli/issues)
