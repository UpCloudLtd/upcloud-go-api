package request

import (
	"encoding/json"
	"fmt"

	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud"
)

type GetLoadBalancersRequest struct{}

func (r *GetLoadBalancersRequest) RequestURL() string {
	return "/loadbalancer"
}

type GetLoadBalancerRequest struct {
	UUID string
}

func (r *GetLoadBalancerRequest) RequestURL() string {
	return fmt.Sprintf("/loadbalancer/%s", r.UUID)
}

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
	return "/loadbalancer"
}

type ModifyLoadBalancerRequest struct {
	UUID             string `json:"-"`
	Name             string `json:"name,omitempty"`
	Plan             string `json:"plan,omitempty"`
	ConfiguredStatus string `json:"configured_status,omitempty"`
}

func (r *ModifyLoadBalancerRequest) RequestURL() string {
	return fmt.Sprintf("/loadbalancer/%s", r.UUID)
}

type DeleteLoadBalancerRequest struct {
	UUID string `json:"-"`
}

func (r *DeleteLoadBalancerRequest) RequestURL() string {
	return fmt.Sprintf("/loadbalancer/%s", r.UUID)
}

type GetLoadBalancerBackendsRequest struct {
	ServiceUUID string `json:"-"`
}

func (r *GetLoadBalancerBackendsRequest) RequestURL() string {
	return fmt.Sprintf("/loadbalancer/%s/backends", r.ServiceUUID)
}

// BalancerBackend represents the payload for CreateLoadBalancerBackendRequest
type LoadBalancerBackend struct {
	Name     string                      `json:"name"`
	Resolver string                      `json:"resolver,omitempty"`
	Members  []LoadBalancerBackendMember `json:"members"`
}

type CreateLoadBalancerBackendRequest struct {
	ServiceUUID string `json:"-"`
	Backend     LoadBalancerBackend
}

func (r *CreateLoadBalancerBackendRequest) RequestURL() string {
	return fmt.Sprintf("/loadbalancer/%s/backends", r.ServiceUUID)
}

func (r *CreateLoadBalancerBackendRequest) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.Backend)
}

type GetLoadBalancerBackendRequest struct {
	ServiceUUID string `json:"-"`
	Name        string `json:"-"`
}

func (r *GetLoadBalancerBackendRequest) RequestURL() string {
	return fmt.Sprintf("/loadbalancer/%s/backends/%s", r.ServiceUUID, r.Name)
}

// ModifyLoadBalancerBackend represents the payload for ModifyLoadBalancerBackendRequest
type ModifyLoadBalancerBackend struct {
	Name     string `json:"name,omitempty"`
	Resolver string `json:"resolver,omitempty"`
}

type ModifyLoadBalancerBackendRequest struct {
	ServiceUUID string `json:"-"`
	Name        string `json:"-"`
	Backend     ModifyLoadBalancerBackend
}

func (r *ModifyLoadBalancerBackendRequest) RequestURL() string {
	return fmt.Sprintf("/loadbalancer/%s/backends/%s", r.ServiceUUID, r.Name)
}

func (r *ModifyLoadBalancerBackendRequest) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.Backend)
}

type DeleteLoadBalancerBackendRequest struct {
	ServiceUUID string `json:"-"`
	Name        string `json:"-"`
}

func (r *DeleteLoadBalancerBackendRequest) RequestURL() string {
	return fmt.Sprintf("/loadbalancer/%s/backends/%s", r.ServiceUUID, r.Name)
}

// LoadBalancerBackendMember represents the payload for CreateLoadBalancerBackendMemberRequest
type LoadBalancerBackendMember struct {
	Name        string                                `json:"name,omitempty"`
	Weight      int                                   `json:"weight,omitempty"`
	MaxSessions int                                   `json:"max_sessions,omitempty"`
	Enabled     bool                                  `json:"enabled,omitempty"`
	Type        upcloud.LoadBalancerBackendMemberType `json:"type,omitempty"`
	IP          string                                `json:"ip,omitempty"`
	Port        int                                   `json:"port,omitempty"`
}

type CreateLoadBalancerBackendMemberRequest struct {
	ServiceUUID string `json:"-"`
	BackendName string `json:"-"`
	Member      LoadBalancerBackendMember
}

func (r *CreateLoadBalancerBackendMemberRequest) RequestURL() string {
	return fmt.Sprintf("/loadbalancer/%s/backends/%s/members", r.ServiceUUID, r.BackendName)
}

func (r *CreateLoadBalancerBackendMemberRequest) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.Member)
}

type GetLoadBalancerBackendMembersRequest struct {
	ServiceUUID string `json:"-"`
	BackendName string `json:"-"`
}

func (r *GetLoadBalancerBackendMembersRequest) RequestURL() string {
	return fmt.Sprintf("/loadbalancer/%s/backends/%s/members", r.ServiceUUID, r.BackendName)
}

type GetLoadBalancerBackendMemberRequest struct {
	ServiceUUID string `json:"-"`
	BackendName string `json:"-"`
	Name        string `json:"-"`
}

func (r *GetLoadBalancerBackendMemberRequest) RequestURL() string {
	return fmt.Sprintf("/loadbalancer/%s/backends/%s/members/%s", r.ServiceUUID, r.BackendName, r.Name)
}

type ModifyLoadBalancerBackendMemberRequest struct {
	ServiceUUID string `json:"-"`
	BackendName string `json:"-"`
	Name        string `json:"-"`
	Member      LoadBalancerBackendMember
}

func (r *ModifyLoadBalancerBackendMemberRequest) RequestURL() string {
	return fmt.Sprintf("/loadbalancer/%s/backends/%s/members/%s", r.ServiceUUID, r.BackendName, r.Name)
}

func (r *ModifyLoadBalancerBackendMemberRequest) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.Member)
}

type DeleteLoadBalancerBackendMemberRequest struct {
	ServiceUUID string `json:"-"`
	BackendName string `json:"-"`
	Name        string `json:"-"`
}

func (r *DeleteLoadBalancerBackendMemberRequest) RequestURL() string {
	return fmt.Sprintf("/loadbalancer/%s/backends/%s/members/%s", r.ServiceUUID, r.BackendName, r.Name)
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

type CreateLoadBalancerResolverRequest struct {
	ServiceUUID string `json:"-"`
	Resolver    LoadBalancerResolver
}

func (r *CreateLoadBalancerResolverRequest) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.Resolver)
}

func (r *CreateLoadBalancerResolverRequest) RequestURL() string {
	return fmt.Sprintf("/loadbalancer/%s/resolvers", r.ServiceUUID)
}

type ModifyLoadBalancerRevolverRequest struct {
	ServiceUUID string `json:"-"`
	Name        string `json:"-"`
	Resolver    LoadBalancerResolver
}

func (r *ModifyLoadBalancerRevolverRequest) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.Resolver)
}

func (r *ModifyLoadBalancerRevolverRequest) RequestURL() string {
	return fmt.Sprintf("/loadbalancer/%s/resolvers/%s", r.ServiceUUID, r.Name)
}

type GetLoadBalancerResolversRequest struct {
	ServiceUUID string `json:"-"`
}

func (r *GetLoadBalancerResolversRequest) RequestURL() string {
	return fmt.Sprintf("/loadbalancer/%s/resolvers", r.ServiceUUID)
}

type GetLoadBalancerResolverRequest struct {
	ServiceUUID string `json:"-"`
	Name        string `json:"-"`
}

func (r *GetLoadBalancerResolverRequest) RequestURL() string {
	return fmt.Sprintf("/loadbalancer/%s/resolvers/%s", r.ServiceUUID, r.Name)
}

type DeleteLoadBalancerResolverRequest struct {
	ServiceUUID string `json:"-"`
	Name        string `json:"-"`
}

func (r *DeleteLoadBalancerResolverRequest) RequestURL() string {
	return fmt.Sprintf("/loadbalancer/%s/resolvers/%s", r.ServiceUUID, r.Name)
}

type GetLoadBalancerPlansRequest struct{}

func (r *GetLoadBalancerPlansRequest) RequestURL() string {
	return "/loadbalancer/plans"
}

type GetLoadBalancerFrontendsRequest struct {
	ServiceUUID string `json:"-"`
}

func (r *GetLoadBalancerFrontendsRequest) RequestURL() string {
	return fmt.Sprintf("/loadbalancer/%s/frontends", r.ServiceUUID)
}

type GetLoadBalancerFrontendRequest struct {
	ServiceUUID string `json:"-"`
	Name        string `json:"-"`
}

func (r *GetLoadBalancerFrontendRequest) RequestURL() string {
	return fmt.Sprintf("/loadbalancer/%s/frontends/%s", r.ServiceUUID, r.Name)
}

// LoadBalancerFrontendRule represents frontend rule payload
type LoadBalancerFrontendRule struct {
	Name     string                        `json:"name,omitempty"`
	Priority int                           `json:"priority,omitempty"`
	Matchers []upcloud.LoadBalancerMatcher `json:"matchers,omitempty"`
	Actions  []upcloud.LoadBalancerAction  `json:"actions,omitempty"`
}

// LoadBalancerTLSConfig represents TLS config payload
type LoadBalancerTLSConfig struct {
	Name                  string `json:"name,omitempty"`
	CertificateBundleUUID string `json:"certificate_bundle_uuid,omitempty"`
}

// LoadBalancerFrontend represents frontend payload
type LoadBalancerFrontend struct {
	Name           string                     `json:"name,omitempty"`
	Mode           upcloud.LoadBalancerMode   `json:"mode,omitempty"`
	Port           int                        `json:"port,omitempty"`
	DefaultBackend string                     `json:"default_backend,omitempty"`
	Rules          []LoadBalancerFrontendRule `json:"rules,omitempty"`
	TLSConfigs     []LoadBalancerTLSConfig    `json:"tls_configs,omitempty"`
}

type CreateLoadBalancerFrontendRequest struct {
	ServiceUUID string
	Frontend    LoadBalancerFrontend
}

func (r *CreateLoadBalancerFrontendRequest) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.Frontend)
}

func (r *CreateLoadBalancerFrontendRequest) RequestURL() string {
	return fmt.Sprintf("/loadbalancer/%s/frontends", r.ServiceUUID)
}

type ModifyLoadBalancerFrontend struct {
	Name           string                   `json:"name,omitempty"`
	Mode           upcloud.LoadBalancerMode `json:"mode,omitempty"`
	Port           int                      `json:"port,omitempty"`
	DefaultBackend string                   `json:"default_backend,omitempty"`
}

type ModifyLoadBalancerFrontendRequest struct {
	ServiceUUID string `json:"-"`
	Name        string `json:"-"`
	Frontend    ModifyLoadBalancerFrontend
}

func (r *ModifyLoadBalancerFrontendRequest) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.Frontend)
}

func (r *ModifyLoadBalancerFrontendRequest) RequestURL() string {
	return fmt.Sprintf("/loadbalancer/%s/frontends/%s", r.ServiceUUID, r.Name)
}

type DeleteLoadBalancerFrontendRequest struct {
	ServiceUUID string `json:"-"`
	Name        string `json:"-"`
}

func (r *DeleteLoadBalancerFrontendRequest) RequestURL() string {
	return fmt.Sprintf("/loadbalancer/%s/frontends/%s", r.ServiceUUID, r.Name)
}
