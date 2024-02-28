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
	// ManagedObjectStorageOperationalStateStopped indicates that service is down
	ManagedObjectStorageOperationalStateStopped ManagedObjectStorageOperationalState = "stopped"
)

const (
	// ManagedObjectStorageUserAccessKeyStatusActive indicates an active access key
	ManagedObjectStorageUserAccessKeyStatusActive ManagedObjectStorageUserAccessKeyStatus = "Active"
	// ManagedObjectStorageUserAccessKeyStatusInactive indicates an inactive access key
	ManagedObjectStorageUserAccessKeyStatusInactive ManagedObjectStorageUserAccessKeyStatus = "Inactive"
)

type (
	// ManagedObjectStorageConfiguredStatus indicates the service's current intended status. Managed by the customer
	ManagedObjectStorageConfiguredStatus string
	// ManagedObjectStorageOperationalState indicates the service's current operational, effective state. Managed by the system
	ManagedObjectStorageOperationalState string
	// ManagedObjectStorageUserAccessKeyStatus indicates the access key's current status. Managed by the customer
	ManagedObjectStorageUserAccessKeyStatus string
)

// ManagedObjectStorage represents a Managed Object Storage service
type ManagedObjectStorage struct {
	ConfiguredStatus ManagedObjectStorageConfiguredStatus `json:"configured_status"`
	CreatedAt        time.Time                            `json:"created_at"`
	Endpoints        []ManagedObjectStorageEndpoint       `json:"endpoints"`
	Labels           []Label                              `json:"labels"`
	Name             string                               `json:"name,omitempty"`
	Networks         []ManagedObjectStorageNetwork        `json:"networks"`
	OperationalState ManagedObjectStorageOperationalState `json:"operational_state"`
	Region           string                               `json:"region"`
	UpdatedAt        time.Time                            `json:"updated_at"`
	UUID             string                               `json:"uuid"`
}

// ManagedObjectStorageEndpoint represents an endpoint for accessing the Managed Object Storage service
type ManagedObjectStorageEndpoint struct {
	DomainName string `json:"domain_name"`
	Type       string `json:"type"`
	IAMURL     string `json:"iam_url"`
	STSURL     string `json:"sts_url"`
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
	AccessKeys []ManagedObjectStorageUserAccessKey `json:"access_keys"`
	Arn        string                              `json:"arn"`
	CreatedAt  time.Time                           `json:"created_at"`
	Policies   []ManagedObjectStoragePolicy        `json:"policies"`
	Username   string                              `json:"username"`
}

// ManagedObjectStoragePolicy represents a policy for the Managed Object Storage service
type ManagedObjectStoragePolicy struct {
	ARN              string    `json:"arn"`
	AttachmentCount  int       `json:"attachment_count"`
	CreatedAt        time.Time `json:"created_at"`
	DefaultVersionID string    `json:"default_version_id"`
	Description      string    `json:"description"`
	Document         string    `json:"document"`
	Name             string    `json:"name"`
	System           bool      `json:"system"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// ManagedObjectStorageUserPolicy represents a policy attached to a Managed Object Storage user
type ManagedObjectStorageUserPolicy struct {
	ARN  string `json:"arn"`
	Name string `json:"name"`
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
	AccessKeyID     string                                  `json:"access_key_id"`
	CreatedAt       time.Time                               `json:"created_at"`
	LastUsedAt      time.Time                               `json:"last_used_at"`
	SecretAccessKey *string                                 `json:"secret_access_key,omitempty"`
	Status          ManagedObjectStorageUserAccessKeyStatus `json:"status"`
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
