package service

import (
	"context"
	"testing"

	"github.com/dnaeon/go-vcr/recorder"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
)

const testFiHel1Zone string = "fi-hel1"

// TestGetZones  tests that the GetZones() function returns proper data
func TestGetZones(t *testing.T) {
	record(t, "getzones", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service) {
		zones, err := svc.GetZones(ctx)
		require.NoError(t, err)
		assert.NotEmpty(t, zones.Zones)

		var found bool
		for _, z := range zones.Zones {
			if z.Description == "Helsinki #1" && z.ID == testFiHel1Zone {
				found = true
				assert.True(t, z.Public.Bool())
				break
			}
		}
		assert.True(t, found)
	})
}

// TestGetPriceZones tests that GetPriceZones() function returns proper data
func TestGetPriceZones(t *testing.T) {
	record(t, "getpricezones", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service) {
		zones, err := svc.GetPriceZones(ctx)
		require.NoError(t, err)
		assert.NotEmpty(t, zones.PriceZones)

		var found bool
		var zone upcloud.PriceZone //nolint: staticcheck // To be removed in v9
		for _, z := range zones.PriceZones {
			if z.Name == testFiHel1Zone {
				found = true
				zone = z
				break
			}
		}
		assert.True(t, found)
		assert.NotZero(t, zone.Firewall.Amount)
		assert.NotZero(t, zone.Firewall.Price)
		assert.NotZero(t, zone.IPv4Address.Amount)
		assert.NotZero(t, zone.IPv4Address.Price)
	})
}

// TestGetTimeZones ensures that the GetTimeZones() function returns proper data
func TestGetTimeZones(t *testing.T) {
	record(t, "gettimezones", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service) {
		zones, err := svc.GetTimeZones(ctx)
		require.NoError(t, err)
		assert.NotEmpty(t, zones.TimeZones)
		assert.Contains(t, zones.TimeZones, "Pacific/Wallis")
	})
}

// TestGetPlans ensures that the GetPlans() functions returns proper data
func TestGetPlans(t *testing.T) {
	record(t, "getplans", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service) {
		plans, err := svc.GetPlans(ctx)
		require.NoError(t, err)
		assert.NotEmpty(t, plans.Plans)

		var found bool
		var plan upcloud.Plan
		for _, p := range plans.Plans {
			if p.Name == "1xCPU-1GB" {
				found = true
				plan = p
				break
			}
		}
		assert.True(t, found)

		assert.Equal(t, 1, plan.CoreNumber)
		assert.Equal(t, 1024, plan.MemoryAmount)
		assert.Equal(t, 1024, plan.PublicTrafficOut)
		assert.Equal(t, 25, plan.StorageSize)
		assert.Equal(t, upcloud.StorageTierMaxIOPS, plan.StorageTier)
	})
}
