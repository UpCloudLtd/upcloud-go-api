package request

import (
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
}

func (r *CreateGatewayRequest) RequestURL() string {
	return gatewayBaseURL
}

type ModifyGatewayRequest struct {
	UUID             string                          `json:"-"`
	Name             string                          `json:"name,omitempty"`
	ConfiguredStatus upcloud.GatewayConfiguredStatus `json:"configured_status,omitempty"`
	Labels           []upcloud.Label                 `json:"labels,omitempty"`
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
