package upcloud

/**
Represents a /plan response
*/
type Plans struct {
	Plans []Plan `xml:"plan"`
}

/**
Represents a pre-configured server configuration plan
*/
type Plan struct {
	CoreNumber       int    `xml:"core_number"`
	MemoryAmount     int    `xml:"memory_amount"`
	Name             string `xml:"name"`
	PublicTrafficOut int    `xml:"public_traffic_out"`
	StorageSize      int    `xml:"storage_size"`
	StorageTier      string `xml:"storage_tier"`
}
