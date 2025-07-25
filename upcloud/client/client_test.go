package client

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClientBaseURL(t *testing.T) {
	t.Parallel()

	assert.Equal(t, APIBaseURL, clientBaseURL(""))
	assert.Equal(t, APIBaseURL, clientBaseURL("127.0.0.1"))
	assert.Equal(t, APIBaseURL, clientBaseURL("http://"))
	assert.Equal(t, "http://127.0.0.1", clientBaseURL("http://127.0.0.1"))
	assert.Equal(t, "https://127.0.0.1", clientBaseURL("https://127.0.0.1"))
}

func TestClientTimeout(t *testing.T) {
	t.Parallel()

	var u, p string
	c := New(u, p)
	assert.Equal(t, time.Duration(0), c.config.httpClient.Timeout)

	c = New(u, p, WithTimeout(5*time.Second))
	assert.Equal(t, 5*time.Second, c.config.httpClient.Timeout)
}

func TestAddDefaultHeaders(t *testing.T) {
	t.Parallel()

	wantUsername := "user"
	wantPassword := "pass"
	c := New(wantUsername, wantPassword)
	r, err := c.createRequest(context.TODO(), http.MethodGet, "/", nil)
	require.NoError(t, err)
	c.addDefaultHeaders(r)
	assert.Equal(t, "application/json", r.Header.Get("Accept"))
	assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
	assert.Equal(t, c.UserAgent, r.Header.Get("User-Agent"))
	r.Header.Set("Accept", "text/plain")
	c.addDefaultHeaders(r)
	assert.Equal(t, "text/plain", r.Header.Get("Accept"))
	assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
	assert.Equal(t, c.UserAgent, r.Header.Get("User-Agent"))
	gotUsername, gotPassword, ok := r.BasicAuth()
	assert.True(t, ok)
	assert.Equal(t, wantUsername, gotUsername)
	assert.Equal(t, wantPassword, gotPassword)
}

func TestClientCreateRequest(t *testing.T) {
	t.Parallel()

	c := New("", "")
	wantBody := []byte("test content")
	r, err := c.createRequest(context.TODO(), http.MethodPost, "/foo/bar", wantBody)
	require.NoError(t, err)
	assert.Equal(t, "POST", r.Method)
	assert.Equal(t, int64(len(wantBody)), r.ContentLength)
	assert.Equal(t, c.createRequestURL("/foo/bar"), r.URL.String())
	gotBody, err := io.ReadAll(r.Body)
	require.NoError(t, err)
	assert.Equal(t, wantBody, gotBody)
}

func TestClientUserAgent(t *testing.T) {
	t.Parallel()

	var u, p string
	c1 := New(u, p)
	assert.Equal(t, fmt.Sprintf("upcloud-go-api/%s", Version), c1.UserAgent)
}

func TestClientGet(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, fmt.Sprintf("/%s%s", APIVersion, "/test"), r.URL.Path)
		fmt.Fprint(w, string("ok"))
	}))
	defer srv.Close()
	c := New("", "", WithBaseURL(srv.URL))
	res, err := c.Get(context.TODO(), "/test")
	require.NoError(t, err)
	assert.Equal(t, "ok", string(res))
}

func TestClientPut(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPut, r.Method)
		assert.Equal(t, fmt.Sprintf("/%s%s", APIVersion, "/test"), r.URL.Path)
		fmt.Fprint(w, string("ok"))
	}))
	defer srv.Close()
	c := New("", "", WithBaseURL(srv.URL))
	res, err := c.Put(context.TODO(), "/test", nil)
	require.NoError(t, err)
	assert.Equal(t, "ok", string(res))
}

func TestClientPatch(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPatch, r.Method)
		assert.Equal(t, fmt.Sprintf("/%s%s", APIVersion, "/test"), r.URL.Path)
		fmt.Fprint(w, string("ok"))
	}))
	defer srv.Close()
	c := New("", "", WithBaseURL(srv.URL))
	res, err := c.Patch(context.TODO(), "/test", nil)
	require.NoError(t, err)
	assert.Equal(t, "ok", string(res))
}

func TestClientDelete(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		assert.Equal(t, fmt.Sprintf("/%s%s", APIVersion, "/test"), r.URL.Path)
		fmt.Fprint(w, string("ok"))
	}))
	defer srv.Close()
	c := New("", "", WithBaseURL(srv.URL))
	res, err := c.Delete(context.TODO(), "/test")
	require.NoError(t, err)
	assert.Equal(t, "ok", string(res))
}

func TestClientDo(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, string("ok"))
		assert.Equal(t, "/test", r.URL.Path)
		// test that we don't leak credentials when calling something else than baseURL
		_, _, ok := r.BasicAuth()
		assert.False(t, ok)
	}))
	defer srv.Close()
	c := New("", "")
	req, err := http.NewRequestWithContext(
		context.Background(),
		http.MethodGet,
		fmt.Sprintf("%s/test", srv.URL),
		nil,
	)
	require.NoError(t, err)
	res, err := c.Do(req)
	require.NoError(t, err)
	assert.Equal(t, "ok", string(res))
}

func TestClientPost(t *testing.T) {
	t.Parallel()

	timeout, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	wantUsername := "user"
	wantPassword := "pass"
	wantBody := []byte("test body")
	wantPath := "/some/path"
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// check common headers
		assert.Equal(t, "application/json", r.Header.Get("Accept"))
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, userAgent(), r.Header.Get("User-Agent"))
		// check auth header
		gotUsername, gotPassword, ok := r.BasicAuth()
		assert.True(t, ok)
		assert.Equal(t, wantUsername, gotUsername)
		assert.Equal(t, wantPassword, gotPassword)
		// check body
		gotBody, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		assert.Equal(t, wantBody, gotBody)
		// check URL
		assert.Equal(t, fmt.Sprintf("/%s%s", APIVersion, wantPath), r.URL.Path)
		fmt.Fprint(w, string("ok"))
	}))
	defer srv.Close()
	c := New(wantUsername, wantPassword, WithBaseURL(srv.URL))
	res, err := c.Post(timeout, wantPath, wantBody)
	require.NoError(t, err)
	assert.Equal(t, "ok", string(res))
}

func TestClientGetContextDeadline(t *testing.T) {
	t.Parallel()

	deadline, cancel := context.WithDeadline(context.Background(), time.Now().Add(2*time.Second))
	defer cancel()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(3 * time.Second)
		fmt.Fprint(w, string("ok"))
	}))
	defer srv.Close()
	c := New("", "", WithBaseURL(srv.URL))
	_, err := c.Get(deadline, "/")
	require.True(t, errors.Is(err, context.DeadlineExceeded))
}

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
		fmt.Fprintf(w, `{"method": "%s", "path": "%s"}`, r.Method, r.URL.Path)
	}))
	defer srv.Close()

	c := New("username", "password", WithBaseURL(srv.URL), WithLogger(logger.DebugContext))
	_, err := c.Get(context.TODO(), "/test")
	require.NoError(t, err)

	expected := fmt.Sprintf(`{"time":"2 Minutes to Midnight","level":"DEBUG","msg":"Sending request to UpCloud API","url":"http://server/1.3/test","method":"GET","headers":{"Accept":["application/json"],"Authorization":["Basic [REDACTED]"],"Content-Type":["application/json"],"User-Agent":["upcloud-go-api/%s"]},"body":""}
{"time":"2 Minutes to Midnight","level":"DEBUG","msg":"Received response from UpCloud API","url":"http://server/1.3/test","status":"200 OK","headers":{"Content-Length":["38"],"Content-Type":["text/plain; charset=utf-8"],"Date":["Fri, 11 Oct 2024 23:58:00 GMT"]},"body":"{\n  \"method\": \"GET\",\n  \"path\": \"/1.3/test\"\n}"}
`, Version)
	assert.Equal(t, expected, output.String())
}

func ExampleWithTimeout() {
	New(os.Getenv("UPCLOUD_USERNAME"), os.Getenv("UPCLOUD_PASSWORD"), WithTimeout(10*time.Second))
}

func ExampleWithHTTPClient() {
	httpClient := &http.Client{
		// setup custom HTTP client
	}
	New(os.Getenv("UPCLOUD_USERNAME"), os.Getenv("UPCLOUD_PASSWORD"), WithHTTPClient(httpClient))
}

func TestNewFromEnv(t *testing.T) {
	t.Run("token authentication", func(t *testing.T) {
		t.Setenv(EnvToken, "test-token")
		t.Setenv(EnvUsername, "")
		t.Setenv(EnvPassword, "")

		client, err := NewFromEnv()
		require.NoError(t, err)
		assert.Equal(t, "test-token", client.config.token)
		assert.Empty(t, client.config.username)
		assert.Empty(t, client.config.password)
	})

	t.Run("basic auth", func(t *testing.T) {
		t.Setenv(EnvToken, "")
		t.Setenv(EnvUsername, "test-user")
		t.Setenv(EnvPassword, "test-pass")

		client, err := NewFromEnv()
		require.NoError(t, err)
		assert.Empty(t, client.config.token)
		assert.Equal(t, "test-user", client.config.username)
		assert.Equal(t, "test-pass", client.config.password)
	})

	t.Run("both token and basic auth provided - should error", func(t *testing.T) {
		t.Setenv(EnvToken, "test-token")
		t.Setenv(EnvUsername, "test-user")
		t.Setenv(EnvPassword, "test-pass")

		client, err := NewFromEnv()
		assert.Error(t, err)
		assert.Nil(t, client)
		assert.Contains(t, err.Error(), "only one authentication method")
	})

	t.Run("no credentials provided - should error", func(t *testing.T) {
		t.Setenv(EnvToken, "")
		t.Setenv(EnvUsername, "")
		t.Setenv(EnvPassword, "")

		client, err := NewFromEnv()
		assert.Error(t, err)
		assert.Nil(t, client)
		assert.Contains(t, err.Error(), "authentication credentials must be provided")
	})

	t.Run("only username provided - should error", func(t *testing.T) {
		t.Setenv(EnvToken, "")
		t.Setenv(EnvUsername, "test-user")
		t.Setenv(EnvPassword, "")

		client, err := NewFromEnv()
		assert.Error(t, err)
		assert.Nil(t, client)
		assert.Contains(t, err.Error(), "authentication credentials must be provided")
	})

	t.Run("only password provided - should error", func(t *testing.T) {
		t.Setenv(EnvToken, "")
		t.Setenv(EnvUsername, "")
		t.Setenv(EnvPassword, "test-pass")

		client, err := NewFromEnv()
		assert.Error(t, err)
		assert.Nil(t, client)
		assert.Contains(t, err.Error(), "authentication credentials must be provided")
	})

	t.Run("with config functions", func(t *testing.T) {
		t.Setenv(EnvToken, "test-token")
		t.Setenv(EnvUsername, "")
		t.Setenv(EnvPassword, "")

		client, err := NewFromEnv(WithTimeout(5 * time.Second))
		require.NoError(t, err)
		assert.Equal(t, "test-token", client.config.token)
		assert.Equal(t, 5*time.Second, client.config.httpClient.Timeout)
	})
}
