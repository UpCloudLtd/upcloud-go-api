package upcloud

import "encoding/json"

// Plans represents a /plan response
type Plans struct {
	Plans []Plan `xml:"plan" json:"plans"`
}

// UnmarshalJSON is a custom unmarshaller that deals with
// deeply embedded values.
func (s *Plans) UnmarshalJSON(b []byte) error {
	type planWrapper struct {
		Plans []Plan `json:"plan"`
	}

	v := struct {
		Plans planWrapper `json:"plans"`
	}{}
	err := json.Unmarshal(b, &v)
	if err != nil {
		return err
	}

	s.Plans = v.Plans.Plans

	return nil
}

// Plan represents a pre-configured server configuration plan
type Plan struct {
	CoreNumber       int    `xml:"core_number" json:"core_number"`
	MemoryAmount     int    `xml:"memory_amount" json:"memory_amount"`
	Name             string `xml:"name" json:"name"`
	PublicTrafficOut int    `xml:"public_traffic_out" json:"public_traffic_out"`
	StorageSize      int    `xml:"storage_size" json:"storage_size"`
	StorageTier      string `xml:"storage_tier" json:"storage_tier"`
}
