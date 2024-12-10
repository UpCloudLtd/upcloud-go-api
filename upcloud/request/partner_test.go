package request

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreatePartnerAccountRequestWithoutContactDetails(t *testing.T) {
	request := CreatePartnerAccountRequest{
		Username: "someuser",
		Password: "superSecret123",
	}
	expectedJSON := `
		{
			"username" : "someuser",
			"password" : "superSecret123"
		}
	`
	actualJSON, err := json.Marshal(&request)
	assert.NoError(t, err)
	assert.JSONEq(t, expectedJSON, string(actualJSON))
	assert.Equal(t, "/partner/accounts", request.RequestURL())
}

func TestCreatePartnerAccountRequestWithMinimalContactDetails(t *testing.T) {
	request := CreatePartnerAccountRequest{
		Username: "someuser",
		Password: "superSecret123",
		ContactDetails: &CreatePartnerAccountContactDetails{
			FirstName: "Some",
			LastName:  "User",
			Country:   "FIN",
			Email:     "some.user@somecompany.com",
			Phone:     "+358.91234567",
		},
	}
	expectedJSON := `
		{
			"username" : "someuser",
			"password" : "superSecret123",
			"contact_details" : {
				"first_name" : "Some",
				"last_name" : "User",
				"country" : "FIN",
				"email" : "some.user@somecompany.com",
				"phone" : "+358.91234567"
			}
		}
	`
	actualJSON, err := json.Marshal(&request)
	assert.NoError(t, err)
	assert.JSONEq(t, expectedJSON, string(actualJSON))
	assert.Equal(t, "/partner/accounts", request.RequestURL())
}
