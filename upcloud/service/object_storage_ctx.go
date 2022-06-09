package service

import (
	"context"

	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud/request"
)

type ObjectStorageContext interface {
	GetObjectStorages(ctx context.Context) (*upcloud.ObjectStorages, error)
	GetObjectStorageDetails(ctx context.Context, r *request.GetObjectStorageDetailsRequest) (*upcloud.ObjectStorageDetails, error)
	CreateObjectStorage(ctx context.Context, r *request.CreateObjectStorageRequest) (*upcloud.ObjectStorageDetails, error)
	ModifyObjectStorage(ctx context.Context, r *request.ModifyObjectStorageRequest) (*upcloud.ObjectStorageDetails, error)
	DeleteObjectStorage(ctx context.Context, r *request.DeleteObjectStorageRequest) error
}

// GetObjectStorages returns the available objects storages
func (s *ServiceContext) GetObjectStorages(ctx context.Context) (*upcloud.ObjectStorages, error) {
	objectStorages := upcloud.ObjectStorages{}
	return &objectStorages, s.get(ctx, "/object-storage", &objectStorages)
}

// GetObjectStorageDetails returns extended details about the specified Object Storage
func (s *ServiceContext) GetObjectStorageDetails(ctx context.Context, r *request.GetObjectStorageDetailsRequest) (*upcloud.ObjectStorageDetails, error) {
	objectStorageDetails := upcloud.ObjectStorageDetails{}
	return &objectStorageDetails, s.get(ctx, r.RequestURL(), &objectStorageDetails)
}

// CreateObjectStorage creates a Object Storage and return the Object Storage details for the newly created device
func (s *ServiceContext) CreateObjectStorage(ctx context.Context, r *request.CreateObjectStorageRequest) (*upcloud.ObjectStorageDetails, error) {
	objectStorageDetails := upcloud.ObjectStorageDetails{}
	return &objectStorageDetails, s.create(ctx, r, &objectStorageDetails)
}

// ModifyObjectStorage modifies the configuration of an existing Object Storage
func (s *ServiceContext) ModifyObjectStorage(ctx context.Context, r *request.ModifyObjectStorageRequest) (*upcloud.ObjectStorageDetails, error) {
	objectStorageDetails := upcloud.ObjectStorageDetails{}
	return &objectStorageDetails, s.modify(ctx, r, &objectStorageDetails)
}

// DeleteObjectStorage deletes the specific Object Storage
func (s *ServiceContext) DeleteObjectStorage(ctx context.Context, r *request.DeleteObjectStorageRequest) error {
	return s.delete(ctx, r)
}
