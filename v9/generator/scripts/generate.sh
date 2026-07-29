#!/usr/bin/env bash
set -euo pipefail

echo "== Go SDK GENERATION =="

# Preprocess spec to fix enum null values
./generator/scripts/preprocess-spec.sh ./generator/spec.json ./generator/spec.processed.json

# Run oapi-codegen with processed spec
# Note: go generate runs from generator/ directory, so path should be relative to that
SPEC_FILE=./generator/spec.processed.json

# Generate common server URL file
go tool oapi-codegen -config ./generator/config/server-urls.yaml -o ./pkg/upcloud/server_urls.gen.go "$SPEC_FILE"

# Loop through all tags in the spec and generate models and client for each tag
for tag in $(jq -r '.tags[]?.name' "$SPEC_FILE"); do
  echo "Generating $tag"
  go tool oapi-codegen -config ./generator/config/models.yaml -include-tags "$tag" -o ./pkg/upcloud/"$tag"_models.gen.go "$SPEC_FILE"
  go tool oapi-codegen -config ./generator/config/client.yaml -include-tags "$tag" -o ./pkg/upcloud/"$tag"_client.gen.go "$SPEC_FILE"
done

# Clean up processed spec
rm $SPEC_FILE

echo "Generation complete."