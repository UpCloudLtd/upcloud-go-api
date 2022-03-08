package upcloud

import (
	"time"

	"github.com/google/uuid"
)

// LoadBalancerPlan represents load balancer plan details
type LoadBalancerPlan struct {
	Name                 string
	PerServerMaxSessions int
	ServerNumber         int
}

type LoadBalancerFrontend struct {
	Name string
	Mode string
	Port int
}

type LoadBalancerRule struct {
	Name     string
	Priority int
	Matchers []*LoadBalancerMatchers
	Actions  []*LoadBalancerActions
}

type LoadBalancerMatchers struct {
}

type LoadBalancerActions struct {
}

type LoadBalancerBackend struct {
	Name     string               `json:"name"`
	Members  []LoadBalancerMember `json:"members"`
	Resolver string               `json:"resolver,omitempty"`
}

type LoadBalancerMember struct {
	ServerUuid  uuid.UUID `json:"server_uuid,omitempty"`
	Name        string    `json:"name"`
	Ip          string    `json:"ip"`
	Port        int       `json:"port"`
	Weight      int       `json:"weight"`
	MaxSessions int       `json:"max_sessions"`
	Type        string    `json:"type"`
	Enabled     bool      `json:"enabled"`
}

type LoadBalancerResolver struct {
	Name        string
	NameServers []string
}

type LoadBalancer struct {
	Uuid                     uuid.UUID              `json:"uuid"`
	Name                     string                 `json:"name"`
	Zone                     string                 `json:"zone"`
	Plan                     string                 `json:"plan"`
	PlanCoreNumber           int                    `json:"plan_core_number"`
	PlanMemoryAmount         int                    `json:"plan_memory_amount"`
	PlanPerServerMaxSessions int                    `json:"plan_per_server_max_sessions"`
	MainAccountId            int                    `json:"main_account_id,omitempty"`
	NetworkUuid              uuid.UUID              `json:"network_uuid"`
	FloatingIps              []string               `json:"floating_ips,omitempty"`
	DnsName                  string                 `json:"dns_name"`
	ConfiguredStatus         string                 `json:"configured_status"`
	OperationalState         string                 `json:"operational_state"`
	Frontends                []LoadBalancerFrontend `json:"frontends"`
	Backends                 []LoadBalancerBackend  `json:"backends"`
	Resolvers                []LoadBalancerResolver `json:"resolvers"`
	Deleted                  bool                   `json:"deleted"`
	CreatedAt                time.Time              `json:"created_at"`
	UpdatedAt                time.Time              `json:"updated_at"`
}
