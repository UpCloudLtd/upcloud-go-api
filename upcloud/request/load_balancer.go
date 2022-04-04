package request

import (
	"encoding/json"
	"fmt"

	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud"
)

const loadBalancerCertificateBundleBaseURL string = "/load-balancer/certificate-bundles"

// GetLoadBalancersRequest represents a request to list load balancers
// List size can be filtered using optional Page object
type GetLoadBalancersRequest struct {
	Page *Page
}

func (r *GetLoadBalancersRequest) RequestURL() string {
	if r.Page != nil {
		return fmt.Sprintf("/load-balancer?%s", r.Page.String())
	}
	return "/load-balancer"
}

// GetLoadBalancerRequest represents a request to get load balancer details
type GetLoadBalancerRequest struct {
	UUID string
}

func (r *GetLoadBalancerRequest) RequestURL() string {
	return fmt.Sprintf("/load-balancer/%s", r.UUID)
}

// CreateLoadBalancerRequest represents a request to create load balancer
type CreateLoadBalancerRequest struct {
	Name             string                               `json:"name,omitempty"`
	Plan             string                               `json:"plan,omitempty"`
	Zone             string                               `json:"zone,omitempty"`
	NetworkUUID      string                               `json:"network_uuid,omitempty"`
	ConfiguredStatus upcloud.LoadBalancerConfiguredStatus `json:"configured_status,omitempty"`
	Frontends        []LoadBalancerFrontend               `json:"frontends"`
	Backends         []LoadBalancerBackend                `json:"backends"`
	Resolvers        []LoadBalancerResolver               `json:"resolvers"`
}

func (r *CreateLoadBalancerRequest) RequestURL() string {
	return "/load-balancer"
}

// CreateLoadBalancerRequest represents a request to modify load balancer
type ModifyLoadBalancerRequest struct {
	UUID             string `json:"-"`
	Name             string `json:"name,omitempty"`
	Plan             string `json:"plan,omitempty"`
	ConfiguredStatus string `json:"configured_status,omitempty"`
}

func (r *ModifyLoadBalancerRequest) RequestURL() string {
	return fmt.Sprintf("/load-balancer/%s", r.UUID)
}

// DeleteLoadBalancerRequest represents a request to delete load balancer
type DeleteLoadBalancerRequest struct {
	UUID string `json:"-"`
}

func (r *DeleteLoadBalancerRequest) RequestURL() string {
	return fmt.Sprintf("/load-balancer/%s", r.UUID)
}

// GetLoadBalancerBackendsRequest represents a request to list load balancer backends
type GetLoadBalancerBackendsRequest struct {
	ServiceUUID string `json:"-"`
}

func (r *GetLoadBalancerBackendsRequest) RequestURL() string {
	return fmt.Sprintf("/load-balancer/%s/backends", r.ServiceUUID)
}

// BalancerBackend represents the payload for CreateLoadBalancerBackendRequest
type LoadBalancerBackend struct {
	Name     string                      `json:"name"`
	Resolver string                      `json:"resolver,omitempty"`
	Members  []LoadBalancerBackendMember `json:"members"`
}

// CreateLoadBalancerBackendRequest represents a request to create load balancer backend
type CreateLoadBalancerBackendRequest struct {
	ServiceUUID string `json:"-"`
	Backend     LoadBalancerBackend
}

func (r *CreateLoadBalancerBackendRequest) RequestURL() string {
	return fmt.Sprintf("/load-balancer/%s/backends", r.ServiceUUID)
}

func (r *CreateLoadBalancerBackendRequest) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.Backend)
}

// GetLoadBalancerBackendRequest represents a request to get load balancer backend details
type GetLoadBalancerBackendRequest struct {
	ServiceUUID string `json:"-"`
	Name        string `json:"-"`
}

func (r *GetLoadBalancerBackendRequest) RequestURL() string {
	return fmt.Sprintf("/load-balancer/%s/backends/%s", r.ServiceUUID, r.Name)
}

// ModifyLoadBalancerBackend represents the payload for ModifyLoadBalancerBackendRequest
type ModifyLoadBalancerBackend struct {
	Name     string `json:"name,omitempty"`
	Resolver string `json:"resolver,omitempty"`
}

// ModifyLoadBalancerBackendRequest represents a request to modify load balancer backend
type ModifyLoadBalancerBackendRequest struct {
	ServiceUUID string `json:"-"`
	Name        string `json:"-"`
	Backend     ModifyLoadBalancerBackend
}

func (r *ModifyLoadBalancerBackendRequest) RequestURL() string {
	return fmt.Sprintf("/load-balancer/%s/backends/%s", r.ServiceUUID, r.Name)
}

func (r *ModifyLoadBalancerBackendRequest) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.Backend)
}

// DeleteLoadBalancerBackendRequest represents a request to delete load balancer backend
type DeleteLoadBalancerBackendRequest struct {
	ServiceUUID string `json:"-"`
	Name        string `json:"-"`
}

func (r *DeleteLoadBalancerBackendRequest) RequestURL() string {
	return fmt.Sprintf("/load-balancer/%s/backends/%s", r.ServiceUUID, r.Name)
}

// LoadBalancerBackendMember represents the payload for backend member request
type LoadBalancerBackendMember struct {
	Name        string                                `json:"name,omitempty"`
	Weight      int                                   `json:"weight"`
	MaxSessions int                                   `json:"max_sessions"`
	Enabled     bool                                  `json:"enabled"`
	Type        upcloud.LoadBalancerBackendMemberType `json:"type,omitempty"`
	IP          string                                `json:"ip,omitempty"`
	Port        int                                   `json:"port,omitempty"`
}

// CreateLoadBalancerBackendMemberRequest represents a request to create load balancer backend member
type CreateLoadBalancerBackendMemberRequest struct {
	ServiceUUID string `json:"-"`
	BackendName string `json:"-"`
	Member      LoadBalancerBackendMember
}

func (r *CreateLoadBalancerBackendMemberRequest) RequestURL() string {
	return fmt.Sprintf("/load-balancer/%s/backends/%s/members", r.ServiceUUID, r.BackendName)
}

func (r *CreateLoadBalancerBackendMemberRequest) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.Member)
}

// GetLoadBalancerBackendMembersRequest represents a request to get load balancer backend member list
type GetLoadBalancerBackendMembersRequest struct {
	ServiceUUID string `json:"-"`
	BackendName string `json:"-"`
}

func (r *GetLoadBalancerBackendMembersRequest) RequestURL() string {
	return fmt.Sprintf("/load-balancer/%s/backends/%s/members", r.ServiceUUID, r.BackendName)
}

// GetLoadBalancerBackendMemberRequest represents a request to get load balancer backend member details
type GetLoadBalancerBackendMemberRequest struct {
	ServiceUUID string `json:"-"`
	BackendName string `json:"-"`
	Name        string `json:"-"`
}

func (r *GetLoadBalancerBackendMemberRequest) RequestURL() string {
	return fmt.Sprintf("/load-balancer/%s/backends/%s/members/%s", r.ServiceUUID, r.BackendName, r.Name)
}

// ModifyLoadBalancerBackendMemberRequest represents a request to modify load balancer backend member
type ModifyLoadBalancerBackendMemberRequest struct {
	ServiceUUID string `json:"-"`
	BackendName string `json:"-"`
	Name        string `json:"-"`
	Member      LoadBalancerBackendMember
}

func (r *ModifyLoadBalancerBackendMemberRequest) RequestURL() string {
	return fmt.Sprintf("/load-balancer/%s/backends/%s/members/%s", r.ServiceUUID, r.BackendName, r.Name)
}

func (r *ModifyLoadBalancerBackendMemberRequest) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.Member)
}

// DeleteLoadBalancerBackendMemberRequest represents a request to delete load balancer backend member
type DeleteLoadBalancerBackendMemberRequest struct {
	ServiceUUID string `json:"-"`
	BackendName string `json:"-"`
	Name        string `json:"-"`
}

func (r *DeleteLoadBalancerBackendMemberRequest) RequestURL() string {
	return fmt.Sprintf("/load-balancer/%s/backends/%s/members/%s", r.ServiceUUID, r.BackendName, r.Name)
}

// LoadBalancerResolver represents resolver payload
type LoadBalancerResolver struct {
	Name         string   `json:"name,omitempty"`
	Nameservers  []string `json:"nameservers,omitempty"`
	Retries      int      `json:"retries,omitempty"`
	Timeout      int      `json:"timeout,omitempty"`
	TimeoutRetry int      `json:"timeout_retry,omitempty"`
	CacheValid   int      `json:"cache_valid,omitempty"`
	CacheInvalid int      `json:"cache_invalid,omitempty"`
}

// CreateLoadBalancerResolverRequest represents a request to create load balancer resolver
type CreateLoadBalancerResolverRequest struct {
	ServiceUUID string `json:"-"`
	Resolver    LoadBalancerResolver
}

func (r *CreateLoadBalancerResolverRequest) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.Resolver)
}

func (r *CreateLoadBalancerResolverRequest) RequestURL() string {
	return fmt.Sprintf("/load-balancer/%s/resolvers", r.ServiceUUID)
}

// ModifyLoadBalancerRevolverRequest represents a request to modify load balancer resolver
type ModifyLoadBalancerRevolverRequest struct {
	ServiceUUID string `json:"-"`
	Name        string `json:"-"`
	Resolver    LoadBalancerResolver
}

func (r *ModifyLoadBalancerRevolverRequest) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.Resolver)
}

func (r *ModifyLoadBalancerRevolverRequest) RequestURL() string {
	return fmt.Sprintf("/load-balancer/%s/resolvers/%s", r.ServiceUUID, r.Name)
}

// GetLoadBalancerResolversRequest represents a request to get load balancer resolver list
type GetLoadBalancerResolversRequest struct {
	ServiceUUID string `json:"-"`
}

func (r *GetLoadBalancerResolversRequest) RequestURL() string {
	return fmt.Sprintf("/load-balancer/%s/resolvers", r.ServiceUUID)
}

// GetLoadBalancerResolverRequest represents a request to get load balancer resolver details
type GetLoadBalancerResolverRequest struct {
	ServiceUUID string `json:"-"`
	Name        string `json:"-"`
}

func (r *GetLoadBalancerResolverRequest) RequestURL() string {
	return fmt.Sprintf("/load-balancer/%s/resolvers/%s", r.ServiceUUID, r.Name)
}

// DeleteLoadBalancerResolverRequest represents a request to delete load balancer resolver
type DeleteLoadBalancerResolverRequest struct {
	ServiceUUID string `json:"-"`
	Name        string `json:"-"`
}

func (r *DeleteLoadBalancerResolverRequest) RequestURL() string {
	return fmt.Sprintf("/load-balancer/%s/resolvers/%s", r.ServiceUUID, r.Name)
}

// GetLoadBalancerPlansRequest represents a request to list load balancer plans
type GetLoadBalancerPlansRequest struct{}

func (r *GetLoadBalancerPlansRequest) RequestURL() string {
	return "/load-balancer/plans"
}

// GetLoadBalancerFrontendsRequest represents a request to list load balancer frontends
type GetLoadBalancerFrontendsRequest struct {
	ServiceUUID string `json:"-"`
}

func (r *GetLoadBalancerFrontendsRequest) RequestURL() string {
	return fmt.Sprintf("/load-balancer/%s/frontends", r.ServiceUUID)
}

// GetLoadBalancerFrontendRequest represents a request to get load balancer frontend details
type GetLoadBalancerFrontendRequest struct {
	ServiceUUID string `json:"-"`
	Name        string `json:"-"`
}

func (r *GetLoadBalancerFrontendRequest) RequestURL() string {
	return fmt.Sprintf("/load-balancer/%s/frontends/%s", r.ServiceUUID, r.Name)
}

// LoadBalancerFrontendRule represents frontend rule payload
type LoadBalancerFrontendRule struct {
	Name     string                        `json:"name"`
	Priority int                           `json:"priority"`
	Matchers []upcloud.LoadBalancerMatcher `json:"matchers"`
	Actions  []upcloud.LoadBalancerAction  `json:"actions"`
}

// LoadBalancerFrontend represents frontend payload
type LoadBalancerFrontend struct {
	Name           string                          `json:"name,omitempty"`
	Mode           upcloud.LoadBalancerMode        `json:"mode,omitempty"`
	Port           int                             `json:"port,omitempty"`
	DefaultBackend string                          `json:"default_backend,omitempty"`
	Rules          []LoadBalancerFrontendRule      `json:"rules,omitempty"`
	TLSConfigs     []LoadBalancerFrontendTLSConfig `json:"tls_configs,omitempty"`
}

// CreateLoadBalancerFrontendRequest represents a request to create load balancer frontend
type CreateLoadBalancerFrontendRequest struct {
	ServiceUUID string
	Frontend    LoadBalancerFrontend
}

func (r *CreateLoadBalancerFrontendRequest) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.Frontend)
}

func (r *CreateLoadBalancerFrontendRequest) RequestURL() string {
	return fmt.Sprintf("/load-balancer/%s/frontends", r.ServiceUUID)
}

// ModifyLoadBalancerFrontend represents payload to modify frontend
type ModifyLoadBalancerFrontend struct {
	Name           string                   `json:"name,omitempty"`
	Mode           upcloud.LoadBalancerMode `json:"mode,omitempty"`
	Port           int                      `json:"port,omitempty"`
	DefaultBackend string                   `json:"default_backend,omitempty"`
}

// ModifyLoadBalancerFrontendRequest represents a request to modify load balancer frontend
type ModifyLoadBalancerFrontendRequest struct {
	ServiceUUID string `json:"-"`
	Name        string `json:"-"`
	Frontend    ModifyLoadBalancerFrontend
}

func (r *ModifyLoadBalancerFrontendRequest) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.Frontend)
}

func (r *ModifyLoadBalancerFrontendRequest) RequestURL() string {
	return fmt.Sprintf("/load-balancer/%s/frontends/%s", r.ServiceUUID, r.Name)
}

// DeleteLoadBalancerFrontendRequest represents a request to delete load balancer frontend
type DeleteLoadBalancerFrontendRequest struct {
	ServiceUUID string `json:"-"`
	Name        string `json:"-"`
}

func (r *DeleteLoadBalancerFrontendRequest) RequestURL() string {
	return fmt.Sprintf("/load-balancer/%s/frontends/%s", r.ServiceUUID, r.Name)
}

// GetLoadBalancerFrontendRulesRequest represents a request to list frontend rules
type GetLoadBalancerFrontendRulesRequest struct {
	ServiceUUID  string `json:"-"`
	FrontendName string `json:"-"`
}

func (r *GetLoadBalancerFrontendRulesRequest) RequestURL() string {
	return fmt.Sprintf("/load-balancer/%s/frontends/%s/rules", r.ServiceUUID, r.FrontendName)
}

// GetLoadBalancerFrontendRuleRequest represents a request to get frontend rule details
type GetLoadBalancerFrontendRuleRequest struct {
	ServiceUUID  string `json:"-"`
	FrontendName string `json:"-"`
	Name         string `json:"-"`
}

func (r *GetLoadBalancerFrontendRuleRequest) RequestURL() string {
	return fmt.Sprintf("/load-balancer/%s/frontends/%s/rules/%s", r.ServiceUUID, r.FrontendName, r.Name)
}

// CreateLoadBalancerFrontendRuleRequest represents a request to create frontend rule
type CreateLoadBalancerFrontendRuleRequest struct {
	ServiceUUID  string `json:"-"`
	FrontendName string `json:"-"`
	Rule         LoadBalancerFrontendRule
}

func (r *CreateLoadBalancerFrontendRuleRequest) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.Rule)
}

func (r *CreateLoadBalancerFrontendRuleRequest) RequestURL() string {
	return fmt.Sprintf("/load-balancer/%s/frontends/%s/rules", r.ServiceUUID, r.FrontendName)
}

// ReplaceLoadBalancerFrontendRuleRequest represents a request to replace frontend rule
type ReplaceLoadBalancerFrontendRuleRequest struct {
	ServiceUUID  string `json:"-"`
	FrontendName string `json:"-"`
	Name         string `json:"-"`
	Rule         LoadBalancerFrontendRule
}

func (r *ReplaceLoadBalancerFrontendRuleRequest) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.Rule)
}

func (r *ReplaceLoadBalancerFrontendRuleRequest) RequestURL() string {
	return fmt.Sprintf("/load-balancer/%s/frontends/%s/rules/%s", r.ServiceUUID, r.FrontendName, r.Name)
}

// ModifyLoadBalancerFrontendRule represents frontend rule modification payload
type ModifyLoadBalancerFrontendRule struct {
	Name     string `json:"name,omitempty"`
	Priority int    `json:"priority,omitempty"`
}

// ModifyLoadBalancerFrontendRuleRequest represents a request to modify frontend rule
type ModifyLoadBalancerFrontendRuleRequest struct {
	ServiceUUID  string `json:"-"`
	FrontendName string `json:"-"`
	Name         string `json:"-"`
	Rule         ModifyLoadBalancerFrontendRule
}

func (r *ModifyLoadBalancerFrontendRuleRequest) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.Rule)
}

func (r *ModifyLoadBalancerFrontendRuleRequest) RequestURL() string {
	return fmt.Sprintf("/load-balancer/%s/frontends/%s/rules/%s", r.ServiceUUID, r.FrontendName, r.Name)
}

// DeleteLoadBalancerFrontendRuleRequest represents a request to delete frontend rule
type DeleteLoadBalancerFrontendRuleRequest struct {
	ServiceUUID  string `json:"-"`
	FrontendName string `json:"-"`
	Name         string `json:"-"`
}

func (r *DeleteLoadBalancerFrontendRuleRequest) RequestURL() string {
	return fmt.Sprintf("/load-balancer/%s/frontends/%s/rules/%s", r.ServiceUUID, r.FrontendName, r.Name)
}

// LoadBalancerFrontendTLSConfig represents TLS config payload
type LoadBalancerFrontendTLSConfig struct {
	Name                  string `json:"name,omitempty"`
	CertificateBundleUUID string `json:"certificate_bundle_uuid,omitempty"`
}

// GetLoadBalancerFrontendTLSConfigsRequest represents a request to get frontend TLS configs
type GetLoadBalancerFrontendTLSConfigsRequest struct {
	ServiceUUID  string `json:"-"`
	FrontendName string `json:"-"`
}

func (r *GetLoadBalancerFrontendTLSConfigsRequest) RequestURL() string {
	return fmt.Sprintf("/load-balancer/%s/frontends/%s/tls-configs", r.ServiceUUID, r.FrontendName)
}

// GetLoadBalancerFrontendTLSConfigRequest represents a request to get frontend TLS config
type GetLoadBalancerFrontendTLSConfigRequest struct {
	ServiceUUID  string `json:"-"`
	FrontendName string `json:"-"`
	Name         string `json:"-"`
}

func (r *GetLoadBalancerFrontendTLSConfigRequest) RequestURL() string {
	return fmt.Sprintf("/load-balancer/%s/frontends/%s/tls-configs/%s", r.ServiceUUID, r.FrontendName, r.Name)
}

// CreateLoadBalancerFrontendTLSConfigRequest represents a request to create frontend TLS config
type CreateLoadBalancerFrontendTLSConfigRequest struct {
	ServiceUUID  string `json:"-"`
	FrontendName string `json:"-"`
	Config       LoadBalancerFrontendTLSConfig
}

func (r *CreateLoadBalancerFrontendTLSConfigRequest) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.Config)
}

func (r *CreateLoadBalancerFrontendTLSConfigRequest) RequestURL() string {
	return fmt.Sprintf("/load-balancer/%s/frontends/%s/tls-configs", r.ServiceUUID, r.FrontendName)
}

// ModifyLoadBalancerFrontendTLSConfigRequest represents a request to modify frontend TLS config
type ModifyLoadBalancerFrontendTLSConfigRequest struct {
	ServiceUUID  string `json:"-"`
	FrontendName string `json:"-"`
	Name         string `json:"-"`
	Config       LoadBalancerFrontendTLSConfig
}

func (r *ModifyLoadBalancerFrontendTLSConfigRequest) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.Config)
}

func (r *ModifyLoadBalancerFrontendTLSConfigRequest) RequestURL() string {
	return fmt.Sprintf("/load-balancer/%s/frontends/%s/tls-configs/%s", r.ServiceUUID, r.FrontendName, r.Name)
}

// DeleteLoadBalancerFrontendTLSConfigRequest represents a request to delete frontend TLS config
type DeleteLoadBalancerFrontendTLSConfigRequest struct {
	ServiceUUID  string `json:"-"`
	FrontendName string `json:"-"`
	Name         string `json:"-"`
}

func (r *DeleteLoadBalancerFrontendTLSConfigRequest) RequestURL() string {
	return fmt.Sprintf("/load-balancer/%s/frontends/%s/tls-configs/%s", r.ServiceUUID, r.FrontendName, r.Name)
}

// CreateLoadBalancerCertificateBundleRequest represents a request to create certificate bundle
type CreateLoadBalancerCertificateBundleRequest struct {
	Type upcloud.LoadBalancerCertificateBundleType `json:"type,omitempty"`

	Name          string   `json:"name,omitempty"`
	Certificate   string   `json:"certificate,omitempty"`
	Intermediates string   `json:"intermediates,omitempty"`
	PrivateKey    string   `json:"private_key,omitempty"`
	KeyType       string   `json:"key_type,omitempty"`
	Hostnames     []string `json:"hostnames,omitempty"`
}

func (r *CreateLoadBalancerCertificateBundleRequest) RequestURL() string {
	return loadBalancerCertificateBundleBaseURL
}

// ModifyLoadBalancerCertificateBundleRequest represents a request to modify certificate bundle
type ModifyLoadBalancerCertificateBundleRequest struct {
	UUID          string   `json:"-"`
	Name          string   `json:"name,omitempty"`
	Certificate   string   `json:"certificate,omitempty"`
	Intermediates string   `json:"intermediates,omitempty"`
	PrivateKey    string   `json:"private_key,omitempty"`
	Hostnames     []string `json:"hostnames,omitempty"`
}

func (r *ModifyLoadBalancerCertificateBundleRequest) RequestURL() string {
	return fmt.Sprintf("/load-balancer/certificate-bundles/%s", r.UUID)
}

// GetLoadBalancerCertificateBundlesRequest represents a request to list certificate bundles
type GetLoadBalancerCertificateBundlesRequest struct{}

func (r *GetLoadBalancerCertificateBundlesRequest) RequestURL() string {
	return loadBalancerCertificateBundleBaseURL
}

// GetLoadBalancerCertificateBundleRequest represents a request to get certificate bundle details
type GetLoadBalancerCertificateBundleRequest struct {
	UUID string `json:"-"`
}

func (r *GetLoadBalancerCertificateBundleRequest) RequestURL() string {
	return fmt.Sprintf("/load-balancer/certificate-bundles/%s", r.UUID)
}

// DeleteLoadBalancerCertificateBundleRequest represents a request to delete certificate bundle
type DeleteLoadBalancerCertificateBundleRequest struct {
	UUID string `json:"-"`
}

func (r *DeleteLoadBalancerCertificateBundleRequest) RequestURL() string {
	return fmt.Sprintf("/load-balancer/certificate-bundles/%s", r.UUID)
}
