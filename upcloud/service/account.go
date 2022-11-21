package service

import (
	"context"

	"github.com/UpCloudLtd/upcloud-go-api/v5/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v5/upcloud/request"
)

type Account interface {
	GetAccountList(ctx context.Context) (upcloud.AccountList, error)
	GetAccount(ctx context.Context) (*upcloud.Account, error)
	GetAccountDetails(ctx context.Context, r *request.GetAccountDetailsRequest) (*upcloud.AccountDetails, error)
	CreateSubaccount(ctx context.Context, r *request.CreateSubaccountRequest) (*upcloud.AccountDetails, error)
	ModifySubaccount(ctx context.Context, r *request.ModifySubaccountRequest) (*upcloud.AccountDetails, error)
	DeleteSubaccount(ctx context.Context, r *request.DeleteSubaccountRequest) error
}

// GetAccount returns the current user's account
func (s *Service) GetAccount(ctx context.Context) (*upcloud.Account, error) {
	account := upcloud.Account{}
	return &account, s.get(ctx, "/account", &account)
}

// GetAccountList returns the account list
func (s *Service) GetAccountList(ctx context.Context) (upcloud.AccountList, error) {
	accountList := make(upcloud.AccountList, 0)
	return accountList, s.get(ctx, "/account/list", &accountList)
}

// GetAccountDetails returns account details
func (s *Service) GetAccountDetails(ctx context.Context, r *request.GetAccountDetailsRequest) (*upcloud.AccountDetails, error) {
	account := upcloud.AccountDetails{}
	return &account, s.get(ctx, r.RequestURL(), &account)
}

// ModifySubaccount modifies a sub account
func (s *Service) ModifySubaccount(ctx context.Context, r *request.ModifySubaccountRequest) (*upcloud.AccountDetails, error) {
	if err := s.replace(ctx, r, nil); err != nil {
		return nil, err
	}
	return s.GetAccountDetails(ctx, &request.GetAccountDetailsRequest{Username: r.Username})
}

// CreateSubaccount creates a new sub account
func (s *Service) CreateSubaccount(ctx context.Context, r *request.CreateSubaccountRequest) (*upcloud.AccountDetails, error) {
	if err := s.create(ctx, r, nil); err != nil {
		return nil, err
	}
	return s.GetAccountDetails(ctx, &request.GetAccountDetailsRequest{Username: r.Subaccount.Username})
}

// DeleteSubaccount deletes a sub account
func (s *Service) DeleteSubaccount(ctx context.Context, r *request.DeleteSubaccountRequest) error {
	return s.delete(ctx, r)
}
