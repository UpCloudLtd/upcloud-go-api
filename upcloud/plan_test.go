package upcloud

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestUnmarshalPlans tests that Plans and Plan objects unmarshal correctly
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
            },
            {
              "core_number" : 8,
              "gpu_amount" : 1,
              "gpu_model" : "NVIDIA L40S",
              "memory_amount" : 65536,
              "name" : "GPU-8xCPU-64GB-1xL40S",
              "public_traffic_out" : 12,
              "storage_size" : 0,
              "storage_tier" : null
            }
          ]
        }
      }
    `

	plans := Plans{}
	err := json.Unmarshal([]byte(originalJSON), &plans)
	assert.Nil(t, err)
	assert.Len(t, plans.Plans, 3)

	testData := []Plan{
		{
			CoreNumber:       1,
			MemoryAmount:     2048,
			Name:             "1xCPU-2GB",
			PublicTrafficOut: 2048,
			StorageSize:      50,
			StorageTier:      "maxiops",
			GPUAmount:        0,
			GPUModel:         "",
		},
		{
			CoreNumber:       2,
			MemoryAmount:     4096,
			Name:             "2xCPU-4GB",
			PublicTrafficOut: 4096,
			StorageSize:      80,
			StorageTier:      "maxiops",
			GPUAmount:        0,
			GPUModel:         "",
		},
		{
			CoreNumber:       8,
			MemoryAmount:     65536,
			Name:             "GPU-8xCPU-64GB-1xL40S",
			PublicTrafficOut: 12,
			StorageSize:      0,
			StorageTier:      "",
			GPUAmount:        1,
			GPUModel:         "NVIDIA L40S",
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
		assert.Equal(t, p.GPUAmount, plan.GPUAmount)
		assert.Equal(t, p.GPUModel, plan.GPUModel)
	}
}
