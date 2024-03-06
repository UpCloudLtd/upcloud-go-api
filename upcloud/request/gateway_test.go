package request

import (
	"encoding/json"
	"testing"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
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
	t.Parallel()

	r := GetGatewayRequest{UUID: "fake"}
	assert.Equal(t, gatewayBaseURL+"/fake", r.RequestURL())
}

func TestDeleteGatewayRequest(t *testing.T) {
	t.Parallel()

	r := DeleteGatewayRequest{UUID: "fake"}
	assert.Equal(t, gatewayBaseURL+"/fake", r.RequestURL())
}

func TestCreateGatewayRequest(t *testing.T) {
	t.Parallel()

	const want string = `
	{
		"name": "example-gateway",
		"zone": "fi-hel1",
		"features": [
		  "nat"
		],
		"routers": [
		  {
			"uuid": "0485d477-8d8f-4c97-9bef-731933187538"
		  }
		],
		"configured_status": "started",
		"plan": "random-plan",
		"addresses": [
		  {
			"name": "my-public-ip"
		  }
		],
		"connections": [
			{
				"name": "example-connection",
      			"type": "ipsec",
      			"local_routes": [
        			{
						"name": "upcloud-example-route",
						"type": "static",
						"static_network": "10.0.0.0/24"
					}
				],
				"remote_routes": [
					{
						"name": "remote-example-route",
						"type": "static",
						"static_network": "10.0.1.0/24"
					}
				],
				"tunnels": [
					{
						"name": "example-tunnel-1",
						"local_address": {
							"name": "public-ip-1"
						},
						"remote_address": {
							"address": "100.10.0.111"
						},
						"ipsec": {
							"authentication": {
								"authentication": "psk",
								"psk": "pskpsksk"
							}
						}
					}
				]

			}
		]
	}
	`
	r := CreateGatewayRequest{
		Name:     "example-gateway",
		Zone:     "fi-hel1",
		Features: []upcloud.GatewayFeature{upcloud.GatewayFeatureNAT},
		Routers: []GatewayRouter{
			{
				UUID: "0485d477-8d8f-4c97-9bef-731933187538",
			},
		},
		ConfiguredStatus: upcloud.GatewayConfiguredStatusStarted,
		Plan:             "random-plan",
		Addresses: []upcloud.GatewayAddress{
			{
				Name: "my-public-ip",
			},
		},
		Connections: []GatewayConnection{
			{
				Name: "example-connection",
				Type: upcloud.GatewayConnectionTypeIPSec,
				LocalRoutes: []upcloud.GatewayRoute{
					{
						Name:          "upcloud-example-route",
						Type:          upcloud.GatewayRouteTypeStatic,
						StaticNetwork: "10.0.0.0/24",
					},
				},
				RemoteRoutes: []upcloud.GatewayRoute{
					{
						Name:          "remote-example-route",
						Type:          upcloud.GatewayRouteTypeStatic,
						StaticNetwork: "10.0.1.0/24",
					},
				},
				Tunnels: []GatewayTunnel{
					{
						Name: "example-tunnel-1",
						LocalAddress: upcloud.GatewayTunnelLocalAddress{
							Name: "public-ip-1",
						},
						RemoteAddress: upcloud.GatewayTunnelRemoteAddress{
							Address: "100.10.0.111",
						},
						IPSec: upcloud.GatewayTunnelIPSec{
							Authentication: upcloud.GatewayTunnelIPSecAuth{
								Authentication: upcloud.GatewayTunnelIPSecAuthTypePSK,
								PSK:            "pskpsksk",
							},
						},
					},
				},
			},
		},
	}

	got, err := json.Marshal(&r)
	require.NoError(t, err)
	assert.Equal(t, gatewayBaseURL, r.RequestURL())
	assert.JSONEq(t, want, string(got))
}

func TestModifyGatewayRequest(t *testing.T) {
	t.Parallel()

	want := `
	{
		"name": "example-gateway",
		"configured_status": "started"
	}
	`
	r := ModifyGatewayRequest{
		UUID:             "fake",
		Name:             "example-gateway",
		ConfiguredStatus: upcloud.GatewayConfiguredStatusStarted,
	}
	got, err := json.Marshal(&r)
	require.NoError(t, err)
	assert.JSONEq(t, want, string(got))

	want = `
	{
		"name": "example-gateway"
	}
	`
	r = ModifyGatewayRequest{
		UUID: "fake",
		Name: "example-gateway",
	}
	got, err = json.Marshal(&r)
	require.NoError(t, err)
	assert.JSONEq(t, want, string(got))

	want = `
	{
		"configured_status": "started"
	}
	`
	r = ModifyGatewayRequest{
		UUID:             "fake",
		ConfiguredStatus: upcloud.GatewayConfiguredStatusStarted,
	}
	got, err = json.Marshal(&r)
	require.NoError(t, err)
	assert.JSONEq(t, want, string(got))

	assert.Equal(t, gatewayBaseURL+"/fake", r.RequestURL())
}
