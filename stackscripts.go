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
	e, err := c.StackScripts.Endpoint()
	if err != nil {
		return nil, err
	}

	resp, err := c.R().
		SetResult(&LinodeStackscriptsPagedResponse{}).
		Get(e)
	if err != nil {
		return nil, err
	}

	return resp.Result().(*LinodeStackscriptsPagedResponse).Data, nil
}

// GetStackscript returns a stackscript with specified id
func (c *Client) GetStackscript(id int) (*LinodeStackscript, error) {
	e, err := c.StackScripts.Endpoint()
	if err != nil {
		return nil, err
	}
	e = fmt.Sprintf("%s/%d", e, id)

	resp, err := c.R().
		SetResult(&LinodeStackscriptsPagedResponse{}).
		Get(e)

	if err != nil {
		return nil, err
	}

	return resp.Result().(*LinodeStackscriptsPagedResponse).Data[0], nil
}
