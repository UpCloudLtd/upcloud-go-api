package upcloud

import (
	"github.com/google/uuid"
	"time"
)

type LoadBalancers struct {
	LoadBalancers []LoadBalancer `json:"-"`
}

type LoadBalancerDetails struct {
	LoadBalancer
}

type Frontend struct {
	Name string
	Mode string
	Port int
}

type Rule struct {
	Name     string
	Priority int
	Matchers []*Matchers
	Actions  []*Actions
}

type Matchers struct {
}

type Actions struct {
}

type Backend struct {
	Name     string   `json:"name"`
	Members  []Member `json:"members"`
	Resolver string   `json:"resolver,omitempty"`
}

type Member struct {
	ServerUuid  uuid.UUID `json:"server_uuid,omitempty"`
	Name        string    `json:"name"`
	Ip          string    `json:"ip"`
	Port        int       `json:"port"`
	Weight      int       `json:"weight"`
	MaxSessions int       `json:"max_sessions"`
	Type        string    `json:"type"`
	Enabled     bool      `json:"enabled"`
}

type Resolver struct {
	Name        string
	NameServers []string
}

type LoadBalancer struct {
	Uuid                     uuid.UUID  `json:"uuid"`
	Name                     string     `json:"name"`
	Zone                     string     `json:"zone"`
	Plan                     string     `json:"plan"`
	PlanCoreNumber           int        `json:"plan_core_number"`
	PlanMemoryAmount         int        `json:"plan_memory_amount"`
	PlanPerServerMaxSessions int        `json:"plan_per_server_max_sessions"`
	MainAccountId            int        `json:"main_account_id,omitempty"`
	NetworkUuid              uuid.UUID  `json:"network_uuid"`
	FloatingIps              []string   `json:"floating_ips,omitempty"`
	DnsName                  string     `json:"dns_name"`
	ConfiguredStatus         string     `json:"configured_status"`
	OperationalState         string     `json:"operational_state"`
	Frontends                []Frontend `json:"frontends"`
	Backends                 []Backend  `json:"backends"`
	Resolvers                []Resolver `json:"resolvers"`
	Deleted                  bool       `json:"deleted"`
	CreatedAt                time.Time  `json:"created_at"`
	UpdatedAt                time.Time  `json:"updated_at"`
}
