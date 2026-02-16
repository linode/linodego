package unit

import (
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/linode/linodego"
	"github.com/linode/linodego/internal/testutil"
	"github.com/stretchr/testify/mock"
)

type MockResponse struct {
	StatusCode int
	Body       interface{}
}

// ClientBaseCase provides a base for unit tests
type ClientBaseCase struct {
	Client  *linodego.Client
	Mock    *mock.Mock
	BaseURL string // Base URL up to /v4
}

// SetUp initializes the Linode client using the mock HTTP client
func (c *ClientBaseCase) SetUp(t *testing.T) {
	c.Mock = &mock.Mock{}
	c.Client = testutil.CreateMockClient(t, linodego.NewClient)
	c.BaseURL = "https://api.linode.com/v4/"
}

func (c *ClientBaseCase) TearDown(t *testing.T) {
	httpmock.DeactivateAndReset() // Reset HTTP mock after tests
	c.Mock.AssertExpectations(t)
}

// MockGet mocks a GET request to the client
func (c *ClientBaseCase) MockGet(path string, response interface{}) {
	fullURL := c.BaseURL + path
	httpmock.RegisterResponder("GET", fullURL, httpmock.NewJsonResponderOrPanic(http.StatusOK, response))
}

// MockPost mocks a POST request for a given path with the provided response body
func (c *ClientBaseCase) MockPost(path string, response interface{}) {
	fullURL := c.BaseURL + path
	httpmock.RegisterResponder("POST", fullURL, httpmock.NewJsonResponderOrPanic(http.StatusOK, response))
}

// MockPut mocks a PUT request for a given path with the provided response body
func (c *ClientBaseCase) MockPut(path string, response interface{}) {
	fullURL := c.BaseURL + path
	httpmock.RegisterResponder("PUT", fullURL, httpmock.NewJsonResponderOrPanic(http.StatusOK, response))
}

// MockDelete mocks a DELETE request for a given path with the provided response body
func (c *ClientBaseCase) MockDelete(path string, response interface{}) {
	fullURL := c.BaseURL + path
	httpmock.RegisterResponder("DELETE", fullURL, httpmock.NewJsonResponderOrPanic(http.StatusOK, response))
}

// MonitorClientBaseCase provides a base for unit tests
type MonitorClientBaseCase struct {
	MonitorClient *linodego.MonitorClient
	Mock          *mock.Mock
	BaseURL       string // Base monitor-api URL
}

// SetUp initializes the Monitor client using the mock HTTP client
func (c *MonitorClientBaseCase) SetUp(t *testing.T) {
	c.Mock = &mock.Mock{}
	c.MonitorClient = testutil.CreateMockClient(t, linodego.NewMonitorClient)
	c.BaseURL = "https://monitor-api.linode.com/v2beta/"
}

func (c *MonitorClientBaseCase) TearDown(t *testing.T) {
	httpmock.DeactivateAndReset() // Reset HTTP mock after tests
	c.Mock.AssertExpectations(t)
}

// MockPost mocks a POST request for a given path with the provided response body
func (c *MonitorClientBaseCase) MockPost(path string, response interface{}) {
	fullURL := c.BaseURL + path
	httpmock.RegisterResponder("POST", fullURL, httpmock.NewJsonResponderOrPanic(http.StatusOK, response))
}
