package service

import (
	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud/request"
)

type Permission interface {
	GetPermissions(*request.GetPermissionsRequest) (upcloud.Permissions, error)
	GrantPermission(*request.GrantPermissionRequest) (*upcloud.Permission, error)
	RevokePermission(*request.RevokePermissionRequest) error
}

var _ Permission = (*Service)(nil)

func (s *Service) GetPermissions(r *request.GetPermissionsRequest) (upcloud.Permissions, error) {
	p := make(upcloud.Permissions, 0)
	return p, s.get(r.RequestURL(), &p)
}

func (s *Service) GrantPermission(r *request.GrantPermissionRequest) (*upcloud.Permission, error) {
	p := upcloud.Permission{}
	resp := struct{ Permission *upcloud.Permission }{Permission: &p}
	return &p, s.create(r, &resp)
}

func (s *Service) RevokePermission(r *request.RevokePermissionRequest) error {
	return s.create(r, nil)
}
