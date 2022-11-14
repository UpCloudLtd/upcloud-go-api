package service

import (
	"context"

	"github.com/UpCloudLtd/upcloud-go-api/v5/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v5/upcloud/request"
)

type PermissionContext interface {
	GetPermissions(context.Context, *request.GetPermissionsRequest) (upcloud.Permissions, error)
	GrantPermission(context.Context, *request.GrantPermissionRequest) (*upcloud.Permission, error)
	RevokePermission(context.Context, *request.RevokePermissionRequest) error
}

func (s *Service) GetPermissions(ctx context.Context, r *request.GetPermissionsRequest) (upcloud.Permissions, error) {
	p := make(upcloud.Permissions, 0)
	return p, s.get(ctx, r.RequestURL(), &p)
}

func (s *Service) GrantPermission(ctx context.Context, r *request.GrantPermissionRequest) (*upcloud.Permission, error) {
	p := upcloud.Permission{}
	resp := struct{ Permission *upcloud.Permission }{Permission: &p}
	return &p, s.create(ctx, r, &resp)
}

func (s *Service) RevokePermission(ctx context.Context, r *request.RevokePermissionRequest) error {
	return s.create(ctx, r, nil)
}
