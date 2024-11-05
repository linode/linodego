package linodego

import (
	"context"
)

// ListInstanceNodeBalancers lists NodeBalancers
func (c *Client) ListInstanceNodeBalancers(ctx context.Context, linodeID int, opts *ListOptions) ([]NodeBalancer, error) {
	return getPaginatedResults[NodeBalancer](ctx, c, formatAPIPath("linode/instances/%d/nodebalancers", linodeID), opts)
}
