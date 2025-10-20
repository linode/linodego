package linodego

import (
	"context"
	"encoding/json"
	"time"

	"github.com/linode/linodego/internal/parseabletime"
)

type NetLoadBalancer struct {
	// This NetLoadBalancer's unique ID.
	ID int `json:"id"`
	// This NetLoadBalancer's label. These must be unique on your Account.
	Label string `json:"label"`
	// The Region where this NetLoadBalancer is located.
	Region string `json:"region"`
	// The IPv4 address of this NetLoadBalancer.
	AddressV4 string `json:"address_v4"`
	// The IPv6 address of this NetLoadBalancer.
	AddressV6 string `json:"address_v6"`
	// The status of this NetLoadBalancer.
	Status string `json:"status"`
	// This NetLoadBalancer's date and time of creation.
	Created *time.Time `json:"-"`
	// This NetLoadBalancer's date and time of last update.
	Updated *time.Time `json:"-"`
	// This NetLoadBalancer's date and time of last composite update.
	LastCompositeUpdated *time.Time `json:"-"`
	// An array of listeners for this NetLoadBalancer.
	Listeners []NetLoadBalancerListener `json:"listeners"`
}

type NetLoadBalancerCreateOptions struct {
	// This NetLoadBalancer's label. These must be unique on your Account.
	Label string `json:"label"`
	// The Region where this NetLoadBalancer is located.
	Region string `json:"region"`
	// An array of listeners for this NetLoadBalancer.
	Listeners []NetLoadBalancerListenerCreateOptions `json:"listeners,omitempty"`
}

type NetLoadBalancerUpdateOptions struct {
	// This NetLoadBalancer's label. These must be unique on your Account.
	Label string `json:"label,omitempty"`
	// An array of listeners for this NetLoadBalancer.
	Listeners []NetLoadBalancerListenerUpdateOptions `json:"listeners,omitempty"`
}

func (i *NetLoadBalancer) UnmarshalJSON(b []byte) error {
	type Mask NetLoadBalancer

	p := struct {
		*Mask
		Created              *parseabletime.ParseableTime `json:"created"`
		Updated              *parseabletime.ParseableTime `json:"updated"`
		LastCompositeUpdated *parseabletime.ParseableTime `json:"last_composite_updated"`
	}{
		Mask: (*Mask)(i),
	}

	if err := json.Unmarshal(b, &p); err != nil {
		return err
	}

	i.Created = (*time.Time)(p.Created)
	i.Updated = (*time.Time)(p.Updated)
	i.LastCompositeUpdated = (*time.Time)(p.LastCompositeUpdated)

	return nil
}

func (i *NetLoadBalancer) GetCreateOptions() NetLoadBalancerCreateOptions {
	opts := make([]NetLoadBalancerListenerCreateOptions, len(i.Listeners))
	for i, listener := range i.Listeners {
		opts[i] = listener.GetCreateOptions()
	}
	return NetLoadBalancerCreateOptions{
		Label:     i.Label,
		Region:    i.Region,
		Listeners: opts,
	}
}

func (i *NetLoadBalancer) GetUpdateOptions() NetLoadBalancerUpdateOptions {
	opts := make([]NetLoadBalancerListenerUpdateOptions, len(i.Listeners))
	for i, listener := range i.Listeners {
		opts[i] = listener.GetUpdateOptions()
	}
	return NetLoadBalancerUpdateOptions{
		Label:     i.Label,
		Listeners: opts,
	}
}

// ListNetLoadBalancers retrieves a list of NetLoadBalancers
func (c *Client) ListNetLoadBalancers(ctx context.Context, opts *ListOptions) ([]NetLoadBalancer, error) {
	return getPaginatedResults[NetLoadBalancer](ctx, c, "netloadbalancers", opts)
}

// GetNetLoadBalancer retrieves a NetLoadBalancer by ID
func (c *Client) GetNetLoadBalancer(ctx context.Context, netloadbalancerID int) (*NetLoadBalancer, error) {
	e := formatAPIPath("netloadbalancers/%d", netloadbalancerID)
	return doGETRequest[NetLoadBalancer](ctx, c, e)
}

// CreateNetLoadBalancer creates a new NetLoadBalancer
func (c *Client) CreateNetLoadBalancer(ctx context.Context, opts NetLoadBalancerCreateOptions) (*NetLoadBalancer, error) {
	return doPOSTRequest[NetLoadBalancer](ctx, c, "netloadbalancers", opts)
}

// UpdateNetLoadBalancer updates a NetLoadBalancer
func (c *Client) UpdateNetLoadBalancer(ctx context.Context, netloadbalancerID int, opts NetLoadBalancerUpdateOptions) (*NetLoadBalancer, error) {
	e := formatAPIPath("netloadbalancers/%d", netloadbalancerID)
	return doPUTRequest[NetLoadBalancer](ctx, c, e, opts)
}

// DeleteNetLoadBalancer deletes a NetLoadBalancer
func (c *Client) DeleteNetLoadBalancer(ctx context.Context, netloadbalancerID int) error {
	e := formatAPIPath("netloadbalancers/%d", netloadbalancerID)
	return doDELETERequest(ctx, c, e)
}
