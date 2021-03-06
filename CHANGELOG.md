# Change log

All notable changes to this project will be documented in this file.
See updating [Changelog example here](https://keepachangelog.com/en/1.0.0/)

## 4.0.0

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

## 3.0.0

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

## 2.0.0

- moved project to UpCloud's own GitHub organization
- raise the minimum required Go version to 1.7

## 1.1.0

- improve documentation
- remove ability to override the API base URL and version

## 1.0.2

- remove credentials related getters and setters from `Client`
- hopefully fix CD-ROM integration tests for good by performing all operations while the server is stopped

## 1.0.1

- minor tweaks to the integration tests
- correct the package name in the README installation instructions

## 1.0.0

First stable release
