package upcloud

import (
	"time"
)

type TokenType string

const (
	TokenTypePAT       = "pat"
	TokenTypeWorkspace = "workspace"
)

type Tokens []Token

type Token struct {
	APIToken           string     `json:"token,omitempty"` // APIToken is the API token. Returned only when creating a new token.
	ID                 string     `json:"id"`
	Name               string     `json:"name"`
	Type               string     `json:"type"`
	Created            time.Time  `json:"created_at"`
	Expires            time.Time  `json:"expires_at"`
	LastUsed           *time.Time `json:"last_used_at,omitempty"`
	CanCreateSubTokens bool       `json:"can_create_tokens"`
	AllowedIPRanges    []string   `json:"allowed_ip_ranges"`
}
