package upcloud

import "encoding/json"

// PriceZones represents a /price response
type PriceZones struct {
	PriceZones []PriceZone `xml:"prices" json:"prices"`
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
type PriceZone struct {
	Name string `xml:"name" json:"name"`

	Firewall               *Price `xml:"firewall" json:"firewall"`
	IORequestBackup        *Price `xml:"io_request_backup" json:"io_request_backup"`
	IORequestMaxIOPS       *Price `xml:"io_request_maxiops" json:"io_request_maxiops"`
	IPv4Address            *Price `xml:"ipv4_address" json:"ipv4_address"`
	IPv6Address            *Price `xml:"ipv6_address" json:"ipv6_address"`
	PublicIPv4BandwidthIn  *Price `xml:"public_ipv4_bandwidth_in" json:"public_ipv4_bandwidth_in"`
	PublicIPv4BandwidthOut *Price `xml:"public_ipv4_bandwidth_out" json:"public_ipv4_bandwidth_out"`
	PublicIPv6BandwidthIn  *Price `xml:"public_ipv6_bandwidth_in" json:"public_ipv6_bandwidth_in"`
	PublicIPv6BandwidthOut *Price `xml:"public_ipv6_bandwidth_out" json:"public_ipv6_bandwidth_out"`
	ServerCore             *Price `xml:"server_core" json:"server_core"`
	ServerMemory           *Price `xml:"server_memory" json:"server_memory"`
	ServerPlan1xCPU1GB     *Price `xml:"server_plan_1xCPU-1GB" json:"server_plan_1xCPU-1GB"`
	ServerPlan2xCPU2GB     *Price `xml:"server_plan_2xCPU-2GB" json:"server_plan_1xCPU-2GB"`
	ServerPlan4xCPU4GB     *Price `xml:"server_plan_4xCPU-4GB" json:"server_plan_4xCPU-4GB"`
	ServerPlan6xCPU8GB     *Price `xml:"server_plan_6xCPU-8GB" json:"server_plan_6xCPU-8GB"`
	StorageBackup          *Price `xml:"storage_backup" json:"storage_backup"`
	StorageMaxIOPS         *Price `xml:"storage_maxiops" json:"storage_maxiops"`
	StorageTemplate        *Price `xml:"storage_template" json:"storage_template"`
}

// Price represents a price
type Price struct {
	Amount int     `xml:"amount" json:"amount"`
	Price  float64 `xml:"price" json:"price"`
}
