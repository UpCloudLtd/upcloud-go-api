package upcloud

// Account represents an account
type Account struct {
	Credits        float64        `xml:"credits"`
	UserName       string         `xml:"username"`
	ResourceLimits resourceLimits `xml:"resource_limits"`
}

type resourceLimits struct {
	Cores      int `xml:"cores"`
	Memory     int `xml:"memory"`
	Networks   int `xml:"networks"`
	PublicIpv4 int `xml:"public_ipv4"`
	PublicIpv6 int `xml:"public_ipv6"`
	StorageHdd int `xml:"storage_hdd"`
	StorageSsd int `xml:"storage_ssd"`
}

type AccountList struct {
	Accounts []AccountListAccount `xml:"account"`
}

type AccountListAccount struct {
	Roles    []string `xml:"roles>role"`
	Type     string   `xml:"type"`
	UserName string   `xml:"username"`
}

type AccountDetails struct {
	FirstName              string   `xml:"first_name"`
	LastName               string   `xml:"last_name"`
	Phone                  string   `xml:"phone"`
	Email                  string   `xml:"email"`
	Company                string   `xml:"company"`
	Address                string   `xml:"address"`
	City                   string   `xml:"city"`
	State                  string   `xml:"state"`
	PostalCode             int      `xml:"postal_code"`
	Country                string   `xml:"country"`
	Timezone               string   `xml:"timezone"`
	Currency               string   `xml:"currency"`
	AllowApi               string   `xml:"allow_api"`
	Enable3rdPartyServices string   `xml:"enable_3rd_party_services"`
	Language               string   `xml:"language"`
	Campaigns              []string `xml:"campaigns>campaign"`
	Roles                  []string `xml:"roles>role"`
	SimpleBackup           string   `xml:"simple_backup"`
	Type                   string   `xml:"type"`
	UserName               string   `xml:"username"`
	VatNumber              string   `xml:"vat_number"`
	NetworkAccess          []string `xml:"network_access>network"`
	ServerAccess           []server `xml:"server_access>server"`
	StorageAccess          []string `xml:"storage_access>storage"`
	TagAccess              []tag    `xml:"tag_access>tag"`
	IpFilter               []string `xml:"ip_filters>ip_filter"`
}

type server struct {
	Storage string `xml:"storage"`
	Uuid    string `xml:"uuid"`
}

type tag struct {
	Name    string `xml:"name"`
	Storage string `xml:"storage"`
}
