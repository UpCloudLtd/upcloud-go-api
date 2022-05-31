package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud/request"
)

type StorageContext interface {
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
func (s *ServiceContext) GetStorages(ctx context.Context, r *request.GetStoragesRequest) (*upcloud.Storages, error) {
	storages := upcloud.Storages{}
	return &storages, s.get(ctx, r.RequestURL(), &storages)
}

// GetStorageDetails returns extended details about the specified piece of storage
func (s *ServiceContext) GetStorageDetails(ctx context.Context, r *request.GetStorageDetailsRequest) (*upcloud.StorageDetails, error) {
	storageDetails := upcloud.StorageDetails{}
	return &storageDetails, s.get(ctx, r.RequestURL(), &storageDetails)
}

// CreateStorage creates the specified storage
func (s *ServiceContext) CreateStorage(ctx context.Context, r *request.CreateStorageRequest) (*upcloud.StorageDetails, error) {
	storageDetails := upcloud.StorageDetails{}
	return &storageDetails, s.create(ctx, r, &storageDetails)
}

// ModifyStorage modifies the specified storage device
func (s *ServiceContext) ModifyStorage(ctx context.Context, r *request.ModifyStorageRequest) (*upcloud.StorageDetails, error) {
	storageDetails := upcloud.StorageDetails{}
	return &storageDetails, s.replace(ctx, r, &storageDetails)
}

// AttachStorage attaches the specified storage to the specified server
func (s *ServiceContext) AttachStorage(ctx context.Context, r *request.AttachStorageRequest) (*upcloud.ServerDetails, error) {
	serverDetails := upcloud.ServerDetails{}
	return &serverDetails, s.create(ctx, r, &serverDetails)
}

// DetachStorage detaches the specified storage from the specified server
func (s *ServiceContext) DetachStorage(ctx context.Context, r *request.DetachStorageRequest) (*upcloud.ServerDetails, error) {
	serverDetails := upcloud.ServerDetails{}
	return &serverDetails, s.create(ctx, r, &serverDetails)
}

// DeleteStorage deletes the specified storage device
func (s *ServiceContext) DeleteStorage(ctx context.Context, r *request.DeleteStorageRequest) error {
	return s.delete(ctx, r)
}

// CloneStorage detaches the specified storage from the specified server
func (s *ServiceContext) CloneStorage(ctx context.Context, r *request.CloneStorageRequest) (*upcloud.StorageDetails, error) {
	storageDetails := upcloud.StorageDetails{}
	return &storageDetails, s.create(ctx, r, &storageDetails)
}

// TemplatizeStorage detaches the specified storage from the specified server
func (s *ServiceContext) TemplatizeStorage(ctx context.Context, r *request.TemplatizeStorageRequest) (*upcloud.StorageDetails, error) {
	storageDetails := upcloud.StorageDetails{}
	return &storageDetails, s.create(ctx, r, &storageDetails)
}

// WaitForStorageState blocks execution until the specified storage device has entered the specified state. If the
// state changes favorably, the new storage details is returned. The method will give up after the specified timeout
func (s *ServiceContext) WaitForStorageState(ctx context.Context, r *request.WaitForStorageStateRequest) (*upcloud.StorageDetails, error) {
	attempts := 0
	sleepDuration := time.Second * 5

	for {
		attempts++

		storageDetails, err := s.GetStorageDetails(ctx, &request.GetStorageDetailsRequest{
			UUID: r.UUID,
		})

		if err != nil {
			return nil, err
		}

		if storageDetails.State == r.DesiredState {
			return storageDetails, nil
		}

		time.Sleep(sleepDuration)

		if time.Duration(attempts)*sleepDuration >= r.Timeout {
			return nil, fmt.Errorf("timeout reached while waiting for storage to enter state \"%s\"", r.DesiredState)
		}
	}
}

// LoadCDROM loads a storage as a CD-ROM in the CD-ROM device of a server
func (s *ServiceContext) LoadCDROM(ctx context.Context, r *request.LoadCDROMRequest) (*upcloud.ServerDetails, error) {
	serverDetails := upcloud.ServerDetails{}
	return &serverDetails, s.create(ctx, r, &serverDetails)
}

// EjectCDROM ejects the storage from the CD-ROM device of a server
func (s *ServiceContext) EjectCDROM(ctx context.Context, r *request.EjectCDROMRequest) (*upcloud.ServerDetails, error) {
	serverDetails := upcloud.ServerDetails{}
	return &serverDetails, s.create(ctx, r, &serverDetails)
}

// CreateBackup creates a backup of the specified storage
func (s *ServiceContext) CreateBackup(ctx context.Context, r *request.CreateBackupRequest) (*upcloud.StorageDetails, error) {
	storageDetails := upcloud.StorageDetails{}
	return &storageDetails, s.create(ctx, r, &storageDetails)
}

// RestoreBackup creates a backup of the specified storage
func (s *ServiceContext) RestoreBackup(ctx context.Context, r *request.RestoreBackupRequest) error {
	return s.create(ctx, r, nil)
}

// CreateStorageImport begins the process of importing an image onto a storage device. A `upcloud.StorageImportSourceHTTPImport` source
// will import from an HTTP source. `upcloud.StorageImportSourceDirectUpload` will directly upload the file specified in `SourceLocation`.
func (s *ServiceContext) CreateStorageImport(ctx context.Context, r *request.CreateStorageImportRequest) (*upcloud.StorageImportDetails, error) {
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
func (s *ServiceContext) doCreateStorageImport(ctx context.Context, r *request.CreateStorageImportRequest) (*upcloud.StorageImportDetails, error) {
	storageImport := upcloud.StorageImportDetails{}
	return &storageImport, s.create(ctx, r, &storageImport)
}

// directStorageImport handles the direct upload logic including getting the upload URL and PUT the file data
// to that endpoint.
func (s *ServiceContext) directStorageImport(ctx context.Context, r *request.CreateStorageImportRequest) (*upcloud.StorageImportDetails, error) {
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
		defer f.Close()
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

	s.client.AddRequestHeaders(req)
	req.Header.Set("Content-Type", r.ContentType)
	if _, err := s.client.PerformRequest(req); err != nil {
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
func (s *ServiceContext) GetStorageImportDetails(ctx context.Context, r *request.GetStorageImportDetailsRequest) (*upcloud.StorageImportDetails, error) {
	storageDetails := upcloud.StorageImportDetails{}
	return &storageDetails, s.get(ctx, r.RequestURL(), &storageDetails)
}

// WaitForStorageImportCompletion waits for the importing storage to complete.
func (s *ServiceContext) WaitForStorageImportCompletion(ctx context.Context, r *request.WaitForStorageImportCompletionRequest) (*upcloud.StorageImportDetails, error) {
	attempts := 0
	sleepDuration := time.Second * 5

	for {
		attempts++

		storageImportDetails, err := s.GetStorageImportDetails(ctx, &request.GetStorageImportDetailsRequest{
			UUID: r.StorageUUID,
		})

		if err != nil {
			return nil, err
		}

		switch storageImportDetails.State {
		case upcloud.StorageImportStateCompleted:
			return storageImportDetails, nil
		case upcloud.StorageImportStateCancelled,
			upcloud.StorageImportStateCancelling,
			upcloud.StorageImportStateFailed:
			if storageImportDetails.ErrorCode != "" || storageImportDetails.ErrorMessage != "" {
				return storageImportDetails, &upcloud.Error{
					ErrorCode:    storageImportDetails.ErrorCode,
					ErrorMessage: storageImportDetails.ErrorMessage,
				}
			}
			return storageImportDetails, &upcloud.Error{
				ErrorCode:    storageImportDetails.State,
				ErrorMessage: "Storage Import Failed",
			}
		default:
			if time.Duration(attempts)*sleepDuration >= r.Timeout {
				return nil, errors.New("timeout reached while waiting for import to complete")
			}

			time.Sleep(sleepDuration)
		}
	}
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
func (s *ServiceContext) ResizeStorageFilesystem(ctx context.Context, r *request.ResizeStorageFilesystemRequest) (*upcloud.ResizeStorageFilesystemBackup, error) {
	resizeBackup := upcloud.ResizeStorageFilesystemBackup{}
	return &resizeBackup, s.create(ctx, r, &resizeBackup)
}
