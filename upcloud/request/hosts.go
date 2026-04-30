package request

import (
	"encoding/json"
	"fmt"
)

// GetHostDetailsRequest represents the request for the details of a
// single private host
type GetHostDetailsRequest struct {
	// Deprecated: Use HostID instead.
	ID     int
	HostID int64
}

// RequestURL implements the Request interface
func (r *GetHostDetailsRequest) RequestURL() string {
	return fmt.Sprintf("/host/%d", hostIDValue(r.ID, r.HostID))
}

// ModifyHostRequest represents the request to modify a private host
type ModifyHostRequest struct {
	// Deprecated: Use HostID instead.
	ID          int    `json:"-"`
	HostID      int64  `json:"-"`
	Description string `json:"description"`
}

// RequestURL implements the Request interface
func (r *ModifyHostRequest) RequestURL() string {
	return fmt.Sprintf("/host/%d", hostIDValue(r.ID, r.HostID))
}

// MarshalJSON is a custom marshaller that deals with
// deeply embedded values.
func (r ModifyHostRequest) MarshalJSON() ([]byte, error) {
	type localModifyHostRequest ModifyHostRequest
	v := struct {
		ModifyHostRequest localModifyHostRequest `json:"host"`
	}{}
	v.ModifyHostRequest = localModifyHostRequest(r)

	return json.Marshal(&v)
}
