package client

import (
	"bytes"
	"fmt"
	"github.com/blang/semver"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

const API_VERSION = "1.2.3"
const API_BASE_URL = "https://api.upcloud.com"

/**
The default timeout (in seconds)
*/
const DEFAULT_TIMEOUT = 10

/**
Represents an API client
*/
type Client struct {
	userName   string
	password   string
	httpClient *http.Client
}

/**
Creates ands returns a new client configured with the specified user and password
*/
func New(userName, password string) *Client {
	client := Client{}

	client.SetUserName(userName)
	client.SetPassword(password)
	client.httpClient = http.DefaultClient
	client.SetTimeout(time.Second * DEFAULT_TIMEOUT)

	return &client
}

/**
Sets the client user to the specified user
*/
func (c *Client) SetUserName(userName string) {
	c.userName = userName
}

/**
Sets the client password to the specified password
*/
func (c *Client) SetPassword(password string) {
	c.password = password
}

/**
Sets the client timeout to the specified amount of seconds
*/
func (c *Client) SetTimeout(timeout time.Duration) {
	c.httpClient.Timeout = timeout
}

/**
Returns a complete request URL for the specified API location
*/
func (c *Client) CreateRequestUrl(location string) string {
	return fmt.Sprintf("%s%s", getBaseUrl(), location)
}

/**
Performs a GET request to the specified URL and returns the response body and eventual errors
*/
func (c *Client) PerformGetRequest(url string) ([]byte, error) {
	request, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return nil, err
	}

	return c.performRequest(request)
}

/**
Performs a POST request to the specified URL and returns the response body and eventual errors
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
Performs a PUT request to the specified URL and returns the response body and eventual errors
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
Performs a DELETE request to the specified URL and returns the response body and eventual errors
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
func getBaseUrl() string {
	urlVersion, _ := semver.Make(API_VERSION)

	return fmt.Sprintf("%s/%d.%d", API_BASE_URL, urlVersion.Major, urlVersion.Minor)
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
