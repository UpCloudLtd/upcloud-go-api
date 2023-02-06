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

func TestParseJSONServiceLegacyError(t *testing.T) {
	parsed := parseJSONServiceError(&client.Error{
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

	ucErr, ok := parsed.(*upcloud.Error)
	assert.True(t, ok)
	assert.Equal(t, http.StatusNotFound, ucErr.Status)
	assert.Equal(t, "CODE", ucErr.Code)
	assert.Equal(t, "msg", ucErr.Message)
	assert.False(t, ucErr.IsProblem())

	problem, ok := ucErr.Problem()
	assert.Nil(t, problem)
	assert.False(t, ok)
}

func TestParseJSONServiceProblemError(t *testing.T) {
	parsed := parseJSONServiceError(&client.Error{
		ErrorCode: http.StatusBadRequest,
		ResponseBody: []byte(`
		{
			"type": "https://developers.upcloud.com/1.3/errors#ERROR_INVALID_REQUEST",
			"title": "Validation error",
			"correlation_id": "01GRKDQBGD7FA84MGR9373F093",
			"invalid_params": [
				{
					"name": "plan",
					"reason": "Plan doesn't exist"
				}
			],
			"status": 400
		}
		`),
		Type: client.ErrorTypeProblem,
	})

	ucErr, ok := parsed.(*upcloud.Error)
	assert.True(t, ok)
	assert.Equal(t, http.StatusBadRequest, ucErr.Status)
	assert.Equal(t, "INVALID_REQUEST", ucErr.Code)
	assert.Equal(t, "Validation error", ucErr.Message)
	assert.True(t, ucErr.IsProblem())

	problem, ok := ucErr.Problem()
	assert.NotNil(t, problem)
	assert.True(t, ok)
	assert.Equal(t, "https://developers.upcloud.com/1.3/errors#ERROR_INVALID_REQUEST", problem.Type)
	assert.Equal(t, "Validation error", problem.Title)
	assert.Equal(t, "01GRKDQBGD7FA84MGR9373F093", problem.CorrelationID)
	assert.Equal(t, http.StatusBadRequest, problem.Status)
	assert.Len(t, problem.InvalidParams, 1)
	assert.Equal(t, "plan", problem.InvalidParams[0].Name)
	assert.Equal(t, "Plan doesn't exist", problem.InvalidParams[0].Reason)
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
