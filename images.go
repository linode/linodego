package golinode

// LinodeImagesPagedResponse represents a linode API response for listing of images
type LinodeImagesPagedResponse struct {
	Page, Pages, Results int
	Data                 []*LinodeImage
}

// LinodeImage represents a linode image object
type LinodeImage struct {
	Architecture string
	DiskMinimum  int `json:"disk_minimum"`
	Deprecated   bool
	Label        string
	ID           string
	Updated      string
	Vendor       string
}

// ListImages will list linode distributions
func (c *Client) ListImages() ([]*LinodeImage, error) {
	e, err := c.Images.Endpoint()
	if err != nil {
		return nil, err
	}
	resp, err := c.R().
		SetResult(&LinodeImagesPagedResponse{}).
		Get(e)
	if err != nil {
		return nil, err
	}
	list := resp.Result().(*LinodeImagesPagedResponse).Data
	return list, nil
}
