package request

import (
	"encoding/json"
	"testing"

	"github.com/UpCloudLtd/upcloud-go-api/v7/upcloud"
	"github.com/stretchr/testify/assert"
)

func TestCreateSubaccountRequest(t *testing.T) {
	expectedJSON := `
	{
		"sub_account": {
			"first_name": "first",
			"last_name": "last",
			"company": "my company name",
			"address": "my address",
			"postal_code": "00130",
			"city": "Helsinki",
			"state": "",
			"country": "FIN",
			"phone": "+358.31245434",
			"email": "user@example.com",
			"vat_number": "my vat number",
			"timezone": "Europe/Helsinki",
			"username": "myusername",
			"password": "mysecr3tPassword",
			"currency": "USD",
			"language": "en",
			"roles": {
				"role": [
					"technical",
					"billing",
					"aux_billing"
				]
			},
			"allow_gui": "no",
			"allow_api": "yes",
			"network_access": {
				"network": ["127.0.0.1"]
			},
			"server_access": {
				"server": [{
					"uuid": "*",
					"storage": "no"
				}]
			},
			"storage_access": {
				"storage": ["*"]
			},
			"tag_access": {
				"tag": [{
					"name": "test_tag",
					"storage": "yes"
				}]
			},
			"ip_filters": {
				"ip_filter": ["127.0.0.2"]
			}
		}
	}
	`
	r := CreateSubaccountRequest{
		Subaccount: CreateSubaccount{
			Username:   "myusername",
			Password:   "mysecr3tPassword",
			FirstName:  "first",
			LastName:   "last",
			Company:    "my company name",
			Address:    "my address",
			PostalCode: "00130",
			City:       "Helsinki",
			Email:      "user@example.com",
			Phone:      "+358.31245434",
			State:      "",
			Country:    "FIN",
			Currency:   "USD",
			Language:   "en",
			VATNnumber: "my vat number",
			Timezone:   "Europe/Helsinki",
			AllowAPI:   upcloud.True,
			AllowGUI:   upcloud.False,
			TagAccess: upcloud.AccountTagAccess{
				Tag: []upcloud.AccountTag{
					{Name: "test_tag", Storage: upcloud.True},
				},
			},
			Roles: upcloud.AccountRoles{
				Role: []string{"technical", "billing", "aux_billing"},
			},
			ServerAccess: upcloud.AccountServerAccess{
				Server: []upcloud.AccountServer{
					{UUID: "*", Storage: upcloud.False},
				},
			},
			StorageAccess: upcloud.AccountStorageAccess{
				Storage: []string{"*"},
			},
			NetworkAccess: upcloud.AccountNetworkAccess{
				Network: []string{"127.0.0.1"},
			},
			IPFilters: upcloud.AccountIPFilters{
				IPFilter: []string{"127.0.0.2"},
			},
		},
	}

	actualJSON, err := json.Marshal(&r)
	assert.NoError(t, err)
	assert.JSONEq(t, expectedJSON, string(actualJSON))
}

func TestModifySubaccountRequest(t *testing.T) {
	expectedJSON := `
	{
		"account": {
			"first_name": "first",
			"last_name": "last",
			"company": "my company name",
			"address": "my address",
			"postal_code": "00130",
			"city": "Helsinki",
			"state": "s",
			"country": "FIN",
			"phone": "+358.31245434",
			"email": "user@example.com",
			"vat_number": "my vat number",
			"timezone": "Europe/Helsinki",
			"currency": "USD",
			"language": "en",
			"roles": {
				"role": [
					"technical",
					"billing",
					"aux_billing"
				]
			},
			"allow_gui": "no",
			"allow_api": "yes",
			"network_access": {
				"network": ["127.0.0.1"]
			},
			"server_access": {
				"server": [{
					"uuid": "*",
					"storage": "no"
				}]
			},
			"storage_access": {
				"storage": ["*"]
			},
			"tag_access": {
				"tag": [{
					"name": "test_tag",
					"storage": "yes"
				}]
			},
			"ip_filters": {
				"ip_filter": ["127.0.0.2"]
			}
		}
	}
	`
	r := ModifySubaccountRequest{
		Username: "username",
		Subaccount: ModifySubaccount{
			FirstName:  "first",
			LastName:   "last",
			Company:    "my company name",
			Address:    "my address",
			PostalCode: "00130",
			City:       "Helsinki",
			Email:      "user@example.com",
			Phone:      "+358.31245434",
			State:      "s",
			Country:    "FIN",
			Currency:   "USD",
			Language:   "en",
			VATNnumber: "my vat number",
			Timezone:   "Europe/Helsinki",
			AllowAPI:   upcloud.True,
			AllowGUI:   upcloud.False,
			TagAccess: upcloud.AccountTagAccess{
				Tag: []upcloud.AccountTag{
					{Name: "test_tag", Storage: upcloud.True},
				},
			},
			Roles: upcloud.AccountRoles{
				Role: []string{"technical", "billing", "aux_billing"},
			},
			ServerAccess: upcloud.AccountServerAccess{
				Server: []upcloud.AccountServer{
					{UUID: "*", Storage: upcloud.False},
				},
			},
			StorageAccess: upcloud.AccountStorageAccess{
				Storage: []string{"*"},
			},
			NetworkAccess: upcloud.AccountNetworkAccess{
				Network: []string{"127.0.0.1"},
			},
			IPFilters: upcloud.AccountIPFilters{
				IPFilter: []string{"127.0.0.2"},
			},
		},
	}

	// Check marshaling
	actualJSON, err := json.Marshal(&r)
	assert.NoError(t, err)
	assert.JSONEq(t, expectedJSON, string(actualJSON))
}
