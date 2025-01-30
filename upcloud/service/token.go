package service

import (
	"context"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
)

type Token interface {
	CreateToken(context.Context, *request.CreateTokenRequest) (*upcloud.Token, error)
	GetTokenDetails(ctx context.Context, r *request.GetTokenDetailsRequest) (*upcloud.Token, error)
	GetTokens(context.Context, *request.GetTokensRequest) (*upcloud.Tokens, error)
	DeleteToken(context.Context, *request.DeleteTokenRequest) error
}

// CreateToken creates a new token.
func (s *Service) CreateToken(ctx context.Context, r *request.CreateTokenRequest) (*upcloud.Token, error) {
	token := upcloud.Token{}
	return &token, s.create(ctx, r, &token)
}

// GetTokenDetails returns the details for the specified token.
func (s *Service) GetTokenDetails(ctx context.Context, r *request.GetTokenDetailsRequest) (*upcloud.Token, error) {
	token := upcloud.Token{}
	return &token, s.get(ctx, r.RequestURL(), &token)
}

// GetTokens returns the all the available networks
func (s *Service) GetTokens(ctx context.Context, req *request.GetTokensRequest) (*upcloud.Tokens, error) {
	tokens := upcloud.Tokens{}
	return &tokens, s.get(ctx, req.RequestURL(), &tokens)
}

// DeleteToken deletes the specified token.
func (s *Service) DeleteToken(ctx context.Context, r *request.DeleteTokenRequest) error {
	return s.delete(ctx, r)
}
