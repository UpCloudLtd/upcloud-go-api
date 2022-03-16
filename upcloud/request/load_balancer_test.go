package request

import (
	"encoding/json"
	"testing"

	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateLoadBalancerRequest(t *testing.T) {
	expected := `
	{
		"name": "example-service",
		"plan": "development",
		"zone": "fi-hel1",
		"network_uuid": "03631160-d57a-4926-ad48-a2f828229dcb",
		"configured_status": "started",
		"frontends": [
			{
				"name": "example-frontend",
				"mode": "http",
				"port": 443,
				"default_backend": "example-backend-1"
			}
		],
		"backends": [
			{
				"name": "example-backend-1",
				"members": [
					{
						"name": "example-member-1",
						"ip": "172.16.1.4",
						"port": 8000,
						"type": "static",
						"weight": 100,
						"max_sessions": 1000,
						"enabled": true
					}
				]
			}
		],
		"resolvers": [
			{
				"name": "example-resolver",
				"nameservers": [
					"172.16.1.4:53"
				],
				"retries": 5,
				"timeout": 30,
				"timeout_retry": 10,
				"cache_valid": 180,
				"cache_invalid": 10
			}
		]
	}
	`
	r := CreateLoadBalancerRequest{
		Name:             "example-service",
		Plan:             "development",
		Zone:             "fi-hel1",
		NetworkUUID:      "03631160-d57a-4926-ad48-a2f828229dcb",
		ConfiguredStatus: upcloud.LoadBalancerConfiguredStatusStarted,
		Frontends: []LoadBalancerFrontend{{
			Name:           "example-frontend",
			Mode:           upcloud.LoadBalancerModeHTTP,
			Port:           443,
			DefaultBackend: "example-backend-1",
		}},
		Backends: []LoadBalancerBackend{{
			Name: "example-backend-1",
			Members: []LoadBalancerBackendMember{{
				Name:        "example-member-1",
				Weight:      100,
				MaxSessions: 1000,
				Type:        upcloud.LoadBalancerBackendMemberTypeStatic,
				IP:          "172.16.1.4",
				Port:        8000,
				Enabled:     true,
			}},
		}},
		Resolvers: []LoadBalancerResolver{{
			Name:         "example-resolver",
			Nameservers:  []string{"172.16.1.4:53"},
			Retries:      5,
			Timeout:      30,
			TimeoutRetry: 10,
			CacheValid:   180,
			CacheInvalid: 10,
		}},
	}
	actual, err := json.Marshal(&r)
	assert.NoError(t, err)
	assert.JSONEq(t, expected, string(actual))
	assert.Equal(t, "/loadbalancer", r.RequestURL())
}

func TestCreateLoadBalancerBackendRequest(t *testing.T) {
	r := CreateLoadBalancerBackendRequest{
		ServiceUUID: "lb",
		Backend: LoadBalancerBackend{
			Name:     "sesese",
			Members:  []LoadBalancerBackendMember{},
			Resolver: "testresolver",
		},
	}

	expectedJson := `
	{
		"name": "sesese",
		"resolver": "testresolver",
		"members": []
	}`

	actualJson, err := json.Marshal(&r)

	require.NoError(t, err)
	assert.Exactly(t, "/loadbalancer/lb/backends", r.RequestURL())
	assert.JSONEq(t, expectedJson, string(actualJson))
}

func TestGetLoadBalancerBackendsRequest(t *testing.T) {
	r := GetLoadBalancerBackendsRequest{
		ServiceUUID: "lb",
	}

	expectedJson := "{}"
	actualJson, err := json.Marshal(&r)

	require.NoError(t, err)
	assert.Exactly(t, "/loadbalancer/lb/backends", r.RequestURL())
	assert.JSONEq(t, expectedJson, string(actualJson))
}

func TestGetLoadBalancerBackendRequest(t *testing.T) {
	r := GetLoadBalancerBackendRequest{
		ServiceUUID: "lb",
		Name:        "be",
	}

	expectedJson := "{}"
	actualJson, err := json.Marshal(&r)

	require.NoError(t, err)
	assert.Exactly(t, "/loadbalancer/lb/backends/be", r.RequestURL())
	assert.JSONEq(t, expectedJson, string(actualJson))
}

func TestModifyLoadBalancerBackendRequest(t *testing.T) {
	r := ModifyLoadBalancerBackendRequest{
		ServiceUUID: "lb",
		Name:        "be",
		Backend: ModifyLoadBalancerBackend{
			Name:     "newnew",
			Resolver: "newresolver",
		},
	}

	expectedJson := `
	{
		"name": "newnew",
		"resolver": "newresolver"	
	}`

	actualJson, err := json.Marshal(&r)

	require.NoError(t, err)
	assert.Exactly(t, "/loadbalancer/lb/backends/be", r.RequestURL())
	assert.JSONEq(t, expectedJson, string(actualJson))
}

func TestDeleteLoadBalancerBackendRequest(t *testing.T) {
	r := DeleteLoadBalancerBackendRequest{
		ServiceUUID: "lb",
		Name:        "be",
	}

	expectedJson := "{}"
	actualJson, err := json.Marshal(&r)

	require.NoError(t, err)
	assert.Exactly(t, "/loadbalancer/lb/backends/be", r.RequestURL())
	assert.JSONEq(t, expectedJson, string(actualJson))
}

func TestCreateLoadBalancerBackendMember(t *testing.T) {
	r := CreateLoadBalancerBackendMemberRequest{
		ServiceUUID: "lb",
		BackendName: "be",
		Member: LoadBalancerBackendMember{
			Name:        "mem",
			Weight:      100,
			MaxSessions: 5,
			Enabled:     true,
			Type:        "static",
			IP:          "10.0.0.1",
			Port:        80,
		},
	}

	expectedJson := `
	{
		"name": "mem",
		"weight": 100,
		"max_sessions": 5,
		"enabled": true,
		"type": "static",
		"ip": "10.0.0.1",
		"port": 80
	}`

	actualJson, err := json.Marshal(&r)

	require.NoError(t, err)
	assert.Exactly(t, "/loadbalancer/lb/backends/be/members", r.RequestURL())
	assert.JSONEq(t, expectedJson, string(actualJson))
}

func TestGetLoadBalancerBackendMembersRequest(t *testing.T) {
	r := GetLoadBalancerBackendMembersRequest{
		ServiceUUID: "lb",
		BackendName: "be",
	}

	expectedJson := "{}"
	actualJson, err := json.Marshal(&r)

	require.NoError(t, err)
	assert.Exactly(t, "/loadbalancer/lb/backends/be/members", r.RequestURL())
	assert.JSONEq(t, expectedJson, string(actualJson))
}

func TestGetLoadBalancerBackendMemberRequest(t *testing.T) {
	r := GetLoadBalancerBackendMemberRequest{
		ServiceUUID: "lb",
		BackendName: "be",
		Name:        "mem",
	}

	expectedJson := "{}"
	actualJson, err := json.Marshal(&r)

	require.NoError(t, err)
	assert.Exactly(t, "/loadbalancer/lb/backends/be/members/mem", r.RequestURL())
	assert.JSONEq(t, expectedJson, string(actualJson))
}

func TestModifyLoadBalancerBackendMemberRequest(t *testing.T) {
	r := ModifyLoadBalancerBackendMemberRequest{
		ServiceUUID: "lb",
		BackendName: "be",
		Name:        "mem",
		Member: LoadBalancerBackendMember{
			Name:        "newmem",
			Weight:      100,
			MaxSessions: 5,
			Enabled:     true,
			Type:        "static",
			IP:          "10.0.0.1",
			Port:        80,
		},
	}

	expectedJson := `
	{
		"name": "newmem",
		"weight": 100,
		"max_sessions": 5,
		"enabled": true,
		"type": "static",
		"ip": "10.0.0.1",
		"port": 80
	}`

	actualJson, err := json.Marshal(&r)

	require.NoError(t, err)
	assert.Exactly(t, "/loadbalancer/lb/backends/be/members/mem", r.RequestURL())
	assert.JSONEq(t, expectedJson, string(actualJson))
}

func TestDeleteLoadBalancerBackendMemberRequest(t *testing.T) {
	r := DeleteLoadBalancerBackendMemberRequest{
		ServiceUUID: "lb",
		BackendName: "be",
		Name:        "mem",
	}

	expectedJson := "{}"
	actualJson, err := json.Marshal(&r)

	require.NoError(t, err)
	assert.Exactly(t, "/loadbalancer/lb/backends/be/members/mem", r.RequestURL())
	assert.JSONEq(t, expectedJson, string(actualJson))
}

func TestCreateLoadBalancerResolverRequest(t *testing.T) {
	r := CreateLoadBalancerResolverRequest{
		ServiceUUID: "service-uuid",
		Resolver: LoadBalancerResolver{
			Name:         "testname",
			Nameservers:  []string{"10.0.0.0", "10.0.0.1"},
			Retries:      5,
			TimeoutRetry: 10,
			Timeout:      20,
			CacheValid:   123,
			CacheInvalid: 321},
	}

	expectedJson := `
	{
		"name":"testname",
		"nameservers":["10.0.0.0","10.0.0.1"],
		"retries":5,
		"timeout":20,
		"timeout_retry":10,
		"cache_valid":123,
		"cache_invalid":321
	}`

	actualJson, err := json.Marshal(&r)

	require.NoError(t, err)
	assert.EqualValues(t, "/loadbalancer/service-uuid/resolvers", r.RequestURL())
	assert.JSONEq(t, expectedJson, string(actualJson))
}

func TestGetLoadBalancerResolversRequest(t *testing.T) {
	r := GetLoadBalancerResolversRequest{
		ServiceUUID: "service-uuid",
	}

	expectedJson := "{}"
	actualJson, err := json.Marshal(&r)

	require.NoError(t, err)
	assert.EqualValues(t, "/loadbalancer/service-uuid/resolvers", r.RequestURL())
	assert.JSONEq(t, expectedJson, string(actualJson))
}

func TestGetLoadBalancerResolverRequest(t *testing.T) {
	r := GetLoadBalancerResolverRequest{
		ServiceUUID: "service-uuid",
		Name:        "sesese",
	}

	expectedJson := "{}"
	actualJson, err := json.Marshal(&r)

	require.NoError(t, err)
	assert.EqualValues(t, "/loadbalancer/service-uuid/resolvers/sesese", r.RequestURL())
	assert.JSONEq(t, expectedJson, string(actualJson))
}

func TestModifyLoadBalancerResolverRequest(t *testing.T) {
	r := ModifyLoadBalancerRevolverRequest{
		ServiceUUID: "service-uuid",
		Name:        "sesese",
		Resolver: LoadBalancerResolver{
			Name:         "testname",
			Nameservers:  []string{"10.0.0.0", "10.0.0.1"},
			Retries:      5,
			TimeoutRetry: 10,
			Timeout:      20,
			CacheValid:   123,
			CacheInvalid: 321},
	}

	expectedJson := `
	{
		"name":"testname",
		"nameservers":["10.0.0.0","10.0.0.1"],
		"retries":5,
		"timeout":20,
		"timeout_retry":10,
		"cache_valid":123,
		"cache_invalid":321
	}`

	actualJson, err := json.Marshal(&r)

	require.NoError(t, err)
	assert.EqualValues(t, "/loadbalancer/service-uuid/resolvers/sesese", r.RequestURL())
	assert.JSONEq(t, expectedJson, string(actualJson))
}

func TestDeleteLoadBalancerResolverRequest(t *testing.T) {
	r := DeleteLoadBalancerResolverRequest{
		ServiceUUID: "service-uuid",
		Name:        "sesese",
	}

	expectedJson := "{}"
	actualJson, err := json.Marshal(&r)

	require.NoError(t, err)
	assert.EqualValues(t, "/loadbalancer/service-uuid/resolvers/sesese", r.RequestURL())
	assert.JSONEq(t, expectedJson, string(actualJson))
}

func TestGetLoadBalancerPlansRequest(t *testing.T) {
	r := GetLoadBalancerPlansRequest{}
	assert.Equal(t, "/loadbalancer/plans", r.RequestURL())
}

func TestGetLoadBalancerFrontendsRequest(t *testing.T) {
	r := GetLoadBalancerFrontendsRequest{"sid"}
	assert.Equal(t, "/loadbalancer/sid/frontends", r.RequestURL())
}

func TestGetLoadBalancerFrontendRequest(t *testing.T) {
	r := GetLoadBalancerFrontendRequest{
		ServiceUUID: "sid",
		Name:        "be_name",
	}
	assert.Equal(t, "/loadbalancer/sid/frontends/be_name", r.RequestURL())
}

func TestCreateLoadBalancerFrontendRequest(t *testing.T) {
	expected := `
	{
		"name": "example-frontend",
		"mode": "http",
		"port": 443,
		"default_backend": "example-backend",
		"rules": [
			{
				"name": "example-rule-1",
				"priority": 100,
				"matchers": [
					{
						"type": "path",
						"match_path": {
							"method": "exact",
							"value": "/app"
						}
					}
				],
				"actions": [
					{
						"type": "use_backend",
						"action_use_backend": {
							"backend": "example-backend-2"
						}
					}
				]
			}
		],
		"tls_configs": [
			{
				"name": "example-tls-config",
				"certificate_bundle_uuid": "0aded5c1-c7a3-498a-b9c8-a871611c47a2"
			}
		]
	}
	`
	r := CreateLoadBalancerFrontendRequest{
		ServiceUUID: "sid",
		Frontend: LoadBalancerFrontend{
			Name:           "example-frontend",
			Mode:           upcloud.LoadBalancerModeHTTP,
			Port:           443,
			DefaultBackend: "example-backend",
			Rules: []LoadBalancerFrontendRule{{
				Name:     "example-rule-1",
				Priority: 100,
				Matchers: []upcloud.LoadBalancerMatcher{{
					Type: upcloud.LoadBalancerMatcherTypePath,
					Path: &upcloud.LoadBalancerMatcherString{
						Method: upcloud.LoadBalancerStringMatcherMethodExact,
						Value:  "/app",
					},
				}},
				Actions: []upcloud.LoadBalancerAction{{
					Type: upcloud.LoadBalancerActionTypeUseBackend,
					UseBackend: &upcloud.LoadBalancerActionUseBackend{
						Backend: "example-backend-2",
					},
				}},
			}},
			TLSConfigs: []LoadBalancerTLSConfig{{
				Name:                  "example-tls-config",
				CertificateBundleUUID: "0aded5c1-c7a3-498a-b9c8-a871611c47a2",
			}},
		},
	}
	actual, err := json.Marshal(&r)
	assert.NoError(t, err)
	assert.JSONEq(t, expected, string(actual))
	assert.Equal(t, "/loadbalancer/sid/frontends", r.RequestURL())
}

func TestModifyLoadBalancerFrontendRequest(t *testing.T) {
	expected := `
	{
		"name": "example-frontend",
		"mode": "http",
		"port": 443,
		"default_backend": "example-backend"
	}`
	r := ModifyLoadBalancerFrontendRequest{
		ServiceUUID: "sid",
		Name:        "example",
		Frontend: ModifyLoadBalancerFrontend{
			Name:           "example-frontend",
			Mode:           upcloud.LoadBalancerModeHTTP,
			Port:           443,
			DefaultBackend: "example-backend"},
	}
	actual, err := json.Marshal(&r)
	assert.NoError(t, err)
	assert.JSONEq(t, expected, string(actual))
	assert.Equal(t, "/loadbalancer/sid/frontends/example", r.RequestURL())
}

func TestDeleteLoadBalancerFrontendRequest(t *testing.T) {
	r := DeleteLoadBalancerFrontendRequest{
		ServiceUUID: "sid",
		Name:        "example",
	}
	assert.Equal(t, "/loadbalancer/sid/frontends/example", r.RequestURL())
}
