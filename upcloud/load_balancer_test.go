package upcloud

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMarshalLoadBalancer(t *testing.T) {
	lbString := `
	[
		{
			"name": "example-service",
			"network_uuid": "03631160-d57a-4926-ad48-a2f828229dcb",
			"operational_state": "running",
			"plan": "development",
			"configured_status": "started",
			"created_at": "2021-12-07T13:58:30.817272Z",
			"dns_name": "lb-0aff6dac143c43009b33ee2756f6592d.upcloudlb.com",
			"updated_at": "2022-02-11T17:33:59.898714Z",
			"uuid": "0aff6dac-143c-4300-9b33-ee2756f6592d",
			"zone": "fi-hel1",
			"resolvers": [
				{
					"cache_invalid": 30,
					"cache_valid": 180,
					"created_at": "2022-02-11T16:39:55.321306Z",
					"name": "example-resolver",
					"nameservers": [
						"172.16.1.250"
					],
					"retries": 5,
					"timeout": 30,
					"timeout_retry": 10,
					"updated_at": "2022-02-11T17:33:08.490581Z"
				}
			],
			"backends": [
				{
					"name": "example-backend-1",
					"resolver": "example-resolver",
					"updated_at": "2022-02-11T17:33:08.490581Z",
					"created_at": "2021-12-07T13:58:30.817272Z",
					"members": [
						{
							"created_at": "2021-12-07T13:58:30.817272Z",
							"enabled": true,
							"ip": "172.16.1.4",
							"max_sessions": 1000,
							"name": "example-member-1",
							"port": 8000,
							"type": "static",
							"updated_at": "2022-02-11T17:33:08.490581Z",
							"weight": 100
						}
					]
				}
			],
			"frontends": [
				{
					"default_backend": "example-backend-1",
					"mode": "http",
					"name": "example-frontend",
					"port": 443,
					"created_at": "2021-12-07T13:58:30.817272Z",
					"updated_at": "2022-02-11T17:33:08.490581Z",
					"rules": [
						{
							"name": "example-rule-1",
							"priority": 100,
							"created_at": "2021-12-07T13:58:30.817272Z",
							"updated_at": "2022-02-11T17:33:08.490581Z",
							"actions": [
								{
									"type": "use_backend",
									"action_use_backend": {
										"backend": "example-backend-2"
									}
								}
							],
							"matchers": [
								{
									"type": "path",
									"match_path": {
										"method": "exact",
										"value": "/app"
									}
								}
							]
						}
					],
					"tls_configs": [
						{
							"certificate_bundle_uuid": "0aded5c1-c7a3-498a-b9c8-a871611c47a2",
							"created_at": "2022-02-11T17:33:08.490581Z",
							"name": "example-tls-config",
							"updated_at": "2022-02-11T17:33:08.490581Z"
						}
					]
				}
			]
		}
	]`
	lbs := []LoadBalancer{{
		UUID:             "0aff6dac-143c-4300-9b33-ee2756f6592d",
		Name:             "example-service",
		Zone:             "fi-hel1",
		Plan:             "development",
		NetworkUUID:      "03631160-d57a-4926-ad48-a2f828229dcb",
		DNSName:          "lb-0aff6dac143c43009b33ee2756f6592d.upcloudlb.com",
		ConfiguredStatus: LoadBalancerConfiguredStatusStarted,
		OperationalState: LoadBalancerOperationalStateRunning,
		CreatedAt:        timeParse("2021-12-07T13:58:30.817272Z"),
		UpdatedAt:        timeParse("2022-02-11T17:33:59.898714Z"),
		Frontends: []LoadBalancerFrontend{
			{
				Name:           "example-frontend",
				Mode:           LoadBalancerModeHTTP,
				Port:           443,
				DefaultBackend: "example-backend-1",
				Rules: []LoadBalancerRule{
					{
						Name:      "example-rule-1",
						Priority:  100,
						CreatedAt: timeParse("2021-12-07T13:58:30.817272Z"),
						UpdatedAt: timeParse("2022-02-11T17:33:08.490581Z"),
						Matchers: []LoadBalancerMatcher{{
							Type: LoadBalancerMatcherTypePath,
							Path: &LoadBalancerMatcherString{
								Method: LoadBalancerStringMatcherMethodExact,
								Value:  "/app",
							},
						}},
						Actions: []LoadBalancerAction{
							{
								Type: LoadBalancerActionTypeUseBackend,
								UseBackend: &LoadBalancerActionUseBackend{
									Backend: "example-backend-2",
								},
							},
						},
					},
				},
				TLSConfigs: []TLSConfig{{
					Name:                  "example-tls-config",
					CertificateBundleUUID: "0aded5c1-c7a3-498a-b9c8-a871611c47a2",
					CreatedAt:             timeParse("2022-02-11T17:33:08.490581Z"),
					UpdatedAt:             timeParse("2022-02-11T17:33:08.490581Z"),
				}},
				CreatedAt: timeParse("2021-12-07T13:58:30.817272Z"),
				UpdatedAt: timeParse("2022-02-11T17:33:08.490581Z"),
			},
		},
		Backends: []LoadBalancerBackend{
			{
				Name:      "example-backend-1",
				Resolver:  "example-resolver",
				CreatedAt: timeParse("2021-12-07T13:58:30.817272Z"),
				UpdatedAt: timeParse("2022-02-11T17:33:08.490581Z"),
				Members: []LoadBalancerMember{
					{
						Name:        "example-member-1",
						Ip:          "172.16.1.4",
						Port:        8000,
						Weight:      100,
						MaxSessions: 1000,
						Type:        LoadBalancerMemberTypeStatic,
						Enabled:     true,
						CreatedAt:   timeParse("2021-12-07T13:58:30.817272Z"),
						UpdatedAt:   timeParse("2022-02-11T17:33:08.490581Z"),
					},
				},
			},
		},
		Resolvers: []LoadBalancerResolver{
			{
				Name:         "example-resolver",
				Nameservers:  []string{"172.16.1.250"},
				Retries:      5,
				Timeout:      30,
				TimeoutRetry: 10,
				CacheValid:   180,
				CacheInvalid: 30,
				CreatedAt:    timeParse("2022-02-11T16:39:55.321306Z"),
				UpdatedAt:    timeParse("2022-02-11T17:33:08.490581Z"),
			},
		},
	}}
	actual, err := json.Marshal(lbs)
	assert.NoError(t, err)
	assert.JSONEq(t, lbString, string(actual))
	l := []LoadBalancer{}
	err = json.Unmarshal([]byte(lbString), &l)
	assert.NoError(t, err)
	assert.Equal(t, lbs, l)
}
