package service

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/dnaeon/go-vcr/recorder"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud/request"
)

// TestCreateModifyDeleteStorage performs the following actions:
//
// - creates a new storage disk
// - modifies the storage
// - deletes the storage
func TestCreateModifyDeleteStorage(t *testing.T) {
	t.Parallel()

	record(t, "createmodifydeletestorage", func(t *testing.T, rec *recorder.Recorder, svc *Service) {
		// Create some storage
		storageDetails, err := createStorage(svc)
		require.NoError(t, err)
		t.Logf("Storage %s with UUID %s created", storageDetails.Title, storageDetails.UUID)

		// Modify the storage
		t.Log("Modifying the storage ...")

		newTitle := "New fancy title"
		storageDetails, err = svc.ModifyStorage(&request.ModifyStorageRequest{
			UUID:  storageDetails.UUID,
			Title: newTitle,
		})
		require.NoError(t, err)
		assert.Equal(t, newTitle, storageDetails.Title)
		t.Logf("Storage with UUID %s modified successfully, new title is %s", storageDetails.UUID, storageDetails.Title)

		// Delete the storage
		t.Log("Deleting the storage ...")
		err = deleteStorage(svc, storageDetails.UUID)
		require.NoError(t, err)
		t.Log("Storage is now deleted")
	})
}

// TestAttachDetachStorage performs the following actions:
//
// - creates a server
// - stops the server
// - creates a new storage disk
// - attaches the storage
// - detaches the storage
// - deletes the storage
// - deletes the server
func TestAttachDetachStorage(t *testing.T) {
	t.Parallel()

	record(t, "attachdetachstorage", func(t *testing.T, rec *recorder.Recorder, svc *Service) {
		// Create a server
		serverDetails, err := createServer(svc, "TestAttachDetachStorage")
		require.NoError(t, err)
		t.Logf("Server %s with UUID %s created", serverDetails.Title, serverDetails.UUID)

		// Stop the server
		t.Logf("Stopping server with UUID %s ...", serverDetails.UUID)
		err = stopServer(svc, serverDetails.UUID)
		require.NoError(t, err)
		t.Log("Server is now stopped")

		// Create some storage
		storageDetails, err := createStorage(svc)
		require.NoError(t, err)
		t.Logf("Storage %s with UUID %s created", storageDetails.Title, storageDetails.UUID)

		// Attach the storage
		t.Logf("Attaching storage %s", storageDetails.UUID)

		serverDetails, err = svc.AttachStorage(&request.AttachStorageRequest{
			StorageUUID: storageDetails.UUID,
			ServerUUID:  serverDetails.UUID,
			Type:        upcloud.StorageTypeDisk,
			Address:     "scsi:0:0",
		})
		require.NoError(t, err)
		t.Logf("Storage attached to server with UUID %s", serverDetails.UUID)

		// Detach the storage
		t.Logf("Detaching storage %s", storageDetails.UUID)

		_, err = svc.DetachStorage(&request.DetachStorageRequest{
			ServerUUID: serverDetails.UUID,
			Address:    "scsi:0:0",
		})
		require.NoError(t, err)
		t.Logf("Storage %s detached", storageDetails.UUID)
	})
}

// TestCloneStorage performs the following actions:
//
// - creates a storage device
// - clones the storage device
// - deletes the clone and the storage device
func TestCloneStorage(t *testing.T) {
	t.Parallel()

	record(t, "clonestorage", func(t *testing.T, rec *recorder.Recorder, svc *Service) {
		// Create storage
		storageDetails, err := createStorage(svc)
		require.NoError(t, err)
		t.Logf("Storage %s with UUID %s created", storageDetails.Title, storageDetails.UUID)

		// Clone the storage
		t.Log("Cloning storage ...")

		clonedStorageDetails, err := svc.CloneStorage(&request.CloneStorageRequest{
			UUID:  storageDetails.UUID,
			Title: "Cloned storage",
			Zone:  "fi-hel2",
			Tier:  upcloud.StorageTierMaxIOPS,
		})
		require.NoError(t, err)
		err = waitForStorageOnline(svc, clonedStorageDetails.UUID)
		require.NoError(t, err)
		t.Logf("Storage cloned as %s", clonedStorageDetails.UUID)
	})
}

// TestTemplatizeServerStorage performs the following actions:
//
// - creates a server
// - templatizes the server's storage
// - deletes the new storage
// - stops and deletes the server
func TestTemplatizeServerStorage(t *testing.T) {
	t.Parallel()

	record(t, "templatizeserverstorage", func(t *testing.T, rec *recorder.Recorder, svc *Service) {
		// Create server
		serverDetails, err := createServer(svc, "TestTemplatizeServerStorage")
		require.NoError(t, err)
		t.Logf("Server %s with UUID %s created", serverDetails.Title, serverDetails.UUID)

		// Stop the server
		t.Logf("Stopping server with UUID %s ...", serverDetails.UUID)
		err = stopServer(svc, serverDetails.UUID)
		require.NoError(t, err)
		t.Log("Server is now stopped")

		// Get extended service details
		serverDetails, err = svc.GetServerDetails(&request.GetServerDetailsRequest{
			UUID: serverDetails.UUID,
		})
		require.NoError(t, err)

		// Templatize the server's first storage device
		require.NotEmpty(t, serverDetails.StorageDevices)
		t.Log("Templatizing storage ...")

		storageDetails, err := svc.TemplatizeStorage(&request.TemplatizeStorageRequest{
			UUID:  serverDetails.StorageDevices[0].UUID,
			Title: "Templatized storage",
		})
		require.NoErrorf(t, err, "Error: %#v", err)

		err = waitForStorageOnline(svc, storageDetails.UUID)
		require.NoError(t, err)
		t.Logf("Storage templatized as %s", storageDetails.UUID)
	})
}

// TestLoadEjectCDROM performs the following actions:
//
// - creates a server
// - stops the server
// - attaches a CD-ROM device
// - loads a CD-ROM
// - ejects the CD-ROM
func TestLoadEjectCDROM(t *testing.T) {
	t.Parallel()

	record(t, "loadejectcdrom", func(t *testing.T, rec *recorder.Recorder, svc *Service) {
		// Create the server
		serverDetails, err := createServer(svc, "TestLoadEjectCDROM")
		require.NoError(t, err)
		t.Logf("Server %s with UUID %s created", serverDetails.Title, serverDetails.UUID)

		// Stop the server
		t.Logf("Stopping server with UUID %s ...", serverDetails.UUID)
		err = stopServer(svc, serverDetails.UUID)
		require.NoError(t, err)
		t.Log("Server is now stopped")

		// Attach CD-ROM device
		t.Logf("Attaching CD-ROM device to server with UUID %s", serverDetails.UUID)
		_, err = svc.AttachStorage(&request.AttachStorageRequest{
			ServerUUID: serverDetails.UUID,
			Type:       upcloud.StorageTypeCDROM,
		})
		require.NoError(t, err)
		t.Log("CD-ROM is now attached")

		// Load the CD-ROM
		t.Log("Loading CD-ROM into CD-ROM device")
		_, err = svc.LoadCDROM(&request.LoadCDROMRequest{
			ServerUUID:  serverDetails.UUID,
			StorageUUID: "01000000-0000-4000-8000-000030060101",
		})
		require.NoError(t, err)
		t.Log("CD-ROM is now loaded")

		// Eject the CD-ROM
		t.Log("Ejecting CD-ROM from CD-ROM device")
		_, err = svc.EjectCDROM(&request.EjectCDROMRequest{
			ServerUUID: serverDetails.UUID,
		})
		require.NoError(t, err)
		t.Log("CD-ROM is now ejected")
	})
}

// TestCreateRestoreBackup performs the following actions:
//
// - creates a storage device
// - creates a backup of the storage device
// - gets backup storage details
// - restores the backup
//
func TestCreateRestoreBackup(t *testing.T) {
	t.Parallel()

	record(t, "createrestorebackup", func(t *testing.T, rec *recorder.Recorder, svc *Service) {
		// Create the storage
		storageDetails, err := createStorage(svc)
		require.NoError(t, err)
		t.Logf("Storage %s with UUID %s created", storageDetails.Title, storageDetails.UUID)

		// Create a backup
		t.Logf("Creating backup of storage with UUID %s ...", storageDetails.UUID)

		timeBeforeBackup, err := utcTimeWithSecondPrecision()
		require.NoError(t, err)

		// Because we are recording the API tests we need to store the 'before'
		// time for the later check. We're storing it in the Title field.
		backupDetails, err := svc.CreateBackup(&request.CreateBackupRequest{
			UUID:  storageDetails.UUID,
			Title: fmt.Sprintf("backup-%d", timeBeforeBackup.UnixNano()),
		})
		require.NoError(t, err)

		err = waitForStorageOnline(svc, storageDetails.UUID)
		require.NoError(t, err)

		t.Logf("Created backup with UUID %s", backupDetails.UUID)

		// Get backup storage details
		t.Logf("Getting details of backup storage with UUID %s ...", backupDetails.UUID)

		backupStorageDetails, err := svc.GetStorageDetails(&request.GetStorageDetailsRequest{
			UUID: backupDetails.UUID,
		})
		require.NoError(t, err)

		assert.Equalf(
			t,
			backupStorageDetails.Origin,
			storageDetails.UUID,
			"The origin UUID %s of backup storage UUID %s does not match the actual origin UUID %s",
			backupStorageDetails.Origin,
			backupDetails.UUID,
			storageDetails.UUID,
		)
		t.Logf("Backup storage origin UUID OK")

		err = svc.RestoreBackup(&request.RestoreBackupRequest{
			UUID: backupDetails.UUID,
		})
		assert.NoError(t, err)

		err = waitForStorageOnline(svc, backupDetails.Origin)
		require.NoError(t, err)
	})
}

func TestStorageImport(t *testing.T) {
	record(t, "storageimport", func(t *testing.T, rec *recorder.Recorder, svc *Service) {
		storage, err := svc.CreateStorage(&request.CreateStorageRequest{
			Size:  10,
			Tier:  upcloud.StorageTierMaxIOPS,
			Zone:  "fi-hel2",
			Title: "Alpine Linux (test)",
		})
		require.NoError(t, err)

		const sha256sum string = "fd805e748f1950a34e354dc8fdfdf2f883237d65f5cdb8bcb47c64b0561d97a5"

		_, err = svc.CreateStorageImport(&request.CreateStorageImportRequest{
			StorageUUID:    storage.UUID,
			Source:         upcloud.StorageImportSourceHTTPImport,
			SourceLocation: "http://dl-cdn.alpinelinux.org/alpine/v3.12/releases/x86/alpine-standard-3.12.0-x86.iso",
		})
		require.NoError(t, err)

		afterStorageImportDetails, err := svc.WaitForStorageImportCompletion(&request.WaitForStorageImportCompletionRequest{
			StorageUUID: storage.UUID,
			Timeout:     15 * time.Minute,
		})
		require.NoError(t, err)
		require.Equal(t, upcloud.StorageImportStateCompleted, afterStorageImportDetails.State)
		require.Equal(t, sha256sum, afterStorageImportDetails.SHA256Sum)
	})
}

func TestDirectUploadStorageImport(t *testing.T) {
	record(t, "directuploadstorageimport", func(t *testing.T, rec *recorder.Recorder, svc *Service) {
		storage, err := svc.CreateStorage(&request.CreateStorageRequest{
			Size:  10,
			Tier:  upcloud.StorageTierMaxIOPS,
			Zone:  "fi-hel2",
			Title: "Direct Upload (test)",
		})
		require.NoError(t, err)

		// Test for an error if SourceLocation is missing
		_, err = svc.CreateStorageImport(&request.CreateStorageImportRequest{
			StorageUUID: storage.UUID,
			Source:      upcloud.StorageImportSourceDirectUpload,
		})
		require.Error(t, err)
		assert.EqualError(t, err, "SourceLocation must be specified")

		// Test for an error if file doesn't exist
		_, err = svc.CreateStorageImport(&request.CreateStorageImportRequest{
			StorageUUID:    storage.UUID,
			Source:         upcloud.StorageImportSourceDirectUpload,
			SourceLocation: "/this/file/doesnt/exists.txt",
		})
		require.Error(t, err)
		assert.EqualError(t, err, "unable to open SourceLocation: open /this/file/doesnt/exists.txt: no such file or directory")

		// Make temporary file
		buf := make([]byte, 100000000)
		sum := sha256.Sum256(buf)
		sha256sum := hex.EncodeToString(sum[:])

		tempf, err := ioutil.TempFile(os.TempDir(), "temp_file.txt")
		require.NoError(t, err)
		defer func() {
			if err := tempf.Close(); err == nil {
				os.Remove(tempf.Name())
			}
		}()

		_, err = svc.CreateStorageImport(&request.CreateStorageImportRequest{
			StorageUUID:    storage.UUID,
			Source:         upcloud.StorageImportSourceDirectUpload,
			SourceLocation: tempf.Name(),
		})
		require.NoError(t, err)

		afterStorageImportDetails, err := svc.WaitForStorageImportCompletion(&request.WaitForStorageImportCompletionRequest{
			StorageUUID: storage.UUID,
			Timeout:     15 * time.Minute,
		})

		require.NoError(t, err)
		require.Equal(t, upcloud.StorageImportStateCompleted, afterStorageImportDetails.State)
		require.Equal(t, sha256sum, afterStorageImportDetails.SHA256Sum)
	})
}

// TestResizeStorageFilesystem performs the following actions:
// - creates a server
// - stops the server
// - resizes the storage disk
// - resizes the storage
// - cleanup
func TestResizeStorageFilesystem(t *testing.T) {
	t.Parallel()

	record(t, "resizestoragefilesystem", func(t *testing.T, rec *recorder.Recorder, svc *Service) {
		// start server
		serverDetails, err := createMinimalServer(svc, "TestResizeStorageFilesystem")
		require.NoError(t, err)

		// stop server
		require.NoError(t, stopServer(svc, serverDetails.UUID))

		// modify disk size
		_, err = svc.ModifyStorage(&request.ModifyStorageRequest{
			UUID: serverDetails.StorageDevices[0].UUID,
			Size: serverDetails.StorageDevices[0].Size + 10,
		})
		require.NoError(t, err)

		// wait disk to become back online
		storageDetails, err := svc.WaitForStorageState(&request.WaitForStorageStateRequest{
			UUID:         serverDetails.StorageDevices[0].UUID,
			DesiredState: upcloud.StorageStateOnline,
			Timeout:      600,
		})
		require.NoError(t, err)

		// resize storage to populate new disk size
		resizeBackup, err := svc.ResizeStorageFilesystem(&request.ResizeStorageFilesystemRequest{
			UUID: storageDetails.UUID,
		})
		require.NoError(t, err)
		assert.Equal(t, storageDetails.Size, resizeBackup.Size)
		assert.Equal(t, storageDetails.UUID, resizeBackup.Origin)
		assert.Equal(t, upcloud.StorageStateOnline, resizeBackup.State)

		// cleanup
		assert.NoError(t, svc.DeleteStorage(&request.DeleteStorageRequest{UUID: resizeBackup.UUID}))
		assert.NoError(t, svc.DeleteServerAndStorages(
			&request.DeleteServerAndStoragesRequest{UUID: serverDetails.UUID}))
	})
}
