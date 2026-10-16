package upcloud

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

type generatedErrorResponse interface {
	StatusCode() int
	ContentType() string
	GetBody() []byte
}

type generatedErrorCall func(context.Context, *ClientWithResponses) (generatedErrorResponse, error)

func TestRepresentativeGeneratedErrorFamilies(t *testing.T) {
	legacyBody := `{"error":{"error_code":"FIXTURE_ERROR","error_message":"fixture failure"}}`
	testCases := []struct {
		name        string
		path        string
		status      int
		contentType string
		body        string
		call        generatedErrorCall
		assertTyped func(*testing.T, generatedErrorResponse)
	}{
		{
			name:        "server legacy 403",
			path:        "/1.3/server",
			status:      http.StatusForbidden,
			contentType: "application/json",
			body:        `{"error":{"code":"SERVER_FORBIDDEN","message":"fixture denied"}}`,
			call: func(ctx context.Context, client *ClientWithResponses) (generatedErrorResponse, error) {
				return client.ListServersWithResponse(ctx, nil)
			},
			assertTyped: func(t *testing.T, response generatedErrorResponse) {
				serverResponse, ok := response.(*ListServersResp)
				require.True(t, ok)
				require.NotNil(t, serverResponse.JSON403)
				require.Equal(t, "SERVER_FORBIDDEN", serverResponse.JSON403.Error.Code)
				require.Equal(t, "fixture denied", serverResponse.JSON403.Error.Message)
			},
		},
		{
			name:        "network peering problem default",
			path:        "/1.3/network-peering",
			status:      http.StatusNotFound,
			contentType: "application/problem+json",
			body:        legacyBody,
			call: func(ctx context.Context, client *ClientWithResponses) (generatedErrorResponse, error) {
				return client.ListNetworkPeeringsWithResponse(ctx, nil)
			},
			assertTyped: func(t *testing.T, response generatedErrorResponse) {
				networkResponse, ok := response.(*ListNetworkPeeringsResp)
				require.True(t, ok)
				require.NotNil(t, networkResponse.ApplicationproblemJSONDefault)
				require.Equal(t, "FIXTURE_ERROR", networkResponse.ApplicationproblemJSONDefault.Error.ErrorCode)
			},
		},
		{
			name:        "zone problem default",
			path:        "/1.3/zone",
			status:      http.StatusInternalServerError,
			contentType: "application/problem+json",
			body:        legacyBody,
			call: func(ctx context.Context, client *ClientWithResponses) (generatedErrorResponse, error) {
				return client.ListZonesWithResponse(ctx)
			},
			assertTyped: func(t *testing.T, response generatedErrorResponse) {
				zoneResponse, ok := response.(*ListZonesResp)
				require.True(t, ok)
				require.NotNil(t, zoneResponse.ApplicationproblemJSONDefault)
				require.Equal(t, "fixture failure", zoneResponse.ApplicationproblemJSONDefault.Error.ErrorMessage)
			},
		},
		{
			name:        "database error 503",
			path:        "/1.3/database/plans",
			status:      http.StatusServiceUnavailable,
			contentType: "application/problem+json",
			body:        legacyBody,
			call: func(ctx context.Context, client *ClientWithResponses) (generatedErrorResponse, error) {
				return client.ListDatabasePlansWithResponse(ctx)
			},
			assertTyped: func(t *testing.T, response generatedErrorResponse) {
				databaseResponse, ok := response.(*ListDatabasePlansResp)
				require.True(t, ok)
				require.NotNil(t, databaseResponse.ApplicationproblemJSONDefault)
				require.NotNil(t, databaseResponse.ApplicationproblemJSONDefault.Error)
				require.Equal(t, "FIXTURE_ERROR", databaseResponse.ApplicationproblemJSONDefault.Error.ErrorCode)
			},
		},
		{
			name:        "file storage problem 404",
			path:        "/1.3/file-storage",
			status:      http.StatusNotFound,
			contentType: "application/problem+json",
			body:        `{"correlation_id":"fixture-correlation","status":404,"title":"missing file storage","type":"https://example.test/errors/not-found"}`,
			call: func(ctx context.Context, client *ClientWithResponses) (generatedErrorResponse, error) {
				return client.ListFileStoragesWithResponse(ctx, nil)
			},
			assertTyped: func(t *testing.T, response generatedErrorResponse) {
				fileStorageResponse, ok := response.(*ListFileStoragesResp)
				require.True(t, ok)
				require.NotNil(t, fileStorageResponse.ApplicationproblemJSONDefault)
				require.Equal(t, "missing file storage", fileStorageResponse.ApplicationproblemJSONDefault.Title)
				require.Equal(t, int32(http.StatusNotFound), fileStorageResponse.ApplicationproblemJSONDefault.Status)
			},
		},
		{
			name:        "firewall problem 409",
			path:        "/1.3/firewall-ruleset",
			status:      http.StatusConflict,
			contentType: "application/problem+json",
			body:        `{"correlation_id":"fixture-correlation","status":409,"title":"fixture conflict","type":"https://example.test/errors/conflict"}`,
			call: func(ctx context.Context, client *ClientWithResponses) (generatedErrorResponse, error) {
				return client.ListFirewallRulesetsWithResponse(ctx, nil)
			},
			assertTyped: func(t *testing.T, response generatedErrorResponse) {
				firewallResponse, ok := response.(*ListFirewallRulesetsResp)
				require.True(t, ok)
				require.NotNil(t, firewallResponse.ApplicationproblemJSONDefault)
				require.Equal(t, "fixture conflict", firewallResponse.ApplicationproblemJSONDefault.Title)
				require.Equal(t, int32(http.StatusConflict), firewallResponse.ApplicationproblemJSONDefault.Status)
			},
		},
		{
			name:        "load balancer error 502",
			path:        "/1.3/load-balancer",
			status:      http.StatusBadGateway,
			contentType: "application/problem+json",
			body:        legacyBody,
			call: func(ctx context.Context, client *ClientWithResponses) (generatedErrorResponse, error) {
				return client.ListLoadBalancersWithResponse(ctx, nil)
			},
			assertTyped: func(t *testing.T, response generatedErrorResponse) {
				loadBalancerResponse, ok := response.(*ListLoadBalancersResp)
				require.True(t, ok)
				require.NotNil(t, loadBalancerResponse.ApplicationproblemJSONDefault)
				require.NotNil(t, loadBalancerResponse.ApplicationproblemJSONDefault.Error)
				require.Equal(t, "FIXTURE_ERROR", loadBalancerResponse.ApplicationproblemJSONDefault.Error.ErrorCode)
			},
		},
		{
			name:        "object storage problem 401",
			path:        "/1.3/object-storage-2",
			status:      http.StatusUnauthorized,
			contentType: "application/problem+json",
			body:        `{"correlation_id":"fixture-correlation","status":401,"title":"fixture unauthorized","type":"https://example.test/errors/unauthorized"}`,
			call: func(ctx context.Context, client *ClientWithResponses) (generatedErrorResponse, error) {
				return client.ListObjectStoragesWithResponse(ctx, nil)
			},
			assertTyped: func(t *testing.T, response generatedErrorResponse) {
				objectStorageResponse, ok := response.(*ListObjectStoragesResp)
				require.True(t, ok)
				require.NotNil(t, objectStorageResponse.ApplicationproblemJSONDefault)
				require.Equal(t, "fixture unauthorized", objectStorageResponse.ApplicationproblemJSONDefault.Title)
				require.Equal(t, int32(http.StatusUnauthorized), objectStorageResponse.ApplicationproblemJSONDefault.Status)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			requests := make(chan *http.Request, 1)
			server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
				requests <- request
				writer.Header().Set("Content-Type", testCase.contentType)
				writer.WriteHeader(testCase.status)
				_, _ = fmt.Fprint(writer, testCase.body)
			}))
			defer server.Close()

			client, err := NewClientWithResponses(server.URL)
			require.NoError(t, err)
			response, err := testCase.call(t.Context(), client)
			require.NoError(t, err)
			require.NotNil(t, response)
			observedRequest := <-requests
			require.Equal(t, http.MethodGet, observedRequest.Method)
			require.Equal(t, testCase.path, observedRequest.URL.Path)
			require.Equal(t, testCase.status, response.StatusCode())
			require.Equal(t, testCase.contentType, response.ContentType())
			require.JSONEq(t, testCase.body, string(response.GetBody()))
			require.NotNil(t, testCase.assertTyped, "error test case must assert its generated typed response")
			testCase.assertTyped(t, response)
		})
	}
}

func TestGetRouterWithResponseDecodesTypedNotFound(t *testing.T) {
	const responseBody = `{"error":{"error_code":"RESOURCE_NOT_FOUND","error_message":"router not found"}}`
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusNotFound)
		_, _ = fmt.Fprint(writer, responseBody)
	}))
	defer server.Close()

	client, err := NewClientWithResponses(server.URL)
	require.NoError(t, err)
	response, err := client.GetRouterWithResponse(t.Context(), uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"))
	require.NoError(t, err)
	require.Equal(t, http.StatusNotFound, response.StatusCode())
	require.NotNil(t, response.JSON404)
	require.Equal(t, "RESOURCE_NOT_FOUND", response.JSON404.Error.ErrorCode)
	require.Equal(t, "router not found", response.JSON404.Error.ErrorMessage)
	require.JSONEq(t, responseBody, string(response.GetBody()))
}

func TestGetRouterWithResponseReturnsParseErrorForMalformedSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(writer, `{"router":`)
	}))
	defer server.Close()

	client, err := NewClientWithResponses(server.URL)
	require.NoError(t, err)
	response, err := client.GetRouterWithResponse(t.Context(), uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"))
	require.Nil(t, response)
	require.Error(t, err)
}
