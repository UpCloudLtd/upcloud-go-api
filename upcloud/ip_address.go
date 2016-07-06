package upcloud

const (
	IPAddressFamilyIPv4 = "IPv4"
	IPAddressFamilyIPv6 = "IPv6"

	IPAddressAccessPrivate = "private"
	IPAddressAccessPublic  = "public"
)

/**
Represents a /ip_address response
*/
type IPAddresses struct {
	IPAddresses []IPAddress `xml:"ip_address"`
}

/**
Represents an IP address
*/
type IPAddress struct {
	Access  string `xml:"access"`
	Address string `xml:"address"`
	Family  string `xml:"family"`
	// TODO: Convert to boolean
	PartOfPlan string `xml:"part_of_plan"`
	PTRRecord  string `xml:"ptr_record"`
	ServerUUID string `xml:"server"`
}
