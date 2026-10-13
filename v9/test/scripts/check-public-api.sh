#!/usr/bin/env bash
set -euo pipefail
export GOWORK=off GOTOOLCHAIN=local

APIDIFF='golang.org/x/exp/cmd/apidiff@v0.0.0-20250911091902-df9299821621'
SDK_ROOT=$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)
BASELINE="$SDK_ROOT/test/api/public-api.export"
MODULE='github.com/UpCloudLtd/upcloud-go-api/v9'
PACKAGE="$MODULE/pkg/upcloud"

usage() {
    printf '%s\n' \
        'Usage: check-public-api.sh [check|update] [--base-ref REF|--base-sha SHA] [--report-dir DIR]' \
        'check (default): fail on any unreviewed baseline API difference.' \
        'update: replace the baseline with export data from the v9 package.' \
        'BASE_REF: additionally compare the working API with git merge-base HEAD REF.' \
        '--base-sha SHA: compare directly with that commit (for push reports).' \
        'Reports include additions and incompatible changes, not only removals.' \
        'Merge-base differences are reported for review, not automatically rejected.' \
        'A matching reviewed baseline does not establish nonbreaking PR changes.' \
        "Immutable tool pin: $APIDIFF (never resolve latest)." \
        'Export data is tool-generated binary data; use the same pin to read/write it.'
}

MODE=check
BASE_REF=''
BASE_SHA=''
REPORT_DIR=''
if [[ ${1:-} == check || ${1:-} == update ]]; then
    MODE=$1
    shift
fi
while (( $# )); do
    case "$1" in
        --base-ref|--base-sha|--report-dir)
            if (( $# < 2 )) || [[ -z $2 ]]; then
                usage >&2
                exit 2
            fi
            if [[ $1 == --base-ref ]]; then
                BASE_REF=$2
            elif [[ $1 == --base-sha ]]; then
                BASE_SHA=$2
            else
                REPORT_DIR=$2
            fi
            shift 2
            ;;
        -h|--help) usage; exit 0 ;;
        -*) usage >&2; exit 2 ;;
        *)
            if [[ -n $BASE_REF ]]; then
                usage >&2
                exit 2
            fi
            BASE_REF=$1
            shift
            ;;
    esac
done
if [[ -n $BASE_REF && -n $BASE_SHA ]]; then
    printf '%s\n' 'Choose either a base ref or an exact base commit.' >&2
    exit 2
fi
if [[ $MODE == update && ( -n $BASE_REF || -n $BASE_SHA ) ]]; then
    printf '%s\n' 'update cannot be combined with a base ref; run check separately.' >&2
    exit 2
fi
if [[ -n $REPORT_DIR ]]; then
    mkdir -p "$REPORT_DIR"
    REPORT_DIR=$(cd "$REPORT_DIR" && pwd)
fi
TEMP_DIR=$(mktemp -d)
trap 'rm -rf "$TEMP_DIR"' EXIT
trap 'exit 130' INT
trap 'exit 143' TERM

export_api() {
    local root=$1 output=$2
    local packages
    packages=$(cd "$root" && go list -buildvcs=false -mod=readonly -f '{{if and (ne .Name "main") (gt (len .GoFiles) 0)}}{{.ImportPath}}{{end}}' ./... | sed '/^$/d')
    if [[ $packages != "$PACKAGE" ]]; then
        printf 'Public package inventory changed; extend the baseline and comparisons before proceeding:\n%s\n' "$packages" >&2
        return 1
    fi
    (cd "$root" && go run "$APIDIFF" -w "$output" "$PACKAGE")
    if [[ ! -s $output ]]; then
        printf '%s\n' 'apidiff produced empty export data.' >&2
        return 1
    fi
}

compare_api() {
    local old=$1 new=$2 name=$3
    local report="$TEMP_DIR/$name.diff"
    if ! (cd "$SDK_ROOT" && go run "$APIDIFF" "$old" "$new") > "$report"; then
        printf 'apidiff failed for %s; this is not an API comparison result.\n' "$name" >&2
        exit 1
    fi
    printf '\nPublic API comparison: %s\n' "$name"
    cat "$report"
    if [[ -n $REPORT_DIR ]]; then
        cp "$report" "$REPORT_DIR/$name.diff"
    fi
    if [[ -s $report ]]; then
        if [[ $name == baseline ]]; then
            CHANGED=1
        else
            printf '%s\n' '::warning::Public API differs from the comparison base; review the API report and migration policy.'
        fi
    else
        printf '%s\n' 'No API differences.'
    fi
}

export_api "$SDK_ROOT" "$TEMP_DIR/current.export"
if [[ $MODE == update ]]; then
    mkdir -p "$(dirname "$BASELINE")"
    STAGED_BASELINE=$(mktemp "${BASELINE}.XXXXXX")
    trap 'rm -rf "$TEMP_DIR"; rm -f "${STAGED_BASELINE:-}"' EXIT
    cp "$TEMP_DIR/current.export" "$STAGED_BASELINE"
    mv "$STAGED_BASELINE" "$BASELINE"
    printf 'Updated %s using %s\n' "$BASELINE" "$APIDIFF"
    exit 0
fi
if [[ ! -s $BASELINE ]]; then
    printf 'Missing or empty baseline: %s\n' "$BASELINE" >&2
    exit 1
fi
CHANGED=0
compare_api "$BASELINE" "$TEMP_DIR/current.export" baseline
if [[ -n $BASE_REF || -n $BASE_SHA ]]; then
    REPO_ROOT=$(git -C "$SDK_ROOT" rev-parse --show-toplevel)
    COMPARISON=merge-base
    if [[ -n $BASE_SHA ]]; then
        BASE_COMMIT=$(git -C "$REPO_ROOT" rev-parse --verify "${BASE_SHA}^{commit}")
        COMPARISON=base-sha
        printf '\nExact baseline commit: %s\n' "$BASE_COMMIT"
    else
        BASE_COMMIT=$(git -C "$REPO_ROOT" merge-base HEAD "$BASE_REF")
        printf '\nMerge-base: %s (HEAD versus %s)\n' "$BASE_COMMIT" "$BASE_REF"
    fi
    if [[ -n $REPORT_DIR ]]; then
        printf '%s\n' "$BASE_COMMIT" > "$REPORT_DIR/$COMPARISON.txt"
    fi
    mkdir "$TEMP_DIR/base"
    git -C "$REPO_ROOT" archive "$BASE_COMMIT" | tar -x -C "$TEMP_DIR/base"
    if [[ -f "$TEMP_DIR/base/v9/go.mod" ]]; then
        export_api "$TEMP_DIR/base/v9" "$TEMP_DIR/base.export"
    elif git -C "$REPO_ROOT" cat-file -e "$BASE_COMMIT:v9" 2>/dev/null; then
        printf '%s\n' 'Merge-base contains v9 but no v9/go.mod; cannot compare its API.' >&2
        exit 1
    else
        printf '%s\n' 'Merge-base has no v9; comparing against an empty v9 API.'
        mkdir -p "$TEMP_DIR/empty/pkg/upcloud"
        printf 'module %s\n\ngo 1.25.0\n' "$MODULE" > "$TEMP_DIR/empty/go.mod"
        printf 'package upcloud\n' > "$TEMP_DIR/empty/pkg/upcloud/api.go"
        export_api "$TEMP_DIR/empty" "$TEMP_DIR/base.export"
    fi
    compare_api "$TEMP_DIR/base.export" "$TEMP_DIR/current.export" "$COMPARISON"
fi
if (( CHANGED )); then
    printf '\nPublic API differences detected.\n' >&2
    exit 1
fi