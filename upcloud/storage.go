package upcloud

/**
Constants
*/
const (
	StorageTypeDisk     = "disk"
	StorageTypeCDROM    = "cdrom"
	StorageTypeTemplate = "template"
	StorageTypeBackup   = "backup"

	StorageTierHDD     = "hdd"
	StorageTierMaxIOPS = "maxiops"

	StorageAccessPublic  = "public"
	StorageAccessPrivate = "private"

	StorageStateOnline      = "online"
	StorageStateMaintenance = "maintenance"
	StorageStateCloning     = "cloning"
	StorageStateBackuping   = "backuping"
	StorageStateError       = "error"
)

/**
Storages represents a /storage response
*/
type Storages struct {
	Storages []Storage `xml:"storage"`
}

/**
Storage represents a storage device
*/
type Storage struct {
	Action  string  `xml:"action"`
	Access  string  `xml:"access"`
	Address string  `xml:"address"`
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
}

/**
StorageDetails represents detailed information about a piece of storage
*/
type StorageDetails struct {
	Storage

	BackupRule  *BackupRule `xml:"backup_rule"`
	ServerUUIDs []string    `xml:"servers>server"`
}

/**
BackupRule represents a backup rule
*/
type BackupRule struct {
	Interval  string `xml:"interval"`
	Time      string `xml:"time"`
	Retention string `xml:"retention"`
}
