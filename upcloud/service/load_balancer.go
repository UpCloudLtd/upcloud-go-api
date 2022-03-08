package service

import (
	"encoding/json"
	"fmt"

	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud/request"
)

type LoadBalancer interface {
	GetLoadBalancers(r *request.GetLoadBalancersRequest) ([]*upcloud.LoadBalancer, error)
	GetLoadBalancerDetails(r *request.GetLoadBalancerDetailsRequest) (*upcloud.LoadBalancer, error)
	CreateLoadBalancer(r *request.CreateLoadBalancerRequest) (*upcloud.LoadBalancer, error)
	ModifyLoadBalancer(r *request.ModifyLoadBalancerRequest) (*upcloud.LoadBalancer, error)
	DeleteLoadBalancer(r *request.DeleteLoadBalancerRequest) error
	GetLoadBalancerBackends(r *request.GetLoadBalancerBackendsRequest) ([]*upcloud.LoadBalancerBackend, error)
	GetLoadBalancerBackendDetails(r *request.GetLoadBalancerBackendDetailsRequest) (*upcloud.LoadBalancerBackend, error)
	CreateLoadBalancerBackend(r *request.CreateLoadBalancerBackendRequest) (*upcloud.LoadBalancerBackend, error)
	ModifyLoadBalancerBackend(r *request.ModifyLoadBalancerBackendRequest) (*upcloud.LoadBalancerBackend, error)
	DeleteLoadBalancerBackend(r *request.DeleteLoadBalancerBackendRequest) error
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

// GetLoadBalancerDetails retrieves details of a load balancer.
func (s *Service) GetLoadBalancerDetails(r *request.GetLoadBalancerDetailsRequest) (*upcloud.LoadBalancer, error) {
	loadBalancerDetails := upcloud.LoadBalancer{}
	res, err := s.basicGetRequest(r.RequestURL())
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(res, &loadBalancerDetails)
	if err != nil {
		return nil, err
	}

	return &loadBalancerDetails, nil
}

// CreateLoadBalancer creates a new load balancer.
func (s *Service) CreateLoadBalancer(r *request.CreateLoadBalancerRequest) (*upcloud.LoadBalancer, error) {
	loadBalancerDetails := upcloud.LoadBalancer{}
	requestBody, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	fmt.Printf("%s\n", s.client.CreateRequestURL(r.RequestURL()))

	res, err := s.client.PerformJSONPostRequest(s.client.CreateRequestURL(r.RequestURL()), requestBody)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(res, &loadBalancerDetails)
	if err != nil {
		return nil, err
	}

	return &loadBalancerDetails, nil
}

// ModifyLoadBalancer modifies an existing load balancer.
func (s *Service) ModifyLoadBalancer(r *request.ModifyLoadBalancerRequest) (*upcloud.LoadBalancer, error) {
	loadBalancerDetails := upcloud.LoadBalancer{}
	requestBody, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	res, err := s.client.PerformJSONPatchRequest(s.client.CreateRequestURL(r.RequestURL()), requestBody)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(res, &loadBalancerDetails)
	if err != nil {
		return nil, err
	}

	return &loadBalancerDetails, nil
}

// DeleteLoadBalancer deletes an existing load balancer.
func (s *Service) DeleteLoadBalancer(r *request.DeleteLoadBalancerRequest) error {
	err := s.client.PerformJSONDeleteRequest(s.client.CreateRequestURL(r.RequestURL()))

	if err != nil {
		return err
	}

	return nil
}

func (s *Service) GetLoadBalancerBackends(r *request.GetLoadBalancerBackendsRequest) ([]*upcloud.LoadBalancerBackend, error) {
	var backends []*upcloud.LoadBalancerBackend
	res, err := s.basicGetRequest(r.RequestURL())
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(res, &backends)
	return backends, err
}

func (s *Service) GetLoadBalancerBackendDetails(r *request.GetLoadBalancerBackendDetailsRequest) (*upcloud.LoadBalancerBackend, error) {
	var backendDetails upcloud.LoadBalancerBackend
	res, err := s.basicGetRequest(r.RequestURL())
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(res, &backendDetails)
	return &backendDetails, err
}

func (s *Service) CreateLoadBalancerBackend(r *request.CreateLoadBalancerBackendRequest) (*upcloud.LoadBalancerBackend, error) {
	var backendDetails upcloud.LoadBalancerBackend

	reqBody, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	res, err := s.client.PerformJSONPostRequest(s.client.CreateRequestURL(r.RequestURL()), reqBody)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(res, &backendDetails)
	return &backendDetails, err
}

func (s *Service) ModifyLoadBalancerBackend(r *request.ModifyLoadBalancerBackendRequest) (*upcloud.LoadBalancerBackend, error) {
	var backendDetails upcloud.LoadBalancerBackend

	reqBody, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	res, err := s.client.PerformJSONPostRequest(s.client.CreateRequestURL(r.RequestURL()), reqBody)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(res, &backendDetails)
	return &backendDetails, err
}

func (s *Service) DeleteLoadBalancerBackend(r *request.DeleteLoadBalancerBackendRequest) error {
	return s.client.PerformJSONDeleteRequest(s.client.CreateRequestURL(r.RequestURL()))
}
