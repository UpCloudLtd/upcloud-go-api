package upcloud

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
	IPnetworks []IPNetwork     `xml:"ip_networks>ip_network"`
	Name       string          `xml:"name"`
	Type       string          `xml:"type"`
	UUID       string          `xml:"uuid"`
	Zone       string          `xml:"zone"`
	Servers    []NetworkServer `xml:"servers>server"`
}

type NetworkServer struct {
	Title string `xml:"title"`
	UUID  string `xml:"uuid"`
}
