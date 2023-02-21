package upcloud

import (
	"testing"
)

func TestNetworkGateway(t *testing.T) {
	t.Parallel()

	jsonStr := `
	{
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

	gateway := &NetworkGateway{
		ConfiguredStatus: NetworkGatewayStatusStarted,
		CreatedAt:        timeParse("2022-12-01T09:04:08.529138Z"),
		Features: []NetworkGatewayFeature{
			NetworkGatewayFeatureNAT,
		},
		Name:             "example-gateway",
		OperationalState: "running",
		Routers: []NetworkGatewayRouter{
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

	testJSON(t, &NetworkGateway{}, gateway, jsonStr)
}
