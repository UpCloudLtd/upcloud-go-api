package request

import (
	"encoding/json"
	"fmt"

	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud"
)

const (
	serverGroupBasePath string = "/server-group"
)

// GetServerGroupsRequest represents a request to list server groups
type GetServerGroupsRequest struct{}

func (s GetServerGroupsRequest) RequestURL() string {
	return serverGroupBasePath
}

type ServerGroupFilter interface {
	ToQueryParam() string
}

// GetServerGroupsWithFiltersRequest represents a request to get
// all server groups using labels or label keys as filters.
// Using multiple filters returns only groups that match all.
type GetServerGroupsWithFiltersRequest struct {
	Filters []ServerGroupFilter
}

// RequestURL implements the Request interface.
func (r *GetServerGroupsWithFiltersRequest) RequestURL() string {
	if len(r.Filters) == 0 {
		return serverGroupBasePath
	}

	params := ""
	for _, v := range r.Filters {
		if len(params) > 0 {
			params += "&"
		}
		params += v.ToQueryParam()
	}

	return fmt.Sprintf("%s?%s", serverGroupBasePath, params)
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
	Labels       *upcloud.LabelSlice     `json:"labels,omitempty"`
	Members      upcloud.ServerUUIDSlice `json:"servers,omitempty"`
	AntiAffinity upcloud.Boolean         `json:"anti_affinity,omitempty"`
	Title        string                  `json:"title,omitempty"`
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
	Labels       *upcloud.LabelSlice      `json:"labels,omitempty"`
	Members      *upcloud.ServerUUIDSlice `json:"servers,omitempty"`
	AntiAffinity upcloud.Boolean          `json:"anti_affinity,omitempty"`
	Title        string                   `json:"title,omitempty"`
	UUID         string                   `json:"-"`
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
