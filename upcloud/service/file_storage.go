package service

import (
	"context"
	"errors"
	"net/http"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
)

type FileStorage interface {
	GetFileStorages(ctx context.Context, r *request.GetFileStoragesRequest) ([]upcloud.FileStorage, error)
	CreateFileStorage(ctx context.Context, r *request.CreateFileStorageRequest) (*upcloud.FileStorage, error)
	GetFileStorage(ctx context.Context, r *request.GetFileStorageRequest) (*upcloud.FileStorage, error)
	ReplaceFileStorage(ctx context.Context, r *request.ReplaceFileStorageRequest) (*upcloud.FileStorage, error)
	ModifyFileStorage(ctx context.Context, r *request.ModifyFileStorageRequest) (*upcloud.FileStorage, error)
	DeleteFileStorage(ctx context.Context, r *request.DeleteFileStorageRequest) error
	GetFileStorageNetworks(ctx context.Context, r *request.GetFileStorageNetworksRequest) ([]upcloud.FileStorageNetwork, error)
	CreateFileStorageNetwork(ctx context.Context, r *request.CreateFileStorageNetworkRequest) (*upcloud.FileStorageNetwork, error)
	GetFileStorageNetwork(ctx context.Context, r *request.GetFileStorageNetworkRequest) (*upcloud.FileStorageNetwork, error)
	ModifyFileStorageNetwork(ctx context.Context, r *request.ModifyFileStorageNetworkRequest) (*upcloud.FileStorageNetwork, error)
	DeleteFileStorageNetwork(ctx context.Context, r *request.DeleteFileStorageNetworkRequest) error
	GetFileStorageShares(ctx context.Context, r *request.GetFileStorageSharesRequest) ([]upcloud.FileStorageShare, error)
	CreateFileStorageShare(ctx context.Context, r *request.CreateFileStorageShareRequest) (*upcloud.FileStorageShare, error)
	GetFileStorageShare(ctx context.Context, r *request.GetFileStorageShareRequest) (*upcloud.FileStorageShare, error)
	ModifyFileStorageShare(ctx context.Context, r *request.ModifyFileStorageShareRequest) (*upcloud.FileStorageShare, error)
	DeleteFileStorageShare(ctx context.Context, r *request.DeleteFileStorageShareRequest) error
	GetFileStorageLabels(ctx context.Context, r *request.GetFileStorageLabelsRequest) ([]upcloud.Label, error)
	CreateFileStorageLabel(ctx context.Context, r *request.CreateFileStorageLabelRequest) (*upcloud.Label, error)
	GetFileStorageLabel(ctx context.Context, r *request.GetFileStorageLabelRequest) (*upcloud.Label, error)
	ModifyFileStorageLabel(ctx context.Context, r *request.ModifyFileStorageLabelRequest) (*upcloud.Label, error)
	DeleteFileStorageLabel(ctx context.Context, r *request.DeleteFileStorageLabelRequest) error
}

// GetFileStorages retrieves a list of file storages. (EXPERIMENTAL)
func (s *Service) GetFileStorages(ctx context.Context, r *request.GetFileStoragesRequest) ([]upcloud.FileStorage, error) {
	fileStorages := make([]upcloud.FileStorage, 0)
	if r.Page != nil {
		return fileStorages, s.get(ctx, r.RequestURL(), &fileStorages)
	}

	// copy request value so that we are not altering original request
	req := *r

	// use default page size and get all available records
	req.Page = request.DefaultPage

	// loop until max result is reached or until response doesn't fill our page anymore
	for len(fileStorages) <= request.PageResultMaxSize {
		fss := make([]upcloud.FileStorage, 0)
		if err := s.get(ctx, req.RequestURL(), &fss); err != nil || len(fss) < 1 {
			return fileStorages, err
		}

		fileStorages = append(fileStorages, fss...)
		if len(fss) < req.Page.Size {
			return fileStorages, nil
		}

		req.Page = req.Page.Next()
	}

	return fileStorages, nil
}

// CreateFileStorage creates a new file storage. (EXPERIMENTAL)
func (s *Service) CreateFileStorage(ctx context.Context, r *request.CreateFileStorageRequest) (*upcloud.FileStorage, error) {
	fileStorage := upcloud.FileStorage{}
	return &fileStorage, s.create(ctx, r, &fileStorage)
}

// GetFileStorage retrieves details of a file storage. (EXPERIMENTAL)
func (s *Service) GetFileStorage(ctx context.Context, r *request.GetFileStorageRequest) (*upcloud.FileStorage, error) {
	fileStorage := upcloud.FileStorage{}
	return &fileStorage, s.get(ctx, r.RequestURL(), &fileStorage)
}

// ReplaceFileStorage replaces an existing file storage. (EXPERIMENTAL)
func (s *Service) ReplaceFileStorage(ctx context.Context, r *request.ReplaceFileStorageRequest) (*upcloud.FileStorage, error) {
	fileStorage := upcloud.FileStorage{}
	return &fileStorage, s.replace(ctx, r, &fileStorage)
}

// ModifyFileStorage modifies properties of an existing file storage. (EXPERIMENTAL)
func (s *Service) ModifyFileStorage(ctx context.Context, r *request.ModifyFileStorageRequest) (*upcloud.FileStorage, error) {
	fileStorage := upcloud.FileStorage{}
	return &fileStorage, s.modify(ctx, r, &fileStorage)
}

// DeleteFileStorage deletes a file storage. (EXPERIMENTAL)
func (s *Service) DeleteFileStorage(ctx context.Context, r *request.DeleteFileStorageRequest) error {
	return s.delete(ctx, r)
}

// WaitForFileStorageDeletion blocks execution until the specified Managed Object Storage service has been deleted.
func (s *Service) WaitForFileStorageDeletion(ctx context.Context, r *request.WaitForFileStorageDeletionRequest) error {
	_, err := retry(ctx, func(_ int, c context.Context) (*upcloud.FileStorage, error) {
		details, err := s.GetFileStorage(c, &request.GetFileStorageRequest{
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

// WaitForFileStorageOperationalState blocks execution until the specified file storage instance has entered the
// specified state. If the state changes favorably, the file storage details are returned. The method will give up
// after the specified timeout. (EXPERIMENTAL)
func (s *Service) WaitForFileStorageOperationalState(ctx context.Context, r *request.WaitForFileStorageOperationalStateRequest) (*upcloud.FileStorage, error) {
	return retry(ctx, func(_ int, c context.Context) (*upcloud.FileStorage, error) {
		details, err := s.GetFileStorage(c, &request.GetFileStorageRequest{
			UUID: r.UUID,
		})
		if err != nil {
			return nil, err
		}

		if details.OperationalState == string(r.DesiredState) {
			return details, nil
		}
		return nil, nil
	}, nil)
}

// GetFileStorageNetworks retrieves a list of file storage networks. (EXPERIMENTAL)
func (s *Service) GetFileStorageNetworks(ctx context.Context, r *request.GetFileStorageNetworksRequest) ([]upcloud.FileStorageNetwork, error) {
	fileStorageNetworks := make([]upcloud.FileStorageNetwork, 0)
	return fileStorageNetworks, s.get(ctx, r.RequestURL(), &fileStorageNetworks)
}

// CreateFileStorageNetwork creates a new file storage network. (EXPERIMENTAL)
func (s *Service) CreateFileStorageNetwork(ctx context.Context, r *request.CreateFileStorageNetworkRequest) (*upcloud.FileStorageNetwork, error) {
	fileStorageNetwork := upcloud.FileStorageNetwork{}
	return &fileStorageNetwork, s.create(ctx, r, &fileStorageNetwork)
}

// GetFileStorageNetwork retrieves details of a file storage network. (EXPERIMENTAL)
func (s *Service) GetFileStorageNetwork(ctx context.Context, r *request.GetFileStorageNetworkRequest) (*upcloud.FileStorageNetwork, error) {
	fileStorageNetwork := upcloud.FileStorageNetwork{}
	return &fileStorageNetwork, s.get(ctx, r.RequestURL(), &fileStorageNetwork)
}

// ModifyFileStorageNetwork modifies properties of an existing file storage network. (EXPERIMENTAL)
func (s *Service) ModifyFileStorageNetwork(ctx context.Context, r *request.ModifyFileStorageNetworkRequest) (*upcloud.FileStorageNetwork, error) {
	fileStorageNetwork := upcloud.FileStorageNetwork{}
	return &fileStorageNetwork, s.modify(ctx, r, &fileStorageNetwork)
}

// DeleteFileStorageNetwork deletes a file storage network. (EXPERIMENTAL)
func (s *Service) DeleteFileStorageNetwork(ctx context.Context, r *request.DeleteFileStorageNetworkRequest) error {
	return s.delete(ctx, r)
}

// GetFileStorageShares retrieves a list of file storage shares. (EXPERIMENTAL)
func (s *Service) GetFileStorageShares(ctx context.Context, r *request.GetFileStorageSharesRequest) ([]upcloud.FileStorageShare, error) {
	fileStorageShares := make([]upcloud.FileStorageShare, 0)
	return fileStorageShares, s.get(ctx, r.RequestURL(), &fileStorageShares)
}

// CreateFileStorageShare creates a new file storage share. (EXPERIMENTAL)
func (s *Service) CreateFileStorageShare(ctx context.Context, r *request.CreateFileStorageShareRequest) (*upcloud.FileStorageShare, error) {
	fileStorageShare := upcloud.FileStorageShare{}
	return &fileStorageShare, s.create(ctx, r, &fileStorageShare)
}

// GetFileStorageShare retrieves details of a file storage share. (EXPERIMENTAL)
func (s *Service) GetFileStorageShare(ctx context.Context, r *request.GetFileStorageShareRequest) (*upcloud.FileStorageShare, error) {
	fileStorageShare := upcloud.FileStorageShare{}
	return &fileStorageShare, s.get(ctx, r.RequestURL(), &fileStorageShare)
}

// ModifyFileStorageShare modifies properties of an existing file storage share. (EXPERIMENTAL)
func (s *Service) ModifyFileStorageShare(ctx context.Context, r *request.ModifyFileStorageShareRequest) (*upcloud.FileStorageShare, error) {
	fileStorageShare := upcloud.FileStorageShare{}
	return &fileStorageShare, s.modify(ctx, r, &fileStorageShare)
}

// DeleteFileStorageShare deletes a file storage share. (EXPERIMENTAL)
func (s *Service) DeleteFileStorageShare(ctx context.Context, r *request.DeleteFileStorageShareRequest) error {
	return s.delete(ctx, r)
}

// GetFileStorageLabels retrieves a list of file storage labels. (EXPERIMENTAL)
func (s *Service) GetFileStorageLabels(ctx context.Context, r *request.GetFileStorageLabelsRequest) ([]upcloud.Label, error) {
	labels := make([]upcloud.Label, 0)
	return labels, s.get(ctx, r.RequestURL(), &labels)
}

// CreateFileStorageLabel creates a new file storage label. (EXPERIMENTAL)
func (s *Service) CreateFileStorageLabel(ctx context.Context, r *request.CreateFileStorageLabelRequest) (*upcloud.Label, error) {
	label := upcloud.Label{}
	return &label, s.create(ctx, r, &label)
}

func (s *Service) GetFileStorageLabel(ctx context.Context, r *request.GetFileStorageLabelRequest) (*upcloud.Label, error) {
	label := upcloud.Label{}
	return &label, s.get(ctx, r.RequestURL(), &label)
}

// ModifyFileStorageLabel modifies properties of an existing file storage label. (EXPERIMENTAL)
func (s *Service) ModifyFileStorageLabel(ctx context.Context, r *request.ModifyFileStorageLabelRequest) (*upcloud.Label, error) {
	label := upcloud.Label{}
	return &label, s.modify(ctx, r, &label)
}

// DeleteFileStorageLabel deletes a file storage label. (EXPERIMENTAL)
func (s *Service) DeleteFileStorageLabel(ctx context.Context, r *request.DeleteFileStorageLabelRequest) error {
	return s.delete(ctx, r)
}
