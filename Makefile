GO ?= go
NPM ?= npm
COVERAGE_MIN ?= 80

GORELEASER ?= go run github.com/goreleaser/goreleaser/v2@v2.17.1

.PHONY: build test coverage coverage-check lint ci web-install web-test web-build release-snapshot release-snapshot-docker

build:
	$(GO) build -o bin/agc ./cmd/agc

test:
	$(GO) test ./...

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

ci: lint coverage-check web-test web-build

release-snapshot:
	$(GORELEASER) release --snapshot --clean --skip=docker

release-snapshot-docker:
	$(GORELEASER) release --snapshot --clean
