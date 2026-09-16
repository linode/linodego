package linodego

import (
	"context"
	"encoding/json"
	"time"

	"github.com/linode/linodego/v2/internal/parseabletime"
)

// NFSMTLSMode is the mutual TLS mode for an NFS Space access policy.
type NFSMTLSMode string

const (
	// NFSMTLSModeRequired requires mutual TLS.
	NFSMTLSModeRequired NFSMTLSMode = "required"
	// NFSMTLSModeOptional permits mutual TLS.
	NFSMTLSModeOptional NFSMTLSMode = "optional"
	// NFSMTLSModeDisabled disables mutual TLS.
	NFSMTLSModeDisabled NFSMTLSMode = "disabled"
)

// NFSAccessPolicyStatus is the lifecycle state of an NFS access policy.
type NFSAccessPolicyStatus string

const (
	// NFSAccessPolicyStatusCreating indicates that an NFS access policy is being created.
	NFSAccessPolicyStatusCreating NFSAccessPolicyStatus = "creating"
	// NFSAccessPolicyStatusActive indicates that an NFS access policy is active.
	NFSAccessPolicyStatusActive NFSAccessPolicyStatus = "active"
	// NFSAccessPolicyStatusUpdating indicates that an NFS access policy is being updated.
	NFSAccessPolicyStatusUpdating NFSAccessPolicyStatus = "updating"
	// NFSAccessPolicyStatusDeleting indicates that an NFS access policy is being deleted.
	NFSAccessPolicyStatusDeleting NFSAccessPolicyStatus = "deleting"
	// NFSAccessPolicyStatusError indicates that an NFS access policy reached an error state.
	NFSAccessPolicyStatusError NFSAccessPolicyStatus = "error"
)

// NFSSquashPolicy is the user-ID squash policy for an NFS Filesystem access policy.
type NFSSquashPolicy string

const (
	// NFSSquashPolicyNone disables user-ID squashing.
	NFSSquashPolicyNone NFSSquashPolicy = "none"
	// NFSSquashPolicyRootSquash squashes the root user ID.
	NFSSquashPolicyRootSquash NFSSquashPolicy = "root_squash"
	// NFSSquashPolicyAllSquash squashes all user IDs.
	NFSSquashPolicyAllSquash NFSSquashPolicy = "all_squash"
)

// NFSSpaceAccessPolicyVPC represents a VPC allowed by an NFS Space access policy.
type NFSSpaceAccessPolicyVPC struct {
	ID      int                             `json:"id"`
	Label   *string                         `json:"label"`
	URL     *string                         `json:"url"`
	Range   *string                         `json:"range"`
	Subnets []NFSSpaceAccessPolicyVPCSubnet `json:"subnets"`
}

// NFSSpaceAccessPolicyVPCSubnet represents a subnet allowed by an NFS Space access policy.
type NFSSpaceAccessPolicyVPCSubnet struct {
	ID    int     `json:"id"`
	Label *string `json:"label"`
	URL   *string `json:"url"`
	Range *string `json:"range"`
}

// NFSSpaceAccessPolicyVPCOptions identifies a VPC and its optional subnets for an NFS Space access policy update.
type NFSSpaceAccessPolicyVPCOptions struct {
	ID      int   `json:"id"`
	Subnets []int `json:"subnets,omitzero"`
}

// NFSFilesystemAccessPolicyLinode represents a Linode allowed by an NFS Filesystem access policy.
type NFSFilesystemAccessPolicyLinode struct {
	ID    int     `json:"id"`
	Label *string `json:"label"`
	URL   *string `json:"url"`
	IPv6  *string `json:"ipv6"`
}

// NFSSpaceAccessPolicy represents an NFS Space access policy.
type NFSSpaceAccessPolicy struct {
	SpaceID    int                       `json:"space_id"`
	Label      string                    `json:"label"`
	Enabled    bool                      `json:"enabled"`
	VPCACL     []NFSSpaceAccessPolicyVPC `json:"vpc_acl"`
	MTLSCACert *string                   `json:"mtls_ca_cert"`
	MTLSMode   NFSMTLSMode               `json:"mtls_mode"`
	Status     NFSAccessPolicyStatus     `json:"status"`
	Created    *time.Time                `json:"-"`
	Updated    *time.Time                `json:"-"`
}

// NFSSpaceAccessPolicyUpdateOptions contains fields accepted when updating an NFS Space access policy.
type NFSSpaceAccessPolicyUpdateOptions struct {
	Label      *string                           `json:"label,omitzero"`
	Enabled    *bool                             `json:"enabled,omitzero"`
	VPCs       *[]NFSSpaceAccessPolicyVPCOptions `json:"vpcs,omitzero"`
	MTLSCACert *string                           `json:"mtls_ca_cert,omitzero"`
	MTLSMode   *NFSMTLSMode                      `json:"mtls_mode,omitzero"`
}

// NFSFilesystemAccessPolicy represents an NFS Filesystem access policy.
type NFSFilesystemAccessPolicy struct {
	FilesystemID int                               `json:"filesystem_id"`
	Label        string                            `json:"label"`
	Enabled      bool                              `json:"enabled"`
	LinodeACL    []NFSFilesystemAccessPolicyLinode `json:"linode_acl"`
	SquashPolicy NFSSquashPolicy                   `json:"squash_policy"`
	Protocols    []NFSProtocolVersion              `json:"protocols"`
	Status       NFSAccessPolicyStatus             `json:"status"`
	Created      *time.Time                        `json:"-"`
	Updated      *time.Time                        `json:"-"`
}

// NFSFilesystemAccessPolicyUpdateOptions contains fields accepted when updating an NFS Filesystem access policy.
type NFSFilesystemAccessPolicyUpdateOptions struct {
	Label        *string               `json:"label,omitzero"`
	Enabled      *bool                 `json:"enabled,omitzero"`
	LinodeIDs    *[]int                `json:"linode_ids,omitzero"`
	SquashPolicy *NFSSquashPolicy      `json:"squash_policy,omitzero"`
	Protocols    *[]NFSProtocolVersion `json:"protocols,omitzero"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (n *NFSSpaceAccessPolicy) UnmarshalJSON(b []byte) error {
	type mask NFSSpaceAccessPolicy

	payload := struct {
		*mask

		Created *parseabletime.ParseableTime `json:"created"`
		Updated *parseabletime.ParseableTime `json:"updated"`
	}{
		mask: (*mask)(n),
	}

	if err := json.Unmarshal(b, &payload); err != nil {
		return err
	}

	n.Created = (*time.Time)(payload.Created)
	n.Updated = (*time.Time)(payload.Updated)

	return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (n *NFSFilesystemAccessPolicy) UnmarshalJSON(b []byte) error {
	type mask NFSFilesystemAccessPolicy

	payload := struct {
		*mask

		Created *parseabletime.ParseableTime `json:"created"`
		Updated *parseabletime.ParseableTime `json:"updated"`
	}{
		mask: (*mask)(n),
	}

	if err := json.Unmarshal(b, &payload); err != nil {
		return err
	}

	n.Created = (*time.Time)(payload.Created)
	n.Updated = (*time.Time)(payload.Updated)

	return nil
}

// GetNFSSpaceAccessPolicy returns an NFS Space access policy. This endpoint requires the v4beta API.
func (c *Client) GetNFSSpaceAccessPolicy(ctx context.Context, spaceID int) (*NFSSpaceAccessPolicy, error) {
	return doGETRequest[NFSSpaceAccessPolicy](ctx, c, formatAPIPath("nfs/spaces/%d/access-policy", spaceID))
}

// UpdateNFSSpaceAccessPolicy updates an NFS Space access policy. This endpoint requires the v4beta API.
func (c *Client) UpdateNFSSpaceAccessPolicy(
	ctx context.Context,
	spaceID int,
	opts NFSSpaceAccessPolicyUpdateOptions,
) (*NFSSpaceAccessPolicy, error) {
	return doPUTRequest[NFSSpaceAccessPolicy](ctx, c, formatAPIPath("nfs/spaces/%d/access-policy", spaceID), opts)
}

// GetNFSFilesystemAccessPolicy returns an NFS Filesystem access policy. This endpoint requires the v4beta API.
func (c *Client) GetNFSFilesystemAccessPolicy(
	ctx context.Context,
	spaceID int,
	filesystemID int,
) (*NFSFilesystemAccessPolicy, error) {
	return doGETRequest[NFSFilesystemAccessPolicy](
		ctx,
		c,
		formatAPIPath("nfs/spaces/%d/filesystems/%d/access-policy", spaceID, filesystemID),
	)
}

// UpdateNFSFilesystemAccessPolicy updates an NFS Filesystem access policy. This endpoint requires the v4beta API.
func (c *Client) UpdateNFSFilesystemAccessPolicy(
	ctx context.Context,
	spaceID int,
	filesystemID int,
	opts NFSFilesystemAccessPolicyUpdateOptions,
) (*NFSFilesystemAccessPolicy, error) {
	return doPUTRequest[NFSFilesystemAccessPolicy](
		ctx,
		c,
		formatAPIPath("nfs/spaces/%d/filesystems/%d/access-policy", spaceID, filesystemID),
		opts,
	)
}
