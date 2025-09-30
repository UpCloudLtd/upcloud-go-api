# Changelog

All notable changes to this project will be documented in this file.
See updating [Changelog example here](https://keepachangelog.com/en/1.0.0/)

## [Unreleased]

### Added

- database: add `AdditionalDiskSpaceGiB` field to managed databases

## [8.28.0]

### Added

- database: add GetAllManagedDatabases to retrieve all managed databases regardless of paging

### Changed

- network: ExcludeBySource, FilterByDestination and FilterByRouteType in EffectiveRoutesAutoPopulation struct changed to pointers

## [8.27.0]

### Added

- load-balancer: support for floating IP addresses
- ip-address: support for release policy

## [8.26.0]

### Added

- network: support for dhcp_routes_configuration

### Changed

- Go version bump to 1.24
## [8.25.0]

### Added

- kubernetes: support for customising node groups utilising GPU and Cloud Native plans

## [8.24.0]

### Added

- gateway: support for getting metrics

## [8.23.0]

### Added

- account: add GPUs to resource limits
- server: add GPUs to server plans

## [8.22.0]

### Added
- account: support for getting billing summary
- network: support for assigning and deleting network interface IP addresses

## [8.21.0]

### Added
- loadbalancer: `WaitForLoadBalancerOperationalState` helper for waiting for the load balancer to achieve a desired operational state.
- loadbalancer: `WaitForLoadBalancerDeletion` blocks execution until the specified load balancer instance has been deleted.
- `GetDevicesAvailability` function for listing device availability per zone.

## [8.20.0]

### Added
- audit logs: export support
- client: streaming Get/Do variants

### Changed
- client: set minimum TLS version to 1.3

## [8.19.0]

### Added
- NewFromEnv client constructor for implicit client authentication through environment variables `UPCLOUD_USERNAME`, `UPCLOUD_PASSWORD` and `UPCLOUD_TOKEN`

### Fixed
- loadbalancer: optimize slice operations in load balancer request

### Changed
- Go version bump to 1.23

## [8.18.0]

### Added
- kubernetes: add support for cluster upgrade end-points

## [8.17.0]

### Added
- server: add relocation support

## [8.16.1]

### Added
- kubernetes: add `Deprecated` field to `KubernetesPlan`

## [8.16.0]

### Added
- Experimental support for token-based authentication in the client and functions for token management.

## [8.15.0]

### Added
- managed load balancer: support for redirect rule HTTP status

## [8.14.0]

### Added
- support for Partner API

## [8.13.0]

### Added
- managed databases: add support for termination protection

### Fixed
- server: add storage labels to server storage devices list items

## [8.12.0]

### Added
- managed databases: add support for Valkey

### Deprecated
- managed databases: deprecate Redis

## [8.11.0]

### Added
- server: add `Index` field to `request.CreateServerInterface`
- managed load balancer: `http_status`, `request_header`, and `response_header` rule matchers
- managed load balancer: `set_request_header`, and `set_response_header` rule actions

## [8.10.0]

### Added
- client: add option to configure logger for logging requests and responses

## [8.9.0]

### Added
- managed load balancer: `MatchingCondition` field to frontend rule for controlling which operator to use when combining matchers
- managed object storage: support for creating and deleting buckets
- managed object storage: `Deleted` field to `upcloud.ManagedObjectStorageBucketMetrics`

## [8.8.1]

### Changed
- Go version bump to 1.22

## [8.8.0]

### Added
- managed object storage: support for custom domains
- managed load balancer: support for reading DNS challenge domain

## [8.7.1]

### Added
- kubernetes: constant for standard storage tier

## [8.7.0]

### Added
- router: add `type` field to static routes
- managed databases: add support for labels

## [8.6.2]

### Added
- storage: `TemplateType` field to `upcloud.Storage`

### Fixed
- storage: typo in `StorageEncryptionDataAtRest` constant

### Deprecated
- `StorageEncryptionDataAtReset` in favor of `StorageEncryptionDataAtRest`

## [8.6.1]

### Added
- storage: constant for standard storage tier
- storage: constant for disabling storage encryption

## [8.6.0]

### Added
- **Experimental**, gateway: add UUID support for VPN gateway connections and tunnels
- service: a more elaborate error message if get request returns an error on unmarshalling json array

## [8.5.0]

### Added
- kubernetes: add support for node group custom plans
- kubernetes: add support for data at rest encryption in node groups

## [8.4.0]

### Added
- Cloud: `ParentZone` field to `Zone` struct (only available for private zones)

### Changed
- Go version bump to 1.21

## [8.3.0]

### Added
- **Experimental**, Gateway: support for VPN feature. Note that VPN feature is currently in beta, you can learn more about it on the [product page](https://upcloud.com/resources/docs/networking#nat-and-vpn-gateways)

## [8.2.0]

### Added
- Network peering: add `WaitForNetworkPeeringState` helper

## [8.1.0]

### Added
- Managed Database: add support managing attached SDN networks via `networks` field.
- Managed Database: add paging support to `GetManagedDatabases` method.

## [8.0.0]

### Added
- Managed Object Storage: `ManagedObjectStoragePolicy` struct
- Managed Object Storage: `ManagedObjectStorageUserPolicy` struct
- Managed Object Storage: `IAMURL` field to `ManagedObjectStorageEndpoint`
- Managed Object Storage: `STSURL` field to `ManagedObjectStorageEndpoint`
- Managed Object Storage: `ARN` field to `ManagedObjectStorageUser`
- Managed Object Storage: `Policies` field to `ManagedObjectStorageUser`
- Managed Object Storage: `Status` field to `ManagedObjectStorageUserAccessKey`

### Removed
- **Breaking**, Managed Object Storage: `Users` field removed from `ManagedObjectStorage`
- **Breaking**, Managed Object Storage: `ARN` field removed from `ManagedObjectStorageUser`
- **Breaking**, Managed Object Storage: `OperationalState` field removed from `ManagedObjectStorageUser`
- **Breaking**, Managed Object Storage: `Enabled` field removed from `ManagedObjectStorageUserAccessKey`
- **Breaking**, Managed Object Storage: `Name` field removed from `ManagedObjectStorageUserAccessKey`
- **Breaking**, Managed Object Storage: `UpdatedAt` field removed from `ManagedObjectStorageUserAccessKey`

### Changed
- **Breaking**, Managed Object Storage: `AccessKeyId` field in `ManagedObjectStorageUserAccessKey` renamed to `AccessKeyID`

## [7.0.0]

### Added
- Managed Load Balancer: `MaintenanceDOW` and `MaintenanceTime` fields for controlling maintenance window occurrence
- Kubernetes: support for cluster labels.

### Changed
- **Breaking**, Managed Database: `ManagedDatabaseUserOpenSearchAccessControl` fields changed to pointers
- **Breaking**, Managed Database: `ManagedDatabaseUserPGAccessControl` fields changed to pointers
- **Breaking**, Managed Database: `ManagedDatabaseUserRedisAccessControl` fields changed to pointers
- **Breaking**, Managed Load Balancer: `LoadBalancerFrontendProperties` field `InboundProxyProtocol` to pointer
- **Breaking**, Managed Object Storage: `CreateManagedObjectStorageUserAccessKeyRequest` field `Enabled` to pointer
- **Breaking**, Managed Object Storage: `ModifyManagedObjectStorageUserAccessKeyRequest` field `Enabled` to pointer
- **Breaking** Kubernetes: the `ControlPlaneIPFilter` of `ModifyKubernetesCluster` is changed from `[]string` to `*[]string`.

### Removed
- **Breaking**, Managed Database: connection related methods in favor of session
- **Breaking**: remove `Timeout` option from `WaitFor*` methods. Use `context.WithTimeout` to define a timeout for these functions.
- **Breaking**: Managed Database: `Type` field from `CloneManagedDatabaseRequest` and `ModifyManagedDatabaseRequest`

## [6.12.0]

### Added
- Server groups: Add `AddServerToServerGroup` and `RemoveServerFromServerGroup` methods.
- Storages: Add support for encryption at rest

### Fixed
- Managed Object Storage: use correct path for `GetManagedObjectStorageBucketMetricsRequest`

## [6.11.0]

### Added
- Managed Database sub-properties support. E.g., PostgreSQL property `timescaledb` is of type `object` and has `max_background_workers` sub-property.
- Managed Object Storage: add `name` property
- Account: add `ManagedObjectStorages` field to `ResourceLimits`

### Fixed
- Managed Object Storage: omit empty labels slice when creating managed object storage instance

## [6.10.0]

### Added
- Managed Load Balancer: `TLSConfigs` field to `LoadBalancerBackend` to control backend TLS configurations
- Managed Load Balancer: `TLSEnabled`,  `TLSUseSystemCA`, `TLSVerify` & `HTTP2Enabled` fields to `LoadBalancerBackendProperties`
- Managed Load Balancer: `HTTP2Enabled` field to `LoadBalancerFrontendProperties`

## [6.9.0]
### Added
- kubernetes: add `Version` field to `request.CreateKubernetesClusterRequest` and `upcloud.KubernetesCluster`

### Changed
- **Breaking**, kubernetes: update `GetKubernetesVersions` return value from `[]string` to `[]upcloud.KubernetesVersion`. (No major version bump, because this end-point has not been included in the API docs)

## [6.8.3]
### Added
- kubernetes: `WaitForKubernetesNodeGroupState` method for waiting the node group to achieve a desired state

### Changed
- kubernetes: sleep before `GET` request in `WaitForKubernetesClusterState`

## [6.8.2]
### Added
- account: `NetworkPeerings`, `NTPExcessGiB`, `StorageMaxIOPS`, and `LoadBalancers` fields to the `ResourceLimits` struct.

## [6.8.1]
### Changed
- kubernetes: `WaitForKubernetesClusterState` ignores two first 404 responses to avoid errors caused by possible false _not found_ responses right after cluster has been created.

## [6.8.0]
- Managed Object Storage support

## [6.7.1]

### Changed
- `.gitignore` editorial commit

## [6.7.0]

### Added
- network: `dhcp_routes` field to IP network for additional DHCP classless static routes to be delivered if the DHCP is enabled
- network: `static_routes` field to router for defining static routes

## [6.6.0]

### Added
- kubernetes: `control_plane_ip_filter` field to cluster for configuring and reading IP addresses and ranges allowed to access cluster control-plane
- gateway: `addresses` list to provide IP addresses assigned to the gateway

## [6.5.0]

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
- raise minimum required go version to 1.13
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

[Unreleased]: https://github.com/UpCloudLtd/upcloud-go-api/compare/v8.28.0...HEAD
[8.28.0]: https://github.com/UpCloudLtd/upcloud-go-api/compare/v8.27.0...v8.28.0
[8.27.0]: https://github.com/UpCloudLtd/upcloud-go-api/compare/v8.26.0...v8.27.0
[8.26.0]: https://github.com/UpCloudLtd/upcloud-go-api/compare/v8.25.0...v8.26.0
[8.25.0]: https://github.com/UpCloudLtd/upcloud-go-api/compare/v8.24.0...v8.25.0
[8.24.0]: https://github.com/UpCloudLtd/upcloud-go-api/compare/v8.23.0...v8.24.0
[8.23.0]: https://github.com/UpCloudLtd/upcloud-go-api/compare/v8.22.0...v8.23.0
[8.22.0]: https://github.com/UpCloudLtd/upcloud-go-api/compare/v8.21.0...v8.22.0
[8.21.0]: https://github.com/UpCloudLtd/upcloud-go-api/compare/v8.20.0...v8.21.0
[8.20.0]: https://github.com/UpCloudLtd/upcloud-go-api/compare/v8.19.0...v8.20.0
[8.19.0]: https://github.com/UpCloudLtd/upcloud-go-api/compare/v8.18.0...v8.19.0
[8.18.0]: https://github.com/UpCloudLtd/upcloud-go-api/compare/v8.17.0...v8.18.0
[8.17.0]: https://github.com/UpCloudLtd/upcloud-go-api/compare/v8.16.1...v8.17.0
[8.16.1]: https://github.com/UpCloudLtd/upcloud-go-api/compare/v8.16.0...v8.16.1
[8.16.0]: https://github.com/UpCloudLtd/upcloud-go-api/compare/v8.15.0...v8.16.0
[8.15.0]: https://github.com/UpCloudLtd/upcloud-go-api/compare/v8.14.0...v8.15.0
[8.14.0]: https://github.com/UpCloudLtd/upcloud-go-api/compare/v8.13.0...v8.14.0
[8.13.0]: https://github.com/UpCloudLtd/upcloud-go-api/compare/v8.12.0...v8.13.0
[8.12.0]: https://github.com/UpCloudLtd/upcloud-go-api/compare/v8.11.0...v8.12.0
[8.11.0]: https://github.com/UpCloudLtd/upcloud-go-api/compare/v8.10.0...v8.11.0
[8.10.0]: https://github.com/UpCloudLtd/upcloud-go-api/compare/v8.9.0...v8.10.0
[8.9.0]: https://github.com/UpCloudLtd/upcloud-go-api/compare/v8.8.1...v8.9.0
[8.8.1]: https://github.com/UpCloudLtd/upcloud-go-api/compare/v8.8.0...v8.8.1
[8.8.0]: https://github.com/UpCloudLtd/upcloud-go-api/compare/v8.7.1...v8.8.0
[8.7.1]: https://github.com/UpCloudLtd/upcloud-go-api/compare/v8.7.0...v8.7.1
[8.7.0]: https://github.com/UpCloudLtd/upcloud-go-api/compare/v8.6.2...v8.7.0
[8.6.2]: https://github.com/UpCloudLtd/upcloud-go-api/compare/v8.6.1...v8.6.2
[8.6.1]: https://github.com/UpCloudLtd/upcloud-go-api/compare/v8.6.0...v8.6.1
[8.6.0]: https://github.com/UpCloudLtd/upcloud-go-api/compare/v8.5.0...v8.6.0
[8.5.0]: https://github.com/UpCloudLtd/upcloud-go-api/compare/v8.4.0...v8.5.0
[8.4.0]: https://github.com/UpCloudLtd/upcloud-go-api/compare/v8.3.0...v8.4.0
[8.3.0]: https://github.com/UpCloudLtd/upcloud-go-api/compare/v8.2.0...v8.3.0
[8.2.0]: https://github.com/UpCloudLtd/upcloud-go-api/compare/v8.1.0...v8.2.0
[8.1.0]: https://github.com/UpCloudLtd/upcloud-go-api/compare/v8.0.0...v8.1.0
[8.0.0]: https://github.com/UpCloudLtd/upcloud-go-api/compare/v7.0.0...v8.0.0
[7.0.0]: https://github.com/UpCloudLtd/upcloud-go-api/compare/v6.12.0...v7.0.0
[6.12.0]: https://github.com/UpCloudLtd/upcloud-go-api/compare/v6.11.0...v6.12.0
[6.11.0]: https://github.com/UpCloudLtd/upcloud-go-api/compare/v6.10.0...v6.11.0
[6.10.0]: https://github.com/UpCloudLtd/upcloud-go-api/compare/v6.9.0...v6.10.0
[6.9.0]: https://github.com/UpCloudLtd/upcloud-go-api/compare/v6.8.3...v6.9.0
[6.8.3]: https://github.com/UpCloudLtd/upcloud-go-api/compare/v6.8.2...v6.8.3
[6.8.2]: https://github.com/UpCloudLtd/upcloud-go-api/compare/v6.8.1...v6.8.2
[6.8.1]: https://github.com/UpCloudLtd/upcloud-go-api/compare/v6.8.0...v6.8.1
[6.8.0]: https://github.com/UpCloudLtd/upcloud-go-api/compare/v6.7.1...v6.8.0
[6.7.1]: https://github.com/UpCloudLtd/upcloud-go-api/compare/v6.7.0...v6.7.1
[6.7.0]: https://github.com/UpCloudLtd/upcloud-go-api/compare/v6.6.0...v6.7.0
[6.6.0]: https://github.com/UpCloudLtd/upcloud-go-api/compare/v6.5.0...v6.6.0
[6.5.0]: https://github.com/UpCloudLtd/upcloud-go-api/compare/v6.4.0...v6.5.0
[6.4.0]: https://github.com/UpCloudLtd/upcloud-go-api/compare/v6.3.2...v6.4.0
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
