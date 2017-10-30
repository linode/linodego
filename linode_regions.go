package golinode

import (
	"fmt"
)

// LinodeRegion represents a linode distribution object
//
type LinodeRegion struct {
	ID      string
	Country string
}

// LinodeRegionPagedResponse represents a linode API response for listing
//
type LinodeRegionPagedResponse struct {
	Page    int
	Pages   int
	Results int
	Data    []*LinodeRegion
}

const (
	region = "regions"
)

func (c *Client) ListRegions() ([]*LinodeRegion, error) {
	req := c.R().SetResult(&LinodeRegionPagedResponse{})

	resp, err := req.Get(region)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() >= 400 {
		return nil, fmt.Errorf("Got bad status code: %d", resp.StatusCode())
	}

	list := resp.Result().(*LinodeRegionPagedResponse)

	return list.Data, nil
}
