package service

import (
	"context"
	"encoding/json"
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
	var result []upcloud.FileStorage
	response, err := s.client.Get(ctx, r.RequestURL())
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(response, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// CreateFileStorage creates a new file storage. (EXPERIMENTAL)
func (s *Service) CreateFileStorage(ctx context.Context, r *request.CreateFileStorageRequest) (*upcloud.FileStorage, error) {
	var result upcloud.FileStorage
	payload, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}
	response, err := s.client.Post(ctx, r.RequestURL(), payload)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(response, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetFileStorage retrieves details of a file storage. (EXPERIMENTAL)
func (s *Service) GetFileStorage(ctx context.Context, r *request.GetFileStorageRequest) (*upcloud.FileStorage, error) {
	var result upcloud.FileStorage
	response, err := s.client.Get(ctx, r.RequestURL())
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(response, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ReplaceFileStorage replaces an existing file storage. (EXPERIMENTAL)
func (s *Service) ReplaceFileStorage(ctx context.Context, r *request.ReplaceFileStorageRequest) (*upcloud.FileStorage, error) {
	var result upcloud.FileStorage
	payload, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}
	response, err := s.client.Put(ctx, r.RequestURL(), payload)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(response, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ModifyFileStorage modifies properties of an existing file storage. (EXPERIMENTAL)
func (s *Service) ModifyFileStorage(ctx context.Context, r *request.ModifyFileStorageRequest) (*upcloud.FileStorage, error) {
	var result upcloud.FileStorage
	payload, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}
	response, err := s.client.Patch(ctx, r.RequestURL(), payload)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(response, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// DeleteFileStorage deletes a file storage. (EXPERIMENTAL)
func (s *Service) DeleteFileStorage(ctx context.Context, r *request.DeleteFileStorageRequest) error {
	_, err := s.client.Delete(ctx, r.RequestURL())
	return err
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
	var result []upcloud.FileStorageNetwork
	response, err := s.client.Get(ctx, r.RequestURL())
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(response, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// CreateFileStorageNetwork creates a new file storage network. (EXPERIMENTAL)
func (s *Service) CreateFileStorageNetwork(ctx context.Context, r *request.CreateFileStorageNetworkRequest) (*upcloud.FileStorageNetwork, error) {
	var result upcloud.FileStorageNetwork
	payload, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}
	response, err := s.client.Post(ctx, r.RequestURL(), payload)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(response, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetFileStorageNetwork retrieves details of a file storage network. (EXPERIMENTAL)
func (s *Service) GetFileStorageNetwork(ctx context.Context, r *request.GetFileStorageNetworkRequest) (*upcloud.FileStorageNetwork, error) {
	var result upcloud.FileStorageNetwork
	response, err := s.client.Get(ctx, r.RequestURL())
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(response, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ModifyFileStorageNetwork modifies properties of an existing file storage network. (EXPERIMENTAL)
func (s *Service) ModifyFileStorageNetwork(ctx context.Context, r *request.ModifyFileStorageNetworkRequest) (*upcloud.FileStorageNetwork, error) {
	var result upcloud.FileStorageNetwork
	payload, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}
	response, err := s.client.Patch(ctx, r.RequestURL(), payload)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(response, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// DeleteFileStorageNetwork deletes a file storage network. (EXPERIMENTAL)
func (s *Service) DeleteFileStorageNetwork(ctx context.Context, r *request.DeleteFileStorageNetworkRequest) error {
	_, err := s.client.Delete(ctx, r.RequestURL())
	return err
}

// GetFileStorageShares retrieves a list of file storage shares. (EXPERIMENTAL)
func (s *Service) GetFileStorageShares(ctx context.Context, r *request.GetFileStorageSharesRequest) ([]upcloud.FileStorageShare, error) {
	var result []upcloud.FileStorageShare
	response, err := s.client.Get(ctx, r.RequestURL())
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(response, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// CreateFileStorageShare creates a new file storage share. (EXPERIMENTAL)
func (s *Service) CreateFileStorageShare(ctx context.Context, r *request.CreateFileStorageShareRequest) (*upcloud.FileStorageShare, error) {
	var result upcloud.FileStorageShare
	payload, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}
	response, err := s.client.Post(ctx, r.RequestURL(), payload)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(response, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetFileStorageShare retrieves details of a file storage share. (EXPERIMENTAL)
func (s *Service) GetFileStorageShare(ctx context.Context, r *request.GetFileStorageShareRequest) (*upcloud.FileStorageShare, error) {
	var result upcloud.FileStorageShare
	response, err := s.client.Get(ctx, r.RequestURL())
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(response, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ModifyFileStorageShare modifies properties of an existing file storage share. (EXPERIMENTAL)
func (s *Service) ModifyFileStorageShare(ctx context.Context, r *request.ModifyFileStorageShareRequest) (*upcloud.FileStorageShare, error) {
	var result upcloud.FileStorageShare
	payload, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}
	response, err := s.client.Patch(ctx, r.RequestURL(), payload)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(response, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// DeleteFileStorageShare deletes a file storage share. (EXPERIMENTAL)
func (s *Service) DeleteFileStorageShare(ctx context.Context, r *request.DeleteFileStorageShareRequest) error {
	_, err := s.client.Delete(ctx, r.RequestURL())
	return err
}

// GetFileStorageLabels retrieves a list of file storage labels. (EXPERIMENTAL)
func (s *Service) GetFileStorageLabels(ctx context.Context, r *request.GetFileStorageLabelsRequest) ([]upcloud.Label, error) {
	var result []upcloud.Label
	response, err := s.client.Get(ctx, r.RequestURL())
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(response, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// CreateFileStorageLabel creates a new file storage label. (EXPERIMENTAL)
func (s *Service) CreateFileStorageLabel(ctx context.Context, r *request.CreateFileStorageLabelRequest) (*upcloud.Label, error) {
	var result upcloud.Label
	payload, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}
	response, err := s.client.Post(ctx, r.RequestURL(), payload)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(response, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *Service) GetFileStorageLabel(ctx context.Context, r *request.GetFileStorageLabelRequest) (*upcloud.Label, error) {
	var result upcloud.Label
	response, err := s.client.Get(ctx, r.RequestURL())
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(response, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ModifyFileStorageLabel modifies properties of an existing file storage label. (EXPERIMENTAL)
func (s *Service) ModifyFileStorageLabel(ctx context.Context, r *request.ModifyFileStorageLabelRequest) (*upcloud.Label, error) {
	var result upcloud.Label
	payload, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}
	response, err := s.client.Patch(ctx, r.RequestURL(), payload)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(response, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// DeleteFileStorageLabel deletes a file storage label. (EXPERIMENTAL)
func (s *Service) DeleteFileStorageLabel(ctx context.Context, r *request.DeleteFileStorageLabelRequest) error {
	_, err := s.client.Delete(ctx, r.RequestURL())
	return err
}
