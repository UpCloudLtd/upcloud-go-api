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
				"uuid" : "server_group_uuid",
				"anti_affinity": 1,
				"anti_affinity_status": [
					{
						"uuid": "x",
						"status": "met"
					},
					{
						"uuid": "y",
						"status": "unmet"
					}
				]
			}
		}
	`
	actual := ServerGroup{}
	err := json.Unmarshal([]byte(originalJSON), &actual)
	assert.NoError(t, err)

	expected := ServerGroup{
		Labels:       LabelSlice{Label{Key: "managedBy", Value: "upcloud-go-sdk-unit-test"}, Label{Key: "env", Value: "test"}},
		Members:      []string{"x", "y"},
		Title:        "my group",
		UUID:         "server_group_uuid",
		AntiAffinity: ServerGroupAntiAffinityYes,
		AntiAffinityStatus: []ServerGroupMemberAntiAffinityStatus{
			{
				ServerUUID: "x",
				Status:     ServerAntiAffinityStatusMet,
			},
			{
				ServerUUID: "y",
				Status:     ServerAntiAffinityStatusUnmet,
			},
		},
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
					"uuid" : "id",
					"anti_affinity": 0
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
					"uuid" : "id",
					"anti_affinity": 1,
					"anti_affinity_status": [
						{
							"uuid": "a",
							"status": "met"
						},
						{
							"uuid": "b",
							"status": "met"
						},
						{
							"uuid": "c",
							"status": "unmet"
						}
					]
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
			Labels:       LabelSlice{Label{Key: "managedBy", Value: "upcloud-go-sdk-unit-test"}},
			Members:      []string{"x"},
			Title:        "my group 1",
			UUID:         "id",
			AntiAffinity: ServerGroupAntiAffinityNo,
		},
		{
			Labels:       LabelSlice{Label{Key: "managedBy", Value: "upcloud-go-sdk-unit-test"}, Label{Key: "isSecondTestCase", Value: "true"}},
			Members:      []string{"a", "b", "c"},
			Title:        "my group 2",
			UUID:         "id",
			AntiAffinity: ServerGroupAntiAffinityStrict,
			AntiAffinityStatus: []ServerGroupMemberAntiAffinityStatus{
				{
					ServerUUID: "a",
					Status:     ServerAntiAffinityStatusMet,
				},
				{
					ServerUUID: "b",
					Status:     ServerAntiAffinityStatusMet,
				},
				{
					ServerUUID: "c",
					Status:     ServerAntiAffinityStatusUnmet,
				},
			},
		},
	}

	assert.Equal(t, expected, actual)
}
