package upcloud

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProblem(t *testing.T) {
	p := Problem{}
	err := json.Unmarshal([]byte(`
	{
		"type": "https://developers.upcloud.com/1.3/errors#ERROR_INVALID_REQUEST",
		"title": "Validation error.",
		"invalid_params": [
			{
				"name": "default_backend",
				"reason": "Backend doesn't exist."
			}
		],
		"correlation_id": "01FY8RP81GDE07BAVYY7V4DKRY",
		"status": 400
	}
	
	`), &p)
	assert.NoError(t, err)
	assert.Equal(t, "Validation error.", p.Title)
	assert.Equal(t, "01FY8RP81GDE07BAVYY7V4DKRY", p.CorrelationID)
	assert.Equal(t, 400, p.Status)
	assert.Equal(t, "Backend doesn't exist.", p.InvalidParams[0].Reason)
	assert.Equal(t, "default_backend", p.InvalidParams[0].Name)
}
