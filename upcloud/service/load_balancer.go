package service

import (
	"encoding/json"

	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud/request"
)

type LoadBalancer interface {
	GetLoadBalancers(r *request.GetLoadBalancersRequest) ([]upcloud.LoadBalancer, error)
	GetLoadBalancer(r *request.GetLoadBalancerRequest) (*upcloud.LoadBalancer, error)
	CreateLoadBalancer(r *request.CreateLoadBalancerRequest) (*upcloud.LoadBalancer, error)
	ModifyLoadBalancer(r *request.ModifyLoadBalancerRequest) (*upcloud.LoadBalancer, error)
	DeleteLoadBalancer(r *request.DeleteLoadBalancerRequest) error
	// Backends
	GetLoadBalancerBackends(r *request.GetLoadBalancerBackendsRequest) ([]upcloud.LoadBalancerBackend, error)
	GetLoadBalancerBackend(r *request.GetLoadBalancerBackendRequest) (*upcloud.LoadBalancerBackend, error)
	CreateLoadBalancerBackend(r *request.CreateLoadBalancerBackendRequest) (*upcloud.LoadBalancerBackend, error)
	ModifyLoadBalancerBackend(r *request.ModifyLoadBalancerBackendRequest) (*upcloud.LoadBalancerBackend, error)
	DeleteLoadBalancerBackend(r *request.DeleteLoadBalancerBackendRequest) error
	// Backend members
	GetLoadBalancerBackendMembers(r *request.GetLoadBalancerBackendMembersRequest) ([]upcloud.LoadBalancerBackendMember, error)
	GetLoadBalancerBackendMember(r *request.GetLoadBalancerBackendMemberRequest) (*upcloud.LoadBalancerBackendMember, error)
	CreateLoadBalancerBackendMember(r *request.CreateLoadBalancerBackendMemberRequest) (*upcloud.LoadBalancerBackendMember, error)
	ModifyLoadBalancerBackendMember(r *request.ModifyLoadBalancerBackendMemberRequest) (*upcloud.LoadBalancerBackendMember, error)
	DeleteLoadBalancerBackendMember(r *request.DeleteLoadBalancerBackendMemberRequest) error
	// Resolvers
	GetLoadBalancerResolvers(r *request.GetLoadBalancerResolversRequest) ([]upcloud.LoadBalancerResolver, error)
	CreateLoadBalancerResolver(r *request.CreateLoadBalancerResolverRequest) (*upcloud.LoadBalancerResolver, error)
	GetLoadBalancerResolver(r *request.GetLoadBalancerResolverRequest) (*upcloud.LoadBalancerResolver, error)
	ModifyLoadBalancerResolver(r *request.ModifyLoadBalancerRevolverRequest) (*upcloud.LoadBalancerResolver, error)
	DeleteLoadBalancerResolver(r *request.DeleteLoadBalancerResolverRequest) error
	// Plans
	GetLoadBalancerPlans(r *request.GetLoadBalancerPlansRequest) ([]upcloud.LoadBalancerPlan, error)
	// Frontends
	GetLoadBalancerFrontends(r *request.GetLoadBalancerFrontendsRequest) ([]upcloud.LoadBalancerFrontend, error)
	GetLoadBalancerFrontend(r *request.GetLoadBalancerFrontendRequest) (*upcloud.LoadBalancerFrontend, error)
	CreateLoadBalancerFrontend(r *request.CreateLoadBalancerFrontendRequest) (*upcloud.LoadBalancerFrontend, error)
	ModifyLoadBalancerFrontend(r *request.ModifyLoadBalancerFrontendRequest) (*upcloud.LoadBalancerFrontend, error)
	DeleteLoadBalancerFrontend(r *request.DeleteLoadBalancerFrontendRequest) error
	// Frontend rules
	GetLoadBalancerFrontendRules(r *request.GetLoadBalancerFrontendRulesRequest) ([]upcloud.LoadBalancerFrontendRule, error)
	GetLoadBalancerFrontendRule(r *request.GetLoadBalancerFrontendRuleRequest) (*upcloud.LoadBalancerFrontendRule, error)
	CreateLoadBalancerFrontendRule(r *request.CreateLoadBalancerFrontendRuleRequest) (*upcloud.LoadBalancerFrontendRule, error)
	ModifyLoadBalancerFrontendRule(r *request.ModifyLoadBalancerFrontendRuleRequest) (*upcloud.LoadBalancerFrontendRule, error)
	ReplaceLoadBalancerFrontendRule(r *request.ReplaceLoadBalancerFrontendRuleRequest) (*upcloud.LoadBalancerFrontendRule, error)
	DeleteLoadBalancerFrontendRule(r *request.DeleteLoadBalancerFrontendRuleRequest) error
	// TLS Config
	GetLoadBalancerFrontendTLSConfigs(r *request.GetLoadBalancerFrontendTLSConfigsRequest) ([]upcloud.LoadBalancerFrontendTLSConfig, error)
	GetLoadBalancerFrontendTLSConfig(r *request.GetLoadBalancerFrontendTLSConfigRequest) (*upcloud.LoadBalancerFrontendTLSConfig, error)
	CreateLoadBalancerFrontendTLSConfig(r *request.CreateLoadBalancerFrontendTLSConfigRequest) (*upcloud.LoadBalancerFrontendTLSConfig, error)
	ModifyLoadBalancerFrontendTLSConfig(r *request.ModifyLoadBalancerFrontendTLSConfigRequest) (*upcloud.LoadBalancerFrontendTLSConfig, error)
	DeleteLoadBalancerFrontendTLSConfig(r *request.DeleteLoadBalancerFrontendTLSConfigRequest) error
	// Certificate bundles
	GetLoadBalancerCertificateBundles(r *request.GetLoadBalancerCertificateBundlesRequest) ([]upcloud.LoadBalancerCertificateBundle, error)
	GetLoadBalancerCertificateBundle(r *request.GetLoadBalancerCertificateBundleRequest) (*upcloud.LoadBalancerCertificateBundle, error)
	CreateLoadBalancerCertificateBundle(r *request.CreateLoadBalancerCertificateBundleRequest) (*upcloud.LoadBalancerCertificateBundle, error)
	ModifyLoadBalancerCertificateBundle(r *request.ModifyLoadBalancerCertificateBundleRequest) (*upcloud.LoadBalancerCertificateBundle, error)
	DeleteLoadBalancerCertificateBundle(r *request.DeleteLoadBalancerCertificateBundleRequest) error
}

var _ LoadBalancer = (*Service)(nil)

// GetLoadBalancers retrieves a list of load balancers.
func (s *Service) GetLoadBalancers(r *request.GetLoadBalancersRequest) ([]upcloud.LoadBalancer, error) {
	loadBalancers := make([]upcloud.LoadBalancer, 0)
	if r.Page != nil {
		return loadBalancers, s.get(r.RequestURL(), &loadBalancers)
	}

	// copy request value so that we are not altering original request
	req := *r

	// use default page size and get all available records
	req.Page = request.DefaultPage

	// loop until max result is reached or until response doesn't fill our page anymore
	for len(loadBalancers) <= request.PageResultMaxSize {
		lbs := make([]upcloud.LoadBalancer, 0)
		if err := s.get(req.RequestURL(), &lbs); err != nil || len(lbs) < 1 {
			return loadBalancers, err
		}

		loadBalancers = append(loadBalancers, lbs...)
		if len(lbs) < req.Page.SizeInt() {
			return loadBalancers, nil
		}

		req.Page = req.Page.Next()
	}

	return loadBalancers, nil
}

// GetLoadBalancer retrieves details of a load balancer.
func (s *Service) GetLoadBalancer(r *request.GetLoadBalancerRequest) (*upcloud.LoadBalancer, error) {
	loadBalancer := upcloud.LoadBalancer{}
	res, err := s.basicGetRequest(r.RequestURL())
	if err != nil {
		return nil, parseJSONServiceError(err)
	}

	err = json.Unmarshal(res, &loadBalancer)
	if err != nil {
		return nil, err
	}

	return &loadBalancer, nil
}

// CreateLoadBalancer creates a new load balancer.
func (s *Service) CreateLoadBalancer(r *request.CreateLoadBalancerRequest) (*upcloud.LoadBalancer, error) {
	loadBalancer := upcloud.LoadBalancer{}
	requestBody, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	res, err := s.client.PerformJSONPostRequest(s.client.CreateRequestURL(r.RequestURL()), requestBody)
	if err != nil {
		return nil, parseJSONServiceError(err)
	}

	err = json.Unmarshal(res, &loadBalancer)
	if err != nil {
		return nil, err
	}

	return &loadBalancer, nil
}

// ModifyLoadBalancer modifies an existing load balancer.
func (s *Service) ModifyLoadBalancer(r *request.ModifyLoadBalancerRequest) (*upcloud.LoadBalancer, error) {
	loadBalancer := upcloud.LoadBalancer{}
	requestBody, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	res, err := s.client.PerformJSONPatchRequest(s.client.CreateRequestURL(r.RequestURL()), requestBody)
	if err != nil {
		return nil, parseJSONServiceError(err)
	}

	err = json.Unmarshal(res, &loadBalancer)
	if err != nil {
		return nil, err
	}

	return &loadBalancer, nil
}

// DeleteLoadBalancer deletes an existing load balancer.
func (s *Service) DeleteLoadBalancer(r *request.DeleteLoadBalancerRequest) error {
	if err := s.client.PerformJSONDeleteRequest(s.client.CreateRequestURL(r.RequestURL())); err != nil {
		return parseJSONServiceError(err)
	}

	return nil
}

// GetLoadBalancerBackends retrieves a list of load balancer backends.
func (s *Service) GetLoadBalancerBackends(r *request.GetLoadBalancerBackendsRequest) ([]upcloud.LoadBalancerBackend, error) {
	backends := make([]upcloud.LoadBalancerBackend, 0)
	res, err := s.basicGetRequest(r.RequestURL())
	if err != nil {
		return nil, parseJSONServiceError(err)
	}

	err = json.Unmarshal(res, &backends)
	return backends, err
}

// GetLoadBalancerBackend retrieves details of a load balancer backend.
func (s *Service) GetLoadBalancerBackend(r *request.GetLoadBalancerBackendRequest) (*upcloud.LoadBalancerBackend, error) {
	var backend upcloud.LoadBalancerBackend
	res, err := s.basicGetRequest(r.RequestURL())
	if err != nil {
		return nil, parseJSONServiceError(err)
	}

	err = json.Unmarshal(res, &backend)
	return &backend, err
}

// CreateLoadBalancerBackend creates a new load balancer backend.
func (s *Service) CreateLoadBalancerBackend(r *request.CreateLoadBalancerBackendRequest) (*upcloud.LoadBalancerBackend, error) {
	var backend upcloud.LoadBalancerBackend

	reqBody, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	res, err := s.client.PerformJSONPostRequest(s.client.CreateRequestURL(r.RequestURL()), reqBody)
	if err != nil {
		return nil, parseJSONServiceError(err)
	}

	err = json.Unmarshal(res, &backend)
	return &backend, err
}

// ModifyLoadBalancerBackend modifies an existing load balancer backend.
func (s *Service) ModifyLoadBalancerBackend(r *request.ModifyLoadBalancerBackendRequest) (*upcloud.LoadBalancerBackend, error) {
	var backend upcloud.LoadBalancerBackend

	reqBody, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	res, err := s.client.PerformJSONPatchRequest(s.client.CreateRequestURL(r.RequestURL()), reqBody)
	if err != nil {
		return nil, parseJSONServiceError(err)
	}

	err = json.Unmarshal(res, &backend)
	return &backend, err
}

// DeleteLoadBalancerBackend deletes an existing load balancer backend.
func (s *Service) DeleteLoadBalancerBackend(r *request.DeleteLoadBalancerBackendRequest) error {
	if err := s.client.PerformJSONDeleteRequest(s.client.CreateRequestURL(r.RequestURL())); err != nil {
		return parseJSONServiceError(err)
	}
	return nil
}

// GetLoadBalancerBackendMembers retrieves a list of load balancer backend members.
func (s *Service) GetLoadBalancerBackendMembers(r *request.GetLoadBalancerBackendMembersRequest) ([]upcloud.LoadBalancerBackendMember, error) {
	members := make([]upcloud.LoadBalancerBackendMember, 0)
	res, err := s.basicGetRequest(r.RequestURL())
	if err != nil {
		return nil, parseJSONServiceError(err)
	}

	err = json.Unmarshal(res, &members)
	return members, err
}

// GetLoadBalancerBackendMember retrieves details of a load balancer balancer backend member.
func (s *Service) GetLoadBalancerBackendMember(r *request.GetLoadBalancerBackendMemberRequest) (*upcloud.LoadBalancerBackendMember, error) {
	var member upcloud.LoadBalancerBackendMember
	res, err := s.basicGetRequest(r.RequestURL())
	if err != nil {
		return nil, parseJSONServiceError(err)
	}

	err = json.Unmarshal(res, &member)
	return &member, err
}

// CreateLoadBalancerBackendMember creates a new load balancer balancer backend member.
func (s *Service) CreateLoadBalancerBackendMember(r *request.CreateLoadBalancerBackendMemberRequest) (*upcloud.LoadBalancerBackendMember, error) {
	var member upcloud.LoadBalancerBackendMember

	reqBody, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	res, err := s.client.PerformJSONPostRequest(s.client.CreateRequestURL(r.RequestURL()), reqBody)
	if err != nil {
		return nil, parseJSONServiceError(err)
	}

	err = json.Unmarshal(res, &member)
	return &member, err
}

// ModifyLoadBalancerBackendMember modifies an existing load balancer backend member.
func (s *Service) ModifyLoadBalancerBackendMember(r *request.ModifyLoadBalancerBackendMemberRequest) (*upcloud.LoadBalancerBackendMember, error) {
	var member upcloud.LoadBalancerBackendMember

	reqBody, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	res, err := s.client.PerformJSONPatchRequest(s.client.CreateRequestURL(r.RequestURL()), reqBody)
	if err != nil {
		return nil, parseJSONServiceError(err)
	}

	err = json.Unmarshal(res, &member)
	return &member, err
}

// DeleteLoadBalancerBackendMember deletes an existing load balancer backend member.
func (s *Service) DeleteLoadBalancerBackendMember(r *request.DeleteLoadBalancerBackendMemberRequest) error {
	if err := s.client.PerformJSONDeleteRequest(s.client.CreateRequestURL(r.RequestURL())); err != nil {
		return parseJSONServiceError(err)
	}
	return nil
}

// CreateLoadBalancerResolver creates a new load balancer resolver.
func (s *Service) CreateLoadBalancerResolver(r *request.CreateLoadBalancerResolverRequest) (*upcloud.LoadBalancerResolver, error) {
	var resolver upcloud.LoadBalancerResolver

	reqBody, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	res, err := s.client.PerformJSONPostRequest(s.client.CreateRequestURL(r.RequestURL()), reqBody)
	if err != nil {
		return nil, parseJSONServiceError(err)
	}

	err = json.Unmarshal(res, &resolver)
	return &resolver, err
}

// GetLoadBalancerResolvers retrieves a list of load balancer resolvers.
func (s *Service) GetLoadBalancerResolvers(r *request.GetLoadBalancerResolversRequest) ([]upcloud.LoadBalancerResolver, error) {
	resolvers := make([]upcloud.LoadBalancerResolver, 0)

	res, err := s.basicGetRequest(r.RequestURL())
	if err != nil {
		return nil, parseJSONServiceError(err)
	}

	err = json.Unmarshal(res, &resolvers)
	return resolvers, err
}

// GetLoadBalancerResolver retrieves details of a load balancer resolver.
func (s *Service) GetLoadBalancerResolver(r *request.GetLoadBalancerResolverRequest) (*upcloud.LoadBalancerResolver, error) {
	var resolver upcloud.LoadBalancerResolver

	res, err := s.basicGetRequest(r.RequestURL())
	if err != nil {
		return nil, parseJSONServiceError(err)
	}

	err = json.Unmarshal(res, &resolver)
	return &resolver, err
}

// ModifyLoadBalancerResolver modifies an existing load balancer resolver.
func (s *Service) ModifyLoadBalancerResolver(r *request.ModifyLoadBalancerRevolverRequest) (*upcloud.LoadBalancerResolver, error) {
	var resolver upcloud.LoadBalancerResolver

	reqBody, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	res, err := s.client.PerformJSONPatchRequest(s.client.CreateRequestURL(r.RequestURL()), reqBody)
	if err != nil {
		return nil, parseJSONServiceError(err)
	}

	err = json.Unmarshal(res, &resolver)
	return &resolver, err
}

// DeleteLoadBalancerResolver deletes an existing load balancer resolver.
func (s *Service) DeleteLoadBalancerResolver(r *request.DeleteLoadBalancerResolverRequest) error {
	if err := s.client.PerformJSONDeleteRequest(s.client.CreateRequestURL(r.RequestURL())); err != nil {
		return parseJSONServiceError(err)
	}
	return nil
}

// GetLoadBalancerPlans retrieves a list of load balancer plans.
func (s *Service) GetLoadBalancerPlans(r *request.GetLoadBalancerPlansRequest) ([]upcloud.LoadBalancerPlan, error) {
	plans := make([]upcloud.LoadBalancerPlan, 0)
	res, err := s.basicGetRequest(r.RequestURL())
	if err != nil {
		return nil, parseJSONServiceError(err)
	}

	err = json.Unmarshal(res, &plans)
	return plans, err
}

// GetLoadBalancerFrontends retrieves a list of load balancer frontends.
func (s *Service) GetLoadBalancerFrontends(r *request.GetLoadBalancerFrontendsRequest) ([]upcloud.LoadBalancerFrontend, error) {
	fes := make([]upcloud.LoadBalancerFrontend, 0)
	res, err := s.basicGetRequest(r.RequestURL())
	if err != nil {
		return nil, parseJSONServiceError(err)
	}

	err = json.Unmarshal(res, &fes)
	if err != nil {
		return nil, err
	}

	return fes, nil
}

// GetLoadBalancerFrontend retrieves details of a load balancer frontend.
func (s *Service) GetLoadBalancerFrontend(r *request.GetLoadBalancerFrontendRequest) (*upcloud.LoadBalancerFrontend, error) {
	var fe upcloud.LoadBalancerFrontend
	res, err := s.basicGetRequest(r.RequestURL())
	if err != nil {
		return nil, parseJSONServiceError(err)
	}

	err = json.Unmarshal(res, &fe)
	return &fe, err
}

// CreateLoadBalancerFrontend creates a new load balancer frontend.
func (s *Service) CreateLoadBalancerFrontend(r *request.CreateLoadBalancerFrontendRequest) (*upcloud.LoadBalancerFrontend, error) {
	var fe upcloud.LoadBalancerFrontend
	reqBody, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	res, err := s.client.PerformJSONPostRequest(s.client.CreateRequestURL(r.RequestURL()), reqBody)
	if err != nil {
		return nil, parseJSONServiceError(err)
	}

	err = json.Unmarshal(res, &fe)
	return &fe, err
}

// ModifyLoadBalancerFrontend modifies an existing load balancer frontend.
func (s *Service) ModifyLoadBalancerFrontend(r *request.ModifyLoadBalancerFrontendRequest) (*upcloud.LoadBalancerFrontend, error) {
	var fe upcloud.LoadBalancerFrontend

	reqBody, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	res, err := s.client.PerformJSONPatchRequest(s.client.CreateRequestURL(r.RequestURL()), reqBody)
	if err != nil {
		return nil, parseJSONServiceError(err)
	}

	err = json.Unmarshal(res, &fe)
	return &fe, err
}

// DeleteLoadBalancerFrontend deletes an existing load balancer frontend.
func (s *Service) DeleteLoadBalancerFrontend(r *request.DeleteLoadBalancerFrontendRequest) error {
	if err := s.client.PerformJSONDeleteRequest(s.client.CreateRequestURL(r.RequestURL())); err != nil {
		return parseJSONServiceError(err)
	}
	return nil
}

// GetLoadBalancerFrontendRules retrieves a list of load balancer frontend rules.
func (s *Service) GetLoadBalancerFrontendRules(r *request.GetLoadBalancerFrontendRulesRequest) ([]upcloud.LoadBalancerFrontendRule, error) {
	rules := make([]upcloud.LoadBalancerFrontendRule, 0)

	res, err := s.basicGetRequest(r.RequestURL())
	if err != nil {
		return nil, parseJSONServiceError(err)
	}

	err = json.Unmarshal(res, &rules)
	return rules, err
}

// GetLoadBalancerFrontendRule retrieves details of a load balancer frontend rule.
func (s *Service) GetLoadBalancerFrontendRule(r *request.GetLoadBalancerFrontendRuleRequest) (*upcloud.LoadBalancerFrontendRule, error) {
	var rule upcloud.LoadBalancerFrontendRule
	res, err := s.basicGetRequest(r.RequestURL())
	if err != nil {
		return nil, parseJSONServiceError(err)
	}

	err = json.Unmarshal(res, &rule)
	return &rule, err
}

// CreateLoadBalancerFrontendRule creates a new load balancer frontend rule.
func (s *Service) CreateLoadBalancerFrontendRule(r *request.CreateLoadBalancerFrontendRuleRequest) (*upcloud.LoadBalancerFrontendRule, error) {
	var rule upcloud.LoadBalancerFrontendRule

	reqBody, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	res, err := s.client.PerformJSONPostRequest(s.client.CreateRequestURL(r.RequestURL()), reqBody)
	if err != nil {
		return nil, parseJSONServiceError(err)
	}

	err = json.Unmarshal(res, &rule)
	return &rule, err
}

// ModifyLoadBalancerFrontendRule modifies an existing load balancer frontend rule.
func (s *Service) ModifyLoadBalancerFrontendRule(r *request.ModifyLoadBalancerFrontendRuleRequest) (*upcloud.LoadBalancerFrontendRule, error) {
	var rule upcloud.LoadBalancerFrontendRule

	reqBody, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	res, err := s.client.PerformJSONPatchRequest(s.client.CreateRequestURL(r.RequestURL()), reqBody)
	if err != nil {
		return nil, parseJSONServiceError(err)
	}

	err = json.Unmarshal(res, &rule)
	return &rule, err
}

// ReplaceLoadBalancerFrontendRule replaces an existing load balancer frontend rule.
func (s *Service) ReplaceLoadBalancerFrontendRule(r *request.ReplaceLoadBalancerFrontendRuleRequest) (*upcloud.LoadBalancerFrontendRule, error) {
	var rule upcloud.LoadBalancerFrontendRule

	reqBody, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	res, err := s.client.PerformJSONPutRequest(s.client.CreateRequestURL(r.RequestURL()), reqBody)
	if err != nil {
		return nil, parseJSONServiceError(err)
	}

	err = json.Unmarshal(res, &rule)
	return &rule, err
}

// DeleteLoadBalancerFrontendRule deletes an existing load balancer frontend rule.
func (s *Service) DeleteLoadBalancerFrontendRule(r *request.DeleteLoadBalancerFrontendRuleRequest) error {
	if err := s.client.PerformJSONDeleteRequest(s.client.CreateRequestURL(r.RequestURL())); err != nil {
		return parseJSONServiceError(err)
	}
	return nil
}

// GetLoadBalancerFrontendTLSConfigs retrieves a list of load balancer frontend TLS configs.
func (s *Service) GetLoadBalancerFrontendTLSConfigs(r *request.GetLoadBalancerFrontendTLSConfigsRequest) ([]upcloud.LoadBalancerFrontendTLSConfig, error) {
	configs := make([]upcloud.LoadBalancerFrontendTLSConfig, 0)

	res, err := s.basicGetRequest(r.RequestURL())
	if err != nil {
		return nil, parseJSONServiceError(err)
	}

	err = json.Unmarshal(res, &configs)
	return configs, err
}

// GetLoadBalancerFrontendTLSConfig retrieves details of a load balancer frontend TLS config.
func (s *Service) GetLoadBalancerFrontendTLSConfig(r *request.GetLoadBalancerFrontendTLSConfigRequest) (*upcloud.LoadBalancerFrontendTLSConfig, error) {
	var config upcloud.LoadBalancerFrontendTLSConfig
	res, err := s.basicGetRequest(r.RequestURL())
	if err != nil {
		return nil, parseJSONServiceError(err)
	}

	err = json.Unmarshal(res, &config)
	return &config, err
}

// CreateLoadBalancerFrontendTLSConfig creates a new load balancer frontend TLS config.
func (s *Service) CreateLoadBalancerFrontendTLSConfig(r *request.CreateLoadBalancerFrontendTLSConfigRequest) (*upcloud.LoadBalancerFrontendTLSConfig, error) {
	var config upcloud.LoadBalancerFrontendTLSConfig

	reqBody, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	res, err := s.client.PerformJSONPostRequest(s.client.CreateRequestURL(r.RequestURL()), reqBody)
	if err != nil {
		return nil, parseJSONServiceError(err)
	}

	err = json.Unmarshal(res, &config)
	return &config, err
}

// ModifyLoadBalancerFrontendTLSConfig modifies an existing load balancer frontend TLS Config.
func (s *Service) ModifyLoadBalancerFrontendTLSConfig(r *request.ModifyLoadBalancerFrontendTLSConfigRequest) (*upcloud.LoadBalancerFrontendTLSConfig, error) {
	var config upcloud.LoadBalancerFrontendTLSConfig

	reqBody, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	res, err := s.client.PerformJSONPatchRequest(s.client.CreateRequestURL(r.RequestURL()), reqBody)
	if err != nil {
		return nil, parseJSONServiceError(err)
	}

	err = json.Unmarshal(res, &config)
	return &config, err
}

// DeleteLoadBalancerFrontendTLSConfig deletes an existing load balancer frontend TLS config.
func (s *Service) DeleteLoadBalancerFrontendTLSConfig(r *request.DeleteLoadBalancerFrontendTLSConfigRequest) error {
	if err := s.client.PerformJSONDeleteRequest(s.client.CreateRequestURL(r.RequestURL())); err != nil {
		return parseJSONServiceError(err)
	}
	return nil
}

// GetLoadBalancerCertificateBundles retrieves details of a load balancer certificate bundles.
func (s *Service) GetLoadBalancerCertificateBundles(r *request.GetLoadBalancerCertificateBundlesRequest) ([]upcloud.LoadBalancerCertificateBundle, error) {
	certs := make([]upcloud.LoadBalancerCertificateBundle, 0)
	res, err := s.basicGetRequest(r.RequestURL())
	if err != nil {
		return nil, parseJSONServiceError(err)
	}

	err = json.Unmarshal(res, &certs)
	if err != nil {
		return nil, err
	}

	return certs, nil
}

// GetLoadBalancerCertificateBundle retrieves details of a load balancer certificate bundle.
func (s *Service) GetLoadBalancerCertificateBundle(r *request.GetLoadBalancerCertificateBundleRequest) (*upcloud.LoadBalancerCertificateBundle, error) {
	var cert upcloud.LoadBalancerCertificateBundle
	res, err := s.basicGetRequest(r.RequestURL())
	if err != nil {
		return nil, parseJSONServiceError(err)
	}

	err = json.Unmarshal(res, &cert)
	return &cert, err
}

// CreateLoadBalancerCertificateBundle creates a new load balancer certificate bundle.
func (s *Service) CreateLoadBalancerCertificateBundle(r *request.CreateLoadBalancerCertificateBundleRequest) (*upcloud.LoadBalancerCertificateBundle, error) {
	var cert upcloud.LoadBalancerCertificateBundle
	reqBody, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	res, err := s.client.PerformJSONPostRequest(s.client.CreateRequestURL(r.RequestURL()), reqBody)
	if err != nil {
		return nil, parseJSONServiceError(err)
	}

	err = json.Unmarshal(res, &cert)
	return &cert, err
}

// ModifyLoadBalancerCertificateBundle modifies an existing load balancer certificate bundle.
func (s *Service) ModifyLoadBalancerCertificateBundle(r *request.ModifyLoadBalancerCertificateBundleRequest) (*upcloud.LoadBalancerCertificateBundle, error) {
	var cert upcloud.LoadBalancerCertificateBundle
	reqBody, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	res, err := s.client.PerformJSONPatchRequest(s.client.CreateRequestURL(r.RequestURL()), reqBody)
	if err != nil {
		return nil, parseJSONServiceError(err)
	}

	err = json.Unmarshal(res, &cert)
	return &cert, err
}

// DeleteLoadBalancerCertificateBundle deletes an existing load balancer certificate bundle.
func (s *Service) DeleteLoadBalancerCertificateBundle(r *request.DeleteLoadBalancerCertificateBundleRequest) error {
	if err := s.client.PerformJSONDeleteRequest(s.client.CreateRequestURL(r.RequestURL())); err != nil {
		return parseJSONServiceError(err)
	}
	return nil
}
