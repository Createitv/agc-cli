---
name: agc-cli
description: Use the agc CLI to manage AppGallery Connect workflows, discover API endpoints, configure profiles, and safely preview or invoke requests.
---

# agc-cli agent skill

This skill is the complete operational guide for using `agc-cli`; it does not depend on other repository files being present. The executable is named `agc`. It is for Huawei AppGallery Connect (AGC), not Google Cloud or Huawei Cloud.

## Safety rules

- Never ask a user to paste credentials, private keys, or access tokens into chat. Do not print or expose secrets in command output, logs, files, issues, or reports. Treat raw authentication responses as sensitive: API Client output may contain `clientKey`.
- Do not perform a remote check or send an API request unless the user authorized that specific remote action. Requests are dry-run by default; preserve that default until the user explicitly approves execution.
- Before any update, submission, reply, or other mutation, inspect the endpoint, required parameters, request body, target app/project, and dry-run preview. Do not infer business intent or fill missing data with guesses.
- Do not repeat a mutation automatically after a timeout or ambiguous response. First determine whether the remote action may have succeeded.
- Endpoint registration and a successful HTTP response do not prove business prerequisites, approval, or publication.

## Check the CLI and discover commands

The CLI can be installed from a release archive or package manager, or with Go 1.22+:

```bash
go install github.com/Createitv/agc-cli/cmd/agc@latest
agc version
```

Discover commands and endpoint contracts without credentials or network access:

```bash
agc --help
agc capabilities --output table
agc endpoints --pretty
agc publishing endpoints --output table
agc publishing app-info-query --pretty
agc publishing app-info-query --help
agc openapi --pretty
```

The API families include `publishing`, `upload`, `provisioning`, `domains`, `testing`, `reports`, `projects`, `comments`, `pms`, `gameplay`, `game-items`, `resources`, and `cicd`. Discover the installed CLI's current registry rather than relying on a fixed endpoint count.

An endpoint definition shows its method, path, parameter contract, official Huawei `sourceUrl`, REST links, and command affordances. Read the definition before constructing a request. Use its official reference for permissions, field semantics, and business prerequisites. Affordance strings are templates, not verified plans.

## Global flags and output

Global flags can be placed before the command:

| Flag | Purpose |
| --- | --- |
| `--output json\|table\|markdown` | Choose output format; JSON is the default and preferred for agents/scripts. |
| `--pretty` | Pretty-print JSON. |
| `--timeout 60s` | Set request timeout; default is 60 seconds. |
| `--profile NAME` | Select a saved credential profile. |
| `--project DIR` | Select project context directory; default is `.`. |

Parse JSON structurally and check command exit status. Do not assume every response has the same `data` shape. Preserve warnings and errors, and avoid echoing returned credentials, tokens, user data, or other sensitive values. Use table/markdown output only when it helps a person inspect listings.

## Credentials and project context

Credentials are stored locally, by default at `~/.agc/credentials.json`. `AGC_CREDENTIALS_PATH` overrides that path; the `auth` commands also accept `--credentials-path`. Profiles are selected in this order: explicit `--profile`, the project profile, globally active profile, or the only profile.

For credential setup, have the user configure the secret locally or make it available through their approved secret manager. Do not request the secret value in chat. Examples (the values must already exist in the user's protected environment):

```bash
agc auth login --service-account-file PATH --name production
agc auth login --client-id "$AGC_CLIENT_ID" --client-key "$AGC_CLIENT_KEY" --name staging
```

The CLI does not automatically read `AGC_CLIENT_ID` or `AGC_CLIENT_KEY`; shell expansion in the example supplies the values. Service Account JSON needs `key_id`, `private_key`, and `sub_account`; login stores its file path, so that local file must remain available. Avoid putting literal secrets in shell history or command transcripts.

Useful local credential commands:

```bash
agc auth list
agc --profile staging auth check
agc auth token
```

`auth check` inspects local credential configuration; it does not prove Huawei accepts it. `auth check --remote` contacts Huawei and may perform an app/project read to check permissions. `auth token` creates a real token. Use either only with explicit authorization, and never display or log their raw output.

Bind a project to an app and optionally a default profile:

```bash
agc init --app-id APP_ID --default-profile production
agc init --app-id APP_ID --project-id PROJECT_ID --package-name PACKAGE_NAME
```

This writes `.agc/project.json` under the selected `--project` directory. It selects project context/profile; it does not automatically supply every endpoint parameter. Pass `appId`, `projectId`, package name, and other required values explicitly when the endpoint contract requires them. `AGC_ACCESS_TOKEN` can provide an access token for endpoint requests; never print or expose it.

## Inspect, preview, and invoke an endpoint

Without `--invoke`, an endpoint command displays its definition. With `--invoke`, the CLI constructs a dry-run request by default. Only `--dry-run=false` sends it.

| Flag | Usage |
| --- | --- |
| `--param key=value` | Repeat for path parameters. |
| `--query key=value` | Repeat for query parameters. |
| `--header key=value` | Repeat for required HTTP headers, such as `client_id`. |
| `--field key=value` | Repeat for JSON fields; values are encoded as strings. |
| `--body FILE` | Read a JSON object request body; takes precedence over `--field`. Use for arrays, objects, numbers, or booleans. |
| `--token TOKEN` | Explicit request token; handle as a secret. |
| `--out FILE` | Save raw response bytes to a file. |
| `--base-url URL` | Override the API base URL, primarily for an explicitly selected compatible service/test environment. |
| `--dry-run=false` | Send the request; use only after explicit authorization and review. |

Example: inspect, then preview a read-only app information query:

```bash
agc publishing app-info-query --pretty
agc publishing app-info-query --invoke \
  --query appId=APP_ID --query lang=en-US --pretty
```

The preview identifies the method, URL, and `dryRun: true`; it does not contact Huawei and does not validate all business fields. After reviewing the preview and confirming authorization, append `--dry-run=false` to send. Supply `--header client_id=CLIENT_ID` when the endpoint contract requires it; this header is not automatically populated.

For JSON request bodies, prepare the file locally, review its contents, and keep secrets out of it:

```bash
agc publishing app-info-update --invoke --body request.json --pretty
```

If the endpoint requires a parameter, the CLI reports a missing `--param`, `--query`, `--header`, or `--field` value. Consult the endpoint's displayed parameter contract instead of guessing.

## Common workflows and limitations

- **App release:** Discover and invoke the separate query, upload handoff, metadata, and submission endpoints. The CLI does not orchestrate a full release. Binary/multipart file transfer is not a complete built-in workflow; use a separate tool that implements Huawei's documented upload protocol. `--body` does not upload binaries.
- **Reports:** Some report endpoints return raw responses or download URLs. `--out` preserves response bytes; it does not convert them into CSV/Excel or necessarily fetch a second-stage download URL.
- **Local REST:** `agc web-server` starts a local REST API (default listen address `127.0.0.1:8421`). REST invocation defaults to dry-run and does not load saved CLI credential profiles. Real REST invocation requires an explicit token and required headers in the request. Non-loopback listening requires `AGC_SERVER_TOKEN`; do not expose the server or token.
- **OpenAPI:** `agc openapi --pretty` exports the local REST and endpoint invocation contract.
- **Skills:** `agc skills add --agent copilot|claude|codex|all` installs this skill into the selected project's `.github/skills`, `.claude/skills`, or `.agents/skills` directory. Existing files are not overwritten unless `--force` is supplied.

## Troubleshooting

- `agc: command not found`: install the CLI, ensure its binary directory is on `PATH`, and run `agc version`.
- `credential profile ... not found`: inspect local profile names and select the intended profile or project.
- `auth check` passes but an API call fails: local check is not remote verification; review credentials locally, account permissions, app/project ownership, required headers, and the endpoint contract.
- A command shows only an endpoint definition or preview: actual invocation requires `--invoke`; sending also requires `--dry-run=false`.
- `appId` is still required after `agc init`: project setup does not necessarily populate endpoint parameters; pass the required value explicitly.

When reporting an error, include CLI version, OS/architecture, a redacted command, and the error. Never include credentials, tokens, private keys, or unredacted auth output.
