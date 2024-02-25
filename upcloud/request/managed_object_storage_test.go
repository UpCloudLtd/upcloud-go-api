package request

import (
	"encoding/json"
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

func TestCreateManagedObjectStorageRequest_MarshalJSON(t *testing.T) {
	t.Run("TestMinimal", func(t *testing.T) {
		req := CreateManagedObjectStorageRequest{
			Region: "europe-1",
		}
		d, err := json.Marshal(&req)
		assert.NoError(t, err)

		const expected = `{
			"configured_status":"",
			"networks":null,
			"region":"europe-1"
		}`
		assert.JSONEq(t, expected, string(d))
	})

	t.Run("TestWithName", func(t *testing.T) {
		req := CreateManagedObjectStorageRequest{
			Name:   "test-objsto-name",
			Region: "europe-1",
		}
		d, err := json.Marshal(&req)
		assert.NoError(t, err)

		const expected = `{
			"configured_status":"",
			"name":"test-objsto-name",
			"networks":null,
			"region":"europe-1"
		}`
		assert.JSONEq(t, expected, string(d))
	})
}

func TestGetManagedObjectStoragesRequest_RequestURL(t *testing.T) {
	req := &GetManagedObjectStoragesRequest{}
	assert.Equal(t, "/object-storage-2", req.RequestURL())
}

func TestGetManagedObjectStorageDetailsRequest_RequestURL(t *testing.T) {
	req := &GetManagedObjectStorageRequest{
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
	assert.Equal(t, "/object-storage-2/service/buckets?limit=100&offset=0", req.RequestURL())
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
		AccessKeyId: "access",
	}
	assert.Equal(t, "/object-storage-2/service/users/user/access-keys/access", req.RequestURL())
}

func TestModifyManagedObjectStorageUserAccessKeyRequest_RequestURL(t *testing.T) {
	req := &ModifyManagedObjectStorageUserAccessKeyRequest{
		ServiceUUID: "service",
		Username:    "user",
		AccessKeyId: "access",
	}
	assert.Equal(t, "/object-storage-2/service/users/user/access-keys/access", req.RequestURL())
}

func TestDeleteManagedObjectStorageUserAccessKeyRequest_RequestURL(t *testing.T) {
	req := &DeleteManagedObjectStorageUserAccessKeyRequest{
		ServiceUUID: "service",
		Username:    "user",
		AccessKeyId: "access",
	}
	assert.Equal(t, "/object-storage-2/service/users/user/access-keys/access", req.RequestURL())
}

func TestCreateManagedObjectStoragePolicyRequest_RequestURL(t *testing.T) {
	req := &CreateManagedObjectStoragePolicyRequest{
		ServiceUUID: "service",
	}
	assert.Equal(t, "/object-storage-2/service/policies", req.RequestURL())
}

func TestGetManagedObjectStoragePoliciesRequest_RequestURL(t *testing.T) {
	req := &GetManagedObjectStoragePoliciesRequest{
		ServiceUUID: "service",
	}
	assert.Equal(t, "/object-storage-2/service/policies", req.RequestURL())
}

func TestGetManagedObjectStoragePolicyRequest_RequestURL(t *testing.T) {
	req := &GetManagedObjectStoragePolicyRequest{
		ServiceUUID: "service",
		Name:        "policy",
	}
	assert.Equal(t, "/object-storage-2/service/policies/policy", req.RequestURL())
}

func TestDeleteManagedObjectStoragePolicyRequest_RequestURL(t *testing.T) {
	req := &DeleteManagedObjectStoragePolicyRequest{
		ServiceUUID: "service",
		Name:        "policy",
	}
	assert.Equal(t, "/object-storage-2/service/policies/policy", req.RequestURL())
}

func TestAttachManagedObjectStorageUserPolicyRequest_RequestURL(t *testing.T) {
	req := &AttachManagedObjectStorageUserPolicyRequest{
		ServiceUUID: "service",
		Username:    "user",
		Name:        "policy",
	}
	assert.Equal(t, "/object-storage-2/service/users/user/policies", req.RequestURL())
}

func TestGetManagedObjectStorageUserPoliciesRequest_RequestURL(t *testing.T) {
	req := &GetManagedObjectStorageUserPoliciesRequest{
		ServiceUUID: "service",
		Username:    "user",
	}
	assert.Equal(t, "/object-storage-2/service/users/user/policies", req.RequestURL())
}

func TestDetachManagedObjectStorageUserPolicyRequest_RequestURL(t *testing.T) {
	req := &DetachManagedObjectStorageUserPolicyRequest{
		ServiceUUID: "service",
		Username:    "user",
		Name:        "policy",
	}
	assert.Equal(t, "/object-storage-2/service/users/user/policies/policy", req.RequestURL())
}
