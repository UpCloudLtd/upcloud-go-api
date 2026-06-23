package upcloud

import (
	"encoding/json"
	"time"
)

// Hosts represents a GetHosts response
type Hosts struct {
	Hosts []Host `json:"hosts"`
}

// UnmarshalJSON is a custom unmarshaller that deals with
// deeply embedded values.
func (n *Hosts) UnmarshalJSON(b []byte) error {
	type hostWrapper struct {
		Hosts []Host `json:"host"`
	}

	v := struct {
		Hosts hostWrapper `json:"hosts"`
	}{}
	err := json.Unmarshal(b, &v)
	if err != nil {
		return err
	}

	n.Hosts = append(n.Hosts, v.Hosts.Hosts...)

	return nil
}

// StatSlice is a slice of Stat structs
// This exsits to support a custom unmarshaller
type StatSlice []Stat

// UnmarshalJSON is a custom unmarshaller that deals with
// deeply embedded values.
func (t *StatSlice) UnmarshalJSON(b []byte) error {
	v := struct {
		Networks []Stat `json:"stat"`
	}{}
	err := json.Unmarshal(b, &v)
	if err != nil {
		return err
	}

	(*t) = v.Networks

	return nil
}

// Host represents an individual Host in a response
type Host struct {
	// Deprecated: Use HostID instead.
	ID             int       `json:"id"`
	HostID         int64     `json:"-"`
	Description    string    `json:"description"`
	Zone           string    `json:"zone"`
	WindowsEnabled Boolean   `json:"windows_enabled"`
	Stats          StatSlice `json:"stats"`
}

// UnmarshalJSON is a custom unmarshaller that deals with
// deeply embedded values.
func (s *Host) UnmarshalJSON(b []byte) error {
	type localHost Host
	type hostWrapper struct {
		localHost

		ID int64 `json:"id"`
	}

	v := struct {
		Host *hostWrapper `json:"host"`
	}{}
	err := json.Unmarshal(b, &v)
	if err != nil {
		return err
	}

	host := hostWrapper{}
	if v.Host != nil {
		host = *v.Host
	} else {
		err = json.Unmarshal(b, &host)
		if err != nil {
			return err
		}
	}

	*s = Host(host.localHost)
	s.setHostID(host.ID)

	return nil
}

func (s *Host) setHostID(hostID int64) {
	s.HostID = hostID
	s.ID = 0
	if int64FitsInt(hostID) {
		s.ID = int(hostID)
	}
}

// Stat represents Host stats in a response
type Stat struct {
	Name      string    `json:"name"`
	Timestamp time.Time `json:"timestamp"`
	Value     float64   `json:"value"`
}
