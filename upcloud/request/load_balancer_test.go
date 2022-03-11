package request

import (
	"encoding/json"
	"testing"

	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateLoadBalancerBackendRequest(t *testing.T) {
	r := CreateLoadBalancerBackendRequest{
		ServiceUUID: "lb",
		Payload: CreateLoadBalancerBackend{
			Name:     "sesese",
			Members:  []upcloud.LoadBalancerBackendMember{},
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
	actualJson, err := json.Marshal(r)

	require.NoError(t, err)
	assert.Exactly(t, "/loadbalancer/lb/backends", r.RequestURL())
	assert.JSONEq(t, expectedJson, string(actualJson))
}

func TestGetLoadBalancerBackendDetailsRequest(t *testing.T) {
	r := GetLoadBalancerBackendDetailsRequest{
		ServiceUUID: "lb",
		BackendName: "be",
	}

	expectedJson := "{}"
	actualJson, err := json.Marshal(r)

	require.NoError(t, err)
	assert.Exactly(t, "/loadbalancer/lb/backends/be", r.RequestURL())
	assert.JSONEq(t, expectedJson, string(actualJson))
}

func TestModifyLoadBalancerBackendRequest(t *testing.T) {
	r := ModifyLoadBalancerBackendRequest{
		ServiceUUID: "lb",
		Name:        "be",
		Payload: ModifyLoadBalancerBackend{
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
		BackendName: "be",
	}

	expectedJson := "{}"
	actualJson, err := json.Marshal(r)

	require.NoError(t, err)
	assert.Exactly(t, "/loadbalancer/lb/backends/be", r.RequestURL())
	assert.JSONEq(t, expectedJson, string(actualJson))
}

func TestCreateLoadBalancerBackendMember(t *testing.T) {
	r := CreateLoadBalancerBackendMemberRequest{
		ServiceUUID:       "lb",
		BackendName:       "be",
		MemberName:        "mem",
		MemberWeight:      100,
		MemberMaxSessions: 5,
		MemberEnabled:     true,
		MemberType:        "static",
		MemberIP:          "10.0.0.1",
		MemberPort:        80,
		MemberServerUUID:  "serv",
	}

	expectedJson := `
	{
		"name": "mem",
		"weight": 100,
		"max_sessions": 5,
		"enabled": true,
		"type": "static",
		"ip": "10.0.0.1",
		"port": 80,
		"server_uuid": "serv"
	}`

	actualJson, err := json.Marshal(r)

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
	actualJson, err := json.Marshal(r)

	require.NoError(t, err)
	assert.Exactly(t, "/loadbalancer/lb/backends/be/members", r.RequestURL())
	assert.JSONEq(t, expectedJson, string(actualJson))
}

func TestGetLoadBalancerBackendMemberDetailsRequest(t *testing.T) {
	r := GetLoadBalancerBackendMemberDetailsRequest{
		ServiceUUID: "lb",
		BackendName: "be",
		MemberName:  "mem",
	}

	expectedJson := "{}"
	actualJson, err := json.Marshal(r)

	require.NoError(t, err)
	assert.Exactly(t, "/loadbalancer/lb/backends/be/members/mem", r.RequestURL())
	assert.JSONEq(t, expectedJson, string(actualJson))
}

func TestModifyLoadBalancerBackendMemberRequest(t *testing.T) {
	r := ModifyLoadBalancerBackendMemberRequest{
		ServiceUUID:       "lb",
		BackendName:       "be",
		MemberName:        "mem",
		NewMemberName:     "newmem",
		MemberWeight:      100,
		MemberMaxSessions: 5,
		MemberEnabled:     true,
		MemberType:        "static",
		MemberIP:          "10.0.0.1",
		MemberPort:        80,
		MemberServerUUID:  "serv",
	}

	expectedJson := `
	{
		"name": "newmem",
		"weight": 100,
		"max_sessions": 5,
		"enabled": true,
		"type": "static",
		"ip": "10.0.0.1",
		"port": 80,
		"server_uuid": "serv"
	}`

	actualJson, err := json.Marshal(r)

	require.NoError(t, err)
	assert.Exactly(t, "/loadbalancer/lb/backends/be/members/mem", r.RequestURL())
	assert.JSONEq(t, expectedJson, string(actualJson))
}

func TestDeleteLoadBalancerBackendMemberRequest(t *testing.T) {
	r := DeleteLoadBalancerBackendMemberRequest{
		ServiceUUID: "lb",
		BackendName: "be",
		MemberName:  "mem",
	}

	expectedJson := "{}"
	actualJson, err := json.Marshal(r)

	require.NoError(t, err)
	assert.Exactly(t, "/loadbalancer/lb/backends/be/members/mem", r.RequestURL())
	assert.JSONEq(t, expectedJson, string(actualJson))
}

func TestCreateLoadBalancerResolverRequest(t *testing.T) {
	r := CreateLoadBalancerResolverRequest{
		ServiceUUID:  "service-uuid",
		Name:         "testname",
		Nameservers:  []string{"10.0.0.0", "10.0.0.1"},
		Retries:      5,
		TimeoutRetry: 10,
		Timeout:      20,
		CacheValid:   123,
		CacheInvalid: 321,
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

	actualJson, err := json.Marshal(r)

	require.NoError(t, err)
	assert.EqualValues(t, "/loadbalancer/service-uuid/resolvers", r.RequestURL())
	assert.JSONEq(t, expectedJson, string(actualJson))
}

func TestGetLoadBalancerResolversRequest(t *testing.T) {
	r := GetLoadBalancerResolversRequest{
		ServiceUUID: "service-uuid",
	}

	expectedJson := "{}"
	actualJson, err := json.Marshal(r)

	require.NoError(t, err)
	assert.EqualValues(t, "/loadbalancer/service-uuid/resolvers", r.RequestURL())
	assert.JSONEq(t, expectedJson, string(actualJson))
}

func TestGetLoadBalancerResolverDetailsRequest(t *testing.T) {
	r := GetLoadBalancerResolverDetailsRequest{
		ServiceUUID:  "service-uuid",
		ResolverName: "sesese",
	}

	expectedJson := "{}"
	actualJson, err := json.Marshal(r)

	require.NoError(t, err)
	assert.EqualValues(t, "/loadbalancer/service-uuid/resolvers/sesese", r.RequestURL())
	assert.JSONEq(t, expectedJson, string(actualJson))
}

func TestModifyLoadBalancerResolverRequest(t *testing.T) {
	r := ModifyLoadBalancerRevolverRequest{
		ServiceUUID:     "service-uuid",
		ResolverName:    "sesese",
		NewResolverName: "testname",
		Nameservers:     []string{"10.0.0.0", "10.0.0.1"},
		Retries:         5,
		TimeoutRetry:    10,
		Timeout:         20,
		CacheValid:      123,
		CacheInvalid:    321,
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

	actualJson, err := json.Marshal(r)

	require.NoError(t, err)
	assert.EqualValues(t, "/loadbalancer/service-uuid/resolvers/sesese", r.RequestURL())
	assert.JSONEq(t, expectedJson, string(actualJson))
}

func TestDeleteLoadBalancerResolverRequest(t *testing.T) {
	r := DeleteLoadBalancerResolverRequest{
		ServiceUUID:  "service-uuid",
		ResolverName: "sesese",
	}

	expectedJson := "{}"
	actualJson, err := json.Marshal(r)

	require.NoError(t, err)
	assert.EqualValues(t, "/loadbalancer/service-uuid/resolvers/sesese", r.RequestURL())
	assert.JSONEq(t, expectedJson, string(actualJson))
}
