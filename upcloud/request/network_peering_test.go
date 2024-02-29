package request

import (
	"encoding/json"
	"testing"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetNetworkPeeringsRequest(t *testing.T) {
	r := GetNetworkPeeringsRequest{}
	assert.Equal(t, networkPeeringBaseURL, r.RequestURL())

	r = GetNetworkPeeringsRequest{
		Filters: []QueryFilter{
			FilterLabel{
				Label: upcloud.Label{
					Key:   "color",
					Value: "green",
				},
			},
			FilterLabelKey{Key: "size"},
		},
	}
	assert.Equal(t, networkPeeringBaseURL+"?label=color%3Dgreen&label=size", r.RequestURL())
}

func TestGetNetworkPeeringRequest(t *testing.T) {
	r := GetNetworkPeeringRequest{"id"}
	assert.Equal(t, "/network-peering/id", r.RequestURL())
}

func TestCreateNetworkPeeringRequest(t *testing.T) {
	r := CreateNetworkPeeringRequest{
		Name:             "Peering A->B",
		ConfiguredStatus: upcloud.NetworkPeeringConfiguredStatusActive,
		Network: NetworkPeeringNetwork{
			UUID: "03126dc1-a69f-4bc2-8b24-e31c22d64712",
		},
		PeerNetwork: NetworkPeeringNetwork{
			UUID: "03585987-bf7d-4544-8e9b-5a1b4d74a333",
		},
	}
	assert.Equal(t, "/network-peering", r.RequestURL())
	js, err := json.Marshal(&r)
	require.NoError(t, err)
	assert.JSONEq(t, `
	{
		"network_peering": {
		  "configured_status": "active",
		  "name": "Peering A->B",
		  "network": {
			"uuid": "03126dc1-a69f-4bc2-8b24-e31c22d64712"
		  },
		  "peer_network": {
			"uuid": "03585987-bf7d-4544-8e9b-5a1b4d74a333"
		  }
		}
	  }
	`, string(js))
}

func TestCreateNetworkPeeringLabelsRequest(t *testing.T) {
	r := CreateNetworkPeeringRequest{
		Name:             "Peering A->B",
		ConfiguredStatus: upcloud.NetworkPeeringConfiguredStatusActive,
		Network: NetworkPeeringNetwork{
			UUID: "03126dc1-a69f-4bc2-8b24-e31c22d64712",
		},
		PeerNetwork: NetworkPeeringNetwork{
			UUID: "03585987-bf7d-4544-8e9b-5a1b4d74a333",
		},
		Labels: []upcloud.Label{
			{
				Key:   "managedBy",
				Value: "upcloud-go-sdk-unit-test",
			},
		},
	}
	assert.Equal(t, "/network-peering", r.RequestURL())
	js, err := json.Marshal(&r)
	require.NoError(t, err)
	assert.JSONEq(t, `
	{
		"network_peering": {
		  "configured_status": "active",
		  "name": "Peering A->B",
		  "network": {
			"uuid": "03126dc1-a69f-4bc2-8b24-e31c22d64712"
		  },
		  "peer_network": {
			"uuid": "03585987-bf7d-4544-8e9b-5a1b4d74a333"
		  },
		  "labels": [
			{
				"key": "managedBy",
				"value": "upcloud-go-sdk-unit-test"
			}
		  ]
		}
	  }
	`, string(js))
}

func TestModifyNetworkPeeringRequest(t *testing.T) {
	r := ModifyNetworkPeeringRequest{
		UUID: "id",
		NetworkPeering: ModifyNetworkPeering{
			Name:             "Peering A->B modified",
			ConfiguredStatus: upcloud.NetworkPeeringConfiguredStatusDisabled,
		},
	}
	assert.Equal(t, "/network-peering/id", r.RequestURL())
	js, err := json.Marshal(&r)
	require.NoError(t, err)
	assert.JSONEq(t, `
	{
		"network_peering": {
		  "configured_status": "disabled",
		  "name": "Peering A->B modified"
		}
	}
	`, string(js))

	r = ModifyNetworkPeeringRequest{
		UUID: "id",
		NetworkPeering: ModifyNetworkPeering{
			ConfiguredStatus: upcloud.NetworkPeeringConfiguredStatusActive,
		},
	}
	js, err = json.Marshal(&r)
	require.NoError(t, err)
	assert.JSONEq(t, `
	{
		"network_peering": {
		  "configured_status": "active"
		}
	}
	`, string(js))

	r = ModifyNetworkPeeringRequest{
		UUID: "id",
		NetworkPeering: ModifyNetworkPeering{
			Name: "Peering A->B modified",
		},
	}
	js, err = json.Marshal(&r)
	require.NoError(t, err)
	assert.JSONEq(t, `
	{
		"network_peering": {
			"name": "Peering A->B modified"
		}
	}
	`, string(js))
}

func TestModifyNetworkPeeringLabelsRequest(t *testing.T) {
	r := ModifyNetworkPeeringRequest{
		UUID:           "id",
		NetworkPeering: ModifyNetworkPeering{},
	}
	assert.Equal(t, "/network-peering/id", r.RequestURL())
	js, err := json.Marshal(&r)
	require.NoError(t, err)
	assert.JSONEq(t, `{"network_peering":{}}`, string(js))

	r = ModifyNetworkPeeringRequest{
		UUID: "id",
		NetworkPeering: ModifyNetworkPeering{
			Labels: &[]upcloud.Label{
				{
					Key:   "managedBy",
					Value: "upcloud-go-sdk-unit-test",
				},
			},
		},
	}
	js, err = json.Marshal(&r)
	require.NoError(t, err)
	assert.JSONEq(t, `
	{
		"network_peering": {
			"labels": [
				{
					"key": "managedBy",
					"value": "upcloud-go-sdk-unit-test"
				}
			]
		}
	}
	`, string(js))
}

func TestDeleteNetworkPeeringRequest(t *testing.T) {
	r := DeleteNetworkPeeringRequest{"id"}
	assert.Equal(t, "/network-peering/id", r.RequestURL())
}
