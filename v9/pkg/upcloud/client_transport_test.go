package upcloud

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/zalando/go-keyring"
)

type capturedRequestDoer struct {
	request *http.Request
	calls   int
}

func (doer *capturedRequestDoer) Do(request *http.Request) (*http.Response, error) {
	doer.request = request
	doer.calls++
	return &http.Response{
		StatusCode: http.StatusOK,
		Status:     "200 OK",
		Header:     make(http.Header),
		Body:       io.NopCloser(strings.NewReader("")),
		Request:    request,
	}, nil
}

type errorRequestDoer struct {
	err   error
	calls int
}

func (doer *errorRequestDoer) Do(*http.Request) (*http.Response, error) {
	doer.calls++
	return nil, doer.err
}

type contextBlockingDoer struct {
	started chan struct{}
}

type recordingRoundTripper struct {
	base  http.RoundTripper
	calls int
}

func (roundTripper *recordingRoundTripper) RoundTrip(request *http.Request) (*http.Response, error) {
	roundTripper.calls++
	return roundTripper.base.RoundTrip(request)
}

type failingBodyReader struct {
	err error
}

func (reader failingBodyReader) Read([]byte) (int, error) {
	return 0, reader.err
}

type responseBodyErrorDoer struct {
	failingBodyReader
}

type trackedResponseBody struct {
	reader     io.Reader
	readErr    error
	closeErr   error
	closeCalls int
}

type trackedResponseDoer struct {
	body  *trackedResponseBody
	calls int
}

func (doer responseBodyErrorDoer) Do(request *http.Request) (*http.Response, error) {
	return &http.Response{
		StatusCode: http.StatusOK,
		Status:     "200 OK",
		Header:     make(http.Header),
		Body:       io.NopCloser(doer.failingBodyReader),
		Request:    request,
	}, nil
}

func (body *trackedResponseBody) Read(buffer []byte) (int, error) {
	if body.readErr != nil {
		return 0, body.readErr
	}
	return body.reader.Read(buffer)
}

func (body *trackedResponseBody) Close() error {
	body.closeCalls++
	return body.closeErr
}

func (doer *trackedResponseDoer) Do(request *http.Request) (*http.Response, error) {
	doer.calls++
	return &http.Response{
		StatusCode: http.StatusOK,
		Status:     "200 OK",
		Header:     http.Header{"Content-Type": []string{"text/plain"}},
		Body:       doer.body,
		Request:    request,
	}, nil
}

func (doer *contextBlockingDoer) Do(request *http.Request) (*http.Response, error) {
	close(doer.started)
	<-request.Context().Done()
	return nil, request.Context().Err()
}

func TestClientRequestEditorOrderAndBaseURL(t *testing.T) {
	doer := &capturedRequestDoer{}
	var editorOrder []string
	requestContextKey := struct{ name string }{name: "editor-context"}
	requestContext := context.WithValue(t.Context(), requestContextKey, "request-value")
	var editorContextValues []any
	client, err := NewClientWithResponses(
		"https://api.example.test/base",
		WithHTTPClient(doer),
		WithRequestEditorFn(func(ctx context.Context, request *http.Request) error {
			editorOrder = append(editorOrder, "client-first")
			editorContextValues = append(editorContextValues, ctx.Value(requestContextKey))
			request.Header.Set("X-Editor-Order", "first")
			return nil
		}),
		WithRequestEditorFn(func(ctx context.Context, request *http.Request) error {
			editorOrder = append(editorOrder, "client-second")
			editorContextValues = append(editorContextValues, ctx.Value(requestContextKey))
			request.Header.Set("X-Editor-Order", request.Header.Get("X-Editor-Order")+",second")
			return nil
		}),
	)
	require.NoError(t, err)

	response, err := client.ListRouters(requestContext, nil, func(ctx context.Context, request *http.Request) error {
		editorOrder = append(editorOrder, "per-call")
		editorContextValues = append(editorContextValues, ctx.Value(requestContextKey))
		request.Header.Set("X-Editor-Order", request.Header.Get("X-Editor-Order")+",per-call")
		return nil
	})
	require.NoError(t, err)
	require.NoError(t, response.Body.Close())

	require.Equal(t, 1, doer.calls)
	require.Equal(t, []string{"client-first", "client-second", "per-call"}, editorOrder)
	require.Equal(t, []any{"request-value", "request-value", "request-value"}, editorContextValues)
	require.Equal(t, "https://api.example.test/base/1.3/router", doer.request.URL.String())
	require.Equal(t, "first,second,per-call", doer.request.Header.Get("X-Editor-Order"))
	t.Logf("transport calls: %d method=%s url=%s", doer.calls, doer.request.Method, doer.request.URL)
}

func TestClientWithLoggerRejectsCustomDoer(t *testing.T) {
	for _, test := range []struct {
		name    string
		options []ClientOption
	}{
		{
			name: "logger after custom doer",
			options: []ClientOption{
				WithHTTPClient(&capturedRequestDoer{}),
				WithLogger(slog.Default().DebugContext),
			},
		},
		{
			name: "custom doer after logger",
			options: []ClientOption{
				WithLogger(slog.Default().DebugContext),
				WithHTTPClient(&capturedRequestDoer{}),
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			client, err := NewClientWithResponses("https://api.example.test", test.options...)
			require.Nil(t, client)
			require.ErrorContains(t, err, "WithLogger requires an *http.Client")
		})
	}
}

func TestClientRequestEditorErrorPreventsTransport(t *testing.T) {
	editorErr := errors.New("editor failed")
	doer := &capturedRequestDoer{}
	var editorOrder []string
	client, err := NewClientWithResponses(
		"https://api.example.test",
		WithHTTPClient(doer),
		WithRequestEditorFn(func(context.Context, *http.Request) error {
			editorOrder = append(editorOrder, "before-error")
			return nil
		}),
		WithRequestEditorFn(func(context.Context, *http.Request) error {
			editorOrder = append(editorOrder, "error")
			return editorErr
		}),
		WithRequestEditorFn(func(context.Context, *http.Request) error {
			editorOrder = append(editorOrder, "after-error")
			return nil
		}),
	)
	require.NoError(t, err)

	_, err = client.ListRouters(t.Context(), nil, func(context.Context, *http.Request) error {
		editorOrder = append(editorOrder, "per-call")
		return nil
	})
	require.ErrorIs(t, err, editorErr)
	require.Equal(t, []string{"before-error", "error"}, editorOrder)
	require.Zero(t, doer.calls)
	t.Logf("transport calls after editor failure: %d", doer.calls)
}

func TestClientPropagatesTransportError(t *testing.T) {
	transportErr := errors.New("transport failed")
	doer := &errorRequestDoer{err: transportErr}
	client, err := NewClientWithResponses("https://api.example.test", WithHTTPClient(doer))
	require.NoError(t, err)

	_, err = client.ListRouters(t.Context(), nil)
	require.ErrorIs(t, err, transportErr)
	require.Equal(t, 1, doer.calls)
}

func TestClientRejectsInvalidOptions(t *testing.T) {
	optionErr := errors.New("invalid option")
	client, err := NewClientWithResponses("https://api.example.test", func(*Client) error {
		return optionErr
	})
	require.Nil(t, client)
	require.ErrorIs(t, err, optionErr)

	client, err = NewClientWithResponses("https://api.example.test", WithBaseURL("http://%zz"))
	require.Nil(t, client)
	require.Error(t, err)
	var urlErr *url.Error
	require.ErrorAs(t, err, &urlErr)
}

func TestClientPreservesCustomTransportLoggerAndTimeout(t *testing.T) {
	requestStarted := make(chan struct{})
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		close(requestStarted)
		<-request.Context().Done()
	}))
	defer server.Close()

	customTransport := &recordingRoundTripper{base: http.DefaultTransport}
	customClient := &http.Client{Transport: customTransport, Timeout: 100 * time.Millisecond}
	var logs bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&logs, &slog.HandlerOptions{Level: slog.LevelDebug}))
	client, err := NewClientWithResponses(
		server.URL,
		WithHTTPClient(customClient),
		WithLogger(logger.DebugContext),
	)
	require.NoError(t, err)

	_, err = client.ListRouters(t.Context(), nil)
	require.Error(t, err)
	require.ErrorIs(t, err, context.DeadlineExceeded)
	select {
	case <-requestStarted:
	case <-time.After(time.Second):
		t.Fatal("custom transport did not reach the local server")
	}
	require.Equal(t, 1, customTransport.calls)
	require.Same(t, customClient, client.Client.Client)
	require.Equal(t, 100*time.Millisecond, customClient.Timeout)
	require.Contains(t, logs.String(), "Sending request")
}

func TestClientLoggerBeforeHTTPClientPreservesCustomTransport(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
		writer.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	customTransport := &recordingRoundTripper{base: http.DefaultTransport}
	customClient := &http.Client{Transport: customTransport, Timeout: time.Second}
	var logs bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&logs, &slog.HandlerOptions{Level: slog.LevelDebug}))
	client, err := NewClientWithResponses(
		server.URL,
		WithLogger(logger.DebugContext),
		WithHTTPClient(customClient),
	)
	require.NoError(t, err)

	response, err := client.ListRouters(t.Context(), nil)
	require.NoError(t, err)
	require.NoError(t, response.Body.Close())
	require.Same(t, customClient, client.Client.Client)
	require.Equal(t, 1, customTransport.calls)
	require.Contains(t, logs.String(), "Sending request")
}

func TestClientRejectsPreCancelledContextBeforeNetwork(t *testing.T) {
	var calls int
	server := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		calls++
	}))
	defer server.Close()
	client, err := NewClientWithResponses(server.URL)
	require.NoError(t, err)
	ctx, cancel := context.WithCancel(t.Context())
	cancel()

	_, err = client.ListRouters(ctx, nil)
	require.ErrorIs(t, err, context.Canceled)
	require.Zero(t, calls)
}

func TestClientPropagatesRequestBodyReadFailure(t *testing.T) {
	bodyErr := errors.New("request body read failed")
	server := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
	defer server.Close()
	client, err := NewClientWithResponses(server.URL)
	require.NoError(t, err)

	_, err = client.CreateObjectStorageWithBody(t.Context(), "application/json", failingBodyReader{err: bodyErr})
	require.Error(t, err)
	require.ErrorIs(t, err, bodyErr)
}

func TestClientPropagatesResponseBodyReadFailure(t *testing.T) {
	bodyErr := errors.New("response body read failed")
	client, err := NewClientWithResponses(
		"https://api.example.test",
		WithHTTPClient(responseBodyErrorDoer{failingBodyReader{err: bodyErr}}),
	)
	require.NoError(t, err)

	_, err = client.ListRoutersWithResponse(t.Context(), nil)
	require.ErrorIs(t, err, bodyErr)
}

func TestClientDefaultAndTrailingSlashURLsAndHeaders(t *testing.T) {
	tests := []struct {
		name   string
		server string
		want   string
	}{
		{name: "default URL", server: ServerUrlHttpsapiUpcloudCom, want: "https://api.upcloud.com/1.3/router"},
		{name: "base URL without slash", server: "https://api.example.test/base", want: "https://api.example.test/base/1.3/router"},
		{name: "base URL with slash", server: "https://api.example.test/base/", want: "https://api.example.test/base/1.3/router"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			doer := &capturedRequestDoer{}
			client, err := NewClientWithResponses(test.server, WithHTTPClient(doer))
			require.NoError(t, err)
			response, err := client.ListRouters(t.Context(), nil, func(_ context.Context, request *http.Request) error {
				request.Header.Set("Accept", "application/json")
				return nil
			})
			require.NoError(t, err)
			require.NoError(t, response.Body.Close())
			require.Equal(t, test.want, doer.request.URL.String())
			require.Equal(t, "application/json", doer.request.Header.Get("Accept"))
			require.Empty(t, doer.request.Header.Get("Content-Type"))

			response, err = client.CreateObjectStorageWithBody(t.Context(), "application/vnd.upcloud.test+json", strings.NewReader(`{"name":"test"}`))
			require.NoError(t, err)
			require.NoError(t, response.Body.Close())
			require.Equal(t, "POST", doer.request.Method)
			require.Equal(t, "application/vnd.upcloud.test+json", doer.request.Header.Get("Content-Type"))
			require.Empty(t, doer.request.Header.Get("Accept"))
		})
	}
}

func TestNewFromEnvIncompleteCredentialsUseIsolatedKeyring(t *testing.T) {
	if os.Getenv("UPCLOUD_CREDENTIALS_KEYRING_HELPER") == "1" {
		keyring.MockInit()
		_, err := NewFromEnv()
		if err == nil || !strings.Contains(err.Error(), "credentials not found") {
			t.Fatalf("NewFromEnv with incomplete credentials: got %v", err)
		}
		return
	}

	for _, test := range []struct {
		name     string
		username string
		password string
	}{
		{name: "no environment credentials"},
		{name: "username only", username: "synthetic-user"},
		{name: "password only", password: "synthetic-password"},
	} {
		t.Run(test.name, func(t *testing.T) {
			executable, err := os.Executable()
			require.NoError(t, err)
			command := exec.CommandContext(t.Context(), executable, "-test.run=^TestNewFromEnvIncompleteCredentialsUseIsolatedKeyring$")
			command.Env = []string{
				"UPCLOUD_CREDENTIALS_KEYRING_HELPER=1",
				"UPCLOUD_USERNAME=" + test.username,
				"UPCLOUD_PASSWORD=" + test.password,
				"UPCLOUD_TOKEN=",
			}
			output, err := command.CombinedOutput()
			require.NoError(t, err, fmt.Sprintf("child test output: %s", output))
		})
	}
}

func TestNewFromEnvUsesPinnedKeyringCredentials(t *testing.T) {
	const helperEnv = "UPCLOUD_CREDENTIALS_KEYRING_AUTH_HELPER"
	if os.Getenv(helperEnv) != "1" {
		executable, err := os.Executable()
		require.NoError(t, err)
		for _, test := range []struct {
			name     string
			username string
			password string
			token    string
		}{
			{name: "token"},
			{name: "basic", username: "synthetic-keyring-user"},
			{
				name:     "environment token takes precedence over keyring",
				username: "synthetic-keyring-user",
				token:    strings.Join([]string{"synthetic-environment", "token"}, "-"),
			},
		} {
			t.Run(test.name, func(t *testing.T) {
				command := exec.CommandContext(t.Context(), executable, "-test.run=^TestNewFromEnvUsesPinnedKeyringCredentials$")
				command.Env = []string{
					helperEnv + "=1",
					"UPCLOUD_KEYRING_AUTH_KIND=" + test.name,
					"UPCLOUD_USERNAME=" + test.username,
					"UPCLOUD_PASSWORD=" + test.password,
					"UPCLOUD_TOKEN=" + test.token,
					EnvDebugAPIBaseURL + "=",
					EnvDebugSkipCertificateVerify + "=",
				}
				output, err := command.CombinedOutput()
				require.NoError(t, err, "child test output: %s", output)
			})
		}
		return
	}

	keyring.MockInit()
	wantAuthorization := "Bearer synthetic-keyring-token"
	switch os.Getenv("UPCLOUD_KEYRING_AUTH_KIND") {
	case "basic":
		require.NoError(t, keyring.Set("UpCloud", "synthetic-keyring-user", "synthetic-keyring-password"))
		wantAuthorization = "Basic " + base64.StdEncoding.EncodeToString([]byte("synthetic-keyring-user:synthetic-keyring-password"))
	case "environment token takes precedence over keyring":
		require.NoError(t, keyring.Set("UpCloud", "synthetic-keyring-user", "synthetic-keyring-password"))
		require.NoError(t, keyring.Set("UpCloud", "", "synthetic-keyring-token"))
		wantAuthorization = "Bearer " + strings.Join([]string{"synthetic-environment", "token"}, "-")
	default:
		require.NoError(t, keyring.Set("UpCloud", "", "synthetic-keyring-token"))
	}

	authorization := make(chan string, 1)
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		authorization <- request.Header.Get("Authorization")
		writer.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client, err := NewFromEnv(WithBaseURL(server.URL))
	require.NoError(t, err)
	response, err := client.ListRouters(t.Context(), nil)
	require.NoError(t, err)
	require.NoError(t, response.Body.Close())
	require.Equal(t, wantAuthorization, <-authorization)
}

func TestClientPropagatesCancellationDuringRequest(t *testing.T) {
	doer := &contextBlockingDoer{started: make(chan struct{})}
	client, err := NewClientWithResponses("https://api.example.test", WithHTTPClient(doer))
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	result := make(chan error, 1)
	go func() {
		_, requestErr := client.ListRouters(ctx, nil)
		result <- requestErr
	}()

	select {
	case <-doer.started:
	case <-time.After(time.Second):
		t.Fatal("request did not reach the doer")
	}
	cancel()

	select {
	case err := <-result:
		require.ErrorIs(t, err, context.Canceled)
	case <-time.After(time.Second):
		t.Fatal("request did not return after context cancellation")
	}
}

func TestClientCancellationDuringHeaderWait(t *testing.T) {
	requestReceived := make(chan struct{})
	server := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, request *http.Request) {
		close(requestReceived)
		<-request.Context().Done()
	}))
	defer server.Close()

	client, err := NewClientWithResponses(server.URL)
	require.NoError(t, err)
	ctx, cancel := context.WithCancel(t.Context())
	result := make(chan error, 1)
	go func() {
		_, requestErr := client.ListRouters(ctx, nil)
		result <- requestErr
	}()

	select {
	case <-requestReceived:
	case <-time.After(time.Second):
		t.Fatal("server did not receive request before header wait")
	}
	cancel()
	select {
	case err := <-result:
		require.ErrorIs(t, err, context.Canceled)
	case <-time.After(time.Second):
		t.Fatal("request did not return after cancellation during header wait")
	}
}

func TestClientCancellationDuringResponseBodyRead(t *testing.T) {
	bodyStarted := make(chan struct{})
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(writer, `{"routers":[`)
		writer.(http.Flusher).Flush()
		close(bodyStarted)
		<-request.Context().Done()
	}))
	defer server.Close()

	client, err := NewClientWithResponses(server.URL)
	require.NoError(t, err)
	ctx, cancel := context.WithCancel(t.Context())
	result := make(chan error, 1)
	go func() {
		_, requestErr := client.ListRoutersWithResponse(ctx, nil)
		result <- requestErr
	}()

	select {
	case <-bodyStarted:
	case <-time.After(time.Second):
		t.Fatal("server did not begin response body before body-read wait")
	}
	cancel()
	select {
	case err := <-result:
		require.ErrorIs(t, err, context.Canceled)
	case <-time.After(time.Second):
		t.Fatal("response body read did not return after cancellation")
	}
}

func TestClientResponseBodyReadAndCloseOutcomes(t *testing.T) {
	readErr := errors.New("response body read failed")
	closeErr := errors.New("response body close failed")
	for _, test := range []struct {
		name     string
		body     *trackedResponseBody
		wantBody string
		wantErr  error
	}{
		{
			name:     "close error is ignored after successful read",
			body:     &trackedResponseBody{reader: strings.NewReader("response"), closeErr: closeErr},
			wantBody: "response",
		},
		{
			name:    "read error is returned and body is closed",
			body:    &trackedResponseBody{readErr: readErr, closeErr: closeErr},
			wantErr: readErr,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			doer := &trackedResponseDoer{body: test.body}
			client, err := NewClientWithResponses("https://api.example.test", WithHTTPClient(doer))
			require.NoError(t, err)

			response, err := client.ListRoutersWithResponse(t.Context(), nil)
			if test.wantErr != nil {
				require.Nil(t, response)
				require.ErrorIs(t, err, test.wantErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, test.wantBody, string(response.Body))
			}
			require.Equal(t, 1, test.body.closeCalls)
			require.Equal(t, 1, doer.calls)
		})
	}
}

func TestNewFromEnvInsecureTLSIsOptIn(t *testing.T) {
	server := httptest.NewTLSServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client, err := New("synthetic-token", WithBaseURL(server.URL))
	require.NoError(t, err)
	_, err = client.ListRouters(t.Context(), nil)
	require.Error(t, err)
	require.ErrorContains(t, err, "certificate")

	t.Setenv("UPCLOUD_USERNAME", "")
	t.Setenv("UPCLOUD_PASSWORD", "")
	t.Setenv("UPCLOUD_TOKEN", "synthetic-token")
	t.Setenv(EnvDebugAPIBaseURL, server.URL)
	t.Setenv(EnvDebugSkipCertificateVerify, "1")

	debugClient, err := NewFromEnv()
	require.NoError(t, err)
	response, err := debugClient.ListRouters(t.Context(), nil)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, response.StatusCode)
	require.NoError(t, response.Body.Close())
}

func TestClientDoesNotForwardTokenAcrossHostRedirect(t *testing.T) {
	var authorization string
	redirectedServer := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		authorization = request.Header.Get("Authorization")
		writer.WriteHeader(http.StatusOK)
	}))
	defer redirectedServer.Close()
	redirectURL := strings.Replace(redirectedServer.URL, "127.0.0.1", "localhost", 1) + "/redirected"

	originServer := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		require.Equal(t, "Bearer synthetic-token", request.Header.Get("Authorization"))
		http.Redirect(writer, request, redirectURL, http.StatusFound)
	}))
	defer originServer.Close()

	client, err := New("synthetic-token", WithBaseURL(originServer.URL))
	require.NoError(t, err)
	response, err := client.ListRouters(t.Context(), nil)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, response.StatusCode)
	require.NoError(t, response.Body.Close())
	require.Empty(t, authorization)
}

func TestNewSetsTokenAndDefaultUserAgent(t *testing.T) {
	doer := &capturedRequestDoer{}
	client, err := New("synthetic-token", WithHTTPClient(doer))
	require.NoError(t, err)

	response, err := client.ListRouters(t.Context(), nil)
	require.NoError(t, err)
	require.NoError(t, response.Body.Close())
	require.Equal(t, "Bearer synthetic-token", doer.request.Header.Get("Authorization"))
	require.Equal(t, defaultUserAgent(), doer.request.Header.Get("User-Agent"))

	response, err = client.ListRouters(t.Context(), nil, func(_ context.Context, request *http.Request) error {
		request.Header.Set("User-Agent", "consumer-agent/1.0")
		return nil
	})
	require.NoError(t, err)
	require.NoError(t, response.Body.Close())
	require.Equal(t, "consumer-agent/1.0", doer.request.Header.Get("User-Agent"))
	require.Equal(t, 2, doer.calls)
}

func TestNewFromEnvAuthentication(t *testing.T) {
	tests := []struct {
		name              string
		username          string
		password          string
		token             string
		wantAuthorization string
	}{
		{
			name:              "token",
			token:             "synthetic-token",
			wantAuthorization: "Bearer synthetic-token",
		},
		{
			name:              "basic auth",
			username:          "synthetic-user",
			password:          "synthetic-password",
			wantAuthorization: "Basic " + base64.StdEncoding.EncodeToString([]byte("synthetic-user:synthetic-password")),
		},
		{
			name:              "token takes precedence over basic auth",
			username:          "synthetic-user",
			password:          "synthetic-password",
			token:             "synthetic-token",
			wantAuthorization: "Bearer synthetic-token",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Setenv("UPCLOUD_USERNAME", test.username)
			t.Setenv("UPCLOUD_PASSWORD", test.password)
			t.Setenv("UPCLOUD_TOKEN", test.token)
			t.Setenv(EnvDebugAPIBaseURL, "https://api.example.test/base")
			t.Setenv(EnvDebugSkipCertificateVerify, "")

			doer := &capturedRequestDoer{}
			client, err := NewFromEnv(WithHTTPClient(doer))
			require.NoError(t, err)

			response, err := client.ListRouters(t.Context(), nil)
			require.NoError(t, err)
			require.NoError(t, response.Body.Close())

			require.Equal(t, 1, doer.calls)
			require.Equal(t, "https://api.example.test/base/1.3/router", doer.request.URL.String())
			require.Equal(t, test.wantAuthorization, doer.request.Header.Get("Authorization"))
			require.Equal(t, defaultUserAgent(), doer.request.Header.Get("User-Agent"))
		})
	}
}
