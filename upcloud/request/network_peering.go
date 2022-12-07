package request

import (
	"encoding/json"
	"fmt"

	"github.com/UpCloudLtd/upcloud-go-api/v5/upcloud"
)

const networkPeeringBaseURL string = "/network-peering"

type GetNetworkPeeringRequest struct {
	UUID string
}

func (r *GetNetworkPeeringRequest) RequestURL() string {
	return fmt.Sprintf("%s/%s", networkPeeringBaseURL, r.UUID)
}

type NetworkPeeringNetwork struct {
	UUID string `json:"uuid,omitempty"`
}

type CreateNetworkPeeringRequest struct {
	Name             string                                 `json:"name,omitempty"`
	ConfiguredStatus upcloud.NetworkPeeringConfiguredStatus `json:"configured_status,omitempty"`
	Network          NetworkPeeringNetwork                  `json:"network,omitempty"`
	PeerNetwork      NetworkPeeringNetwork                  `json:"peer_network,omitempty"`
}

func (r *CreateNetworkPeeringRequest) MarshalJSON() ([]byte, error) {
	type rt CreateNetworkPeeringRequest
	v := struct {
		NetworkPeering rt `json:"network_peering,omitempty"`
	}{
		NetworkPeering: rt(*r),
	}
	return json.Marshal(&v)
}

func (r *CreateNetworkPeeringRequest) RequestURL() string {
	return networkPeeringBaseURL
}

type ModifyNetworkPeering struct {
	Name             string                                 `json:"name,omitempty"`
	ConfiguredStatus upcloud.NetworkPeeringConfiguredStatus `json:"configured_status,omitempty"`
}

type ModifyNetworkPeeringRequest struct {
	UUID           string               `json:"-"`
	NetworkPeering ModifyNetworkPeering `json:"network_peering,omitempty"`
}

func (r *ModifyNetworkPeeringRequest) RequestURL() string {
	return fmt.Sprintf("%s/%s", networkPeeringBaseURL, r.UUID)
}

type DeleteNetworkPeeringRequest GetNetworkPeeringRequest

func (r *DeleteNetworkPeeringRequest) RequestURL() string {
	return fmt.Sprintf("%s/%s", networkPeeringBaseURL, r.UUID)
}
