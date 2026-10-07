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
	python3 scripts/test_ci_report.py
	python3 scripts/test_api_suite.py
	python3 scripts/test_live_lifecycle.py

# 每个注册接口验证真实本地 HTTP 收发，8 路并发。
test-api-transport:
	$(GO) test ./cmd/agc/command -run TestEveryEndpointHTTPTransport -count=1 -parallel=8

coverage:
	$(GO) test ./... -coverprofile=coverage.out
	$(GO) tool cover -func=coverage.out

coverage-check: coverage
	@total=$$($(GO) tool cover -func=coverage.out | awk '/^total:/ { gsub("%","",$$3); print $$3 }'); \
	awk -v total="$$total" -v min="$(COVERAGE_MIN)" 'BEGIN { if (total + 0 < min + 0) { printf("coverage %.1f%% is below %.1f%%\n", total, min); exit 1 } printf("coverage %.1f%% >= %.1f%%\n", total, min) }'

lint:
	$(GO) vet ./...

web-install:
	$(NPM) --prefix apps/web install

web-test:
	$(NPM) --prefix apps/web test -- --run

web-build:
	$(NPM) --prefix apps/web run build

check-docs:
	python3 scripts/check-docs.py

ci: check-docs test-scripts lint coverage-check web-test web-build

release-snapshot:
	$(GORELEASER) release --snapshot --clean --skip=docker

release-snapshot-docker:
	$(GORELEASER) release --snapshot --clean
