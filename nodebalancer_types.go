package linodego

import (
	"context"
)

// NodeBalancerType represents a single valid NodeBalancer type.
type NodeBalancerType struct {
	baseType
}

// ListNodeBalancerTypes lists NodeBalancer types. This endpoint is cached by default.
func (c *Client) ListNodeBalancerTypes(ctx context.Context, opts *ListOptions) ([]NodeBalancerType, error) {
	e := "nodebalancers/types"

	endpoint, err := generateListCacheURL(e, opts)
	if err != nil {
		return nil, err
	}

	if result := c.getCachedResponse(endpoint); result != nil {
		return result.([]NodeBalancerType), nil
	}

	response, err := getPaginatedResults[NodeBalancerType](ctx, c, e, opts)
	if err != nil {
		return nil, err
	}

	c.addCachedResponse(endpoint, response, &cacheExpiryTime)

	return response, nil
}
