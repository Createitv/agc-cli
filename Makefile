GO ?= go
NPM ?= npm
COVERAGE_MIN ?= 80

GORELEASER ?= go run github.com/goreleaser/goreleaser/v2@v2.17.1

.PHONY: check-docs build test test-scripts test-api-transport coverage coverage-check lint ci web-install web-test web-build release-snapshot release-snapshot-docker

build:
	$(GO) build -o bin/agc ./cmd/agc

test:
	$(GO) test ./...

# 离线测试：不读取真实华为授权，不调用华为服务器。
test-scripts:
	python3 scripts/test_coverage_isolation.py
	python3 scripts/test_ci_report.py
	python3 scripts/test_api_suite.py
	python3 scripts/test_live_lifecycle.py
	python3 scripts/test_live_package.py
	python3 scripts/test_certificate_fingerprint.py
	python3 scripts/test_live_multipart.py
	python3 scripts/test_goal_progress.py
	python3 scripts/test_publishing_audit.py
	python3 scripts/test_publishing_models.py

# 每个注册接口验证真实本地 HTTP 收发，8 路并发。
test-api-transport:
	$(GO) test ./cmd/agc/command -run TestEveryEndpointHTTPTransport -count=1 -parallel=8

COVERAGE_RUNNER := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))scripts/check-coverage.py

coverage:
	python3 "$(COVERAGE_RUNNER)" --go "$(GO)"

coverage-check:
	python3 "$(COVERAGE_RUNNER)" --go "$(GO)" --minimum "$(COVERAGE_MIN)"

lint:
	$(GO) vet ./...

web-install:
	$(NPM) --prefix apps/web install

web-test:
	$(NPM) --prefix apps/web test -- --run
	node --test scripts/test_browser_capture.mjs

web-build:
	$(NPM) --prefix apps/web run build

check-docs:
	python3 scripts/check-docs.py

ci: check-docs test-scripts lint coverage-check web-test web-build

release-snapshot:
	$(GORELEASER) release --snapshot --clean --skip=docker

release-snapshot-docker:
	$(GORELEASER) release --snapshot --clean
