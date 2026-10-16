#!/usr/bin/env bash
set -euo pipefail
export GOWORK=off GOTOOLCHAIN=local

SDK_ROOT=$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)
CONSUMER_DIR=$(mktemp -d)
trap 'rm -rf "$CONSUMER_DIR"' EXIT
cp "$SDK_ROOT/test/consumer/go.mod" "$SDK_ROOT/test/consumer/consumer_test.go" "$CONSUMER_DIR/"
cp "$SDK_ROOT/go.sum" "$CONSUMER_DIR/go.sum"
cd "$CONSUMER_DIR"
go mod edit -replace "github.com/UpCloudLtd/upcloud-go-api/v9=$SDK_ROOT"
go mod tidy
go test -mod=readonly -count=1 -timeout=2m ./...