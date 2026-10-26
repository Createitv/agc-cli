---
name: agc-cli
description: Discover and invoke AppGallery Connect workflows with the agc CLI.
---

# AGC CLI

Use this skill for AppGallery Connect API discovery, request previews, and explicitly authorized operations.

## Discover before invoking

- Run `agc version --pretty` and `agc capabilities --pretty` to inspect the installed CLI.
- Run `agc endpoints --pretty` to discover endpoint names; do not invent commands.
- Read `agc <family> <endpoint> --docs --pretty` and `agc <family> <endpoint> --help` before constructing a request.
- Registered endpoints are not proof of successful live validation. Resolve required fields and prerequisites from endpoint documentation.

## Credentials and project context

- Inspect local credential configuration with `agc auth check`; this does not prove remote authorization.
- Use the user's explicitly selected profile via `--profile NAME` and project directory via `--project PATH`.
- Never print, commit, or include credentials in reports. Do not request secrets in chat.
- Provide endpoint appId/projectId parameters explicitly when required; project context does not automatically fill them.

## Preview, authorize, execute, verify

- Build a preview with `agc <family> <endpoint> --invoke` and documented `--query`, `--header`, or body flags. Invocation defaults to dry-run.
- Only add `--dry-run=false` when the requested live operation is authorized. Do not use publishing or deletion as a smoke test.
- For uploads, follow the documented upload URL, transfer, confirmation and parsing steps; do not assume a file ID means an upload is ready.
- Inspect machine-readable errors and documented prerequisites before retrying.
- Read back resulting state. Distinguish request acceptance, upload processing, review submission, approval and availability.

## Skill access

`agc skills list --pretty` lists the guides bundled into this binary.
`agc skills show agc-cli --output markdown` prints this guide as Markdown.
These commands work offline without credentials and do not install files or query a remote skill registry.
