package service

import (
	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud/request"
)

type ServerGroup interface {
	GetServerGroups(r *request.GetServerGroupsRequest) (upcloud.ServerGroups, error)
	GetServerGroup(r *request.GetServerGroupRequest) (*upcloud.ServerGroup, error)
	CreateServerGroup(r *request.CreateServerGroupRequest) (*upcloud.ServerGroup, error)
	ModifyServerGroup(r *request.ModifyServerGroupRequest) (*upcloud.ServerGroup, error)
	DeleteServerGroup(r *request.DeleteServerGroupRequest) error
}

var _ ServerGroup = (*Service)(nil)

// GetServerGroups retrieves a list of server groups (EXPERIMENTAL).
func (s *Service) GetServerGroups(r *request.GetServerGroupsRequest) (upcloud.ServerGroups, error) {
	groups := upcloud.ServerGroups{}
	return groups, s.get(r.RequestURL(), &groups)
}

// GetServerGroupsWithFilters retrieves a list of server groups with filters (EXPERIMENTAL).
func (s *Service) GetServerGroupsWithFilters(r *request.GetServerGroupsWithFiltersRequest) (upcloud.ServerGroups, error) {
	groups := upcloud.ServerGroups{}
	return groups, s.get(r.RequestURL(), &groups)
}

// GetServerGroup retrieves details of a server group (EXPERIMENTAL).
func (s *Service) GetServerGroup(r *request.GetServerGroupRequest) (*upcloud.ServerGroup, error) {
	group := upcloud.ServerGroup{}
	return &group, s.get(r.RequestURL(), &group)
}

// CreateServerGroup creates a new server group (EXPERIMENTAL).
func (s *Service) CreateServerGroup(r *request.CreateServerGroupRequest) (*upcloud.ServerGroup, error) {
	var group upcloud.ServerGroup
	return &group, s.create(r, &group)
}

// ModifyServerGroup modifies an existing server group (EXPERIMENTAL).
func (s *Service) ModifyServerGroup(r *request.ModifyServerGroupRequest) (*upcloud.ServerGroup, error) {
	var group upcloud.ServerGroup
	return &group, s.modify(r, &group)
}

// DeleteServerGroup deletes an existing server group (EXPERIMENTAL).
func (s *Service) DeleteServerGroup(r *request.DeleteServerGroupRequest) error {
	return s.delete(r)
}
