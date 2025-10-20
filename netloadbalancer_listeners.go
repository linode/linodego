package linodego

import (
	"context"
	"encoding/json"
	"time"

	"github.com/linode/linodego/internal/parseabletime"
)

type NetLoadBalancerListener struct {
	// This NetLoadBalancerListener's unique ID.
	ID int `json:"id"`
	// The protocol of this NetLoadBalancerListener.
	Protocol string `json:"protocol"`
	// The port of this NetLoadBalancerListener.
	Port int `json:"port"`
	// This NetLoadBalancerListener's label.
	Label string `json:"label"`
	// This NetLoadBalancerListener's date and time of creation.
	Created *time.Time `json:"-"`
	// This NetLoadBalancerListener's date and time of last update.
	Updated *time.Time `json:"-"`
}

// Need a create options and update options for the NetLoadBalancerListener
type NetLoadBalancerListenerCreateOptions struct {
	// The protocol of this NetLoadBalancerListener.
	Protocol string `json:"protocol,omitempty"`
	// The port of this NetLoadBalancerListener.
	Port int `json:"port,omitempty"`
	// The label of this NetLoadBalancerListener.
	Label string `json:"label,omitempty"`
	// The nodes of this NetLoadBalancerListener.
	Nodes []NetLoadBalancerNodeCreateOptions `json:"nodes,omitempty"`
}

type NetLoadBalancerListenerUpdateOptions struct {
	// The protocol of this NetLoadBalancerListener.
	Protocol string `json:"protocol,omitempty"`
	// The port of this NetLoadBalancerListener.
	Port int `json:"port"`
	// The label of this NetLoadBalancerListener.
	Label string `json:"label,omitempty"`
	// The nodes of this NetLoadBalancerListener.
	Nodes []NetLoadBalancerNodeUpdateOptions `json:"nodes,omitempty"`
}

type NetLoadBalancerListenerNodeWeightsUpdateOptions struct {
	// The nodes of this NetLoadBalancerListener.
	Nodes []NetLoadBalancerListenerNodeWeightUpdateOptions `json:"nodes,omitempty"`
}

type NetLoadBalancerListenerNodeWeightUpdateOptions struct {
	// The ID of the node to update.
	ID int `json:"id"`
	// The weight of the node.
	Weight int `json:"weight"`
}

func (i *NetLoadBalancerListener) UnmarshalJSON(b []byte) error {
	type Mask NetLoadBalancerListener

	p := struct {
		*Mask
		Created *parseabletime.ParseableTime `json:"created"`
		Updated *parseabletime.ParseableTime `json:"updated"`
	}{
		Mask: (*Mask)(i),
	}

	if err := json.Unmarshal(b, &p); err != nil {
		return err
	}

	i.Created = (*time.Time)(p.Created)
	i.Updated = (*time.Time)(p.Updated)

	return nil
}

// GetCreateOptions returns the create options for the NetLoadBalancerListener but does not include the nodes
func (i *NetLoadBalancerListener) GetCreateOptions() NetLoadBalancerListenerCreateOptions {
	return NetLoadBalancerListenerCreateOptions{
		Protocol: i.Protocol,
		Port:     i.Port,
		Label:    i.Label,
	}
}

func (i *NetLoadBalancerListener) GetUpdateOptions() NetLoadBalancerListenerUpdateOptions {
	return NetLoadBalancerListenerUpdateOptions{
		Label: i.Label,
	}
}

// CreateNetLoadBalancerListener creates a new NetLoadBalancerListener
func (c *Client) CreateNetLoadBalancerListener(ctx context.Context, netloadbalancerID int, opts NetLoadBalancerListenerCreateOptions) (*NetLoadBalancerListener, error) {
	e := formatAPIPath("netloadbalancers/%d/listeners", netloadbalancerID)
	return doPOSTRequest[NetLoadBalancerListener](ctx, c, e, opts)
}

// ListNetLoadBalancerListeners retrieves a list of NetLoadBalancerListeners
func (c *Client) ListNetLoadBalancerListeners(ctx context.Context, netloadbalancerID int, opts *ListOptions) ([]NetLoadBalancerListener, error) {
	e := formatAPIPath("netloadbalancers/%d/listeners", netloadbalancerID)
	return getPaginatedResults[NetLoadBalancerListener](ctx, c, e, opts)
}

// GetNetLoadBalancerListener retrieves a NetLoadBalancerListener by ID
func (c *Client) GetNetLoadBalancerListener(ctx context.Context, netloadbalancerID int, listenerID int) (*NetLoadBalancerListener, error) {
	e := formatAPIPath("netloadbalancers/%d/listeners/%d", netloadbalancerID, listenerID)
	return doGETRequest[NetLoadBalancerListener](ctx, c, e)
}

// UpdateNetLoadBalancerListener updates a NetLoadBalancerListener
func (c *Client) UpdateNetLoadBalancerListener(ctx context.Context, netloadbalancerID int, listenerID int, opts NetLoadBalancerListenerUpdateOptions) (*NetLoadBalancerListener, error) {
	e := formatAPIPath("netloadbalancers/%d/listeners/%d", netloadbalancerID, listenerID)
	return doPUTRequest[NetLoadBalancerListener](ctx, c, e, opts)
}

// DeleteNetLoadBalancerListener deletes a NetLoadBalancerListener
func (c *Client) DeleteNetLoadBalancerListener(ctx context.Context, netloadbalancerID int, listenerID int) error {
	e := formatAPIPath("netloadbalancers/%d/listeners/%d", netloadbalancerID, listenerID)
	return doDELETERequest(ctx, c, e)
}

// UpdateNetLoadBalancerListenerNodeWeights updates the weights of the nodes of a NetLoadBalancerListener
// Use this to update the weights of the nodes of a NetLoadBalancerListener in case of frequent changes
// High frequency updates are allowed. No response is returned.
func (c *Client) UpdateNetLoadBalancerListenerNodeWeights(ctx context.Context, netloadbalancerID int, listenerID int, opts NetLoadBalancerListenerNodeWeightsUpdateOptions) error {
	e := formatAPIPath("netloadbalancers/%d/listeners/%d/node-weights", netloadbalancerID, listenerID)
	return doPOSTRequestNoResponseBody(ctx, c, e, opts)
}
