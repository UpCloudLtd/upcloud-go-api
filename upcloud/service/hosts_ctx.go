package service

import (
	"context"

	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud/request"
)

type HostContext interface {
	GetHosts(ctx context.Context) (*upcloud.Hosts, error)
	GetHostDetails(ctx context.Context, r *request.GetHostDetailsRequest) (*upcloud.Host, error)
	ModifyHost(ctx context.Context, r *request.ModifyHostRequest) (*upcloud.Host, error)
}

// GetHosts returns the all the available private hosts
func (s *ServiceContext) GetHosts(ctx context.Context) (*upcloud.Hosts, error) {
	hosts := upcloud.Hosts{}
	return &hosts, s.get(ctx, "/host", &hosts)
}

// GetHostDetails returns the details for a single private host
func (s *ServiceContext) GetHostDetails(ctx context.Context, r *request.GetHostDetailsRequest) (*upcloud.Host, error) {
	host := upcloud.Host{}
	return &host, s.get(ctx, r.RequestURL(), &host)
}

// ModifyHost modifies the configuration of an existing host.
func (s *ServiceContext) ModifyHost(ctx context.Context, r *request.ModifyHostRequest) (*upcloud.Host, error) {
	host := upcloud.Host{}
	return &host, s.modify(ctx, r, &host)
}
