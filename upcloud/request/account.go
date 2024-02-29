package request

import (
	"fmt"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
)

// CreateSubaccount represents data required to create a sub account
type CreateSubaccount struct {
	Username   string `json:"username"`
	Password   string `json:"password"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Company    string `json:"company"`
	Address    string `json:"address"`
	PostalCode string `json:"postal_code"`
	City       string `json:"city"`
	Email      string `json:"email"`

	// Phone number in international format, country code and national part separated by a period
	Phone string `json:"phone"`

	// U.S. state if applicable
	State string `json:"state"`

	// ISO 3166-1 three character country code
	Country       string                       `json:"country"`
	Currency      string                       `json:"currency"`
	Language      string                       `json:"language"`
	VATNnumber    string                       `json:"vat_number"`
	Timezone      string                       `json:"timezone"`
	AllowAPI      upcloud.Boolean              `json:"allow_api"`
	AllowGUI      upcloud.Boolean              `json:"allow_gui"`
	TagAccess     upcloud.AccountTagAccess     `json:"tag_access"`
	Roles         upcloud.AccountRoles         `json:"roles"`
	ServerAccess  upcloud.AccountServerAccess  `json:"server_access"`
	StorageAccess upcloud.AccountStorageAccess `json:"storage_access"`
	NetworkAccess upcloud.AccountNetworkAccess `json:"network_access"`
	IPFilters     upcloud.AccountIPFilters     `json:"ip_filters"`
}

// CreateSubaccountRequest represents a request to create a sub account
type CreateSubaccountRequest struct {
	Subaccount CreateSubaccount `json:"sub_account"`
}

func (r CreateSubaccountRequest) RequestURL() string {
	return "/account/sub"
}

// ModifySubaccount represents data required to modify a Subaccount
type ModifySubaccount struct {
	FirstName  string `json:"first_name,omitempty"`
	LastName   string `json:"last_name,omitempty"`
	Company    string `json:"company,omitempty"`
	Address    string `json:"address,omitempty"`
	PostalCode string `json:"postal_code,omitempty"`
	City       string `json:"city,omitempty"`
	Email      string `json:"email,omitempty"`

	// Phone number in international format, country code and national part separated by a period
	Phone string `json:"phone,omitempty"`

	// U.S. state if applicable
	State string `json:"state"`

	// ISO 3166-1 three character country code
	Country       string                       `json:"country,omitempty"`
	Currency      string                       `json:"currency,omitempty"`
	Language      string                       `json:"language,omitempty"`
	VATNnumber    string                       `json:"vat_number"`
	Timezone      string                       `json:"timezone,omitempty"`
	AllowAPI      upcloud.Boolean              `json:"allow_api,omitempty"`
	AllowGUI      upcloud.Boolean              `json:"allow_gui,omitempty"`
	TagAccess     upcloud.AccountTagAccess     `json:"tag_access,omitempty"`
	Roles         upcloud.AccountRoles         `json:"roles,omitempty"`
	ServerAccess  upcloud.AccountServerAccess  `json:"server_access,omitempty"`
	StorageAccess upcloud.AccountStorageAccess `json:"storage_access,omitempty"`
	NetworkAccess upcloud.AccountNetworkAccess `json:"network_access,omitempty"`
	IPFilters     upcloud.AccountIPFilters     `json:"ip_filters,omitempty"`
}

// ModifySubaccountRequest represents a request to modify a Subaccount
type ModifySubaccountRequest struct {
	Username   string           `json:"-"`
	Subaccount ModifySubaccount `json:"account"`
}

func (r ModifySubaccountRequest) RequestURL() string {
	return fmt.Sprintf("/account/sub/%s", r.Username)
}

// DeleteSubaccountRequest represents a request to delete a subaccount
type DeleteSubaccountRequest struct {
	Username string
}

// RequestURL implements the Request interface
func (r *DeleteSubaccountRequest) RequestURL() string {
	return fmt.Sprintf("/account/sub/%s", r.Username)
}

// GetAccountDetailsRequest represents a request to get account details
type GetAccountDetailsRequest struct {
	Username string
}

// RequestURL implements the Request interface
func (r *GetAccountDetailsRequest) RequestURL() string {
	return fmt.Sprintf("/account/details/%s", r.Username)
}
