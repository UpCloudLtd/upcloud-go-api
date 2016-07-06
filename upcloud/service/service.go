package service

import (
	"encoding/xml"
	"fmt"
	"github.com/jalle19/upcloud-go-sdk/upcloud"
	"github.com/jalle19/upcloud-go-sdk/upcloud/client"
	"github.com/jalle19/upcloud-go-sdk/upcloud/request"
	"time"
)

/**
Represents the API service. The specified client is used to communicate with the API
*/
type Service struct {
	client *client.Client
}

/**
Constructs and returns a new service object configured with the specified client
*/
func New(client *client.Client) *Service {
	service := Service{}
	service.client = client

	return &service
}

/**
Returns the current user's account
*/
func (s *Service) GetAccount() (*upcloud.Account, error) {
	account := upcloud.Account{}
	response, err := s.basicGetRequest("/account")

	if err != nil {
		return nil, err
	}

	xml.Unmarshal(response, &account)

	return &account, nil
}

/**
Returns the available zones
*/
func (s *Service) GetZones() (*upcloud.Zones, error) {
	zones := upcloud.Zones{}
	response, err := s.basicGetRequest("/zone")

	if err != nil {
		return nil, err
	}

	xml.Unmarshal(response, &zones)

	return &zones, nil
}

/**
Returns the available price zones and their corresponding prices
*/
func (s *Service) GetPriceZones() (*upcloud.PrizeZones, error) {
	zones := upcloud.PrizeZones{}
	response, err := s.basicGetRequest("/price")

	if err != nil {
		return nil, err
	}

	xml.Unmarshal(response, &zones)

	return &zones, nil
}

/**
Returns the available timezones
*/
func (s *Service) GetTimeZones() (*upcloud.TimeZones, error) {
	zones := upcloud.TimeZones{}
	response, err := s.basicGetRequest("/timezone")

	if err != nil {
		return nil, err
	}

	xml.Unmarshal(response, &zones)

	return &zones, nil
}

/**
Returns the available service plans
*/
func (s *Service) GetPlans() (*upcloud.Plans, error) {
	plans := upcloud.Plans{}
	response, err := s.basicGetRequest("/plan")

	if err != nil {
		return nil, err
	}

	xml.Unmarshal(response, &plans)

	return &plans, nil
}

/**
Returns the available pre-configured server configurations
*/
func (s *Service) GetServerConfigurations() (*upcloud.ServerConfigurations, error) {
	serverConfigurations := upcloud.ServerConfigurations{}
	response, err := s.basicGetRequest("/server_size")

	if err != nil {
		return nil, err
	}

	xml.Unmarshal(response, &serverConfigurations)

	return &serverConfigurations, nil
}

/**
Returns the available servers
*/
func (s *Service) GetServers() (*upcloud.Servers, error) {
	servers := upcloud.Servers{}
	response, err := s.basicGetRequest("/server")

	if err != nil {
		return nil, err
	}

	xml.Unmarshal(response, &servers)

	return &servers, nil
}

/**
Returns extended details about the specified server
*/
func (s *Service) GetServerDetails(r *request.GetServerDetailsRequest) (*upcloud.ServerDetails, error) {
	serverDetails := upcloud.ServerDetails{}
	response, err := s.basicGetRequest(r.RequestURL())

	if err != nil {
		return nil, err
	}

	xml.Unmarshal(response, &serverDetails)

	return &serverDetails, nil
}

/**
Creates a server and returns the server details for the newly created server
*/
func (s *Service) CreateServer(r *request.CreateServerRequest) (*upcloud.ServerDetails, error) {
	serverDetails := upcloud.ServerDetails{}
	requestBody, _ := xml.Marshal(r)
	response, err := s.client.PerformPostRequest(s.client.CreateRequestUrl(r.RequestURL()), requestBody)

	if err != nil {
		return nil, parseServiceError(err)
	}

	xml.Unmarshal(response, &serverDetails)

	return &serverDetails, nil
}

/**
Blocks execution until the specified server has entered the specified state. The method will give up after the
specified timeout
*/
func (s *Service) WaitForServerState(r *request.WaitForServerStateRequest) error {
	attempts := 0
	sleepDuration := time.Second * 5

	for {
		attempts++

		serverDetails, err := s.GetServerDetails(&request.GetServerDetailsRequest{
			UUID: r.UUID,
		})

		if err != nil {
			return err
		}

		if serverDetails.State == r.DesiredState {
			return nil
		}

		time.Sleep(sleepDuration)

		if time.Duration(attempts)*sleepDuration >= r.Timeout {
			return fmt.Errorf("Timeout reached while waiting for server to enter state \"%s\"", r.DesiredState)
		}
	}
}

/**
Starts the specified server
*/
func (s *Service) StartServer(r *request.StartServerRequest) (*upcloud.ServerDetails, error) {
	// Increase the client timeout to match the request timeout
	s.client.SetTimeout(r.Timeout)

	serverDetails := upcloud.ServerDetails{}
	response, err := s.client.PerformPostRequest(s.client.CreateRequestUrl(r.RequestURL()), nil)

	if err != nil {
		return nil, parseServiceError(err)
	}

	xml.Unmarshal(response, &serverDetails)

	return &serverDetails, nil
}

/**
Stops the specified server
*/
func (s *Service) StopServer(r *request.StopServerRequest) (*upcloud.ServerDetails, error) {
	// Increase the client timeout to match the request timeout
	s.client.SetTimeout(r.Timeout)

	serverDetails := upcloud.ServerDetails{}
	requestBody, _ := xml.Marshal(r)
	response, err := s.client.PerformPostRequest(s.client.CreateRequestUrl(r.RequestURL()), requestBody)

	if err != nil {
		return nil, parseServiceError(err)
	}

	xml.Unmarshal(response, &serverDetails)

	return &serverDetails, nil
}

/**
Restarts the specified server
*/
func (s *Service) RestartServer(r *request.RestartServerRequest) (*upcloud.ServerDetails, error) {
	// Increase the client timeout to match the request timeout
	s.client.SetTimeout(r.Timeout)

	serverDetails := upcloud.ServerDetails{}
	requestBody, _ := xml.Marshal(r)
	response, err := s.client.PerformPostRequest(s.client.CreateRequestUrl(r.RequestURL()), requestBody)

	if err != nil {
		return nil, parseServiceError(err)
	}

	xml.Unmarshal(response, &serverDetails)

	return &serverDetails, nil
}

/**
Modifies the configuration of an existing server. Attaching and detaching storages as well as assigning and releasing
IP addresses have their own separate operations
*/
func (s *Service) ModifyServer(r *request.ModifyServerRequest) (*upcloud.ServerDetails, error) {
	serverDetails := upcloud.ServerDetails{}
	requestBody, _ := xml.Marshal(r)
	response, err := s.client.PerformPutRequest(s.client.CreateRequestUrl(r.RequestURL()), requestBody)

	if err != nil {
		return nil, parseServiceError(err)
	}

	xml.Unmarshal(response, &serverDetails)

	return &serverDetails, nil
}

/**
Deletes the specified server
*/
func (s *Service) DeleteServer(r *request.DeleteServerRequest) error {
	err := s.client.PerformDeleteRequest(s.client.CreateRequestUrl(r.RequestURL()))

	if err != nil {
		return parseServiceError(err)
	}

	return nil
}

/**
Returns all available storages
*/
func (s *Service) GetStorages(r *request.GetStoragesRequest) (*upcloud.Storages, error) {
	storages := upcloud.Storages{}
	response, err := s.basicGetRequest(r.RequestURL())

	if err != nil {
		return nil, err
	}

	xml.Unmarshal(response, &storages)

	return &storages, nil
}

/**
Returns extended details about the specified piece of storage
*/
func (s *Service) GetStorageDetails(r *request.GetStorageDetailsRequest) (*upcloud.StorageDetails, error) {
	storageDetails := upcloud.StorageDetails{}
	response, err := s.basicGetRequest(r.RequestURL())

	if err != nil {
		return nil, err
	}

	xml.Unmarshal(response, &storageDetails)

	return &storageDetails, nil
}

/**
Creates the specified storage
*/
func (s *Service) CreateStorage(r *request.CreateStorageRequest) (*upcloud.StorageDetails, error) {
	storageDetails := upcloud.StorageDetails{}
	requestBody, _ := xml.Marshal(r)
	response, err := s.client.PerformPostRequest(s.client.CreateRequestUrl(r.RequestURL()), requestBody)

	if err != nil {
		return nil, parseServiceError(err)
	}

	xml.Unmarshal(response, &storageDetails)

	return &storageDetails, nil
}

/**
Modifies the specified storage device
*/
func (s *Service) ModifyStorage(r *request.ModifyStorageRequest) (*upcloud.StorageDetails, error) {
	storageDetails := upcloud.StorageDetails{}
	requestBody, _ := xml.Marshal(r)
	response, err := s.client.PerformPutRequest(s.client.CreateRequestUrl(r.RequestURL()), requestBody)

	if err != nil {
		return nil, parseServiceError(err)
	}

	xml.Unmarshal(response, &storageDetails)

	return &storageDetails, nil
}

/**
Attaches the specified storage to the specified server
*/
func (s *Service) AttachStorageRequest(r *request.AttachStorageRequest) (*upcloud.ServerDetails, error) {
	serverDetails := upcloud.ServerDetails{}
	requestBody, _ := xml.Marshal(r)
	response, err := s.client.PerformPostRequest(s.client.CreateRequestUrl(r.RequestURL()), requestBody)

	if err != nil {
		return nil, parseServiceError(err)
	}

	xml.Unmarshal(response, &serverDetails)

	return &serverDetails, nil
}

/**
Detaches the specified storage from the specified server
*/
func (s *Service) DetachStorageRequest(r *request.DetachStorageRequest) (*upcloud.ServerDetails, error) {
	serverDetails := upcloud.ServerDetails{}
	requestBody, _ := xml.Marshal(r)
	response, err := s.client.PerformPostRequest(s.client.CreateRequestUrl(r.RequestURL()), requestBody)

	if err != nil {
		return nil, parseServiceError(err)
	}

	xml.Unmarshal(response, &serverDetails)

	return &serverDetails, nil
}

/**
Deletes the specified storage device
*/
func (s *Service) DeleteStorage(r *request.DeleteStorageRequest) error {
	err := s.client.PerformDeleteRequest(s.client.CreateRequestUrl(r.RequestURL()))

	if err != nil {
		return parseServiceError(err)
	}

	return nil
}

/**
Returns all IP addresses associated with the account
*/
func (s *Service) GetIPAddresses() (*upcloud.IPAddresses, error) {
	ipAddresses := upcloud.IPAddresses{}
	response, err := s.basicGetRequest("/ip_address")

	if err != nil {
		return nil, err
	}

	xml.Unmarshal(response, &ipAddresses)

	return &ipAddresses, nil
}

/**
Returns extended details about the specified IP address
*/
func (s *Service) GetIPAddressDetails(r *request.GetIPAddressDetailsRequest) (*upcloud.IPAddress, error) {
	ipAddress := upcloud.IPAddress{}
	response, err := s.basicGetRequest(r.RequestURL())

	if err != nil {
		return nil, err
	}

	xml.Unmarshal(response, &ipAddress)

	return &ipAddress, nil
}

/**
Returns the firewall rules for the specified server
*/
func (s *Service) GetServerFirewallRules(r *request.GetServerFirewallRulesRequest) (*upcloud.FirewallRules, error) {
	firewallRules := upcloud.FirewallRules{}
	response, err := s.basicGetRequest(r.RequestURL())

	if err != nil {
		return nil, err
	}

	xml.Unmarshal(response, &firewallRules)

	return &firewallRules, nil
}

/**
Returns all tags
*/
func (s *Service) GetTags() (*upcloud.Tags, error) {
	tags := upcloud.Tags{}
	response, err := s.basicGetRequest("/tag")

	if err != nil {
		return nil, err
	}

	xml.Unmarshal(response, &tags)

	return &tags, nil
}

/**
Wrapper that performs a GET request to the specified location and returns the response or a service error
*/
func (s *Service) basicGetRequest(location string) ([]byte, error) {
	requestUrl := s.client.CreateRequestUrl(location)
	response, err := s.client.PerformGetRequest(requestUrl)

	if err != nil {
		return nil, parseServiceError(err)
	}

	return response, nil
}

/**
Parses an error returned from the client into a service error object
*/
func parseServiceError(err error) error {
	// Parse service errors
	if clientError, ok := err.(*client.Error); ok {
		serviceError := upcloud.Error{}
		responseBody := clientError.ResponseBody
		xml.Unmarshal(responseBody, &serviceError)

		return &serviceError
	}

	return err
}
