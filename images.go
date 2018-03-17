package golinode

import (
	"time"
)

// LinodeImagesPagedResponse represents a linode API response for listing of images
type LinodeImagesPagedResponse struct {
	Page, Pages, Results int
	Data                 []*LinodeImage
}

// LinodeImage represents a linode image object
type LinodeImage struct {
	CreatedStr  string `json:"created"`
	UpdatedStr  string `json:"updated"`
	ID          string
	Label       string
	Description string
	Type        string
	IsPublic    bool
	Size        int
	Vendor      string
	Deprecated  bool

	CreatedBy string     `json:"created_by"`
	Created   *time.Time `json:"-"`
	Updated   *time.Time `json:"-"`
}

func (l *LinodeImage) fixDates() *LinodeImage {
	l.Created, _ = parseDates(l.CreatedStr)
	l.Updated, _ = parseDates(l.UpdatedStr)
	return l
}

// ListImages will list linode images
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
	for _, el := range list {
		el.fixDates()
	}

	return list, nil
}
