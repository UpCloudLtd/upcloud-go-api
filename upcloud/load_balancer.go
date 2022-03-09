package upcloud

import (
	"time"
)

type LoadBalancerMode string
type LoadBalancerMatcherType string
type LoadBalancerActionType string
type LoadBalancerStringMatcherMethod string
type LoadBalancerHTTPMatcherMethod string
type LoadBalancerIntegerMatcherMethod string
type LoadBalancerBackendMemberType string
type LoadBalancerOperationalState string
type LoadBalancerConfiguredStatus string

const (
	LoadBalancerModeHTTP LoadBalancerMode = "http"
	LoadBalancerModeTCP  LoadBalancerMode = "tcp"

	LoadBalancerBackendMemberTypeStatic  LoadBalancerBackendMemberType = "static"
	LoadBalancerBackendMemberTypeDynamic LoadBalancerBackendMemberType = "dynamic"

	LoadBalancerConfiguredStatusStarted LoadBalancerConfiguredStatus = "started"
	LoadBalancerConfiguredStatusStopped LoadBalancerConfiguredStatus = "stopped"

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

	LoadBalancerMatcherTypeSrcIP        LoadBalancerMatcherType = "src_ip "
	LoadBalancerMatcherTypeSrcPort      LoadBalancerMatcherType = "src_port"
	LoadBalancerMatcherTypeBodySize     LoadBalancerMatcherType = "body_size"
	LoadBalancerMatcherTypePath         LoadBalancerMatcherType = "path"
	LoadBalancerMatcherTypeURL          LoadBalancerMatcherType = "url"
	LoadBalancerMatcherTypeURLQuery     LoadBalancerMatcherType = "url_query"
	LoadBalancerMatcherTypeHost         LoadBalancerMatcherType = "host"
	LoadBalancerMatcherTypeHTTPMethod   LoadBalancerMatcherType = "http_method"
	LoadBalancerMatcherTypeCookie       LoadBalancerMatcherType = "cookie"
	LoadBalancerMatcherTypeHeader       LoadBalancerMatcherType = "header"
	LoadBalancerMatcherTypeURLParam     LoadBalancerMatcherType = "url_param"
	LoadBalancerMatcherTypeNumMembersUP LoadBalancerMatcherType = "num_members_up"

	LoadBalancerActionTypeUseBackend   LoadBalancerActionType = "use_backend"
	LoadBalancerActionTypeTCPReject    LoadBalancerActionType = "tcp_reject"
	LoadBalancerActionTypeHTTPReturn   LoadBalancerActionType = "http_return"
	LoadBalancerActionTypeHTTPRedirect LoadBalancerActionType = "http_redirect"

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

	LoadBalancerIntegerMatcherMethodEqual               LoadBalancerIntegerMatcherMethod = "equal"
	LoadBalancerIntegerMatcherMethodEqualGreaterOrEqual LoadBalancerIntegerMatcherMethod = "greater_or_equal"
	LoadBalancerIntegerMatcherMethodEqualGreater        LoadBalancerIntegerMatcherMethod = "greater"
	LoadBalancerIntegerMatcherMethodEqualLessOrEqual    LoadBalancerIntegerMatcherMethod = "less_or_equal"
	LoadBalancerIntegerMatcherMethodEqualLess           LoadBalancerIntegerMatcherMethod = "less"
)

// LoadBalancerPlan represents load balancer plan details
type LoadBalancerPlan struct {
	Name                 string `json:"name,omitempty"`
	PerServerMaxSessions int    `json:"per_server_max_sessions,omitempty"`
	ServerNumber         int    `json:"server_number,omitempty"`
}

type LoadBalancerFrontend struct {
	Name           string             `json:"name,omitempty"`
	Mode           LoadBalancerMode   `json:"mode,omitempty"`
	Port           int                `json:"port,omitempty"`
	DefaultBackend string             `json:"default_backend,omitempty"`
	Rules          []LoadBalancerRule `json:"rules,omitempty"`
	TLSConfigs     []TLSConfig        `json:"tls_configs,omitempty"`
	CreatedAt      time.Time          `json:"created_at,omitempty"`
	UpdatedAt      time.Time          `json:"updated_at,omitempty"`
}

type LoadBalancerRule struct {
	Name      string                `json:"name,omitempty"`
	Priority  int                   `json:"priority,omitempty"`
	Matchers  []LoadBalancerMatcher `json:"matchers,omitempty"`
	Actions   []LoadBalancerAction  `json:"actions,omitempty"`
	CreatedAt time.Time             `json:"created_at,omitempty"`
	UpdatedAt time.Time             `json:"updated_at,omitempty"`
}

type TLSConfig struct {
	Name                  string    `json:"name,omitempty"`
	CertificateBundleUUID string    `json:"certificate_bundle_uuid,omitempty"`
	CreatedAt             time.Time `json:"created_at,omitempty"`
	UpdatedAt             time.Time `json:"updated_at,omitempty"`
}

type LoadBalancerBackend struct {
	Name      string                      `json:"name"`
	Members   []LoadBalancerBackendMember `json:"members"`
	Resolver  string                      `json:"resolver,omitempty"`
	CreatedAt time.Time                   `json:"created_at,omitempty"`
	UpdatedAt time.Time                   `json:"updated_at,omitempty"`
}

type LoadBalancerBackendMember struct {
	Name        string                        `json:"name"`
	Ip          string                        `json:"ip"`
	Port        int                           `json:"port"`
	Weight      int                           `json:"weight"`
	MaxSessions int                           `json:"max_sessions"`
	Type        LoadBalancerBackendMemberType `json:"type"`
	Enabled     bool                          `json:"enabled"`
	CreatedAt   time.Time                     `json:"created_at,omitempty"`
	UpdatedAt   time.Time                     `json:"updated_at,omitempty"`
}

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

type LoadBalancer struct {
	UUID             string                       `json:"uuid,omitempty"`
	Name             string                       `json:"name,omitempty"`
	Zone             string                       `json:"zone,omitempty"`
	Plan             string                       `json:"plan,omitempty"`
	NetworkUUID      string                       `json:"network_uuid,omitempty"`
	DNSName          string                       `json:"dns_name,omitempty"`
	ConfiguredStatus LoadBalancerConfiguredStatus `json:"configured_status,omitempty"`
	OperationalState LoadBalancerOperationalState `json:"operational_state,omitempty"`
	Frontends        []LoadBalancerFrontend       `json:"frontends,omitempty"`
	Backends         []LoadBalancerBackend        `json:"backends,omitempty"`
	Resolvers        []LoadBalancerResolver       `json:"resolvers,omitempty"`
	CreatedAt        time.Time                    `json:"created_at,omitempty"`
	UpdatedAt        time.Time                    `json:"updated_at,omitempty"`
}

type LoadBalancerMatcher struct {
	Type         LoadBalancerMatcherType                `json:"type,omitempty"`
	SrcIP        *LoadBalancerMatcherSourceIP           `json:"match_src_ip,omitempty"`
	SrcPort      *LoadBalancerMatcherInteger            `json:"match_src_port,omitempty"`
	BodySize     *LoadBalancerMatcherInteger            `json:"match_body_size,omitempty"`
	Path         *LoadBalancerMatcherString             `json:"match_path,omitempty"`
	URL          *LoadBalancerMatcherString             `json:"match_url,omitempty"`
	URLQuery     *LoadBalancerMatcherString             `json:"match_url_query,omitempty"`
	Host         *LoadBalancerMatcherHost               `json:"match_host,omitempty"`
	HTTPMethod   *LoadBalancerMatcherHTTPMethod         `json:"match_http_method,omitempty"`
	Cookie       *LoadBalancerMatcherStringWithArgument `json:"match_cookie,omitempty"`
	Header       *LoadBalancerMatcherStringWithArgument `json:"match_header,omitempty"`
	URLParam     *LoadBalancerMatcherStringWithArgument `json:"match_url_param,omitempty"`
	NumMembersUP *LoadBalancerMatcherBackend            `json:"match_num_members_up,omitempty"`
}

type LoadBalancerMatcherStringWithArgument struct {
	Method     LoadBalancerStringMatcherMethod `json:"method,omitempty"`
	Name       string                          `json:"name,omitempty"`
	Value      string                          `json:"value,omitempty"`
	IgnoreCase *bool                           `json:"ignore_case,omitempty"`
}

type LoadBalancerMatcherHost struct {
	Value string `json:"value,omitempty"`
}

type LoadBalancerMatcherBackend struct {
	Method  LoadBalancerIntegerMatcherMethod `json:"method,omitempty"`
	Value   int                              `json:"value,omitempty"`
	Backend string
}

type LoadBalancerMatcherHTTPMethod struct {
	Value LoadBalancerHTTPMatcherMethod `json:"value,omitempty"`
}

type LoadBalancerMatcherInteger struct {
	Method LoadBalancerIntegerMatcherMethod `json:"method,omitempty"`
	Value  int                              `json:"value,omitempty"`
}

type LoadBalancerMatcherString struct {
	Method     LoadBalancerStringMatcherMethod `json:"method,omitempty"`
	Value      string                          `json:"value,omitempty"`
	IgnoreCase *bool                           `json:"ignore_case,omitempty"`
}

type LoadBalancerMatcherSourceIP struct {
	Value string `json:"value,omitempty"`
}

type LoadBalancerAction struct {
	Type         LoadBalancerActionType          `json:"type,omitempty"`
	UseBackend   *LoadBalancerActionUseBackend   `json:"action_use_backend,omitempty"`
	TCPReject    *LoadBalancerActionTCPReject    `json:"action_tcp_reject,omitempty"`
	HTTPReturn   *LoadBalancerActionHTTPReturn   `json:"action_http_return,omitempty"`
	HTTPRedirect *LoadBalancerActionHTTPRedirect `json:"action_http_redirect,omitempty"`
}

type LoadBalancerActionUseBackend struct {
	Backend string `json:"backend,omitempty"`
}

type LoadBalancerActionTCPReject struct {
}

type LoadBalancerActionHTTPReturn struct {
	Status      int    `json:"status,omitempty"`
	ContentType string `json:"content_type,omitempty"`
	Payload     string `json:"payload,omitempty"`
}

type LoadBalancerActionHTTPRedirect struct {
	Location string `json:"location,omitempty"`
}
