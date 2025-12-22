package request

import (
	"encoding/json"
	"fmt"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
)

const (
	managedObjectStorageBasePath string = "/object-storage-2"
)

// ManagedObjectStorageUser represents a user
type ManagedObjectStorageUser struct {
	Username string `json:"username"`
}

// GetManagedObjectStorageRegionsRequest represents a request for retrieving Managed Object Storage regions
type GetManagedObjectStorageRegionsRequest struct{}

// RequestURL implements the Request interface
func (r *GetManagedObjectStorageRegionsRequest) RequestURL() string {
	return fmt.Sprintf("%s/regions", managedObjectStorageBasePath)
}

// GetManagedObjectStorageRegionRequest represents a request for retrieving details about a Managed Object Storage region
type GetManagedObjectStorageRegionRequest struct {
	Name string `json:"-"`
}

// RequestURL implements the Request interface
func (r *GetManagedObjectStorageRegionRequest) RequestURL() string {
	return fmt.Sprintf("%s/regions/%s", managedObjectStorageBasePath, r.Name)
}

// CreateManagedObjectStorageRequest represents a request for creating a new Managed Object Storage service
type CreateManagedObjectStorageRequest struct {
	ConfiguredStatus upcloud.ManagedObjectStorageConfiguredStatus `json:"configured_status"`
	Labels           []upcloud.Label                              `json:"labels,omitempty"`
	Name             string                                       `json:"name,omitempty"`
	Networks         []upcloud.ManagedObjectStorageNetwork        `json:"networks"`
	Region           string                                       `json:"region"`
}

// RequestURL implements the Request interface
func (r *CreateManagedObjectStorageRequest) RequestURL() string {
	return managedObjectStorageBasePath
}

// GetManagedObjectStoragesRequest represents a request to list Managed Object Storage services
// List size can be filtered using optional Page object
type GetManagedObjectStoragesRequest struct {
	Page *Page `json:"-"`
}

func (r *GetManagedObjectStoragesRequest) RequestURL() string {
	if r.Page != nil {
		return fmt.Sprintf("%s?%s", managedObjectStorageBasePath, r.Page.String())
	}
	return managedObjectStorageBasePath
}

// GetManagedObjectStorageRequest represents a request for retrieving details about a Managed Object Storage service
type GetManagedObjectStorageRequest struct {
	UUID string `json:"-"`
}

// RequestURL implements the Request interface
func (r *GetManagedObjectStorageRequest) RequestURL() string {
	return fmt.Sprintf("%s/%s", managedObjectStorageBasePath, r.UUID)
}

// ReplaceManagedObjectStorageRequest represents a request to replace a Managed Object Storage service
type ReplaceManagedObjectStorageRequest struct {
	ConfiguredStatus upcloud.ManagedObjectStorageConfiguredStatus `json:"configured_status"`
	Labels           []upcloud.Label                              `json:"labels,omitempty"`
	Name             string                                       `json:"name,omitempty"`
	Networks         []upcloud.ManagedObjectStorageNetwork        `json:"networks"`
	UUID             string                                       `json:"-"`
}

func (r *ReplaceManagedObjectStorageRequest) RequestURL() string {
	return fmt.Sprintf("%s/%s", managedObjectStorageBasePath, r.UUID)
}

// ModifyManagedObjectStorageRequest represents a request to modify a Object Storage
type ModifyManagedObjectStorageRequest struct {
	ConfiguredStatus *upcloud.ManagedObjectStorageConfiguredStatus `json:"configured_status,omitempty"`
	Labels           *[]upcloud.Label                              `json:"labels,omitempty"`
	Name             *string                                       `json:"name,omitempty"`
	Networks         *[]upcloud.ManagedObjectStorageNetwork        `json:"networks,omitempty"`
	UUID             string                                        `json:"-"`
}

// RequestURL implements the Request interface
func (r *ModifyManagedObjectStorageRequest) RequestURL() string {
	return fmt.Sprintf("%s/%s", managedObjectStorageBasePath, r.UUID)
}

// DeleteManagedObjectStorageRequest represents a request to delete a Managed Object Storage service
type DeleteManagedObjectStorageRequest struct {
	UUID string `json:"-"`
}

// RequestURL implements the Request interface
func (r *DeleteManagedObjectStorageRequest) RequestURL() string {
	return fmt.Sprintf("%s/%s", managedObjectStorageBasePath, r.UUID)
}

// GetManagedObjectStorageMetricsRequest represents a request for retrieving metrics
type GetManagedObjectStorageMetricsRequest struct {
	ServiceUUID string `json:"-"`
}

// RequestURL implements the Request interface
func (r *GetManagedObjectStorageMetricsRequest) RequestURL() string {
	return fmt.Sprintf("%s/%s/metrics", managedObjectStorageBasePath, r.ServiceUUID)
}

// GetManagedObjectStorageBucketMetricsRequest represents a request for retrieving buckets' metrics
type GetManagedObjectStorageBucketMetricsRequest struct {
	Page        *Page  `json:"-"`
	ServiceUUID string `json:"-"`
}

// RequestURL implements the Request interface
func (r *GetManagedObjectStorageBucketMetricsRequest) RequestURL() string {
	path := fmt.Sprintf("%s/%s/buckets", managedObjectStorageBasePath, r.ServiceUUID)
	if r.Page != nil {
		return fmt.Sprintf("%s?%s", path, r.Page.String())
	}

	return path
}

type CreateManagedObjectStorageBucketRequest struct {
	ServiceUUID string `json:"-"`
	Name        string `json:"name"`
}

// RequestURL implements the Request interface
func (r *CreateManagedObjectStorageBucketRequest) RequestURL() string {
	return fmt.Sprintf("%s/%s/buckets", managedObjectStorageBasePath, r.ServiceUUID)
}

type DeleteManagedObjectStorageBucketRequest struct {
	ServiceUUID string `json:"-"`
	Name        string `json:"-"`
}

// RequestURL implements the Request interface
func (r *DeleteManagedObjectStorageBucketRequest) RequestURL() string {
	return fmt.Sprintf("%s/%s/buckets/%s", managedObjectStorageBasePath, r.ServiceUUID, r.Name)
}

// CreateManagedObjectStorageNetworkRequest represents a request for creating a network
type CreateManagedObjectStorageNetworkRequest struct {
	Family      string `json:"family"`
	Name        string `json:"name"`
	ServiceUUID string `json:"-"`
	Type        string `json:"type"`
	UUID        string `json:"uuid,omitempty"`
}

// RequestURL implements the Request interface
func (r *CreateManagedObjectStorageNetworkRequest) RequestURL() string {
	return fmt.Sprintf("%s/%s/networks", managedObjectStorageBasePath, r.ServiceUUID)
}

// GetManagedObjectStorageNetworksRequest represents a request for retrieving networks
type GetManagedObjectStorageNetworksRequest struct {
	ServiceUUID string `json:"-"`
}

// RequestURL implements the Request interface
func (r *GetManagedObjectStorageNetworksRequest) RequestURL() string {
	return fmt.Sprintf("%s/%s/networks", managedObjectStorageBasePath, r.ServiceUUID)
}

// GetManagedObjectStorageNetworkRequest represents a request for retrieving details about a network
type GetManagedObjectStorageNetworkRequest struct {
	ServiceUUID string `json:"-"`
	NetworkName string `json:"-"`
}

// RequestURL implements the Request interface
func (r *GetManagedObjectStorageNetworkRequest) RequestURL() string {
	return fmt.Sprintf("%s/%s/networks/%s", managedObjectStorageBasePath, r.ServiceUUID, r.NetworkName)
}

// DeleteManagedObjectStorageNetworkRequest represents a request to delete a network
type DeleteManagedObjectStorageNetworkRequest struct {
	ServiceUUID string `json:"-"`
	NetworkName string `json:"-"`
}

// RequestURL implements the Request interface
func (r *DeleteManagedObjectStorageNetworkRequest) RequestURL() string {
	return fmt.Sprintf("%s/%s/networks/%s", managedObjectStorageBasePath, r.ServiceUUID, r.NetworkName)
}

// CreateManagedObjectStorageUserRequest represents a request for creating a user
type CreateManagedObjectStorageUserRequest struct {
	Username    string `json:"username"`
	ServiceUUID string `json:"-"`
}

// RequestURL implements the Request interface
func (r *CreateManagedObjectStorageUserRequest) RequestURL() string {
	return fmt.Sprintf("%s/%s/users", managedObjectStorageBasePath, r.ServiceUUID)
}

// GetManagedObjectStorageUsersRequest represents a request for retrieving users
type GetManagedObjectStorageUsersRequest struct {
	ServiceUUID string `json:"-"`
}

// RequestURL implements the Request interface
func (r *GetManagedObjectStorageUsersRequest) RequestURL() string {
	return fmt.Sprintf("%s/%s/users", managedObjectStorageBasePath, r.ServiceUUID)
}

// GetManagedObjectStorageUserRequest represents a request for retrieving details about a user
type GetManagedObjectStorageUserRequest struct {
	ServiceUUID string `json:"-"`
	Username    string `json:"-"`
}

// RequestURL implements the Request interface
func (r *GetManagedObjectStorageUserRequest) RequestURL() string {
	return fmt.Sprintf("%s/%s/users/%s", managedObjectStorageBasePath, r.ServiceUUID, r.Username)
}

// DeleteManagedObjectStorageUserRequest represents a request to delete a user
type DeleteManagedObjectStorageUserRequest struct {
	ServiceUUID string `json:"-"`
	Username    string `json:"-"`
}

// RequestURL implements the Request interface
func (r *DeleteManagedObjectStorageUserRequest) RequestURL() string {
	return fmt.Sprintf("%s/%s/users/%s", managedObjectStorageBasePath, r.ServiceUUID, r.Username)
}

// CreateManagedObjectStorageUserAccessKeyRequest represents a request for creating an access key
type CreateManagedObjectStorageUserAccessKeyRequest struct {
	Username    string `json:"-"`
	ServiceUUID string `json:"-"`
}

// RequestURL implements the Request interface
func (r *CreateManagedObjectStorageUserAccessKeyRequest) RequestURL() string {
	return fmt.Sprintf("%s/%s/users/%s/access-keys", managedObjectStorageBasePath, r.ServiceUUID, r.Username)
}

// GetManagedObjectStorageUserAccessKeysRequest represents a request for retrieving access keys
type GetManagedObjectStorageUserAccessKeysRequest struct {
	ServiceUUID string `json:"-"`
	Username    string `json:"-"`
}

// RequestURL implements the Request interface
func (r *GetManagedObjectStorageUserAccessKeysRequest) RequestURL() string {
	return fmt.Sprintf("%s/%s/users/%s/access-keys", managedObjectStorageBasePath, r.ServiceUUID, r.Username)
}

// GetManagedObjectStorageUserAccessKeyRequest represents a request for retrieving details about an access key
type GetManagedObjectStorageUserAccessKeyRequest struct {
	ServiceUUID string `json:"-"`
	Username    string `json:"-"`
	AccessKeyID string `json:"-"`
}

// RequestURL implements the Request interface
func (r *GetManagedObjectStorageUserAccessKeyRequest) RequestURL() string {
	return fmt.Sprintf("%s/%s/users/%s/access-keys/%s", managedObjectStorageBasePath, r.ServiceUUID, r.Username, r.AccessKeyID)
}

// ModifyManagedObjectStorageUserAccessKeyRequest represents a request for creating an access key
type ModifyManagedObjectStorageUserAccessKeyRequest struct {
	Username    string                                          `json:"-"`
	ServiceUUID string                                          `json:"-"`
	AccessKeyID string                                          `json:"-"`
	Status      upcloud.ManagedObjectStorageUserAccessKeyStatus `json:"status,omitempty"`
}

// RequestURL implements the Request interface
func (r *ModifyManagedObjectStorageUserAccessKeyRequest) RequestURL() string {
	return fmt.Sprintf("%s/%s/users/%s/access-keys/%s", managedObjectStorageBasePath, r.ServiceUUID, r.Username, r.AccessKeyID)
}

// DeleteManagedObjectStorageUserAccessKeyRequest represents a request to delete a Managed Object Storage service
type DeleteManagedObjectStorageUserAccessKeyRequest struct {
	ServiceUUID string `json:"-"`
	Username    string `json:"-"`
	AccessKeyID string `json:"-"`
}

// RequestURL implements the Request interface
func (r *DeleteManagedObjectStorageUserAccessKeyRequest) RequestURL() string {
	return fmt.Sprintf("%s/%s/users/%s/access-keys/%s", managedObjectStorageBasePath, r.ServiceUUID, r.Username, r.AccessKeyID)
}

// CreateManagedObjectStoragePolicyRequest represents a request for creating a policy
type CreateManagedObjectStoragePolicyRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Document    string `json:"document"`
	ServiceUUID string `json:"-"`
}

// RequestURL implements the Request interface
func (r *CreateManagedObjectStoragePolicyRequest) RequestURL() string {
	return fmt.Sprintf("%s/%s/policies", managedObjectStorageBasePath, r.ServiceUUID)
}

// GetManagedObjectStoragePoliciesRequest represents a request for retrieving policies
type GetManagedObjectStoragePoliciesRequest struct {
	ServiceUUID string `json:"-"`
}

// RequestURL implements the Request interface
func (r *GetManagedObjectStoragePoliciesRequest) RequestURL() string {
	return fmt.Sprintf("%s/%s/policies", managedObjectStorageBasePath, r.ServiceUUID)
}

// GetManagedObjectStoragePolicyRequest represents a request for retrieving details about a policy
type GetManagedObjectStoragePolicyRequest struct {
	Name        string `json:"-"`
	ServiceUUID string `json:"-"`
}

// RequestURL implements the Request interface
func (r *GetManagedObjectStoragePolicyRequest) RequestURL() string {
	return fmt.Sprintf("%s/%s/policies/%s", managedObjectStorageBasePath, r.ServiceUUID, r.Name)
}

// DeleteManagedObjectStoragePolicyRequest represents a request to delete a policy
type DeleteManagedObjectStoragePolicyRequest struct {
	ServiceUUID string `json:"-"`
	Name        string `json:"-"`
}

// RequestURL implements the Request interface
func (r *DeleteManagedObjectStoragePolicyRequest) RequestURL() string {
	return fmt.Sprintf("%s/%s/policies/%s", managedObjectStorageBasePath, r.ServiceUUID, r.Name)
}

// GetManagedObjectStorageUserPoliciesRequest represents a request for retrieving policies attached to a user
type GetManagedObjectStorageUserPoliciesRequest struct {
	ServiceUUID string `json:"-"`
	Username    string `json:"-"`
}

// RequestURL implements the Request interface
func (r *GetManagedObjectStorageUserPoliciesRequest) RequestURL() string {
	return fmt.Sprintf("%s/%s/users/%s/policies", managedObjectStorageBasePath, r.ServiceUUID, r.Username)
}

// AttachManagedObjectStorageUserPolicyRequest represents a request for attaching a policy to a user
type AttachManagedObjectStorageUserPolicyRequest struct {
	ServiceUUID string `json:"-"`
	Username    string `json:"-"`
	Name        string `json:"name"`
}

// RequestURL implements the Request interface
func (r *AttachManagedObjectStorageUserPolicyRequest) RequestURL() string {
	return fmt.Sprintf("%s/%s/users/%s/policies", managedObjectStorageBasePath, r.ServiceUUID, r.Username)
}

// DetachManagedObjectStorageUserPolicyRequest represents a request for detaching a policy to a user
type DetachManagedObjectStorageUserPolicyRequest struct {
	ServiceUUID string `json:"-"`
	Username    string `json:"-"`
	Name        string `json:"-"`
}

// RequestURL implements the Request interface
func (r *DetachManagedObjectStorageUserPolicyRequest) RequestURL() string {
	return fmt.Sprintf("%s/%s/users/%s/policies/%s", managedObjectStorageBasePath, r.ServiceUUID, r.Username, r.Name)
}

// CreateManagedObjectStorageCustomDomainRequest represents a request for creating a policy
type CreateManagedObjectStorageCustomDomainRequest struct {
	DomainName  string `json:"domain_name"`
	Type        string `json:"type"`
	ServiceUUID string `json:"-"`
}

// RequestURL implements the Request interface
func (r *CreateManagedObjectStorageCustomDomainRequest) RequestURL() string {
	return fmt.Sprintf("%s/%s/custom-domains", managedObjectStorageBasePath, r.ServiceUUID)
}

// GetManagedObjectStorageCustomDomainsRequest represents a request for retrieving policies
type GetManagedObjectStorageCustomDomainsRequest struct {
	ServiceUUID string `json:"-"`
}

// RequestURL implements the Request interface
func (r *GetManagedObjectStorageCustomDomainsRequest) RequestURL() string {
	return fmt.Sprintf("%s/%s/custom-domains", managedObjectStorageBasePath, r.ServiceUUID)
}

// GetManagedObjectStorageCustomDomainRequest represents a request for retrieving details about a policy
type GetManagedObjectStorageCustomDomainRequest struct {
	DomainName  string `json:"-"`
	ServiceUUID string `json:"-"`
}

// RequestURL implements the Request interface
func (r *GetManagedObjectStorageCustomDomainRequest) RequestURL() string {
	return fmt.Sprintf("%s/%s/custom-domains/%s", managedObjectStorageBasePath, r.ServiceUUID, r.DomainName)
}

type ModifyCustomDomain struct {
	DomainName string `json:"domain_name"`
	Type       string `json:"type"`
}

// ModifyManagedObjectStorageCustomDomainRequest represents a request for retrieving details about a policy
type ModifyManagedObjectStorageCustomDomainRequest struct {
	DomainName   string `json:"-"`
	ServiceUUID  string `json:"-"`
	CustomDomain ModifyCustomDomain
}

// RequestURL implements the Request interface
func (r *ModifyManagedObjectStorageCustomDomainRequest) RequestURL() string {
	return fmt.Sprintf("%s/%s/custom-domains/%s", managedObjectStorageBasePath, r.ServiceUUID, r.DomainName)
}

func (r *ModifyManagedObjectStorageCustomDomainRequest) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.CustomDomain)
}

// DeleteManagedObjectStorageCustomDomainRequest represents a request to delete a policy
type DeleteManagedObjectStorageCustomDomainRequest struct {
	ServiceUUID string `json:"-"`
	DomainName  string `json:"-"`
}

// RequestURL implements the Request interface
func (r *DeleteManagedObjectStorageCustomDomainRequest) RequestURL() string {
	return fmt.Sprintf("%s/%s/custom-domains/%s", managedObjectStorageBasePath, r.ServiceUUID, r.DomainName)
}

// WaitForManagedObjectStorageOperationalStateRequest represents a request to wait for a Managed Object Storage service
// to enter a desired state
type WaitForManagedObjectStorageOperationalStateRequest struct {
	DesiredState upcloud.ManagedObjectStorageOperationalState `json:"-"`
	UUID         string                                       `json:"-"`
}

// RequestURL implements the Request interface
func (r *WaitForManagedObjectStorageOperationalStateRequest) RequestURL() string {
	return fmt.Sprintf("%s/%s", managedObjectStorageBasePath, r.UUID)
}

// WaitForManagedObjectStorageDeletionRequest represents a request to wait for a Managed Object Storage service
// to be deleted
type WaitForManagedObjectStorageDeletionRequest struct {
	DesiredState upcloud.ManagedObjectStorageOperationalState `json:"-"` // Deprecated: The managed object storage instance has no state after being deleted.
	UUID         string                                       `json:"-"`
}

// RequestURL implements the Request interface
func (r *WaitForManagedObjectStorageDeletionRequest) RequestURL() string {
	return fmt.Sprintf("%s/%s", managedObjectStorageBasePath, r.UUID)
}

// WaitForManagedObjectStorageDeletionRequest represents a request to wait for a Managed Object Storage service
// to be deleted
type WaitForManagedObjectStorageBucketDeletionRequest struct {
	ServiceUUID string `json:"-"`
	Name        string `json:"-"`
}

// RequestURL implements the Request interface
func (r *WaitForManagedObjectStorageBucketDeletionRequest) RequestURL() string {
	return fmt.Sprintf("%s/%s/buckets", managedObjectStorageBasePath, r.ServiceUUID)
}

type CreateManagedObjectStoragePolicyVersionRequest struct {
	ServiceUUID string `json:"-"`
	Name        string `json:"-"`
	Document    string `json:"document"`
	IsDefault   bool   `json:"is_default,omitempty"`
}

func (r *CreateManagedObjectStoragePolicyVersionRequest) RequestURL() string {
	return fmt.Sprintf(
		"%s/%s/policies/%s/versions",
		managedObjectStorageBasePath,
		r.ServiceUUID,
		r.Name,
	)
}
