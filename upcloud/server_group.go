package upcloud

import "encoding/json"

// ServerAntiAffinityStatus represents the current status of anti affinity setting for a single server. Can be "met" or "unmet"
type ServerAntiAffinityStatus string

const (
	ServerAntiAffinityStatusMet   ServerAntiAffinityStatus = "met"
	ServerAntiAffinityStatusUnmet ServerAntiAffinityStatus = "unmet"
)

// ServerGroupAntiAffinity represents the anti affinity setting for a server groups. Can be "strict", "yes" or "no"
type ServerGroupAntiAffinity string

const (
	// ServerGroupAntiAffinityStrict doesn't allow servers in the same server group to be on the same host
	ServerGroupAntiAffinityStrict ServerGroupAntiAffinity = "strict"
	// ServerGroupAntiAffinityYes tries to put servers on different hosts, but this is not guaranteed
	ServerGroupAntiAffinityYes ServerGroupAntiAffinity = "yes"
	// ServerGroupAntiAffinityNo doesn't affect server host affinity
	ServerGroupAntiAffinityNo ServerGroupAntiAffinity = "no"
)

// ServerGroupMemberAntiAffinityStatus represents all the data related to an anti affinity status for a single member within the server group
type ServerGroupMemberAntiAffinityStatus struct {
	ServerUUID string                   `json:"uuid"`
	Status     ServerAntiAffinityStatus `json:"status"`
}

// ServerGroup represents server group
type ServerGroup struct {
	Labels             LabelSlice                            `json:"labels,omitempty"`
	Members            ServerUUIDSlice                       `json:"servers,omitempty"`
	Title              string                                `json:"title,omitempty"`
	UUID               string                                `json:"uuid,omitempty"`
	AntiAffinity       ServerGroupAntiAffinity               `json:"anti_affinity,omitempty"`
	AntiAffinityStatus []ServerGroupMemberAntiAffinityStatus `json:"anti_affinity_status,omitempty"`
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
