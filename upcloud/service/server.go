package service

import (
	"context"

	"github.com/UpCloudLtd/upcloud-go-api/v6/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v6/upcloud/request"
)

type Server interface {
	GetServerConfigurations(ctx context.Context) (*upcloud.ServerConfigurations, error)
	GetServers(ctx context.Context) (*upcloud.Servers, error)
	GetServerDetails(ctx context.Context, r *request.GetServerDetailsRequest) (*upcloud.ServerDetails, error)
	CreateServer(ctx context.Context, r *request.CreateServerRequest) (*upcloud.ServerDetails, error)
	WaitForServerState(ctx context.Context, r *request.WaitForServerStateRequest) (*upcloud.ServerDetails, error)
	StartServer(ctx context.Context, r *request.StartServerRequest) (*upcloud.ServerDetails, error)
	StopServer(ctx context.Context, r *request.StopServerRequest) (*upcloud.ServerDetails, error)
	RestartServer(ctx context.Context, r *request.RestartServerRequest) (*upcloud.ServerDetails, error)
	ModifyServer(ctx context.Context, r *request.ModifyServerRequest) (*upcloud.ServerDetails, error)
	DeleteServer(ctx context.Context, r *request.DeleteServerRequest) error
	DeleteServerAndStorages(ctx context.Context, r *request.DeleteServerAndStoragesRequest) error
}

// GetServerConfigurations returns the available pre-configured server configurations
func (s *Service) GetServerConfigurations(ctx context.Context) (*upcloud.ServerConfigurations, error) {
	serverConfigurations := upcloud.ServerConfigurations{}
	return &serverConfigurations, s.get(ctx, "/server_size", &serverConfigurations)
}

// GetServers returns the available servers
func (s *Service) GetServers(ctx context.Context) (*upcloud.Servers, error) {
	servers := upcloud.Servers{}
	return &servers, s.get(ctx, "/server", &servers)
}

// GetServersWithFilters returns the all the available servers using given filters.
func (s *Service) GetServersWithFilters(ctx context.Context, r *request.GetServersWithFiltersRequest) (*upcloud.Servers, error) {
	servers := upcloud.Servers{}
	return &servers, s.get(ctx, r.RequestURL(), &servers)
}

// GetServerDetails returns extended details about the specified server
func (s *Service) GetServerDetails(ctx context.Context, r *request.GetServerDetailsRequest) (*upcloud.ServerDetails, error) {
	serverDetails := upcloud.ServerDetails{}
	return &serverDetails, s.get(ctx, r.RequestURL(), &serverDetails)
}

// CreateServer creates a server and returns the server details for the newly created server
func (s *Service) CreateServer(ctx context.Context, r *request.CreateServerRequest) (*upcloud.ServerDetails, error) {
	serverDetails := upcloud.ServerDetails{}
	return &serverDetails, s.create(ctx, r, &serverDetails)
}

// WaitForServerState blocks execution until the specified server has entered the specified state. If the state changes
// favorably, the new server details are returned. The method will give up after the specified timeout
func (s *Service) WaitForServerState(ctx context.Context, r *request.WaitForServerStateRequest) (*upcloud.ServerDetails, error) {
	return retry(ctx, func(i int, c context.Context) (*upcloud.ServerDetails, error) {
		details, err := s.GetServerDetails(c, &request.GetServerDetailsRequest{
			UUID: r.UUID,
		})
		if err != nil {
			return nil, err
		}

		// Either wait for the server to enter the desired state or wait for it to leave the undesired state
		if r.DesiredState != "" && details.State == r.DesiredState {
			return details, nil
		} else if r.UndesiredState != "" && details.State != r.UndesiredState {
			return details, nil
		}

		return nil, nil
	}, nil)
}

// StartServer starts the specified server
func (s *Service) StartServer(ctx context.Context, r *request.StartServerRequest) (*upcloud.ServerDetails, error) {
	serverDetails := upcloud.ServerDetails{}
	return &serverDetails, s.create(ctx, r, &serverDetails)
}

// StopServer stops the specified server
func (s *Service) StopServer(ctx context.Context, r *request.StopServerRequest) (*upcloud.ServerDetails, error) {
	serverDetails := upcloud.ServerDetails{}
	if r.Timeout > 0 {
		timeoutCtx, cancel := context.WithTimeout(ctx, r.Timeout)
		defer cancel()
		return &serverDetails, s.create(timeoutCtx, r, &serverDetails)
	}
	return &serverDetails, s.create(ctx, r, &serverDetails)
}

// RestartServer restarts the specified server
func (s *Service) RestartServer(ctx context.Context, r *request.RestartServerRequest) (*upcloud.ServerDetails, error) {
	serverDetails := upcloud.ServerDetails{}
	if r.Timeout > 0 {
		timeoutCtx, cancel := context.WithTimeout(ctx, r.Timeout)
		defer cancel()
		return &serverDetails, s.create(timeoutCtx, r, &serverDetails)
	}
	return &serverDetails, s.create(ctx, r, &serverDetails)
}

// ModifyServer modifies the configuration of an existing server. Attaching and detaching storages as well as assigning
// and releasing IP addresses have their own separate operations.
func (s *Service) ModifyServer(ctx context.Context, r *request.ModifyServerRequest) (*upcloud.ServerDetails, error) {
	serverDetails := upcloud.ServerDetails{}
	return &serverDetails, s.replace(ctx, r, &serverDetails)
}

// DeleteServer deletes the specified server
func (s *Service) DeleteServer(ctx context.Context, r *request.DeleteServerRequest) error {
	return s.delete(ctx, r)
}

// DeleteServerAndStorages deletes the specified server and all attached storages
func (s *Service) DeleteServerAndStorages(ctx context.Context, r *request.DeleteServerAndStoragesRequest) error {
	return s.delete(ctx, r)
}
