package service

import (
	"context"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
)

type Partner interface {
	CreatePartnerAccount(ctx context.Context, r *request.CreatePartnerAccountRequest) (*upcloud.PartnerAccount, error)
	GetPartnerAccounts(ctx context.Context) ([]upcloud.PartnerAccount, error)
}

// CreatePartnerAccount creates new main account for partner
func (s *Service) CreatePartnerAccount(ctx context.Context, r *request.CreatePartnerAccountRequest) (*upcloud.PartnerAccount, error) {
	partnerAccount := upcloud.PartnerAccount{}
	return &partnerAccount, s.create(ctx, r, &partnerAccount)
}

// GetPartnerAccounts lists accounts associated with partner
func (s *Service) GetPartnerAccounts(ctx context.Context) ([]upcloud.PartnerAccount, error) {
	accounts := make([]upcloud.PartnerAccount, 0)
	return accounts, s.get(ctx, "/partner/accounts", &accounts)
}
