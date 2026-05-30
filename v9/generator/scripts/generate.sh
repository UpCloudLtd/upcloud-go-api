#!/usr/bin/env bash
set -euo pipefail

echo "== Go SDK GENERATION =="

# Preprocess spec to fix enum null values
./generator/scripts/preprocess-spec.sh ./generator/spec.json ./generator/spec.processed.json

# Run oapi-codegen with processed spec
# Note: go generate runs from generator/ directory, so path should be relative to that
SPEC_FILE=./spec.processed.json go generate ./generator/generator.go

# Clean up processed spec
rm -f ./generator/spec.processed.json

echo "Generation complete."