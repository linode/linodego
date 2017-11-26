package golinode

import (
	"errors"
	"fmt"
	"os"

	"github.com/go-resty/resty"
)

const (
	APIHost    = "api.linode.com"
	APIVersion = "v4"
	APIProto   = "https"
	// Version of golinode
	Version   = "1.0.0"
	apiEnvVar = "LINODE_API_KEY"
)

// Client is a wrapper around the Resty client
type Client struct {
	apiKey string
	Resty  *resty.Client
}

// LinodePagedResponse is an abstraction of the paged response data from the Linode API
type LinodePagedResponse struct {
	Page    int
	Pages   int
	Results int
}

type LinodeResponsePager interface {
	Page() int
	Pages() int
	Results() interface{}
}

// R wraps resty's R method
func (c *Client) R() *resty.Request {
	return c.Resty.R()
}

// SetDebug sets the debug on resty's client
func (c *Client) SetDebug(debug bool) *Client {
	c.Resty.SetDebug(debug)
	return c
}

// NewClient factory to create new Client struct
func NewClient(apiKey string) (*Client, error) {
	envAPIKey, ok := os.LookupEnv(apiEnvVar)
	if ok {
		apiKey = envAPIKey
	}
	if len(apiKey) == 0 || apiKey == "" {
		return nil, errors.New("No API key was provided or LINODE_API_KEY was not set")
	}

	restyClient := resty.New().
		SetHostURL(fmt.Sprintf("%s://%s/%s", APIProto, APIHost, APIVersion)).
		SetAuthToken(apiKey).
		SetHeader("User-Agent", fmt.Sprintf("go-linode %s https://github.com/chiefy/go-linode", Version))

	return &Client{
		apiKey: apiKey,
		Resty:  restyClient,
	}, nil
}
