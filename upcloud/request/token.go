package request

import (
	"fmt"
	"time"
)

// GetTokenDetailsRequest represents a request to get token details. Will not return the actual API token.
type GetTokenDetailsRequest struct {
	ID string
}

func (r *GetTokenDetailsRequest) RequestURL() string {
	return fmt.Sprintf("/account/tokens/%s", r.ID)
}

// GetTokensRequest represents a request to get a list of tokens. Will not return the actual API tokens.
type GetTokensRequest struct{}

func (r *GetTokensRequest) RequestURL() string {
	return "/account/tokens"
}

// CreateTokenRequest represents a request to create a new network.
type CreateTokenRequest struct {
	Name               string    `json:"name"`
	ExpiresAt          time.Time `json:"expires_at"`
	CanCreateSubTokens bool      `json:"can_create_tokens"`
	AllowedIPRanges    []string  `json:"allowed_ip_ranges"`
}

// RequestURL implements the Request interface.
func (r *CreateTokenRequest) RequestURL() string {
	return "/account/tokens"
}

// DeleteTokenRequest represents a request to delete a token.
type DeleteTokenRequest struct {
	ID string
}

// RequestURL implements the Request interface.
func (r *DeleteTokenRequest) RequestURL() string {
	return fmt.Sprintf("/account/tokens/%s", r.ID)
}
