package upcloud

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClientWithLogger(t *testing.T) {
	t.Parallel()

	var output strings.Builder
	logger := slog.New(slog.NewJSONHandler(&output, &slog.HandlerOptions{
		Level: slog.LevelDebug,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			// Replace time with a static value
			if a.Key == slog.TimeKey {
				a.Value = slog.StringValue("2 Minutes to Midnight")
			}

			// Replace URL with a static value as port is random
			if a.Key == "url" {
				re := regexp.MustCompile(`127\.0\.0\.1:\d+`)
				a.Value = slog.StringValue(re.ReplaceAllString(a.Value.String(), "server"))
			}
			return a
		},
	}))

	responseBody := `{"error":{"error_code":"AUTHENTICATION_FAILED","error_message":"authentication failed using the API token"}}`
	var requestCount atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount.Add(1)
		requestBody, readErr := io.ReadAll(r.Body)
		if readErr != nil {
			t.Errorf("read request body: %v", readErr)
		}
		if r.Method == http.MethodPost {
			assert.JSONEq(t, `{"router":{"name":"test"}}`, string(requestBody))
		} else {
			assert.Empty(t, requestBody)
		}
		w.Header().Set("Date", "Fri, 11 Oct 2024 23:58:00 GMT")
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, responseBody)
	}))
	defer srv.Close()

	c, err := New("ucat_testtoken", WithBaseURL(srv.URL), WithLogger(logger.DebugContext))
	require.NoError(t, err)

	listResponse, err := c.ListRouters(t.Context(), &ListRoutersParams{})
	require.NoError(t, err)
	require.Equal(t, http.StatusUnauthorized, listResponse.StatusCode)
	listBody, err := io.ReadAll(listResponse.Body)
	require.NoError(t, err)
	assert.Equal(t, responseBody, string(listBody))
	require.NoError(t, listResponse.Body.Close())
	createResponse, err := c.CreateRouter(t.Context(), CreateRouterJSONRequestBody{
		Router: struct {
			Labels       *[]RouterLabel `json:"labels,omitempty"`
			Name         string         `json:"name"`
			StaticRoutes *[]struct {
				Name    *string                                  `json:"name,omitempty"`
				Nexthop RouterCreate_Router_StaticRoutes_Nexthop `json:"nexthop"`
				Route   RouterIpCidr                             `json:"route"`
			} `json:"static_routes,omitempty"`
		}{Name: "test"},
	})
	require.NoError(t, err)
	require.Equal(t, http.StatusUnauthorized, createResponse.StatusCode)
	createBody, err := io.ReadAll(createResponse.Body)
	require.NoError(t, err)
	assert.Equal(t, responseBody, string(createBody))
	require.NoError(t, createResponse.Body.Close())

	expected := fmt.Sprintf(`{"time":"2 Minutes to Midnight","level":"DEBUG","msg":"Sending request to 127.0.0.1","url":"http://server/1.3/router","method":"GET","headers":{"Authorization":["Bearer [REDACTED]"],"User-Agent":["upcloud-go-api/v9 openapi/%s"]},"body":""}
{"time":"2 Minutes to Midnight","level":"DEBUG","msg":"Received response from 127.0.0.1","url":"http://server/1.3/router","status":"401 Unauthorized","headers":{"Content-Length":["108"],"Content-Type":["text/plain; charset=utf-8"],"Date":["Fri, 11 Oct 2024 23:58:00 GMT"]},"body":"{\n  \"error\": {\n    \"error_code\": \"AUTHENTICATION_FAILED\",\n    \"error_message\": \"authentication failed using the API token\"\n  }\n}"}
{"time":"2 Minutes to Midnight","level":"DEBUG","msg":"Sending request to 127.0.0.1","url":"http://server/1.3/router","method":"POST","headers":{"Authorization":["Bearer [REDACTED]"],"Content-Type":["application/json"],"User-Agent":["upcloud-go-api/v9 openapi/%s"]},"body":"{\n  \"router\": {\n    \"name\": \"test\"\n  }\n}"}
{"time":"2 Minutes to Midnight","level":"DEBUG","msg":"Received response from 127.0.0.1","url":"http://server/1.3/router","status":"401 Unauthorized","headers":{"Content-Length":["108"],"Content-Type":["text/plain; charset=utf-8"],"Date":["Fri, 11 Oct 2024 23:58:00 GMT"]},"body":"{\n  \"error\": {\n    \"error_code\": \"AUTHENTICATION_FAILED\",\n    \"error_message\": \"authentication failed using the API token\"\n  }\n}"}
`, specVersion, specVersion)
	assert.Equal(t, expected, output.String())
	require.EqualValues(t, 2, requestCount.Load())
	t.Logf("redacted logs: %s", output.String())
	t.Log("transport calls: GET=1 POST=1")
}

func TestClientWithLoggerRedactsBasicAuth(t *testing.T) {
	t.Setenv("UPCLOUD_USERNAME", "synthetic-basic-user")
	t.Setenv("UPCLOUD_PASSWORD", "synthetic-basic-password")
	t.Setenv("UPCLOUD_TOKEN", "")

	var authorization string
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		authorization = request.Header.Get("Authorization")
		writer.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	var output strings.Builder
	logger := slog.New(slog.NewJSONHandler(&output, &slog.HandlerOptions{Level: slog.LevelDebug}))
	client, err := NewFromEnv(WithBaseURL(server.URL), WithLogger(logger.DebugContext))
	require.NoError(t, err)
	response, err := client.ListRouters(t.Context(), nil)
	require.NoError(t, err)
	require.NoError(t, response.Body.Close())

	require.True(t, strings.HasPrefix(authorization, "Basic "))
	require.Contains(t, output.String(), "[REDACTED]")
	require.NotContains(t, output.String(), "synthetic-basic-user")
	require.NotContains(t, output.String(), "synthetic-basic-password")
	require.NotContains(t, output.String(), authorization)
	t.Logf("redacted logs: %s", output.String())
	t.Log("transport calls: 1")
}

func TestWithResponseDecodesCuratedFixtures(t *testing.T) {
	t.Parallel()

	const routerFixture = `{"routers":{"router":[{"attached_network_gateways":[],"attached_networks":{"network":[]},"labels":[{"key":"team","value":"platform"}],"name":"edge-router","static_routes":[],"type":"normal","uuid":"1200ecde-db95-4d1c-9133-6508f3232567"}]}}`
	const routerDetailFixture = `{"router":{"attached_network_gateways":[],"attached_networks":{"network":[]},"labels":[{"key":"team","value":"platform"}],"name":"edge-router-detail","static_routes":[],"type":"normal","uuid":"1200ecde-db95-4d1c-9133-6508f3232567"}}`
	const objectStorageFixture = `[{"configured_status":"started","created_at":"2026-01-02T03:04:05Z","labels":[{"key":"team","value":"archives"}],"name":"archive-store","operational_state":"running","region":"fi-hel1","termination_protection":false,"usage":{"total_objects":7,"total_size_bytes":4096},"uuid":"34a6f90a-9909-41d7-a03a-3d4a515e555e"},{"name":"empty-store","region":"fi-hel1","termination_protection":null,"usage":{"total_objects":0}}]`
	const zonesFixture = `{"zones":{"zone":[{"description":"Helsinki","id":"fi-hel1","public":"yes"}]}}`
	const serversFixture = `{"servers":{"server":[{"core_number":"2","created":1760000000,"hostname":"web-1","license":"0","memory_amount":"4096","state":"started","title":"web","uuid":"1200ecde-db95-4d1c-9133-6508f3232567","zone":"fi-hel1"}]}}`
	const databasesFixture = `[{"name":"analytics","type":"pg","uuid":"123e4567-e89b-12d3-a456-426614174000","zone":"fi-hel1"}]`
	const accountFixture = `{"account":{"resource_limits":{"cores":8},"trial_mode":1,"username":"fixture-account","extension":{"enabled":true}}}`

	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		switch request.URL.Path {
		case "/1.3/router":
			_, _ = io.WriteString(writer, routerFixture)
		case "/1.3/router/1200ecde-db95-4d1c-9133-6508f3232567":
			if request.Method == http.MethodDelete {
				writer.Header().Del("Content-Type")
				writer.WriteHeader(http.StatusNoContent)
				return
			}
			_, _ = io.WriteString(writer, routerDetailFixture)
		case "/1.3/object-storage-2":
			writer.Header().Set("UpCloud-Total-Count", "2")
			_, _ = io.WriteString(writer, objectStorageFixture)
		case "/1.3/zone":
			_, _ = io.WriteString(writer, zonesFixture)
		case "/1.3/server":
			_, _ = io.WriteString(writer, serversFixture)
		case "/1.3/database":
			_, _ = io.WriteString(writer, databasesFixture)
		case "/1.3/account":
			_, _ = io.WriteString(writer, accountFixture)
		default:
			t.Errorf("unexpected request path %q", request.URL.Path)
			http.NotFound(writer, request)
		}
	}))
	defer server.Close()

	client, err := NewClientWithResponses(server.URL)
	require.NoError(t, err)

	account, err := client.AccountDetailsWithResponse(t.Context())
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, account.StatusCode())
	require.NotNil(t, account.JSON200)
	require.Equal(t, "fixture-account", account.JSON200.Account.Username)
	require.NotNil(t, account.JSON200.Account.TrialMode)
	require.Equal(t, AccountBoolean01N1, *account.JSON200.Account.TrialMode)
	require.NotNil(t, account.JSON200.Account.ResourceLimits.Cores)
	require.Equal(t, int64(8), *account.JSON200.Account.ResourceLimits.Cores)
	extension, found := account.JSON200.Account.Get("extension")
	require.True(t, found)
	require.Equal(t, map[string]interface{}{"enabled": true}, extension)

	routers, err := client.ListRoutersWithResponse(t.Context(), nil)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, routers.StatusCode())
	require.NotNil(t, routers.JSON200)
	require.Len(t, routers.JSON200.Routers.Router, 1)
	router := routers.JSON200.Routers.Router[0]
	require.Equal(t, "edge-router", router.Name)
	require.Equal(t, "normal", string(router.Type))
	require.Equal(t, "1200ecde-db95-4d1c-9133-6508f3232567", router.Uuid.String())
	require.Equal(t, []RouterLabel{{Key: "team", Value: "platform"}}, router.Labels)

	routerDetails, err := client.GetRouterWithResponse(t.Context(), uuid.MustParse("1200ecde-db95-4d1c-9133-6508f3232567"))
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, routerDetails.StatusCode())
	require.NotNil(t, routerDetails.JSON200)
	require.Equal(t, "edge-router-detail", routerDetails.JSON200.Router.Name)
	require.Equal(t, "1200ecde-db95-4d1c-9133-6508f3232567", routerDetails.JSON200.Router.Uuid.String())
	deletedRouter, err := client.DeleteRouterWithResponse(t.Context(), uuid.MustParse("1200ecde-db95-4d1c-9133-6508f3232567"))
	require.NoError(t, err)
	require.Equal(t, http.StatusNoContent, deletedRouter.StatusCode())
	require.Empty(t, deletedRouter.GetBody())

	services, err := client.ListObjectStoragesWithResponse(t.Context(), nil)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, services.StatusCode())
	require.NotNil(t, services.JSON200)
	require.Len(t, *services.JSON200, 2)
	service := (*services.JSON200)[0]
	require.NotNil(t, service.Name)
	require.Equal(t, "archive-store", *service.Name)
	require.NotNil(t, service.CreatedAt)
	require.Equal(t, "2026-01-02T03:04:05Z", service.CreatedAt.Format("2006-01-02T15:04:05Z"))
	require.NotNil(t, service.TerminationProtection)
	require.False(t, *service.TerminationProtection)
	require.NotNil(t, service.Usage)
	require.NotNil(t, service.Usage.TotalObjects)
	require.Equal(t, 7, *service.Usage.TotalObjects)
	require.NotNil(t, service.Labels)
	require.Len(t, *service.Labels, 1)
	require.Equal(t, "archives", *(*service.Labels)[0].Value)
	require.NotNil(t, services.Headers200)
	require.NotNil(t, services.Headers200.UpCloudTotalCount)
	require.Equal(t, int32(2), *services.Headers200.UpCloudTotalCount)
	emptyService := (*services.JSON200)[1]
	require.Nil(t, emptyService.Labels)
	require.Nil(t, emptyService.TerminationProtection)
	require.NotNil(t, emptyService.Usage)
	require.NotNil(t, emptyService.Usage.TotalObjects)
	require.Zero(t, *emptyService.Usage.TotalObjects)

	zones, err := client.ListZonesWithResponse(t.Context())
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, zones.StatusCode())
	require.NotNil(t, zones.JSON200)
	require.Len(t, zones.JSON200.Zones.Zone, 1)
	zone := zones.JSON200.Zones.Zone[0]
	require.Equal(t, "fi-hel1", zone.Id)
	require.Equal(t, ZoneBooleanYesnoYes, zone.Public)

	servers, err := client.ListServersWithResponse(t.Context(), nil)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, servers.StatusCode())
	require.NotNil(t, servers.JSON200)
	require.NotNil(t, servers.JSON200.Servers.Server)
	require.Len(t, *servers.JSON200.Servers.Server, 1)
	serverDetails := (*servers.JSON200.Servers.Server)[0]
	require.Equal(t, "web-1", serverDetails.Hostname)
	require.Equal(t, "started", string(serverDetails.State))
	require.Equal(t, "1200ecde-db95-4d1c-9133-6508f3232567", serverDetails.Uuid.String())

	databases, err := client.ListDatabasesWithResponse(t.Context(), nil)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, databases.StatusCode())
	require.NotNil(t, databases.JSON200)
	require.Len(t, *databases.JSON200, 1)
	database := (*databases.JSON200)[0]
	require.NotNil(t, database.Name)
	require.Equal(t, "analytics", *database.Name)
	require.NotNil(t, database.Type)
	require.Equal(t, "pg", *database.Type)
	require.NotNil(t, database.Uuid)
	require.Equal(t, "123e4567-e89b-12d3-a456-426614174000", database.Uuid.String())
	require.NotNil(t, database.Zone)
	require.Equal(t, "fi-hel1", *database.Zone)
}

func TestAccountDetailsTrialModeAndAdditionalProperties(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string
		json string
		want AccountBoolean01
	}{
		{name: "disabled", json: `{"resource_limits":{},"trial_mode":0,"username":"fixture"}`, want: AccountBoolean01N0},
		{name: "enabled", json: `{"resource_limits":{},"trial_mode":1,"username":"fixture","extension":{"enabled":true}}`, want: AccountBoolean01N1},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			var account AccountDetails
			require.NoError(t, json.Unmarshal([]byte(testCase.json), &account))
			require.NotNil(t, account.TrialMode)
			require.Equal(t, testCase.want, *account.TrialMode)
			require.True(t, account.TrialMode.Valid())
			if testCase.name == "enabled" {
				extension, found := account.Get("extension")
				require.True(t, found)
				require.Equal(t, map[string]interface{}{"enabled": true}, extension)
			}
		})
	}
}

func TestRouterIPAddressUnionJSONRoundTrip(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name  string
		value string
		set   func(*RouterIpAddress, string) error
		get   func(RouterIpAddress) (string, error)
	}{
		{
			name:  "IPv4",
			value: "192.0.2.10",
			set: func(address *RouterIpAddress, value string) error {
				return address.FromRouterIpAddress0(value)
			},
			get: func(address RouterIpAddress) (string, error) {
				return address.AsRouterIpAddress0()
			},
		},
		{
			name:  "IPv6",
			value: "2001:db8::10",
			set: func(address *RouterIpAddress, value string) error {
				return address.FromRouterIpAddress1(value)
			},
			get: func(address RouterIpAddress) (string, error) {
				return address.AsRouterIpAddress1()
			},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			var address RouterIpAddress
			require.NoError(t, testCase.set(&address, testCase.value))
			encoded, err := json.Marshal(address)
			require.NoError(t, err)
			require.JSONEq(t, `"`+testCase.value+`"`, string(encoded))

			var decoded RouterIpAddress
			require.NoError(t, json.Unmarshal(encoded, &decoded))
			actual, err := testCase.get(decoded)
			require.NoError(t, err)
			require.Equal(t, testCase.value, actual)
		})
	}
}

func TestNewHonorsDebugEnv(t *testing.T) {
	t.Setenv(EnvDebugAPIBaseURL, "http://127.0.0.1:8080")
	t.Setenv(EnvDebugSkipCertificateVerify, "1")

	c, err := New("token")
	require.NoError(t, err)
	require.Equal(t, "http://127.0.0.1:8080/", c.Server)
	tr, ok := c.Client.Client.(*http.Client).Transport.(*http.Transport)
	require.True(t, ok)
	require.NotNil(t, tr.TLSClientConfig)
	require.True(t, tr.TLSClientConfig.InsecureSkipVerify)
}

func TestNewDefaultServer(t *testing.T) {
	t.Setenv(EnvDebugAPIBaseURL, "")
	t.Setenv(EnvDebugSkipCertificateVerify, "")

	c, err := New("token")
	require.NoError(t, err)
	require.Equal(t, ServerUrlHttpsapiUpcloudCom+"/", c.Server)
}
