package golinode

// LinodeRegionsPagedResponse represents a linode API response for listing
type LinodeRegionsPagedResponse struct {
	Page, Pages, Results int
	Data                 []*LinodeRegion
}

// LinodeRegion represents a linode distribution object
type LinodeRegion struct {
	ID      string
	Country string
}

// ListRegions - list all available regions for a Linode instance
func (c *Client) ListRegions() ([]*LinodeRegion, error) {
	resp, err := c.R().
		SetResult(&LinodeRegionsPagedResponse{}).
		Get(regionsEndpoint)

	if err != nil {
		return nil, err
	}
	return resp.Result().(*LinodeRegionsPagedResponse).Data, nil
}
