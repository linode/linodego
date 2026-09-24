package linodego

import (
	"context"
	"encoding/json"
	"time"

	"github.com/linode/linodego/v2/internal/parseabletime"
)

// NFSQuotaStatus is the lifecycle state of an NFS Quota.
type NFSQuotaStatus string

const (
	// NFSQuotaStatusActive indicates that an NFS Quota is active.
	NFSQuotaStatusActive NFSQuotaStatus = "active"
	// NFSQuotaStatusUpdating indicates that an NFS Quota is being updated.
	NFSQuotaStatusUpdating NFSQuotaStatus = "updating"
	// NFSQuotaStatusError indicates that an NFS Quota reached an error state.
	NFSQuotaStatusError NFSQuotaStatus = "error"
)

// NFSQuotaRuleIdentifierType is the identifier type of an NFS Quota rule.
type NFSQuotaRuleIdentifierType string

const (
	// NFSQuotaRuleIdentifierTypeUID indicates the identifier is a user ID.
	NFSQuotaRuleIdentifierTypeUID NFSQuotaRuleIdentifierType = "uid"
	// NFSQuotaRuleIdentifierTypeUsername indicates the identifier is a username.
	NFSQuotaRuleIdentifierTypeUsername NFSQuotaRuleIdentifierType = "username"
	// NFSQuotaRuleIdentifierTypeGID indicates the identifier is a group ID.
	NFSQuotaRuleIdentifierTypeGID NFSQuotaRuleIdentifierType = "gid"
	// NFSQuotaRuleIdentifierTypeGroupname indicates the identifier is a groupname.
	NFSQuotaRuleIdentifierTypeGroupname NFSQuotaRuleIdentifierType = "groupname"
)

// NFSQuota represents an NFS Quota.
type NFSQuota struct {
	ID                int                 `json:"id"`
	FilesystemID      int                 `json:"filesystem_id"`
	Path              string              `json:"path"`
	IsRoot            bool                `json:"is_root"`
	MaxCapacityBytes  *int64              `json:"max_capacity_bytes"`
	UsedCapacityBytes *int64              `json:"used_capacity_bytes"`
	QuotaExceeded     bool                `json:"quota_exceeded"`
	MaxFileCount      *int64              `json:"max_file_count"`
	UsedFileCount     *int64              `json:"used_file_count"`
	Status            NFSQuotaStatus      `json:"status"`
	Created           *time.Time          `json:"-"`
	Updated           *time.Time          `json:"-"`
	CollectedAt       *time.Time          `json:"-"`
	UserGroupConfig   *NFSUserGroupConfig `json:"user_group_config"`
}

// NFSUserGroupConfig represents user and group quota configuration for an NFS Quota.
type NFSUserGroupConfig struct {
	DefaultUserLimit  *NFSCapacityLimit           `json:"default_user_limit"`
	DefaultGroupLimit *NFSCapacityLimit           `json:"default_group_limit"`
	UserLimits        []NFSIdentifiedLimit        `json:"user_limits"`
	GroupLimits       []NFSIdentifiedLimit        `json:"group_limits"`
	UserGroupLimits   []NFSUserGroupCombinedLimit `json:"user_group_limits"`
}

// NFSCapacityLimit represents capacity and file count limits for an NFS Quota.
type NFSCapacityLimit struct {
	MaxCapacityBytes *int64 `json:"max_capacity_bytes"`
	MaxFileCount     *int64 `json:"max_file_count"`
}

// NFSIdentifiedLimit represents capacity and file count limits for a specific user or group.
type NFSIdentifiedLimit struct {
	IdentifierType   NFSQuotaRuleIdentifierType `json:"identifier_type"`
	Identifier       string                     `json:"identifier"`
	MaxCapacityBytes *int64                     `json:"max_capacity_bytes"`
	MaxFileCount     *int64                     `json:"max_file_count"`
}

// NFSUserGroupCombinedLimit represents capacity and file count limits for a combined user and group.
type NFSUserGroupCombinedLimit struct {
	User             NFSUserGroupLimitIdentity `json:"user"`
	Group            NFSUserGroupLimitIdentity `json:"group"`
	MaxCapacityBytes *int64                    `json:"max_capacity_bytes"`
	MaxFileCount     *int64                    `json:"max_file_count"`
}

// NFSUserGroupLimitIdentity identifies a user or group in NFS quota configuration.
type NFSUserGroupLimitIdentity struct {
	IdentifierType NFSQuotaRuleIdentifierType `json:"identifier_type"`
	Identifier     string                     `json:"identifier"`
}

// NFSQuotaCreateOptions contains fields accepted when creating an NFS Quota.
type NFSQuotaCreateOptions struct {
	Path             string                           `json:"path"`
	MaxCapacityBytes int64                            `json:"max_capacity_bytes"`
	MaxFileCount     int64                            `json:"max_file_count"`
	UserGroupConfig  *NFSUserGroupConfigUpdateOptions `json:"user_group_config,omitzero"`
}

// NFSQuotaUpdateOptions contains fields accepted when updating an NFS Quota.
type NFSQuotaUpdateOptions struct {
	MaxCapacityBytes *int64                           `json:"max_capacity_bytes,omitzero"`
	MaxFileCount     *int64                           `json:"max_file_count,omitzero"`
	UserGroupConfig  *NFSUserGroupConfigUpdateOptions `json:"user_group_config,omitzero"`
}

// NFSUserGroupConfigUpdateOptions contains fields for updating user and group quota configuration.
type NFSUserGroupConfigUpdateOptions struct {
	DefaultUserLimit  *NFSCapacityLimitUpdateOptions            `json:"default_user_limit,omitzero"`
	DefaultGroupLimit *NFSCapacityLimitUpdateOptions            `json:"default_group_limit,omitzero"`
	UserLimits        *[]NFSIdentifiedLimitUpdateOptions        `json:"user_limits,omitzero"`
	GroupLimits       *[]NFSIdentifiedLimitUpdateOptions        `json:"group_limits,omitzero"`
	UserGroupLimits   *[]NFSUserGroupCombinedLimitUpdateOptions `json:"user_group_limits,omitzero"`
}

// NFSCapacityLimitUpdateOptions contains fields for updating capacity and file count limits.
type NFSCapacityLimitUpdateOptions struct {
	MaxCapacityBytes *int64 `json:"max_capacity_bytes,omitzero"`
	MaxFileCount     *int64 `json:"max_file_count,omitzero"`
}

// NFSIdentifiedLimitUpdateOptions contains fields for updating limits for a specific user or group.
type NFSIdentifiedLimitUpdateOptions struct {
	IdentifierType   NFSQuotaRuleIdentifierType `json:"identifier_type"`
	Identifier       string                     `json:"identifier"`
	MaxCapacityBytes *int64                     `json:"max_capacity_bytes,omitzero"`
	MaxFileCount     *int64                     `json:"max_file_count,omitzero"`
}

// NFSUserGroupCombinedLimitUpdateOptions contains fields for updating limits for a combined user and group.
type NFSUserGroupCombinedLimitUpdateOptions struct {
	User             NFSUserGroupLimitIdentity `json:"user"`
	Group            NFSUserGroupLimitIdentity `json:"group"`
	MaxCapacityBytes *int64                    `json:"max_capacity_bytes,omitzero"`
	MaxFileCount     *int64                    `json:"max_file_count,omitzero"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (n *NFSQuota) UnmarshalJSON(b []byte) error {
	type Mask NFSQuota

	p := struct {
		*Mask

		Created     *parseabletime.ParseableTime `json:"created"`
		Updated     *parseabletime.ParseableTime `json:"updated"`
		CollectedAt *parseabletime.ParseableTime `json:"collected_at"`
	}{
		Mask: (*Mask)(n),
	}

	if err := json.Unmarshal(b, &p); err != nil {
		return err
	}

	n.Created = (*time.Time)(p.Created)
	n.Updated = (*time.Time)(p.Updated)
	n.CollectedAt = (*time.Time)(p.CollectedAt)

	return nil
}

// ListNFSQuotas lists NFS Quotas for NFS Filesystem. This endpoint requires the v4beta API.
func (c *Client) ListNFSQuotas(ctx context.Context, spaceID int, filesystemID int, opts *ListOptions) ([]NFSQuota, error) {
	return getPaginatedResults[NFSQuota](ctx, c, formatAPIPath("nfs/spaces/%d/filesystems/%d/quotas", spaceID, filesystemID), opts)
}

// CreateNFSQuota creates an NFS Quota for an NFS Filesystem. This endpoint requires the v4beta API.
func (c *Client) CreateNFSQuota(ctx context.Context, spaceID int, filesystemID int, opts NFSQuotaCreateOptions) (*NFSQuota, error) {
	return doPOSTRequest[NFSQuota](ctx, c, formatAPIPath("nfs/spaces/%d/filesystems/%d/quotas", spaceID, filesystemID), opts)
}

// GetNFSQuota returns an NFS Quota by its ID. This endpoint requires the v4beta API.
func (c *Client) GetNFSQuota(ctx context.Context, spaceID int, filesystemID int, quotaID int) (*NFSQuota, error) {
	return doGETRequest[NFSQuota](ctx, c, formatAPIPath("nfs/spaces/%d/filesystems/%d/quotas/%d", spaceID, filesystemID, quotaID))
}

// UpdateNFSQuota updates an NFS Quota for an NFS Filesystem. This endpoint requires the v4beta API.
func (c *Client) UpdateNFSQuota(ctx context.Context, spaceID int, filesystemID int, quotaID int, opts NFSQuotaUpdateOptions) (*NFSQuota, error) {
	return doPUTRequest[NFSQuota](ctx, c, formatAPIPath("nfs/spaces/%d/filesystems/%d/quotas/%d", spaceID, filesystemID, quotaID), opts)
}

// DeleteNFSQuota deletes an NFS Quota for an NFS Filesystem. This endpoint requires the v4beta API.
func (c *Client) DeleteNFSQuota(ctx context.Context, spaceID int, filesystemID int, quotaID int) error {
	return doDELETERequest(ctx, c, formatAPIPath("nfs/spaces/%d/filesystems/%d/quotas/%d", spaceID, filesystemID, quotaID))
}
