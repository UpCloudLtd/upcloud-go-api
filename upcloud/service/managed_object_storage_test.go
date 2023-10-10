package service

import (
	"context"
	"fmt"
	"testing"
	"time"

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
		require.Equal(t, storage.ConfiguredStatus, upcloud.ManagedObjectStorageConfiguredStatusStarted)
		require.Len(t, storage.Labels, 1)
		require.Len(t, storage.Networks, 1)
		require.Len(t, storage.Users, 1)
		require.NotEmpty(t, storage.UUID)

		err = deleteManagedObjectStorage(ctx, svc, storage.UUID)
		require.NoError(t, err)
	})
}

func TestGetManagedObjectStorages(t *testing.T) {
	record(t, "getmanagedobjectstorages", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service) {
		storages, err := svc.GetManagedObjectStorages(ctx, &request.GetManagedObjectStoragesRequest{})
		require.NoError(t, err)
		require.Len(t, storages, 0)

		storage, err := createManagedObjectStorage(ctx, svc)
		require.NoError(t, err)

		storages, err = svc.GetManagedObjectStorages(ctx, &request.GetManagedObjectStoragesRequest{})
		require.NoError(t, err)
		require.Len(t, storages, 1)

		err = deleteManagedObjectStorage(ctx, svc, storage.UUID)
		require.NoError(t, err)
	})
}

func TestGetManagedObjectStorageDetails(t *testing.T) {
	record(t, "getmanagedobjectstoragedetails", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service) {
		storage, err := createManagedObjectStorage(ctx, svc)
		require.NoError(t, err)

		storage, err = svc.GetManagedObjectStorageDetails(ctx, &request.GetManagedObjectStorageDetailsRequest{UUID: storage.UUID})
		require.NoError(t, err)
		require.Equal(t, storage.ConfiguredStatus, upcloud.ManagedObjectStorageConfiguredStatusStarted)
		require.Len(t, storage.Labels, 1)
		require.Len(t, storage.Networks, 1)
		require.Len(t, storage.Users, 1)
		require.NotEmpty(t, storage.UUID)

		err = deleteManagedObjectStorage(ctx, svc, storage.UUID)
		require.NoError(t, err)
	})
}

func TestReplaceManagedObjectStorage(t *testing.T) {
	record(t, "replacemanagedobjectstorage", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service) {
		storage, err := createManagedObjectStorage(ctx, svc)
		require.NoError(t, err)

		storage, err = svc.ReplaceManagedObjectStorage(ctx, &request.ReplaceManagedObjectStorageRequest{
			ConfiguredStatus: upcloud.ManagedObjectStorageConfiguredStatusStopped,
			Networks: []upcloud.ManagedObjectStorageNetwork{{
				Family: "IPv4",
				Name:   "replaced-network",
				Type:   "public",
			}},
			Users: []request.ManagedObjectStorageUser{{Username: storage.Users[0].Username}},
			UUID:  storage.UUID,
		})
		require.NoError(t, err)
		require.Equal(t, upcloud.ManagedObjectStorageConfiguredStatusStopped, storage.ConfiguredStatus)
		require.Equal(t, storage.Networks[0].Name, "replaced-network")
		require.Len(t, storage.Labels, 0)

		err = deleteManagedObjectStorage(ctx, svc, storage.UUID)
		require.NoError(t, err)
	})
}

func TestModifyManagedObjectStorage(t *testing.T) {
	record(t, "modifymanagedobjectstorage", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service) {
		storage, err := createManagedObjectStorage(ctx, svc)
		require.NoError(t, err)

		status := upcloud.ManagedObjectStorageConfiguredStatusStopped
		storage, err = svc.ModifyManagedObjectStorage(ctx, &request.ModifyManagedObjectStorageRequest{
			ConfiguredStatus: &status,
			Labels:           &[]upcloud.Label{},
			UUID:             storage.UUID,
		})
		require.NoError(t, err)
		require.Equal(t, upcloud.ManagedObjectStorageConfiguredStatusStopped, storage.ConfiguredStatus)
		require.Len(t, storage.Labels, 0)

		err = deleteManagedObjectStorage(ctx, svc, storage.UUID)
		require.NoError(t, err)
	})
}

func TestGetManagedObjectStorageMetrics(t *testing.T) {
	record(t, "getmanagedobjectstoragemetrics", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service) {
		storage, err := createManagedObjectStorage(ctx, svc)
		require.NoError(t, err)

		m, err := svc.GetManagedObjectStorageMetrics(ctx, &request.GetManagedObjectStorageMetricsRequest{UUID: storage.UUID})
		require.NoError(t, err)
		require.Equal(t, m.TotalObjects, 0)
		require.Equal(t, m.TotalSizeBytes, 0)

		err = deleteManagedObjectStorage(ctx, svc, storage.UUID)
		require.NoError(t, err)
	})
}

func TestGetManagedObjectStorageBucketMetrics(t *testing.T) {
	record(t, "getmanagedobjectstoragebucketmetrics", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service) {
		storage, err := createManagedObjectStorage(ctx, svc)
		require.NoError(t, err)

		m, err := svc.GetManagedObjectStorageBucketMetrics(ctx, &request.GetManagedObjectStorageBucketMetricsRequest{ServiceUUID: storage.UUID})
		require.NoError(t, err)
		assert.Len(t, m, 0)

		err = deleteManagedObjectStorage(ctx, svc, storage.UUID)
		require.NoError(t, err)
	})
}

func TestCreateManagedObjectStorageNetwork(t *testing.T) {
	record(t, "createmanagedobjectstoragenetwork", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service) {
		storage, err := createManagedObjectStorage(ctx, svc)
		require.NoError(t, err)

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

		err = deleteManagedObjectStorage(ctx, svc, storage.UUID)
		require.NoError(t, err)

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

		networks, err := svc.GetManagedObjectStorageNetworks(ctx, &request.GetManagedObjectStorageNetworksRequest{ServiceUUID: storage.UUID})
		require.NoError(t, err)
		require.Len(t, networks, 1)

		err = deleteManagedObjectStorage(ctx, svc, storage.UUID)
		require.NoError(t, err)
	})
}

func TestGetManagedObjectStorageNetwork(t *testing.T) {
	record(t, "getmanagedobjectstoragenetwork", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service) {
		storage, err := createManagedObjectStorage(ctx, svc)
		require.NoError(t, err)

		networks, err := svc.GetManagedObjectStorageNetworks(ctx, &request.GetManagedObjectStorageNetworksRequest{ServiceUUID: storage.UUID})
		require.NoError(t, err)
		require.Len(t, networks, 1)

		network, err := svc.GetManagedObjectStorageNetwork(ctx, &request.GetManagedObjectStorageNetworkRequest{ServiceUUID: storage.UUID, NetworkName: networks[0].Name})
		require.NoError(t, err)
		require.Equal(t, networks[0].Name, network.Name)

		err = deleteManagedObjectStorage(ctx, svc, storage.UUID)
		require.NoError(t, err)
	})
}

func TestDeleteManagedObjectStorageNetwork(t *testing.T) {
	record(t, "deletemanagedobjectstoragenetwork", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service) {
		storage, err := createManagedObjectStorage(ctx, svc)
		require.NoError(t, err)

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

		err = deleteManagedObjectStorage(ctx, svc, storage.UUID)
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

		account, err := createSubAccount(ctx, svc, "managedobjectstorage2")
		require.NoError(t, err)

		user, err := svc.CreateManagedObjectStorageUser(ctx, &request.CreateManagedObjectStorageUserRequest{
			Username:    account.Username,
			ServiceUUID: storage.UUID,
		})
		require.NoError(t, err)
		require.Equal(t, user.Username, account.Username)

		err = deleteManagedObjectStorage(ctx, svc, storage.UUID)
		require.NoError(t, err)
	})
}

func TestGetManagedObjectStorageUsers(t *testing.T) {
	record(t, "getmanagedobjectstorageusers", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service) {
		storage, err := createManagedObjectStorage(ctx, svc)
		require.NoError(t, err)

		users, err := svc.GetManagedObjectStorageUsers(ctx, &request.GetManagedObjectStorageUsersRequest{ServiceUUID: storage.UUID})
		require.NoError(t, err)
		require.Len(t, users, 1)

		err = deleteManagedObjectStorage(ctx, svc, storage.UUID)
		require.NoError(t, err)
	})
}

func TestGetManagedObjectStorageUser(t *testing.T) {
	record(t, "getmanagedobjectstorageuser", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service) {
		storage, err := createManagedObjectStorage(ctx, svc)
		require.NoError(t, err)

		users, err := svc.GetManagedObjectStorageUsers(ctx, &request.GetManagedObjectStorageUsersRequest{ServiceUUID: storage.UUID})
		require.NoError(t, err)
		require.Len(t, users, 1)

		user, err := svc.GetManagedObjectStorageUser(ctx, &request.GetManagedObjectStorageUserRequest{ServiceUUID: storage.UUID, Username: users[0].Username})
		require.NoError(t, err)
		require.Equal(t, users[0].Username, user.Username)

		err = deleteManagedObjectStorage(ctx, svc, storage.UUID)
		require.NoError(t, err)
	})
}

func TestDeleteManagedObjectStorageUser(t *testing.T) {
	record(t, "deletemanagedobjectstorageuser", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service) {
		storage, err := createManagedObjectStorage(ctx, svc)
		require.NoError(t, err)

		account, err := createSubAccount(ctx, svc, "managedobjectstorage2")
		require.NoError(t, err)

		user, err := svc.CreateManagedObjectStorageUser(ctx, &request.CreateManagedObjectStorageUserRequest{
			Username:    account.Username,
			ServiceUUID: storage.UUID,
		})
		require.NoError(t, err)

		err = svc.DeleteManagedObjectStorageUser(ctx, &request.DeleteManagedObjectStorageUserRequest{
			ServiceUUID: storage.UUID,
			Username:    user.Username,
		})
		require.NoError(t, err)

		err = deleteManagedObjectStorage(ctx, svc, storage.UUID)
		require.NoError(t, err)
	})
}

func TestCreateManagedObjectStorageUserAccessKey(t *testing.T) {
	record(t, "createmanagedobjectstorageuseraccesskey", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service) {
		storage, err := createManagedObjectStorage(ctx, svc)
		require.NoError(t, err)

		account, err := createSubAccount(ctx, svc, "managedobjectstorage2")
		require.NoError(t, err)

		user, err := svc.CreateManagedObjectStorageUser(ctx, &request.CreateManagedObjectStorageUserRequest{
			Username:    account.Username,
			ServiceUUID: storage.UUID,
		})
		require.NoError(t, err)
		require.Equal(t, user.Username, account.Username)

		_, err = svc.WaitForManagedObjectStorageUserOperationalState(context.Background(), &request.WaitForManagedObjectStorageUserOperationalStateRequest{
			ServiceUUID:  storage.UUID,
			Timeout:      time.Second * 90,
			Username:     user.Username,
			DesiredState: upcloud.ManagedObjectStorageUserOperationalStateReady,
		})
		require.NoError(t, err)

		accessKey, err := svc.CreateManagedObjectStorageUserAccessKey(ctx, &request.CreateManagedObjectStorageUserAccessKeyRequest{
			Username:    user.Username,
			ServiceUUID: storage.UUID,
			Name:        "example-access-key",
			Enabled:     false,
		})
		require.NoError(t, err)
		require.Equal(t, "example-access-key", accessKey.Name)
		require.NotEmpty(t, accessKey.AccessKeyId)
		require.NotEmpty(t, accessKey.SecretAccessKey)

		err = deleteManagedObjectStorage(ctx, svc, storage.UUID)
		require.NoError(t, err)
	})
}

func TestGetManagedObjectStorageUserAccesKeys(t *testing.T) {
	record(t, "getmanagedobjectstorageuseraccesskeys", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service) {
		storage, err := createManagedObjectStorage(ctx, svc)
		require.NoError(t, err)

		_, err = svc.WaitForManagedObjectStorageUserOperationalState(context.Background(), &request.WaitForManagedObjectStorageUserOperationalStateRequest{
			ServiceUUID:  storage.UUID,
			Timeout:      time.Second * 90,
			Username:     storage.Users[0].Username,
			DesiredState: upcloud.ManagedObjectStorageUserOperationalStateReady,
		})
		require.NoError(t, err)

		_, err = svc.CreateManagedObjectStorageUserAccessKey(ctx, &request.CreateManagedObjectStorageUserAccessKeyRequest{
			Username:    storage.Users[0].Username,
			ServiceUUID: storage.UUID,
			Name:        "example-access-key",
			Enabled:     false,
		})
		require.NoError(t, err)

		accessKeys, err := svc.GetManagedObjectStorageUserAccessKeys(ctx, &request.GetManagedObjectStorageUserAccessKeysRequest{
			ServiceUUID: storage.UUID,
			Username:    storage.Users[0].Username,
		})
		require.NoError(t, err)
		require.Len(t, accessKeys, 1)

		err = deleteManagedObjectStorage(ctx, svc, storage.UUID)
		require.NoError(t, err)
	})
}

func TestGetManagedObjectStorageUserAccessKey(t *testing.T) {
	record(t, "getmanagedobjectstorageuseraccesskey", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service) {
		storage, err := createManagedObjectStorage(ctx, svc)
		require.NoError(t, err)

		_, err = svc.WaitForManagedObjectStorageUserOperationalState(context.Background(), &request.WaitForManagedObjectStorageUserOperationalStateRequest{
			ServiceUUID:  storage.UUID,
			Timeout:      time.Second * 90,
			Username:     storage.Users[0].Username,
			DesiredState: upcloud.ManagedObjectStorageUserOperationalStateReady,
		})
		require.NoError(t, err)

		_, err = svc.CreateManagedObjectStorageUserAccessKey(ctx, &request.CreateManagedObjectStorageUserAccessKeyRequest{
			Username:    storage.Users[0].Username,
			ServiceUUID: storage.UUID,
			Name:        "example-access-key",
			Enabled:     false,
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

		err = deleteManagedObjectStorage(ctx, svc, storage.UUID)
		require.NoError(t, err)
	})
}

func TestModifyManagedObjectStorageUserAccessKey(t *testing.T) {
	record(t, "modifymanagedobjectstorageuseraccesskey", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service) {
		storage, err := createManagedObjectStorage(ctx, svc)
		require.NoError(t, err)

		_, err = svc.WaitForManagedObjectStorageUserOperationalState(context.Background(), &request.WaitForManagedObjectStorageUserOperationalStateRequest{
			ServiceUUID:  storage.UUID,
			Timeout:      time.Second * 90,
			Username:     storage.Users[0].Username,
			DesiredState: upcloud.ManagedObjectStorageUserOperationalStateReady,
		})
		require.NoError(t, err)

		accessKey, err := svc.CreateManagedObjectStorageUserAccessKey(ctx, &request.CreateManagedObjectStorageUserAccessKeyRequest{
			Username:    storage.Users[0].Username,
			ServiceUUID: storage.UUID,
			Name:        "example-access-key",
			Enabled:     false,
		})
		require.NoError(t, err)

		accessKey, err = svc.ModifyManagedObjectStorageUserAccessKey(ctx, &request.ModifyManagedObjectStorageUserAccessKeyRequest{
			ServiceUUID: storage.UUID,
			Username:    storage.Users[0].Username,
			Name:        accessKey.Name,
			Enabled:     true,
		})
		require.NoError(t, err)
		require.Equal(t, true, accessKey.Enabled)

		err = deleteManagedObjectStorage(ctx, svc, storage.UUID)
		require.NoError(t, err)
	})
}

func TestDeleteManagedObjectStorageUserAccessKey(t *testing.T) {
	record(t, "deletemanagedobjectstorageuseraccesskey", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service) {
		storage, err := createManagedObjectStorage(ctx, svc)
		require.NoError(t, err)

		_, err = svc.WaitForManagedObjectStorageUserOperationalState(context.Background(), &request.WaitForManagedObjectStorageUserOperationalStateRequest{
			ServiceUUID:  storage.UUID,
			Timeout:      time.Second * 90,
			Username:     storage.Users[0].Username,
			DesiredState: upcloud.ManagedObjectStorageUserOperationalStateReady,
		})
		require.NoError(t, err)

		accessKey, err := svc.CreateManagedObjectStorageUserAccessKey(ctx, &request.CreateManagedObjectStorageUserAccessKeyRequest{
			Username:    storage.Users[0].Username,
			ServiceUUID: storage.UUID,
			Name:        "example-access-key",
			Enabled:     false,
		})
		require.NoError(t, err)

		err = svc.DeleteManagedObjectStorageUserAccessKey(ctx, &request.DeleteManagedObjectStorageUserAccessKeyRequest{
			ServiceUUID: storage.UUID,
			Username:    storage.Users[0].Username,
			Name:        accessKey.Name,
		})
		require.NoError(t, err)

		err = deleteManagedObjectStorage(ctx, svc, storage.UUID)
		require.NoError(t, err)
	})
}

func createManagedObjectStorage(ctx context.Context, svc *Service) (*upcloud.ManagedObjectStorage, error) {
	regions, err := svc.GetManagedObjectStorageRegions(ctx, &request.GetManagedObjectStorageRegionsRequest{})
	if err != nil {
		return nil, err
	}

	account, err := createSubAccount(ctx, svc, "managedobjectstorage")
	if err != nil {
		return nil, err
	}

	return svc.CreateManagedObjectStorage(ctx, &request.CreateManagedObjectStorageRequest{
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
		Users:  []request.ManagedObjectStorageUser{{Username: account.Username}},
	})
}

func createSubAccount(ctx context.Context, svc *Service, username string) (*upcloud.AccountDetails, error) {
	return svc.CreateSubaccount(ctx, &request.CreateSubaccountRequest{Subaccount: request.CreateSubaccount{
		Username:      username,
		Email:         fmt.Sprintf("%s@example.com", username),
		Phone:         "+358.123456789",
		Currency:      "EUR",
		Language:      "en",
		Timezone:      "Europe/Helsinki",
		AllowAPI:      upcloud.False,
		AllowGUI:      upcloud.False,
		TagAccess:     upcloud.AccountTagAccess{},
		ServerAccess:  upcloud.AccountServerAccess{},
		StorageAccess: upcloud.AccountStorageAccess{},
	}})
}

func deleteManagedObjectStorage(ctx context.Context, svc *Service, uuid string) error {
	users, err := svc.GetManagedObjectStorageUsers(ctx, &request.GetManagedObjectStorageUsersRequest{ServiceUUID: uuid})
	if err != nil {
		return err
	}

	err = svc.DeleteManagedObjectStorage(ctx, &request.DeleteManagedObjectStorageRequest{UUID: uuid})
	if err != nil {
		return err
	}

	err = svc.WaitForManagedObjectStorageDeletion(context.Background(), &request.WaitForManagedObjectStorageDeletionRequest{
		UUID:    uuid,
		Timeout: time.Second * 90,
	})
	if err != nil {
		return err
	}

	for _, user := range users {
		err = svc.DeleteSubaccount(ctx, &request.DeleteSubaccountRequest{Username: user.Username})
		if err != nil {
			return err
		}
	}

	return nil
}
