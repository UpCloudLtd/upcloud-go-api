package service

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/dnaeon/go-vcr/cassette"
	"github.com/dnaeon/go-vcr/recorder"
	"github.com/hashicorp/go-cleanhttp"
	"github.com/stretchr/testify/require"

	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud/client"
	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud/request"
)

const waitTimeout = time.Minute * 15

type customRoundTripper struct {
	fn func(r *http.Request) (*http.Response, error)
}

func (c *customRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	return c.fn(r)
}

// Configures the test environment.
func getService() *Service {
	user, password := getCredentials()

	c := client.New(user, password)
	c.SetTimeout(time.Second * 300)

	return New(c)
}

// records the API interactions of the test.
func record(t *testing.T, fixture string, f func(*testing.T, *recorder.Recorder, *Service)) {
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

	f(t, r, New(c))
}

// Tears down the test environment by removing all resources.
func teardown() {
	svc := getService()

	// Delete all server groups
	log.Print("Deleting all server groups ...")
	serverGroups, err := svc.GetServerGroups(&request.GetServerGroupsRequest{})
	handleError(err)

	for _, serverGroup := range serverGroups {
		log.Printf("Deleting the server group with UUID %s ...", serverGroup.UUID)
		err = deleteServerGroup(svc, serverGroup.UUID)
		handleError(err)
	}

	log.Print("Deleting all servers ...")
	servers, err := svc.GetServers()
	handleError(err)

	for _, server := range servers.Servers {
		// Try to ensure the server is not in maintenance state
		log.Printf("Waiting for server with UUID %s to leave maintenance state ...", server.UUID)
		serverDetails, err := svc.WaitForServerState(&request.WaitForServerStateRequest{
			UUID:           server.UUID,
			UndesiredState: upcloud.ServerStateMaintenance,
			Timeout:        waitTimeout,
		})
		handleError(err)

		// Stop the server if it's still running
		if serverDetails.State != upcloud.ServerStateStopped {
			log.Printf("Stopping server with UUID %s ...", server.UUID)
			err = stopServer(svc, server.UUID)
			handleError(err)
		}

		// Delete the server
		log.Printf("Deleting the server with UUID %s ...", server.UUID)
		err = deleteServer(svc, server.UUID)
		handleError(err)
	}

	// Delete all private storage devices
	log.Print("Deleting all storage devices ...")
	storages, err := svc.GetStorages(&request.GetStoragesRequest{
		Access: upcloud.StorageAccessPrivate,
	})
	handleError(err)

	for _, storage := range storages.Storages {
		// Wait for the storage to come online so we can delete it
		if storage.State != upcloud.StorageStateOnline {
			log.Printf("Waiting for storage %s to come online ...", storage.UUID)
			_, err = svc.WaitForStorageState(&request.WaitForStorageStateRequest{
				UUID:         storage.UUID,
				DesiredState: upcloud.StorageStateOnline,
				Timeout:      waitTimeout,
			})
			handleError(err)
		}

		log.Printf("Deleting the storage with UUID %s ...", storage.UUID)
		err = deleteStorage(svc, storage.UUID)
		handleError(err)
	}

	// Delete all tags
	log.Print("Deleting all tags ...")
	err = deleteAllTags(svc)
	handleError(err)

	log.Print("Deleting all networks...")
	networks, err := svc.GetNetworks()
	handleError(err)
	var count int
	for _, network := range networks.Networks {
		if strings.Contains(network.Name, "(test)") {
			err := svc.DeleteNetwork(&request.DeleteNetworkRequest{
				UUID: network.UUID,
			})
			count++
			handleError(err)
		}
	}
	log.Printf("Deleted %d networks...", count)

	log.Print("Deleting all routers...")
	routers, err := svc.GetRouters()
	handleError(err)
	count = 0
	for _, router := range routers.Routers {
		if strings.Contains(router.Name, "(test)") {
			err := svc.DeleteRouter(&request.DeleteRouterRequest{
				UUID: router.UUID,
			})
			count++
			handleError(err)
		}
	}
	log.Printf("Deleted %d routers...", count)

	// Delete all object storages
	log.Print("Delete all object storages...")
	objectStorages, err := svc.GetObjectStorages()
	handleError(err)

	for _, objectStorage := range objectStorages.ObjectStorages {
		// Delete the Object Storage
		log.Printf("Deleting the object storage with UUID %s ...", objectStorage.UUID)
		err = deleteObjectStorage(svc, objectStorage.UUID)
		handleError(err)
	}
}

// Creates a server and returns the details about it, panic if creation fails.
func createServer(svc *Service, name string) (*upcloud.ServerDetails, error) {
	return createServerWithNetwork(svc, name, "")
}

func createServerWithNetwork(svc *Service, name string, network string) (*upcloud.ServerDetails, error) {
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
	serverDetails, err := svc.CreateServer(&createServerRequest)
	if err != nil {
		return nil, err
	}

	// Wait for the server to start
	serverDetails, err = svc.WaitForServerState(&request.WaitForServerStateRequest{
		UUID:         serverDetails.UUID,
		DesiredState: upcloud.ServerStateStarted,
		Timeout:      waitTimeout,
	})
	if err != nil {
		return nil, err
	}

	return serverDetails, nil
}

// Creates a server with recorder and returns the details about it, panic if creation fails.
func createServerWithRecorder(rec *recorder.Recorder, svc *Service, name string) (*upcloud.ServerDetails, error) {
	return createServerWithNetworkWithRecorder(rec, svc, name, "")
}

func createServerWithNetworkWithRecorder(rec *recorder.Recorder, svc *Service, name, network string) (*upcloud.ServerDetails, error) {
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
	serverDetails, err := svc.CreateServer(&createServerRequest)
	if err != nil {
		return nil, err
	}

	if rec.Mode() == recorder.ModeRecording {
		rec.AddPassthrough(func(h *http.Request) bool {
			return true
		})

		// Wait for the server to start
		_, err = svc.WaitForServerState(&request.WaitForServerStateRequest{
			UUID:         serverDetails.UUID,
			DesiredState: upcloud.ServerStateStarted,
			Timeout:      time.Minute * 15,
		})
		if err != nil {
			return nil, err
		}

		rec.Passthroughs = nil
	}

	return svc.GetServerDetails(&request.GetServerDetailsRequest{
		UUID: serverDetails.UUID,
	})
}

func createMinimalServer(rec *recorder.Recorder, svc *Service, name string) (*upcloud.ServerDetails, error) {
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
			},
		},
	}

	// Create the server and block until it has started
	serverDetails, err := svc.CreateServer(&createServerRequest)
	if err != nil {
		return nil, err
	}

	// Wait for the server to start
	if rec.Mode() == recorder.ModeRecording {
		rec.AddPassthrough(func(h *http.Request) bool {
			return true
		})

		serverDetails, err = svc.WaitForServerState(&request.WaitForServerStateRequest{
			UUID:         serverDetails.UUID,
			DesiredState: upcloud.ServerStateStarted,
			Timeout:      waitTimeout,
		})
		if err != nil {
			return nil, err
		}

		rec.Passthroughs = nil
	}

	return serverDetails, nil
}

// Stops the specified server (forcibly).
func stopServer(svc *Service, uuid string) error {
	serverDetails, err := svc.StopServer(&request.StopServerRequest{
		UUID:     uuid,
		Timeout:  waitTimeout,
		StopType: request.ServerStopTypeHard,
	})
	if err != nil {
		return err
	}

	_, err = svc.WaitForServerState(&request.WaitForServerStateRequest{
		UUID:         serverDetails.UUID,
		DesiredState: upcloud.ServerStateStopped,
		Timeout:      waitTimeout,
	})
	if err != nil {
		return err
	}

	return nil
}

// Stops the specified server with recorder (forcibly).
func stopServerWithRecorder(rec *recorder.Recorder, svc *Service, uuid string) error {
	serverDetails, err := svc.StopServer(&request.StopServerRequest{
		UUID:     uuid,
		Timeout:  waitTimeout,
		StopType: request.ServerStopTypeHard,
	})
	if err != nil {
		return err
	}

	if rec.Mode() == recorder.ModeRecording {
		rec.AddPassthrough(func(h *http.Request) bool {
			return true
		})

		_, err = svc.WaitForServerState(&request.WaitForServerStateRequest{
			UUID:         serverDetails.UUID,
			DesiredState: upcloud.ServerStateStopped,
			Timeout:      waitTimeout,
		})
		if err != nil {
			return err
		}

		rec.Passthroughs = nil
	}

	return nil
}

// Deletes the specified server.
func deleteServer(svc *Service, uuid string) error {
	err := svc.DeleteServer(&request.DeleteServerRequest{
		UUID: uuid,
	})

	return err
}

// Deletes the specified server and storages.
func deleteServerAndStorages(svc *Service, uuid string) error {
	err := svc.DeleteServerAndStorages(&request.DeleteServerAndStoragesRequest{
		UUID: uuid,
	})

	return err
}

// Deletes the specified server group.
func deleteServerGroup(svc *Service, uuid string) error {
	err := svc.DeleteServerGroup(&request.DeleteServerGroupRequest{
		UUID: uuid,
	})

	return err
}

// Creates a piece of storage and returns the details about it, panic if creation fails.
func createStorage(svc *Service) (*upcloud.StorageDetails, error) {
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

	storageDetails, err := svc.CreateStorage(&createStorageRequest)
	if err != nil {
		return nil, err
	}

	return storageDetails, nil
}

// Deletes the specified storage.
func deleteStorage(svc *Service, uuid string) error {
	err := svc.DeleteStorage(&request.DeleteStorageRequest{
		UUID: uuid,
	})

	return err
}

// deleteAllTags deletes all existing tags.
func deleteAllTags(svc *Service) error {
	tags, err := svc.GetTags()
	if err != nil {
		return err
	}

	for _, tagDetails := range tags.Tags {
		err = svc.DeleteTag(&request.DeleteTagRequest{
			Name: tagDetails.Name,
		})

		if err != nil {
			return err
		}
	}

	return nil
}

// Waits for the specified storage to come online.
func waitForStorageOnline(svc *Service, uuid string) error {
	_, err := svc.WaitForStorageState(&request.WaitForStorageStateRequest{
		UUID:         uuid,
		DesiredState: upcloud.StorageStateOnline,
		Timeout:      waitTimeout,
	})

	return err
}

// Creates an Object Storage and returns the details about it, panic if creation fails.
func createObjectStorage(svc *Service, name, description, zone string, size int) (*upcloud.ObjectStorageDetails, error) {
	createObjectStorageRequest := request.CreateObjectStorageRequest{
		Name:        "go-test-" + name,
		Description: description,
		Zone:        zone,
		Size:        size,
		AccessKey:   "UCOB5HE4NVTVFMXXRBQ2",
		SecretKey:   "ssLDVHvTRjHaEAPRcMiFep3HItcqdNUNtql3DcLx",
	}

	// Create the Object Storage and block until it has started
	objectStorageDetails, err := svc.CreateObjectStorage(&createObjectStorageRequest)
	if err != nil {
		return nil, err
	}

	return objectStorageDetails, nil
}

// Deletes the specific Object Storage.
func deleteObjectStorage(svc *Service, uuid string) error {
	err := svc.DeleteObjectStorage(&request.DeleteObjectStorageRequest{
		UUID: uuid,
	})

	return err
}

// Returns the current UTC time with second precision (milliseconds truncated).
// This is the format we usually get from the UpCloud API.
func utcTimeWithSecondPrecision() (time.Time, error) {
	utc, err := time.LoadLocation("UTC")
	if err != nil {
		return time.Time{}, err
	}

	t := time.Now().In(utc).Truncate(time.Second)

	return t, err
}

// Handles the error by panicing, thus stopping the test execution.
func handleError(err error) {
	if err != nil {
		panic(err)
	}
}

// Reads the API username and password from the environment, panics if they are not available.
func getCredentials() (string, string) {
	if os.Getenv("UPCLOUD_GO_SDK_TEST_NO_CREDENTIALS") == "yes" {
		return "username", "password"
	}

	user := os.Getenv("UPCLOUD_GO_SDK_TEST_USER")
	password := os.Getenv("UPCLOUD_GO_SDK_TEST_PASSWORD")

	if user == "" || password == "" {
		panic("Unable to retrieve credentials from the environment, ensure UPCLOUD_GO_SDK_TEST_USER and UPCLOUD_GO_SDK_TEST_PASSWORD are exported")
	}

	return user, password
}

func createLoadBalancerAndNetwork(svc *Service, zone, addr string) (*upcloud.LoadBalancer, error) {
	n, err := createLoadBalancerPrivateNetwork(svc, zone, addr)
	if err != nil {
		return nil, err
	}
	return createLoadBalancer(svc, n.UUID, zone)
}

func createLoadBalancer(svc *Service, networkUUID, zone string) (*upcloud.LoadBalancer, error) {
	createLoadBalancerRequest := request.CreateLoadBalancerRequest{
		Name:             fmt.Sprintf("go-test-lb-%s-%d", zone, time.Now().Unix()),
		Zone:             zone,
		Plan:             "development",
		NetworkUUID:      networkUUID,
		ConfiguredStatus: "started",
		Frontends:        []request.LoadBalancerFrontend{},
		Backends:         []request.LoadBalancerBackend{},
		Resolvers:        []request.LoadBalancerResolver{},
	}

	loadBalancerDetails, err := svc.CreateLoadBalancer(&createLoadBalancerRequest)
	if err != nil {
		return nil, err
	}

	return loadBalancerDetails, nil
}

func waitLoadBalancerToShutdown(svc *Service, lb *upcloud.LoadBalancer) error {
	const maxRetries int = 100
	// wait delete request
	for i := 0; i <= maxRetries; i++ {
		_, err := svc.GetLoadBalancer(&request.GetLoadBalancerRequest{UUID: lb.UUID})
		if err != nil {
			if svcErr, ok := err.(*upcloud.Problem); ok && svcErr.Status == http.StatusNotFound {
				return nil
			}
		}
		time.Sleep(2 * time.Second)
	}
	return errors.New("max retries reached while waiting for load balancer instance to shutdown")
}

func deleteLoadBalancer(svc *Service, lb *upcloud.LoadBalancer, waitShutdown bool) error {
	netID := lb.NetworkUUID
	if err := svc.DeleteLoadBalancer(&request.DeleteLoadBalancerRequest{UUID: lb.UUID}); err != nil {
		return err
	}
	if waitShutdown {
		if err := waitLoadBalancerToShutdown(svc, lb); err != nil {
			return fmt.Errorf("unable to shutdown LB '%s' (%s) (check dangling networks)", lb.UUID, lb.Name)
		}
	}
	return svc.DeleteNetwork(&request.DeleteNetworkRequest{UUID: netID})
}

func createLoadBalancerBackend(svc *Service, lbUUID string) (*upcloud.LoadBalancerBackend, error) {
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

	return svc.CreateLoadBalancerBackend(&req)
}

func createLoadBalancerPrivateNetwork(svc *Service, zone, addr string) (*upcloud.Network, error) {
	return svc.CreateNetwork(&request.CreateNetworkRequest{
		Name: fmt.Sprintf("go-test-lb-%d", time.Now().Unix()),
		Zone: zone,
		IPNetworks: []upcloud.IPNetwork{
			{
				Address: addr,
				DHCP:    upcloud.True,
				Family:  upcloud.IPAddressFamilyIPv4,
			},
		},
	})
}
