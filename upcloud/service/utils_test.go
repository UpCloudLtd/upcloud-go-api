package service

import (
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/client"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/dnaeon/go-vcr/cassette"
	"github.com/dnaeon/go-vcr/recorder"
	"github.com/hashicorp/go-cleanhttp"
	"github.com/stretchr/testify/require"
)

// Configures the test environment
func getService() *Service {
	user, password := getCredentials()

	c := client.New(user, password)
	c.SetTimeout(time.Second * 300)

	return New(c)
}

// records the API interactions of the test
func record(t *testing.T, fixture string, f func(*testing.T, *Service)) {
	r, err := recorder.New("fixtures/" + fixture)
	require.NoError(t, err)

	r.AddFilter(func(i *cassette.Interaction) error {
		delete(i.Request.Headers, "Authorization")
		return nil
	})

	defer func() {
		err := r.Stop()
		require.NoError(t, err)
	}()

	user, password := getCredentials()

	httpClient := cleanhttp.DefaultClient()
	httpClient.Transport = r

	c := client.NewWithHTTPClient(user, password, httpClient)
	c.SetTimeout(time.Second * 300)

	f(t, New(c))
}

// Tears down the test environment by removing all resources
func teardown() {
	svc := getService()

	log.Print("Deleting all servers ...")
	servers, err := svc.GetServers()
	handleError(err)

	for _, server := range servers.Servers {
		// Try to ensure the server is not in maintenance state
		log.Printf("Waiting for server with UUID %s to leave maintenance state ...", server.UUID)
		serverDetails, err := svc.WaitForServerState(&request.WaitForServerStateRequest{
			UUID:           server.UUID,
			UndesiredState: upcloud.ServerStateMaintenance,
			Timeout:        time.Minute * 15,
		})
		handleError(err)

		// Stop the server if it's still running
		if serverDetails.State != upcloud.ServerStateStopped {
			log.Printf("Stopping server with UUID %s ...", server.UUID)
			stopServer(svc, server.UUID)
		}

		// Delete the server
		log.Printf("Deleting the server with UUID %s ...", server.UUID)
		deleteServer(svc, server.UUID)
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
				Timeout:      time.Minute * 15,
			})
			handleError(err)
		}

		log.Printf("Deleting the storage with UUID %s ...", storage.UUID)
		deleteStorage(svc, storage.UUID)
	}

	// Delete all tags
	log.Print("Deleting all tags ...")
	deleteAllTags(svc)

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
}

// Creates a server and returns the details about it, panic if creation fails
func createServer(svc *Service, name string) (*upcloud.ServerDetails, error) {
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
				Storage: "01000000-0000-4000-8000-000030060200",
				Title:   "disk1",
				Size:    30,
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
		Timeout:      time.Minute * 15,
	})
	if err != nil {
		return nil, err
	}

	return serverDetails, nil
}

// Stops the specified server (forcibly)
func stopServer(svc *Service, uuid string) error {
	serverDetails, err := svc.StopServer(&request.StopServerRequest{
		UUID:     uuid,
		Timeout:  time.Minute * 15,
		StopType: request.ServerStopTypeHard,
	})
	if err != nil {
		return err
	}

	_, err = svc.WaitForServerState(&request.WaitForServerStateRequest{
		UUID:         serverDetails.UUID,
		DesiredState: upcloud.ServerStateStopped,
		Timeout:      time.Minute * 15,
	})
	if err != nil {
		return err
	}

	return nil
}

// Deletes the specified server
func deleteServer(svc *Service, uuid string) error {
	err := svc.DeleteServer(&request.DeleteServerRequest{
		UUID: uuid,
	})

	return err
}

// Deletes the specified server and storages
func deleteServerAndStorages(svc *Service, uuid string) error {
	err := svc.DeleteServerAndStorages(&request.DeleteServerAndStoragesRequest{
		UUID: uuid,
	})

	return err
}

// Creates a piece of storage and returns the details about it, panic if creation fails
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

// Deletes the specified storage
func deleteStorage(svc *Service, uuid string) error {
	err := svc.DeleteStorage(&request.DeleteStorageRequest{
		UUID: uuid,
	})

	return err
}

// deleteAllTags deletes all existing tags
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

// Waits for the specified storage to come online
func waitForStorageOnline(svc *Service, uuid string) error {
	_, err := svc.WaitForStorageState(&request.WaitForStorageStateRequest{
		UUID:         uuid,
		DesiredState: upcloud.StorageStateOnline,
		Timeout:      time.Minute * 15,
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

// Handles the error by panicing, thus stopping the test execution
func handleError(err error) {
	if err != nil {
		panic(err)
	}
}

// Reads the API username and password from the environment, panics if they are not available
func getCredentials() (string, string) {
	user := os.Getenv("UPCLOUD_GO_SDK_TEST_USER")
	password := os.Getenv("UPCLOUD_GO_SDK_TEST_PASSWORD")

	if user == "" || password == "" {
		panic("Unable to retrieve credentials from the environment, ensure UPCLOUD_GO_SDK_TEST_USER and UPCLOUD_GO_SDK_TEST_PASSWORD are exported")
	}

	return user, password
}
