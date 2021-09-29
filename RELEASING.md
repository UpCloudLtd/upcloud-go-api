# Releasing

1. Merge all your changes to the stable branch
1. Update CHANGELOG.md
    1. Add new heading with the correct version e.g. `## [2.3.3]`
    1. Update links at the bottom of the page
    1. Leave "Unreleased” section at the top empty
1. Visit the repo [GitHub releases-page](https://github.com/UpCloudLtd/upcloud-go-api/releases) and draft a new release
1. Tag the release `X.Y.Z` (e.g. `2.3.3`)
1. Select the stable branch
1. Title the release “X.Y.Z”
1. In the description of the release, paste the changes from CHANGELOG.md for this version release
1. Publish the release when you are ready
