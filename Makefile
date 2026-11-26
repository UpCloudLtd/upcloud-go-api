VERSION ?= $(shell git describe --tags --always --dirty --match=v* 2> /dev/null || \
			cat $(CURDIR)/.version 2> /dev/null || echo v0)

.PHONY: test
test: check-test-env
	go test ./... -parallel 8

.PHONY: lint
lint:
	golangci-lint run

check-test-env:
ifndef UPCLOUD_GO_SDK_TEST_USER
	$(error UPCLOUD_GO_SDK_TEST_USER is undefined)
endif
ifndef UPCLOUD_GO_SDK_TEST_PASSWORD
	$(error UPCLOUD_GO_SDK_TEST_PASSWORD is undefined)
endif

.PHONY: version
version:
	@echo $(VERSION)

.PHONY: release-notes
release-notes: CHANGELOG ?= CHANGELOG.md
release-notes: CHANGELOG_HEADER = ^\#\# \[
release-notes: CHANGELOG_VERSION = $(subst v,,$(VERSION))
release-notes:
	@awk -v version='$(CHANGELOG_VERSION)' \
		'index($$0, "## [" version "]") == 1 { flag = 1; next } \
		/${CHANGELOG_HEADER}/ { if ( flag ) { exit; } } \
		flag { lines[++n] = $$0 } \
		END { while (n && lines[n] == "") n--; for (i = 1; i <= n; i++) print lines[i] }' \
		$(CHANGELOG)
