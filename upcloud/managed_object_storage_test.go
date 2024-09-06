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
			CustomDomains: []ManagedObjectStorageCustomDomain{
				{
					DomainName: "objects.example.com",
					Type:       "public",
				},
			},
			Endpoints: []ManagedObjectStorageEndpoint{
				{
					DomainName: "7mf5k.upbucket.com",
					Type:       "public",
					IAMURL:     "https://7mf5k.upbucket.com:4443/iam",
					STSURL:     "https://7mf5k.upbucket.com:4443/sts",
				},
				{
					DomainName: "7mf5k-private.upbucket.com",
					Type:       "private",
					IAMURL:     "https://7mf5k-private.upbucket.com:4443/iam",
					STSURL:     "https://7mf5k-private.upbucket.com:4443/sts",
				},
			},
			Labels: []Label{{
				Key:   "managedBy",
				Value: "upcloud-go-sdk-unit-test",
			}},
			Name: "go-sdk-test-objsto",
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
			UUID:             "1200ecde-db95-4d1c-9133-6508f3232567",
		},
		`
		{
			"configured_status": "started",
			"created_at": "2023-05-07T15:55:24.655776Z",
			"custom_domains": [
				{
					"domain_name": "objects.example.com",
					"type": "public"
				}
			],
			"endpoints": [
				{
					"domain_name": "7mf5k.upbucket.com",
	                "iam_url": "https://7mf5k.upbucket.com:4443/iam",
    	            "sts_url": "https://7mf5k.upbucket.com:4443/sts",
					"type": "public"
				},
				{
					"domain_name": "7mf5k-private.upbucket.com",
					"iam_url": "https://7mf5k-private.upbucket.com:4443/iam",
					"sts_url": "https://7mf5k-private.upbucket.com:4443/sts",
					"type": "private"
				}
			],
			"labels": [
				{
					"key": "managedBy",
					"value": "upcloud-go-sdk-unit-test"
				}
			],
			"name": "go-sdk-test-objsto",
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
			"uuid": "1200ecde-db95-4d1c-9133-6508f3232567"
		}
		`,
	)
}
