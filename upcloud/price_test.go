package upcloud

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

const testdata = `
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
        "name": "de-fra1",
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

func TestUnmarshalPriceZones(t *testing.T) {
	priceZones := PriceZones{}
	err := json.Unmarshal([]byte(testdata), &priceZones)
	assert.Nil(t, err)
	assert.Len(t, priceZones.PriceZones, 2)

	zone := priceZones.PriceZones[0]
	assert.Equal(t, 1, zone.Firewall.Amount)
	assert.Equal(t, 0.56, zone.Firewall.Price)
	assert.Equal(t, 1000000, zone.IORequestBackup.Amount)
	assert.Equal(t, 10.0, zone.IORequestBackup.Price)
}

func TestUnmarshalPricingByZone(t *testing.T) {
	pricing := PricingByZone{}
	err := json.Unmarshal([]byte(testdata), &pricing)
	assert.Nil(t, err)
	assert.Len(t, pricing, 2)

	// Test de-fra1 zone
	deFramItems, ok := pricing["de-fra1"]
	assert.True(t, ok)
	assert.NotNil(t, deFramItems)

	// Test specific pricing items in de-fra1
	firewall := deFramItems["firewall"]
	assert.Equal(t, 2, firewall.Amount)
	assert.Equal(t, 0.5, firewall.Price)

	ioRequest := deFramItems["io_request_backup"]
	assert.Equal(t, 2000000, ioRequest.Amount)
	assert.Equal(t, 10.0, ioRequest.Price)

	serverPlan := deFramItems["server_plan_1xCPU-2GB"]
	assert.Equal(t, 1, serverPlan.Amount)
	assert.Equal(t, 2.2321, serverPlan.Price)

	serverPlan2 := deFramItems["server_plan_2xCPU-4GB"]
	assert.Equal(t, 1, serverPlan2.Amount)
	assert.Equal(t, 4.4642, serverPlan2.Price)

	storageMax := deFramItems["storage_maxiops"]
	assert.Equal(t, 1, storageMax.Amount)
	assert.Equal(t, 0.028, storageMax.Price)

	// Test fi-hel1 zone
	helItems, ok := pricing["fi-hel1"]
	assert.True(t, ok)
	assert.NotNil(t, helItems)

	helFirewall := helItems["firewall"]
	assert.Equal(t, 1, helFirewall.Amount)
	assert.Equal(t, 0.56, helFirewall.Price)

	helServerPlan := helItems["server_plan_1xCPU-2GB"]
	assert.Equal(t, 1, helServerPlan.Amount)
	assert.Equal(t, 2.2321, helServerPlan.Price)

	assert.Len(t, deFramItems, 17)
	assert.Len(t, helItems, 17)
}
