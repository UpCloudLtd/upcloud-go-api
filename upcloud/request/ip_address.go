package request

import "fmt"

/**
GetIPAddressDetailsRequest represents a request to retrieve details about a specific IP address
*/
type GetIPAddressDetailsRequest struct {
	Address string
}

/**
RequestURL() implements the Request interface
*/
func (r *GetIPAddressDetailsRequest) RequestURL() string {
	return fmt.Sprintf("/ip_address/%s", r.Address)
}
