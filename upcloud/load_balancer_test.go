package upcloud

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMarshalLoadBalancer(t *testing.T) {
	t.Parallel()
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
			"labels": [
				{
					"key": "managedby",
					"value": "upcloud-go-sdk-unit-test"
				}
			],
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
					],
					"tls_configs": [
						{
							"certificate_bundle_uuid": "0aded5c1-c7a3-498a-b9c8-a871611c47a3",
							"created_at": "2023-02-11T17:33:08.490581Z",
							"name": "example-tls-config",
							"updated_at": "2023-02-11T17:33:08.490581Z"
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
			],
			"nodes": [
            	{
					"operational_state": "running",
					"networks": [
						{
							"ip_addresses": [
								{
									"address": "100.127.5.146",
									"listen": true
								},
								{
									"address": "100.127.5.140",
									"listen": false
								}
							],
							"name": "example-network-1",
							"type": "public"
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
		Labels: []Label{
			{
				Key:   "managedby",
				Value: "upcloud-go-sdk-unit-test",
			},
		},
		Frontends: []LoadBalancerFrontend{
			{
				Name:           "example-frontend",
				Mode:           LoadBalancerModeHTTP,
				Port:           443,
				DefaultBackend: "example-backend-1",
				Rules: []LoadBalancerFrontendRule{
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
				TLSConfigs: []LoadBalancerFrontendTLSConfig{{
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
				Members: []LoadBalancerBackendMember{
					{
						Name:        "example-member-1",
						IP:          "172.16.1.4",
						Port:        8000,
						Weight:      100,
						MaxSessions: 1000,
						Type:        LoadBalancerBackendMemberTypeStatic,
						Enabled:     true,
						CreatedAt:   timeParse("2021-12-07T13:58:30.817272Z"),
						UpdatedAt:   timeParse("2022-02-11T17:33:08.490581Z"),
					},
				},
				TLSConfigs: []LoadBalancerBackendTLSConfig{{
					Name:                  "example-tls-config",
					CertificateBundleUUID: "0aded5c1-c7a3-498a-b9c8-a871611c47a3",
					CreatedAt:             timeParse("2023-02-11T17:33:08.490581Z"),
					UpdatedAt:             timeParse("2023-02-11T17:33:08.490581Z"),
				}},
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
		Nodes: []LoadBalancerNode{{
			OperationalState: LoadBalancerNodeOperationalStateRunning,
			Networks: []LoadBalancerNodeNetwork{{
				Name: "example-network-1",
				Type: LoadBalancerNetworkTypePublic,
				IPAddresses: []LoadBalancerIPAddress{
					{
						Address: "100.127.5.146",
						Listen:  true,
					},
					{
						Address: "100.127.5.140",
						Listen:  false,
					},
				},
			}},
		}},
	}}
	actual, err := json.Marshal(lbs)
	assert.NoError(t, err)
	assert.JSONEq(t, lbString, string(actual))
	l := []LoadBalancer{}
	err = json.Unmarshal([]byte(lbString), &l)
	assert.NoError(t, err)
	assert.Equal(t, lbs, l)
}

func TestLoadBalancerPlan(t *testing.T) {
	t.Parallel()
	testJSON(t,
		&LoadBalancerPlan{},
		&LoadBalancerPlan{
			Name:                 "development",
			PerServerMaxSessions: 10000,
			ServerNumber:         1,
		},
		`
		{
			"name": "development",
			"per_server_max_sessions": 10000,
			"server_number": 1
		}
		`,
	)
}

func TestLoadBalancerFrontend(t *testing.T) {
	t.Parallel()
	testJSON(t,
		&LoadBalancerFrontend{},
		&LoadBalancerFrontend{
			Name:           "example-frontend",
			Mode:           LoadBalancerModeTCP,
			Port:           443,
			DefaultBackend: "example-backend",
			Properties: &LoadBalancerFrontendProperties{
				TimeoutClient: 10,
			},
			CreatedAt: timeParse("2021-12-07T13:58:30.817272Z"),
			UpdatedAt: timeParse("2022-02-11T17:33:08.490581Z"),
		},
		`
		{
			"name": "example-frontend",
			"mode": "tcp",
			"port": 443,
			"default_backend": "example-backend",
			"properties": {
				"timeout_client": 10,
				"inbound_proxy_protocol": false
			},
			"created_at": "2021-12-07T13:58:30.817272Z",
			"updated_at": "2022-02-11T17:33:08.490581Z"
		}
		`,
	)
	testJSON(t,
		&LoadBalancerFrontend{},
		&LoadBalancerFrontend{
			Networks: []LoadBalancerFrontendNetwork{
				{
					Name: "PublicNet",
				},
				{
					Name: "PrivateNet",
				},
			},
			CreatedAt: timeParse("2021-12-07T13:58:30.817272Z"),
			UpdatedAt: timeParse("2022-02-11T17:33:08.490581Z"),
		},
		`
		{
			"networks": [
				{
					"name": "PublicNet"
				},
				{
					"name": "PrivateNet"
				}
			],
			"created_at": "2021-12-07T13:58:30.817272Z",
			"updated_at": "2022-02-11T17:33:08.490581Z"
		}
		`,
	)
}

func TestLoadBalancerFrontendProperties(t *testing.T) {
	t.Parallel()
	testJSON(t,
		&LoadBalancerFrontendProperties{},
		&LoadBalancerFrontendProperties{
			TimeoutClient: 10,
		},
		`
		{
			"timeout_client": 10,
			"inbound_proxy_protocol": false
		}
		`,
	)
	testJSON(t,
		&LoadBalancerFrontendProperties{},
		&LoadBalancerFrontendProperties{
			TimeoutClient:        10,
			InboundProxyProtocol: true,
		},
		`
		{
			"timeout_client": 10,
			"inbound_proxy_protocol": true
		}
		`,
	)
	testJSON(t,
		&LoadBalancerFrontendProperties{},
		&LoadBalancerFrontendProperties{
			TimeoutClient:        10,
			InboundProxyProtocol: false,
		},
		`
		{
			"timeout_client": 10,
			"inbound_proxy_protocol": false
		}
		`,
	)
	testJSON(t,
		&LoadBalancerFrontendProperties{},
		&LoadBalancerFrontendProperties{
			TimeoutClient:        10,
			InboundProxyProtocol: false,
			HTTP2Enabled:         BoolPtr(false),
		},
		`
		{
			"timeout_client": 10,
			"inbound_proxy_protocol": false,
			"http2_enabled": false
		}
		`,
	)
	testJSON(t,
		&LoadBalancerFrontendProperties{},
		&LoadBalancerFrontendProperties{
			TimeoutClient:        10,
			InboundProxyProtocol: false,
			HTTP2Enabled:         BoolPtr(true),
		},
		`
		{
			"timeout_client": 10,
			"inbound_proxy_protocol": false,
			"http2_enabled": true
		}
		`,
	)
}

func TestLoadBalancerRule(t *testing.T) {
	t.Parallel()
	testJSON(t,
		&LoadBalancerFrontendRule{},
		&LoadBalancerFrontendRule{
			Name:     "example-rule-1",
			Priority: 100,
			Actions: []LoadBalancerAction{
				{
					Type:      LoadBalancerActionTypeTCPReject,
					TCPReject: &LoadBalancerActionTCPReject{},
				},
				{
					Type: LoadBalancerActionTypeHTTPRedirect,
					HTTPRedirect: &LoadBalancerActionHTTPRedirect{
						Location: "/new",
					},
				},
				{
					Type: LoadBalancerActionTypeHTTPReturn,
					HTTPReturn: &LoadBalancerActionHTTPReturn{
						Status:      200,
						ContentType: "text/html",
						Payload:     "PGgxPmFwcGxlYmVlPC9oMT4K",
					},
				},
				{
					Type: LoadBalancerActionTypeUseBackend,
					UseBackend: &LoadBalancerActionUseBackend{
						Backend: "example-backend",
					},
				},
				{
					Type:                LoadBalancerActionTypeSetForwardedHeaders,
					SetForwardedHeaders: &LoadBalancerActionSetForwardedHeaders{},
				},
			},
			CreatedAt: timeParse("2021-12-07T13:58:30.817272Z"),
			UpdatedAt: timeParse("2022-02-11T17:33:08.490581Z"),
		},
		`
		{
			"name": "example-rule-1",
			"priority": 100,
			"created_at": "2021-12-07T13:58:30.817272Z",
			"updated_at": "2022-02-11T17:33:08.490581Z",
			"actions": [
				{
					"type": "tcp_reject",
					"action_tcp_reject": {}
				},
				{
					"type": "http_redirect",
					"action_http_redirect": {
						"location": "/new"
					}
				},
				{
					"type": "http_return",
					"action_http_return": {
							"status": 200,
							"content_type": "text/html",
							"payload": "PGgxPmFwcGxlYmVlPC9oMT4K"
					}
				},
				{
					"type": "use_backend",
					"action_use_backend": {
							"backend": "example-backend"
					}
				},
				{
					"type": "set_forwarded_headers",
					"action_set_forwarded_headers": {}
				}
			]
		}
		`,
	)
}

func TestLoadBalancerFrontendTLSConfig(t *testing.T) {
	t.Parallel()
	testJSON(t,
		&LoadBalancerFrontendTLSConfig{},
		&LoadBalancerFrontendTLSConfig{
			Name:                  "example-tls-config",
			CertificateBundleUUID: "0aded5c1-c7a3-498a-b9c8-a871611c47a2",
			CreatedAt:             timeParse("2022-02-11T17:33:08.490581Z"),
			UpdatedAt:             timeParse("2022-02-11T17:33:08.490581Z"),
		},
		`
		{
			"certificate_bundle_uuid": "0aded5c1-c7a3-498a-b9c8-a871611c47a2",
			"name": "example-tls-config",
			"created_at": "2022-02-11T17:33:08.490581Z",
			"updated_at": "2022-02-11T17:33:08.490581Z"
		}
		`,
	)
}

func TestLoadBalancerBackend(t *testing.T) {
	t.Parallel()
	testJSON(t,
		&LoadBalancerBackend{},
		&LoadBalancerBackend{
			Name:    "example-backend-2",
			Members: []LoadBalancerBackendMember{},
			Properties: &LoadBalancerBackendProperties{
				TimeoutServer: 30,
				TimeoutTunnel: 3600,
			},
			CreatedAt: timeParse("2022-02-11T17:33:08.490581Z"),
			UpdatedAt: timeParse("2022-02-11T17:33:08.490581Z"),
		},
		`
		{
			"name": "example-backend-2",
			"properties": {
				"timeout_server": 30,
				"timeout_tunnel": 3600
			},
			"created_at": "2022-02-11T17:33:08.490581Z",
			"updated_at": "2022-02-11T17:33:08.490581Z",
			"members": []
		}
		`,
	)
}

func TestLoadBalancerBackendTLSConfig(t *testing.T) {
	t.Parallel()
	testJSON(t,
		&LoadBalancerBackendTLSConfig{},
		&LoadBalancerBackendTLSConfig{
			Name:                  "example-tls-config",
			CertificateBundleUUID: "0aded5c1-c7a3-498a-b9c8-a871611c47a3",
			CreatedAt:             timeParse("2023-02-11T17:33:08.490581Z"),
			UpdatedAt:             timeParse("2023-02-11T17:33:08.490581Z"),
		},
		`
		{
			"certificate_bundle_uuid": "0aded5c1-c7a3-498a-b9c8-a871611c47a3",
			"name": "example-tls-config",
			"created_at": "2023-02-11T17:33:08.490581Z",
			"updated_at": "2023-02-11T17:33:08.490581Z"
		}
		`,
	)
}

func TestLoadBalancerBackendProperties(t *testing.T) {
	t.Parallel()
	testJSON(t,
		&LoadBalancerBackendProperties{},
		&LoadBalancerBackendProperties{
			TimeoutServer:             30,
			TimeoutTunnel:             3600,
			HealthCheckTLSVerify:      BoolPtr(false),
			HealthCheckType:           LoadBalancerHealthCheckTypeHTTP,
			HealthCheckInterval:       20,
			HealthCheckFall:           3,
			HealthCheckRise:           3,
			HealthCheckURL:            "/health",
			HealthCheckExpectedStatus: 200,
			StickySessionCookieName:   "SERVERID",
			OutboundProxyProtocol:     LoadBalancerProxyProtocolVersion1,
		},
		`
		{
			"timeout_server": 30,
			"timeout_tunnel": 3600,
			"health_check_type": "http",
			"health_check_tls_verify": false,
			"health_check_interval": 20,
			"health_check_fall": 3,
			"health_check_rise": 3,
			"health_check_url": "/health",
			"health_check_expected_status": 200,
			"sticky_session_cookie_name": "SERVERID",
			"outbound_proxy_protocol": "v1"
		}
		`,
	)
	testJSON(t,
		&LoadBalancerBackendProperties{},
		&LoadBalancerBackendProperties{
			TLSVerify:      BoolPtr(true),
			TLSEnabled:     BoolPtr(true),
			TLSUseSystemCA: BoolPtr(true),
			HTTP2Enabled:   BoolPtr(true),
		},
		`{
			"tls_verify": true,
			"tls_enabled": true,
			"tls_use_system_ca": true,
			"http2_enabled": true
		}
		`,
	)
}

func TestLoadBalancerBackendMember(t *testing.T) {
	t.Parallel()
	members := []LoadBalancerBackendMember{
		{
			Name:        "member.example.com",
			Port:        8000,
			Weight:      100,
			MaxSessions: 1000,
			Type:        LoadBalancerBackendMemberTypeDynamic,
			Enabled:     true,
			CreatedAt:   timeParse("2022-02-11T16:39:55.321306Z"),
			UpdatedAt:   timeParse("2022-02-11T17:33:08.490581Z"),
		},
		{
			Name:        "example-member-1",
			Port:        8000,
			Weight:      100,
			MaxSessions: 1000,
			IP:          "172.16.1.4",
			Type:        LoadBalancerBackendMemberTypeStatic,
			Enabled:     true,
			CreatedAt:   timeParse("2021-12-07T13:58:30.817272Z"),
			UpdatedAt:   timeParse("2022-02-11T17:33:08.490581Z"),
		},
	}
	membersEmpty := make([]LoadBalancerBackendMember, 0)
	testJSON(t,
		&membersEmpty,
		&members,
		`
		[
			{
				"created_at": "2022-02-11T16:39:55.321306Z",
				"enabled": true,
				"ip": "",
				"max_sessions": 1000,
				"name": "member.example.com",
				"port": 8000,
				"type": "dynamic",
				"updated_at": "2022-02-11T17:33:08.490581Z",
				"weight": 100
			},
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
		`,
	)
}

func TestLoadBalancerResolver(t *testing.T) {
	t.Parallel()
	testJSON(t,
		&LoadBalancerResolver{},
		&LoadBalancerResolver{
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
		`
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
		`,
	)
}

func TestLoadBalancerMatcherStringWithArgument(t *testing.T) {
	t.Parallel()
	tv := true
	testJSON(t,
		&LoadBalancerMatcher{},
		&LoadBalancerMatcher{
			Type: LoadBalancerMatcherTypeURLParam,
			URLParam: &LoadBalancerMatcherStringWithArgument{
				Method:     LoadBalancerStringMatcherMethodExact,
				Name:       "status",
				Value:      "active",
				IgnoreCase: &tv,
			},
		},
		`
		{
			"type": "url_param",
			"match_url_param": {
				"method": "exact",
				"name": "status",
				"value": "active",
				"ignore_case": true
			}
		}
		`,
	)
}

func TestLoadBalancerMatcherHost(t *testing.T) {
	t.Parallel()
	testJSON(t,
		&LoadBalancerMatcher{},
		&LoadBalancerMatcher{
			Type: LoadBalancerMatcherTypeHost,
			Host: &LoadBalancerMatcherHost{
				Value: "example.com",
			},
		},
		`
		{
			"type": "host",
			"match_host": {
				"value": "example.com"
			}
		}
		`,
	)
}

func TestLoadBalancerMatcherNumMembersUp(t *testing.T) {
	t.Parallel()
	testJSON(t,
		&LoadBalancerMatcher{},
		&LoadBalancerMatcher{
			Type: LoadBalancerMatcherTypeNumMembersUp,
			NumMembersUp: &LoadBalancerMatcherNumMembersUp{
				Method:  LoadBalancerIntegerMatcherMethodLess,
				Value:   1,
				Backend: "example-fallback-backend",
			},
		},
		`
		{
			"type": "num_members_up",
			"match_num_members_up": {
				"method": "less",
				"value": 1,
				"backend": "example-fallback-backend"
			}
		}
		`,
	)
}

func TestLoadBalancerMatcherHTTPMethod(t *testing.T) {
	t.Parallel()
	testJSON(t,
		&LoadBalancerMatcher{},
		&LoadBalancerMatcher{
			Type: LoadBalancerMatcherTypeHTTPMethod,
			HTTPMethod: &LoadBalancerMatcherHTTPMethod{
				Value: LoadBalancerHTTPMatcherMethodPatch,
			},
		},
		`
		{
			"type": "http_method",
			"match_http_method": {
				"value": "PATCH"
			}
		}
		`,
	)
}

func TestLoadBalancerMatcherInteger(t *testing.T) {
	t.Parallel()
	testJSON(t,
		&LoadBalancerMatcher{},
		&LoadBalancerMatcher{
			Type: LoadBalancerMatcherTypeSrcPort,
			SrcPort: &LoadBalancerMatcherInteger{
				Method: LoadBalancerIntegerMatcherMethodEqual,
				Value:  8000,
			},
		},
		`
		{
			"type": "src_port",
			"match_src_port": {
				"method": "equal",
				"value": 8000
			}
		}
		`,
	)
}

func TestLoadBalancerMatcherIntegerRange(t *testing.T) {
	t.Parallel()
	testJSON(t,
		&LoadBalancerMatcher{},
		&LoadBalancerMatcher{
			Type: LoadBalancerMatcherTypeSrcPort,
			SrcPort: &LoadBalancerMatcherInteger{
				Method:     LoadBalancerIntegerMatcherMethodRange,
				RangeStart: 8000,
				RangeEnd:   9000,
			},
		},
		`
		{
			"type": "src_port",
			"match_src_port": {
				"method": "range",
				"range_start": 8000,
				"range_end": 9000
			}
		}
		`,
	)
}

func TestLoadBalancerMatcherString(t *testing.T) {
	t.Parallel()
	tv := true
	testJSON(t,
		&LoadBalancerMatcher{},
		&LoadBalancerMatcher{
			Type: LoadBalancerMatcherTypePath,
			Path: &LoadBalancerMatcherString{
				Method:     LoadBalancerStringMatcherMethodStarts,
				Value:      "/application",
				IgnoreCase: &tv,
			},
		},
		`
		{
			"type": "path",
			"match_path": {
				"method": "starts",
				"value": "/application",
				"ignore_case": true
			}
		}
		`,
	)
}

func TestLoadBalancerMatcherSourceIP(t *testing.T) {
	t.Parallel()
	testJSON(t,
		&LoadBalancerMatcher{},
		&LoadBalancerMatcher{
			Type: LoadBalancerMatcherTypeSrcIP,
			SrcIP: &LoadBalancerMatcherSourceIP{
				Value: "213.3.44.11",
			},
		},
		`
		{
			"type": "src_ip",
			"match_src_ip": {
				"value": "213.3.44.11"
			}
		}
		`,
	)
}

func TestLoadBalancerMatcherSourceIP_Inverse(t *testing.T) {
	t.Parallel()
	testJSON(t,
		&LoadBalancerMatcher{},
		&LoadBalancerMatcher{
			Type:    LoadBalancerMatcherTypeSrcIP,
			Inverse: BoolPtr(true),
			SrcIP: &LoadBalancerMatcherSourceIP{
				Value: "213.3.44.11",
			},
		},
		`
		{
			"type": "src_ip",
			"inverse": true,
			"match_src_ip": {
				"value": "213.3.44.11"
			}
		}
		`,
	)
}

func TestLoadBalancerActionUseBackend(t *testing.T) {
	t.Parallel()
	testJSON(t,
		&LoadBalancerAction{},
		&LoadBalancerAction{
			Type: LoadBalancerActionTypeUseBackend,
			UseBackend: &LoadBalancerActionUseBackend{
				Backend: "example-backend",
			},
		},
		`
		{
			"type": "use_backend",
			"action_use_backend": {
				"backend": "example-backend"
			}
		}
		`,
	)
}

func TestLoadBalancerActionTCPReject(t *testing.T) {
	t.Parallel()
	testJSON(t,
		&LoadBalancerAction{},
		&LoadBalancerAction{
			Type:      LoadBalancerActionTypeTCPReject,
			TCPReject: &LoadBalancerActionTCPReject{},
		},
		`
		{
			"type": "tcp_reject",
			"action_tcp_reject": {}
		}
		`,
	)
}

func TestLoadBalancerActionHTTPReturn(t *testing.T) {
	t.Parallel()
	testJSON(t,
		&LoadBalancerAction{},
		&LoadBalancerAction{
			Type: LoadBalancerActionTypeHTTPReturn,
			HTTPReturn: &LoadBalancerActionHTTPReturn{
				Status:      200,
				ContentType: "text/html",
				Payload:     "PGgxPmFwcGxlYmVlPC9oMT4K",
			},
		},
		`
		{
			"type": "http_return",
			"action_http_return": {
				"status": 200,
				"content_type": "text/html",
				"payload": "PGgxPmFwcGxlYmVlPC9oMT4K"
			}
		}
		`,
	)
}

func TestLoadBalancerActionHTTPRedirect(t *testing.T) {
	t.Parallel()
	testJSON(t,
		&LoadBalancerAction{},
		&LoadBalancerAction{
			Type: LoadBalancerActionTypeHTTPRedirect,
			HTTPRedirect: &LoadBalancerActionHTTPRedirect{
				Location: "https://internal.example.com",
			},
		},
		`
		{
			"type": "http_redirect",
			"action_http_redirect": {
				"location": "https://internal.example.com"
			}
		}
		`,
	)
}

func TestLoadBalancerCertificateBundle(t *testing.T) {
	t.Parallel()
	testJSON(t,
		&LoadBalancerCertificateBundle{},
		&LoadBalancerCertificateBundle{
			UUID:             "bf571589-7378-41f8-879e-5505613b070d",
			Type:             LoadBalancerCertificateBundleTypeManual,
			OperationalState: LoadBalancerCertificateBundleOperationalStateIdle,
			Certificate:      "LS0LS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUNIVENDQWNPZ0F3SUJBZ0lVSWlNbzg1cGd0b25kUmVESU1McVR4YjhncHI0d0NnWUlLb1pJemowRUF3SXcKWkRFTE1Ba0dBMVVFQmhNQ1FWVXhFekFSQmdOVkJBZ01DbE52YldVdFUzUmhkR1V4SVRBZkJnTlZCQW9NR0VsdQpkR1Z5Ym1WMElGZHBaR2RwZEhNZ1VIUjVJRXgwWkRFZE1Cc0dBMVVFQXd3VVpHVjJMblZ3YkdJdWRYQmpiRzkxClpDNWpiMjB3SGhjTk1qRXhNREl5TVRJeE1ETTJXaGNOTXpFeE1ESXdNVEl4TURNMldqQmtNUXN3Q1FZRFZRUUcKRXdKQlZURVRNQkVHQTFVRUNBd0tVMjl0WlMxVGRHRjBaVEVoTUI4R0ExVUVDZ3dZU1c1MFpYSnVaWFFnVjJsawpaMmwwY3lCUWRIa2dUSFJrTVIwd0d3WURWUVFEREJSa1pYWXVkWEJzWWk1MWNHTnNiM1ZrTG1OdmJUQlpNQk1HCkJ5cUdTTTQ5QWdFR0NDcUdTTTQ5QXdFSEEwSUFCQmpVcUgyaVNuMFV2ZkU3UDdkT0QrSDBoKytxUWpnTG9OeWQKSTFwMmkvdlJPZmhMa0hCUjIxZ2JCSUdENjllYU1WWnZ4RWNxVWlKVWYwcmxLU3FzKzIralV6QlJNQjBHQTFVZApEZ1FXQkJTYTFaU3V1NkxJczMrc2lSSUJ5MHRXL3RnamZEQWZCZ05WSFNNRUdEQVdnQlNhMVpTdXU2TElzMytzCmlSSUJ5MHRXL3RnamZEQVBCZ05WSFJNQkFmOEVCVEFEQVFIL01Bb0dDQ3FHU000OUJBTUNBMGdBTUVVQ0lRQ3IKWXA5dHc2TmVXTHZGOGwrWm9rSE9QUzUzaGc2SDM0OHNMSjEvNit4YXN3SWdWVmN6WkFDc3JyUWt3TnVBZEVCeQo5TkxJR1VrWlhqeWgwdVFCS2x4Si9Wdz0KLS0tLS1FTkQgQ0VSVElGSUNBVEUtLS0tLQo=",
			Hostnames:        []string{"application.example.com"},
			KeyType:          "ecdsa",
			Name:             "example-manual-certificate",
			NotAfter:         timeParse("2031-10-20T12:10:36Z"),
			NotBefore:        timeParse("2021-10-22T12:10:36Z"),
			CreatedAt:        timeParse("2021-11-09T08:07:39.749472Z"),
			UpdatedAt:        timeParse("2021-11-09T08:07:39.749472Z"),
		},
		`
		{
			"certificate": "LS0LS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUNIVENDQWNPZ0F3SUJBZ0lVSWlNbzg1cGd0b25kUmVESU1McVR4YjhncHI0d0NnWUlLb1pJemowRUF3SXcKWkRFTE1Ba0dBMVVFQmhNQ1FWVXhFekFSQmdOVkJBZ01DbE52YldVdFUzUmhkR1V4SVRBZkJnTlZCQW9NR0VsdQpkR1Z5Ym1WMElGZHBaR2RwZEhNZ1VIUjVJRXgwWkRFZE1Cc0dBMVVFQXd3VVpHVjJMblZ3YkdJdWRYQmpiRzkxClpDNWpiMjB3SGhjTk1qRXhNREl5TVRJeE1ETTJXaGNOTXpFeE1ESXdNVEl4TURNMldqQmtNUXN3Q1FZRFZRUUcKRXdKQlZURVRNQkVHQTFVRUNBd0tVMjl0WlMxVGRHRjBaVEVoTUI4R0ExVUVDZ3dZU1c1MFpYSnVaWFFnVjJsawpaMmwwY3lCUWRIa2dUSFJrTVIwd0d3WURWUVFEREJSa1pYWXVkWEJzWWk1MWNHTnNiM1ZrTG1OdmJUQlpNQk1HCkJ5cUdTTTQ5QWdFR0NDcUdTTTQ5QXdFSEEwSUFCQmpVcUgyaVNuMFV2ZkU3UDdkT0QrSDBoKytxUWpnTG9OeWQKSTFwMmkvdlJPZmhMa0hCUjIxZ2JCSUdENjllYU1WWnZ4RWNxVWlKVWYwcmxLU3FzKzIralV6QlJNQjBHQTFVZApEZ1FXQkJTYTFaU3V1NkxJczMrc2lSSUJ5MHRXL3RnamZEQWZCZ05WSFNNRUdEQVdnQlNhMVpTdXU2TElzMytzCmlSSUJ5MHRXL3RnamZEQVBCZ05WSFJNQkFmOEVCVEFEQVFIL01Bb0dDQ3FHU000OUJBTUNBMGdBTUVVQ0lRQ3IKWXA5dHc2TmVXTHZGOGwrWm9rSE9QUzUzaGc2SDM0OHNMSjEvNit4YXN3SWdWVmN6WkFDc3JyUWt3TnVBZEVCeQo5TkxJR1VrWlhqeWgwdVFCS2x4Si9Wdz0KLS0tLS1FTkQgQ0VSVElGSUNBVEUtLS0tLQo=",
			"created_at": "2021-11-09T08:07:39.749472Z",
			"hostnames": [
				"application.example.com"
			],
			"key_type": "ecdsa",
			"name": "example-manual-certificate",
			"not_after": "2031-10-20T12:10:36Z",
			"not_before": "2021-10-22T12:10:36Z",
			"operational_state": "idle",
			"type": "manual",
			"updated_at": "2021-11-09T08:07:39.749472Z",
			"uuid": "bf571589-7378-41f8-879e-5505613b070d"
		}
		`,
	)
	testJSON(t,
		&LoadBalancerCertificateBundle{},
		&LoadBalancerCertificateBundle{
			UUID:             "0aded5c1-c7a3-498a-b9c8-a871611c47a2",
			Certificate:      "LStLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUN3ekNDQW1tZ0F3SUJBZ0lRS0oxNWtaSFVtb1llUUN2ZGE2YVRzVEFLQmdncWhrak9QUVFEQWpBc01Rd3cKQ2dZRFZRUUtFd05rWlhZeEhEQWFCZ05WQkFNVEUyUmxkaUJKYm5SbGNtMWxaR2xoZEdVZ1EwRXdIaGNOTWpJdwpNakV5TURVek1qUXpXaGNOTWpJd01qRXpNRFV6TXpReldqQUFNSUlCSWpBTkJna3Foa2lHOXcwQkFRRUZBQU9DCkFROEFNSUlCQ2dLQ0FRRUExTlhmUzhIdEhVbURncjB1VkptbTl0SWRZTHpYS0NrSklGejFJNDI1WVZGTlB5cGYKZkJ5MDdxMUpYejIyQSs5WFdrZEVLY21yZ3N3QU50YVZkbHFDRzU5SzVQbXhCTVRLREFTYXc5REtZcGVTNUNlMwpWVG1PRnB3T3lMbUluMFJTZmdzdU5DOWRTR1BvemdhR0VLNXVVS29SRU4rNzRzNGNkcmVVME5pM0lQTUpzeHpmClMvVXJQRTcydzU5eG5jaEs4dVNUL0pzRzBsQUt4TkRsMHRqRkw3K25zNmYvNzQxRE9YSVVRZERMekhVYktheEsKWVNKTlQ2dGNlVlJqM01RNEtaaDRlQ3pNRmJlc0V2M1UxRVRRL1ZYYkNhQXBBOTB6UWdJQURqbGNNZTV0UTFZcwplNDIvWTkrWGVYZ25ISXFvR2RzdlF3VlNXYkdlWmdsZ1orZlNCd0lEQVFBQm80SE5NSUhLTUE0R0ExVWREd0VCCi93UUVBd0lGb0RBZEJnTlZIU1VFRmpBVUJnZ3JCZ0VGQlFjREFRWUlLd1lCQlFVSEF3SXdIUVlEVlIwT0JCWUUKRkh1NVVOZGlkb0VyN0FrTTR4Um1NQzR4NUdxTU1COEdBMVVkSXdRWU1CYUFGTldtZnpER1NHOXBwZWlQWE93TgoydUx6RnVjbk1Eb0dBMVVkRVFFQi93UXdNQzZDTEd4aUxUQmhabVkyWkdGak1UUXpZelF6TURBNVlqTXpaV1V5Ck56VTJaalkxT1RKa0xtRndjR3hsWW1WbE1CMEdEQ3NHQVFRQmdxUmt4aWhBQVFRTk1Bc0NBUVlFQkdGamJXVUUKQURBS0JnZ3Foa2pPUFFRREFnTklBREJGQWlFQTBUcXhjNlVVejREVlh1OFJwSFhxS2R3VHVoSkNkSGxOMVA1agpNekkvMTFjQ0lEb2ZoelkwSFcxUHlMbEgrckVydW90cC9FYm5IMElLaDEzQTM4dndTQnM2Ci0tLS0tRU5EIENFUlRJRklDQVRFLS0tLS0K",
			Intermediates:    "LS0tS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUJ0VENDQVZ1Z0F3SUJBZ0lSQU5wSDZzV0ZtQzErWkdnUzFMWllVZGN3Q2dZSUtvWkl6ajBFQXdJd0pERU0KTUFvR0ExVUVDaE1EWkdWMk1SUXdFZ1lEVlFRREV3dGtaWFlnVW05dmRDQkRRVEFlRncweU1URXlNRGt4TXpVMwpNREZhRncwek1URXlNRGN4TXpVM01ERmFNQ3d4RERBS0JnTlZCQW9UQTJSbGRqRWNNQm9HQTFVRUF4TVRaR1YyCklFbHVkR1Z5YldWa2FXRjBaU0JEUVRCWk1CTUdCeXFHU000OUFnRUdDQ3FHU000OUF3RUhBMElBQkswbGMzNmcKN01TaDJTaXd3MUdDUjkvL3lSODR6S1VuNml6SmdCUkpFTlBxbmNXcjQzTi8rNktJR1EraERaazhRWHZ6RmExYQp2dFloc3JEVGtnRm9EV0tqWmpCa01BNEdBMVVkRHdFQi93UUVBd0lCQmpBU0JnTlZIUk1CQWY4RUNEQUdBUUgvCkFnRUFNQjBHQTFVZERnUVdCQlRWcG44d3hraHZhYVhvajF6c0Rkcmk4eGJuSnpBZkJnTlZIU01FR0RBV2dCU2oKckgwV0pubDdUSUJtc3NESGVveENFTVZyRmpBS0JnZ3Foa2pPUFFRREFnTklBREJGQWlBa3NhUXdPMkFESGhBLwppRVR1SVY1dTlNV3hFTU5BVGlVODFIZjc0cGVhWlFJaEFLMnJDRmhVVnQxbFlzR1o3dFdjWGFHVDhyU1k2cU1YClBmK3dnUXFnNXUyVAotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0tCg==",
			Hostnames:        []string{"web.example.com"},
			KeyType:          "rsa",
			Name:             "example-dynamic-certificate",
			NotAfter:         timeParse("2022-02-13T05:33:43Z"),
			NotBefore:        timeParse("2022-02-12T05:32:43Z"),
			CreatedAt:        timeParse("2022-02-11T17:31:38.202398Z"),
			UpdatedAt:        timeParse("2022-02-12T08:13:47.877562Z"),
			Type:             LoadBalancerCertificateBundleTypeDynamic,
			OperationalState: LoadBalancerCertificateBundleOperationalStateIdle,
		},
		`
		{
			"certificate": "LStLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUN3ekNDQW1tZ0F3SUJBZ0lRS0oxNWtaSFVtb1llUUN2ZGE2YVRzVEFLQmdncWhrak9QUVFEQWpBc01Rd3cKQ2dZRFZRUUtFd05rWlhZeEhEQWFCZ05WQkFNVEUyUmxkaUJKYm5SbGNtMWxaR2xoZEdVZ1EwRXdIaGNOTWpJdwpNakV5TURVek1qUXpXaGNOTWpJd01qRXpNRFV6TXpReldqQUFNSUlCSWpBTkJna3Foa2lHOXcwQkFRRUZBQU9DCkFROEFNSUlCQ2dLQ0FRRUExTlhmUzhIdEhVbURncjB1VkptbTl0SWRZTHpYS0NrSklGejFJNDI1WVZGTlB5cGYKZkJ5MDdxMUpYejIyQSs5WFdrZEVLY21yZ3N3QU50YVZkbHFDRzU5SzVQbXhCTVRLREFTYXc5REtZcGVTNUNlMwpWVG1PRnB3T3lMbUluMFJTZmdzdU5DOWRTR1BvemdhR0VLNXVVS29SRU4rNzRzNGNkcmVVME5pM0lQTUpzeHpmClMvVXJQRTcydzU5eG5jaEs4dVNUL0pzRzBsQUt4TkRsMHRqRkw3K25zNmYvNzQxRE9YSVVRZERMekhVYktheEsKWVNKTlQ2dGNlVlJqM01RNEtaaDRlQ3pNRmJlc0V2M1UxRVRRL1ZYYkNhQXBBOTB6UWdJQURqbGNNZTV0UTFZcwplNDIvWTkrWGVYZ25ISXFvR2RzdlF3VlNXYkdlWmdsZ1orZlNCd0lEQVFBQm80SE5NSUhLTUE0R0ExVWREd0VCCi93UUVBd0lGb0RBZEJnTlZIU1VFRmpBVUJnZ3JCZ0VGQlFjREFRWUlLd1lCQlFVSEF3SXdIUVlEVlIwT0JCWUUKRkh1NVVOZGlkb0VyN0FrTTR4Um1NQzR4NUdxTU1COEdBMVVkSXdRWU1CYUFGTldtZnpER1NHOXBwZWlQWE93TgoydUx6RnVjbk1Eb0dBMVVkRVFFQi93UXdNQzZDTEd4aUxUQmhabVkyWkdGak1UUXpZelF6TURBNVlqTXpaV1V5Ck56VTJaalkxT1RKa0xtRndjR3hsWW1WbE1CMEdEQ3NHQVFRQmdxUmt4aWhBQVFRTk1Bc0NBUVlFQkdGamJXVUUKQURBS0JnZ3Foa2pPUFFRREFnTklBREJGQWlFQTBUcXhjNlVVejREVlh1OFJwSFhxS2R3VHVoSkNkSGxOMVA1agpNekkvMTFjQ0lEb2ZoelkwSFcxUHlMbEgrckVydW90cC9FYm5IMElLaDEzQTM4dndTQnM2Ci0tLS0tRU5EIENFUlRJRklDQVRFLS0tLS0K",
			"created_at": "2022-02-11T17:31:38.202398Z",
			"hostnames": [
				"web.example.com"
			],
			"intermediates": "LS0tS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUJ0VENDQVZ1Z0F3SUJBZ0lSQU5wSDZzV0ZtQzErWkdnUzFMWllVZGN3Q2dZSUtvWkl6ajBFQXdJd0pERU0KTUFvR0ExVUVDaE1EWkdWMk1SUXdFZ1lEVlFRREV3dGtaWFlnVW05dmRDQkRRVEFlRncweU1URXlNRGt4TXpVMwpNREZhRncwek1URXlNRGN4TXpVM01ERmFNQ3d4RERBS0JnTlZCQW9UQTJSbGRqRWNNQm9HQTFVRUF4TVRaR1YyCklFbHVkR1Z5YldWa2FXRjBaU0JEUVRCWk1CTUdCeXFHU000OUFnRUdDQ3FHU000OUF3RUhBMElBQkswbGMzNmcKN01TaDJTaXd3MUdDUjkvL3lSODR6S1VuNml6SmdCUkpFTlBxbmNXcjQzTi8rNktJR1EraERaazhRWHZ6RmExYQp2dFloc3JEVGtnRm9EV0tqWmpCa01BNEdBMVVkRHdFQi93UUVBd0lCQmpBU0JnTlZIUk1CQWY4RUNEQUdBUUgvCkFnRUFNQjBHQTFVZERnUVdCQlRWcG44d3hraHZhYVhvajF6c0Rkcmk4eGJuSnpBZkJnTlZIU01FR0RBV2dCU2oKckgwV0pubDdUSUJtc3NESGVveENFTVZyRmpBS0JnZ3Foa2pPUFFRREFnTklBREJGQWlBa3NhUXdPMkFESGhBLwppRVR1SVY1dTlNV3hFTU5BVGlVODFIZjc0cGVhWlFJaEFLMnJDRmhVVnQxbFlzR1o3dFdjWGFHVDhyU1k2cU1YClBmK3dnUXFnNXUyVAotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0tCg==",
			"key_type": "rsa",
			"name": "example-dynamic-certificate",
			"not_after": "2022-02-13T05:33:43Z",
			"not_before": "2022-02-12T05:32:43Z",
			"operational_state": "idle",
			"type": "dynamic",
			"updated_at": "2022-02-12T08:13:47.877562Z",
			"uuid": "0aded5c1-c7a3-498a-b9c8-a871611c47a2"
		}
		`,
	)
}

func TestLoadBalancerNetworks(t *testing.T) {
	t.Parallel()
	testJSON(t,
		&LoadBalancer{},
		&LoadBalancer{
			Networks: []LoadBalancerNetwork{
				{
					UUID:      "0aded5c1-c7a3-498a-b9c8-a871611c47a2",
					Name:      "PublicNet",
					DNSName:   "public.name",
					Type:      LoadBalancerNetworkTypePublic,
					Family:    LoadBalancerAddressFamilyIPv4,
					CreatedAt: timeParse("2021-12-07T13:58:30.817272Z"),
					UpdatedAt: timeParse("2022-02-11T17:33:08.490581Z"),
				},
				{
					UUID:      "bf571589-7378-41f8-879e-5505613b070d",
					Name:      "PrivateNet",
					DNSName:   "private.name",
					Type:      LoadBalancerNetworkTypePrivate,
					Family:    LoadBalancerAddressFamilyIPv4,
					CreatedAt: timeParse("2021-12-07T13:58:30.817272Z"),
					UpdatedAt: timeParse("2022-02-11T17:33:08.490581Z"),
				},
			},
			CreatedAt: timeParse("2021-12-07T13:58:30.817272Z"),
			UpdatedAt: timeParse("2022-02-11T17:33:08.490581Z"),
		},
		`
		{
			"networks": [
				{
					"uuid": "0aded5c1-c7a3-498a-b9c8-a871611c47a2",
					"name": "PublicNet",
					"dns_name": "public.name",
					"type": "public",
					"family": "IPv4",
					"created_at": "2021-12-07T13:58:30.817272Z",
					"updated_at": "2022-02-11T17:33:08.490581Z"
				},
				{
					"uuid": "bf571589-7378-41f8-879e-5505613b070d",
					"name": "PrivateNet",
					"dns_name": "private.name",
					"type": "private",
					"family": "IPv4",
					"created_at": "2021-12-07T13:58:30.817272Z",
					"updated_at": "2022-02-11T17:33:08.490581Z"
				}
			],
			"created_at": "2021-12-07T13:58:30.817272Z",
			"updated_at": "2022-02-11T17:33:08.490581Z"
		}`)
}

func testJSON(t *testing.T, unMarshall, marshall interface{}, want string) {
	got, err := json.Marshal(marshall)
	assert.NoError(t, err)
	assert.JSONEq(t, want, string(got))
	err = json.Unmarshal([]byte(want), unMarshall)
	assert.NoError(t, err)
	assert.Equal(t, marshall, unMarshall)
}
