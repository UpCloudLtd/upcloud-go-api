package request

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
)

func TestMarshalGetNetworks(t *testing.T) {
	request := GetNetworksRequest{}
	assert.Equal(t, "/network", request.RequestURL())

	request = GetNetworksRequest{Filters: []QueryFilter{
		FilterLabel{Label: upcloud.Label{
			Key:   "env",
			Value: "test",
		}},
		FilterLabelKey{Key: "managedBy"},
	}}

	assert.Equal(t, "/network?label=env%3Dtest&label=managedBy", request.RequestURL())
}

// TestMarshalGetNetworksInZoneRequest tests that GetNetworksInZoneRequest behaves correctly
func TestMarshalGetNetworksInZoneRequest(t *testing.T) {
	request := GetNetworksInZoneRequest{
		Zone: "foo",
	}

	assert.Equal(t, "/network/?zone=foo", request.RequestURL())

	requestWithFilters := GetNetworksInZoneRequest{
		Zone: "fi-hel1",
		Filters: []QueryFilter{
			FilterLabel{Label: upcloud.Label{
				Key:   "env",
				Value: "test",
			}},
			FilterLabelKey{
				Key: "managed",
			},
		},
	}

	assert.Equal(t, "/network?zone=fi-hel1&label=env%3Dtest&label=managed", requestWithFilters.RequestURL())
}

// TestMarshalGetNetworkDetailsRequest tests that GetNetworkDetailsRequest behaves correctly
func TestMarshalGetNetworkDetailsRequest(t *testing.T) {
	request := GetNetworkDetailsRequest{
		UUID: "foo",
	}

	assert.Equal(t, "/network/foo", request.RequestURL())
}

// TestMarshalCreateNetworkRequest tests that CreateNetworkRequest behaves correctly
func TestMarshalCreateNetworkRequest(t *testing.T) {
	request := CreateNetworkRequest{
		Name:   "Test private net",
		Zone:   "uk-lon1",
		Router: "04c0df35-2658-4b0c-8ac7-962090f4e92a",
		Labels: []upcloud.Label{
			{
				Key:   "env",
				Value: "test",
			},
		},
		IPNetworks: []upcloud.IPNetwork{
			{
				Address:          "172.16.0.0/22",
				DHCP:             upcloud.True,
				DHCPDefaultRoute: upcloud.False,
				DHCPDns: []string{
					"172.16.0.10",
					"172.16.1.10",
				},
				DHCPRoutes: []string{
					"192.168.0.0/24",
					"192.168.100.100/32",
				},
				Family:  upcloud.IPAddressFamilyIPv4,
				Gateway: "172.16.0.1",
				DHCPRoutesConfiguration: upcloud.DHCPRoutesConfiguration{
					EffectiveRoutesAutoPopulation: upcloud.EffectiveRoutesAutoPopulation{
						Enabled:             upcloud.True,
						ExcludeBySource:     []upcloud.NetworkRouteSource{"static-route"},
						FilterByDestination: []string{"172.16.0.0/22"},
						FilterByRouteType:   []upcloud.NetworkRouteType{"service"},
					},
				},
			},
		},
	}

	expectedJSON := `
	  {
		"network": {
		  "name": "Test private net",
		  "zone": "uk-lon1",
		  "router": "04c0df35-2658-4b0c-8ac7-962090f4e92a",
		  "labels": [
			{
				"key": "env",
				"value": "test"
			}
		  ],
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
				"dhcp_routes" : [
				  "192.168.0.0/24",
				  "192.168.100.100/32"
				],
				"dhcp_routes_configuration": {
					"effective_routes_auto_population": {
						"enabled": "yes",
						"exclude_by_source": [
							"static-route"
						],
						"filter_by_destination": [
							"172.16.0.0/22"
						],
						"filter_by_route_type": [
							"service"
						]
					}
				},
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

// TestMarshalModifyNetworkRequest tests that ModifyNetworkRequest behaves correctly
func TestMarshalModifyNetworkRequest(t *testing.T) {
	request := ModifyNetworkRequest{
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
			"ip_networks": {
			"ip_network": [
				{
				"dhcp": "no",
				"dhcp_default_route": "no",
				"dhcp_routes_configuration": {
					"effective_routes_auto_population": {
					"enabled": "no"
					}
				},
				"family": "IPv4"
				}
			]
			},
			"name": "My private network"
		}
	}
	`

	actualJSON, err := json.Marshal(&request)
	assert.NoError(t, err)
	assert.Equal(t, "/network/foo", request.RequestURL())
	assert.JSONEq(t, expectedJSON, string(actualJSON))

	request = ModifyNetworkRequest{
		UUID: "foo",
		Name: "supername",
		Labels: &[]upcloud.Label{
			{
				Key:   "env",
				Value: "test",
			},
		},
	}

	expectedJSON = `
	  {
		"network": {
		  "name": "supername",
		  "labels": [
			{
				"key": "env",
				"value": "test"
			}
		  ]
	  	}
	  }	
	`

	actualJSON, err = json.Marshal(&request)
	assert.NoError(t, err)
	assert.Equal(t, "/network/foo", request.RequestURL())
	assert.JSONEq(t, expectedJSON, string(actualJSON))
}

// TestMarshalDeleteNetwork tests the DeleteNetworkRequest behaves correctly
func TestMarshalDeleteNetwork(t *testing.T) {
	request := DeleteNetworkRequest{
		UUID: "foo",
	}

	assert.Equal(t, "/network/foo", request.RequestURL())
}

// TestMarshalAttachNetworkRouterRequest tests that AttachNetworkRouterRequest behaves correctly.
func TestMarshalAttachNetworkRouterRequest(t *testing.T) {
	request := AttachNetworkRouterRequest{
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
	request := DetachNetworkRouterRequest{
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

// TestMarshalGetServerNetworksRequest tests the GetServerNetworksRequest behaves correctly
func TestMarshalGetServerNetworksRequest(t *testing.T) {
	request := GetServerNetworksRequest{
		ServerUUID: "foo",
	}

	assert.Equal(t, "/server/foo/networking", request.RequestURL())
}

// TestMarshalCreateNetworkInterfaceRequest tests that CreateNetworkInterfaceRequest behaves correctly.
func TestMarshalCreateNetworkInterfaceRequest(t *testing.T) {
	request := CreateNetworkInterfaceRequest{
		ServerUUID:        "foo",
		Type:              upcloud.IPAddressAccessPrivate,
		NetworkUUID:       "0374ce47-4303-4490-987d-32dc96cfd79b",
		SourceIPFiltering: upcloud.True,
		IPAddresses: []CreateNetworkInterfaceIPAddress{
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
	request := DeleteNetworkInterfaceRequest{
		ServerUUID: "foo",
		Index:      1,
	}

	assert.Equal(t, "/server/foo/networking/interface/1", request.RequestURL())
}

// TestMarshalModifyNetworkInterfaceRequest tests that ModifyNetworkInterfaceRequest behaves correctly.
func TestMarshalModifyNetworkInterfaceRequest(t *testing.T) {
	request := ModifyNetworkInterfaceRequest{
		ServerUUID:   "foo",
		CurrentIndex: 99,
		NewIndex:     101,

		Type:              upcloud.IPAddressAccessPrivate,
		NetworkUUID:       "0374ce47-4303-4490-987d-32dc96cfd79b",
		SourceIPFiltering: upcloud.True,
		IPAddresses: []CreateNetworkInterfaceIPAddress{
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
	request := GetRouterDetailsRequest{
		UUID: "foo",
	}

	assert.Equal(t, "/router/foo", request.RequestURL())
}

// TestMarshalCreateRouterRequest tests that CreateRouterRequest behaves correctly.
func TestMarshalCreateRouterRequest(t *testing.T) {
	request := CreateRouterRequest{
		Name: "Example router",
		StaticRoutes: []upcloud.StaticRoute{
			{
				Name:    "example_static_route",
				Route:   "0.0.0.0/0",
				Nexthop: "10.0.0.100",
			},
		},
	}

	expectedJSON := `
	  {
		"router": {
		  "name": "Example router",
		  "static_routes": [
            {
			  "route": "0.0.0.0/0",
			  "nexthop": "10.0.0.100",
			  "name": "example_static_route"
            }
          ]
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
	request := ModifyRouterRequest{
		Name: "Modified router",
		StaticRoutes: &[]upcloud.StaticRoute{
			{
				Name:    "example_static_route",
				Route:   "0.0.0.0/0",
				Nexthop: "10.0.0.100",
			},
		},
		UUID: "foo",
	}

	expectedJSON := `
	  {
		"router": {
		  "name": "Modified router",
		  "static_routes": [
            {
			  "route": "0.0.0.0/0",
			  "nexthop": "10.0.0.100",
			  "name": "example_static_route"
            }
		  ]
		}
	  }
	`

	actualJSON, err := json.Marshal(&request)
	assert.NoError(t, err)

	assert.JSONEq(t, expectedJSON, string(actualJSON))

	assert.Equal(t, "/router/foo", request.RequestURL())

	request = ModifyRouterRequest{
		UUID:   "",
		Name:   "Modified router name",
		Labels: &[]upcloud.Label{},
	}

	expectedJSON = `
	  {
		"router": {
		  "name": "Modified router name",
		  "labels": []
		}
	  }
	`

	actualJSON, err = json.Marshal(&request)
	assert.NoError(t, err)

	assert.JSONEq(t, expectedJSON, string(actualJSON))
}
