package request

import (
	"encoding/json"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestMarshalGetHostDetailsRequest tests that GetHostDetailsRequest behaves correctly
func TestMarshalGetHostDetailsRequest(t *testing.T) {
	request := GetHostDetailsRequest{
		ID: 1234,
	}

	assert.Equal(t, "/host/1234", request.RequestURL())
}

func TestMarshalGetHostDetailsRequest_WithHostID(t *testing.T) {
	request := GetHostDetailsRequest{
		HostID: math.MaxInt64,
	}

	assert.Equal(t, "/host/9223372036854775807", request.RequestURL())
}

// TestMarshalModifyHostRequest tests that ModifyHostRequest behaves correctly
func TestMarshalModifyHostRequest(t *testing.T) {
	request := ModifyHostRequest{
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

func TestMarshalModifyHostRequest_WithHostID(t *testing.T) {
	request := ModifyHostRequest{
		HostID:      math.MaxInt64,
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

	assert.Equal(t, "/host/9223372036854775807", request.RequestURL())
}
