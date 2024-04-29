package linodego

import "context"

// LKEClusterControlPlaneACLAddresses describes the
// allowed IP ranges for an LKE cluster's control plane.
type LKEClusterControlPlaneACLAddresses struct {
	IPv4 []string `json:"ipv4,omitempty"`
	IPv6 []string `json:"ipv6,omitempty"`
}

// LKEClusterControlPlaneACL describes the ACL configuration
// for an LKE cluster's control plane.
type LKEClusterControlPlaneACL struct {
	Enabled   bool                               `json:"enabled,omitempty"`
	Addresses LKEClusterControlPlaneACLAddresses `json:"addresses,omitempty"`
}

// LKEClusterControlPlaneACLUpdateOptions represents the options
// available when updating the ACL configuration of an LKE cluster's
// control plane.
type LKEClusterControlPlaneACLUpdateOptions struct {
	ACL LKEClusterControlPlaneACL `json:"acl"`
}

// LKEClusterControlPlaneACLResponse represents the response structure
// for the Client.GetLKEClusterControlPlaneACL(...) method.
type LKEClusterControlPlaneACLResponse struct {
	ACL LKEClusterControlPlaneACL `json:"acl"`
}

// GetLKEClusterControlPlaneACL gets the ACL configuration for the
// given cluster's control plane.
func (c *Client) GetLKEClusterControlPlaneACL(ctx context.Context, clusterID int) (*LKEClusterControlPlaneACLResponse, error) {
	return doGETRequest[LKEClusterControlPlaneACLResponse](
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
) (*LKEClusterControlPlaneACLResponse, error) {
	return doPUTRequest[LKEClusterControlPlaneACLResponse](
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
