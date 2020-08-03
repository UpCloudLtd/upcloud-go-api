package upcloud

import (
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
	Storages []Storage `xml:"storage"`
}

// Storage represents a storage device
type Storage struct {
	Access  string  `xml:"access"`
	License float64 `xml:"license"`
	// TODO: Convert to boolean
	PartOfPlan string `xml:"part_of_plan"`
	Size       int    `xml:"size"`
	State      string `xml:"state"`
	Tier       string `xml:"tier"`
	Title      string `xml:"title"`
	Type       string `xml:"type"`
	UUID       string `xml:"uuid"`
	Zone       string `xml:"zone"`
	// Only for type "backup":
	Origin  string    `xml:"origin"`
	Created time.Time `xml:"created"`
}

// StorageDetails represents detailed information about a piece of storage
type StorageDetails struct {
	Storage

	BackupRule  *BackupRule `xml:"backup_rule"`
	BackupUUIDs []string    `xml:"backups>backup"`
	ServerUUIDs []string    `xml:"servers>server"`
}

// BackupRule represents a backup rule
type BackupRule struct {
	XMLName  xml.Name `xml:"backup_rule"`
	Interval string   `xml:"interval"`
	// Time should be in the format "hhmm", e.g. "0430"
	Time      string `xml:"time"`
	Retention int    `xml:"retention"`
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
