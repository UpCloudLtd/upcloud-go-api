package service

import (
	"context"
	"testing"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
	"github.com/dnaeon/go-vcr/recorder"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetManagedObjectStorageRegions(t *testing.T) {
	t.Parallel()

	record(t, "getmanagedobjectstorageregions", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service) {
		regions, err := svc.GetManagedObjectStorageRegions(ctx, &request.GetManagedObjectStorageRegionsRequest{})
		require.NoError(t, err)
		require.NotEmpty(t, regions)
	})
}

func TestGetManagedObjectStorageRegion(t *testing.T) {
	t.Parallel()

	record(t, "getmanagedobjectstorageregion", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service) {
		regions, err := svc.GetManagedObjectStorageRegions(ctx, &request.GetManagedObjectStorageRegionsRequest{})
		require.NoError(t, err)
		require.NotEmpty(t, regions)

		region, err := svc.GetManagedObjectStorageRegion(ctx, &request.GetManagedObjectStorageRegionRequest{Name: regions[0].Name})
		require.NoError(t, err)
		require.Equal(t, regions[0], *region)
	})
}

func TestCreateManagedObjectStorage(t *testing.T) {
	record(t, "createmanagedobjectstorage", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service) {
		storage, err := createManagedObjectStorage(ctx, svc)
		require.NoError(t, err)

		defer func(uuid string) {
			err = deleteManagedObjectStorageAndUsers(ctx, svc, uuid)
			require.NoError(t, err)
		}(storage.UUID)

		require.Equal(t, storage.ConfiguredStatus, upcloud.ManagedObjectStorageConfiguredStatusStarted)
		require.Len(t, storage.Labels, 1)
		require.Len(t, storage.Networks, 1)
		require.NotEmpty(t, storage.UUID)
	})
}

func TestGetManagedObjectStorages(t *testing.T) {
	record(t, "getmanagedobjectstorages", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service) {
		storages, err := svc.GetManagedObjectStorages(ctx, &request.GetManagedObjectStoragesRequest{})
		require.NoError(t, err)
		require.Len(t, storages, 0)

		storage, err := createManagedObjectStorage(ctx, svc)
		require.NoError(t, err)

		defer func(uuid string) {
			err = deleteManagedObjectStorageAndUsers(ctx, svc, uuid)
			require.NoError(t, err)
		}(storage.UUID)

		storages, err = svc.GetManagedObjectStorages(ctx, &request.GetManagedObjectStoragesRequest{})
		require.NoError(t, err)
		require.Len(t, storages, 1)
	})
}

func TestGetManagedObjectStorageDetails(t *testing.T) {
	record(t, "getmanagedobjectstoragedetails", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service) {
		storage, err := createManagedObjectStorage(ctx, svc)
		require.NoError(t, err)

		defer func(uuid string) {
			err = deleteManagedObjectStorageAndUsers(ctx, svc, uuid)
			require.NoError(t, err)
		}(storage.UUID)

		storage, err = svc.GetManagedObjectStorage(ctx, &request.GetManagedObjectStorageRequest{UUID: storage.UUID})
		require.NoError(t, err)
		require.Equal(t, storage.ConfiguredStatus, upcloud.ManagedObjectStorageConfiguredStatusStarted)
		require.Len(t, storage.Labels, 1)
		require.Len(t, storage.Networks, 1)
		require.NotEmpty(t, storage.UUID)
		require.NotEmpty(t, storage.Name)
		require.NotEmpty(t, storage.CreatedAt)
		require.NotEmpty(t, storage.Endpoints)
		require.NotEmpty(t, storage.Region)
	})
}

func TestReplaceManagedObjectStorage(t *testing.T) {
	record(t, "replacemanagedobjectstorage", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service) {
		storage, err := createManagedObjectStorage(ctx, svc)
		require.NoError(t, err)

		defer func(uuid string) {
			err = deleteManagedObjectStorageAndUsers(ctx, svc, uuid)
			require.NoError(t, err)
		}(storage.UUID)

		storage, err = svc.ReplaceManagedObjectStorage(ctx, &request.ReplaceManagedObjectStorageRequest{
			ConfiguredStatus: upcloud.ManagedObjectStorageConfiguredStatusStopped,
			Networks: []upcloud.ManagedObjectStorageNetwork{{
				Family: "IPv4",
				Name:   "replaced-network",
				Type:   "public",
			}},
			UUID: storage.UUID,
			Name: "test2",
		})
		require.NoError(t, err)
		require.Equal(t, upcloud.ManagedObjectStorageConfiguredStatusStopped, storage.ConfiguredStatus)
		require.Equal(t, storage.Networks[0].Name, "replaced-network")
		require.Len(t, storage.Labels, 0)
	})
}

func TestModifyManagedObjectStorage(t *testing.T) {
	record(t, "modifymanagedobjectstorage", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service) {
		storage, err := createManagedObjectStorage(ctx, svc)
		require.NoError(t, err)

		defer func(uuid string) {
			err = deleteManagedObjectStorageAndUsers(ctx, svc, uuid)
			require.NoError(t, err)
		}(storage.UUID)

		status := upcloud.ManagedObjectStorageConfiguredStatusStopped
		storage, err = svc.ModifyManagedObjectStorage(ctx, &request.ModifyManagedObjectStorageRequest{
			ConfiguredStatus: &status,
			Labels:           &[]upcloud.Label{},
			UUID:             storage.UUID,
		})
		require.NoError(t, err)
		require.Equal(t, upcloud.ManagedObjectStorageConfiguredStatusStopped, storage.ConfiguredStatus)
		require.Len(t, storage.Labels, 0)
	})
}

func TestGetManagedObjectStorageMetrics(t *testing.T) {
	record(t, "getmanagedobjectstoragemetrics", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service) {
		storage, err := createManagedObjectStorage(ctx, svc)
		require.NoError(t, err)

		defer func(uuid string) {
			err = deleteManagedObjectStorageAndUsers(ctx, svc, uuid)
			require.NoError(t, err)
		}(storage.UUID)

		m, err := svc.GetManagedObjectStorageMetrics(ctx, &request.GetManagedObjectStorageMetricsRequest{ServiceUUID: storage.UUID})
		require.NoError(t, err)
		require.Equal(t, m.TotalObjects, 0)
		require.Equal(t, m.TotalSizeBytes, 0)
	})
}

func TestManagedObjectStorageBucketOperations(t *testing.T) {
	record(t, "managedobjectstoragebucketoperations", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service) {
		objsto, err := createManagedObjectStorage(ctx, svc)
		require.NoError(t, err)

		defer func(uuid string) {
			err = deleteManagedObjectStorageAndUsers(ctx, svc, uuid)
			require.NoError(t, err)
		}(objsto.UUID)

		m, err := svc.GetManagedObjectStorageBucketMetrics(ctx, &request.GetManagedObjectStorageBucketMetricsRequest{ServiceUUID: objsto.UUID})
		require.NoError(t, err)
		assert.Len(t, m, 0)

		_, err = svc.CreateManagedObjectStorageBucket(ctx, &request.CreateManagedObjectStorageBucketRequest{Name: "test", ServiceUUID: objsto.UUID})
		require.NoError(t, err)

		m, err = svc.GetManagedObjectStorageBucketMetrics(ctx, &request.GetManagedObjectStorageBucketMetricsRequest{ServiceUUID: objsto.UUID})
		require.NoError(t, err)
		assert.Len(t, m, 1)
		assert.Equal(t, "test", m[0].Name)
		assert.False(t, m[0].Deleted)

		err = svc.DeleteManagedObjectStorageBucket(ctx, &request.DeleteManagedObjectStorageBucketRequest{Name: "test", ServiceUUID: objsto.UUID})
		require.NoError(t, err)

		m, err = svc.GetManagedObjectStorageBucketMetrics(ctx, &request.GetManagedObjectStorageBucketMetricsRequest{ServiceUUID: objsto.UUID})
		require.NoError(t, err)
		assert.Len(t, m, 1)
		assert.Equal(t, "test", m[0].Name)
		assert.True(t, m[0].Deleted)

		err = svc.WaitForManagedObjectStorageBucketDeletion(ctx, &request.WaitForManagedObjectStorageBucketDeletionRequest{Name: "test", ServiceUUID: objsto.UUID})
		require.NoError(t, err)
	})
}

func TestCreateManagedObjectStorageNetwork(t *testing.T) {
	record(t, "createmanagedobjectstoragenetwork", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service) {
		storage, err := createManagedObjectStorage(ctx, svc)
		require.NoError(t, err)

		defer func(uuid string) {
			err = deleteManagedObjectStorageAndUsers(ctx, svc, uuid)
			require.NoError(t, err)
		}(storage.UUID)

		router, err := svc.CreateRouter(ctx, &request.CreateRouterRequest{
			Name: "managed-object-storage-router",
		})
		require.NoError(t, err)

		network, err := svc.CreateNetwork(ctx, &request.CreateNetworkRequest{
			Name:   "managed-object-storage",
			Zone:   "fi-hel1",
			Router: router.UUID,
			IPNetworks: upcloud.IPNetworkSlice{upcloud.IPNetwork{
				Address:          "172.18.1.0/24",
				DHCP:             0,
				DHCPDefaultRoute: 0,
				DHCPDns:          nil,
				DHCPRoutes:       nil,
				Family:           "IPv4",
				Gateway:          "172.18.1.1",
			}},
		})
		require.NoError(t, err)

		storageNetwork, err := svc.CreateManagedObjectStorageNetwork(ctx, &request.CreateManagedObjectStorageNetworkRequest{
			Family:      "IPv4",
			Name:        "private-network",
			ServiceUUID: storage.UUID,
			Type:        "private",
			UUID:        network.UUID,
		})
		require.NoError(t, err)
		require.Equal(t, network.UUID, *storageNetwork.UUID)

		err = svc.DeleteNetwork(ctx, &request.DeleteNetworkRequest{UUID: network.UUID})
		require.NoError(t, err)

		err = svc.DeleteRouter(ctx, &request.DeleteRouterRequest{UUID: router.UUID})
		require.NoError(t, err)
	})
}

func TestGetManagedObjectStorageNetworks(t *testing.T) {
	record(t, "getmanagedobjectstoragenetworks", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service) {
		storage, err := createManagedObjectStorage(ctx, svc)
		require.NoError(t, err)

		defer func(uuid string) {
			err = deleteManagedObjectStorageAndUsers(ctx, svc, uuid)
			require.NoError(t, err)
		}(storage.UUID)

		networks, err := svc.GetManagedObjectStorageNetworks(ctx, &request.GetManagedObjectStorageNetworksRequest{ServiceUUID: storage.UUID})
		require.NoError(t, err)
		require.Len(t, networks, 1)
	})
}

func TestGetManagedObjectStorageNetwork(t *testing.T) {
	record(t, "getmanagedobjectstoragenetwork", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service) {
		storage, err := createManagedObjectStorage(ctx, svc)
		require.NoError(t, err)

		defer func(uuid string) {
			err = deleteManagedObjectStorageAndUsers(ctx, svc, uuid)
			require.NoError(t, err)
		}(storage.UUID)

		networks, err := svc.GetManagedObjectStorageNetworks(ctx, &request.GetManagedObjectStorageNetworksRequest{ServiceUUID: storage.UUID})
		require.NoError(t, err)
		require.Len(t, networks, 1)

		network, err := svc.GetManagedObjectStorageNetwork(ctx, &request.GetManagedObjectStorageNetworkRequest{ServiceUUID: storage.UUID, NetworkName: networks[0].Name})
		require.NoError(t, err)
		require.Equal(t, networks[0].Name, network.Name)
	})
}

func TestDeleteManagedObjectStorageNetwork(t *testing.T) {
	record(t, "deletemanagedobjectstoragenetwork", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service) {
		storage, err := createManagedObjectStorage(ctx, svc)
		require.NoError(t, err)

		defer func(uuid string) {
			err = deleteManagedObjectStorageAndUsers(ctx, svc, uuid)
			require.NoError(t, err)
		}(storage.UUID)

		router, err := svc.CreateRouter(ctx, &request.CreateRouterRequest{
			Name: "managed-object-storage-router",
		})
		require.NoError(t, err)

		network, err := svc.CreateNetwork(ctx, &request.CreateNetworkRequest{
			Name:   "managed-object-storage",
			Zone:   "fi-hel1",
			Router: router.UUID,
			IPNetworks: upcloud.IPNetworkSlice{upcloud.IPNetwork{
				Address:          "172.18.1.0/24",
				DHCP:             0,
				DHCPDefaultRoute: 0,
				DHCPDns:          nil,
				DHCPRoutes:       nil,
				Family:           "IPv4",
				Gateway:          "172.18.1.1",
			}},
		})
		require.NoError(t, err)

		storageNetwork, err := svc.CreateManagedObjectStorageNetwork(ctx, &request.CreateManagedObjectStorageNetworkRequest{
			Family:      "IPv4",
			Name:        "private-network",
			ServiceUUID: storage.UUID,
			Type:        "private",
			UUID:        network.UUID,
		})
		require.NoError(t, err)
		require.Equal(t, network.UUID, *storageNetwork.UUID)

		err = svc.DeleteManagedObjectStorageNetwork(ctx, &request.DeleteManagedObjectStorageNetworkRequest{ServiceUUID: storage.UUID, NetworkName: "private-network"})
		require.NoError(t, err)

		err = svc.DeleteNetwork(ctx, &request.DeleteNetworkRequest{UUID: network.UUID})
		require.NoError(t, err)

		err = svc.DeleteRouter(ctx, &request.DeleteRouterRequest{UUID: router.UUID})
		require.NoError(t, err)
	})
}

func TestCreateManagedObjectStorageUser(t *testing.T) {
	record(t, "createmanagedobjectstorageuser", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service) {
		storage, err := createManagedObjectStorage(ctx, svc)
		require.NoError(t, err)

		defer func(uuid string) {
			err = deleteManagedObjectStorageAndUsers(ctx, svc, uuid)
			require.NoError(t, err)
		}(storage.UUID)

		user, err := svc.CreateManagedObjectStorageUser(ctx, &request.CreateManagedObjectStorageUserRequest{
			Username:    "testuser",
			ServiceUUID: storage.UUID,
		})
		require.NoError(t, err)
		require.Equal(t, user.Username, "testuser")
	})
}

func TestGetManagedObjectStorageUsers(t *testing.T) {
	record(t, "getmanagedobjectstorageusers", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service) {
		storage, err := createManagedObjectStorage(ctx, svc)
		require.NoError(t, err)

		defer func(uuid string) {
			err = deleteManagedObjectStorageAndUsers(ctx, svc, uuid)
			require.NoError(t, err)
		}(storage.UUID)

		user, err := svc.CreateManagedObjectStorageUser(ctx, &request.CreateManagedObjectStorageUserRequest{
			Username:    "testuser",
			ServiceUUID: storage.UUID,
		})
		require.NoError(t, err)
		require.Equal(t, user.Username, "testuser")
		require.NotEmpty(t, user.ARN)

		users, err := svc.GetManagedObjectStorageUsers(ctx, &request.GetManagedObjectStorageUsersRequest{ServiceUUID: storage.UUID})
		require.NoError(t, err)
		require.Len(t, users, 1)
	})
}

func TestGetManagedObjectStorageUser(t *testing.T) {
	record(t, "getmanagedobjectstorageuser", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service) {
		storage, err := createManagedObjectStorage(ctx, svc)
		require.NoError(t, err)

		defer func(uuid string) {
			err = deleteManagedObjectStorageAndUsers(ctx, svc, uuid)
			require.NoError(t, err)
		}(storage.UUID)

		user, err := svc.CreateManagedObjectStorageUser(ctx, &request.CreateManagedObjectStorageUserRequest{
			Username:    "testuser",
			ServiceUUID: storage.UUID,
		})
		require.NoError(t, err)
		require.Equal(t, user.Username, "testuser")
		require.NotEmpty(t, user.ARN)

		_, err = svc.GetManagedObjectStorageUser(ctx, &request.GetManagedObjectStorageUserRequest{ServiceUUID: storage.UUID, Username: user.Username})
		require.NoError(t, err)
	})
}

func TestDeleteManagedObjectStorageUser(t *testing.T) {
	record(t, "deletemanagedobjectstorageuser", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service) {
		storage, err := createManagedObjectStorage(ctx, svc)
		require.NoError(t, err)

		defer func(uuid string) {
			err = deleteManagedObjectStorageAndUsers(ctx, svc, uuid)
			require.NoError(t, err)
		}(storage.UUID)

		user, err := svc.CreateManagedObjectStorageUser(ctx, &request.CreateManagedObjectStorageUserRequest{
			Username:    "testuser",
			ServiceUUID: storage.UUID,
		})
		require.NoError(t, err)

		err = svc.DeleteManagedObjectStorageUser(ctx, &request.DeleteManagedObjectStorageUserRequest{
			ServiceUUID: storage.UUID,
			Username:    user.Username,
		})
		require.NoError(t, err)
	})
}

func TestCreateManagedObjectStorageUserAccessKey(t *testing.T) {
	record(t, "createmanagedobjectstorageuseraccesskey", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service) {
		storage, err := createManagedObjectStorage(ctx, svc)
		require.NoError(t, err)

		defer func(uuid string) {
			err = deleteManagedObjectStorageAndUsers(ctx, svc, uuid)
			require.NoError(t, err)
		}(storage.UUID)

		user, err := svc.CreateManagedObjectStorageUser(ctx, &request.CreateManagedObjectStorageUserRequest{
			Username:    "testuser",
			ServiceUUID: storage.UUID,
		})
		require.NoError(t, err)
		require.Equal(t, user.Username, "testuser")

		accessKey, err := svc.CreateManagedObjectStorageUserAccessKey(ctx, &request.CreateManagedObjectStorageUserAccessKeyRequest{
			Username:    user.Username,
			ServiceUUID: storage.UUID,
		})
		require.NoError(t, err)
		require.Equal(t, upcloud.ManagedObjectStorageUserAccessKeyStatusActive, accessKey.Status)
		require.NotEmpty(t, accessKey.AccessKeyID)
		require.NotEmpty(t, accessKey.CreatedAt)
		require.NotEmpty(t, accessKey.SecretAccessKey)
	})
}

func TestGetManagedObjectStorageUserAccessKeys(t *testing.T) {
	record(t, "getmanagedobjectstorageuseraccesskeys", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service) {
		storage, err := createManagedObjectStorage(ctx, svc)
		require.NoError(t, err)

		defer func(uuid string) {
			err = deleteManagedObjectStorageAndUsers(ctx, svc, uuid)
			require.NoError(t, err)
		}(storage.UUID)

		user, err := svc.CreateManagedObjectStorageUser(ctx, &request.CreateManagedObjectStorageUserRequest{
			Username:    "testuser",
			ServiceUUID: storage.UUID,
		})
		require.NoError(t, err)
		require.Equal(t, user.Username, "testuser")

		_, err = svc.CreateManagedObjectStorageUserAccessKey(ctx, &request.CreateManagedObjectStorageUserAccessKeyRequest{
			Username:    user.Username,
			ServiceUUID: storage.UUID,
		})
		require.NoError(t, err)

		_, err = svc.CreateManagedObjectStorageUserAccessKey(ctx, &request.CreateManagedObjectStorageUserAccessKeyRequest{
			Username:    user.Username,
			ServiceUUID: storage.UUID,
		})
		require.NoError(t, err)

		accessKeys, err := svc.GetManagedObjectStorageUserAccessKeys(ctx, &request.GetManagedObjectStorageUserAccessKeysRequest{
			ServiceUUID: storage.UUID,
			Username:    user.Username,
		})
		require.NoError(t, err)
		require.Len(t, accessKeys, 2)
	})
}

func TestGetManagedObjectStorageUserAccessKey(t *testing.T) {
	record(t, "getmanagedobjectstorageuseraccesskey", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service) {
		storage, err := createManagedObjectStorage(ctx, svc)
		require.NoError(t, err)

		defer func(uuid string) {
			err = deleteManagedObjectStorageAndUsers(ctx, svc, uuid)
			require.NoError(t, err)
		}(storage.UUID)

		user, err := svc.CreateManagedObjectStorageUser(ctx, &request.CreateManagedObjectStorageUserRequest{
			Username:    "testuser",
			ServiceUUID: storage.UUID,
		})
		require.NoError(t, err)
		require.Equal(t, user.Username, "testuser")

		accessKey, err := svc.CreateManagedObjectStorageUserAccessKey(ctx, &request.CreateManagedObjectStorageUserAccessKeyRequest{
			Username:    user.Username,
			ServiceUUID: storage.UUID,
		})
		require.NoError(t, err)

		accessKeys, err := svc.GetManagedObjectStorageUserAccessKeys(ctx, &request.GetManagedObjectStorageUserAccessKeysRequest{
			ServiceUUID: storage.UUID,
			Username:    user.Username,
		})
		require.NoError(t, err)
		require.Len(t, accessKeys, 1)

		accessKey, err = svc.GetManagedObjectStorageUserAccessKey(ctx, &request.GetManagedObjectStorageUserAccessKeyRequest{
			ServiceUUID: storage.UUID,
			Username:    user.Username,
			AccessKeyID: accessKey.AccessKeyID,
		})
		require.NoError(t, err)
		require.Empty(t, accessKey.SecretAccessKey)
	})
}

func TestModifyManagedObjectStorageUserAccessKey(t *testing.T) {
	record(t, "modifymanagedobjectstorageuseraccesskey", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service) {
		storage, err := createManagedObjectStorage(ctx, svc)
		require.NoError(t, err)

		defer func(uuid string) {
			err = deleteManagedObjectStorageAndUsers(ctx, svc, uuid)
			require.NoError(t, err)
		}(storage.UUID)

		user, err := svc.CreateManagedObjectStorageUser(ctx, &request.CreateManagedObjectStorageUserRequest{
			Username:    "testuser",
			ServiceUUID: storage.UUID,
		})
		require.NoError(t, err)
		require.Equal(t, user.Username, "testuser")

		accessKey, err := svc.CreateManagedObjectStorageUserAccessKey(ctx, &request.CreateManagedObjectStorageUserAccessKeyRequest{
			Username:    user.Username,
			ServiceUUID: storage.UUID,
		})
		require.NoError(t, err)
		require.Equal(t, upcloud.ManagedObjectStorageUserAccessKeyStatusActive, accessKey.Status)

		accessKey, err = svc.ModifyManagedObjectStorageUserAccessKey(ctx, &request.ModifyManagedObjectStorageUserAccessKeyRequest{
			ServiceUUID: storage.UUID,
			Username:    user.Username,
			AccessKeyID: accessKey.AccessKeyID,
			Status:      upcloud.ManagedObjectStorageUserAccessKeyStatusInactive,
		})
		require.NoError(t, err)
		require.Equal(t, upcloud.ManagedObjectStorageUserAccessKeyStatusInactive, accessKey.Status)

		accessKey, err = svc.ModifyManagedObjectStorageUserAccessKey(ctx, &request.ModifyManagedObjectStorageUserAccessKeyRequest{
			ServiceUUID: storage.UUID,
			Username:    user.Username,
			AccessKeyID: accessKey.AccessKeyID,
			Status:      upcloud.ManagedObjectStorageUserAccessKeyStatusActive,
		})
		require.NoError(t, err)
		require.Equal(t, upcloud.ManagedObjectStorageUserAccessKeyStatusActive, accessKey.Status)
	})
}

func TestDeleteManagedObjectStorageUserAccessKey(t *testing.T) {
	record(t, "deletemanagedobjectstorageuseraccesskey", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service) {
		storage, err := createManagedObjectStorage(ctx, svc)
		require.NoError(t, err)

		defer func(uuid string) {
			err = deleteManagedObjectStorageAndUsers(ctx, svc, uuid)
			require.NoError(t, err)
		}(storage.UUID)

		user, err := svc.CreateManagedObjectStorageUser(ctx, &request.CreateManagedObjectStorageUserRequest{
			Username:    "testuser",
			ServiceUUID: storage.UUID,
		})
		require.NoError(t, err)
		require.Equal(t, user.Username, "testuser")

		accessKey, err := svc.CreateManagedObjectStorageUserAccessKey(ctx, &request.CreateManagedObjectStorageUserAccessKeyRequest{
			Username:    user.Username,
			ServiceUUID: storage.UUID,
		})
		require.NoError(t, err)

		err = svc.DeleteManagedObjectStorageUserAccessKey(ctx, &request.DeleteManagedObjectStorageUserAccessKeyRequest{
			ServiceUUID: storage.UUID,
			Username:    user.Username,
			AccessKeyID: accessKey.AccessKeyID,
		})
		require.NoError(t, err)
	})
}

func TestCreateManagedObjectStoragePolicy(t *testing.T) {
	record(t, "createmanagedobjectstoragepolicy", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service) {
		storage, err := createManagedObjectStorage(ctx, svc)
		require.NoError(t, err)

		defer func(uuid string) {
			err = deleteManagedObjectStorageAndUsers(ctx, svc, uuid)
			require.NoError(t, err)
		}(storage.UUID)

		policy, err := svc.CreateManagedObjectStoragePolicy(ctx, &request.CreateManagedObjectStoragePolicyRequest{
			Name:        "testpolicy",
			Description: "description2",
			Document:    "%7B%22Version%22%3A%20%222012-10-17%22%2C%20%20%22Statement%22%3A%20%5B%7B%22Action%22%3A%20%5B%22iam%3AGetUser%22%5D%2C%20%22Resource%22%3A%20%22%2A%22%2C%20%22Effect%22%3A%20%22Allow%22%2C%20%22Sid%22%3A%20%22editor%22%7D%5D%7D",
			ServiceUUID: storage.UUID,
		})
		require.NoError(t, err)
		assert.Equal(t, policy.Name, "testpolicy")
		assert.Equal(t, policy.Description, "description2")

		err = svc.DeleteManagedObjectStoragePolicy(ctx, &request.DeleteManagedObjectStoragePolicyRequest{
			ServiceUUID: storage.UUID,
			Name:        policy.Name,
		})
		require.NoError(t, err)
	})
}

func TestGetManagedObjectStoragePolicies(t *testing.T) {
	record(t, "getmanagedobjectstoragepolicies", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service) {
		storage, err := createManagedObjectStorage(ctx, svc)
		require.NoError(t, err)

		defer func(uuid string) {
			err = deleteManagedObjectStorageAndUsers(ctx, svc, uuid)
			require.NoError(t, err)
		}(storage.UUID)

		policy, err := svc.CreateManagedObjectStoragePolicy(ctx, &request.CreateManagedObjectStoragePolicyRequest{
			Name:        "testpolicy",
			Description: "description2",
			Document:    "%7B%22Version%22%3A%20%222012-10-17%22%2C%20%20%22Statement%22%3A%20%5B%7B%22Action%22%3A%20%5B%22iam%3AGetUser%22%5D%2C%20%22Resource%22%3A%20%22%2A%22%2C%20%22Effect%22%3A%20%22Allow%22%2C%20%22Sid%22%3A%20%22editor%22%7D%5D%7D",
			ServiceUUID: storage.UUID,
		})
		require.NoError(t, err)

		policies, err := svc.GetManagedObjectStoragePolicies(ctx, &request.GetManagedObjectStoragePoliciesRequest{ServiceUUID: storage.UUID})
		assert.NoError(t, err)
		assert.Len(t, policies, 6)

		err = svc.DeleteManagedObjectStoragePolicy(ctx, &request.DeleteManagedObjectStoragePolicyRequest{
			ServiceUUID: storage.UUID,
			Name:        policy.Name,
		})
		require.NoError(t, err)
	})
}

func TestGetManagedObjectStoragePolicy(t *testing.T) {
	record(t, "getmanagedobjectstoragepolicy", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service) {
		storage, err := createManagedObjectStorage(ctx, svc)
		require.NoError(t, err)

		defer func(uuid string) {
			err = deleteManagedObjectStorageAndUsers(ctx, svc, uuid)
			require.NoError(t, err)
		}(storage.UUID)

		policy, err := svc.CreateManagedObjectStoragePolicy(ctx, &request.CreateManagedObjectStoragePolicyRequest{
			Name:        "testpolicy",
			ServiceUUID: storage.UUID,
			Document:    "%7B%22Version%22%3A%20%222012-10-17%22%2C%20%20%22Statement%22%3A%20%5B%7B%22Action%22%3A%20%5B%22iam%3AGetUser%22%5D%2C%20%22Resource%22%3A%20%22%2A%22%2C%20%22Effect%22%3A%20%22Allow%22%2C%20%22Sid%22%3A%20%22editor%22%7D%5D%7D",
		})
		require.NoError(t, err)
		assert.Equal(t, policy.Name, "testpolicy")

		_, err = svc.GetManagedObjectStoragePolicy(ctx, &request.GetManagedObjectStoragePolicyRequest{ServiceUUID: storage.UUID, Name: policy.Name})
		assert.NoError(t, err)

		err = svc.DeleteManagedObjectStoragePolicy(ctx, &request.DeleteManagedObjectStoragePolicyRequest{
			ServiceUUID: storage.UUID,
			Name:        policy.Name,
		})
		require.NoError(t, err)
	})
}

func TestDeleteManagedObjectStoragePolicy(t *testing.T) {
	record(t, "deletemanagedobjectstoragepolicy", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service) {
		storage, err := createManagedObjectStorage(ctx, svc)
		require.NoError(t, err)

		defer func(uuid string) {
			err = deleteManagedObjectStorageAndUsers(ctx, svc, uuid)
			require.NoError(t, err)
		}(storage.UUID)

		policy, err := svc.CreateManagedObjectStoragePolicy(ctx, &request.CreateManagedObjectStoragePolicyRequest{
			Name:        "testpolicy",
			ServiceUUID: storage.UUID,
			Document:    "%7B%22Version%22%3A%20%222012-10-17%22%2C%20%20%22Statement%22%3A%20%5B%7B%22Action%22%3A%20%5B%22iam%3AGetUser%22%5D%2C%20%22Resource%22%3A%20%22%2A%22%2C%20%22Effect%22%3A%20%22Allow%22%2C%20%22Sid%22%3A%20%22editor%22%7D%5D%7D",
		})
		require.NoError(t, err)

		err = svc.DeleteManagedObjectStoragePolicy(ctx, &request.DeleteManagedObjectStoragePolicyRequest{
			ServiceUUID: storage.UUID,
			Name:        policy.Name,
		})
		require.NoError(t, err)
	})
}

func TestAttachManagedObjectStorageUserPolicy(t *testing.T) {
	record(t, "attachmanagedobjectstorageuserpolicy", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service) {
		storage, err := createManagedObjectStorage(ctx, svc)
		require.NoError(t, err)

		defer func(uuid string) {
			err = deleteManagedObjectStorageAndUsers(ctx, svc, uuid)
			require.NoError(t, err)
		}(storage.UUID)

		policy, err := svc.CreateManagedObjectStoragePolicy(ctx, &request.CreateManagedObjectStoragePolicyRequest{
			Name:        "testpolicy",
			Description: "description2",
			Document:    "%7B%22Version%22%3A%20%222012-10-17%22%2C%20%20%22Statement%22%3A%20%5B%7B%22Action%22%3A%20%5B%22iam%3AGetUser%22%5D%2C%20%22Resource%22%3A%20%22%2A%22%2C%20%22Effect%22%3A%20%22Allow%22%2C%20%22Sid%22%3A%20%22editor%22%7D%5D%7D",
			ServiceUUID: storage.UUID,
		})
		require.NoError(t, err)

		user, err := svc.CreateManagedObjectStorageUser(ctx, &request.CreateManagedObjectStorageUserRequest{
			Username:    "testuser",
			ServiceUUID: storage.UUID,
		})
		assert.NoError(t, err)

		err = svc.AttachManagedObjectStorageUserPolicy(ctx, &request.AttachManagedObjectStorageUserPolicyRequest{
			Name:        policy.Name,
			ServiceUUID: storage.UUID,
			Username:    user.Username,
		})
		assert.NoError(t, err)

		err = svc.DetachManagedObjectStorageUserPolicy(ctx, &request.DetachManagedObjectStorageUserPolicyRequest{
			ServiceUUID: storage.UUID,
			Username:    user.Username,
			Name:        policy.Name,
		})
		assert.NoError(t, err)

		err = svc.DeleteManagedObjectStoragePolicy(ctx, &request.DeleteManagedObjectStoragePolicyRequest{
			ServiceUUID: storage.UUID,
			Name:        policy.Name,
		})
		require.NoError(t, err)
	})
}

func TestGetManagedObjectStorageUserPolicies(t *testing.T) {
	record(t, "getmanagedobjectstorageuserpolicies", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service) {
		storage, err := createManagedObjectStorage(ctx, svc)
		require.NoError(t, err)

		defer func(uuid string) {
			err = deleteManagedObjectStorageAndUsers(ctx, svc, uuid)
			require.NoError(t, err)
		}(storage.UUID)

		policy, err := svc.CreateManagedObjectStoragePolicy(ctx, &request.CreateManagedObjectStoragePolicyRequest{
			Name:        "testpolicy",
			Description: "description2",
			Document:    "%7B%22Version%22%3A%20%222012-10-17%22%2C%20%20%22Statement%22%3A%20%5B%7B%22Action%22%3A%20%5B%22iam%3AGetUser%22%5D%2C%20%22Resource%22%3A%20%22%2A%22%2C%20%22Effect%22%3A%20%22Allow%22%2C%20%22Sid%22%3A%20%22editor%22%7D%5D%7D",
			ServiceUUID: storage.UUID,
		})
		require.NoError(t, err)

		user, err := svc.CreateManagedObjectStorageUser(ctx, &request.CreateManagedObjectStorageUserRequest{
			Username:    "testuser",
			ServiceUUID: storage.UUID,
		})
		assert.NoError(t, err)

		err = svc.AttachManagedObjectStorageUserPolicy(ctx, &request.AttachManagedObjectStorageUserPolicyRequest{
			Name:        policy.Name,
			ServiceUUID: storage.UUID,
			Username:    user.Username,
		})
		assert.NoError(t, err)

		policies, err := svc.GetManagedObjectStorageUserPolicies(ctx, &request.GetManagedObjectStorageUserPoliciesRequest{
			ServiceUUID: storage.UUID,
			Username:    user.Username,
		})
		assert.NoError(t, err)
		assert.Len(t, policies, 1)

		err = svc.DetachManagedObjectStorageUserPolicy(ctx, &request.DetachManagedObjectStorageUserPolicyRequest{
			ServiceUUID: storage.UUID,
			Username:    user.Username,
			Name:        policy.Name,
		})
		assert.NoError(t, err)

		err = svc.DeleteManagedObjectStoragePolicy(ctx, &request.DeleteManagedObjectStoragePolicyRequest{
			ServiceUUID: storage.UUID,
			Name:        policy.Name,
		})
		require.NoError(t, err)
	})
}

func TestDetachManagedObjectStorageUserPolicy(t *testing.T) {
	record(t, "detachmanagedobjectstorageuserpolicy", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service) {
		storage, err := createManagedObjectStorage(ctx, svc)
		require.NoError(t, err)

		defer func(uuid string) {
			err = deleteManagedObjectStorageAndUsers(ctx, svc, uuid)
			require.NoError(t, err)
		}(storage.UUID)

		policy, err := svc.CreateManagedObjectStoragePolicy(ctx, &request.CreateManagedObjectStoragePolicyRequest{
			Name:        "testpolicy",
			Description: "description2",
			Document:    "%7B%22Version%22%3A%20%222012-10-17%22%2C%20%20%22Statement%22%3A%20%5B%7B%22Action%22%3A%20%5B%22iam%3AGetUser%22%5D%2C%20%22Resource%22%3A%20%22%2A%22%2C%20%22Effect%22%3A%20%22Allow%22%2C%20%22Sid%22%3A%20%22editor%22%7D%5D%7D",
			ServiceUUID: storage.UUID,
		})
		require.NoError(t, err)

		user, err := svc.CreateManagedObjectStorageUser(ctx, &request.CreateManagedObjectStorageUserRequest{
			Username:    "testuser",
			ServiceUUID: storage.UUID,
		})
		assert.NoError(t, err)

		err = svc.AttachManagedObjectStorageUserPolicy(ctx, &request.AttachManagedObjectStorageUserPolicyRequest{
			Name:        policy.Name,
			ServiceUUID: storage.UUID,
			Username:    user.Username,
		})
		assert.NoError(t, err)

		err = svc.DetachManagedObjectStorageUserPolicy(ctx, &request.DetachManagedObjectStorageUserPolicyRequest{
			ServiceUUID: storage.UUID,
			Username:    user.Username,
			Name:        policy.Name,
		})
		assert.NoError(t, err)

		err = svc.DeleteManagedObjectStoragePolicy(ctx, &request.DeleteManagedObjectStoragePolicyRequest{
			ServiceUUID: storage.UUID,
			Name:        policy.Name,
		})
		require.NoError(t, err)
	})
}

func TestManagedObjectStorageCustomDomains(t *testing.T) {
	record(t, "managedobjectstoragecustomdomains", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service) {
		storage, err := createManagedObjectStorage(ctx, svc)
		require.NoError(t, err)

		defer func(uuid string) {
			err = deleteManagedObjectStorageAndUsers(ctx, svc, uuid)
			require.NoError(t, err)
		}(storage.UUID)

		domainName := "obj.example.com"
		err = svc.CreateManagedObjectStorageCustomDomain(ctx, &request.CreateManagedObjectStorageCustomDomainRequest{
			DomainName:  domainName,
			Type:        "public",
			ServiceUUID: storage.UUID,
		})
		require.NoError(t, err)

		domains, err := svc.GetManagedObjectStorageCustomDomains(ctx, &request.GetManagedObjectStorageCustomDomainsRequest{
			ServiceUUID: storage.UUID,
		})
		assert.NoError(t, err)
		assert.Len(t, domains, 1)
		assert.Equal(t, domains[0].DomainName, domainName)

		objsto, err := svc.GetManagedObjectStorage(ctx, &request.GetManagedObjectStorageRequest{
			UUID: storage.UUID,
		})
		assert.NoError(t, err)
		assert.Len(t, objsto.CustomDomains, 1)

		domain, err := svc.GetManagedObjectStorageCustomDomain(ctx, &request.GetManagedObjectStorageCustomDomainRequest{
			ServiceUUID: storage.UUID,
			DomainName:  domainName,
		})
		assert.NoError(t, err)
		assert.Equal(t, domain.DomainName, domainName)

		modifiedDomainName := "objects.example.com"
		domain, err = svc.ModifyManagedObjectStorageCustomDomain(ctx, &request.ModifyManagedObjectStorageCustomDomainRequest{
			ServiceUUID: storage.UUID,
			DomainName:  domainName,
			CustomDomain: request.ModifyCustomDomain{
				Type:       "public",
				DomainName: modifiedDomainName,
			},
		})
		assert.NoError(t, err)
		assert.Equal(t, domain.DomainName, modifiedDomainName)

		err = svc.DeleteManagedObjectStorageCustomDomain(ctx, &request.DeleteManagedObjectStorageCustomDomainRequest{
			ServiceUUID: storage.UUID,
			DomainName:  modifiedDomainName,
		})
		require.NoError(t, err)

		domains, err = svc.GetManagedObjectStorageCustomDomains(ctx, &request.GetManagedObjectStorageCustomDomainsRequest{
			ServiceUUID: storage.UUID,
		})
		assert.NoError(t, err)
		assert.Len(t, domains, 0)

		objsto, err = svc.GetManagedObjectStorage(ctx, &request.GetManagedObjectStorageRequest{
			UUID: storage.UUID,
		})
		assert.NoError(t, err)
		assert.Len(t, objsto.CustomDomains, 0)
	})
}

func createManagedObjectStorage(ctx context.Context, svc *Service) (*upcloud.ManagedObjectStorage, error) {
	regions, err := svc.GetManagedObjectStorageRegions(ctx, &request.GetManagedObjectStorageRegionsRequest{})
	if err != nil {
		return nil, err
	}

	return svc.CreateManagedObjectStorage(ctx, &request.CreateManagedObjectStorageRequest{
		Name:             "go-sdk-integration-test",
		ConfiguredStatus: upcloud.ManagedObjectStorageConfiguredStatusStarted,
		Labels: []upcloud.Label{
			{
				Key:   "example-key",
				Value: "example-value",
			},
		},
		Networks: []upcloud.ManagedObjectStorageNetwork{
			{
				Family: "IPv4",
				Name:   "example-public-network",
				Type:   "public",
			},
		},
		Region: regions[0].Name,
	})
}

func deleteManagedObjectStorageAndUsers(ctx context.Context, svc *Service, uuid string) error {
	users, err := svc.GetManagedObjectStorageUsers(ctx, &request.GetManagedObjectStorageUsersRequest{ServiceUUID: uuid})
	if err != nil {
		return err
	}

	for _, user := range users {
		errDelete := svc.DeleteManagedObjectStorageUser(ctx, &request.DeleteManagedObjectStorageUserRequest{ServiceUUID: uuid, Username: user.Username})
		if errDelete != nil {
			return errDelete
		}
	}

	return svc.DeleteManagedObjectStorage(ctx, &request.DeleteManagedObjectStorageRequest{
		UUID: uuid,
	})
}
