package service

import (
	"strings"
	"testing"

	"github.com/dnaeon/go-vcr/recorder"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud/request"
)

// TestGetObjectStorages tests that the GetObjectStorages() function returns proper data
func TestGetObjectStorages(t *testing.T) {
	record(t, "getobjectstorages", func(t *testing.T, rec *recorder.Recorder, svc *Service) {
		objectStorages, err := svc.GetObjectStorages()
		require.NoError(t, err)
		assert.NotEmpty(t, objectStorages.ObjectStorages)

		var found bool
		for _, o := range objectStorages.ObjectStorages {
			if strings.Contains(o.Name, "getobjectstorages") {
				found = true
				break
			}
		}
		assert.True(t, found)
	})
}

// TestGetObjectStorageDetails ensures that the GetObjectStorageDetails() function returns proper data
func TestGetObjectStorageDetails(t *testing.T) {
	record(t, "getobjectstoragedetails", func(t *testing.T, rec *recorder.Recorder, svc *Service) {
		d, err := createObjectStorage(svc, "getobjectstoragedetails", "App object storage", "fi-hel2", 500)
		require.NoError(t, err)

		objectStorageDetails, err := svc.GetObjectStorageDetails(&request.GetObjectStorageDetailsRequest{
			UUID: d.UUID,
		})
		require.NoError(t, err)

		assert.Contains(t, objectStorageDetails.Name, "getobjectstoragedetails")
		assert.Equal(t, "fi-hel2", objectStorageDetails.Zone)
	})
}
