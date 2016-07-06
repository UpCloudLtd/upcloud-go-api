# upcloud-go-sdk

[![Build Status](https://travis-ci.org/Jalle19/upcloud-go-sdk.svg?branch=master)](https://travis-ci.org/Jalle19/upcloud-go-sdk)
[![Go Report Card](https://goreportcard.com/badge/github.com/jalle19/upcloud-go-sdk)](https://goreportcard.com/report/github.com/jalle19/upcloud-go-sdk)

This is an SDK for interfacing with Upcloud's API using the Go programming language. It is loosely based on similar 
SDKs such as https://github.com/aws/aws-sdk-go.

## Installation and requirements

You'll need Go 1.6 or higher to use the SDK. You can use the following command to retrieve the SDK:

```
go get -u github.com/jalle19/upcloud-go-sdk
```

## Usage

The general pattern for using the SDK goes like this:

* Create a `client.Client`
* Create a `service.Service` by passing the newly created client object to it
* Interface with the API using the various methods of the service object. Methods that take parameters wrap them in 
request objects.

### Creating the client and the service

```go
// Upcloud doesn't use dedicated API keys, instead you pass your account login credentials to the client
c := client.New(user, password)

// It is generally a good idea to override the default timeout of the underlying HTTP client since some requests block for longer periods of time
c.SetTimeout(time.Second * 30)

// Create the service object
svc := New(c)
```

The rest of these examples assume you already have a service object configured and named `svc`.

### Retrieving a list of servers

```go
// Retrieve the list of servers
servers, err := svc.GetServers()

if err != nil {
	panic(err)
}

// Print the UUID and hostname of each server
for _, server := range servers.Servers {
	fmt.Println(fmt.Sprintf("UUID: %s, hostname: %s", server.UUID, server.Hostname))
}
```

## License

This SDK is distributed under the [MIT License](https://opensource.org/licenses/MIT), see LICENSE.txt for more information.
