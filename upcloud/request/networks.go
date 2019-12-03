package request

import (
	"encoding/xml"
	"fmt"
)

// GetNetworksRequest represents a request to retrieve all networks
type GetNetworksRequest struct {
}

// RequestURL implements the Request interface
func (r *GetNetworksRequest) RequestURL() string {
	return fmt.Sprint("/network")
}

// GetNetworksRequest represents a request to retrieve all networks in a zone
type GetNetworksInZoneRequest struct {
	Zone string `xml:"zone"`
}

// RequestURL implements the Request interface
func (r *GetNetworksInZoneRequest) RequestURL() string {
	return fmt.Sprintf("/network/?zone=%s", r.Zone)
}

// FIXME due to what appears to be a bug, cant use XML, temporary JSON instead
// type IPNetworks struct {
// 	IPNetwork IPNetwork
// }

// type IPNetworks struct {
// XMLName          xml.Name `xml:"ip_network"`
// Address          string   `xml:"address"`
// DHCP             string   `xml:"dhcp"`
// DHCPDefaultRoute string   `xml:"dhcp_default_route,omitempty"`
// DHCPDNS          []string `xml:"dhcp_dns,omitempty"`
// Family           string   `xml:"family"`
// Gateway          string   `xml:"gateway,omitempty"`
// }

// GetIPAddressDetailsRequest represents a request to retrieve details about a specific IP address
// type CreateSDNPrivateNetworkRequest struct {
// XMLName    xml.Name          `xml:"network",json:"network"`
// Name       string            `xml:"name", json:"name"`
// Zone       string            `xml:"zone", json:"zone"`
// IPNetworks []CreateIPNetwork `xml:"ip_networks>ip_network", json:"ip_networks"`
// }

type IPNetworks struct {
	IPNetwork []IPNetwork `json:"ip_network"`
}

type IPNetwork struct {
	Address          string   `json:"address,omitempty"`
	DHCP             string   `json:"dhcp,omitempty"`
	DHCPDefaultRoute string   `json:"dhcp_default_route,omitempty"`
	DHCPDNS          []string `json:"dhcp_dns,omitempty"`
	Family           string   `json:"family,omitempty"`
	Gateway          string   `json:"gateway,omitempty"`
}

type Network struct {
	Name       string     `json:"name,omitempty"`
	Zone       string     `json:"zone,omitempty"`
	IPNetworks IPNetworks `json:"ip_networks"`
}

type CreateSDNPrivateNetworkRequest struct {
	Network Network `json:"network"`
}

// RequestURL implements the Request interface
func (r *CreateSDNPrivateNetworkRequest) RequestURL() string {
	return fmt.Sprint("/network")
}

type GetNetworkDetailsRequest struct {
	UUID string `xml:"uuid"`
}

// RequestURL implements the Request interface
func (r *GetNetworkDetailsRequest) RequestURL() string {
	return fmt.Sprintf("/network/%s", r.UUID)
}

// FIXME due to what appears to be a bug, cant use XML, temporary JSON instead
// type ModifyIPNetwork struct {
// 	DHCP             string   `xml:"dhcp"`
// 	DHCPDefaultRoute string   `xml:"dhcp_default_route,omitempty"`
// 	DHCPDNS          []string `xml:"dhcp_dns,omitempty"`
// 	Family           string   `xml:"family"`
// 	Gateway          string   `xml:"gateway,omitempty"`
// }

// type ModifyNetworkDetailsRequest struct {
// 	XMLName    xml.Name        `xml:"network"`
// 	Name       string          `xml:"name,omitempty"`
// 	UUID       string          `xml:"-"`
// 	IPNetworks ModifyIPNetwork `xml:"ip_networks>ip_network"`
// }

type ModifyNetworkDetailsRequest struct {
	Network Network `json:"network"`
	UUID    string  `json:"-"`
}

// RequestURL implements the Request interface
func (r *ModifyNetworkDetailsRequest) RequestURL() string {
	return fmt.Sprintf("/network/%s", r.UUID)
}

type DeleteNetworkRequest struct {
	UUID string `xml:"uuid"`
}

// RequestURL implements the Request interface
func (r *DeleteNetworkRequest) RequestURL() string {
	return fmt.Sprintf("/network/%s", r.UUID)
}

type CreateNetworkInterfaceRequest struct {
	Interface  Interface `json:"interface"`
	ServerUUID string    `json:"-"`
}
type IPAddress struct {
	Family  string `xml:"family"`
	Address string `xml:"address"`
}
type IPAddresses struct {
	IPAddress []IPAddress `xml:"ip_address,omitempty"`
}
type Interface struct {
	XMLName           xml.Name     `xml:"interface"`
	Type              string       `xml:"type,omitempty"`
	NetworkUUID       string       `xml:"network,omitempty"`
	IPAddresses       *IPAddresses `xml:"ip_addresses,omitempty"`
	SourceIPFiltering string       `xml:"source_ip_filtering,omitempty"`
}

func (r *CreateNetworkInterfaceRequest) RequestURL() string {
	return fmt.Sprintf("/server/%s/networking/interface", r.ServerUUID)
}

type ListServerNetworks struct {
	Interface  Interface `xml:"interface"`
	ServerUUID string    `xml:"-"`
}

func (r *ListServerNetworks) RequestURL() string {
	return fmt.Sprintf("/server/%s/networking", r.ServerUUID)
}

type ModifyNetworkInterfaceRequest struct {
	Interface  Interface `xml:"interface"`
	ServerUUID string    `xml:"-"`
	Index      int       `xml:"-"`
}

func (r *ModifyNetworkInterfaceRequest) RequestURL() string {
	return fmt.Sprintf("/server/%s/networking/interface/%d", r.ServerUUID, r.Index)
}

type DeleteNetworkInterfaceRequest struct {
	Index      int    `xml:"-"`
	ServerUUID string `xml:"-"`
}

func (r *DeleteNetworkInterfaceRequest) RequestURL() string {
	return fmt.Sprintf("/server/%s/networking/interface/%d", r.ServerUUID, r.Index)
}
