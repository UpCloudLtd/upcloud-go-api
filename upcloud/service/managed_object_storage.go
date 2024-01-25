package service

import (
	"context"
	"errors"
	"net/http"

	"github.com/UpCloudLtd/upcloud-go-api/v6/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v6/upcloud/request"
)

type ManagedObjectStorage interface {
	GetManagedObjectStorageRegions(ctx context.Context, r *request.GetManagedObjectStorageRegionsRequest) ([]upcloud.ManagedObjectStorageRegion, error)
	GetManagedObjectStorageRegion(ctx context.Context, r *request.GetManagedObjectStorageRegionRequest) (*upcloud.ManagedObjectStorageRegion, error)
	CreateManagedObjectStorage(ctx context.Context, r *request.CreateManagedObjectStorageRequest) (*upcloud.ManagedObjectStorage, error)
	GetManagedObjectStorages(ctx context.Context, r *request.GetManagedObjectStoragesRequest) ([]upcloud.ManagedObjectStorage, error)
	GetManagedObjectStorage(ctx context.Context, r *request.GetManagedObjectStorageRequest) (*upcloud.ManagedObjectStorage, error)
	ReplaceManagedObjectStorage(ctx context.Context, r *request.ReplaceManagedObjectStorageRequest) (*upcloud.ManagedObjectStorage, error)
	ModifyManagedObjectStorage(ctx context.Context, r *request.ModifyManagedObjectStorageRequest) (*upcloud.ManagedObjectStorage, error)
	DeleteManagedObjectStorage(ctx context.Context, r *request.DeleteManagedObjectStorageRequest) error
	GetManagedObjectStorageMetrics(ctx context.Context, r *request.GetManagedObjectStorageMetricsRequest) (*upcloud.ManagedObjectStorageMetrics, error)
	GetManagedObjectStorageBucketMetrics(ctx context.Context, r *request.GetManagedObjectStorageBucketMetricsRequest) ([]upcloud.ManagedObjectStorageBucketMetrics, error)
	CreateManagedObjectStorageNetwork(ctx context.Context, r *request.CreateManagedObjectStorageNetworkRequest) (*upcloud.ManagedObjectStorageNetwork, error)
	GetManagedObjectStorageNetworks(ctx context.Context, r *request.GetManagedObjectStorageNetworksRequest) ([]upcloud.ManagedObjectStorageNetwork, error)
	GetManagedObjectStorageNetwork(ctx context.Context, r *request.GetManagedObjectStorageNetworkRequest) (*upcloud.ManagedObjectStorageNetwork, error)
	DeleteManagedObjectStorageNetwork(ctx context.Context, r *request.DeleteManagedObjectStorageNetworkRequest) error
	CreateManagedObjectStorageUser(ctx context.Context, r *request.CreateManagedObjectStorageUserRequest) (*upcloud.ManagedObjectStorageUser, error)
	GetManagedObjectStorageUsers(ctx context.Context, r *request.GetManagedObjectStorageUsersRequest) ([]upcloud.ManagedObjectStorageUser, error)
	GetManagedObjectStorageUser(ctx context.Context, r *request.GetManagedObjectStorageUserRequest) (*upcloud.ManagedObjectStorageUser, error)
	DeleteManagedObjectStorageUser(ctx context.Context, r *request.DeleteManagedObjectStorageUserRequest) error
	CreateManagedObjectStorageUserAccessKey(ctx context.Context, r *request.CreateManagedObjectStorageUserAccessKeyRequest) (*upcloud.ManagedObjectStorageUserAccessKey, error)
	GetManagedObjectStorageUserAccessKeys(ctx context.Context, r *request.GetManagedObjectStorageUserAccessKeysRequest) ([]upcloud.ManagedObjectStorageUserAccessKey, error)
	GetManagedObjectStorageUserAccessKey(ctx context.Context, r *request.GetManagedObjectStorageUserAccessKeyRequest) (*upcloud.ManagedObjectStorageUserAccessKey, error)
	ModifyManagedObjectStorageUserAccessKey(ctx context.Context, r *request.ModifyManagedObjectStorageUserAccessKeyRequest) (*upcloud.ManagedObjectStorageUserAccessKey, error)
	DeleteManagedObjectStorageUserAccessKey(ctx context.Context, r *request.DeleteManagedObjectStorageUserAccessKeyRequest) error
	WaitForManagedObjectStorageOperationalState(ctx context.Context, r *request.WaitForManagedObjectStorageOperationalStateRequest) (*upcloud.ManagedObjectStorage, error)
	WaitForManagedObjectStorageDeletion(ctx context.Context, r *request.WaitForManagedObjectStorageDeletionRequest) error
}

func (s *Service) GetManagedObjectStorageRegions(ctx context.Context, r *request.GetManagedObjectStorageRegionsRequest) ([]upcloud.ManagedObjectStorageRegion, error) {
	regions := make([]upcloud.ManagedObjectStorageRegion, 0)
	return regions, s.get(ctx, r.RequestURL(), &regions)
}

func (s *Service) GetManagedObjectStorageRegion(ctx context.Context, r *request.GetManagedObjectStorageRegionRequest) (*upcloud.ManagedObjectStorageRegion, error) {
	region := upcloud.ManagedObjectStorageRegion{}
	return &region, s.get(ctx, r.RequestURL(), &region)
}

func (s *Service) CreateManagedObjectStorage(ctx context.Context, r *request.CreateManagedObjectStorageRequest) (*upcloud.ManagedObjectStorage, error) {
	storage := upcloud.ManagedObjectStorage{}
	return &storage, s.create(ctx, r, &storage)
}

func (s *Service) GetManagedObjectStorages(ctx context.Context, r *request.GetManagedObjectStoragesRequest) ([]upcloud.ManagedObjectStorage, error) {
	storages := make([]upcloud.ManagedObjectStorage, 0)
	return storages, s.get(ctx, r.RequestURL(), &storages)
}

func (s *Service) GetManagedObjectStorage(ctx context.Context, r *request.GetManagedObjectStorageRequest) (*upcloud.ManagedObjectStorage, error) {
	storage := upcloud.ManagedObjectStorage{}
	return &storage, s.get(ctx, r.RequestURL(), &storage)
}

func (s *Service) ReplaceManagedObjectStorage(ctx context.Context, r *request.ReplaceManagedObjectStorageRequest) (*upcloud.ManagedObjectStorage, error) {
	storage := upcloud.ManagedObjectStorage{}
	return &storage, s.replace(ctx, r, &storage)
}

func (s *Service) ModifyManagedObjectStorage(ctx context.Context, r *request.ModifyManagedObjectStorageRequest) (*upcloud.ManagedObjectStorage, error) {
	storage := upcloud.ManagedObjectStorage{}
	return &storage, s.modify(ctx, r, &storage)
}

func (s *Service) DeleteManagedObjectStorage(ctx context.Context, r *request.DeleteManagedObjectStorageRequest) error {
	return s.delete(ctx, r)
}

func (s *Service) GetManagedObjectStorageMetrics(ctx context.Context, r *request.GetManagedObjectStorageMetricsRequest) (*upcloud.ManagedObjectStorageMetrics, error) {
	metrics := upcloud.ManagedObjectStorageMetrics{}
	return &metrics, s.get(ctx, r.RequestURL(), &metrics)
}

func (s *Service) GetManagedObjectStorageBucketMetrics(ctx context.Context, r *request.GetManagedObjectStorageBucketMetricsRequest) ([]upcloud.ManagedObjectStorageBucketMetrics, error) {
	bucketMetrics := make([]upcloud.ManagedObjectStorageBucketMetrics, 0)
	return bucketMetrics, s.get(ctx, r.RequestURL(), &bucketMetrics)
}

func (s *Service) CreateManagedObjectStorageNetwork(ctx context.Context, r *request.CreateManagedObjectStorageNetworkRequest) (*upcloud.ManagedObjectStorageNetwork, error) {
	network := upcloud.ManagedObjectStorageNetwork{}
	return &network, s.create(ctx, r, &network)
}

func (s *Service) GetManagedObjectStorageNetworks(ctx context.Context, r *request.GetManagedObjectStorageNetworksRequest) ([]upcloud.ManagedObjectStorageNetwork, error) {
	networks := make([]upcloud.ManagedObjectStorageNetwork, 0)
	return networks, s.get(ctx, r.RequestURL(), &networks)
}

func (s *Service) GetManagedObjectStorageNetwork(ctx context.Context, r *request.GetManagedObjectStorageNetworkRequest) (*upcloud.ManagedObjectStorageNetwork, error) {
	network := upcloud.ManagedObjectStorageNetwork{}
	return &network, s.get(ctx, r.RequestURL(), &network)
}

func (s *Service) DeleteManagedObjectStorageNetwork(ctx context.Context, r *request.DeleteManagedObjectStorageNetworkRequest) error {
	return s.delete(ctx, r)
}

func (s *Service) CreateManagedObjectStorageUser(ctx context.Context, r *request.CreateManagedObjectStorageUserRequest) (*upcloud.ManagedObjectStorageUser, error) {
	user := upcloud.ManagedObjectStorageUser{}
	return &user, s.create(ctx, r, &user)
}

func (s *Service) GetManagedObjectStorageUsers(ctx context.Context, r *request.GetManagedObjectStorageUsersRequest) ([]upcloud.ManagedObjectStorageUser, error) {
	users := make([]upcloud.ManagedObjectStorageUser, 0)
	return users, s.get(ctx, r.RequestURL(), &users)
}

func (s *Service) GetManagedObjectStorageUser(ctx context.Context, r *request.GetManagedObjectStorageUserRequest) (*upcloud.ManagedObjectStorageUser, error) {
	user := upcloud.ManagedObjectStorageUser{}
	return &user, s.get(ctx, r.RequestURL(), &user)
}

func (s *Service) DeleteManagedObjectStorageUser(ctx context.Context, r *request.DeleteManagedObjectStorageUserRequest) error {
	return s.delete(ctx, r)
}

func (s *Service) CreateManagedObjectStorageUserAccessKey(ctx context.Context, r *request.CreateManagedObjectStorageUserAccessKeyRequest) (*upcloud.ManagedObjectStorageUserAccessKey, error) {
	accessKey := upcloud.ManagedObjectStorageUserAccessKey{}
	return &accessKey, s.create(ctx, r, &accessKey)
}

func (s *Service) GetManagedObjectStorageUserAccessKeys(ctx context.Context, r *request.GetManagedObjectStorageUserAccessKeysRequest) ([]upcloud.ManagedObjectStorageUserAccessKey, error) {
	accessKeys := make([]upcloud.ManagedObjectStorageUserAccessKey, 0)
	return accessKeys, s.get(ctx, r.RequestURL(), &accessKeys)
}

func (s *Service) GetManagedObjectStorageUserAccessKey(ctx context.Context, r *request.GetManagedObjectStorageUserAccessKeyRequest) (*upcloud.ManagedObjectStorageUserAccessKey, error) {
	accessKey := upcloud.ManagedObjectStorageUserAccessKey{}
	return &accessKey, s.get(ctx, r.RequestURL(), &accessKey)
}

func (s *Service) ModifyManagedObjectStorageUserAccessKey(ctx context.Context, r *request.ModifyManagedObjectStorageUserAccessKeyRequest) (*upcloud.ManagedObjectStorageUserAccessKey, error) {
	accessKey := upcloud.ManagedObjectStorageUserAccessKey{}
	return &accessKey, s.modify(ctx, r, &accessKey)
}

func (s *Service) DeleteManagedObjectStorageUserAccessKey(ctx context.Context, r *request.DeleteManagedObjectStorageUserAccessKeyRequest) error {
	return s.delete(ctx, r)
}

// WaitForManagedObjectStorageOperationalState blocks execution until the specified Managed Object Storage service
// has entered the specified state. If the state changes favorably, service details is returned. The method will give up
// after the specified timeout
func (s *Service) WaitForManagedObjectStorageOperationalState(ctx context.Context, r *request.WaitForManagedObjectStorageOperationalStateRequest) (*upcloud.ManagedObjectStorage, error) {
	return retry(ctx, func(_ int, c context.Context) (*upcloud.ManagedObjectStorage, error) {
		details, err := s.GetManagedObjectStorage(c, &request.GetManagedObjectStorageRequest{
			UUID: r.UUID,
		})
		if err != nil {
			return nil, err
		}

		if details.OperationalState == r.DesiredState {
			return details, nil
		}
		return nil, nil
	}, nil)
}

// WaitForManagedObjectStorageUserOperationalState blocks execution until the specified Managed Object Storage service
// user has entered the specified state. If the state changes favorably, user details is returned. The method will give up
// after the specified timeout
func (s *Service) WaitForManagedObjectStorageUserOperationalState(ctx context.Context, r *request.WaitForManagedObjectStorageUserOperationalStateRequest) (*upcloud.ManagedObjectStorageUser, error) {
	return retry(ctx, func(_ int, c context.Context) (*upcloud.ManagedObjectStorageUser, error) {
		details, err := s.GetManagedObjectStorageUser(c, &request.GetManagedObjectStorageUserRequest{
			ServiceUUID: r.ServiceUUID,
			Username:    r.Username,
		})
		if err != nil {
			return nil, err
		}

		if details.OperationalState == r.DesiredState {
			return details, nil
		}
		return nil, nil
	}, nil)
}

// WaitForManagedObjectStorageDeletion blocks execution until the specified Managed Object Storage service
// has been deleted. The method will give upafter the specified timeout
func (s *Service) WaitForManagedObjectStorageDeletion(ctx context.Context, r *request.WaitForManagedObjectStorageDeletionRequest) error {
	_, err := retry(ctx, func(_ int, c context.Context) (*upcloud.ManagedObjectStorage, error) {
		details, err := s.GetManagedObjectStorage(c, &request.GetManagedObjectStorageRequest{
			UUID: r.UUID,
		})
		if err != nil {
			var ucErr *upcloud.Problem
			if errors.As(err, &ucErr) && ucErr.Status == http.StatusNotFound {
				return nil, nil
			}

			return nil, err
		}

		return details, err
	}, &retryConfig{inverse: true})
	return err
}
