package upcloud

import (
	"testing"
)

func TestManagedObjectStorage(t *testing.T) {
	t.Parallel()
	testJSON(t,
		&ManagedObjectStorage{},
		&ManagedObjectStorage{
			ConfiguredStatus: ManagedObjectStorageConfiguredStatusStarted,
			CreatedAt:        timeParse("2023-05-07T15:55:24.655776Z"),
			Endpoints: []ManagedObjectStorageEndpoint{
				{
					DomainName: "7mf5k.upbucket.com",
					Type:       "public",
				},
				{
					DomainName: "7mf5k-private.upbucket.com",
					Type:       "private",
				},
			},
			Labels: []Label{{
				Key:   "managedBy",
				Value: "upcloud-go-sdk-unit-test",
			}},
			Networks: []ManagedObjectStorageNetwork{
				{
					Family: "IPv4",
					Name:   "example-public-network",
					Type:   "public",
				},
				{
					Family: "IPv4",
					Name:   "example-private-network",
					Type:   "private",
					UUID:   StringPtr("03aa7245-2ff9-49c8-9f0e-7ca0270d71a4"),
				},
			},
			OperationalState: ManagedObjectStorageOperationalStateRunning,
			Region:           "europe-1",
			UpdatedAt:        timeParse("2023-05-07T21:38:15.757405Z"),
			Users: []ManagedObjectStorageUser{
				{
					AccessKeys: []ManagedObjectStorageUserAccessKey{
						{
							AccessKeyId:     "AKIA63F41D01345BB477",
							CreatedAt:       timeParse("2023-05-07T20:52:19.705405Z"),
							Enabled:         true,
							LastUsedAt:      timeParse("2023-05-07T20:52:17Z"),
							Name:            "example-access-key",
							SecretAccessKey: nil,
							UpdatedAt:       timeParse("2023-05-07T21:06:18.81511Z"),
						},
					},
					CreatedAt:        timeParse("2023-05-07T15:55:24.655776Z"),
					OperationalState: ManagedObjectStorageUserOperationalStateReady,
					UpdatedAt:        timeParse("2023-05-07T16:48:14.744079Z"),
					Username:         "example-user",
				},
			},
			UUID: "1200ecde-db95-4d1c-9133-6508f3232567",
		},
		`
		{
			"configured_status": "started",
			"created_at": "2023-05-07T15:55:24.655776Z",
			"endpoints": [
				{
					"domain_name": "7mf5k.upbucket.com",
					"type": "public"
				},
				{
					"domain_name": "7mf5k-private.upbucket.com",
					"type": "private"
				}
			],
			"labels": [
				{
					"key": "managedBy",
					"value": "upcloud-go-sdk-unit-test"
				}
			],
			"networks": [
				{
					"family": "IPv4",
					"name": "example-public-network",
					"type": "public"
				},
				{
					"family": "IPv4",
					"name": "example-private-network",
					"type": "private",
					"uuid": "03aa7245-2ff9-49c8-9f0e-7ca0270d71a4"
				}
			],
			"operational_state": "running",
			"region": "europe-1",
			"updated_at": "2023-05-07T21:38:15.757405Z",
			"users": [
				{
					"access_keys": [
						{
							"access_key_id": "AKIA63F41D01345BB477",
							"created_at": "2023-05-07T20:52:19.705405Z",
							"enabled": true,
							"last_used_at": "2023-05-07T20:52:17Z",
							"name": "example-access-key",
							"updated_at": "2023-05-07T21:06:18.81511Z"
						}
					],
					"created_at": "2023-05-07T15:55:24.655776Z",
					"operational_state": "ready",
					"updated_at": "2023-05-07T16:48:14.744079Z",
					"username": "example-user"
				}
			],
			"uuid": "1200ecde-db95-4d1c-9133-6508f3232567"
		}
		`,
	)
}
