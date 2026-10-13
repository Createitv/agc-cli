---
name: agc-cli-usage
description: Use agc-cli to discover, preview, and invoke Huawei AppGallery Connect operations from terminals, CI, or agent workflows. Apply when a user asks an agent to inspect AGC capabilities, configure a project/profile, construct an API request, or perform and verify an authorized remote AGC change.
---

# agc-cli usage

Use `agc` as a discoverable, JSON-first interface to AppGallery Connect. Prefer the installed command's help and endpoint definition over remembered flags.

## Discover the command before acting

Start with the narrowest useful read-only command:

```bash
agc --help
agc capabilities --output table
agc <family> --help
agc <family> endpoints --output table
agc <family> <operation> --pretty
agc <family> <operation> --help
```

The operation definition includes its method, path, `sourceUrl`, `_links`, and command `affordances`. Treat affordances as command templates, not proof that remote prerequisites or permissions are satisfied. Consult the linked Huawei reference when request fields or business preconditions are unclear.

## Resolve credentials and project context

Inspect the selected profile before remote work:

```bash
agc auth check
agc auth list --output table
agc --profile staging auth check
```

Credential selection is `--profile`, then the profile in `.agc/project.json`, then the active or only saved profile. Bind stable project context when appropriate:

```bash
agc init --app-id "APP_ID" --project-id "PROJECT_ID" \
  --package-name "com.example.app" --default-profile "staging"
```

Use `agc auth login --help` before importing a service-account or API-client credential. Keep service-account JSON, private keys, client keys, credential files, and access tokens out of prompts, logs, issues, commits, and pull requests. `agc auth token` returns a secret; do not print or persist it unless the user's workflow explicitly requires a protected destination.

## Build and preview requests

Registered operations do not send a request unless `--invoke` is present. Invoked operations still default to dry-run:

```bash
agc publishing app-info-query \
  --invoke \
  --query appId="APP_ID" \
  --query lang="zh-CN" \
  --pretty
```

Use repeatable flags for simple values:

- `--param key=value` for path parameters
- `--query key=value` for query parameters
- `--header key=value` for required headers
- `--field key=value` for string-valued JSON fields

Use `--body request.json` for arrays, objects, numbers, booleans, or a reviewed full request body. Use `--out response.json` only when the raw response is an intended artifact. Dry-run output intentionally omits secrets and may not validate every Huawei business rule.

## Perform authorized remote changes

Only add `--dry-run=false` after the user has authorized the remote action and the preview matches the intended app, profile, endpoint, parameters, and body:

```bash
agc publishing app-info-update \
  --invoke \
  --body app-info.json \
  --dry-run=false
```

Before a mutation:

1. Resolve the exact app/project/profile instead of relying on an ambiguous active context.
2. Read the operation definition and help.
3. Preview with the default dry-run and review the target URL and inputs.
4. Confirm that uploads, metadata, permissions, and other business prerequisites are actually satisfied.
5. Send only the authorized operation.

Do not infer permission for submission, publication, deletion, replies, product changes, or other external mutations from a request to inspect or prepare them.

## Verify the result

Treat transport success, AGC request acceptance, review approval, and public availability as separate states. After a write, use a relevant read endpoint or the AppGallery Connect console to verify the remote state. Report the selected app/profile, operation, dry-run or live mode, response status, verification evidence, and anything still unverified.

For application releases, keep the sequence explicit: query app, upload externally where required, associate the file, update metadata, submit, then poll the remote status. agc-cli exposes these operations but does not currently provide a single transactional release command. A successful submit response is not approval or store publication.

## Output and automation

- JSON is the default and is preferred for agents and CI.
- Use `--pretty` only for readable JSON.
- Use `--output table` or `--output markdown` for human review.
- Keep data on stdout and diagnostics separate when wrapping commands.
- Pass `--timeout` explicitly for long requests when needed.
- Use `--project` or `--profile` explicitly in unattended multi-project jobs.

The local REST server and OpenAPI document expose the same registered operations:

```bash
agc web-server --addr 127.0.0.1:8421
agc openapi --pretty
```

REST invocations also default to dry-run, but they do not automatically reuse CLI credential profiles. Do not place a bearer token into a request unless the destination is trusted and the user authorized that transmission.

## Current boundaries

- Full binary or multipart upload orchestration may require an external uploader that follows Huawei's protocol.
- Endpoint definitions do not guarantee that the active account has the required role or that the app is ready for the operation.
- Use the installed CLI and its `sourceUrl` references when this skill and the current command surface differ.
