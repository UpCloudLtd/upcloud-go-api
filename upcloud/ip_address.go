package upcloud

import "encoding/json"

// Constants
const (
	IPAddressFamilyIPv4 = "IPv4"
	IPAddressFamilyIPv6 = "IPv6"

	IPAddressAccessPrivate = "private"
	IPAddressAccessPublic  = "public"
)

// IPAddresses represents a /ip_address response
type IPAddresses struct {
	IPAddresses []IPAddress `xml:"ip_address" json:"ip_addresses"`
}

// UnmarshalJSON is a custom unmarshaller that deals with
// deeply embedded values.
func (s *IPAddresses) UnmarshalJSON(b []byte) error {
	type localIPAddress IPAddress
	type ipAddressWrapper struct {
		IPAddresses []localIPAddress `json:"ip_address"`
	}

	v := struct {
		IPAddresses ipAddressWrapper `json:"ip_addresses"`
	}{}
	err := json.Unmarshal(b, &v)
	if err != nil {
		return err
	}

	for _, ip := range v.IPAddresses.IPAddresses {
		s.IPAddresses = append(s.IPAddresses, IPAddress(ip))
	}

	return nil
}

// IPAddress represents an IP address
type IPAddress struct {
	Access  string `xml:"access" json:"access"`
	Address string `xml:"address" json:"address"`
	Family  string `xml:"family" json:"family"`
	// TODO: Convert to boolean
	PartOfPlan string `xml:"part_of_plan" json:"part_of_plan"`
	PTRRecord  string `xml:"ptr_record" json:"ptr_record"`
	ServerUUID string `xml:"server" json:"server"`
}

// UnmarshalJSON is a custom unmarshaller that deals with
// deeply embedded values.
func (s *IPAddress) UnmarshalJSON(b []byte) error {
	type localIPAddress IPAddress

	v := struct {
		IPAddress localIPAddress `json:"ip_address"`
	}{}
	err := json.Unmarshal(b, &v)
	if err != nil {
		return err
	}

	(*s) = IPAddress(v.IPAddress)

	return nil
}
