package service

import (
	"context"

	"github.com/UpCloudLtd/upcloud-go-api/v6/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v6/upcloud/request"
)

type IPAddress interface {
	GetIPAddresses(ctx context.Context) (*upcloud.IPAddresses, error)
	GetIPAddressDetails(ctx context.Context, r *request.GetIPAddressDetailsRequest) (*upcloud.IPAddress, error)
	AssignIPAddress(ctx context.Context, r *request.AssignIPAddressRequest) (*upcloud.IPAddress, error)
	ModifyIPAddress(ctx context.Context, r *request.ModifyIPAddressRequest) (*upcloud.IPAddress, error)
	ReleaseIPAddress(ctx context.Context, r *request.ReleaseIPAddressRequest) error
}

// GetIPAddresses returns all IP addresses associated with the account
func (s *Service) GetIPAddresses(ctx context.Context) (*upcloud.IPAddresses, error) {
	ipAddresses := upcloud.IPAddresses{}
	return &ipAddresses, s.get(ctx, "/ip_address", &ipAddresses)
}

// GetIPAddressDetails returns extended details about the specified IP address
func (s *Service) GetIPAddressDetails(ctx context.Context, r *request.GetIPAddressDetailsRequest) (*upcloud.IPAddress, error) {
	ipAddress := upcloud.IPAddress{}
	return &ipAddress, s.get(ctx, r.RequestURL(), &ipAddress)
}

// AssignIPAddress assigns the specified IP address to the specified server
func (s *Service) AssignIPAddress(ctx context.Context, r *request.AssignIPAddressRequest) (*upcloud.IPAddress, error) {
	ipAddress := upcloud.IPAddress{}
	return &ipAddress, s.create(ctx, r, &ipAddress)
}

// ModifyIPAddress modifies the specified IP address
func (s *Service) ModifyIPAddress(ctx context.Context, r *request.ModifyIPAddressRequest) (*upcloud.IPAddress, error) {
	ipAddress := upcloud.IPAddress{}
	return &ipAddress, s.modify(ctx, r, &ipAddress)
}

// ReleaseIPAddress releases the specified IP address from the server it is attached to
func (s *Service) ReleaseIPAddress(ctx context.Context, r *request.ReleaseIPAddressRequest) error {
	return s.delete(ctx, r)
}
