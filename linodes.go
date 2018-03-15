package golinode

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-resty/resty"
)

/*
 * https://developers.linode.com/v4/reference/endpoints/linode/instances
 */

// LinodeDisk represents a linode disk
type LinodeDisk struct {
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

func (l *LinodeDisk) fixDates() *LinodeDisk {
	l.Created, _ = parseDates(l.CreatedStr)
	l.Updated, _ = parseDates(l.UpdatedStr)
	return l
}

// LinodeAlert represents a metric alert
type LinodeAlert struct {
	CPU           int
	IO            int
	NetworkIn     int
	NetworkOut    int
	TransferQuote int
}

// LinodeSpec represents a linode spec
type LinodeSpec struct {
	Disk     int
	Memory   int
	VCPUs    int
	Transfer int
}

// LinodeInstance represents a linode distribution object
type LinodeInstance struct {
	CreatedStr string `json:"created"`
	UpdatedStr string `json:"updated"`

	ID           int
	Created      *time.Time `json:"-"`
	Updated      *time.Time `json:"-"`
	Region       string
	Alerts       *LinodeAlert
	Backups      *LinodeBackup
	Snapshot     *LinodeBackup
	Distribution string
	Group        string
	IPv4         []string
	IPv6         string
	Label        string
	Type         string
	Status       string
	Hypervisor   string
	Specs        *LinodeSpec
}

func (l *LinodeInstance) fixDates() *LinodeInstance {
	l.Created, _ = parseDates(l.CreatedStr)
	l.Updated, _ = parseDates(l.UpdatedStr)
	return l
}

// LinodeInstancesPagedResponse represents a linode API response for listing
type LinodeInstancesPagedResponse struct {
	Page, Pages, Results int
	Data                 []*LinodeInstance
}

// LinodeCloneOptions is an options struct when sending a clone request to the API
type LinodeCloneOptions struct {
	Region         string
	Type           string
	LinodeID       int
	Label          string
	Group          string
	BackupsEnabled bool
	Disks          []string
	Configs        []string
}

// ListInstances lists linode instances
func (c *Client) ListInstances() ([]*LinodeInstance, error) {
	e, err := c.Instances.Endpoint()
	if err != nil {
		return nil, err
	}
	resp, err := c.R().
		SetResult(&LinodeInstancesPagedResponse{}).
		Get(e)
	if err != nil {
		return nil, err
	}
	l := resp.Result().(*LinodeInstancesPagedResponse).Data
	for _, el := range l {
		el.fixDates()
	}
	return l, nil
}

// GetInstance gets the instance with the provided ID
func (c *Client) GetInstance(linodeID int) (*LinodeInstance, error) {
	e, err := c.Instances.Endpoint()
	if err != nil {
		return nil, err
	}
	e = fmt.Sprintf("%s/%d", e, linodeID)
	resp, err := c.R().
		SetResult(&LinodeInstance{}).
		Get(e)
	if err != nil {
		return nil, err
	}
	i := resp.Result().(*LinodeInstance).fixDates()
	return i, nil
}

// BootInstance will boot a new linode instance
func (c *Client) BootInstance(id int, configID int) (bool, error) {
	bodyStr := ""

	if configID != 0 {
		bodyMap := map[string]string{"config_id": string(configID)}
		bodyJSON, err := json.Marshal(bodyMap)
		if err != nil {
			return false, err
		}
		bodyStr = string(bodyJSON)
	}

	e, err := c.Instances.Endpoint()
	if err != nil {
		return false, err
	}

	e = fmt.Sprintf("%s/%d/boot", e, id)
	resp, err := c.R().
		SetHeader("Content-Type", "application/json").
		SetBody(bodyStr).
		Post(e)

	return settleBoolResponseOrError(resp, err)
}

// CloneInstance clones a Linode instance
func (c *Client) CloneInstance(id int, options *LinodeCloneOptions) (*LinodeInstance, error) {
	var body string
	e, err := c.Instances.Endpoint()
	if err != nil {
		return nil, err
	}
	e = fmt.Sprintf("%s/%d/clone", e, id)

	req := c.R().SetResult(&LinodeInstance{})

	if bodyData, err := json.Marshal(options); err == nil {
		body = string(bodyData)
	} else {
		return nil, err
	}

	resp, err := req.
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		Post(e)

	if err != nil {
		return nil, err
	}

	return resp.Result().(*LinodeInstance).fixDates(), nil
}

// RebootInstance reboots a Linode instance
func (c *Client) RebootInstance(id int, configID int) (bool, error) {
	body := fmt.Sprintf("{\"config_id\":\"%d\"}", configID)

	e, err := c.Instances.Endpoint()
	if err != nil {
		return false, err
	}

	e = fmt.Sprintf("%s/%d/reboot", e, id)

	resp, err := c.R().
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		Post(e)

	return settleBoolResponseOrError(resp, err)
}

// ResizeInstance resizes an instance to new Linode type
func (c *Client) ResizeInstance(id int, linodeType string) (bool, error) {
	body := fmt.Sprintf("{\"type\":\"%s\"}", linodeType)

	e, err := c.Instances.Endpoint()
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

// ShutdownInstance - Shutdown an instance
func (c *Client) ShutdownInstance(id int) (bool, error) {
	e, err := c.Instances.Endpoint()
	if err != nil {
		return false, err
	}
	e = fmt.Sprintf("%s/%d/resize", e, id)
	return settleBoolResponseOrError(c.R().Post(e))
}

// ListInstanceVolumes lists volumes attached to a linode instance
func (c *Client) ListInstanceVolumes(id int) ([]*LinodeVolume, error) {
	e, err := c.Instances.Endpoint()
	e = fmt.Sprintf("%s/%d/volumes", e, id)
	if err != nil {
		return nil, err
	}
	resp, err := c.R().
		SetResult(&LinodeVolumesPagedResponse{}).
		Get(e)
	if err != nil {
		return nil, err
	}
	l := resp.Result().(*LinodeVolumesPagedResponse).Data
	for _, el := range l {
		el.fixDates()
	}
	return l, nil
}

func settleBoolResponseOrError(resp *resty.Response, err error) (bool, error) {
	if err != nil {
		return false, err
	}
	return true, nil
}
