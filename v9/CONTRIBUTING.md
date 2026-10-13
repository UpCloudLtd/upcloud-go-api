# Contributing

This SDK is generated from the UpCloud OpenAPI specification using [oapi-codegen](https://github.com/oapi-codegen/oapi-codegen).

Regeneration is typically triggered via workflow on new spec release.

## Regenerating locally

```bash
mise run generate
```

## Testing

Run `make test` from this directory with Go 1.25.0 or newer. Run `make check`
for all local checks, including generator input-safety tests and lint
(golangci-lint 2.10.1). See [test/README.md](test/README.md) for individual
targets, generation safety, and API review guidance.
