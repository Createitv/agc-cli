---
name: agc-testing
description: Manage AppGallery Connect open, internal, and invitation testing with agc-cli. Use for HarmonyOS test versions, test groups, testers, invitation codes, public links, feedback, self-test reports, and promoting or stopping a test release.
---

# AppGallery Connect testing

Discover the exact installed testing operations before constructing a workflow:

```bash
agc testing endpoints --output table
agc testing <operation> --pretty
agc testing <operation> --help
```

Common operation families include test-version creation and submission, test groups and users, invitation codes, public links, tester feedback, detection reports, notifications, stopping a test, and promoting an effective test version.

## Safe workflow

1. Query the existing test versions or groups and identify the exact IDs.
2. Display the selected operation contract and Huawei `sourceUrl`.
3. Preview with `--invoke` and required `--param`, `--query`, `--header`, or `--body` inputs.
4. Verify group type, target users, validity window, download limits, and package/version IDs.
5. Send with `--dry-run=false` only after the specific mutation is authorized.
6. Query the resulting version, group, users, link, feedback, or report state.

Example shape:

```bash
agc --profile staging testing test-api-get-test-grouplist \
  --invoke --query appId="APP_ID" --pretty
```

Use the operation definition rather than copying this example when the installed contract differs.

## Boundaries

- Creating a test version is not the same as uploading its package or submitting it for review.
- Sending tester notifications communicates with third parties; require explicit authorization and verify recipients first.
- Invitation codes and public links can grant access. Do not expose them in public logs or issues.
- Stopping, deleting, or promoting a test version changes live distribution and must not be inferred from a request to inspect it.
- A submitted test version is not necessarily approved, effective, or installable. Poll and report each state separately.
