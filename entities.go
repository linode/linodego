package linodego

import "context"

type LinodeEntity struct {
	ID    int    `json:"id"`
	Label string `json:"label"`
	Type  string `json:"type"`
}

func (c *Client) ListEntities(ctx context.Context, opts *ListOptions) ([]LinodeEntity, error) {
	return getPaginatedResults[LinodeEntity](ctx, c, "entities", opts)
}
