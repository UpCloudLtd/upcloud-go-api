package upcloud

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestUnmarshalIPAddresses tests that IPAddresses and IPAddress structs are unmarshaled correctly
func TestUnmarshalIPAddresses(t *testing.T) {
	originalJSON := `
      {
        "ip_addresses": {
          "ip_address": [
            {
              "access": "private",
              "address": "10.0.0.0",
              "family": "IPv4",
              "ptr_record": "",
              "server": "0053cd80-5945-4105-9081-11192806a8f7"
            },
            {
              "access": "private",
              "address": "10.0.0.1",
              "family": "IPv4",
              "ptr_record": "",
              "server": "006b6701-55d2-4374-ac40-56cc1501037f"
            },
            {
              "access": "public",
              "address": "x.x.x.x",
              "family": "IPv4",
              "part_of_plan": "yes",
              "ptr_record": "x-x-x-x.zone.upcloud.host",
              "server": "0053cd80-5945-4105-9081-11192806a8f7"
            },
            {
              "access": "public",
              "address": "xxxx:xxxx:xxxx:xxxx:xxxx:xxxx:xxxx:xxxx",
              "family": "IPv6",
              "ptr_record": "xxxx-xxxx-xxxx-xxxx.v6.zone.upcloud.host",
              "server": "006b6701-55d2-4374-ac40-56cc1501037f"
            }
          ]
        }
      }
    `

	ipAddresses := IPAddresses{}
	err := json.Unmarshal([]byte(originalJSON), &ipAddresses)
	assert.NoError(t, err)
	assert.Len(t, ipAddresses.IPAddresses, 4)

	testData := []IPAddress{
		{
			Access:     IPAddressAccessPrivate,
			Address:    "10.0.0.0",
			Family:     IPAddressFamilyIPv4,
			PTRRecord:  "",
			ServerUUID: "0053cd80-5945-4105-9081-11192806a8f7",
		},
		{
			Access:     IPAddressAccessPrivate,
			Address:    "10.0.0.1",
			Family:     IPAddressFamilyIPv4,
			PTRRecord:  "",
			ServerUUID: "006b6701-55d2-4374-ac40-56cc1501037f",
		},
		{
			Access:     IPAddressAccessPublic,
			Address:    "x.x.x.x",
			Family:     IPAddressFamilyIPv4,
			PartOfPlan: "yes",
			PTRRecord:  "x-x-x-x.zone.upcloud.host",
			ServerUUID: "0053cd80-5945-4105-9081-11192806a8f7",
		},
		{
			Access:     IPAddressAccessPublic,
			Address:    "xxxx:xxxx:xxxx:xxxx:xxxx:xxxx:xxxx:xxxx",
			Family:     IPAddressFamilyIPv6,
			PTRRecord:  "xxxx-xxxx-xxxx-xxxx.v6.zone.upcloud.host",
			ServerUUID: "006b6701-55d2-4374-ac40-56cc1501037f",
		},
	}

	for i, v := range testData {
		address := ipAddresses.IPAddresses[i]
		assert.Equal(t, v.Access, address.Access)
		assert.Equal(t, v.Address, address.Address)
		assert.Equal(t, v.Family, address.Family)
		assert.Equal(t, v.PTRRecord, address.PTRRecord)
		assert.Equal(t, v.ServerUUID, address.ServerUUID)
	}
}

// TestUnmarshalIPAddress tests that IPAddress is unmarshaled correctly on its own
func TestUnmarshalIPAddress(t *testing.T) {
	originalJSON := `
      {
        "ip_address" : {
           "access" : "public",
           "address" : "94.237.104.58",
           "family" : "IPv4",
           "ptr_record" : "94-237-104-58.fi-hel2.upcloud.host",
           "server" : "0028ab30-491a-4696-a601-91e810d154a8"
        }
      }
    `

	ipAddress := IPAddress{}
	err := json.Unmarshal([]byte(originalJSON), &ipAddress)
	assert.NoError(t, err)

	assert.Equal(t, IPAddressAccessPublic, ipAddress.Access)
	assert.Equal(t, "94.237.104.58", ipAddress.Address)
	assert.Equal(t, IPAddressFamilyIPv4, ipAddress.Family)
	assert.Equal(t, "94-237-104-58.fi-hel2.upcloud.host", ipAddress.PTRRecord)
	assert.Equal(t, "0028ab30-491a-4696-a601-91e810d154a8", ipAddress.ServerUUID)
}
