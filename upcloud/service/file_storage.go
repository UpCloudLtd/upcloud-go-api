package service

import (
	"context"
	"encoding/json"

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

func (s *Service) DeleteFileStorage(ctx context.Context, r *request.DeleteFileStorageRequest) error {
	_, err := s.client.Delete(ctx, r.RequestURL())
	return err
}

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

func (s *Service) DeleteFileStorageNetwork(ctx context.Context, r *request.DeleteFileStorageNetworkRequest) error {
	_, err := s.client.Delete(ctx, r.RequestURL())
	return err
}

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

func (s *Service) DeleteFileStorageShare(ctx context.Context, r *request.DeleteFileStorageShareRequest) error {
	_, err := s.client.Delete(ctx, r.RequestURL())
	return err
}

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

func (s *Service) ModifyFileStorageLabel(ctx context.Context, r *request.ModifyFileStorageLabelRequest) (*upcloud.Label, error) {
	var result upcloud.Label
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

func (s *Service) DeleteFileStorageLabel(ctx context.Context, r *request.DeleteFileStorageLabelRequest) error {
	_, err := s.client.Delete(ctx, r.RequestURL())
	return err
}
