package request

import (
	"encoding/json"
	"fmt"
)

// GetObjectStorageDetailsRequest represents a request for retrieving details about an Object Storage device.
//
// Deprecated: non-managed object storage service has been decommissioned.
type GetObjectStorageDetailsRequest struct {
	UUID string
}

// RequestURL implements the Request interface
func (r *GetObjectStorageDetailsRequest) RequestURL() string {
	return fmt.Sprintf("/object-storage/%s", r.UUID)
}

// CreateObjectStorageRequest represents a request for creating a new Object Storage device
//
// Deprecated: non-managed object storage service has been decommissioned.
type CreateObjectStorageRequest struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Zone        string `json:"zone"`
	AccessKey   string `json:"access_key"` //gosec:disable G117 -- struct field for API credential, not a hardcoded secret
	SecretKey   string `json:"secret_key"`
	Size        int    `json:"size"`
}

// MarshalJSON is a custom marshaller that deals with
// deeply embedded values.
func (r CreateObjectStorageRequest) MarshalJSON() ([]byte, error) {
	type localCreateObjectStorageRequest CreateObjectStorageRequest
	v := struct {
		ObjectStorage localCreateObjectStorageRequest `json:"object_storage"`
	}{}
	v.ObjectStorage = localCreateObjectStorageRequest(r)

	return json.Marshal(&v)
}

// RequestURL implements the Request interface
func (r *CreateObjectStorageRequest) RequestURL() string {
	return "/object-storage"
}

// ModifyObjectStorageRequest represents a request to modify an Object Storage.
//
// Deprecated: non-managed object storage service has been decommissioned.
type ModifyObjectStorageRequest struct {
	UUID        string `json:"-"`
	Description string `json:"description,omitempty"`
	AccessKey   string `json:"access_key,omitempty"` //gosec:disable G117 -- struct field for API credential, not a hardcoded secret
	SecretKey   string `json:"secret_key,omitempty"`
	Size        int    `json:"size,omitempty"`
}

// MarshalJSON is a custom marshaller that deals with
// deeply embedded values.
func (r ModifyObjectStorageRequest) MarshalJSON() ([]byte, error) {
	type localModifyObjectStorageRequest ModifyObjectStorageRequest
	v := struct {
		ModifyObjectStorageRequest localModifyObjectStorageRequest `json:"object_storage"`
	}{}
	v.ModifyObjectStorageRequest = localModifyObjectStorageRequest(r)

	return json.Marshal(&v)
}

// RequestURL implements the Request interface
func (r *ModifyObjectStorageRequest) RequestURL() string {
	return fmt.Sprintf("/object-storage/%s", r.UUID)
}

// DeleteObjectStorageRequest represents a request to delete an Object Storage.
//
// Deprecated: non-managed object storage service has been decommissioned.
type DeleteObjectStorageRequest struct {
	UUID string
}

// RequestURL implements the Request interface
func (r *DeleteObjectStorageRequest) RequestURL() string {
	return fmt.Sprintf("/object-storage/%s", r.UUID)
}
