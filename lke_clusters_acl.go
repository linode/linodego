package linodego

import "context"

// LKEClusterControlPlaneACLAddresses describes the
// allowed IP ranges for an LKE cluster's control plane.
type LKEClusterControlPlaneACLAddresses struct {
	IPv4 []string `json:"ipv4"`
	IPv6 []string `json:"ipv6"`
}

// LKEClusterControlPlaneACL describes the ACL configuration
// for an LKE cluster's control plane.
type LKEClusterControlPlaneACL struct {
	Enabled   bool                               `json:"enabled"`
	Addresses LKEClusterControlPlaneACLAddresses `json:"addresses"`
}

// LKEClusterControlPlaneACLUpdateOptions represents the options
// available when updating the ACL configuration of an LKE cluster's
// control plane.
type LKEClusterControlPlaneACLUpdateOptions struct {
	ACL *LKEClusterControlPlaneACL `json:"acl"`
}

// GetLKEClusterControlPlaneACL gets the ACL configuration for the
// given cluster's control plane.
func (c *Client) GetLKEClusterControlPlaneACL(ctx context.Context, clusterID int) (*LKEClusterControlPlaneACL, error) {
	return doGETRequest[LKEClusterControlPlaneACL](
		ctx,
		c,
		formatAPIPath("lke/clusters/%d/control_plane_acl", clusterID),
	)
}

// UpdateLKEClusterControlPlaneACL gets the ACL configuration for the
// given cluster's control plane.
func (c *Client) UpdateLKEClusterControlPlaneACL(
	ctx context.Context,
	clusterID int,
	opts LKEClusterControlPlaneACLUpdateOptions,
) (*LKEClusterControlPlaneACL, error) {
	return doPUTRequest[LKEClusterControlPlaneACL](
		ctx,
		c,
		formatAPIPath("lke/clusters/%d/control_plane_acl", clusterID),
		opts,
	)
}

// DeleteLKEClusterControlPlaneACL deletes the ACL configuration for the
// given cluster's control plane.
func (c *Client) DeleteLKEClusterControlPlaneACL(
	ctx context.Context,
	clusterID int,
) error {
	return doDELETERequest(
		ctx,
		c,
		formatAPIPath("lke/clusters/%d/control_plane_acl", clusterID),
	)
}
