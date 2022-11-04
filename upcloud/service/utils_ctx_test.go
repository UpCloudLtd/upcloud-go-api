package service

import (
	"context"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud/client"
	"github.com/dnaeon/go-vcr/cassette"
	"github.com/dnaeon/go-vcr/recorder"
	"github.com/hashicorp/go-cleanhttp"
	"github.com/stretchr/testify/require"
)

// records the API interactions of the test. Function provides both services to test cases so that old utility functions can be used to initialize environment.
func recordWithContext(t *testing.T, fixture string, f func(context.Context, *testing.T, *recorder.Recorder, *Service, *ServiceContext)) {
	if testing.Short() {
		t.Skip("Skipping recorded test in short mode")
	}

	r, err := recorder.New("fixtures/" + fixture)
	require.NoError(t, err)

	r.AddFilter(func(i *cassette.Interaction) error {
		delete(i.Request.Headers, "Authorization")
		if i.Request.Method == http.MethodPut && strings.Contains(i.Request.URL, "uploader") {
			// We will remove the body from the upload to reduce fixture size
			i.Request.Body = ""
		}
		return nil
	})

	defer func() {
		err := r.Stop()
		require.NoError(t, err)
	}()

	user, password := getCredentials()

	httpClient := cleanhttp.DefaultClient()
	origTransport := httpClient.Transport
	r.SetTransport(origTransport)
	httpClient.Transport = r

	c := client.NewWithHTTPClient(user, password, httpClient)
	c.SetTimeout(time.Second * 300)

	customAPI := os.Getenv("UPCLOUD_GO_SDK_API_HOST")
	if customAPI != "" {
		// Override api host after the go-vcr to maintain consistent test fixtures
		r.SetTransport(&customRoundTripper{fn: func(r *http.Request) (*http.Response, error) {
			clone := r.Clone(r.Context())
			clone.URL.Host = customAPI
			clone.Host = customAPI
			return origTransport.RoundTrip(clone)
		}})
	}

	// just some random timeout value. High enough that it won't be reached during normal test.
	ctx, cancel := context.WithTimeout(context.Background(), waitTimeout*4)
	defer cancel()
	f(ctx, t, r, New(c), NewWithContext(client.NewWithHTTPClientContext(user, password, httpClient)))
}
