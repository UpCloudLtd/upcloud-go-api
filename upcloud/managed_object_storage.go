package upcloud

import (
	"time"
)

const (
	// ManagedObjectStorageConfiguredStatusStarted indicates that service is running
	ManagedObjectStorageConfiguredStatusStarted ManagedObjectStorageConfiguredStatus = "started"
	// ManagedObjectStorageConfiguredStatusStopped indicates that service is stopped
	ManagedObjectStorageConfiguredStatusStopped ManagedObjectStorageConfiguredStatus = "stopped"
)

const (
	// ManagedObjectStorageOperationalStateDeleteDNS indicates that DNS records are being removed
	ManagedObjectStorageOperationalStateDeleteDNS ManagedObjectStorageOperationalState = "delete-dns"
	// ManagedObjectStorageOperationalStateDeleteNetwork indicates that network is being reconfigured
	ManagedObjectStorageOperationalStateDeleteNetwork ManagedObjectStorageOperationalState = "delete-network"
	// ManagedObjectStorageOperationalStateDeleteService indicates that service is being deleted
	ManagedObjectStorageOperationalStateDeleteService ManagedObjectStorageOperationalState = "delete-service"
	// ManagedObjectStorageOperationalStateDeleteUser indicates that users are being deleted
	ManagedObjectStorageOperationalStateDeleteUser ManagedObjectStorageOperationalState = "delete-user"
	// ManagedObjectStorageOperationalStatePending indicates newly created service or that started reconfiguration
	ManagedObjectStorageOperationalStatePending ManagedObjectStorageOperationalState = "started"
	// ManagedObjectStorageOperationalStateRunning indicates that service is up and running
	ManagedObjectStorageOperationalStateRunning ManagedObjectStorageOperationalState = "running"
	// ManagedObjectStorageOperationalStateSetupCheckup indicates that service configuration and state are being verified
	ManagedObjectStorageOperationalStateSetupCheckup ManagedObjectStorageOperationalState = "setup-checkup"
	// ManagedObjectStorageOperationalStateSetupDNS indicates that DNS records are being updated
	ManagedObjectStorageOperationalStateSetupDNS ManagedObjectStorageOperationalState = "setup-dns"
	// ManagedObjectStorageOperationalStateSetupNetwork indicates that network is being configured
	ManagedObjectStorageOperationalStateSetupNetwork ManagedObjectStorageOperationalState = "setup-network"
	// ManagedObjectStorageOperationalStateSetupService indicates that service is being configured
	ManagedObjectStorageOperationalStateSetupService ManagedObjectStorageOperationalState = "setup-service"
	// ManagedObjectStorageOperationalStateSetupUser indicates that users are being configured
	ManagedObjectStorageOperationalStateSetupUser ManagedObjectStorageOperationalState = "setup-user"
	// ManagedObjectStorageOperationalStateStopped indicates that service is down
	ManagedObjectStorageOperationalStateStopped ManagedObjectStorageOperationalState = "stopped"
)

const (
	// ManagedObjectStorageUserOperationalStatePending indicates a newly attached user
	ManagedObjectStorageUserOperationalStatePending ManagedObjectStorageUserOperationalState = "pending"
	// ManagedObjectStorageUserOperationalStateReady indicates that the user is configured and ready for access keys issuing
	ManagedObjectStorageUserOperationalStateReady ManagedObjectStorageUserOperationalState = "ready"
)

type (
	// ManagedObjectStorageConfiguredStatus indicates the service's current intended status. Managed by the customer
	ManagedObjectStorageConfiguredStatus string
	// ManagedObjectStorageOperationalState indicates the service's current operational, effective state. Managed by the system
	ManagedObjectStorageOperationalState string
	// ManagedObjectStorageUserOperationalState indicates the user's current operational, effective state. Managed by the system
	ManagedObjectStorageUserOperationalState string
)

// ManagedObjectStorage represents a Managed Object Storage service
type ManagedObjectStorage struct {
	ConfiguredStatus ManagedObjectStorageConfiguredStatus `json:"configured_status"`
	CreatedAt        time.Time                            `json:"created_at"`
	Endpoints        []ManagedObjectStorageEndpoint       `json:"endpoints"`
	Labels           []Label                              `json:"labels"`
	Networks         []ManagedObjectStorageNetwork        `json:"networks"`
	OperationalState ManagedObjectStorageOperationalState `json:"operational_state"`
	Region           string                               `json:"region"`
	UpdatedAt        time.Time                            `json:"updated_at"`
	Users            []ManagedObjectStorageUser           `json:"users"`
	UUID             string                               `json:"uuid"`
}

// ManagedObjectStorageEndpoint represents an endpoint for accessing the Managed Object Storage service
type ManagedObjectStorageEndpoint struct {
	DomainName string `json:"domain_name"`
	Type       string `json:"type"`
}

// ManagedObjectStorageNetwork represents a network from where object storage can be used. Private networks must reside in object storage region
type ManagedObjectStorageNetwork struct {
	Family string  `json:"family"`
	Name   string  `json:"name"`
	Type   string  `json:"type"`
	UUID   *string `json:"uuid,omitempty"`
}

// ManagedObjectStorageUser represents a user for the Managed Object Storage service
type ManagedObjectStorageUser struct {
	AccessKeys       []ManagedObjectStorageUserAccessKey      `json:"access_keys"`
	CreatedAt        time.Time                                `json:"created_at"`
	OperationalState ManagedObjectStorageUserOperationalState `json:"operational_state"`
	UpdatedAt        time.Time                                `json:"updated_at"`
	Username         string                                   `json:"username"`
}

// ManagedObjectStorageRegion represents a region where Managed Object Storage service can be hosted
type ManagedObjectStorageRegion struct {
	Name        string                           `json:"name"`
	PrimaryZone string                           `json:"primary_zone"`
	Zones       []ManagedObjectStorageRegionZone `json:"zones"`
}

// ManagedObjectStorageRegionZone represents a zone within the Managed Object Storage service region
type ManagedObjectStorageRegionZone struct {
	Name string `json:"name"`
}

// ManagedObjectStorageUserAccessKey represents Access Key details for a Managed Object Storage service user
type ManagedObjectStorageUserAccessKey struct {
	AccessKeyId     string    `json:"access_key_id"`
	CreatedAt       time.Time `json:"created_at"`
	Enabled         bool      `json:"enabled"`
	LastUsedAt      time.Time `json:"last_used_at"`
	Name            string    `json:"name"`
	SecretAccessKey *string   `json:"secret_access_key,omitempty"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// ManagedObjectStorageBucketMetrics represents metrics for a Managed Object Storage service bucket
type ManagedObjectStorageBucketMetrics struct {
	Name           string `json:"name"`
	TotalObjects   int    `json:"total_objects"`
	TotalSizeBytes int    `json:"total_size_bytes"`
}

// ManagedObjectStorageMetrics represents metrics for a Managed Object Storage service
type ManagedObjectStorageMetrics struct {
	TotalObjects   int `json:"total_objects"`
	TotalSizeBytes int `json:"total_size_bytes"`
}
