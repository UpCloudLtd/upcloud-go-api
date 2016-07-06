package request

import (
	"encoding/xml"
	"fmt"
	"time"
)

const (
	PasswordDeliveryNone  = "none"
	PasswordDeliveryEmail = "email"
	PasswordDeliverySMS   = "sms"

	VideoModelVGA    = "vga"
	VideoModelCirrus = "cirrus"

	CreateStorageDeviceActionCreate = "create"
	CreateStorageDeviceActionClone  = "clone"
	CreateStorageDeviceActionAttach = "attach"

	CreateStorageDeviceTierHDD     = "hdd"
	CreateStorageDeviceTierMaxIOPS = "maxiops"

	ServerStopTypeSoft = "soft"
	ServerStopTypeHard = "hard"

	RestartTimeoutActionDestroy = "destroy"
	RestartTimeoutActionIgnore  = "ignore"
)

/**
Represents a request for retrieving details about a server
*/
type GetServerDetailsRequest struct {
	UUID string
}

func (r *GetServerDetailsRequest) RequestURL() string {
	return fmt.Sprintf("/server/%s", r.UUID)
}

/**
Represents a request for creating a new server
*/
type CreateServerRequest struct {
	XMLName xml.Name `xml:"server"`

	AvoidHost string `xml:"avoid_host,omitempty"`
	BootOrder string `xml:"boot_order,omitempty"`
	// TODO: Investigate correct type and format
	CoreNumber string `xml:"core_number,omitempty"`
	// TODO: Convert to boolean
	Firewall    string                  `xml:"firewall,omitempty"`
	Hostname    string                  `xml:"hostname"`
	IPAddresses []CreateServerIPAddress `xml:"ip_addresses>ip_address"`
	LoginUser   string                  `xml:"login_user,omitempty"`
	// TODO: Investigate correct type and format
	MemoryAmount     string                      `xml:"memory_amount,omitempty"`
	PasswordDelivery string                      `xml:"password_delivery,omitempty"`
	Plan             string                      `xml:"plan,omitempty"`
	StorageDevices   []CreateServerStorageDevice `xml:"storage_devices>storage_device"`
	TimeZone         string                      `xml:"timezone,omitempty"`
	Title            string                      `xml:"title"`
	UserData         string                      `xml:"user_data,omitempty"`
	VideoModel       string                      `xml:"video_model,omitempty"`
	// TODO: Convert to boolean
	VNC         string `xml:"vnc,omitempty"`
	VNCPassword string `xml:"vnc_password,omitempty"`
	Zone        string `xml:"zone"`
}

func (r *CreateServerRequest) RequestURL() string {
	return "/server"
}

/**
Represents a storage device for a CreateServerRequest
*/
type CreateServerStorageDevice struct {
	Action  string `xml:"action"`
	Address string `xml:"address,omitempty"`
	Storage string `xml:"storage"`
	Title   string `xml:"title,omitempty"`
	// Storage size in gigabytes
	Size int    `xml:"size"`
	Tier string `xml:"tier,omitempty"`
}

/**
Represents an IP address for a CreateServerRequest
*/
type CreateServerIPAddress struct {
	Access string `xml:"access"`
	Family string `xml:"family"`
}

/**
Represents a request to wait for a server to enter a specific state
*/
type WaitForServerStateRequest struct {
	UUID         string
	DesiredState string
	Timeout      time.Duration
}

/**
Represents a request to start a server
*/
type StartServerRequest struct {
	UUID string

	// TODO: Start server requests have no timeout in the API
	Timeout time.Duration
}

func (r *StartServerRequest) RequestURL() string {
	return fmt.Sprintf("/server/%s/start", r.UUID)
}

/**
Represents a request to stop a server
*/
type StopServerRequest struct {
	XMLName xml.Name `xml:"stop_server"`

	UUID string `xml:"-"`

	StopType string        `xml:"stop_type,omitempty"`
	Timeout  time.Duration `xml:"timeout,omitempty"`
}

func (r *StopServerRequest) RequestURL() string {
	return fmt.Sprintf("/server/%s/stop", r.UUID)
}

/**
Custom marshaller for StopServerRequest which converts the timeout to seconds
*/
func (r *StopServerRequest) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	type Alias StopServerRequest

	return e.Encode(&struct {
		Timeout int `xml:"timeout,omitempty"`
		*Alias
	}{
		Timeout: int(r.Timeout.Seconds()),
		Alias:   (*Alias)(r),
	})
}

/**
Represents a request to restart a server
*/
type RestartServerRequest struct {
	XMLName xml.Name `xml:"restart_server"`

	UUID string `xml:"-"`

	StopType      string        `xml:"stop_type,omitempty"`
	Timeout       time.Duration `xml:"timeout,omitempty"`
	TimeoutAction string        `xml:"timeout_action,omitempty"`
}

func (r *RestartServerRequest) RequestURL() string {
	return fmt.Sprintf("/server/%s/restart", r.UUID)
}

/**
Custom marshaller for RestartServerRequest which converts the timeout to seconds
*/
func (r *RestartServerRequest) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	type Alias RestartServerRequest

	return e.Encode(&struct {
		Timeout int `xml:"timeout,omitempty"`
		*Alias
	}{
		Timeout: int(r.Timeout.Seconds()),
		Alias:   (*Alias)(r),
	})
}

/**
Represents a request to modify a server
*/
type ModifyServerRequest struct {
	XMLName xml.Name `xml:"server"`

	UUID string `xml:"-"`

	AvoidHost string `xml:"avoid_host,omitempty"`
	BootOrder string `xml:"boot_order,omitempty"`
	// TODO: Investigate correct type and format
	CoreNumber string `xml:"core_number,omitempty"`
	// TODO: Convert to boolean
	Firewall string `xml:"firewall,omitempty"`
	Hostname string `xml:"hostname,omitempty"`
	// TODO: Investigate correct type and format
	MemoryAmount string `xml:"memory_amount,omitempty"`
	Plan         string `xml:"plan,omitempty"`
	TimeZone     string `xml:"timezone,omitempty"`
	Title        string `xml:"title,omitempty"`
	VideoModel   string `xml:"video_model,omitempty"`
	// TODO: Convert to boolean
	VNC         string `xml:"vnc,omitempty"`
	VNCPassword string `xml:"vnc_password,omitempty"`
}

func (r *ModifyServerRequest) RequestURL() string {
	return fmt.Sprintf("/server/%s", r.UUID)
}

/**
Represents a request to delete a server
*/
type DeleteServerRequest struct {
	UUID string
}

func (r *DeleteServerRequest) RequestURL() string {
	return fmt.Sprintf("/server/%s", r.UUID)
}
