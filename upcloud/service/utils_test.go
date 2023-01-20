package service

import (
	"context"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/UpCloudLtd/upcloud-go-api/v5/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v5/upcloud/client"
	"github.com/UpCloudLtd/upcloud-go-api/v5/upcloud/request"
	"github.com/dnaeon/go-vcr/cassette"
	"github.com/dnaeon/go-vcr/recorder"
	"github.com/hashicorp/go-cleanhttp"
	"github.com/stretchr/testify/require"
)

const waitTimeout = time.Minute * 15

type customRoundTripper struct {
	fn func(r *http.Request) (*http.Response, error)
}

func (c *customRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	return c.fn(r)
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

// Handles the error by panicing, thus stopping the test execution.
func handleError(err error) {
	if err != nil {
		panic(err)
	}
}

// records the API interactions of the test. Function provides both services to test cases so that old utility functions can be used to initialize environment.
func record(t *testing.T, fixture string, f func(context.Context, *testing.T, *recorder.Recorder, *Service)) {
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
	ctx, cancel := context.WithTimeout(context.Background(), waitTimeout*4)
	defer cancel()
	f(ctx, t, r, New(client.New(user, password, client.WithHTTPClient(httpClient))))
}

// Tears down the test environment by removing all resources.
func teardown() {
	svc := getService()

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(1*time.Hour))
	defer cancel()

	// Delete all server groups
	log.Print("Deleting all server groups ...")
	serverGroups, err := svc.GetServerGroups(ctx, &request.GetServerGroupsRequest{})
	handleError(err)

	for _, serverGroup := range serverGroups {
		log.Printf("Deleting the server group with UUID %s ...", serverGroup.UUID)
		err = deleteServerGroup(ctx, svc, serverGroup.UUID)
		handleError(err)
	}

	log.Print("Deleting all servers ...")
	servers, err := svc.GetServers(ctx)
	handleError(err)

	for _, server := range servers.Servers {
		// Try to ensure the server is not in maintenance state
		log.Printf("Waiting for server with UUID %s to leave maintenance state ...", server.UUID)
		serverDetails, err := svc.WaitForServerState(ctx, &request.WaitForServerStateRequest{
			UUID:           server.UUID,
			UndesiredState: upcloud.ServerStateMaintenance,
			Timeout:        waitTimeout,
		})
		handleError(err)

		// Stop the server if it's still running
		if serverDetails.State != upcloud.ServerStateStopped {
			log.Printf("Stopping server with UUID %s ...", server.UUID)
			err = stopServerWithoutRecorder(ctx, svc, server.UUID)
			handleError(err)
		}

		// Delete the server
		log.Printf("Deleting the server with UUID %s ...", server.UUID)
		err = deleteServer(ctx, svc, server.UUID)
		handleError(err)
	}

	// Delete all private storage devices
	log.Print("Deleting all storage devices ...")
	storages, err := svc.GetStorages(ctx, &request.GetStoragesRequest{
		Access: upcloud.StorageAccessPrivate,
	})
	handleError(err)

	for _, storage := range storages.Storages {
		// Wait for the storage to come online so we can delete it
		if storage.State != upcloud.StorageStateOnline {
			log.Printf("Waiting for storage %s to come online ...", storage.UUID)
			_, err = svc.WaitForStorageState(ctx, &request.WaitForStorageStateRequest{
				UUID:         storage.UUID,
				DesiredState: upcloud.StorageStateOnline,
				Timeout:      waitTimeout,
			})
			handleError(err)
		}

		log.Printf("Deleting the storage with UUID %s ...", storage.UUID)
		err = deleteStorage(ctx, svc, storage.UUID)
		handleError(err)
	}

	// Delete all tags
	log.Print("Deleting all tags ...")
	err = deleteAllTags(ctx, svc)
	handleError(err)

	log.Print("Deleting all networks...")
	networks, err := svc.GetNetworks(ctx)
	handleError(err)
	var count int
	for _, network := range networks.Networks {
		if strings.Contains(network.Name, "(test)") {
			err := svc.DeleteNetwork(ctx, &request.DeleteNetworkRequest{
				UUID: network.UUID,
			})
			count++
			handleError(err)
		}
	}
	log.Printf("Deleted %d networks...", count)

	log.Print("Deleting all routers...")
	routers, err := svc.GetRouters(ctx)
	handleError(err)
	count = 0
	for _, router := range routers.Routers {
		if strings.Contains(router.Name, "(test)") {
			err := svc.DeleteRouter(ctx, &request.DeleteRouterRequest{
				UUID: router.UUID,
			})
			count++
			handleError(err)
		}
	}
	log.Printf("Deleted %d routers...", count)

	// Delete all object storages
	log.Print("Delete all object storages...")
	objectStorages, err := svc.GetObjectStorages(ctx)
	handleError(err)

	for _, objectStorage := range objectStorages.ObjectStorages {
		// Delete the Object Storage
		log.Printf("Deleting the object storage with UUID %s ...", objectStorage.UUID)
		err = deleteObjectStorage(ctx, svc, objectStorage.UUID)
		handleError(err)
	}
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

// Returns a mock server with handler for a single endpoint and a new service that targets said mock server
func setupTestServerAndService(handler http.Handler) (*httptest.Server, *Service) {
	srv := httptest.NewServer(handler)
	return srv, New(client.New("user", "pass", client.WithBaseURL(srv.URL)))
}
