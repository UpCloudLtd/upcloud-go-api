package request

import (
	"encoding/json"
	"fmt"

	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud"
)

// GetServerGroupsRequest represents a request to list server groups
type GetServerGroupsRequest struct{}

func (s GetServerGroupsRequest) RequestURL() string {
	return "/server-group"
}

// GetServerGroupsRequest represents a request to get server group details
type GetServerGroupRequest struct {
	UUID string `json:"-"`
}

func (s GetServerGroupRequest) RequestURL() string {
	return fmt.Sprintf("/server-group/%s", s.UUID)
}

// CreateServerGroupRequest represents a request to create server group
type CreateServerGroupRequest struct {
	Title   string                  `json:"title,omitempty"`
	Members upcloud.ServerUUIDSlice `json:"servers,omitempty"`
}

func (s CreateServerGroupRequest) RequestURL() string {
	return "/server-group"
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
	UUID    string                   `json:"-"`
	Title   string                   `json:"title,omitempty"`
	Members *upcloud.ServerUUIDSlice `json:"servers,omitempty"`
}

func (s ModifyServerGroupRequest) RequestURL() string {
	return fmt.Sprintf("/server-group/%s", s.UUID)
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
	return fmt.Sprintf("/server-group/%s", s.UUID)
}
