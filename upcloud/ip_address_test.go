package upcloud

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestUnmarshalIPAddresses tests that IPAddresses and IPAddress structs are unmarshaled correctly.
func TestUnmarshalIPAddresses(t *testing.T) {
	originalJSON := `
	  {
		"ip_addresses": {
		  "ip_address": [
			{
			  "access": "utility",
			  "address": "10.0.0.0",
			  "family": "IPv4",
			  "ptr_record": "",
			  "server": "0053cd80-5945-4105-9081-11192806a8f7",
			  "mac": "mm:mm:mm:mm:mm:m1",
			  "floating": "no",
			  "zone": "fi-hel2"
			},
			{
			  "access": "utility",
			  "address": "10.0.0.1",
			  "family": "IPv4",
			  "ptr_record": "",
			  "server": "006b6701-55d2-4374-ac40-56cc1501037f",
			  "mac": "mm:mm:mm:mm:mm:m2",
			  "floating": "no",
			  "zone": "de-fra1"
			},
			{
			  "access": "public",
			  "address": "x.x.x.x",
			  "family": "IPv4",
			  "part_of_plan": "yes",
			  "ptr_record": "x-x-x-x.zone.upcloud.host",
			  "server": "0053cd80-5945-4105-9081-11192806a8f7",
			  "mac": "mm:mm:mm:mm:mm:m1",
			  "floating": "yes",
			  "zone": "de-fra1"
			},
			{
			  "access": "public",
			  "address": "xxxx:xxxx:xxxx:xxxx:xxxx:xxxx:xxxx:xxxx",
			  "family": "IPv6",
			  "ptr_record": "xxxx-xxxx-xxxx-xxxx.v6.zone.upcloud.host",
			  "server": "006b6701-55d2-4374-ac40-56cc1501037f",
			  "mac": "mm:mm:mm:mm:mm:m3",
			  "floating": "no",
			  "zone": "fi-hel1"
			},
			{
			  "access": "public",
			  "address": "y.y.y.y",
			  "family": "IPv4",
			  "ptr_record": "y.y.y.y.zone.host.upcloud.com",
			  "floating": "yes",
			  "zone": "nl-ams1"
			}
		  ]
		}
	  }
	`

	ipAddresses := IPAddresses{}
	err := json.Unmarshal([]byte(originalJSON), &ipAddresses)
	assert.NoError(t, err)
	assert.Len(t, ipAddresses.IPAddresses, 5)

	testData := []IPAddress{
		{
			Access:     IPAddressAccessUtility,
			Address:    "10.0.0.0",
			Family:     IPAddressFamilyIPv4,
			PTRRecord:  "",
			ServerUUID: "0053cd80-5945-4105-9081-11192806a8f7",
			Floating:   False,
			MAC:        "mm:mm:mm:mm:mm:m1",
			Zone:       "fi-hel2",
		},
		{
			Access:     IPAddressAccessUtility,
			Address:    "10.0.0.1",
			Family:     IPAddressFamilyIPv4,
			PTRRecord:  "",
			ServerUUID: "006b6701-55d2-4374-ac40-56cc1501037f",
			Floating:   False,
			MAC:        "mm:mm:mm:mm:mm:m2",
			Zone:       "de-fra1",
		},
		{
			Access:     IPAddressAccessPublic,
			Address:    "x.x.x.x",
			Family:     IPAddressFamilyIPv4,
			PartOfPlan: True,
			PTRRecord:  "x-x-x-x.zone.upcloud.host",
			ServerUUID: "0053cd80-5945-4105-9081-11192806a8f7",
			Floating:   True,
			MAC:        "mm:mm:mm:mm:mm:m1",
			Zone:       "de-fra1",
		},
		{
			Access:     IPAddressAccessPublic,
			Address:    "xxxx:xxxx:xxxx:xxxx:xxxx:xxxx:xxxx:xxxx",
			Family:     IPAddressFamilyIPv6,
			PTRRecord:  "xxxx-xxxx-xxxx-xxxx.v6.zone.upcloud.host",
			ServerUUID: "006b6701-55d2-4374-ac40-56cc1501037f",
			Floating:   False,
			MAC:        "mm:mm:mm:mm:mm:m3",
			Zone:       "fi-hel1",
		},
		{
			Access:    IPAddressAccessPublic,
			Address:   "y.y.y.y",
			Family:    IPAddressFamilyIPv4,
			PTRRecord: "y.y.y.y.zone.host.upcloud.com",
			Floating:  True,
			Zone:      "nl-ams1",
		},
	}

	for i, v := range testData {
		address := ipAddresses.IPAddresses[i]
		assert.Equal(t, v, address)
	}
}

// TestUnmarshalIPAddress tests that IPAddress is unmarshaled correctly on its own.
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
