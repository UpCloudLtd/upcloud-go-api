package request_test

import (
	"encoding/json"
	"testing"

	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"

	"github.com/stretchr/testify/assert"

	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
)

// TestMarshalGetNetworksInZoneRequest tests that GetNetworksInZoneRequest behaves correctly.
func TestMarshalGetNetworksInZoneRequest(t *testing.T) {
	request := request.GetNetworksInZoneRequest{
		Zone: "foo",
	}

	assert.Equal(t, "/network/?zone=foo", request.RequestURL())
}

// TestMarshalGetNetworkDetailsRequest tests that GetNetworkDetailsRequest behaves correctly.
func TestMarshalGetNetworkDetailsRequest(t *testing.T) {
	request := request.GetNetworkDetailsRequest{
		UUID: "foo",
	}

	assert.Equal(t, "/network/foo", request.RequestURL())
}

// TestMarshalCreateNetworkRequest tests that CreateNetworkRequest behaves correctly.
func TestMarshalCreateNetworkRequest(t *testing.T) {
	request := request.CreateNetworkRequest{
		Name:   "Test private net",
		Zone:   "uk-lon1",
		Router: "04c0df35-2658-4b0c-8ac7-962090f4e92a",
		IPNetworks: []upcloud.IPNetwork{
			{
				Address:          "172.16.0.0/22",
				DHCP:             upcloud.True,
				DHCPDefaultRoute: upcloud.False,
				DHCPDns: []string{
					"172.16.0.10",
					"172.16.1.10",
				},
				Family:  upcloud.IPAddressFamilyIPv4,
				Gateway: "172.16.0.1",
			},
		},
	}

	expectedJSON := `
	  {
		"network": {
		  "name": "Test private net",
		  "zone": "uk-lon1",
		  "router": "04c0df35-2658-4b0c-8ac7-962090f4e92a",
		  "ip_networks" : {
			"ip_network" : [
			  {
				"address" : "172.16.0.0/22",
				"dhcp" : "yes",
				"dhcp_default_route" : "no",
				"dhcp_dns" : [
				  "172.16.0.10",
				  "172.16.1.10"
				],
				"family" : "IPv4",
				"gateway" : "172.16.0.1"
			  }
			]
		  }
		}
	  }
	`

	actualJSON, err := json.Marshal(&request)
	assert.NoError(t, err)
	assert.Equal(t, "/network/", request.RequestURL())
	assert.JSONEq(t, expectedJSON, string(actualJSON))
}

// TestMarshalModifyNetworkRequest tests that ModifyNetworkRequest behaves correctly.
func TestMarshalModifyNetworkRequest(t *testing.T) {
	request := request.ModifyNetworkRequest{
		UUID: "foo",

		Name: "My private network",
		IPNetworks: []upcloud.IPNetwork{
			{
				DHCP:   upcloud.False,
				Family: upcloud.IPAddressFamilyIPv4,
			},
		},
	}

	expectedJSON := `
	  {
		"network": {
		  "name": "My private network",
			"ip_networks": {
			  "ip_network": [
				{
				  "dhcp": "no",
				  "dhcp_default_route": "no",
				  "family" : "IPv4"
				}
			  ]
			}
		  }
	  }
	`

	actualJSON, err := json.Marshal(&request)
	assert.NoError(t, err)
	assert.Equal(t, "/network/foo", request.RequestURL())
	assert.JSONEq(t, expectedJSON, string(actualJSON))
}

// TestMarshalDeleteNetwork tests the DeleteNetworkRequest behaves correctly.
func TestMarshalDeleteNetwork(t *testing.T) {
	request := request.DeleteNetworkRequest{
		UUID: "foo",
	}

	assert.Equal(t, "/network/foo", request.RequestURL())
}

// TestMarshalAttachNetworkRouterRequest tests that AttachNetworkRouterRequest behaves correctly.
func TestMarshalAttachNetworkRouterRequest(t *testing.T) {
	request := request.AttachNetworkRouterRequest{
		NetworkUUID: "mocknetworkuuid",
		RouterUUID:  "mockrouteruuid",
	}

	expectedJSON := `
	  {
		"network": {
		  "router": "mockrouteruuid"
		}
	  }
	`

	actualJSON, err := json.Marshal(&request)
	assert.NoError(t, err)
	assert.JSONEq(t, expectedJSON, string(actualJSON))
	assert.Equal(t, "/network/mocknetworkuuid", request.RequestURL())
}

// TestMarshalDetachNetworkRouterRequest tests that DetachNetworkRouterRequest behaves correctly.
func TestMarshalDetachNetworkRouterRequest(t *testing.T) {
	request := request.DetachNetworkRouterRequest{
		NetworkUUID: "mocknetworkuuid",
	}

	expectedJSON := `
	  {
		"network": {
		  "router": null
		}
	  }
	`

	actualJSON, err := json.Marshal(&request)
	assert.NoError(t, err)
	assert.JSONEq(t, expectedJSON, string(actualJSON))
	assert.Equal(t, "/network/mocknetworkuuid", request.RequestURL())
}

// TestMarshalGetServerNetworksRequest tests the GetServerNetworksRequest behaves correctly.
func TestMarshalGetServerNetworksRequest(t *testing.T) {
	request := request.GetServerNetworksRequest{
		ServerUUID: "foo",
	}

	assert.Equal(t, "/server/foo/networking", request.RequestURL())
}

// TestMarshalCreateNetworkInterfaceRequest tests that CreateNetworkInterfaceRequest behaves correctly.
func TestMarshalCreateNetworkInterfaceRequest(t *testing.T) {
	request := request.CreateNetworkInterfaceRequest{
		ServerUUID:        "foo",
		Type:              upcloud.IPAddressAccessPrivate,
		NetworkUUID:       "0374ce47-4303-4490-987d-32dc96cfd79b",
		SourceIPFiltering: upcloud.True,
		IPAddresses: []request.CreateNetworkInterfaceIPAddress{
			{
				Address: "10.0.0.20",
				Family:  upcloud.IPAddressFamilyIPv4,
			},
		},
	}

	expectedJSON := `
	  {
		"interface": {
		  "type": "private",
		  "network": "0374ce47-4303-4490-987d-32dc96cfd79b",
		  "ip_addresses": {
			"ip_address": [
			  {
				"family": "IPv4",
				"address": "10.0.0.20"
			  }
			]
		  },
		  "source_ip_filtering": "yes"
		}
	  }
	`

	actualJSON, err := json.Marshal(&request)
	assert.NoError(t, err)
	assert.JSONEq(t, expectedJSON, string(actualJSON))
	assert.Equal(t, "/server/foo/networking/interface", request.RequestURL())
}

// TestMarshalDeleteNetworkInterfaceRequest tests that DeleteNetworkInterfaceRequest behaves correctly.
func TestMarshalDeleteNetworkInterfaceRequest(t *testing.T) {
	request := request.DeleteNetworkInterfaceRequest{
		ServerUUID: "foo",
		Index:      1,
	}

	assert.Equal(t, "/server/foo/networking/interface/1", request.RequestURL())
}

// TestMarshalModifyNetworkInterfaceRequest tests that ModifyNetworkInterfaceRequest behaves correctly.
func TestMarshalModifyNetworkInterfaceRequest(t *testing.T) {
	request := request.ModifyNetworkInterfaceRequest{
		ServerUUID:   "foo",
		CurrentIndex: 99,
		NewIndex:     101,

		Type:              upcloud.IPAddressAccessPrivate,
		NetworkUUID:       "0374ce47-4303-4490-987d-32dc96cfd79b",
		SourceIPFiltering: upcloud.True,
		IPAddresses: []request.CreateNetworkInterfaceIPAddress{
			{
				Address: "10.0.0.20",
				Family:  upcloud.IPAddressFamilyIPv4,
			},
		},
	}

	expectedJSON := `
	  {
		"interface": {
		  "index": 101,
		  "type": "private",
		  "network": "0374ce47-4303-4490-987d-32dc96cfd79b",
		  "ip_addresses": {
			"ip_address": [
			  {
				"family": "IPv4",
				"address": "10.0.0.20"
			  }
			]
		  },
		  "source_ip_filtering": "yes"
		}
	  }
	`

	actualJSON, err := json.Marshal(&request)
	assert.NoError(t, err)
	assert.JSONEq(t, expectedJSON, string(actualJSON))

	assert.Equal(t, "/server/foo/networking/interface/99", request.RequestURL())
}

// TestMarshalGetRouterDetailsRequest tests that GetRouterDetailsRequest behaves correctly.
func TestMarshalGetRouterDetailsRequest(t *testing.T) {
	request := request.GetRouterDetailsRequest{
		UUID: "foo",
	}

	assert.Equal(t, "/router/foo", request.RequestURL())
}

// TestMarshalCreateRouterRequest tests that CreateRouterRequest behaves correctly.
func TestMarshalCreateRouterRequest(t *testing.T) {
	request := request.CreateRouterRequest{
		Name: "Example router",
	}

	expectedJSON := `
	  {
		"router": {
		  "name": "Example router"
		}
	  }
	`

	actualJSON, err := json.Marshal(&request)
	assert.NoError(t, err)

	assert.JSONEq(t, expectedJSON, string(actualJSON))

	assert.Equal(t, "/router", request.RequestURL())
}

// TestMarshalModifyRouterRequest tests that ModifyRouterRequest behaves correctly.
func TestMarshalModifyRouterRequest(t *testing.T) {
	request := request.ModifyRouterRequest{
		Name: "Modified router",
		UUID: "foo",
	}

	expectedJSON := `
	  {
		"router": {
		  "name": "Modified router"
		}
	  }
	`

	actualJSON, err := json.Marshal(&request)
	assert.NoError(t, err)

	assert.JSONEq(t, expectedJSON, string(actualJSON))

	assert.Equal(t, "/router/foo", request.RequestURL())
}
