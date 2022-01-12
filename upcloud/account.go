package upcloud

import (
	"encoding/json"
)

type AccountType string

const (
	AccountTypeMain       AccountType = "main"
	AccountTypeSubaccount AccountType = "sub"
)

// Account represents an account
type Account struct {
	Credits        float64        `json:"credits"`
	UserName       string         `json:"username"`
	ResourceLimits ResourceLimits `json:"resource_limits"`
}

// ResourceLimits represents an account's resource limits
type ResourceLimits struct {
	Cores               int `json:"cores,omitempty"`
	DetachedFloatingIps int `json:"detached_floating_ips,omitempty"`
	Memory              int `json:"memory,omitempty"`
	Networks            int `json:"networks,omitempty"`
	PublicIPv4          int `json:"public_ipv4,omitempty"`
	PublicIPv6          int `json:"public_ipv6,omitempty"`
	StorageHDD          int `json:"storage_hdd,omitempty"`
	StorageSSD          int `json:"storage_ssd,omitempty"`
}

// UnmarshalJSON is a custom unmarshaller that deals with
// deeply embedded values.
func (s *Account) UnmarshalJSON(b []byte) error {
	type localAccount Account

	v := struct {
		Account localAccount `json:"account"`
	}{}
	err := json.Unmarshal(b, &v)
	if err != nil {
		return err
	}

	(*s) = Account(v.Account)

	return nil
}

// AccountCampaigns represents campaigns associated to account
type AccountCampaigns struct {
	Campaign []string `json:"campaign"`
}

// AccountRoles represents roles associated to account
// Roles for the account; billing, aux_billing, or technical.
type AccountRoles struct {
	Role []string `json:"role"`
}

// AccountTag represents tag and storage access permisission
type AccountTag struct {
	Name    string  `json:"name"`
	Storage Boolean `json:"storage"`
}

// AccountTagAccess represents tags associated to account
type AccountTagAccess struct {
	Tag []AccountTag `json:"tag"`
}

// AccountServer represents server (UUID) and storage access permisission
type AccountServer struct {
	UUID    string  `json:"uuid"`
	Storage Boolean `json:"storage"`
}

// AccountServerAccess represents servers the account is allowed to manage
type AccountServerAccess struct {
	Server []AccountServer `json:"server"`
}

// AccountStorageAccess represents UUIDs of storages the account is allowed to manage, or * for all
type AccountStorageAccess struct {
	Storage []string `json:"storage"`
}

// AccountNetworkAccess represents UUIDs of networks the account is allowed to manage, * means all
type AccountNetworkAccess struct {
	Network []string `json:"network"`
}

// AccountIPFilters represents IP address restrictions on API access; if set, allowed only from the specified ranges.
// Ranges can be specified in CIDR format, ranges separated by a dash, or as single IP addresses
type AccountIPFilters struct {
	IPFilter []string `json:"ip_filter"`
}

// AccountDetails represents detailed information about an account
type AccountDetails struct {
	MainAccount string      `json:"main_account"`
	Type        AccountType `json:"type"`
	Username    string      `json:"username"`
	FirstName   string      `json:"first_name"`
	LastName    string      `json:"last_name"`
	Company     string      `json:"company"`
	Address     string      `json:"address"`
	PostalCode  string      `json:"postal_code"`
	City        string      `json:"city"`
	Email       string      `json:"email"`

	// Phone number in international format, country code and national part separated by a period
	Phone string `json:"phone"`

	// U.S. state if applicable
	State string `json:"state"`

	// ISO 3166-1 three character country code
	Country       string               `json:"country"`
	Currency      string               `json:"currency"`
	Language      string               `json:"language"`
	VATNnumber    string               `json:"vat_number"`
	Timezone      string               `json:"timezone"`
	AllowAPI      Boolean              `json:"allow_api"`
	AllowGUI      Boolean              `json:"allow_gui"`
	SimpleBackup  Boolean              `json:"simple_backup"`
	TagAccess     AccountTagAccess     `json:"tag_access"`
	Campaigns     AccountCampaigns     `json:"campaigns"`
	Roles         AccountRoles         `json:"roles"`
	ServerAccess  AccountServerAccess  `json:"server_access"`
	StorageAccess AccountStorageAccess `json:"storage_access"`
	NetworkAccess AccountNetworkAccess `json:"network_access"`
	IPFilters     AccountIPFilters     `json:"ip_filters"`

	// Whether 3rd party services are allowed in the account's context when logged in the UpCloud control panel.
	// Consult the complete description in the control panel.
	Enable3rdPartyServices Boolean `json:"enable_3rd_party_services"`
}

// IsSubaccount checks if account is subaccount
func (a AccountDetails) IsSubaccount() bool {
	return a.Type == AccountTypeSubaccount
}

func (a *AccountDetails) UnmarshalJSON(b []byte) error {
	type ad AccountDetails
	v := struct {
		Account ad `json:"account"`
	}{}

	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}

	*a = AccountDetails(v.Account)

	return nil
}

// AccountListItem represents account list item
type AccountListItem struct {
	Type     string       `json:"type"`
	Username string       `json:"username"`
	Roles    AccountRoles `json:"roles"`
}

// AccountList represents account list
type AccountList []AccountListItem

func (a *AccountList) UnmarshalJSON(b []byte) error {
	v := struct {
		Accounts struct {
			Item []AccountListItem `json:"account"`
		} `json:"accounts"`
	}{}

	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}

	*a = append(*a, v.Accounts.Item...)
	return nil
}
