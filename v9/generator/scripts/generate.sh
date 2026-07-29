#!/usr/bin/env bash
set -euo pipefail

echo "== Go SDK GENERATION =="

# Remove existing generated files
rm ./pkg/upcloud/server_urls.gen.go || true
rm ./pkg/upcloud/*_models.gen.go || true
rm ./pkg/upcloud/*_client.gen.go || true

# Preprocess spec to fix enum null values
./generator/scripts/preprocess-spec.sh ./generator/spec.json ./generator/spec.processed.json

# Run oapi-codegen with processed spec
# Note: go generate runs from generator/ directory, so path should be relative to that
SPEC_FILE=./generator/spec.processed.json
TAG_SPEC=
trap 'rm -f "$SPEC_FILE" "${TAG_SPEC:-}"' EXIT

# Generate common server URL file
go tool oapi-codegen -config ./generator/config/server-urls.yaml -o ./pkg/upcloud/server_urls.gen.go "$SPEC_FILE"

# Loop through all tags in the spec and generate models and client for each tag
for tag in $(jq -r '.tags[]?.name' "$SPEC_FILE"); do
  echo "Generating $tag"
  TAG_SPEC=$(mktemp)
  jq --arg tag "$tag" -f ./generator/scripts/filter-spec-by-tag.jq "$SPEC_FILE" > "$TAG_SPEC"

  go tool oapi-codegen -config ./generator/config/models.yaml -include-tags "$tag" -o ./pkg/upcloud/"$tag"_models.gen.go "$TAG_SPEC"
  go tool oapi-codegen -config ./generator/config/client.yaml -include-tags "$tag" -o ./pkg/upcloud/"$tag"_client.gen.go "$TAG_SPEC"

  rm -f "$TAG_SPEC"
  TAG_SPEC=
done

# Clean up processed spec
rm -f "$SPEC_FILE"

# De-duplicate any duplicate types/constants across tag-filtered model files
echo "De-duplicating generated types and constants"
go run ./generator/scripts/dedup.go

echo "Generation complete."