package upcloud

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTokenUnmarshal(t *testing.T) {
	want := Tokens{
		{
			ID:                 "deadbeef-dead-beef-dead-beefdeadbee1",
			Name:               "token_workspace",
			Type:               TokenTypeWorkspace,
			CreatedAt:          timeParse("2024-12-19T11:46:09.888763Z"),
			ExpiresAt:          timeParse("2024-12-19T12:16:09.888531Z"),
			LastUsed:           nil,
			CanCreateSubTokens: true,
			AllowedIPRanges:    []string{"0.0.0.0/0", "::/0"},
		},
		{
			ID:                 "deadbeef-dead-beef-dead-beefdeadbee2",
			Name:               "token_pat",
			Type:               TokenTypePAT,
			CreatedAt:          timeParse("2024-12-19T11:57:33.40507Z"),
			ExpiresAt:          timeParse("2024-12-19T12:27:33.404897Z"),
			LastUsed:           TimePtr(timeParse("2024-12-19T12:01:13.538016Z")),
			CanCreateSubTokens: true,
			AllowedIPRanges:    []string{"0.0.0.0/0", "::/0"},
		},
	}
	got := Tokens{}
	err := json.Unmarshal([]byte(`
		[
		  {
			"id": "deadbeef-dead-beef-dead-beefdeadbee1",
			"name": "token_workspace",
			"type": "workspace",
			"created_at": "2024-12-19T11:46:09.888763Z",
			"expires_at": "2024-12-19T12:16:09.888531Z",
			"can_create_tokens": true,
			"allowed_ip_ranges": [
			  "0.0.0.0/0",
			  "::/0"
			]
		  },
		  {
			"id": "deadbeef-dead-beef-dead-beefdeadbee2",
			"name": "token_pat",
			"type": "pat",
			"created_at": "2024-12-19T11:57:33.40507Z",
			"expires_at": "2024-12-19T12:27:33.404897Z",
			"last_used_at": "2024-12-19T12:01:13.538016Z",
			"can_create_tokens": true,
			"allowed_ip_ranges": [
			  "0.0.0.0/0",
			  "::/0"
			]
		  }
		]
`), &got)
	assert.NoError(t, err)
	assert.Equal(t, want, got)
}

func TestTokenMarshal(t *testing.T) {
	want := `
		{
		  "token": "ucat_01DEADBEEFDEADBEEFDEADBEEF",
		  "id": "deadbeef-dead-beef-dead-beefdeadbeef",
		  "name": "test_token",
		  "type": "workspace",
		  "created_at": "2024-12-19T11:46:09.888763Z",
		  "expires_at": "2024-12-19T12:16:09.888531Z",
		  "can_create_tokens": true,
		  "allowed_ip_ranges": [
			"0.0.0.0/0",
			"::/0"
		  ]
		}
	`
	got, err := json.Marshal(&Token{
		APIToken:           "ucat_01DEADBEEFDEADBEEFDEADBEEF",
		ID:                 "deadbeef-dead-beef-dead-beefdeadbeef",
		Name:               "test_token",
		Type:               TokenTypeWorkspace,
		CreatedAt:          timeParse("2024-12-19T11:46:09.888763Z"),
		ExpiresAt:          timeParse("2024-12-19T12:16:09.888531Z"),
		LastUsed:           nil,
		CanCreateSubTokens: true,
		AllowedIPRanges:    []string{"0.0.0.0/0", "::/0"},
	})
	assert.NoError(t, err)
	assert.JSONEq(t, want, string(got))
}
