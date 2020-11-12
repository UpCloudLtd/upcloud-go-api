package request

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestGetObjectStorageDetailsRequest tests that GetObjectStorageDetailsRequest objects behave correctly
func TestGetObjectStorageDetailsRequest(t *testing.T) {
	request := GetObjectStorageDetailsRequest{
		UUID: "foo",
	}

	assert.Equal(t, "/object-storage/foo", request.RequestURL())
}

// TestCreateObjectStorageRequest tests that CreateObjectStorageRequest objects behave correctly
func TestCreateObjectStorageRequest(t *testing.T) {
	request := CreateObjectStorageRequest{
		Name:        "app-object-storage",
		Description: "App object storage",
		Zone:        "fi-hel2",
		Size:        500,
		AccessKey:   "UCOB5HE4NVTVFMXXRBQ2",
		SecretKey:   "ssLDVHvTRjHaEAPRcMiFep3HItcqdNUNtql3DcLx",
	}

	expectedJSON := `
		{
			"object_storage": {
				"access_key": "UCOB5HE4NVTVFMXXRBQ2",
				"description": "App object storage",
				"name": "app-object-storage",
				"secret_key": "ssLDVHvTRjHaEAPRcMiFep3HItcqdNUNtql3DcLx",
				"zone": "fi-hel2",
				"size": 500
			}
		}
	`

	actualJSON, err := json.Marshal(&request)
	assert.NoError(t, err)
	assert.JSONEq(t, expectedJSON, string(actualJSON))
	assert.Equal(t, "/object-storage", request.RequestURL())
}

// TestModifyObjectStorageRequest tests that ModifyObjectStorageRequest objects behave correctly
func TestModifyObjectStorageRequest(t *testing.T) {
	request := ModifyObjectStorageRequest{
		UUID:        "foo",
		Description: "Modified object storage",
		AccessKey:   "UCOB5HE4NVTVFMXXRBQ2",
		SecretKey:   "ssLDVHvTRjHaEAPRcMiFep3HItcqdNUNtql3DcLx",
	}

	expectedJSON := `
		{
			"object_storage": {
				"access_key": "UCOB5HE4NVTVFMXXRBQ2",
				"description": "Modified object storage",
				"secret_key": "ssLDVHvTRjHaEAPRcMiFep3HItcqdNUNtql3DcLx"
			}
		}
	`

	actualJSON, err := json.Marshal(&request)
	assert.Nil(t, err)
	assert.JSONEq(t, expectedJSON, string(actualJSON))
	assert.Equal(t, "/object-storage/foo", request.RequestURL())
}

// TestDeleteObjectStorageRequest tests that DeleteObjectStorageRequest objects behave correctly
func TestDeleteObjectStorageRequest(t *testing.T) {
	request := DeleteObjectStorageRequest{
		UUID: "foo",
	}

	assert.Equal(t, "/object-storage/foo", request.RequestURL())
}
