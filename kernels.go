package golinode

import (
	"fmt"
	"strconv"
)

// LinodeKernel represents a linode kernel object
type LinodeKernel struct {
	ID           string
	Label        string
	Version      string
	KVM          bool
	XEN          bool
	Architecture string
	PVOPS        bool
}

// LinodeKernelsPagedResponse represents a linode kernels API response for listing
type LinodeKernelsPagedResponse struct {
	*PageOptions
	Data []*LinodeKernel
}

// ListKernels lists linode kernels
func (c *Client) ListKernels(opts *ListOptions) ([]*LinodeKernel, error) {
	e, err := c.Kernels.Endpoint()
	if err != nil {
		return nil, err
	}
	req := c.R().SetResult(&LinodeKernelsPagedResponse{})

	if opts != nil {
		req.SetQueryParam("page", strconv.Itoa(opts.Page))
	}

	r, err := req.Get(e)
	if err != nil {
		return nil, err
	}

	data := r.Result().(*LinodeKernelsPagedResponse).Data
	pages := r.Result().(*LinodeKernelsPagedResponse).Pages
	results := r.Result().(*LinodeKernelsPagedResponse).Results

	if opts == nil {
		for page := 2; page <= pages; page = page + 1 {
			next, _ := c.ListKernels(&ListOptions{PageOptions: &PageOptions{Page: page}})
			data = append(data, next...)
		}
	} else {
		opts.Results = results
		opts.Pages = pages
	}

	return data, nil
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
