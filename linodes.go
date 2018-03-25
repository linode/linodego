package golinode

import (
	"fmt"
)

/*
 * https://developers.linode.com/v4/reference/endpoints/linode/instances
 */

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

// GetKernel gets the kernel with the provided ID
func (c *Client) GetKernel(kernelID string) (*LinodeKernel, error) {
	e, err := c.Kernels.Endpoint()
	if err != nil {
		return nil, err
	}
	e = fmt.Sprintf("%s/%s", e, kernelID)
	r, err := c.R().
		SetResult(&LinodeKernel{}).
		Get(e)
	if err != nil {
		return nil, err
	}
	return r.Result().(*LinodeKernel), nil
}

// GetType gets the type with the provided ID
func (c *Client) GetType(typeID string) (*LinodeType, error) {
	e, err := c.Types.Endpoint()
	if err != nil {
		return nil, err
	}
	e = fmt.Sprintf("%s/%s", e, typeID)
	r, err := c.R().
		SetResult(&LinodeType{}).
		Get(e)
	if err != nil {
		return nil, err
	}
	return r.Result().(*LinodeType), nil
}
