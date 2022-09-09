package service

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud/client"
	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud/request"
	"github.com/dnaeon/go-vcr/cassette"
	"github.com/dnaeon/go-vcr/recorder"
	"github.com/hashicorp/go-cleanhttp"
	"github.com/stretchr/testify/require"
)

// records the API interactions of the test. Function provides both services to test cases so that old utility functions can be used to initialize environment.
func recordWithContext(t *testing.T, fixture string, f func(context.Context, *testing.T, *recorder.Recorder, *Service, *ServiceContext)) {
	if testing.Short() {
		t.Skip("Skipping recorded test in short mode")
	}

	r, err := recorder.New("fixtures/" + fixture)
	require.NoError(t, err)

	r.AddFilter(func(i *cassette.Interaction) error {
		delete(i.Request.Headers, "Authorization")
		if i.Request.Method == http.MethodPut && strings.Contains(i.Request.URL, "uploader") {
			// We will remove the body from the upload to reduce fixture size
			i.Request.Body = ""
		}
		return nil
	})

	defer func() {
		err := r.Stop()
		require.NoError(t, err)
	}()

	user, password := getCredentials()

	httpClient := cleanhttp.DefaultClient()
	origTransport := httpClient.Transport
	r.SetTransport(origTransport)
	httpClient.Transport = r

	c := client.NewWithHTTPClient(user, password, httpClient)
	c.SetTimeout(time.Second * 300)

	customAPI := os.Getenv("UPCLOUD_GO_SDK_API_HOST")
	if customAPI != "" {
		// Override api host after the go-vcr to maintain consistent test fixtures
		r.SetTransport(&customRoundTripper{fn: func(r *http.Request) (*http.Response, error) {
			clone := r.Clone(r.Context())
			clone.URL.Host = customAPI
			clone.Host = customAPI
			return origTransport.RoundTrip(clone)
		}})
	}

	// just some random timeout value. High enough that it won't be reached during normal test.
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*60)
	defer cancel()
	f(ctx, t, r, New(c), NewWithContext(client.NewWithHTTPClientContext(user, password, httpClient)))
}

// Deletes the specified server and storages
func deleteServerAndStoragesWithContext(ctx context.Context, svc *ServiceContext, uuid string) error {
	err := svc.DeleteServerAndStorages(ctx, &request.DeleteServerAndStoragesRequest{
		UUID: uuid,
	})
	return err
}

// Creates a server and returns the details about it, panic if creation fails
func createServerWithContext(ctx context.Context, svc *ServiceContext, name string) (*upcloud.ServerDetails, error) {
	return createServerWithNetworkWithContext(ctx, svc, name, "")
}

func createServerWithNetworkWithContext(ctx context.Context, svc *ServiceContext, name string, network string) (*upcloud.ServerDetails, error) {
	title := "uploud-go-sdk-integration-test-" + name
	hostname := strings.ToLower(title + ".example.com")

	createServerRequest := request.CreateServerRequest{
		Zone:             "fi-hel2",
		Title:            title,
		Hostname:         hostname,
		PasswordDelivery: request.PasswordDeliveryNone,
		StorageDevices: []request.CreateServerStorageDevice{
			{
				Action:  request.CreateServerStorageDeviceActionClone,
				Storage: "01000000-0000-4000-8000-000020060100",
				Title:   "disk1",
				Size:    10,
				Tier:    upcloud.StorageTierMaxIOPS,
			},
		},
		Networking: &request.CreateServerNetworking{
			Interfaces: []request.CreateServerInterface{
				{
					IPAddresses: []request.CreateServerIPAddress{
						{
							Family: upcloud.IPAddressFamilyIPv4,
						},
					},
					Type: upcloud.NetworkTypeUtility,
				},
				{
					IPAddresses: []request.CreateServerIPAddress{
						{
							Family: upcloud.IPAddressFamilyIPv4,
						},
					},
					Type: upcloud.NetworkTypePublic,
				},
				{
					IPAddresses: []request.CreateServerIPAddress{
						{
							Family: upcloud.IPAddressFamilyIPv6,
						},
					},
					Type: upcloud.NetworkTypePublic,
				},
			},
		},
		Labels: &upcloud.LabelSlice{
			upcloud.Label{
				Key:   "managedBy",
				Value: "upcloud-sdk-integration-test",
			},
			upcloud.Label{
				Key:   "testName",
				Value: name,
			},
		},
	}

	if network != "" {
		createServerRequest.Networking.Interfaces = append(createServerRequest.Networking.Interfaces,
			request.CreateServerInterface{
				IPAddresses: []request.CreateServerIPAddress{
					{
						Family: upcloud.IPAddressFamilyIPv4,
					},
				},
				Type:    upcloud.NetworkTypePrivate,
				Network: network,
			})
	}

	// Create the server and block until it has started
	serverDetails, err := svc.CreateServer(ctx, &createServerRequest)
	if err != nil {
		return nil, err
	}

	// Wait for the server to start
	serverDetails, err = svc.WaitForServerState(ctx, &request.WaitForServerStateRequest{
		UUID:         serverDetails.UUID,
		DesiredState: upcloud.ServerStateStarted,
		Timeout:      time.Minute * 15,
	})
	if err != nil {
		return nil, err
	}

	return serverDetails, nil
}

// Creates a piece of storage and returns the details about it, panic if creation fails
func createStorageWithContext(ctx context.Context, svc *ServiceContext) (*upcloud.StorageDetails, error) {
	createStorageRequest := request.CreateStorageRequest{
		Tier:  upcloud.StorageTierMaxIOPS,
		Title: "Test storage",
		Size:  10,
		Zone:  "fi-hel2",
		BackupRule: &upcloud.BackupRule{
			Interval:  upcloud.BackupRuleIntervalDaily,
			Time:      "0430",
			Retention: 30,
		},
	}

	storageDetails, err := svc.CreateStorage(ctx, &createStorageRequest)
	if err != nil {
		return nil, err
	}

	return storageDetails, nil
}

// Deletes the specified storage
func deleteStorageWithContext(ctx context.Context, svc *ServiceContext, uuid string) error {
	err := svc.DeleteStorage(ctx, &request.DeleteStorageRequest{
		UUID: uuid,
	})

	return err
}

func createLoadBalancerBackendContext(ctx context.Context, svc *ServiceContext, lbUUID string) (*upcloud.LoadBalancerBackend, error) {
	req := request.CreateLoadBalancerBackendRequest{
		ServiceUUID: lbUUID,
		Backend: request.LoadBalancerBackend{
			Name: fmt.Sprintf("go-test-lb-backend-%d", time.Now().Unix()),
			Properties: &upcloud.LoadBalancerBackendProperties{
				TimeoutServer: 30,
			},
			Members: []request.LoadBalancerBackendMember{
				{
					Name:        "default-lb-backend-member",
					Type:        "dynamic",
					Weight:      100,
					MaxSessions: 1000,
					Enabled:     true,
					Port:        8000,
					IP:          "196.123.123.123",
				},
			},
		},
	}

	return svc.CreateLoadBalancerBackend(ctx, &req)
}
