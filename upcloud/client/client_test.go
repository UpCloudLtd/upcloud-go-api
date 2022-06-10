package client

import (
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	globals "github.com/UpCloudLtd/upcloud-go-api/v4/internal"
	"github.com/stretchr/testify/assert"
)

func TestClientBaseURL(t *testing.T) {
	assert.Equal(t, DefaultAPIBaseURL, clientBaseURL(""))
	assert.Equal(t, DefaultAPIBaseURL, clientBaseURL("127.0.0.1"))
	assert.Equal(t, DefaultAPIBaseURL, clientBaseURL("http://"))
	assert.Equal(t, "http://127.0.0.1", clientBaseURL("http://127.0.0.1"))
	assert.Equal(t, "https://127.0.0.1", clientBaseURL("https://127.0.0.1"))
}

func TestClientTimeout(t *testing.T) {
	var u, p string
	c1 := New(u, p)
	assert.Equal(t, time.Second*DefaultTimeout, c1.GetTimeout())
	c2 := NewWithHTTPClient(u, p, http.DefaultClient)
	assert.Equal(t, time.Second*DefaultTimeout, c2.GetTimeout())
}

func TestClientUserAgent(t *testing.T) {
	var u, p string
	c1 := New(u, p)
	assert.Equal(t, fmt.Sprintf("upcloud-go-api/%s", globals.Version), c1.UserAgent)
}

func ExampleNew() {
	client := New(os.Getenv("UPCLOUD_USERNAME"), os.Getenv("UPCLOUD_PASSWORD"))
	client.SetTimeout(10 * time.Second)
}

func ExampleNewWithHTTPClient() {
	client := NewWithHTTPClient(os.Getenv("UPCLOUD_USERNAME"), os.Getenv("UPCLOUD_PASSWORD"), &http.Client{
		// setup custom HTTP client
	})
	client.SetTimeout(10 * time.Second)
}
