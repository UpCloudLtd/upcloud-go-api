package request

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
		"timout":20,
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
		"timout":20,
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
