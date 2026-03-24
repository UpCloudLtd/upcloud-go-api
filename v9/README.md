# UpCloud Go API client library

[![Go Report Card](https://goreportcard.com/badge/github.com/UpCloudLtd/upcloud-go-api)](https://goreportcard.com/report/github.com/UpCloudLtd/upcloud-go-api)
[![OpenSSF Scorecard](https://api.scorecard.dev/projects/github.com/UpCloudLtd/upcloud-go-api/badge)](https://scorecard.dev/viewer/?uri=github.com%2FUpCloudLtd%2Fupcloud-go-api)
[![Go Reference](https://pkg.go.dev/badge/github.com/UpCloudLtd/upcloud-go-api/v9/pkg/upcloud.svg)](https://pkg.go.dev/github.com/UpCloudLtd/upcloud-go-api/v9/pkg/upcloud)

This is the official client for interfacing with UpCloud's API using the Go programming language. The features allows for easy and quick development and integration when using Go.

## Usage

### Quickstart

Add the module to your project:

```shell
go get github.com/UpCloudLtd/upcloud-go-api/v9
```

Create a client and call the API. Use an API token, or load credentials from the environment or system keyring (see below).

```go
package main

import (
	"context"
	
	"github.com/UpCloudLtd/upcloud-go-api/v9/pkg/upcloud"
)

func main() {
	client, err := upcloud.New("your_api_token")
	if err != nil {
		return
	}

	// Or: client, err := upcloud.NewFromEnv()
	// That uses UPCLOUD_TOKEN or UPCLOUD_USERNAME / UPCLOUD_PASSWORD (and optional keyring),
	// same as github.com/UpCloudLtd/upcloud-go-api/credentials.

	ctx := context.Background()
	resp, err := client.ListObjectStoragesWithResponse(ctx, nil)
	if err != nil {
		return
	}
	if resp.JSON200 == nil {
		// See "Error handling" below.
		return
	}
	_ = resp.JSON200 // use the decoded success payload
}
```

For integration tests or local debugging you can point the client at another API base URL and optionally skip TLS verification (not for production):

- `UPCLOUD_DEBUG_API_BASE_URL` — override the server URL
- `UPCLOUD_DEBUG_SKIP_CERTIFICATE_VERIFY=1` — only when using `NewFromEnv` (see `pkg/upcloud/client.go`)

### Error handling

Methods named `*WithResponse` return a typed response struct and an `error`. The `error` is non-nil for transport-level failures (network, TLS, timeouts). A successful HTTP exchange is still represented with a non-nil response value; inspect the typed fields for the status code and decoded body.

- Success bodies are in fields such as `JSON200`, `JSON201`, depending on the operation.
- API error bodies are often parsed into `ApplicationproblemJSONDefault` (and similarly named fields on each response type), typically as `ObjectStorage2ErrorResponse` with `Status`, `Title`, `CorrelationId`, `InvalidParams`, and `Type` (see the [UpCloud API documentation](https://developers.upcloud.com/1.3) for error semantics).

```go
resp, err := client.ListObjectStoragesWithResponse(ctx, nil)
if err != nil {
	// Connection or client-side failure.
	return err
}
if resp.JSON200 != nil {
	// Success: use *resp.JSON200
	return nil
}
if resp.ApplicationproblemJSONDefault != nil {
	p := resp.ApplicationproblemJSONDefault
	_ = p.Status
	_ = p.Title
	_ = p.CorrelationId
	_ = p.Type
	if p.InvalidParams != nil {
		for _, inv := range *p.InvalidParams {
			_, _ = inv.Name, inv.Reason
		}
	}
	return nil
}
// Unexpected status: check resp.StatusCode() and resp.HTTPResponse.
```

The `pkg/upcloud` package also provides helpers that poll until Object Storage reaches a given operational state or is deleted (`WaitForObjectStorageOperationalState`, `WaitForObjectStorageDeletion`, `WaitForObjectStorageBucketDeletion`).

### Packages

- `pkg/upcloud` — generated `Client` / `ClientWithResponses`, request and response types, and small extensions (`client.go`, `object_storage.go`, `retry.go`) for authentication and wait helpers.

### Examples

The following assumes you already have a `*upcloud.ClientWithResponses` named `client`.

#### Listing object storage services

```go
resp, err := client.ListObjectStoragesWithResponse(ctx, nil)
if err != nil {
	panic(err)
}
if resp.JSON200 == nil {
	panic("unexpected response")
}
for _, svc := range *resp.JSON200 {
	if svc.Name != nil {
		println(*svc.Name)
	}
}
```

A full create–wait–delete flow, including long-running waits, is in [examples/object_storage.go](examples/object_storage.go).

## Documentation

- Go package reference: [pkg.go.dev](https://pkg.go.dev/github.com/UpCloudLtd/upcloud-go-api/v9)
- Human-oriented API reference: [UpCloud API docs](https://developers.upcloud.com/1.3)

Regenerating the SDK and release mechanics are described in [CONTRIBUTING.md](CONTRIBUTING.md).

## License

This client is distributed under the [MIT License](https://opensource.org/licenses/MIT); see [LICENSE](LICENSE) for the full text.
