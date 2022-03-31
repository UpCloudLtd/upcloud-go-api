package request

import (
	"encoding/json"
	"testing"

	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud"
	"github.com/stretchr/testify/assert"
)

func TestMatcheresAndActionsHelper(t *testing.T) {
	isTrue := true

	expected := `
	{
		"name": "rule-name",
		"priority": 0,
		"matchers": [
			{
				"type": "num_members_up",
				"match_num_members_up": {
					"method": "less",
					"value": 1,
					"backend": "example-fallback-backend"
				}
			},
			{
				"type": "path",
				"match_path": {
					"method": "exact",
					"value": "/app"
				}
			},
			{
				"type": "url_param",
				"match_url_param": {
					"method": "exact",
					"name": "status",
					"value": "active",
					"ignore_case": true
				}
			},
			{
				"type": "header",
				"match_header": {
					"method": "exists",
					"name": "X-Test"
				}
			},
			{
				"type": "cookie",
				"match_cookie": {
					"method": "exists",
					"name": "x-session-id"
				}
			},
			{
				"type": "http_method",
				"match_http_method": {
					"value": "PATCH"
				}
			},
			{
				"type": "host",
				"match_host": {
					"value": "example.com"
				}
			},
			{
				"type": "url_query",
				"match_url_query": {
					"method": "exists"
				}
			},
			{
				"type": "url",
				"match_url": {
					"method": "starts",
					"value": "/app",
					"ignore_case": true
				}
			},
			{
				"type": "src_port",
				"match_src_port": {
					"method": "range",
					"range_start": 8000,
					"range_end": 9000
				}
			},
			{
				"type": "src_port",
				"match_src_port": {
					"method": "equal",
					"value": 8000
				}
			},
			{
				"type": "body_size",
				"match_body_size": {
					"method": "range",
					"range_start": 8000,
					"range_end": 9000
				}
			},
			{
				"type": "body_size",
				"match_body_size": {
					"method": "equal",
					"value": 8000
				}
			},
			{
				"type": "src_ip",
				"match_src_ip": {
					"value": "127.0.0.1"
				}
			}
		],
		"actions": [
			{
				"type": "use_backend",
				"action_use_backend": {
					"backend": "example-backend-2"
				}
			},
			{
				"type": "tcp_reject",
				"action_tcp_reject": {}
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
				"type": "http_redirect",
				"action_http_redirect": {
					"location": "https://internal.example.com"
				}
			}	
		]
	}
	`
	r := LoadBalancerFrontendRule{
		Name:     "rule-name",
		Priority: 0,
		Matchers: []upcloud.LoadBalancerMatcher{
			NewLoadBalancerNumMembersUPMatcher(upcloud.LoadBalancerIntegerMatcherMethodLess, 1, "example-fallback-backend"),
			NewLoadBalancerPathMatcher(upcloud.LoadBalancerStringMatcherMethodExact, "/app", nil),
			NewLoadBalancerURLParamMatcher(upcloud.LoadBalancerStringMatcherMethodExact, "status", "active", &isTrue),
			NewLoadBalancerHeaderMatcher(upcloud.LoadBalancerStringMatcherMethodExists, "X-Test", "", nil),
			NewLoadBalancerCookieMatcher(upcloud.LoadBalancerStringMatcherMethodExists, "x-session-id", "", nil),
			NewLoadBalancerHTTPMethodMatcher(upcloud.LoadBalancerHTTPMatcherMethodPatch),
			NewLoadBalancerHostMatcher("example.com"),
			NewLoadBalancerURLQueryMatcher(upcloud.LoadBalancerStringMatcherMethodExists, "", nil),
			NewLoadBalancerURLMatcher(upcloud.LoadBalancerStringMatcherMethodStarts, "/app", &isTrue),
			NewLoadBalancerSrcPortRangeMatcher(8000, 9000),
			NewLoadBalancerSrcPortMatcher(upcloud.LoadBalancerIntegerMatcherMethodEqual, 8000),
			NewLoadBalancerBodySizeRangeMatcher(8000, 9000),
			NewLoadBalancerBodySizeMatcher(upcloud.LoadBalancerIntegerMatcherMethodEqual, 8000),
			NewLoadBalancerSrcIPMatcher("127.0.0.1"),
		},
		Actions: []upcloud.LoadBalancerAction{
			NewLoadBalancerUseBackendAction("example-backend-2"),
			NewLoadBalancerTCPRejectAction(),
			NewLoadBalancerHTTPReturnAction(200, "text/html", "PGgxPmFwcGxlYmVlPC9oMT4K"),
			NewLoadBalancerHTTPRedirectAction("https://internal.example.com"),
		},
	}
	actual, err := json.Marshal(&r)
	assert.NoError(t, err)
	assert.JSONEq(t, expected, string(actual))
}

func TestNewLoadBalancerBackendMemberHelper(t *testing.T) {
	expected := `
	[
		{
			"name": "mem",
			"weight": 0,
			"max_sessions": 0,
			"enabled": true,
			"type": "static",
			"ip": "10.0.0.1",
			"port": 80
		},
		{
			"name": "mem",
			"weight": 0,
			"max_sessions": 0,
			"enabled": true,
			"type": "dynamic",
			"ip": "10.0.0.1",
			"port": 80
		}
	]`
	r := []LoadBalancerBackendMember{
		NewLoadBalancerStaticBackendMember("mem", 0, 0, true, "10.0.0.1", 80),
		NewLoadBalancerDynamicBackendMember("mem", 0, 0, true, "10.0.0.1", 80),
	}
	actual, err := json.Marshal(&r)
	assert.NoError(t, err)
	assert.JSONEq(t, expected, string(actual))
}

func TestCreateLoadBalancerCertificateBundleRequest(t *testing.T) {
	expected := `
	[
		{
			"name": "man",
			"type": "manual",
			"certificate": "x",
			"intermediates": "x",
			"private_key": "x"
		},
		{
			"name": "dyn",
			"type": "dynamic",
			"hostnames": [
				"example.com",
				"app.example.com"
			],
			"key_type": "rsa"
		}
	]
	`
	r := []CreateLoadBalancerCertificateBundleRequest{
		NewCreateLoadBalancerManualCertificateBundleRequest("man", "x", "x", "x"),
		NewCreateLoadBalancerDynamicCertificateBundleRequest("dyn", "rsa", []string{"example.com", "app.example.com"}),
	}
	actual, err := json.Marshal(&r)
	assert.NoError(t, err)
	assert.JSONEq(t, expected, string(actual))
}
