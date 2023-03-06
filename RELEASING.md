# Releasing

1. Merge all your changes to the stable branch
1. If releasing a new major version, ensure that package name has been updated, e.g. if new version is `v6` package name in go.mod and every import should be `github.com/UpCloudLtd/upcloud-go-api/v6`
1. Update CHANGELOG.md
    1. Add new heading with the correct version e.g. `## [v2.3.5]`
    1. Update links at the bottom of the page
    1. Leave “Unreleased” section at the top empty
1. Update `Version` in [upcloud/client/client.go](./upcloud/client/client.go)
1. Visit the repo [GitHub releases-page](https://github.com/UpCloudLtd/upcloud-go-api/releases) and draft a new release
1. Tag the release `vX.Y.Z` (e.g. `v2.3.5`)
1. Select the stable branch
1. Title the release “vX.Y.Z”
1. In the description of the release, paste the changes from CHANGELOG.md for this version release
1. Publish the release when you are ready
