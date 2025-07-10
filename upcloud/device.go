package upcloud

// DevicesAvailability represents the availability of devices in different zones. Zone ID is used as the key.
type DevicesAvailability map[string]Devices

type Devices struct {
	// GPUPlans contains the available GPU plans. Plan name is used as the key.
	GPUPlans map[string]DeviceAvailability `json:"gpu_plans"`
}

type DeviceAvailability struct {
	Amount int `json:"amount"`
}
