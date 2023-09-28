# Releasing

1. Merge all your changes to the stable branch
2. If releasing a new major version, ensure that package name has been updated, e.g. if new version is `v6` package name in go.mod and every import should be `github.com/UpCloudLtd/upcloud-go-api/v6`
3. Update CHANGELOG.md
   1. Add new heading with the correct version e.g. `## [6.7.0]`
   2. Update links at the bottom of the page
   3. Leave `## Unreleased` section at the top empty
4. Update `Version` in [upcloud/client/client.go](./upcloud/client/client.go)
5. Test GoReleaser config with `goreleaser check`
6. Tag a commit with the version you want to release e.g. `v6.7.0`
7. Push the tag & commit to GitHub
   - GitHub action automatically
      - sets the version based on the tag
      - creates a draft release to GitHub
      - populates the release notes from `CHANGELOG.md` with `make release-notes`
8. Verify that [release notes](https://github.com/UpCloudLtd/upcloud-go-api/releases) are in line with `CHANGELOG.MD`
9. Publish the drafted release
