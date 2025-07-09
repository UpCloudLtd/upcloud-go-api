package service

import (
	"context"
	"errors"
	"net/http"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
)

type LoadBalancer interface {
	GetLoadBalancers(ctx context.Context, r *request.GetLoadBalancersRequest) ([]upcloud.LoadBalancer, error)
	GetLoadBalancer(ctx context.Context, r *request.GetLoadBalancerRequest) (*upcloud.LoadBalancer, error)
	CreateLoadBalancer(ctx context.Context, r *request.CreateLoadBalancerRequest) (*upcloud.LoadBalancer, error)
	ModifyLoadBalancer(ctx context.Context, r *request.ModifyLoadBalancerRequest) (*upcloud.LoadBalancer, error)
	DeleteLoadBalancer(ctx context.Context, r *request.DeleteLoadBalancerRequest) error
	WaitForLoadBalancerOperationalState(ctx context.Context, r *request.WaitForLoadBalancerOperationalStateRequest) (*upcloud.LoadBalancer, error)
	WaitForLoadBalancerDeletion(ctx context.Context, r *request.WaitForLoadBalancerDeletionRequest) error
	// Backends
	GetLoadBalancerBackends(ctx context.Context, r *request.GetLoadBalancerBackendsRequest) ([]upcloud.LoadBalancerBackend, error)
	GetLoadBalancerBackend(ctx context.Context, r *request.GetLoadBalancerBackendRequest) (*upcloud.LoadBalancerBackend, error)
	CreateLoadBalancerBackend(ctx context.Context, r *request.CreateLoadBalancerBackendRequest) (*upcloud.LoadBalancerBackend, error)
	ModifyLoadBalancerBackend(ctx context.Context, r *request.ModifyLoadBalancerBackendRequest) (*upcloud.LoadBalancerBackend, error)
	DeleteLoadBalancerBackend(ctx context.Context, r *request.DeleteLoadBalancerBackendRequest) error
	// Backend members
	GetLoadBalancerBackendMembers(ctx context.Context, r *request.GetLoadBalancerBackendMembersRequest) ([]upcloud.LoadBalancerBackendMember, error)
	GetLoadBalancerBackendMember(ctx context.Context, r *request.GetLoadBalancerBackendMemberRequest) (*upcloud.LoadBalancerBackendMember, error)
	CreateLoadBalancerBackendMember(ctx context.Context, r *request.CreateLoadBalancerBackendMemberRequest) (*upcloud.LoadBalancerBackendMember, error)
	ModifyLoadBalancerBackendMember(ctx context.Context, r *request.ModifyLoadBalancerBackendMemberRequest) (*upcloud.LoadBalancerBackendMember, error)
	DeleteLoadBalancerBackendMember(ctx context.Context, r *request.DeleteLoadBalancerBackendMemberRequest) error
	// Backend TLS Config
	GetLoadBalancerBackendTLSConfigs(ctx context.Context, r *request.GetLoadBalancerBackendTLSConfigsRequest) ([]upcloud.LoadBalancerBackendTLSConfig, error)
	GetLoadBalancerBackendTLSConfig(ctx context.Context, r *request.GetLoadBalancerBackendTLSConfigRequest) (*upcloud.LoadBalancerBackendTLSConfig, error)
	CreateLoadBalancerBackendTLSConfig(ctx context.Context, r *request.CreateLoadBalancerBackendTLSConfigRequest) (*upcloud.LoadBalancerBackendTLSConfig, error)
	ModifyLoadBalancerBackendTLSConfig(ctx context.Context, r *request.ModifyLoadBalancerBackendTLSConfigRequest) (*upcloud.LoadBalancerBackendTLSConfig, error)
	DeleteLoadBalancerBackendTLSConfig(ctx context.Context, r *request.DeleteLoadBalancerBackendTLSConfigRequest) error
	// Resolvers
	GetLoadBalancerResolvers(ctx context.Context, r *request.GetLoadBalancerResolversRequest) ([]upcloud.LoadBalancerResolver, error)
	CreateLoadBalancerResolver(ctx context.Context, r *request.CreateLoadBalancerResolverRequest) (*upcloud.LoadBalancerResolver, error)
	GetLoadBalancerResolver(ctx context.Context, r *request.GetLoadBalancerResolverRequest) (*upcloud.LoadBalancerResolver, error)
	ModifyLoadBalancerResolver(ctx context.Context, r *request.ModifyLoadBalancerResolverRequest) (*upcloud.LoadBalancerResolver, error)
	DeleteLoadBalancerResolver(ctx context.Context, r *request.DeleteLoadBalancerResolverRequest) error
	// Plans
	GetLoadBalancerPlans(ctx context.Context, r *request.GetLoadBalancerPlansRequest) ([]upcloud.LoadBalancerPlan, error)
	// Frontends
	GetLoadBalancerFrontends(ctx context.Context, r *request.GetLoadBalancerFrontendsRequest) ([]upcloud.LoadBalancerFrontend, error)
	GetLoadBalancerFrontend(ctx context.Context, r *request.GetLoadBalancerFrontendRequest) (*upcloud.LoadBalancerFrontend, error)
	CreateLoadBalancerFrontend(ctx context.Context, r *request.CreateLoadBalancerFrontendRequest) (*upcloud.LoadBalancerFrontend, error)
	ModifyLoadBalancerFrontend(ctx context.Context, r *request.ModifyLoadBalancerFrontendRequest) (*upcloud.LoadBalancerFrontend, error)
	DeleteLoadBalancerFrontend(ctx context.Context, r *request.DeleteLoadBalancerFrontendRequest) error
	// Frontend rules
	GetLoadBalancerFrontendRules(ctx context.Context, r *request.GetLoadBalancerFrontendRulesRequest) ([]upcloud.LoadBalancerFrontendRule, error)
	GetLoadBalancerFrontendRule(ctx context.Context, r *request.GetLoadBalancerFrontendRuleRequest) (*upcloud.LoadBalancerFrontendRule, error)
	CreateLoadBalancerFrontendRule(ctx context.Context, r *request.CreateLoadBalancerFrontendRuleRequest) (*upcloud.LoadBalancerFrontendRule, error)
	ModifyLoadBalancerFrontendRule(ctx context.Context, r *request.ModifyLoadBalancerFrontendRuleRequest) (*upcloud.LoadBalancerFrontendRule, error)
	ReplaceLoadBalancerFrontendRule(ctx context.Context, r *request.ReplaceLoadBalancerFrontendRuleRequest) (*upcloud.LoadBalancerFrontendRule, error)
	DeleteLoadBalancerFrontendRule(ctx context.Context, r *request.DeleteLoadBalancerFrontendRuleRequest) error
	// Frontend TLS Config
	GetLoadBalancerFrontendTLSConfigs(ctx context.Context, r *request.GetLoadBalancerFrontendTLSConfigsRequest) ([]upcloud.LoadBalancerFrontendTLSConfig, error)
	GetLoadBalancerFrontendTLSConfig(ctx context.Context, r *request.GetLoadBalancerFrontendTLSConfigRequest) (*upcloud.LoadBalancerFrontendTLSConfig, error)
	CreateLoadBalancerFrontendTLSConfig(ctx context.Context, r *request.CreateLoadBalancerFrontendTLSConfigRequest) (*upcloud.LoadBalancerFrontendTLSConfig, error)
	ModifyLoadBalancerFrontendTLSConfig(ctx context.Context, r *request.ModifyLoadBalancerFrontendTLSConfigRequest) (*upcloud.LoadBalancerFrontendTLSConfig, error)
	DeleteLoadBalancerFrontendTLSConfig(ctx context.Context, r *request.DeleteLoadBalancerFrontendTLSConfigRequest) error
	// Certificate bundles
	GetLoadBalancerCertificateBundles(ctx context.Context, r *request.GetLoadBalancerCertificateBundlesRequest) ([]upcloud.LoadBalancerCertificateBundle, error)
	GetLoadBalancerCertificateBundle(ctx context.Context, r *request.GetLoadBalancerCertificateBundleRequest) (*upcloud.LoadBalancerCertificateBundle, error)
	CreateLoadBalancerCertificateBundle(ctx context.Context, r *request.CreateLoadBalancerCertificateBundleRequest) (*upcloud.LoadBalancerCertificateBundle, error)
	ModifyLoadBalancerCertificateBundle(ctx context.Context, r *request.ModifyLoadBalancerCertificateBundleRequest) (*upcloud.LoadBalancerCertificateBundle, error)
	DeleteLoadBalancerCertificateBundle(ctx context.Context, r *request.DeleteLoadBalancerCertificateBundleRequest) error
	// Networks
	ModifyLoadBalancerNetwork(ctx context.Context, r *request.ModifyLoadBalancerNetworkRequest) (*upcloud.LoadBalancerNetwork, error)
	GetLoadBalancerDNSChallengeDomain(ctx context.Context, r *request.GetLoadBalancerDNSChallengeDomainRequest) (*upcloud.LoadBalancerDNSChallengeDomain, error)
}

// GetLoadBalancers retrieves a list of load balancers.
func (s *Service) GetLoadBalancers(ctx context.Context, r *request.GetLoadBalancersRequest) ([]upcloud.LoadBalancer, error) {
	loadBalancers := make([]upcloud.LoadBalancer, 0)
	if r.Page != nil {
		return loadBalancers, s.get(ctx, r.RequestURL(), &loadBalancers)
	}

	// copy request value so that we are not altering original request
	req := *r

	// use default page size and get all available records
	req.Page = request.DefaultPage

	// loop until max result is reached or until response doesn't fill our page anymore
	for len(loadBalancers) <= request.PageResultMaxSize {
		lbs := make([]upcloud.LoadBalancer, 0)
		if err := s.get(ctx, req.RequestURL(), &lbs); err != nil || len(lbs) < 1 {
			return loadBalancers, err
		}

		loadBalancers = append(loadBalancers, lbs...)
		if len(lbs) < req.Page.Size {
			return loadBalancers, nil
		}

		req.Page = req.Page.Next()
	}

	return loadBalancers, nil
}

// GetLoadBalancer retrieves details of a load balancer.
func (s *Service) GetLoadBalancer(ctx context.Context, r *request.GetLoadBalancerRequest) (*upcloud.LoadBalancer, error) {
	loadBalancer := upcloud.LoadBalancer{}
	return &loadBalancer, s.get(ctx, r.RequestURL(), &loadBalancer)
}

// CreateLoadBalancer creates a new load balancer.
func (s *Service) CreateLoadBalancer(ctx context.Context, r *request.CreateLoadBalancerRequest) (*upcloud.LoadBalancer, error) {
	loadBalancer := upcloud.LoadBalancer{}
	return &loadBalancer, s.create(ctx, r, &loadBalancer)
}

// ModifyLoadBalancer modifies an existing load balancer.
func (s *Service) ModifyLoadBalancer(ctx context.Context, r *request.ModifyLoadBalancerRequest) (*upcloud.LoadBalancer, error) {
	loadBalancer := upcloud.LoadBalancer{}
	return &loadBalancer, s.modify(ctx, r, &loadBalancer)
}

// DeleteLoadBalancer deletes an existing load balancer.
func (s *Service) DeleteLoadBalancer(ctx context.Context, r *request.DeleteLoadBalancerRequest) error {
	return s.delete(ctx, r)
}

// WaitForLoadBalancerOperationalState blocks execution until the specified load balancer instance has entered the
// specified state. If the state changes favorably, the new load balancer details is returned. The method will give up
// after the specified timeout
func (s *Service) WaitForLoadBalancerOperationalState(ctx context.Context, r *request.WaitForLoadBalancerOperationalStateRequest) (*upcloud.LoadBalancer, error) {
	return retry(ctx, func(_ int, c context.Context) (*upcloud.LoadBalancer, error) {
		details, err := s.GetLoadBalancer(c, &request.GetLoadBalancerRequest{
			UUID: r.UUID,
		})
		if err != nil {
			return nil, err
		}

		if details.OperationalState == r.DesiredState {
			return details, nil
		}
		return nil, nil
	}, nil)
}

// WaitForLoadBalancerDeletion blocks execution until the specified load balancer instance has been deleted.
// nil error is returned when the load balancer has been deleted, otherwise an error is returned.
func (s *Service) WaitForLoadBalancerDeletion(ctx context.Context, r *request.WaitForLoadBalancerDeletionRequest) error {

	_, err := retry(ctx, func(_ int, c context.Context) (*upcloud.LoadBalancer, error) {
		details, err := s.GetLoadBalancer(c, &request.GetLoadBalancerRequest{
			UUID: r.UUID,
		})
		if err != nil {
			var ucErr *upcloud.Problem
			if errors.As(err, &ucErr) && ucErr.Status == http.StatusNotFound {
				return nil, nil
			}

			return nil, err
		}

		return details, err

	}, &retryConfig{inverse: true})
	return err
}

// GetLoadBalancerBackends retrieves a list of load balancer backends.
func (s *Service) GetLoadBalancerBackends(ctx context.Context, r *request.GetLoadBalancerBackendsRequest) ([]upcloud.LoadBalancerBackend, error) {
	backends := make([]upcloud.LoadBalancerBackend, 0)
	return backends, s.get(ctx, r.RequestURL(), &backends)
}

// GetLoadBalancerBackend retrieves details of a load balancer backend.
func (s *Service) GetLoadBalancerBackend(ctx context.Context, r *request.GetLoadBalancerBackendRequest) (*upcloud.LoadBalancerBackend, error) {
	var backend upcloud.LoadBalancerBackend
	return &backend, s.get(ctx, r.RequestURL(), &backend)
}

// CreateLoadBalancerBackend creates a new load balancer backend.
func (s *Service) CreateLoadBalancerBackend(ctx context.Context, r *request.CreateLoadBalancerBackendRequest) (*upcloud.LoadBalancerBackend, error) {
	var backend upcloud.LoadBalancerBackend
	return &backend, s.create(ctx, r, &backend)
}

// ModifyLoadBalancerBackend modifies an existing load balancer backend.
func (s *Service) ModifyLoadBalancerBackend(ctx context.Context, r *request.ModifyLoadBalancerBackendRequest) (*upcloud.LoadBalancerBackend, error) {
	var backend upcloud.LoadBalancerBackend
	return &backend, s.modify(ctx, r, &backend)
}

// DeleteLoadBalancerBackend deletes an existing load balancer backend.
func (s *Service) DeleteLoadBalancerBackend(ctx context.Context, r *request.DeleteLoadBalancerBackendRequest) error {
	return s.delete(ctx, r)
}

// GetLoadBalancerBackendTLSConfigs retrieves a list of load balancer backend TLS configs.
func (s *Service) GetLoadBalancerBackendTLSConfigs(ctx context.Context, r *request.GetLoadBalancerBackendTLSConfigsRequest) ([]upcloud.LoadBalancerBackendTLSConfig, error) {
	configs := make([]upcloud.LoadBalancerBackendTLSConfig, 0)
	return configs, s.get(ctx, r.RequestURL(), &configs)
}

// GetLoadBalancerBackendTLSConfig retrieves details of a load balancer backend TLS config.
func (s *Service) GetLoadBalancerBackendTLSConfig(ctx context.Context, r *request.GetLoadBalancerBackendTLSConfigRequest) (*upcloud.LoadBalancerBackendTLSConfig, error) {
	var config upcloud.LoadBalancerBackendTLSConfig
	return &config, s.get(ctx, r.RequestURL(), &config)
}

// CreateLoadBalancerBackendTLSConfig creates a new load balancer backend TLS config.
func (s *Service) CreateLoadBalancerBackendTLSConfig(ctx context.Context, r *request.CreateLoadBalancerBackendTLSConfigRequest) (*upcloud.LoadBalancerBackendTLSConfig, error) {
	var config upcloud.LoadBalancerBackendTLSConfig
	return &config, s.create(ctx, r, &config)
}

// ModifyLoadBalancerBackendTLSConfig modifies an existing load balancer backend TLS Config.
func (s *Service) ModifyLoadBalancerBackendTLSConfig(ctx context.Context, r *request.ModifyLoadBalancerBackendTLSConfigRequest) (*upcloud.LoadBalancerBackendTLSConfig, error) {
	var config upcloud.LoadBalancerBackendTLSConfig
	return &config, s.modify(ctx, r, &config)
}

// DeleteLoadBalancerBackendTLSConfig deletes an existing load balancer backend TLS config.
func (s *Service) DeleteLoadBalancerBackendTLSConfig(ctx context.Context, r *request.DeleteLoadBalancerBackendTLSConfigRequest) error {
	return s.delete(ctx, r)
}

// GetLoadBalancerBackendMembers retrieves a list of load balancer backend members.
func (s *Service) GetLoadBalancerBackendMembers(ctx context.Context, r *request.GetLoadBalancerBackendMembersRequest) ([]upcloud.LoadBalancerBackendMember, error) {
	members := make([]upcloud.LoadBalancerBackendMember, 0)
	return members, s.get(ctx, r.RequestURL(), &members)
}

// GetLoadBalancerBackendMember retrieves details of a load balancer backend member.
func (s *Service) GetLoadBalancerBackendMember(ctx context.Context, r *request.GetLoadBalancerBackendMemberRequest) (*upcloud.LoadBalancerBackendMember, error) {
	var member upcloud.LoadBalancerBackendMember
	return &member, s.get(ctx, r.RequestURL(), &member)
}

// CreateLoadBalancerBackendMember creates a new load balancer backend member.
func (s *Service) CreateLoadBalancerBackendMember(ctx context.Context, r *request.CreateLoadBalancerBackendMemberRequest) (*upcloud.LoadBalancerBackendMember, error) {
	var member upcloud.LoadBalancerBackendMember
	return &member, s.create(ctx, r, &member)
}

// ModifyLoadBalancerBackendMember modifies an existing load balancer backend member.
func (s *Service) ModifyLoadBalancerBackendMember(ctx context.Context, r *request.ModifyLoadBalancerBackendMemberRequest) (*upcloud.LoadBalancerBackendMember, error) {
	var member upcloud.LoadBalancerBackendMember
	return &member, s.modify(ctx, r, &member)
}

// DeleteLoadBalancerBackendMember deletes an existing load balancer backend member.
func (s *Service) DeleteLoadBalancerBackendMember(ctx context.Context, r *request.DeleteLoadBalancerBackendMemberRequest) error {
	return s.delete(ctx, r)
}

// CreateLoadBalancerResolver creates a new load balancer resolver.
func (s *Service) CreateLoadBalancerResolver(ctx context.Context, r *request.CreateLoadBalancerResolverRequest) (*upcloud.LoadBalancerResolver, error) {
	var resolver upcloud.LoadBalancerResolver
	return &resolver, s.create(ctx, r, &resolver)
}

// GetLoadBalancerResolvers retrieves a list of load balancer resolvers.
func (s *Service) GetLoadBalancerResolvers(ctx context.Context, r *request.GetLoadBalancerResolversRequest) ([]upcloud.LoadBalancerResolver, error) {
	resolvers := make([]upcloud.LoadBalancerResolver, 0)
	return resolvers, s.get(ctx, r.RequestURL(), &resolvers)
}

// GetLoadBalancerResolver retrieves details of a load balancer resolver.
func (s *Service) GetLoadBalancerResolver(ctx context.Context, r *request.GetLoadBalancerResolverRequest) (*upcloud.LoadBalancerResolver, error) {
	var resolver upcloud.LoadBalancerResolver
	return &resolver, s.get(ctx, r.RequestURL(), &resolver)
}

// ModifyLoadBalancerResolver modifies an existing load balancer resolver.
func (s *Service) ModifyLoadBalancerResolver(ctx context.Context, r *request.ModifyLoadBalancerResolverRequest) (*upcloud.LoadBalancerResolver, error) {
	var resolver upcloud.LoadBalancerResolver
	return &resolver, s.modify(ctx, r, &resolver)
}

// DeleteLoadBalancerResolver deletes an existing load balancer resolver.
func (s *Service) DeleteLoadBalancerResolver(ctx context.Context, r *request.DeleteLoadBalancerResolverRequest) error {
	return s.delete(ctx, r)
}

// GetLoadBalancerPlans retrieves a list of load balancer plans.
func (s *Service) GetLoadBalancerPlans(ctx context.Context, r *request.GetLoadBalancerPlansRequest) ([]upcloud.LoadBalancerPlan, error) {
	plans := make([]upcloud.LoadBalancerPlan, 0)

	if r.Page != nil {
		return plans, s.get(ctx, r.RequestURL(), &plans)
	}

	// copy request value so that we are not altering original request
	req := *r

	// use default page size and get all available records
	req.Page = request.DefaultPage

	// loop until max result is reached or until response doesn't fill our page anymore
	for len(plans) <= request.PageResultMaxSize {
		p := make([]upcloud.LoadBalancerPlan, 0)
		if err := s.get(ctx, req.RequestURL(), &p); err != nil || len(p) < 1 {
			return plans, err
		}

		plans = append(plans, p...)
		if len(p) < req.Page.Size {
			return plans, nil
		}

		req.Page = req.Page.Next()
	}

	return plans, nil
}

// GetLoadBalancerFrontends retrieves a list of load balancer frontends.
func (s *Service) GetLoadBalancerFrontends(ctx context.Context, r *request.GetLoadBalancerFrontendsRequest) ([]upcloud.LoadBalancerFrontend, error) {
	fes := make([]upcloud.LoadBalancerFrontend, 0)
	return fes, s.get(ctx, r.RequestURL(), &fes)
}

// GetLoadBalancerFrontend retrieves details of a load balancer frontend.
func (s *Service) GetLoadBalancerFrontend(ctx context.Context, r *request.GetLoadBalancerFrontendRequest) (*upcloud.LoadBalancerFrontend, error) {
	var fe upcloud.LoadBalancerFrontend
	return &fe, s.get(ctx, r.RequestURL(), &fe)
}

// CreateLoadBalancerFrontend creates a new load balancer frontend.
func (s *Service) CreateLoadBalancerFrontend(ctx context.Context, r *request.CreateLoadBalancerFrontendRequest) (*upcloud.LoadBalancerFrontend, error) {
	var fe upcloud.LoadBalancerFrontend
	return &fe, s.create(ctx, r, &fe)
}

// ModifyLoadBalancerFrontend modifies an existing load balancer frontend.
func (s *Service) ModifyLoadBalancerFrontend(ctx context.Context, r *request.ModifyLoadBalancerFrontendRequest) (*upcloud.LoadBalancerFrontend, error) {
	var fe upcloud.LoadBalancerFrontend
	return &fe, s.modify(ctx, r, &fe)
}

// DeleteLoadBalancerFrontend deletes an existing load balancer frontend.
func (s *Service) DeleteLoadBalancerFrontend(ctx context.Context, r *request.DeleteLoadBalancerFrontendRequest) error {
	return s.delete(ctx, r)
}

// GetLoadBalancerFrontendRules retrieves a list of load balancer frontend rules.
func (s *Service) GetLoadBalancerFrontendRules(ctx context.Context, r *request.GetLoadBalancerFrontendRulesRequest) ([]upcloud.LoadBalancerFrontendRule, error) {
	rules := make([]upcloud.LoadBalancerFrontendRule, 0)
	return rules, s.get(ctx, r.RequestURL(), &rules)
}

// GetLoadBalancerFrontendRule retrieves details of a load balancer frontend rule.
func (s *Service) GetLoadBalancerFrontendRule(ctx context.Context, r *request.GetLoadBalancerFrontendRuleRequest) (*upcloud.LoadBalancerFrontendRule, error) {
	var rule upcloud.LoadBalancerFrontendRule
	return &rule, s.get(ctx, r.RequestURL(), &rule)
}

// CreateLoadBalancerFrontendRule creates a new load balancer frontend rule.
func (s *Service) CreateLoadBalancerFrontendRule(ctx context.Context, r *request.CreateLoadBalancerFrontendRuleRequest) (*upcloud.LoadBalancerFrontendRule, error) {
	var rule upcloud.LoadBalancerFrontendRule
	return &rule, s.create(ctx, r, &rule)
}

// ModifyLoadBalancerFrontendRule modifies an existing load balancer frontend rule.
func (s *Service) ModifyLoadBalancerFrontendRule(ctx context.Context, r *request.ModifyLoadBalancerFrontendRuleRequest) (*upcloud.LoadBalancerFrontendRule, error) {
	var rule upcloud.LoadBalancerFrontendRule
	return &rule, s.modify(ctx, r, &rule)
}

// ReplaceLoadBalancerFrontendRule replaces an existing load balancer frontend rule.
func (s *Service) ReplaceLoadBalancerFrontendRule(ctx context.Context, r *request.ReplaceLoadBalancerFrontendRuleRequest) (*upcloud.LoadBalancerFrontendRule, error) {
	var rule upcloud.LoadBalancerFrontendRule
	return &rule, s.replace(ctx, r, &rule)
}

// DeleteLoadBalancerFrontendRule deletes an existing load balancer frontend rule.
func (s *Service) DeleteLoadBalancerFrontendRule(ctx context.Context, r *request.DeleteLoadBalancerFrontendRuleRequest) error {
	return s.delete(ctx, r)
}

// GetLoadBalancerFrontendTLSConfigs retrieves a list of load balancer frontend TLS configs.
func (s *Service) GetLoadBalancerFrontendTLSConfigs(ctx context.Context, r *request.GetLoadBalancerFrontendTLSConfigsRequest) ([]upcloud.LoadBalancerFrontendTLSConfig, error) {
	configs := make([]upcloud.LoadBalancerFrontendTLSConfig, 0)
	return configs, s.get(ctx, r.RequestURL(), &configs)
}

// GetLoadBalancerFrontendTLSConfig retrieves details of a load balancer frontend TLS config.
func (s *Service) GetLoadBalancerFrontendTLSConfig(ctx context.Context, r *request.GetLoadBalancerFrontendTLSConfigRequest) (*upcloud.LoadBalancerFrontendTLSConfig, error) {
	var config upcloud.LoadBalancerFrontendTLSConfig
	return &config, s.get(ctx, r.RequestURL(), &config)
}

// CreateLoadBalancerFrontendTLSConfig creates a new load balancer frontend TLS config.
func (s *Service) CreateLoadBalancerFrontendTLSConfig(ctx context.Context, r *request.CreateLoadBalancerFrontendTLSConfigRequest) (*upcloud.LoadBalancerFrontendTLSConfig, error) {
	var config upcloud.LoadBalancerFrontendTLSConfig
	return &config, s.create(ctx, r, &config)
}

// ModifyLoadBalancerFrontendTLSConfig modifies an existing load balancer frontend TLS Config.
func (s *Service) ModifyLoadBalancerFrontendTLSConfig(ctx context.Context, r *request.ModifyLoadBalancerFrontendTLSConfigRequest) (*upcloud.LoadBalancerFrontendTLSConfig, error) {
	var config upcloud.LoadBalancerFrontendTLSConfig
	return &config, s.modify(ctx, r, &config)
}

// DeleteLoadBalancerFrontendTLSConfig deletes an existing load balancer frontend TLS config.
func (s *Service) DeleteLoadBalancerFrontendTLSConfig(ctx context.Context, r *request.DeleteLoadBalancerFrontendTLSConfigRequest) error {
	return s.delete(ctx, r)
}

// GetLoadBalancerCertificateBundles retrieves details of a load balancer certificate bundles.
func (s *Service) GetLoadBalancerCertificateBundles(ctx context.Context, r *request.GetLoadBalancerCertificateBundlesRequest) ([]upcloud.LoadBalancerCertificateBundle, error) {
	certs := make([]upcloud.LoadBalancerCertificateBundle, 0)

	if r.Page != nil {
		return certs, s.get(ctx, r.RequestURL(), &certs)
	}

	// copy request value so that we are not altering original request
	req := *r

	// use default page size and get all available records
	req.Page = request.DefaultPage

	// loop until max result is reached or until response doesn't fill our page anymore
	for len(certs) <= request.PageResultMaxSize {
		c := make([]upcloud.LoadBalancerCertificateBundle, 0)
		if err := s.get(ctx, req.RequestURL(), &c); err != nil || len(c) < 1 {
			return certs, err
		}

		certs = append(certs, c...)
		if len(c) < req.Page.Size {
			return certs, nil
		}

		req.Page = req.Page.Next()
	}

	return certs, nil
}

// GetLoadBalancerCertificateBundle retrieves details of a load balancer certificate bundle.
func (s *Service) GetLoadBalancerCertificateBundle(ctx context.Context, r *request.GetLoadBalancerCertificateBundleRequest) (*upcloud.LoadBalancerCertificateBundle, error) {
	var cert upcloud.LoadBalancerCertificateBundle
	return &cert, s.get(ctx, r.RequestURL(), &cert)
}

// CreateLoadBalancerCertificateBundle creates a new load balancer certificate bundle.
func (s *Service) CreateLoadBalancerCertificateBundle(ctx context.Context, r *request.CreateLoadBalancerCertificateBundleRequest) (*upcloud.LoadBalancerCertificateBundle, error) {
	var cert upcloud.LoadBalancerCertificateBundle
	return &cert, s.create(ctx, r, &cert)
}

// ModifyLoadBalancerCertificateBundle modifies an existing load balancer certificate bundle.
func (s *Service) ModifyLoadBalancerCertificateBundle(ctx context.Context, r *request.ModifyLoadBalancerCertificateBundleRequest) (*upcloud.LoadBalancerCertificateBundle, error) {
	var cert upcloud.LoadBalancerCertificateBundle
	return &cert, s.modify(ctx, r, &cert)
}

// DeleteLoadBalancerCertificateBundle deletes an existing load balancer certificate bundle.
func (s *Service) DeleteLoadBalancerCertificateBundle(ctx context.Context, r *request.DeleteLoadBalancerCertificateBundleRequest) error {
	return s.delete(ctx, r)
}

func (s *Service) ModifyLoadBalancerNetwork(ctx context.Context, r *request.ModifyLoadBalancerNetworkRequest) (*upcloud.LoadBalancerNetwork, error) {
	n := upcloud.LoadBalancerNetwork{}
	return &n, s.modify(ctx, r, &n)
}

func (s *Service) GetLoadBalancerDNSChallengeDomain(ctx context.Context, r *request.GetLoadBalancerDNSChallengeDomainRequest) (*upcloud.LoadBalancerDNSChallengeDomain, error) {
	var domain upcloud.LoadBalancerDNSChallengeDomain
	return &domain, s.get(ctx, r.RequestURL(), &domain)
}
