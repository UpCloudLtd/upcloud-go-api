package request

import (
	"fmt"
	"time"
)

const basePath = "/account/tokens"

// GetTokenDetailsRequest (EXPERIMENTAL) represents a request to get token details. Will not return the actual API token.
type GetTokenDetailsRequest struct {
	ID string
}

// RequestURL (EXPERIMENTAL) implements the Request interface.
func (r *GetTokenDetailsRequest) RequestURL() string {
	return fmt.Sprintf("%s/%s", basePath, r.ID)
}

// GetTokensRequest (EXPERIMENTAL) represents a request to get a list of tokens. Will not return the actual API tokens.
type GetTokensRequest struct {
	Page *Page
}

// RequestURL (EXPERIMENTAL) implements the Request interface.
func (r *GetTokensRequest) RequestURL() string {
	if r.Page != nil {
		f := make([]QueryFilter, 0)
		f = append(f, r.Page)
		return fmt.Sprintf("%s?%s", basePath, encodeQueryFilters(f))
	}

	return basePath
}

// CreateTokenRequest (EXPERIMENTAL) represents a request to create a new network.
type CreateTokenRequest struct {
	Name               string    `json:"name"`
	ExpiresAt          time.Time `json:"expires_at"`
	CanCreateSubTokens bool      `json:"can_create_tokens"`
	AllowedIPRanges    []string  `json:"allowed_ip_ranges"`
}

// RequestURL (EXPERIMENTAL) implements the Request interface.
func (r *CreateTokenRequest) RequestURL() string {
	return basePath
}

// DeleteTokenRequest (EXPERIMENTAL) represents a request to delete a token.
type DeleteTokenRequest struct {
	ID string
}

// RequestURL (EXPERIMENTAL) implements the Request interface.
func (r *DeleteTokenRequest) RequestURL() string {
	return fmt.Sprintf("%s/%s", basePath, r.ID)
}
