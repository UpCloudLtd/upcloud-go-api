package client

import (
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestClientBaseURL(t *testing.T) {
	assert.Equal(t, DefaultAPIBaseURL, clientBaseURL(""))
	assert.Equal(t, DefaultAPIBaseURL, clientBaseURL("127.0.0.1"))
	assert.Equal(t, DefaultAPIBaseURL, clientBaseURL("http://"))
	assert.Equal(t, "http://127.0.0.1", clientBaseURL("http://127.0.0.1"))
	assert.Equal(t, "https://127.0.0.1", clientBaseURL("https://127.0.0.1"))
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
