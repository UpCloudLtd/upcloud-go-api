package upcloud

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestUnmarshalObjectStorages tests that Object Storages are unmarshaled correctly.
func TestUnmarshalObjectStorages(t *testing.T) {
	originalJSON := `
		{
			"object_storages": {
				"object_storage": [
					{
						"created": "2020-07-23T05:06:35Z",
						"description": "Example object storage",
						"name": "example-object-storage",
						"size": 250,
						"state": "started",
						"url": "https://example-object-storage.nl-ams1.upcloudobjects.com/",
						"uuid": "06832a75-be7b-4d23-be05-130dc3dfd9e7",
						"zone": "uk-lon1"
					}
				]
			}
		}
	`

	objectStorages := ObjectStorages{}
	err := json.Unmarshal([]byte(originalJSON), &objectStorages)
	assert.Nil(t, err)
	assert.Len(t, objectStorages.ObjectStorages, 1)

	objectStorage := objectStorages.ObjectStorages[0]
	assert.Equal(t, "2020-07-23T05:06:35Z", objectStorage.Created)
	assert.Equal(t, "Example object storage", objectStorage.Description)
	assert.Equal(t, "example-object-storage", objectStorage.Name)
	assert.Equal(t, 250, objectStorage.Size)
	assert.Equal(t, "started", objectStorage.State)
	assert.Equal(t, "https://example-object-storage.nl-ams1.upcloudobjects.com/", objectStorage.URL)
	assert.Equal(t, "06832a75-be7b-4d23-be05-130dc3dfd9e7", objectStorage.UUID)
	assert.Equal(t, "uk-lon1", objectStorage.Zone)
}
