package service

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/UpCloudLtd/upcloud-go-api/v5/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v5/upcloud/client"
	"github.com/stretchr/testify/assert"
)

func TestParseJSONServiceError(t *testing.T) {
	want := &upcloud.Error{
		ErrorCode:    "CODE",
		ErrorMessage: "msg",
		Status:       http.StatusNotFound,
	}
	got := parseJSONServiceError(&client.Error{
		ErrorCode: http.StatusNotFound,
		ResponseBody: []byte(`
		{
			"error": {
			  "error_message": "msg",
			  "error_code": "CODE"
			}
		  }
		`),
		Type: client.ErrorTypeError,
	})
	assert.Equal(t, want, got)
}

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
	svc := New(client.New("<username>", "<password>"))
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

	return New(c)
}
