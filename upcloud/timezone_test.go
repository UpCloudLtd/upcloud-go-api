package upcloud_test

import (
	"encoding/json"
	"testing"

	"github.com/UpCloudLtd/upcloud-go-api/upcloud"

	"github.com/stretchr/testify/assert"
)

// TestUnmarshalTimeZone tests that the TimeZones struct is correctly marshaled.
func TestUnmarshalTimeZones(t *testing.T) {
	originalJSON := `
{
	"timezones": {
	  "timezone": [
		"Africa/Abidjan",
		"Africa/Accra",
		"UTC"
	  ]
	}
  }
`

	timeZones := upcloud.TimeZones{}
	err := json.Unmarshal([]byte(originalJSON), &timeZones)
	assert.NoError(t, err)

	timezoneData := []string{
		"Africa/Abidjan",
		"Africa/Accra",
		"UTC",
	}
	assert.Len(t, timeZones.TimeZones, 3)

	for i, tz := range timezoneData {
		assert.Equal(t, tz, timeZones.TimeZones[i])
	}
}
