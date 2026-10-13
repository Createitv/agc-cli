---
name: agc-cicd
description: Build unattended multi-platform CI/CD workflows around agc-cli. Use for Windows, macOS, or Linux installation, protected credential setup, dry-run gates, HarmonyOS Hvigor builds, release automation, and post-action verification.
---

# agc-cli CI/CD

Install the pinned version needed by the pipeline and verify it before work:

```bash
npm install --global @createitv/agc-cli@VERSION
agc version
```

Use an explicit profile and project directory in unattended jobs. Supply credentials through the CI provider's protected secret/file mechanism; never write secrets into the repository, workflow output, cache keys, or artifacts.

## Pipeline shape

1. Install and verify the pinned CLI version.
2. Materialize protected credentials with the narrowest permissions and delete ephemeral files in the runner cleanup phase.
3. Run `agc auth check` locally; use `--remote` only when an authorized live check is intended.
4. Resolve app/project/profile explicitly.
5. Discover the operation contract from the installed CLI.
6. Produce and retain a redacted dry-run artifact for review or policy gates.
7. Execute only approved live operations with `--dry-run=false`.
8. Query the remote state until the required terminal state, and preserve redacted evidence.

For HarmonyOS local builds, inspect the current local CI/CD operation first:

```bash
agc cicd endpoints --output table
agc cicd build-hvigor --pretty
```

## Reliability rules

- Pin stable versions; do not let production jobs float on an unverified latest release.
- Keep Windows, macOS, and Linux path/quoting differences inside shell-specific steps.
- Use JSON output for automation and parse it structurally.
- Do not treat a queued job, successful build, accepted upload, review submission, approval, and public availability as the same state.
- Do not blindly retry mutations after timeouts. Query remote state and make retries idempotent where Huawei provides a stable identifier.
- Make test resources and production resources distinct. A smoke test must not publish an app or alter an existing production resource.
