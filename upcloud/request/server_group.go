package request

import (
	"encoding/json"
	"fmt"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
)

const (
	serverGroupBasePath string = "/server-group"
)

// GetServerGroupsRequest represents a request to list server groups
type GetServerGroupsRequest struct {
	Filters []QueryFilter
}

func (s GetServerGroupsRequest) RequestURL() string {
	if len(s.Filters) == 0 {
		return serverGroupBasePath
	}

	return fmt.Sprintf("%s?%s", serverGroupBasePath, encodeQueryFilters(s.Filters))
}

// Deprecated: ServerGroupFilter filter is deprecated. Use QueryFilter instead.
type ServerGroupFilter = QueryFilter

// Deprecated: GetServerGroupsWithFiltersRequest is deprecated. Use GetServerGroupsRequest instead.
type GetServerGroupsWithFiltersRequest struct {
	Filters []QueryFilter
}

func (r *GetServerGroupsWithFiltersRequest) RequestURL() string {
	if len(r.Filters) == 0 {
		return serverGroupBasePath
	}
	return fmt.Sprintf("%s?%s", serverGroupBasePath, encodeQueryFilters(r.Filters))
}

// GetServerGroupsRequest represents a request to get server group details
type GetServerGroupRequest struct {
	UUID string `json:"-"`
}

func (s GetServerGroupRequest) RequestURL() string {
	return fmt.Sprintf("%s/%s", serverGroupBasePath, s.UUID)
}

// CreateServerGroupRequest represents a request to create server group
type CreateServerGroupRequest struct {
	Labels             *upcloud.LabelSlice                   `json:"labels,omitempty"`
	Members            upcloud.ServerUUIDSlice               `json:"servers,omitempty"`
	AntiAffinityPolicy upcloud.ServerGroupAntiAffinityPolicy `json:"anti_affinity,omitempty"`
	Title              string                                `json:"title,omitempty"`
}

func (s CreateServerGroupRequest) RequestURL() string {
	return serverGroupBasePath
}

// MarshalJSON is a custom marshaller that deals with deeply embedded values.
func (r CreateServerGroupRequest) MarshalJSON() ([]byte, error) {
	type c CreateServerGroupRequest
	v := struct {
		ServerGroup c `json:"server_group"`
	}{}

	v.ServerGroup = c(r)

	return json.Marshal(&v)
}

// ModifyServerGroupRequest represents a request to modify server group
type ModifyServerGroupRequest struct {
	Labels             *upcloud.LabelSlice                   `json:"labels,omitempty"`
	Members            *upcloud.ServerUUIDSlice              `json:"servers,omitempty"`
	AntiAffinityPolicy upcloud.ServerGroupAntiAffinityPolicy `json:"anti_affinity,omitempty"`
	Title              string                                `json:"title,omitempty"`
	UUID               string                                `json:"-"`
}

func (s ModifyServerGroupRequest) RequestURL() string {
	return fmt.Sprintf("%s/%s", serverGroupBasePath, s.UUID)
}

// MarshalJSON is a custom marshaller that deals with deeply embedded values.
func (r ModifyServerGroupRequest) MarshalJSON() ([]byte, error) {
	type c ModifyServerGroupRequest
	v := struct {
		ServerGroup c `json:"server_group"`
	}{}

	v.ServerGroup = c(r)

	return json.Marshal(&v)
}

// DeleteServerGroupRequest represents a request to delete server group
type DeleteServerGroupRequest struct {
	UUID string `json:"-"`
}

func (s DeleteServerGroupRequest) RequestURL() string {
	return fmt.Sprintf("%s/%s", serverGroupBasePath, s.UUID)
}

// AddServerToServerGroupRequest represents a request to add server to a server group
type AddServerToServerGroupRequest struct {
	ServerUUID string `json:"uuid,omitempty"`
	UUID       string `json:"-"`
}

func (s AddServerToServerGroupRequest) RequestURL() string {
	return fmt.Sprintf("%s/%s/servers", serverGroupBasePath, s.UUID)
}

// MarshalJSON is a custom marshaller that deals with deeply embedded values.
func (r AddServerToServerGroupRequest) MarshalJSON() ([]byte, error) {
	type c AddServerToServerGroupRequest
	v := struct {
		Server c `json:"server"`
	}{}

	v.Server = c(r)

	return json.Marshal(&v)
}

// RemoveServerFromServerGroupRequest represents a request to remove server from a server group
type RemoveServerFromServerGroupRequest struct {
	ServerUUID string `json:"-"`
	UUID       string `json:"-"`
}

func (s RemoveServerFromServerGroupRequest) RequestURL() string {
	return fmt.Sprintf("%s/%s/servers/%s", serverGroupBasePath, s.UUID, s.ServerUUID)
}
