#!/usr/bin/env bash
set -euo pipefail
export GOWORK=off GOSUMDB=off
export GOPROXY="file://$(go env GOMODCACHE)/cache/download"

SCRIPT_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)
FIXTURES="$SCRIPT_DIR/../api/fixtures"
TEST_ROOT=$(mktemp -d)
trap 'status=$?; if (( status )) && [[ -f "$TEST_ROOT/output" ]]; then cat "$TEST_ROOT/output" >&2; fi; rm -rf "$TEST_ROOT"' EXIT
mkdir "$TEST_ROOT/tmp" "$TEST_ROOT/reports"
export TMPDIR="$TEST_ROOT/tmp"
REPO="$TEST_ROOT/repo"
git init -q "$REPO"
git -C "$REPO" config user.name 'Public API Test'
git -C "$REPO" config user.email 'public-api-test@example.invalid'
git -C "$REPO" -c commit.gpgsign=false commit -q --allow-empty -m 'Before v9'
INITIAL=$(git -C "$REPO" rev-parse HEAD)
mkdir -p "$REPO/v9/test/scripts" "$REPO/v9/pkg/upcloud"
cp "$SCRIPT_DIR/check-public-api.sh" "$REPO/v9/test/scripts/"
cp "$FIXTURES/go.mod.txt" "$REPO/v9/go.mod"
cp "$FIXTURES/base.go.txt" "$REPO/v9/pkg/upcloud/api.go"
SCRIPT="$REPO/v9/test/scripts/check-public-api.sh"

run() {
    bash "$SCRIPT" "$@" > "$TEST_ROOT/output" 2>&1
}

expect_failure() {
    if run "$@"; then
        printf 'Expected failure: %s\n' "$*" >&2
        exit 1
    fi
}

expect_output() {
    if ! grep -Fq "$1" "$TEST_ROOT/output"; then
        cat "$TEST_ROOT/output" >&2
        printf 'Missing output: %s\n' "$1" >&2
        exit 1
    fi
}

assert_clean() {
    if [[ -n $(find "$TMPDIR" -mindepth 1 -print -quit) ]]; then
        printf '%s\n' 'Temporary export/archive files were not cleaned.' >&2
        exit 1
    fi
}

run update
cp "$REPO/v9/test/api/public-api.export" "$TEST_ROOT/first.export"
run update
run check --report-dir "$TEST_ROOT/reports"
[[ ! -s "$TEST_ROOT/reports/baseline.diff" ]]
assert_clean
git -C "$REPO" add v9
git -C "$REPO" -c commit.gpgsign=false commit -q -m 'Base public API'
BASE=$(git -C "$REPO" rev-parse HEAD)

run check "$INITIAL" --report-dir "$TEST_ROOT/reports"
expect_output 'Merge-base has no v9'
expect_output 'Existing: added'
[[ -s "$TEST_ROOT/reports/merge-base.diff" ]]
[[ $(cat "$TEST_ROOT/reports/merge-base.txt") == "$INITIAL" ]]
assert_clean

cp "$FIXTURES/addition.go.txt" "$REPO/v9/pkg/upcloud/api.go"
expect_failure check
expect_output 'Added: added'
run update
run check
git -C "$REPO" add v9
git -C "$REPO" -c commit.gpgsign=false commit -q -m 'Add and update baseline'
run check --base-ref "$BASE" --report-dir "$TEST_ROOT/reports"
expect_output 'Added: added'
[[ ! -s "$TEST_ROOT/reports/baseline.diff" ]]
grep -Fq 'Added: added' "$TEST_ROOT/reports/merge-base.diff"
[[ $(cat "$TEST_ROOT/reports/merge-base.txt") == "$BASE" ]]
assert_clean

cp "$FIXTURES/incompatible.go.txt" "$REPO/v9/pkg/upcloud/api.go"
expect_failure check --base-ref "$BASE" --report-dir "$TEST_ROOT/reports"
expect_output 'Incompatible changes:'
expect_output 'Added: removed'
grep -Fq 'Existing:' "$TEST_ROOT/reports/merge-base.diff"
assert_clean

cp "$FIXTURES/base.go.txt" "$REPO/v9/pkg/upcloud/api.go"
run update
run check --base-sha HEAD --report-dir "$TEST_ROOT/reports"
expect_output 'Added: removed'
grep -Fq 'Added: removed' "$TEST_ROOT/reports/base-sha.diff"
[[ $(cat "$TEST_ROOT/reports/base-sha.txt") == $(git -C "$REPO" rev-parse HEAD) ]]
assert_clean
expect_failure check --base-sha HEAD --base-ref "$BASE"
expect_failure check --base-ref refs/heads/does-not-exist
assert_clean
cp "$FIXTURES/invalid.go.txt" "$REPO/v9/pkg/upcloud/api.go"
expect_failure update
cmp "$TEST_ROOT/first.export" "$REPO/v9/test/api/public-api.export"
assert_clean
cp "$FIXTURES/base.go.txt" "$REPO/v9/pkg/upcloud/api.go"
mkdir "$REPO/v9/extra"
cp "$FIXTURES/base.go.txt" "$REPO/v9/extra/api.go"
expect_failure check
expect_output 'Public package inventory changed'
assert_clean
rm -rf "$REPO/v9/extra"
mkdir -p "$REPO/v9/testonly"
printf 'package testonly\n\nimport "testing"\n\nfunc TestOnly(t *testing.T) {}\n' > "$REPO/v9/testonly/api_test.go"
run check
expect_output 'No API differences.'
assert_clean
rm -rf "$REPO/v9/testonly"
cp "$FIXTURES/base.go.txt" "$REPO/v9/test/api/public-api.export"
expect_failure check
expect_output 'apidiff failed for baseline'
assert_clean
rm "$REPO/v9/test/api/public-api.export"
expect_failure check
expect_output 'Missing or empty baseline'
assert_clean

printf '%s\n' 'PASS: semantic self-check, addition, incompatible change, removal,' \
    'reviewed baseline versus merge-base, exact push baseline, pre-v9 history, package inventory, tool errors, and cleanup (offline).'