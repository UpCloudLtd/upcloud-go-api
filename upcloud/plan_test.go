package upcloud

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestUnmarshalPlans tests that Plans and Plan objects unmarshal correctly.
func TestUnmarshalPlans(t *testing.T) {
	originalJSON := `
      {
        "plans" : {
          "plan" : [
            {
              "core_number" : 1,
              "memory_amount" : 2048,
              "name" : "1xCPU-2GB",
              "public_traffic_out" : 2048,
              "storage_size" : 50,
              "storage_tier" : "maxiops"
            },
            {
              "core_number" : 2,
              "memory_amount" : 4096,
              "name" : "2xCPU-4GB",
              "public_traffic_out" : 4096,
              "storage_size" : 80,
              "storage_tier" : "maxiops"
            }
          ]
        }
      }
    `

	plans := Plans{}
	err := json.Unmarshal([]byte(originalJSON), &plans)
	assert.Nil(t, err)
	assert.Len(t, plans.Plans, 2)

	testData := []Plan{
		{
			CoreNumber:       1,
			MemoryAmount:     2048,
			Name:             "1xCPU-2GB",
			PublicTrafficOut: 2048,
			StorageSize:      50,
			StorageTier:      "maxiops",
		},
		{
			CoreNumber:       2,
			MemoryAmount:     4096,
			Name:             "2xCPU-4GB",
			PublicTrafficOut: 4096,
			StorageSize:      80,
			StorageTier:      "maxiops",
		},
	}

	for i, p := range testData {
		plan := plans.Plans[i]
		assert.Equal(t, p.CoreNumber, plan.CoreNumber)
		assert.Equal(t, p.MemoryAmount, plan.MemoryAmount)
		assert.Equal(t, p.Name, plan.Name)
		assert.Equal(t, p.PublicTrafficOut, plan.PublicTrafficOut)
		assert.Equal(t, p.StorageSize, plan.StorageSize)
		assert.Equal(t, p.StorageTier, plan.StorageTier)
	}
}
