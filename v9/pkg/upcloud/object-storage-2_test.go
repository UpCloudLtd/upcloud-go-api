package upcloud

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"testing/synctest"
	"time"

	"github.com/stretchr/testify/require"
)

const (
	deletionWait = "deletion"
	bucketWait   = "bucket"
)

type waitResponseDoer struct {
	status   int
	body     string
	calls    int
	sequence []waitHTTPResponse
	err      error
}

type waitHTTPResponse struct {
	status int
	body   string
}

func (doer *waitResponseDoer) Do(request *http.Request) (*http.Response, error) {
	doer.calls++
	if doer.err != nil {
		return nil, doer.err
	}
	status, body := doer.status, doer.body
	if len(doer.sequence) > 0 {
		response := doer.sequence[min(doer.calls-1, len(doer.sequence)-1)]
		status, body = response.status, response.body
	}
	contentType := "application/json"
	if status >= http.StatusBadRequest {
		contentType = "application/problem+json"
	}
	return &http.Response{
		StatusCode: status,
		Status:     fmt.Sprintf("%d %s", status, http.StatusText(status)),
		Header:     http.Header{"Content-Type": {contentType}},
		Body:       io.NopCloser(strings.NewReader(body)),
		Request:    request,
	}, nil
}

func TestObjectStorageWaitOutcomes(t *testing.T) {
	transportError := errors.New("transport failed")
	type waitCase struct {
		name           string
		wait           string
		responses      []waitHTTPResponse
		transportError error
		wantError      string
		wantCalls      int
	}
	tests := make([]waitCase, 0, 29)
	tests = append(tests, []waitCase{
		{name: "state ready", wait: "state", responses: []waitHTTPResponse{{200, `{"operational_state":"started"}`}}, wantCalls: 1},
		{name: "state transition", wait: "state", responses: []waitHTTPResponse{{200, `{"operational_state":"starting"}`}, {200, `{"operational_state":"started"}`}}, wantCalls: 2},
		{name: "state not found", wait: "state", responses: []waitHTTPResponse{{404, `{}`}}, wantError: "404", wantCalls: 1},
		{name: "state not found malformed error body", wait: "state", responses: []waitHTTPResponse{{404, `{`}}, wantError: "404", wantCalls: 1},
		{name: "state missing field", wait: "state", responses: []waitHTTPResponse{{200, `{}`}}, wantError: "missing operational state", wantCalls: 1},
		{name: "service absent", wait: "deletion", responses: []waitHTTPResponse{{404, `{}`}}, wantCalls: 1},
		{name: "service absent malformed error body", wait: "deletion", responses: []waitHTTPResponse{{404, `{`}}, wantCalls: 1},
		{name: "service deletion transition", wait: "deletion", responses: []waitHTTPResponse{{200, `{}`}, {404, `{}`}}, wantCalls: 2},
		{name: "bucket service absent", wait: "bucket", responses: []waitHTTPResponse{{404, `{}`}}, wantCalls: 1},
		{name: "bucket absent malformed error body", wait: "bucket", responses: []waitHTTPResponse{{404, `{`}}, wantCalls: 1},
		{name: "bucket absent", wait: "bucket", responses: []waitHTTPResponse{{200, `[]`}}, wantCalls: 1},
		{name: "bucket null entry", wait: "bucket", responses: []waitHTTPResponse{{200, `[null]`}}, wantError: "missing bucket name", wantCalls: 1},
		{name: "bucket missing name", wait: "bucket", responses: []waitHTTPResponse{{200, `[{}]`}}, wantError: "missing bucket name", wantCalls: 1},
		{name: "bucket null name", wait: "bucket", responses: []waitHTTPResponse{{200, `[{"name":null}]`}}, wantError: "missing bucket name", wantCalls: 1},
		{name: "bucket empty name", wait: "bucket", responses: []waitHTTPResponse{{200, `[{"name":""}]`}}, wantError: "missing bucket name", wantCalls: 1},
		{name: "other bucket preserved", wait: "bucket", responses: []waitHTTPResponse{{200, `[{"name":"other"}]`}}, wantCalls: 1},
		{name: "bucket deletion transition", wait: "bucket", responses: []waitHTTPResponse{{200, `[{"name":"test-bucket"}]`}, {200, `[]`}}, wantCalls: 2},
	}...)
	for _, wait := range []string{"state", deletionWait, bucketWait} {
		tests = append(tests,
			waitCase{name: wait + " null payload", wait: wait, responses: []waitHTTPResponse{{200, `null`}}, wantError: "missing", wantCalls: 1},
			waitCase{name: wait + " malformed payload", wait: wait, responses: []waitHTTPResponse{{200, `{`}}, wantError: "unexpected end", wantCalls: 1},
			waitCase{name: wait + " transport failure", wait: wait, transportError: transportError, wantError: transportError.Error(), wantCalls: 1},
		)
		finalBody := `{"operational_state":"started"}`
		finalStatus := 200
		if wait == deletionWait {
			finalStatus, finalBody = 404, `{}`
		}
		if wait == bucketWait {
			finalBody = `[]`
		}
		tests = append(tests,
			waitCase{name: wait + " retry recovery", wait: wait, responses: []waitHTTPResponse{{500, `{`}, {502, `not-json`}, {503, ``}, {finalStatus, finalBody}}, wantCalls: 4},
			waitCase{name: wait + " retry exhaustion", wait: wait, responses: []waitHTTPResponse{{503, `{`}}, wantError: "503", wantCalls: 4},
		)
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				doer := &waitResponseDoer{sequence: test.responses, err: test.transportError}
				client, err := NewClientWithResponses("http://localhost", WithHTTPClient(doer))
				require.NoError(t, err)
				ctx, cancel := context.WithTimeout(t.Context(), time.Minute)
				defer cancel()
				const serviceUUID = "00000000-0000-4000-8000-000000000001"
				switch test.wait {
				case "state":
					var result *GetObjectStorage200
					result, err = client.WaitForObjectStorageOperationalState(ctx, serviceUUID, "started")
					if test.wantError == "" {
						require.NotNil(t, result)
						require.NotNil(t, result.OperationalState)
						require.Equal(t, "started", string(*result.OperationalState))
					}
				case deletionWait:
					err = client.WaitForObjectStorageDeletion(ctx, serviceUUID)
				case bucketWait:
					err = client.WaitForObjectStorageBucketDeletion(ctx, serviceUUID, "test-bucket")
				}
				if test.wantError == "" {
					require.NoError(t, err)
				} else {
					require.ErrorContains(t, err, test.wantError)
				}
				require.Equal(t, test.wantCalls, doer.calls)
			})
		})
	}
}

func TestObjectStorageWaitCancellation(t *testing.T) {
	testCases := []struct {
		name string
		body string
		wait func(context.Context, *ClientWithResponses) error
	}{
		{
			name: "operational state",
			body: `{"operational_state":"starting"}`,
			wait: func(ctx context.Context, client *ClientWithResponses) error {
				_, err := client.WaitForObjectStorageOperationalState(ctx, "00000000-0000-4000-8000-000000000001", "started")
				return err
			},
		},
		{
			name: "service deletion",
			body: `{}`,
			wait: func(ctx context.Context, client *ClientWithResponses) error {
				return client.WaitForObjectStorageDeletion(ctx, "00000000-0000-4000-8000-000000000001")
			},
		},
		{
			name: "bucket deletion",
			body: `[{"name":"test-bucket"}]`,
			wait: func(ctx context.Context, client *ClientWithResponses) error {
				return client.WaitForObjectStorageBucketDeletion(ctx, "00000000-0000-4000-8000-000000000001", "test-bucket")
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				doer := &waitResponseDoer{status: http.StatusOK, body: testCase.body}
				client, err := NewClientWithResponses("http://localhost", WithHTTPClient(doer))
				require.NoError(t, err)

				ctx, cancel := context.WithCancel(t.Context())
				cancel()
				require.ErrorIs(t, testCase.wait(ctx, client), context.Canceled)
				require.Zero(t, doer.calls)

				ctx, cancel = context.WithTimeout(t.Context(), 6*time.Second)
				defer cancel()
				require.ErrorIs(t, testCase.wait(ctx, client), context.DeadlineExceeded)
				require.Equal(t, 1, doer.calls)
			})
		})
	}
}

func TestObjectStorageWaitHTTPFailures(t *testing.T) {
	for _, status := range []int{400, 401, 403, 409, 429, 204, 302} {
		for _, wait := range []string{"state", deletionWait, bucketWait} {
			t.Run(fmt.Sprintf("%s/%d", wait, status), func(t *testing.T) {
				synctest.Test(t, func(t *testing.T) {
					doer := &waitResponseDoer{status: status, body: `{}`}
					client, err := NewClientWithResponses("http://localhost", WithHTTPClient(doer))
					require.NoError(t, err)
					ctx, cancel := context.WithTimeout(t.Context(), 11*time.Second)
					defer cancel()
					const serviceUUID = "00000000-0000-4000-8000-000000000001"
					switch wait {
					case "state":
						_, err = client.WaitForObjectStorageOperationalState(ctx, serviceUUID, "started")
					case deletionWait:
						err = client.WaitForObjectStorageDeletion(ctx, serviceUUID)
					case bucketWait:
						err = client.WaitForObjectStorageBucketDeletion(ctx, serviceUUID, "test-bucket")
					}
					require.Error(t, err)
					require.False(t, errors.Is(err, context.DeadlineExceeded), "HTTP failure must terminate, not poll until timeout")
					require.Contains(t, err.Error(), strconv.Itoa(status))
					require.Equal(t, 1, doer.calls)
				})
			})
		}
	}
}

func TestObjectStorageWaitHelpersRejectInvalidUUIDWithoutRequest(t *testing.T) {
	testCases := []struct {
		name string
		wait func(*ClientWithResponses) error
	}{
		{
			name: "operational state",
			wait: func(client *ClientWithResponses) error {
				_, err := client.WaitForObjectStorageOperationalState(t.Context(), "invalid-uuid", "started")
				return err
			},
		},
		{
			name: "service deletion",
			wait: func(client *ClientWithResponses) error {
				return client.WaitForObjectStorageDeletion(t.Context(), "invalid-uuid")
			},
		},
		{
			name: "bucket deletion",
			wait: func(client *ClientWithResponses) error {
				return client.WaitForObjectStorageBucketDeletion(t.Context(), "invalid-uuid", "fixture-bucket")
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			doer := &waitResponseDoer{status: http.StatusOK, body: `{}`}
			client, err := NewClientWithResponses("http://localhost", WithHTTPClient(doer))
			require.NoError(t, err)
			require.Error(t, testCase.wait(client))
			require.Zero(t, doer.calls, "invalid UUID must fail before making a request")
		})
	}
}
