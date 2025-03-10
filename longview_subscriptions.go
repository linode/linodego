package linodego

import (
	"context"
)

// LongviewSubscription represents a LongviewSubscription object
type LongviewSubscription struct {
	ID              string       `json:"id"`
	Label           string       `json:"label"`
	ClientsIncluded int          `json:"clients_included"`
	Price           *LinodePrice `json:"price"`
	// UpdatedStr string `json:"updated"`
	// Updated *time.Time `json:"-"`
}

// ListLongviewSubscriptions lists LongviewSubscriptions
func (c *Client) ListLongviewSubscriptions(ctx context.Context, opts *ListOptions) ([]LongviewSubscription, error) {
	return getPaginatedResults[LongviewSubscription](ctx, c, "longview/subscriptions", opts)
}

// GetLongviewSubscription gets the template with the provided ID
func (c *Client) GetLongviewSubscription(ctx context.Context, templateID string) (*LongviewSubscription, error) {
	e := formatAPIPath("longview/subscriptions/%s", templateID)
	return doGETRequest[LongviewSubscription](ctx, c, e)
}
