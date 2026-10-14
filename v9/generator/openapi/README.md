# Canonical OpenAPI Input

`spec.json` is the durable, versioned source used by `mise run generate` and the
request-conformance inventory. Do not place it at the old temporary
`generator/spec.json` path.

The regeneration workflow downloads the release asset into runner-temporary
storage, verifies the extracted member's SHA-256 and `info.version`, then replaces
`spec.json` and writes `provenance.json` before generating. The resulting PR must
include the spec, provenance, generated SDK, and spec-version changes together.
Workflow cleanup removes only the temporary download; it must never restore or
delete this directory.

For a manual source update, use the exact release asset and archive member, record
their version and SHA-256 in `provenance.json`, and review the canonical spec diff
with generated output. Never infer provenance from the version string alone.