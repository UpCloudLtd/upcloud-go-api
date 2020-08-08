package upcloud

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func timeParse(s string) time.Time {
	t, _ := time.Parse(time.RFC3339, s)
	return t
}

// TestUnmarshalHosts tests that multiple Host and Stat are unmarshalled correcltly.
func TestUnmarshalHosts(t *testing.T) {
	originalJSON := `
	  {
		"hosts": {
		  "host": [
			{
			  "id": 7653311107,
			  "description": "My Host #1",
			  "zone": "private-zone-id",
			  "windows_enabled": "no",
			  "stats": {
				"stat": [
				  {
					"name": "cpu_idle",
					"timestamp": "2019-08-09T12:46:57Z",
					"value": 95.2
				  },
				  {
					"name": "memory_free",
					"timestamp": "2019-08-09T12:46:57Z",
					"value": 102
				  }
				]
			  }
			},
			{
			  "id": 8055964291,
			  "description": "My Host #2",
			  "zone": "private-zone-id",
			  "windows_enabled": "no",
			  "stats": {
				"stat": [
				  {
					"name": "cpu_idle",
					"timestamp": "2019-08-09T12:46:57Z",
					"value": 80.1
				  },
				  {
					"name": "memory_free",
					"timestamp": "2019-08-09T12:46:57Z",
					"value": 61
				  }
				]
			  }
			}
		  ]
		}
	  }
	`

	var hosts Hosts
	err := json.Unmarshal([]byte(originalJSON), &hosts)
	assert.NoError(t, err)

	testsHosts := []Host{
		{
			ID:             7653311107,
			Description:    "My Host #1",
			Zone:           "private-zone-id",
			WindowsEnabled: false,
			Stats: []Stat{
				{
					Name:      "cpu_idle",
					Timestamp: timeParse("2019-08-09T12:46:57Z"),
					Value:     95.2,
				},
				{
					Name:      "memory_free",
					Timestamp: timeParse("2019-08-09T12:46:57Z"),
					Value:     102,
				},
			},
		},
		{
			ID:             8055964291,
			Description:    "My Host #2",
			Zone:           "private-zone-id",
			WindowsEnabled: false,
			Stats: []Stat{
				{
					Name:      "cpu_idle",
					Timestamp: timeParse("2019-08-09T12:46:57Z"),
					Value:     80.1,
				},
				{
					Name:      "memory_free",
					Timestamp: timeParse("2019-08-09T12:46:57Z"),
					Value:     61,
				},
			},
		},
	}

	for i, h := range testsHosts {
		assert.Equal(t, h, hosts.Hosts[i])
	}
}

// TestUnmarshalHosts tests that a single Host and Stat are unmarshalled correcltly.
func TestUnmarshalHost(t *testing.T) {
	originalJSON := `
	{
		"host": {
		  "id": 7653311107,
		  "description": "My Host #1",
		  "zone": "private-zone-id",
		  "windows_enabled": "no",
		  "stats": {
			"stat": [
			  {
				"name": "cpu_idle",
				"timestamp": "2019-08-09T12:46:57Z",
				"value": 95.2
			  },
			  {
				"name": "memory_free",
				"timestamp": "2019-08-09T12:46:57Z",
				"value": 102
			  }
			]
		  }
		}
	  }
	`

	var host Host
	err := json.Unmarshal([]byte(originalJSON), &host)
	assert.NoError(t, err)

	testHost := Host{
		ID:             7653311107,
		Description:    "My Host #1",
		Zone:           "private-zone-id",
		WindowsEnabled: false,
		Stats: []Stat{
			{
				Name:      "cpu_idle",
				Timestamp: timeParse("2019-08-09T12:46:57Z"),
				Value:     95.2,
			},
			{
				Name:      "memory_free",
				Timestamp: timeParse("2019-08-09T12:46:57Z"),
				Value:     102,
			},
		},
	}

	assert.Equal(t, testHost, host)
}
