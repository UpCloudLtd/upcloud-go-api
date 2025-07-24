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
	Cores                 int `json:"cores,omitempty"`
	DetachedFloatingIps   int `json:"detached_floating_ips,omitempty"`
	ManagedObjectStorages int `json:"managed_object_storages,omitempty"`
	Memory                int `json:"memory,omitempty"`
	NetworkPeerings       int `json:"network_peerings,omitempty"`
	Networks              int `json:"networks,omitempty"`
	NTPExcessGiB          int `json:"ntp_excess_gib,omitempty"`
	PublicIPv4            int `json:"public_ipv4,omitempty"`
	PublicIPv6            int `json:"public_ipv6,omitempty"`
	StorageHDD            int `json:"storage_hdd,omitempty"`
	StorageMaxIOPS        int `json:"storage_maxiops,omitempty"`
	StorageSSD            int `json:"storage_ssd,omitempty"`
	LoadBalancers         int `json:"load_balancers,omitempty"`
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
	Type     AccountType  `json:"type"`
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

// BillingSummary represents billing summary for a specific month
type BillingSummary struct {
	Currency              string           `json:"currency"`
	TotalAmount           float64          `json:"total_amount"`
	Servers               *BillingCategory `json:"servers,omitempty"`
	ManagedDatabases      *BillingCategory `json:"managed_databases,omitempty"`
	ManagedObjectStorages *BillingCategory `json:"managed_object_storages,omitempty"`
	ManagedLoadbalancers  *BillingCategory `json:"managed_loadbalancers,omitempty"`
	ManagedKubernetes     *BillingCategory `json:"managed_kubernetes,omitempty"`
	NetworkGateways       *BillingCategory `json:"network_gateways,omitempty"`
}

// BillingCategory represents a billing category with its resources
type BillingCategory struct {
	TotalAmount          float64               `json:"total_amount"`
	Server               *BillingResourceGroup `json:"server,omitempty"`
	ManagedDatabase      *BillingResourceGroup `json:"managed_database,omitempty"`
	ManagedObjectStorage *BillingResourceGroup `json:"managed_object_storage,omitempty"`
	ManagedLoadbalancer  *BillingResourceGroup `json:"managed_loadbalancers,omitempty"`
	ManagedKubernetes    *BillingResourceGroup `json:"managed_kubernetes,omitempty"`
	NetworkGateway       *BillingResourceGroup `json:"network_gateway,omitempty"`
}

// BillingResourceGroup represents a group of resources with their details
type BillingResourceGroup struct {
	Resources   []BillingResource `json:"resources"`
	TotalAmount float64           `json:"total_amount"`
}

// BillingResource represents a billable resource
type BillingResource struct {
	ResourceID string                  `json:"resource_id"`
	Amount     float64                 `json:"amount"`
	Hours      int                     `json:"hours"`
	Details    []BillingResourceDetail `json:"details"`
}

// BillingResourceDetail represents detailed billing information for a resource
type BillingResourceDetail struct {
	Amount          float64 `json:"amount"`
	Hours           int     `json:"hours"`
	Plan            string  `json:"plan,omitempty"`
	Zone            string  `json:"zone,omitempty"`
	Size            int     `json:"size,omitempty"`
	Cores           int     `json:"cores,omitempty"`
	Memory          int     `json:"memory,omitempty"`
	Firewall        float64 `json:"firewall,omitempty"`
	Licenses        float64 `json:"licenses,omitempty"`
	SimpleBackup    float64 `json:"simple_backup,omitempty"`
	BillableSizeGiB int     `json:"billable_size_gib,omitempty"`
	Labels          []Label `json:"labels,omitempty"`
}

// UnmarshalJSON is a custom unmarshaller for BillingSummary
func (b *BillingSummary) UnmarshalJSON(data []byte) error {
	type localBillingSummary BillingSummary

	v := struct {
		*localBillingSummary
	}{
		localBillingSummary: (*localBillingSummary)(b),
	}

	return json.Unmarshal(data, &v)
}
