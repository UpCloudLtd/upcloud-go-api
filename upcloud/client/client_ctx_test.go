package client

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	globals "github.com/UpCloudLtd/upcloud-go-api/v4/internal"
	"github.com/stretchr/testify/assert"
)

func TestClientContextTimeout(t *testing.T) {
	var u, p string
	c1 := NewWithContext(u, p)
	assert.Equal(t, time.Duration(0), c1.GetTimeout())
	c2 := NewWithHTTPClientContext(u, p, http.DefaultClient)
	assert.Equal(t, time.Duration(0), c2.GetTimeout())
}

func TestClientContextUserAgent(t *testing.T) {
	var u, p string
	c1 := NewWithContext(u, p)
	assert.Equal(t, fmt.Sprintf("upcloud-go-api/%s", globals.Version), c1.UserAgent)
}
