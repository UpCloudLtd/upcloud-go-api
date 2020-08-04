package upcloud

import (
	"encoding/xml"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestUnmarshalTimeZone tests that the TimeZones struct is correctly marshaled
func TestUnmarshalTimeZones(t *testing.T) {
	originalXML := `<?xml version="1.0" encoding="utf-8"?>
<timezones>
    <timezone>Africa/Abidjan</timezone>
    <timezone>Africa/Accra</timezone>
    <timezone>UTC</timezone>
</timezones>`

	timeZones := TimeZones{}
	err := xml.Unmarshal([]byte(originalXML), &timeZones)

	assert.Nil(t, err)
	assert.Len(t, timeZones.TimeZones, 3)
	assert.Equal(t, "Africa/Abidjan", timeZones.TimeZones[0])
}
