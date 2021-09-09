package upcloud_test

import (
	"encoding/json"
	"testing"

	"github.com/UpCloudLtd/upcloud-go-api/upcloud"

	"github.com/stretchr/testify/assert"
)

// TestUnmarshalNetworks ensures that the unmarshalling of a Networks response
// behaves correctly.
func TestUnmarshalNetworks(t *testing.T) {
	t.Parallel()
	originalJSON := `
	{
		"networks": {
		  "network": [
			{
			  "ip_networks" : {
				"ip_network": [
				  {
					"address": "80.69.172.0/22",
					"dhcp": "yes",
					"dhcp_default_route": "yes",
					"dhcp_dns": [
					  "94.237.127.9",
					  "94.237.40.9"
					],
					"family": "IPv4",
					"gateway": "80.69.172.1"
				  }
				]
			  },
			  "name": "Public 80.69.172.0/22",
			  "type": "public",
			  "uuid": "03000000-0000-4000-8001-000000000000",
			  "zone": "fi-hel1",
			  "servers": {
				"server": [
				  {"uuid": "007e3200-268f-4848-8b45-bd88c44555d2", "title": "Helsinki server #1"},
				  {"uuid": "00c8f13a-945a-48b8-bf5c-db2d7a3a37fe", "title": "Helsinki server #2"}
				]
			  }
			},
			{
			  "ip_networks" : {
				"ip_network": [
				  {
					"address": "80.69.173.0/22",
					"dhcp": "no",
					"dhcp_default_route": "no",
					"dhcp_dns": [
					  "94.237.17.9",
					  "94.237.4.9"
					],
					"family": "IPv6",
					"gateway": "80.6.172.1"
				  }
				]
			  },
			  "name": "Public 80.69.173.0/22",
			  "type": "utility",
			  "uuid": "03000011-0000-4000-8001-000000000000",
			  "zone": "fi-hel2",
			  "servers": {
				"server": [
				  {"uuid": "008e3200-268f-4848-8b45-bd88c44555d2", "title": "Helsinki server #1"},
				  {"uuid": "00d8f13a-945a-48b8-bf5c-db2d7a3a37fe", "title": "Helsinki server #2"}
				]
			  }
			}
		  ]
		}
	  }
	`

	var networks upcloud.Networks
	err := json.Unmarshal([]byte(originalJSON), &networks)
	assert.NoError(t, err)
	assert.Len(t, networks.Networks, 2)

	testNetworks := []upcloud.Network{
		{
			IPNetworks: []upcloud.IPNetwork{
				{
					Address:          "80.69.172.0/22",
					DHCP:             upcloud.True,
					DHCPDefaultRoute: upcloud.True,
					DHCPDns: []string{
						"94.237.127.9",
						"94.237.40.9",
					},
					Family:  upcloud.IPAddressFamilyIPv4,
					Gateway: "80.69.172.1",
				},
			},
			Name: "Public 80.69.172.0/22",
			Type: upcloud.NetworkTypePublic,
			UUID: "03000000-0000-4000-8001-000000000000",
			Zone: "fi-hel1",
			Servers: []upcloud.NetworkServer{
				{
					ServerUUID:  "007e3200-268f-4848-8b45-bd88c44555d2",
					ServerTitle: "Helsinki server #1",
				},
				{
					ServerUUID:  "00c8f13a-945a-48b8-bf5c-db2d7a3a37fe",
					ServerTitle: "Helsinki server #2",
				},
			},
		},
		{
			IPNetworks: []upcloud.IPNetwork{
				{
					Address:          "80.69.173.0/22",
					DHCP:             upcloud.False,
					DHCPDefaultRoute: upcloud.False,
					DHCPDns: []string{
						"94.237.17.9",
						"94.237.4.9",
					},
					Family:  upcloud.IPAddressFamilyIPv6,
					Gateway: "80.6.172.1",
				},
			},
			Name: "Public 80.69.173.0/22",
			Type: upcloud.NetworkTypeUtility,
			UUID: "03000011-0000-4000-8001-000000000000",
			Zone: "fi-hel2",
			Servers: []upcloud.NetworkServer{
				{
					ServerUUID:  "008e3200-268f-4848-8b45-bd88c44555d2",
					ServerTitle: "Helsinki server #1",
				},
				{
					ServerUUID:  "00d8f13a-945a-48b8-bf5c-db2d7a3a37fe",
					ServerTitle: "Helsinki server #2",
				},
			},
		},
	}

	for i, n := range testNetworks {
		assert.Equal(t, n, networks.Networks[i])
	}
}

// TestUnmarshalNetwork ensures that the unmarshalling of a single Network response
// behaves correctly.
func TestUnmarshalNetwork(t *testing.T) {
	t.Parallel()
	originalJSON := `
	  {
		"network": {
		  "name": "Test private net",
		  "type": "private",
		  "uuid": "034c12bc-cf15-4b19-97b2-0ab4e51bb98d",
		  "zone": "uk-lon1",
		  "router": "04c0df35-2658-4b0c-8ac7-962090f4e92a",
		  "ip_networks": {
			"ip_network": [
			  {
				"address": "172.16.0.0/22",
				"dhcp": "yes",
				"dhcp_default_route": "no",
				"dhcp_dns" : [
				  "172.16.0.10",
				  "172.16.1.10"
				],
				"family": "IPv4",
				"gateway": "172.16.0.1"
			  }
			]
		  },
		  "servers": {
			"server": [
			  {"uuid": "009d64ef-31d1-4684-a26b-c86c955cbf46", "title": "London server #1"},
			  {"uuid": "0035079f-9d66-42d5-aa74-12090e7b4ed1", "title": "London server #2"}
			]
		  }
		}
	  }
	`

	var network upcloud.Network
	err := json.Unmarshal([]byte(originalJSON), &network)
	assert.NoError(t, err)

	testNetwork := upcloud.Network{
		Name:   "Test private net",
		Type:   upcloud.NetworkTypePrivate,
		UUID:   "034c12bc-cf15-4b19-97b2-0ab4e51bb98d",
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
		Servers: []upcloud.NetworkServer{
			{
				ServerUUID:  "009d64ef-31d1-4684-a26b-c86c955cbf46",
				ServerTitle: "London server #1",
			},
			{
				ServerUUID:  "0035079f-9d66-42d5-aa74-12090e7b4ed1",
				ServerTitle: "London server #2",
			},
		},
	}

	assert.Equal(t, testNetwork, network)
}

// TestUnmarshalServerNetworks ensures that the unmarshalling of a ServerNetworks response
// behaves correctly.
func TestUnmarshalServerNetworks(t *testing.T) {
	t.Parallel()
	originalJSON := `
	  {
		"networking": {
		  "interfaces": {
			"interface": [
			  {
				"index": 2,
				"ip_addresses": {
				  "ip_address": [
					{
					  "address": "94.237.0.207",
					  "family": "IPv4"
					}
				  ]
				},
				"mac": "de:ff:ff:ff:66:89",
				"network": "037fcf2a-6745-45dd-867e-f9479ea8c044",
				"source_ip_filtering": "yes",
				"type": "public",
				"bootable": "no"
			  },
			  {
				"index": 3,
				"ip_addresses": {
				  "ip_address": [
					{
					  "address": "10.199.3.15",
					  "family": "IPv4"
					}
				  ]
				},
				"mac": "de:ff:ff:ff:ed:85",
				"network": "03c93fd8-cc60-4849-91b8-6e404b228e2a",
				"source_ip_filtering": "yes",
				"type": "utility",
				"bootable": "no"
			  },
			  {
				"index": 4,
				"ip_addresses": {
				  "ip_address": [
					{
					  "address": "10.0.0.20",
					  "family": "IPv4"
					}
				  ]
				},
				"mac": "de:ff:ff:ff:cc:20",
				"network": "0374ce47-4303-4490-987d-32dc96cfd79b",
				"source_ip_filtering": "yes",
				"type": "private",
				"bootable": "no"
			  }
			]
		  }
		}
	  }
	`

	var networking upcloud.Networking
	err := json.Unmarshal([]byte(originalJSON), &networking)
	assert.NoError(t, err)

	testNetworking := upcloud.Networking{
		Interfaces: []upcloud.ServerInterface{
			{
				Index: 2,
				IPAddresses: []upcloud.IPAddress{
					{
						Address: "94.237.0.207",
						Family:  upcloud.IPAddressFamilyIPv4,
					},
				},
				MAC:               "de:ff:ff:ff:66:89",
				Network:           "037fcf2a-6745-45dd-867e-f9479ea8c044",
				SourceIPFiltering: upcloud.True,
				Type:              upcloud.NetworkTypePublic,
				Bootable:          upcloud.False,
			},
			{
				Index: 3,
				IPAddresses: []upcloud.IPAddress{
					{
						Address: "10.199.3.15",
						Family:  upcloud.IPAddressFamilyIPv4,
					},
				},
				MAC:               "de:ff:ff:ff:ed:85",
				Network:           "03c93fd8-cc60-4849-91b8-6e404b228e2a",
				SourceIPFiltering: upcloud.True,
				Type:              upcloud.NetworkTypeUtility,
				Bootable:          upcloud.False,
			},
			{
				Index: 4,
				IPAddresses: []upcloud.IPAddress{
					{
						Address: "10.0.0.20",
						Family:  upcloud.IPAddressFamilyIPv4,
					},
				},
				MAC:               "de:ff:ff:ff:cc:20",
				Network:           "0374ce47-4303-4490-987d-32dc96cfd79b",
				SourceIPFiltering: upcloud.True,
				Type:              upcloud.NetworkTypePrivate,
				Bootable:          upcloud.False,
			},
		},
	}

	assert.Equal(t, testNetworking, networking)
}

// TestUnmarshalInterface ensures that the unmarshalling of an Interface response
// behaves correctly.
func TestUnmarshalInterface(t *testing.T) {
	t.Parallel()
	originalJSON := `
	  {
		"interface": {
		  "index": 4,
		  "ip_addresses" : {
			"ip_address" : [
			  {
			   "address" : "10.0.0.20",
			   "family" : "IPv4",
			   "floating" : "no"
			  }
			]
		  },
		  "mac": "de:ff:ff:ff:86:cf",
		  "network": "0374ce47-4303-4490-987d-32dc96cfd79b",
		  "source_ip_filtering": "yes",
		  "type": "private",
		  "bootable": "no"
		}
	  }
	`

	var iface upcloud.Interface
	err := json.Unmarshal([]byte(originalJSON), &iface)
	assert.NoError(t, err)

	testIface := upcloud.Interface{
		Index: 4,
		IPAddresses: []upcloud.IPAddress{
			{
				Address:  "10.0.0.20",
				Family:   upcloud.IPAddressFamilyIPv4,
				Floating: upcloud.False,
			},
		},
		MAC:               "de:ff:ff:ff:86:cf",
		Network:           "0374ce47-4303-4490-987d-32dc96cfd79b",
		SourceIPFiltering: upcloud.True,
		Type:              upcloud.NetworkTypePrivate,
		Bootable:          upcloud.False,
	}

	assert.Equal(t, testIface, iface)
}

// TestUnmarshalRouters ensures that the unmarshalling of an Routers response
// behaves correctly.
func TestUnmarshalRouters(t *testing.T) {
	t.Parallel()
	originalJSON := `
	  {
		"routers": {
		  "router": [
			{
			  "attached_networks": {
				"network" : [
					{
					   "uuid" : "03206c92-50f2-40b0-ad05-5f02cab2e932"
					},
					{
					   "uuid" : "03d781c3-65e3-4a7a-b6cd-a7ce7e23b8c5"
					},
					{
					   "uuid" : "03854a28-fafe-428f-964f-760cd1b83f1f"
					},
					{
					   "uuid" : "03bae36d-a30a-4640-9a35-f2ccb0e2cddc"
					}
				 ]
			  },
			  "name": "Example router",
			  "type": "normal",
			  "uuid": "04c0df35-2658-4b0c-8ac7-962090f4e92a"
			}
		  ]
		}
	  }
	`

	testRouters := []upcloud.Router{
		{
			AttachedNetworks: []upcloud.RouterNetwork{
				{
					NetworkUUID: "03206c92-50f2-40b0-ad05-5f02cab2e932",
				},
				{
					NetworkUUID: "03d781c3-65e3-4a7a-b6cd-a7ce7e23b8c5",
				},
				{
					NetworkUUID: "03854a28-fafe-428f-964f-760cd1b83f1f",
				},
				{
					NetworkUUID: "03bae36d-a30a-4640-9a35-f2ccb0e2cddc",
				},
			},
			Name: "Example router",
			Type: "normal",
			UUID: "04c0df35-2658-4b0c-8ac7-962090f4e92a",
		},
	}

	routers := upcloud.Routers{}
	err := json.Unmarshal([]byte(originalJSON), &routers)
	assert.NoError(t, err)

	for i, r := range testRouters {
		assert.Equal(t, r, routers.Routers[i])
	}
}

// TestUnmarshalRouters ensures that the unmarshalling of a single Router response
// behaves correctly.
func TestUnmarshalRouter(t *testing.T) {
	t.Parallel()
	originalJSON := `
	  {
		"router": {
		  "attached_networks": {
			"network": [
				{
					"uuid" : "03206c92-50f2-40b0-ad05-5f02cab2e932"
				},
				{
					"uuid" : "03d781c3-65e3-4a7a-b6cd-a7ce7e23b8c5"
				},
				{
					"uuid" : "03854a28-fafe-428f-964f-760cd1b83f1f"
				},
				{
					"uuid" : "03bae36d-a30a-4640-9a35-f2ccb0e2cddc"
				}
			]
		  },
		  "name": "Example router",
		  "type": "normal",
		  "uuid": "04c0df35-2658-4b0c-8ac7-962090f4e92a"
		}
	  }
	`

	var router upcloud.Router
	err := json.Unmarshal([]byte(originalJSON), &router)
	assert.NoError(t, err)

	testRouter := upcloud.Router{
		AttachedNetworks: []upcloud.RouterNetwork{
			{
				NetworkUUID: "03206c92-50f2-40b0-ad05-5f02cab2e932",
			},
			{
				NetworkUUID: "03d781c3-65e3-4a7a-b6cd-a7ce7e23b8c5",
			},
			{
				NetworkUUID: "03854a28-fafe-428f-964f-760cd1b83f1f",
			},
			{
				NetworkUUID: "03bae36d-a30a-4640-9a35-f2ccb0e2cddc",
			},
		},
		Name: "Example router",
		Type: "normal",
		UUID: "04c0df35-2658-4b0c-8ac7-962090f4e92a",
	}

	assert.Equal(t, testRouter, router)
}
