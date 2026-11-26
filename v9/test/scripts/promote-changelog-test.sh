#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)
SCRIPT="$SCRIPT_DIR/../../generator/scripts/promote-changelog.sh"
TEST_ROOT=$(mktemp -d)
trap 'rm -rf "$TEST_ROOT"' EXIT
LOG="$TEST_ROOT/CHANGELOG.md"
TODAY=$(date -u +%F)

fail() {
	printf 'FAIL: %s\n' "$1" >&2
	cat "$LOG" >&2
	exit 1
}

write_log() {
	printf '# Changelog\n\n## [Unreleased]\n\n%s## [9.0.0-alpha.0.0.108.0] - 2026-10-05\n\n- old\n' "$1" > "$LOG"
}

# Promotes pending entries and leaves an empty Unreleased section.
write_log $'### Changed\n\n- new\n\n'
bash "$SCRIPT" 0.0.118 "$LOG" > /dev/null
expected="# Changelog

## [Unreleased]

## [9.0.0-alpha.0.0.118.0] - $TODAY

<!-- markdownlint-disable-next-line MD024 -->
### Changed

- new

## [9.0.0-alpha.0.0.108.0] - 2026-10-05

- old"
[[ "$(cat "$LOG")" == "$expected" ]] || fail 'promote 0.0.118'

# Skips a release that already has a dated heading.
write_log $'- new\n\n'
before=$(cat "$LOG")
bash "$SCRIPT" 0.0.108 "$LOG" | grep -q 'skip promotion' || fail 'skip 0.0.108'
[[ "$(cat "$LOG")" == "$before" ]] || fail 'skip 0.0.108 changed the file'

# Adds a default entry when Unreleased has no entries.
default="# Changelog

## [Unreleased]

## [9.0.0-alpha.0.0.118.0] - $TODAY

<!-- markdownlint-disable-next-line MD024 -->
### Changed

- sdk: regenerate from OpenAPI spec 0.0.118.

## [9.0.0-alpha.0.0.108.0] - 2026-10-05

- old"
for body in '' $'### Changed\n\n'; do
	write_log "$body"
	bash "$SCRIPT" 0.0.118 "$LOG" > /dev/null
	[[ "$(cat "$LOG")" == "$default" ]] || fail "default entry for body '$body'"
done

# Rejects spec versions that do not map to a v9 alpha tag.
write_log $'- new\n\n'
if bash "$SCRIPT" 1.2.3 "$LOG" 2> /dev/null; then fail 'version 1.2.3'; fi

# changelog-version.sh prints the newest release heading.
VERSION_SCRIPT="$SCRIPT_DIR/../../generator/scripts/changelog-version.sh"
write_log $'- new\n\n'
[[ "$(bash "$VERSION_SCRIPT" "$LOG")" == 9.0.0-alpha.0.0.108.0 ]] || fail 'version 108'
printf '# Changelog\n\n## [Unreleased]\n' > "$LOG"
[[ -z "$(bash "$VERSION_SCRIPT" "$LOG")" ]] || fail 'no release heading'
printf '## [9.0.0; rm -rf /]\n' > "$LOG"
if bash "$VERSION_SCRIPT" "$LOG" 2> /dev/null; then fail 'bad version'; fi
printf '## [9.0.0\n' > "$LOG"
[[ -z "$(bash "$VERSION_SCRIPT" "$LOG")" ]] || fail 'heading without ]'

# release-notes in the root Makefile matches the full version only.
write_log $'- new\n\n'
notes() { make -s -C "$SCRIPT_DIR/../../.." release-notes CHANGELOG="$LOG" VERSION="$1"; }
[[ "$(notes v9.0.0-alpha.0.0.108.0)" == $'\n- old' ]] || fail 'notes 108'
[[ -z "$(notes v9.0.0)" ]] || fail 'notes v9.0.0 matched a prefix'

printf 'promote-changelog tests passed\n'
