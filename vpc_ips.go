package linodego

import (
	"context"
)

// ListVPCIPAddresses gets the list of all IP addresses of all VPCs in the Linode account.
func (c *Client) ListVPCIPAddresses(ctx context.Context, opts *ListOptions) ([]VPCIP, error) {
	return getPaginatedResults[VPCIP](ctx, c, "vpcs/ips", opts)
}
