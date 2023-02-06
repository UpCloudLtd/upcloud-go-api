package service

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/UpCloudLtd/upcloud-go-api/v5/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v5/upcloud/client"
)

type Client interface {
	// Get performs a GET request to the specified path and returns the response body.
	Get(ctx context.Context, path string) ([]byte, error)
	// Post performs a POST request to the specified path and returns the response body.
	Post(ctx context.Context, path string, body []byte) ([]byte, error)
	// Put performs a PUT request to the specified path and returns the response body.
	Put(ctx context.Context, path string, body []byte) ([]byte, error)
	// Patch performs a PATCH request to the specified path and returns the response body.
	Patch(ctx context.Context, path string, body []byte) ([]byte, error)
	// Delete performs a DELETE request to the specified path and returns the response body.
	Delete(ctx context.Context, path string) ([]byte, error)
	// Do performs a HTTP request using custom request object and returns the response body.
	Do(r *http.Request) ([]byte, error)
}

type requestable interface {
	RequestURL() string
}

type service interface {
	Cloud
	Account
	Firewall
	Host
	IPAddress
	LoadBalancer
	ServerGroup
	Network
	Tag
	Server
	Storage
	ObjectStorage
	ManagedDatabaseServiceManager
	ManagedDatabaseUserManager
	ManagedDatabaseLogicalDatabaseManager
	Permission
	Kubernetes
}

var _ service = (*Service)(nil)

// Service represents the API service with context support. The specified client is used to communicate with the API
type Service struct {
	client Client
}

// Get performs a GET request to the specified location with context and stores the result in the value pointed to by v.
func (s *Service) get(ctx context.Context, location string, v interface{}) error {
	res, err := s.client.Get(ctx, location)
	if err != nil {
		return parseJSONServiceError(err)
	}
	if v == nil {
		return nil
	}
	return json.Unmarshal(res, v)
}

// Create performs a POST request to the specified location with context and stores the response in the value pointed to by v.
func (s *Service) create(ctx context.Context, r requestable, v interface{}) error {
	payload, err := json.Marshal(r)
	if err != nil {
		return err
	}

	res, err := s.client.Post(ctx, r.RequestURL(), payload)
	if err != nil {
		return parseJSONServiceError(err)
	}
	if v == nil {
		return nil
	}
	return json.Unmarshal(res, v)
}

// Modify performs a PATCH request to the specified location with context and stores the response in the value pointed to by v.
func (s *Service) modify(ctx context.Context, r requestable, v interface{}) error {
	payload, err := json.Marshal(r)
	if err != nil {
		return err
	}

	res, err := s.client.Patch(ctx, r.RequestURL(), payload)
	if err != nil {
		return parseJSONServiceError(err)
	}
	if v == nil {
		return nil
	}
	return json.Unmarshal(res, v)
}

// Modify performs a PUT request to the specified location with context and stores the response in the value pointed to by v.
func (s *Service) replace(ctx context.Context, r requestable, v interface{}) error {
	payload, err := json.Marshal(r)
	if err != nil {
		return err
	}

	res, err := s.client.Put(ctx, r.RequestURL(), payload)
	if err != nil {
		return parseJSONServiceError(err)
	}
	if v == nil {
		return nil
	}
	return json.Unmarshal(res, v)
}

// Delete performs a DELETE request to the specified location with context
func (s *Service) delete(ctx context.Context, r requestable) error {
	_, err := s.client.Delete(ctx, r.RequestURL())
	if err != nil {
		return parseJSONServiceError(err)
	}
	return nil
}

func New(client Client) *Service {
	return &Service{client}
}

// Parses an error returned from the client into corresponding error type
func parseJSONServiceError(err error) error {
	if clientError, ok := err.(*client.Error); ok {
		return upcloud.NewError(clientError)
	}
	return err
}
