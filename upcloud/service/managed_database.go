package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/UpCloudLtd/upcloud-go-api/v6/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v6/upcloud/request"
)

var ErrCancelManagedDatabaseConnection = errors.New("managed database connection cancellation failed")

type ManagedDatabaseServiceManager interface {
	CancelManagedDatabaseConnection(ctx context.Context, r *request.CancelManagedDatabaseConnection) error
	CloneManagedDatabase(ctx context.Context, r *request.CloneManagedDatabaseRequest) (*upcloud.ManagedDatabase, error)
	CreateManagedDatabase(ctx context.Context, r *request.CreateManagedDatabaseRequest) (*upcloud.ManagedDatabase, error)
	GetManagedDatabase(ctx context.Context, r *request.GetManagedDatabaseRequest) (*upcloud.ManagedDatabase, error)
	GetManagedDatabases(ctx context.Context, r *request.GetManagedDatabasesRequest) ([]upcloud.ManagedDatabase, error)
	GetManagedDatabaseAccessControl(ctx context.Context, r *request.GetManagedDatabaseAccessControlRequest) (*upcloud.ManagedDatabaseAccessControl, error)
	GetManagedDatabaseConnections(ctx context.Context, r *request.GetManagedDatabaseConnectionsRequest) ([]upcloud.ManagedDatabaseConnection, error)
	GetManagedDatabaseIndex(ctx context.Context, r *request.GetManagedDatabaseIndexRequest) (*upcloud.ManagedDatabaseIndex, error)
	GetManagedDatabaseIndices(ctx context.Context, r *request.GetManagedDatabaseIndicesRequest) ([]upcloud.ManagedDatabaseIndex, error)
	GetManagedDatabaseMetrics(ctx context.Context, r *request.GetManagedDatabaseMetricsRequest) (*upcloud.ManagedDatabaseMetrics, error)
	GetManagedDatabaseLogs(ctx context.Context, r *request.GetManagedDatabaseLogsRequest) (*upcloud.ManagedDatabaseLogs, error)
	GetManagedDatabaseQueryStatisticsMySQL(ctx context.Context, r *request.GetManagedDatabaseQueryStatisticsRequest) ([]upcloud.ManagedDatabaseQueryStatisticsMySQL, error)
	GetManagedDatabaseQueryStatisticsPostgreSQL(ctx context.Context, r *request.GetManagedDatabaseQueryStatisticsRequest) ([]upcloud.ManagedDatabaseQueryStatisticsPostgreSQL, error)
	DeleteManagedDatabase(ctx context.Context, r *request.DeleteManagedDatabaseRequest) error
	DeleteManagedDatabaseIndex(ctx context.Context, r *request.DeleteManagedDatabaseIndexRequest) error
	ModifyManagedDatabase(ctx context.Context, r *request.ModifyManagedDatabaseRequest) (*upcloud.ManagedDatabase, error)
	ModifyManagedDatabaseAccessControl(ctx context.Context, r *request.ModifyManagedDatabaseAccessControlRequest) (*upcloud.ManagedDatabaseAccessControl, error)
	UpgradeManagedDatabaseVersion(ctx context.Context, r *request.UpgradeManagedDatabaseVersionRequest) (*upcloud.ManagedDatabase, error)
	GetManagedDatabaseVersions(ctx context.Context, r *request.GetManagedDatabaseVersionsRequest) ([]string, error)
	StartManagedDatabase(ctx context.Context, r *request.StartManagedDatabaseRequest) (*upcloud.ManagedDatabase, error)
	ShutdownManagedDatabase(ctx context.Context, r *request.ShutdownManagedDatabaseRequest) (*upcloud.ManagedDatabase, error)
	WaitForManagedDatabaseState(ctx context.Context, r *request.WaitForManagedDatabaseStateRequest) (*upcloud.ManagedDatabase, error)
	GetManagedDatabaseServiceType(ctx context.Context, r *request.GetManagedDatabaseServiceTypeRequest) (*upcloud.ManagedDatabaseType, error)
	GetManagedDatabaseServiceTypes(ctx context.Context, r *request.GetManagedDatabaseServiceTypesRequest) (map[string]upcloud.ManagedDatabaseType, error)
}

type ManagedDatabaseUserManager interface {
	CreateManagedDatabaseUser(ctx context.Context, r *request.CreateManagedDatabaseUserRequest) (*upcloud.ManagedDatabaseUser, error)
	GetManagedDatabaseUser(ctx context.Context, r *request.GetManagedDatabaseUserRequest) (*upcloud.ManagedDatabaseUser, error)
	GetManagedDatabaseUsers(ctx context.Context, r *request.GetManagedDatabaseUsersRequest) ([]upcloud.ManagedDatabaseUser, error)
	DeleteManagedDatabaseUser(ctx context.Context, r *request.DeleteManagedDatabaseUserRequest) error
	ModifyManagedDatabaseUser(ctx context.Context, r *request.ModifyManagedDatabaseUserRequest) (*upcloud.ManagedDatabaseUser, error)
	ModifyManagedDatabaseUserAccessControl(ctx context.Context, r *request.ModifyManagedDatabaseUserAccessControlRequest) (*upcloud.ManagedDatabaseUser, error)
}

type ManagedDatabaseLogicalDatabaseManager interface {
	CreateManagedDatabaseLogicalDatabase(ctx context.Context, r *request.CreateManagedDatabaseLogicalDatabaseRequest) (*upcloud.ManagedDatabaseLogicalDatabase, error)
	GetManagedDatabaseLogicalDatabases(ctx context.Context, r *request.GetManagedDatabaseLogicalDatabasesRequest) ([]upcloud.ManagedDatabaseLogicalDatabase, error)
	DeleteManagedDatabaseLogicalDatabase(ctx context.Context, r *request.DeleteManagedDatabaseLogicalDatabaseRequest) error
}

/* Service Management */

// CancelManagedDatabaseConnection cancels a current query of a database connection or terminates it entirely.
// In case of the server is unable to cancel the query or terminate the connection ErrCancelManagedDatabaseConnection
// is returned.
func (s *Service) CancelManagedDatabaseConnection(ctx context.Context, r *request.CancelManagedDatabaseConnection) error {
	res := struct {
		Success bool `json:"success"`
	}{}
	response, err := s.client.Delete(ctx, r.RequestURL())
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

// CloneManagedDatabase clones en existing managed database instance
func (s *Service) CloneManagedDatabase(ctx context.Context, r *request.CloneManagedDatabaseRequest) (*upcloud.ManagedDatabase, error) {
	managedDatabaseDetails := upcloud.ManagedDatabase{}
	return &managedDatabaseDetails, s.create(ctx, r, &managedDatabaseDetails)
}

// CreateManagedDatabase creates a new managed database instance
func (s *Service) CreateManagedDatabase(ctx context.Context, r *request.CreateManagedDatabaseRequest) (*upcloud.ManagedDatabase, error) {
	managedDatabaseDetails := upcloud.ManagedDatabase{}
	return &managedDatabaseDetails, s.create(ctx, r, &managedDatabaseDetails)
}

// GetManagedDatabase gets details of an existing managed database instance
func (s *Service) GetManagedDatabase(ctx context.Context, r *request.GetManagedDatabaseRequest) (*upcloud.ManagedDatabase, error) {
	managedDatabaseDetails := upcloud.ManagedDatabase{}
	return &managedDatabaseDetails, s.get(ctx, r.RequestURL(), &managedDatabaseDetails)
}

// GetManagedDatabases returns a slice of all managed database instances within an account
func (s *Service) GetManagedDatabases(ctx context.Context, r *request.GetManagedDatabasesRequest) ([]upcloud.ManagedDatabase, error) {
	var services []upcloud.ManagedDatabase
	return services, s.get(ctx, r.RequestURL(), &services)
}

// GetManagedDatabaseConnections returns a slice of connections from an existing managed database instance
func (s *Service) GetManagedDatabaseConnections(ctx context.Context, r *request.GetManagedDatabaseConnectionsRequest) ([]upcloud.ManagedDatabaseConnection, error) {
	conns := make([]upcloud.ManagedDatabaseConnection, 0)
	return conns, s.get(ctx, r.RequestURL(), &conns)
}

// GetManagedDatabaseMetrics returns metrics collection for the selected period
func (s *Service) GetManagedDatabaseMetrics(ctx context.Context, r *request.GetManagedDatabaseMetricsRequest) (*upcloud.ManagedDatabaseMetrics, error) {
	metrics := upcloud.ManagedDatabaseMetrics{}
	return &metrics, s.get(ctx, r.RequestURL(), &metrics)
}

// GetManagedDatabaseLogs returns logs of a managed database instance
func (s *Service) GetManagedDatabaseLogs(ctx context.Context, r *request.GetManagedDatabaseLogsRequest) (*upcloud.ManagedDatabaseLogs, error) {
	logs := upcloud.ManagedDatabaseLogs{}
	return &logs, s.get(ctx, r.RequestURL(), &logs)
}

// GetManagedDatabaseQueryStatisticsMySQL returns MySQL query statistics of a managed database instance
func (s *Service) GetManagedDatabaseQueryStatisticsMySQL(ctx context.Context, r *request.GetManagedDatabaseQueryStatisticsRequest) ([]upcloud.ManagedDatabaseQueryStatisticsMySQL, error) {
	var parsed struct {
		Mysql []upcloud.ManagedDatabaseQueryStatisticsMySQL
	}
	if err := s.get(ctx, r.RequestURL(), &parsed); err != nil {
		return nil, err
	}
	return parsed.Mysql, nil
}

// GetManagedDatabaseQueryStatisticsPostgreSQL returns PostgreSQL query statistics of a managed database instance
func (s *Service) GetManagedDatabaseQueryStatisticsPostgreSQL(ctx context.Context, r *request.GetManagedDatabaseQueryStatisticsRequest) ([]upcloud.ManagedDatabaseQueryStatisticsPostgreSQL, error) {
	var parsed struct {
		Pg []upcloud.ManagedDatabaseQueryStatisticsPostgreSQL
	}
	if err := s.get(ctx, r.RequestURL(), &parsed); err != nil {
		return nil, err
	}
	return parsed.Pg, nil
}

// DeleteManagedDatabase deletes an existing managed database instance
func (s *Service) DeleteManagedDatabase(ctx context.Context, r *request.DeleteManagedDatabaseRequest) error {
	return s.delete(ctx, r)
}

// ModifyManagedDatabase modifies an existing managed database instance
func (s *Service) ModifyManagedDatabase(ctx context.Context, r *request.ModifyManagedDatabaseRequest) (*upcloud.ManagedDatabase, error) {
	managedDatabaseDetails := upcloud.ManagedDatabase{}
	return &managedDatabaseDetails, s.modify(ctx, r, &managedDatabaseDetails)
}

// UpgradeManagedDatabaseVersion upgrades the version of the database service;
// for the list of available versions use GetManagedDatabaseVersions function
func (s *Service) UpgradeManagedDatabaseVersion(ctx context.Context, r *request.UpgradeManagedDatabaseVersionRequest) (*upcloud.ManagedDatabase, error) {
	managedDatabaseDetails := upcloud.ManagedDatabase{}
	return &managedDatabaseDetails, s.create(ctx, r, &managedDatabaseDetails)
}

// GetManagedDatabaseVersions available versions of the specific Managed Database service
func (s *Service) GetManagedDatabaseVersions(ctx context.Context, r *request.GetManagedDatabaseVersionsRequest) ([]string, error) {
	versions := make([]string, 0)
	return versions, s.get(ctx, r.RequestURL(), &versions)
}

// WaitForManagedDatabaseState blocks execution until the specified managed database instance has entered the
// specified state. If the state changes favorably, the new managed database details is returned. The method will give up
// after the specified timeout
func (s *Service) WaitForManagedDatabaseState(ctx context.Context, r *request.WaitForManagedDatabaseStateRequest) (*upcloud.ManagedDatabase, error) {
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

// StartManagedDatabase starts a shut down existing managed database instance
func (s *Service) StartManagedDatabase(ctx context.Context, r *request.StartManagedDatabaseRequest) (*upcloud.ManagedDatabase, error) {
	managedDatabaseDetails := upcloud.ManagedDatabase{}
	return &managedDatabaseDetails, s.modify(ctx, r, &managedDatabaseDetails)
}

// ShutdownManagedDatabase shuts down existing managed database instance. Only a service which has at least one
// full backup can be shut down.
func (s *Service) ShutdownManagedDatabase(ctx context.Context, r *request.ShutdownManagedDatabaseRequest) (*upcloud.ManagedDatabase, error) {
	managedDatabaseDetails := upcloud.ManagedDatabase{}
	return &managedDatabaseDetails, s.modify(ctx, r, &managedDatabaseDetails)
}

/* User Management */

// CreateManagedDatabaseUser creates a new normal user to an existing managed database instance
func (s *Service) CreateManagedDatabaseUser(ctx context.Context, r *request.CreateManagedDatabaseUserRequest) (*upcloud.ManagedDatabaseUser, error) {
	userDetails := upcloud.ManagedDatabaseUser{}
	return &userDetails, s.create(ctx, r, &userDetails)
}

// GetManagedDatabaseUser returns details of an existing user of an existing managed database instance
func (s *Service) GetManagedDatabaseUser(ctx context.Context, r *request.GetManagedDatabaseUserRequest) (*upcloud.ManagedDatabaseUser, error) {
	userDetails := upcloud.ManagedDatabaseUser{}
	return &userDetails, s.get(ctx, r.RequestURL(), &userDetails)
}

// GetManagedDatabaseUsers returns a slice of all users of an existing managed database instance
func (s *Service) GetManagedDatabaseUsers(ctx context.Context, r *request.GetManagedDatabaseUsersRequest) ([]upcloud.ManagedDatabaseUser, error) {
	userList := make([]upcloud.ManagedDatabaseUser, 0)
	return userList, s.get(ctx, r.RequestURL(), &userList)
}

// DeleteManagedDatabaseUser deletes an existing user of an existing managed database instance
func (s *Service) DeleteManagedDatabaseUser(ctx context.Context, r *request.DeleteManagedDatabaseUserRequest) error {
	return s.delete(ctx, r)
}

// ModifyManagedDatabaseUser modifies an existing user of an existing managed database instance
func (s *Service) ModifyManagedDatabaseUser(ctx context.Context, r *request.ModifyManagedDatabaseUserRequest) (*upcloud.ManagedDatabaseUser, error) {
	userDetails := upcloud.ManagedDatabaseUser{}
	return &userDetails, s.modify(ctx, r, &userDetails)
}

/* Logical Database Management */

// CreateManagedDatabaseLogicalDatabase creates a new logical database to an existing managed database instance
func (s *Service) CreateManagedDatabaseLogicalDatabase(ctx context.Context, r *request.CreateManagedDatabaseLogicalDatabaseRequest) (*upcloud.ManagedDatabaseLogicalDatabase, error) {
	dbDetails := upcloud.ManagedDatabaseLogicalDatabase{}
	return &dbDetails, s.create(ctx, r, &dbDetails)
}

// GetManagedDatabaseLogicalDatabases returns a slice of all logical databases of an existing managed database instance
func (s *Service) GetManagedDatabaseLogicalDatabases(ctx context.Context, r *request.GetManagedDatabaseLogicalDatabasesRequest) ([]upcloud.ManagedDatabaseLogicalDatabase, error) {
	var dbList []upcloud.ManagedDatabaseLogicalDatabase
	return dbList, s.get(ctx, r.RequestURL(), &dbList)
}

// DeleteManagedDatabaseLogicalDatabase deletes an existing logical database of an existing managed database instance
func (s *Service) DeleteManagedDatabaseLogicalDatabase(ctx context.Context, r *request.DeleteManagedDatabaseLogicalDatabaseRequest) error {
	return s.delete(ctx, r)
}

// GetManagedDatabaseServiceType returns details of requested service type
func (s *Service) GetManagedDatabaseServiceType(ctx context.Context, r *request.GetManagedDatabaseServiceTypeRequest) (*upcloud.ManagedDatabaseType, error) {
	var serviceType upcloud.ManagedDatabaseType
	return &serviceType, s.get(ctx, r.RequestURL(), &serviceType)
}

// GetManagedDatabaseServiceTypes returns a map of available database service types
func (s *Service) GetManagedDatabaseServiceTypes(ctx context.Context, r *request.GetManagedDatabaseServiceTypesRequest) (map[string]upcloud.ManagedDatabaseType, error) {
	serviceTypes := make(map[string]upcloud.ManagedDatabaseType)
	return serviceTypes, s.get(ctx, r.RequestURL(), &serviceTypes)
}

// ModifyManagedDatabaseUserAccessControl modifies access control for an existing user
func (s *Service) ModifyManagedDatabaseUserAccessControl(ctx context.Context, r *request.ModifyManagedDatabaseUserAccessControlRequest) (*upcloud.ManagedDatabaseUser, error) {
	userDetails := upcloud.ManagedDatabaseUser{}
	return &userDetails, s.modify(ctx, r, &userDetails)
}

// GetManagedDatabaseAccessControl returns details of access control for an existing user
func (s *Service) GetManagedDatabaseAccessControl(ctx context.Context, r *request.GetManagedDatabaseAccessControlRequest) (*upcloud.ManagedDatabaseAccessControl, error) {
	ac := upcloud.ManagedDatabaseAccessControl{}
	return &ac, s.get(ctx, r.RequestURL(), &ac)
}

// ModifyManagedDatabaseAccessControl modifies access control settings for an existing OpenSearch Managed Database service
func (s *Service) ModifyManagedDatabaseAccessControl(ctx context.Context, r *request.ModifyManagedDatabaseAccessControlRequest) (*upcloud.ManagedDatabaseAccessControl, error) {
	ac := upcloud.ManagedDatabaseAccessControl{}
	return &ac, s.modify(ctx, r, &ac)
}

// GetManagedDatabaseIndices returns indices of an existing OpenSearch Managed Database service
func (s *Service) GetManagedDatabaseIndices(ctx context.Context, r *request.GetManagedDatabaseIndicesRequest) ([]upcloud.ManagedDatabaseIndex, error) {
	indices := make([]upcloud.ManagedDatabaseIndex, 0)
	return indices, s.get(ctx, r.RequestURL(), &indices)
}

// GetManagedDatabaseIndex returns details of an index existing OpenSearch Managed Database service
func (s *Service) GetManagedDatabaseIndex(ctx context.Context, r *request.GetManagedDatabaseIndexRequest) (*upcloud.ManagedDatabaseIndex, error) {
	index := upcloud.ManagedDatabaseIndex{}
	return &index, s.get(ctx, r.RequestURL(), &index)
}

// DeleteManagedDatabaseIndex returns details of an index existing OpenSearch Managed Database service
func (s *Service) DeleteManagedDatabaseIndex(ctx context.Context, r *request.DeleteManagedDatabaseIndexRequest) error {
	return s.delete(ctx, r)
}
