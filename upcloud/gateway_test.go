package upcloud

import (
	"testing"
)

func TestGateway(t *testing.T) {
	t.Parallel()

	jsonStr := `
	{
		"addresses": [
			{
				"address": "192.0.2.96",
				"name": "public-ip-1"
			}
		],
		"configured_status": "started",
		"created_at": "2022-12-01T09:04:08.529138Z",
		"features": [
			"nat"
		],
		"name": "example-gateway",
		"operational_state": "running",
		"routers": [
			{
				"created_at": "2022-12-01T09:04:08.529138Z",
				"uuid": "0485d477-8d8f-4c97-9bef-731933187538"
			}
		],
		"labels": [
			{
				"key":"env",
				"value":"testing"
			}
		],
		"updated_at": "2022-12-01T19:04:08.529138Z",
		"uuid": "10c153e0-12e4-4dea-8748-4f34850ff76d",
		"zone": "fi-hel1"
	}
	`

	gateway := &Gateway{
		Addresses: []GatewayAddress{{
			Address: "192.0.2.96",
			Name:    "public-ip-1",
		}},
		ConfiguredStatus: GatewayConfiguredStatusStarted,
		CreatedAt:        timeParse("2022-12-01T09:04:08.529138Z"),
		Features: []GatewayFeature{
			GatewayFeatureNAT,
		},
		Name:             "example-gateway",
		OperationalState: "running",
		Routers: []GatewayRouter{
			{
				CreatedAt: timeParse("2022-12-01T09:04:08.529138Z"),
				UUID:      "0485d477-8d8f-4c97-9bef-731933187538",
			},
		},
		Labels: []Label{
			{Key: "env", Value: "testing"},
		},
		UpdatedAt: timeParse("2022-12-01T19:04:08.529138Z"),
		UUID:      "10c153e0-12e4-4dea-8748-4f34850ff76d",
		Zone:      "fi-hel1",
	}

	testJSON(t, &Gateway{}, gateway, jsonStr)
}

func TestGatewayPlan(t *testing.T) {
	t.Parallel()

	jsonStr := `
	{
		"name": "advanced",
		"per_gateway_bandwidth_mbps": 10000,
		"per_gateway_max_connections": 100000,
		"server_number": 2,
		"supported_features": [
			"nat",
			"vpn"
		],
		"vpn_tunnel_amount": 10
	}
	`

	plan := &GatewayPlan{
		Name:                     "advanced",
		PerGatewayBandwidthMbps:  10000,
		PerGatewayMaxConnections: 100000,
		ServerNumber:             2,
		SupportedFeatures:        []GatewayFeature{GatewayFeatureNAT, GatewayFeatureVPN},
		VPNTunnelAmount:          10,
	}

	testJSON(t, &GatewayPlan{}, plan, jsonStr)
}
