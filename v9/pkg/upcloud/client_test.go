package upcloud

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
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

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Date", "Fri, 11 Oct 2024 23:58:00 GMT")
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, `{"error":{"error_code":"AUTHENTICATION_FAILED","error_message":"authentication failed using the API token"}}`)
	}))
	defer srv.Close()

	c, err := New("ucat_testtoken", WithBaseURL(srv.URL), WithLogger(logger.DebugContext))
	require.NoError(t, err)

	c.CreateRouter(context.TODO(), CreateRouterJSONRequestBody{
		Name: "test",
	})

	expected := fmt.Sprintf(`{"time":"2 Minutes to Midnight","level":"DEBUG","msg":"Sending request to UpCloud API","url":"http://server/1.3/router","method":"POST","headers":{"Authorization":["Bearer [REDACTED]"],"Content-Type":["application/json"],"User-Agent":["upcloud-go-api/v9 openapi/%s"]},"body":"{\n  \"name\": \"test\"\n}"}
{"time":"2 Minutes to Midnight","level":"DEBUG","msg":"Received response from UpCloud API","url":"http://server/1.3/router","status":"401 Unauthorized","headers":{"Content-Length":["108"],"Content-Type":["text/plain; charset=utf-8"],"Date":["Fri, 11 Oct 2024 23:58:00 GMT"]},"body":"{\n  \"error\": {\n    \"error_code\": \"AUTHENTICATION_FAILED\",\n    \"error_message\": \"authentication failed using the API token\"\n  }\n}"}
`, specVersion)
	assert.Equal(t, expected, output.String())
}
