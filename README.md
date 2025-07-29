# UpCloud Go API client library

![Build Status](https://github.com/UpCloudLtd/upcloud-go-api/workflows/Upcloud%20go%20api%20test/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/UpCloudLtd/upcloud-go-api)](https://goreportcard.com/report/github.com/UpCloudLtd/upcloud-go-api)
[![OpenSSF Scorecard](https://api.scorecard.dev/projects/github.com/UpCloudLtd/upcloud-go-api/badge)](https://scorecard.dev/viewer/?uri=github.com%2FUpCloudLtd%2Fupcloud-go-api)
[![GoDoc](https://godoc.org/github.com/UpCloudLtd/upcloud-go-api?status.svg)](https://godoc.org/github.com/UpCloudLtd/upcloud-go-api)

This is the official client for interfacing with UpCloud's API using the Go programming language. The features allows for easy and quick development and integration when using Go.

## Usage

### Quickstart

Add the library to your project:
```shell
go get github.com/UpCloudLtd/upcloud-go-api/v8
```

Next in your code:
```go
import (
	"time"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/client"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/service"
)

func main() {
	// First instantiate a client with your username and password.
	// `client.New` accepts config functions that allow you to customise the returned client behaviour.
	//  Config functions are exported from the `client` package.
	clnt := client.New("my_username", "password123", client.WithTimeout(time.Second * 30))

	// Next instantiate new service using the created client
	svc := service.New(clnt)

	// Validate that everything is set up correctly
	account, err := svc.GetAccount(context.Background())
}
```

### Error handling

```go
import (
	"context"
	"errors"
	"fmt"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/client"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/service"
)

func main() {
	svc := service.New(client.New("my_username", "password123"))
	_, err := svc.GetAccount(context.Background())

	if err != nil {
		// `upcloud.Problem` is the error object returned by all of the `Service` methods.
		//  You can differentiate between generic connection errors (like the API not being reachable) and service errors, which are errors returned in the response body by the API;
		//	this is useful for gracefully recovering from certain types of errors;
		var problem *upcloud.Problem

		if errors.As(err, &problem) {
			fmt.Println(problem.Status)        // HTTP status code returned by the API
			fmt.Print(problem.Title)           // Short, human-readable description of the problem
			fmt.Println(problem.CorrelationID) // Unique string that identifies the request that caused the problem; note that this field is not always populated
			fmt.Println(problem.InvalidParams) // List of invalid request parameters

			for _, invalidParam := range problem.InvalidParams {
				fmt.Println(invalidParam.Name)   // Path to the request field that is invalid
				fmt.Println(invalidParam.Reason) // Human-readable description of the problem with that particular field
			}

			// You can also check against the specific api error codes to programatically react to certain situations.
			// Base `upcloud` package exports all the error codes that API can return.
			// You can check which error code is return in which situation in UpCloud API docs -> https://developers.upcloud.com/1.3
			if problem.ErrorCode() == upcloud.ErrCodeResourceAlreadyExists {
				fmt.Println("Looks like we don't need to create this")
			}

			// `upcloud.Problem` implements the Error interface, so you can also just use it as any other error
			fmt.Println(fmt.Errorf("we got an error from the UpCloud API: %w", problem))
		} else {
			// This means you got an error, but it does not come from the API itself. This can happen, for example, if you have some connection issues,
			// or if the UpCloud API is unreachable for some other reason
			fmt.Println("We got a generic error!")
		}
	}
}
```

### Packages

UpCloud Go SDK includes the following packages:
- `upcloud` package - contains type definitions for all UpCloud API objects like servers, storages, load balancers, Kubernetes clusters, errors, etc. It also has a lot of constants that allow you, for example, to compare state, status and other properties of various objects.
- `client` package - contains functions that allow you to create and customise HTTP client that will be used to make requests to UpCloud API. The returned client does expose some methods for making requests, but you shouldn't really use them directly, client should only be used to instantiate a new `Service`
- `service` package - contains the `Service` type, which exposes all the methods to interact with UpCloud API. This is the package you will probably use most frequently. All `Service` methods accept `context.Context` as firt parameter. _Most_ `Service` methods accept a `request` object as the second parameter (see package below).
- `request` package - contains various `request` objects. Those objects should always be used as an argument for a `Service` method and allow you to provide additional params for the request URL or body. For example, when fetching details of a speficic server, you would use a request object to speficy the server UUID. Similarly, when creating server you would use request object to specify server properties, like CPU, memory, OS, login method, etc.

### Examples

All of these examples assume you already have a service object configured and named `svc`.

#### Retrieving a list of servers

The following example will retrieve a list of servers the account has access to.

```go
// Retrieve the list of servers
servers, err := svc.GetServers(context.Background())

if err != nil {
	panic(err)
}

// Print the UUID and hostname of each server
for _, server := range servers.Servers {
	fmt.Println(fmt.Sprintf("UUID: %s, hostname: %s", server.UUID, server.Hostname))
}
```

#### Creating a new server

Since the request for creating a new server is asynchronous, the server will report its status as "maintenance" until the deployment has been fully completed.

```go
// Create the server
serverDetails, err := svc.CreateServer(context.Background(), &request.CreateServerRequest{
	Zone:             "fi-hel2",
	Title:            "My new server",
	Hostname:         "server.example.com",
	PasswordDelivery: request.PasswordDeliveryNone,
	StorageDevices: []request.CreateServerStorageDevice{
		{
			Action:  request.CreateStorageDeviceActionClone,
			Storage: "01000000-0000-4000-8000-000020060100",
			Title:   "disk1",
			Size:    10,
			Tier:    upcloud.StorageTierMaxIOPS,
		},
	},
	IPAddresses: []request.CreateServerIPAddress{
		{
			Access: upcloud.IPAddressAccessPrivate,
			Family: upcloud.IPAddressFamilyIPv4,
		},
		{
			Access: upcloud.IPAddressAccessPublic,
			Family: upcloud.IPAddressFamilyIPv4,
		},
		{
			Access: upcloud.IPAddressAccessPublic,
			Family: upcloud.IPAddressFamilyIPv6,
		},
	},
})

if err != nil {
	panic(err)
}

fmt.Println(fmt.Sprintf("Server %s with UUID %s created", serverDetails.Title, serverDetails.UUID))

// Block for up to five minutes until the server has entered the "started" state
err = svc.WaitForServerState(context.Background(), &request.WaitForServerStateRequest{
	UUID:         serverDetails.UUID,
	DesiredState: upcloud.ServerStateStarted,
})

if err != nil {
	panic(err)
}

fmt.Println("Server is now started")
```

### Templatizing a server's storage device

In this example, we assume that there is a server represented by the variable `serverDetails` and that the server state is `stopped`. The next piece of code allows you to templatize the first storage device of the server.

```go
// Loop through the storage devices
for i, storage := range serverDetails.StorageDevices {
	// Find the first device
	if i == 0 {
		// Templatize the storage
		storageDetails, err := svc.TemplatizeStorage(context.Background(), &request.TemplatizeStorageRequest{
			UUID:  storage.UUID,
			Title: "Templatized storage",
		})

		if err != nil {
			panic(err)
		}

		fmt.Println(fmt.Sprintf("Storage templatized as %s", storageDetails.UUID))
		break
	}
}
```

### Create a manual backup

In this example, we assume that there is a storage device represented by `storageDetails` and that if it is attached to any server, the server is stopped.

```go
backupDetails, err := svc.CreateBackup(context.Background(), &request.CreateBackupRequest{
	UUID:  storageDetails.UUID,
	Title: "Backup",
})

if err != nil {
    panic(err)
}

fmt.Println(fmt.Sprintf("Backup of %s created as %s", storageDetails.UUID, backupDetails.UUID))
```

### Create a new firewall rule

In this example, we assume that there is a server represented by the variable `serverDetails`.

```go
firewallRule, err := svc.CreateFirewallRule(context.Background(), &request.CreateFirewallRuleRequest{
	ServerUUID: serverDetails.UUID,
	FirewallRule: upcloud.FirewallRule{
		Direction: upcloud.FirewallRuleDirectionIn,
		Action:    upcloud.FirewallRuleActionAccept,
		Family:    upcloud.IPAddressFamilyIPv4,
		Protocol:  upcloud.FirewallRuleProtocolTCP,
		Position:  1,
		Comment:   "Accept all TCP input on IPv4",
	},
})

if err != nil {
    panic(err)
}
```

For more examples, please consult the service integration test suite (`upcloud/service/service_test.go`).

## License

This client is distributed under the [MIT License](https://opensource.org/licenses/MIT), see LICENSE.txt for more information.
