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
			"load_balancers": 50
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
