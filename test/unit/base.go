package unit

import (
	"github.com/jarcoal/httpmock"
	"github.com/linode/linodego"
	"github.com/linode/linodego/internal/testutil"
	"github.com/stretchr/testify/mock"
	"testing"
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

// MockGet mocks a GET request to the client.
func (c *ClientBaseCase) MockGet(path string, response interface{}) {
	fullURL := c.BaseURL + path
	httpmock.RegisterResponder("GET", fullURL, httpmock.NewJsonResponderOrPanic(200, response))
}
