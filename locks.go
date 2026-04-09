package linodego

import (
	"context"
)

// LockType represents the type of lock that can be applied to a resource
type LockType string

// LockType enums
const (
	LockTypeCannotDelete                 LockType = "cannot_delete"
	LockTypeCannotDeleteWithSubresources LockType = "cannot_delete_with_subresources"
)

// LockedEntity represents the entity that is locked
type LockedEntity struct {
	ID    int        `json:"id"`
	Type  EntityType `json:"type"`
	Label string     `json:"label"`
	URL   string     `json:"url"`
}

// Lock represents a resource lock
type Lock struct {
	ID       int          `json:"id"`
	LockType LockType     `json:"lock_type"`
	Entity   LockedEntity `json:"entity"`
}

// LockCreateOptions fields are those accepted by CreateLock
type LockCreateOptions struct {
	EntityType EntityType `json:"entity_type"`
	EntityID   int        `json:"entity_id"`
	LockType   LockType   `json:"lock_type"`
}

// ListLocks returns a paginated list of Locks
func (c *Client) ListLocks(ctx context.Context, opts *ListOptions) ([]Lock, error) {
	return getPaginatedResults[Lock](ctx, c, "locks", opts)
}

// GetLock gets a single Lock with the provided ID
func (c *Client) GetLock(ctx context.Context, lockID int) (*Lock, error) {
	e := formatAPIPath("locks/%d", lockID)
	return doGETRequest[Lock](ctx, c, e)
}

// CreateLock creates a lock for a resource
func (c *Client) CreateLock(ctx context.Context, opts LockCreateOptions) (*Lock, error) {
	return doPOSTRequest[Lock](ctx, c, "locks", opts)
}

// DeleteLock deletes a single Lock with the provided ID
func (c *Client) DeleteLock(ctx context.Context, lockID int) error {
	e := formatAPIPath("locks/%d", lockID)
	return doDELETERequest(ctx, c, e)
}
