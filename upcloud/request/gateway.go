package request

import (
	"encoding/json"
	"fmt"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
)

const gatewayBaseURL string = "/gateway"

type GetGatewayPlansRequest struct{}

func (r *GetGatewayPlansRequest) RequestURL() string {
	return fmt.Sprintf("%s/plans", gatewayBaseURL)
}

type GetGatewaysRequest struct {
	Filters []QueryFilter
}

func (r *GetGatewaysRequest) RequestURL() string {
	if len(r.Filters) == 0 {
		return gatewayBaseURL
	}

	return fmt.Sprintf("%s?%s", gatewayBaseURL, encodeQueryFilters(r.Filters))
}

type GetGatewayRequest struct {
	UUID string
}

func (r *GetGatewayRequest) RequestURL() string {
	return fmt.Sprintf("%s/%s", gatewayBaseURL, r.UUID)
}

type CreateGatewayRequest struct {
	Name             string                          `json:"name,omitempty"`
	Zone             string                          `json:"zone,omitempty"`
	Features         []upcloud.GatewayFeature        `json:"features,omitempty"`
	Routers          []GatewayRouter                 `json:"routers,omitempty"`
	Labels           []upcloud.Label                 `json:"labels,omitempty"`
	ConfiguredStatus upcloud.GatewayConfiguredStatus `json:"configured_status,omitempty"`
	Plan             string                          `json:"plan,omitempty"`
	Addresses        []upcloud.GatewayAddress        `json:"addresses,omitempty"`
	Connections      []GatewayConnection             `json:"connections,omitempty"`
}

func (r *CreateGatewayRequest) RequestURL() string {
	return gatewayBaseURL
}

type ModifyGatewayRequest struct {
	UUID             string                          `json:"-"`
	Name             string                          `json:"name,omitempty"`
	Plan             string                          `json:"plan,omitempty"`
	ConfiguredStatus upcloud.GatewayConfiguredStatus `json:"configured_status,omitempty"`
	Labels           []upcloud.Label                 `json:"labels,omitempty"`
	Connections      []GatewayConnection             `json:"connections,omitempty"`
}

func (r *ModifyGatewayRequest) RequestURL() string {
	return fmt.Sprintf("%s/%s", gatewayBaseURL, r.UUID)
}

type DeleteGatewayRequest struct {
	UUID string
}

func (r *DeleteGatewayRequest) RequestURL() string {
	return fmt.Sprintf("%s/%s", gatewayBaseURL, r.UUID)
}

type GatewayRouter struct {
	UUID string `json:"uuid,omitempty"`
}

type GatewayConnection struct {
	Name         string                        `json:"name,omitempty"`
	Type         upcloud.GatewayConnectionType `json:"type,omitempty"`
	LocalRoutes  []upcloud.GatewayRoute        `json:"local_routes,omitempty"`
	RemoteRoutes []upcloud.GatewayRoute        `json:"remote_routes,omitempty"`
	Tunnels      []GatewayTunnel               `json:"tunnels,omitempty"`
}

type GatewayTunnel struct {
	Name             string                                `json:"name,omitempty"`
	LocalAddress     upcloud.GatewayTunnelLocalAddress     `json:"local_address,omitempty"`
	RemoteAddress    upcloud.GatewayTunnelRemoteAddress    `json:"remote_address,omitempty"`
	IPSec            upcloud.GatewayTunnelIPSec            `json:"ipsec,omitempty"`
	OperationalState upcloud.GatewayTunnelOperationalState `json:"operational_state,omitempty"`
}

type GetGatewayConnectionsRequest struct {
	ServiceUUID string `json:"-"`
}

func (r *GetGatewayConnectionsRequest) RequestURL() string {
	return fmt.Sprintf("%s/%s/connections", gatewayBaseURL, r.ServiceUUID)
}

type GetGatewayConnectionRequest struct {
	ServiceUUID string `json:"-"`
	Name        string `json:"-"`
}

func (r *GetGatewayConnectionRequest) RequestURL() string {
	return fmt.Sprintf("%s/%s/connections/%s", gatewayBaseURL, r.ServiceUUID, r.Name)
}

type CreateGatewayConnectionRequest struct {
	ServiceUUID string `json:"-"`
	Connection  GatewayConnection
}

func (r *CreateGatewayConnectionRequest) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.Connection)
}

func (r *CreateGatewayConnectionRequest) RequestURL() string {
	return fmt.Sprintf("%s/%s/connections/", gatewayBaseURL, r.ServiceUUID)
}

type DeleteGatewayConnectionRequest struct {
	ServiceUUID string `json:"-"`
	Name        string `json:"-"`
}

func (r *DeleteGatewayConnectionRequest) RequestURL() string {
	return fmt.Sprintf("%s/%s/connections/%s", gatewayBaseURL, r.ServiceUUID, r.Name)
}
