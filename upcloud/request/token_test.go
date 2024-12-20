package request

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetTokenDetailsRequest(t *testing.T) {
	assert.Equal(t, "/account/tokens/foo", (&GetTokenDetailsRequest{ID: "foo"}).RequestURL())
}

func TestGetTokensRequest(t *testing.T) {
	assert.Equal(t, "/account/tokens", (&GetTokensRequest{}).RequestURL())
	assert.Equal(t, "/account/tokens", (&GetTokensRequest{}).RequestURL())
}

func TestDeleteTokenRequest(t *testing.T) {
	assert.Equal(t, "/account/tokens/foo", (&DeleteTokenRequest{ID: "foo"}).RequestURL())
}

func TestCreateTokenRequest(t *testing.T) {
	want := `
	{
		"name": "my_1st_token",
		"expires_at": "2025-01-01T00:00:00Z",
		"can_create_tokens": true,
		"allowed_ip_ranges": ["0.0.0.0/0", "::/0"]
	}
	`
	req := &CreateTokenRequest{
		Name:               "my_1st_token",
		ExpiresAt:          time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
		CanCreateSubTokens: true,
		AllowedIPRanges:    []string{"0.0.0.0/0", "::/0"},
	}
	got, err := json.Marshal(req)
	assert.NoError(t, err)
	assert.JSONEq(t, want, string(got))
	assert.Equal(t, "/account/tokens", req.RequestURL())
}
