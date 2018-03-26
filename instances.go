package golinode

import (
	"encoding/json"
	"fmt"
	"net"
	"time"

	"github.com/go-resty/resty"
)

// LinodeSpec represents a linode spec
type LinodeSpec struct {
	Disk     int
	Memory   int
	VCPUs    int
	Transfer int
}

// LinodeAlert represents a metric alert
type LinodeAlert struct {
	CPU           int
	IO            int
	NetworkIn     int
	NetworkOut    int
	TransferQuote int
}

// LinodeInstanceDisk represents a linode disk
type LinodeInstanceDisk struct {
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

func (l *LinodeInstanceDisk) fixDates() *LinodeInstanceDisk {
	l.Created, _ = parseDates(l.CreatedStr)
	l.Updated, _ = parseDates(l.UpdatedStr)
	return l
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
	Image      string
	Group      string
	IPv4       []*net.IP
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

type LinodeInstanceConfigDevice struct {
	DiskID   int `json:"disk_id"`
	VolumeID int `json:"volume_id"`
}

type LinodeInstanceConfigDeviceMap struct {
	SDA *LinodeInstanceConfigDevice
	SDB *LinodeInstanceConfigDevice
	SDC *LinodeInstanceConfigDevice
	SDD *LinodeInstanceConfigDevice
	SDE *LinodeInstanceConfigDevice
	SDF *LinodeInstanceConfigDevice
	SDG *LinodeInstanceConfigDevice
	SDH *LinodeInstanceConfigDevice
}

type LinodeInstanceConfigHelpers struct {
	UpdateDBDisabled  bool `json:"updatedb_disabled"`
	Distro            bool
	ModulesDep        bool `json:"modules_dep"`
	Network           bool
	DevTmpFsAutomount bool `json:"devtmpfs_automount"`
}

type LinodeInstanceConfig struct {
	CreatedStr string `json:"created"`
	UpdatedStr string `json:"updated"`

	ID          int
	Label       string
	Comments    string
	Devices     *LinodeInstanceConfigDeviceMap
	Helpers     *LinodeInstanceConfigHelpers
	MemoryLimit int `json:"memory_limit"`
	Kernel      string
	InitRD      int
	RootDevice  string     `json:"root_device"`
	RunLevel    string     `json:"run_level"`
	VirtMode    string     `json:"virt_mode"`
	Created     *time.Time `json:"-"`
	Updated     *time.Time `json:"-"`
}

func (l *LinodeInstanceConfig) fixDates() *LinodeInstanceConfig {
	l.Created, _ = parseDates(l.CreatedStr)
	l.Updated, _ = parseDates(l.UpdatedStr)
	return l
}

// LinodeInstancesPagedResponse represents a linode API response for listing
type LinodeInstancesPagedResponse struct {
	Page, Pages, Results int
	Data                 []*LinodeInstance
}

// LinodeDisksPagedResponse represents a linode API response for listing
type LinodeInstanceDisksPagedResponse struct {
	Page, Pages, Results int
	Data                 []*LinodeInstanceDisk
}

// LinodeConfigsPagedResponse represents a linode API response for listing
type LinodeInstanceConfigsPagedResponse struct {
	Page, Pages, Results int
	Data                 []*LinodeInstanceConfig
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

// ListDisks lists linode disks
func (c *Client) ListInstanceDisks(linodeID int) ([]*LinodeInstanceDisk, error) {
	e, err := c.Instances.Endpoint()
	if err != nil {
		return nil, err
	}
	e = fmt.Sprintf("%s/%d/disks", e, linodeID)
	r, err := c.R().
		SetResult(&LinodeInstanceDisksPagedResponse{}).
		Get(e)
	if err != nil {
		return nil, err
	}
	l := r.Result().(*LinodeInstanceDisksPagedResponse).Data
	for _, el := range l {
		el.fixDates()
	}
	return l, nil
}

// ListConfigs lists linode configs
func (c *Client) ListInstanceConfigs(linodeID int) ([]*LinodeInstanceConfig, error) {
	e, err := c.Instances.Endpoint()
	if err != nil {
		return nil, err
	}
	e = fmt.Sprintf("%s/%d/configs", e, linodeID)
	r, err := c.R().
		SetResult(&LinodeInstanceConfigsPagedResponse{}).
		Get(e)
	if err != nil {
		return nil, err
	}
	l := r.Result().(*LinodeInstanceConfigsPagedResponse).Data
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
	r, err := c.R().
		SetResult(&LinodeInstance{}).
		Get(e)
	if err != nil {
		return nil, err
	}
	return r.Result().(*LinodeInstance).fixDates(), nil
}

// GetDisk gets the linode disk with the provided ID
func (c *Client) GetInstanceDisk(linodeID int, diskID int) (*LinodeInstanceDisk, error) {
	e, err := c.Instances.Endpoint()
	if err != nil {
		return nil, err
	}
	e = fmt.Sprintf("%s/%d/disks/%d", e, linodeID, diskID)
	r, err := c.R().
		SetResult(&LinodeInstanceDisk{}).
		Get(e)
	if err != nil {
		return nil, err
	}
	return r.Result().(*LinodeInstanceDisk).fixDates(), nil
}

// GetConfig gets the linode config with the provided ID
func (c *Client) GetInstanceConfig(linodeID int, configID int) (*LinodeInstanceConfig, error) {
	e, err := c.Instances.Endpoint()
	if err != nil {
		return nil, err
	}
	e = fmt.Sprintf("%s/%d/configs/%d", e, linodeID, configID)
	r, err := c.R().
		SetResult(&LinodeInstanceConfig{}).
		Get(e)
	if err != nil {
		return nil, err
	}
	return r.Result().(*LinodeInstanceConfig).fixDates(), nil
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
