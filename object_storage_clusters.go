package linodego

import (
	"context"
	"fmt"

	"github.com/go-resty/resty/v2"
)

// ObjectStorageCluster represents a linode object storage cluster object
type ObjectStorageCluster struct {
	ID               string `json:"id"`
	Domain           string `json:"domain"`
	Status           string `json:"status"`
	Region           string `json:"region"`
	StaticSiteDomain string `json:"static_site_domain"`
}

// ObjectStorageClustersPagedResponse represents a linode API response for listing
type ObjectStorageClustersPagedResponse struct {
	*PageOptions
	Data []ObjectStorageCluster `json:"data"`
}

// endpoint gets the endpoint URL for ObjectStorageCluster
func (ObjectStorageClustersPagedResponse) endpoint(c *Client, _ ...any) string {
	endpoint, err := c.ObjectStorageClusters.Endpoint()
	if err != nil {
		panic(err)
	}
	return endpoint
}

func (resp *ObjectStorageClustersPagedResponse) castResult(r *resty.Request, e string) (int, int, error) {
	res, err := coupleAPIErrors(r.SetResult(ObjectStorageClustersPagedResponse{}).Get(e))
	if err != nil {
		return 0, 0, err
	}
	castedRes := res.Result().(*ObjectStorageClustersPagedResponse)
	resp.Data = append(resp.Data, castedRes.Data...)
	return castedRes.Pages, castedRes.Results, nil
}

// ListObjectStorageClusters lists ObjectStorageClusters
func (c *Client) ListObjectStorageClusters(ctx context.Context, opts *ListOptions) ([]ObjectStorageCluster, error) {
	response := ObjectStorageClustersPagedResponse{}
	err := c.listHelper(ctx, &response, opts)
	if err != nil {
		return nil, err
	}
	return response.Data, nil
}

// GetObjectStorageCluster gets the template with the provided ID
func (c *Client) GetObjectStorageCluster(ctx context.Context, id string) (*ObjectStorageCluster, error) {
	e, err := c.ObjectStorageClusters.Endpoint()
	if err != nil {
		return nil, err
	}
	e = fmt.Sprintf("%s/%s", e, id)
	r, err := coupleAPIErrors(c.R(ctx).SetResult(&ObjectStorageCluster{}).Get(e))
	if err != nil {
		return nil, err
	}
	return r.Result().(*ObjectStorageCluster), nil
}
