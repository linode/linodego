package golinode

import (
	"fmt"
)

// LinodeInstancesPagedResponse represents a linode API response for listing
type LinodeRegionsPagedResponse struct {
	Page, Pages, Results int
	data                 []*LinodeRegion
}

// Data returns data collection from paged response
func (r LinodeRegionsPagedResponse) Data() ([]*LinodeRegion, error) {
	return r.data, nil
}

// LinodeRegion represents a linode distribution object
type LinodeRegion struct {
	ID      string
	Country string
}

const (
	regionEndpoint = "regions"
)

// ListRegions - list all available regions for a Linode instance
func (c *Client) ListRegions() ([]*LinodeRegion, error) {
	req := c.R().SetResult(&LinodeRegionsPagedResponse{})
	resp, err := req.Get(regionEndpoint)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() >= 400 {
		return nil, fmt.Errorf("Got bad status code: %d", resp.StatusCode())
	}
	list := resp.Result().(*LinodeRegionsPagedResponse)
	return list.Data()
}
