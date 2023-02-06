package upcloud

import (
	"net/http"
	"testing"

	"github.com/UpCloudLtd/upcloud-go-api/v5/upcloud/client"
	"github.com/stretchr/testify/assert"
)

func TestNewError(t *testing.T) {
	rawErr := NewError(&client.Error{
		Type:      client.ErrorTypeError,
		ErrorCode: http.StatusNotFound,
		ResponseBody: []byte(`
			{
				"error": {
				  "error_message": "The server 00af0f73-7082-4283-b925-811d1585774b does not exist.",
				  "error_code": "SERVER_NOT_FOUND"
				}
			}
		`),
	})

	// fmt.Printf("\033[33m Err: %s \033[0m\n", rawErr.Error())

	err, ok := rawErr.(*Error)
	assert.True(t, ok)
	assert.Equal(t, "The server 00af0f73-7082-4283-b925-811d1585774b does not exist.", err.Message)
	assert.Equal(t, "SERVER_NOT_FOUND", err.Code)
	assert.Nil(t, err.problemError)
}

func TestNewErrorWithProblem(t *testing.T) {
	rawErr := NewError(&client.Error{
		Type:      client.ErrorTypeProblem,
		ErrorCode: http.StatusBadRequest,
		ResponseBody: []byte(`
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
		`),
	})

	err, ok := rawErr.(*Error)
	assert.True(t, ok)

	// Check basic (legacy) error properties
	assert.Equal(t, "INVALID_REQUEST", err.Code)
	assert.Equal(t, "Validation error.", err.Message)

	// Check json+problem properties
	assert.Equal(t, "Validation error.", err.problemError.Title)
	assert.Equal(t, "https://developers.upcloud.com/1.3/errors#ERROR_INVALID_REQUEST", err.problemError.Type)
	assert.Equal(t, "01GRKDQBGD7FA84MGR9373F093", err.problemError.CorrelationID)
	assert.Equal(t, 400, err.problemError.Status)
	assert.Len(t, err.problemError.InvalidParams, 1)
	assert.Equal(t, err.problemError.InvalidParams[0].Name, "plan")
	assert.Equal(t, err.problemError.InvalidParams[0].Reason, "Plan doesn't exist")
}
