package upcloud

import "encoding/json"

// PricingByZone represents pricing information organized by zone and item name.
// The structure is map[zoneName]map[itemName]Price.
type PricingByZone map[string]map[string]Price

// UnmarshalJSON converts the API response to the PricingByZone format.
func (p *PricingByZone) UnmarshalJSON(b []byte) error {
	// Parse the raw JSON into a generic structure
	var raw struct {
		Prices struct {
			Zone []map[string]interface{} `json:"zone"`
		} `json:"prices"`
	}

	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}

	// Initialize the map
	result := make(map[string]map[string]Price)

	// Process each zone
	for _, zoneData := range raw.Prices.Zone {
		// Extract zone name
		zoneName, ok := zoneData["name"].(string)
		if !ok {
			continue
		}

		// Remove the name from the map
		delete(zoneData, "name")

		// Create map for this zone
		zoneItems := make(map[string]Price)

		for itemName, itemData := range zoneData {
			// Convert itemData to JSON and unmarshal into Price
			itemJSON, err := json.Marshal(itemData)
			if err != nil {
				continue
			}

			var price Price
			if err := json.Unmarshal(itemJSON, &price); err != nil {
				continue
			}

			zoneItems[itemName] = price
		}

		result[zoneName] = zoneItems
	}

	*p = result
	return nil
}

// PriceZones represents a /price response
//
// Deprecated: Use PricingByZone instead.
type PriceZones struct {
	PriceZones []PriceZone `json:"prices"`
}

// UnmarshalJSON is a custom unmarshaller that deals with
// deeply embedded values.
func (s *PriceZones) UnmarshalJSON(b []byte) error {
	type serverWrapper struct {
		PriceZones []PriceZone `json:"zone"`
	}

	v := struct {
		PriceZones serverWrapper `json:"prices"`
	}{}
	err := json.Unmarshal(b, &v)
	if err != nil {
		return err
	}

	s.PriceZones = v.PriceZones.PriceZones

	return nil
}

// PriceZone represents a price zone. A prize zone consists of multiple items that each have a price.
//
// Deprecated: Use PricingByZone instead.
type PriceZone struct {
	Name string `json:"name"`

	Firewall               *Price `json:"firewall"`
	IORequestBackup        *Price `json:"io_request_backup"`
	IORequestMaxIOPS       *Price `json:"io_request_maxiops"`
	IPv4Address            *Price `json:"ipv4_address"`
	IPv6Address            *Price `json:"ipv6_address"`
	PublicIPv4BandwidthIn  *Price `json:"public_ipv4_bandwidth_in"`
	PublicIPv4BandwidthOut *Price `json:"public_ipv4_bandwidth_out"`
	PublicIPv6BandwidthIn  *Price `json:"public_ipv6_bandwidth_in"`
	PublicIPv6BandwidthOut *Price `json:"public_ipv6_bandwidth_out"`
	ServerCore             *Price `json:"server_core"`
	ServerMemory           *Price `json:"server_memory"`
	ServerPlan1xCPU1GB     *Price `json:"server_plan_1xCPU-1GB"`
	ServerPlan2xCPU2GB     *Price `json:"server_plan_1xCPU-2GB"`
	ServerPlan4xCPU4GB     *Price `json:"server_plan_4xCPU-4GB"`
	ServerPlan6xCPU8GB     *Price `json:"server_plan_6xCPU-8GB"`
	StorageBackup          *Price `json:"storage_backup"`
	StorageMaxIOPS         *Price `json:"storage_maxiops"`
	StorageTemplate        *Price `json:"storage_template"`
}

// Price represents a price
type Price struct {
	Amount int     `json:"amount"`
	Price  float64 `json:"price"`
}
