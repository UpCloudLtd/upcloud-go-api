#!/usr/bin/env bash
# Move the ## [Unreleased] entries in the v9 changelog under a release heading.
set -euo pipefail

if [[ $# -ne 2 ]]; then
	printf 'Usage: %s <spec version> <changelog>\n' "$0" >&2
	exit 2
fi

SPEC_VERSION=$1
CHANGELOG=$2

if [[ ! "$SPEC_VERSION" =~ ^0\.0\.[0-9]+$ ]]; then
	printf 'ERROR: version must match 0.0.<patch>, got %s\n' "$SPEC_VERSION" >&2
	exit 1
fi
if ! grep -q '^## \[Unreleased\]$' "$CHANGELOG"; then
	printf 'ERROR: %s has no ## [Unreleased] heading\n' "$CHANGELOG" >&2
	exit 1
fi

# Assumption: tag patch equals spec patch (0.0.118 -> v9.0.0-alpha.0.0.118.0).
release="9.0.0-alpha.0.0.${SPEC_VERSION##*.}.0"
if grep -q "^## \[${release//./\\.}\]" "$CHANGELOG"; then
	printf 'Release %s already in changelog, skip promotion\n' "$release"
	exit 0
fi
empty=0
if ! awk '
	/^## \[Unreleased\]$/ { in_unreleased = 1; next }
	in_unreleased && /^## \[/ { exit }
	in_unreleased && NF && !/^#/ && !/^<!--/ { found = 1; exit }
	END { exit !found }' "$CHANGELOG"; then
	empty=1
fi

tmp=$(mktemp)
trap 'rm -f "$tmp"' EXIT
# Empty Unreleased: drop its leftover headings and add one default entry.
awk -v rel="$release" -v d="$(date -u +%F)" -v spec="$SPEC_VERSION" -v empty="$empty" '
	!done && /^## \[Unreleased\]$/ {
		print; print ""; print "## [" rel "] - " d; print ""
		print "<!-- markdownlint-disable-next-line MD024 -->"
		done = 1
		if (empty) {
			print "### Changed"; print ""
			print "- sdk: regenerate from OpenAPI spec " spec "."; print ""
			drop = 1
		} else {
			skip = 1
		}
		next
	}
	drop && !/^## \[/ { next }
	drop { drop = 0 }
	skip && /^$/ { skip = 0; next }
	{ skip = 0; print }' "$CHANGELOG" > "$tmp"
cat "$tmp" > "$CHANGELOG"
