package service

import (
	"context"
	"testing"

	"github.com/UpCloudLtd/upcloud-go-api/v6/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v6/upcloud/request"
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

		defer func() {
			err = svc.DeleteManagedObjectStorage(ctx, &request.DeleteManagedObjectStorageRequest{UUID: storage.UUID})
			require.NoError(t, err)
		}()

		require.Equal(t, storage.ConfiguredStatus, upcloud.ManagedObjectStorageConfiguredStatusStarted)
		require.Len(t, storage.Labels, 1)
		require.Len(t, storage.Networks, 1)
		require.Len(t, storage.Users, 1)
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

		defer func() {
			err = svc.DeleteManagedObjectStorage(ctx, &request.DeleteManagedObjectStorageRequest{UUID: storage.UUID})
			require.NoError(t, err)
		}()

		storages, err = svc.GetManagedObjectStorages(ctx, &request.GetManagedObjectStoragesRequest{})
		require.NoError(t, err)
		require.Len(t, storages, 1)
	})
}

func TestGetManagedObjectStorageDetails(t *testing.T) {
	record(t, "getmanagedobjectstoragedetails", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service) {
		storage, err := createManagedObjectStorage(ctx, svc)
		require.NoError(t, err)

		defer func() {
			err = svc.DeleteManagedObjectStorage(ctx, &request.DeleteManagedObjectStorageRequest{UUID: storage.UUID})
			require.NoError(t, err)
		}()

		storage, err = svc.GetManagedObjectStorage(ctx, &request.GetManagedObjectStorageRequest{UUID: storage.UUID})
		require.NoError(t, err)
		require.Equal(t, storage.ConfiguredStatus, upcloud.ManagedObjectStorageConfiguredStatusStarted)
		require.Len(t, storage.Labels, 1)
		require.Len(t, storage.Networks, 1)
		require.Len(t, storage.Users, 1)
		require.NotEmpty(t, storage.UUID)
	})
}

func TestReplaceManagedObjectStorage(t *testing.T) {
	record(t, "replacemanagedobjectstorage", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service) {
		storage, err := createManagedObjectStorage(ctx, svc)
		require.NoError(t, err)

		defer func() {
			err = svc.DeleteManagedObjectStorage(ctx, &request.DeleteManagedObjectStorageRequest{UUID: storage.UUID})
			require.NoError(t, err)
		}()

		storage, err = svc.ReplaceManagedObjectStorage(ctx, &request.ReplaceManagedObjectStorageRequest{
			ConfiguredStatus: upcloud.ManagedObjectStorageConfiguredStatusStopped,
			Networks: []upcloud.ManagedObjectStorageNetwork{{
				Family: "IPv4",
				Name:   "replaced-network",
				Type:   "public",
			}},
			Users: []request.ManagedObjectStorageUser{{Username: storage.Users[0].Username}},
			UUID:  storage.UUID,
			Name:  "test2",
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

		defer func() {
			err = svc.DeleteManagedObjectStorage(ctx, &request.DeleteManagedObjectStorageRequest{UUID: storage.UUID})
			require.NoError(t, err)
		}()

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

		defer func() {
			err = svc.DeleteManagedObjectStorage(ctx, &request.DeleteManagedObjectStorageRequest{UUID: storage.UUID})
			require.NoError(t, err)
		}()

		m, err := svc.GetManagedObjectStorageMetrics(ctx, &request.GetManagedObjectStorageMetricsRequest{ServiceUUID: storage.UUID})
		require.NoError(t, err)
		require.Equal(t, m.TotalObjects, 0)
		require.Equal(t, m.TotalSizeBytes, 0)
	})
}

func TestGetManagedObjectStorageBucketMetrics(t *testing.T) {
	record(t, "getmanagedobjectstoragebucketmetrics", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service) {
		storage, err := createManagedObjectStorage(ctx, svc)
		require.NoError(t, err)

		defer func() {
			err = svc.DeleteManagedObjectStorage(ctx, &request.DeleteManagedObjectStorageRequest{UUID: storage.UUID})
			require.NoError(t, err)
		}()

		m, err := svc.GetManagedObjectStorageBucketMetrics(ctx, &request.GetManagedObjectStorageBucketMetricsRequest{ServiceUUID: storage.UUID})
		require.NoError(t, err)
		assert.Len(t, m, 0)
	})
}

func TestCreateManagedObjectStorageNetwork(t *testing.T) {
	record(t, "createmanagedobjectstoragenetwork", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service) {
		storage, err := createManagedObjectStorage(ctx, svc)
		require.NoError(t, err)

		defer func() {
			err = svc.DeleteManagedObjectStorage(ctx, &request.DeleteManagedObjectStorageRequest{UUID: storage.UUID})
			require.NoError(t, err)
		}()

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

		defer func() {
			err = svc.DeleteManagedObjectStorage(ctx, &request.DeleteManagedObjectStorageRequest{UUID: storage.UUID})
			require.NoError(t, err)
		}()

		networks, err := svc.GetManagedObjectStorageNetworks(ctx, &request.GetManagedObjectStorageNetworksRequest{ServiceUUID: storage.UUID})
		require.NoError(t, err)
		require.Len(t, networks, 1)
	})
}

func TestGetManagedObjectStorageNetwork(t *testing.T) {
	record(t, "getmanagedobjectstoragenetwork", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service) {
		storage, err := createManagedObjectStorage(ctx, svc)
		require.NoError(t, err)

		defer func() {
			err = svc.DeleteManagedObjectStorage(ctx, &request.DeleteManagedObjectStorageRequest{UUID: storage.UUID})
			require.NoError(t, err)
		}()

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

		defer func() {
			err = svc.DeleteManagedObjectStorage(ctx, &request.DeleteManagedObjectStorageRequest{UUID: storage.UUID})
			require.NoError(t, err)
		}()

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

		defer func() {
			err = svc.DeleteManagedObjectStorage(ctx, &request.DeleteManagedObjectStorageRequest{UUID: storage.UUID})
			require.NoError(t, err)
		}()

		user, err := svc.CreateManagedObjectStorageUser(ctx, &request.CreateManagedObjectStorageUserRequest{
			Username:    "test2",
			ServiceUUID: storage.UUID,
		})
		require.NoError(t, err)
		require.Equal(t, user.Username, "test2")
	})
}

func TestGetManagedObjectStorageUsers(t *testing.T) {
	record(t, "getmanagedobjectstorageusers", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service) {
		storage, err := createManagedObjectStorage(ctx, svc)
		require.NoError(t, err)

		defer func() {
			err = svc.DeleteManagedObjectStorage(ctx, &request.DeleteManagedObjectStorageRequest{UUID: storage.UUID})
			require.NoError(t, err)
		}()

		users, err := svc.GetManagedObjectStorageUsers(ctx, &request.GetManagedObjectStorageUsersRequest{ServiceUUID: storage.UUID})
		require.NoError(t, err)
		require.Len(t, users, 1)
	})
}

func TestGetManagedObjectStorageUser(t *testing.T) {
	record(t, "getmanagedobjectstorageuser", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service) {
		storage, err := createManagedObjectStorage(ctx, svc)
		require.NoError(t, err)

		defer func() {
			err = svc.DeleteManagedObjectStorage(ctx, &request.DeleteManagedObjectStorageRequest{UUID: storage.UUID})
			require.NoError(t, err)
		}()

		user, err := svc.CreateManagedObjectStorageUser(ctx, &request.CreateManagedObjectStorageUserRequest{
			Username:    "test2",
			ServiceUUID: storage.UUID,
		})
		require.NoError(t, err)
		require.Equal(t, user.Username, "test2")

		_, err = svc.GetManagedObjectStorageUser(ctx, &request.GetManagedObjectStorageUserRequest{ServiceUUID: storage.UUID, Username: user.Username})
		require.NoError(t, err)
	})
}

func TestDeleteManagedObjectStorageUser(t *testing.T) {
	record(t, "deletemanagedobjectstorageuser", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service) {
		storage, err := createManagedObjectStorage(ctx, svc)
		require.NoError(t, err)

		defer func() {
			err = svc.DeleteManagedObjectStorage(ctx, &request.DeleteManagedObjectStorageRequest{UUID: storage.UUID})
			require.NoError(t, err)
		}()

		user, err := svc.CreateManagedObjectStorageUser(ctx, &request.CreateManagedObjectStorageUserRequest{
			Username:    "test2",
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

		defer func() {
			err = svc.DeleteManagedObjectStorage(ctx, &request.DeleteManagedObjectStorageRequest{UUID: storage.UUID})
			require.NoError(t, err)
		}()

		user, err := svc.CreateManagedObjectStorageUser(ctx, &request.CreateManagedObjectStorageUserRequest{
			Username:    "test2",
			ServiceUUID: storage.UUID,
		})
		require.NoError(t, err)
		require.Equal(t, user.Username, "test2")

		_, err = svc.WaitForManagedObjectStorageUserOperationalState(context.Background(), &request.WaitForManagedObjectStorageUserOperationalStateRequest{
			ServiceUUID:  storage.UUID,
			Username:     user.Username,
			DesiredState: upcloud.ManagedObjectStorageUserOperationalStateReady,
		})
		require.NoError(t, err)

		accessKey, err := svc.CreateManagedObjectStorageUserAccessKey(ctx, &request.CreateManagedObjectStorageUserAccessKeyRequest{
			Username:    user.Username,
			ServiceUUID: storage.UUID,
			Name:        "example-access-key",
			Enabled:     upcloud.BoolPtr(false),
		})
		require.NoError(t, err)
		require.Equal(t, "example-access-key", accessKey.Name)
		require.Equal(t, false, accessKey.Enabled)
		require.NotEmpty(t, accessKey.AccessKeyId)
		require.NotEmpty(t, accessKey.SecretAccessKey)
	})
}

func TestGetManagedObjectStorageUserAccesKeys(t *testing.T) {
	record(t, "getmanagedobjectstorageuseraccesskeys", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service) {
		storage, err := createManagedObjectStorage(ctx, svc)
		require.NoError(t, err)

		defer func() {
			err = svc.DeleteManagedObjectStorage(ctx, &request.DeleteManagedObjectStorageRequest{UUID: storage.UUID})
			require.NoError(t, err)
		}()

		_, err = svc.WaitForManagedObjectStorageUserOperationalState(context.Background(), &request.WaitForManagedObjectStorageUserOperationalStateRequest{
			ServiceUUID:  storage.UUID,
			Username:     storage.Users[0].Username,
			DesiredState: upcloud.ManagedObjectStorageUserOperationalStateReady,
		})
		require.NoError(t, err)

		_, err = svc.CreateManagedObjectStorageUserAccessKey(ctx, &request.CreateManagedObjectStorageUserAccessKeyRequest{
			Username:    storage.Users[0].Username,
			ServiceUUID: storage.UUID,
			Name:        "example-access-key",
			Enabled:     upcloud.BoolPtr(false),
		})
		require.NoError(t, err)

		accessKeys, err := svc.GetManagedObjectStorageUserAccessKeys(ctx, &request.GetManagedObjectStorageUserAccessKeysRequest{
			ServiceUUID: storage.UUID,
			Username:    storage.Users[0].Username,
		})
		require.NoError(t, err)
		require.Len(t, accessKeys, 1)
	})
}

func TestGetManagedObjectStorageUserAccessKey(t *testing.T) {
	record(t, "getmanagedobjectstorageuseraccesskey", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service) {
		storage, err := createManagedObjectStorage(ctx, svc)
		require.NoError(t, err)

		defer func() {
			err = svc.DeleteManagedObjectStorage(ctx, &request.DeleteManagedObjectStorageRequest{UUID: storage.UUID})
			require.NoError(t, err)
		}()

		_, err = svc.WaitForManagedObjectStorageUserOperationalState(context.Background(), &request.WaitForManagedObjectStorageUserOperationalStateRequest{
			ServiceUUID:  storage.UUID,
			Username:     storage.Users[0].Username,
			DesiredState: upcloud.ManagedObjectStorageUserOperationalStateReady,
		})
		require.NoError(t, err)

		_, err = svc.CreateManagedObjectStorageUserAccessKey(ctx, &request.CreateManagedObjectStorageUserAccessKeyRequest{
			Username:    storage.Users[0].Username,
			ServiceUUID: storage.UUID,
			Name:        "example-access-key",
			Enabled:     upcloud.BoolPtr(false),
		})
		require.NoError(t, err)

		accessKeys, err := svc.GetManagedObjectStorageUserAccessKeys(ctx, &request.GetManagedObjectStorageUserAccessKeysRequest{
			ServiceUUID: storage.UUID,
			Username:    storage.Users[0].Username,
		})
		require.NoError(t, err)
		require.Len(t, accessKeys, 1)

		accessKey, err := svc.GetManagedObjectStorageUserAccessKey(ctx, &request.GetManagedObjectStorageUserAccessKeyRequest{
			ServiceUUID: storage.UUID,
			Username:    storage.Users[0].Username,
			Name:        "example-access-key",
		})
		require.NoError(t, err)
		require.Empty(t, accessKey.SecretAccessKey)
	})
}

func TestModifyManagedObjectStorageUserAccessKey(t *testing.T) {
	record(t, "modifymanagedobjectstorageuseraccesskey", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service) {
		storage, err := createManagedObjectStorage(ctx, svc)
		require.NoError(t, err)

		defer func() {
			err = svc.DeleteManagedObjectStorage(ctx, &request.DeleteManagedObjectStorageRequest{UUID: storage.UUID})
			require.NoError(t, err)
		}()

		_, err = svc.WaitForManagedObjectStorageUserOperationalState(context.Background(), &request.WaitForManagedObjectStorageUserOperationalStateRequest{
			ServiceUUID:  storage.UUID,
			Username:     storage.Users[0].Username,
			DesiredState: upcloud.ManagedObjectStorageUserOperationalStateReady,
		})
		require.NoError(t, err)

		accessKey, err := svc.CreateManagedObjectStorageUserAccessKey(ctx, &request.CreateManagedObjectStorageUserAccessKeyRequest{
			Username:    storage.Users[0].Username,
			ServiceUUID: storage.UUID,
			Name:        "example-access-key",
			Enabled:     upcloud.BoolPtr(false),
		})
		require.NoError(t, err)

		accessKey, err = svc.ModifyManagedObjectStorageUserAccessKey(ctx, &request.ModifyManagedObjectStorageUserAccessKeyRequest{
			ServiceUUID: storage.UUID,
			Username:    storage.Users[0].Username,
			Name:        accessKey.Name,
			Enabled:     upcloud.BoolPtr(true),
		})
		require.NoError(t, err)
		require.Equal(t, true, accessKey.Enabled)

		accessKey, err = svc.ModifyManagedObjectStorageUserAccessKey(ctx, &request.ModifyManagedObjectStorageUserAccessKeyRequest{
			ServiceUUID: storage.UUID,
			Username:    storage.Users[0].Username,
			Name:        accessKey.Name,
			Enabled:     upcloud.BoolPtr(false),
		})
		require.NoError(t, err)
		require.Equal(t, false, accessKey.Enabled)
	})
}

func TestDeleteManagedObjectStorageUserAccessKey(t *testing.T) {
	record(t, "deletemanagedobjectstorageuseraccesskey", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service) {
		storage, err := createManagedObjectStorage(ctx, svc)
		require.NoError(t, err)

		defer func() {
			err = svc.DeleteManagedObjectStorage(ctx, &request.DeleteManagedObjectStorageRequest{UUID: storage.UUID})
			require.NoError(t, err)
		}()

		_, err = svc.WaitForManagedObjectStorageUserOperationalState(context.Background(), &request.WaitForManagedObjectStorageUserOperationalStateRequest{
			ServiceUUID:  storage.UUID,
			Username:     storage.Users[0].Username,
			DesiredState: upcloud.ManagedObjectStorageUserOperationalStateReady,
		})
		require.NoError(t, err)

		accessKey, err := svc.CreateManagedObjectStorageUserAccessKey(ctx, &request.CreateManagedObjectStorageUserAccessKeyRequest{
			Username:    storage.Users[0].Username,
			ServiceUUID: storage.UUID,
			Name:        "example-access-key",
			Enabled:     upcloud.BoolPtr(false),
		})
		require.NoError(t, err)

		err = svc.DeleteManagedObjectStorageUserAccessKey(ctx, &request.DeleteManagedObjectStorageUserAccessKeyRequest{
			ServiceUUID: storage.UUID,
			Username:    storage.Users[0].Username,
			Name:        accessKey.Name,
		})
		require.NoError(t, err)
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
		Users:  []request.ManagedObjectStorageUser{{Username: "test"}},
	})
}
