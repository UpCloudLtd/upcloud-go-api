package upcloud

import "encoding/json"

type PermissionTarget string

const (
	PermissionTargetServer              PermissionTarget = "server"
	PermissionTargetStorage             PermissionTarget = "storage"
	PermissionTargetNetwork             PermissionTarget = "network"
	PermissionTargetRouter              PermissionTarget = "router"
	PermissionTargetObjectStorage       PermissionTarget = "object_storage"
	PermissionTargetManagedDatabase     PermissionTarget = "managed_database"
	PermissionTargetManagedLoadbalancer PermissionTarget = "managed_loadbalancer"
	PermissionTargetTagAccess           PermissionTarget = "tag_access"
)

type PermissionOptions struct {
	Storage Boolean `json:"storage,omitempty"`
}

type Permission struct {
	TargetIdentifier string             `json:"target_identifier,omitempty"`
	TargetType       PermissionTarget   `json:"target_type,omitempty"`
	User             string             `json:"user,omitempty"`
	Options          *PermissionOptions `json:"options,omitempty"`
}

type Permissions []Permission

// UnmarshalJSON is a custom unmarshaller that deals with deeply embedded values.
func (p *Permissions) UnmarshalJSON(b []byte) error {
	v := struct {
		Permissions struct {
			Permission []Permission `json:"permission"`
		} `json:"permissions"`
	}{}

	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}

	*p = Permissions(v.Permissions.Permission)
	return nil
}
