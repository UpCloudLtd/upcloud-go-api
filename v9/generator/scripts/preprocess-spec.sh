#!/usr/bin/env bash
set -euo pipefail

# Remove null from enum arrays while setting nullable: true
# This fixes the issue where oapi-codegen generates invalid Go code like `<nil>` for null enum values

SPEC_IN="$1"
SPEC_OUT="${2:-${SPEC_IN}.processed}"

jq '
  # Walk through all objects in the spec
  walk(
    if type == "object" and .enum? and (.enum | type) == "array" then
      # If the enum contains null
      if (.enum | any(. == null)) then
        # Set nullable to true and remove null from enum
        .nullable = true |
        .enum = (.enum | map(select(. != null)))
      else
        .
      end
    else
      .
    end
  )
' "$SPEC_IN" > "$SPEC_OUT"

echo "Preprocessed spec written to: $SPEC_OUT"
