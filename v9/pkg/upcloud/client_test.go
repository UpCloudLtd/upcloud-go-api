package upcloud

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"sync/atomic"
	"testing"

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
