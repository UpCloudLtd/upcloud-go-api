#!/usr/bin/env bash
set -euo pipefail

echo "== Go SDK GENERATION =="

# Remove existing generated files
rm ./pkg/upcloud/server_urls.gen.go || true
rm ./pkg/upcloud/*_models.gen.go || true
rm ./pkg/upcloud/*_client.gen.go || true

SPEC_FILE=./generator/spec.json
TAG_SPEC=
trap 'rm -f "${TAG_SPEC:-}"' EXIT

# Generate common server URL file
go tool oapi-codegen -config ./generator/config/server-urls.yaml -o ./pkg/upcloud/server_urls.gen.go "$SPEC_FILE"

# schemas_to_exclude maps a tag name to a comma-separated list of schema names that should
# not be generated for that tag because they are already defined in another tag's model file.
# Use this when the same schema is reachable from multiple tags but should only be declared once.
declare -A schemas_to_exclude=(
  # networkPeeringName and networkPeeringState are referenced by networkDetails (network tag)
  # but are owned by the network-peering tag.
  ["network"]="networkPeeringName,networkPeeringState"
)

# Loop through all tags in the spec and generate models and client for each tag
for tag in $(jq -r '.tags[]?.name' "$SPEC_FILE" | sort); do
  echo "Generating $tag"
  TAG_SPEC=$(mktemp)
  jq --arg tag "$tag" -f ./generator/scripts/filter-spec-by-tag.jq "$SPEC_FILE" > "$TAG_SPEC"

  exclude_flag=""
  if [[ -n "${schemas_to_exclude[$tag]+set}" ]]; then
    exclude_flag="-exclude-schemas ${schemas_to_exclude[$tag]}"
  fi

  # shellcheck disable=SC2086
  go tool oapi-codegen -config ./generator/config/models.yaml $exclude_flag -include-tags "$tag" -o ./pkg/upcloud/"$tag"_models.gen.go "$TAG_SPEC"
  go tool oapi-codegen -config ./generator/config/client.yaml -include-tags "$tag" -o ./pkg/upcloud/"$tag"_client.gen.go "$TAG_SPEC"

  rm -f "$TAG_SPEC"
  TAG_SPEC=
done

echo "Generation complete."