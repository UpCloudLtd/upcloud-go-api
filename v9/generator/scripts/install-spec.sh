#!/usr/bin/env bash
set -euo pipefail

if [[ $# -ne 4 ]]; then
	printf 'Usage: %s <archive> <member> <version> <sha256>\n' "$0" >&2
	exit 2
fi

ARCHIVE=$1
SPEC_MEMBER=$2
SPEC_VERSION=$3
SPEC_SHA256=$4

if [[ ! -r "$ARCHIVE" ]]; then
	printf 'ERROR: Spec archive is not readable: %s\n' "$ARCHIVE" >&2
	exit 1
fi
if [[ ! "$SPEC_SHA256" =~ ^[[:xdigit:]]{64}$ ]]; then
	printf 'ERROR: spec_sha256 must be 64 hexadecimal characters\n' >&2
	exit 1
fi

STAGING_DIR=$(mktemp -d)
trap 'rm -rf "$STAGING_DIR"; [[ -z "${TEMP_SPEC:-}" ]] || rm -f "$TEMP_SPEC"' EXIT
STAGED_SPEC="$STAGING_DIR/openapi.json"
tar -xOzf "$ARCHIVE" -- "$SPEC_MEMBER" > "$STAGED_SPEC"
if ! printf '%s  %s\n' "$SPEC_SHA256" "$STAGED_SPEC" | sha256sum --check --status; then
	printf 'ERROR: Spec member does not match the supplied SHA-256\n' >&2
	exit 1
fi

ACTUAL_SPEC_VERSION=$(jq -er '.info.version | strings' "$STAGED_SPEC")
if [[ "$ACTUAL_SPEC_VERSION" != "$SPEC_VERSION" ]]; then
	printf 'ERROR: Requested version %s does not match spec version %s\n' "$SPEC_VERSION" "$ACTUAL_SPEC_VERSION" >&2
	exit 1
fi

SCRIPT_DIR=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)
V9_ROOT=$(cd -- "$SCRIPT_DIR/../.." && pwd)
CANONICAL_SPEC="$V9_ROOT/generator/openapi/spec.json"
mkdir -p "$(dirname "$CANONICAL_SPEC")"
TEMP_SPEC=$(mktemp "${CANONICAL_SPEC}.tmp.XXXXXX")
cp -- "$STAGED_SPEC" "$TEMP_SPEC"
mv -f -- "$TEMP_SPEC" "$CANONICAL_SPEC"
TEMP_SPEC=