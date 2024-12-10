package request

// CreatePartnerAccountContactDetails represents optional contact details in CreatePartnerAccountRequest
type CreatePartnerAccountContactDetails struct {
	FirstName  string `json:"first_name,omitempty"`
	LastName   string `json:"last_name,omitempty"`
	Country    string `json:"country,omitempty"`
	State      string `json:"state,omitempty"`
	Email      string `json:"email,omitempty"`
	Phone      string `json:"phone,omitempty"`
	Company    string `json:"company,omitempty"`
	Address    string `json:"address,omitempty"`
	PostalCode string `json:"postal_code,omitempty"`
	City       string `json:"city,omitempty"`
	VATNumber  string `json:"vat_number,omitempty"`
}

// CreatePartnerAccountRequest represents a request to create new main account for partner
type CreatePartnerAccountRequest struct {
	Username       string                              `json:"username"`
	Password       string                              `json:"password"`
	ContactDetails *CreatePartnerAccountContactDetails `json:"contact_details,omitempty"`
}

func (r *CreatePartnerAccountRequest) RequestURL() string {
	return "/partner/accounts"
}
