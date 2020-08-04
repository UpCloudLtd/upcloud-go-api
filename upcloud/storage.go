package upcloud

import (
	"encoding/json"
	"encoding/xml"
	"time"
)

// Constants
const (
	StorageTypeBackup   = "backup"
	StorageTypeCDROM    = "cdrom"
	StorageTypeDisk     = "disk"
	StorageTypeNormal   = "normal"
	StorageTypeTemplate = "template"

	StorageTierHDD     = "hdd"
	StorageTierMaxIOPS = "maxiops"

	StorageAccessPublic  = "public"
	StorageAccessPrivate = "private"

	StorageStateOnline      = "online"
	StorageStateMaintenance = "maintenance"
	StorageStateCloning     = "cloning"
	StorageStateBackuping   = "backuping"
	StorageStateError       = "error"

	BackupRuleIntervalDaily     = "daily"
	BackupRuleIntervalMonday    = "mon"
	BackupRuleIntervalTuesday   = "tue"
	BackupRuleIntervalWednesday = "wed"
	BackupRuleIntervalThursday  = "thu"
	BackupRuleIntervalFriday    = "fri"
	BackupRuleIntervalSaturday  = "sat"
	BackupRuleIntervalSunday    = "sun"

	CreateServerStorageDeviceActionCreate = "create"
	CreateServerStorageDeviceActionClone  = "clone"
	CreateServerStorageDeviceActionAttach = "attach"
)

// Storages represents a /storage response
type Storages struct {
	Storages []Storage `xml:"storage" json:"storages"`
}

// UnmarshalJSON is a custom unmarshaller that deals with
// deeply embedded values.
func (s *Storages) UnmarshalJSON(b []byte) error {
	type storageWrapper struct {
		Storages []Storage `json:"storage"`
	}

	v := struct {
		Storages storageWrapper `json:"storages"`
	}{}
	err := json.Unmarshal(b, &v)
	if err != nil {
		return err
	}

	s.Storages = v.Storages.Storages

	return nil
}

// Storage represents a storage device
type Storage struct {
	Access  string  `xml:"access" json:"access"`
	License float64 `xml:"license" json:"license"`
	// TODO: Convert to boolean
	PartOfPlan string `xml:"part_of_plan" json:"part_of_plan"`
	Size       int    `xml:"size" json:"size"`
	State      string `xml:"state" json:"state"`
	Tier       string `xml:"tier" json:"tier"`
	Title      string `xml:"title" json:"title"`
	Type       string `xml:"type" json:"type"`
	UUID       string `xml:"uuid" json:"uuid"`
	Zone       string `xml:"zone" json:"zone"`
	// Only for type "backup":
	Origin  string    `xml:"origin" json:"origin"`
	Created time.Time `xml:"created" json:"created"`
}

// BackupUUIDSlice is a slice of string.
// It exists to allow for a custom JSON unmarshaller.
type BackupUUIDSlice []string

// UnmarshalJSON is a custom unmarshaller that deals with
// deeply embedded values.
func (s *BackupUUIDSlice) UnmarshalJSON(b []byte) error {
	v := struct {
		BackupUUIDs []string `json:"backup"`
	}{}
	err := json.Unmarshal(b, &v)
	if err != nil {
		return err
	}

	(*s) = v.BackupUUIDs

	return nil
}

// ServerUUIDSlice is a slice of string.
// It exists to allow for a custom JSON unmarshaller.
type ServerUUIDSlice []string

// UnmarshalJSON is a custom unmarshaller that deals with
// deeply embedded values.
func (s *ServerUUIDSlice) UnmarshalJSON(b []byte) error {
	v := struct {
		ServerUUIDs []string `json:"server"`
	}{}
	err := json.Unmarshal(b, &v)
	if err != nil {
		return err
	}

	(*s) = v.ServerUUIDs

	return nil
}

// StorageDetails represents detailed information about a piece of storage
type StorageDetails struct {
	Storage

	BackupRule  *BackupRule     `xml:"backup_rule" json:"backup_rule"`
	BackupUUIDs BackupUUIDSlice `xml:"backups>backup" json:"backups"`
	ServerUUIDs ServerUUIDSlice `xml:"servers>server" json:"servers"`
}

// UnmarshalJSON is a custom unmarshaller that deals with
// deeply embedded values.
func (s *StorageDetails) UnmarshalJSON(b []byte) error {
	type localStorageDetails StorageDetails

	v := struct {
		StorageDetails localStorageDetails `json:"storage"`
	}{}
	err := json.Unmarshal(b, &v)
	if err != nil {
		return err
	}

	(*s) = StorageDetails(v.StorageDetails)

	return nil
}

// BackupRule represents a backup rule
type BackupRule struct {
	XMLName  xml.Name `xml:"backup_rule" json:"-"`
	Interval string   `xml:"interval" json:"interval"`
	// Time should be in the format "hhmm", e.g. "0430"
	Time      string `xml:"time" json:"time"`
	Retention int    `xml:"retention" json:"retention,string"`
}

// ServerStorageDevice represents a storage device in the context of server requests or server details
type ServerStorageDevice struct {
	XMLName xml.Name `xml:"storage_device"`

	Address string `xml:"address" json:"address"`
	// TODO: Convert to boolean
	PartOfPlan string `xml:"part_of_plan" json:"part_of_plan"`
	UUID       string `xml:"storage" json:"storage"`
	Size       int    `xml:"storage_size" json:"storage_size"`
	Title      string `xml:"storage_title" json:"storage_title"`
	Type       string `xml:"type" json:"type"`
	BootDisk   int    `xml:"-" json:"boot_disk,string"`
}

// CreateServerStorageDevice represents a storage device for a CreateServerRequest
type CreateServerStorageDevice struct {
	XMLName xml.Name `xml:"storage_device" json:"-"`

	Action  string `xml:"action" json:"action"`
	Address string `xml:"address,omitempty" json:"address,omitempty"`
	Storage string `xml:"storage" json:"storage"`
	Title   string `xml:"title,omitempty" json:"title,omitempty"`
	// Storage size in gigabytes
	Size int    `xml:"size" json:"size"`
	Tier string `xml:"tier,omitempty" json:"tier,omitempty"`
	Type string `xml:"type,omitempty" json:"type,omitempty"`
}
