# agc-cli

[![Website](https://img.shields.io/badge/website-agccli.app-f04444)](https://agccli.app)
[![Coverage](https://img.shields.io/badge/coverage-80%25%20gate-brightgreen)](#testing--coverage)
[![License](https://img.shields.io/badge/license-MIT-green)](LICENSE)

AppGallery Connect, from your terminal. `agc-cli` is a Go command center for Huawei AppGallery Connect automation across CLI, CI, local REST, and a web Command Center.

## Install from source

```bash
make build
install -m 0755 bin/agc /usr/local/bin/agc
```

```bash
agc version
```

Release archives, Homebrew, and container publishing are configured but require a real repository owner and release target before use.

## Quick Start

```bash
agc auth login \
  --service-account-file ~/.agc/service-account.json \
  --name production
agc auth check

agc init \
  --app-id <app-id> \
  --project-id <project-id> \
  --package-name com.example.app \
  --default-profile production
agc capabilities --output table
agc web-server --addr :8421
```

Credential selection is deterministic: `--profile` overrides the profile stored in `.agc/project.json`, which overrides the globally active credential. This lets a repository bind to `production` once while still allowing `agc --profile staging ...` for a single command.

`agc` defaults to JSON output and includes affordances so agents can discover the next legal command without hard-coding the command tree.

```bash
agc endpoints --pretty
agc openapi --pretty
agc publishing endpoints --output table
agc publishing app-info-query --invoke --query appId=123 --query lang=zh-CN --dry-run
agc publishing add-packageurl --invoke --field appId=123 --field packageUrl=https://example.com/app.app --dry-run
agc reports appdownloadexport --invoke --param appId=123 --query from=2026-08-01 --query to=2026-08-11 --out downloads.csv --dry-run=false
```

## API Coverage

| API family | Command | Registered interfaces |
| --- | --- | --- |
| Publishing API | `agc publishing` | 14 |
| Upload Management API | `agc upload` | 6 |
| Provisioning API | `agc provisioning` | 17 |
| Domain Management API | `agc domains` | 5 |
| Testing API | `agc testing` | 27 |
| Reports API | `agc reports` | 12 |
| Project Management API | `agc projects` | 8 |
| Comments API | `agc comments` | 8 |
| PMS API | `agc pms` | 40 |
| 在玩服务 | `agc gameplay` | 8 |
| 游戏道具商城 | `agc game-items` | 2 |
| 资源包预下载 | `agc resources` | 8 |
| CI/CD 平台 | `agc cicd` | 1 |

The repository now registers 156 interface entries: 153 Huawei official AppGallery Connect Reference endpoints/callbacks, 2 upload URL handoff operations, and 1 local Hvigor bridge. Each entry has an ID, method, path template, official document URL, parameter locations (`path`, `query`, `header`, `body`, `file`), CLI command, REST route, OpenAPI operation, and dry-run/invoke path through `pkg/agcapi`.

## Local REST & OpenAPI

```bash
agc web-server --addr :8421
curl http://localhost:8421/api/v1/capabilities
curl http://localhost:8421/api/v1/openapi.json
```

Every registered endpoint also has a REST invocation route:

```bash
curl -X POST http://localhost:8421/api/v1/publishing/endpoints/app-info-query/invoke \
  -H 'Content-Type: application/json' \
  -d '{
    "baseUrl": "https://connect-api.cloud.huawei.com",
    "query": {"appId": "123", "lang": "zh-CN"},
    "headers": {"client_id": "<client-id>"},
    "dryRun": true
  }'
```

`dryRun` defaults to `true` on the REST API. Set `"dryRun": false` and pass `"token": "<access-token>"` only when you want the local server to send the request to AppGallery Connect. The REST layer validates required `params`, `query`, `headers`, `fields`, or raw JSON `body` before calling the shared `pkg/agcapi` invoker.

## Web Command Center

```bash
npm --prefix apps/web install
npm --prefix apps/web run dev
```

Open the Vite URL and start `agc web-server --addr :8421` for live local REST data. If the server is not running, the web app uses demo data.

Production website: [agccli.app](https://agccli.app)

## Documentation

- [CLI usage guide](docs/CLI_USAGE.md)
- [Full implementation plan](docs/AGC_CLI_FULL_PLAN.md)
- `docs/AGC_CLI_使用手册.docx` for the formatted Word handbook

## Testing & Coverage

```bash
make test
make coverage
make coverage-check
make ci
```

The CI gate requires at least 80% Go test coverage. Every command must have help/argument coverage, success output coverage, failure coverage, and REST coverage when exposed over `/api/v1`.

## Release

Tags named `v*` trigger GoReleaser:

- GitHub Release archives and checksums.
- Homebrew formula update.
- Docker image pushed to GHCR.
- Static web documentation deployed by GitHub Pages.

Local release checks:

```bash
make release-snapshot
make release-snapshot-docker # requires Docker
```

## Security

Use service account or API client credentials. Service Account credentials generate the official PS256 JWT bearer token from `key_id`, `private_key`, and `sub_account`; API Client credentials call `/api/oauth2/v1/token`. Do not store Huawei account passwords, browser cookies, or captcha-backed sessions in this CLI. `.agc/project.json` stores project context only; secrets belong in `~/.agc/credentials.json` or your CI secret store.
