package linodego

import (
	"context"
	"encoding/json"
	"time"

	"github.com/linode/linodego/v2/internal/parseabletime"
)

// NFSSpaceStatus is the lifecycle state of an NFS Space.
type NFSSpaceStatus string

const (
	// NFSSpaceStatusCreating indicates that an NFS Space is being created.
	NFSSpaceStatusCreating NFSSpaceStatus = "creating"
	// NFSSpaceStatusActive indicates that an NFS Space is active.
	NFSSpaceStatusActive NFSSpaceStatus = "active"
	// NFSSpaceStatusDeleting indicates that an NFS Space is being deleted.
	NFSSpaceStatusDeleting NFSSpaceStatus = "deleting"
	// NFSSpaceStatusError indicates that an NFS Space reached an error state.
	NFSSpaceStatusError NFSSpaceStatus = "error"
)

// NFSSpace represents an NFS Space.
type NFSSpace struct {
	ID                int            `json:"id"`
	Label             string         `json:"label"`
	Description       *string        `json:"description"`
	Status            NFSSpaceStatus `json:"status"`
	MaxCapacityBytes  int64          `json:"max_capacity_bytes"`
	UsedCapacityBytes int64          `json:"used_capacity_bytes"`
	Created           *time.Time     `json:"-"`
	Updated           *time.Time     `json:"-"`
	Tags              []string       `json:"tags"`
}

// NFSSpaceCreateOptions contains fields accepted when creating an NFS Space.
type NFSSpaceCreateOptions struct {
	Label       string    `json:"label"`
	Description **string  `json:"description,omitzero"`
	Tags        *[]string `json:"tags,omitzero"`
}

// NFSSpaceUpdateOptions contains fields accepted when updating an NFS Space.
type NFSSpaceUpdateOptions struct {
	Label       *string   `json:"label,omitzero"`
	Description **string  `json:"description,omitzero"`
	Tags        *[]string `json:"tags,omitzero"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (n *NFSSpace) UnmarshalJSON(b []byte) error {
	type mask NFSSpace

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

// GetCreateOptions returns the NFS Space fields accepted by CreateNFSSpace.
func (n NFSSpace) GetCreateOptions() NFSSpaceCreateOptions {
	result := NFSSpaceCreateOptions{
		Label:       n.Label,
		Description: Pointer(n.Description),
	}

	if n.Tags != nil {
		result.Tags = Pointer(n.Tags)
	}

	return result
}

// GetUpdateOptions returns the NFS Space fields accepted by UpdateNFSSpace.
func (n NFSSpace) GetUpdateOptions() NFSSpaceUpdateOptions {
	result := NFSSpaceUpdateOptions{
		Label:       Pointer(n.Label),
		Description: Pointer(n.Description),
	}

	if n.Tags != nil {
		result.Tags = Pointer(n.Tags)
	}

	return result
}

// ListNFSSpaces lists NFS Spaces. This endpoint requires the v4beta API.
func (c *Client) ListNFSSpaces(ctx context.Context, opts *ListOptions) ([]NFSSpace, error) {
	return getPaginatedResults[NFSSpace](ctx, c, "nfs/spaces", opts)
}

// GetNFSSpace returns an NFS Space. This endpoint requires the v4beta API.
func (c *Client) GetNFSSpace(ctx context.Context, spaceID int) (*NFSSpace, error) {
	return doGETRequest[NFSSpace](ctx, c, formatAPIPath("nfs/spaces/%d", spaceID))
}

// CreateNFSSpace creates an NFS Space. This endpoint requires the v4beta API.
func (c *Client) CreateNFSSpace(ctx context.Context, opts NFSSpaceCreateOptions) (*NFSSpace, error) {
	return doPOSTRequest[NFSSpace](ctx, c, "nfs/spaces", opts)
}

// UpdateNFSSpace updates an NFS Space. This endpoint requires the v4beta API.
func (c *Client) UpdateNFSSpace(ctx context.Context, spaceID int, opts NFSSpaceUpdateOptions) (*NFSSpace, error) {
	return doPUTRequest[NFSSpace](ctx, c, formatAPIPath("nfs/spaces/%d", spaceID), opts)
}

// DeleteNFSSpace deletes an NFS Space. This endpoint requires the v4beta API.
func (c *Client) DeleteNFSSpace(ctx context.Context, spaceID int) error {
	return doDELETERequest(ctx, c, formatAPIPath("nfs/spaces/%d", spaceID))
}
