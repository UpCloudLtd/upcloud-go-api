package client

import (
	"bytes"
	"fmt"
	"github.com/blang/semver"
	"github.com/hashicorp/go-cleanhttp"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

/**
Constants
*/
const (
	DEFAULT_API_VERSION = "1.2.3"
	DEFAULT_API_BASEURL = "https://api.upcloud.com"

	// The default timeout (in seconds)
	DEFAULT_TIMEOUT = 10
)

/**
Client represents an API client
*/
type Client struct {
	userName   string
	password   string
	httpClient *http.Client

	apiVersion string
	apiBaseUrl string
}

/**
New creates ands returns a new client configured with the specified user and password
*/
func New(userName, password string) *Client {
	client := Client{}

	client.SetUserName(userName)
	client.SetPassword(password)
	client.httpClient = cleanhttp.DefaultClient()
	client.SetTimeout(time.Second * DEFAULT_TIMEOUT)

	client.SetAPIVersion(DEFAULT_API_VERSION)
	client.SetAPIBaseUrl(DEFAULT_API_BASEURL)

	return &client
}

/**
GetUserName returns the user name the client uses
*/
func (c *Client) GetUserName() string {
	return c.userName
}

/**
SetUserName sets the client user to the specified user
*/
func (c *Client) SetUserName(userName string) {
	c.userName = userName
}

/**
SetPassword sets the client password to the specified password
*/
func (c *Client) SetPassword(password string) {
	c.password = password
}

/**
SetTimeout sets the client timeout to the specified amount of seconds
*/
func (c *Client) SetTimeout(timeout time.Duration) {
	c.httpClient.Timeout = timeout
}

/**
SetAPIVersion tells the client which API version to use
*/
func (c *Client) SetAPIVersion(version string) {
	c.apiVersion = version
}

/**
SetAPIBaseUrl tells the client which API URL to use
*/
func (c *Client) SetAPIBaseUrl(url string) {
	c.apiBaseUrl = url
}

/**
CreateRequestUrl creates and returns a complete request URL for the specified API location
*/
func (c *Client) CreateRequestUrl(location string) string {
	return fmt.Sprintf("%s%s", c.getBaseUrl(), location)
}

/**
PerformGetRequest performs a GET request to the specified URL and returns the response body and eventual errors
*/
func (c *Client) PerformGetRequest(url string) ([]byte, error) {
	request, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return nil, err
	}

	return c.performRequest(request)
}

/**
PerformPostRequest performs a POST request to the specified URL and returns the response body and eventual errors
*/
func (c *Client) PerformPostRequest(url string, requestBody []byte) ([]byte, error) {
	var bodyReader io.Reader = nil

	if requestBody != nil {
		bodyReader = bytes.NewBuffer(requestBody)
	}

	request, err := http.NewRequest("POST", url, bodyReader)

	if err != nil {
		return nil, err
	}

	return c.performRequest(request)
}

/**
PerformPutRequest performs a PUT request to the specified URL and returns the response body and eventual errors
*/
func (c *Client) PerformPutRequest(url string, requestBody []byte) ([]byte, error) {
	var bodyReader io.Reader = nil

	if requestBody != nil {
		bodyReader = bytes.NewBuffer(requestBody)
	}

	request, err := http.NewRequest("PUT", url, bodyReader)

	if err != nil {
		return nil, err
	}

	return c.performRequest(request)
}

/**
PerformDeleteRequest performs a DELETE request to the specified URL and returns the response body and eventual errors
*/
func (c *Client) PerformDeleteRequest(url string) error {
	request, err := http.NewRequest("DELETE", url, nil)

	if err != nil {
		return err
	}

	_, err = c.performRequest(request)
	return err
}

/**
Adds common headers to the specified request
*/
func (c *Client) addRequestHeaders(request *http.Request) *http.Request {
	request.SetBasicAuth(c.userName, c.password)
	request.Header.Add("Accept", "application/xml")
	request.Header.Add("Content-Type", "application/xml")

	return request
}

/**
Performs the specified HTTP request and returns the response through handleResponse()
*/
func (c *Client) performRequest(request *http.Request) ([]byte, error) {
	c.addRequestHeaders(request)
	response, err := c.httpClient.Do(request)

	if err != nil {
		return nil, err
	}

	return handleResponse(response)
}

/**
Returns the base URL to use for API requests
*/
func (c *Client) getBaseUrl() string {
	urlVersion, _ := semver.Make(c.apiVersion)

	return fmt.Sprintf("%s/%d.%d", c.apiBaseUrl, urlVersion.Major, urlVersion.Minor)
}

/**
Parses the response and returns either the response body or an error
*/
func handleResponse(response *http.Response) ([]byte, error) {
	defer response.Body.Close()

	// Return an error on unsuccessful requests
	if response.StatusCode < 200 || response.StatusCode > 299 {
		errorBody, _ := ioutil.ReadAll(response.Body)

		return nil, &Error{response.StatusCode, response.Status, errorBody}
	}

	responseBody, err := ioutil.ReadAll(response.Body)

	return responseBody, err
}
