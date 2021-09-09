package request_test

import (
	"encoding/json"
	"testing"

	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"

	"github.com/stretchr/testify/assert"
)

// TestMarshalGetHostDetailsRequest tests that GetHostDetailsRequest behaves correctly.
func TestMarshalGetHostDetailsRequest(t *testing.T) {
	request := request.GetHostDetailsRequest{
		ID: 1234,
	}

	assert.Equal(t, "/host/1234", request.RequestURL())
}

// TestMarshalModifyHostRequest tests that ModifyHostRequest behaves correctly.
func TestMarshalModifyHostRequest(t *testing.T) {
	request := request.ModifyHostRequest{
		ID:          1234,
		Description: "My New Host",
	}

	expectedJSON := `
	  {
		"host" : {
		  "description": "My New Host"
		}
	  }
	`

	actualJSON, err := json.Marshal(&request)
	assert.NoError(t, err)
	assert.JSONEq(t, expectedJSON, string(actualJSON))

	assert.Equal(t, "/host/1234", request.RequestURL())
}
