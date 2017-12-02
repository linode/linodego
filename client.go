package golinode

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/go-resty/resty"
)

const (
	// APIHost Linode API hostname
	APIHost = "api.linode.com"
	// APIVersion Linode API version
	APIVersion = "v4"
	// APIProto connect to API with http(s)
	APIProto = "https"
	// Version version of golinode
	Version = "1.0.0"
	// APIEnvVar environment var to check for API key
	APIEnvVar = "LINODE_API_KEY"
)

// Client is a wrapper around the Resty client
type Client struct {
	apiKey string
	resty  *resty.Client
}

// R wraps resty's R method
func (c *Client) R() *resty.Request {
	return c.resty.R()
}

// SetDebug sets the debug on resty's client
func (c *Client) SetDebug(debug bool) *Client {
	c.resty.SetDebug(debug)
	return c
}

// NewClient factory to create new Client struct
func NewClient(codeAPIKey *string, transport http.RoundTripper) (*Client, error) {
	linodeAPIKey := ""

	if codeAPIKey != nil {
		linodeAPIKey = *codeAPIKey
	} else if envAPIKey, ok := os.LookupEnv(APIEnvVar); ok {
		linodeAPIKey = envAPIKey
	}
	if len(linodeAPIKey) == 0 || linodeAPIKey == "" {
		return nil, errors.New("No API key was provided or LINODE_API_KEY was not set")
	}

	restyClient := resty.New().
		SetHostURL(fmt.Sprintf("%s://%s/%s", APIProto, APIHost, APIVersion)).
		SetAuthToken(linodeAPIKey).
		SetTransport(transport).
		SetHeader("User-Agent", fmt.Sprintf("go-linode %s https://github.com/chiefy/go-linode", Version))

	return &Client{
		apiKey: linodeAPIKey,
		resty:  restyClient,
	}, nil
}
