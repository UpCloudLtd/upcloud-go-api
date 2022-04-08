package service

import (
	"encoding/json"
	"fmt"

	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud/client"
)

type PriceZones interface {
	GetPriceZones() (*upcloud.PriceZones, error)
}

var _ PriceZones = (*Service)(nil)

type TimeZones interface {
	GetTimeZones() (*upcloud.TimeZones, error)
}

var _ TimeZones = (*Service)(nil)

type Plans interface {
	GetPlans() (*upcloud.Plans, error)
}

var _ Plans = (*Service)(nil)

// Service represents the API service. The specified client is used to communicate with the API
type Service struct {
	client *client.Client
}

// New constructs and returns a new service object configured with the specified client
func New(client *client.Client) *Service {
	service := Service{}
	service.client = client

	return &service
}

// GetZones returns the available zones
func (s *Service) GetZones() (*upcloud.Zones, error) {
	zones := upcloud.Zones{}
	response, err := s.basicGetRequest("/zone")

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(response, &zones)
	if err != nil {
		return nil, err
	}

	return &zones, nil
}

// GetPriceZones returns the available price zones and their corresponding prices
func (s *Service) GetPriceZones() (*upcloud.PriceZones, error) {
	zones := upcloud.PriceZones{}
	response, err := s.basicGetRequest("/price")

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(response, &zones)
	if err != nil {
		return nil, err
	}

	return &zones, nil
}

// GetTimeZones returns the available timezones
func (s *Service) GetTimeZones() (*upcloud.TimeZones, error) {
	zones := upcloud.TimeZones{}
	response, err := s.basicGetRequest("/timezone")

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(response, &zones)
	if err != nil {
		return nil, err
	}

	return &zones, nil
}

// GetPlans returns the available service plans
func (s *Service) GetPlans() (*upcloud.Plans, error) {
	plans := upcloud.Plans{}
	response, err := s.basicGetRequest("/plan")

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(response, &plans)
	if err != nil {
		return nil, err
	}

	return &plans, nil
}

// Wrapper that performs a GET request to the specified location and returns the response or a service error
func (s *Service) basicGetRequest(location string) ([]byte, error) {
	requestURL := s.client.CreateRequestURL(location)

	response, err := s.client.PerformJSONGetRequest(requestURL)

	if err != nil {
		return nil, parseJSONServiceError(err)
	}

	return response, nil
}

// Performs a GET request to the specified location and stores the result in the value pointed to by v.
func (s *Service) get(location string, v interface{}) error {
	res, err := s.basicGetRequest(location)
	if err != nil {
		return err
	}
	return json.Unmarshal(res, v)
}

// Parses an error returned from the client into corresponding error type
func parseJSONServiceError(err error) error {
	if clientError, ok := err.(*client.Error); ok {
		var serviceError error
		switch clientError.Type {
		case client.ErrorTypeProblem:
			serviceError = &upcloud.Problem{}
		default:
			serviceError = &upcloud.Error{}
		}
		if err := json.Unmarshal(clientError.ResponseBody, serviceError); err != nil {
			return fmt.Errorf("received malformed client error: %s", string(clientError.ResponseBody))
		}
		return serviceError
	}
	return err
}
