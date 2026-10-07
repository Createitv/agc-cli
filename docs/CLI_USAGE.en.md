# agc CLI user guide

[Project README](../README.en.md) · [简体中文](CLI_USAGE.md)

## Quick start

You need an AppGallery Connect app ID and either a Service Account JSON file or API Client credentials. Find your app ID in the app information section of the [AppGallery Connect console](https://developer.huawei.com/consumer/cn/service/josp/agc/index.html).

### 1. Install

On macOS, with Homebrew already installed:

```bash
brew tap createitv/tap
brew install agc-cli
agc version
```

For Windows and Linux, see [installation and upgrades](#installation-and-upgrades). Release binaries require neither Go nor Node.js.

### 2. Save credentials

Keep your Service Account JSON at a stable local path:

```bash
agc auth login --service-account-file ~/.agc/service-account.json --name production
agc auth check
```

`auth check` displays local credential configuration. It does not verify that Huawei accepts those credentials. See [authentication and multiple accounts](#authentication-and-multiple-accounts) for alternatives.

### 3. Bind your project

In your app's project directory, replace the example ID:

```bash
agc init --app-id YOUR_APP_ID --default-profile production
```

This writes `.agc/project.json`. Subsequent commands select the project's credential profile automatically. **Endpoint requests still require explicit appId, projectId, and package parameters where applicable.**

### 4. Query your first app

Preview the request without contacting Huawei:

```bash
agc publishing app-info-query --invoke --query appId=YOUR_APP_ID --query lang=en-US --pretty
```

The response shows `data.dryRun: true`, the HTTP method, and the target URL. Send the actual query after inspection:

```bash
agc publishing app-info-query --invoke --query appId=YOUR_APP_ID --query lang=en-US --dry-run=false --pretty
```

If the endpoint requires a `client_id` header, add `--header client_id=YOUR_CLIENT_ID`. The CLI currently does not populate this header automatically.

## Features and implementation status

The registry contains **159 entries**: 156 Huawei official endpoints/callbacks, two upload URL handoff operations, and one local Hvigor bridge entry. Registration means discovery and generic request construction are available; it does not mean every endpoint has been verified in production.

| Task | Command | Entries |
| --- | --- | ---: |
| Query/update app metadata, localizations, and submission requests | `agc publishing` | 14 |
| Discover upload URL, multipart, and upload confirmation endpoints | `agc upload` | 6 |
| Manage HarmonyOS certificates, provisioning profiles, devices, and fingerprints | `agc provisioning` | 17 |
| Query/update atomic service domains | `agc domains` | 5 |
| Manage test versions, testers, groups, and feedback | `agc testing` | 27 |
| Request download, sales, and other reports | `agc reports` | 12 |
| Query teams, projects, apps, and SDK configuration | `agc projects` | 8 |
| Query ratings/reviews and manage replies | `agc comments` | 8 |
| Manage products, subscriptions, prices, and promotions | `agc pms` | 40 |
| Discover game playing and game item callback contracts | `agc gameplay` / `agc game-items` | 8 / 2 |
| Discover resource predownload endpoints | `agc resources` | 8 |
| Inspect the local Hvigor build contract | `agc cicd` | 1 |

Discovery, credential selection, request construction, generic HTTP invocation, local REST, and OpenAPI are implemented. Binary/multipart file uploads are not yet orchestrated as a complete workflow. The Hvigor entry is not yet connected to a local build executor.

## Documentation by task

| Task | Start here |
| --- | --- |
| Install, log in, and run your first query | [Quick start](#quick-start) |
| Accounts, project configuration, and environment variables | [Authentication](#authentication-and-multiple-accounts) |
| Upload and submit an app | [App release workflow](#app-release-workflow) |
| Discover commands, flags, and official references | [Command discovery](#command-discovery-and-request-flags) |
| Integrate scripts or AI agents | [JSON and agents](#json-and-agents) |
| Use the browser UI or local API | [Web, REST, and OpenAPI](#web-rest-and-openapi) |
| Troubleshoot | [Common questions](#common-questions) |
| Develop, test, and release the CLI itself | [Development](#development) |

The [CLI guide (Chinese)](../docs/CLI_USAGE.md) contains additional examples. The [implementation plan (English)](../docs/AGC_CLI_FULL_PLAN.md) describes architecture and planned work; current source and the status notes above determine what is implemented.

## Installation and upgrades

### macOS: Homebrew

```bash
brew tap createitv/tap
brew install agc-cli
# Upgrade an existing installation
brew update
brew upgrade agc-cli
```

### Windows: Scoop

With Scoop already installed, run in PowerShell:

```powershell
scoop bucket add createitv https://github.com/Createitv/scoop-bucket
scoop install agc-cli
agc version
# Upgrade an existing installation
scoop update agc-cli
```

See the [Winget manifest PR](https://github.com/microsoft/winget-pkgs/pull/415361) for its publishing history. Check availability with `winget search --id Createitv.AgcCli -e` before using it.

### macOS / Linux / Windows: release archives

Choose the archive matching your OS and architecture from [Releases](https://github.com/Createitv/agc-cli/releases). Verify it against that release's `checksums.txt`, extract it, place `agc` (`agc.exe` on Windows) on PATH, and run `agc version`. Replace the binary to upgrade.

The release configuration targets macOS/Linux amd64 and arm64, plus Windows amd64. Consult each release for its actual assets.

### Install with Go

Requires Go 1.22 or later:

```bash
go install github.com/Createitv/agc-cli/cmd/agc@latest
```

Add Go's binary installation directory to PATH, normally `$(go env GOPATH)/bin` by default. `agc version` may show `dev` with this method; release archives receive version metadata from GoReleaser.

## Authentication and multiple accounts

Service Account JSON must include `key_id`, `private_key`, and `sub_account`. The CLI signs a PS256 JWT using this file. Login stores the file path, so the file must remain available locally.

Alternatively, save API Client credentials using variables you set beforehand:

```bash
agc auth login --client-id "$AGC_CLIENT_ID" --client-key "$AGC_CLIENT_KEY" --name staging
agc --profile staging auth check
```

These variables are shell argument examples. The CLI does not automatically read `AGC_CLIENT_ID` or `AGC_CLIENT_KEY`.

Profile precedence: `--profile` → project `.agc/project.json` profile → globally active account (or the only account). Each login makes that account active. Override an account or select another project:

```bash
agc --profile staging publishing endpoints
agc --project ../another-app auth check
```

`agc init` also accepts `--project-id` and `--package-name` to store project context.

| Environment variable | Purpose |
| --- | --- |
| `AGC_ACCESS_TOKEN` | Bearer token for endpoint requests; below `--token`, above profile credentials |
| `AGC_CREDENTIALS_PATH` | Override `~/.agc/credentials.json`; below `--credentials-path` |

`agc auth token` outputs a token. Currently, API Client output from `auth login`, `auth list`, and `auth check` may include `clientKey`. Do not publish that output in issues or public CI logs. Keep credential files, Service Account private keys, and access tokens locally or in CI secret storage.

## Command discovery and request flags

Inspect definitions without logging in:

```bash
agc --help
agc capabilities --output table
agc publishing endpoints --output table
agc publishing app-info-query --pretty
agc publishing app-info-query --help
agc endpoints --pretty
```

Each endpoint's `sourceUrl` links to its official Huawei reference, such as [query app information](https://developer.huawei.com/consumer/cn/doc/AppGallery-connect-References/agcapi-app-info-query-0000001158365045). See [endpoint_catalog.go](../pkg/domain/endpoint_catalog.go) for all mappings.

| Invocation | Behavior |
| --- | --- |
| Without `--invoke` | Display the endpoint definition |
| `--invoke` | Build a dry-run request; output method, URL, and `dryRun` |
| `--invoke --dry-run=false` | Send the request, including GET queries |

| Flag | Purpose |
| --- | --- |
| `--param key=value` | Repeatable path parameter |
| `--query key=value` | Repeatable query parameter |
| `--header key=value` | Repeatable HTTP header |
| `--field key=value` | Repeatable JSON field, encoded as a string |
| `--body request.json` | Read the full request body; takes precedence over `--field` |
| `--token TOKEN` | Explicit Bearer token |
| `--out response.json` | Save the raw response without format conversion |
| `--timeout 120s` | Request timeout; defaults to 60 seconds |

Use `--body` for arrays, objects, numbers, or booleans. Dry-run does not display full headers/body or validate all Huawei business fields. Inspect request files, permissions, and the official protocol before sending.

## App release workflow

A release usually involves **querying the app → uploading files → updating app file information → completing metadata → submitting for review → checking remote state**.

The CLI exposes individual endpoints; it does not yet orchestrate a one-command release. Inspect the definitions first:

```bash
agc upload endpoints --output table
agc publishing app-file-info --pretty
agc publishing app-info-update --pretty
agc publishing language-info-update --pretty
agc publishing app-submit --pretty
```

Binary/multipart uploads currently require an external tool implementing Huawei's protocol. `--body` does not replace that upload workflow. Prepare JSON files following the corresponding official references, then preview update/submission requests:

```bash
agc publishing app-file-info --invoke --body app-file-info.json --pretty
agc publishing app-info-update --invoke --body app-info.json --pretty
agc publishing app-submit --invoke --body submission.json --pretty
```

These commands do not send requests by default. Confirm uploads, file associations, and metadata, then add `--dry-run=false` to the step you intend to execute, together with required headers.

After submission succeeds, query AppGallery Connect for remote state. A successful request, approval, and publication are separate stages.

## JSON and agents

JSON is the default for `jq`, scripts, and AI agents. Use `--output table` or `--output markdown` for human-readable output.

```bash
agc capabilities --pretty
agc publishing app-info-query --pretty
```

Endpoint definitions include `_links` to local REST routes and `affordances` with next-command templates. Agents can inspect the definition and `sourceUrl`, fill in parameters, and preview requests. Templates do not confirm that business prerequisites are satisfied.

## Web, REST, and OpenAPI

Browse static endpoint references at [agccli.app](https://agccli.app/). For local REST data, start the API and Web development server from the source checkout:

```bash
agc web-server --addr :8421
# In another terminal; requires Node.js 20+
npm --prefix apps/web ci
npm --prefix apps/web run dev
```

Open the local URL printed by Vite. It proxies `/api` to `127.0.0.1:8421`; without the local API, the page shows reference data.

```bash
curl http://localhost:8421/api/v1/capabilities
curl http://localhost:8421/api/v1/endpoints
curl http://localhost:8421/api/v1/openapi.json
agc openapi --pretty
```

Preview a REST invocation:

```bash
curl -X POST http://localhost:8421/api/v1/publishing/endpoints/app-info-query/invoke \
  -H 'Content-Type: application/json' \
  -d '{"query":{"appId":"YOUR_APP_ID","lang":"en-US"},"dryRun":true}'
```

REST defaults to dry-run. Real requests require `"dryRun": false`, a `token` in the JSON payload, and any required `headers`. REST invocation does not automatically load saved CLI credential profiles.

## Common questions

| Problem | What to check |
| --- | --- |
| `agc: command not found` | Add the binary directory to PATH, reopen the terminal, then run `agc version` |
| Only a URL appears, no query result | Include `--invoke --dry-run=false` |
| `credential profile ... not found` | Check local profile names; override with `--profile` or update project configuration |
| `auth check` succeeds but requests fail | It only checks local configuration; inspect key files, permissions, project ownership, and headers |
| appId is still required after init | Only profile selection is automatic; pass app parameters explicitly |
| Downloaded output is not CSV/Excel | `--out` saves the raw response; some report endpoints return download URLs to fetch separately |
| Web shows reference mode | Start the API and Vite servers, then visit the Vite URL |

Report issues with `agc version`, OS/architecture, a redacted command, error details, and expected behavior in [GitHub Issues](https://github.com/Createitv/agc-cli/issues). Exclude keys and access tokens.

## Development

Requires Go 1.22+ and Node.js 20+ for Web development:

```bash
git clone https://github.com/Createitv/agc-cli.git
cd agc-cli
make build
./bin/agc version
npm --prefix apps/web ci
make ci
```

`make ci` runs Go vet, Go tests, the 80% coverage gate, Web tests, and the Web build. Individual targets include `make test`, `make coverage-check`, `make web-test`, and `make web-build`.

Key directories: `cmd/agc/command` (CLI), `pkg/agcapi` (auth/HTTP), `pkg/domain` (registry), `pkg/server` (local REST), and `apps/web` (website). Contributions through [pull requests](https://github.com/Createitv/agc-cli/pulls) are welcome.

For releases of the CLI itself, see [.goreleaser.yaml](../.goreleaser.yaml) and the [Release workflow](../.github/workflows/release.yml). A `v*` tag triggers publishing. Configuration includes release archives/checksums, Homebrew, Scoop, a Winget manifest PR, and GHCR images. Homebrew/Scoop/Winget publishing requires `TAP_GITHUB_TOKEN`. Actions and release pages show actual publishing results.

## License

[MIT](../LICENSE)
