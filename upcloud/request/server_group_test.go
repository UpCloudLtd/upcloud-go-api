package request

import (
	"encoding/json"
	"testing"

	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud"
	"github.com/stretchr/testify/assert"
)

// ServerGroup
func TestGetServerGroupsRequest(t *testing.T) {
	r := GetServerGroupsRequest{}
	assert.Equal(t, "/server-group", r.RequestURL())
}

func TestGetServerGroupRequest(t *testing.T) {
	r := GetServerGroupRequest{UUID: "id"}
	assert.Equal(t, "/server-group/id", r.RequestURL())
}

func TestDeleteServerGroupRequest(t *testing.T) {
	r := DeleteServerGroupRequest{UUID: "id"}
	assert.Equal(t, "/server-group/id", r.RequestURL())
}

func TestCreateServerGroupRequest(t *testing.T) {
	expected := `
	{
		"server_group": {
			"title": "test"
		}
	}	
	`
	r := CreateServerGroupRequest{
		Title:   "test",
		Members: upcloud.ServerUUIDSlice{},
	}
	actual, err := json.Marshal(&r)
	assert.NoError(t, err)
	assert.JSONEq(t, expected, string(actual))

	expected = `
	{
		"server_group": {
			"title": "test",
			"servers": {
				"server": ["x", "y"]
			}
		}
	}	
	`
	r = CreateServerGroupRequest{
		Title:   "test",
		Members: upcloud.ServerUUIDSlice{"x", "y"},
	}
	actual, err = json.Marshal(&r)
	assert.NoError(t, err)
	assert.JSONEq(t, expected, string(actual))
	assert.Equal(t, "/server-group", r.RequestURL())
}

func TestModifyServerGroupRequest(t *testing.T) {
	expected := `
	{
		"server_group": {
			"title": "test"
		}
	}	
	`
	r := ModifyServerGroupRequest{
		UUID:  "id",
		Title: "test",
	}
	actual, err := json.Marshal(&r)
	assert.NoError(t, err)
	assert.JSONEq(t, expected, string(actual))

	expected = `
	{
		"server_group": {
			"title": "test",
			"servers": {
				"server": []
			}
		}
	}	
	`
	r = ModifyServerGroupRequest{
		UUID:    "id",
		Title:   "test",
		Members: &upcloud.ServerUUIDSlice{},
	}
	actual, err = json.Marshal(&r)
	assert.NoError(t, err)
	assert.JSONEq(t, expected, string(actual))

	expected = `
	{
		"server_group": {
			"title": "test",
			"servers": {
				"server": ["x"]
			}
		}
	}	
	`
	r = ModifyServerGroupRequest{
		UUID:    "id",
		Title:   "test",
		Members: &upcloud.ServerUUIDSlice{"x"},
	}
	actual, err = json.Marshal(&r)
	assert.NoError(t, err)
	assert.JSONEq(t, expected, string(actual))

	expected = `
	{
		"server_group": {
			"servers": {
				"server": ["x"]
			}
		}
	}	
	`
	r = ModifyServerGroupRequest{
		UUID:    "id",
		Members: &upcloud.ServerUUIDSlice{"x"},
	}
	actual, err = json.Marshal(&r)
	assert.NoError(t, err)
	assert.JSONEq(t, expected, string(actual))

	assert.Equal(t, "/server-group/id", r.RequestURL())
}
