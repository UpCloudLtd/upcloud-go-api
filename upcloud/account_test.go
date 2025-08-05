package upcloud

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestUnmarshalAccount tests that Account objects unmarshal correctly
func TestUnmarshalAccount(t *testing.T) {
	originalJSON := `
	  {
		"account": {
		  "credits": 9972.2324,
		  "username": "username",
		  "resource_limits": {
			"cores": 200,
			"detached_floating_ips": 10,
			"managed_object_storages": 7,
			"memory": 1048576,
			"network_peerings": 100,
			"networks": 100,
			"ntp_excess_gib": 20000,
			"public_ipv4": 100,
			"public_ipv6": 100,
			"storage_hdd": 10240,
			"storage_maxiops": 10240,
			"storage_ssd": 10240,
			"load_balancers": 50,
			"gpus": 8
		  }
		}
	  }
	`

	account := Account{}
	err := json.Unmarshal([]byte(originalJSON), &account)
	assert.NoError(t, err)
	assert.Equal(t, 9972.2324, account.Credits)
	assert.Equal(t, "username", account.UserName)
	assert.Equal(t, 200, account.ResourceLimits.Cores)
	assert.Equal(t, 10, account.ResourceLimits.DetachedFloatingIps)
	assert.Equal(t, 7, account.ResourceLimits.ManagedObjectStorages)
	assert.Equal(t, 1048576, account.ResourceLimits.Memory)
	assert.Equal(t, 100, account.ResourceLimits.NetworkPeerings)
	assert.Equal(t, 100, account.ResourceLimits.Networks)
	assert.Equal(t, 20000, account.ResourceLimits.NTPExcessGiB)
	assert.Equal(t, 100, account.ResourceLimits.PublicIPv4)
	assert.Equal(t, 100, account.ResourceLimits.PublicIPv6)
	assert.Equal(t, 10240, account.ResourceLimits.StorageHDD)
	assert.Equal(t, 10240, account.ResourceLimits.StorageMaxIOPS)
	assert.Equal(t, 10240, account.ResourceLimits.StorageSSD)
	assert.Equal(t, 50, account.ResourceLimits.LoadBalancers)
	assert.Equal(t, 8, account.ResourceLimits.GPUs)
}

// TestMarshalAccount tests that Account objects marshal correctly
func TestMarshalAccount(t *testing.T) {
	request := Account{
		Credits:  100,
		UserName: "username",
		ResourceLimits: ResourceLimits{
			Memory: 123,
		},
	}

	expectedJSON := `
	  {
      "username": "username",
      "credits": 100,
      "resource_limits": {
        "memory": 123
      }
	  }
	`

	actualJSON, err := json.Marshal(&request)
	println(string(actualJSON))
	assert.NoError(t, err)
	assert.JSONEq(t, expectedJSON, string(actualJSON))
}

// TestUnmarshalAccountDetails tests that AccountDetails objects unmarshal correctly
func TestUnmarshalAccountDetails(t *testing.T) {
	originalJSON := []byte(`
	{
		"account": {
			"main_account": "mymain",
			"type": "sub",
			"username": "my_sub_account",
			"first_name": "first",
			"last_name": "last",
			"company": "UpCloud Ltd",
			"address": "my address",
			"postal_code": "00130",
			"city": "Helsinki",
			"state": "",
			"country": "FIN",
			"currency": "USD",
			"language": "fi",
			"phone": "+358.31245434",
			"email": "test@myhost.mydomain",
			"vat_number": "FI24315605",
			"timezone": "UTC",
			"campaigns": {
				"campaign": [
					"test1",
					"test2"
				]
			},
			"roles": {
				"role": [
					"billing",
					"aux_billing",
					"technical"
				]
			},
			"allow_api": "yes",
			"allow_gui": "no",
			"enable_3rd_party_services": "yes",
			"network_access": {
				"network": [
					"*"
				]
			},
			"server_access": {
				"server": [
					{
						"storage": "no",
						"uuid": "*"
					}
				]
			},
			"storage_access": {
				"storage": [
					"*"
				]
			},
			"tag_access": {
				"tag": [
					{
						"name": "mytag",
						"storage": "yes"
					},
					{
						"name": "mytag2",
						"storage": "no"
					}
				]
			},
			"ip_filters": {
				"ip_filter": [
					"10.0.0.1-255.255.255.255"
				]
			}
		}
	}
	`)

	a := AccountDetails{}
	err := json.Unmarshal(originalJSON, &a)
	assert.NoError(t, err)
	assert.True(t, a.IsSubaccount())
	assert.Len(t, a.Campaigns.Campaign, 2)
	assert.Equal(t, "test1", a.Campaigns.Campaign[0])
	assert.Len(t, a.Roles.Role, 3)
	assert.Equal(t, "billing", a.Roles.Role[0])
	assert.Len(t, a.TagAccess.Tag, 2)
	assert.Equal(t, "mytag", a.TagAccess.Tag[0].Name)
	assert.Len(t, a.ServerAccess.Server, 1)
	assert.Equal(t, "*", a.ServerAccess.Server[0].UUID)
	assert.Len(t, a.StorageAccess.Storage, 1)
	assert.Equal(t, "*", a.StorageAccess.Storage[0])
	assert.Len(t, a.NetworkAccess.Network, 1)
	assert.Equal(t, "*", a.NetworkAccess.Network[0])
	assert.Len(t, a.IPFilters.IPFilter, 1)
	assert.Equal(t, "10.0.0.1-255.255.255.255", a.IPFilters.IPFilter[0])
	assert.Equal(t, a.MainAccount, "mymain")
	assert.Equal(t, a.Type, AccountType("sub"))
	assert.Equal(t, a.Username, "my_sub_account")
	assert.Equal(t, a.FirstName, "first")
	assert.Equal(t, a.LastName, "last")
	assert.Equal(t, a.Company, "UpCloud Ltd")
	assert.Equal(t, a.Address, "my address")
	assert.Equal(t, a.PostalCode, "00130")
	assert.Equal(t, a.City, "Helsinki")
	assert.Equal(t, a.State, "")
	assert.Equal(t, a.Country, "FIN")
	assert.Equal(t, a.Currency, "USD")
	assert.Equal(t, a.Language, "fi")
	assert.Equal(t, a.Phone, "+358.31245434")
	assert.Equal(t, a.Email, "test@myhost.mydomain")
	assert.Equal(t, a.VATNnumber, "FI24315605")
	assert.Equal(t, a.Timezone, "UTC")
}

// TestUnmarshalAccountList tests that Accounts slice unmarshal correctly
func TestUnmarshalAccountList(t *testing.T) {
	originalJSON := []byte(`
	{
		"accounts": {
			"account": [
				{
					"roles": {
						"role": [
							"technical"
						]
					},
					"type": "mymain",
					"username": "test"
				},
				{
					"roles": {
						"role": [
							"billing"
						]
					},
					"type": "sub",
					"username": "my_billing_account"
				}
			]
		}
	}
	`)
	a := make(AccountList, 0)
	err := json.Unmarshal(originalJSON, &a)
	assert.NoError(t, err)
	assert.Len(t, a, 2)
	assert.Equal(t, "billing", a[1].Roles.Role[0])
	assert.Equal(t, AccountType("sub"), a[1].Type)
	assert.Equal(t, "my_billing_account", a[1].Username)
}

// TestUnmarshalBillingSummary tests that BillingSummary unmarshal correctly
func TestUnmarshalBillingSummary(t *testing.T) {
	originalJSON := []byte(`
{
  "currency": "EUR",
  "managed_databases": {
    "managed_database": {
      "resources": [
        {
          "amount": 15.10223,
          "details": [
            {
              "amount": 15.10223,
              "hours": 420,
              "plan": "1x1xCPU-2GB-25GB",
              "zone": "fi-hel2"
            }
          ],
          "hours": 420,
          "resource_id": "09001a4d-b525-4ab8-835d-000000000000"
        }
      ],
      "total_amount": 15.10223
    },
    "total_amount": 15.10223
  },
  "managed_object_storages": {
    "managed_object_storage": {
      "resources": [
        {
          "amount": 5.03384,
          "details": [
            {
              "amount": 5.03384,
             "billable_size_gib": 500,
              "hours": 420,
              "zone": "fi-hel2"
            }
          ],
          "hours": 420,
          "resource_id": "127ee42a-8304-477c-84e1-000000000000"
        }
      ],
      "total_amount": 5.03384
    },
    "total_amount": 5.03384
  },
  "servers": {
    "server": {
      "resources": [
        {
          "amount": 16.18091,
          "details": [
            {
              "amount": 16.18091,
              "hours": 420,
              "labels": [
                {
                  "key": "test",
                  "value": ""
                }
              ],
              "plan": "2xCPU-4GB",
              "zone": "fi-hel2"
            }
          ],
          "hours": 420,
          "resource_id": "001c2caa-00dc-4eff-9321-000000000000"
        },
        {
          "amount": 16.10386,
          "details": [
            {
              "amount": 16.10386,
              "hours": 418,
              "labels": [
                {
                  "key": "test2",
                  "value": ""
                }
              ],
              "plan": "2xCPU-4GB",
              "zone": "fi-hel2"
            },
            {
              "amount": 0.9288,
              "hours": 12,
              "labels": [
                {
                  "key": "test2",
                  "value": ""
                }
              ],
              "plan": "4xCPU-8GB",
              "zone": "fi-hel2"
            }
          ],
          "hours": 430,
          "resource_id": "0010cf04-3608-4512-b4de-000000000000"
        }
      ],
      "total_amount": 33.21357
    },
    "total_amount": 33.21357
  },
  "total_amount": 53.34964
}
  `)
	var b BillingSummary
	err := json.Unmarshal(originalJSON, &b)
	assert.NoError(t, err)

	assert.Equal(t, "EUR", b.Currency)
	assert.Equal(t, 53.34964, b.TotalAmount)

	assert.NotNil(t, b.Servers)
	assert.Equal(t, 33.21357, b.Servers.TotalAmount)
	assert.NotNil(t, b.Servers.Server)
	assert.Equal(t, 33.21357, b.Servers.Server.TotalAmount)
	assert.Len(t, b.Servers.Server.Resources, 2)

	assert.Equal(t, "001c2caa-00dc-4eff-9321-000000000000", b.Servers.Server.Resources[0].ResourceID)
	assert.Equal(t, 16.18091, b.Servers.Server.Resources[0].Amount)
	assert.Equal(t, 420, b.Servers.Server.Resources[0].Hours)
	assert.Equal(t, []BillingResourceDetail{
		{
			Amount: 16.18091,
			Hours:  420,
			Plan:   "2xCPU-4GB",
			Zone:   "fi-hel2",
			Labels: []Label{
				{Key: "test", Value: ""},
			},
		},
	}, b.Servers.Server.Resources[0].Details)

	assert.Equal(t, "0010cf04-3608-4512-b4de-000000000000", b.Servers.Server.Resources[1].ResourceID)
	assert.Equal(t, 16.10386, b.Servers.Server.Resources[1].Amount)
	assert.Equal(t, 430, b.Servers.Server.Resources[1].Hours)
	assert.Equal(t, []BillingResourceDetail{
		{
			Amount: 16.10386,
			Hours:  418,
			Plan:   "2xCPU-4GB",
			Zone:   "fi-hel2",
			Labels: []Label{
				{Key: "test2", Value: ""},
			},
		},
		{
			Amount: 0.9288,
			Hours:  12,
			Plan:   "4xCPU-8GB",
			Zone:   "fi-hel2",
			Labels: []Label{
				{Key: "test2", Value: ""},
			},
		},
	}, b.Servers.Server.Resources[1].Details)

	assert.NotNil(t, b.ManagedDatabases)
	assert.Equal(t, 15.10223, b.ManagedDatabases.TotalAmount)
	assert.NotNil(t, b.ManagedDatabases.ManagedDatabase)
	assert.Equal(t, 15.10223, b.ManagedDatabases.ManagedDatabase.TotalAmount)
	assert.Len(t, b.ManagedDatabases.ManagedDatabase.Resources, 1)
	assert.Equal(t, "09001a4d-b525-4ab8-835d-000000000000", b.ManagedDatabases.ManagedDatabase.Resources[0].ResourceID)
	assert.Equal(t, 15.10223, b.ManagedDatabases.ManagedDatabase.Resources[0].Amount)
	assert.Equal(t, 420, b.ManagedDatabases.ManagedDatabase.Resources[0].Hours)
	assert.Equal(t, []BillingResourceDetail{
		{
			Amount: 15.10223,
			Hours:  420,
			Plan:   "1x1xCPU-2GB-25GB",
			Zone:   "fi-hel2",
		},
	}, b.ManagedDatabases.ManagedDatabase.Resources[0].Details)

	assert.NotNil(t, b.ManagedObjectStorages)
	assert.Equal(t, 5.03384, b.ManagedObjectStorages.TotalAmount)
	assert.NotNil(t, b.ManagedObjectStorages.ManagedObjectStorage)
	assert.Equal(t, 5.03384, b.ManagedObjectStorages.ManagedObjectStorage.TotalAmount)
	assert.Len(t, b.ManagedObjectStorages.ManagedObjectStorage.Resources, 1)
	assert.Equal(t, "127ee42a-8304-477c-84e1-000000000000", b.ManagedObjectStorages.ManagedObjectStorage.Resources[0].ResourceID)
	assert.Equal(t, 5.03384, b.ManagedObjectStorages.ManagedObjectStorage.Resources[0].Amount)
	assert.Equal(t, 420, b.ManagedObjectStorages.ManagedObjectStorage.Resources[0].Hours)
	assert.Equal(t, []BillingResourceDetail{
		{
			Amount:          5.03384,
			Hours:           420,
			Zone:            "fi-hel2",
			BillableSizeGiB: 500,
		},
	}, b.ManagedObjectStorages.ManagedObjectStorage.Resources[0].Details)
}
