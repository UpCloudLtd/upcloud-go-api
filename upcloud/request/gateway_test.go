package request

import (
	"encoding/json"
	"testing"

	"github.com/UpCloudLtd/upcloud-go-api/v5/upcloud"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetGatewaysRequest(t *testing.T) {
	r := GetGatewaysRequest{}
	assert.Equal(t, "/gateway", r.RequestURL())

	r = GetGatewaysRequest{
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

func TestGetGatewayRequest(t *testing.T) {
	r := GetGatewayRequest{"id"}
	assert.Equal(t, "/gateway/id", r.RequestURL())
}

func TestCreateGatewayRequest(t *testing.T) {
	r := CreateGatewayRequest{
		Name:             "test-create",
		Zone:             "fi-hel1",
		Features:         []upcloud.GatewayFeature{upcloud.GatewayFeatureNAT},
		Routers:          []GatewayRouter{{UUID: "router-uuid"}},
		Labels:           []upcloud.Label{{Key: "test", Value: "Create request"}},
		ConfiguredStatus: upcloud.GatewayStatusStarted,
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

func TestModifyGatewayRequest(t *testing.T) {
	r := ModifyGatewayRequest{
		UUID:             "id",
		Name:             "test-modify",
		ConfiguredStatus: upcloud.GatewayStatusStopped,
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

func TestDeleteGatewayRequest(t *testing.T) {
	r := DeleteGatewayRequest{"id"}
	assert.Equal(t, "/gateway/id", r.RequestURL())
}
