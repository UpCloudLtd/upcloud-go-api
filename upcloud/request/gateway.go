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

type GatewayRouter struct {
	UUID string `json:"uuid,omitempty"`
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

type GatewayConnection struct {
	Name         string                        `json:"name,omitempty"`
	Type         upcloud.GatewayConnectionType `json:"type,omitempty"`
	LocalRoutes  []upcloud.GatewayRoute        `json:"local_routes,omitempty"`
	RemoteRoutes []upcloud.GatewayRoute        `json:"remote_routes,omitempty"`
	Tunnels      []GatewayTunnel               `json:"tunnels,omitempty"`
}

type GetGatewayConnectionsRequest struct {
	ServiceUUID string `json:"-"`
}

func (r *GetGatewayConnectionsRequest) RequestURL() string {
	return fmt.Sprintf("%s/%s/connections", gatewayBaseURL, r.ServiceUUID)
}

type GetGatewayConnectionRequest struct {
	ServiceUUID string `json:"-"`
	UUID        string `json:"-"`
}

func (r *GetGatewayConnectionRequest) RequestURL() string {
	return fmt.Sprintf("%s/%s/connections/%s", gatewayBaseURL, r.ServiceUUID, r.UUID)
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

type ModifyGatewayConnection struct {
	LocalRoutes  []upcloud.GatewayRoute `json:"local_routes,omitempty"`
	RemoteRoutes []upcloud.GatewayRoute `json:"remote_routes,omitempty"`
	Tunnels      []GatewayTunnel        `json:"tunnels,omitempty"`
}

type ModifyGatewayConnectionRequest struct {
	ServiceUUID string `json:"-"`
	UUID        string `json:"-"`
	Connection  ModifyGatewayConnection
}

func (r *ModifyGatewayConnectionRequest) RequestURL() string {
	return fmt.Sprintf("%s/%s/connections/%s", gatewayBaseURL, r.ServiceUUID, r.UUID)
}

func (r *ModifyGatewayConnectionRequest) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.Connection)
}

type DeleteGatewayConnectionRequest struct {
	ServiceUUID string `json:"-"`
	UUID        string `json:"-"`
}

func (r *DeleteGatewayConnectionRequest) RequestURL() string {
	return fmt.Sprintf("%s/%s/connections/%s", gatewayBaseURL, r.ServiceUUID, r.UUID)
}

type GatewayTunnel struct {
	Name             string                                `json:"name,omitempty"`
	LocalAddress     upcloud.GatewayTunnelLocalAddress     `json:"local_address,omitempty"`
	RemoteAddress    upcloud.GatewayTunnelRemoteAddress    `json:"remote_address,omitempty"`
	IPSec            upcloud.GatewayTunnelIPSec            `json:"ipsec,omitempty"`
	OperationalState upcloud.GatewayTunnelOperationalState `json:"operational_state,omitempty"`
}

type GetGatewayConnectionTunnelsRequest struct {
	ServiceUUID    string `json:"-"`
	ConnectionUUID string `json:"-"`
}

func (r *GetGatewayConnectionTunnelsRequest) RequestURL() string {
	return fmt.Sprintf("%s/%s/connections/%s/tunnels", gatewayBaseURL, r.ServiceUUID, r.ConnectionUUID)
}

type GetGatewayConnectionTunnelRequest struct {
	ServiceUUID    string `json:"-"`
	ConnectionUUID string `json:"-"`
	UUID           string `json:"-"`
}

func (r *GetGatewayConnectionTunnelRequest) RequestURL() string {
	return fmt.Sprintf("%s/%s/connections/%s/tunnels/%s", gatewayBaseURL, r.ServiceUUID, r.ConnectionUUID, r.UUID)
}

type CreateGatewayConnectionTunnelRequest struct {
	ServiceUUID    string `json:"-"`
	ConnectionUUID string `json:"-"`
	Tunnel         GatewayTunnel
}

func (r *CreateGatewayConnectionTunnelRequest) RequestURL() string {
	return fmt.Sprintf("%s/%s/connections/%s/tunnels", gatewayBaseURL, r.ServiceUUID, r.ConnectionUUID)
}

func (r *CreateGatewayConnectionTunnelRequest) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.Tunnel)
}

type ModifyGatewayTunnel struct {
	Name          string                              `json:"name,omitempty"`
	LocalAddress  *upcloud.GatewayTunnelLocalAddress  `json:"local_address,omitempty"`
	RemoteAddress *upcloud.GatewayTunnelRemoteAddress `json:"remote_address,omitempty"`
	IPSec         *upcloud.GatewayTunnelIPSec         `json:"ipsec,omitempty"`
}

type ModifyGatewayConnectionTunnelRequest struct {
	ServiceUUID    string `json:"-"`
	ConnectionUUID string `json:"-"`
	UUID           string `json:"-"`
	Tunnel         ModifyGatewayTunnel
}

func (r *ModifyGatewayConnectionTunnelRequest) RequestURL() string {
	return fmt.Sprintf("%s/%s/connections/%s/tunnels/%s", gatewayBaseURL, r.ServiceUUID, r.ConnectionUUID, r.UUID)
}

func (r *ModifyGatewayConnectionTunnelRequest) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.Tunnel)
}

type DeleteGatewayConnectionTunnelRequest struct {
	ServiceUUID    string `json:"-"`
	ConnectionUUID string `json:"-"`
	UUID           string `json:"-"`
}

func (r *DeleteGatewayConnectionTunnelRequest) RequestURL() string {
	return fmt.Sprintf("%s/%s/connections/%s/tunnels/%s", gatewayBaseURL, r.ServiceUUID, r.ConnectionUUID, r.UUID)
}
