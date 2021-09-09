package upcloud_test

import (
	"encoding/json"
	"testing"

	"github.com/UpCloudLtd/upcloud-go-api/upcloud"

	"github.com/stretchr/testify/assert"
)

// TestUnmarshalFirewallRules tests the FirewallRules and FirewallRule are unmarshaled correctly.
func TestError(t *testing.T) {
	t.Parallel()
	originalJSON := `
        {
            "error": {
              "error_message": "The server 00af0f73-7082-4283-b925-811d1585774b does not exist.",
              "error_code": "SERVER_NOT_FOUND"
            }
        }
    `

	e := upcloud.Error{}
	err := json.Unmarshal([]byte(originalJSON), &e)
	assert.NoError(t, err)
	assert.Equal(t, "The server 00af0f73-7082-4283-b925-811d1585774b does not exist.", e.ErrorMessage)
	assert.Equal(t, "SERVER_NOT_FOUND", e.ErrorCode)
}
