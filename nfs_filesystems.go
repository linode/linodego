package linodego

import (
	"context"
	"encoding/json"
	"time"

	"github.com/linode/linodego/v2/internal/parseabletime"
)

// NFSProtocolVersion is an NFS protocol version accepted by an NFS Filesystem.
type NFSProtocolVersion string

const (
	// NFSProtocolVersionV4 is NFS version 4.
	NFSProtocolVersionV4 NFSProtocolVersion = "nfsv4"
)

// NFSFilesystemStatus is the lifecycle state of an NFS Filesystem.
type NFSFilesystemStatus string

const (
	// NFSFilesystemStatusCreating indicates that an NFS Filesystem is being created.
	NFSFilesystemStatusCreating NFSFilesystemStatus = "creating"
	// NFSFilesystemStatusActive indicates that an NFS Filesystem is active.
	NFSFilesystemStatusActive NFSFilesystemStatus = "active"
	// NFSFilesystemStatusUpdating indicates that an NFS Filesystem is being updated.
	NFSFilesystemStatusUpdating NFSFilesystemStatus = "updating"
	// NFSFilesystemStatusDeleting indicates that an NFS Filesystem is being deleted.
	NFSFilesystemStatusDeleting NFSFilesystemStatus = "deleting"
	// NFSFilesystemStatusError indicates that an NFS Filesystem reached an error state.
	NFSFilesystemStatusError NFSFilesystemStatus = "error"
)

// NFSFilesystem represents an NFS Filesystem.
type NFSFilesystem struct {
	ID                 int                  `json:"id"`
	SpaceID            int                  `json:"space_id"`
	Label              string               `json:"label"`
	Region             string               `json:"region"`
	ProtocolVersions   []NFSProtocolVersion `json:"protocol_versions"`
	Status             NFSFilesystemStatus  `json:"status"`
	MountTargetIPs     []string             `json:"mount_target_ips"`
	MountTargetFQDN    *string              `json:"mount_target_fqdn"`
	MaxCapacityBytes   int64                `json:"max_capacity_bytes"`
	SnapshotUsageBytes *int64               `json:"snapshot_usage_bytes"`
	LDAPConfigID       *string              `json:"ldap_config_id"`
	SourceSnapshotID   *int                 `json:"source_snapshot_id"`
	Created            *time.Time           `json:"-"`
	Updated            *time.Time           `json:"-"`
	Tags               []string             `json:"tags"`
	Stats              NFSFilesystemStats   `json:"stats"`
}

// NFSFilesystemCreateOptions contains fields accepted when creating an NFS Filesystem.
type NFSFilesystemCreateOptions struct {
	Label            string                `json:"label"`
	Region           string                `json:"region"`
	MaxCapacityBytes int64                 `json:"max_capacity_bytes"`
	ProtocolVersions *[]NFSProtocolVersion `json:"protocol_versions,omitzero"`
	Tags             *[]string             `json:"tags,omitzero"`
}

// NFSFilesystemUpdateOptions contains fields accepted when updating an NFS Filesystem.
type NFSFilesystemUpdateOptions struct {
	Label            *string   `json:"label,omitzero"`
	MaxCapacityBytes *int64    `json:"max_capacity_bytes,omitzero"`
	Tags             *[]string `json:"tags,omitzero"`
}

// NFSFilesystemStats represents collected NFS Filesystem usage statistics.
type NFSFilesystemStats struct {
	UsedCapacityBytes *int64     `json:"used_capacity_bytes"`
	UsedFileCount     *int64     `json:"used_file_count"`
	CollectedAt       *time.Time `json:"-"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (n *NFSFilesystem) UnmarshalJSON(b []byte) error {
	type mask NFSFilesystem

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
func (n *NFSFilesystemStats) UnmarshalJSON(b []byte) error {
	type mask NFSFilesystemStats

	payload := struct {
		*mask

		CollectedAt *parseabletime.ParseableTime `json:"collected_at"`
	}{
		mask: (*mask)(n),
	}

	if err := json.Unmarshal(b, &payload); err != nil {
		return err
	}

	n.CollectedAt = (*time.Time)(payload.CollectedAt)

	return nil
}

// GetCreateOptions returns the NFS Filesystem fields accepted by CreateNFSFilesystem.
func (n NFSFilesystem) GetCreateOptions() NFSFilesystemCreateOptions {
	result := NFSFilesystemCreateOptions{
		Label:            n.Label,
		Region:           n.Region,
		MaxCapacityBytes: n.MaxCapacityBytes,
	}

	if n.ProtocolVersions != nil {
		result.ProtocolVersions = Pointer(n.ProtocolVersions)
	}

	if n.Tags != nil {
		result.Tags = Pointer(n.Tags)
	}

	return result
}

// GetUpdateOptions returns the NFS Filesystem fields accepted by UpdateNFSFilesystem.
func (n NFSFilesystem) GetUpdateOptions() NFSFilesystemUpdateOptions {
	result := NFSFilesystemUpdateOptions{
		Label:            Pointer(n.Label),
		MaxCapacityBytes: Pointer(n.MaxCapacityBytes),
	}

	if n.Tags != nil {
		result.Tags = Pointer(n.Tags)
	}

	return result
}

// ListNFSFilesystems lists NFS Filesystems in an NFS Space. This endpoint requires the v4beta API.
func (c *Client) ListNFSFilesystems(ctx context.Context, spaceID int, opts *ListOptions) ([]NFSFilesystem, error) {
	return getPaginatedResults[NFSFilesystem](ctx, c, formatAPIPath("nfs/spaces/%d/filesystems", spaceID), opts)
}

// GetNFSFilesystem returns an NFS Filesystem in an NFS Space. This endpoint requires the v4beta API.
func (c *Client) GetNFSFilesystem(ctx context.Context, spaceID int, filesystemID int) (*NFSFilesystem, error) {
	return doGETRequest[NFSFilesystem](
		ctx,
		c,
		formatAPIPath("nfs/spaces/%d/filesystems/%d", spaceID, filesystemID),
	)
}

// GetNFSFilesystemByID returns an NFS Filesystem by its ID. This endpoint requires the v4beta API.
func (c *Client) GetNFSFilesystemByID(ctx context.Context, filesystemID int) (*NFSFilesystem, error) {
	return doGETRequest[NFSFilesystem](ctx, c, formatAPIPath("nfs/filesystems/%d", filesystemID))
}

// CreateNFSFilesystem creates an NFS Filesystem in an NFS Space. This endpoint requires the v4beta API.
func (c *Client) CreateNFSFilesystem(
	ctx context.Context,
	spaceID int,
	opts NFSFilesystemCreateOptions,
) (*NFSFilesystem, error) {
	return doPOSTRequest[NFSFilesystem](ctx, c, formatAPIPath("nfs/spaces/%d/filesystems", spaceID), opts)
}

// UpdateNFSFilesystem updates an NFS Filesystem in an NFS Space. This endpoint requires the v4beta API.
func (c *Client) UpdateNFSFilesystem(
	ctx context.Context,
	spaceID int,
	filesystemID int,
	opts NFSFilesystemUpdateOptions,
) (*NFSFilesystem, error) {
	return doPUTRequest[NFSFilesystem](
		ctx,
		c,
		formatAPIPath("nfs/spaces/%d/filesystems/%d", spaceID, filesystemID),
		opts,
	)
}

// DeleteNFSFilesystem deletes an NFS Filesystem in an NFS Space. This endpoint requires the v4beta API.
func (c *Client) DeleteNFSFilesystem(ctx context.Context, spaceID int, filesystemID int) error {
	return doDELETERequest(ctx, c, formatAPIPath("nfs/spaces/%d/filesystems/%d", spaceID, filesystemID))
}
