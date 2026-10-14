GO ?= go
GOLANGCI_LINT ?= golangci-lint
TEST_TIMEOUT ?= 5m
API_BASE_REF ?=
API_REPORT_DIR ?=

export GOTOOLCHAIN := local
export GOWORK := off

.PHONY: build test vet lint test-consumer test-credentials test-generator generate-request-adapters check-request-adapters test-request-conformance test-public-api test-public-api-gate update-public-api check

build:
	$(GO) build -mod=readonly ./...
	$(GO) test -mod=readonly -run '^$$' -count=1 -timeout=$(TEST_TIMEOUT) ./...

test:
	$(GO) test -mod=readonly -count=1 -timeout=$(TEST_TIMEOUT) ./...

vet:
	$(GO) vet -mod=readonly ./...

lint:
	$(GOLANGCI_LINT) run --config ../.golangci.yaml --timeout=$(TEST_TIMEOUT)

test-consumer:
	bash test/scripts/check-consumer.sh

test-credentials:
	cd ../credentials && $(GO) test -mod=readonly -count=1 -timeout=$(TEST_TIMEOUT) ./...

test-generator:
	$(GO) test -mod=readonly -count=1 -timeout=$(TEST_TIMEOUT) ./generator/tests

generate-request-adapters:
	$(GO) run -mod=readonly ./test/requestconformance/generate

check-request-adapters:
	$(GO) run -mod=readonly ./test/requestconformance/generate -check

test-request-conformance: check-request-adapters
	$(GO) test -mod=readonly -count=1 -timeout=$(TEST_TIMEOUT) ./test/requestconformance

test-public-api:
	bash test/scripts/check-public-api.sh check $(if $(API_BASE_REF),--base-ref "$(API_BASE_REF)") $(if $(API_REPORT_DIR),--report-dir "$(API_REPORT_DIR)")

test-public-api-gate:
	bash test/scripts/check-public-api-test.sh

update-public-api:
	bash test/scripts/check-public-api.sh update

check: build test vet lint test-consumer test-credentials test-generator test-request-conformance test-public-api test-public-api-gate