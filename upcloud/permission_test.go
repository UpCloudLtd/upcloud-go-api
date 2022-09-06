package upcloud

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPermissionsUnmarshal(t *testing.T) {
	want := Permissions{
		{
			TargetIdentifier: "0ad9408c-8563-4abf-b862-dbde5b581123",
			TargetType:       PermissionTargetManagedLoadbalancer,
			User:             "sub_account_user1",
			Options: &PermissionOptions{
				Storage: FromBool(true),
			},
		},
		{
			TargetIdentifier: "0603a187-3ede-4aae-883e-85ea3e69babc",
			TargetType:       PermissionTargetObjectStorage,
			User:             "sub_account_user2",
		},
	}
	got := Permissions{}
	err := json.Unmarshal([]byte(`
		{
			"permissions": {
				"permission": [
					{
						"target_identifier": "0ad9408c-8563-4abf-b862-dbde5b581123",
						"target_type": "managed_loadbalancer",
						"user": "sub_account_user1",
						"options": {
							"storage": "yes"
						}
					},
					{
						"target_identifier": "0603a187-3ede-4aae-883e-85ea3e69babc",
						"target_type": "object_storage",
						"user": "sub_account_user2"
					}
				]
			}
		}
	`), &got)
	assert.NoError(t, err)
	assert.Equal(t, want, got)
}

func TestPermissionMarshal(t *testing.T) {
	want := `
	{
		"target_identifier": "0ad9408c-8563-4abf-b862-dbde5b581123",
		"target_type": "managed_loadbalancer",
		"user": "sub_account_user1"
	
	}
	`
	got, err := json.Marshal(&Permission{
		TargetIdentifier: "0ad9408c-8563-4abf-b862-dbde5b581123",
		TargetType:       PermissionTargetManagedLoadbalancer,
		User:             "sub_account_user1",
	})
	assert.NoError(t, err)
	assert.JSONEq(t, want, string(got))

	want = `
	{
		"target_identifier": "0ad9408c-8563-4abf-b862-dbde5b581123",
		"target_type": "managed_loadbalancer",
		"user": "sub_account_user1",
		"options": {
			"storage": "yes"
		}
	}
	`
	got, err = json.Marshal(&Permission{
		TargetIdentifier: "0ad9408c-8563-4abf-b862-dbde5b581123",
		TargetType:       PermissionTargetManagedLoadbalancer,
		User:             "sub_account_user1",
		Options: &PermissionOptions{
			Storage: FromBool(true),
		},
	})
	assert.NoError(t, err)
	assert.JSONEq(t, want, string(got))
}
