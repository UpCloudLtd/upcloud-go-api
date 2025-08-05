package service

import (
	"context"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
)

type Gateway interface {
	GetGatewayPlans(ctx context.Context) ([]upcloud.GatewayPlan, error)
	GetGateways(ctx context.Context, f ...request.QueryFilter) ([]upcloud.Gateway, error)
	GetGateway(ctx context.Context, r *request.GetGatewayRequest) (*upcloud.Gateway, error)
	CreateGateway(ctx context.Context, r *request.CreateGatewayRequest) (*upcloud.Gateway, error)
	ModifyGateway(ctx context.Context, r *request.ModifyGatewayRequest) (*upcloud.Gateway, error)
	DeleteGateway(ctx context.Context, r *request.DeleteGatewayRequest) error

	GetGatewayConnections(ctx context.Context, r *request.GetGatewayConnectionsRequest) ([]upcloud.GatewayConnection, error)
	GetGatewayConnection(ctx context.Context, r *request.GetGatewayConnectionRequest) (*upcloud.GatewayConnection, error)
	CreateGatewayConnection(ctx context.Context, r *request.CreateGatewayConnectionRequest) (*upcloud.GatewayConnection, error)
	ModifyGatewayConnection(ctx context.Context, r *request.ModifyGatewayConnectionRequest) (*upcloud.GatewayConnection, error)
	DeleteGatewayConnection(ctx context.Context, r *request.DeleteGatewayConnectionRequest) error

	GetGatewayConnectionTunnels(ctx context.Context, r *request.GetGatewayConnectionTunnelsRequest) ([]upcloud.GatewayTunnel, error)
	GetGatewayConnectionTunnel(ctx context.Context, r *request.GetGatewayConnectionTunnelRequest) (*upcloud.GatewayTunnel, error)
	CreateGatewayConnectionTunnel(ctx context.Context, r *request.CreateGatewayConnectionTunnelRequest) (*upcloud.GatewayTunnel, error)
	// ModifyGatewayConnectionTunnel(ctx context.Context, r *request.ModifyGatewayConnectionTunnelRequest) (*upcloud.GatewayTunnel, error)
	DeleteGatewayConnectionTunnel(ctx context.Context, r *request.DeleteGatewayConnectionTunnelRequest) error
	GetGatewayMetrics(ctx context.Context, r *request.GetGatewayMetricsRequest) (*upcloud.GatewayMetrics, error)
}

// GetGatewayPlans retrieves a list of all available plans for network gateway service
func (s *Service) GetGatewayPlans(ctx context.Context) ([]upcloud.GatewayPlan, error) {
	r := request.GetGatewayPlansRequest{}
	p := []upcloud.GatewayPlan{}
	return p, s.get(ctx, r.RequestURL(), &p)
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

// GetGatewayConnections retrieves a list of specific gateway connections
func (s *Service) GetGatewayConnections(ctx context.Context, r *request.GetGatewayConnectionsRequest) ([]upcloud.GatewayConnection, error) {
	p := []upcloud.GatewayConnection{}
	return p, s.get(ctx, r.RequestURL(), &p)
}

// GetGatewayConnection retrieves details of a specific network gateway connection
func (s *Service) GetGatewayConnection(ctx context.Context, r *request.GetGatewayConnectionRequest) (*upcloud.GatewayConnection, error) {
	p := upcloud.GatewayConnection{}
	return &p, s.get(ctx, r.RequestURL(), &p)
}

// CreateGatewayConnection creates a new connection for a specific gateway
func (s *Service) CreateGatewayConnection(ctx context.Context, r *request.CreateGatewayConnectionRequest) (*upcloud.GatewayConnection, error) {
	p := upcloud.GatewayConnection{}
	return &p, s.create(ctx, r, &p)
}

func (s *Service) ModifyGatewayConnection(ctx context.Context, r *request.ModifyGatewayConnectionRequest) (*upcloud.GatewayConnection, error) {
	p := upcloud.GatewayConnection{}
	return &p, s.modify(ctx, r, &p)
}

// DeleteGatewayConnection deletes a specific connection of a network gateway
func (s *Service) DeleteGatewayConnection(ctx context.Context, r *request.DeleteGatewayConnectionRequest) error {
	return s.delete(ctx, r)
}

// GetGatewayConnectionTunnels retrieves tunnels for specific connection of specific gateway
func (s *Service) GetGatewayConnectionTunnels(ctx context.Context, r *request.GetGatewayConnectionTunnelsRequest) ([]upcloud.GatewayTunnel, error) {
	p := []upcloud.GatewayTunnel{}
	return p, s.get(ctx, r.RequestURL(), &p)
}

// GetGatewayConnectionTunnel retrieves a single tunnel details for specific connection of specific gateway
func (s *Service) GetGatewayConnectionTunnel(ctx context.Context, r *request.GetGatewayConnectionTunnelRequest) (*upcloud.GatewayTunnel, error) {
	p := upcloud.GatewayTunnel{}
	return &p, s.get(ctx, r.RequestURL(), &p)
}

// CreateGatewayConnectionTunnel creates a tunnel for specific connection of specific gateway
func (s *Service) CreateGatewayConnectionTunnel(ctx context.Context, r *request.CreateGatewayConnectionTunnelRequest) (*upcloud.GatewayTunnel, error) {
	p := upcloud.GatewayTunnel{}
	return &p, s.create(ctx, r, &p)
}

// ModifyGatewayConnectionTunnel modifies a single tunnel for specific connection of specific gateway
func (s *Service) ModifyGatewayConnectionTunnel(ctx context.Context, r *request.ModifyGatewayConnectionTunnelRequest) (*upcloud.GatewayTunnel, error) {
	p := upcloud.GatewayTunnel{}
	return &p, s.modify(ctx, r, &p)
}

// DeleteGatewayConnectionTunnel deletes a tunnel for specific connection of specific gateway
func (s *Service) DeleteGatewayConnectionTunnel(ctx context.Context, r *request.DeleteGatewayConnectionTunnelRequest) error {
	return s.delete(ctx, r)
}

// GetGatewayMetrics retrieves metrics for a specific gateway service
func (s *Service) GetGatewayMetrics(ctx context.Context, r *request.GetGatewayMetricsRequest) (*upcloud.GatewayMetrics, error) {
	metrics := upcloud.GatewayMetrics{}
	return &metrics, s.get(ctx, r.RequestURL(), &metrics)
}
