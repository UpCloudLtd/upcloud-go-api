# SDK Tests

Run targets from `v9`. Go 1.25.0 or newer, Bash, Make and Git are required;
lint additionally requires golangci-lint 2.10.1. Tools and module dependencies
may need downloading on the first run. The tests themselves use local HTTP only,
with no UpCloud credentials required.

The Makefile disables automatic Go toolchain upgrades and parent workspaces using
`GOTOOLCHAIN=local` and `GOWORK=off`. Select the desired Go installation on `PATH`;
the helper scripts also resolve `go` through `PATH`.

## Targets

| Target | Result |
| --- | --- |
| `make build` | Compile all packages, examples and test code without running tests. |
| `make test` | Run the SDK module's tests. Nested consumer tests have their own target. |
| `make vet` | Analyze generated code, extensions, tests and examples. |
| `make lint` | Run the established root lint configuration against v9. |
| `make test-consumer` | Compile and run a disposable external consumer module. |
| `make test-generator` | Exercise generator input and output safety in disposable fixtures. |
| `make test-public-api` | Check the reviewed public API baseline. |
| `make test-public-api-gate` | Exercise the API check's own regression fixtures. |
| `make update-public-api` | Explicitly replace the API baseline for review. |
| `make check` | Run every check above except the baseline update. |

`TEST_TIMEOUT` defaults to `5m` for build/test/lint. Override the linter executable
with `GOLANGCI_LINT` if needed. Generated files and examples are excluded from lint
findings by the existing configuration; build and vet still cover them.

The consumer copies a small module outside the SDK tree, replaces only the SDK
with the candidate checkout, and exercises public typed calls against a local
server. Its tests have a separate two-minute timeout. This verifies checkout
consumption, not a published release.

## Public API Review

The API tool is pinned to
`golang.org/x/exp/cmd/apidiff@v0.0.0-20250911091902-df9299821621`.
[api/public-api.export](api/public-api.export) is its binary export baseline,
not a byte-stable snapshot across toolchains. Comparisons are semantic and fail
on any unreviewed difference, including additions. Adding another public package
also fails until the comparison inventory covers it.

To compare with the PR target and write a report:

```bash
make test-public-api API_BASE_REF=origin/v9 API_REPORT_DIR=/tmp/v9-api-report
```

For an intentional API change, inspect that report, run `make update-public-api`,
then rerun the comparison. Review the baseline change and report together.
Merge-base differences remain visible even after a baseline update; they are
warnings rather than automatic failures because deliberate additions and alpha
changes may be valid. Incompatible changes need migration guidance and review
under the v9 release policy. Do not update the baseline to hide accidental breaks.

The gate's regression fixtures check additions, incompatibilities, reviewed
baseline updates, direct commit comparisons, missing packages, tool failures and
cleanup. They create Git history only in a disposable fixture repository.

For a direct comparison rather than a merge base:

```bash
bash test/scripts/check-public-api.sh check --base-sha <commit> --report-dir /tmp/v9-api-report
```

## GitHub Checks

The root [workflow](../../.github/workflows/v9-tests.yml) emits separate checks:

- `V9 Build (Go 1.25.0)`, `V9 Build (Go 1.25.x)`, `V9 Build (Go 1.26.x)`.
- `V9 Consumer (Go 1.25.0)`, `V9 Consumer (Go 1.25.x)`, `V9 Consumer (Go 1.26.x)`.
- `V9 Vet`, `V9 Lint`, `V9 Public API`, and `V9 Generator Input (Go 1.25.0)` / `V9 Generator Input (Go 1.26.x)`.

Checks run on default PR events and pushes to v9, without
cloud secrets or path filters. Each job owns a separate result; aggregate local
validation does not replace those checks. These checks do not yet establish full
operation or v8 behavioral parity.

PR API reports use the target merge base. Push reports compare directly with the
previous SHA, fetching it if needed; an unavailable commit fails rather than
silently omitting the report. Both also check the reviewed baseline.

GitHub rulesets must separately require the emitted check names. Workflow YAML
alone cannot enforce merge blocking. Regeneration PRs created with `GITHUB_TOKEN`
need an App/PAT event to trigger these checks because of GitHub's recursion guard.

## Generation Safety

Run generation from `v9` using `mise run generate`. The input must already exist
at `generator/spec.json`. The regeneration workflow downloads the requested
release archive, copies the requested member to this path, and removes it before
opening its PR. The SDK does not need this file at runtime.

Generation requires Bash, jq and the Go toolchain/dependencies in the module.
The input check requires one JSON OpenAPI 3.0/3.1 document with nonempty title,
version, paths and operations. Each operation must have exactly one declared
tag. Tag names must be unique and safe as filenames (letters, digits,
underscores and hyphens, starting with a letter). Single-tag operations are a
constraint of this tag-based generator, not a restriction in OpenAPI generally.

The jq check is not a complete OpenAPI semantic validator. oapi-codegen performs
its own loading, reference and generation checks afterward. Outputs are staged
in a temporary directory. Existing generated files are removed/replaced only
after every generation command succeeds; handwritten files are never selected
for replacement. The temporary directory is removed on ordinary success or
failure. Final publication is not transactional against filesystem errors or
interruption, so run regeneration in an isolated checkout.

Run the generator input-safety tests with `make test-generator` from `v9`. They
invoke the real script, jq and pinned oapi-codegen in temporary fixtures; they
need no private release spec or cloud credentials. Coverage includes invalid
and unsafe inputs, unsupported tag layouts, unresolved references, a late
template failure, output preservation, and successful generation. The success
case parses generated Go files but does not claim full regeneration
reproducibility or consumer compilation.

### Spec Provenance Follow-Up

The latest committed change to `pkg/upcloud/spec_version.go` is regeneration
commit `4b4c922`, identifying version `0.0.108` (PR #621). The current release
changelog names `9.0.0-alpha.0.0.108.0`. The previous `118` in the release
heading was a typo, corrected after maintainer confirmation. No generated code
or existing Git tag was changed by that correction.

Before claiming reproducibility, recover the original regeneration workflow run
and record its release/asset URL, archive member, archive and extracted-spec
SHA-256 values, supplied version, and generator/configuration revisions. These
historical asset/checksum values are not recorded in this checkout. Recover them
in a separate provenance PR; do not invent checksums or substitute a sibling
SDK's spec.