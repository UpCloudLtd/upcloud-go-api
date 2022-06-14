package client

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"os"
)

// ClientContext represents an API client with context support
type ClientContext struct {
	*Client
}

// PerformJSONGetRequest performs a GET request to the specified URL and returns the response body and eventual errors
func (c *ClientContext) PerformJSONGetRequest(ctx context.Context, url string) ([]byte, error) {
	return c.PerformJSONRequest(ctx, http.MethodGet, url, nil)
}

// PerformJSONPostRequest performs a POST request to the specified URL and returns the response body and eventual errors
func (c *ClientContext) PerformJSONPostRequest(ctx context.Context, url string, requestBody []byte) ([]byte, error) {
	return c.performJSONRequestWithPayload(ctx, http.MethodPost, url, requestBody)
}

// PerformJSONPutRequest performs a PUT request to the specified URL and returns the response body and eventual errors
func (c *ClientContext) PerformJSONPutRequest(ctx context.Context, url string, requestBody []byte) ([]byte, error) {
	return c.performJSONRequestWithPayload(ctx, http.MethodPut, url, requestBody)
}

// PerformJSONPatchRequest performs a PATCH request to the specified URL and returns the response body and eventual errors
func (c *ClientContext) PerformJSONPatchRequest(ctx context.Context, url string, requestBody []byte) ([]byte, error) {
	return c.performJSONRequestWithPayload(ctx, http.MethodPatch, url, requestBody)
}

// PerformJSONDeleteRequest performs a DELETE request to the specified URL and returns eventual errors
func (c *ClientContext) PerformJSONDeleteRequest(ctx context.Context, url string) error {
	_, err := c.PerformJSONRequest(ctx, http.MethodDelete, url, nil)
	return err
}

// PerformJSONDeleteRequestWithResponseBody performs a DELETE request to the specified URL and returns
// the response body and eventual errors
func (c *ClientContext) PerformJSONDeleteRequestWithResponseBody(ctx context.Context, url string) ([]byte, error) {
	return c.PerformJSONRequest(ctx, http.MethodDelete, url, nil)
}

// PerformJSONPutUploadRequest performs a PUT request to the specified URL with an io.Reader
// and returns the response body and eventual errors
func (c *ClientContext) PerformJSONPutUploadRequest(ctx context.Context, url string, requestBody io.Reader) ([]byte, error) {
	return c.PerformJSONRequest(ctx, http.MethodPut, url, requestBody)
}

func (c *ClientContext) performJSONRequestWithPayload(ctx context.Context, method, url string, body []byte) ([]byte, error) {
	var bodyReader io.Reader

	if body != nil {
		bodyReader = bytes.NewBuffer(body)
	}
	return c.PerformJSONRequest(ctx, method, url, bodyReader)
}

// Performs the specified HTTP request with context and returns the response through handleResponse()
func (c *ClientContext) PerformJSONRequest(ctx context.Context, method, url string, body io.Reader) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}
	c.AddRequestHeaders(req)
	return c.PerformRequest(req)
}

// NewWithContext creates ands returns a new client with context support configured with the specified user and password
func NewWithContext(userName, password string) *ClientContext {
	return NewWithHTTPClientContext(userName, password, httpClient())
}

// NewWithHTTPClientContext creates ands returns a new client with context support configured with the specified user and password and
// using a supplied `http.Client`.
func NewWithHTTPClientContext(userName string, password string, httpClient *http.Client) *ClientContext {
	return &ClientContext{
		&Client{
			userName:   userName,
			password:   password,
			httpClient: httpClient,
			baseURL:    clientBaseURL(os.Getenv(EnvDebugAPIBaseURL)),
			UserAgent:  userAgent(true),
		},
	}
}
