package golinode

// LinodeDistributionsPagedResponse represents a linode API response for listing
type LinodeDistributionsPagedResponse struct {
	Page, Pages, Results int
	Data                 []*LinodeDistribution
}

// LinodeDistribution represents a linode distribution object
type LinodeDistribution struct {
	Architecture string
	DiskMinimum  int `json:"disk_minimum"`
	Deprecated   bool
	Label        string
	ID           string
	Updated      string
	Vendor       string
}

// ListDistributions will list linode distributions
func (c *Client) ListDistributions() ([]*LinodeDistribution, error) {
	e, err := c.Distributions.Endpoint()
	if err != nil {
		return nil, err
	}
	resp, err := c.R().
		SetResult(&LinodeDistributionsPagedResponse{}).
		Get(e)
	if err != nil {
		return nil, err
	}
	list := resp.Result().(*LinodeDistributionsPagedResponse).Data
	return list, nil
}
