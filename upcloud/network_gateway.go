package upcloud

import "time"

// NetworkGatewayConfiguredStatus represents a desired status of the service
type NetworkGatewayConfiguredStatus string

const (
	// NetworkGatewayStatusStarted represents a network gateway instance in started state
	NetworkGatewayStatusStarted NetworkGatewayConfiguredStatus = "started"
	// NetworkGatewayStatusStarted represents a network gateway instance in stopped state
	NetworkGatewayStatusStopped NetworkGatewayConfiguredStatus = "stopped"
)

// NetworkGatewayFeature represents a feature of the service
type NetworkGatewayFeature string

const (
	// NetworkGatewayFeatureNAT is the network address translation (NAT) feature of the network gateway
	NetworkGatewayFeatureNAT NetworkGatewayFeature = "nat"
)

type NetworkGatewayRouter struct {
	CreatedAt time.Time `json:"created_at,omitempty"`
	UUID      string    `json:"uuid,omitempty"`
}

type NetworkGateway struct {
	ConfiguredStatus NetworkGatewayConfiguredStatus `json:"configured_status,omitempty"`
	CreatedAt        time.Time                      `json:"created_at,omitempty"`
	Features         []NetworkGatewayFeature        `json:"features,omitempty"`
	Name             string                         `json:"name,omitempty"`
	OperationalState string                         `json:"operational_state,omitempty"`
	Routers          []NetworkGatewayRouter         `json:"routers,omitempty"`
	Labels           []Label                        `json:"labels,omitempty"`
	UpdatedAt        time.Time                      `json:"updated_at,omitempty"`
	UUID             string                         `json:"uuid,omitempty"`
	Zone             string                         `json:"zone,omitempty"`
}
