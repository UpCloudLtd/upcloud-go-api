#!/usr/bin/env bash
# Print the newest release version in a changelog, e.g. 9.0.0-alpha.0.0.118.0.
# Print nothing if there is no release heading.
set -euo pipefail

if [[ $# -ne 1 ]]; then
	printf 'Usage: %s <changelog>\n' "$0" >&2
	exit 2
fi

version=$(awk '
	/^## \[Unreleased\]/ { next }
	/^## \[[^]]+\]/ { sub(/^## \[/, ""); sub(/\].*/, ""); print; exit }' "$1")
if [[ -n "$version" && ! "$version" =~ ^[0-9]+\.[0-9]+\.[0-9]+(-[0-9A-Za-z.]+)?$ ]]; then
	printf 'ERROR: bad release version in %s: %s\n' "$1" "$version" >&2
	exit 1
fi
printf '%s' "$version"
