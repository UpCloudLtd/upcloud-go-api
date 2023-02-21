package request

import (
	"encoding/json"
	"testing"

	"github.com/UpCloudLtd/upcloud-go-api/v5/upcloud"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetNetworkGatewaysRequest(t *testing.T) {
	r := GetNetworkGatewaysRequest{}
	assert.Equal(t, "/gateway", r.RequestURL())

	r = GetNetworkGatewaysRequest{
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
	assert.Equal(t, "/gateway?label=color%3Dgreen&label=size", r.RequestURL())
}

func TestGetNetworkGatewayRequest(t *testing.T) {
	r := GetNetworkGatewayRequest{"id"}
	assert.Equal(t, "/gateway/id", r.RequestURL())
}

func TestCreateNetworkGatewayRequest(t *testing.T) {
	r := CreateNetworkGatewayRequest{
		Name:             "test-create",
		Zone:             "fi-hel1",
		Features:         []upcloud.NetworkGatewayFeature{upcloud.NetworkGatewayFeatureNAT},
		Routers:          []NetworkGatewayRouter{{UUID: "router-uuid"}},
		Labels:           []upcloud.Label{{Key: "test", Value: "Create request"}},
		ConfiguredStatus: upcloud.NetworkGatewayStatusStarted,
	}
	assert.Equal(t, "/gateway", r.RequestURL())
	js, err := json.Marshal(&r)
	require.NoError(t, err)
	assert.JSONEq(t, `
	{
		"name": "test-create",
		"zone": "fi-hel1",
		"features": ["nat"],
		"routers": [{ "uuid": "router-uuid" }],
		"labels": [{ "key": "test", "value": "Create request" }],
		"configured_status": "started"
	}
	`, string(js))
}

func TestModifyNetworkGatewayRequest(t *testing.T) {
	r := ModifyNetworkGatewayRequest{
		UUID:             "id",
		Name:             "test-modify",
		ConfiguredStatus: upcloud.NetworkGatewayStatusStopped,
		Labels:           []upcloud.Label{{Key: "test", Value: "Modify request"}},
	}
	assert.Equal(t, "/gateway/id", r.RequestURL())
	js, err := json.Marshal(&r)
	require.NoError(t, err)
	assert.JSONEq(t, `
	{
		"name": "test-modify",
		"configured_status": "stopped",
		"labels": [{ "key": "test", "value": "Modify request" }]
	}
	`, string(js))
}

func TestDeleteNetworkGatewayRequest(t *testing.T) {
	r := DeleteNetworkGatewayRequest{"id"}
	assert.Equal(t, "/gateway/id", r.RequestURL())
}
