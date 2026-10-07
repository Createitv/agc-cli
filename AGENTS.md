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
