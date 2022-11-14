package request

import "github.com/UpCloudLtd/upcloud-go-api/v5/upcloud"

// GetPermissionsRequest represents a request to get permissions
type GetPermissionsRequest struct{}

func (r *GetPermissionsRequest) RequestURL() string {
	return "/permission"
}

// GrantPermissionRequest represents a request to grant permission
type GrantPermissionRequest struct {
	Permission upcloud.Permission `json:"permission,omitempty"`
}

func (r *GrantPermissionRequest) RequestURL() string {
	return "/permission/grant"
}

// RevokePermissionRequest represents a request to revoke permission
type RevokePermissionRequest struct {
	Permission upcloud.Permission `json:"permission,omitempty"`
}

func (r *RevokePermissionRequest) RequestURL() string {
	return "/permission/revoke"
}
