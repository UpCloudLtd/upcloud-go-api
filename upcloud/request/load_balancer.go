package request

import (
	"fmt"

	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud"
)

type GetLoadBalancersRequest struct{}

type GetLoadBalancerDetailsRequest struct {
	UUID string
}

type CreateLoadBalancerRequest struct {
	Name             string                         `json:"name"`
	Plan             string                         `json:"plan"`
	Zone             string                         `json:"zone"`
	NetworkUuid      string                         `json:"network_uuid"`
	ConfiguredStatus string                         `json:"configured_status"`
	Frontends        []upcloud.LoadBalancerFrontend `json:"frontends"`
	Backends         []upcloud.LoadBalancerBackend  `json:"backends"`
	// Resolvers        []*upcloud.Resolver `json:"resolvers"` // TODO explore omit empty
}

type ModifyLoadBalancerRequest struct {
	UUID string `json:"-"`

	Name             string `json:"name,omitempty"`
	Plan             string `json:"plan,omitempty"`
	ConfiguredStatus string `json:"configured_status,omitempty"`
}

type DeleteLoadBalancerRequest struct {
	UUID string `json:"-"`
}

func (r *GetLoadBalancersRequest) RequestURL() string {
	return "/loadbalancer"
}

func (r *GetLoadBalancerDetailsRequest) RequestURL() string {
	return fmt.Sprintf("/loadbalancer/%s", r.UUID)
}

func (r *CreateLoadBalancerRequest) RequestURL() string {
	return "/loadbalancer"
}

func (r *ModifyLoadBalancerRequest) RequestURL() string {
	return fmt.Sprintf("/loadbalancer/%s", r.UUID)
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

type CreateLoadBalancerBackendRequest struct {
	ServiceUUID string                              `json:"-"`
	Name        string                              `json:"name"`
	Resolver    string                              `json:"resolver,omitempty"`
	Members     []upcloud.LoadBalancerBackendMember `json:"members"`
}

func (r *CreateLoadBalancerBackendRequest) RequestURL() string {
	return fmt.Sprintf("/loadbalancer/%s/backends", r.ServiceUUID)
}

type GetLoadBalancerBackendDetailsRequest struct {
	ServiceUUID string `json:"-"`
	BackendName string `json:"-"`
}

func (r *GetLoadBalancerBackendDetailsRequest) RequestURL() string {
	return fmt.Sprintf("/loadbalancer/%s/backends/%s", r.ServiceUUID, r.BackendName)
}

type ModifyLoadBalancerBackendRequest struct {
	ServiceUUID    string `json:"-"`
	BackendName    string `json:"-"`
	NewBackendName string `json:"name,omitempty"`
	Resolver       string `json:"resolver,omitempty"`
}

func (r *ModifyLoadBalancerBackendRequest) RequestURL() string {
	return fmt.Sprintf("/loadbalancer/%s/backends/%s", r.ServiceUUID, r.BackendName)
}

type DeleteLoadBalancerBackendRequest struct {
	ServiceUUID string `json:"-"`
	BackendName string `json:"-"`
}

func (r *DeleteLoadBalancerBackendRequest) RequestURL() string {
	return fmt.Sprintf("/loadbalancer/%s/backends/%s", r.ServiceUUID, r.BackendName)
}

type CreateLoadBalancerBackendMemberRequest struct {
	ServiceUUID       string `json:"-"`
	BackendName       string `json:"-"`
	MemberName        string `json:"name"`
	MemberWeight      int    `json:"weight"`
	MemberMaxSessions int    `json:"max_sessions"`
	MemberEnabled     bool   `json:"enabled"`
	MemberType        string `json:"type"`
	MemberIP          string `json:"ip,omitempty"`
	MemberPort        int    `json:"port,omitempty"`
	MemberServerUUID  string `json:"server_uuid,omitempty"`
}

func (r *CreateLoadBalancerBackendMemberRequest) RequestURL() string {
	return fmt.Sprintf("/loadbalancer/%s/backends/%s/members", r.ServiceUUID, r.BackendName)
}

type GetLoadBalancerBackendMembersRequest struct {
	ServiceUUID string `json:"-"`
	BackendName string `json:"-"`
}

func (r *GetLoadBalancerBackendMembersRequest) RequestURL() string {
	return fmt.Sprintf("/loadbalancer/%s/backends/%s/members", r.ServiceUUID, r.BackendName)
}

type GetLoadBalancerBackendMemberDetailsRequest struct {
	ServiceUUID string `json:"-"`
	BackendName string `json:"-"`
	MemberName  string `json:"-"`
}

func (r *GetLoadBalancerBackendMemberDetailsRequest) RequestURL() string {
	return fmt.Sprintf("/loadbalancer/%s/backends/%s/members/%s", r.ServiceUUID, r.BackendName, r.MemberName)
}

type ModifyLoadBalancerBackendMemberRequest struct {
	ServiceUUID       string `json:"-"`
	BackendName       string `json:"-"`
	MemberName        string `json:"-"`
	NewMemberName     string `json:"name,omitempty"`
	MemberWeight      int    `json:"weight,omitempty"`
	MemberMaxSessions int    `json:"max_sessions,omitempty"`
	MemberEnabled     bool   `json:"enabled,omitempty"`
	MemberIP          string `json:"ip,omitempty"`
	MemberPort        int    `json:"port,omitempty"`
	MemberType        string `json:"type,omitempty"`
	MemberServerUUID  string `json:"server_uuid,omitempty"`
}

func (r *ModifyLoadBalancerBackendMemberRequest) RequestURL() string {
	return fmt.Sprintf("/loadbalancer/%s/backends/%s/members/%s", r.ServiceUUID, r.BackendName, r.MemberName)
}

type DeleteLoadBalancerBackendMemberRequest struct {
	ServiceUUID string `json:"-"`
	BackendName string `json:"-"`
	MemberName  string `json:"-"`
}

func (r *DeleteLoadBalancerBackendMemberRequest) RequestURL() string {
	return fmt.Sprintf("/loadbalancer/%s/backends/%s/members/%s", r.ServiceUUID, r.BackendName, r.MemberName)
}

type CreateLoadBalancerResolverRequest struct {
	ServiceUUID  string   `json:"-"`
	Name         string   `json:"name,omitempty"`
	Nameservers  []string `json:"nameservers,omitempty"`
	Retries      int      `json:"retries,omitempty"`
	Timeout      int      `json:"timeout,omitempty"`
	TimeoutRetry int      `json:"timeout_retry,omitempty"`
	CacheValid   int      `json:"cache_valid,omitempty"`
	CacheInvalid int      `json:"cache_invalid,omitempty"`
}

func (r *CreateLoadBalancerResolverRequest) RequestURL() string {
	return fmt.Sprintf("/loadbalancer/%s/resolvers", r.ServiceUUID)
}

type GetLoadBalancerResolversRequest struct {
	ServiceUUID string `json:"-"`
}

func (r *GetLoadBalancerResolversRequest) RequestURL() string {
	return fmt.Sprintf("/loadbalancer/%s/resolvers", r.ServiceUUID)
}

type GetLoadBalancerResolverDetailsRequest struct {
	ServiceUUID  string `json:"-"`
	ResolverName string `json:"-"`
}

func (r *GetLoadBalancerResolverDetailsRequest) RequestURL() string {
	return fmt.Sprintf("/loadbalancer/%s/resolvers/%s", r.ServiceUUID, r.ResolverName)
}

type ModifyLoadBalancerRevolverRequest struct {
	ServiceUUID     string   `json:"-"`
	ResolverName    string   `json:"-"`
	NewResolverName string   `json:"name,omitempty"`
	Nameservers     []string `json:"nameservers,omitempty"`
	Retries         int      `json:"retries,omitempty"`
	Timeout         int      `json:"timeout,omitempty"`
	TimeoutRetry    int      `json:"timeout_retry,omitempty"`
	CacheValid      int      `json:"cache_valid,omitempty"`
	CacheInvalid    int      `json:"cache_invalid,omitempty"`
}

func (r *ModifyLoadBalancerRevolverRequest) RequestURL() string {
	return fmt.Sprintf("/loadbalancer/%s/resolvers/%s", r.ServiceUUID, r.ResolverName)
}

type DeleteLoadBalancerResolverRequest struct {
	ServiceUUID  string `json:"-"`
	ResolverName string `json:"-"`
}

func (r *DeleteLoadBalancerResolverRequest) RequestURL() string {
	return fmt.Sprintf("/loadbalancer/%s/resolvers/%s", r.ServiceUUID, r.ResolverName)
}
