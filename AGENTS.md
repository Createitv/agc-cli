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

## Official API rules and Agent usability

- Review the actual official page for each affected endpoint, using the browser when requested. Follow linked request/response models and relevant appendices; an endpoint list or a URL alone does not establish a verified contract.
- Preserve browser-read documentation in `pkg/domain/reference/`, including parameter tables, notes, examples, source paragraph URLs and the review date. Keep the capture process reproducible with `scripts/browser-capture-reference.mjs`; preserve merged table cells so device-specific rules retain their context.
- Expose reviewed documentation through endpoint discovery and `agc <family> <endpoint> --docs --pretty`. Include relevant supplemental models and appendices so an Agent can resolve prerequisites without guessing.
- Encode verified parameter locations, types, requiredness, nested objects, array item types, enums and explicit bounds in domain contracts. Update CLI serialization, REST validation, OpenAPI and independent regression tests together.
- Follow nested model links recursively. Validate child fields with their full paths (for example, `body.languages[0].appDesc`), including conditional child requirements and per-item limits. A valid parent object or array does not establish that its contents satisfy the API contract.
- Preserve the documented wire type when a string contains encoded JSON. Validate its decoded structure locally, but serialize it as a string when the API requires one. Encode form values exactly once, preserve integer precision and reject unsupported repeated scalar parameters.
- When official tables, examples or notes conflict, retain the source evidence and identify the unresolved rule. Do not silently choose an array/object shape, wrapper, required field or mutually exclusive branch. Separate documented requirements from locally verified behavior and implementation limitations.
- Distinguish unconditional requirements from conditional requirements and alternative authentication or application-identification branches. Do not require every authentication header at once. Do not interpret `Integer(32)` as a numeric range, dynamic map examples such as `additionalPropX` as literal required keys, or logical request-model names as JSON wrapper properties without checking the actual request example.
- Preserve transport semantics from the official examples: JSON, form data, multipart files and signed binary uploads require different handling. Respect returned upload methods, signed headers, URL expiry, file length and upload/parse prerequisites.
- Treat `inbound-callback` endpoints as callbacks to the developer's server. Do not automatically load Huawei credentials or forward an AGC Bearer token to that server. Preserve the documented form/JSON encoding and signatures; a simulated callback transport test does not prove Huawei signature verification or a real game transaction succeeded.
- Preserve ordered repeated query keys when the reference shows array query parameters. CLI callers repeat `--query key=value`; REST callers can use `queryValues` arrays alongside scalar `query` entries. Keep paired filters and values aligned, validate documented date formats and interval limits, and keep OpenAPI required-field alternatives consistent with both input forms.
- File rules depend on asset kind, language, device and application kind. Check the official file appendix and linked `LangFileInfo`/`FileInfo` models; a top-level array limit does not replace per-device asset counts, formats, sizes, resolutions or video/poster ordering. Do not guess missing application context or pretend a remote object ID proves local file properties.
- Reject invalid inputs before authentication or network calls when the rule can be checked locally. Return machine-readable CLI/REST errors with a stable code, endpoint, affected field path, source URL, concrete correction suggestion and documentation command. Nested violations must identify the array index and child field. Keep JSON stdout parseable and never include credentials or raw credential responses in feedback.
- Separate input violations, authentication failures, transport failures and upstream business errors. Preserve upstream error codes. Guide the Agent to the documented role, application state or prerequisite instead of suggesting blind retries or silently changing business resources.
- Test meaningful documented examples and failure cases before fixing behavior. Use independent source-based assertions for requiredness, encoding and nested shapes; registry-generated smoke fixtures alone cannot prove the registry is correct. Offline tests must not use the user's Huawei authorization.
- Report documentation review, implemented rules, offline verification and live verification accurately. Do not claim every official rule or real-app operation is supported merely because documentation was captured or generic transport tests passed. Continue toward a usable CLI rather than producing a report in place of the requested implementation.
- Rebuild the CLI after source changes and verify the binary actually invoked by the Agent. Distinguish the workspace binary from an installed package; a successful source test does not prove that an older installed CLI contains the fix.

## Version releases and package managers

- Ordinary branch pushes and PR updates run CI. To publish a stable version, push a new `vMAJOR.MINOR.PATCH` tag on the intended, tested commit from the merged main branch. Check existing tags first; do not move or reuse a published tag.
- Tag pushes trigger `.github/workflows/release.yml`. Release requires publisher access checks and the reusable CI workflow to pass before GoReleaser publishes packages.
- GoReleaser automatically updates `Createitv/homebrew-tap` on `master` (`agc-cli.rb`) and `Createitv/scoop-bucket` on `main` (`agc-cli.json`). Maintain these targets in `.goreleaser.yaml` and `scripts/verify-package-publish.py` together. Stable package channels must not be overwritten by prereleases.
- Cross-repository publishing uses the GitHub Actions secret `TAP_GITHUB_TOKEN`; the workflow's default `GITHUB_TOKEN` handles this repository's Release and GHCR. Never put either token in source files, documentation, reports or logs.
- Preserve the pinned GoReleaser version and existing Homebrew Formula compatibility. Validate tool upgrades by generating a release snapshot and inspecting both package manifests before changing the publishing toolchain.
- To run Release manually, select an existing version tag, for example `gh workflow run release.yml --ref vX.Y.Z`. Do not rerun the full publisher for an already published version merely to check its status.
- Use `.github/workflows/publisher-check.yml` for read-only checks of the saved publisher credentials and the latest stable package manifests. It does not publish a version.
- Before reporting a release complete, inspect the exact tag/SHA's Actions run, confirm Release assets and `validation-reports.zip`, and verify both package manifests' versions, download URLs and SHA256 against that Release's `checksums.txt`. A successful build, tag push or credential check alone is not a completed release.
- Updating package manifests does not upgrade installations on users' computers. Distinguish package publication, manifest readback and actual installation/upgrade tests when reporting results.
- See `docs/CLOUD_TESTING.md` for the cloud validation and release workflow.
