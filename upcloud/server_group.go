package upcloud

import "encoding/json"

// ServerGroups represents list of server groups
type ServerGroups []ServerGroup

// UnmarshalJSON is a custom unmarshaller that deals with deeply embedded values.
func (s *ServerGroups) UnmarshalJSON(b []byte) error {
	type sg ServerGroup
	v := struct {
		ServerGroups struct {
			ServerGroup []sg `json:"server_group"`
		} `json:"server_groups"`
	}{}

	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}

	for _, val := range v.ServerGroups.ServerGroup {
		*s = append(*s, ServerGroup(val))
	}

	return nil
}

// ServerGroup represents server group
type ServerGroup struct {
	UUID    string          `json:"uuid,omitempty"`
	Title   string          `json:"title,omitempty"`
	Members ServerUUIDSlice `json:"servers,omitempty"`
}

// UnmarshalJSON is a custom unmarshaller that deals with deeply embedded values.
func (s *ServerGroup) UnmarshalJSON(b []byte) error {
	type sg ServerGroup
	v := struct {
		ServerGroup sg `json:"server_group"`
	}{}

	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}

	*s = ServerGroup(v.ServerGroup)
	return nil
}
