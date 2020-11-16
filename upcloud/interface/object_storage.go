package interfaces

import (
	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
)

type ObjectStorage interface {
	GetObjectStorages() (*upcloud.ObjectStorages, error)
	GetObjectStorageDetails(r *request.GetObjectStorageDetailsRequest) (*upcloud.ObjectStorageDetails, error)
	CreateObjectStorage(r *request.CreateObjectStorageRequest) (*upcloud.ObjectStorageDetails, error)
	ModifyObjectStorage(r *request.ModifyObjectStorageRequest) (*upcloud.ObjectStorageDetails, error)
	DeleteObjectStorage(r *request.DeleteObjectStorageRequest) error
}
