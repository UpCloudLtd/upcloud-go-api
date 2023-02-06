package upcloud

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrorLegacyUnmarshal(t *testing.T) {
	originalJSON := `
        {
            "error": {
              "error_message": "The server 00af0f73-7082-4283-b925-811d1585774b does not exist.",
              "error_code": "SERVER_NOT_FOUND"
            }
        }
    `

	e := Error{}
	err := json.Unmarshal([]byte(originalJSON), &e)
	assert.NoError(t, err)
	assert.Equal(t, "The server 00af0f73-7082-4283-b925-811d1585774b does not exist.", e.Message)
	assert.Equal(t, "SERVER_NOT_FOUND", e.Code)
	assert.Nil(t, e.problemError)
}

func TestErrorProblemUnmarshal(t *testing.T) {
	originalJSON := `
		{
			"type": "https://developers.upcloud.com/1.3/errors#ERROR_INVALID_REQUEST",
			"title": "Validation error.",
			"correlation_id": "01GRKDQBGD7FA84MGR9373F093",
			"invalid_params": [
				{
					"name": "plan",
					"reason": "Plan doesn't exist"
				}
			],
			"status": 400
		}
	`

	e := Error{}
	err := json.Unmarshal([]byte(originalJSON), &e)
	assert.NoError(t, err)

	// Check basic (legacy) error properties
	assert.Equal(t, "INVALID_REQUEST", e.Code)
	assert.Equal(t, "Validation error.", e.Message)

	// Check json+problem properties
	assert.Equal(t, "Validation error.", e.problemError.Title)
	assert.Equal(t, "https://developers.upcloud.com/1.3/errors#ERROR_INVALID_REQUEST", e.problemError.Type)
	assert.Equal(t, "01GRKDQBGD7FA84MGR9373F093", e.problemError.CorrelationID)
	assert.Equal(t, 400, e.problemError.Status)
	assert.Len(t, e.problemError.InvalidParams, 1)
	assert.Equal(t, e.problemError.InvalidParams[0].Name, "plan")
	assert.Equal(t, e.problemError.InvalidParams[0].Reason, "Plan doesn't exist")
}
