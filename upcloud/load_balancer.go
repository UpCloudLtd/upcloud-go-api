package upcloud

import (
	"time"
)

type (
	LoadBalancerMode                              string
	LoadBalancerMatcherType                       string
	LoadBalancerMatchingCondition                 string
	LoadBalancerActionType                        string
	LoadBalancerActionHTTPRedirectScheme          string
	LoadBalancerStringMatcherMethod               string
	LoadBalancerHTTPMatcherMethod                 string
	LoadBalancerIntegerMatcherMethod              string
	LoadBalancerBackendMemberType                 string
	LoadBalancerOperationalState                  string
	LoadBalancerCertificateBundleOperationalState string
	LoadBalancerConfiguredStatus                  string
	LoadBalancerCertificateBundleType             string
	LoadBalancerProxyProtocolVersion              string
	LoadBalancerHealthCheckType                   string
	LoadBalancerNetworkType                       string
	LoadBalancerAddressFamily                     string
	LoadBalancerNodeOperationalState              string
	LoadBalancerMaintenanceDOW                    string
)

const (
	LoadBalancerModeHTTP LoadBalancerMode = "http"
	LoadBalancerModeTCP  LoadBalancerMode = "tcp"

	LoadBalancerBackendMemberTypeStatic  LoadBalancerBackendMemberType = "static"
	LoadBalancerBackendMemberTypeDynamic LoadBalancerBackendMemberType = "dynamic"

	LoadBalancerConfiguredStatusStarted LoadBalancerConfiguredStatus = "started"
	LoadBalancerConfiguredStatusStopped LoadBalancerConfiguredStatus = "stopped"

	LoadBalancerCertificateBundleTypeManual    LoadBalancerCertificateBundleType = "manual"
	LoadBalancerCertificateBundleTypeDynamic   LoadBalancerCertificateBundleType = "dynamic"
	LoadBalancerCertificateBundleTypeAuthority LoadBalancerCertificateBundleType = "authority"

	LoadBalancerOperationalStatePending       LoadBalancerOperationalState = "pending"
	LoadBalancerOperationalStateSetupAgent    LoadBalancerOperationalState = "setup-agent"
	LoadBalancerOperationalStateSetupServer   LoadBalancerOperationalState = "setup-server"
	LoadBalancerOperationalStateSetupNetwork  LoadBalancerOperationalState = "setup-network"
	LoadBalancerOperationalStateSetupLB       LoadBalancerOperationalState = "setup-lb"
	LoadBalancerOperationalStateSetupDNS      LoadBalancerOperationalState = "setup-dns"
	LoadBalancerOperationalStateCheckup       LoadBalancerOperationalState = "checkup"
	LoadBalancerOperationalStateRunning       LoadBalancerOperationalState = "running"
	LoadBalancerOperationalStateDeleteDNS     LoadBalancerOperationalState = "delete-dns"
	LoadBalancerOperationalStateDeleteNetwork LoadBalancerOperationalState = "delete-network"
	LoadBalancerOperationalStateDeleteServer  LoadBalancerOperationalState = "delete-server"
	LoadBalancerOperationalStateDeleteService LoadBalancerOperationalState = "delete-service"

	LoadBalancerCertificateBundleOperationalStateIdle              LoadBalancerCertificateBundleOperationalState = "idle"
	LoadBalancerCertificateBundleOperationalStatePending           LoadBalancerCertificateBundleOperationalState = "pending"
	LoadBalancerCertificateBundleOperationalStateSetupChallenge    LoadBalancerCertificateBundleOperationalState = "setup-challenge"
	LoadBalancerCertificateBundleOperationalStateCompleteChallenge LoadBalancerCertificateBundleOperationalState = "complete-challenge"

	LoadBalancerMatcherTypeSrcIP          LoadBalancerMatcherType = "src_ip"
	LoadBalancerMatcherTypeSrcPort        LoadBalancerMatcherType = "src_port"
	LoadBalancerMatcherTypeBodySize       LoadBalancerMatcherType = "body_size"
	LoadBalancerMatcherTypePath           LoadBalancerMatcherType = "path"
	LoadBalancerMatcherTypeURL            LoadBalancerMatcherType = "url"
	LoadBalancerMatcherTypeURLQuery       LoadBalancerMatcherType = "url_query"
	LoadBalancerMatcherTypeHost           LoadBalancerMatcherType = "host"
	LoadBalancerMatcherTypeHTTPMethod     LoadBalancerMatcherType = "http_method"
	LoadBalancerMatcherTypeHTTPStatus     LoadBalancerMatcherType = "http_status"
	LoadBalancerMatcherTypeCookie         LoadBalancerMatcherType = "cookie"
	LoadBalancerMatcherTypeHeader         LoadBalancerMatcherType = "header"
	LoadBalancerMatcherTypeRequestHeader  LoadBalancerMatcherType = "request_header"
	LoadBalancerMatcherTypeResponseHeader LoadBalancerMatcherType = "response_header"
	LoadBalancerMatcherTypeURLParam       LoadBalancerMatcherType = "url_param"
	LoadBalancerMatcherTypeNumMembersUp   LoadBalancerMatcherType = "num_members_up"

	LoadBalancerMatchingConditionAnd LoadBalancerMatchingCondition = "and"
	LoadBalancerMatchingConditionOr  LoadBalancerMatchingCondition = "or"

	LoadBalancerActionTypeUseBackend          LoadBalancerActionType = "use_backend"
	LoadBalancerActionTypeTCPReject           LoadBalancerActionType = "tcp_reject"
	LoadBalancerActionTypeHTTPReturn          LoadBalancerActionType = "http_return"
	LoadBalancerActionTypeHTTPRedirect        LoadBalancerActionType = "http_redirect"
	LoadBalancerActionTypeSetForwardedHeaders LoadBalancerActionType = "set_forwarded_headers"
	LoadBalancerActionTypeSetRequestHeader    LoadBalancerActionType = "set_request_header"
	LoadBalancerActionTypeSetResponseHeader   LoadBalancerActionType = "set_response_header"

	LoadBalancerActionHTTPRedirectSchemeHTTP  LoadBalancerActionHTTPRedirectScheme = "http"
	LoadBalancerActionHTTPRedirectSchemeHTTPS LoadBalancerActionHTTPRedirectScheme = "https"

	LoadBalancerStringMatcherMethodExact     LoadBalancerStringMatcherMethod = "exact"
	LoadBalancerStringMatcherMethodSubstring LoadBalancerStringMatcherMethod = "substring"
	LoadBalancerStringMatcherMethodRegexp    LoadBalancerStringMatcherMethod = "regexp"
	LoadBalancerStringMatcherMethodStarts    LoadBalancerStringMatcherMethod = "starts"
	LoadBalancerStringMatcherMethodEnds      LoadBalancerStringMatcherMethod = "ends"
	LoadBalancerStringMatcherMethodDomain    LoadBalancerStringMatcherMethod = "domain"
	LoadBalancerStringMatcherMethodIP        LoadBalancerStringMatcherMethod = "ip"
	LoadBalancerStringMatcherMethodExists    LoadBalancerStringMatcherMethod = "exists"

	LoadBalancerHTTPMatcherMethodGet     LoadBalancerHTTPMatcherMethod = "GET"
	LoadBalancerHTTPMatcherMethodHead    LoadBalancerHTTPMatcherMethod = "HEAD"
	LoadBalancerHTTPMatcherMethodPost    LoadBalancerHTTPMatcherMethod = "POST"
	LoadBalancerHTTPMatcherMethodPut     LoadBalancerHTTPMatcherMethod = "PUT"
	LoadBalancerHTTPMatcherMethodPatch   LoadBalancerHTTPMatcherMethod = "PATCH"
	LoadBalancerHTTPMatcherMethodDelete  LoadBalancerHTTPMatcherMethod = "DELETE"
	LoadBalancerHTTPMatcherMethodConnect LoadBalancerHTTPMatcherMethod = "CONNECT"
	LoadBalancerHTTPMatcherMethodOptions LoadBalancerHTTPMatcherMethod = "OPTIONS"
	LoadBalancerHTTPMatcherMethodTrace   LoadBalancerHTTPMatcherMethod = "TRACE"

	LoadBalancerIntegerMatcherMethodEqual          LoadBalancerIntegerMatcherMethod = "equal"
	LoadBalancerIntegerMatcherMethodGreaterOrEqual LoadBalancerIntegerMatcherMethod = "greater_or_equal"
	LoadBalancerIntegerMatcherMethodGreater        LoadBalancerIntegerMatcherMethod = "greater"
	LoadBalancerIntegerMatcherMethodLessOrEqual    LoadBalancerIntegerMatcherMethod = "less_or_equal"
	LoadBalancerIntegerMatcherMethodLess           LoadBalancerIntegerMatcherMethod = "less"
	LoadBalancerIntegerMatcherMethodRange          LoadBalancerIntegerMatcherMethod = "range"

	LoadBalancerProxyProtocolVersion1 LoadBalancerProxyProtocolVersion = "v1"
	LoadBalancerProxyProtocolVersion2 LoadBalancerProxyProtocolVersion = "v2"

	LoadBalancerHealthCheckTypeTCP  LoadBalancerHealthCheckType = "tcp"
	LoadBalancerHealthCheckTypeHTTP LoadBalancerHealthCheckType = "http"

	LoadBalancerNetworkTypePublic  LoadBalancerNetworkType   = "public"
	LoadBalancerNetworkTypePrivate LoadBalancerNetworkType   = "private"
	LoadBalancerAddressFamilyIPv4  LoadBalancerAddressFamily = "IPv4"
	// IPv6 addresses are not supported yet
	// LoadBalancerAddressFamilyIPv6  LoadBalancerAddressFamily = "IPv6"

	LoadBalancerNodeOperationalStatePending               LoadBalancerNodeOperationalState = "pending"
	LoadBalancerNodeOperationalStatePullConfig            LoadBalancerNodeOperationalState = "pull-config"
	LoadBalancerNodeOperationalStateSetupLB               LoadBalancerNodeOperationalState = "setup-lb"
	LoadBalancerNodeOperationalStateRunning               LoadBalancerNodeOperationalState = "running"
	LoadBalancerNodeOperationalStateFailing               LoadBalancerNodeOperationalState = "failing"
	LoadBalancerNodeOperationalStateAgentUpgradeStarting  LoadBalancerNodeOperationalState = "agent-upgrade-starting"
	LoadBalancerNodeOperationalStateAgentUpgradeFinishing LoadBalancerNodeOperationalState = "agent-upgrade-finishing"
	LoadBalancerNodeOperationalStateStopped               LoadBalancerNodeOperationalState = "stopped"
	LoadBalancerNodeOperationalStateNotResponding         LoadBalancerNodeOperationalState = "not-responding"

	LoadBalancerMaintenanceDOWMonday    LoadBalancerMaintenanceDOW = "monday"
	LoadBalancerMaintenanceDOWTuesday   LoadBalancerMaintenanceDOW = "tuesday"
	LoadBalancerMaintenanceDOWWednesday LoadBalancerMaintenanceDOW = "wednesday"
	LoadBalancerMaintenanceDOWThursday  LoadBalancerMaintenanceDOW = "thursday"
	LoadBalancerMaintenanceDOWFriday    LoadBalancerMaintenanceDOW = "friday"
	LoadBalancerMaintenanceOWSaturday   LoadBalancerMaintenanceDOW = "saturday"
	LoadBalancerMaintenanceDOWSunday    LoadBalancerMaintenanceDOW = "sunday"
)

// LoadBalancerPlan represents load balancer plan details
type LoadBalancerPlan struct {
	Name                 string `json:"name,omitempty"`
	PerServerMaxSessions int    `json:"per_server_max_sessions,omitempty"`
	ServerNumber         int    `json:"server_number,omitempty"`
}

// LoadBalancerFrontend represents service frontend
type LoadBalancerFrontend struct {
	Name           string                          `json:"name,omitempty"`
	Mode           LoadBalancerMode                `json:"mode,omitempty"`
	Port           int                             `json:"port,omitempty"`
	Networks       []LoadBalancerFrontendNetwork   `json:"networks,omitempty"`
	DefaultBackend string                          `json:"default_backend,omitempty"`
	Rules          []LoadBalancerFrontendRule      `json:"rules,omitempty"`
	TLSConfigs     []LoadBalancerFrontendTLSConfig `json:"tls_configs,omitempty"`
	Properties     *LoadBalancerFrontendProperties `json:"properties,omitempty"`
	CreatedAt      time.Time                       `json:"created_at,omitempty"`
	UpdatedAt      time.Time                       `json:"updated_at,omitempty"`
}

// LoadBalancerFrontendRule represents frontend rule
type LoadBalancerFrontendRule struct {
	Name              string                        `json:"name,omitempty"`
	Priority          int                           `json:"priority,omitempty"`
	MatchingCondition LoadBalancerMatchingCondition `json:"matching_condition,omitempty"`
	Matchers          []LoadBalancerMatcher         `json:"matchers,omitempty"`
	Actions           []LoadBalancerAction          `json:"actions,omitempty"`
	CreatedAt         time.Time                     `json:"created_at,omitempty"`
	UpdatedAt         time.Time                     `json:"updated_at,omitempty"`
}

// LoadBalancerFrontendTLSConfig represents frontend TLS configuration
type LoadBalancerFrontendTLSConfig struct {
	Name                  string    `json:"name,omitempty"`
	CertificateBundleUUID string    `json:"certificate_bundle_uuid,omitempty"`
	CreatedAt             time.Time `json:"created_at,omitempty"`
	UpdatedAt             time.Time `json:"updated_at,omitempty"`
}

// LoadBalancerFrontendProperties represents frontend properties
type LoadBalancerFrontendProperties struct {
	TimeoutClient        int   `json:"timeout_client,omitempty"`
	InboundProxyProtocol *bool `json:"inbound_proxy_protocol,omitempty"`
	HTTP2Enabled         *bool `json:"http2_enabled,omitempty"`
}

// LoadBalancerBackend represents service backend
type LoadBalancerBackend struct {
	Name       string                         `json:"name"`
	Members    []LoadBalancerBackendMember    `json:"members"`
	Resolver   string                         `json:"resolver,omitempty"`
	Properties *LoadBalancerBackendProperties `json:"properties,omitempty"`
	TLSConfigs []LoadBalancerBackendTLSConfig `json:"tls_configs,omitempty"`
	CreatedAt  time.Time                      `json:"created_at,omitempty"`
	UpdatedAt  time.Time                      `json:"updated_at,omitempty"`
}

// LoadBalancerBackendMember represents backend member
type LoadBalancerBackendMember struct {
	Name        string                        `json:"name"`
	IP          string                        `json:"ip"`
	Port        int                           `json:"port"`
	Weight      int                           `json:"weight"`
	MaxSessions int                           `json:"max_sessions"`
	Type        LoadBalancerBackendMemberType `json:"type"`
	Enabled     bool                          `json:"enabled"`
	CreatedAt   time.Time                     `json:"created_at,omitempty"`
	UpdatedAt   time.Time                     `json:"updated_at,omitempty"`
}

// LoadBalancerBackendTLSConfig represents backend TLS configuration
type LoadBalancerBackendTLSConfig struct {
	Name                  string    `json:"name,omitempty"`
	CertificateBundleUUID string    `json:"certificate_bundle_uuid,omitempty"`
	CreatedAt             time.Time `json:"created_at,omitempty"`
	UpdatedAt             time.Time `json:"updated_at,omitempty"`
}

// LoadBalancerBackendProperties represents backend properties
type LoadBalancerBackendProperties struct {
	TimeoutServer             int                              `json:"timeout_server,omitempty"`
	TimeoutTunnel             int                              `json:"timeout_tunnel,omitempty"`
	HealthCheckTLSVerify      *bool                            `json:"health_check_tls_verify,omitempty"`
	HealthCheckType           LoadBalancerHealthCheckType      `json:"health_check_type,omitempty"`
	HealthCheckInterval       int                              `json:"health_check_interval,omitempty"`
	HealthCheckFall           int                              `json:"health_check_fall,omitempty"`
	HealthCheckRise           int                              `json:"health_check_rise,omitempty"`
	HealthCheckURL            string                           `json:"health_check_url,omitempty"`
	HealthCheckExpectedStatus int                              `json:"health_check_expected_status,omitempty"`
	StickySessionCookieName   string                           `json:"sticky_session_cookie_name,omitempty"`
	OutboundProxyProtocol     LoadBalancerProxyProtocolVersion `json:"outbound_proxy_protocol,omitempty"`
	TLSEnabled                *bool                            `json:"tls_enabled,omitempty"`
	TLSVerify                 *bool                            `json:"tls_verify,omitempty"`
	TLSUseSystemCA            *bool                            `json:"tls_use_system_ca,omitempty"`
	HTTP2Enabled              *bool                            `json:"http2_enabled,omitempty"`
}

// LoadBalancerResolver represents domain name resolver
type LoadBalancerResolver struct {
	Name         string    `json:"name,omitempty"`
	Nameservers  []string  `json:"nameservers,omitempty"`
	Retries      int       `json:"retries,omitempty"`
	Timeout      int       `json:"timeout,omitempty"`
	TimeoutRetry int       `json:"timeout_retry,omitempty"`
	CacheValid   int       `json:"cache_valid,omitempty"`
	CacheInvalid int       `json:"cache_invalid,omitempty"`
	CreatedAt    time.Time `json:"created_at,omitempty"`
	UpdatedAt    time.Time `json:"updated_at,omitempty"`
}

// LoadBalancer service
type LoadBalancer struct {
	UUID             string                          `json:"uuid,omitempty"`
	Name             string                          `json:"name,omitempty"`
	Zone             string                          `json:"zone,omitempty"`
	Plan             string                          `json:"plan,omitempty"`
	NetworkUUID      string                          `json:"network_uuid,omitempty"` // deprecated
	Networks         []LoadBalancerNetwork           `json:"networks,omitempty"`
	IPAddresses      []LoadBalancerFloatingIPAddress `json:"ip_addresses,omitempty"`
	DNSName          string                          `json:"dns_name,omitempty"` // deprecated
	Labels           []Label                         `json:"labels,omitempty"`
	ConfiguredStatus LoadBalancerConfiguredStatus    `json:"configured_status,omitempty"`
	OperationalState LoadBalancerOperationalState    `json:"operational_state,omitempty"`
	Frontends        []LoadBalancerFrontend          `json:"frontends,omitempty"`
	Backends         []LoadBalancerBackend           `json:"backends,omitempty"`
	Resolvers        []LoadBalancerResolver          `json:"resolvers,omitempty"`
	Nodes            []LoadBalancerNode              `json:"nodes,omitempty"`
	CreatedAt        time.Time                       `json:"created_at,omitempty"`
	UpdatedAt        time.Time                       `json:"updated_at,omitempty"`
	MaintenanceDOW   LoadBalancerMaintenanceDOW      `json:"maintenance_dow,omitempty"`
	MaintenanceTime  string                          `json:"maintenance_time,omitempty"`
}

// LoadBalancerMatcher represents rule matcher
type LoadBalancerMatcher struct {
	Type           LoadBalancerMatcherType                `json:"type,omitempty"`
	SrcIP          *LoadBalancerMatcherSourceIP           `json:"match_src_ip,omitempty"`
	SrcPort        *LoadBalancerMatcherInteger            `json:"match_src_port,omitempty"`
	BodySize       *LoadBalancerMatcherInteger            `json:"match_body_size,omitempty"`
	Path           *LoadBalancerMatcherString             `json:"match_path,omitempty"`
	RequestHeader  *LoadBalancerMatcherStringWithArgument `json:"match_request_header,omitempty"`
	ResponseHeader *LoadBalancerMatcherStringWithArgument `json:"match_response_header,omitempty"`
	URL            *LoadBalancerMatcherString             `json:"match_url,omitempty"`
	URLQuery       *LoadBalancerMatcherString             `json:"match_url_query,omitempty"`
	Host           *LoadBalancerMatcherHost               `json:"match_host,omitempty"`
	HTTPMethod     *LoadBalancerMatcherHTTPMethod         `json:"match_http_method,omitempty"`
	HTTPStatus     *LoadBalancerMatcherInteger            `json:"match_http_status,omitempty"`
	Cookie         *LoadBalancerMatcherStringWithArgument `json:"match_cookie,omitempty"`
	Header         *LoadBalancerMatcherStringWithArgument `json:"match_header,omitempty"` // Deprecated: use RequestHeader instead
	URLParam       *LoadBalancerMatcherStringWithArgument `json:"match_url_param,omitempty"`
	NumMembersUp   *LoadBalancerMatcherNumMembersUp       `json:"match_num_members_up,omitempty"`
	Inverse        *bool                                  `json:"inverse,omitempty"`
}

// LoadBalancerMatcherStringWithArgument represents 'string with argument' matcher
type LoadBalancerMatcherStringWithArgument struct {
	Method     LoadBalancerStringMatcherMethod `json:"method,omitempty"`
	Name       string                          `json:"name,omitempty"`
	Value      string                          `json:"value,omitempty"`
	IgnoreCase *bool                           `json:"ignore_case,omitempty"`
}

// LoadBalancerMatcherHost represents 'host' matcher
type LoadBalancerMatcherHost struct {
	Value string `json:"value,omitempty"`
}

// LoadBalancerMatcherNumMembersUp represents 'num_members_up' matcher
type LoadBalancerMatcherNumMembersUp struct {
	Method  LoadBalancerIntegerMatcherMethod `json:"method,omitempty"`
	Value   int                              `json:"value,omitempty"`
	Backend string                           `json:"backend,omitempty"`
}

// LoadBalancerMatcherHTTPMethod represents 'http_method' matcher
type LoadBalancerMatcherHTTPMethod struct {
	Value LoadBalancerHTTPMatcherMethod `json:"value,omitempty"`
}

// LoadBalancerMatcherInteger represents integer matcher
type LoadBalancerMatcherInteger struct {
	Method     LoadBalancerIntegerMatcherMethod `json:"method,omitempty"`
	Value      int                              `json:"value,omitempty"`
	RangeStart int                              `json:"range_start,omitempty"`
	RangeEnd   int                              `json:"range_end,omitempty"`
}

// LoadBalancerMatcherString represents string matcher
type LoadBalancerMatcherString struct {
	Method     LoadBalancerStringMatcherMethod `json:"method,omitempty"`
	Value      string                          `json:"value,omitempty"`
	IgnoreCase *bool                           `json:"ignore_case,omitempty"`
}

// LoadBalancerMatcherSourceIP represents 'src_ip' matcher
type LoadBalancerMatcherSourceIP struct {
	Value string `json:"value,omitempty"`
}

// LoadBalancerAction represents rule action
type LoadBalancerAction struct {
	Type                LoadBalancerActionType                 `json:"type,omitempty"`
	UseBackend          *LoadBalancerActionUseBackend          `json:"action_use_backend,omitempty"`
	TCPReject           *LoadBalancerActionTCPReject           `json:"action_tcp_reject,omitempty"`
	HTTPReturn          *LoadBalancerActionHTTPReturn          `json:"action_http_return,omitempty"`
	HTTPRedirect        *LoadBalancerActionHTTPRedirect        `json:"action_http_redirect,omitempty"`
	SetForwardedHeaders *LoadBalancerActionSetForwardedHeaders `json:"action_set_forwarded_headers,omitempty"`
	SetRequestHeader    *LoadBalancerActionSetHeader           `json:"action_set_request_header,omitempty"`
	SetResponseHeader   *LoadBalancerActionSetHeader           `json:"action_set_response_header,omitempty"`
}

// LoadBalancerActionUseBackend represents 'use_backend' action
type LoadBalancerActionUseBackend struct {
	Backend string `json:"backend,omitempty"`
}

// LoadBalancerActionTCPReject represents 'tcp_reject' action
type LoadBalancerActionTCPReject struct{}

// LoadBalancerActionHTTPReturn represents 'http_return' action
type LoadBalancerActionHTTPReturn struct {
	Status      int    `json:"status,omitempty"`
	ContentType string `json:"content_type,omitempty"`
	Payload     string `json:"payload,omitempty"`
}

// LoadBalancerActionHTTPRedirect represents 'http_redirect' action. Only either Location or Scheme should be defined.
type LoadBalancerActionHTTPRedirect struct {
	Location string                               `json:"location,omitempty"`
	Scheme   LoadBalancerActionHTTPRedirectScheme `json:"scheme,omitempty"`
	Status   int                                  `json:"status,omitempty"`
}

// LoadBalancerActionSetForwardedHeaders represents 'set_forwarded_headers' action
type LoadBalancerActionSetForwardedHeaders struct{}

// LoadBalancerActionSetHeader represents 'set_request_header' and 'set_response_header' actions
type LoadBalancerActionSetHeader struct {
	Header string `json:"header,omitempty"`
	Value  string `json:"value,omitempty"`
}

// LoadBalancerCertificateBundle represents certificate bundle
type LoadBalancerCertificateBundle struct {
	UUID          string    `json:"uuid,omitempty"`
	Certificate   string    `json:"certificate,omitempty"`
	Intermediates string    `json:"intermediates,omitempty"`
	Hostnames     []string  `json:"hostnames,omitempty"`
	KeyType       string    `json:"key_type,omitempty"`
	Name          string    `json:"name,omitempty"`
	NotAfter      time.Time `json:"not_after,omitempty"`
	NotBefore     time.Time `json:"not_before,omitempty"`
	CreatedAt     time.Time `json:"created_at,omitempty"`
	UpdatedAt     time.Time `json:"updated_at,omitempty"`

	Type             LoadBalancerCertificateBundleType             `json:"type,omitempty"`
	OperationalState LoadBalancerCertificateBundleOperationalState `json:"operational_state,omitempty"`
}

// LoadBalancerNetwork represents network attached to loadbalancer
type LoadBalancerNetwork struct {
	UUID        string                    `json:"uuid,omitempty"`
	Name        string                    `json:"name,omitempty"`
	Type        LoadBalancerNetworkType   `json:"type,omitempty"`
	Family      LoadBalancerAddressFamily `json:"family,omitempty"`
	IPAddresses []LoadBalancerIPAddress   `json:"ip_addresses,omitempty"`
	DNSName     string                    `json:"dns_name,omitempty"`
	CreatedAt   time.Time                 `json:"created_at,omitempty"`
	UpdatedAt   time.Time                 `json:"updated_at,omitempty"`
}

type LoadBalancerFloatingIPAddress struct {
	Address     string    `json:"address,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
	NetworkName string    `json:"network_name,omitempty"`
}

// LoadBalancerIPAddress represents IP address inside loadbalancer service
type LoadBalancerIPAddress struct {
	Address string `json:"address,omitempty"`
	Listen  bool   `json:"listen"`
}

// LoadBalancerNode represents loadbalancer node
type LoadBalancerNode struct {
	Networks         []LoadBalancerNodeNetwork        `json:"networks,omitempty"`
	OperationalState LoadBalancerNodeOperationalState `json:"operational_state,omitempty"`
}

// LoadBalancerNodeNetwork represents node network
type LoadBalancerNodeNetwork struct {
	Name        string                  `json:"name,omitempty"`
	Type        LoadBalancerNetworkType `json:"type,omitempty"`
	IPAddresses []LoadBalancerIPAddress `json:"ip_addresses,omitempty"`
}

// LoadBalancerFrontendNetwork represents network attached to loadbalancer
type LoadBalancerFrontendNetwork struct {
	Name string `json:"name,omitempty"`
}

type LoadBalancerDNSChallengeDomain struct {
	Domain string `json:"domain"`
}
