# Codex Can Now Use agc-cli for HarmonyOS App Workflows

![agc-cli NPM and Agent Skills cover](https://typora-1300715298.cos.ap-shanghai.myqcloud.com/uPic/articles/agc-cli/20261013-npm-agents/agc-cli-npm-agents-cover-en.png)

AI coding assistants can write code, update tests, and prepare release notes. But once a workflow reaches AppGallery Connect, developers still spend time finding documentation, selecting endpoints, and assembling parameters by hand.

I wanted Codex and other coding agents to understand more than the existence of a command-line tool. They should know what to inspect first, how to choose the right workflow, how to preview a request, and when to stop for explicit user authorization.

That is why **agc-cli** now has two additional distribution layers:

1. an NPM package for installing the CLI on Windows, macOS, and Linux;
2. six task-oriented Agent Skills for Codex, Claude, and other compatible agents.

The website also publishes structured Schema, `llms.txt`, complete agent guides, and machine-readable API contracts so AI systems can read the project's capabilities without guessing from marketing copy.

- Source: [github.com/Createitv/agc-cli](https://github.com/Createitv/agc-cli)
- Website: [agccli.app](https://agccli.app/)
- NPM: [npmjs.com/package/@createitv/agc-cli](https://www.npmjs.com/package/@createitv/agc-cli)

## One NPM command across three operating systems

With Node.js 18 or later installed, the same command now works on Windows, macOS, and Linux:

```bash
npm install --global @createitv/agc-cli
agc version
```

The NPM package does not rely on a lifecycle script that downloads an unknown binary during installation. It contains checksum-verified release binaries for Windows x64, macOS x64/arm64, and Linux x64/arm64, with no `postinstall` download.

Developer workstations and unattended CI/CD machines can therefore use the same installation entry point. Homebrew, Scoop, and GitHub Releases remain available as alternatives.

## Agents need workflows, not only a command list

![Workflow from NPM installation to Agent Skills](https://typora-1300715298.cos.ap-shanghai.myqcloud.com/uPic/articles/agc-cli/20261013-npm-agents/agc-cli-agent-workflow.png)

Giving an AI model the output of `--help` still leaves it to guess operational boundaries. That is risky for app submissions, tester management, and release workflows.

agc-cli now provides six task-oriented Skills:

| Skill | Responsibility |
| --- | --- |
| `agc-cli-usage` | command discovery, dry runs, authorized execution, and verification |
| `agc-auth` | credential profiles, identity checks, and access diagnosis |
| `agc-publishing` | metadata, packages, review submission, and release state |
| `agc-testing` | test versions, groups, testers, and feedback |
| `agc-reports` | distribution, install, revenue, and retention reports |
| `agc-cicd` | cross-platform installation, protected credentials, and CI gates |

Install all six directly from GitHub:

```bash
npx --yes skills add Createitv/agc-cli --global --agent codex
```

Or install only the publishing workflow:

```bash
npx --yes skills add Createitv/agc-cli \
  --skill agc-publishing \
  --global \
  --agent codex
```

You can then ask Codex:

> Use agc-cli to inspect this HarmonyOS project's release readiness. Discover the endpoints and run dry previews first. Do not submit the app. List missing parameters and the next commands.

The Skill gives the agent a safer path: discover capabilities, read endpoint documentation, check credentials and targets, then preview the request. Real invocation remains an explicit decision.

## Schema explains what the project is; agent docs explain how to use it

![Human-readable pages and agent-readable resources describe agc-cli together](https://typora-1300715298.cos.ap-shanghai.myqcloud.com/uPic/articles/agc-cli/20261013-npm-agents/agc-cli-agent-schema.png)

The website now provides content for both people and AI systems:

- JSON-LD `SoftwareApplication` Schema for the project name, author, license, operating systems, install location, and core capabilities;
- `llms.txt` as a lightweight AI navigation file;
- `llms-full.txt` as a complete single-file project context;
- Chinese and English agent guides for task selection, safety boundaries, and verification;
- `skills.md` for the six Skills and their installation commands;
- `api-reference.json` for machine-readable methods, paths, parameters, and reference URLs.

These files do not guarantee search ranking, and they do not give an agent access to a Huawei account. They solve a more practical problem: Codex, Claude, and other models can retrieve consistent, explicit, and citable facts about agc-cli.

## Give agc-cli to Codex in five minutes

First, install the CLI:

```bash
npm install --global @createitv/agc-cli
agc version
```

Second, install the Skills:

```bash
npx --yes skills add Createitv/agc-cli --global --agent codex
```

Third, open Codex in your HarmonyOS project and describe the goal in natural language:

> Use agc-cli to list the available publishing endpoints and identify the credentials this project needs. Perform read-only checks and dry runs first. Do not upload a package or submit a release.

Service Account credentials still remain on the local machine. Never commit private keys, OAuth secrets, or `~/.agc/credentials.json`, and never paste them into a public conversation.

## An agent can run commands; authorization still belongs to the developer

Making a CLI understandable to an agent does not mean giving the agent unrestricted control of an app.

agc-cli endpoint commands are dry-run first. Uploading packages, changing metadata, submitting reviews, and releasing versions still require a clear target, the correct account permissions, and remote-state verification afterward.

The Skills and agent documentation preserve those boundaries:

- discover before invoking;
- preview before changing;
- keep credentials local;
- let the user explicitly decide submissions and releases;
- verify remote state because a successful command is not the same as a successful review.

If you use Codex, Claude, or another coding agent for HarmonyOS development, you can now bring part of the AppGallery Connect workflow into the same development loop.

Source: [github.com/Createitv/agc-cli](https://github.com/Createitv/agc-cli)

Installation and agent documentation: [agccli.app](https://agccli.app/)

Issues and feature requests: [GitHub Issues](https://github.com/Createitv/agc-cli/issues)
