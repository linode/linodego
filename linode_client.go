package golinode

import (
	"errors"
	"fmt"
	"os"

	"github.com/go-resty/resty"
)

const (
	apiHost    = "api.linode.com"
	apiVersion = "v4"
	apiProto   = "https"
	// Version of golinode
	Version   = "1.0.0"
	apiEnvVar = "LINODE_API_KEY"
)

// Client is a wrapper around the Resty client
type Client struct {
	apiKey string
	Resty  *resty.Client
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

	restyClient := resty.New()
	restyClient = restyClient.SetHostURL(fmt.Sprintf("%s://%s/%s", apiProto, apiHost, apiVersion))
	restyClient = restyClient.SetAuthToken(apiKey)
	restyClient = restyClient.SetHeader("User-Agent", fmt.Sprintf("go-linode %s https://github.com/chiefy/go-linode", Version))

	return &Client{
		apiKey: apiKey,
		Resty:  restyClient,
	}, nil
}
