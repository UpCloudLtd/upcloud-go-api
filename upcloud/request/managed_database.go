package request

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/UpCloudLtd/upcloud-go-api/v5/upcloud"
)

/* Service Management */

// CancelManagedDatabaseConnection represents a request to cancel the current query of a connection or terminate
// the entire connection.
type CancelManagedDatabaseConnection struct {
	// UUID selects a managed database instance to manage queries of
	UUID string
	// Pid selects a connection pid to cancel or terminate
	Pid int
	// Terminate selects whether the connection will be forcefully terminated
	Terminate bool
}

// RequestURL implements the request.Request interface
func (c *CancelManagedDatabaseConnection) RequestURL() string {
	qv := url.Values{}
	if c.Terminate {
		qv.Set("terminate", "true")
	}
	return fmt.Sprintf("/database/%s/connections/%d?%s", c.UUID, c.Pid, qv.Encode())
}

// CloneManagedDatabaseRequest represents a request to cancel
type CloneManagedDatabaseRequest struct {
	// UUID selects an existing managed database instance to clone
	UUID string `json:"-"`
	// CloneTime selects a point-in-time from where to clone the data. Zero value selects the most recent available.
	CloneTime time.Time `json:"clone_time"`

	// Only for Redis. Create a clone of your database service data from the backups by name.
	BackupName string `json:"backup_name,omitempty"`

	HostNamePrefix string                                `json:"hostname_prefix"`
	Maintenance    ManagedDatabaseMaintenanceTimeRequest `json:"maintenance,omitempty"`
	Plan           string                                `json:"plan"`
	Properties     ManagedDatabasePropertiesRequest      `json:"properties,omitempty"`
	Title          string                                `json:"title,omitempty"`
	Type           upcloud.ManagedDatabaseServiceType    `json:"type"`
	Zone           string                                `json:"zone"`
}

// MarshalJSON implements json.Marshaler
func (c CloneManagedDatabaseRequest) MarshalJSON() ([]byte, error) {
	type alias CloneManagedDatabaseRequest
	req := struct {
		alias
		CloneTime   *time.Time                             `json:"clone_time,omitempty"`
		Maintenance *ManagedDatabaseMaintenanceTimeRequest `json:"maintenance,omitempty"`
	}{alias: alias(c)}
	if !req.alias.CloneTime.IsZero() {
		req.CloneTime = &req.alias.CloneTime
	}
	if c.Maintenance.Time != "" || c.Maintenance.DayOfWeek != "" {
		req.Maintenance = &c.Maintenance
	}
	return json.Marshal(&req)
}

// RequestURL implements the request.Request interface
func (c *CloneManagedDatabaseRequest) RequestURL() string {
	return fmt.Sprintf("/database/%s/clone", c.UUID)
}

// CreateManagedDatabaseRequest represents a request to create a new managed database instance
type CreateManagedDatabaseRequest struct {
	HostNamePrefix string                                `json:"hostname_prefix"`
	Maintenance    ManagedDatabaseMaintenanceTimeRequest `json:"maintenance,omitempty"`
	Plan           string                                `json:"plan"`
	Properties     ManagedDatabasePropertiesRequest      `json:"properties,omitempty"`
	Title          string                                `json:"title,omitempty"`
	Type           upcloud.ManagedDatabaseServiceType    `json:"type"`
	Zone           string                                `json:"zone"`
}

// MarshalJSON implements json.Marshaler
func (c CreateManagedDatabaseRequest) MarshalJSON() ([]byte, error) {
	type alias CreateManagedDatabaseRequest
	req := struct {
		alias
		Maintenance *ManagedDatabaseMaintenanceTimeRequest `json:"maintenance,omitempty"`
	}{alias: alias(c)}
	if c.Maintenance.Time != "" || c.Maintenance.DayOfWeek != "" {
		req.Maintenance = &c.Maintenance
	}
	return json.Marshal(&req)
}

// RequestURL implements the request.Request interface
func (c *CreateManagedDatabaseRequest) RequestURL() string {
	return "/database"
}

// DeleteManagedDatabaseRequest represents a request to delete an existing managed database instance
type DeleteManagedDatabaseRequest struct {
	UUID string
}

// RequestURL implements the request.Request interface
func (d *DeleteManagedDatabaseRequest) RequestURL() string {
	return fmt.Sprintf("/database/%s", d.UUID)
}

// GetManagedDatabaseRequest represents a request to get details of an existing managed database instance
type GetManagedDatabaseRequest struct {
	UUID string
}

// RequestURL implements the request.Request interface
func (g *GetManagedDatabaseRequest) RequestURL() string {
	return fmt.Sprintf("/database/%s", g.UUID)
}

// GetManagedDatabasesRequest represents a request to get a slice of existing managed database instances
type GetManagedDatabasesRequest struct{}

// RequestURL implements the request.Request interface
func (g *GetManagedDatabasesRequest) RequestURL() string {
	return "/database"
}

// GetManagedDatabaseConnectionsRequest represents a request to get managed database instance's current connections
type GetManagedDatabaseConnectionsRequest struct {
	// UUID selects a managed database instance to query connections from
	UUID string
	// Limit sets the upper bound how many connections to fetch
	Limit int
	// Offset skips n connections before returning them. It can be used to iteratively fetch connections.
	Offset int
}

// RequestURL implements the request.Request interface
func (g *GetManagedDatabaseConnectionsRequest) RequestURL() string {
	qv := url.Values{}
	if g.Limit != 0 {
		qv.Set("limit", strconv.Itoa(g.Limit))
	}
	if g.Offset != 0 {
		qv.Set("offset", strconv.Itoa(g.Offset))
	}
	return fmt.Sprintf("/database/%s/connections?%s", g.UUID, qv.Encode())
}

// GetManagedDatabaseMetricsRequest represents a request to get managed database instance performance metrics
type GetManagedDatabaseMetricsRequest struct {
	// UUID selects a managed database instance to query metrics from
	UUID string
	// Period selects the observation window. See:
	// 	upcloud.ManagedDatabaseMetricPeriodHour
	// 	upcloud.ManagedDatabaseMetricPeriodDay
	// 	upcloud.ManagedDatabaseMetricPeriodWeek
	// 	upcloud.ManagedDatabaseMetricPeriodMonth
	// 	upcloud.ManagedDatabaseMetricPeriodYear
	Period upcloud.ManagedDatabaseMetricPeriod
}

// RequestURL implements the request.Request interface
func (g *GetManagedDatabaseMetricsRequest) RequestURL() string {
	qv := url.Values{}
	if g.Period != "" {
		qv.Set("period", string(g.Period))
	}
	return fmt.Sprintf("/database/%s/metrics?%s", g.UUID, qv.Encode())
}

// GetManagedDatabaseLogsRequest represents a request to get managed database instance logs
type GetManagedDatabaseLogsRequest struct {
	// UUID selects a managed database instance to query logs from
	UUID string
	// Limit sets the maximum number of logs to query in one go
	Limit int
	// Offset sets the offset from which to query logs onwards
	Offset string
	// Order sets the log sort order. See:
	// 	upcloud.ManagedDatabaseLogOrderAscending
	// 	upcloud.ManagedDatabaseLogOrderDescending
	Order upcloud.ManagedDatabaseLogOrder
}

// RequestURL implements the request.Request interface
func (g *GetManagedDatabaseLogsRequest) RequestURL() string {
	qv := url.Values{}
	if g.Limit != 0 {
		qv.Set("limit", strconv.Itoa(g.Limit))
	}
	if g.Offset != "" {
		qv.Set("offset", g.Offset)
	}
	if g.Order != "" {
		qv.Set("order", string(g.Order))
	}
	if g.Offset == "" && g.Order == upcloud.ManagedDatabaseLogOrderAscending {
		qv.Set("offset", "0")
	}
	return fmt.Sprintf("/database/%s/logs?%s", g.UUID, qv.Encode())
}

type GetManagedDatabaseQueryStatisticsRequest struct {
	// UUID selects a managed database instance to query statistics from
	UUID string
	// Limit sets the upper bound how many query stats to fetch
	Limit int
	// Offset skips n query stat rows before returning them. It can be used to iteratively fetch all.
	Offset int
}

// RequestURL implements the request.Request interface
func (g *GetManagedDatabaseQueryStatisticsRequest) RequestURL() string {
	qv := url.Values{}
	if g.Limit != 0 {
		qv.Set("limit", strconv.Itoa(g.Limit))
	}
	if g.Offset != 0 {
		qv.Set("offset", strconv.Itoa(g.Offset))
	}
	return fmt.Sprintf("/database/%s/query-statistics?%s", g.UUID, qv.Encode())
}

// GetManagedDatabaseServiceTypeRequest represents a request to get details of a database type
type GetManagedDatabaseServiceTypeRequest struct {
	Type string
}

// RequestURL implements the request.Request interface
func (g *GetManagedDatabaseServiceTypeRequest) RequestURL() string {
	return fmt.Sprintf("/database/service-types/%s", g.Type)
}

// GetManagedDatabaseServiceTypesRequest represents a request to get a map of available database types
type GetManagedDatabaseServiceTypesRequest struct{}

// RequestURL implements the request.Request interface
func (g *GetManagedDatabaseServiceTypesRequest) RequestURL() string {
	return "/database/service-types"
}

// ManagedDatabaseMaintenanceTimeRequest represents the set time of week when automatic maintenance operations are allowed
type ManagedDatabaseMaintenanceTimeRequest struct {
	DayOfWeek string `json:"dow,omitempty"`
	Time      string `json:"time,omitempty"`
}

// ManagedDatabasePropertiesRequest is a Properties helper type for CreateManagedDatabaseRequest and ModifyManagedDatabaseRequest.
// It allows initialisation by chaining the Set methods.
//
// For example:
//
//	req := CreateManagedDatabaseRequest{}; req.Properties.SetString("foo", "bar").Set("test", customValue)
type ManagedDatabasePropertiesRequest map[upcloud.ManagedDatabasePropertyKey]interface{}

// Set associates key with a any type of value. The underlying map is initialised if it's nil
func (m *ManagedDatabasePropertiesRequest) Set(name upcloud.ManagedDatabasePropertyKey, value interface{}) *ManagedDatabasePropertiesRequest {
	if m == nil {
		return nil
	}
	if *m == nil {
		*m = make(ManagedDatabasePropertiesRequest)
	}
	(*m)[name] = value
	return m
}

// SetString associates key with a string value. The underlying map is initialised if it's nil
func (m *ManagedDatabasePropertiesRequest) SetString(name upcloud.ManagedDatabasePropertyKey, value string) *ManagedDatabasePropertiesRequest {
	return m.Set(name, value)
}

// SetInt associates key with an integer value. The underlying map is initialised if it's nil
func (m *ManagedDatabasePropertiesRequest) SetInt(name upcloud.ManagedDatabasePropertyKey, value int) *ManagedDatabasePropertiesRequest {
	return m.Set(name, value)
}

// SetBool associates key with a boolean value. The underlying map is initialised if it's nil
func (m *ManagedDatabasePropertiesRequest) SetBool(name upcloud.ManagedDatabasePropertyKey, value bool) *ManagedDatabasePropertiesRequest {
	return m.Set(name, value)
}

// SetStringSlice associates key with a slice of strings. The underlying map is initialised if it's nil
func (m *ManagedDatabasePropertiesRequest) SetStringSlice(name upcloud.ManagedDatabasePropertyKey, value []string) *ManagedDatabasePropertiesRequest {
	return m.Set(name, value)
}

// SetAutoUtilityIPFilter enables or disables automatic utility network ip filtering.
// See upcloud.ManagedDatabasePropertyAutoUtilityIPFilter for more information.
func (m *ManagedDatabasePropertiesRequest) SetAutoUtilityIPFilter(enabled bool) *ManagedDatabasePropertiesRequest {
	return m.SetBool(upcloud.ManagedDatabasePropertyAutoUtilityIPFilter, enabled)
}

// SetIPFilter sets the list of allowed host or networks that can access the service.
//
// Use upcloud.ManagedDatabaseAllIPv4 to enable access from anywhere.
//
// See upcloud.ManagedDatabasePropertyIPFilter for more information.
func (m *ManagedDatabasePropertiesRequest) SetIPFilter(addressOrNetworkWithCIDRMask ...string) *ManagedDatabasePropertiesRequest {
	return m.SetStringSlice(upcloud.ManagedDatabasePropertyIPFilter, addressOrNetworkWithCIDRMask)
}

// SetPublicAccess enables or disables public access from the internet.
// See upcloud.ManagedDatabasePropertyPublicAccess for more information.
func (m *ManagedDatabasePropertiesRequest) SetPublicAccess(enabled bool) *ManagedDatabasePropertiesRequest {
	return m.SetBool(upcloud.ManagedDatabasePropertyPublicAccess, enabled)
}

// Get returns a property value by name. The underlying map is initialised if it's nil
func (m *ManagedDatabasePropertiesRequest) Get(name upcloud.ManagedDatabasePropertyKey) interface{} {
	if m == nil {
		return nil
	}
	if *m == nil {
		*m = make(ManagedDatabasePropertiesRequest)
	}
	return (*m)[name]
}

// GetString returns a string property value.
// The underlying map is initialised if it's nil.
func (m *ManagedDatabasePropertiesRequest) GetString(name upcloud.ManagedDatabasePropertyKey) (string, error) {
	if v, ok := m.Get(name).(string); ok {
		return v, nil
	}
	return "", fmt.Errorf("not a string property %q", name)
}

// GetInt returns an integer property value.
// The underlying map is initialised if it's nil.
func (m *ManagedDatabasePropertiesRequest) GetInt(name upcloud.ManagedDatabasePropertyKey) (int, error) {
	if v, ok := m.Get(name).(int); ok {
		return v, nil
	}
	return 0, fmt.Errorf("not an int property %q", name)
}

// GetBool returns a boolean property value.
// The underlying map is initialised if it's nil.
func (m *ManagedDatabasePropertiesRequest) GetBool(name upcloud.ManagedDatabasePropertyKey) (bool, error) {
	if v, ok := m.Get(name).(bool); ok {
		return v, nil
	}
	return false, fmt.Errorf("not a boolean property %q", name)
}

// GetStringSlice returns a string-slice property value.
// The underlying map is initialised if it's nil.
func (m *ManagedDatabasePropertiesRequest) GetStringSlice(name upcloud.ManagedDatabasePropertyKey) ([]string, error) {
	if slice, ok := m.Get(name).([]string); ok {
		return slice, nil
	}
	slice, ok := m.Get(name).([]interface{})
	if !ok {
		return nil, fmt.Errorf("not a string-slice property %q", name)
	}
	stringSlice := make([]string, len(slice))
	for i, iv := range slice {
		v, ok := iv.(string)
		if !ok {
			return nil, fmt.Errorf("not a string-slice property %q", name)
		}
		stringSlice[i] = v
	}
	return stringSlice, nil
}

// GetAutoUtilityIPFilter returns the state of automatic utility network IP filtering.
// See upcloud.ManagedDatabasePropertyAutoUtilityIPFilter for more information.
func (m *ManagedDatabasePropertiesRequest) GetAutoUtilityIPFilter() bool {
	v, _ := m.GetBool(upcloud.ManagedDatabasePropertyAutoUtilityIPFilter)
	return v
}

// GetIPFilter returns a slice of allowed hosts or networks.
// See upcloud.ManagedDatabasePropertyIPFilter for more information.
func (m *ManagedDatabasePropertiesRequest) GetIPFilter() []string {
	v, _ := m.GetStringSlice(upcloud.ManagedDatabasePropertyIPFilter)
	return v
}

// GetPublicAccess returns the state of public access to the service.
// See upcloud.ManagedDatabasePropertyPublicAccess for more information.
func (m *ManagedDatabasePropertiesRequest) GetPublicAccess() bool {
	v, _ := m.GetBool(upcloud.ManagedDatabasePropertyPublicAccess)
	return v
}

// ModifyManagedDatabaseRequest represents a request to modify an existing managed database instance
type ModifyManagedDatabaseRequest struct {
	Maintenance ManagedDatabaseMaintenanceTimeRequest `json:"maintenance"`
	Plan        string                                `json:"plan,omitempty"`
	Properties  ManagedDatabasePropertiesRequest      `json:"properties,omitempty"`
	Title       string                                `json:"title,omitempty"`
	Type        upcloud.ManagedDatabaseServiceType    `json:"type,omitempty"`
	UUID        string                                `json:"-"`
	Zone        string                                `json:"zone,omitempty"`
}

// MarshalJSON implements json.Marshaler
func (m ModifyManagedDatabaseRequest) MarshalJSON() ([]byte, error) {
	type alias ModifyManagedDatabaseRequest
	req := struct {
		alias
		Maintenance *ManagedDatabaseMaintenanceTimeRequest `json:"maintenance,omitempty"`
	}{alias: alias(m)}
	if m.Maintenance.Time != "" || m.Maintenance.DayOfWeek != "" {
		req.Maintenance = &m.Maintenance
	}
	return json.Marshal(&req)
}

// RequestURL implements the request.Request interface
func (m *ModifyManagedDatabaseRequest) RequestURL() string {
	return fmt.Sprintf("/database/%s", m.UUID)
}

// UpgradeManagedDatabaseVersionRequest represents a request to upgrade an existing managed database version;
// for a list of available versions use GetManagedDatabaseVersionsRequest
type UpgradeManagedDatabaseVersionRequest struct {
	UUID          string `json:"-"`
	TargetVersion string `json:"target_version,omitempty"`
}

// RequestURL implements the request.Request interface
func (r *UpgradeManagedDatabaseVersionRequest) RequestURL() string {
	return fmt.Sprintf("/database/%s/upgrade", r.UUID)
}

// GetManagedDatabaseVersionsRequests represents a request to list available versions of the Managed Database service by its UUID
type GetManagedDatabaseVersionsRequest struct {
	UUID string `json:"-"`
}

// RequestURL implements the request.Request interface
func (r *GetManagedDatabaseVersionsRequest) RequestURL() string {
	return fmt.Sprintf("/database/%s/versions", r.UUID)
}

// WaitForManagedDatabaseStateRequest represents a request to wait for a managed database instance to enter a specific state
type WaitForManagedDatabaseStateRequest struct {
	UUID         string
	DesiredState upcloud.ManagedDatabaseState
	Timeout      time.Duration
}

// StartManagedDatabaseRequest represents a request to start an existing managed database instance
type StartManagedDatabaseRequest struct {
	UUID string
}

// MarshalJSON implements json.Marshaler
func (m *StartManagedDatabaseRequest) MarshalJSON() ([]byte, error) {
	return json.RawMessage(`{"powered":true}`), nil
}

// RequestURL implements the request.Request interface
func (m *StartManagedDatabaseRequest) RequestURL() string {
	return fmt.Sprintf("/database/%s", m.UUID)
}

// ShutdownManagedDatabaseRequest represents a request to shut down an existing managed database instance
type ShutdownManagedDatabaseRequest struct {
	UUID string
}

// MarshalJSON implements json.Marshaler
func (m *ShutdownManagedDatabaseRequest) MarshalJSON() ([]byte, error) {
	return json.RawMessage(`{"powered":false}`), nil
}

// RequestURL implements the request.Request interface
func (m *ShutdownManagedDatabaseRequest) RequestURL() string {
	return fmt.Sprintf("/database/%s", m.UUID)
}

/* User Management */

// CreateManagedDatabaseUserRequest represents a request to create a new normal user to an existing managed database instance
type CreateManagedDatabaseUserRequest struct {
	// ServiceUUID selects a managed database service to modify
	ServiceUUID string `json:"-"`
	Username    string `json:"username"`
	Password    string `json:"password,omitempty"`
	// Authentication selects authentication type for the user. See following constants for more information:
	// 	upcloud.ManagedDatabaseUserAuthenticationCachingSHA2Password
	// 	upcloud.ManagedDatabaseUserAuthenticationMySQLNativePassword
	Authentication     upcloud.ManagedDatabaseUserAuthenticationType  `json:"authentication,omitempty"`
	PGAccessControl    *upcloud.ManagedDatabaseUserPGAccessControl    `json:"pg_access_control,omitempty"`
	RedisAccessControl *upcloud.ManagedDatabaseUserRedisAccessControl `json:"redis_access_control,omitempty"`
}

// RequestURL implements the request.Request interface
func (m *CreateManagedDatabaseUserRequest) RequestURL() string {
	return fmt.Sprintf("/database/%s/users", m.ServiceUUID)
}

// DeleteManagedDatabaseUserRequest represents a request to delete a normal user from an existing managed database instance
type DeleteManagedDatabaseUserRequest struct {
	// ServiceUUID selects a managed database service to modify
	ServiceUUID string `json:"-"`
	// Username selects the username to delete
	Username string `json:"-"`
}

// RequestURL implements the request.Request interface
func (m *DeleteManagedDatabaseUserRequest) RequestURL() string {
	return fmt.Sprintf("/database/%s/users/%s", m.ServiceUUID, m.Username)
}

// GetManagedDatabaseUserRequest represents a request to get details of a user of an existing managed database
// instance. This request also returns the password of the user if it's known by the API.
type GetManagedDatabaseUserRequest struct {
	// ServiceUUID selects a managed database service to query
	ServiceUUID string `json:"-"`
	// Username selects the username to get
	Username string `json:"-"`
}

// RequestURL implements the request.Request interface
func (g *GetManagedDatabaseUserRequest) RequestURL() string {
	return fmt.Sprintf("/database/%s/users/%s", g.ServiceUUID, g.Username)
}

// GetManagedDatabaseUsersRequest represents a request to get a slice of users of an existing managed database instance
// The returned response doesn't contain the passwords of the users.
type GetManagedDatabaseUsersRequest struct {
	// ServiceUUID selects a managed database service to query
	ServiceUUID string `json:"-"`
}

// RequestURL implements the request.Request interface
func (g *GetManagedDatabaseUsersRequest) RequestURL() string {
	return fmt.Sprintf("/database/%s/users", g.ServiceUUID)
}

// ModifyManagedDatabaseUserRequest represents a request to modify an existing user of an existing managed database instance
type ModifyManagedDatabaseUserRequest struct {
	// ServiceUUID selects a managed database service to modify
	ServiceUUID string `json:"-"`
	// Username selects the username to modify. The username itself is immutable. To change it, recreate the user.
	Username string `json:"-"`
	Password string `json:"password,omitempty"`
	// Authentication selects authentication type for the user. See following constants for more information:
	// 	upcloud.ManagedDatabaseUserAuthenticationCachingSHA2Password
	// 	upcloud.ManagedDatabaseUserAuthenticationMySQLNativePassword
	Authentication upcloud.ManagedDatabaseUserAuthenticationType `json:"authentication,omitempty"`
}

// RequestURL implements the request.Request interface
func (m *ModifyManagedDatabaseUserRequest) RequestURL() string {
	return fmt.Sprintf("/database/%s/users/%s", m.ServiceUUID, m.Username)
}

// ModifyManagedDatabaseUserRequest represents a request to modify an existing user of an existing managed database instance
type ModifyManagedDatabaseUserAccessControlRequest struct {
	// ServiceUUID selects a managed database service to modify
	ServiceUUID string `json:"-"`
	// Username selects the username to modify. The username itself is immutable. To change it, recreate the user.
	Username           string                                         `json:"-"`
	PGAccessControl    *upcloud.ManagedDatabaseUserPGAccessControl    `json:"pg_access_control,omitempty"`
	RedisAccessControl *upcloud.ManagedDatabaseUserRedisAccessControl `json:"redis_access_control,omitempty"`
}

// RequestURL implements the request.Request interface
func (m *ModifyManagedDatabaseUserAccessControlRequest) RequestURL() string {
	return fmt.Sprintf("/database/%s/users/%s/access-control", m.ServiceUUID, m.Username)
}

/* Logical Database Management */

// CreateManagedDatabaseLogicalDatabaseRequest represents a request to create a new logical database to an existing
// managed database instance
type CreateManagedDatabaseLogicalDatabaseRequest struct {
	// ServiceUUID selects a managed database service to modify
	ServiceUUID string `json:"-"`
	Name        string `json:"name"`
	// LCCollate represents a default string sort order of a logical database
	LCCollate string `json:"lc_collate"`
	// LCCType represents a default character classification of a logical database
	LCCType string `json:"lc_ctype"`
}

// RequestURL implements the request.Request interface
func (c *CreateManagedDatabaseLogicalDatabaseRequest) RequestURL() string {
	return fmt.Sprintf("/database/%s/databases", c.ServiceUUID)
}

// GetManagedDatabaseLogicalDatabasesRequest represents a request to get a slice of existing logical databases
// of a managed database instance
type GetManagedDatabaseLogicalDatabasesRequest struct {
	// ServiceUUID selects a managed database service to query
	ServiceUUID string `json:"-"`
}

// RequestURL implements the request.Request interface
func (g *GetManagedDatabaseLogicalDatabasesRequest) RequestURL() string {
	return fmt.Sprintf("/database/%s/databases", g.ServiceUUID)
}

// DeleteManagedDatabaseLogicalDatabaseRequest represents a request to delete a logical database from an existing
// managed database instance
type DeleteManagedDatabaseLogicalDatabaseRequest struct {
	// ServiceUUID selects a managed database service to modify
	ServiceUUID string `json:"-"`
	Name        string `json:"-"`
}

// RequestURL implements the request.Request interface
func (d *DeleteManagedDatabaseLogicalDatabaseRequest) RequestURL() string {
	return fmt.Sprintf("/database/%s/databases/%s", d.ServiceUUID, d.Name)
}
