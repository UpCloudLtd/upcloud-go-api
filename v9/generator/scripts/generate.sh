#!/usr/bin/env bash
set -euo pipefail

echo "== Go SDK GENERATION =="

SPEC_FILE=./generator/openapi/spec.json
if [[ ! -f "$SPEC_FILE" || ! -r "$SPEC_FILE" || ! -s "$SPEC_FILE" ]]; then
  echo "ERROR: Missing, unreadable or empty OpenAPI input: $SPEC_FILE" >&2
  exit 1
fi
PROVENANCE_FILE=./generator/openapi/provenance.json
if [[ ! -f "$PROVENANCE_FILE" || ! -r "$PROVENANCE_FILE" || ! -s "$PROVENANCE_FILE" ]]; then
  echo "ERROR: Missing, unreadable or empty spec provenance: $PROVENANCE_FILE" >&2
  exit 1
fi
if ! EXPECTED_SHA256=$(jq -er '.sha256 | select(type == "string" and test("^[0-9a-fA-F]{64}$"))' "$PROVENANCE_FILE"); then
  echo "ERROR: Invalid spec SHA-256 in $PROVENANCE_FILE" >&2
  exit 1
fi
if ! EXPECTED_VERSION=$(jq -er '.spec_version | select(type == "string" and length > 0)' "$PROVENANCE_FILE"); then
  echo "ERROR: Invalid spec version in $PROVENANCE_FILE" >&2
  exit 1
fi
ACTUAL_SHA256=$(sha256sum -- "$SPEC_FILE" | awk '{print $1}')
if [[ "${EXPECTED_SHA256,,}" != "$ACTUAL_SHA256" ]]; then
  echo "ERROR: OpenAPI input does not match SHA-256 in $PROVENANCE_FILE" >&2
  exit 1
fi
if ! jq -s -e -f ./generator/scripts/validate-spec.jq "$SPEC_FILE" > /dev/null; then
  echo "ERROR: Invalid OpenAPI input or unsupported tag layout: $SPEC_FILE" >&2
  exit 1
fi
ACTUAL_VERSION=$(jq -er '.info.version | strings' "$SPEC_FILE")
if [[ "$EXPECTED_VERSION" != "$ACTUAL_VERSION" ]]; then
  echo "ERROR: OpenAPI input version does not match $PROVENANCE_FILE" >&2
  exit 1
fi

OUTPUT_DIR=$(mktemp -d)
TAG_SPEC="$OUTPUT_DIR/tag-spec.json"
trap 'rm -rf "$OUTPUT_DIR"' EXIT

go tool oapi-codegen -config ./generator/config/server-urls.yaml -o "$OUTPUT_DIR/server_urls.gen.go" "$SPEC_FILE"

# schemas_to_exclude maps a tag name to a comma-separated list of schema names that should
# not be generated for that tag because they are already defined in another tag's model file.
# Use this when the same schema is reachable from multiple tags but should only be declared once.
declare -A schemas_to_exclude=(
  # networkPeeringName and networkPeeringState are referenced by networkDetails (network tag)
  # but are owned by the network-peering tag.
  ["network"]="networkPeeringName,networkPeeringState"
)

for tag in $(jq -r '.tags[]?.name' "$SPEC_FILE" | sort); do
  echo "Generating $tag"
  jq --arg tag "$tag" -f ./generator/scripts/filter-spec-by-tag.jq "$SPEC_FILE" > "$TAG_SPEC"

  exclude_flag=""
  if [[ -n "${schemas_to_exclude[$tag]+set}" ]]; then
    exclude_flag="-exclude-schemas ${schemas_to_exclude[$tag]}"
  fi

  # shellcheck disable=SC2086
  go tool oapi-codegen -config ./generator/config/models.yaml $exclude_flag -include-tags "$tag" -o "$OUTPUT_DIR/${tag}_models.gen.go" "$TAG_SPEC"
  go tool oapi-codegen -config ./generator/config/client.yaml -include-tags "$tag" -o "$OUTPUT_DIR/${tag}_client.gen.go" "$TAG_SPEC"
done

shopt -s nullglob
existing_files=(./pkg/upcloud/server_urls.gen.go ./pkg/upcloud/*_models.gen.go ./pkg/upcloud/*_client.gen.go)
rm -f -- "${existing_files[@]}"
mv -- "$OUTPUT_DIR/"*.gen.go ./pkg/upcloud/

echo "Generation complete."