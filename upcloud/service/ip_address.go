package service

import (
	"encoding/json"

	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud/request"
)

type IpAddress interface {
	GetIPAddresses() (*upcloud.IPAddresses, error)
	GetIPAddressDetails(r *request.GetIPAddressDetailsRequest) (*upcloud.IPAddress, error)
	AssignIPAddress(r *request.AssignIPAddressRequest) (*upcloud.IPAddress, error)
	ModifyIPAddress(r *request.ModifyIPAddressRequest) (*upcloud.IPAddress, error)
	ReleaseIPAddress(r *request.ReleaseIPAddressRequest) error
}

var _ IpAddress = (*Service)(nil)

// GetIPAddresses returns all IP addresses associated with the account
func (s *Service) GetIPAddresses() (*upcloud.IPAddresses, error) {
	ipAddresses := upcloud.IPAddresses{}
	response, err := s.basicGetRequest("/ip_address")

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(response, &ipAddresses)
	if err != nil {
		return nil, err
	}

	return &ipAddresses, nil
}

// GetIPAddressDetails returns extended details about the specified IP address
func (s *Service) GetIPAddressDetails(r *request.GetIPAddressDetailsRequest) (*upcloud.IPAddress, error) {
	ipAddress := upcloud.IPAddress{}
	response, err := s.basicGetRequest(r.RequestURL())

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(response, &ipAddress)
	if err != nil {
		return nil, err
	}

	return &ipAddress, nil
}

// AssignIPAddress assigns the specified IP address to the specified server
func (s *Service) AssignIPAddress(r *request.AssignIPAddressRequest) (*upcloud.IPAddress, error) {
	ipAddress := upcloud.IPAddress{}
	requestBody, _ := json.Marshal(r)
	response, err := s.client.PerformJSONPostRequest(s.client.CreateRequestURL(r.RequestURL()), requestBody)

	if err != nil {
		return nil, parseJSONServiceError(err)
	}

	err = json.Unmarshal(response, &ipAddress)
	if err != nil {
		return nil, err
	}

	return &ipAddress, nil
}

// ModifyIPAddress modifies the specified IP address
func (s *Service) ModifyIPAddress(r *request.ModifyIPAddressRequest) (*upcloud.IPAddress, error) {
	ipAddress := upcloud.IPAddress{}
	requestBody, _ := json.Marshal(r)
	response, err := s.client.PerformJSONPatchRequest(s.client.CreateRequestURL(r.RequestURL()), requestBody)

	if err != nil {
		return nil, parseJSONServiceError(err)
	}

	err = json.Unmarshal(response, &ipAddress)
	if err != nil {
		return nil, err
	}

	return &ipAddress, nil
}

// ReleaseIPAddress releases the specified IP address from the server it is attached to
func (s *Service) ReleaseIPAddress(r *request.ReleaseIPAddressRequest) error {
	err := s.client.PerformJSONDeleteRequest(s.client.CreateRequestURL(r.RequestURL()))

	if err != nil {
		return parseJSONServiceError(err)
	}

	return nil
}
