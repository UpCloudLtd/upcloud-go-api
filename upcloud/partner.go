package upcloud

// PartnerAccount represents details of an account associated with a partner
type PartnerAccount struct {
	Username   string `json:"username"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Country    string `json:"country"`
	State      string `json:"state"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`
	Company    string `json:"company"`
	Address    string `json:"address"`
	PostalCode string `json:"postal_code"`
	City       string `json:"city"`
	VATNumber  string `json:"vat_number"`
}
