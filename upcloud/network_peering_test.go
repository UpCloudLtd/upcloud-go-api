package upcloud

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNetworkPeering(t *testing.T) {
	t.Parallel()

	got := make(NetworkPeerings, 0)
	require.NoError(t, json.Unmarshal([]byte(`
	{
		"network_peerings": {
		  "network_peering": [
			{
			  "configured_status": "active",
			  "name": "Peering A->B",
			  "network": {
				"ip_networks": {
				  "ip_network": [
					{
					  "address": "192.168.0.0/24",
					  "family": "IPv4"
					},
					{
					  "address": "fc02:c4f3::/64",
					  "family": "IPv6"
					}
				  ]
				},
				"uuid": "03126dc1-a69f-4bc2-8b24-e31c22d64712"
			  },
			  "peer_network": {
				"ip_networks": {
				  "ip_network": [
					{
					  "address": "192.168.99.0/24",
					  "family": "IPv4"
					},
					{
					  "address": "fc02:c4f3:99::/64",
					  "family": "IPv6"
					}
				  ]
				},
				"uuid": "03585987-bf7d-4544-8e9b-5a1b4d74a333"
			  },
			  "state": "active",
			  "uuid": "0f7984bc-5d72-4aaf-b587-90e6a8f32efc",
			  "labels": [
				{
					"key": "managedBy",
					"value": "upcloud-go-sdk-unit-test"
				}
			  ]
			}
		  ]
		}
	  }
	`), &got))

	want := NetworkPeerings{NetworkPeering{
		UUID:             "0f7984bc-5d72-4aaf-b587-90e6a8f32efc",
		ConfiguredStatus: NetworkPeeringConfiguredStatusActive,
		Name:             "Peering A->B",
		Network: NetworkPeeringNetwork{
			UUID: "03126dc1-a69f-4bc2-8b24-e31c22d64712",
			IPNetworks: []NetworkPeeringIPNetwork{
				{
					Address: "192.168.0.0/24",
					Family:  NetworkPeeringIPNetworkFamilyIPv4,
				},
				{
					Address: "fc02:c4f3::/64",
					Family:  NetworkPeeringIPNetworkFamilyIPv6,
				},
			},
		},
		PeerNetwork: NetworkPeeringNetwork{
			UUID: "03585987-bf7d-4544-8e9b-5a1b4d74a333",
			IPNetworks: []NetworkPeeringIPNetwork{
				{
					Address: "192.168.99.0/24",
					Family:  NetworkPeeringIPNetworkFamilyIPv4,
				},
				{
					Address: "fc02:c4f3:99::/64",
					Family:  NetworkPeeringIPNetworkFamilyIPv6,
				},
			},
		},
		State: NetworkPeeringStateActive,
		Labels: []Label{
			{
				Key:   "managedBy",
				Value: "upcloud-go-sdk-unit-test",
			},
		},
	}}
	assert.Equal(t, want, got)
}
