package upcloud_test

import (
	"encoding/json"
	"testing"

	"github.com/UpCloudLtd/upcloud-go-api/upcloud"

	"github.com/stretchr/testify/assert"
)

// TestUnmarshalZones tests that the Zone and Zones structs are correctly marshaled.
func TestUnmarshalZones(t *testing.T) {
	originalJSON := `
{
    "zones": {
      "zone": [
        {
          "description" : "Frankfurt #1",
          "id" : "de-fra1"
        },
        {
          "description": "Helsinki #2",
          "id": "fi-hel2"
        },
        {
          "description": "London #1",
          "id": "uk-lon1"
        },
        {
          "description" : "Chicago #1",
          "id" : "us-chi1"
        }
      ]
    }
  }
`

	zones := upcloud.Zones{}
	err := json.Unmarshal([]byte(originalJSON), &zones)

	assert.Nil(t, err)
	assert.Len(t, zones.Zones, 4)

	zoneData := []struct {
		Description string
		ID          string
	}{
		{
			ID:          "de-fra1",
			Description: "Frankfurt #1",
		},
		{
			ID:          "fi-hel2",
			Description: "Helsinki #2",
		},
		{
			ID:          "uk-lon1",
			Description: "London #1",
		},
		{
			ID:          "us-chi1",
			Description: "Chicago #1",
		},
	}

	for i, d := range zoneData {
		z := zones.Zones[i]
		assert.Equal(t, d.Description, z.Description)
		assert.Equal(t, d.ID, z.ID)
	}
}
