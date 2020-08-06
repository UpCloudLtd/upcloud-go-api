package request

import (
	"encoding/json"
	"testing"

	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/stretchr/testify/assert"
)

// TestCreateTagRequest tests that CreateTagRequest behaves correctly
func TestCreateTagRequest(t *testing.T) {
	request := CreateTagRequest{
		Tag: upcloud.Tag{
			Name:        "DEV",
			Description: "Development servers",
			Servers: []string{
				"server1",
				"server2",
			},
		},
	}

	// Check the request URL
	assert.Equal(t, "/tag", request.RequestURL())

	// Check marshaling
	actualJSON, err := json.Marshal(&request)
	assert.NoError(t, err)

	expectedJSON := `
	  {
		"tag": {
		  "name": "DEV",
		  "description": "Development servers",
		  "servers": {
			"server": [
				"server1",
				"server2"
			]
		  }
		}
	  }
	`
	assert.JSONEq(t, expectedJSON, string(actualJSON))
}

// TestCreateTagRequest_OmittedElements tests that an empty array
// is present in the marshalled JSON.
func TestCreateTagRequest_OmittedElements(t *testing.T) {
	// Test with omitted elements
	request := CreateTagRequest{
		Tag: upcloud.Tag{
			Name: "foo",
		},
	}

	actualJSON, err := json.Marshal(&request)
	assert.NoError(t, err)

	expectedJSON := `
	{
		"tag": {
			"name": "foo",
			"servers": {
				"server": []
			}
		}
	}
	`
	assert.JSONEq(t, expectedJSON, string(actualJSON))
}

// TestModifyTagRequest tests that ModifyTagRequest marshals correctly
func TestModifyTagRequest(t *testing.T) {
	request := ModifyTagRequest{
		Name: "foo",
		Tag: upcloud.Tag{
			Name: "bar",
			Servers: []string{
				"foo1",
				"foo2",
			},
		},
	}

	// Check the request URL
	assert.Equal(t, "/tag/foo", request.RequestURL())

	// Check marshaling
	actualJSON, err := json.Marshal(&request)
	assert.NoError(t, err)

	expectedJSON := `
	  {
		  "tag": {
			  "name": "bar",
			  "servers": {
				  "server": [
					  "foo1",
					  "foo2"
				  ]
			  }
		  }
	  }
	`
	assert.JSONEq(t, expectedJSON, string(actualJSON))
}

// TestModifyTagRequest_OmitServers tests that an empty array is present in
// the marshalled JSON if there are no servers specified.
func TestModifyTagRequest_OmitServers(t *testing.T) {
	request := ModifyTagRequest{
		Name: "foo",
		Tag: upcloud.Tag{
			Name: "bar",
		},
	}

	// Check the request URL
	assert.Equal(t, "/tag/foo", request.RequestURL())

	// Check marshaling
	actualJSON, err := json.Marshal(&request)
	assert.NoError(t, err)

	expectedJSON := `
	  {
		  "tag": {
			  "name": "bar",
			  "servers": {
				  "server": []
			  }
		  }
	  }
	`
	assert.JSONEq(t, expectedJSON, string(actualJSON))
}

// TestDeleteTagRequest tests that DeleteTagRequest behaves correctly
func TestDeleteTagRequest(t *testing.T) {
	request := DeleteTagRequest{
		Name: "foo",
	}

	// Check the request URL
	assert.Equal(t, "/tag/foo", request.RequestURL())
}
