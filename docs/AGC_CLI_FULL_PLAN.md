# agc-cli Full Implementation Plan

## Summary

Build a Go-based `agc` CLI for Huawei AppGallery Connect automation. The project includes CLI commands, reusable Go packages, local REST APIs, a React/Vite Command Center, a console/reference website, tests, coverage gates, GitHub Actions, GoReleaser, Homebrew, Docker, and GitHub Pages deployment.

## Scope

Official Connect API families to expose:

- Publishing API
- Upload Management API
- Provisioning API
- Domain Management API
- Testing API
- Reports API
- Project Management API
- Comments API
- PMS API
- 在玩服务
- 游戏道具商城
- 资源包预下载
- CI/CD 平台

Each public API family gets a CLI command tree, REST route, web navigation entry, feature documentation, and tests. The current implementation registers 156 interface entries across the 13 API families: 153 Huawei official AppGallery Connect Reference endpoints/callbacks, 2 upload URL handoff operations, and 1 local Hvigor bridge. If Huawei does not expose a public endpoint for a behavior, return `unsupportedReason` and document the limitation instead of scraping browser sessions.

## Architecture

- `cmd/agc`: command entrypoint.
- `pkg/domain`: pure models, semantic state, affordances.
- `pkg/agcapi`: Huawei Connect API HTTP clients, token flow, upload management, pagination, error mapping.
- `pkg/domain/endpoint.go`: endpoint registry for all registered Connect API interfaces.
- `pkg/agcapi/token.go`: Service Account JWT and API Client token providers.
- `pkg/agcapi/invoke.go`: generic endpoint invoker with path/query/body/file input mapping.
- Endpoint commands support `--out` for raw CSV/Excel/PDF/binary responses.
- `pkg/project`: `.agc/project.json` context.
- `pkg/server`: local REST server for web and agent workflows, including per-endpoint dry-run/invoke routes.
- `pkg/domain/openapi.go`: OpenAPI 3.1 export for the local REST API and all registered endpoint invocation routes.
- `pkg/hvigor`: DevEco/Hvigor build orchestration.
- `apps/web`: homepage, Command Center, and console.
- `docs/features`: feature-level docs.

## Auth

Default to service account or API client authentication:

```bash
agc auth login --service-account-file ~/.agc/service-account.json --name prod
agc auth login --client-id <id> --client-key <key> --name legacy
agc auth token
```

Service Account mode generates the official PS256 JWT bearer token from `key_id`, `private_key`, and `sub_account`. API Client mode exchanges `client_id/client_key` through `/api/oauth2/v1/token`.

Credential precedence:

1. Explicit token: `--token`, then `AGC_ACCESS_TOKEN`
2. Credential selector: `--profile`
3. Project-bound default: `.agc/project.json` profile
4. Globally active credential in `~/.agc/credentials.json`

`agc init --default-profile <name>` binds a credential profile to the current project. Authentication checks, token generation, and endpoint invocation all use that bound profile automatically unless `--profile` overrides it.

Never save Huawei account passwords or browser cookies.

## Web

The web surface follows the same product shape as `asccli.app`:

- Independent Release Rail visual system using an AppGallery-red accent, technical dark surfaces, and `AUTH → PROFILE → PACKAGE → RELEASE` workflow semantics.
- Homepage with real commands, profile precedence, API coverage, dry-run safety, and local REST examples.
- Interactive Command Center with left-side API family navigation and capability detail.
- Console/reference area connected to `agc web-server`, with built-in reference mode when local REST is unavailable.
- SEO metadata, manifest, sitemap, security headers, responsive layout, keyboard focus, and reduced-motion support for `agccli.app`.

## REST & OpenAPI

Local discovery endpoints:

- `GET /api/v1`
- `GET /api/v1/capabilities`
- `GET /api/v1/endpoints`
- `GET /api/v1/openapi.json`
- `GET /api/v1/{family}`
- `GET /api/v1/{family}/endpoints`
- `GET /api/v1/{family}/endpoints/{endpointID}`

Every registered AppGallery Connect interface also has:

- `POST /api/v1/{family}/endpoints/{endpointID}/invoke`

The invocation payload accepts `baseUrl`, `params`, `query`, `fields`, raw JSON `body`, `token`, and `dryRun`. REST `dryRun` defaults to `true`; production calls require an explicit token in the payload. This keeps the local web surface useful for request construction without silently sending privileged requests.

## Testing

Coverage gate: 80%.

Required tests:

- CLI help and argument validation for every command.
- Endpoint registry coverage for every API family.
- CLI endpoint dry-run and request-building tests.
- Domain model JSON and affordances.
- AGC client request and error mapping with `httptest`.
- REST routes and `_links`.
- Web component render tests.
- Release/build configuration smoke checks.

## CI/CD

GitHub Actions:

- `ci.yml`: Go vet, Go tests, coverage gate, web tests, web build.
- `release.yml`: GoReleaser on `v*` tags, including GitHub Release archives, Homebrew formula, and GHCR image publish.
- `pages.yml`: web docs deployment.

Local smoke checks:

- `make release-snapshot`: archives and checksums without requiring Docker.
- `make release-snapshot-docker`: full snapshot including Docker image build.

GoReleaser outputs:

- macOS arm64/amd64
- Linux arm64/amd64
- Windows amd64
- checksums
- Homebrew formula
- Docker image
