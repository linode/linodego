package golinode

import (
	"fmt"
)

// LinodeDistribution represents a linode distribution object
//
type LinodeDistribution struct {
	Architecture string
	DiskMinimum  int `json:"disk_minimum"`
	Deprecated   bool
	Label        string
	ID           string
	Updated      string
	Vendor       string
}

// LinodeDistributionPagedResponse represents a linode API response for listing
//
type LinodeDistributionPagedResponse struct {
	Page    int
	Pages   int
	Results int
	Data    []*LinodeDistribution
}

const (
	distribution = "linode/distributions"
)

func (c *Client) ListDistributions() ([]*LinodeDistribution, error) {
	req := c.R().SetResult(&LinodeDistributionPagedResponse{})

	resp, err := req.Get(distribution)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() >= 400 {
		return nil, fmt.Errorf("Got bad status code: %d", resp.StatusCode())
	}

	list := resp.Result().(*LinodeDistributionPagedResponse)

	return list.Data, nil
}
