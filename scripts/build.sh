#!/usr/bin/env bash

#
# Stop build script if there are any errors
#
set -x
set -e

PROJECT_HOME=$(git rev-parse --show-toplevel)
if [[ -z "$GOBIN" ]]; then
  GOBIN="$(go env GOPATH)/bin"
fi

# Get build dependencies
pushd /tmp
if [[ ! -e $GOBIN/staticcheck ]]; then
  go get honnef.co/go/tools/cmd/staticcheck
fi
if [[ ! -e $GOBIN/goimports ]]; then
  go get golang.org/x/tools/cmd/goimports
fi
popd

# Go Mod Download
go mod download

# Check format of Go files
$PROJECT_HOME/scripts/check-gofmt.sh .

# Check format of Imports
$PROJECT_HOME/scripts/check-goimports.sh .

# Check module tidiness
#$PROJECT_HOME/scripts/check-gomodtidy.sh

# Check for common problems
go vet ./...

# Lint
staticcheck ./...

# Run tests
go test ./... -v -parallel 1 -timeout 60m

