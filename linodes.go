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

// LinodeInstance represents a linode object
type LinodeInstance struct {
	CreatedStr string `json:"created"`
	UpdatedStr string `json:"updated"`

	ID         int
	Created    *time.Time `json:"-"`
	Updated    *time.Time `json:"-"`
	Region     string
	Alerts     *LinodeAlert
	Backups    *LinodeBackup
	Snapshot   *LinodeBackup
	Image      string
	Group      string
	IPv4       []string
	IPv6       string
	Label      string
	Type       string
	Status     string
	Hypervisor string
	Specs      *LinodeSpec
}

func (l *LinodeInstance) fixDates() *LinodeInstance {
	l.Created, _ = parseDates(l.CreatedStr)
	l.Updated, _ = parseDates(l.UpdatedStr)
	return l
}

// LinodeKernel represents a linode kernel object
type LinodeKernel struct {
	ID           int
	Label        string
	Version      string
	KVM          bool
	XEN          bool
	Architecture string
	PVOPS        bool
}

// LinodePrice represents a linode type price object
type LinodePrice struct {
	Hourly  float32
	Monthly float32
}

// LinodeBackupsAddon represents a linode backups addon object
type LinodeBackupsAddon struct {
	Price *LinodePrice
}

// LinodeAddons represent the linode addons object
type LinodeAddons struct {
	Backups *LinodeBackupsAddon
}

// LinodeType represents a linode type object
type LinodeType struct {
	ID         int
	Disk       int
	Class      string // enum: nanode, standard, highmem
	Price      *LinodePrice
	Label      string
	Addons     *LinodeAddons
	NetworkOut int `json:"network_out"`
	Memory     int
	Transfer   int
	VCPUs      int
}

// LinodeInstancesPagedResponse represents a linode API response for listing
type LinodeInstancesPagedResponse struct {
	Page, Pages, Results int
	Data                 []*LinodeInstance
}

// LinodeKernelsPagedResponse represents a linode kernels API response for listing
type LinodeKernelsPagedResponse struct {
	Page, Pages, Results int
	Data                 []*LinodeKernel
}

// LinodeTypesPagedResponse represents a linode types API response for listing
type LinodeTypesPagedResponse struct {
	Page, Pages, Results int
	Data                 []*LinodeType
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
	r, err := c.R().
		SetResult(&LinodeInstancesPagedResponse{}).
		Get(e)
	if err != nil {
		return nil, err
	}
	l := r.Result().(*LinodeInstancesPagedResponse).Data
	for _, el := range l {
		el.fixDates()
	}
	return l, nil
}

// ListKernels lists linode kernels
func (c *Client) ListKernels() ([]*LinodeKernel, error) {
	e, err := c.Kernels.Endpoint()
	if err != nil {
		return nil, err
	}
	r, err := c.R().
		SetResult(&LinodeKernelsPagedResponse{}).
		Get(e)
	if err != nil {
		return nil, err
	}
	l := r.Result().(*LinodeKernelsPagedResponse).Data
	return l, nil
}

// ListTypes lists linode types
func (c *Client) ListTypes() ([]*LinodeType, error) {
	e, err := c.Types.Endpoint()
	if err != nil {
		return nil, err
	}
	r, err := c.R().
		SetResult(&LinodeTypesPagedResponse{}).
		Get(e)
	if err != nil {
		return nil, err
	}
	l := r.Result().(*LinodeTypesPagedResponse).Data
	return l, nil
}

// GetInstance gets the instance with the provided ID
func (c *Client) GetInstance(linodeID int) (*LinodeInstance, error) {
	e, err := c.Instances.Endpoint()
	if err != nil {
		return nil, err
	}
	e = fmt.Sprintf("%s/%d", e, linodeID)
	r, err := c.R().
		SetResult(&LinodeInstance{}).
		Get(e)
	if err != nil {
		return nil, err
	}
	return r.Result().(*LinodeInstance).fixDates(), nil
}

// GetKernel gets the kernel with the provided ID
func (c *Client) GetKernel(kernelID int) (*LinodeKernel, error) {
	e, err := c.Kernels.Endpoint()
	if err != nil {
		return nil, err
	}
	e = fmt.Sprintf("%s/%d", e, kernelID)
	r, err := c.R().
		SetResult(&LinodeKernel{}).
		Get(e)
	if err != nil {
		return nil, err
	}
	return r.Result().(*LinodeKernel), nil
}

// GetType gets the type with the provided ID
func (c *Client) GetType(typeID int) (*LinodeType, error) {
	e, err := c.Types.Endpoint()
	if err != nil {
		return nil, err
	}
	e = fmt.Sprintf("%s/%d", e, typeID)
	r, err := c.R().
		SetResult(&LinodeType{}).
		Get(e)
	if err != nil {
		return nil, err
	}
	return r.Result().(*LinodeType), nil
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
	r, err := c.R().
		SetHeader("Content-Type", "application/json").
		SetBody(bodyStr).
		Post(e)

	return settleBoolResponseOrError(r, err)
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

	r, err := req.
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		Post(e)

	if err != nil {
		return nil, err
	}

	return r.Result().(*LinodeInstance).fixDates(), nil
}

// RebootInstance reboots a Linode instance
func (c *Client) RebootInstance(id int, configID int) (bool, error) {
	body := fmt.Sprintf("{\"config_id\":\"%d\"}", configID)

	e, err := c.Instances.Endpoint()
	if err != nil {
		return false, err
	}

	e = fmt.Sprintf("%s/%d/reboot", e, id)

	r, err := c.R().
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		Post(e)

	return settleBoolResponseOrError(r, err)
}

// MutateInstance Upgrades a Linode to its next generation.
func (c *Client) MutateInstance(id int) (bool, error) {
	e, err := c.Instances.Endpoint()
	if err != nil {
		return false, err
	}
	e = fmt.Sprintf("%s/%d/mutate", e, id)

	r, err := c.R().Post(e)
	return settleBoolResponseOrError(r, err)
}

// RebuildInstanceOptions is a struct representing the options to send to the rebuild linode endpoint
type RebuildInstanceOptions struct {
	Image           string
	RootPass        string
	AuthorizedKeys  []string
	StackscriptID   int
	StackscriptData map[string]string
	Booted          bool
}

// RebuildInstance Deletes all Disks and Configs on this Linode,
// then deploys a new Image to this Linode with the given attributes.
func (c *Client) RebuildInstance(id int, opts *RebuildInstanceOptions) (*LinodeInstance, error) {
	o, err := json.Marshal(opts)
	if err != nil {
		return nil, err
	}
	b := string(o)
	e, err := c.Instances.Endpoint()
	if err != nil {
		return nil, err
	}
	e = fmt.Sprintf("%s/%d/rebuild", e, id)
	r, err := c.R().
		SetHeader("Content-Type", "application/json").
		SetBody(b).
		SetResult(&LinodeInstance{}).
		Post(e)
	if err != nil {
		return nil, err
	}
	return r.Result().(*LinodeInstance).fixDates(), nil
}

// ResizeInstance resizes an instance to new Linode type
func (c *Client) ResizeInstance(id int, linodeType string) (bool, error) {
	body := fmt.Sprintf("{\"type\":\"%s\"}", linodeType)

	e, err := c.Instances.Endpoint()
	if err != nil {
		return false, err
	}
	e = fmt.Sprintf("%s/%d/resize", e, id)

	r, err := c.R().
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		Post(e)

	return settleBoolResponseOrError(r, err)
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
