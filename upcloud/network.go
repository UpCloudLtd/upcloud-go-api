package upcloud

import "encoding/xml"

// Constants
const (
	NetworkAccessPrivate = "private"
	NetworkAccessUtility = "utility"
	NetworkAccessPublic  = "public"
)

type Networks struct {
	Networks []Network `xml:"network"`
}

type IPNetwork struct {
	Address          string   `xml:"address"`
	DHCP             string   `xml:"dhcp"`
	DHCPDefaultRoute string   `xml:"dhcp_default_route"`
	DHCPDNS          []string `xml:"dhcp_dns"`
	Family           string   `xml:"family"`
	Gateway          string   `xml:"gateway"`
}

type Network struct {
	IPnetworks []IPNetwork     `xml:"ip_networks>ip_network,omitempty"`
	Name       string          `xml:"name"`
	Type       string          `xml:"type"`
	UUID       string          `xml:"uuid"`
	Zone       string          `xml:"zone"`
	Servers    []NetworkServer `xml:"servers>server"`
}

type NetworkResponse struct {
	Network Network `xml:"network"`
}

type NetworkServer struct {
	Title string `xml:"title"`
	UUID  string `xml:"uuid"`
}
type Interface struct {
	XmlName           xml.Name    `xml:"interface"`
	Index             int         `xml:"index"`
	IPAddresses       IPAddresses `xml:"ip_addresses>ip_address"`
	Mac               string      `xml:"mac"`
	Network           string      `xml:"network"`
	SourceIPFiltering string      `xml:"source_ip_filtering"`
	Type              string      `xml:"type"`
}

type ServerNetworkresponse struct {
	Networking Networking `xml:"networking"`
}
type Interfaces struct {
	Interface []Interface `xml:"interface"`
}
type Networking struct {
	Interfaces Interfaces `xml:"interfaces"`
}
