package upcloud

import "encoding/json"

type NetworkPeeringConfiguredStatus string

type NetworkPeeringIPNetworkFamily string

type NetworkPeeringState string

const (
	NetworkPeeringConfiguredStatusActive   NetworkPeeringConfiguredStatus = "active"
	NetworkPeeringConfiguredStatusDisabled NetworkPeeringConfiguredStatus = "disabled"

	NetworkPeeringStateActive             NetworkPeeringState = "active"
	NetworkPeeringStatePendingPeer        NetworkPeeringState = "pending-peer"
	NetworkPeeringStateProvisioning       NetworkPeeringState = "provisioning"
	NetworkPeeringStateConflictSubnet     NetworkPeeringState = "conflict-subnet"
	NetworkPeeringStateMissingLocalRouter NetworkPeeringState = "missing-local-router"
	NetworkPeeringStateMissingPeerRouter  NetworkPeeringState = "missing-peer-router"
	NetworkPeeringStateDeletedPeerNetwork NetworkPeeringState = "deleted-peer-network"
	NetworkPeeringStateDisabled           NetworkPeeringState = "disabled"
	NetworkPeeringStatePeerDisabled       NetworkPeeringState = "peer-disabled"
	NetworkPeeringStateError              NetworkPeeringState = "error"

	NetworkPeeringIPNetworkFamilyIPv4 NetworkPeeringIPNetworkFamily = "IPv4"
	NetworkPeeringIPNetworkFamilyIPv6 NetworkPeeringIPNetworkFamily = "IPv6"
)

type NetworkPeerings []NetworkPeering

// UnmarshalJSON is a custom unmarshaller that deals with deeply embedded values.
func (n *NetworkPeerings) UnmarshalJSON(b []byte) error {
	type np NetworkPeering
	v := struct {
		NetworkPeerings struct {
			NetworkPeering []np `json:"network_peering,omitempty"`
		} `json:"network_peerings,omitempty"`
	}{}

	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	p := make([]NetworkPeering, 0)
	for _, r := range v.NetworkPeerings.NetworkPeering {
		p = append(p, NetworkPeering(r))
	}
	*n = p
	return nil
}

type NetworkPeering struct {
	UUID             string                         `json:"uuid,omitempty"`
	ConfiguredStatus NetworkPeeringConfiguredStatus `json:"configured_status,omitempty"`
	Name             string                         `json:"name,omitempty"`
	Network          NetworkPeeringNetwork          `json:"network,omitempty"`
	PeerNetwork      NetworkPeeringNetwork          `json:"peer_network,omitempty"`
	State            NetworkPeeringState            `json:"state,omitempty"`
}

// UnmarshalJSON is a custom unmarshaller that deals with deeply embedded values.
func (n *NetworkPeering) UnmarshalJSON(b []byte) error {
	type np NetworkPeering
	v := struct {
		NetworkPeering np `json:"network_peering"`
	}{}

	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}

	*n = NetworkPeering(v.NetworkPeering)
	return nil
}

type NetworkPeeringNetwork struct {
	UUID       string `json:"uuid,omitempty"`
	IPNetworks []NetworkPeeringIPNetwork
}

// UnmarshalJSON is a custom unmarshaller that deals with deeply embedded values.
func (n *NetworkPeeringNetwork) UnmarshalJSON(b []byte) error {
	v := struct {
		UUID       string
		IPNetworks struct {
			IPNetwork []NetworkPeeringIPNetwork `json:"ip_network,omitempty"`
		} `json:"ip_networks,omitempty"`
	}{}

	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}

	*n = NetworkPeeringNetwork{UUID: v.UUID, IPNetworks: v.IPNetworks.IPNetwork}
	return nil
}

type NetworkPeeringIPNetwork struct {
	Address string                        `json:"address,omitempty"`
	Family  NetworkPeeringIPNetworkFamily `json:"family,omitempty"`
}
