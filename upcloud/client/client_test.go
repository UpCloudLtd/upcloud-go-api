package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClientBaseURL(t *testing.T) {
	assert.Equal(t, DefaultAPIBaseURL, clientBaseURL(""))
	assert.Equal(t, DefaultAPIBaseURL, clientBaseURL("127.0.0.1"))
	assert.Equal(t, DefaultAPIBaseURL, clientBaseURL("http://"))
	assert.Equal(t, "http://127.0.0.1", clientBaseURL("http://127.0.0.1"))
	assert.Equal(t, "https://127.0.0.1", clientBaseURL("https://127.0.0.1"))
}
