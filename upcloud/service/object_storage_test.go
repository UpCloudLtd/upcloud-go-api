package service

import (
	"context"
	"strings"
	"testing"

	"github.com/dnaeon/go-vcr/recorder"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/UpCloudLtd/upcloud-go-api/v5/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v5/upcloud/request"
)

// TestGetObjectStorages tests that the GetObjectStorages() function returns proper data
func TestGetObjectStorages(t *testing.T) {
	t.Parallel()
	record(t, "getobjectstorages", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service) {
		objectStorages, err := svc.GetObjectStorages(ctx)
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
	t.Parallel()
	record(t, "getobjectstoragedetails", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service) {
		d, err := createObjectStorage(ctx, svc, "getobjectstoragedetails", "App object storage", "fi-hel2", 500)
		require.NoError(t, err)

		objectStorageDetails, err := svc.GetObjectStorageDetails(ctx, &request.GetObjectStorageDetailsRequest{
			UUID: d.UUID,
		})
		require.NoError(t, err)

		assert.Contains(t, objectStorageDetails.Name, "getobjectstoragedetails")
		assert.Equal(t, "fi-hel2", objectStorageDetails.Zone)
	})
}

// Creates an Object Storage and returns the details about it, panic if creation fails
func createObjectStorage(ctx context.Context, svc *Service, name string, description string, zone string, size int) (*upcloud.ObjectStorageDetails, error) {
	createObjectStorageRequest := request.CreateObjectStorageRequest{
		Name:        "go-test-" + name,
		Description: description,
		Zone:        zone,
		Size:        size,
		AccessKey:   "UCOB5HE4NVTVFMXXRBQ2",
		SecretKey:   "ssLDVHvTRjHaEAPRcMiFep3HItcqdNUNtql3DcLx",
	}

	// Create the Object Storage and block until it has started
	objectStorageDetails, err := svc.CreateObjectStorage(ctx, &createObjectStorageRequest)
	if err != nil {
		return nil, err
	}

	return objectStorageDetails, nil
}

// Deletes the specific Object Storage.
func deleteObjectStorage(ctx context.Context, svc *Service, uuid string) error {
	err := svc.DeleteObjectStorage(ctx, &request.DeleteObjectStorageRequest{
		UUID: uuid,
	})

	return err
}
