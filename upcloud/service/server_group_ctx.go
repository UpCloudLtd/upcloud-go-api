package service

import (
	"context"

	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud/request"
)

type ServerGroupContext interface {
	GetServerGroups(ctx context.Context, r *request.GetServerGroupsRequest) (upcloud.ServerGroups, error)
	GetServerGroup(ctx context.Context, r *request.GetServerGroupRequest) (*upcloud.ServerGroup, error)
	CreateServerGroup(ctx context.Context, r *request.CreateServerGroupRequest) (*upcloud.ServerGroup, error)
	ModifyServerGroup(ctx context.Context, r *request.ModifyServerGroupRequest) (*upcloud.ServerGroup, error)
	DeleteServerGroup(ctx context.Context, r *request.DeleteServerGroupRequest) error
}

// GetServerGroups retrieves a list of server groups with context (EXPERIMENTAL).
func (s *ServiceContext) GetServerGroups(ctx context.Context, r *request.GetServerGroupsRequest) (upcloud.ServerGroups, error) {
	groups := upcloud.ServerGroups{}
	return groups, s.get(ctx, r.RequestURL(), &groups)
}

// GetServerGroupsWithFilters retrieves a list of server groups with filters (EXPERIMENTAL).
func (s *ServiceContext) GetServerGroupsWithFilters(ctx context.Context, r *request.GetServerGroupsWithFiltersRequest) (upcloud.ServerGroups, error) {
	groups := upcloud.ServerGroups{}
	return groups, s.get(ctx, r.RequestURL(), &groups)
}

// GetServerGroup retrieves details of a server group  with context (EXPERIMENTAL).
func (s *ServiceContext) GetServerGroup(ctx context.Context, r *request.GetServerGroupRequest) (*upcloud.ServerGroup, error) {
	group := upcloud.ServerGroup{}
	return &group, s.get(ctx, r.RequestURL(), &group)
}

// CreateServerGroup creates a new server group  with context (EXPERIMENTAL).
func (s *ServiceContext) CreateServerGroup(ctx context.Context, r *request.CreateServerGroupRequest) (*upcloud.ServerGroup, error) {
	var group upcloud.ServerGroup
	return &group, s.create(ctx, r, &group)
}

// ModifyServerGroup modifies an existing server group  with context (EXPERIMENTAL).
func (s *ServiceContext) ModifyServerGroup(ctx context.Context, r *request.ModifyServerGroupRequest) (*upcloud.ServerGroup, error) {
	var group upcloud.ServerGroup
	return &group, s.modify(ctx, r, &group)
}

// DeleteServerGroup deletes an existing server group  with context (EXPERIMENTAL).
func (s *ServiceContext) DeleteServerGroup(ctx context.Context, r *request.DeleteServerGroupRequest) error {
	return s.delete(ctx, r)
}
