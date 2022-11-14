package service

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/UpCloudLtd/upcloud-go-api/v5/upcloud/client"
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

func ExampleNew() {
	client := client.NewWithContext("<username>", "<password>")
	svc := New(client)
	zones, err := svc.GetZones(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%d", len(zones.Zones))
}

func ExampleNewWithContext() {
	client := client.NewWithContext("<username>", "<password>")
	svc := New(client)
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

	c := client.NewWithContext(user, password)
	c.SetTimeout(time.Minute * 5)

	return New(c)
}
