package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud/request"
)

type ManagedDatabaseServiceManager interface {
	CancelManagedDatabaseConnection(r *request.CancelManagedDatabaseConnection) error
	CloneManagedDatabase(r *request.CloneManagedDatabaseRequest) (*upcloud.ManagedDatabase, error)
	CreateManagedDatabase(r *request.CreateManagedDatabaseRequest) (*upcloud.ManagedDatabase, error)
	GetManagedDatabase(r *request.GetManagedDatabaseRequest) (*upcloud.ManagedDatabase, error)
	GetManagedDatabases(r *request.GetManagedDatabasesRequest) ([]upcloud.ManagedDatabase, error)
	GetManagedDatabaseConnections(r *request.GetManagedDatabaseConnectionsRequest) ([]upcloud.ManagedDatabaseConnection, error)
	GetManagedDatabaseMetrics(r *request.GetManagedDatabaseMetricsRequest) (*upcloud.ManagedDatabaseMetrics, error)
	GetManagedDatabaseLogs(r *request.GetManagedDatabaseLogsRequest) (*upcloud.ManagedDatabaseLogs, error)
	GetManagedDatabaseQueryStatisticsMySQL(r *request.GetManagedDatabaseQueryStatisticsRequest) ([]upcloud.ManagedDatabaseQueryStatisticsMySQL, error)
	GetManagedDatabaseQueryStatisticsPostgreSQL(r *request.GetManagedDatabaseQueryStatisticsRequest) ([]upcloud.ManagedDatabaseQueryStatisticsPostgreSQL, error)
	GetManagedDatabaseServiceType(r *request.GetManagedDatabaseServiceTypeRequest) (*upcloud.ManagedDatabaseType, error)
	GetManagedDatabaseServiceTypes(r *request.GetManagedDatabaseServiceTypesRequest) (map[string]upcloud.ManagedDatabaseType, error)
	DeleteManagedDatabase(r *request.DeleteManagedDatabaseRequest) error
	ModifyManagedDatabase(r *request.ModifyManagedDatabaseRequest) (*upcloud.ManagedDatabase, error)
	UpgradeManagedDatabaseVersion(r *request.UpgradeManagedDatabaseVersionRequest) (*upcloud.ManagedDatabase, error)
	GetManagedDatabaseVersions(r *request.GetManagedDatabaseVersionsRequest) ([]string, error)
	StartManagedDatabase(r *request.StartManagedDatabaseRequest) (*upcloud.ManagedDatabase, error)
	ShutdownManagedDatabase(r *request.ShutdownManagedDatabaseRequest) (*upcloud.ManagedDatabase, error)
	WaitForManagedDatabaseState(r *request.WaitForManagedDatabaseStateRequest) (*upcloud.ManagedDatabase, error)
}

type ManagedDatabaseUserManager interface {
	CreateManagedDatabaseUser(r *request.CreateManagedDatabaseUserRequest) (*upcloud.ManagedDatabaseUser, error)
	GetManagedDatabaseUser(r *request.GetManagedDatabaseUserRequest) (*upcloud.ManagedDatabaseUser, error)
	GetManagedDatabaseUsers(r *request.GetManagedDatabaseUsersRequest) ([]upcloud.ManagedDatabaseUser, error)
	DeleteManagedDatabaseUser(r *request.DeleteManagedDatabaseUserRequest) error
	ModifyManagedDatabaseUser(r *request.ModifyManagedDatabaseUserRequest) (*upcloud.ManagedDatabaseUser, error)
}

type ManagedDatabaseLogicalDatabaseManager interface {
	CreateManagedDatabaseLogicalDatabase(r *request.CreateManagedDatabaseLogicalDatabaseRequest) (*upcloud.ManagedDatabaseLogicalDatabase, error)
	GetManagedDatabaseLogicalDatabases(r *request.GetManagedDatabaseLogicalDatabasesRequest) ([]upcloud.ManagedDatabaseLogicalDatabase, error)
	DeleteManagedDatabaseLogicalDatabase(r *request.DeleteManagedDatabaseLogicalDatabaseRequest) error
}

var (
	_ ManagedDatabaseServiceManager         = (*Service)(nil)
	_ ManagedDatabaseUserManager            = (*Service)(nil)
	_ ManagedDatabaseLogicalDatabaseManager = (*Service)(nil)
)

var ErrCancelManagedDatabaseConnection = errors.New("managed database connection cancellation failed")

/* Service Management */

// CancelManagedDatabaseConnection (EXPERIMENTAL) cancels a current query of a database connection or terminates it entirely.
// In case of the server is unable to cancel the query or terminate the connection ErrCancelManagedDatabaseConnection
// is returned.
func (s *Service) CancelManagedDatabaseConnection(r *request.CancelManagedDatabaseConnection) error {
	res := struct {
		Success bool `json:"success"`
	}{}
	response, err := s.client.PerformJSONDeleteRequestWithResponseBody(s.client.CreateRequestURL(r.RequestURL()))
	if err != nil {
		return parseJSONServiceError(err)
	}

	err = json.Unmarshal(response, &res)
	if err != nil {
		return fmt.Errorf("unable to unmarshal JSON: %w", err)
	}
	if !res.Success {
		return ErrCancelManagedDatabaseConnection
	}
	return nil
}

// CloneManagedDatabase (EXPERIMENTAL) clones en existing managed database instance
func (s *Service) CloneManagedDatabase(r *request.CloneManagedDatabaseRequest) (*upcloud.ManagedDatabase, error) {
	managedDatabaseDetails := upcloud.ManagedDatabase{}
	requestBody, _ := json.Marshal(r)
	response, err := s.client.PerformJSONPostRequest(s.client.CreateRequestURL(r.RequestURL()), requestBody)
	if err != nil {
		return nil, parseJSONServiceError(err)
	}

	err = json.Unmarshal(response, &managedDatabaseDetails)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal JSON: %w", err)
	}

	return &managedDatabaseDetails, nil
}

// CreateManagedDatabase (EXPERIMENTAL) creates a new managed database instance
func (s *Service) CreateManagedDatabase(r *request.CreateManagedDatabaseRequest) (*upcloud.ManagedDatabase, error) {
	managedDatabaseDetails := upcloud.ManagedDatabase{}
	requestBody, _ := json.Marshal(r)
	response, err := s.client.PerformJSONPostRequest(s.client.CreateRequestURL(r.RequestURL()), requestBody)
	if err != nil {
		return nil, parseJSONServiceError(err)
	}

	err = json.Unmarshal(response, &managedDatabaseDetails)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal JSON: %w", err)
	}

	return &managedDatabaseDetails, nil
}

// GetManagedDatabase (EXPERIMENTAL) gets details of an existing managed database instance
func (s *Service) GetManagedDatabase(r *request.GetManagedDatabaseRequest) (*upcloud.ManagedDatabase, error) {
	managedDatabaseDetails := upcloud.ManagedDatabase{}
	response, err := s.client.PerformJSONGetRequest(s.client.CreateRequestURL(r.RequestURL()))
	if err != nil {
		return nil, parseJSONServiceError(err)
	}

	err = json.Unmarshal(response, &managedDatabaseDetails)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal JSON: %w", err)
	}

	return &managedDatabaseDetails, nil
}

// GetManagedDatabases (EXPERIMENTAL) returns a slice of all managed database instances within an account
func (s *Service) GetManagedDatabases(r *request.GetManagedDatabasesRequest) ([]upcloud.ManagedDatabase, error) {
	var services []upcloud.ManagedDatabase
	response, err := s.client.PerformJSONGetRequest(s.client.CreateRequestURL(r.RequestURL()))
	if err != nil {
		return nil, parseJSONServiceError(err)
	}

	err = json.Unmarshal(response, &services)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal JSON: %w", err)
	}

	return services, nil
}

// GetManagedDatabaseConnections (EXPERIMENTAL) returns a slice of connections from an existing managed database instance
func (s *Service) GetManagedDatabaseConnections(r *request.GetManagedDatabaseConnectionsRequest) ([]upcloud.ManagedDatabaseConnection, error) {
	var conns []upcloud.ManagedDatabaseConnection
	response, err := s.client.PerformJSONGetRequest(s.client.CreateRequestURL(r.RequestURL()))
	if err != nil {
		return nil, parseJSONServiceError(err)
	}

	err = json.Unmarshal(response, &conns)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal JSON: %w", err)
	}

	return conns, nil
}

// GetManagedDatabaseMetrics (EXPERIMENTAL) returns metrics collection for the selected period
func (s *Service) GetManagedDatabaseMetrics(r *request.GetManagedDatabaseMetricsRequest) (*upcloud.ManagedDatabaseMetrics, error) {
	metrics := upcloud.ManagedDatabaseMetrics{}
	response, err := s.client.PerformJSONGetRequest(s.client.CreateRequestURL(r.RequestURL()))
	if err != nil {
		return nil, parseJSONServiceError(err)
	}

	err = json.Unmarshal(response, &metrics)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal JSON: %w", err)
	}

	return &metrics, nil
}

// GetManagedDatabaseLogs (EXPERIMENTAL) returns logs of a managed database instance
func (s *Service) GetManagedDatabaseLogs(r *request.GetManagedDatabaseLogsRequest) (*upcloud.ManagedDatabaseLogs, error) {
	logs := upcloud.ManagedDatabaseLogs{}
	response, err := s.client.PerformJSONGetRequest(s.client.CreateRequestURL(r.RequestURL()))
	if err != nil {
		return nil, parseJSONServiceError(err)
	}

	err = json.Unmarshal(response, &logs)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal JSON: %w", err)
	}

	return &logs, nil
}

// GetManagedDatabaseServiceType (EXPERIMENTAL) returns details of requested service type
func (s *Service) GetManagedDatabaseServiceType(r *request.GetManagedDatabaseServiceTypeRequest) (*upcloud.ManagedDatabaseType, error) {
	var serviceType upcloud.ManagedDatabaseType
	return &serviceType, s.get(r.RequestURL(), &serviceType)
}

// GetManagedDatabaseServiceTypes (EXPERIMENTAL) returns a map of available database service types
func (s *Service) GetManagedDatabaseServiceTypes(r *request.GetManagedDatabaseServiceTypesRequest) (map[string]upcloud.ManagedDatabaseType, error) {
	serviceTypes := make(map[string]upcloud.ManagedDatabaseType)
	return serviceTypes, s.get(r.RequestURL(), &serviceTypes)
}

// GetManagedDatabaseQueryStatisticsMySQL (EXPERIMENTAL) returns MySQL query statistics of a managed database instance
func (s *Service) GetManagedDatabaseQueryStatisticsMySQL(r *request.GetManagedDatabaseQueryStatisticsRequest) ([]upcloud.ManagedDatabaseQueryStatisticsMySQL, error) {
	var parsed struct {
		Mysql []upcloud.ManagedDatabaseQueryStatisticsMySQL
	}
	response, err := s.client.PerformJSONGetRequest(s.client.CreateRequestURL(r.RequestURL()))
	if err != nil {
		return nil, parseJSONServiceError(err)
	}

	err = json.Unmarshal(response, &parsed)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal JSON: %w", err)
	}

	return parsed.Mysql, nil
}

// GetManagedDatabaseQueryStatisticsPostgres (EXPERIMENTAL) returns PostgreSQL query statistics of a managed database instance
func (s *Service) GetManagedDatabaseQueryStatisticsPostgreSQL(r *request.GetManagedDatabaseQueryStatisticsRequest) ([]upcloud.ManagedDatabaseQueryStatisticsPostgreSQL, error) {
	var parsed struct {
		Pg []upcloud.ManagedDatabaseQueryStatisticsPostgreSQL
	}
	response, err := s.client.PerformJSONGetRequest(s.client.CreateRequestURL(r.RequestURL()))
	if err != nil {
		return nil, parseJSONServiceError(err)
	}

	err = json.Unmarshal(response, &parsed)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal JSON: %w", err)
	}

	return parsed.Pg, nil
}

// DeleteManagedDatabase (EXPERIMENTAL) deletes an existing managed database instance
func (s *Service) DeleteManagedDatabase(r *request.DeleteManagedDatabaseRequest) error {
	err := s.client.PerformJSONDeleteRequest(s.client.CreateRequestURL(r.RequestURL()))
	if err != nil {
		return parseJSONServiceError(err)
	}

	return nil
}

// ModifyManagedDatabase (EXPERIMENTAL) modifies an existing managed database instance
func (s *Service) ModifyManagedDatabase(r *request.ModifyManagedDatabaseRequest) (*upcloud.ManagedDatabase, error) {
	managedDatabaseDetails := upcloud.ManagedDatabase{}
	requestBody, _ := json.Marshal(r)
	response, err := s.client.PerformJSONPatchRequest(s.client.CreateRequestURL(r.RequestURL()), requestBody)
	if err != nil {
		return nil, parseJSONServiceError(err)
	}

	err = json.Unmarshal(response, &managedDatabaseDetails)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal JSON: %w", err)
	}

	return &managedDatabaseDetails, nil
}

// UpgradeManagedDatabaseServiceVersion upgrades the version of the database service;
// for the list of available versions use GetManagedDatabaseVersions function
func (s *Service) UpgradeManagedDatabaseVersion(r *request.UpgradeManagedDatabaseVersionRequest) (*upcloud.ManagedDatabase, error) {
	managedDatabaseDetails := upcloud.ManagedDatabase{}

	reqBody, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	response, err := s.client.PerformJSONPostRequest(s.client.CreateRequestURL(r.RequestURL()), reqBody)
	if err != nil {
		return nil, parseJSONServiceError(err)
	}

	err = json.Unmarshal(response, &managedDatabaseDetails)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal JSON: %w", err)
	}

	return &managedDatabaseDetails, nil
}

// GetManagedDatabaseVersions available versions of the specific Managed Database service
func (s *Service) GetManagedDatabaseVersions(r *request.GetManagedDatabaseVersionsRequest) ([]string, error) {
	versions := []string{}

	res, err := s.basicGetRequest(r.RequestURL())
	if err != nil {
		return nil, parseJSONServiceError(err)
	}

	err = json.Unmarshal(res, &versions)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal JSON: %w", err)
	}

	return versions, nil
}

// WaitForManagedDatabaseState (EXPERIMENTAL) blocks execution until the specified managed database instance has entered the
// specified state. If the state changes favorably, the new managed database details is returned. The method will give up
// after the specified timeout
func (s *Service) WaitForManagedDatabaseState(r *request.WaitForManagedDatabaseStateRequest) (*upcloud.ManagedDatabase, error) {
	attempts := 0
	sleepDuration := time.Second * 5

	for {
		attempts++

		details, err := s.GetManagedDatabase(&request.GetManagedDatabaseRequest{
			UUID: r.UUID,
		})
		if err != nil {
			return nil, err
		}

		if details.State == r.DesiredState {
			return details, nil
		}

		time.Sleep(sleepDuration)

		if time.Duration(attempts)*sleepDuration >= r.Timeout {
			return nil, fmt.Errorf("timeout reached while waiting for managed database instance to enter state \"%s\"", r.DesiredState)
		}
	}
}

// StartManagedDatabase (EXPERIMENTAL) starts a shut down existing managed database instance
func (s *Service) StartManagedDatabase(r *request.StartManagedDatabaseRequest) (*upcloud.ManagedDatabase, error) {
	managedDatabaseDetails := upcloud.ManagedDatabase{}
	requestBody, _ := json.Marshal(r)
	response, err := s.client.PerformJSONPatchRequest(s.client.CreateRequestURL(r.RequestURL()), requestBody)
	if err != nil {
		return nil, parseJSONServiceError(err)
	}

	err = json.Unmarshal(response, &managedDatabaseDetails)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal JSON: %w", err)
	}

	return &managedDatabaseDetails, nil
}

// ShutdownManagedDatabase (EXPERIMENTAL) shuts down existing managed database instance. Only a service which has at least one
// full backup can be shut down.
func (s *Service) ShutdownManagedDatabase(r *request.ShutdownManagedDatabaseRequest) (*upcloud.ManagedDatabase, error) {
	managedDatabaseDetails := upcloud.ManagedDatabase{}
	requestBody, _ := json.Marshal(r)
	response, err := s.client.PerformJSONPatchRequest(s.client.CreateRequestURL(r.RequestURL()), requestBody)
	if err != nil {
		return nil, parseJSONServiceError(err)
	}

	err = json.Unmarshal(response, &managedDatabaseDetails)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal JSON: %w", err)
	}

	return &managedDatabaseDetails, nil
}

/* User Management */

// CreateManagedDatabaseUser (EXPERIMENTAL) creates a new normal user to an existing managed database instance
func (s *Service) CreateManagedDatabaseUser(r *request.CreateManagedDatabaseUserRequest) (*upcloud.ManagedDatabaseUser, error) {
	userDetails := upcloud.ManagedDatabaseUser{}
	requestBody, _ := json.Marshal(r)
	response, err := s.client.PerformJSONPostRequest(s.client.CreateRequestURL(r.RequestURL()), requestBody)
	if err != nil {
		return nil, parseJSONServiceError(err)
	}

	err = json.Unmarshal(response, &userDetails)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal JSON: %w", err)
	}

	return &userDetails, nil
}

// GetManagedDatabaseUser (EXPERIMENTAL) returns details of an existing user of an existing managed database instance
func (s *Service) GetManagedDatabaseUser(r *request.GetManagedDatabaseUserRequest) (*upcloud.ManagedDatabaseUser, error) {
	userDetails := upcloud.ManagedDatabaseUser{}
	response, err := s.client.PerformJSONGetRequest(s.client.CreateRequestURL(r.RequestURL()))
	if err != nil {
		return nil, parseJSONServiceError(err)
	}

	err = json.Unmarshal(response, &userDetails)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal JSON: %w", err)
	}

	return &userDetails, nil
}

// GetManagedDatabaseUsers (EXPERIMENTAL) returns a slice of all users of an existing managed database instance
func (s *Service) GetManagedDatabaseUsers(r *request.GetManagedDatabaseUsersRequest) ([]upcloud.ManagedDatabaseUser, error) {
	var userList []upcloud.ManagedDatabaseUser
	response, err := s.client.PerformJSONGetRequest(s.client.CreateRequestURL(r.RequestURL()))
	if err != nil {
		return nil, parseJSONServiceError(err)
	}

	err = json.Unmarshal(response, &userList)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal JSON: %w", err)
	}

	return userList, nil
}

// DeleteManagedDatabaseUser (EXPERIMENTAL) deletes an existing user of an existing managed database instance
func (s *Service) DeleteManagedDatabaseUser(r *request.DeleteManagedDatabaseUserRequest) error {
	err := s.client.PerformJSONDeleteRequest(s.client.CreateRequestURL(r.RequestURL()))
	if err != nil {
		return parseJSONServiceError(err)
	}

	return nil
}

// ModifyManagedDatabaseUser (EXPERIMENTAL) modifies an existing user of an existing managed database instance
func (s *Service) ModifyManagedDatabaseUser(r *request.ModifyManagedDatabaseUserRequest) (*upcloud.ManagedDatabaseUser, error) {
	userDetails := upcloud.ManagedDatabaseUser{}
	requestBody, _ := json.Marshal(r)
	response, err := s.client.PerformJSONPatchRequest(s.client.CreateRequestURL(r.RequestURL()), requestBody)
	if err != nil {
		return nil, parseJSONServiceError(err)
	}

	err = json.Unmarshal(response, &userDetails)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal JSON: %w", err)
	}

	return &userDetails, nil
}

/* Logical Database Management */

// CreateManagedDatabaseLogicalDatabase (EXPERIMENTAL) creates a new logical database to an existing managed database instance
func (s *Service) CreateManagedDatabaseLogicalDatabase(r *request.CreateManagedDatabaseLogicalDatabaseRequest) (*upcloud.ManagedDatabaseLogicalDatabase, error) {
	dbDetails := upcloud.ManagedDatabaseLogicalDatabase{}
	requestBody, _ := json.Marshal(r)
	response, err := s.client.PerformJSONPostRequest(s.client.CreateRequestURL(r.RequestURL()), requestBody)
	if err != nil {
		return nil, parseJSONServiceError(err)
	}

	err = json.Unmarshal(response, &dbDetails)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal JSON: %w", err)
	}

	return &dbDetails, nil
}

// GetManagedDatabaseLogicalDatabases (EXPERIMENTAL) returns a slice of all logical databases of an existing managed database instance
func (s *Service) GetManagedDatabaseLogicalDatabases(r *request.GetManagedDatabaseLogicalDatabasesRequest) ([]upcloud.ManagedDatabaseLogicalDatabase, error) {
	var dbList []upcloud.ManagedDatabaseLogicalDatabase
	response, err := s.client.PerformJSONGetRequest(s.client.CreateRequestURL(r.RequestURL()))
	if err != nil {
		return nil, parseJSONServiceError(err)
	}

	err = json.Unmarshal(response, &dbList)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal JSON: %w", err)
	}

	return dbList, nil
}

// DeleteManagedDatabaseLogicalDatabase (EXPERIMENTAL) deletes an existing logical database of an existing managed database instance
func (s *Service) DeleteManagedDatabaseLogicalDatabase(r *request.DeleteManagedDatabaseLogicalDatabaseRequest) error {
	err := s.client.PerformJSONDeleteRequest(s.client.CreateRequestURL(r.RequestURL()))
	if err != nil {
		return parseJSONServiceError(err)
	}

	return nil
}
