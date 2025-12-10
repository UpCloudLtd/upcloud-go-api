package upcloud

import "time"

const (
	FileStorageConfiguredStatusStarted FileStorageConfiguredStatus = "started"
	FileStorageConfiguredStatusStopped FileStorageConfiguredStatus = "stopped"
)

const (
	FileStorageOperationalStateRunning FileStorageOperationalState = "running"
)

const (
	FileStorageACLPermissionReadOnly  FileStorageACLPermission = "ro"
	FileStorageACLPermissionReadWrite FileStorageACLPermission = "rw"
)

type (
	FileStorageConfiguredStatus string
	FileStorageOperationalState string
	FileStorageACLPermission    string
)

type FileStorage struct {
	UUID             string                      `json:"uuid"`
	Name             string                      `json:"name"`
	Zone             string                      `json:"zone"`
	ConfiguredStatus FileStorageConfiguredStatus `json:"configured_status"`
	OperationalState string                      `json:"operational_state,omitempty"`
	SizeGiB          int                         `json:"size_gib"`
	Networks         []FileStorageNetwork        `json:"networks"`
	Shares           []FileStorageShare          `json:"shares"`
	Labels           []Label                     `json:"labels"`
	StateMessages    []FileStorageStateMessage   `json:"state_messages"`
	CreatedAt        time.Time                   `json:"created_at"`
	UpdatedAt        time.Time                   `json:"updated_at"`
}

type FileStorageStateMessage struct {
	OperationalState string `json:"operational_state"`
	Message          string `json:"message"`
	Code             string `json:"code"`
}

type FileStorageNetwork struct {
	UUID      string `json:"uuid"`
	Name      string `json:"name"`
	Family    string `json:"family"`
	IPAddress string `json:"ip_address,omitempty"`
}

type FileStorageShare struct {
	Name     string           `json:"name"`
	Path     string           `json:"path"`
	ACL      []FileStorageACL `json:"acl"`
	Deleting bool             `json:"deleting"`
}

type FileStorageACL struct {
	Target     string                   `json:"target"`
	Permission FileStorageACLPermission `json:"permission"`
}
