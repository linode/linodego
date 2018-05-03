package golinode

import (
	"fmt"
	"time"

	"github.com/go-resty/resty"
)

type InstanceConfig struct {
	CreatedStr string `json:"created"`
	UpdatedStr string `json:"updated"`

	ID          int
	Label       string
	Comments    string
	Devices     *InstanceConfigDeviceMap
	Helpers     *InstanceConfigHelpers
	MemoryLimit int `json:"memory_limit"`
	Kernel      string
	InitRD      int
	RootDevice  string     `json:"root_device"`
	RunLevel    string     `json:"run_level"`
	VirtMode    string     `json:"virt_mode"`
	Created     *time.Time `json:"-"`
	Updated     *time.Time `json:"-"`
}

type InstanceConfigDevice struct {
	DiskID   int `json:"disk_id"`
	VolumeID int `json:"volume_id"`
}

type InstanceConfigDeviceMap struct {
	SDA *InstanceConfigDevice
	SDB *InstanceConfigDevice
	SDC *InstanceConfigDevice
	SDD *InstanceConfigDevice
	SDE *InstanceConfigDevice
	SDF *InstanceConfigDevice
	SDG *InstanceConfigDevice
	SDH *InstanceConfigDevice
}

type InstanceConfigHelpers struct {
	UpdateDBDisabled  bool `json:"updatedb_disabled"`
	Distro            bool
	ModulesDep        bool `json:"modules_dep"`
	Network           bool
	DevTmpFsAutomount bool `json:"devtmpfs_automount"`
}

// InstanceConfigsPagedResponse represents a paginated InstanceConfig API response
type InstanceConfigsPagedResponse struct {
	*PageOptions
	Data []*InstanceConfig
}

// Endpoint gets the endpoint URL for InstanceConfig
func (InstanceConfigsPagedResponse) EndpointWithID(c *Client, id int) string {
	endpoint, err := c.InstanceConfigs.EndpointWithID(id)
	if err != nil {
		panic(err)
	}
	return endpoint
}

// AppendData appends InstanceConfigs when processing paginated InstanceConfig responses
func (resp *InstanceConfigsPagedResponse) AppendData(r *InstanceConfigsPagedResponse) {
	(*resp).Data = append(resp.Data, r.Data...)
}

// SetResult sets the Resty response type of InstanceConfig
func (InstanceConfigsPagedResponse) SetResult(r *resty.Request) {
	r.SetResult(InstanceConfigsPagedResponse{})
}

// ListInstanceConfigs lists InstanceConfigs
func (c *Client) ListInstanceConfigs(linodeID int, opts *ListOptions) ([]*InstanceConfig, error) {
	response := InstanceConfigsPagedResponse{}
	err := c.ListHelperWithID(response, linodeID, opts)
	for _, el := range response.Data {
		el.fixDates()
	}
	return response.Data, err
}

// fixDates converts JSON timestamps to Go time.Time values
func (v *InstanceConfig) fixDates() *InstanceConfig {
	v.Created, _ = parseDates(v.CreatedStr)
	v.Updated, _ = parseDates(v.UpdatedStr)
	return v
}

// GetInstanceConfig gets the template with the provided ID
func (c *Client) GetInstanceConfig(linodeID int, configID int) (*InstanceConfig, error) {
	e, err := c.InstanceConfigs.EndpointWithID(linodeID)
	if err != nil {
		return nil, err
	}
	e = fmt.Sprintf("%s/%d", e, configID)
	r, err := c.R().SetResult(&InstanceConfig{}).Get(e)
	if err != nil {
		return nil, err
	}
	return r.Result().(*InstanceConfig).fixDates(), nil
}
