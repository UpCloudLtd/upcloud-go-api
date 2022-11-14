package service

import (
	"encoding/json"

	"github.com/UpCloudLtd/upcloud-go-api/v5/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v5/upcloud/request"
)

type Account interface {
	GetAccountList() (upcloud.AccountList, error)
	GetAccount() (*upcloud.Account, error)
	GetAccountDetails(r *request.GetAccountDetailsRequest) (*upcloud.AccountDetails, error)
	CreateSubaccount(r *request.CreateSubaccountRequest) (*upcloud.AccountDetails, error)
	ModifySubaccount(r *request.ModifySubaccountRequest) (*upcloud.AccountDetails, error)
	DeleteSubaccount(r *request.DeleteSubaccountRequest) error
}

var _ Account = (*Service)(nil)

// GetAccount returns the current user's account
func (s *Service) GetAccount() (*upcloud.Account, error) {
	account := upcloud.Account{}
	response, err := s.basicGetRequest("/account")
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(response, &account)
	if err != nil {
		return nil, err
	}

	return &account, nil
}

// GetAccountList returns the account list
func (s *Service) GetAccountList() (upcloud.AccountList, error) {
	accountList := make(upcloud.AccountList, 0)

	response, err := s.basicGetRequest("/account/list")
	if err != nil {
		return accountList, err
	}

	if err = json.Unmarshal(response, &accountList); err != nil {
		return accountList, err
	}

	return accountList, nil
}

// GetAccountDetails returns account details
func (s *Service) GetAccountDetails(r *request.GetAccountDetailsRequest) (*upcloud.AccountDetails, error) {
	account := &upcloud.AccountDetails{}
	response, err := s.basicGetRequest(r.RequestURL())
	if err != nil {
		return account, err
	}

	if err = json.Unmarshal(response, account); err != nil {
		return account, err
	}

	return account, nil
}

// ModifySubaccount modifies a sub account
func (s *Service) ModifySubaccount(r *request.ModifySubaccountRequest) (*upcloud.AccountDetails, error) {
	requestBody, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	if _, err = s.client.PerformJSONPutRequest(s.client.CreateRequestURL(r.RequestURL()), requestBody); err != nil {
		return nil, parseJSONServiceError(err)
	}

	return s.GetAccountDetails(&request.GetAccountDetailsRequest{Username: r.Username})
}

// CreateSubaccount creates a new sub account
func (s *Service) CreateSubaccount(r *request.CreateSubaccountRequest) (*upcloud.AccountDetails, error) {
	account := &upcloud.AccountDetails{}
	requestBody, err := json.Marshal(r)
	if err != nil {
		return account, err
	}

	if _, err := s.client.PerformJSONPostRequest(s.client.CreateRequestURL(r.RequestURL()), requestBody); err != nil {
		return account, parseJSONServiceError(err)
	}

	return s.GetAccountDetails(&request.GetAccountDetailsRequest{Username: r.Subaccount.Username})
}

// DeleteSubaccount deletes a sub account
func (s *Service) DeleteSubaccount(r *request.DeleteSubaccountRequest) error {
	if err := s.client.PerformJSONDeleteRequest(s.client.CreateRequestURL(r.RequestURL())); err != nil {
		return parseJSONServiceError(err)
	}
	return nil
}
