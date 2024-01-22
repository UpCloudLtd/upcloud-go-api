package upcloud

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestUnmarshalStorage tests that Storages and Storage struct are unmarshaled correctly
func TestUnmarshalStorage(t *testing.T) {
	originalJSON := `
{
  "storages": {
    "storage": [
      {
        "access": "private",
        "encrypted": "yes",
        "license": 0,
        "size": 10,
        "state": "online",
        "tier": "hdd",
        "title": "Operating system disk",
        "type": "normal",
        "uuid": "01eff7ad-168e-413e-83b0-054f6a28fa23",
        "zone": "uk-lon1"
      },
      {
        "access" : "private",
        "encrypted": "no",
        "created" : "2019-09-17T14:34:43Z",
        "license" : 0,
        "origin" : "01eff7ad-168e-413e-83b0-054f6a28fa23",
        "size" : 10,
        "state" : "online",
        "title" : "On demand backup",
        "type" : "backup",
        "uuid" : "01287ad1-496c-4b5f-bb67-0fc2e3494740",
        "zone" : "uk-lon1"
      },
      {
        "access": "private",
        "license": 0,
        "part_of_plan": "yes",
        "size": 50,
        "state": "online",
        "tier": "maxiops",
        "title": "Databases",
        "type": "normal",
        "uuid": "01f3286c-a5ea-4670-8121-d0b9767d625b",
        "zone": "fi-hel1",
		"labels": [
			{
				"key": "managedBy",
				"value": "upcloud-go-sdk"
			}
		]
      }
    ]
  }
}`

	storages := Storages{}
	err := json.Unmarshal([]byte(originalJSON), &storages)
	assert.NoError(t, err)
	assert.Len(t, storages.Storages, 3)

	testData := []Storage{
		{
			Access:    StorageAccessPrivate,
			Encrypted: FromBool(true),
			License:   0.0,
			Size:      10,
			State:     StorageStateOnline,
			Tier:      StorageTierHDD,
			Title:     "Operating system disk",
			Type:      StorageTypeNormal,
			UUID:      "01eff7ad-168e-413e-83b0-054f6a28fa23",
			Zone:      "uk-lon1",
		},
		{
			Access:  StorageAccessPrivate,
			License: 0.0,
			Origin:  "01eff7ad-168e-413e-83b0-054f6a28fa23",
			Size:    10,
			State:   StorageStateOnline,
			Title:   "On demand backup",
			Type:    StorageTypeBackup,
			UUID:    "01287ad1-496c-4b5f-bb67-0fc2e3494740",
			Zone:    "uk-lon1",
		},
		{
			Access:     StorageAccessPrivate,
			License:    0.0,
			PartOfPlan: "yes",
			Size:       50,
			State:      StorageStateOnline,
			Tier:       StorageTierMaxIOPS,
			Title:      "Databases",
			Type:       StorageTypeNormal,
			UUID:       "01f3286c-a5ea-4670-8121-d0b9767d625b",
			Zone:       "fi-hel1",
			Labels: []Label{{
				Key:   "managedBy",
				Value: "upcloud-go-sdk",
			}},
		},
	}

	for i, data := range testData {
		storage := storages.Storages[i]
		assert.Equal(t, data.Access, storage.Access)
		assert.Equal(t, data.License, storage.License)
		assert.Equal(t, data.Size, storage.Size)
		assert.Equal(t, data.Title, storage.Title)
		assert.Equal(t, data.Type, storage.Type)
		assert.Equal(t, data.UUID, storage.UUID)
		assert.Equal(t, data.PartOfPlan, storage.PartOfPlan)
		assert.Equal(t, data.State, storage.State)
		assert.Equal(t, data.Tier, storage.Tier)
		assert.Equal(t, data.Zone, storage.Zone)
	}
}

// TestUnmarshalStorageDetails tests that StorageDetails struct is unmarshaled correctly
func TestUnmarshalStorageDetails(t *testing.T) {
	originalJSON := `
	{
		"storage": {
		  "access": "private",
		  "backup_rule": {
			  "interval": "daily",
			  "time": "0400",
			  "retention": "1"
		  },
		  "backups": {
			"backup": [
              "37c96670-9c02-4d5d-8f60-291d38f9a80c",
              "ecfda9f2-e071-4bbb-b38f-079ed26eb32a"
			]
		  },
		  "license": 0,
		  "servers": {
			"server": [
			  "00798b85-efdc-41ca-8021-f6ef457b8531"
			]
		  },
		  "size": 10,
		  "state": "online",
		  "tier": "maxiops",
		  "title": "Operating system disk",
		  "type": "normal",
		  "uuid": "01d4fcd4-e446-433b-8a9c-551a1284952e",
		  "zone": "fi-hel1",
		  "labels": [
			{
				"key": "managedBy",
				"value": "upcloud-go-sdk"
			}
		  ]
		}
	  }
	`

	storageDeviceDetails := StorageDetails{}
	err := json.Unmarshal([]byte(originalJSON), &storageDeviceDetails)
	assert.NoError(t, err)

	assert.Equal(t, StorageAccessPrivate, storageDeviceDetails.Access)
	assert.Equal(t, 0.0, storageDeviceDetails.License)
	assert.Equal(t, 10, storageDeviceDetails.Size)
	assert.Equal(t, StorageStateOnline, storageDeviceDetails.State)
	assert.Equal(t, StorageTierMaxIOPS, storageDeviceDetails.Tier)
	assert.Equal(t, "Operating system disk", storageDeviceDetails.Title)
	assert.Equal(t, StorageTypeNormal, storageDeviceDetails.Type)
	assert.Equal(t, "01d4fcd4-e446-433b-8a9c-551a1284952e", storageDeviceDetails.UUID)
	assert.Equal(t, "fi-hel1", storageDeviceDetails.Zone)

	assert.Equal(t, BackupRuleIntervalDaily, storageDeviceDetails.BackupRule.Interval)
	assert.Equal(t, 1, storageDeviceDetails.BackupRule.Retention)
	assert.Equal(t, "0400", storageDeviceDetails.BackupRule.Time)

	assert.Equal(t, 2, len(storageDeviceDetails.BackupUUIDs))
	assert.Equal(t, "37c96670-9c02-4d5d-8f60-291d38f9a80c", storageDeviceDetails.BackupUUIDs[0])
	assert.Equal(t, "ecfda9f2-e071-4bbb-b38f-079ed26eb32a", storageDeviceDetails.BackupUUIDs[1])

	assert.Equal(t, 1, len(storageDeviceDetails.ServerUUIDs))
	assert.Equal(t, "00798b85-efdc-41ca-8021-f6ef457b8531", storageDeviceDetails.ServerUUIDs[0])
	assert.Equal(t, 1, len(storageDeviceDetails.Labels))
	assert.Equal(t, "managedBy", storageDeviceDetails.Labels[0].Key)
	assert.Equal(t, "upcloud-go-sdk", storageDeviceDetails.Labels[0].Value)
}

// TestUnmarshalStorageImport tests that StorageImport struct is unmarshaled correctly
func TestUnmarshalStorageImport(t *testing.T) {
	originalJSON := `
	  {
		"storage_import": {
		  "client_content_length": 1,
		  "client_content_type": "abc",
		  "completed": "",
		  "created": "2020-06-26T08:51:07Z",
		  "direct_upload_url": "https://fi-hel1.img.upcloud.com/uploader/session/07a6c9a3-300e-4d0e-b935-624f3dbdff3f",
		  "error_code": "ghi",
		  "error_message": "jkl",
		  "md5sum": "mno",
		  "read_bytes": 2,
		  "sha256sum": "pqr",
		  "source": "direct_upload",
		  "state": "prepared",
		  "uuid": "07a6c9a3-300e-4d0e-b935-624f3dbdff3f",
		  "written_bytes": 3 
		}
	  }
	`

	storageImport := StorageImportDetails{}
	err := json.Unmarshal([]byte(originalJSON), &storageImport)
	assert.NoError(t, err)

	testStorageImport := StorageImportDetails{
		ClientContentLength: 1,
		ClientContentType:   "abc",
		Completed:           time.Time{},
		Created:             timeParse("2020-06-26T08:51:07Z"),
		DirectUploadURL:     "https://fi-hel1.img.upcloud.com/uploader/session/07a6c9a3-300e-4d0e-b935-624f3dbdff3f",
		ErrorCode:           "ghi",
		ErrorMessage:        "jkl",
		MD5Sum:              "mno",
		ReadBytes:           2,
		SHA256Sum:           "pqr",
		Source:              StorageImportSourceDirectUpload,
		State:               "prepared",
		UUID:                "07a6c9a3-300e-4d0e-b935-624f3dbdff3f",
		WrittenBytes:        3,
	}

	assert.Equal(t, testStorageImport, storageImport)
}

// TestUnmarshalResizeStorageFilesystemBackup tests that ResizeStorageFilesystemBackup struct is unmarshaled correctly
func TestUnmarshalResizeStorageFilesystemBackup(t *testing.T) {
	originalJSON := `
	{
		"resize_backup" : {
		   "access" : "private",
		   "created" : "2021-12-03T06:25:15Z",
		   "license" : 0,
		   "origin" : "017ca4cc-def2-458d-a797-7782959b30a7",
		   "servers" : {
			  "server" : [
				  "117ca4cc-def2-458d-a797-7782959b30a7"
			  ]
		   },
		   "size" : 10,
		   "state" : "online",
		   "title" : "Resize Backup",
		   "type" : "backup",
		   "uuid" : "01beec3a-14ac-4f71-9c63-3338341121c3",
		   "zone" : "uk-lon1"
		}
	}
	`

	resizeBackup := ResizeStorageFilesystemBackup{}
	err := json.Unmarshal([]byte(originalJSON), &resizeBackup)
	assert.NoError(t, err)

	testResizeBackup := ResizeStorageFilesystemBackup{
		Access:  StorageAccessPrivate,
		Created: timeParse("2021-12-03T06:25:15Z"),
		License: 0,
		Origin:  "017ca4cc-def2-458d-a797-7782959b30a7",
		Servers: ServerUUIDSlice{"117ca4cc-def2-458d-a797-7782959b30a7"},
		Size:    10,
		State:   StorageStateOnline,
		Title:   "Resize Backup",
		Type:    StorageTypeBackup,
		UUID:    "01beec3a-14ac-4f71-9c63-3338341121c3",
		Zone:    "uk-lon1",
	}

	assert.Equal(t, testResizeBackup, resizeBackup)
}
