package service

import (
	"context"

	"github.com/UpCloudLtd/upcloud-go-api/v6/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v6/upcloud/request"
)

type Gateway interface {
	GetGateways(ctx context.Context, f ...request.QueryFilter) ([]upcloud.Gateway, error)
	GetGateway(ctx context.Context, r *request.GetGatewayRequest) (*upcloud.Gateway, error)
	CreateGateway(ctx context.Context, r *request.CreateGatewayRequest) (*upcloud.Gateway, error)
	ModifyGateway(ctx context.Context, r *request.ModifyGatewayRequest) (*upcloud.Gateway, error)
	DeleteGateway(ctx context.Context, r *request.DeleteGatewayRequest) error
}

// GetGateways retrieves a list of network gateways within an account.
func (s *Service) GetGateways(ctx context.Context, f ...request.QueryFilter) ([]upcloud.Gateway, error) {
	r := request.GetGatewaysRequest{Filters: f}
	p := []upcloud.Gateway{}
	return p, s.get(ctx, r.RequestURL(), &p)
}

// GetGateway retrieves details of a network gateway.
func (s *Service) GetGateway(ctx context.Context, r *request.GetGatewayRequest) (*upcloud.Gateway, error) {
	p := upcloud.Gateway{}
	return &p, s.get(ctx, r.RequestURL(), &p)
}

// CreateGateway creates a new network gateway.
func (s *Service) CreateGateway(ctx context.Context, r *request.CreateGatewayRequest) (*upcloud.Gateway, error) {
	p := upcloud.Gateway{}
	return &p, s.create(ctx, r, &p)
}

// ModifyGateway modifies an existing network gateway.
func (s *Service) ModifyGateway(ctx context.Context, r *request.ModifyGatewayRequest) (*upcloud.Gateway, error) {
	p := upcloud.Gateway{}
	return &p, s.modify(ctx, r, &p)
}

// DeleteGateway deletes a network gateway.
func (s *Service) DeleteGateway(ctx context.Context, r *request.DeleteGatewayRequest) error {
	return s.delete(ctx, r)
}
