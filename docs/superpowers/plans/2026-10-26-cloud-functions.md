# Cloud Functions Implementation Plan

> **For agentic workers:** Use superpowers:executing-plans inline. User authorized implementation and opted out of routine questions; preserve all pre-existing dirty changes. No commits or deployment of business resources for tests.

**Goal:** Add native Go cloud-function create/upload, code/config update, publish, list/query/version and invoke commands, with discoverable local REST/OpenAPI.

**Architecture:** Domain requests and source-derived endpoint metadata live in pkg/domain; HTTP multipart, digest, OAuth headers and revision workflows live in pkg/cloudfunctions; CLI and REST translate shared request types. Source: Huawei-published @hw-agconnect/agc-cli 1.12.0, independently inspected with file hashes. Distinguish vendor package source review from a public REST specification or live verification.

**Tech Stack:** Go standard library, existing Cobra and JSON envelopes; no dependency on a Node subprocess at runtime.

## Global Constraints

Keep native agc dry-run default, JSON stdout, local API authentication and project context. Never automatically forward Service Account credentials to an OAuth-only endpoint. Session file explicitly supplies user OAuth uid/teamId/productId/accessToken, with expiry check and project mismatch rejection. Never print session values or function environment values. Preserve TLS verification and refuse cross-origin redirects for authenticated requests.

## Review Focus

- ZIP upload uses actual bytes and correct digest, length and multipart fields.
- Code/config update obtains revision and preserves environment, layers and handler; concurrent modification errors stop without blind retry.
- Failed update must never publish. Publish failure must report that update may already have succeeded.
- Dry-run never reads authorization or contacts Huawei; invalid local input fails before auth/network.
- Invoke payload is arbitrary valid JSON; headers such as function execution status must remain observable independently of HTTP success.

## Tasks

1. Source and domain contracts: write failing tests for functions discovery, docs provenance and OpenAPI. Capture reproducible package-source evidence; register 9 operations and shared input types.
2. Transport: write httptest failures for digests, read-before-write, config preservation, OAuth isolation, payload fidelity, conflict/no-publish, failed invocation and dry-run. Implement native client with bounded ZIP streaming/temp files, revision workflows and redacted responses.
3. CLI and REST: failing tests for official-style flags, session-file validation, no automatic credentials, structured errors, arbitrary JSON payload and shared local API behavior. Implement adapters and metadata schemas.
4. Verify and document: source-backed parameter cases, full Go suite/race, make ci; build bin/agc and read back real binary help/discovery/dry-run. Record any existing suite failures separately. No cloud write claims without dedicated-resource live readback.
