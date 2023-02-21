package request

import (
	"fmt"

	"github.com/UpCloudLtd/upcloud-go-api/v5/upcloud"
)

const networkGatewayBaseURL string = "/gateway"

type GetNetworkGatewaysRequest struct {
	Filters []QueryFilter
}

func (r *GetNetworkGatewaysRequest) RequestURL() string {
	if len(r.Filters) == 0 {
		return networkGatewayBaseURL
	}

	return fmt.Sprintf("%s?%s", networkGatewayBaseURL, encodeQueryFilters(r.Filters))
}

type GetNetworkGatewayRequest struct {
	UUID string
}

func (r *GetNetworkGatewayRequest) RequestURL() string {
	return fmt.Sprintf("%s/%s", networkGatewayBaseURL, r.UUID)
}

type NetworkGatewayRouter struct {
	UUID string `json:"uuid,omitempty"`
}

type CreateNetworkGatewayRequest struct {
	Name             string                                 `json:"name,omitempty"`
	Zone             string                                 `json:"zone,omitempty"`
	Features         []upcloud.NetworkGatewayFeature        `json:"features,omitempty"`
	Routers          []NetworkGatewayRouter                 `json:"routers,omitempty"`
	Labels           []upcloud.Label                        `json:"labels,omitempty"`
	ConfiguredStatus upcloud.NetworkGatewayConfiguredStatus `json:"configured_status,omitempty"`
}

func (r *CreateNetworkGatewayRequest) RequestURL() string {
	return networkGatewayBaseURL
}

type ModifyNetworkGatewayRequest struct {
	UUID             string                                 `json:"-"`
	Name             string                                 `json:"name,omitempty"`
	ConfiguredStatus upcloud.NetworkGatewayConfiguredStatus `json:"configured_status,omitempty"`
	Labels           []upcloud.Label                        `json:"labels,omitempty"`
}

func (r *ModifyNetworkGatewayRequest) RequestURL() string {
	return fmt.Sprintf("%s/%s", networkGatewayBaseURL, r.UUID)
}

type DeleteNetworkGatewayRequest GetNetworkGatewayRequest

func (r *DeleteNetworkGatewayRequest) RequestURL() string {
	return fmt.Sprintf("%s/%s", networkGatewayBaseURL, r.UUID)
}
