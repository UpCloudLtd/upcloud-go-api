package upcloud

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestUnmarshalServerGroup tests that ServerGroup struct is unmarshaled correctly
func TestUnmarshalServerGroup(t *testing.T) {
	originalJSON := `
	{
		"server_group" : {
		   "servers" : {
			  "server" : [
				 "x",
				 "y"
			  ]
		   },
		   "title" : "my group",
		   "uuid" : "server_uuid"
		}
	 }	 
	`
	actual := ServerGroup{}
	err := json.Unmarshal([]byte(originalJSON), &actual)
	assert.NoError(t, err)

	expected := ServerGroup{
		UUID:    "server_uuid",
		Title:   "my group",
		Members: []string{"x", "y"},
	}

	assert.Equal(t, expected, actual)
}

// TestUnmarshalServerGroup tests that ServerGroup struct is unmarshaled correctly
func TestUnmarshalServerGroups(t *testing.T) {
	originalJSON := `
	{
		"server_groups" : {
			"server_group" : [
				{
					"servers" : {
						"server" : [
							"x"
						]
					},
					"title" : "my group 1",
					"uuid" : "id"
				},
				{
					"servers" : {
						"server" : [
							"a",
							"b",
							"c"
						]
					},
					"title" : "my group 2",
					"uuid" : "id"
				}
			]
		}
	}
	`
	actual := ServerGroups{}
	err := json.Unmarshal([]byte(originalJSON), &actual)
	assert.NoError(t, err)

	expected := ServerGroups{
		{
			UUID:    "id",
			Title:   "my group 1",
			Members: []string{"x"},
		},
		{
			UUID:    "id",
			Title:   "my group 2",
			Members: []string{"a", "b", "c"},
		},
	}

	assert.Equal(t, expected, actual)
}
