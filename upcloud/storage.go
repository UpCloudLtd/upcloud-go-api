package upcloud

const (
	StorageTypeDisk     = "disk"
	StorageTypeCDROM    = "cdrom"
	StorageTypeTemplate = "template"
	StorageTypeBackup   = "backup"

	StorageTierHDD     = "hdd"
	StorageTierMaxIOPS = "maxiops"

	StorageAccessPublic  = "public"
	StorageAccessPrivate = "private"
)

/**
Represents a /storage response
*/
type Storages struct {
	Storages []Storage `xml:"storage"`
}

/**
Represents a storage device
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
Represents detailed information about a piece of storage
*/
type StorageDetails struct {
	Storage

	BackupRule  *BackupRule `xml:"backup_rule"`
	ServerUUIDs []string    `xml:"servers>server"`
}

/**
Represents a backup rule
*/
type BackupRule struct {
	Interval  string `xml:"interval"`
	Time      string `xml:"time"`
	Retention string `xml:"retention"`
}
