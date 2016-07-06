package request

import (
	"encoding/xml"
	"fmt"
	"github.com/jalle19/upcloud-go-sdk/upcloud"
)

/**
Represents a request for retrieving all or some storages
*/
type GetStoragesRequest struct {
	// If specified, only storages with this access type will be retrieved
	Access string
	// If specified, only storages with this type will be retrieved
	Type string
	// If specified, only storages marked as favorite will be retrieved
	Favorite bool
}

func (r *GetStoragesRequest) RequestURL() string {
	if r.Access != "" {
		return fmt.Sprintf("/storage/%s", r.Access)
	}

	if r.Type != "" {
		return fmt.Sprintf("/storage/%s", r.Type)
	}

	if r.Favorite {
		return "/storage/favorite"
	}

	return "/storage"
}

/**
Represents a request for retrieving details about a piece of storage
*/
type GetStorageDetailsRequest struct {
	UUID string
}

func (r *GetStorageDetailsRequest) RequestURL() string {
	return fmt.Sprintf("/storage/%s", r.UUID)
}

/**
Represents a request to create a storage device
*/
type CreateStorageRequest struct {
	XMLName xml.Name `xml:"storage"`

	Size       int                 `xml:"size"`
	Tier       string              `xml:"tier,omitempty"`
	Title      string              `xml:"title"`
	Zone       string              `xml:"zone"`
	BackupRule *upcloud.BackupRule `xml:"backup_rule,omitempty"`
}

func (r *CreateStorageRequest) RequestURL() string {
	return "/storage"
}

/**
Represents a request to modify a storage device
*/
type ModifyStorageRequest struct {
	XMLName xml.Name `xml:"storage"`
	UUID    string   `xml:"-"`

	Title      string              `xml:"title,omitempty"`
	Size       int                 `xml:"size,omitempty"`
	BackupRule *upcloud.BackupRule `xml:"backup_rule,omitempty"`
}

func (r *ModifyStorageRequest) RequestURL() string {
	return fmt.Sprintf("/storage/%s", r.UUID)
}

/**
Represents a request to attach a storage device to a server
*/
type AttachStorageRequest struct {
	XMLName    xml.Name `xml:"storage_device"`
	ServerUUID string   `xml:"-"`

	Type        string `xml:"type,omitempty"`
	Address     string `xml:"address,omitempty"`
	StorageUUID string `xml:"storage,omitempty"`
}

func (r *AttachStorageRequest) RequestURL() string {
	return fmt.Sprintf("/server/%s/storage/attach", r.ServerUUID)
}

/**
Represents a request to detch a storage device from a server
*/
type DetachStorageRequest struct {
	XMLName    xml.Name `xml:"storage_device"`
	ServerUUID string   `xml:"-"`

	Address string `xml:"address"`
}

func (r *DetachStorageRequest) RequestURL() string {
	return fmt.Sprintf("/server/%s/storage/detach", r.ServerUUID)
}

/**
Represents a request to delete a storage device
*/
type DeleteStorageRequest struct {
	UUID string
}

func (r *DeleteStorageRequest) RequestURL() string {
	return fmt.Sprintf("/storage/%s", r.UUID)
}
