package service

import (
	"context"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
)

type Cloud interface {
	GetZones(ctx context.Context) (*upcloud.Zones, error)
	GetPriceZones(ctx context.Context) (*upcloud.PriceZones, error) //nolint: staticcheck // To be removed in v9
	GetPricingByZone(ctx context.Context) (*upcloud.PricingByZone, error)
	GetTimeZones(ctx context.Context) (*upcloud.TimeZones, error)
	GetPlans(ctx context.Context) (*upcloud.Plans, error)
	GetDevicesAvailability(ctx context.Context) (*upcloud.DevicesAvailability, error)
}

// GetZones returns the available zones
func (s *Service) GetZones(ctx context.Context) (*upcloud.Zones, error) {
	zones := upcloud.Zones{}
	return &zones, s.get(ctx, "/zone", &zones)
}

// GetPriceZones returns the available price zones and their corresponding prices
//
// Deprecated: Use GetPricingByZone instead.
func (s *Service) GetPriceZones(ctx context.Context) (*upcloud.PriceZones, error) {
	zones := upcloud.PriceZones{}
	return &zones, s.get(ctx, "/price", &zones)
}

// GetPricingByZone returns pricing information organized by zone and item name.
func (s *Service) GetPricingByZone(ctx context.Context) (*upcloud.PricingByZone, error) {
	pricing := upcloud.PricingByZone{}
	return &pricing, s.get(ctx, "/price", &pricing)
}

// GetTimeZones returns the available timezones
func (s *Service) GetTimeZones(ctx context.Context) (*upcloud.TimeZones, error) {
	zones := upcloud.TimeZones{}
	return &zones, s.get(ctx, "/timezone", &zones)
}

// GetPlans returns the available service plans
func (s *Service) GetPlans(ctx context.Context) (*upcloud.Plans, error) {
	plans := upcloud.Plans{}
	return &plans, s.get(ctx, "/plan", &plans)
}

func (s *Service) GetDevicesAvailability(ctx context.Context) (*upcloud.DevicesAvailability, error) {
	r := upcloud.DevicesAvailability{}
	return &r, s.get(ctx, "/device/availability", &r)
}
