package request

import "fmt"

/**
Represents a request to retrieve details about a specific IP address
*/
type GetIPAddressDetailsRequest struct {
	Address string
}

func (r *GetIPAddressDetailsRequest) RequestURL() string {
	return fmt.Sprintf("/ip_address/%s", r.Address)
}
