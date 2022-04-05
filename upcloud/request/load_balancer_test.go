package request

import (
	"encoding/json"
	"testing"

	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetLoadBalancersRequest(t *testing.T) {
	assert.Equal(t, "/load-balancer", (&GetLoadBalancersRequest{}).RequestURL())
	assert.Equal(t, "/load-balancer?limit=50&offset=450", (&GetLoadBalancersRequest{Page: &Page{Number: 10, Size: 50}}).RequestURL())
}

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
	assert.Equal(t, "/load-balancer", r.RequestURL())
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
	assert.Exactly(t, "/load-balancer/lb/backends", r.RequestURL())
	assert.JSONEq(t, expectedJson, string(actualJson))
}

func TestGetLoadBalancerBackendsRequest(t *testing.T) {
	r := GetLoadBalancerBackendsRequest{
		ServiceUUID: "lb",
	}

	expectedJson := "{}"
	actualJson, err := json.Marshal(&r)

	require.NoError(t, err)
	assert.Exactly(t, "/load-balancer/lb/backends", r.RequestURL())
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
	assert.Exactly(t, "/load-balancer/lb/backends/be", r.RequestURL())
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
	assert.Exactly(t, "/load-balancer/lb/backends/be", r.RequestURL())
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
	assert.Exactly(t, "/load-balancer/lb/backends/be", r.RequestURL())
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
			Enabled:     false,
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
		"enabled": false,
		"type": "static",
		"ip": "10.0.0.1",
		"port": 80
	}`

	actualJson, err := json.Marshal(&r)

	require.NoError(t, err)
	assert.Exactly(t, "/load-balancer/lb/backends/be/members", r.RequestURL())
	assert.JSONEq(t, expectedJson, string(actualJson))

	r = CreateLoadBalancerBackendMemberRequest{
		ServiceUUID: "lb",
		BackendName: "be",
		Member: LoadBalancerBackendMember{
			Name:        "mem",
			Weight:      0,
			MaxSessions: 0,
			Enabled:     true,
			Type:        "static",
			IP:          "10.0.0.1",
			Port:        80,
		},
	}

	expectedJson = `
	{
		"name": "mem",
		"weight": 0,
		"max_sessions": 0,
		"enabled": true,
		"type": "static",
		"ip": "10.0.0.1",
		"port": 80
	}`

	actualJson, err = json.Marshal(&r)

	require.NoError(t, err)
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
	assert.Exactly(t, "/load-balancer/lb/backends/be/members", r.RequestURL())
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
	assert.Exactly(t, "/load-balancer/lb/backends/be/members/mem", r.RequestURL())
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
	assert.Exactly(t, "/load-balancer/lb/backends/be/members/mem", r.RequestURL())
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
	assert.Exactly(t, "/load-balancer/lb/backends/be/members/mem", r.RequestURL())
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
	assert.EqualValues(t, "/load-balancer/service-uuid/resolvers", r.RequestURL())
	assert.JSONEq(t, expectedJson, string(actualJson))
}

func TestGetLoadBalancerResolversRequest(t *testing.T) {
	r := GetLoadBalancerResolversRequest{
		ServiceUUID: "service-uuid",
	}

	expectedJson := "{}"
	actualJson, err := json.Marshal(&r)

	require.NoError(t, err)
	assert.EqualValues(t, "/load-balancer/service-uuid/resolvers", r.RequestURL())
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
	assert.EqualValues(t, "/load-balancer/service-uuid/resolvers/sesese", r.RequestURL())
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
	assert.EqualValues(t, "/load-balancer/service-uuid/resolvers/sesese", r.RequestURL())
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
	assert.EqualValues(t, "/load-balancer/service-uuid/resolvers/sesese", r.RequestURL())
	assert.JSONEq(t, expectedJson, string(actualJson))
}

func TestGetLoadBalancerPlansRequest(t *testing.T) {
	r := GetLoadBalancerPlansRequest{}
	assert.Equal(t, "/load-balancer/plans", r.RequestURL())
	r = GetLoadBalancerPlansRequest{Page: DefaultPage}
	assert.Equal(t, "/load-balancer/plans?limit=100&offset=0", r.RequestURL())
}

func TestGetLoadBalancerFrontendsRequest(t *testing.T) {
	r := GetLoadBalancerFrontendsRequest{"sid"}
	assert.Equal(t, "/load-balancer/sid/frontends", r.RequestURL())
}

func TestGetLoadBalancerFrontendRequest(t *testing.T) {
	r := GetLoadBalancerFrontendRequest{
		ServiceUUID: "sid",
		Name:        "be_name",
	}
	assert.Equal(t, "/load-balancer/sid/frontends/be_name", r.RequestURL())
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
			TLSConfigs: []LoadBalancerFrontendTLSConfig{{
				Name:                  "example-tls-config",
				CertificateBundleUUID: "0aded5c1-c7a3-498a-b9c8-a871611c47a2",
			}},
		},
	}
	actual, err := json.Marshal(&r)
	assert.NoError(t, err)
	assert.JSONEq(t, expected, string(actual))
	assert.Equal(t, "/load-balancer/sid/frontends", r.RequestURL())
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
	assert.Equal(t, "/load-balancer/sid/frontends/example", r.RequestURL())
}

func TestDeleteLoadBalancerFrontendRequest(t *testing.T) {
	r := DeleteLoadBalancerFrontendRequest{
		ServiceUUID: "sid",
		Name:        "example",
	}
	assert.Equal(t, "/load-balancer/sid/frontends/example", r.RequestURL())
}

func TestGetLoadBalancerFrontendRulesRequest(t *testing.T) {
	r := GetLoadBalancerFrontendRulesRequest{
		ServiceUUID:  "sid",
		FrontendName: "fename",
	}
	assert.Equal(t, "/load-balancer/sid/frontends/fename/rules", r.RequestURL())
}

func TestGetLoadBalancerFrontendRuleRequest(t *testing.T) {
	r := GetLoadBalancerFrontendRuleRequest{
		ServiceUUID:  "sid",
		FrontendName: "fename",
		Name:         "name",
	}
	assert.Equal(t, "/load-balancer/sid/frontends/fename/rules/name", r.RequestURL())
}

func TestCreateLoadBalancerFrontendRuleRequest(t *testing.T) {
	expected := `
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
	`
	r := CreateLoadBalancerFrontendRuleRequest{
		ServiceUUID:  "sid",
		FrontendName: "fename",
		Rule: LoadBalancerFrontendRule{
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
		},
	}
	actual, err := json.Marshal(&r)
	assert.NoError(t, err)
	assert.JSONEq(t, expected, string(actual))
	assert.Equal(t, "/load-balancer/sid/frontends/fename/rules", r.RequestURL())
}

func TestReplaceLoadBalancerFrontendRuleRequest(t *testing.T) {
	expected := `
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
	`
	r := ReplaceLoadBalancerFrontendRuleRequest{
		ServiceUUID:  "sid",
		FrontendName: "fename",
		Name:         "example-rule-1",
		Rule: LoadBalancerFrontendRule{
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
		},
	}
	actual, err := json.Marshal(&r)
	assert.NoError(t, err)
	assert.JSONEq(t, expected, string(actual))
	assert.Equal(t, "/load-balancer/sid/frontends/fename/rules/example-rule-1", r.RequestURL())
}

func TestModifyLoadBalancerFrontendRuleRequest(t *testing.T) {
	expected := `
	{
		"name": "example-rule-2",
		"priority": 100
	}
	`
	r := ModifyLoadBalancerFrontendRuleRequest{
		ServiceUUID:  "sid",
		FrontendName: "fename",
		Name:         "example-rule-1",
		Rule: ModifyLoadBalancerFrontendRule{
			Name:     "example-rule-2",
			Priority: 100,
		},
	}
	actual, err := json.Marshal(&r)
	assert.NoError(t, err)
	assert.JSONEq(t, expected, string(actual))
	assert.Equal(t, "/load-balancer/sid/frontends/fename/rules/example-rule-1", r.RequestURL())
}

func TestDeleteLoadBalancerFrontendRuleRequest(t *testing.T) {
	r := DeleteLoadBalancerFrontendRuleRequest{
		ServiceUUID:  "sid",
		FrontendName: "fename",
		Name:         "name",
	}
	assert.Equal(t, "/load-balancer/sid/frontends/fename/rules/name", r.RequestURL())
}

func TestGetLoadBalancerFrontendTLSConfigsRequest(t *testing.T) {
	r := GetLoadBalancerFrontendTLSConfigsRequest{
		ServiceUUID:  "sid",
		FrontendName: "fename",
	}
	assert.Equal(t, "/load-balancer/sid/frontends/fename/tls-configs", r.RequestURL())
}

func TestGetLoadBalancerFrontendTLSConfigRequest(t *testing.T) {
	r := GetLoadBalancerFrontendTLSConfigRequest{
		ServiceUUID:  "sid",
		FrontendName: "fename",
		Name:         "cfg",
	}
	assert.Equal(t, "/load-balancer/sid/frontends/fename/tls-configs/cfg", r.RequestURL())
}

func TestCreateLoadBalancerFrontendTLSConfigRequest(t *testing.T) {
	expected := `
	{
		"name": "example-tls-config",
		"certificate_bundle_uuid": "0aded5c1-c7a3-498a-b9c8-a871611c47a2"
	}
	`
	r := CreateLoadBalancerFrontendTLSConfigRequest{
		ServiceUUID:  "sid",
		FrontendName: "fename",
		Config: LoadBalancerFrontendTLSConfig{
			Name:                  "example-tls-config",
			CertificateBundleUUID: "0aded5c1-c7a3-498a-b9c8-a871611c47a2",
		},
	}
	actual, err := json.Marshal(&r)
	assert.NoError(t, err)
	assert.JSONEq(t, expected, string(actual))
	assert.Equal(t, "/load-balancer/sid/frontends/fename/tls-configs", r.RequestURL())
}

func TestModifyLoadBalancerFrontendTLSConfigRequest(t *testing.T) {
	expected := `
	{
		"name": "example-tls-config",
		"certificate_bundle_uuid": "0aded5c1-c7a3-498a-b9c8-a871611c47a2"
	}
	`
	r := ModifyLoadBalancerFrontendTLSConfigRequest{
		ServiceUUID:  "sid",
		FrontendName: "fename",
		Name:         "cfg",
		Config: LoadBalancerFrontendTLSConfig{
			Name:                  "example-tls-config",
			CertificateBundleUUID: "0aded5c1-c7a3-498a-b9c8-a871611c47a2",
		},
	}
	actual, err := json.Marshal(&r)
	assert.NoError(t, err)
	assert.JSONEq(t, expected, string(actual))
	assert.Equal(t, "/load-balancer/sid/frontends/fename/tls-configs/cfg", r.RequestURL())
}

func TestDeleteLoadBalancerFrontendTLSConfigRequest(t *testing.T) {
	r := GetLoadBalancerFrontendTLSConfigRequest{
		ServiceUUID:  "sid",
		FrontendName: "fename",
		Name:         "cfg",
	}
	assert.Equal(t, "/load-balancer/sid/frontends/fename/tls-configs/cfg", r.RequestURL())
}

func TestCreateLoadBalancerManualCertificateBundleRequest(t *testing.T) {
	expected := `
	{
		"name": "example-manual-certificate",
		"type": "manual",
		"certificate": "LS0LS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUNIVENDQWNPZ0F3SUJBZ0lVSWlNbzg1cGd0b25kUmVESU1McVR4YjhncHI0d0NnWUlLb1pJemowRUF3SXcKWkRFTE1Ba0dBMVVFQmhNQ1FWVXhFekFSQmdOVkJBZ01DbE52YldVdFUzUmhkR1V4SVRBZkJnTlZCQW9NR0VsdQpkR1Z5Ym1WMElGZHBaR2RwZEhNZ1VIUjVJRXgwWkRFZE1Cc0dBMVVFQXd3VVpHVjJMblZ3YkdJdWRYQmpiRzkxClpDNWpiMjB3SGhjTk1qRXhNREl5TVRJeE1ETTJXaGNOTXpFeE1ESXdNVEl4TURNMldqQmtNUXN3Q1FZRFZRUUcKRXdKQlZURVRNQkVHQTFVRUNBd0tVMjl0WlMxVGRHRjBaVEVoTUI4R0ExVUVDZ3dZU1c1MFpYSnVaWFFnVjJsawpaMmwwY3lCUWRIa2dUSFJrTVIwd0d3WURWUVFEREJSa1pYWXVkWEJzWWk1MWNHTnNiM1ZrTG1OdmJUQlpNQk1HCkJ5cUdTTTQ5QWdFR0NDcUdTTTQ5QXdFSEEwSUFCQmpVcUgyaVNuMFV2ZkU3UDdkT0QrSDBoKytxUWpnTG9OeWQKSTFwMmkvdlJPZmhMa0hCUjIxZ2JCSUdENjllYU1WWnZ4RWNxVWlKVWYwcmxLU3FzKzIralV6QlJNQjBHQTFVZApEZ1FXQkJTYTFaU3V1NkxJczMrc2lSSUJ5MHRXL3RnamZEQWZCZ05WSFNNRUdEQVdnQlNhMVpTdXU2TElzMytzCmlSSUJ5MHRXL3RnamZEQVBCZ05WSFJNQkFmOEVCVEFEQVFIL01Bb0dDQ3FHU000OUJBTUNBMGdBTUVVQ0lRQ3IKWXA5dHc2TmVXTHZGOGwrWm9rSE9QUzUzaGc2SDM0OHNMSjEvNit4YXN3SWdWVmN6WkFDc3JyUWt3TnVBZEVCeQo5TkxJR1VrWlhqeWgwdVFCS2x4Si9Wdz0KLS0tLS1FTkQgQ0VSVElGSUNBVEUtLS0tLQo=",
		"intermediates": "LS0tLS1CRdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUJ0VENDQVZ1Z0F4SUJBZ0lSQU5wSDZzV0ZtQzErWkdnUzFMWllVZGN3Q2dZSUtvWkl6ajBFQXdJd0pERU0KTUFvR0ExVUVDaE1EWkdWMk1SUXdFZ1lEVlFRREV3dGtaWFlnVW05dmRDQkRRVEFlRncweU1URXlNRGt4TXpVMwpNREZhRncwek1URXlNRGN4TXpVM01ERmFNQ3d4RERBS0JnTlZCQW9UQTJSbGRqRWNNQm9HQTFVRUF4TVRaR1YyCklFbHVkR1Z5YldWa2FXRjBaU0JEUVRCWk1CTUdCeXFHU000OUFnRUdDQ3FHU000OUF3RUhBMElBQkswbGMzNmcKN01TaDJTaXd3MUdDUjkvL3lSODR6S1VuNml6SmdCUkpFTlBxbmNXcjQzTi8rNktJR1EraERaazhRWHZ6RmExYQp2dFloc3JEVGtnRm9EV0tqWmpCa01BNEdBMVVkRHdFQi93UUVBd0lCQmpBU0JnTlZIUk1CQWY4RUNEQUdBUUgvCkFnRUFNQjBHQTFVZERnUVdCQlRWcG44d3hraHZhYVhvajF6c0Rkcmk4eGJuSnpBZkJnTlZIU01FR0RBV2dCU2oKckgwV0pubDdUSUJtc3NESGVveENFTVZyRmpBS0JnZ3Foa2pPUFFRREFnTklBREJGQWlBa3NhUXdPMkFESGhBLwppRVR1SVY1dTlNV3hFTU5BVGlVODFIZjc0cGVhWlFJaEFLMnJDRmhVVnQxbFlzR1o3dFdjWGFHVDhyU1k2cU1YClBmK3dnUXFnNXUyVAotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0tCg==",
		"private_key": "LS0tL1CRUdJTiBQUklWQVRFIEtFWS0tLS0tCk1JR0hBZ0VBTUJNR0J5cUdTTTQ5QWdFR0NDcUdTTTQ5QXdFSEJHMHdhd0lCQVFRZ3NQMzI2RlIxcmNwL0xybmcKNFBCT3BLRjIzSUNaM01GdGNrZFJuWkFESnRlaFJBTkNBQVFZMUtoOW9rcDlGTDN4T3orM1RnL2g5SWZ2cWtJNApDNkRjblNOYWRvdjcwVG40UzVCd1VkdFlHd1NCZyt2WG1qRldiOFJIS2xJaVZIOUs1U2txclB0dgotLS0tLUVORCBQUklWQVRFIEtFWS0tLS0tCg=="
	}
	`
	r := CreateLoadBalancerCertificateBundleRequest{
		Type:          upcloud.LoadBalancerCertificateBundleTypeManual,
		Name:          "example-manual-certificate",
		Certificate:   "LS0LS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUNIVENDQWNPZ0F3SUJBZ0lVSWlNbzg1cGd0b25kUmVESU1McVR4YjhncHI0d0NnWUlLb1pJemowRUF3SXcKWkRFTE1Ba0dBMVVFQmhNQ1FWVXhFekFSQmdOVkJBZ01DbE52YldVdFUzUmhkR1V4SVRBZkJnTlZCQW9NR0VsdQpkR1Z5Ym1WMElGZHBaR2RwZEhNZ1VIUjVJRXgwWkRFZE1Cc0dBMVVFQXd3VVpHVjJMblZ3YkdJdWRYQmpiRzkxClpDNWpiMjB3SGhjTk1qRXhNREl5TVRJeE1ETTJXaGNOTXpFeE1ESXdNVEl4TURNMldqQmtNUXN3Q1FZRFZRUUcKRXdKQlZURVRNQkVHQTFVRUNBd0tVMjl0WlMxVGRHRjBaVEVoTUI4R0ExVUVDZ3dZU1c1MFpYSnVaWFFnVjJsawpaMmwwY3lCUWRIa2dUSFJrTVIwd0d3WURWUVFEREJSa1pYWXVkWEJzWWk1MWNHTnNiM1ZrTG1OdmJUQlpNQk1HCkJ5cUdTTTQ5QWdFR0NDcUdTTTQ5QXdFSEEwSUFCQmpVcUgyaVNuMFV2ZkU3UDdkT0QrSDBoKytxUWpnTG9OeWQKSTFwMmkvdlJPZmhMa0hCUjIxZ2JCSUdENjllYU1WWnZ4RWNxVWlKVWYwcmxLU3FzKzIralV6QlJNQjBHQTFVZApEZ1FXQkJTYTFaU3V1NkxJczMrc2lSSUJ5MHRXL3RnamZEQWZCZ05WSFNNRUdEQVdnQlNhMVpTdXU2TElzMytzCmlSSUJ5MHRXL3RnamZEQVBCZ05WSFJNQkFmOEVCVEFEQVFIL01Bb0dDQ3FHU000OUJBTUNBMGdBTUVVQ0lRQ3IKWXA5dHc2TmVXTHZGOGwrWm9rSE9QUzUzaGc2SDM0OHNMSjEvNit4YXN3SWdWVmN6WkFDc3JyUWt3TnVBZEVCeQo5TkxJR1VrWlhqeWgwdVFCS2x4Si9Wdz0KLS0tLS1FTkQgQ0VSVElGSUNBVEUtLS0tLQo=",
		Intermediates: "LS0tLS1CRdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUJ0VENDQVZ1Z0F4SUJBZ0lSQU5wSDZzV0ZtQzErWkdnUzFMWllVZGN3Q2dZSUtvWkl6ajBFQXdJd0pERU0KTUFvR0ExVUVDaE1EWkdWMk1SUXdFZ1lEVlFRREV3dGtaWFlnVW05dmRDQkRRVEFlRncweU1URXlNRGt4TXpVMwpNREZhRncwek1URXlNRGN4TXpVM01ERmFNQ3d4RERBS0JnTlZCQW9UQTJSbGRqRWNNQm9HQTFVRUF4TVRaR1YyCklFbHVkR1Z5YldWa2FXRjBaU0JEUVRCWk1CTUdCeXFHU000OUFnRUdDQ3FHU000OUF3RUhBMElBQkswbGMzNmcKN01TaDJTaXd3MUdDUjkvL3lSODR6S1VuNml6SmdCUkpFTlBxbmNXcjQzTi8rNktJR1EraERaazhRWHZ6RmExYQp2dFloc3JEVGtnRm9EV0tqWmpCa01BNEdBMVVkRHdFQi93UUVBd0lCQmpBU0JnTlZIUk1CQWY4RUNEQUdBUUgvCkFnRUFNQjBHQTFVZERnUVdCQlRWcG44d3hraHZhYVhvajF6c0Rkcmk4eGJuSnpBZkJnTlZIU01FR0RBV2dCU2oKckgwV0pubDdUSUJtc3NESGVveENFTVZyRmpBS0JnZ3Foa2pPUFFRREFnTklBREJGQWlBa3NhUXdPMkFESGhBLwppRVR1SVY1dTlNV3hFTU5BVGlVODFIZjc0cGVhWlFJaEFLMnJDRmhVVnQxbFlzR1o3dFdjWGFHVDhyU1k2cU1YClBmK3dnUXFnNXUyVAotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0tCg==",
		PrivateKey:    "LS0tL1CRUdJTiBQUklWQVRFIEtFWS0tLS0tCk1JR0hBZ0VBTUJNR0J5cUdTTTQ5QWdFR0NDcUdTTTQ5QXdFSEJHMHdhd0lCQVFRZ3NQMzI2RlIxcmNwL0xybmcKNFBCT3BLRjIzSUNaM01GdGNrZFJuWkFESnRlaFJBTkNBQVFZMUtoOW9rcDlGTDN4T3orM1RnL2g5SWZ2cWtJNApDNkRjblNOYWRvdjcwVG40UzVCd1VkdFlHd1NCZyt2WG1qRldiOFJIS2xJaVZIOUs1U2txclB0dgotLS0tLUVORCBQUklWQVRFIEtFWS0tLS0tCg==",
	}
	actual, err := json.Marshal(&r)
	assert.NoError(t, err)
	assert.JSONEq(t, expected, string(actual))
	assert.Equal(t, "/load-balancer/certificate-bundles", r.RequestURL())
}

func TestCreateLoadBalancerDynamicCertificateBundleRequest(t *testing.T) {
	expected := `
	{
		"name": "example-dynamic-certificate",
		"type": "dynamic",
		"hostnames": [
			"example.com",
			"app.example.com"
		],
		"key_type": "rsa"
	}
	`
	r := CreateLoadBalancerCertificateBundleRequest{
		Type: upcloud.LoadBalancerCertificateBundleTypeDynamic,
		Name: "example-dynamic-certificate",
		Hostnames: []string{
			"example.com",
			"app.example.com",
		},
		KeyType: "rsa",
	}
	actual, err := json.Marshal(&r)
	assert.NoError(t, err)
	assert.JSONEq(t, expected, string(actual))
	assert.Equal(t, "/load-balancer/certificate-bundles", r.RequestURL())
}

func TestModifyLoadBalancerManualCertificateBundleRequest(t *testing.T) {
	expected := `
	{
		"name": "example-manual-certificate",
		"certificate": "LS0LS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUNIVENDQWNPZ0F3SUJBZ0lVSWlNbzg1cGd0b25kUmVESU1McVR4YjhncHI0d0NnWUlLb1pJemowRUF3SXcKWkRFTE1Ba0dBMVVFQmhNQ1FWVXhFekFSQmdOVkJBZ01DbE52YldVdFUzUmhkR1V4SVRBZkJnTlZCQW9NR0VsdQpkR1Z5Ym1WMElGZHBaR2RwZEhNZ1VIUjVJRXgwWkRFZE1Cc0dBMVVFQXd3VVpHVjJMblZ3YkdJdWRYQmpiRzkxClpDNWpiMjB3SGhjTk1qRXhNREl5TVRJeE1ETTJXaGNOTXpFeE1ESXdNVEl4TURNMldqQmtNUXN3Q1FZRFZRUUcKRXdKQlZURVRNQkVHQTFVRUNBd0tVMjl0WlMxVGRHRjBaVEVoTUI4R0ExVUVDZ3dZU1c1MFpYSnVaWFFnVjJsawpaMmwwY3lCUWRIa2dUSFJrTVIwd0d3WURWUVFEREJSa1pYWXVkWEJzWWk1MWNHTnNiM1ZrTG1OdmJUQlpNQk1HCkJ5cUdTTTQ5QWdFR0NDcUdTTTQ5QXdFSEEwSUFCQmpVcUgyaVNuMFV2ZkU3UDdkT0QrSDBoKytxUWpnTG9OeWQKSTFwMmkvdlJPZmhMa0hCUjIxZ2JCSUdENjllYU1WWnZ4RWNxVWlKVWYwcmxLU3FzKzIralV6QlJNQjBHQTFVZApEZ1FXQkJTYTFaU3V1NkxJczMrc2lSSUJ5MHRXL3RnamZEQWZCZ05WSFNNRUdEQVdnQlNhMVpTdXU2TElzMytzCmlSSUJ5MHRXL3RnamZEQVBCZ05WSFJNQkFmOEVCVEFEQVFIL01Bb0dDQ3FHU000OUJBTUNBMGdBTUVVQ0lRQ3IKWXA5dHc2TmVXTHZGOGwrWm9rSE9QUzUzaGc2SDM0OHNMSjEvNit4YXN3SWdWVmN6WkFDc3JyUWt3TnVBZEVCeQo5TkxJR1VrWlhqeWgwdVFCS2x4Si9Wdz0KLS0tLS1FTkQgQ0VSVElGSUNBVEUtLS0tLQo=",
		"intermediates": "LS0tLS1CRdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUJ0VENDQVZ1Z0F4SUJBZ0lSQU5wSDZzV0ZtQzErWkdnUzFMWllVZGN3Q2dZSUtvWkl6ajBFQXdJd0pERU0KTUFvR0ExVUVDaE1EWkdWMk1SUXdFZ1lEVlFRREV3dGtaWFlnVW05dmRDQkRRVEFlRncweU1URXlNRGt4TXpVMwpNREZhRncwek1URXlNRGN4TXpVM01ERmFNQ3d4RERBS0JnTlZCQW9UQTJSbGRqRWNNQm9HQTFVRUF4TVRaR1YyCklFbHVkR1Z5YldWa2FXRjBaU0JEUVRCWk1CTUdCeXFHU000OUFnRUdDQ3FHU000OUF3RUhBMElBQkswbGMzNmcKN01TaDJTaXd3MUdDUjkvL3lSODR6S1VuNml6SmdCUkpFTlBxbmNXcjQzTi8rNktJR1EraERaazhRWHZ6RmExYQp2dFloc3JEVGtnRm9EV0tqWmpCa01BNEdBMVVkRHdFQi93UUVBd0lCQmpBU0JnTlZIUk1CQWY4RUNEQUdBUUgvCkFnRUFNQjBHQTFVZERnUVdCQlRWcG44d3hraHZhYVhvajF6c0Rkcmk4eGJuSnpBZkJnTlZIU01FR0RBV2dCU2oKckgwV0pubDdUSUJtc3NESGVveENFTVZyRmpBS0JnZ3Foa2pPUFFRREFnTklBREJGQWlBa3NhUXdPMkFESGhBLwppRVR1SVY1dTlNV3hFTU5BVGlVODFIZjc0cGVhWlFJaEFLMnJDRmhVVnQxbFlzR1o3dFdjWGFHVDhyU1k2cU1YClBmK3dnUXFnNXUyVAotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0tCg==",
		"private_key": "LS0tL1CRUdJTiBQUklWQVRFIEtFWS0tLS0tCk1JR0hBZ0VBTUJNR0J5cUdTTTQ5QWdFR0NDcUdTTTQ5QXdFSEJHMHdhd0lCQVFRZ3NQMzI2RlIxcmNwL0xybmcKNFBCT3BLRjIzSUNaM01GdGNrZFJuWkFESnRlaFJBTkNBQVFZMUtoOW9rcDlGTDN4T3orM1RnL2g5SWZ2cWtJNApDNkRjblNOYWRvdjcwVG40UzVCd1VkdFlHd1NCZyt2WG1qRldiOFJIS2xJaVZIOUs1U2txclB0dgotLS0tLUVORCBQUklWQVRFIEtFWS0tLS0tCg=="
	}
	`
	r := ModifyLoadBalancerCertificateBundleRequest{
		UUID:          "id",
		Name:          "example-manual-certificate",
		Certificate:   "LS0LS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUNIVENDQWNPZ0F3SUJBZ0lVSWlNbzg1cGd0b25kUmVESU1McVR4YjhncHI0d0NnWUlLb1pJemowRUF3SXcKWkRFTE1Ba0dBMVVFQmhNQ1FWVXhFekFSQmdOVkJBZ01DbE52YldVdFUzUmhkR1V4SVRBZkJnTlZCQW9NR0VsdQpkR1Z5Ym1WMElGZHBaR2RwZEhNZ1VIUjVJRXgwWkRFZE1Cc0dBMVVFQXd3VVpHVjJMblZ3YkdJdWRYQmpiRzkxClpDNWpiMjB3SGhjTk1qRXhNREl5TVRJeE1ETTJXaGNOTXpFeE1ESXdNVEl4TURNMldqQmtNUXN3Q1FZRFZRUUcKRXdKQlZURVRNQkVHQTFVRUNBd0tVMjl0WlMxVGRHRjBaVEVoTUI4R0ExVUVDZ3dZU1c1MFpYSnVaWFFnVjJsawpaMmwwY3lCUWRIa2dUSFJrTVIwd0d3WURWUVFEREJSa1pYWXVkWEJzWWk1MWNHTnNiM1ZrTG1OdmJUQlpNQk1HCkJ5cUdTTTQ5QWdFR0NDcUdTTTQ5QXdFSEEwSUFCQmpVcUgyaVNuMFV2ZkU3UDdkT0QrSDBoKytxUWpnTG9OeWQKSTFwMmkvdlJPZmhMa0hCUjIxZ2JCSUdENjllYU1WWnZ4RWNxVWlKVWYwcmxLU3FzKzIralV6QlJNQjBHQTFVZApEZ1FXQkJTYTFaU3V1NkxJczMrc2lSSUJ5MHRXL3RnamZEQWZCZ05WSFNNRUdEQVdnQlNhMVpTdXU2TElzMytzCmlSSUJ5MHRXL3RnamZEQVBCZ05WSFJNQkFmOEVCVEFEQVFIL01Bb0dDQ3FHU000OUJBTUNBMGdBTUVVQ0lRQ3IKWXA5dHc2TmVXTHZGOGwrWm9rSE9QUzUzaGc2SDM0OHNMSjEvNit4YXN3SWdWVmN6WkFDc3JyUWt3TnVBZEVCeQo5TkxJR1VrWlhqeWgwdVFCS2x4Si9Wdz0KLS0tLS1FTkQgQ0VSVElGSUNBVEUtLS0tLQo=",
		Intermediates: "LS0tLS1CRdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUJ0VENDQVZ1Z0F4SUJBZ0lSQU5wSDZzV0ZtQzErWkdnUzFMWllVZGN3Q2dZSUtvWkl6ajBFQXdJd0pERU0KTUFvR0ExVUVDaE1EWkdWMk1SUXdFZ1lEVlFRREV3dGtaWFlnVW05dmRDQkRRVEFlRncweU1URXlNRGt4TXpVMwpNREZhRncwek1URXlNRGN4TXpVM01ERmFNQ3d4RERBS0JnTlZCQW9UQTJSbGRqRWNNQm9HQTFVRUF4TVRaR1YyCklFbHVkR1Z5YldWa2FXRjBaU0JEUVRCWk1CTUdCeXFHU000OUFnRUdDQ3FHU000OUF3RUhBMElBQkswbGMzNmcKN01TaDJTaXd3MUdDUjkvL3lSODR6S1VuNml6SmdCUkpFTlBxbmNXcjQzTi8rNktJR1EraERaazhRWHZ6RmExYQp2dFloc3JEVGtnRm9EV0tqWmpCa01BNEdBMVVkRHdFQi93UUVBd0lCQmpBU0JnTlZIUk1CQWY4RUNEQUdBUUgvCkFnRUFNQjBHQTFVZERnUVdCQlRWcG44d3hraHZhYVhvajF6c0Rkcmk4eGJuSnpBZkJnTlZIU01FR0RBV2dCU2oKckgwV0pubDdUSUJtc3NESGVveENFTVZyRmpBS0JnZ3Foa2pPUFFRREFnTklBREJGQWlBa3NhUXdPMkFESGhBLwppRVR1SVY1dTlNV3hFTU5BVGlVODFIZjc0cGVhWlFJaEFLMnJDRmhVVnQxbFlzR1o3dFdjWGFHVDhyU1k2cU1YClBmK3dnUXFnNXUyVAotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0tCg==",
		PrivateKey:    "LS0tL1CRUdJTiBQUklWQVRFIEtFWS0tLS0tCk1JR0hBZ0VBTUJNR0J5cUdTTTQ5QWdFR0NDcUdTTTQ5QXdFSEJHMHdhd0lCQVFRZ3NQMzI2RlIxcmNwL0xybmcKNFBCT3BLRjIzSUNaM01GdGNrZFJuWkFESnRlaFJBTkNBQVFZMUtoOW9rcDlGTDN4T3orM1RnL2g5SWZ2cWtJNApDNkRjblNOYWRvdjcwVG40UzVCd1VkdFlHd1NCZyt2WG1qRldiOFJIS2xJaVZIOUs1U2txclB0dgotLS0tLUVORCBQUklWQVRFIEtFWS0tLS0tCg==",
	}
	actual, err := json.Marshal(&r)
	assert.NoError(t, err)
	assert.JSONEq(t, expected, string(actual))
	assert.Equal(t, "/load-balancer/certificate-bundles/id", r.RequestURL())
}

func TestModifyLoadBalancerDynamicCertificateBundleRequest(t *testing.T) {
	expected := `
	{
		"name": "example-dynamic-certificate",
		"hostnames": [
			"example.com",
			"app.example.com"
		]
	}
	`
	r := ModifyLoadBalancerCertificateBundleRequest{
		UUID: "id",
		Name: "example-dynamic-certificate",
		Hostnames: []string{
			"example.com",
			"app.example.com",
		},
	}
	actual, err := json.Marshal(&r)
	assert.NoError(t, err)
	assert.JSONEq(t, expected, string(actual))
	assert.Equal(t, "/load-balancer/certificate-bundles/id", r.RequestURL())
}

func TestGetLoadBalancerCertificateBundlesRequest(t *testing.T) {
	r := GetLoadBalancerCertificateBundlesRequest{}
	assert.Equal(t, "/load-balancer/certificate-bundles", r.RequestURL())
}

func TestGetLoadBalancerCertificateBundleRequest(t *testing.T) {
	r := GetLoadBalancerCertificateBundleRequest{UUID: "id"}
	assert.Equal(t, "/load-balancer/certificate-bundles/id", r.RequestURL())
}

func TestDeleteLoadBalancerCertificateBundleRequest(t *testing.T) {
	r := GetLoadBalancerCertificateBundleRequest{UUID: "id"}
	assert.Equal(t, "/load-balancer/certificate-bundles/id", r.RequestURL())
}
