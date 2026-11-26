# Changelog

All notable changes to the v9 major version of this project will be documented in this file.
See updating [Changelog example here](https://keepachangelog.com/en/1.0.0/)

This changelog tracks v9 pre-releases while v9 is developed on the `v9` branch.
It will be merged into the root `CHANGELOG.md` when v9 becomes stable and lands on `main`.

## [Unreleased]

### Changed

- go: minimum Go version is 1.26.11.
- generator: upgrade oapi-codegen to v2.9.0. Client templates follow its
  upstream v2.9.0 templates. Local edits live in
  `generator/config/templates/patches`.
- gateway: `GatewayLabelCreateRequest.Value` drops its pointer.
- database: date fields such as
  `DatabaseMaintenancePendingUpdatesResponse.Deadline` change from `*string`
  to `*time.Time`.
- See `make test-public-api` for the full list of API changes.

## [9.0.0-alpha.0.0.108.0] - 2026-10-05

<!-- markdownlint-disable-next-line MD024 -->
### Changed

- Initial v9 alpha release.
