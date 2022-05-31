package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud/request"
)

type ManagedDatabaseServiceManagerContext interface {
	CancelManagedDatabaseConnection(ctx context.Context, r *request.CancelManagedDatabaseConnection) error
	CloneManagedDatabase(ctx context.Context, r *request.CloneManagedDatabaseRequest) (*upcloud.ManagedDatabase, error)
	CreateManagedDatabase(ctx context.Context, r *request.CreateManagedDatabaseRequest) (*upcloud.ManagedDatabase, error)
	GetManagedDatabase(ctx context.Context, r *request.GetManagedDatabaseRequest) (*upcloud.ManagedDatabase, error)
	GetManagedDatabases(ctx context.Context, r *request.GetManagedDatabasesRequest) ([]upcloud.ManagedDatabase, error)
	GetManagedDatabaseConnections(ctx context.Context, r *request.GetManagedDatabaseConnectionsRequest) ([]upcloud.ManagedDatabaseConnection, error)
	GetManagedDatabaseMetrics(ctx context.Context, r *request.GetManagedDatabaseMetricsRequest) (*upcloud.ManagedDatabaseMetrics, error)
	GetManagedDatabaseLogs(ctx context.Context, r *request.GetManagedDatabaseLogsRequest) (*upcloud.ManagedDatabaseLogs, error)
	GetManagedDatabaseQueryStatisticsMySQL(ctx context.Context, r *request.GetManagedDatabaseQueryStatisticsRequest) ([]upcloud.ManagedDatabaseQueryStatisticsMySQL, error)
	GetManagedDatabaseQueryStatisticsPostgreSQL(ctx context.Context, r *request.GetManagedDatabaseQueryStatisticsRequest) ([]upcloud.ManagedDatabaseQueryStatisticsPostgreSQL, error)
	DeleteManagedDatabase(ctx context.Context, r *request.DeleteManagedDatabaseRequest) error
	ModifyManagedDatabase(ctx context.Context, r *request.ModifyManagedDatabaseRequest) (*upcloud.ManagedDatabase, error)
	UpgradeManagedDatabaseVersion(ctx context.Context, r *request.UpgradeManagedDatabaseVersionRequest) (*upcloud.ManagedDatabase, error)
	GetManagedDatabaseVersions(ctx context.Context, r *request.GetManagedDatabaseVersionsRequest) ([]string, error)
	StartManagedDatabase(ctx context.Context, r *request.StartManagedDatabaseRequest) (*upcloud.ManagedDatabase, error)
	ShutdownManagedDatabase(ctx context.Context, r *request.ShutdownManagedDatabaseRequest) (*upcloud.ManagedDatabase, error)
	WaitForManagedDatabaseState(ctx context.Context, r *request.WaitForManagedDatabaseStateRequest) (*upcloud.ManagedDatabase, error)
}

type ManagedDatabaseUserManagerContext interface {
	CreateManagedDatabaseUser(ctx context.Context, r *request.CreateManagedDatabaseUserRequest) (*upcloud.ManagedDatabaseUser, error)
	GetManagedDatabaseUser(ctx context.Context, r *request.GetManagedDatabaseUserRequest) (*upcloud.ManagedDatabaseUser, error)
	GetManagedDatabaseUsers(ctx context.Context, r *request.GetManagedDatabaseUsersRequest) ([]upcloud.ManagedDatabaseUser, error)
	DeleteManagedDatabaseUser(ctx context.Context, r *request.DeleteManagedDatabaseUserRequest) error
	ModifyManagedDatabaseUser(ctx context.Context, r *request.ModifyManagedDatabaseUserRequest) (*upcloud.ManagedDatabaseUser, error)
}

type ManagedDatabaseLogicalDatabaseManagerContext interface {
	CreateManagedDatabaseLogicalDatabase(ctx context.Context, r *request.CreateManagedDatabaseLogicalDatabaseRequest) (*upcloud.ManagedDatabaseLogicalDatabase, error)
	GetManagedDatabaseLogicalDatabases(ctx context.Context, r *request.GetManagedDatabaseLogicalDatabasesRequest) ([]upcloud.ManagedDatabaseLogicalDatabase, error)
	DeleteManagedDatabaseLogicalDatabase(ctx context.Context, r *request.DeleteManagedDatabaseLogicalDatabaseRequest) error
}

/* Service Management */

// CancelManagedDatabaseConnection (EXPERIMENTAL) cancels a current query of a database connection or terminates it entirely.
// In case of the server is unable to cancel the query or terminate the connection ErrCancelManagedDatabaseConnection
// is returned.
func (s *ServiceContext) CancelManagedDatabaseConnection(ctx context.Context, r *request.CancelManagedDatabaseConnection) error {
	res := struct {
		Success bool `json:"success"`
	}{}
	response, err := s.client.PerformJSONDeleteRequestWithResponseBody(ctx, s.client.CreateRequestURL(r.RequestURL()))

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
func (s *ServiceContext) CloneManagedDatabase(ctx context.Context, r *request.CloneManagedDatabaseRequest) (*upcloud.ManagedDatabase, error) {
	managedDatabaseDetails := upcloud.ManagedDatabase{}
	return &managedDatabaseDetails, s.create(ctx, r, &managedDatabaseDetails)
}

// CreateManagedDatabase (EXPERIMENTAL) creates a new managed database instance
func (s *ServiceContext) CreateManagedDatabase(ctx context.Context, r *request.CreateManagedDatabaseRequest) (*upcloud.ManagedDatabase, error) {
	managedDatabaseDetails := upcloud.ManagedDatabase{}
	return &managedDatabaseDetails, s.create(ctx, r, &managedDatabaseDetails)
}

// GetManagedDatabase (EXPERIMENTAL) gets details of an existing managed database instance
func (s *ServiceContext) GetManagedDatabase(ctx context.Context, r *request.GetManagedDatabaseRequest) (*upcloud.ManagedDatabase, error) {
	managedDatabaseDetails := upcloud.ManagedDatabase{}
	return &managedDatabaseDetails, s.get(ctx, r.RequestURL(), &managedDatabaseDetails)
}

// GetManagedDatabases (EXPERIMENTAL) returns a slice of all managed database instances within an account
func (s *ServiceContext) GetManagedDatabases(ctx context.Context, r *request.GetManagedDatabasesRequest) ([]upcloud.ManagedDatabase, error) {
	var services []upcloud.ManagedDatabase
	return services, s.get(ctx, r.RequestURL(), &services)
}

// GetManagedDatabaseConnections (EXPERIMENTAL) returns a slice of connections from an existing managed database instance
func (s *ServiceContext) GetManagedDatabaseConnections(ctx context.Context, r *request.GetManagedDatabaseConnectionsRequest) ([]upcloud.ManagedDatabaseConnection, error) {
	conns := make([]upcloud.ManagedDatabaseConnection, 0)
	return conns, s.get(ctx, r.RequestURL(), &conns)
}

// GetManagedDatabaseMetrics (EXPERIMENTAL) returns metrics collection for the selected period
func (s *ServiceContext) GetManagedDatabaseMetrics(ctx context.Context, r *request.GetManagedDatabaseMetricsRequest) (*upcloud.ManagedDatabaseMetrics, error) {
	metrics := upcloud.ManagedDatabaseMetrics{}
	return &metrics, s.get(ctx, r.RequestURL(), &metrics)
}

// GetManagedDatabaseLogs (EXPERIMENTAL) returns logs of a managed database instance
func (s *ServiceContext) GetManagedDatabaseLogs(ctx context.Context, r *request.GetManagedDatabaseLogsRequest) (*upcloud.ManagedDatabaseLogs, error) {
	logs := upcloud.ManagedDatabaseLogs{}
	return &logs, s.get(ctx, r.RequestURL(), &logs)
}

// GetManagedDatabaseQueryStatisticsMySQL (EXPERIMENTAL) returns MySQL query statistics of a managed database instance
func (s *ServiceContext) GetManagedDatabaseQueryStatisticsMySQL(ctx context.Context, r *request.GetManagedDatabaseQueryStatisticsRequest) ([]upcloud.ManagedDatabaseQueryStatisticsMySQL, error) {
	var parsed struct {
		Mysql []upcloud.ManagedDatabaseQueryStatisticsMySQL
	}
	if err := s.get(ctx, r.RequestURL(), &parsed); err != nil {
		return nil, err
	}
	return parsed.Mysql, nil
}

// GetManagedDatabaseQueryStatisticsPostgres (EXPERIMENTAL) returns PostgreSQL query statistics of a managed database instance
func (s *ServiceContext) GetManagedDatabaseQueryStatisticsPostgreSQL(ctx context.Context, r *request.GetManagedDatabaseQueryStatisticsRequest) ([]upcloud.ManagedDatabaseQueryStatisticsPostgreSQL, error) {
	var parsed struct {
		Pg []upcloud.ManagedDatabaseQueryStatisticsPostgreSQL
	}
	if err := s.get(ctx, r.RequestURL(), &parsed); err != nil {
		return nil, err
	}
	return parsed.Pg, nil
}

// DeleteManagedDatabase (EXPERIMENTAL) deletes an existing managed database instance
func (s *ServiceContext) DeleteManagedDatabase(ctx context.Context, r *request.DeleteManagedDatabaseRequest) error {
	return s.delete(ctx, r)
}

// ModifyManagedDatabase (EXPERIMENTAL) modifies an existing managed database instance
func (s *ServiceContext) ModifyManagedDatabase(ctx context.Context, r *request.ModifyManagedDatabaseRequest) (*upcloud.ManagedDatabase, error) {
	managedDatabaseDetails := upcloud.ManagedDatabase{}
	return &managedDatabaseDetails, s.modify(ctx, r, &managedDatabaseDetails)
}

// UpgradeManagedDatabaseServiceVersion upgrades the version of the database service;
// for the list of available versions use GetManagedDatabaseVersions function
func (s *ServiceContext) UpgradeManagedDatabaseVersion(ctx context.Context, r *request.UpgradeManagedDatabaseVersionRequest) (*upcloud.ManagedDatabase, error) {
	managedDatabaseDetails := upcloud.ManagedDatabase{}
	return &managedDatabaseDetails, s.create(ctx, r, &managedDatabaseDetails)
}

// GetManagedDatabaseVersions available versions of the specific Managed Database service
func (s *ServiceContext) GetManagedDatabaseVersions(ctx context.Context, r *request.GetManagedDatabaseVersionsRequest) ([]string, error) {
	versions := make([]string, 0)
	return versions, s.get(ctx, r.RequestURL(), &versions)
}

// WaitForManagedDatabaseState (EXPERIMENTAL) blocks execution until the specified managed database instance has entered the
// specified state. If the state changes favorably, the new managed database details is returned. The method will give up
// after the specified timeout
func (s *ServiceContext) WaitForManagedDatabaseState(ctx context.Context, r *request.WaitForManagedDatabaseStateRequest) (*upcloud.ManagedDatabase, error) {
	attempts := 0
	sleepDuration := time.Second * 5

	for {
		attempts++

		details, err := s.GetManagedDatabase(ctx, &request.GetManagedDatabaseRequest{
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
func (s *ServiceContext) StartManagedDatabase(ctx context.Context, r *request.StartManagedDatabaseRequest) (*upcloud.ManagedDatabase, error) {
	managedDatabaseDetails := upcloud.ManagedDatabase{}
	return &managedDatabaseDetails, s.modify(ctx, r, &managedDatabaseDetails)
}

// ShutdownManagedDatabase (EXPERIMENTAL) shuts down existing managed database instance. Only a service which has at least one
// full backup can be shut down.
func (s *ServiceContext) ShutdownManagedDatabase(ctx context.Context, r *request.ShutdownManagedDatabaseRequest) (*upcloud.ManagedDatabase, error) {
	managedDatabaseDetails := upcloud.ManagedDatabase{}
	return &managedDatabaseDetails, s.modify(ctx, r, &managedDatabaseDetails)
}

/* User Management */

// CreateManagedDatabaseUser (EXPERIMENTAL) creates a new normal user to an existing managed database instance
func (s *ServiceContext) CreateManagedDatabaseUser(ctx context.Context, r *request.CreateManagedDatabaseUserRequest) (*upcloud.ManagedDatabaseUser, error) {
	userDetails := upcloud.ManagedDatabaseUser{}
	return &userDetails, s.create(ctx, r, &userDetails)
}

// GetManagedDatabaseUser (EXPERIMENTAL) returns details of an existing user of an existing managed database instance
func (s *ServiceContext) GetManagedDatabaseUser(ctx context.Context, r *request.GetManagedDatabaseUserRequest) (*upcloud.ManagedDatabaseUser, error) {
	userDetails := upcloud.ManagedDatabaseUser{}
	return &userDetails, s.get(ctx, r.RequestURL(), &userDetails)
}

// GetManagedDatabaseUsers (EXPERIMENTAL) returns a slice of all users of an existing managed database instance
func (s *ServiceContext) GetManagedDatabaseUsers(ctx context.Context, r *request.GetManagedDatabaseUsersRequest) ([]upcloud.ManagedDatabaseUser, error) {
	userList := make([]upcloud.ManagedDatabaseUser, 0)
	return userList, s.get(ctx, r.RequestURL(), &userList)
}

// DeleteManagedDatabaseUser (EXPERIMENTAL) deletes an existing user of an existing managed database instance
func (s *ServiceContext) DeleteManagedDatabaseUser(ctx context.Context, r *request.DeleteManagedDatabaseUserRequest) error {
	return s.delete(ctx, r)
}

// ModifyManagedDatabaseUser (EXPERIMENTAL) modifies an existing user of an existing managed database instance
func (s *ServiceContext) ModifyManagedDatabaseUser(ctx context.Context, r *request.ModifyManagedDatabaseUserRequest) (*upcloud.ManagedDatabaseUser, error) {
	userDetails := upcloud.ManagedDatabaseUser{}
	return &userDetails, s.modify(ctx, r, &userDetails)
}

/* Logical Database Management */

// CreateManagedDatabaseLogicalDatabase (EXPERIMENTAL) creates a new logical database to an existing managed database instance
func (s *ServiceContext) CreateManagedDatabaseLogicalDatabase(ctx context.Context, r *request.CreateManagedDatabaseLogicalDatabaseRequest) (*upcloud.ManagedDatabaseLogicalDatabase, error) {
	dbDetails := upcloud.ManagedDatabaseLogicalDatabase{}
	return &dbDetails, s.create(ctx, r, &dbDetails)
}

// GetManagedDatabaseLogicalDatabases (EXPERIMENTAL) returns a slice of all logical databases of an existing managed database instance
func (s *ServiceContext) GetManagedDatabaseLogicalDatabases(ctx context.Context, r *request.GetManagedDatabaseLogicalDatabasesRequest) ([]upcloud.ManagedDatabaseLogicalDatabase, error) {
	var dbList []upcloud.ManagedDatabaseLogicalDatabase
	return dbList, s.get(ctx, r.RequestURL(), &dbList)
}

// DeleteManagedDatabaseLogicalDatabase (EXPERIMENTAL) deletes an existing logical database of an existing managed database instance
func (s *ServiceContext) DeleteManagedDatabaseLogicalDatabase(ctx context.Context, r *request.DeleteManagedDatabaseLogicalDatabaseRequest) error {
	return s.delete(ctx, r)
}
