---
name: agc-auth
description: Authenticate agc-cli, manage credential profiles, identify the active Huawei AppGallery Connect account, and diagnose authorization failures. Use for login, profile switching, auth checks, whoami, OAuth team lookup, or CI credential setup.
---

# agc auth

Use `agc auth` for local credential profiles and identity inspection. Never ask for or expose client keys, service-account private keys, access tokens, credential files, or raw authentication responses.

## Inspect before changing credentials

```bash
agc auth list --output table
agc auth check
agc auth whoami --pretty
agc --profile production auth check
```

`auth check` validates local configuration unless `--remote` is explicitly supplied. `auth whoami` reports the selected profile, non-secret client or service-account identity, and project context. It does not invent a personal developer name from credentials that do not contain one.

With a user-authorized Huawei OAuth token, query official team accounts without saving the token:

```bash
AGC_OAUTH2_TOKEN="$AGC_OAUTH2_TOKEN" agc auth whoami --teams --pretty
```

Treat the token as a secret. Do not echo it, put it in command examples with a literal value, persist it in project files, or include it in logs.

## Login and profile selection

Inspect the installed command first:

```bash
agc auth login --help
```

Then use a credential file or protected environment values already available on the machine:

```bash
agc auth login --service-account-file /protected/agc-service-account.json --name production
agc auth login --api-client-file /protected/agc-api-client.json --name staging
agc auth login --client-id "$AGC_CLIENT_ID" --client-key "$AGC_CLIENT_KEY" --name ci
```

Use `--profile NAME` on unattended commands. Project context may provide a default profile, but an explicit profile is easier to audit in multi-environment CI.

## Remote checks

`agc auth check --remote` contacts Huawei and may perform an app or project read. Run it only when the user has authorized the remote check and the intended profile/project is explicit. Token creation is not proof that the profile can access a particular app.

For 401/403 failures, report the selected profile and redacted command, then check credential mode, project/app ownership, required Huawei roles, and endpoint-specific headers. Never include secrets in the report.
