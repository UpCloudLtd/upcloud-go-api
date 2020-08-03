package upcloud

import "encoding/json"

// Zones represents a /zone response
type Zones struct {
	Zones []Zone `xml:"zone" json:"zone"`
}

// UnmarshalJSON is a custom unmarshaller that deals with
// deeply embedded values.
func (s *Zones) UnmarshalJSON(b []byte) error {
	type zoneWrapper struct {
		Zones []Zone `json:"zone"`
	}

	v := struct {
		Zones zoneWrapper `json:"zones"`
	}{}
	err := json.Unmarshal(b, &v)
	if err != nil {
		return err
	}

	s.Zones = v.Zones.Zones

	return nil
}

// Zone represents a zone
type Zone struct {
	ID          string `xml:"id" json:"id"`
	Description string `xml:"description" json:"description"`
}
