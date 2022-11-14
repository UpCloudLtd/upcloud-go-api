package service

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/UpCloudLtd/upcloud-go-api/v5/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v5/upcloud/request"
)

type Server interface {
	GetServerConfigurations() (*upcloud.ServerConfigurations, error)
	GetServers() (*upcloud.Servers, error)
	GetServerDetails(r *request.GetServerDetailsRequest) (*upcloud.ServerDetails, error)
	CreateServer(r *request.CreateServerRequest) (*upcloud.ServerDetails, error)
	WaitForServerState(r *request.WaitForServerStateRequest) (*upcloud.ServerDetails, error)
	StartServer(r *request.StartServerRequest) (*upcloud.ServerDetails, error)
	StopServer(r *request.StopServerRequest) (*upcloud.ServerDetails, error)
	RestartServer(r *request.RestartServerRequest) (*upcloud.ServerDetails, error)
	ModifyServer(r *request.ModifyServerRequest) (*upcloud.ServerDetails, error)
	DeleteServer(r *request.DeleteServerRequest) error
	DeleteServerAndStorages(r *request.DeleteServerAndStoragesRequest) error
}

var _ Server = (*Service)(nil)

// GetServerConfigurations returns the available pre-configured server configurations
func (s *Service) GetServerConfigurations() (*upcloud.ServerConfigurations, error) {
	serverConfigurations := upcloud.ServerConfigurations{}
	response, err := s.basicGetRequest("/server_size")
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(response, &serverConfigurations)
	if err != nil {
		return nil, err
	}

	return &serverConfigurations, nil
}

// GetServers returns the available servers
func (s *Service) GetServers() (*upcloud.Servers, error) {
	servers := upcloud.Servers{}
	response, err := s.basicGetRequest("/server")
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(response, &servers)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal JSON: %s, %w", string(response), err)
	}

	return &servers, nil
}

// GetServersWithFilters returns the available servers that match all the given filters
func (s *Service) GetServersWithFilters(r *request.GetServersWithFiltersRequest) (*upcloud.Servers, error) {
	servers := upcloud.Servers{}
	return &servers, s.get(r.RequestURL(), &servers)
}

// GetServerDetails returns extended details about the specified server
func (s *Service) GetServerDetails(r *request.GetServerDetailsRequest) (*upcloud.ServerDetails, error) {
	serverDetails := upcloud.ServerDetails{}
	response, err := s.basicGetRequest(r.RequestURL())
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(response, &serverDetails)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal JSON: %s, %w", string(response), err)
	}

	return &serverDetails, nil
}

// CreateServer creates a server and returns the server details for the newly created server
func (s *Service) CreateServer(r *request.CreateServerRequest) (*upcloud.ServerDetails, error) {
	serverDetails := upcloud.ServerDetails{}
	requestBody, _ := json.Marshal(r)
	response, err := s.client.PerformJSONPostRequest(s.client.CreateRequestURL(r.RequestURL()), requestBody)
	if err != nil {
		return nil, parseJSONServiceError(err)
	}

	err = json.Unmarshal(response, &serverDetails)
	if err != nil {
		return nil, err
	}

	return &serverDetails, nil
}

// WaitForServerState blocks execution until the specified server has entered the specified state. If the state changes
// favorably, the new server details are returned. The method will give up after the specified timeout
func (s *Service) WaitForServerState(r *request.WaitForServerStateRequest) (*upcloud.ServerDetails, error) {
	attempts := 0
	sleepDuration := time.Second * 5

	for {
		// Always wait for one attempt period before querying the state the first time. Newly created servers
		// may not immediately switch to "maintenance" upon creation, triggering a false positive from this
		// method
		attempts++
		time.Sleep(sleepDuration)

		serverDetails, err := s.GetServerDetails(&request.GetServerDetailsRequest{
			UUID: r.UUID,
		})
		if err != nil {
			return nil, err
		}

		// Either wait for the server to enter the desired state or wait for it to leave the undesired state
		if r.DesiredState != "" && serverDetails.State == r.DesiredState {
			return serverDetails, nil
		} else if r.UndesiredState != "" && serverDetails.State != r.UndesiredState {
			return serverDetails, nil
		}

		if time.Duration(attempts)*sleepDuration >= r.Timeout {
			return nil, fmt.Errorf("timeout reached while waiting for server to enter state \"%s\"", r.DesiredState)
		}
	}
}

// StartServer starts the specified server
func (s *Service) StartServer(r *request.StartServerRequest) (*upcloud.ServerDetails, error) {
	// Save previous timeout
	prevTimeout := s.client.GetTimeout()

	// Increase the client timeout to match the request timeout
	s.client.SetTimeout(r.Timeout)

	requestBody, _ := json.Marshal(r)
	response, err := s.client.PerformJSONPostRequest(s.client.CreateRequestURL(r.RequestURL()), requestBody)

	// Restore previous timout
	s.client.SetTimeout(prevTimeout)

	if err != nil {
		return nil, parseJSONServiceError(err)
	}

	serverDetails := upcloud.ServerDetails{}
	err = json.Unmarshal(response, &serverDetails)
	if err != nil {
		return nil, err
	}

	return &serverDetails, nil
}

// StopServer stops the specified server
func (s *Service) StopServer(r *request.StopServerRequest) (*upcloud.ServerDetails, error) {
	// Save previous timeout
	prevTimeout := s.client.GetTimeout()

	// Increase the client timeout to match the request timeout
	// Allow ten seconds to give the API a chance to respond with an error
	s.client.SetTimeout(r.Timeout + 10*time.Second)

	serverDetails := upcloud.ServerDetails{}
	requestBody, _ := json.Marshal(r)
	response, err := s.client.PerformJSONPostRequest(s.client.CreateRequestURL(r.RequestURL()), requestBody)

	// Restore previous timeout
	s.client.SetTimeout(prevTimeout)

	if err != nil {
		return nil, parseJSONServiceError(err)
	}

	err = json.Unmarshal(response, &serverDetails)
	if err != nil {
		return nil, err
	}

	return &serverDetails, nil
}

// RestartServer restarts the specified server
func (s *Service) RestartServer(r *request.RestartServerRequest) (*upcloud.ServerDetails, error) {
	// Save previous timeout
	prevTimeout := s.client.GetTimeout()

	// Increase the client timeout to match the request timeout
	// Allow ten seconds to give the API a chance to respond with an error
	s.client.SetTimeout(r.Timeout + 10*time.Second)

	serverDetails := upcloud.ServerDetails{}
	requestBody, _ := json.Marshal(r)
	response, err := s.client.PerformJSONPostRequest(s.client.CreateRequestURL(r.RequestURL()), requestBody)

	// Restore previous timeout
	s.client.SetTimeout(prevTimeout)

	if err != nil {
		return nil, parseJSONServiceError(err)
	}

	err = json.Unmarshal(response, &serverDetails)
	if err != nil {
		return nil, err
	}

	return &serverDetails, nil
}

// ModifyServer modifies the configuration of an existing server. Attaching and detaching storages as well as assigning
// and releasing IP addresses have their own separate operations.
func (s *Service) ModifyServer(r *request.ModifyServerRequest) (*upcloud.ServerDetails, error) {
	serverDetails := upcloud.ServerDetails{}
	requestBody, _ := json.Marshal(r)
	response, err := s.client.PerformJSONPutRequest(s.client.CreateRequestURL(r.RequestURL()), requestBody)
	if err != nil {
		return nil, parseJSONServiceError(err)
	}

	err = json.Unmarshal(response, &serverDetails)
	if err != nil {
		return nil, err
	}

	return &serverDetails, nil
}

// DeleteServer deletes the specified server
func (s *Service) DeleteServer(r *request.DeleteServerRequest) error {
	err := s.client.PerformJSONDeleteRequest(s.client.CreateRequestURL(r.RequestURL()))
	if err != nil {
		return parseJSONServiceError(err)
	}

	return nil
}

// DeleteServerAndStorages deletes the specified server and all attached storages
func (s *Service) DeleteServerAndStorages(r *request.DeleteServerAndStoragesRequest) error {
	err := s.client.PerformJSONDeleteRequest(s.client.CreateRequestURL(r.RequestURL()))
	if err != nil {
		return parseJSONServiceError(err)
	}

	return nil
}
