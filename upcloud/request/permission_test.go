package request

import (
	"encoding/json"
	"testing"

	"github.com/UpCloudLtd/upcloud-go-api/v6/upcloud"
	"github.com/stretchr/testify/assert"
)

func TestGetPermissionsRequest(t *testing.T) {
	assert.Equal(t, "/permission", (&GetPermissionsRequest{}).RequestURL())
}

func TestGrantPermissionsRequest(t *testing.T) {
	want := `
	{
		"permission": {
			"target_identifier": "0ad9408c-8563-4abf-b862-dbde5b581123",
			"target_type": "managed_loadbalancer",
			"user": "sub_account_user1",
			"options": {
				"storage": "yes"
			}
		}
	}
	`
	got, err := json.Marshal(&GrantPermissionRequest{
		Permission: upcloud.Permission{
			TargetIdentifier: "0ad9408c-8563-4abf-b862-dbde5b581123",
			TargetType:       upcloud.PermissionTargetManagedLoadbalancer,
			User:             "sub_account_user1",
			Options: &upcloud.PermissionOptions{
				Storage: upcloud.FromBool(true),
			},
		},
	})
	assert.NoError(t, err)
	assert.JSONEq(t, want, string(got))
	assert.Equal(t, "/permission/grant", (&GrantPermissionRequest{}).RequestURL())
}

func TestRevokePermissionsRequest(t *testing.T) {
	want := `
	{
		"permission": {
			"target_identifier": "0ad9408c-8563-4abf-b862-dbde5b581123",
			"target_type": "managed_loadbalancer",
			"user": "sub_account_user1"
		}
	}
	`
	got, err := json.Marshal(&RevokePermissionRequest{
		Permission: upcloud.Permission{
			TargetIdentifier: "0ad9408c-8563-4abf-b862-dbde5b581123",
			TargetType:       upcloud.PermissionTargetManagedLoadbalancer,
			User:             "sub_account_user1",
		},
	})
	assert.NoError(t, err)
	assert.JSONEq(t, want, string(got))
	assert.Equal(t, "/permission/revoke", (&RevokePermissionRequest{}).RequestURL())
}
