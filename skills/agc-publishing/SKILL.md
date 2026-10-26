---
name: agc-publishing
description: Manage AppGallery Connect application metadata, packages, staged rollout, and review submission with agc-cli. Use when preparing, previewing, submitting, or checking a Huawei or HarmonyOS app release.
---

# AppGallery Connect publishing

Use the installed endpoint registry as the source of truth:

```bash
agc publishing endpoints --output table
agc publishing <operation> --pretty
agc publishing <operation> --help
```

An operation without `--invoke` shows its contract. `--invoke` still defaults to dry-run. Inspect the returned method, URL, parameters, `sourceUrl`, and affordances before sending anything.

## Release workflow

1. Resolve the exact profile, app, project, package name, language, and release type.
2. Query current app/version state with the matching read endpoint.
3. Prepare metadata and package inputs outside the command when required.
4. Preview every request with `--invoke`; keep the default dry-run.
5. Add `--dry-run=false` only for the operation the user authorized.
6. Read back the app/version state and continue polling when processing or review is asynchronous.

Example preview:

```bash
agc --profile production publishing app-info-query \
  --invoke --query appId="APP_ID" --query lang="zh-CN" --pretty
```

For structured payloads, prefer a reviewed JSON file:

```bash
agc --profile production publishing app-info-update \
  --invoke --body app-info.json --pretty
```

Only after authorization and a matching preview:

```bash
agc --profile production publishing app-info-update \
  --invoke --body app-info.json --dry-run=false --pretty
```

## Important boundaries

- Package upload may require Huawei's external upload protocol or a reachable download URL; `--body` is not a binary uploader.
- Updating metadata, associating a package, submitting review, approval, staged rollout, and public availability are separate states.
- A successful HTTP response does not prove approval or publication.
- Do not automatically retry a submission or other mutation after a timeout or ambiguous response. Query the remote state first.
- Deletion, submission, rollout changes, and publication each require their own user authorization.
