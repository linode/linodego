package linodego

import (
	"context"
	"encoding/json"
	"time"

	"github.com/linode/linodego/internal/parseabletime"
)

type NetLoadBalancerNode struct {
	// This NetLoadBalancerNode's unique ID.
	ID int `json:"id"`
	// The ID of the Linode this NetLoadBalancerNode is associated with.
	LinodeID int `json:"linode_id"`
	// The IPv6 address of this NetLoadBalancerNode.
	AddressV6 string `json:"address_v6"`
	// This NetLoadBalancerNode's label.
	Label string `json:"label"`
	// The weight of this NetLoadBalancerNode.
	Weight int `json:"weight"`
	// This NetLoadBalancerNode's date and time of creation.
	Created *time.Time `json:"-"`
	// This NetLoadBalancerNode's date and time of last update.
	Updated *time.Time `json:"-"`
	// This NetLoadBalancerNode's date and time of last weight update.
	WeightUpdated *time.Time `json:"-"`
}

type NetLoadBalancerNodeCreateOptions struct {
	// The label of the node.
	Label string `json:"label"`
	// The IPv6 address of the node.
	AddressV6 string `json:"address_v6"`
	// The weight of the node.
	Weight int `json:"weight,omitempty"`
}

type NetLoadBalancerNodeUpdateOptions struct {
	// The label of the node.
	Label string `json:"label"`
	// The IPv6 address of the node.
	AddressV6 string `json:"address_v6"`
	// The weight of the node.
	Weight int `json:"weight,omitempty"`
}

type NetLoadBalancerNodeLabelUpdateOptions struct {
	// The label of the node.
	Label string `json:"label"`
}

func (i *NetLoadBalancerNode) UnmarshalJSON(b []byte) error {
	type Mask NetLoadBalancerNode

	p := struct {
		*Mask
		Created       *parseabletime.ParseableTime `json:"created"`
		Updated       *parseabletime.ParseableTime `json:"updated"`
		WeightUpdated *parseabletime.ParseableTime `json:"weight_updated"`
	}{
		Mask: (*Mask)(i),
	}

	if err := json.Unmarshal(b, &p); err != nil {
		return err
	}

	i.Created = (*time.Time)(p.Created)
	i.Updated = (*time.Time)(p.Updated)
	i.WeightUpdated = (*time.Time)(p.WeightUpdated)

	return nil
}

func (i *NetLoadBalancerNode) GetCreateOptions() NetLoadBalancerNodeCreateOptions {
	return NetLoadBalancerNodeCreateOptions{
		Label:     i.Label,
		AddressV6: i.AddressV6,
		Weight:    i.Weight,
	}
}

func (i *NetLoadBalancerNode) GetUpdateOptions() NetLoadBalancerNodeUpdateOptions {
	return NetLoadBalancerNodeUpdateOptions{
		Label: i.Label,
		AddressV6: i.AddressV6,
		Weight: i.Weight,
	}
}

// CreateNetLoadBalancerNode creates a new NetLoadBalancerNode
func (c *Client) CreateNetLoadBalancerNode(ctx context.Context, netloadbalancerID int, listenerID int, opts NetLoadBalancerNodeCreateOptions) (*NetLoadBalancerNode, error) {
	e := formatAPIPath("netloadbalancers/%d/listeners/%d/nodes", netloadbalancerID, listenerID)
	return doPOSTRequest[NetLoadBalancerNode](ctx, c, e, opts)
}

// ListNetLoadBalancerNodes retrieves a list of NetLoadBalancerNodes
func (c *Client) ListNetLoadBalancerNodes(ctx context.Context, netloadbalancerID int, listenerID int, opts *ListOptions) ([]NetLoadBalancerNode, error) {
	e := formatAPIPath("netloadbalancers/%d/listeners/%d/nodes", netloadbalancerID, listenerID)
	return getPaginatedResults[NetLoadBalancerNode](ctx, c, e, opts)
}

// GetNetLoadBalancerNode retrieves a NetLoadBalancerNode by ID
func (c *Client) GetNetLoadBalancerNode(ctx context.Context, netloadbalancerID int, listenerID int, nodeID int) (*NetLoadBalancerNode, error) {
	e := formatAPIPath("netloadbalancers/%d/listeners/%d/nodes/%d", netloadbalancerID, listenerID, nodeID)
	return doGETRequest[NetLoadBalancerNode](ctx, c, e)
}

// UpdateNetLoadBalancerNode updates a NetLoadBalancerNode
func (c *Client) UpdateNetLoadBalancerNode(ctx context.Context, netloadbalancerID int, listenerID int, nodeID int, opts NetLoadBalancerNodeLabelUpdateOptions) (*NetLoadBalancerNode, error) {
	e := formatAPIPath("netloadbalancers/%d/listeners/%d/nodes/%d", netloadbalancerID, listenerID, nodeID)
	return doPUTRequest[NetLoadBalancerNode](ctx, c, e, opts)
}

// DeleteNetLoadBalancerNode deletes a NetLoadBalancerNode
func (c *Client) DeleteNetLoadBalancerNode(ctx context.Context, netloadbalancerID int, listenerID int, nodeID int) error {
	e := formatAPIPath("netloadbalancers/%d/listeners/%d/nodes/%d", netloadbalancerID, listenerID, nodeID)
	return doDELETERequest(ctx, c, e)
}
