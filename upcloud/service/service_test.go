package service

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/dnaeon/go-vcr/recorder"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud/client"
)

// TestMain is the main test method
func TestMain(m *testing.M) {
	retCode := m.Run()

	// Optionally perform teardown
	deleteResources := os.Getenv("UPCLOUD_GO_SDK_TEST_DELETE_RESOURCES") == "yes"
	noCredentials := os.Getenv("UPCLOUD_GO_SDK_TEST_NO_CREDENTIALS") == "yes"
	if deleteResources && !noCredentials {
		log.Print("UPCLOUD_GO_SDK_TEST_DELETE_RESOURCES defined, deleting all resources ...")
		teardown()
	}

	os.Exit(retCode)
}

// TestGetAccount tests that the GetAccount() method returns proper data
func TestGetAccount(t *testing.T) {
	if os.Getenv("UPCLOUD_GO_SDK_TEST_NO_CREDENTIALS") == "yes" || testing.Short() {
		t.Skip("Skipping TestGetAccount...")
	}

	svc := getService()

	account, err := svc.GetAccount()
	require.NoError(t, err)

	username, _ := getCredentials()

	if account.UserName != username {
		t.Errorf("TestGetAccount expected %s, got %s", username, account.UserName)
	}

	assert.NotZero(t, account.ResourceLimits.Cores)
	assert.NotZero(t, account.ResourceLimits.Memory)
	assert.NotZero(t, account.ResourceLimits.Networks)
	assert.NotZero(t, account.ResourceLimits.PublicIPv6)
	assert.NotZero(t, account.ResourceLimits.StorageHDD)
	assert.NotZero(t, account.ResourceLimits.StorageSSD)
}

// TestGetZones tests that the GetZones() function returns proper data
func TestGetZones(t *testing.T) {
	record(t, "getzones", func(t *testing.T, rec *recorder.Recorder, svc *Service) {
		zones, err := svc.GetZones()
		require.NoError(t, err)
		assert.NotEmpty(t, zones.Zones)

		var found bool
		for _, z := range zones.Zones {
			if z.Description == "Helsinki #1" && z.ID == "fi-hel1" {
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
	record(t, "getpricezones", func(t *testing.T, rec *recorder.Recorder, svc *Service) {
		zones, err := svc.GetPriceZones()
		require.NoError(t, err)
		assert.NotEmpty(t, zones.PriceZones)

		var found bool
		var zone upcloud.PriceZone
		for _, z := range zones.PriceZones {
			if z.Name == "fi-hel1" {
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
	record(t, "gettimezones", func(t *testing.T, rec *recorder.Recorder, svc *Service) {
		zones, err := svc.GetTimeZones()
		require.NoError(t, err)
		assert.NotEmpty(t, zones.TimeZones)

		var found bool
		for _, z := range zones.TimeZones {
			if z == "Pacific/Wallis" {
				found = true
				break
			}
		}
		assert.True(t, found)
	})
}

// TestGetPlans ensures that the GetPlans() functions returns proper data
func TestGetPlans(t *testing.T) {
	record(t, "getplans", func(t *testing.T, rec *recorder.Recorder, svc *Service) {
		plans, err := svc.GetPlans()
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

func ExampleNew() {
	client := client.New("<username>", "<password>")
	svc := New(client)
	zones, err := svc.GetZones()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%d", len(zones.Zones))
}

func ExampleNewWithContext() {
	client := client.NewWithContext("<username>", "<password>")
	svc := NewWithContext(client)
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(10*time.Second))
	defer cancel()
	zones, err := svc.GetZones(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%d", len(zones.Zones))
}

// Configures the test environment.
func getService() *Service {
	user, password := getCredentials()

	c := client.New(user, password)
	c.SetTimeout(time.Minute * 5)

	return New(c)
}
