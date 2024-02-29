# Developing

If you spot an issue, or have an idea for a feature that would make your work with this SDK easier - don't hesitate to open an issue or make a pull requests. We are grateful for all feedback and contributions.

## Code style & linting

We use [golangci-lint](https://github.com/golangci/golangci-lint) for linting and formatting. Please run `golangci-lint run ./...` before making a PR.

## Commit Messages

Please follow [conventional commits](https://www.conventionalcommits.org/en/v1.0.0/).

## Testing

To be able to run the test suite you'll need to export the following environment variables with their corresponding
values:

* `UPCLOUD_GO_SDK_TEST_USER` (the API username)
* `UPCLOUD_GO_SDK_TEST_PASSWORD` (the API password)
* `UPCLOUD_GO_SDK_TEST_DELETE_RESOURCES` (either `yes` or `no`)

To run the test suite, run `go test ./... -v -parallel 8`. If `UPCLOUD_GO_SDK_TEST_DELETE_RESOURCES` is set to `yes`,
all resources will be stopped and/or deleted after the test suite has run. Be careful which account you use for
testing so you don't accidentally delete or your production resources!

You can skip running the integration tests and just run the unit tests by passing `-short` to the test command.

## Debugging

Environment variables `UPCLOUD_DEBUG_API_BASE_URL` and `UPCLOUD_DEBUG_SKIP_CERTIFICATE_VERIFY` can be used for HTTP client debugging purposes.
* `UPCLOUD_DEBUG_API_BASE_URL` overrides static base URL. This can be used with local server to debug request problems.  
E.g. `UPCLOUD_DEBUG_API_BASE_URL=http://127.0.0.1:8080`
* `UPCLOUD_DEBUG_SKIP_CERTIFICATE_VERIFY` skips server's certificate verification. If set to `1`, API client accepts any certificate presented by the server and any host name in that certificate.  
E.g. `UPCLOUD_DEBUG_SKIP_CERTIFICATE_VERIFY=1`
