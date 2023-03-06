package request

import (
	"encoding/json"
	"testing"

	"github.com/UpCloudLtd/upcloud-go-api/v6/upcloud"
	"github.com/stretchr/testify/assert"
)

// ServerGroup
func TestGetServerGroupsRequest(t *testing.T) {
	r := GetServerGroupsRequest{}
	assert.Equal(t, "/server-group", r.RequestURL())

	rWithFilters := GetServerGroupsRequest{
		Filters: []QueryFilter{
			FilterLabel{Label: upcloud.Label{
				Key:   "env",
				Value: "test",
			}},
			FilterLabelKey{Key: "managed"},
		},
	}

	assert.Equal(t, "/server-group?label=env%3Dtest&label=managed", rWithFilters.RequestURL())
}

// TestGetServerGroupsWithFiltersRequest tests that GetServerGroupsWithFiltersRequest objects behave correctly
func TestGetServerGroupsWithFiltersRequest(t *testing.T) {
	request := GetServerGroupsWithFiltersRequest{
		Filters: []ServerGroupFilter{
			FilterLabelKey{"onlyKey1"},
			FilterLabelKey{"onlyKey2"},
			FilterLabel{Label: upcloud.Label{
				Key:   "pairKey1",
				Value: "pairValue1",
			}},
			FilterLabel{Label: upcloud.Label{
				Key:   "pairKey2",
				Value: "pairValue2",
			}},
		},
	}

	assert.Equal(
		t,
		"/server-group?label=onlyKey1&label=onlyKey2&label=pairKey1%3DpairValue1&label=pairKey2%3DpairValue2",
		request.RequestURL(),
	)
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
			"labels": {
				"label": [
					{
						"key": "managedBy",
						"value": "upcloud-go-sdk-unit-test"
					}
				]
			},
			"servers": {
				"server": ["x", "y"]
			},
			"title": "test",
			"anti_affinity": "yes"
		}
	}	
	`
	r = CreateServerGroupRequest{
		Labels: &upcloud.LabelSlice{
			upcloud.Label{
				Key:   "managedBy",
				Value: "upcloud-go-sdk-unit-test",
			},
		},
		Members:      upcloud.ServerUUIDSlice{"x", "y"},
		Title:        "test",
		AntiAffinity: upcloud.True,
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
			},
			"anti_affinity": "no"
		}
	}	
	`
	r = ModifyServerGroupRequest{
		UUID:         "id",
		Title:        "test",
		Members:      &upcloud.ServerUUIDSlice{"x"},
		AntiAffinity: upcloud.False,
	}
	actual, err = json.Marshal(&r)
	assert.NoError(t, err)
	assert.JSONEq(t, expected, string(actual))

	expected = `
	{
		"server_group": {
			"servers": {
				"server": ["x"]
			},
			"anti_affinity": "yes"
		}
	}	
	`
	r = ModifyServerGroupRequest{
		UUID:         "id",
		Members:      &upcloud.ServerUUIDSlice{"x"},
		AntiAffinity: upcloud.True,
	}
	actual, err = json.Marshal(&r)
	assert.NoError(t, err)
	assert.JSONEq(t, expected, string(actual))

	assert.Equal(t, "/server-group/id", r.RequestURL())
}
