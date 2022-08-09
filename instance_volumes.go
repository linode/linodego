package linodego

import (
	"context"
)

// InstanceVolumesPagedResponse represents a paginated InstanceVolume API response
type InstanceVolumesPagedResponse struct {
	*PageOptions
	Data []Volume `json:"data"`
}

// endpoint gets the endpoint URL for InstanceVolume
func (InstanceVolumesPagedResponse) endpoint(c *Client, ids ...any) string {
	id := ids[0].(int)
	endpoint, err := c.InstanceVolumes.endpointWithParams(id)
	if err != nil {
		panic(err)
	}
	return endpoint
}

// ListInstanceVolumes lists InstanceVolumes
func (c *Client) ListInstanceVolumes(ctx context.Context, linodeID int, opts *ListOptions) ([]Volume, error) {
	response := InstanceVolumesPagedResponse{}
	err := c.listHelper(ctx, &response, opts, linodeID)
	if err != nil {
		return nil, err
	}
	return response.Data, nil
}
