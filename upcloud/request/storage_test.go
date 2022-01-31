package request

import (
	"encoding/json"
	"testing"

	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud"
	"github.com/stretchr/testify/assert"
)

// TestGetStoragesRequest tests that GetStoragesRequest objects behave correctly
func TestGetStoragesRequest(t *testing.T) {
	request := GetStoragesRequest{}

	assert.Equal(t, "/storage", request.RequestURL())
	request.Access = upcloud.StorageAccessPublic
	assert.Equal(t, "/storage/public", request.RequestURL())
	request.Access = ""
	request.Favorite = true
	assert.Equal(t, "/storage/favorite", request.RequestURL())
	request.Favorite = false
	request.Type = upcloud.StorageTypeDisk
	assert.Equal(t, "/storage/disk", request.RequestURL())
}

// TestGetStorageDetailsRequest tests that GetStorageDetailsRequest objects behave correctly
func TestGetStorageDetailsRequest(t *testing.T) {
	request := GetStorageDetailsRequest{
		UUID: "foo",
	}

	assert.Equal(t, "/storage/foo", request.RequestURL())
}

// TestCreateStorageRequest tests that CreateStorageRequest objects behave correctly
func TestCreateStorageRequest(t *testing.T) {
	request := CreateStorageRequest{
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

	expectedJSON := `
	  {
		"storage": {
		  "size": "10",
		  "tier": "maxiops",
		  "title": "Test storage",
		  "zone": "fi-hel2",
		  "backup_rule": {
			"interval": "daily",
			"time": "0430",
			"retention": "30"
		  }
		}
	  }
	`
	actualJSON, err := json.MarshalIndent(&request, "", "  ")
	assert.NoError(t, err)
	assert.JSONEq(t, expectedJSON, string(actualJSON))
	assert.Equal(t, "/storage", request.RequestURL())
}

// TestModifyStorageRequest tests that ModifyStorageRequest objects behave correctly
func TestModifyStorageRequest(t *testing.T) {
	request := ModifyStorageRequest{
		UUID:  "foo",
		Title: "A larger storage",
		Size:  20,
	}

	expectedJSON := `
	  {
        "storage": {
          "size": "20",
          "title": "A larger storage"
        }
      }
	`
	actualJSON, err := json.Marshal(&request)
	assert.NoError(t, err)
	assert.JSONEq(t, expectedJSON, string(actualJSON))
	assert.Equal(t, "/storage/foo", request.RequestURL())
}

// TestAttachStorageRequest tests that AttachStorageRequest objects behave correctly
func TestAttachStorageRequest(t *testing.T) {
	request := AttachStorageRequest{
		StorageUUID: "foo",
		ServerUUID:  "bar",
		Type:        upcloud.StorageTypeDisk,
		Address:     "scsi:0:0",
		BootDisk:    1,
	}

	expectedJSON := `
	{
		"storage_device": {
		  "type": "disk",
		  "address": "scsi:0:0",
		  "storage": "foo",
		  "boot_disk": "1"
		}
	}
	`
	actualJSON, err := json.Marshal(&request)
	assert.NoError(t, err)
	assert.JSONEq(t, expectedJSON, string(actualJSON))
	assert.Equal(t, "/server/bar/storage/attach", request.RequestURL())
}

// TestDetachStorageRequest tests that DetachStorageRequest objects behave correctly
func TestDetachStorageRequest(t *testing.T) {
	request := DetachStorageRequest{
		ServerUUID: "foo",
		Address:    "scsi:0:0",
	}

	expectedJSON := `
	  {
        "storage_device": {
          "address": "scsi:0:0"
        }
      }
	`
	actualJSON, err := json.Marshal(&request)
	assert.NoError(t, err)
	assert.JSONEq(t, expectedJSON, string(actualJSON))
	assert.Equal(t, "/server/foo/storage/detach", request.RequestURL())
}

// TestDeleteStorageRequest tests that DeleteStorageRequest objects behave correctly
func TestDeleteStorageRequest(t *testing.T) {
	request := DeleteStorageRequest{
		UUID: "foo",
	}

	assert.Equal(t, "/storage/foo", request.RequestURL())
}

// TestCloneStorageRequest testa that CloneStorageRequest objects behave correctly
func TestCloneStorageRequest(t *testing.T) {
	request := CloneStorageRequest{
		UUID:  "foo",
		Title: "Clone of operating system disk",
		Zone:  "fi-hel1",
		Tier:  upcloud.StorageTierMaxIOPS,
	}

	expectedJSON := `
	{
      "storage": {
        "zone": "fi-hel1",
        "tier": "maxiops",
        "title": "Clone of operating system disk"
      }
    }
	`
	actualJSON, err := json.Marshal(&request)
	assert.NoError(t, err)
	assert.JSONEq(t, expectedJSON, string(actualJSON))
	assert.Equal(t, "/storage/foo/clone", request.RequestURL())
}

// TestTemplatizeStorageRequest tests that TemplatizeStorageRequest objects behave correctly
func TestTemplatizeStorageRequest(t *testing.T) {
	request := TemplatizeStorageRequest{
		UUID:  "foo",
		Title: "Templatized storage",
	}

	expectedJSON := `
	  {
        "storage": {
          "title": "Templatized storage"
        }
      }
	`
	actualJSON, err := json.Marshal(&request)
	assert.NoError(t, err)

	assert.JSONEq(t, expectedJSON, string(actualJSON))
	assert.Equal(t, "/storage/foo/templatize", request.RequestURL())
}

// TestLoadCDROMRequest tests that LoadCDROMRequest objects behave correctly
func TestLoadCDROMRequest(t *testing.T) {
	request := LoadCDROMRequest{
		ServerUUID:  "foo",
		StorageUUID: "bar",
	}

	expectedJSON := `
	  {
		"storage_device": {
		  "storage": "bar"
		}
	  }
	`
	actualJSON, err := json.Marshal(&request)
	assert.Nil(t, err)
	assert.JSONEq(t, expectedJSON, string(actualJSON))
	assert.Equal(t, "/server/foo/cdrom/load", request.RequestURL())
}

// TestEjectCDROMRequest tests that EjectCDROMRequest objects behave correctly
func TestEjectCDROMRequest(t *testing.T) {
	request := EjectCDROMRequest{
		ServerUUID: "foo",
	}

	assert.Equal(t, "/server/foo/cdrom/eject", request.RequestURL())
}

// TestCreateBackupRequest tests that CreateBackupRequest objects behave correctly
func TestCreateBackupRequest(t *testing.T) {
	request := CreateBackupRequest{
		UUID:  "foo",
		Title: "Manually created backup",
	}

	expectedJSON := `
	  {
		"storage": {
		  "title": "Manually created backup"
		}
	  }
	`
	actualJSON, err := json.Marshal(&request)
	assert.NoError(t, err)
	assert.JSONEq(t, expectedJSON, string(actualJSON))
	assert.Equal(t, "/storage/foo/backup", request.RequestURL())
}

// TestRestoreBackupRequest tests that RestoreBackupRequest objects behave correctly
func TestRestoreBackupRequest(t *testing.T) {
	request := RestoreBackupRequest{
		UUID: "foo",
	}

	assert.Equal(t, "/storage/foo/restore", request.RequestURL())
}

// TestStorageImportRequest tests that StorageImportRequest marshals correctly
func TestStorageImportRequest(t *testing.T) {
	request := CreateStorageImportRequest{
		StorageUUID:    "foo",
		Source:         StorageImportSourceHTTPImport,
		SourceLocation: "http://somewhere.com",
	}

	expectedJSON := `
	  {
		  "storage_import": {
			  "source": "http_import",
			  "source_location": "http://somewhere.com"
		  }
	  }
	`

	actualJSON, err := json.Marshal(&request)
	assert.NoError(t, err)

	assert.JSONEq(t, expectedJSON, string(actualJSON))

	assert.Equal(t, "/storage/foo/import", request.RequestURL())
}

// TestGetStorageImportDetails tests that GetStorageImportDetails objects behave correctly
func TestGetStorageImportDetails(t *testing.T) {
	request := GetStorageImportDetailsRequest{
		UUID: "foo",
	}

	assert.Equal(t, "/storage/foo/import", request.RequestURL())
}

// TestResizeFilesystemRequest tests that ResizeFilesystemRequest objects behave correctly
func TestResizeFilesystemRequest(t *testing.T) {
	request := ResizeFilesystemRequest{
		UUID: "foo",
	}

	assert.Equal(t, "/storage/foo/resize", request.RequestURL())
}
