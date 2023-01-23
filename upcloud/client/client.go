package client

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/hashicorp/go-cleanhttp"
)

const (
	Version    string = "5.3.0"
	APIVersion string = "1.3"
	APIBaseURL string = "https://api.upcloud.com"

	EnvDebugAPIBaseURL            string = "UPCLOUD_DEBUG_API_BASE_URL"
	EnvDebugSkipCertificateVerify string = "UPCLOUD_DEBUG_SKIP_CERTIFICATE_VERIFY"
)

type config struct {
	username   string
	password   string
	baseURL    string
	httpClient *http.Client
}

// Client represents an API client
type Client struct {
	UserAgent string
	config    config
}

// Get performs a GET request to the specified path and returns the response body.
func (c *Client) Get(ctx context.Context, path string) ([]byte, error) {
	r, err := c.createRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	return c.Do(r)
}

// Post performs a POST request to the specified path and returns the response body.
func (c *Client) Post(ctx context.Context, path string, body []byte) ([]byte, error) {
	r, err := c.createRequest(ctx, http.MethodPost, path, body)
	if err != nil {
		return nil, err
	}
	return c.Do(r)
}

// Put performs a PUT request to the specified path and returns the response body.
func (c *Client) Put(ctx context.Context, path string, body []byte) ([]byte, error) {
	r, err := c.createRequest(ctx, http.MethodPut, path, body)
	if err != nil {
		return nil, err
	}
	return c.Do(r)
}

// Patch performs a PATCH request to the specified path and returns the response body.
func (c *Client) Patch(ctx context.Context, path string, body []byte) ([]byte, error) {
	r, err := c.createRequest(ctx, http.MethodPatch, path, body)
	if err != nil {
		return nil, err
	}
	return c.Do(r)
}

// Delete performs a DELETE request to the specified path and returns the response body.
func (c *Client) Delete(ctx context.Context, path string) ([]byte, error) {
	r, err := c.createRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}
	return c.Do(r)
}

// Do performs HTTP request and returns the response body.
func (c *Client) Do(r *http.Request) ([]byte, error) {
	c.addDefaultHeaders(r)
	response, err := c.config.httpClient.Do(r)
	if err != nil {
		return nil, err
	}

	return handleResponse(response)
}

func (c *Client) createRequest(ctx context.Context, method, path string, body []byte) (*http.Request, error) {
	var bodyReader io.Reader

	if body != nil {
		bodyReader = bytes.NewBuffer(body)
	}
	req, err := http.NewRequestWithContext(ctx, method, c.createRequestURL(path), bodyReader)
	if err != nil {
		return nil, err
	}
	return req, err
}

func (c *Client) addDefaultHeaders(r *http.Request) {
	const (
		accept        string = "Accept"
		contentType   string = "Content-Type"
		userAgent     string = "User-Agent"
		authorization string = "Authorization"
	)
	if _, ok := r.Header[accept]; !ok {
		r.Header.Set(accept, "application/json")
	}
	if _, ok := r.Header[contentType]; !ok {
		r.Header.Set(contentType, "application/json")
	}
	if _, ok := r.Header[userAgent]; !ok {
		r.Header.Set(userAgent, c.UserAgent)
	}
	if _, ok := r.Header[authorization]; !ok && strings.HasPrefix(r.URL.String(), c.config.baseURL) {
		r.SetBasicAuth(c.config.username, c.config.password)
	}
}

// createRequestURL creates and returns a complete request URL for the specified API location using a newer API version.
func (c *Client) createRequestURL(location string) string {
	return fmt.Sprintf("%s%s", c.getBaseURL(), location)
}

// Returns the base URL to use for API requests
func (c *Client) getBaseURL() string {
	if c.config.baseURL == "" {
		c.config.baseURL = clientBaseURL(os.Getenv(EnvDebugAPIBaseURL))
	}
	return fmt.Sprintf("%s/%s", c.config.baseURL, APIVersion)
}

type configFn func(o *config)

func WithBaseURL(baseURL string) configFn {
	return func(c *config) {
		c.baseURL = baseURL
	}
}

func WithHTTPClient(httpClient *http.Client) configFn {
	return func(c *config) {
		c.httpClient = httpClient
	}
}

func WithTimeout(timeout time.Duration) configFn {
	return func(c *config) {
		c.httpClient.Timeout = timeout
	}
}

// New creates and returns a new client configured with the specified user and password
func New(username, password string, c ...configFn) *Client {
	config := config{
		username:   username,
		password:   password,
		baseURL:    clientBaseURL(os.Getenv(EnvDebugAPIBaseURL)),
		httpClient: httpClient(),
	}
	for _, fn := range c {
		fn(&config)
	}
	return &Client{
		UserAgent: userAgent(),
		config:    config,
	}
}

func httpClient() *http.Client {
	var client *http.Client
	if os.Getenv(EnvDebugSkipCertificateVerify) == "1" {
		client = &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true, //nolint
				},
			},
		}
	} else {
		client = cleanhttp.DefaultClient()
	}
	return client
}

func userAgent() string {
	return fmt.Sprintf("upcloud-go-api/%s", Version)
}

func clientBaseURL(URL string) string {
	if URL != "" {
		if u, err := url.Parse(URL); err != nil || u.Scheme == "" || u.Host == "" {
			return APIBaseURL
		}
		return URL
	}
	return APIBaseURL
}

// Parses the response and returns either the response body or an error
func handleResponse(response *http.Response) ([]byte, error) {
	defer response.Body.Close()

	// Return an error on unsuccessful requests
	if response.StatusCode < 200 || response.StatusCode > 299 {
		errorBody, _ := io.ReadAll(response.Body)
		var errorType ErrorType
		switch response.Header.Get("Content-Type") {
		case "application/problem+json":
			errorType = ErrorTypeProblem
		default:
			errorType = ErrorTypeError
		}
		return nil, &Error{response.StatusCode, response.Status, errorBody, errorType}
	}

	responseBody, err := io.ReadAll(response.Body)

	return responseBody, err
}
