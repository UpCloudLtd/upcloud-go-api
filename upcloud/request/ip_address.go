package request

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
)

// GetIPAddressDetailsRequest represents a request to retrieve details about a specific IP address
type GetIPAddressDetailsRequest struct {
	Address string
}

// RequestURL implements the Request interface
func (r *GetIPAddressDetailsRequest) RequestURL() string {
	return fmt.Sprintf("/ip_address/%s", r.Address)
}

// AssignIPAddressRequest represents a request to assign a new IP address to a server
type AssignIPAddressRequest struct {
	XMLName xml.Name `xml:"ip_address" json:"-"`

	Access     string `xml:"access" json:"access"`
	Family     string `xml:"family,omitempty" json:"family,omitempty"`
	ServerUUID string `xml:"server" json:"server"`
}

// RequestURL implements the Request interface
func (r *AssignIPAddressRequest) RequestURL() string {
	return "/ip_address"
}

// MarshalJSON is a custom marshaller that deals with
// deeply embedded values.
func (r AssignIPAddressRequest) MarshalJSON() ([]byte, error) {
	type localAssignIPAddressRequest AssignIPAddressRequest
	v := struct {
		AssignIPAddressRequest localAssignIPAddressRequest `json:"ip_address"`
	}{}
	v.AssignIPAddressRequest = localAssignIPAddressRequest(r)

	return json.Marshal(&v)
}

// ModifyIPAddressRequest represents a request to modify the PTR DNS record of a specific IP address
type ModifyIPAddressRequest struct {
	XMLName   xml.Name `xml:"ip_address" json:"-"`
	IPAddress string   `xml:"-" json:"-"`

	PTRRecord string `xml:"ptr_record" json:"ptr_record"`
}

// RequestURL implements the Request interface
func (r *ModifyIPAddressRequest) RequestURL() string {
	return fmt.Sprintf("/ip_address/%s", r.IPAddress)
}

// MarshalJSON is a custom marshaller that deals with
// deeply embedded values.
func (r ModifyIPAddressRequest) MarshalJSON() ([]byte, error) {
	type localModifyIPAddressRequest ModifyIPAddressRequest
	v := struct {
		ModifyIPAddressRequest localModifyIPAddressRequest `json:"ip_address"`
	}{}
	v.ModifyIPAddressRequest = localModifyIPAddressRequest(r)

	return json.Marshal(&v)
}

// ReleaseIPAddressRequest represents a request to remove a specific IP address from server
type ReleaseIPAddressRequest struct {
	IPAddress string
}

// RequestURL implements the Request interface
func (r *ReleaseIPAddressRequest) RequestURL() string {
	return fmt.Sprintf("/ip_address/%s", r.IPAddress)
}
