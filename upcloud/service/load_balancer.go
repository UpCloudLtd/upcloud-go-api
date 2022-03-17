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
}

var _ LoadBalancer = (*Service)(nil)

// GetLoadBalancers retrieves a list of load balancers.
func (s *Service) GetLoadBalancers(r *request.GetLoadBalancersRequest) ([]upcloud.LoadBalancer, error) {
	loadBalancers := make([]upcloud.LoadBalancer, 0)
	res, err := s.basicGetRequest(r.RequestURL())
	if err != nil {
		return nil, parseJSONServiceError(err)
	}

	err = json.Unmarshal(res, &loadBalancers)
	if err != nil {
		return nil, err
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

func (s *Service) GetLoadBalancerBackends(r *request.GetLoadBalancerBackendsRequest) ([]upcloud.LoadBalancerBackend, error) {
	backends := make([]upcloud.LoadBalancerBackend, 0)
	res, err := s.basicGetRequest(r.RequestURL())
	if err != nil {
		return nil, parseJSONServiceError(err)
	}

	err = json.Unmarshal(res, &backends)
	return backends, err
}

func (s *Service) GetLoadBalancerBackend(r *request.GetLoadBalancerBackendRequest) (*upcloud.LoadBalancerBackend, error) {
	var backend upcloud.LoadBalancerBackend
	res, err := s.basicGetRequest(r.RequestURL())
	if err != nil {
		return nil, parseJSONServiceError(err)
	}

	err = json.Unmarshal(res, &backend)
	return &backend, err
}

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

func (s *Service) DeleteLoadBalancerBackend(r *request.DeleteLoadBalancerBackendRequest) error {
	if err := s.client.PerformJSONDeleteRequest(s.client.CreateRequestURL(r.RequestURL())); err != nil {
		return parseJSONServiceError(err)
	}
	return nil
}

func (s *Service) GetLoadBalancerBackendMembers(r *request.GetLoadBalancerBackendMembersRequest) ([]upcloud.LoadBalancerBackendMember, error) {
	members := make([]upcloud.LoadBalancerBackendMember, 0)
	res, err := s.basicGetRequest(r.RequestURL())
	if err != nil {
		return nil, parseJSONServiceError(err)
	}

	err = json.Unmarshal(res, &members)
	return members, err
}

func (s *Service) GetLoadBalancerBackendMember(r *request.GetLoadBalancerBackendMemberRequest) (*upcloud.LoadBalancerBackendMember, error) {
	var member upcloud.LoadBalancerBackendMember
	res, err := s.basicGetRequest(r.RequestURL())
	if err != nil {
		return nil, parseJSONServiceError(err)
	}

	err = json.Unmarshal(res, &member)
	return &member, err
}

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

func (s *Service) DeleteLoadBalancerBackendMember(r *request.DeleteLoadBalancerBackendMemberRequest) error {
	if err := s.client.PerformJSONDeleteRequest(s.client.CreateRequestURL(r.RequestURL())); err != nil {
		return parseJSONServiceError(err)
	}
	return nil
}

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

func (s *Service) GetLoadBalancerResolvers(r *request.GetLoadBalancerResolversRequest) ([]upcloud.LoadBalancerResolver, error) {
	resolvers := make([]upcloud.LoadBalancerResolver, 0)

	res, err := s.basicGetRequest(r.RequestURL())
	if err != nil {
		return nil, parseJSONServiceError(err)
	}

	err = json.Unmarshal(res, &resolvers)
	return resolvers, err
}

func (s *Service) GetLoadBalancerResolver(r *request.GetLoadBalancerResolverRequest) (*upcloud.LoadBalancerResolver, error) {
	var resolver upcloud.LoadBalancerResolver

	res, err := s.basicGetRequest(r.RequestURL())
	if err != nil {
		return nil, parseJSONServiceError(err)
	}

	err = json.Unmarshal(res, &resolver)
	return &resolver, err
}

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

func (s *Service) DeleteLoadBalancerResolver(r *request.DeleteLoadBalancerResolverRequest) error {
	if err := s.client.PerformJSONDeleteRequest(s.client.CreateRequestURL(r.RequestURL())); err != nil {
		return parseJSONServiceError(err)
	}
	return nil
}

func (s *Service) GetLoadBalancerPlans(r *request.GetLoadBalancerPlansRequest) ([]upcloud.LoadBalancerPlan, error) {
	plans := make([]upcloud.LoadBalancerPlan, 0)
	res, err := s.basicGetRequest(r.RequestURL())
	if err != nil {
		return nil, parseJSONServiceError(err)
	}

	err = json.Unmarshal(res, &plans)
	return plans, err
}

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

func (s *Service) GetLoadBalancerFrontend(r *request.GetLoadBalancerFrontendRequest) (*upcloud.LoadBalancerFrontend, error) {
	var fe upcloud.LoadBalancerFrontend
	res, err := s.basicGetRequest(r.RequestURL())
	if err != nil {
		return nil, parseJSONServiceError(err)
	}

	err = json.Unmarshal(res, &fe)
	return &fe, err
}

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

func (s *Service) DeleteLoadBalancerFrontend(r *request.DeleteLoadBalancerFrontendRequest) error {
	if err := s.client.PerformJSONDeleteRequest(s.client.CreateRequestURL(r.RequestURL())); err != nil {
		return parseJSONServiceError(err)
	}
	return nil
}

func (s *Service) GetLoadBalancerFrontendRules(r *request.GetLoadBalancerFrontendRulesRequest) ([]upcloud.LoadBalancerFrontendRule, error) {
	rules := make([]upcloud.LoadBalancerFrontendRule, 0)

	res, err := s.basicGetRequest(r.RequestURL())
	if err != nil {
		return nil, parseJSONServiceError(err)
	}

	err = json.Unmarshal(res, &rules)
	return rules, err
}

func (s *Service) GetLoadBalancerFrontendRule(r *request.GetLoadBalancerFrontendRuleRequest) (*upcloud.LoadBalancerFrontendRule, error) {
	var rule upcloud.LoadBalancerFrontendRule
	res, err := s.basicGetRequest(r.RequestURL())
	if err != nil {
		return nil, parseJSONServiceError(err)
	}

	err = json.Unmarshal(res, &rule)
	return &rule, err
}

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

func (s *Service) DeleteLoadBalancerFrontendRule(r *request.DeleteLoadBalancerFrontendRuleRequest) error {
	if err := s.client.PerformJSONDeleteRequest(s.client.CreateRequestURL(r.RequestURL())); err != nil {
		return parseJSONServiceError(err)
	}
	return nil
}
