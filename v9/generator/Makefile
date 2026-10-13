GO ?= go

export GOTOOLCHAIN := local
export GOWORK := off

.PHONY: test
test:
	$(GO) test -mod=readonly -count=1 -timeout=5m ./tests