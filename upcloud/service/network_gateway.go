package service

import (
	"context"

	"github.com/UpCloudLtd/upcloud-go-api/v5/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v5/upcloud/request"
)

type NetworkGateway interface {
	GetNetworkGateways(ctx context.Context, f ...request.QueryFilter) ([]upcloud.NetworkGateway, error)
	GetNetworkGateway(ctx context.Context, r *request.GetNetworkGatewayRequest) (*upcloud.NetworkGateway, error)
	CreateNetworkGateway(ctx context.Context, r *request.CreateNetworkGatewayRequest) (*upcloud.NetworkGateway, error)
	ModifyNetworkGateway(ctx context.Context, r *request.ModifyNetworkGatewayRequest) (*upcloud.NetworkGateway, error)
	DeleteNetworkGateway(ctx context.Context, r *request.DeleteNetworkGatewayRequest) error
}

// GetNetworkGateways retrieves a list of network gateways within an account.
func (s *Service) GetNetworkGateways(ctx context.Context, f ...request.QueryFilter) ([]upcloud.NetworkGateway, error) {
	r := request.GetNetworkGatewaysRequest{Filters: f}
	p := []upcloud.NetworkGateway{}
	return p, s.get(ctx, r.RequestURL(), &p)
}

// GetNetworkGateway retrieves details of a network gateway.
func (s *Service) GetNetworkGateway(ctx context.Context, r *request.GetNetworkGatewayRequest) (*upcloud.NetworkGateway, error) {
	p := upcloud.NetworkGateway{}
	return &p, s.get(ctx, r.RequestURL(), &p)
}

// CreateNetworkGateway creates a new network gateway.
func (s *Service) CreateNetworkGateway(ctx context.Context, r *request.CreateNetworkGatewayRequest) (*upcloud.NetworkGateway, error) {
	p := upcloud.NetworkGateway{}
	return &p, s.create(ctx, r, &p)
}

// ModifyNetworkGateway modifies an existing network gateway.
func (s *Service) ModifyNetworkGateway(ctx context.Context, r *request.ModifyNetworkGatewayRequest) (*upcloud.NetworkGateway, error) {
	p := upcloud.NetworkGateway{}
	return &p, s.modify(ctx, r, &p)
}

// DeleteNetworkGateway deletes a gateway. Peering can be deleted only when the state is disabled.
func (s *Service) DeleteNetworkGateway(ctx context.Context, r *request.DeleteNetworkGatewayRequest) error {
	return s.delete(ctx, r)
}
