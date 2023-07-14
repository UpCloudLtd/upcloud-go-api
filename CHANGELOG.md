# Changelog

All notable changes to this project will be documented in this file.
See updating [Changelog example here](https://keepachangelog.com/en/1.0.0/)

## [Unreleased]

### Added
- kubernetes: `utility_network_access` field to node group for configuring utility network access on the given group
- Managed Database session support, including methods `service.GetManagedDatabaseSessions` & `service.CancelManagedDatabaseSession`.

### Deprecated
- `service.GetManagedDatabaseConnections` and `service.CancelManagedDatabaseConnection` in favor of `service.GetManagedDatabaseSessions` and `service.CancelManagedDatabaseSession`

## [6.4.0]

### Added
- client functions `NewDefaultHTTPClient` and `NewDefaultHTTPTransport` to provide HTTP client default properties
- kubernetes: experimental support for deleting nodes from node groups
- kubernetes: consts for `scaling-up` and `scaling-down` node-group states

### Changed
- `service.GetKubernetesNodeGroup` method to return `upcloud.KubernetesNodeGroupDetails` type which is extended version of the previous `upcloud.KubernetesNodeGroup`

### Fixed
- `request.ModifyServerRequest` does not set boolean properties `Metadata` and `RemoteAccessEnabled` to `"no"` by default.

## [6.3.2]

### Added
- Managed Load Balancer `health_check_tls_verify` field to backend member properties to control certificate validation for health checks utilising HTTPS

## [6.3.1]

### Added
- Managed Load Balancer inverse rule matcher

## [6.3.0]

### Added
- ServerGroup `AntiAffinityPolicy` field to support strict, best-effort and off policies. This replaces `AntiAffinity`

### Removed
- ServerGroup `AntiAffinity` boolean field in favor of `AntiAffinityPolicy` string enum field

### Changed
- GetManagedDatabaseIndices method to return a slice of structs instead of pointers 

## [6.2.0]

### Added
- Managed Database OpenSearch support
- Support for defining NIC model upon creating or modifying a server. Also exported constants for each support NIC model.
- Support for `PrivateNodeGroups` property of Kubernetes clusters.

### Changed
- client: overwrite the HTTP Client Transport accordingly when `UPCLOUD_DEBUG_SKIP_CERTIFICATE_VERIFY` is set to `1`

## [6.1.1]

### Added
- method for fetching available kubernetes plans

## [6.1.0]

### Added
- new `Backups` field for `DeleteServerAndStoragesRequest` that controls if backups should be kept or deleted while deleting the server
- Kubernetes cluster `Plan` field

## [6.0.0]

### Added
- labels support for storages
- support for network gateways
- `State` field to kubernetes node-groups

### Changed
- errors: all service method now return `Problem` type in case of errors (*BREAKING CHANGE*)
- Type for `NodeGroups` field `CreateKubernetesClusterRequest` is now `request.KubernetesNodeGroup` instead of `upcloud.KubernetesNodeGroup`

### Removed
- errors: `Error` type was removed (*BREAKING CHANGE*)

## [5.4.0]

### Added
- kubernetes: experimental support for anti-affinity node groups

## [5.3.1]

### Deprecated
- request: `GetServerGroupsWithFiltersRequest` is now deprecated; use `GetServerGroupsRequest` with `Filters` field instead

## [5.3.0]

### Added
- storage: possibility to control what to do with backups related to the storage that is about to be deleted
- labels support for: load-balancers, networks, network-peerings, routers

## Deprecated
- request: `ServerFilter` and `ServerGroupFilter` types will be replaced with `QueryFilter` type in the future

## [5.2.1]

### Added
- kubernetes: experimental support for node group CRUD operations

## [5.2.0]

### Added
- experimental support for Managed Redis Database
- load-balancer: add `Scheme` property to front-end rule HTTP redirect action.

## [5.1.0]

### Added
- servers: allow adding server to a server group during server creation
- servers: add server group UUID to server details
- HTTP status code field to `upcloud.Error` type
- experimental support for internal network peering

## [5.0.0]

### Changed
- context-aware type `client.ClientContext` renamed to `client.Client`
- context-aware type `service.ServiceContext` renamed to `service.Service`
- `service.Service` accepts any client that implements `service.Client` interface 
- `client.New` constructor allows overwrite `baseURL`, `httpClient` and `httpClient.Timeout` client properties
- bump Go version to 1.18
- bump module version from `v4` to `v5`

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

[Unreleased]: https://github.com/UpCloudLtd/upcloud-go-api/compare/v6.3.2...HEAD
[6.3.2]: https://github.com/UpCloudLtd/upcloud-go-api/compare/v6.3.1...v6.3.2
[6.3.1]: https://github.com/UpCloudLtd/upcloud-go-api/compare/v6.3.0...v6.3.1
[6.3.0]: https://github.com/UpCloudLtd/upcloud-go-api/compare/v6.2.0...v6.3.0
[6.2.0]: https://github.com/UpCloudLtd/upcloud-go-api/compare/v6.1.1...v6.2.0
[6.1.1]: https://github.com/UpCloudLtd/upcloud-go-api/compare/v6.1.0...v6.1.1
[6.1.0]: https://github.com/UpCloudLtd/upcloud-go-api/compare/v6.0.0...v6.1.0
[6.0.0]: https://github.com/UpCloudLtd/upcloud-go-api/compare/v5.4.0...v6.0.0
[5.4.0]: https://github.com/UpCloudLtd/upcloud-go-api/compare/v5.3.1...v5.4.0
[5.3.1]: https://github.com/UpCloudLtd/upcloud-go-api/compare/v5.3.0...v5.3.1
[5.3.0]: https://github.com/UpCloudLtd/upcloud-go-api/compare/v5.2.1...v5.3.0
[5.2.1]: https://github.com/UpCloudLtd/upcloud-go-api/compare/v5.2.0...v5.2.1
[5.2.0]: https://github.com/UpCloudLtd/upcloud-go-api/compare/v5.1.0...v5.2.0
[5.1.0]: https://github.com/UpCloudLtd/upcloud-go-api/compare/v5.0.0...v5.1.0
[5.0.0]: https://github.com/UpCloudLtd/upcloud-go-api/compare/v4.10.0...v5.0.0
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
