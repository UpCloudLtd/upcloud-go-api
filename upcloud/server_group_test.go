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
				"labels" : {
					"label" : [
						{
							"key" : "managedBy",
							"value" : "upcloud-go-sdk-unit-test"
						},
						{
							"key" : "env",
							"value" : "test"
						}
					]
				},
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
		Labels:  LabelSlice{Label{Key: "managedBy", Value: "upcloud-go-sdk-unit-test"}, Label{Key: "env", Value: "test"}},
		Members: []string{"x", "y"},
		Title:   "my group",
		UUID:    "server_uuid",
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
					"labels" : {
						"label" : [
							{
								"key" : "managedBy",
								"value" : "upcloud-go-sdk-unit-test"
							}
						]
					},
					"servers" : {
						"server" : [
							"x"
						]
					},
					"title" : "my group 1",
					"uuid" : "id"
				},
				{
					"labels" : {
						"label" : [
							{
								"key" : "managedBy",
								"value" : "upcloud-go-sdk-unit-test"
							},
							{
								"key" : "isSecondTestCase",
								"value" : "true"
							}
						]
					},
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
			Labels:  LabelSlice{Label{Key: "managedBy", Value: "upcloud-go-sdk-unit-test"}},
			Members: []string{"x"},
			Title:   "my group 1",
			UUID:    "id",
		},
		{
			Labels:  LabelSlice{Label{Key: "managedBy", Value: "upcloud-go-sdk-unit-test"}, Label{Key: "isSecondTestCase", Value: "true"}},
			Members: []string{"a", "b", "c"},
			Title:   "my group 2",
			UUID:    "id",
		},
	}

	assert.Equal(t, expected, actual)
}
