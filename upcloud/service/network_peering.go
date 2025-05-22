package service

import (
	"context"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
)

type NetworkPeering interface {
	GetNetworkPeerings(ctx context.Context, f ...request.QueryFilter) (upcloud.NetworkPeerings, error)
	GetNetworkPeering(ctx context.Context, r *request.GetNetworkPeeringRequest) (*upcloud.NetworkPeering, error)
	CreateNetworkPeering(ctx context.Context, r *request.CreateNetworkPeeringRequest) (*upcloud.NetworkPeering, error)
	ModifyNetworkPeering(ctx context.Context, r *request.ModifyNetworkPeeringRequest) (*upcloud.NetworkPeering, error)
	DeleteNetworkPeering(ctx context.Context, r *request.DeleteNetworkPeeringRequest) error
	WaitForNetworkPeeringState(ctx context.Context, r *request.WaitForNetworkPeeringStateRequest) (*upcloud.NetworkPeering, error)
}

// GetNetworkPeerings (EXPERIMENTAL) retrieves a list of network peerings within an account.
func (s *Service) GetNetworkPeerings(ctx context.Context, f ...request.QueryFilter) (upcloud.NetworkPeerings, error) {
	r := request.GetNetworkPeeringsRequest{Filters: f}
	p := upcloud.NetworkPeerings{}
	return p, s.get(ctx, r.RequestURL(), &p)
}

// GetNetworkPeering (EXPERIMENTAL) retrieves details of a network peering.
func (s *Service) GetNetworkPeering(ctx context.Context, r *request.GetNetworkPeeringRequest) (*upcloud.NetworkPeering, error) {
	p := upcloud.NetworkPeering{}
	return &p, s.get(ctx, r.RequestURL(), &p)
}

// CreateNetworkPeering (EXPERIMENTAL) creates a new network peering.
func (s *Service) CreateNetworkPeering(ctx context.Context, r *request.CreateNetworkPeeringRequest) (*upcloud.NetworkPeering, error) {
	p := upcloud.NetworkPeering{}
	return &p, s.create(ctx, r, &p)
}

// ModifyNetworkPeering (EXPERIMENTAL) modifies an existing network peering.
func (s *Service) ModifyNetworkPeering(ctx context.Context, r *request.ModifyNetworkPeeringRequest) (*upcloud.NetworkPeering, error) {
	p := upcloud.NetworkPeering{}
	return &p, s.modify(ctx, r, &p)
}

// DeleteNetworkPeering (EXPERIMENTAL) deletes a peering. Peering can be deleted only when the state is disabled.
func (s *Service) DeleteNetworkPeering(ctx context.Context, r *request.DeleteNetworkPeeringRequest) error {
	return s.delete(ctx, r)
}

func (s *Service) WaitForNetworkPeeringState(ctx context.Context, r *request.WaitForNetworkPeeringStateRequest) (*upcloud.NetworkPeering, error) {
	return wait(ctx, func(_ int, c context.Context) (*upcloud.NetworkPeering, error) {
		details, err := s.GetNetworkPeering(c, &request.GetNetworkPeeringRequest{
			UUID: r.UUID,
		})
		if err != nil {
			return nil, err
		}

		if details.State == r.DesiredState {
			return details, nil
		}
		return nil, nil
	}, nil)
}
