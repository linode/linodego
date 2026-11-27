package linodego

import "context"

// ObjectStorageEndpointType constants start with Notification and include all known Linode API Notification Types.
type ObjectStorageEndpointType string

// NotificationType constants represent the actions that cause a Notification. New types may be added in the future.
const (
	ObjectStorageEndpointE0  ObjectStorageEndpointType = "E0"
	ObjectStorageEndpointE1  ObjectStorageEndpointType = "E1"
	ObjectStorageEndpointE2  ObjectStorageEndpointType = "E2"
	ObjectStorageEndpointE3  ObjectStorageEndpointType = "E3"
	ObjectStorageEndpointIA0 ObjectStorageEndpointType = "IA0"
)

// ObjectStorageEndpoint represents a linode object storage endpoint object
type ObjectStorageEndpoint struct {
	Region       string                    `json:"region"`
	S3Endpoint   *string                   `json:"s3_endpoint"`
	EndpointType ObjectStorageEndpointType `json:"endpoint_type"`

	// The API returns these fields only for accounts with Infrequent Access enabled.
	SupportedKeyClasses []ObjectStorageKeyClass `json:"supported_key_classes"`
	Country             string                  `json:"country"`
	RegionLabel         string                  `json:"region_label"`
}

// ListObjectStorageEndpoints lists all endpoints in all regions
func (c *Client) ListObjectStorageEndpoints(ctx context.Context, opts *ListOptions) ([]ObjectStorageEndpoint, error) {
	return getPaginatedResults[ObjectStorageEndpoint](ctx, c, "object-storage/endpoints", opts)
}
