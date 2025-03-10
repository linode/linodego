package linodego

import (
	"context"
)

type ObjectStorageBucketCertV2 struct {
	SSL *bool `json:"ssl"`
}

type ObjectStorageBucketCertUploadOptions struct {
	Certificate string `json:"certificate"`
	PrivateKey  string `json:"private_key"`
}

// UploadObjectStorageBucketCert uploads a TLS/SSL Cert to be used with an Object Storage Bucket.
func (c *Client) UploadObjectStorageBucketCertV2(ctx context.Context, clusterOrRegionID, bucket string, opts ObjectStorageBucketCertUploadOptions) (*ObjectStorageBucketCertV2, error) {
	e := formatAPIPath("object-storage/buckets/%s/%s/ssl", clusterOrRegionID, bucket)
	return doPOSTRequest[ObjectStorageBucketCertV2](ctx, c, e, opts)
}

// GetObjectStorageBucketCertV2 gets an ObjectStorageBucketCert
func (c *Client) GetObjectStorageBucketCertV2(ctx context.Context, clusterOrRegionID, bucket string) (*ObjectStorageBucketCertV2, error) {
	e := formatAPIPath("object-storage/buckets/%s/%s/ssl", clusterOrRegionID, bucket)
	return doGETRequest[ObjectStorageBucketCertV2](ctx, c, e)
}

// DeleteObjectStorageBucketCert deletes an ObjectStorageBucketCert
func (c *Client) DeleteObjectStorageBucketCert(ctx context.Context, clusterOrRegionID, bucket string) error {
	e := formatAPIPath("object-storage/buckets/%s/%s/ssl", clusterOrRegionID, bucket)
	return doDELETERequest(ctx, c, e)
}
