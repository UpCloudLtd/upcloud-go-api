package request

import (
	"encoding/json"
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

// CreatLoadBalancerBackend represents the payload for CreateLoadBalancerBackendRequest
type CreateLoadBalancerBackend struct {
	Name     string                            `json:"name"`
	Resolver string                            `json:"resolver,omitempty"`
	Members  []CreateLoadBalancerBackendMember `json:"members"`
}

type CreateLoadBalancerBackendRequest struct {
	ServiceUUID string `json:"-"`
	Payload     CreateLoadBalancerBackend
}

func (r *CreateLoadBalancerBackendRequest) RequestURL() string {
	return fmt.Sprintf("/loadbalancer/%s/backends", r.ServiceUUID)
}

func (r *CreateLoadBalancerBackendRequest) MarshalJSON() ([]byte, error) {
	return json.Marshal(&CreateLoadBalancerBackend{
		Name:     r.Payload.Name,
		Resolver: r.Payload.Resolver,
		Members:  r.Payload.Members,
	})
}

type GetLoadBalancerBackendDetailsRequest struct {
	ServiceUUID string `json:"-"`
	BackendName string `json:"-"`
}

func (r *GetLoadBalancerBackendDetailsRequest) RequestURL() string {
	return fmt.Sprintf("/loadbalancer/%s/backends/%s", r.ServiceUUID, r.BackendName)
}

// ModifyLoadBalancerBackend represents the payload for ModifyLoadBalancerBackendRequest
type ModifyLoadBalancerBackend struct {
	Name     string `json:"name,omitempty"`
	Resolver string `json:"resolver,omitempty"`
}

type ModifyLoadBalancerBackendRequest struct {
	ServiceUUID string `json:"-"`
	Name        string `json:"-"`
	Payload     ModifyLoadBalancerBackend
}

func (r *ModifyLoadBalancerBackendRequest) RequestURL() string {
	return fmt.Sprintf("/loadbalancer/%s/backends/%s", r.ServiceUUID, r.Name)
}

func (r *ModifyLoadBalancerBackendRequest) MarshalJSON() ([]byte, error) {
	return json.Marshal(&ModifyLoadBalancerBackend{
		Name:     r.Payload.Name,
		Resolver: r.Payload.Resolver,
	})
}

type DeleteLoadBalancerBackendRequest struct {
	ServiceUUID string `json:"-"`
	BackendName string `json:"-"`
}

func (r *DeleteLoadBalancerBackendRequest) RequestURL() string {
	return fmt.Sprintf("/loadbalancer/%s/backends/%s", r.ServiceUUID, r.BackendName)
}

// CreateLoadBalancerBackendMember represents the payload for CreateLoadBalancerBackendMemberRequest
type CreateLoadBalancerBackendMember struct {
	Name        string `json:"name"`
	Weight      int    `json:"weight"`
	MaxSessions int    `json:"max_sessions"`
	Enabled     bool   `json:"enabled"`
	Type        string `json:"type"`
	IP          string `json:"ip,omitempty"`
	Port        int    `json:"port,omitempty"`
	ServerUUID  string `json:"server_uuid,omitempty"`
}

type CreateLoadBalancerBackendMemberRequest struct {
	ServiceUUID string `json:"-"`
	BackendName string `json:"-"`
	Payload     CreateLoadBalancerBackendMember
}

func (r *CreateLoadBalancerBackendMemberRequest) RequestURL() string {
	return fmt.Sprintf("/loadbalancer/%s/backends/%s/members", r.ServiceUUID, r.BackendName)
}

func (r *CreateLoadBalancerBackendMemberRequest) MarshalJSON() ([]byte, error) {
	return json.Marshal(&CreateLoadBalancerBackendMember{
		Name:        r.Payload.Name,
		Weight:      r.Payload.Weight,
		MaxSessions: r.Payload.MaxSessions,
		Enabled:     r.Payload.Enabled,
		Type:        r.Payload.Type,
		IP:          r.Payload.IP,
		Port:        r.Payload.Port,
		ServerUUID:  r.Payload.ServerUUID,
	})
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

type ModifyLoadBalancerBackendMember = CreateLoadBalancerBackendMember

type ModifyLoadBalancerBackendMemberRequest struct {
	ServiceUUID string `json:"-"`
	BackendName string `json:"-"`
	Name        string `json:"-"`
	Payload     ModifyLoadBalancerBackendMember
}

func (r *ModifyLoadBalancerBackendMemberRequest) RequestURL() string {
	return fmt.Sprintf("/loadbalancer/%s/backends/%s/members/%s", r.ServiceUUID, r.BackendName, r.Name)
}

func (r *ModifyLoadBalancerBackendMemberRequest) MarshalJSON() ([]byte, error) {
	return json.Marshal(&ModifyLoadBalancerBackendMember{
		Name:        r.Payload.Name,
		Weight:      r.Payload.Weight,
		MaxSessions: r.Payload.MaxSessions,
		Enabled:     r.Payload.Enabled,
		Type:        r.Payload.Type,
		IP:          r.Payload.IP,
		Port:        r.Payload.Port,
		ServerUUID:  r.Payload.ServerUUID,
	})
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
