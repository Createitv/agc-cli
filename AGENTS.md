# Codex repository instructions

- Do not append Co-Authored-By trailers to commits.
- Preserve existing uncommitted work. Stage explicit paths.
- Keep domain contracts separate from HTTP, credential storage and CLI presentation.
- Reproduce a bug with a meaningful regression test before fixing it.
- Run `make ci` and `go test -race ./...` before proposing a merge. Never claim cloud checks passed without inspecting the GitHub Actions run for the exact commit.
- Never commit credentials, tokens, private keys or raw credential responses. Offline tests must isolate user authorization and must not call Huawei.
- Live tests must use explicitly selected profiles and dedicated test resources. Do not publish applications or alter existing business resources as a smoke test.
- Update API contracts, OpenAPI documentation and tests together when parameters change.
- Include validation evidence and remaining limitations in PR descriptions.
