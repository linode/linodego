package golinode

import (
	"encoding/json"
	"fmt"

	"github.com/go-resty/resty"
)

/*
 * https://developers.linode.com/v4/reference/endpoints/linode/instances
 */

type LinodeLinodesClient interface {
	ListInstances() ([]*LinodeInstance, error)
	GetInstance(int) (*LinodeInstance, error)
	BootInstance(int, int) (bool, error)
	CloneInstance(int, *LinodeCloneOptions) (bool, error)
	RebootInstance(int, int) (bool, error)
	ResizeInstance(int, string) (bool, error)
	ShutdownInstance(int) (bool, error)
}

// LinodeSnapshot represents a linode backup snapshot
type LinodeSnapshot struct {
	ID       int
	Label    string
	Status   string
	Type     string
	Created  string
	Updated  string
	Finished string
	Configs  []string
}

// LinodeDisk represents a linode disk
type LinodeDisk struct {
	ID         int
	Label      string
	Status     string
	Size       int
	Filesystem string
	Created    string
	Updated    string
}
type LinodeBackup struct {
	Enabled      bool
	Availability string
	Schedule     struct {
		Day    string
		Window string
	}
	LastBackup *LinodeSnapshot
	Disks      []*LinodeDisk
}

type LinodeAlert struct {
	CPU           int
	IO            int
	NetworkIn     int
	NetworkOut    int
	TransferQuote int
}

type LinodeSpec struct {
	Disk     int
	Memory   int
	VCPUs    int
	Transfer int
}

// LinodeInstance represents a linode distribution object
type LinodeInstance struct {
	ID           int
	Created      string
	Updated      string
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

// LinodeInstancesPagedResponse represents a linode API response for listing
type LinodeInstancesPagedResponse struct {
	Page    int
	Pages   int
	Results int
	Data    []*LinodeInstance
}

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

const (
	instanceEndpoint = "linode/instances"
)

// ListInstances lists linode instances
func (c *Client) ListInstances() ([]*LinodeInstance, error) {
	req := c.R().SetResult(&LinodeInstancesPagedResponse{})

	resp, err := req.Get(instanceEndpoint)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() >= 400 {
		return nil, fmt.Errorf("Got bad status code: %d", resp.StatusCode())
	}

	list := resp.Result().(*LinodeInstancesPagedResponse)

	return list.Data, nil
}

// GetInstance gets the instance with the provided ID
func (c *Client) GetInstance(linodeID int) (*LinodeInstance, error) {
	req := c.R().SetResult(&LinodeInstance{})

	resp, err := req.Get(fmt.Sprintf("%s/%d", instanceEndpoint, linodeID))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() >= 400 {
		return nil, fmt.Errorf("Got bad status code: %d", resp.StatusCode())
	}

	return resp.Result().(*LinodeInstance), nil
}

// BootInstance will boot a new linode instance
func (c *Client) BootInstance(id int, configID int) (bool, error) {
	var body string

	if configID != 0 {
		body = fmt.Sprintf("{\"config_id\": \"%d\"}", configID)
	}

	resp, err := c.R().
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		Post(fmt.Sprintf("%s/%d/boot", instanceEndpoint, id))

	return settleBoolResponseOrError(resp, err)
}

// CloneInstance - Clones a Linode instance
func (c *Client) CloneInstance(id int, options *LinodeCloneOptions) (*LinodeInstance, error) {
	var body string

	req := c.R().SetResult(&LinodeInstance{})
	endpoint := fmt.Sprintf("%s/%d/clone", instanceEndpoint, id)

	if bodyData, err := json.Marshal(options); err == nil {
		body = string(bodyData)
	} else {
		return nil, err
	}

	resp, err := req.
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		Post(endpoint)

	if err != nil {
		return nil, err
	}
	if resp.StatusCode() >= 400 {
		return nil, fmt.Errorf("Expected 20x got %d", resp.StatusCode())
	}

	return resp.Result().(*LinodeInstance), nil
}

// RebootInstance - Reboots a Linode instance
func (c *Client) RebootInstance(id int, configID int) (bool, error) {
	body := fmt.Sprintf("{\"config_id\":\"%d\"}", configID)
	endpoint := fmt.Sprintf("%s/%d/reboot", instanceEndpoint, id)
	resp, err := c.R().
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		Post(endpoint)
	return settleBoolResponseOrError(resp, err)
}

// ResizeInstance - Resize an instance to new Linode type
func (c *Client) ResizeInstance(id int, linodeType string) (bool, error) {
	body := fmt.Sprintf("{\"type\":\"%s\"}", linodeType)
	endpoint := fmt.Sprintf("%s/%d/resize", instanceEndpoint, id)
	resp, err := c.R().
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		Post(endpoint)
	return settleBoolResponseOrError(resp, err)
}

// ShutdownInstance - Shutdown an instance
func (c *Client) ShutdownInstance(id int) (bool, error) {
	endpoint := fmt.Sprintf("%s/%d/resize", instanceEndpoint, id)
	return settleBoolResponseOrError(c.R().Post(endpoint))
}

func settleBoolResponseOrError(resp *resty.Response, err error) (bool, error) {
	if err != nil {
		return false, err
	}
	if resp.StatusCode() >= 400 {
		return false, fmt.Errorf("Expected a 20x, got %d", resp.StatusCode())
	}
	return true, nil
}
