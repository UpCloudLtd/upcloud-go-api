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
			"memory": 1048576,
			"networks": 100,
			"public_ipv4": 100,
			"public_ipv6": 100,
			"storage_hdd": 10240,
			"storage_ssd": 10240
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
	assert.Equal(t, 1048576, account.ResourceLimits.Memory)
	assert.Equal(t, 100, account.ResourceLimits.Networks)
	assert.Equal(t, 100, account.ResourceLimits.PublicIPv4)
	assert.Equal(t, 100, account.ResourceLimits.PublicIPv6)
	assert.Equal(t, 10240, account.ResourceLimits.StorageHDD)
	assert.Equal(t, 10240, account.ResourceLimits.StorageSSD)
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
