#!/usr/bin/env bash
set -euo pipefail

echo "== Go SDK GENERATION =="

# Run oapi-codegen
go generate ./generator/generator.go

echo "Generation complete."