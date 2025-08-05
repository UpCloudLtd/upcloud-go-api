package upcloud

import "time"

// GatewayConfiguredStatus represents a desired status of the service
type GatewayConfiguredStatus string

// GatewayConfiguredStatus represents current, actual status of the service
type GatewayOperationalState string

// GatewayFeature represents a feature of the service
type GatewayFeature string

// GatewayTunnelOperationalState represents current, actual status of the tunnel
type GatewayTunnelOperationalState string

type (
	GatewayConnectionType          string
	GatewayRouteType               string
	GatewayIPSecAuthType           string
	GatewayIPSecAlgorithm          string
	GatewayIPSecIntegrityAlgorithm string
)

const (
	GatewayConfiguredStatusStarted GatewayConfiguredStatus = "started"
	GatewayConfiguredStatusStopped GatewayConfiguredStatus = "stopped"

	GatewayOperationalStatePending           GatewayOperationalState = "pending"
	GatewayOperationalStateSetupAgent        GatewayOperationalState = "setup-agent"
	GatewayOperationalStateSetupLinkNetwork  GatewayOperationalState = "setup-link-network"
	GatewayOperationalStateSetupServer       GatewayOperationalState = "setup-server"
	GatewayOperationalStateSetupNetwork      GatewayOperationalState = "setup-network"
	GatewayOperationalStateSetupGW           GatewayOperationalState = "setup-gw"
	GatewayOperationalStateSetupDNS          GatewayOperationalState = "setup-dns"
	GatewayOperationalStateCheckup           GatewayOperationalState = "checkup"
	GatewayOperationalStateRunning           GatewayOperationalState = "running"
	GatewayOperationalStateDeleteDNS         GatewayOperationalState = "delete-dns"
	GatewayOperationalStateDeleteNetwork     GatewayOperationalState = "delete-network"
	GatewayOperationalStateDeleteServer      GatewayOperationalState = "delete-server"
	GatewayOperationalStateDeleteLinkNetwork GatewayOperationalState = "delete-link-network"
	GatewayOperationalStateDeleteService     GatewayOperationalState = "delete-service"

	GatewayTunnelOperationalStateUninitialized GatewayTunnelOperationalState = "uninitialized"
	GatewayTunnelOperationalStateCreated       GatewayTunnelOperationalState = "created"
	GatewayTunnelOperationalStateConnecting    GatewayTunnelOperationalState = "connecting"
	GatewayTunnelOperationalStateEstabilished  GatewayTunnelOperationalState = "established"
	GatewayTunnelOperationalStateRekeying      GatewayTunnelOperationalState = "rekeying"
	GatewayTunnelOperationalStateRekeyed       GatewayTunnelOperationalState = "rekeyed"
	GatewayTunnelOperationalStateDeleting      GatewayTunnelOperationalState = "deleting"
	GatewayTunnelOperationalStateDestroying    GatewayTunnelOperationalState = "destroying"
	GatewayTunnelOperationalStateUnknown       GatewayTunnelOperationalState = "unknown"

	// GatewayFeatureNAT is a Network Address Translation (NAT) service that offers a way for cloud servers in SDN private networks to connect to the Internet through the public IP assigned to the network gateway service
	GatewayFeatureNAT GatewayFeature = "nat"

	// GatewayFeatureVPN is a Virtual Private Network (VPN) service used to establish an encrypted network connection when using public networks
	// Please note that VPN feature is currently in beta. You can learn more about it on its [product page]
	// Also note that VPN is available only in some of the gateway plans. To check which plans support VPN, you can use the GetGatewayPlans method.
	//
	// [product page]: https://upcloud.com/resources/docs/networking#nat-and-vpn-gateways
	GatewayFeatureVPN GatewayFeature = "vpn"

	GatewayConnectionTypeIPSec GatewayConnectionType = "ipsec"

	GatewayRouteTypeStatic GatewayRouteType = "static"

	GatewayTunnelIPSecAuthTypePSK GatewayIPSecAuthType = "psk"

	GatewayIPSecAlgorithm_aes128gcm16  GatewayIPSecAlgorithm = "aes128gcm16"
	GatewayIPSecAlgorithm_aes128gcm128 GatewayIPSecAlgorithm = "aes128gcm128"
	GatewayIPSecAlgorithm_aes192gcm16  GatewayIPSecAlgorithm = "aes192gcm16"
	GatewayIPSecAlgorithm_aes192gcm128 GatewayIPSecAlgorithm = "aes192gcm128"
	GatewayIPSecAlgorithm_aes256gcm16  GatewayIPSecAlgorithm = "aes256gcm16"
	GatewayIPSecAlgorithm_aes256gcm128 GatewayIPSecAlgorithm = "aes256gcm128"
	GatewayIPSecAlgorithm_aes128       GatewayIPSecAlgorithm = "aes128"
	GatewayIPSecAlgorithm_aes192       GatewayIPSecAlgorithm = "aes192"
	GatewayIPSecAlgorithm_aes256       GatewayIPSecAlgorithm = "aes256"

	GatewayIPSecIntegrityAlgorithm_aes128gmac GatewayIPSecIntegrityAlgorithm = "aes128gmac"
	GatewayIPSecIntegrityAlgorithm_aes256gmac GatewayIPSecIntegrityAlgorithm = "aes256gmac"
	GatewayIPSecIntegrityAlgorithm_sha1       GatewayIPSecIntegrityAlgorithm = "sha1"
	GatewayIPSecIntegrityAlgorithm_sha256     GatewayIPSecIntegrityAlgorithm = "sha256"
	GatewayIPSecIntegrityAlgorithm_sha384     GatewayIPSecIntegrityAlgorithm = "sha384"
	GatewayIPSecIntegrityAlgorithm_sha512     GatewayIPSecIntegrityAlgorithm = "sha512"
)

type Gateway struct {
	UUID             string                  `json:"uuid,omitempty"`
	Name             string                  `json:"name,omitempty"`
	Zone             string                  `json:"zone,omitempty"`
	Plan             string                  `json:"plan,omitempty"`
	Labels           []Label                 `json:"labels,omitempty"`
	ConfiguredStatus GatewayConfiguredStatus `json:"configured_status,omitempty"`
	OperationalState GatewayOperationalState `json:"operational_state,omitempty"`
	Features         []GatewayFeature        `json:"features,omitempty"`
	Routers          []GatewayRouter         `json:"routers,omitempty"`
	CreatedAt        time.Time               `json:"created_at,omitempty"`
	UpdatedAt        time.Time               `json:"updated_at,omitempty"`
	Addresses        []GatewayAddress        `json:"addresses,omitempty"`
	Connections      []GatewayConnection     `json:"connections,omitempty"`
}

type GatewayAddress struct {
	Address string `json:"address,omitempty"`
	Name    string `json:"name,omitempty"`
}

type GatewayRouter struct {
	CreatedAt time.Time `json:"created_at,omitempty"`
	UUID      string    `json:"uuid,omitempty"`
}

type GatewayConnection struct {
	UUID         string                `json:"uuid"`
	Name         string                `json:"name,omitempty"`
	Type         GatewayConnectionType `json:"type,omitempty"`
	LocalRoutes  []GatewayRoute        `json:"local_routes,omitempty"`
	RemoteRoutes []GatewayRoute        `json:"remote_routes,omitempty"`
	Tunnels      []GatewayTunnel       `json:"tunnels,omitempty"`
	CreatedAt    time.Time             `json:"created_at,omitempty"`
	UpdatedAt    time.Time             `json:"updated_at,omitempty"`
}

type GatewayRoute struct {
	Name          string           `json:"name,omitempty"`
	StaticNetwork string           `json:"static_network,omitempty"`
	Type          GatewayRouteType `json:"type,omitempty"`
}

type GatewayTunnel struct {
	UUID             string                        `json:"uuid,omitempty"`
	Name             string                        `json:"name,omitempty"`
	LocalAddress     GatewayTunnelLocalAddress     `json:"local_address,omitempty"`
	RemoteAddress    GatewayTunnelRemoteAddress    `json:"remote_address,omitempty"`
	IPSec            GatewayTunnelIPSec            `json:"ipsec,omitempty"`
	OperationalState GatewayTunnelOperationalState `json:"operational_state,omitempty"`
	CreatedAt        time.Time                     `json:"created_at,omitempty"`
	UpdatedAt        time.Time                     `json:"updated_at,omitempty"`
}

type GatewayTunnelLocalAddress struct {
	// Name of the UpCloud gateway address; should correspond to the name of one of the gateway address structs
	Name string `json:"name,omitempty"`
}

type GatewayTunnelRemoteAddress struct {
	// Address is a remote peer address VPN will connect to; must be global non-private unicast IP address.
	Address string `json:"address,omitempty"`
}

type GatewayTunnelIPSec struct {
	// Tunnel IPSec authentication object
	Authentication GatewayTunnelIPSecAuth `json:"authentication,omitempty"`
	// IKE SA rekey time in seconds
	RekeyTime int `json:"rekey_time,omitempty"`
	// IKE child SA rekey time in seconds
	ChildRekeyTime int `json:"child_rekey_time,omitempty"`
	// Delay before sending Dead Peer Detection packets if no traffic is detected, in seconds
	DPDDelay int `json:"dpd_delay,omitempty"`
	// Timeout period for DPD reply before considering the peer to be dead, in seconds
	DPDTimeout int `json:"dpd_timeout,omitempty"`
	// Maximum IKE SA lifetime in seconds
	IKELifetime int `json:"ike_lifetime,omitempty"`
	// List of Phase 1: Proposal algorithms
	Phase1Algorithms []GatewayIPSecAlgorithm `json:"phase1_algorithms,omitempty"`
	// List of Phase 1 integrity algorithms
	Phase1IntegrityAlgorithms []GatewayIPSecIntegrityAlgorithm `json:"phase1_integrity_algorithms,omitempty"`
	// List of Phase 1 Diffie-Hellman group numbers
	Phase1DHGroupNumbers []int `json:"phase1_dh_group_numbers,omitempty"`
	// List of Phase 2: Security Association algorithms
	Phase2Algorithms []GatewayIPSecAlgorithm `json:"phase2_algorithms,omitempty"`
	// List of Phase 2 integrity algorithms
	Phase2IntegrityAlgorithms []GatewayIPSecIntegrityAlgorithm `json:"phase2_integrity_algorithms,omitempty"`
	// List of Phase 2 Diffie-Hellman group numbers
	Phase2DHGroupNumbers []int `json:"phase2_dh_group_numbers,omitempty"`
}

type GatewayTunnelIPSecAuth struct {
	Authentication GatewayIPSecAuthType `json:"authentication,omitempty"`
	// PSK is a user-provided pre-shared key.
	// Note that this field is only meant to be used when providing API with your pre-shared key; it will always be empty in API responses
	PSK string `json:"psk,omitempty"`
}

type GatewayPlan struct {
	Name                     string           `json:"name,omitempty"`
	PerGatewayBandwidthMbps  int              `json:"per_gateway_bandwidth_mbps,omitempty"`
	PerGatewayMaxConnections int              `json:"per_gateway_max_connections,omitempty"`
	ServerNumber             int              `json:"server_number,omitempty"`
	SupportedFeatures        []GatewayFeature `json:"supported_features,omitempty"`
	VPNTunnelAmount          int              `json:"vpn_tunnel_amount,omitempty"`
}

// GatewayMetrics represents the response from the gateway metrics endpoint
type GatewayMetrics struct {
	Gateways     []GatewayMetricsGateway `json:"gateways,omitempty"`
	IPSecMetrics *GatewayIPSecMetrics    `json:"ipsec_metrics,omitempty"`
}

// GatewayMetricsGateway represents individual gateway metrics
type GatewayMetricsGateway struct {
	ActiveConnections        int       `json:"active_connections,omitempty"`
	CreatedAt                time.Time `json:"created_at,omitempty"`
	Name                     string    `json:"name,omitempty"`
	TotalAcceptedConnections int       `json:"total_accepted_connections,omitempty"`
	TotalRejectedSessions    int       `json:"total_rejected_sessions,omitempty"`
	UpdatedAt                time.Time `json:"updated_at,omitempty"`
}

// GatewayIPSecMetrics represents IPSec metrics
type GatewayIPSecMetrics struct {
	IKESAs []GatewayIKESA `json:"ike_sas,omitempty"`
}

// GatewayIKESA represents an IKE Security Association
type GatewayIKESA struct {
	ChildSAs         []GatewayChildSA       `json:"child_sas,omitempty"`
	CreatedAt        time.Time              `json:"created_at,omitempty"`
	Established      int                    `json:"established,omitempty"`
	Initiator        bool                   `json:"initator,omitempty"` // Note: API has typo "initator"
	LocalHost        string                 `json:"local_host,omitempty"`
	LocalID          string                 `json:"local_id,omitempty"`
	Name             string                 `json:"name,omitempty"`
	ReauthTime       int                    `json:"reauth_time,omitempty"`
	RekeyTime        int                    `json:"rekey_time,omitempty"`
	RemoteHost       string                 `json:"remote_host,omitempty"`
	RemoteID         string                 `json:"remote_id,omitempty"`
	OperationalState string                 `json:"operational_state,omitempty"`
	InternalState    string                 `json:"internal_state,omitempty"`
	TunnelID         int                    `json:"tunnel_id,omitempty"`
	UniqueID         string                 `json:"unique_id,omitempty"`
	UpdatedAt        time.Time              `json:"updated_at,omitempty"`
	Version          int                    `json:"version,omitempty"`
	HeuristicState   *GatewayHeuristicState `json:"heuristic_state,omitempty"`
}

// GatewayChildSA represents a Child Security Association
type GatewayChildSA struct {
	Name                   string    `json:"name,omitempty"`
	TunnelID               int       `json:"tunnel_id,omitempty"`
	UniqueID               string    `json:"unique_id,omitempty"`
	State                  string    `json:"state,omitempty"`
	SPIIn                  string    `json:"spi_in,omitempty"`
	SPIOut                 string    `json:"spi_out,omitempty"`
	BytesIn                int       `json:"bytes_in,omitempty"`
	BytesOut               int       `json:"bytes_out,omitempty"`
	PacketsIn              int       `json:"packets_in,omitempty"`
	PacketsOut             int       `json:"packets_out,omitempty"`
	RekeyTime              int       `json:"rekey_time,omitempty"`
	LifeTime               int       `json:"life_time,omitempty"`
	InstallTime            int       `json:"install_time,omitempty"`
	LocalTrafficSelectors  []string  `json:"local_traffic_selectors,omitempty"`
	RemoteTrafficSelectors []string  `json:"remote_traffic_selectors,omitempty"`
	CreatedAt              time.Time `json:"created_at,omitempty"`
	UpdatedAt              time.Time `json:"updated_at,omitempty"`
}

// GatewayHeuristicState represents heuristic tunnel state information
type GatewayHeuristicState struct {
	TunnelUp                 bool      `json:"tunnel_up,omitempty"`
	TunnelHealthy            bool      `json:"tunnel_healthy,omitempty"`
	LastDownMessage          string    `json:"last_down_message,omitempty"`
	UpdatedAt                time.Time `json:"updated_at,omitempty"`
	LastDownMessageUpdatedAt time.Time `json:"last_down_message_updated_at,omitempty"`
	UpEvents                 int       `json:"up_events,omitempty"`
	DownEvents               int       `json:"down_events,omitempty"`
	LogMessageBadEvents      int       `json:"log_message_bad_events,omitempty"`
}
