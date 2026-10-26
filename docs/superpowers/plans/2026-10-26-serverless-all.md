# Serverless full integration and verification

Scope: all HTTP operations exposed by the six official Serverless navigation groups, plus the nine native function workflows already specified in cloud-functions.md. SDK class methods are not automatically public HTTP management contracts. Template callback operations remain inbound callbacks. No implicit live credentials or business-resource writes.

1. Capture endpoint bodies, nested models, examples, errors and authorization prerequisites through Computer Use; preserve capture tables and paragraph URLs.
2. Add source-based failing tests for discovery, nested requiredness, auth alternatives, encoded JSON, binary/signed uploads, callbacks and native function revision workflows.
3. Register domain contracts and shared serialization/validation; expose CLI, REST, OpenAPI and --docs.
4. Verify every operation through isolated local HTTP transports and negative input tests; build and inspect workspace binary; make ci and race suite.
5. Live operations require an explicitly selected profile/project and dedicated resources. Record each operation as passed/failed/not-run with reason. Never equate simulation with live verification.

## Ledger

- Initial: existing dirty checkout retained; pre-change status in /tmp/agc-serverless-before-status.txt. Existing TestCloudFunctionsDiscoveryAndNativeContract fails because native functions are missing.
- Computer Use: 65 official rendered pages captured, including 46 REST/callback endpoint documents, request/response models and status-code appendices. Full capture saved locally via TextEdit, not by direct HTTP scraping.
- Ruling: implement native Go function transport from Huawei AGCLI 1.12.0 source, preserving provenance distinct from public REST specifications. Do not reproduce its disabled TLS validation.
- Live selection: question pending for explicit profile/project and dedicated test resources; independent offline work continues.
- Implemented: 46 documented HTTP/callback operations and 9 native function workflows, with 71 browser-read endpoint/model/prerequisite pages and hashed vendor source evidence. Scope boundaries are explicit in docs/SERVERLESS.md; private console-only resource lifecycles are not advertised as complete.
- Independent final review: reproduced and fixed JSON-file download handling, arbitrary function results, CPU allocation, description-only updates, signed auth requiredness, native error envelopes, resumable/digest checks and response overflow. Each finding has a regression test.
- CI observation: make ci stops in the untouched scripts/test_goal_progress.py test_domain_precheck_requires_real_code_file_readback_and_cleanup. Continued gates initially found missing Web dependencies; npm ci restored lockfile dependencies. Final checks and operation matrix follow.
- Final evidence supersedes earlier transient failures: make ci passes; full race suite passes with CGO_ENABLED=1 GOARCH=arm64 on this Apple Silicon host. The default installed Intel Go toolchain caused Rosetta startup stalls in earlier runs. All 55 CLI, REST and dry-run operations passed; report generation requires both entire packages to pass and never marks cloud validation from offline events.
- Independent binary readback: native arm64 bin/agc, all 55 endpoint discoveries and --docs parsed, dry-run with nonexistent session succeeded without reading authorization. /usr/local/bin/agc remains the pre-existing v0.1.0 installation; no package publication or system upgrade performed.
- Live remains not-run pending explicit profile/project/resources and function user OAuth session selection. No commits, pushes, cloud writes or business-resource changes.
