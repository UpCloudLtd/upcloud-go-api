package request

import (
	"encoding/json"
	"testing"

	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/stretchr/testify/assert"
)

// TestGetIPAddressDetailsRequest tests that GetIPAddressDetailsRequest behaves correctly
func TestGetIPAddressDetailsRequest(t *testing.T) {
	request := GetIPAddressDetailsRequest{
		Address: "0.0.0.0",
	}

	assert.Equal(t, "/ip_address/0.0.0.0", request.RequestURL())
}

// TestMarshalAssignIPAddressRequest tests that AssignIPAddressRequest structs are marshaled correctly
func TestMarshalAssignIPAddressRequest(t *testing.T) {
	request := AssignIPAddressRequest{
		Access:     upcloud.IPAddressAccessPublic,
		Family:     upcloud.IPAddressFamilyIPv4,
		ServerUUID: "009d64ef-31d1-4684-a26b-c86c955cbf46",
	}

	actualJSON, err := json.Marshal(&request)
	assert.NoError(t, err)
	expectedJSON := `
	  {
		"ip_address": {
          "access": "public",
		  "family": "IPv4",
		  "server": "009d64ef-31d1-4684-a26b-c86c955cbf46"
		}
	  }
	`

	assert.JSONEq(t, expectedJSON, string(actualJSON))
}

// TestMarshalAssignIPAddressRequest_OmitFields tests that AssignIPAddressRequest structs are marshaled correctly
// when optional fields are left out
func TestMarshalAssignIPAddressRequest_OmitFields(t *testing.T) {
	request := AssignIPAddressRequest{
		Access:     upcloud.IPAddressAccessPublic,
		ServerUUID: "009d64ef-31d1-4684-a26b-c86c955cbf46",
	}

	actualJSON, err := json.Marshal(&request)
	assert.NoError(t, err)

	expectedJSON := `
	  {
		"ip_address": {
          "access": "public",
		  "server": "009d64ef-31d1-4684-a26b-c86c955cbf46"
		}
	  }
	`

	assert.JSONEq(t, expectedJSON, string(actualJSON))
}

// TestModifyIPAddressRequest tests that ModifyIPAddressRequest structs are marshaled correctly and that their URLs
// are correct
func TestModifyIPAddressRequest(t *testing.T) {
	request := ModifyIPAddressRequest{
		IPAddress: "0.0.0.0",
		PTRRecord: "hostname.example.com",
	}

	actualJSON, err := json.Marshal(&request)
	assert.NoError(t, err)

	expectedJSON := `
	  {
		"ip_address": {
		  "ptr_record": "hostname.example.com"
		}
	  }
	`
	assert.JSONEq(t, expectedJSON, string(actualJSON))
	assert.Equal(t, "/ip_address/0.0.0.0", request.RequestURL())
}

// TestReleaseIPAddressRequest tests that ReleaseIPAddressRequest's URL is correct
func TestReleaseIPAddressRequest(t *testing.T) {
	request := ReleaseIPAddressRequest{
		IPAddress: "0.0.0.0",
	}

	assert.Equal(t, "/ip_address/0.0.0.0", request.RequestURL())
}
