package service

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud/request"
	"github.com/dnaeon/go-vcr/recorder"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestGetServerConfigurations ensures that the GetServerConfigurations() function returns proper data
func TestGetServerConfigurations(t *testing.T) {
	t.Parallel()

	record(t, "getserverconfigurations", func(t *testing.T, rec *recorder.Recorder, svc *Service) {
		configurations, err := svc.GetServerConfigurations()
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

// TestGetServersWithFilters ensures that the GetServersWithFilters() function returns proper data
func TestGetServersWithFilters(t *testing.T) {
	record(t, "getserverswithfilters", func(t *testing.T, rec *recorder.Recorder, svc *Service) {
		name := "getserverswithfilters"
		createdServer, err := createServer(svc, "getserverswithfilters")
		require.NoError(t, err)

		servers, err := svc.GetServersWithFilters(&request.GetServersWithFiltersRequest{
			Filters: []request.ServerFilter{
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

// TestGetServerDetails ensures that the GetServerDetails() function returns proper data
func TestGetServerDetails(t *testing.T) {
	t.Parallel()

	record(t, "getserverdetails", func(t *testing.T, rec *recorder.Recorder, svc *Service) {
		d, err := createServer(svc, "getserverdetails")
		require.NoError(t, err)

		serverDetails, err := svc.GetServerDetails(&request.GetServerDetailsRequest{
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

	record(t, "createstartstopserver", func(t *testing.T, rec *recorder.Recorder, svc *Service) {
		d, err := createServer(svc, "createstartstopserver")
		require.NoError(t, err)

		stopServerDetails, err := svc.StopServer(&request.StopServerRequest{
			UUID:     d.UUID,
			Timeout:  15 * time.Minute,
			StopType: upcloud.StopTypeHard,
		})
		require.NoError(t, err)
		assert.Contains(t, stopServerDetails.Title, "createstartstopserver")
		assert.Equal(t, "fi-hel2", stopServerDetails.Zone)
		// We shouldn't have transitioned state yet.
		assert.Equal(t, upcloud.ServerStateStarted, stopServerDetails.State)

		waitServerDetails, err := svc.WaitForServerState(&request.WaitForServerStateRequest{
			UUID:         d.UUID,
			DesiredState: upcloud.ServerStateStopped,
			Timeout:      15 * time.Minute,
		})
		require.NoError(t, err)
		assert.Contains(t, waitServerDetails.Title, "createstartstopserver")
		assert.Equal(t, "fi-hel2", waitServerDetails.Zone)
		assert.Equal(t, upcloud.ServerStateStopped, waitServerDetails.State)

		startServerDetails, err := svc.StartServer(&request.StartServerRequest{
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

	record(t, "startavoidhost", func(t *testing.T, rec *recorder.Recorder, svc *Service) {
		serverDetails, err := createServer(svc, "TestStartAvoidHost")
		require.NoError(t, err)
		assert.NotZero(t, serverDetails.Host)

		_, err = svc.StopServer(&request.StopServerRequest{
			UUID:     serverDetails.UUID,
			StopType: upcloud.StopTypeHard,
		})
		require.NoError(t, err)

		_, err = svc.WaitForServerState(&request.WaitForServerStateRequest{
			UUID:         serverDetails.UUID,
			DesiredState: upcloud.ServerStateStopped,
			Timeout:      15 * time.Minute,
		})
		require.NoError(t, err)

		postServerDetails, err := svc.StartServer(&request.StartServerRequest{
			UUID:      serverDetails.UUID,
			AvoidHost: serverDetails.Host,
		})
		require.NoError(t, err)
		assert.NotZero(t, postServerDetails.Host)
		assert.NotEqual(t, serverDetails.Host, postServerDetails.Host)
	})
}

// TestCreateRestartServer ensures that RestartServer() behaves as expect and returns
// proper data
// The test:
//   - Creates a server
//   - Restarts the server
//   - Checks the details of the restarted server and that it is in the
//     correct state.
func TestCreateRestartServer(t *testing.T) {
	t.Parallel()

	record(t, "createrestartserver", func(t *testing.T, rec *recorder.Recorder, svc *Service) {
		d, err := createServer(svc, "createrestartserver")
		require.NoError(t, err)

		restartServerDetails, err := svc.RestartServer(&request.RestartServerRequest{
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

		waitServerDetails, err := svc.WaitForServerState(&request.WaitForServerStateRequest{
			UUID:           d.UUID,
			UndesiredState: upcloud.ServerStateStarted,
			Timeout:        15 * time.Minute,
		})
		require.NoError(t, err)
		assert.Contains(t, waitServerDetails.Title, "createrestartserver")
		assert.Equal(t, "fi-hel2", waitServerDetails.Zone)

		waitServerDetails2, err := svc.WaitForServerState(&request.WaitForServerStateRequest{
			UUID:         waitServerDetails.UUID,
			DesiredState: upcloud.ServerStateStarted,
			Timeout:      15 * time.Minute,
		})
		require.NoError(t, err)
		assert.Contains(t, waitServerDetails2.Title, "createrestartserver")
		assert.Equal(t, "fi-hel2", waitServerDetails2.Zone)
		assert.Equal(t, upcloud.ServerStateStarted, waitServerDetails2.State)
	})
}

// TestErrorHandling checks that the correct error type is returned from service methods
func TestErrorHandling(t *testing.T) {
	t.Parallel()

	record(t, "errorhandling", func(t *testing.T, rec *recorder.Recorder, svc *Service) {
		// Perform a bogus request that will certainly fail
		_, err := svc.StartServer(&request.StartServerRequest{
			UUID: "invalid",
		})

		// Check that the correct error type is returned
		expectedErrorType := "*upcloud.Error"
		actualErrorType := reflect.TypeOf(err).String()

		if actualErrorType != expectedErrorType {
			t.Errorf("TestErrorHandling expected %s, got %s", expectedErrorType, actualErrorType)
		}
	})
}

// TestCreateModifyDeleteServer performs the following actions:
//
// - creates a server
// - modifies the server
// - stops the server
// - deletes the server
func TestCreateModifyDeleteServer(t *testing.T) {
	t.Parallel()

	record(t, "createmodifydeleteserver", func(t *testing.T, rec *recorder.Recorder, svc *Service) {
		// Create a server
		serverDetails, err := createServer(svc, "TestCreateModifyDeleteServer")
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

		_, err = svc.ModifyServer(&request.ModifyServerRequest{
			Labels: &newLabelSlice,
			UUID:   serverDetails.UUID,
			Title:  newTitle,
		})

		require.NoError(t, err)
		t.Log("Waiting for the server to exit maintenance state ...")

		serverDetails, err = svc.WaitForServerState(&request.WaitForServerStateRequest{
			UUID:         serverDetails.UUID,
			DesiredState: upcloud.ServerStateStarted,
			Timeout:      time.Minute * 15,
		})

		require.NoError(t, err)
		assert.Equal(t, newLabelSlice, serverDetails.Labels)
		assert.Equal(t, newTitle, serverDetails.Title)
		t.Logf("Server is now modified, new title is %s", serverDetails.Title)

		// Stop the server
		t.Logf("Stopping server with UUID %s ...", serverDetails.UUID)
		err = stopServer(svc, serverDetails.UUID)
		require.NoError(t, err)
		t.Log("Server is now stopped")

		// Delete the server
		t.Logf("Deleting the server with UUID %s...", serverDetails.UUID)
		err = deleteServer(svc, serverDetails.UUID)
		require.NoError(t, err)
		t.Log("Server is now deleted")

		// Check if the storage still exists
		storages, err := svc.GetStorages(&request.GetStoragesRequest{
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

// TestCreateDeleteServerAndStorage performs the following actions:
//
// - creates a server
// - deletes the server including storage
func TestCreateDeleteServerAndStorage(t *testing.T) {
	t.Parallel()

	record(t, "createdeleteserverandstorage", func(t *testing.T, rec *recorder.Recorder, svc *Service) {
		// Create a server
		serverDetails, err := createServer(svc, "TestCreateDeleteServerAndStorage")
		require.NoError(t, err)
		t.Logf("Server %s with UUID %s created", serverDetails.Title, serverDetails.UUID)

		// Get details about the storage (UUID is required for testing)
		assert.NotEmptyf(t, serverDetails.StorageDevices, "Server %s with UUID %s has no storages attached", serverDetails.Title, serverDetails.UUID)

		firstStorage := serverDetails.StorageDevices[0]
		storageUUID := firstStorage.UUID
		t.Logf("First storage of server with UUID %s has UUID %s", serverDetails.UUID, storageUUID)

		// Stop the server
		t.Logf("Stopping server with UUID %s ...", serverDetails.UUID)
		err = stopServer(svc, serverDetails.UUID)
		require.NoError(t, err)
		t.Log("Server is now stopped")

		// Delete the server and storage
		t.Logf("Deleting the server with UUID %s, including storages...", serverDetails.UUID)
		err = deleteServerAndStorages(svc, serverDetails.UUID)
		require.NoError(t, err)
		t.Log("Server is now deleted")

		// Check if the storage was deleted
		storages, err := svc.GetStorages(&request.GetStoragesRequest{
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
