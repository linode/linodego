package golinode

import (
	"encoding/json"
	"fmt"
	"time"
)

// LinodeImagesPagedResponse represents a linode API response for listing of images
type LinodeVolumesPagedResponse struct {
	Page, Pages, Results int
	Data                 []*LinodeVolume
}

func (l *LinodeVolume) fixDates() *LinodeVolume {
	l.Created, _ = parseDates(l.CreatedStr)
	l.Updated, _ = parseDates(l.UpdatedStr)
	return l
}

// LinodeVolume represents a linode volume object
type LinodeVolume struct {
	CreatedStr string `json:"created"`
	UpdatedStr string `json:"updated"`

	ID       int
	Label    string
	Status   string
	Region   string
	Size     int
	LinodeID int        `json:"linode_id"`
	Created  *time.Time `json:"-"`
	Updated  *time.Time `json:"-"`
}

type LinodeVolumeAttachOptions struct {
	LinodeID int
	ConfigID int
}

// ListVolumes will list linode volumes
func (c *Client) ListVolumes() ([]*LinodeVolume, error) {
	e, err := c.Volumes.Endpoint()
	if err != nil {
		return nil, err
	}
	resp, err := c.R().
		SetResult(&LinodeVolumesPagedResponse{}).
		Get(e)
	if err != nil {
		return nil, err
	}
	list := resp.Result().(*LinodeVolumesPagedResponse).Data
	for _, el := range list {
		el.fixDates()
	}
	return list, nil
}

// GetInstance gets the instance with the provided ID
func (c *Client) GetVolume(volumeID int) (*LinodeVolume, error) {
	e, err := c.Volumes.Endpoint()
	if err != nil {
		return nil, err
	}
	e = fmt.Sprintf("%s/%d", e, volumeID)
	resp, err := c.R().
		SetResult(&LinodeVolume{}).
		Get(e)
	if err != nil {
		return nil, err
	}
	i := resp.Result().(*LinodeVolume).fixDates()
	return i, nil
}

// BootInstance will boot a new linode instance
func (c *Client) AttachVolume(id int, options *LinodeVolumeAttachOptions) (bool, error) {
	body := ""
	if bodyData, err := json.Marshal(options); err == nil {
		body = string(bodyData)
	} else {
		return false, err
	}

	e, err := c.Volumes.Endpoint()
	if err != nil {
		return false, err
	}

	e = fmt.Sprintf("%s/%d/attach", e, id)
	resp, err := c.R().
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		Post(e)

	return settleBoolResponseOrError(resp, err)
}

// CloneInstance clones a Linode instance
func (c *Client) CloneVolume(id int, label string) (*LinodeVolume, error) {
	body := fmt.Sprintf("{\"label\":\"%d\"}", label)

	e, err := c.Volumes.Endpoint()
	if err != nil {
		return nil, err
	}
	e = fmt.Sprintf("%s/%d/clone", e, id)

	req := c.R().SetResult(&LinodeVolume{})

	resp, err := req.
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		Post(e)

	if err != nil {
		return nil, err
	}

	return resp.Result().(*LinodeVolume).fixDates(), nil
}

// DetachVolume detaches a Linode instance
func (c *Client) DetachVolume(id int) (bool, error) {
	body := ""

	e, err := c.Volumes.Endpoint()
	if err != nil {
		return false, err
	}

	e = fmt.Sprintf("%s/%d/detach", e, id)

	resp, err := c.R().
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		Post(e)

	return settleBoolResponseOrError(resp, err)
}

// ResizeInstance resizes an instance to new Linode type
func (c *Client) ResizeVolume(id int, size int) (bool, error) {
	body := fmt.Sprintf("{\"size\":\"%d\"}", size)

	e, err := c.Volumes.Endpoint()
	if err != nil {
		return false, err
	}
	e = fmt.Sprintf("%s/%d/resize", e, id)

	resp, err := c.R().
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		Post(e)

	return settleBoolResponseOrError(resp, err)
}
