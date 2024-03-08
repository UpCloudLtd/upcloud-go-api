package upcloud

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
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
			"nat",
			"vpn"
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
		"zone": "fi-hel1",
		"plan": "advanced",
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
								"authentication": "psk"
							},
							"child_rekey_time": 1440,
							"dpd_delay": 30,
							"dpd_timeout": 120,
							"ike_lifetime": 86400,
							"phase1_algorithms": ["aes128gcm128", "aes256gcm128"],
							"phase1_dh_group_numbers": [14, 16, 18, 19, 20, 21],
							"phase1_integrity_algorithms": ["aes128gmac", "aes256gmac", "sha256", "sha384", "sha512"],
							"phase2_algorithms": ["aes128gcm128", "aes256gcm128"],
							"phase2_dh_group_numbers": [14, 16, 18, 19, 20, 21],
							"phase2_integrity_algorithms": ["aes128gmac", "aes256gmac", "sha256", "sha384", "sha512"],
							"rekey_time": 14400
						},
						"operational_state": "established",
						"created_at": "2022-12-01T09:04:08.529138Z",
						"updated_at": "2022-12-01T09:04:08.529138Z"
				  	}
				],
				"created_at": "2022-12-01T09:04:08.529138Z",
				"updated_at": "2022-12-01T09:04:08.529138Z"
			}
		]
	}
	`

	timestamp, err := time.Parse(time.RFC3339, "2022-12-01T09:04:08.529138Z")
	require.NoError(t, err)

	gateway := &Gateway{
		Addresses: []GatewayAddress{{
			Address: "192.0.2.96",
			Name:    "public-ip-1",
		}},
		ConfiguredStatus: GatewayConfiguredStatusStarted,
		CreatedAt:        timeParse("2022-12-01T09:04:08.529138Z"),
		Features: []GatewayFeature{
			GatewayFeatureNAT,
			GatewayFeatureVPN,
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
		Plan:      "advanced",
		Connections: []GatewayConnection{
			{
				Name: "example-connection",
				Type: GatewayConnectionTypeIPSec,
				LocalRoutes: []GatewayRoute{
					{
						Name:          "upcloud-example-route",
						Type:          GatewayRouteTypeStatic,
						StaticNetwork: "10.0.0.0/24",
					},
				},
				RemoteRoutes: []GatewayRoute{
					{
						Name:          "remote-example-route",
						Type:          GatewayRouteTypeStatic,
						StaticNetwork: "10.0.1.0/24",
					},
				},
				Tunnels: []GatewayTunnel{
					{
						Name: "example-tunnel-1",
						LocalAddress: GatewayTunnelLocalAddress{
							Name: "public-ip-1",
						},
						RemoteAddress: GatewayTunnelRemoteAddress{
							Address: "100.10.0.111",
						},
						IPSec: GatewayTunnelIPSec{
							Authentication: GatewayTunnelIPSecAuth{
								Authentication: GatewayTunnelIPSecAuthTypePSK,
							},
							ChildRekeyTime: 1440,
							DPDDelay:       30,
							DPDTimeout:     120,
							IKELifetime:    86400,
							Phase1Algortihms: []GatewayIPSecAlgorithm{
								GatewayIPSecAlgorithm_aes128gcm128,
								GatewayIPSecAlgorithm_aes256gcm128,
							},
							Phase1DHGroupNumbers: []int{14, 16, 18, 19, 20, 21},
							Phase1IntegrityAlgorithms: []GatewayIPSecIntegrityAlgorithm{
								GatewayIPSecIntegrityAlgorithm_aes128gmac,
								GatewayIPSecIntegrityAlgorithm_aes256gmac,
								GatewayIPSecIntegrityAlgorithm_sha256,
								GatewayIPSecIntegrityAlgorithm_sha384,
								GatewayIPSecIntegrityAlgorithm_sha512,
							},
							Phase2Algortihms: []GatewayIPSecAlgorithm{
								GatewayIPSecAlgorithm_aes128gcm128,
								GatewayIPSecAlgorithm_aes256gcm128,
							},
							Phase2DHGroupNumbers: []int{14, 16, 18, 19, 20, 21},
							Phase2IntegrityAlgorithms: []GatewayIPSecIntegrityAlgorithm{
								GatewayIPSecIntegrityAlgorithm_aes128gmac,
								GatewayIPSecIntegrityAlgorithm_aes256gmac,
								GatewayIPSecIntegrityAlgorithm_sha256,
								GatewayIPSecIntegrityAlgorithm_sha384,
								GatewayIPSecIntegrityAlgorithm_sha512,
							},
							RekeyTime: 14400,
						},
						OperationalState: GatewayTunnelOperationalStateEstabilished,
						CreatedAt:        timestamp,
						UpdatedAt:        timestamp,
					},
				},
				CreatedAt: timestamp,
				UpdatedAt: timestamp,
			},
		},
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
