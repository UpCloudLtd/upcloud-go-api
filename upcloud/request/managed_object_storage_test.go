package request

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetManagedObjectStorageRegionsRequest_RequestURL(t *testing.T) {
	req := &GetManagedObjectStorageRegionsRequest{}
	assert.Equal(t, "/object-storage-2/regions", req.RequestURL())
}

func TestGetManagedObjectStorageRegionRequest_RequestURL(t *testing.T) {
	req := &GetManagedObjectStorageRegionRequest{
		Name: "region",
	}
	assert.Equal(t, "/object-storage-2/regions/region", req.RequestURL())
}

func TestCreateManagedObjectStorageRequest_RequestURL(t *testing.T) {
	req := &CreateManagedObjectStorageRequest{}
	assert.Equal(t, "/object-storage-2", req.RequestURL())
}

func TestGetManagedObjectStoragesRequest_RequestURL(t *testing.T) {
	req := &GetManagedObjectStoragesRequest{}
	assert.Equal(t, "/object-storage-2", req.RequestURL())
}

func TestGetManagedObjectStorageDetailsRequest_RequestURL(t *testing.T) {
	req := &GetManagedObjectStorageDetailsRequest{
		UUID: "service",
	}
	assert.Equal(t, "/object-storage-2/service", req.RequestURL())
}

func TestReplaceManagedObjectStorageRequest_RequestURL(t *testing.T) {
	req := &ReplaceManagedObjectStorageRequest{
		UUID: "service",
	}
	assert.Equal(t, "/object-storage-2/service", req.RequestURL())
}

func TestModifyManagedObjectStorageRequest_RequestURL(t *testing.T) {
	req := &ModifyManagedObjectStorageRequest{
		UUID: "service",
	}
	assert.Equal(t, "/object-storage-2/service", req.RequestURL())
}

func TestDeleteManagedObjectStorageRequest_RequestURL(t *testing.T) {
	req := &DeleteManagedObjectStorageRequest{
		UUID: "service",
	}
	assert.Equal(t, "/object-storage-2/service", req.RequestURL())
}

func TestGetManagedObjectStorageMetricsRequest_RequestURL(t *testing.T) {
	req := &GetManagedObjectStorageMetricsRequest{
		ServiceUUID: "service",
	}
	assert.Equal(t, "/object-storage-2/service/metrics", req.RequestURL())
}

func TestGetManagedObjectStorageBucketMetricsRequest_RequestURL(t *testing.T) {
	req := &GetManagedObjectStorageBucketMetricsRequest{
		ServiceUUID: "service",
		Page:        DefaultPage,
	}
	assert.Equal(t, "/object-storage-2/service/metrics/buckets?limit=100&offset=0", req.RequestURL())
}

func TestCreateManagedObjectStorageNetworkRequest_RequestURL(t *testing.T) {
	req := &CreateManagedObjectStorageNetworkRequest{
		ServiceUUID: "service",
	}
	assert.Equal(t, "/object-storage-2/service/networks", req.RequestURL())
}

func TestGetManagedObjectStorageNetworksRequest_RequestURL(t *testing.T) {
	req := &GetManagedObjectStorageNetworksRequest{
		ServiceUUID: "service",
	}
	assert.Equal(t, "/object-storage-2/service/networks", req.RequestURL())
}

func TestGetManagedObjectStorageNetworkRequest_RequestURL(t *testing.T) {
	req := &GetManagedObjectStorageNetworkRequest{
		ServiceUUID: "service",
		NetworkName: "network",
	}
	assert.Equal(t, "/object-storage-2/service/networks/network", req.RequestURL())
}

func TestDeleteManagedObjectStorageNetworkRequest_RequestURL(t *testing.T) {
	req := &DeleteManagedObjectStorageNetworkRequest{
		ServiceUUID: "service",
		NetworkName: "network",
	}
	assert.Equal(t, "/object-storage-2/service/networks/network", req.RequestURL())
}

func TestCreateManagedObjectStorageUserRequest_RequestURL(t *testing.T) {
	req := &CreateManagedObjectStorageUserRequest{
		ServiceUUID: "service",
	}
	assert.Equal(t, "/object-storage-2/service/users", req.RequestURL())
}

func TestGetManagedObjectStorageUsersRequest_RequestURL(t *testing.T) {
	req := &GetManagedObjectStorageUsersRequest{
		ServiceUUID: "service",
	}
	assert.Equal(t, "/object-storage-2/service/users", req.RequestURL())
}

func TestGetManagedObjectStorageUserRequest_RequestURL(t *testing.T) {
	req := &GetManagedObjectStorageUserRequest{
		ServiceUUID: "service",
		Username:    "user",
	}
	assert.Equal(t, "/object-storage-2/service/users/user", req.RequestURL())
}

func TestDeleteManagedObjectStorageUserRequest_RequestURL(t *testing.T) {
	req := &DeleteManagedObjectStorageUserRequest{
		ServiceUUID: "service",
		Username:    "user",
	}
	assert.Equal(t, "/object-storage-2/service/users/user", req.RequestURL())
}

func TestCreateManagedObjectStorageUserAccessKeyRequest_RequestURL(t *testing.T) {
	req := &CreateManagedObjectStorageUserAccessKeyRequest{
		ServiceUUID: "service",
		Username:    "user",
	}
	assert.Equal(t, "/object-storage-2/service/users/user/access-keys", req.RequestURL())
}

func TestGetManagedObjectStorageUserAccessKeysRequest_RequestURL(t *testing.T) {
	req := &GetManagedObjectStorageUserAccessKeysRequest{
		ServiceUUID: "service",
		Username:    "user",
	}
	assert.Equal(t, "/object-storage-2/service/users/user/access-keys", req.RequestURL())
}

func TestGetManagedObjectStorageUserAccessKeyRequest_RequestURL(t *testing.T) {
	req := &GetManagedObjectStorageUserAccessKeyRequest{
		ServiceUUID: "service",
		Username:    "user",
		Name:        "access",
	}
	assert.Equal(t, "/object-storage-2/service/users/user/access-keys/access", req.RequestURL())
}

func TestModifyManagedObjectStorageUserAccessKeyRequest_RequestURL(t *testing.T) {
	req := &ModifyManagedObjectStorageUserAccessKeyRequest{
		ServiceUUID: "service",
		Username:    "user",
		Name:        "access",
	}
	assert.Equal(t, "/object-storage-2/service/users/user/access-keys/access", req.RequestURL())
}

func TestDeleteManagedObjectStorageUserAccessKeyRequest_RequestURL(t *testing.T) {
	req := &DeleteManagedObjectStorageUserAccessKeyRequest{
		ServiceUUID: "service",
		Username:    "user",
		Name:        "access",
	}
	assert.Equal(t, "/object-storage-2/service/users/user/access-keys/access", req.RequestURL())
}
