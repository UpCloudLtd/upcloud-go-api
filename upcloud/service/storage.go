package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
)

type Storage interface {
	GetStorages(ctx context.Context, r *request.GetStoragesRequest) (*upcloud.Storages, error)
	GetStorageDetails(ctx context.Context, r *request.GetStorageDetailsRequest) (*upcloud.StorageDetails, error)
	CreateStorage(ctx context.Context, r *request.CreateStorageRequest) (*upcloud.StorageDetails, error)
	ModifyStorage(ctx context.Context, r *request.ModifyStorageRequest) (*upcloud.StorageDetails, error)
	AttachStorage(ctx context.Context, r *request.AttachStorageRequest) (*upcloud.ServerDetails, error)
	DetachStorage(ctx context.Context, r *request.DetachStorageRequest) (*upcloud.ServerDetails, error)
	CloneStorage(ctx context.Context, r *request.CloneStorageRequest) (*upcloud.StorageDetails, error)
	TemplatizeStorage(ctx context.Context, r *request.TemplatizeStorageRequest) (*upcloud.StorageDetails, error)
	WaitForStorageState(ctx context.Context, r *request.WaitForStorageStateRequest) (*upcloud.StorageDetails, error)
	LoadCDROM(ctx context.Context, r *request.LoadCDROMRequest) (*upcloud.ServerDetails, error)
	EjectCDROM(ctx context.Context, r *request.EjectCDROMRequest) (*upcloud.ServerDetails, error)
	CreateBackup(ctx context.Context, r *request.CreateBackupRequest) (*upcloud.StorageDetails, error)
	RestoreBackup(ctx context.Context, r *request.RestoreBackupRequest) error
	CreateStorageImport(ctx context.Context, r *request.CreateStorageImportRequest) (*upcloud.StorageImportDetails, error)
	GetStorageImportDetails(ctx context.Context, r *request.GetStorageImportDetailsRequest) (*upcloud.StorageImportDetails, error)
	WaitForStorageImportCompletion(ctx context.Context, r *request.WaitForStorageImportCompletionRequest) (*upcloud.StorageImportDetails, error)
	DeleteStorage(ctx context.Context, r *request.DeleteStorageRequest) error
	ResizeStorageFilesystem(ctx context.Context, r *request.ResizeStorageFilesystemRequest) (*upcloud.ResizeStorageFilesystemBackup, error)
}

// GetStorages returns all available storages
func (s *Service) GetStorages(ctx context.Context, r *request.GetStoragesRequest) (*upcloud.Storages, error) {
	storages := upcloud.Storages{}
	return &storages, s.get(ctx, r.RequestURL(), &storages)
}

// GetStorageDetails returns extended details about the specified piece of storage
func (s *Service) GetStorageDetails(ctx context.Context, r *request.GetStorageDetailsRequest) (*upcloud.StorageDetails, error) {
	storageDetails := upcloud.StorageDetails{}
	return &storageDetails, s.get(ctx, r.RequestURL(), &storageDetails)
}

// CreateStorage creates the specified storage
func (s *Service) CreateStorage(ctx context.Context, r *request.CreateStorageRequest) (*upcloud.StorageDetails, error) {
	storageDetails := upcloud.StorageDetails{}
	return &storageDetails, s.create(ctx, r, &storageDetails)
}

// ModifyStorage modifies the specified storage device
func (s *Service) ModifyStorage(ctx context.Context, r *request.ModifyStorageRequest) (*upcloud.StorageDetails, error) {
	storageDetails := upcloud.StorageDetails{}
	return &storageDetails, s.replace(ctx, r, &storageDetails)
}

// AttachStorage attaches the specified storage to the specified server
func (s *Service) AttachStorage(ctx context.Context, r *request.AttachStorageRequest) (*upcloud.ServerDetails, error) {
	serverDetails := upcloud.ServerDetails{}
	return &serverDetails, s.create(ctx, r, &serverDetails)
}

// DetachStorage detaches the specified storage from the specified server
func (s *Service) DetachStorage(ctx context.Context, r *request.DetachStorageRequest) (*upcloud.ServerDetails, error) {
	serverDetails := upcloud.ServerDetails{}
	return &serverDetails, s.create(ctx, r, &serverDetails)
}

// DeleteStorage deletes the specified storage device
func (s *Service) DeleteStorage(ctx context.Context, r *request.DeleteStorageRequest) error {
	return s.delete(ctx, r)
}

// CloneStorage detaches the specified storage from the specified server
func (s *Service) CloneStorage(ctx context.Context, r *request.CloneStorageRequest) (*upcloud.StorageDetails, error) {
	storageDetails := upcloud.StorageDetails{}
	return &storageDetails, s.create(ctx, r, &storageDetails)
}

// TemplatizeStorage detaches the specified storage from the specified server
func (s *Service) TemplatizeStorage(ctx context.Context, r *request.TemplatizeStorageRequest) (*upcloud.StorageDetails, error) {
	storageDetails := upcloud.StorageDetails{}
	return &storageDetails, s.create(ctx, r, &storageDetails)
}

// WaitForStorageState blocks execution until the specified storage device has entered the specified state. If the
// state changes favorably, the new storage details is returned. The method will give up after the specified timeout
func (s *Service) WaitForStorageState(ctx context.Context, r *request.WaitForStorageStateRequest) (*upcloud.StorageDetails, error) {
	return retry(ctx, func(i int, c context.Context) (*upcloud.StorageDetails, error) {
		details, err := s.GetStorageDetails(c, &request.GetStorageDetailsRequest{
			UUID: r.UUID,
		})
		if err != nil {
			return nil, err
		}

		if details.State == r.DesiredState {
			return details, nil
		}

		return nil, nil
	}, nil)
}

// LoadCDROM loads a storage as a CD-ROM in the CD-ROM device of a server
func (s *Service) LoadCDROM(ctx context.Context, r *request.LoadCDROMRequest) (*upcloud.ServerDetails, error) {
	serverDetails := upcloud.ServerDetails{}
	return &serverDetails, s.create(ctx, r, &serverDetails)
}

// EjectCDROM ejects the storage from the CD-ROM device of a server
func (s *Service) EjectCDROM(ctx context.Context, r *request.EjectCDROMRequest) (*upcloud.ServerDetails, error) {
	serverDetails := upcloud.ServerDetails{}
	return &serverDetails, s.create(ctx, r, &serverDetails)
}

// CreateBackup creates a backup of the specified storage
func (s *Service) CreateBackup(ctx context.Context, r *request.CreateBackupRequest) (*upcloud.StorageDetails, error) {
	storageDetails := upcloud.StorageDetails{}
	return &storageDetails, s.create(ctx, r, &storageDetails)
}

// RestoreBackup creates a backup of the specified storage
func (s *Service) RestoreBackup(ctx context.Context, r *request.RestoreBackupRequest) error {
	return s.create(ctx, r, nil)
}

// CreateStorageImport begins the process of importing an image onto a storage device. A `upcloud.StorageImportSourceHTTPImport` source
// will import from an HTTP source. `upcloud.StorageImportSourceDirectUpload` will directly upload the file specified in `SourceLocation`.
func (s *Service) CreateStorageImport(ctx context.Context, r *request.CreateStorageImportRequest) (*upcloud.StorageImportDetails, error) {
	if r.Source == request.StorageImportSourceDirectUpload {
		switch r.SourceLocation.(type) {
		case string, io.Reader:
			return s.directStorageImport(ctx, r)
		case nil:
			return nil, errors.New("SourceLocation must be specified")
		default:
			return nil, fmt.Errorf("unsupported storage source location type %T", r.SourceLocation)
		}
	}

	if _, isString := r.SourceLocation.(string); !isString {
		return nil, fmt.Errorf("unsupported storage source location type %T", r.Source)
	}
	return s.doCreateStorageImport(ctx, r)
}

// doCreateStorageImport will POST the CreateStorageImport request and handle the error and normal response.
func (s *Service) doCreateStorageImport(ctx context.Context, r *request.CreateStorageImportRequest) (*upcloud.StorageImportDetails, error) {
	storageImport := upcloud.StorageImportDetails{}
	return &storageImport, s.create(ctx, r, &storageImport)
}

// directStorageImport handles the direct upload logic including getting the upload URL and PUT the file data
// to that endpoint.
func (s *Service) directStorageImport(ctx context.Context, r *request.CreateStorageImportRequest) (*upcloud.StorageImportDetails, error) {
	var bodyReader io.Reader

	switch v := r.SourceLocation.(type) {
	case string:
		if v == "" {
			return nil, errors.New("SourceLocation must be specified")
		}
		f, err := os.Open(v)
		if err != nil {
			return nil, fmt.Errorf("unable to open SourceLocation: %w", err)
		}
		bodyReader = f
	case io.Reader:
		bodyReader = v
	default:
		return nil, fmt.Errorf("unsupported source location type %T", r.SourceLocation)
	}

	r.SourceLocation = ""
	storageImport, err := s.doCreateStorageImport(ctx, r)
	if err != nil {
		return nil, err
	}

	if storageImport.DirectUploadURL == "" {
		return nil, errors.New("no DirectUploadURL found in response")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, storageImport.DirectUploadURL, bodyReader)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", r.ContentType)
	if _, err := s.client.Do(req); err != nil {
		return nil, err
	}

	storageImport, err = s.GetStorageImportDetails(ctx, &request.GetStorageImportDetailsRequest{
		UUID: r.StorageUUID,
	})
	if err != nil {
		return nil, err
	}

	return storageImport, nil
}

// GetStorageImportDetails gets updated details about the specified storage import.
func (s *Service) GetStorageImportDetails(ctx context.Context, r *request.GetStorageImportDetailsRequest) (*upcloud.StorageImportDetails, error) {
	storageDetails := upcloud.StorageImportDetails{}
	return &storageDetails, s.get(ctx, r.RequestURL(), &storageDetails)
}

// WaitForStorageImportCompletion waits for the importing storage to complete.
func (s *Service) WaitForStorageImportCompletion(ctx context.Context, r *request.WaitForStorageImportCompletionRequest) (*upcloud.StorageImportDetails, error) {
	return retry(ctx, func(i int, c context.Context) (*upcloud.StorageImportDetails, error) {
		details, err := s.GetStorageImportDetails(c, &request.GetStorageImportDetailsRequest{
			UUID: r.StorageUUID,
		})
		if err != nil {
			return nil, err
		}

		switch details.State {
		case upcloud.StorageImportStateCompleted:
			return details, nil
		case upcloud.StorageImportStateCancelled,
			upcloud.StorageImportStateCancelling,
			upcloud.StorageImportStateFailed:
			if details.ErrorCode != "" || details.ErrorMessage != "" {
				return details, &upcloud.Problem{
					Type:  details.ErrorCode,
					Title: details.ErrorMessage,
				}
			}
			return details, &upcloud.Problem{
				Type:  details.State,
				Title: "Storage Import Failed",
			}
		default:
			return nil, nil
		}
	}, nil)
}

// ResizeStorageFilesystem resizes the last partition of a storage and the ext3/ext4/XFS/NTFS filesystem
// on that partition if the partition does not extend to the end of the storage yet.
//
// Before the resize is attempted, a backup is taken from the storage. If the resize
// succeeds, backup details are returned. It is advisable to keep the backup until
// you have ensured that everything works after the resize.
//
// If the resize fails, backup is used to restore the storage to the state where it
// was before the resize. After that the backup is deleted automatically.
func (s *Service) ResizeStorageFilesystem(ctx context.Context, r *request.ResizeStorageFilesystemRequest) (*upcloud.ResizeStorageFilesystemBackup, error) {
	resizeBackup := upcloud.ResizeStorageFilesystemBackup{}
	return &resizeBackup, s.create(ctx, r, &resizeBackup)
}
