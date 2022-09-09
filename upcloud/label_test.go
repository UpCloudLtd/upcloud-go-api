package upcloud

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnmarshalLabelSlice(t *testing.T) {
	originalJSON := `
		{
			"labels" : {
				"label" : [
					{
						"key" : "managedBy",
						"value" : "upcloud-go-sdk-unit-test"
					},
					{
						"key" : "test",
						"value" : "true"
					}
				]
			}
		}
	`

	testResource := struct {
		Labels LabelSlice `json:"labels"`
	}{}

	err := json.Unmarshal([]byte(originalJSON), &testResource)
	assert.Nil(t, err)
	assert.Len(t, testResource.Labels, 2)

	testData := []Label{
		{
			Key:   "managedBy",
			Value: "upcloud-go-sdk-unit-test",
		},
		{
			Key:   "test",
			Value: "true",
		},
	}

	for k, v := range testData {
		label := testResource.Labels[k]

		assert.Equal(t, label.Key, v.Key)
		assert.Equal(t, label.Value, v.Value)
	}
}
