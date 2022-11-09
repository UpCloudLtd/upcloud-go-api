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
			err = stopServerWithoutRecorder(svc, server.UUID)
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

func deleteLoadBalancer(svc *Service, lb *upcloud.LoadBalancer) error {
	if err := svc.DeleteLoadBalancer(&request.DeleteLoadBalancerRequest{UUID: lb.UUID}); err != nil {
		return err
	}

	if err := waitLoadBalancerToShutdown(svc, lb); err != nil {
		return fmt.Errorf("unable to shutdown LB '%s' (%s) (check dangling networks)", lb.UUID, lb.Name)
	}

	var errs []error
	if lb.NetworkUUID != "" {
		if err := svc.DeleteNetwork(&request.DeleteNetworkRequest{UUID: lb.NetworkUUID}); err != nil {
			errs = append(errs, err)
		}
	}
	if len(lb.Networks) > 0 {
		for _, n := range lb.Networks {
			if n.Type == upcloud.LoadBalancerNetworkTypePrivate && n.UUID != "" {
				if err := svc.DeleteNetwork(&request.DeleteNetworkRequest{UUID: n.UUID}); err != nil {
					errs = append(errs, err)
				}
			}
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("%s", errs)
	}
	return nil
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
