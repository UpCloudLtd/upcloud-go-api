package service

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/UpCloudLtd/upcloud-go-api/v7/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v7/upcloud/client"
	"github.com/stretchr/testify/assert"
)

func TestParseJSONServiceErrorMinimal(t *testing.T) {
	want := &upcloud.Problem{
		Type:   "CODE",
		Title:  "msg",
		Status: http.StatusNotFound,
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

func TestParseJSONServiceErrorWithProblem(t *testing.T) {
	want := &upcloud.Problem{
		Type:          "typexx",
		Title:         "titlexx",
		Status:        http.StatusBadRequest,
		CorrelationID: "corrxx",
		InvalidParams: []upcloud.ProblemInvalidParam{
			{
				Name:   "namex",
				Reason: "reasonx",
			},
		},
	}

	got := parseJSONServiceError(&client.Error{
		ErrorCode: http.StatusBadRequest,
		Type:      client.ErrorTypeProblem,
		ResponseBody: []byte(`
			{
				"type": "typexx",
				"title": "titlexx",
				"status": 400,
				"correlation_id": "corrxx",
				"invalid_params": [
					{
						"name": "namex",
						"reason": "reasonx"
					}
				]
			}
		`),
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
