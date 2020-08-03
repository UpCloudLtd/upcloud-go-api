package upcloud

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestUmarshalPriceZones tests that PrizeZones, PriceZone and Price are unmarshaled correctly
func TestUmarshalPriceZones(t *testing.T) {
	originalJSON := `
{
    "prices": {
      "zone": [
        {
          "name": "fi-hel1",
          "firewall": {
            "amount": 1,
            "price": 0.56
          },
          "io_request_backup": {
            "amount": 1000000,
            "price": 10
          },
          "io_request_hdd": {
            "amount": 1000000,
            "price": 0
          },
          "io_request_maxiops": {
            "amount": 1000000,
            "price": 0
          },
          "ipv4_address": {
            "amount": 1,
            "price": 0.3
          },
          "ipv6_address": {
            "amount": 1,
            "price": 0
          },
          "public_ipv4_bandwidth_in": {
            "amount": 1,
            "price": 0
          },
          "public_ipv4_bandwidth_out": {
            "amount": 1,
            "price": 5
          },
          "public_ipv6_bandwidth_in": {
            "amount": 1,
            "price": 0
          },
          "public_ipv6_bandwidth_out": {
            "amount": 1,
            "price": 5
          },
          "server_core": {
            "amount": 1,
            "price": 1.3
          },
          "server_memory": {
            "amount": 256,
            "price": 0.45
          },
          "storage_backup": {
            "amount": 1,
            "price": 0.007
          },
          "storage_hdd": {
            "amount": 1,
            "price": 0.013
          },
          "storage_maxiops": {
            "amount": 1,
            "price": 0.028
          },
          "server_plan_1xCPU-2GB": {
            "amount": 1,
            "price": 2.2321
          },
          "server_plan_2xCPU-4GB": {
            "amount": 1,
            "price": 4.4642
          }
        },
        {
          "name": "fi-hel2",
          "firewall": {
            "amount": 2,
            "price": 0.5
          },
          "io_request_backup": {
            "amount": 2000000,
            "price": 10
          },
          "io_request_hdd": {
            "amount": 2000000,
            "price": 1
          },
          "io_request_maxiops": {
            "amount": 2000000,
            "price": 0
          },
          "ipv4_address": {
            "amount": 1,
            "price": 0.3
          },
          "ipv6_address": {
            "amount": 1,
            "price": 0
          },
          "public_ipv4_bandwidth_in": {
            "amount": 1,
            "price": 0
          },
          "public_ipv4_bandwidth_out": {
            "amount": 1,
            "price": 5
          },
          "public_ipv6_bandwidth_in": {
            "amount": 1,
            "price": 0
          },
          "public_ipv6_bandwidth_out": {
            "amount": 1,
            "price": 5
          },
          "server_core": {
            "amount": 1,
            "price": 1.3
          },
          "server_memory": {
            "amount": 256,
            "price": 0.45
          },
          "storage_backup": {
            "amount": 1,
            "price": 0.007
          },
          "storage_hdd": {
            "amount": 1,
            "price": 0.013
          },
          "storage_maxiops": {
            "amount": 1,
            "price": 0.028
          },
          "server_plan_1xCPU-2GB": {
            "amount": 1,
            "price": 2.2321
          },
          "server_plan_2xCPU-4GB": {
            "amount": 1,
            "price": 4.4642
          }
        }
      ]
    }
  }
  `

	priceZones := PriceZones{}
	err := json.Unmarshal([]byte(originalJSON), &priceZones)
	assert.Nil(t, err)
	assert.Len(t, priceZones.PriceZones, 2)

	zone := priceZones.PriceZones[0]
	assert.Equal(t, 1, zone.Firewall.Amount)
	assert.Equal(t, 0.56, zone.Firewall.Price)
	assert.Equal(t, 1000000, zone.IORequestBackup.Amount)
	assert.Equal(t, 10.0, zone.IORequestBackup.Price)

	// TODO: Test the remaining fields
}
