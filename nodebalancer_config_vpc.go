package linodego

import (
	"context"
)

// NodeBalancerVpcConfig objects represent a VPC config for a NodeBalancer
type NodeBalancerVpcConfig struct {
	ID             int    `json:"id"`
	IPv4Range      string `json:"ipv4_range"`
	IPv6Range      string `json:"ipv6_range"`
	NodeBalancerID int    `json:"nodebalancer_id"`
	SubnetID       int    `json:"subnet_id"`
	VPCID          int    `json:"vpc_id"`
}

// ListNodeBalancerVpcConfigs lists NodeBalancer VPC configs
func (c *Client) ListNodeBalancerVpcConfigs(ctx context.Context, nodebalancerID int, opts *ListOptions) ([]NodeBalancerVpcConfig, error) {
	return getPaginatedResults[NodeBalancerVpcConfig](ctx, c, formatAPIPath("nodebalancers/%d/vpcs", nodebalancerID), opts)
}

// GetNodeBalancerVpcConfig gets the NodeBalancer VPC config with the specified id
func (c *Client) GetNodeBalancerVpcConfig(ctx context.Context, nodebalancerID int, vpcID int) (*NodeBalancerVpcConfig, error) {
	e := formatAPIPath("nodebalancers/%d/vpcs/%d", nodebalancerID, vpcID)
	return doGETRequest[NodeBalancerVpcConfig](ctx, c, e)
}
