package upcloud

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnmarshalPartnerAccount(t *testing.T) {
	originalJSON := `
		{
			"address" : "Some street",
			"city" : "Some city",
			"company" : "Some company",
			"country" : "FIN",
			"email" : "some.user@somecompany.com",
			"first_name" : "Some",
			"last_name" : "User",
			"phone" : "+358.91234567",
			"postal_code" : "00100",
			"state" : "",
			"username" : "someuser",
			"vat_number" : ""
		}
	`

	expectedAccount := PartnerAccount{
		Username:   "someuser",
		FirstName:  "Some",
		LastName:   "User",
		Country:    "FIN",
		State:      "",
		Email:      "some.user@somecompany.com",
		Phone:      "+358.91234567",
		Company:    "Some company",
		Address:    "Some street",
		PostalCode: "00100",
		City:       "Some city",
		VATNumber:  "",
	}

	account := PartnerAccount{}
	err := json.Unmarshal([]byte(originalJSON), &account)
	assert.NoError(t, err)
	assert.Equal(t, expectedAccount, account)
}
