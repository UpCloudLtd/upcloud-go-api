package request

import (
	"net/url"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
)

type GetFileStoragesRequest struct {
	Page *Page
	Sort *string
}

func (r *GetFileStoragesRequest) RequestURL() string {
	base := "/file-storage"
	v := url.Values{}
	if r.Page != nil {
		for k, vals := range r.Page.Values() {
			for _, val := range vals {
				v.Add(k, val)
			}
		}
	}
	if r.Sort != nil {
		v.Set("sort", *r.Sort)
	}
	qs := v.Encode()
	if qs == "" {
		return base
	}
	return base + "?" + qs
}

type CreateFileStorageRequest struct {
	Name             string                              `json:"name"`
	Zone             string                              `json:"zone"`
	ConfiguredStatus upcloud.FileStorageConfiguredStatus `json:"configured_status"`
	SizeGiB          int                                 `json:"size_gib"`
	Networks         []upcloud.FileStorageNetwork        `json:"networks,omitempty"`
	Shares           []FileStorageShare                  `json:"shares,omitempty"`
	Labels           []upcloud.Label                     `json:"labels,omitempty"`
}

func (r *CreateFileStorageRequest) RequestURL() string { return "/file-storage" }

type GetFileStorageRequest struct{ UUID string }

func (r *GetFileStorageRequest) RequestURL() string { return "/file-storage/" + r.UUID }

type ReplaceFileStorageRequest struct {
	UUID             string                              `json:"-"`
	Name             string                              `json:"name"`
	ConfiguredStatus upcloud.FileStorageConfiguredStatus `json:"configured_status"`
	SizeGiB          int                                 `json:"size_gib"`
	Networks         []upcloud.FileStorageNetwork        `json:"networks,omitempty"`
	Labels           []upcloud.Label                     `json:"labels,omitempty"`
}

func (r *ReplaceFileStorageRequest) RequestURL() string { return "/file-storage/" + r.UUID }

type ModifyFileStorageRequest struct {
	UUID             string                               `json:"-"`
	Name             *string                              `json:"name,omitempty"`
	ConfiguredStatus *upcloud.FileStorageConfiguredStatus `json:"configured_status,omitempty"`
	SizeGiB          *int                                 `json:"size_gib,omitempty"`
	Networks         *[]upcloud.FileStorageNetwork        `json:"networks,omitempty"`
	Labels           *[]upcloud.Label                     `json:"labels,omitempty"`
}

func (r *ModifyFileStorageRequest) RequestURL() string { return "/file-storage/" + r.UUID }

type DeleteFileStorageRequest struct{ UUID string }

func (r *DeleteFileStorageRequest) RequestURL() string { return "/file-storage/" + r.UUID }

type WaitForFileStorageDeletionRequest struct {
	UUID string `json:"-"`
}

// WaitForFileStorageOperationalStateRequest represents a request to wait for a File Storage service
// to enter a desired state
type WaitForFileStorageOperationalStateRequest struct {
	DesiredState upcloud.FileStorageOperationalState `json:"-"`
	UUID         string                              `json:"-"`
}

type GetFileStorageNetworksRequest struct{ ServiceUUID string }

func (r *GetFileStorageNetworksRequest) RequestURL() string {
	return "/file-storage/" + r.ServiceUUID + "/networks"
}

type CreateFileStorageNetworkRequest struct {
	ServiceUUID string `json:"-"`
	UUID        string `json:"uuid"`
	Name        string `json:"name"`
	Family      string `json:"family"`
	IPAddress   string `json:"ip_address,omitempty"`
}

func (r *CreateFileStorageNetworkRequest) RequestURL() string {
	return "/file-storage/" + r.ServiceUUID + "/networks"
}

type GetFileStorageNetworkRequest struct {
	ServiceUUID string
	NetworkName string
}

func (r *GetFileStorageNetworkRequest) RequestURL() string {
	return "/file-storage/" + r.ServiceUUID + "/networks/" + r.NetworkName
}

type ModifyFileStorageNetworkRequest struct {
	ServiceUUID string  `json:"-"`
	NetworkName string  `json:"-"`
	Family      *string `json:"family,omitempty"`
	IPAddress   *string `json:"ip_address,omitempty"`
}

func (r *ModifyFileStorageNetworkRequest) RequestURL() string {
	return "/file-storage/" + r.ServiceUUID + "/networks/" + r.NetworkName
}

type DeleteFileStorageNetworkRequest struct {
	ServiceUUID string
	NetworkName string
}

func (r *DeleteFileStorageNetworkRequest) RequestURL() string {
	return "/file-storage/" + r.ServiceUUID + "/networks/" + r.NetworkName
}

type GetFileStorageSharesRequest struct{ ServiceUUID string }

func (r *GetFileStorageSharesRequest) RequestURL() string {
	return "/file-storage/" + r.ServiceUUID + "/shares"
}

type CreateFileStorageShareRequest struct {
	ServiceUUID string                   `json:"-"`
	Name        string                   `json:"name"`
	Path        string                   `json:"path"`
	ACL         []upcloud.FileStorageACL `json:"acl"`
}

func (r *CreateFileStorageShareRequest) RequestURL() string {
	return "/file-storage/" + r.ServiceUUID + "/shares"
}

type FileStorageShare struct {
	Name string                   `json:"name"`
	Path string                   `json:"path"`
	ACL  []upcloud.FileStorageACL `json:"acl"`
}

type GetFileStorageShareRequest struct {
	ServiceUUID string
	ShareName   string
}

func (r *GetFileStorageShareRequest) RequestURL() string {
	return "/file-storage/" + r.ServiceUUID + "/shares/" + r.ShareName
}

type ModifyFileStorageShare struct {
	Name *string                   `json:"name,omitempty"`
	ACL  *[]upcloud.FileStorageACL `json:"acl,omitempty"`
}

type ModifyFileStorageShareRequest struct {
	ModifyFileStorageShare

	ServiceUUID string `json:"-"`
	ShareName   string `json:"-"`
}

func (r *ModifyFileStorageShareRequest) RequestURL() string {
	return "/file-storage/" + r.ServiceUUID + "/shares/" + r.ShareName
}

type DeleteFileStorageShareRequest struct {
	ServiceUUID string
	ShareName   string
}

func (r *DeleteFileStorageShareRequest) RequestURL() string {
	return "/file-storage/" + r.ServiceUUID + "/shares/" + r.ShareName
}

type GetFileStorageLabelsRequest struct{ ServiceUUID string }

func (r *GetFileStorageLabelsRequest) RequestURL() string {
	return "/file-storage/" + r.ServiceUUID + "/labels"
}

type CreateFileStorageLabelRequest struct {
	upcloud.Label

	ServiceUUID string `json:"-"`
}

func (r *CreateFileStorageLabelRequest) RequestURL() string {
	return "/file-storage/" + r.ServiceUUID + "/labels"
}

type GetFileStorageLabelRequest struct {
	ServiceUUID string
	LabelKey    string
}

func (r *GetFileStorageLabelRequest) RequestURL() string {
	return "/file-storage/" + r.ServiceUUID + "/labels/" + r.LabelKey
}

type ModifyFileStorageLabelRequest struct {
	upcloud.Label

	ServiceUUID string `json:"-"`
	LabelKey    string `json:"-"`
}

func (r *ModifyFileStorageLabelRequest) RequestURL() string {
	return "/file-storage/" + r.ServiceUUID + "/labels/" + r.LabelKey
}

type DeleteFileStorageLabelRequest struct {
	ServiceUUID string
	LabelKey    string
}

func (r *DeleteFileStorageLabelRequest) RequestURL() string {
	return "/file-storage/" + r.ServiceUUID + "/labels/" + r.LabelKey
}
