package upcloud

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestUnmarshalTags tests that Tags structs are unmarshaled correctly
func TestUnmarshalTags(t *testing.T) {
	originalJSON := `
      {
        "tags": {
          "tag": [
            {
              "description": "Development servers",
              "name": "DEV",
              "servers": {
                "server": [
                  "0077fa3d-32db-4b09-9f5f-30d9e9afb565"
                ]
              }
            },
            {
              "description": "My own servers",
              "name": "private",
              "servers": {
                "server": [
                    "foo1",
                    "foo2"
                ]
              }
            },
            {
              "description": "Production servers",
              "name": "PROD",
              "servers": {
                "server": []
              }
            }
          ]
        }
      }
    `
	tags := Tags{}
	err := json.Unmarshal([]byte(originalJSON), &tags)
	assert.NoError(t, err)
	assert.Len(t, tags.Tags, 3)

	testData := []Tag{
		{
			Description: "Development servers",
			Name:        "DEV",
			Servers: []string{
				"0077fa3d-32db-4b09-9f5f-30d9e9afb565",
			},
		},
		{
			Description: "My own servers",
			Name:        "private",
			Servers: []string{
				"foo1",
				"foo2",
			},
		},
		{
			Description: "Production servers",
			Name:        "PROD",
			Servers:     []string{},
		},
	}

	for i, testTag := range testData {
		tag := tags.Tags[i]
		assert.Equal(t, testTag.Description, tag.Description)
		assert.Equal(t, testTag.Name, tag.Name)
		assert.Len(t, tag.Servers, len(testTag.Servers))
		for j, server := range testTag.Servers {
			assert.Equal(t, server, tag.Servers[j])
		}
	}
}

// TestUnmarshalTag tests that a single Tag struct is unmarshalled correctly.
func TestUnmarshalTag(t *testing.T) {
	originalJSON := `
    {
        "tag" : {
           "description" : "",
           "name" : "tag1",
           "servers" : {
              "server" : []
           }
        }
    }
    `

	tag := Tag{}
	err := json.Unmarshal([]byte(originalJSON), &tag)
	assert.NoError(t, err)

	assert.Equal(t, "", tag.Description)
	assert.Equal(t, "tag1", tag.Name)
	assert.Empty(t, tag.Servers)
}
