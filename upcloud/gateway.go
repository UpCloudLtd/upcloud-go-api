package upcloud

import "time"

// GatewayConfiguredStatus represents a desired status of the service
type GatewayConfiguredStatus string

const (
	// GatewayStatusStarted represents a network gateway instance in started state
	GatewayStatusStarted GatewayConfiguredStatus = "started"
	// GatewayStatusStarted represents a network gateway instance in stopped state
	GatewayStatusStopped GatewayConfiguredStatus = "stopped"
)

// GatewayFeature represents a feature of the service
type GatewayFeature string

const (
	// GatewayFeatureNAT is the network address translation (NAT) feature of the network gateway
	GatewayFeatureNAT GatewayFeature = "nat"
)

type GatewayRouter struct {
	CreatedAt time.Time `json:"created_at,omitempty"`
	UUID      string    `json:"uuid,omitempty"`
}

type Gateway struct {
	ConfiguredStatus GatewayConfiguredStatus `json:"configured_status,omitempty"`
	CreatedAt        time.Time               `json:"created_at,omitempty"`
	Features         []GatewayFeature        `json:"features,omitempty"`
	Name             string                  `json:"name,omitempty"`
	OperationalState string                  `json:"operational_state,omitempty"`
	Routers          []GatewayRouter         `json:"routers,omitempty"`
	Labels           []Label                 `json:"labels,omitempty"`
	UpdatedAt        time.Time               `json:"updated_at,omitempty"`
	UUID             string                  `json:"uuid,omitempty"`
	Zone             string                  `json:"zone,omitempty"`
}
