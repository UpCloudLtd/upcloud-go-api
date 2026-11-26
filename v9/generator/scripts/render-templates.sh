#!/usr/bin/env bash
# Renders the custom oapi-codegen templates from the pinned upstream templates
# and the local patches in generator/config/templates/patches.
#
#   render-templates.sh          write the rendered templates
#   render-templates.sh check    fail if the checked-in templates differ
#
# Patches apply with no fuzz. When a patch fails after an oapi-codegen bump, the
# script keeps the work directory and prints the reject file. Fix the patch, then
# run render again.
set -euo pipefail

MODE=${1:-render}
GO_BIN=${GO:-go}
case "$MODE" in
  render|check) ;;
  *) echo "usage: $0 [render|check]" >&2; exit 2 ;;
esac

cd "$(dirname "$0")/../.."

TEMPLATES=(client.tmpl client-with-responses.tmpl)
OUT_DIR=generator/config/templates
PATCH_DIR="$OUT_DIR/patches"

"$GO_BIN" mod download github.com/oapi-codegen/oapi-codegen/v2
UPSTREAM_DIR="$("$GO_BIN" list -mod=readonly -m -f '{{.Dir}}' github.com/oapi-codegen/oapi-codegen/v2)"
if [[ -z $UPSTREAM_DIR ]]; then
  echo "ERROR: oapi-codegen module directory not found" >&2
  exit 1
fi
UPSTREAM_DIR="$UPSTREAM_DIR/pkg/codegen/templates"
WORK_DIR=$(mktemp -d)
keep_work=no
cleanup() {
  if [[ $keep_work == no ]]; then
    rm -rf "$WORK_DIR"
  fi
}
trap cleanup EXIT

for t in "${TEMPLATES[@]}"; do
  # Module cache files are read-only. Copy with write permission before patching.
  install -m 644 "$UPSTREAM_DIR/$t" "$WORK_DIR/$t"
  if ! patch -s -F0 --no-backup-if-mismatch -r "$WORK_DIR/$t.rej" \
    "$WORK_DIR/$t" < "$PATCH_DIR/$t.patch"; then
    keep_work=yes
    echo "ERROR: $PATCH_DIR/$t.patch does not apply to upstream $t" >&2
    echo "Reject file: $WORK_DIR/$t.rej" >&2
    echo "Upstream copy: $WORK_DIR/$t" >&2
    exit 1
  fi
done

if [[ $MODE == check ]]; then
  status=0
  for t in "${TEMPLATES[@]}"; do
    if ! diff -u "$OUT_DIR/$t" "$WORK_DIR/$t"; then
      echo "ERROR: $OUT_DIR/$t differs from upstream plus $PATCH_DIR/$t.patch" >&2
      echo "Run: bash generator/scripts/render-templates.sh" >&2
      status=1
    fi
  done
  exit $status
fi

for t in "${TEMPLATES[@]}"; do
  install -m 644 "$WORK_DIR/$t" "$OUT_DIR/$t"
done
echo "Rendered: ${TEMPLATES[*]}"
