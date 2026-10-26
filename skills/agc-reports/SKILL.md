---
name: agc-reports
description: Discover and download AppGallery Connect distribution, revenue, install, retention, and HarmonyOS service reports with agc-cli. Use when a user asks for AGC analytics exports, report URLs, or repeatable report collection in CI.
---

# AppGallery Connect reports

Report endpoints differ in app ID placement, date range, file type, and returned data. Always inspect the current operation:

```bash
agc reports endpoints --output table
agc reports <operation> --pretty
agc reports <operation> --help
```

Available exports may cover downloads, install failures, new and retained users, paid downloads, IAP, reservations, coupon or award activity, and HarmonyOS service distribution, users, or widgets.

## Download workflow

1. Select an explicit profile and app/project.
2. Inspect the endpoint contract and official Huawei `sourceUrl` for supported date and format values.
3. Preview the request with `--invoke`.
4. Send the authorized read request with `--dry-run=false`.
5. Distinguish an API response, a returned download URL, and a downloaded report file.
6. Save raw bytes with `--out FILE` only when the file is an intended local artifact.

Example shape:

```bash
agc --profile reporting reports appdownloadexport \
  --invoke --param appId="APP_ID" --query startTime="START" \
  --query endTime="END" --pretty
```

Read the installed operation definition before relying on example parameter names.

## Data handling

Reports may contain business-sensitive metrics or user-related information. Do not paste raw reports into issues, commits, or public logs. Report the requested aggregation and source period, and state whether the result is an API response, a generated link, or an actually downloaded and parsed file.
