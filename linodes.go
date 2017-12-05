package golinode

import (
	"encoding/json"
	"fmt"

	"github.com/go-resty/resty"
)

/*
 * https://developers.linode.com/v4/reference/endpoints/linode/instances
 */

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
	resp, err := c.R().
		SetResult(&LinodeInstancesPagedResponse{}).
		Get(instancesEndpoint)
	if err != nil {
		return nil, err
	}
	return resp.Result().(*LinodeInstancesPagedResponse).Data, nil
}

// GetInstance gets the instance with the provided ID
func (c *Client) GetInstance(linodeID int) (*LinodeInstance, error) {
	resp, err := c.R().
		SetResult(&LinodeInstance{}).
		Get(fmt.Sprintf("%s/%d", instancesEndpoint, linodeID))
	if err != nil {
		return nil, err
	}
	return resp.Result().(*LinodeInstance), nil
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

	resp, err := c.R().
		SetHeader("Content-Type", "application/json").
		SetBody(bodyStr).
		Post(fmt.Sprintf("%s/%d/boot", instancesEndpoint, id))

	return settleBoolResponseOrError(resp, err)
}

// CloneInstance - Clones a Linode instance
func (c *Client) CloneInstance(id int, options *LinodeCloneOptions) (*LinodeInstance, error) {
	var body string

	req := c.R().SetResult(&LinodeInstance{})
	endpoint := fmt.Sprintf("%s/%d/clone", instancesEndpoint, id)

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

	return resp.Result().(*LinodeInstance), nil
}

// RebootInstance - Reboots a Linode instance
func (c *Client) RebootInstance(id int, configID int) (bool, error) {
	body := fmt.Sprintf("{\"config_id\":\"%d\"}", configID)
	endpoint := fmt.Sprintf("%s/%d/reboot", instancesEndpoint, id)
	resp, err := c.R().
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		Post(endpoint)
	return settleBoolResponseOrError(resp, err)
}

// ResizeInstance - Resize an instance to new Linode type
func (c *Client) ResizeInstance(id int, linodeType string) (bool, error) {
	body := fmt.Sprintf("{\"type\":\"%s\"}", linodeType)
	endpoint := fmt.Sprintf("%s/%d/resize", instancesEndpoint, id)
	resp, err := c.R().
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		Post(endpoint)
	return settleBoolResponseOrError(resp, err)
}

// ShutdownInstance - Shutdown an instance
func (c *Client) ShutdownInstance(id int) (bool, error) {
	endpoint := fmt.Sprintf("%s/%d/resize", instancesEndpoint, id)
	return settleBoolResponseOrError(c.R().Post(endpoint))
}

func settleBoolResponseOrError(resp *resty.Response, err error) (bool, error) {
	if err != nil {
		return false, err
	}
	return true, nil
}
