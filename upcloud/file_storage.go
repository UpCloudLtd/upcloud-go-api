package upcloud

import "time"

// FileStorage represents a File Storage service instance
// Fields based on OpenAPI spec

type FileStorage struct {
	UUID             string                    `json:"uuid"`
	Name             string                    `json:"name"`
	Zone             string                    `json:"zone"`
	SizeGiB          int                       `json:"size_gib"`
	ConfiguredStatus string                    `json:"configured_status"`
	OperationalState string                    `json:"operational_state,omitempty"`
	CreatedAt        time.Time                 `json:"created_at"`
	UpdatedAt        time.Time                 `json:"updated_at"`
	Networks         []FileStorageNetwork      `json:"networks"`
	Shares           []FileStorageShare        `json:"shares"`
	Labels           []Label                   `json:"labels"`
	StateMessages    []FileStorageStateMessage `json:"state_messages"`
}

type FileStorageStateMessage struct {
	OperationalState string `json:"operational_state"`
	Message          string `json:"message"`
	Code             string `json:"code"`
}

type FileStorageNetwork struct {
	UUID   string `json:"uuid"`
	Name   string `json:"name"`
	Family string `json:"family"`
	IP     string `json:"ip_address,omitempty"`
}

type FileStorageShare struct {
	Name string           `json:"name"`
	Path string           `json:"path"`
	ACL  []FileStorageACL `json:"acl"`
}

type FileStorageACL struct {
	Target     string `json:"target"`
	Permission string `json:"permission"`
}
