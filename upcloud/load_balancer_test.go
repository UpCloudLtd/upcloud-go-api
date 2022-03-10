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
				TLSConfigs: []LoadBalancerTLSConfig{{
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

func TestLoadBalancerPlan(t *testing.T) {
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
	testJSON(t,
		&LoadBalancerFrontend{},
		&LoadBalancerFrontend{
			Name:           "example-frontend",
			Mode:           LoadBalancerModeTCP,
			Port:           443,
			DefaultBackend: "example-backend",
			CreatedAt:      timeParse("2021-12-07T13:58:30.817272Z"),
			UpdatedAt:      timeParse("2022-02-11T17:33:08.490581Z"),
		},
		`
		{
			"name": "example-frontend",
			"mode": "tcp",
			"port": 443,
			"default_backend": "example-backend",
			"created_at": "2021-12-07T13:58:30.817272Z",
			"updated_at": "2022-02-11T17:33:08.490581Z"
		}
		`,
	)
}

func TestLoadBalancerRule(t *testing.T) {
	testJSON(t,
		&LoadBalancerRule{},
		&LoadBalancerRule{
			Name:     "example-rule-1",
			Priority: 100,
			Actions: []LoadBalancerAction{
				{
					Type:      LoadBalancerActionTypeTCPReject,
					TCPReject: &LoadBalancerActionTCPReject{},
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
				}
			]
		}
		`,
	)
}

func TestTLSConfig(t *testing.T) {
	testJSON(t,
		&LoadBalancerTLSConfig{},
		&LoadBalancerTLSConfig{
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
	testJSON(t,
		&LoadBalancerBackend{},
		&LoadBalancerBackend{
			Name:      "example-backend-2",
			Members:   []LoadBalancerBackendMember{},
			CreatedAt: timeParse("2022-02-11T17:33:08.490581Z"),
			UpdatedAt: timeParse("2022-02-11T17:33:08.490581Z"),
		},
		`
		{
			"name": "example-backend-2",
			"created_at": "2022-02-11T17:33:08.490581Z",
			"updated_at": "2022-02-11T17:33:08.490581Z",
			"members": []
		}
		`,
	)
}

func TestLoadBalancerBackendMember(t *testing.T) {
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

func TestLoadBalancerMatcherNumMembersUP(t *testing.T) {
	testJSON(t,
		&LoadBalancerMatcher{},
		&LoadBalancerMatcher{
			Type: LoadBalancerMatcherTypeNumMembersUP,
			NumMembersUP: &LoadBalancerMatcherNumMembersUP{
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

func TestLoadBalancerActionUseBackend(t *testing.T) {
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

func testJSON(t *testing.T, unMarshall, marshall interface{}, want string) {
	got, err := json.Marshal(marshall)
	assert.NoError(t, err)
	assert.JSONEq(t, want, string(got))
	err = json.Unmarshal([]byte(want), unMarshall)
	assert.NoError(t, err)
	assert.Equal(t, marshall, unMarshall)
}
