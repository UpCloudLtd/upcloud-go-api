package service

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/UpCloudLtd/upcloud-go-api/v5/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v5/upcloud/client"
	"github.com/UpCloudLtd/upcloud-go-api/v5/upcloud/request"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupGatewayTest(handler http.Handler) (*httptest.Server, *Service) {
	srv := httptest.NewServer(handler)
	return srv, New(client.New("user", "pass", client.WithBaseURL(srv.URL)))
}

func TestGetGateways(t *testing.T) {
	t.Parallel()

	srv, svc := setupGatewayTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, fmt.Sprintf("/%s%s", client.APIVersion, "/gateway"), r.URL.Path)
		fmt.Fprintf(w, "[%s]", gatewayCommonResponse)
	}))
	defer srv.Close()
	gw, err := svc.GetGateways(context.TODO())
	assert.NoError(t, err)
	assert.Len(t, gw, 1)
	checkGatewayResponse(t, &gw[0])
}

func TestGetGateway(t *testing.T) {
	t.Parallel()

	srv, svc := setupGatewayTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, fmt.Sprintf("/%s%s", client.APIVersion, "/gateway/_UUID_"), r.URL.Path)
		fmt.Fprint(w, gatewayCommonResponse)
	}))
	defer srv.Close()
	gw, err := svc.GetGateway(context.TODO(), &request.GetGatewayRequest{UUID: "_UUID_"})
	assert.NoError(t, err)
	checkGatewayResponse(t, gw)
}

func TestCreateGateway(t *testing.T) {
	t.Parallel()

	srv, svc := setupGatewayTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, fmt.Sprintf("/%s%s", client.APIVersion, "/gateway"), r.URL.Path)
		body, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		assert.JSONEq(t, `
		{
			"name": "test-create",
			"zone": "fi-hel1",
			"features": ["nat"],
			"routers": [{ "uuid": "router-uuid" }],
			"labels": [{ "key": "test", "value": "Create request" }],
			"configured_status": "started"
		}`, string(body))
		fmt.Fprint(w, gatewayCommonResponse)
	}))
	defer srv.Close()

	p, err := svc.CreateGateway(context.TODO(), &request.CreateGatewayRequest{
		Name:             "test-create",
		Zone:             "fi-hel1",
		Features:         []upcloud.GatewayFeature{upcloud.GatewayFeatureNAT},
		Routers:          []request.GatewayRouter{{UUID: "router-uuid"}},
		Labels:           []upcloud.Label{{Key: "test", Value: "Create request"}},
		ConfiguredStatus: upcloud.GatewayStatusStarted,
	})
	if !assert.NoError(t, err) {
		return
	}
	checkGatewayResponse(t, p)
}

func TestModifyGateway(t *testing.T) {
	t.Parallel()

	srv, svc := setupGatewayTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPatch, r.Method)
		assert.Equal(t, fmt.Sprintf("/%s%s", client.APIVersion, "/gateway/_UUID_"), r.URL.Path)
		body, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		assert.JSONEq(t, `
		{
			"name": "test-modify",
			"configured_status": "stopped",
			"labels": [{ "key": "test", "value": "Modify request" }]
		}
		`, string(body))
		fmt.Fprint(w, gatewayCommonResponse)
	}))
	defer srv.Close()

	p, err := svc.ModifyGateway(context.TODO(), &request.ModifyGatewayRequest{
		UUID:             "_UUID_",
		Name:             "test-modify",
		ConfiguredStatus: upcloud.GatewayStatusStopped,
		Labels:           []upcloud.Label{{Key: "test", Value: "Modify request"}},
	})
	if !assert.NoError(t, err) {
		return
	}
	checkGatewayResponse(t, p)
}

func TestDeleteGateway(t *testing.T) {
	t.Parallel()

	srv, svc := setupGatewayTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		assert.Equal(t, fmt.Sprintf("/%s%s", client.APIVersion, "/gateway/_UUID_"), r.URL.Path)
	}))
	defer srv.Close()
	assert.NoError(t, svc.DeleteGateway(context.TODO(), &request.DeleteGatewayRequest{UUID: "_UUID_"}))
}

func checkGatewayResponse(t *testing.T, gw *upcloud.Gateway) {
	assert.Equal(t, "10c153e0-12e4-4dea-8748-4f34850ff76d", gw.UUID)
	assert.Equal(t, "0485d477-8d8f-4c97-9bef-731933187538", gw.Routers[0].UUID)
	assert.Equal(t, upcloud.GatewayStatusStarted, gw.ConfiguredStatus)
}

const gatewayCommonResponse string = `
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
