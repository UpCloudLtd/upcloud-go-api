# Generator Input Safety

Run generation from `v9` using `mise run generate`. The input must already exist
at `generator/spec.json`. The regeneration workflow downloads the requested
release archive, copies the requested member to this path, and removes it before
opening its PR. The SDK does not need this file at runtime.

Generation requires Bash, jq and the Go toolchain/dependencies in the module.
The input check requires one JSON OpenAPI 3.0/3.1 document with nonempty title,
version, paths, operations and declared tags. Every operation needs declared
tags, and tag names must be unique and safe as filenames (letters, digits,
underscores and hyphens, starting with a letter). These are this tag-based
generator's input requirements, not restrictions on OpenAPI in general.

The jq check is not a complete OpenAPI semantic validator. oapi-codegen performs
its own loading/reference/generation checks afterward. All outputs are staged in
a temporary directory. Existing generated files are removed/replaced only after
every generation command succeeds; handwritten files are never selected for
replacement. The temporary directory is removed on ordinary success or failure.
File replacement is not transactional against filesystem errors or interruption
during the final publication step; run regeneration in an isolated checkout.

## Tests

From `v9`, run:

```bash
make -C generator test
```

Tests invoke the real script, jq and pinned oapi-codegen in temporary fixture
directories. They do not need the private release spec or cloud credentials;
dependencies may need downloading on the first run. They check missing/invalid
inputs, unsafe or incomplete tag layouts, unresolved references, a late template
failure, and successful replacement. The success test parses generated Go files;
it does not claim full SDK regeneration reproducibility or consumer compilation.

## Spec Provenance Follow-Up

The latest committed change to `pkg/upcloud/spec_version.go` is regeneration
commit `4b4c922`, whose message identifies version `0.0.108` (PR #621). The current
release changelog names `9.0.0-alpha.0.0.108.0`, matching the SDK spec constant.
The previous `118` in the release heading was a typo, corrected after maintainer
confirmation. No generated code or existing Git tag was changed by this correction.

Before claiming reproducibility, recover the original regeneration workflow run
and record its release/asset URL, archive member, archive and extracted-spec
SHA-256 values, supplied version, and generator/configuration revisions. Those
historical asset/checksum values are not recorded in this checkout. Recover them
in a separate provenance PR; do not invent checksums or substitute a sibling
SDK's spec.