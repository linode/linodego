package golinode

import (
	"fmt"
	"time"

	"github.com/go-resty/resty"
)

type InstanceDisk struct {
	CreatedStr string `json:"created"`
	UpdatedStr string `json:"updated"`

	ID         int
	Label      string
	Status     string
	Size       int
	Filesystem string
	Created    *time.Time `json:"-"`
	Updated    *time.Time `json:"-"`
}

// InstanceDisksPagedResponse represents a paginated InstanceDisk API response
type InstanceDisksPagedResponse struct {
	*PageOptions
	Data []*InstanceDisk
}

// Endpoint gets the endpoint URL for InstanceDisk
func (InstanceDisksPagedResponse) EndpointWithID(c *Client, id int) string {
	endpoint, err := c.InstanceDisks.EndpointWithID(id)
	if err != nil {
		panic(err)
	}
	return endpoint
}

// AppendData appends InstanceDisks when processing paginated InstanceDisk responses
func (resp *InstanceDisksPagedResponse) AppendData(r *InstanceDisksPagedResponse) {
	(*resp).Data = append(resp.Data, r.Data...)
}

// SetResult sets the Resty response type of InstanceDisk
func (InstanceDisksPagedResponse) SetResult(r *resty.Request) {
	r.SetResult(InstanceDisksPagedResponse{})
}

// ListInstanceDisks lists InstanceDisks
func (c *Client) ListInstanceDisks(linodeID int, opts *ListOptions) ([]*InstanceDisk, error) {
	response := InstanceDisksPagedResponse{}
	err := c.ListHelperWithID(response, linodeID, opts)
	for _, el := range response.Data {
		el.fixDates()
	}
	return response.Data, err
}

// fixDates converts JSON timestamps to Go time.Time values
func (v *InstanceDisk) fixDates() *InstanceDisk {
	v.Created, _ = parseDates(v.CreatedStr)
	v.Updated, _ = parseDates(v.UpdatedStr)
	return v
}

// GetInstanceDisk gets the template with the provided ID
func (c *Client) GetInstanceDisk(linodeID int, configID int) (*InstanceDisk, error) {
	e, err := c.InstanceDisks.EndpointWithID(linodeID)
	if err != nil {
		return nil, err
	}
	e = fmt.Sprintf("%s/%d", e, configID)
	r, err := c.R().SetResult(&InstanceDisk{}).Get(e)
	if err != nil {
		return nil, err
	}
	return r.Result().(*InstanceDisk).fixDates(), nil
}
