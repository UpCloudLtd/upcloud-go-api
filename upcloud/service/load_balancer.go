package service

import (
	"encoding/json"

	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud/request"
)

type LoadBalancer interface {
	GetLoadBalancers(r *request.GetLoadBalancersRequest) ([]*upcloud.LoadBalancer, error)
	GetLoadBalancer(r *request.GetLoadBalancerRequest) (*upcloud.LoadBalancer, error)
	CreateLoadBalancer(r *request.CreateLoadBalancerRequest) (*upcloud.LoadBalancer, error)
	ModifyLoadBalancer(r *request.ModifyLoadBalancerRequest) (*upcloud.LoadBalancer, error)
	DeleteLoadBalancer(r *request.DeleteLoadBalancerRequest) error
	GetLoadBalancerBackends(r *request.GetLoadBalancerBackendsRequest) ([]upcloud.LoadBalancerBackend, error)
	GetLoadBalancerBackend(r *request.GetLoadBalancerBackendRequest) (*upcloud.LoadBalancerBackend, error)
	CreateLoadBalancerBackend(r *request.CreateLoadBalancerBackendRequest) (*upcloud.LoadBalancerBackend, error)
	ModifyLoadBalancerBackend(r *request.ModifyLoadBalancerBackendRequest) (*upcloud.LoadBalancerBackend, error)
	DeleteLoadBalancerBackend(r *request.DeleteLoadBalancerBackendRequest) error
	GetLoadBalancerBackendMembers(r *request.GetLoadBalancerBackendMembersRequest) ([]upcloud.LoadBalancerBackendMember, error)
	GetLoadBalancerBackendMember(r *request.GetLoadBalancerBackendMemberRequest) (*upcloud.LoadBalancerBackendMember, error)
	CreateLoadBalancerBackendMember(r *request.CreateLoadBalancerBackendMemberRequest) (*upcloud.LoadBalancerBackendMember, error)
	ModifyLoadBalancerBackendMember(r *request.ModifyLoadBalancerBackendMemberRequest) (*upcloud.LoadBalancerBackendMember, error)
	DeleteLoadBalancerBackendMember(r *request.DeleteLoadBalancerBackendMemberRequest) error

	CreateLoadBalancerResolver(r *request.CreateLoadBalancerResolverRequest) (*upcloud.LoadBalancerResolver, error)
	GetLoadBalancerResolvers(r *request.GetLoadBalancerResolversRequest) ([]upcloud.LoadBalancerResolver, error)
	GetLoadBalancerResolver(r *request.GetLoadBalancerResolverRequest) (*upcloud.LoadBalancerResolver, error)
	ModifyLoadBalancerResolver(r *request.ModifyLoadBalancerRevolverRequest) (*upcloud.LoadBalancerResolver, error)
	DeleteLoadBalancerResolver(r *request.DeleteLoadBalancerResolverRequest) error
}

var _ LoadBalancer = (*Service)(nil)

// GetLoadBalancers retrieves a list of load balancers.
func (s *Service) GetLoadBalancers(r *request.GetLoadBalancersRequest) ([]*upcloud.LoadBalancer, error) {
	var loadBalancers []*upcloud.LoadBalancer
	res, err := s.basicGetRequest(r.RequestURL())
	if err != nil {
		return nil, err
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
		return nil, err
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
		return nil, err
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
		return nil, err
	}

	err = json.Unmarshal(res, &loadBalancer)
	if err != nil {
		return nil, err
	}

	return &loadBalancer, nil
}

// DeleteLoadBalancer deletes an existing load balancer.
func (s *Service) DeleteLoadBalancer(r *request.DeleteLoadBalancerRequest) error {
	err := s.client.PerformJSONDeleteRequest(s.client.CreateRequestURL(r.RequestURL()))

	if err != nil {
		return err
	}

	return nil
}

func (s *Service) GetLoadBalancerBackends(r *request.GetLoadBalancerBackendsRequest) ([]upcloud.LoadBalancerBackend, error) {
	backends := make([]upcloud.LoadBalancerBackend, 0)
	res, err := s.basicGetRequest(r.RequestURL())
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(res, &backends)
	return backends, err
}

func (s *Service) GetLoadBalancerBackend(r *request.GetLoadBalancerBackendRequest) (*upcloud.LoadBalancerBackend, error) {
	var backend upcloud.LoadBalancerBackend
	res, err := s.basicGetRequest(r.RequestURL())
	if err != nil {
		return nil, err
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
		return nil, err
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
		return nil, err
	}

	err = json.Unmarshal(res, &backend)
	return &backend, err
}

func (s *Service) DeleteLoadBalancerBackend(r *request.DeleteLoadBalancerBackendRequest) error {
	return s.client.PerformJSONDeleteRequest(s.client.CreateRequestURL(r.RequestURL()))
}

func (s *Service) GetLoadBalancerBackendMembers(r *request.GetLoadBalancerBackendMembersRequest) ([]upcloud.LoadBalancerBackendMember, error) {
	members := make([]upcloud.LoadBalancerBackendMember, 0)
	res, err := s.basicGetRequest(r.RequestURL())
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(res, &members)
	return members, err
}

func (s *Service) GetLoadBalancerBackendMember(r *request.GetLoadBalancerBackendMemberRequest) (*upcloud.LoadBalancerBackendMember, error) {
	var member upcloud.LoadBalancerBackendMember
	res, err := s.basicGetRequest(r.RequestURL())
	if err != nil {
		return nil, err
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
		return nil, err
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
		return nil, err
	}

	err = json.Unmarshal(res, &member)
	return &member, err
}

func (s *Service) DeleteLoadBalancerBackendMember(r *request.DeleteLoadBalancerBackendMemberRequest) error {
	return s.client.PerformJSONDeleteRequest(s.client.CreateRequestURL(r.RequestURL()))
}

func (s *Service) CreateLoadBalancerResolver(r *request.CreateLoadBalancerResolverRequest) (*upcloud.LoadBalancerResolver, error) {
	var resolver upcloud.LoadBalancerResolver

	reqBody, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	res, err := s.client.PerformJSONPostRequest(s.client.CreateRequestURL(r.RequestURL()), reqBody)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(res, &resolver)
	return &resolver, err
}

func (s *Service) GetLoadBalancerResolvers(r *request.GetLoadBalancerResolversRequest) ([]upcloud.LoadBalancerResolver, error) {
	resolvers := make([]upcloud.LoadBalancerResolver, 0)

	res, err := s.basicGetRequest(r.RequestURL())
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(res, &resolvers)
	return resolvers, err
}

func (s *Service) GetLoadBalancerResolver(r *request.GetLoadBalancerResolverRequest) (*upcloud.LoadBalancerResolver, error) {
	var resolver upcloud.LoadBalancerResolver

	res, err := s.basicGetRequest(r.RequestURL())
	if err != nil {
		return nil, err
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
		return nil, err
	}

	err = json.Unmarshal(res, &resolver)
	return &resolver, err
}

func (s *Service) DeleteLoadBalancerResolver(r *request.DeleteLoadBalancerResolverRequest) error {
	return s.client.PerformJSONDeleteRequest(s.client.CreateRequestURL(r.RequestURL()))
}
