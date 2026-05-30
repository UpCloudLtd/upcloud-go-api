#!/usr/bin/env bash
set -euo pipefail

# Use SPEC_FILE if set, otherwise default to spec.json
# Paths are relative to generator/ directory
SPEC="${SPEC_FILE:-./spec.json}"

go tool oapi-codegen -config ./config/upcloud.yaml "$SPEC"
