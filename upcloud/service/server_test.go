package service

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/UpCloudLtd/upcloud-go-api/v6/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v6/upcloud/request"
	"github.com/dnaeon/go-vcr/recorder"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestGetServerConfigurations ensures that the GetServerConfigurations() function returns proper data.
func TestGetServerConfigurations(t *testing.T) {
	t.Parallel()

	record(t, "getserverconfigurations", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service) {
		configurations, err := svc.GetServerConfigurations(ctx)
		require.NoError(t, err)
		assert.NotEmpty(t, configurations.ServerConfigurations)

		var found bool
		for _, sc := range configurations.ServerConfigurations {
			if sc.CoreNumber == 1 && sc.MemoryAmount == 1024 {
				found = true
				break
			}
		}
		assert.True(t, found)
	})
}

// TestGetServersWithFilters ensures that the GetServersWithFilters() function returns proper data.
func TestGetServersWithFilters(t *testing.T) {
	record(t, "getserverswithfilters", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service) {
		name := "getserverswithfilters"
		createdServer, err := createServer(ctx, rec, svc, "getserverswithfilters")
		require.NoError(t, err)

		servers, err := svc.GetServersWithFilters(ctx, &request.GetServersWithFiltersRequest{
			Filters: []request.QueryFilter{
				request.FilterLabelKey{Key: "managedBy"},
				request.FilterLabel{Label: upcloud.Label{
					Key:   "testName",
					Value: name,
				}},
			},
		})
		require.NoError(t, err)

		var found bool
		for _, s := range servers.Servers {
			fmt.Println(s.Title)
			if s.Title == createdServer.Title {
				found = true

				break
			}
		}
		assert.True(t, found)
	})
}

// TestGetServerDetails ensures that the GetServerDetails() function returns proper data.
func TestGetServerDetails(t *testing.T) {
	t.Parallel()

	record(t, "getserverdetails", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service) {
		d, err := createServer(ctx, rec, svc, "getserverdetails")
		require.NoError(t, err)

		serverDetails, err := svc.GetServerDetails(ctx, &request.GetServerDetailsRequest{
			UUID: d.UUID,
		})
		require.NoError(t, err)

		assert.Contains(t, serverDetails.Title, "getserverdetails")
		assert.Equal(t, "fi-hel2", serverDetails.Zone)
	})
}

// TestCreateStopStartServer ensures that StartServer() and StopServer() behave
// as expect and return proper data
// The test:
//   - Creates a server
//   - Stops the server
//   - Starts the server
//   - Checks the details of the started server and that it is in the
//     correct state.
func TestCreateStopStartServer(t *testing.T) {
	t.Parallel()

	record(t, "createstartstopserver", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service) {
		d, err := createServer(ctx, rec, svc, "createstartstopserver")
		require.NoError(t, err)

		stopServerDetails, err := svc.StopServer(ctx, &request.StopServerRequest{
			UUID:     d.UUID,
			Timeout:  15 * time.Minute,
			StopType: upcloud.StopTypeHard,
		})
		require.NoError(t, err)
		assert.Contains(t, stopServerDetails.Title, "createstartstopserver")
		assert.Equal(t, "fi-hel2", stopServerDetails.Zone)
		// We shouldn't have transitioned state yet.
		assert.Equal(t, upcloud.ServerStateStarted, stopServerDetails.State)

		err = waitForServerState(ctx, rec, svc, d.UUID, upcloud.ServerStateStopped)
		require.NoError(t, err)

		getServerDetails, err := svc.GetServerDetails(ctx, &request.GetServerDetailsRequest{
			UUID: d.UUID,
		})
		require.NoError(t, err)
		assert.Contains(t, getServerDetails.Title, "createstartstopserver")
		assert.Equal(t, "fi-hel2", getServerDetails.Zone)
		assert.Equal(t, upcloud.ServerStateStopped, getServerDetails.State)

		startServerDetails, err := svc.StartServer(ctx, &request.StartServerRequest{
			UUID: d.UUID,
		})
		require.NoError(t, err)

		assert.Contains(t, startServerDetails.Title, "createstartstopserver")
		assert.Equal(t, "fi-hel2", startServerDetails.Zone)
		assert.Equal(t, upcloud.ServerStateStarted, startServerDetails.State)
	})
}

func TestStartAvoidHost(t *testing.T) {
	t.Parallel()

	record(t, "startavoidhost", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service) {
		serverDetails, err := createServer(ctx, rec, svc, "TestStartAvoidHost")
		require.NoError(t, err)
		assert.NotZero(t, serverDetails.Host)

		_, err = svc.StopServer(ctx, &request.StopServerRequest{
			UUID:     serverDetails.UUID,
			StopType: upcloud.StopTypeHard,
		})
		require.NoError(t, err)

		err = waitForServerState(ctx, rec, svc, serverDetails.UUID, upcloud.ServerStateStopped)
		require.NoError(t, err)

		getServerDetails, err := svc.GetServerDetails(ctx, &request.GetServerDetailsRequest{
			UUID: serverDetails.UUID,
		})
		require.NoError(t, err)
		assert.Equal(t, upcloud.ServerStateStopped, getServerDetails.State)

		postServerDetails, err := svc.StartServer(ctx, &request.StartServerRequest{
			UUID:      serverDetails.UUID,
			AvoidHost: serverDetails.Host,
		})
		require.NoError(t, err)
		assert.NotZero(t, postServerDetails.Host)
		assert.NotEqual(t, serverDetails.Host, postServerDetails.Host)
	})
}

// TestCreateRestartServerContext ensures that RestartServer() behaves as expect and returns
// proper data
// The test:
//   - Creates a server
//   - Restarts the server
//   - Checks the details of the restarted server and that it is in the
//     correct state.
func TestCreateRestartServer(t *testing.T) {
	t.Parallel()

	record(t, "createrestartserver", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service) {
		t.Log("create server")
		d, err := createServer(ctx, rec, svc, "createrestartserver")
		require.NoError(t, err)

		t.Log("restart server, state should be started")
		restartServerDetails, err := svc.RestartServer(ctx, &request.RestartServerRequest{
			UUID:          d.UUID,
			Timeout:       15 * time.Minute,
			StopType:      upcloud.StopTypeHard,
			TimeoutAction: request.RestartTimeoutActionIgnore,
		})
		require.NoError(t, err)
		assert.Contains(t, restartServerDetails.Title, "createrestartserver")
		assert.Equal(t, "fi-hel2", restartServerDetails.Zone)
		// We shouldn't have transitioned state yet.
		assert.Equal(t, upcloud.ServerStateStarted, restartServerDetails.State)

		err = waitForServerState(ctx, rec, svc, d.UUID, upcloud.ServerStateMaintenance)
		require.NoError(t, err)

		t.Log("get server, state should be maintenance")
		getServerDetails, err := svc.GetServerDetails(ctx, &request.GetServerDetailsRequest{UUID: d.UUID})
		require.NoError(t, err)
		assert.Equal(t, upcloud.ServerStateMaintenance, getServerDetails.State)

		err = waitForServerState(ctx, rec, svc, d.UUID, upcloud.ServerStateStarted)
		require.NoError(t, err)

		t.Log("get server, state should be started")
		getServerDetails2, err := svc.GetServerDetails(ctx, &request.GetServerDetailsRequest{UUID: d.UUID})
		require.NoError(t, err)
		assert.Equal(t, upcloud.ServerStateStarted, getServerDetails2.State)
	})
}

// TestErrorHandlingContext checks that the correct error type is returned from service methods.
func TestErrorHandling(t *testing.T) {
	t.Parallel()

	record(t, "errorhandling", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service) {
		// Perform a bogus request that will certainly fail
		_, err := svc.StartServer(ctx, &request.StartServerRequest{
			UUID: "invalid",
		})

		// Check that the correct error type is returned
		expectedErrorType := "*upcloud.Problem"
		actualErrorType := reflect.TypeOf(err).String()

		if actualErrorType != expectedErrorType {
			t.Errorf("TestErrorHandling expected %s, got %s", expectedErrorType, actualErrorType)
		}
	})
}

// TestCreateModifyDeleteServerContext performs the following actions:
//
// - creates a server
// - modifies the server
// - stops the server
// - deletes the server.
func TestCreateModifyDeleteServer(t *testing.T) {
	t.Parallel()

	record(t, "createmodifydeleteserver", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service) {
		// Create a server
		serverDetails, err := createServer(ctx, rec, svc, "TestCreateModifyDeleteServer")
		require.NoError(t, err)
		t.Logf("Server %s with UUID %s created", serverDetails.Title, serverDetails.UUID)

		// Get details about the storage (UUID is required for testing)
		if len(serverDetails.StorageDevices) == 0 {
			t.Errorf("Server %s with UUID %s has no storages attached", serverDetails.Title, serverDetails.UUID)
		}

		firstStorage := serverDetails.StorageDevices[0]
		storageUUID := firstStorage.UUID

		t.Logf("First storage of server with UUID %s has UUID %s", serverDetails.UUID, storageUUID)

		// Modify the server
		t.Log("Modifying the server ...")
		newTitle := "Modified server"
		newLabelSlice := append(serverDetails.Labels, upcloud.Label{Key: "title", Value: newTitle})
		_, err = svc.ModifyServer(ctx, &request.ModifyServerRequest{
			Labels: &newLabelSlice,
			UUID:   serverDetails.UUID,
			Title:  newTitle,
		})

		require.NoError(t, err)
		t.Log("Waiting for the server to exit maintenance state ...")

		err = waitForServerState(ctx, rec, svc, serverDetails.UUID, upcloud.ServerStateStarted)
		require.NoError(t, err)

		getServerDetails, err := svc.GetServerDetails(ctx, &request.GetServerDetailsRequest{UUID: serverDetails.UUID})

		require.NoError(t, err)
		assert.Equal(t, newLabelSlice, getServerDetails.Labels)
		assert.Equal(t, newTitle, getServerDetails.Title)
		t.Logf("Server is now modified, new title is %s", getServerDetails.Title)

		// Stop the server
		t.Logf("Stopping server with UUID %s ...", serverDetails.UUID)
		err = stopServer(ctx, rec, svc, serverDetails.UUID)
		require.NoError(t, err)
		t.Log("Server is now stopped")

		// Delete the server
		t.Logf("Deleting the server with UUID %s...", serverDetails.UUID)
		err = deleteServer(ctx, svc, serverDetails.UUID)
		require.NoError(t, err)
		t.Log("Server is now deleted")

		// Check if the storage still exists
		storages, err := svc.GetStorages(ctx, &request.GetStoragesRequest{
			Access: upcloud.StorageAccessPrivate,
		})
		require.NoError(t, err)

		found := false
		for _, storage := range storages.Storages {
			if storage.UUID == storageUUID {
				found = true
				break
			}
		}
		assert.Truef(t, found, "Storage with UUID %s not found. It should still exist after deleting server with UUID %s", storageUUID, serverDetails.UUID)

		t.Log("Storage still exists")
	})
}

// TestCreateDeleteServerAndStorageContext performs the following actions:
//
// - creates a server
// - deletes the server including storage.
func TestCreateDeleteServerAndStorage(t *testing.T) {
	t.Parallel()

	record(t, "createdeleteserverandstorage", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service) {
		// Create a server
		serverDetails, err := createServer(ctx, rec, svc, "TestCreateDeleteServerAndStorage")
		require.NoError(t, err)
		t.Logf("Server %s with UUID %s created", serverDetails.Title, serverDetails.UUID)

		// Get details about the storage (UUID is required for testing)
		assert.NotEmptyf(t, serverDetails.StorageDevices, "Server %s with UUID %s has no storages attached", serverDetails.Title, serverDetails.UUID)

		firstStorage := serverDetails.StorageDevices[0]
		storageUUID := firstStorage.UUID
		t.Logf("First storage of server with UUID %s has UUID %s", serverDetails.UUID, storageUUID)

		// Stop the server
		t.Logf("Stopping server with UUID %s ...", serverDetails.UUID)
		err = stopServer(ctx, rec, svc, serverDetails.UUID)
		require.NoError(t, err)
		t.Log("Server is now stopped")

		// Delete the server and storage
		t.Logf("Deleting the server with UUID %s, including storages...", serverDetails.UUID)
		err = deleteServerAndStorages(ctx, svc, serverDetails.UUID)
		require.NoError(t, err)
		t.Log("Server is now deleted")

		// Check if the storage was deleted
		storages, err := svc.GetStorages(ctx, &request.GetStoragesRequest{
			Access: upcloud.StorageAccessPrivate,
		})
		require.NoError(t, err)

		found := false
		for _, storage := range storages.Storages {
			if storage.UUID == storageUUID {
				found = true
				break
			}
		}
		assert.Falsef(t, found, "Storage with UUID %s still exists. It should have been deleted with server with UUID %s", storageUUID, serverDetails.UUID)

		t.Log("Storage was deleted, too")
	})
}

// Creates a minimal server with a private utility network interface.
func createMinimalServer(ctx context.Context, rec *recorder.Recorder, svc *Service, name string) (*upcloud.ServerDetails, error) {
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

	// Create the server
	serverDetails, err := svc.CreateServer(ctx, &createServerRequest)
	if err != nil {
		return nil, err
	}

	// Wait for the server to start
	err = waitForServerState(ctx, rec, svc, serverDetails.UUID, upcloud.ServerStateStarted)
	if err != nil {
		return nil, err
	}

	return svc.GetServerDetails(ctx, &request.GetServerDetailsRequest{
		UUID: serverDetails.UUID,
	})
}

// Creates a server and returns the details about it, panic if creation fails.
func createServer(ctx context.Context, rec *recorder.Recorder, svc *Service, name string) (*upcloud.ServerDetails, error) {
	return createServerWithNetwork(ctx, rec, svc, name, "")
}

// Creates a server with a network.
func createServerWithNetwork(ctx context.Context, rec *recorder.Recorder, svc *Service, name, network string) (*upcloud.ServerDetails, error) {
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

	err = waitForServerState(ctx, rec, svc, serverDetails.UUID, upcloud.ServerStateStarted)
	if err != nil {
		return nil, err
	}

	return svc.GetServerDetails(ctx, &request.GetServerDetailsRequest{
		UUID: serverDetails.UUID,
	})
}

// Deletes the specified server and storages.
func deleteServerAndStorages(ctx context.Context, svc *Service, uuid string) error {
	err := svc.DeleteServerAndStorages(ctx, &request.DeleteServerAndStoragesRequest{
		UUID: uuid,
	})
	return err
}

// Stops the specified server (forcibly).
func stopServer(ctx context.Context, rec *recorder.Recorder, svc *Service, uuid string) error {
	_, err := svc.StopServer(ctx, &request.StopServerRequest{
		UUID:     uuid,
		Timeout:  waitTimeout,
		StopType: request.ServerStopTypeHard,
	})
	if err != nil {
		return err
	}

	return waitForServerState(ctx, rec, svc, uuid, upcloud.ServerStateStopped)
}

// Waits for the server to achieve the desired state.
func waitForServerState(ctx context.Context, rec *recorder.Recorder, svc *Service, serverUUID string, desiredState string) error {
	if rec.Mode() != recorder.ModeRecording {
		return nil
	}

	rec.AddPassthrough(func(h *http.Request) bool {
		return true
	})
	defer func() {
		rec.Passthroughs = nil
	}()

	// Wait for the server to start
	_, err := svc.WaitForServerState(ctx, &request.WaitForServerStateRequest{
		UUID:         serverUUID,
		DesiredState: desiredState,
		Timeout:      waitTimeout,
	})
	if err != nil {
		return err
	}

	return nil
}

// Deletes the specified server.
func deleteServer(ctx context.Context, svc *Service, uuid string) error {
	return svc.DeleteServer(ctx, &request.DeleteServerRequest{
		UUID: uuid,
	})
}

// Stops the specified server (forcibly).
func stopServerWithoutRecorder(ctx context.Context, svc *Service, uuid string) error {
	serverDetails, err := svc.StopServer(ctx, &request.StopServerRequest{
		UUID:     uuid,
		Timeout:  waitTimeout,
		StopType: request.ServerStopTypeHard,
	})
	if err != nil {
		return err
	}

	_, err = svc.WaitForServerState(ctx, &request.WaitForServerStateRequest{
		UUID:         serverDetails.UUID,
		DesiredState: upcloud.ServerStateStopped,
		Timeout:      waitTimeout,
	})

	return err
}
