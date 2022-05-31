package service

import (
	"context"

	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud"
)

type CloudContext interface {
	GetZones(ctx context.Context) (*upcloud.Zones, error)
	GetPriceZones(ctx context.Context) (*upcloud.PriceZones, error)
	GetTimeZones(ctx context.Context) (*upcloud.TimeZones, error)
	GetPlans(ctx context.Context) (*upcloud.Plans, error)
}

// GetZones returns the available zones
func (s *ServiceContext) GetZones(ctx context.Context) (*upcloud.Zones, error) {
	zones := upcloud.Zones{}
	return &zones, s.get(ctx, "/zone", &zones)
}

// GetPriceZones returns the available price zones and their corresponding prices
func (s *ServiceContext) GetPriceZones(ctx context.Context) (*upcloud.PriceZones, error) {
	zones := upcloud.PriceZones{}
	return &zones, s.get(ctx, "/price", &zones)
}

// GetTimeZones returns the available timezones
func (s *ServiceContext) GetTimeZones(ctx context.Context) (*upcloud.TimeZones, error) {
	zones := upcloud.TimeZones{}
	return &zones, s.get(ctx, "/timezone", &zones)
}

// GetPlans returns the available service plans
func (s *ServiceContext) GetPlans(ctx context.Context) (*upcloud.Plans, error) {
	plans := upcloud.Plans{}
	return &plans, s.get(ctx, "/plan", &plans)
}
