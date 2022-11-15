# Changelog

All notable changes to this project will be documented in this file.
See updating [Changelog example here](https://keepachangelog.com/en/1.0.0/)

## [Unreleased]

### Changed
- bumped module version from `v4` to `v5`
- context-aware type `client.ClientContext` renamed to `client.Client`
- context-aware type `service.ServiceContext` renamed to `service.Serivce`
- `service.Serivce` accepts any client that implements `service.Client` interface 
- `client.New` constructor allows overwrite `baseURL`, `httpClient` and `httpClient.Timeout` client properties

### Removed
- `service.Service` type without context support
- `client.Client` type without context support
- redundant `Timeout` field from `request.StartServerRequest` type


## [4.10.0]

### Added
- Managed Load Balancer private network support
- Anti-affinity support for server groups (experimental)

### Deprecated
- Managed Load Balancer fields `DNSName` and `NetworkUUID`

### Removed
- separate kubernetes plans

## [4.9.0]

### Added
- label support for server and server groups
- add support for resource permissions
- Experimental support for Kubernetes cluster, plans and versions

## [4.8.0]

### Added
- Managed Load Balancer `SetForwardedHeaders` frontend rule action

## [4.7.0]

### Added
- add context-aware client and service types
- Managed Load Balancer frontend and backend properties

## [4.6.0]

### Added
- add support for fetching available service types for managed database service

### Fixed
- drop IP address family requirement from firewall rule
- typo in `ManagedDatabaseConnection` field `WaitEventType`

## [4.5.2]

### Fixed
- Content-Type HTTP header when importing compressed storage

## [4.5.1]

### Fixed
- Expose `GetZones` through `Zones` interface
- allow time or day of week to be empty in Managed Database maintenance request property

## [4.5.0]

### Added
- add experimental support for Server Groups
- add possibility to overwrite default API URL with `UPCLOUD_DEBUG_API_BASE_URL` environment variable

## [4.4.2]

### Fixed
- Managed Load Balancer resolver request spelling

## [4.4.1]

### Added
- Managed Load Balancer request example

### Fixed
- Managed Load Balancer num_members_up matcher naming convention

## [4.4.0]

### Added
- add support for Managed Load Balancer (beta)

## [4.3.0]

### Added
- add support for fetching available versions for managed database service
- add support for upgrading managed database version

## [4.2.0]

### Added

- add support for resizing storage partition and filesystem

## [4.1.3]

### Added

- add support for subaccount management and account listing

## [4.1.2]

### Fixed

- rewrite import paths for v4

## [4.1.1]

### Fixed

- bump go.mod major version

## [4.1.0]

### Added

- add experimental support for Managed Databases https://github.com/UpCloudLtd/upcloud-go-api/pull/96

### Changed

- We will now publish the release tags with the `v`-prefix to conform with go mod

## [4.0.0]

### Added

- add support for Object Storage
- support io.Reader as SourceLocation in CreateStorageImport
- pass on a custom user agent string in requests
- add storage tier support to ServerStorageDevice
- support for backuprules when creating a server
- add support for explicitly setting IP address for network interfaces (requires special privileges for your UpCloud account)

### Changed

- *BREAKING CHANGE*: changing network router with ModifyNetwork call is no longer supported. Please use AttachNetworkRouter and DetachNetworkRouter from now on.
- bump default Go version to 1.16, keep supporting 1.15
- use a default timeout when no timeout given

### Fixed

- default to original storage size in CreateServerStorageDevice when Size = 0
- *BREAKING CHANGE*: Fix StorageImportDetails.Completed to be a time.Time rather than a string
- don't marshal empty resource limits
- allow empty BackupRules (eg. remove backup rule) to be sent to the backend

## [3.0.0]

### Added

- go.mod
- delete server and all attached storages
- a build script with more code quality checks
- go-vcr to tests
- Host, Firewall, Network, Router resources from UpCloud API 1.3
- Storage import resource

### Changed

- bump Go version to 1.14
- Converted API calls from XML to JSON
- default zone from hel1 to hel2
- only build PRs and merges to master
- raise mininum required go version to 1.13
- changelog format to include different lists
- bump UpCloud API from 1.2 to 1.3 and expand with new functionalities

### Removed

- XML api calls

### Fixed

- Timeout issues

## [2.0.0]

- moved project to UpCloud's own GitHub organization
- raise the minimum required Go version to 1.7

## [1.1.0]

- improve documentation
- remove ability to override the API base URL and version

## [1.0.2]

- remove credentials related getters and setters from `Client`
- hopefully fix CD-ROM integration tests for good by performing all operations while the server is stopped

## [1.0.1]

- minor tweaks to the integration tests
- correct the package name in the README installation instructions

## [1.0.0]

First stable release

[Unreleased]: https://github.com/UpCloudLtd/upcloud-go-api/compare/v4.10.0...HEAD
[4.10.0]: https://github.com/UpCloudLtd/upcloud-go-api/compare/v4.9.0...v4.10.0
[4.9.0]: https://github.com/UpCloudLtd/upcloud-go-api/compare/v4.8.0...v4.9.0
[4.8.0]: https://github.com/UpCloudLtd/upcloud-go-api/compare/v4.7.0...v4.8.0
[4.7.0]: https://github.com/UpCloudLtd/upcloud-go-api/compare/v4.6.0...v4.7.0
[4.6.0]: https://github.com/UpCloudLtd/upcloud-go-api/compare/v4.5.2...v4.6.0
[4.5.2]: https://github.com/UpCloudLtd/upcloud-go-api/compare/v4.5.1...v4.5.2
[4.5.1]: https://github.com/UpCloudLtd/upcloud-go-api/compare/v4.5.0...v4.5.1
[4.5.0]: https://github.com/UpCloudLtd/upcloud-go-api/compare/v4.4.2...v4.5.0
[4.4.2]: https://github.com/UpCloudLtd/upcloud-go-api/compare/v4.4.1...v4.4.2
[4.4.1]: https://github.com/UpCloudLtd/upcloud-go-api/compare/v4.4.0...v4.4.1
[4.4.0]: https://github.com/UpCloudLtd/upcloud-go-api/compare/v4.3.0...v4.4.0
[4.3.0]: https://github.com/UpCloudLtd/upcloud-go-api/compare/v4.2.0...v4.3.0
[4.2.0]: https://github.com/UpCloudLtd/upcloud-go-api/compare/v4.1.3...v4.2.0
[4.1.3]: https://github.com/UpCloudLtd/upcloud-go-api/compare/v4.1.2...v4.1.3
[4.1.2]: https://github.com/UpCloudLtd/upcloud-go-api/compare/v4.1.1...v4.1.2
[4.1.1]: https://github.com/UpCloudLtd/upcloud-go-api/compare/v4.1.0...v4.1.1
[4.1.0]: https://github.com/UpCloudLtd/upcloud-go-api/compare/4.0.0...v4.1.0
[4.0.0]: https://github.com/UpCloudLtd/upcloud-go-api/compare/3.0.0...4.0.0
[3.0.0]: https://github.com/UpCloudLtd/upcloud-go-api/compare/2.0.0...3.0.0
[2.0.0]: https://github.com/UpCloudLtd/upcloud-go-api/compare/1.1.0...2.0.0
[1.1.0]: https://github.com/UpCloudLtd/upcloud-go-api/compare/1.0.2...1.1.0
[1.0.2]: https://github.com/UpCloudLtd/upcloud-go-api/compare/1.0.1...1.0.2
[1.0.1]: https://github.com/UpCloudLtd/upcloud-go-api/compare/1.0.0...1.0.1
[1.0.0]: https://github.com/UpCloudLtd/upcloud-go-api/releases/tag/1.0.0
