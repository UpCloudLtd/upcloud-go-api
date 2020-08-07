package upcloud

import "encoding/json"

// InterfaceSlice is a slice of Interfaces.
// It exists to allow for a custom JSON unmarshaller.
type InterfaceSlice []Interface

// UnmarshalJSON is a custom unmarshaller that deals with
// deeply embedded values.
func (s *InterfaceSlice) UnmarshalJSON(b []byte) error {
	v := struct {
		Interfaces []Interface `json:"interface"`
	}{}
	err := json.Unmarshal(b, &v)
	if err != nil {
		return err
	}

	(*s) = v.Interfaces

	return nil
}

type Networking struct {
	Interfaces InterfaceSlice `json:"interfaces"`
}

type Interface struct {
	Index       int            `json:"index"`
	IPAddresses IPAddressSlice `json:"ip_addresses"`
	MAC         string         `json:"mac"`
	Network     string         `json:"network"`
	Type        string         `json:"type"`
	Bootable    Boolean        `json:"bootable"`
}
