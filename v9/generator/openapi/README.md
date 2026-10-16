# Canonical OpenAPI Input

`spec.json` is the durable, versioned source used by `mise run generate` and the
request-conformance inventory. Do not place it at the old temporary
`generator/spec.json` path.

The regeneration workflow downloads the release asset into runner-temporary
storage, verifies the extracted member's SHA-256 and `info.version`, then replaces
`spec.json` before generating. The resulting PR must include the spec, generated
SDK, and spec-version changes together. Git retains the exact spec for each revision.
Workflow cleanup removes only the temporary download; it must never restore or
delete this directory.

For a manual source update, use the exact release asset and archive member and
verify the extracted spec against the upstream version and SHA-256 using
`generator/scripts/install-spec.sh <archive> <member> <version> <sha256>` from `v9`.
Review the canonical spec diff with generated output. Local generation validates
the committed OpenAPI document without requiring separate provenance metadata.