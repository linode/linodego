package golinode

import "fmt"

// LinodeStackscriptsPagedResponse represents a linode API response for listing
type LinodeStackscriptsPagedResponse struct {
	Page, Pages, Results int
	Data                 []*LinodeStackscript
}

// LinodeStackscript represents a linode stack script
type LinodeStackscript struct {
	ID                int
	Username          string
	Label             string
	Description       string
	Distributions     []*LinodeDistribution
	DeploymentsTotal  int
	DeploymentsActive int
	IsPublic          bool
	Created           string
	Updated           string
	RevNote           string
	UserDefinedFields *map[string]string
}

// ListStackscripts gets all public stackscripts
func (c *Client) ListStackscripts() ([]*LinodeStackscript, error) {
	resp, err := c.R().
		SetResult(&LinodeStackscriptsPagedResponse{}).
		Get(stackscriptsEndpoint)
	if err != nil {
		return nil, err
	}
	return resp.Result().(*LinodeStackscriptsPagedResponse).Data, nil
}

// GetStackscript returns a stackscript with specified id
func (c *Client) GetStackscript(id int) (*LinodeStackscript, error) {
	resp, err := c.R().
		SetResult(&LinodeStackscriptsPagedResponse{}).
		Get(fmt.Sprintf("%s/%d", stackscriptsEndpoint, id))
	if err != nil {
		return nil, err
	}
	return resp.Result().(*LinodeStackscriptsPagedResponse).Data[0], nil
}
