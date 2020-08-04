package upcloud

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestUnmarshalAccount tests that Account objects unmarshal correctly
func TestUnmarshalAccount(t *testing.T) {
	originalJSON := `
	{
      "account": {
        "credits": 9972.2324,
        "username": "username"
      }
    }
	`

	account := Account{}
	err := json.Unmarshal([]byte(originalJSON), &account)
	assert.NoError(t, err)
	assert.Equal(t, 9972.2324, account.Credits)
	assert.Equal(t, "username", account.UserName)
}
