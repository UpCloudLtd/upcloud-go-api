package upcloud

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnmarshalDevicesAvailability(t *testing.T) {
	originalJSON := `{
		"fi-hel2" : {
			"gpu_plans" : {
				"GPU-8xCPU-64GB-1xL40S" : {
					"amount" : 123
				}
			}
		}
	}
	`

	var actual DevicesAvailability
	err := json.Unmarshal([]byte(originalJSON), &actual)
	assert.NoError(t, err)

	expected := DevicesAvailability{
		"fi-hel2": Devices{
			GPUPlans: map[string]DeviceAvailability{
				"GPU-8xCPU-64GB-1xL40S": {
					Amount: 123,
				},
			},
		},
	}

	assert.Equal(t, expected, actual)
}
