package upcloud

import "time"

// GatewayConfiguredStatus represents a desired status of the service
type GatewayConfiguredStatus string

// GatewayConfiguredStatus represents a current actual status of the service
type GatewayOperationalState string

// GatewayFeature represents a feature of the service
type GatewayFeature string

const (
	GatewayConfiguredStatusStarted GatewayConfiguredStatus = "started"
	GatewayConfiguredStatusStopped GatewayConfiguredStatus = "stopped"

	GatewayOperationalStatePending           GatewayOperationalState = "pending"
	GatewayOperationalStateSetupAgent        GatewayOperationalState = "setup-agent"
	GatewayOperationalStateSetupLinkNetwork  GatewayOperationalState = "setup-link-network"
	GatewayOperationalStateSetupServer       GatewayOperationalState = "setup-server"
	GatewayOperationalStateSetupNetwork      GatewayOperationalState = "setup-network"
	GatewayOperationalStateSetupGW           GatewayOperationalState = "setup-gw"
	GatewayOperationalStateSetupDNS          GatewayOperationalState = "setup-dns"
	GatewayOperationalStateCheckup           GatewayOperationalState = "checkup"
	GatewayOperationalStateRunning           GatewayOperationalState = "running"
	GatewayOperationalStateDeleteDNS         GatewayOperationalState = "delete-dns"
	GatewayOperationalStateDeleteNetwork     GatewayOperationalState = "delete-network"
	GatewayOperationalStateDeleteServer      GatewayOperationalState = "delete-server"
	GatewayOperationalStateDeleteLinkNetwork GatewayOperationalState = "delete-link-network"
	GatewayOperationalStateDeleteService     GatewayOperationalState = "delete-service"

	// GatewayFeatureNAT is the network address translation (NAT) feature of the network gateway
	GatewayFeatureNAT GatewayFeature = "nat"
)

type Gateway struct {
	UUID             string                  `json:"uuid,omitempty"`
	Name             string                  `json:"name,omitempty"`
	Zone             string                  `json:"zone,omitempty"`
	Labels           []Label                 `json:"labels,omitempty"`
	ConfiguredStatus GatewayConfiguredStatus `json:"configured_status,omitempty"`
	OperationalState GatewayOperationalState `json:"operational_state,omitempty"`
	Features         []GatewayFeature        `json:"features,omitempty"`
	Routers          []GatewayRouter         `json:"routers,omitempty"`
	CreatedAt        time.Time               `json:"created_at,omitempty"`
	UpdatedAt        time.Time               `json:"updated_at,omitempty"`
	Addresses        []GatewayAddress        `json:"addresses,omitempty"`
}

type GatewayAddress struct {
	Address string `json:"address,omitempty"`
	Name    string `json:"name,omitempty"`
}

type GatewayRouter struct {
	CreatedAt time.Time `json:"created_at,omitempty"`
	UUID      string    `json:"uuid,omitempty"`
}
