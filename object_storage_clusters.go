package linodego

import (
	"context"
)

// ObjectStorageCluster represents a linode object storage cluster object
type ObjectStorageCluster struct {
	ID               string `json:"id"`
	Domain           string `json:"domain"`
	Status           string `json:"status"`
	Region           string `json:"region"`
	StaticSiteDomain string `json:"static_site_domain"`
}

// ListObjectStorageClusters lists ObjectStorageClusters
func (c *Client) ListObjectStorageClusters(ctx context.Context, opts *ListOptions) ([]ObjectStorageCluster, error) {
	return getPaginatedResults[ObjectStorageCluster](ctx, c, "object-storage/clusters", opts)
}

// Deprecated: GetObjectStorageCluster uses a deprecated API endpoint.
// GetObjectStorageCluster gets the template with the provided ID
func (c *Client) GetObjectStorageCluster(ctx context.Context, clusterID string) (*ObjectStorageCluster, error) {
	e := formatAPIPath("object-storage/clusters/%s", clusterID)
	return doGETRequest[ObjectStorageCluster](ctx, c, e)
}
