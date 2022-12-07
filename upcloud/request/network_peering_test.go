package request

import (
	"encoding/json"
	"testing"

	"github.com/UpCloudLtd/upcloud-go-api/v5/upcloud"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
}

func TestDeleteNetworkPeeringRequest(t *testing.T) {
	r := DeleteNetworkPeeringRequest{"id"}
	assert.Equal(t, "/network-peering/id", r.RequestURL())
}
