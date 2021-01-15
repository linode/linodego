package linodego

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/linode/linodego/internal/parseabletime"
)

type VLAN struct {
	ID          int          `json:"id"`
	Description string       `json:"description"`
	Region      string       `json:"region"`
	Linodes     []VLANLinode `json:"linodes"`
	CIDRBlock   string       `json:"cidr_block"`
	Created     *time.Time   `json:"-"`
}

type VLANLinode struct {
	ID          int    `json:"id"`
	MacAddress  string `json:"mac_address"`
	IPv4Address string `json:"ipv4_address"`
}

// VLANCreateOptions fields are those accepted by CreateVLAN
type VLANCreateOptions struct {
	Description string `json:"description,omitempty"`
	Region      string `json:"region"`
	Linodes     []int  `json:"linodes"`
	CIDRBlock   string `json:"cidr_block,omitempty"`
}

// VLANAttachOptions fields are those accepted by VLANAttach
type VLANAttachOptions struct {
	Linodes []int `json:"linodes"`
}

// VLANDetachOptions fields are those accepted by VLANDetach
type VLANDetachOptions struct {
	Linodes []int `json:"linodes"`
}

// UnmarshalJSON for VLAN responses
func (v *VLAN) UnmarshalJSON(b []byte) error {
	type Mask VLAN

	p := struct {
		*Mask
		Created *parseabletime.ParseableTime `json:"created"`
	}{
		Mask: (*Mask)(v),
	}

	if err := json.Unmarshal(b, &p); err != nil {
		return err
	}

	v.Created = (*time.Time)(p.Created)
	return nil
}

// VLANsPagedResponse represents a Linode API response for listing of VLANs
type VLANsPagedResponse struct {
	*PageOptions
	Data []VLAN `json:"data"`
}

func (VLANsPagedResponse) endpoint(c *Client) string {
	endpoint, err := c.VLANs.Endpoint()
	if err != nil {
		panic(err)
	}
	return endpoint
}

func (resp *VLANsPagedResponse) appendData(r *VLANsPagedResponse) {
	resp.Data = append(resp.Data, r.Data...)
}

// ListVLANs returns a paginated list of VLANs
func (c *Client) ListVLANs(ctx context.Context, opts *ListOptions) ([]VLAN, error) {
	response := VLANsPagedResponse{}

	err := c.listHelper(ctx, &response, opts)

	if err != nil {
		return nil, err
	}

	return response.Data, nil
}

// CreateVLAN creates a single VLAN
func (c *Client) CreateVLAN(ctx context.Context, createOpts VLANCreateOptions) (*VLAN, error) {
	var body string
	e, err := c.VLANs.Endpoint()
	if err != nil {
		return nil, err
	}

	req := c.R(ctx).SetResult(&VLAN{})

	if bodyData, err := json.Marshal(createOpts); err == nil {
		body = string(bodyData)
	} else {
		return nil, NewError(err)
	}

	r, err := coupleAPIErrors(req.
		SetBody(body).
		Post(e))

	if err != nil {
		return nil, err
	}
	return r.Result().(*VLAN), nil
}

// GetVLAN gets a single VLAN with the provided ID
func (c *Client) GetVLAN(ctx context.Context, id int) (*VLAN, error) {
	e, err := c.VLANs.Endpoint()
	if err != nil {
		return nil, err
	}

	req := c.R(ctx)

	e = fmt.Sprintf("%s/%d", e, id)
	r, err := coupleAPIErrors(req.SetResult(&VLAN{}).Get(e))
	if err != nil {
		return nil, err
	}

	return r.Result().(*VLAN), nil
}

// DeleteVLAN deletes a single VLAN with the provided ID
func (c *Client) DeleteVLAN(ctx context.Context, id int) error {
	e, err := c.VLANs.Endpoint()
	if err != nil {
		return err
	}

	req := c.R(ctx)

	e = fmt.Sprintf("%s/%d", e, id)
	_, err = coupleAPIErrors(req.Delete(e))
	return err
}

// AttachVLAN attaches the given Linodes to a given VLAN
func (c *Client) AttachVLAN(ctx context.Context, id int, opts VLANAttachOptions) (*VLAN, error) {
	var body string
	e, err := c.VLANs.Endpoint()
	if err != nil {
		return nil, err
	}

	req := c.R(ctx).SetResult(&VLAN{})

	e = fmt.Sprintf("%s/%d/attach", e, id)
	if bodyData, err := json.Marshal(opts); err == nil {
		body = string(bodyData)
	} else {
		return nil, NewError(err)
	}

	r, err := coupleAPIErrors(req.
		SetBody(body).
		Post(e))

	if err != nil {
		return nil, err
	}
	return r.Result().(*VLAN), nil
}

// DetachVLAN detaches the given Linodes from a given VLAN
func (c *Client) DetachVLAN(ctx context.Context, id int, opts VLANDetachOptions) (*VLAN, error) {
	var body string
	e, err := c.VLANs.Endpoint()
	if err != nil {
		return nil, err
	}

	req := c.R(ctx).SetResult(&VLAN{})

	e = fmt.Sprintf("%s/%d/detach", e, id)
	if bodyData, err := json.Marshal(opts); err == nil {
		body = string(bodyData)
	} else {
		return nil, NewError(err)
	}

	r, err := coupleAPIErrors(req.
		SetBody(body).
		Post(e))

	if err != nil {
		return nil, err
	}
	return r.Result().(*VLAN), nil
}
