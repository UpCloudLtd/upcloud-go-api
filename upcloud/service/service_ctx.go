package service

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud/client"
)

type service interface {
	CloudContext
	AccountContext
	FirewallContext
	HostContext
	IpAddressContext
	LoadBalancerContext
	ServerGroupContext
	NetworkContext
	TagContext
	ServerContext
	StorageContext
	ObjectStorageContext
	ManagedDatabaseServiceManagerContext
	ManagedDatabaseUserManagerContext
	ManagedDatabaseLogicalDatabaseManagerContext
	PermissionContext
}

var _ service = (*ServiceContext)(nil)

// Service represents the API service with context support. The specified client is used to communicate with the API
type ServiceContext struct {
	client *client.ClientContext
}

// Get performs a GET request to the specified location with context and stores the result in the value pointed to by v.
func (s *ServiceContext) get(ctx context.Context, location string, v interface{}) error {
	body, err := s.client.PerformJSONRequest(ctx, http.MethodGet, s.client.CreateRequestURL(location), nil)
	if err != nil {
		return parseJSONServiceError(err)
	}
	return json.Unmarshal(body, v)
}

// Create performs a POST request to the specified location with context and stores the response in the value pointed to by v.
func (s *ServiceContext) create(ctx context.Context, r requestable, v interface{}) error {
	payload, err := json.Marshal(r)
	if err != nil {
		return err
	}

	res, err := s.client.PerformJSONPostRequest(ctx, s.client.CreateRequestURL(r.RequestURL()), payload)
	if err != nil {
		return parseJSONServiceError(err)
	}
	if v == nil {
		return nil
	}
	return json.Unmarshal(res, v)
}

// Modify performs a PATCH request to the specified location with context and stores the response in the value pointed to by v.
func (s *ServiceContext) modify(ctx context.Context, r requestable, v interface{}) error {
	payload, err := json.Marshal(r)
	if err != nil {
		return err
	}

	res, err := s.client.PerformJSONPatchRequest(ctx, s.client.CreateRequestURL(r.RequestURL()), payload)
	if err != nil {
		return parseJSONServiceError(err)
	}
	if v == nil {
		return nil
	}
	return json.Unmarshal(res, v)
}

// Modify performs a PUT request to the specified location with context and stores the response in the value pointed to by v.
func (s *ServiceContext) replace(ctx context.Context, r requestable, v interface{}) error {
	payload, err := json.Marshal(r)
	if err != nil {
		return err
	}

	res, err := s.client.PerformJSONPutRequest(ctx, s.client.CreateRequestURL(r.RequestURL()), payload)
	if err != nil {
		return parseJSONServiceError(err)
	}
	if v == nil {
		return nil
	}
	return json.Unmarshal(res, v)
}

// Delete performs a DELETE request to the specified location with context
func (s *ServiceContext) delete(ctx context.Context, r requestable) error {
	err := s.client.PerformJSONDeleteRequest(ctx, s.client.CreateRequestURL(r.RequestURL()))
	if err != nil {
		return parseJSONServiceError(err)
	}
	return nil
}

func NewWithContext(client *client.ClientContext) *ServiceContext {
	return &ServiceContext{client}
}
