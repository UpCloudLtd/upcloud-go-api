package request

import (
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
)

func NewLoadBalancerTCPRejectAction() upcloud.LoadBalancerAction {
	return upcloud.LoadBalancerAction{
		Type:      upcloud.LoadBalancerActionTypeTCPReject,
		TCPReject: &upcloud.LoadBalancerActionTCPReject{},
	}
}

func NewLoadBalancerHTTPReturnAction(statusCode int, contentType, payload string) upcloud.LoadBalancerAction {
	return upcloud.LoadBalancerAction{
		Type: upcloud.LoadBalancerActionTypeHTTPReturn,
		HTTPReturn: &upcloud.LoadBalancerActionHTTPReturn{
			Status:      statusCode,
			ContentType: contentType,
			Payload:     payload,
		},
	}
}

func NewLoadBalancerHTTPRedirectAction(location string) upcloud.LoadBalancerAction {
	return upcloud.LoadBalancerAction{
		Type: upcloud.LoadBalancerActionTypeHTTPRedirect,
		HTTPRedirect: &upcloud.LoadBalancerActionHTTPRedirect{
			Location: location,
		},
	}
}

func NewLoadBalancerHTTPRedirectSchemeAction(scheme upcloud.LoadBalancerActionHTTPRedirectScheme) upcloud.LoadBalancerAction {
	return upcloud.LoadBalancerAction{
		Type: upcloud.LoadBalancerActionTypeHTTPRedirect,
		HTTPRedirect: &upcloud.LoadBalancerActionHTTPRedirect{
			Scheme: scheme,
		},
	}
}

func NewLoadBalancerUseBackendAction(name string) upcloud.LoadBalancerAction {
	return upcloud.LoadBalancerAction{
		Type: upcloud.LoadBalancerActionTypeUseBackend,
		UseBackend: &upcloud.LoadBalancerActionUseBackend{
			Backend: name,
		},
	}
}

func NewLoadBalancerSetForwardedHeadersAction() upcloud.LoadBalancerAction {
	return upcloud.LoadBalancerAction{
		Type:                upcloud.LoadBalancerActionTypeSetForwardedHeaders,
		SetForwardedHeaders: &upcloud.LoadBalancerActionSetForwardedHeaders{},
	}
}

func NewLoadBalancerSetRequestHeaderAction(header, value string) upcloud.LoadBalancerAction {
	return upcloud.LoadBalancerAction{
		Type: upcloud.LoadBalancerActionTypeSetRequestHeader,
		SetRequestHeader: &upcloud.LoadBalancerActionSetHeader{
			Header: header,
			Value:  value,
		},
	}
}

func NewLoadBalancerSetResponseHeaderAction(header, value string) upcloud.LoadBalancerAction {
	return upcloud.LoadBalancerAction{
		Type: upcloud.LoadBalancerActionTypeSetResponseHeader,
		SetResponseHeader: &upcloud.LoadBalancerActionSetHeader{
			Header: header,
			Value:  value,
		},
	}
}

func NewLoadBalancerNumMembersUpMatcher(m upcloud.LoadBalancerIntegerMatcherMethod, count int, backend string) upcloud.LoadBalancerMatcher {
	return upcloud.LoadBalancerMatcher{
		Type: upcloud.LoadBalancerMatcherTypeNumMembersUp,
		NumMembersUp: &upcloud.LoadBalancerMatcherNumMembersUp{
			Method:  m,
			Value:   count,
			Backend: backend,
		},
	}
}

func NewLoadBalancerURLParamMatcher(m upcloud.LoadBalancerStringMatcherMethod, name, value string, ignoreCase *bool) upcloud.LoadBalancerMatcher {
	return upcloud.LoadBalancerMatcher{
		Type: upcloud.LoadBalancerMatcherTypeURLParam,
		URLParam: &upcloud.LoadBalancerMatcherStringWithArgument{
			Method:     m,
			Name:       name,
			Value:      value,
			IgnoreCase: ignoreCase,
		},
	}
}

// Deprecated: Use NewLoadBalancerRequestHeaderMatcher instead
func NewLoadBalancerHeaderMatcher(m upcloud.LoadBalancerStringMatcherMethod, name, value string, ignoreCase *bool) upcloud.LoadBalancerMatcher {
	return upcloud.LoadBalancerMatcher{
		Type: upcloud.LoadBalancerMatcherTypeHeader,
		Header: &upcloud.LoadBalancerMatcherStringWithArgument{
			Method:     m,
			Name:       name,
			Value:      value,
			IgnoreCase: ignoreCase,
		},
	}
}

func NewLoadBalancerRequestHeaderMatcher(m upcloud.LoadBalancerStringMatcherMethod, name, value string, ignoreCase *bool) upcloud.LoadBalancerMatcher {
	return upcloud.LoadBalancerMatcher{
		Type: upcloud.LoadBalancerMatcherTypeRequestHeader,
		RequestHeader: &upcloud.LoadBalancerMatcherStringWithArgument{
			Method:     m,
			Name:       name,
			Value:      value,
			IgnoreCase: ignoreCase,
		},
	}
}

func NewLoadBalancerResponseHeaderMatcher(m upcloud.LoadBalancerStringMatcherMethod, name, value string, ignoreCase *bool) upcloud.LoadBalancerMatcher {
	return upcloud.LoadBalancerMatcher{
		Type: upcloud.LoadBalancerMatcherTypeResponseHeader,
		ResponseHeader: &upcloud.LoadBalancerMatcherStringWithArgument{
			Method:     m,
			Name:       name,
			Value:      value,
			IgnoreCase: ignoreCase,
		},
	}
}

func NewLoadBalancerCookieMatcher(m upcloud.LoadBalancerStringMatcherMethod, name, value string, ignoreCase *bool) upcloud.LoadBalancerMatcher {
	return upcloud.LoadBalancerMatcher{
		Type: upcloud.LoadBalancerMatcherTypeCookie,
		Cookie: &upcloud.LoadBalancerMatcherStringWithArgument{
			Method:     m,
			Name:       name,
			Value:      value,
			IgnoreCase: ignoreCase,
		},
	}
}

func NewLoadBalancerHTTPMethodMatcher(method upcloud.LoadBalancerHTTPMatcherMethod) upcloud.LoadBalancerMatcher {
	return upcloud.LoadBalancerMatcher{
		Type: upcloud.LoadBalancerMatcherTypeHTTPMethod,
		HTTPMethod: &upcloud.LoadBalancerMatcherHTTPMethod{
			Value: method,
		},
	}
}

func NewLoadBalancerHTTPStatusMatcher(m upcloud.LoadBalancerIntegerMatcherMethod, status int) upcloud.LoadBalancerMatcher {
	return upcloud.LoadBalancerMatcher{
		Type: upcloud.LoadBalancerMatcherTypeHTTPStatus,
		HTTPStatus: &upcloud.LoadBalancerMatcherInteger{
			Method: m,
			Value:  status,
		},
	}
}

func NewLoadBalancerHTTPStatusRangeMatcher(start, end int) upcloud.LoadBalancerMatcher {
	return upcloud.LoadBalancerMatcher{
		Type: upcloud.LoadBalancerMatcherTypeHTTPStatus,
		HTTPStatus: &upcloud.LoadBalancerMatcherInteger{
			Method:     upcloud.LoadBalancerIntegerMatcherMethodRange,
			RangeStart: start,
			RangeEnd:   end,
		},
	}
}

func NewLoadBalancerHostMatcher(host string) upcloud.LoadBalancerMatcher {
	return upcloud.LoadBalancerMatcher{
		Type: upcloud.LoadBalancerMatcherTypeHost,
		Host: &upcloud.LoadBalancerMatcherHost{
			Value: host,
		},
	}
}

func NewLoadBalancerURLQueryMatcher(m upcloud.LoadBalancerStringMatcherMethod, query string, ignoreCase *bool) upcloud.LoadBalancerMatcher {
	return upcloud.LoadBalancerMatcher{
		Type: upcloud.LoadBalancerMatcherTypeURLQuery,
		URLQuery: &upcloud.LoadBalancerMatcherString{
			Method:     m,
			Value:      query,
			IgnoreCase: ignoreCase,
		},
	}
}

func NewLoadBalancerURLMatcher(m upcloud.LoadBalancerStringMatcherMethod, URL string, ignoreCase *bool) upcloud.LoadBalancerMatcher {
	return upcloud.LoadBalancerMatcher{
		Type: upcloud.LoadBalancerMatcherTypeURL,
		URL: &upcloud.LoadBalancerMatcherString{
			Method:     m,
			Value:      URL,
			IgnoreCase: ignoreCase,
		},
	}
}

func NewLoadBalancerPathMatcher(m upcloud.LoadBalancerStringMatcherMethod, path string, ignoreCase *bool) upcloud.LoadBalancerMatcher {
	return upcloud.LoadBalancerMatcher{
		Type: upcloud.LoadBalancerMatcherTypePath,
		Path: &upcloud.LoadBalancerMatcherString{
			Method:     m,
			Value:      path,
			IgnoreCase: ignoreCase,
		},
	}
}

func NewLoadBalancerBodySizeRangeMatcher(start, end int) upcloud.LoadBalancerMatcher {
	return upcloud.LoadBalancerMatcher{
		Type: upcloud.LoadBalancerMatcherTypeBodySize,
		BodySize: &upcloud.LoadBalancerMatcherInteger{
			Method:     upcloud.LoadBalancerIntegerMatcherMethodRange,
			RangeStart: start,
			RangeEnd:   end,
		},
	}
}

func NewLoadBalancerBodySizeMatcher(m upcloud.LoadBalancerIntegerMatcherMethod, bodySize int) upcloud.LoadBalancerMatcher {
	return upcloud.LoadBalancerMatcher{
		Type: upcloud.LoadBalancerMatcherTypeBodySize,
		BodySize: &upcloud.LoadBalancerMatcherInteger{
			Method: m,
			Value:  bodySize,
		},
	}
}

func NewLoadBalancerSrcPortRangeMatcher(start, end int) upcloud.LoadBalancerMatcher {
	return upcloud.LoadBalancerMatcher{
		Type: upcloud.LoadBalancerMatcherTypeSrcPort,
		SrcPort: &upcloud.LoadBalancerMatcherInteger{
			Method:     upcloud.LoadBalancerIntegerMatcherMethodRange,
			RangeStart: start,
			RangeEnd:   end,
		},
	}
}

func NewLoadBalancerSrcPortMatcher(m upcloud.LoadBalancerIntegerMatcherMethod, port int) upcloud.LoadBalancerMatcher {
	return upcloud.LoadBalancerMatcher{
		Type: upcloud.LoadBalancerMatcherTypeSrcPort,
		SrcPort: &upcloud.LoadBalancerMatcherInteger{
			Method: m,
			Value:  port,
		},
	}
}

func NewLoadBalancerSrcIPMatcher(IP string) upcloud.LoadBalancerMatcher {
	return upcloud.LoadBalancerMatcher{
		Type: upcloud.LoadBalancerMatcherTypeSrcIP,
		SrcIP: &upcloud.LoadBalancerMatcherSourceIP{
			Value: IP,
		},
	}
}

// NewLoadBalancerInverseMatcher helper converts matcher to inverse matcher.
//
// Usage: inverseIPMatch := NewLoadBalancerInverseMatcher(NewLoadBalancerSrcIPMatcher("127.0.0.2"))
func NewLoadBalancerInverseMatcher(m upcloud.LoadBalancerMatcher) upcloud.LoadBalancerMatcher {
	m.Inverse = upcloud.BoolPtr(true)
	return m
}

func newLoadBalancerBackendMember(t upcloud.LoadBalancerBackendMemberType, name string, weight int, maxSessions int, enabled bool, IP string, port int) LoadBalancerBackendMember {
	return LoadBalancerBackendMember{
		Type:        t,
		Name:        name,
		Weight:      weight,
		MaxSessions: maxSessions,
		Enabled:     enabled,
		IP:          IP,
		Port:        port,
	}
}

func NewLoadBalancerDynamicBackendMember(name string, weight int, maxSessions int, enabled bool, IP string, port int) LoadBalancerBackendMember {
	return newLoadBalancerBackendMember(upcloud.LoadBalancerBackendMemberTypeDynamic, name, weight, maxSessions, enabled, IP, port)
}

func NewLoadBalancerStaticBackendMember(name string, weight int, maxSessions int, enabled bool, IP string, port int) LoadBalancerBackendMember {
	return newLoadBalancerBackendMember(upcloud.LoadBalancerBackendMemberTypeStatic, name, weight, maxSessions, enabled, IP, port)
}

func NewCreateLoadBalancerManualCertificateBundleRequest(name, certificate, intermediates, privateKey string) CreateLoadBalancerCertificateBundleRequest {
	return CreateLoadBalancerCertificateBundleRequest{
		Type:          upcloud.LoadBalancerCertificateBundleTypeManual,
		Name:          name,
		Certificate:   certificate,
		Intermediates: intermediates,
		PrivateKey:    privateKey,
	}
}

func NewCreateLoadBalancerDynamicCertificateBundleRequest(name, keyType string, hostnames []string) CreateLoadBalancerCertificateBundleRequest {
	return CreateLoadBalancerCertificateBundleRequest{
		Type:      upcloud.LoadBalancerCertificateBundleTypeDynamic,
		Name:      name,
		KeyType:   keyType,
		Hostnames: hostnames,
	}
}
