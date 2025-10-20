package upcloud

import (
	"encoding/json"
	"fmt"
	"time"
)

// ManagedDatabaseState represents a current state the service is in
type ManagedDatabaseState string

const (
	// Indicates newly created service or started reconfiguration
	ManagedDatabaseStatePending ManagedDatabaseState = "pending"
	// Configuring network
	ManagedDatabaseStateSetupNetwork ManagedDatabaseState = "setup-network"
	// Check that the network configuration was applied correctly
	ManagedDatabaseStateCheckNetwork ManagedDatabaseState = "check-network"
	// Configuring SDN network peerings if provided
	ManagedDatabaseStateSetupPeer ManagedDatabaseState = "setup-peer"
	// Check SDN network peerings was established if provided
	ManagedDatabaseStateCheckPeer ManagedDatabaseState = "check-peer"
	// Service is being created or reconfigured
	ManagedDatabaseStateSetupService ManagedDatabaseState = "setup-service"
	// Service creation in progress
	ManagedDatabaseStateRebuilding ManagedDatabaseState = "rebuilding"
	// Service is being upgraded or migrated
	ManagedDatabaseStateRebalancing ManagedDatabaseState = "rebalancing"
	// Service encountered an error that requires user action
	ManagedDatabaseStateError ManagedDatabaseState = "error"
	// Indicates service up and running
	ManagedDatabaseStateRunning ManagedDatabaseState = "running"
	// Service is stopped
	ManagedDatabaseStateStopped ManagedDatabaseState = "stopped"
	// Cleaning up service resources
	ManagedDatabaseStateCleanupService ManagedDatabaseState = "cleanup-service"
	// Cleaning up network resources
	ManagedDatabaseStateCleanupNetwork ManagedDatabaseState = "cleanup-network"
	// Deleting the service
	ManagedDatabaseStateDeleteService ManagedDatabaseState = "delete-service"
)

// ManagedDatabaseComponentRoute represents the access route a component is associated with
type ManagedDatabaseComponentRoute string

// ManagedDatabaseComponentUsage represents the logical usage for the component in question
type ManagedDatabaseComponentUsage string

const (
	// ManagedDatabaseComponentRoutePublic component can be reached over public internet
	ManagedDatabaseComponentRoutePublic ManagedDatabaseComponentRoute = "public"
	// ManagedDatabaseComponentRouteDynamic component can be only reached over utility network
	ManagedDatabaseComponentRouteDynamic ManagedDatabaseComponentRoute = "dynamic"

	// ManagedDatabaseComponentUsagePrimary component is a primary (writable) instance in a cluster
	ManagedDatabaseComponentUsagePrimary ManagedDatabaseComponentUsage = "primary"
	// ManagedDatabaseComponentUsageReplica component is a standby (read-only) instance in a cluster
	ManagedDatabaseComponentUsageReplica ManagedDatabaseComponentUsage = "replica"
)

// ManagedDatabaseNodeRole represents the role of a node implementing a service
type ManagedDatabaseNodeRole string

const (
	// ManagedDatabaseNodeRoleMaster node serves read and write requests
	ManagedDatabaseNodeRoleMaster ManagedDatabaseNodeRole = "master"
	// ManagedDatabaseNodeRoleStandby node serves read-only requests and is ready to assume the master role during failure scenario
	ManagedDatabaseNodeRoleStandby ManagedDatabaseNodeRole = "standby"
)

// ManagedDatabaseServiceType represents the type of the service. It mainly refers to the underlying database engine
// that is exposed by the service.
type ManagedDatabaseServiceType string

const (
	// ManagedDatabaseServiceTypePostgreSQL references a PostgreSQL type of database instance
	ManagedDatabaseServiceTypePostgreSQL ManagedDatabaseServiceType = "pg"
	// ManagedDatabaseServiceTypeMySQL references a MySQL type of database instance
	ManagedDatabaseServiceTypeMySQL ManagedDatabaseServiceType = "mysql"
	// Deprecated: Prefer Valkey for new key-value store instances.
	// ManagedDatabaseServiceTypeRedis references a Redis type of database instance
	ManagedDatabaseServiceTypeRedis ManagedDatabaseServiceType = "redis"
	// ManagedDatabaseServiceTypeValkey references a Valkey type of database instance
	ManagedDatabaseServiceTypeValkey ManagedDatabaseServiceType = "valkey"
	// ManagedDatabaseServiceTypeOpenSearch references an OpenSearch type of database instance
	ManagedDatabaseServiceTypeOpenSearch ManagedDatabaseServiceType = "opensearch"
)

// ManagedDatabaseUserOpenSearchAccessControlRulePermission represents a permission for user access control rule in an
// OpenSearch Managed Database service.
type ManagedDatabaseUserOpenSearchAccessControlRulePermission string

const (
	// ManagedDatabaseUserOpenSearchAccessControlRulePermissionAdmin references "admin" permission
	ManagedDatabaseUserOpenSearchAccessControlRulePermissionAdmin ManagedDatabaseUserOpenSearchAccessControlRulePermission = "admin"
	// ManagedDatabaseUserOpenSearchAccessControlRulePermissionDeny references "deny" permission
	ManagedDatabaseUserOpenSearchAccessControlRulePermissionDeny ManagedDatabaseUserOpenSearchAccessControlRulePermission = "deny"
	// ManagedDatabaseUserOpenSearchAccessControlRulePermissionReadWrite references "read-write" permission
	ManagedDatabaseUserOpenSearchAccessControlRulePermissionReadWrite ManagedDatabaseUserOpenSearchAccessControlRulePermission = "readwrite"
	// ManagedDatabaseUserOpenSearchAccessControlRulePermissionRead references "read" permission
	ManagedDatabaseUserOpenSearchAccessControlRulePermissionRead ManagedDatabaseUserOpenSearchAccessControlRulePermission = "read"
	// ManagedDatabaseUserOpenSearchAccessControlRulePermissionWrite references "write" permission
	ManagedDatabaseUserOpenSearchAccessControlRulePermissionWrite ManagedDatabaseUserOpenSearchAccessControlRulePermission = "write"
)

// ManagedDatabaseMetricPeriod represents the observation period of database metrics
type ManagedDatabaseMetricPeriod string

const (
	// ManagedDatabaseMetricPeriodHour represents the observation period of an hour for metrics request
	ManagedDatabaseMetricPeriodHour ManagedDatabaseMetricPeriod = "hour"
	// ManagedDatabaseMetricPeriodDay represents the observation period of a day for metrics request
	ManagedDatabaseMetricPeriodDay ManagedDatabaseMetricPeriod = "day"
	// ManagedDatabaseMetricPeriodWeek represents the observation period of a week for metrics request
	ManagedDatabaseMetricPeriodWeek ManagedDatabaseMetricPeriod = "week"
	// ManagedDatabaseMetricPeriodMonth represents the observation period of a month for metrics request
	ManagedDatabaseMetricPeriodMonth ManagedDatabaseMetricPeriod = "month"
	// ManagedDatabaseMetricPeriodYear represents the observation period of a year for metrics request
	ManagedDatabaseMetricPeriodYear ManagedDatabaseMetricPeriod = "year"
)

// ManagedDatabaseLogOrder represents the order the logs are queried in
type ManagedDatabaseLogOrder string

const (
	// ManagedDatabaseLogOrderAscending can be used to query logs in ascending order
	ManagedDatabaseLogOrderAscending ManagedDatabaseLogOrder = "asc"
	// ManagedDatabaseLogOrderDescending can be used to query logs in descending order
	ManagedDatabaseLogOrderDescending ManagedDatabaseLogOrder = "desc"
)

// ManagedDatabasePropertyKey represents a property name of a service
type ManagedDatabasePropertyKey string

const (
	// ManagedDatabasePropertyAutoUtilityIPFilter enables automatic ip filter generation from utility network
	// within the same zone.
	ManagedDatabasePropertyAutoUtilityIPFilter ManagedDatabasePropertyKey = "automatic_utility_network_ip_filter"
	// ManagedDatabasePropertyIPFilter allows adjusting the custom IP filter of a service. The value should
	// contain a slice of strings representing individual IP addresses or IP addresses with CIDR mask.
	// Currently, IPv4 addresses or networks are supported.
	ManagedDatabasePropertyIPFilter ManagedDatabasePropertyKey = "ip_filter"
	// ManagedDatabasePropertyPublicAccess enables public access via internet to the service. A separate public
	// endpoint DNS name will be available under Components after enabling.
	ManagedDatabasePropertyPublicAccess ManagedDatabasePropertyKey = "public_access"
	// Deprecated: ManagedDatabasePropertyMaxIndexCount in favor of ManagedDatabaseUserOpenSearchAccessControlRule.
	// Allows adjusting the maximum number of indices of an OpenSearch Managed Database service. Use ManagedDatabaseUserOpenSearchAccessControlRule instead.
	ManagedDatabasePropertyMaxIndexCount ManagedDatabasePropertyKey = "max_index_count"

	// ManagedDatabaseAllIPv4 property value can be used together with ManagedDatabasePropertyIPFilter to allow access from all
	// IPv4 hosts.
	ManagedDatabaseAllIPv4 = "0.0.0.0/0"
)

// ManagedDatabaseUserType represents the type of internal database user
type ManagedDatabaseUserType string

const (
	// ManagedDatabaseUserTypePrimary is a type of the primary user of a managed database service. There can be only
	// one primary user per service. The primary user has administrative privileges to manage logical databases and
	// users through the database's native API.
	ManagedDatabaseUserTypePrimary ManagedDatabaseUserType = "primary"
	// ManagedDatabaseUserTypeNormal is a type of normal database user of a managed database service. There can
	// be multiple normal users and the primary user can manage the privileges of these users through the database's
	// native API.
	ManagedDatabaseUserTypeNormal ManagedDatabaseUserType = "normal"
)

// ManagedDatabaseUserAuthenticationType represents the type of authentication method for an internal database user
type ManagedDatabaseUserAuthenticationType string

const (
	// ManagedDatabaseUserAuthenticationCachingSHA2Password selects "caching_sha2_password" type of authentication type.
	// This type is only supported with MySQL services.
	//nolint:gosec // this is not actually a password but an authentication type
	ManagedDatabaseUserAuthenticationCachingSHA2Password ManagedDatabaseUserAuthenticationType = "caching_sha2_password"
	// ManagedDatabaseUserAuthenticationMySQLNativePassword selects "mysql_native_password" type of authentication type.
	// This type is only supported with MySQL services.
	//nolint:gosec // this is not actually a password but an authentication type
	ManagedDatabaseUserAuthenticationMySQLNativePassword ManagedDatabaseUserAuthenticationType = "mysql_native_password"
)

// ManagedDatabaseNetwork represents a network from where database can be used
type ManagedDatabaseNetwork struct {
	Family string  `json:"family"`
	Name   string  `json:"name"`
	Type   string  `json:"type"`
	UUID   *string `json:"uuid,omitempty"`
}

// ManagedDatabase represents an existing managed database instance
type ManagedDatabase struct {
	AdditionalDiskSpaceGiB int                             `json:"additional_disk_space_gib,omitempty"`
	Backups                []ManagedDatabaseBackup         `json:"backups,omitempty"`
	Components             []ManagedDatabaseComponent      `json:"components,omitempty"`
	CreateTime             time.Time                       `json:"create_time,omitempty"`
	Labels                 []Label                         `json:"labels,omitempty"`
	Maintenance            ManagedDatabaseMaintenanceTime  `json:"maintenance,omitempty"`
	Name                   string                          `json:"name,omitempty"`
	Networks               []ManagedDatabaseNetwork        `json:"networks,omitempty"`
	NodeCount              int                             `json:"node_count,omitempty"`
	NodeStates             []ManagedDatabaseNodeState      `json:"node_states,omitempty"`
	Plan                   string                          `json:"plan,omitempty"`
	Powered                bool                            `json:"powered,omitempty"`
	Properties             ManagedDatabaseProperties       `json:"properties,omitempty"`
	State                  ManagedDatabaseState            `json:"state,omitempty"`
	StateError             map[ManagedDatabaseState]string `json:"state_error,omitempty"`
	TerminationProtection  bool                            `json:"termination_protection,omitempty"`
	Title                  string                          `json:"title,omitempty"`
	Type                   ManagedDatabaseServiceType      `json:"type,omitempty"`
	UpdateTime             time.Time                       `json:"update_time,omitempty"`
	ServiceURI             string                          `json:"service_uri,omitempty"`
	ServiceURIParams       ManagedDatabaseServiceURIParams `json:"service_uri_params,omitempty"`
	Users                  []ManagedDatabaseUser           `json:"users,omitempty"`
	UUID                   string                          `json:"uuid,omitempty"`
	Zone                   string                          `json:"zone,omitempty"`
	Metadata               *ManagedDatabaseMetadata        `json:"metadata,omitempty"`
}

// ManagedDatabaseBackup represents a full backup taken at a point in time. It should be noted that both
// MySQL and PostgreSQL support restoring to any point in time between full backups.
type ManagedDatabaseBackup struct {
	BackupName string    `json:"backup_name"`
	BackupTime time.Time `json:"backup_time"`
	DataSize   int       `json:"data_size"`
}

// ManagedDatabaseComponent represents an accessible component within a service. The usage varies between service types
type ManagedDatabaseComponent struct {
	Component string `json:"component"`
	Host      string `json:"host"`
	Port      int    `json:"port"`
	// Route describes how the component can be reached. See following:
	//	upcloud.ManagedDatabaseComponentRoutePublic
	//	upcloud.ManagedDatabaseComponentRouteDynamic
	Route ManagedDatabaseComponentRoute `json:"route"`
	// Usage describes the role of the component. See following:
	//	upcloud.ManagedDatabaseComponentUsagePrimary
	//	upcloud.ManagedDatabaseComponentUsageReplica
	Usage ManagedDatabaseComponentUsage `json:"usage"`
}

// ManagedDatabaseSessions represents sessions in the managed database instance by database instance type.
type ManagedDatabaseSessions struct {
	MySQL      []ManagedDatabaseSessionMySQL      `json:"mysql,omitempty"`
	PostgreSQL []ManagedDatabaseSessionPostgreSQL `json:"pg,omitempty"`
	// Deprecated: Redis support will be removed in favor of Valkey.
	Redis  []ManagedDatabaseSessionRedis  `json:"redis,omitempty"`
	Valkey []ManagedDatabaseSessionValkey `json:"valkey,omitempty"`
}

// ManagedDatabaseSessionMySQL represents a session in a managed MySQL database instance.
type ManagedDatabaseSessionMySQL struct {
	ApplicationName string        `json:"application_name"`
	ClientAddr      string        `json:"client_addr"`
	Datname         string        `json:"datname"`
	Id              string        `json:"id"`
	Query           string        `json:"query"`
	QueryDuration   time.Duration `json:"query_duration"`
	State           string        `json:"state"`
	Usename         string        `json:"usename"`
}

// ManagedDatabaseSessionPostgreSQL represents a session in a managed PostgreSQL database instance.
type ManagedDatabaseSessionPostgreSQL struct {
	ApplicationName string        `json:"application_name"`
	BackendStart    time.Time     `json:"backend_start"`
	BackendType     string        `json:"backend_type"`
	BackendXid      *int          `json:"backend_xid"`
	BackendXmin     *int          `json:"backend_xmin"`
	ClientAddr      string        `json:"client_addr"`
	ClientHostname  *string       `json:"client_hostname"`
	ClientPort      int           `json:"client_port"`
	Datid           int           `json:"datid"`
	Datname         string        `json:"datname"`
	Id              string        `json:"id"`
	Query           string        `json:"query"`
	QueryDuration   time.Duration `json:"query_duration"`
	QueryStart      time.Time     `json:"query_start"`
	State           string        `json:"state"`
	StateChange     time.Time     `json:"state_change"`
	Usename         string        `json:"usename"`
	Usesysid        int           `json:"usesysid"`
	WaitEvent       string        `json:"wait_event"`
	WaitEventType   string        `json:"wait_event_type"`
	XactStart       *time.Time    `json:"xact_start"`
}

// Deprecated: Redis support will be removed in favor of Valkey.
// ManagedDatabaseSessionRedis represents a session in a managed Redis database instance.
type ManagedDatabaseSessionRedis struct {
	ActiveChannelSubscriptions                int           `json:"active_channel_subscriptions"`
	ActiveDatabase                            string        `json:"active_database"`
	ActivePatternMatchingChannelSubscriptions int           `json:"active_pattern_matching_channel_subscriptions"`
	ApplicationName                           string        `json:"application_name"`
	ClientAddr                                string        `json:"client_addr"`
	ConnectionAge                             time.Duration `json:"connection_age"`
	ConnectionIdle                            time.Duration `json:"connection_idle"`
	Flags                                     []string      `json:"flags"`
	FlagsRaw                                  string        `json:"flags_raw"`
	Id                                        string        `json:"id"`
	MultiExecCommands                         int           `json:"multi_exec_commands"`
	OutputBuffer                              int           `json:"output_buffer"`
	OutputBufferMemory                        int           `json:"output_buffer_memory"`
	OutputListLength                          int           `json:"output_list_length"`
	Query                                     string        `json:"query"`
	QueryBuffer                               int           `json:"query_buffer"`
	QueryBufferFree                           int           `json:"query_buffer_free"`
}

// ManagedDatabaseSessionValkey represents a session in a managed Valkey database instance.
type ManagedDatabaseSessionValkey struct {
	ActiveChannelSubscriptions                int           `json:"active_channel_subscriptions"`
	ActiveDatabase                            string        `json:"active_database"`
	ActivePatternMatchingChannelSubscriptions int           `json:"active_pattern_matching_channel_subscriptions"`
	ApplicationName                           string        `json:"application_name"`
	ClientAddr                                string        `json:"client_addr"`
	ConnectionAge                             time.Duration `json:"connection_age"`
	ConnectionIdle                            time.Duration `json:"connection_idle"`
	Flags                                     []string      `json:"flags"`
	FlagsRaw                                  string        `json:"flags_raw"`
	Id                                        string        `json:"id"`
	MultiExecCommands                         int           `json:"multi_exec_commands"`
	OutputBuffer                              int           `json:"output_buffer"`
	OutputBufferMemory                        int           `json:"output_buffer_memory"`
	OutputListLength                          int           `json:"output_list_length"`
	Query                                     string        `json:"query"`
	QueryBuffer                               int           `json:"query_buffer"`
	QueryBufferFree                           int           `json:"query_buffer_free"`
}

// ManagedDatabaseMaintenanceTime represents the set time of week when automatic maintenance operations are allowed
type ManagedDatabaseMaintenanceTime struct {
	DayOfWeek string `json:"dow"`
	Time      string `json:"time"`
}

// ManagedDatabaseMetrics represents managed database service metrics
//
// Metrics are represented in chart form containing a set of columns and two-dimensional slice of rows.
// The inner slice index corresponds a column. If the service consists of multiple nodes, each node gets their
// own column in the chart.
//
// The first column is always a timestamp which denotes the timestamp for the recorded metric
type ManagedDatabaseMetrics struct {
	CPUUsage       ManagedDatabaseMetricsChartFloat64 `json:"cpu_usage"`
	DiskUsage      ManagedDatabaseMetricsChartFloat64 `json:"disk_usage"`
	DiskIOReads    ManagedDatabaseMetricsChartInt     `json:"diskio_reads"`
	DiskIOWrite    ManagedDatabaseMetricsChartInt     `json:"diskio_writes"`
	LoadAverage    ManagedDatabaseMetricsChartFloat64 `json:"load_average"`
	MemoryUsage    ManagedDatabaseMetricsChartFloat64 `json:"mem_usage"`
	NetworkReceive ManagedDatabaseMetricsChartInt     `json:"net_receive"`
	NetworkSend    ManagedDatabaseMetricsChartInt     `json:"net_send"`
}

// ManagedDatabaseMetricsChartFloat64 represents a metric chart with float64 row values
type ManagedDatabaseMetricsChartFloat64 struct {
	ManagedDatabaseMetricsChartHeader

	// Rows contains a slice of values per row. The inner slice has the same indexing as the
	// ManagedDatabaseMetricsChartHeader.Columns
	Rows [][]float64
}

// MarshalJSON implements json.Marshaler for the canonical form of metrics chart
func (m ManagedDatabaseMetricsChartFloat64) MarshalJSON() ([]byte, error) {
	var chart canonicalManagedDatabaseMetricChartForMarshal
	chart.Data.Columns = make([]ManagedDatabaseMetricsColumn, 0, len(m.Columns)+1)
	chart.Data.Columns = append(chart.Data.Columns, ManagedDatabaseMetricsColumn{Label: "time", Type: "date"})
	chart.Data.Columns = append(chart.Data.Columns, m.Columns...)
	chart.Hints.Title = m.Title
	if len(m.Timestamps) != len(m.Rows) {
		return nil, fmt.Errorf("the number of timestamps doesn't match the number of rows: %d != %d",
			len(m.Timestamps), len(m.Rows))
	}
	rows := make([][]any, 0, len(m.Rows))
	for i := range m.Rows {
		row := make([]any, 0, len(m.Rows[i]))
		row = append(row, &m.Timestamps[i])
		if len(m.Rows[i]) != len(m.Columns) {
			return nil, fmt.Errorf("unexpected number of columns at row %d (not %d)", i, len(m.Columns))
		}
		for j := range m.Rows[i] {
			row = append(row, &m.Rows[i][j])
		}
		rows = append(rows, row)
	}
	chart.Data.Rows = rows
	return json.Marshal(&chart)
}

// UnmarshalJSON implements json.Unmarshaler for the canonical form of metrics chart
func (m *ManagedDatabaseMetricsChartFloat64) UnmarshalJSON(d []byte) error {
	var chart canonicalManagedDatabaseMetricChartForUnmarshal
	if err := json.Unmarshal(d, &chart); err != nil {
		return err
	}
	timestamps := make([]time.Time, 0, len(chart.Data.Rows))
	rows := make([][]float64, 0, len(chart.Data.Rows))
	inRows := chart.Data.Rows
	for i := range inRows {
		row := make([]float64, 0, len(inRows)-1)
		var ts time.Time
		if len(inRows[i]) != len(chart.Data.Columns) {
			return fmt.Errorf("unexpected number of columns at row %d (not %d)", i, len(m.Columns))
		}
		if err := json.Unmarshal(inRows[i][0], &ts); err != nil {
			return fmt.Errorf("cannot unmarshal timestamp at row %d: %w", i, err)
		}
		timestamps = append(timestamps, ts)
		for j := 1; j < len(inRows[i]); j++ {
			var iv float64
			if err := json.Unmarshal(inRows[i][j], &iv); err != nil {
				return fmt.Errorf("cannot unmarshal value at row %d, column %d: %w", i, j, err)
			}
			row = append(row, iv)
		}
		rows = append(rows, row)
	}
	m.Timestamps = timestamps
	m.Columns = chart.Data.Columns[1:]
	m.Title = chart.Hints.Title
	m.Rows = rows
	return nil
}

// ManagedDatabaseMetricsChartHeader represents common fields of a metrics chart
type ManagedDatabaseMetricsChartHeader struct {
	// Columns contains a set of columns that describe for what node the corresponding row element belongs to
	// as well as the type of the metric value
	Columns []ManagedDatabaseMetricsColumn
	// Timestamps contains the timestamps of the rows. Its indexing corresponds the rows.
	Timestamps []time.Time
	// Title contains a description of the metrics chart
	Title string
}

type canonicalManagedDatabaseMetricChartForMarshal struct {
	Data struct {
		Columns []ManagedDatabaseMetricsColumn `json:"cols"`
		Rows    [][]any                        `json:"rows"`
	} `json:"data"`
	Hints struct {
		Title string `json:"title"`
	} `json:"hints"`
}

type canonicalManagedDatabaseMetricChartForUnmarshal struct {
	Data struct {
		Columns []ManagedDatabaseMetricsColumn `json:"cols"`
		Rows    [][]json.RawMessage            `json:"rows"`
	} `json:"data"`
	Hints struct {
		Title string `json:"title"`
	} `json:"hints"`
}

// ManagedDatabaseMetricsChartInt represents a metric chart with int row values
type ManagedDatabaseMetricsChartInt struct {
	ManagedDatabaseMetricsChartHeader

	// Rows contains a slice of values per row. The inner slice has the same indexing as the
	// ManagedDatabaseMetricsChartHeader.Columns
	Rows [][]int
}

// MarshalJSON implements json.Marshaler for the canonical form of metrics chart
func (m ManagedDatabaseMetricsChartInt) MarshalJSON() ([]byte, error) {
	var chart canonicalManagedDatabaseMetricChartForMarshal
	chart.Data.Columns = make([]ManagedDatabaseMetricsColumn, 0, len(m.Columns)+1)
	chart.Data.Columns = append(chart.Data.Columns, ManagedDatabaseMetricsColumn{Label: "time", Type: "date"})
	chart.Data.Columns = append(chart.Data.Columns, m.Columns...)
	chart.Hints.Title = m.Title
	if len(m.Timestamps) != len(m.Rows) {
		return nil, fmt.Errorf("the number of timestamps doesn't match the number of rows: %d != %d",
			len(m.Timestamps), len(m.Rows))
	}
	rows := make([][]any, 0, len(m.Rows))
	for i := range m.Rows {
		row := make([]any, 0, len(m.Rows[i]))
		row = append(row, &m.Timestamps[i])
		if len(m.Rows[i]) != len(m.Columns) {
			return nil, fmt.Errorf("unexpected number of columns at row %d (not %d)", i, len(m.Columns))
		}
		for j := range m.Rows[i] {
			row = append(row, &m.Rows[i][j])
		}
		rows = append(rows, row)
	}
	chart.Data.Rows = rows
	return json.Marshal(&chart)
}

// UnmarshalJSON implements json.Unmarshaler for the canonical form of metrics chart
func (m *ManagedDatabaseMetricsChartInt) UnmarshalJSON(d []byte) error {
	var chart canonicalManagedDatabaseMetricChartForUnmarshal
	if err := json.Unmarshal(d, &chart); err != nil {
		return err
	}
	timestamps := make([]time.Time, 0, len(chart.Data.Rows))
	rows := make([][]int, 0, len(chart.Data.Rows))
	inRows := chart.Data.Rows
	for i := range inRows {
		row := make([]int, 0, len(inRows)-1)
		var ts time.Time
		if len(inRows[i]) != len(chart.Data.Columns) {
			return fmt.Errorf("unexpected number of columns at row %d (not %d)", i, len(m.Columns))
		}
		if err := json.Unmarshal(inRows[i][0], &ts); err != nil {
			return fmt.Errorf("cannot unmarshal timestamp at row %d: %w", i, err)
		}
		timestamps = append(timestamps, ts)
		for j := 1; j < len(inRows[i]); j++ {
			var iv int
			if err := json.Unmarshal(inRows[i][j], &iv); err != nil {
				return fmt.Errorf("cannot unmarshal value at row %d, column %d: %w", i, j, err)
			}
			row = append(row, iv)
		}
		rows = append(rows, row)
	}
	m.Timestamps = timestamps
	m.Columns = chart.Data.Columns[1:]
	m.Title = chart.Hints.Title
	m.Rows = rows
	return nil
}

// ManagedDatabaseMetricsColumn represents a single column of a metrics chart
type ManagedDatabaseMetricsColumn struct {
	// Label describes the usage of chart's column
	Label string `json:"label"`
	// Type describes the type of values in chart's column
	Type string `json:"type"`
}

// ManagedDatabaseNodeState represents a database node that is part of the service instance
type ManagedDatabaseNodeState struct {
	// Name field is same as the ManagedDatabase.Name plus a dash plus an index value. The index represents the
	// generation of a node. Certain modifications require re-provisioning of a node.
	Name string `json:"name"`
	// Role represents the role of a node
	Role ManagedDatabaseNodeRole `json:"role,omitempty"`
	// State represents the current state of a node
	State string `json:"state"`
}

// ManagedDatabaseProperties is a Properties helper type for ManagedDatabase
type ManagedDatabaseProperties map[ManagedDatabasePropertyKey]any

// Get returns a property value by name. The underlying map is initialised if it's nil
func (m *ManagedDatabaseProperties) Get(name ManagedDatabasePropertyKey) any {
	if m == nil {
		return nil
	}
	if *m == nil {
		*m = make(ManagedDatabaseProperties)
	}
	return (*m)[name]
}

// GetString returns a string property value.
// The underlying map is initialised if it's nil.
func (m *ManagedDatabaseProperties) GetString(name ManagedDatabasePropertyKey) (string, error) {
	if v, ok := m.Get(name).(string); ok {
		return v, nil
	}
	return "", fmt.Errorf("not a string property %q", name)
}

// GetInt returns an integer property value.
// The underlying map is initialised if it's nil.
func (m *ManagedDatabaseProperties) GetInt(name ManagedDatabasePropertyKey) (int, error) {
	if v, ok := m.Get(name).(int); ok {
		return v, nil
	}
	return 0, fmt.Errorf("not an int property %q", name)
}

// GetBool returns a boolean property value.
// The underlying map is initialised if it's nil.
func (m *ManagedDatabaseProperties) GetBool(name ManagedDatabasePropertyKey) (bool, error) {
	if v, ok := m.Get(name).(bool); ok {
		return v, nil
	}
	return false, fmt.Errorf("not a boolean property %q", name)
}

// GetStringSlice returns a string-slice property value.
// The underlying map is initialised if it's nil.
func (m *ManagedDatabaseProperties) GetStringSlice(name ManagedDatabasePropertyKey) ([]string, error) {
	if slice, ok := m.Get(name).([]string); ok {
		return slice, nil
	}
	slice, ok := m.Get(name).([]any)
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
func (m *ManagedDatabaseProperties) GetAutoUtilityIPFilter() bool {
	v, _ := m.GetBool(ManagedDatabasePropertyAutoUtilityIPFilter)
	return v
}

// GetIPFilter returns a slice of allowed hosts or networks.
// See upcloud.ManagedDatabasePropertyIPFilter for more information.
func (m *ManagedDatabaseProperties) GetIPFilter() []string {
	v, _ := m.GetStringSlice(ManagedDatabasePropertyIPFilter)
	return v
}

// GetPublicAccess returns the state of public access to the service.
// See upcloud.ManagedDatabasePropertyPublicAccess for more information.
func (m *ManagedDatabaseProperties) GetPublicAccess() bool {
	v, _ := m.GetBool(ManagedDatabasePropertyPublicAccess)
	return v
}

// Deprecated: GetMaxIndexCount is deprecated in favor of field upcloud.ManagedDatabasePropertyMaxIndexCount.
// Returns the maximum index count of the service.
func (m *ManagedDatabaseProperties) GetMaxIndexCount() int {
	v, _ := m.GetInt(ManagedDatabasePropertyMaxIndexCount)
	return v
}

type ManagedDatabaseLogs struct {
	// Offset describes the next available offset. Use this to query more logs.
	Offset string `json:"offset"`
	Logs   []ManagedDatabaseLogEntry
}

type ManagedDatabaseLogEntry struct {
	Hostname string    `json:"hostname"`
	Message  string    `json:"msg"`
	Time     time.Time `json:"time"`
	Service  string    `json:"service"`
}

// ManagedDatabaseLogicalDatabase represents a logical database inside a managed database service instance
type ManagedDatabaseLogicalDatabase struct {
	Name string `json:"name"`
	// LCCollate represents a default string sort order of a logical database
	LCCollate string `json:"lc_collate"`
	// LCCType represents a default character classification of a logical database
	LCCType string `json:"lc_ctype"`
}

// ManagedDatabaseServiceURIParams represents individual components of ServiceURI field
type ManagedDatabaseServiceURIParams struct {
	DatabaseName string `json:"dbname"`
	Host         string `json:"host"`
	Password     string `json:"password"`
	Port         string `json:"port"`
	SSLMode      string `json:"ssl_mode"`
	User         string `json:"user"`
}

// ManagedDatabaseUser represents a database internal user
type ManagedDatabaseUser struct {
	// Authentication field represents an allowed authentication type for this user. For more information see:
	// 	upcloud.ManagedDatabaseUserAuthenticationCachingSHA2Password
	// 	upcloud.ManagedDatabaseUserAuthenticationMySQLNativePassword
	Authentication ManagedDatabaseUserAuthenticationType `json:"authentication,omitempty"`
	Type           ManagedDatabaseUserType               `json:"type,omitempty"`
	// Password field is only visible when querying an individual user. It is omitted in main service view and in
	// get all users view.
	Password        string                              `json:"password,omitempty"`
	Username        string                              `json:"username,omitempty"`
	PGAccessControl *ManagedDatabaseUserPGAccessControl `json:"pg_access_control,omitempty"`
	// Deprecated: Redis support will be removed in favor of Valkey.
	RedisAccessControl      *ManagedDatabaseUserRedisAccessControl      `json:"redis_access_control,omitempty"`
	ValkeyAccessControl     *ManagedDatabaseUserValkeyAccessControl     `json:"valkey_access_control,omitempty"`
	OpenSearchAccessControl *ManagedDatabaseUserOpenSearchAccessControl `json:"opensearch_access_control,omitempty"`
}

type ManagedDatabaseUserPGAccessControl struct {
	AllowReplication *bool `json:"allow_replication,omitempty"`
}

// Deprecated: Redis support will be removed in favor of Valkey.
type ManagedDatabaseUserRedisAccessControl struct {
	Categories *[]string `json:"categories,omitempty"`
	Channels   *[]string `json:"channels,omitempty"`
	Commands   *[]string `json:"commands,omitempty"`
	Keys       *[]string `json:"keys,omitempty"`
}

type ManagedDatabaseUserValkeyAccessControl struct {
	Categories *[]string `json:"categories,omitempty"`
	Channels   *[]string `json:"channels,omitempty"`
	Commands   *[]string `json:"commands,omitempty"`
	Keys       *[]string `json:"keys,omitempty"`
}

type ManagedDatabaseUserOpenSearchAccessControl struct {
	Rules *[]ManagedDatabaseUserOpenSearchAccessControlRule `json:"rules,omitempty"`
}

type ManagedDatabaseUserOpenSearchAccessControlRule struct {
	Index      string                                                   `json:"index"`
	Permission ManagedDatabaseUserOpenSearchAccessControlRulePermission `json:"permission"`
}

// ManagedDatabaseQueryStatisticsMySQL represents statistics reported by a MySQL server.
// Statistics are per Digest which is derived from DigestText
type ManagedDatabaseQueryStatisticsMySQL struct {
	AvgTimerWait            time.Duration `json:"avg_timer_wait"`
	CountStar               uint64        `json:"count_star"`
	Digest                  string        `json:"digest"`
	DigestText              string        `json:"digest_text"`
	FirstSeen               time.Time     `json:"first_seen"`
	LastSeen                time.Time     `json:"last_seen"`
	MaxTimerWait            time.Duration `json:"max_timer_wait"`
	MinTimerWait            time.Duration `json:"min_timer_wait"`
	Quantile95              time.Duration `json:"quantile_95"`
	Quantile99              time.Duration `json:"quantile_99"`
	Quantile999             time.Duration `json:"quantile_999"`
	QuerySampleSeen         time.Time     `json:"query_sample_seen"`
	QuerySampleText         string        `json:"query_sample_text"`
	QuerySampleTimerWait    time.Duration `json:"query_sample_timer_wait"`
	SchemaName              string        `json:"schema_name"`
	SumCreatedTmpDiskTables uint64        `json:"sum_created_tmp_disk_tables"`
	SumCreatedTmpTables     uint64        `json:"sum_created_tmp_tables"`
	SumErrors               uint64        `json:"sum_errors"`
	SumLockTime             uint64        `json:"sum_lock_time"`
	SumNoGoodIndexUsed      uint64        `json:"sum_no_good_index_used"`
	SumNoIndexUsed          uint64        `json:"sum_no_index_used"`
	SumRowsAffected         uint64        `json:"sum_rows_affected"`
	SumRowsExamined         uint64        `json:"sum_rows_examined"`
	SumRowsSent             uint64        `json:"sum_rows_sent"`
	SumSelectFullJoin       uint64        `json:"sum_select_full_join"`
	SumSelectFullRangeJoin  uint64        `json:"sum_select_full_range_join"`
	SumSelectRange          uint64        `json:"sum_select_range"`
	SumSelectRangeCheck     uint64        `json:"sum_select_range_check"`
	SumSelectScan           uint64        `json:"sum_select_scan"`
	SumSortMergePasses      uint64        `json:"sum_sort_merge_passes"`
	SumSortRange            uint64        `json:"sum_sort_range"`
	SumSortRows             uint64        `json:"sum_sort_rows"`
	SumSortScan             uint64        `json:"sum_sort_scan"`
	SumTimerWait            time.Duration `json:"sum_timer_wait"`
	SumWarnings             uint64        `json:"sum_warnings"`
}

// ManagedDatabaseQueryStatisticsPostgreSQL represents statistics reported by a PostgreSQL server.
// Statistics are per Query without parameters.
type ManagedDatabaseQueryStatisticsPostgreSQL struct {
	BlockReadTime       time.Duration `json:"blk_read_time"`
	BlockWriteTime      time.Duration `json:"blk_write_time"`
	Calls               uint64        `json:"calls"`
	DatabaseName        string        `json:"database_name"`
	LocalBlocksDirtied  uint64        `json:"local_blks_dirtied"`
	LocalBlocksHit      uint64        `json:"local_blks_hit"`
	LocalBlocksRead     uint64        `json:"local_blks_read"`
	LocalBlocksWritten  uint64        `json:"local_blks_written"`
	MaxTime             time.Duration `json:"max_time"`
	MeanTime            time.Duration `json:"mean_time"`
	MinTime             time.Duration `json:"min_time"`
	Query               string        `json:"query"`
	Rows                uint64        `json:"rows"`
	SharedBlocksDirtied uint64        `json:"shared_blks_dirtied"`
	SharedBlocksHit     uint64        `json:"shared_blks_hit"`
	SharedBlocksRead    uint64        `json:"shared_blks_read"`
	SharedBlocksWritten uint64        `json:"shared_blks_written"`
	StddevTime          time.Duration `json:"stddev_time"`
	TempBlocksRead      uint64        `json:"temp_blks_read"`
	TempBlocksWritten   uint64        `json:"temp_blks_written"`
	TotalTime           time.Duration `json:"total_time"`
	UserName            string        `json:"user_name"`
}

// ManagedDatabaseType represents details of a database service type.
type ManagedDatabaseType struct {
	Name                   string                                    `json:"name"`
	Description            string                                    `json:"description"`
	LatestAvailableVersion string                                    `json:"latest_available_version"`
	ServicePlans           []ManagedDatabaseServicePlan              `json:"service_plans"`
	Properties             map[string]ManagedDatabaseServiceProperty `json:"properties"`
}

// ManagedDatabaseServicePlan represents details of a database service plan.
type ManagedDatabaseServicePlan struct {
	BackupConfig           ManagedDatabaseBackupConfig            `json:"backup_config"`
	BackupConfigMySQL      *ManagedDatabaseBackupConfigMySQL      `json:"backup_config_mysql,omitempty"`
	BackupConfigOpenSearch *ManagedDatabaseBackupConfigOpenSearch `json:"backup_config_opensearch,omitempty"`
	BackupConfigPostgreSQL *ManagedDatabaseBackupConfigPostgreSQL `json:"backup_config_pg,omitempty"`
	// Deprecated: Redis support will be removed in favor of Valkey.
	BackupConfigRedis  *ManagedDatabaseBackupConfigRedis  `json:"backup_config_redis,omitempty"`
	BackupConfigValkey *ManagedDatabaseBackupConfigValkey `json:"backup_config_valkey,omitempty"`
	NodeCount          int                                `json:"node_count"`
	Plan               string                             `json:"plan"`
	CoreNumber         int                                `json:"core_number"`
	StorageSize        int                                `json:"storage_size"`
	MemoryAmount       int                                `json:"memory_amount"`
	Zones              ManagedDatabaseServicePlanZones    `json:"zones"`
}

// Deprecated: ManagedDatabaseBackupConfig is deprecated in favor of service specific ManagedDatabaseBackupConfig<ServiceType> types.
// Represents backup configuration of a database service plan.
type ManagedDatabaseBackupConfig struct {
	Interval     int    `json:"interval"`
	MaxCount     int    `json:"max_count"`
	RecoveryMode string `json:"recovery_mode"`
}

// ManagedDatabaseBackupConfigMySQL represents backup configuration of a MySQL database service plan
type ManagedDatabaseBackupConfigMySQL struct {
	Interval     int    `json:"interval"`
	MaxCount     int    `json:"max_count"`
	RecoveryMode string `json:"recovery_mode"`
}

// ManagedDatabaseBackupConfigOpenSearch represents backup configuration of an OpenSearch database service plan.
type ManagedDatabaseBackupConfigOpenSearch struct {
	FrequentIntervalMinutes    int    `json:"frequent_interval_minutes"`
	FrequentOldestAgeMinutes   int    `json:"frequent_oldest_age_minutes"`
	InfrequentIntervalMinutes  int    `json:"infrequent_interval_minutes"`
	InfrequentOldestAgeMinutes int    `json:"infrequent_oldest_age_minutes"`
	RecoveryMode               string `json:"recovery_mode"`
}

// ManagedDatabaseBackupConfigPostgreSQL represents backup configuration of a PostgreSQL database service plan
type ManagedDatabaseBackupConfigPostgreSQL struct {
	Interval     int    `json:"interval"`
	MaxCount     int    `json:"max_count"`
	RecoveryMode string `json:"recovery_mode"`
}

// Deprecated: Redis support will be removed in favor of Valkey.
// ManagedDatabaseBackupConfigRedis represents backup configuration of a Redis database service plan
type ManagedDatabaseBackupConfigRedis struct {
	Interval     int    `json:"interval"`
	MaxCount     int    `json:"max_count"`
	RecoveryMode string `json:"recovery_mode"`
}

// ManagedDatabaseBackupConfigValkey represents backup configuration of a Valkey database service plan
type ManagedDatabaseBackupConfigValkey struct {
	Interval     int    `json:"interval"`
	MaxCount     int    `json:"max_count"`
	RecoveryMode string `json:"recovery_mode"`
}

// ManagedDatabaseServicePlanZones is a helper for unmarshalling database plan zones.
type ManagedDatabaseServicePlanZones []ManagedDatabaseServicePlanZone

// ManagedDatabaseServicePlanZone represents zone where parent database plan is available
type ManagedDatabaseServicePlanZone struct {
	Name string `json:"name"`
}

// UnmarshalJSON is a custom unmarshaller that deals with deeply embedded values.
func (s *ManagedDatabaseServicePlanZones) UnmarshalJSON(b []byte) error {
	v := struct {
		Zones []ManagedDatabaseServicePlanZone `json:"zone"`
	}{}
	err := json.Unmarshal(b, &v)
	if err != nil {
		return err
	}

	*s = v.Zones

	return nil
}

// ManagedDatabaseServiceProperty contains help for database property usage and validation
type ManagedDatabaseServiceProperty struct {
	CreateOnly  bool                                      `json:"createOnly,omitempty"`
	Default     any                                       `json:"default,omitempty"`
	Example     any                                       `json:"example,omitempty"`
	MaxLength   int                                       `json:"maxLength,omitempty"`
	Maximum     *float64                                  `json:"maximum,omitempty"`
	MinLength   int                                       `json:"minLength,omitempty"`
	Minimum     *float64                                  `json:"minimum,omitempty"`
	Pattern     string                                    `json:"pattern,omitempty"`
	Type        any                                       `json:"type"`
	Title       string                                    `json:"title"`
	Description string                                    `json:"description,omitempty"`
	Enum        any                                       `json:"enum,omitempty"`
	UserError   string                                    `json:"user_error,omitempty"`
	Properties  map[string]ManagedDatabaseServiceProperty `json:"properties,omitempty"`
}

// ManagedDatabaseMetadata contains additional read-only informational data about the managed database
type ManagedDatabaseMetadata struct {
	MaxConnections int    `json:"max_connections,omitempty"`
	PGVersion      string `json:"pg_version,omitempty"`
	MySQLVersion   string `json:"mysql_version,omitempty"`
	// Deprecated: Redis support will be removed in favor of Valkey.
	RedisVersion                string `json:"redis_version,omitempty"`
	ValkeyVersion               string `json:"valkey_version,omitempty"`
	WriteBlockThresholdExceeded *bool  `json:"write_block_threshold_exceeded,omitempty"`
	OpenSearchVersion           string `json:"opensearch_version,omitempty"`
	UpgradeVersion              string `json:"upgrade_version,omitempty"`
}

// ManagedDatabaseIndex represents an index of an OpenSearch Managed Database
type ManagedDatabaseIndex struct {
	CreateTime          time.Time `json:"create_time"`
	Docs                int       `json:"docs"`
	Health              string    `json:"health"`
	IndexName           string    `json:"index_name"`
	NumberOfReplicas    int       `json:"number_of_replicas"`
	NumberOfShards      int       `json:"number_of_shards"`
	ReadOnlyAllowDelete bool      `json:"read_only_allow_delete"`
	Size                int       `json:"size"`
	Status              string    `json:"status"`
}

// ManagedDatabaseAccessControl contains access controls settings for an OpenSearch Managed Database service
type ManagedDatabaseAccessControl struct {
	ACLsEnabled         *bool `json:"access_control"`
	ExtendedACLsEnabled *bool `json:"extended_access_control"`
}
