package service

import (
	"context"

	"github.com/UpCloudLtd/upcloud-go-api/v6/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v6/upcloud/request"
)

type ServerGroup interface {
	GetServerGroups(ctx context.Context, r *request.GetServerGroupsRequest) (upcloud.ServerGroups, error)
	GetServerGroup(ctx context.Context, r *request.GetServerGroupRequest) (*upcloud.ServerGroup, error)
	CreateServerGroup(ctx context.Context, r *request.CreateServerGroupRequest) (*upcloud.ServerGroup, error)
	ModifyServerGroup(ctx context.Context, r *request.ModifyServerGroupRequest) (*upcloud.ServerGroup, error)
	DeleteServerGroup(ctx context.Context, r *request.DeleteServerGroupRequest) error
}

// GetServerGroups retrieves a list of server groups with context (EXPERIMENTAL).
func (s *Service) GetServerGroups(ctx context.Context, r *request.GetServerGroupsRequest) (upcloud.ServerGroups, error) {
	groups := upcloud.ServerGroups{}
	return groups, s.get(ctx, r.RequestURL(), &groups)
}

// Deprecated: GetServerGroupsWithFilters is deprecated. User GetServerGroups instead
func (s *Service) GetServerGroupsWithFilters(ctx context.Context, r *request.GetServerGroupsWithFiltersRequest) (upcloud.ServerGroups, error) {
	groups := upcloud.ServerGroups{}
	return groups, s.get(ctx, r.RequestURL(), &groups)
}

// GetServerGroup retrieves details of a server group  with context (EXPERIMENTAL).
func (s *Service) GetServerGroup(ctx context.Context, r *request.GetServerGroupRequest) (*upcloud.ServerGroup, error) {
	group := upcloud.ServerGroup{}
	return &group, s.get(ctx, r.RequestURL(), &group)
}

// CreateServerGroup creates a new server group  with context (EXPERIMENTAL).
func (s *Service) CreateServerGroup(ctx context.Context, r *request.CreateServerGroupRequest) (*upcloud.ServerGroup, error) {
	var group upcloud.ServerGroup
	return &group, s.create(ctx, r, &group)
}

// ModifyServerGroup modifies an existing server group  with context (EXPERIMENTAL).
func (s *Service) ModifyServerGroup(ctx context.Context, r *request.ModifyServerGroupRequest) (*upcloud.ServerGroup, error) {
	var group upcloud.ServerGroup
	return &group, s.modify(ctx, r, &group)
}

// DeleteServerGroup deletes an existing server group  with context (EXPERIMENTAL).
func (s *Service) DeleteServerGroup(ctx context.Context, r *request.DeleteServerGroupRequest) error {
	return s.delete(ctx, r)
}
