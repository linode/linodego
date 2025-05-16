package integration

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/linode/linodego"
)

func TestObjectStorageQuotas_Get(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestObjectStorageQuotas_Get")
	defer teardown()

	quota, err := client.GetObjectStorageQuota(context.Background(), "obj-objects-us-ord-1.linodeobjects.com")
	assert.NoError(t, err)

	expected := linodego.ObjectStorageQuota{
		QuotaID:        "obj-objects-us-ord-1.linodeobjects.com",
		QuotaName:      "max_objects",
		EndpointType:   "E1",
		S3Endpoint:     "us-ord-1.linodeobjects.com",
		Description:    "Maximum number of objects this customer is allowed to have on this endpoint",
		QuotaLimit:     100000000,
		ResourceMetric: "object",
	}

	assert.Equal(t, expected.QuotaID, quota.QuotaID)
	assert.Equal(t, expected.QuotaName, quota.QuotaName)
	assert.Equal(t, expected.EndpointType, quota.EndpointType)
	assert.Equal(t, expected.S3Endpoint, quota.S3Endpoint)
	assert.Equal(t, expected.Description, quota.Description)
	assert.Equal(t, expected.QuotaLimit, quota.QuotaLimit)
	assert.Equal(t, expected.ResourceMetric, quota.ResourceMetric)
}

func TestObjectStorageQuotas_List(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestObjectStorageQuotas_List")
	defer teardown()

	quotas, err := client.ListObjectStorageQuotas(context.Background(), &linodego.ListOptions{})
	assert.NoError(t, err, "Error listing ObjectStorageQuotas")

	targetQuotaID := "obj-buckets-us-mia-1.linodeobjects.com"
	var foundQuota *linodego.ObjectStorageQuota

	for _, quota := range quotas {
		if quota.QuotaID == targetQuotaID {
			foundQuota = &quota
			break
		}
	}

	if assert.NotNil(t, foundQuota, "Expected quota_id %q not found", targetQuotaID) {
		expected := linodego.ObjectStorageQuota{
			QuotaID:        "obj-buckets-us-mia-1.linodeobjects.com",
			QuotaName:      "max_buckets",
			EndpointType:   "E1",
			S3Endpoint:     "us-mia-1.linodeobjects.com",
			Description:    "Maximum number of buckets this customer is allowed to have on this endpoint",
			QuotaLimit:     1000,
			ResourceMetric: "bucket",
		}

		assert.Equal(t, expected.QuotaID, foundQuota.QuotaID)
		assert.Equal(t, expected.QuotaName, foundQuota.QuotaName)
		assert.Equal(t, expected.EndpointType, foundQuota.EndpointType)
		assert.Equal(t, expected.S3Endpoint, foundQuota.S3Endpoint)
		assert.Equal(t, expected.Description, foundQuota.Description)
		assert.Equal(t, expected.QuotaLimit, foundQuota.QuotaLimit)
		assert.Equal(t, expected.ResourceMetric, foundQuota.ResourceMetric)
	}
}

func TestObjectStorageQuotaUsage_Get(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestObjectStorageQuotaUsage_Get")
	defer teardown()

	quotaUsage, err := client.GetObjectStorageQuotaUsage(context.Background(), "obj-objects-us-ord-1.linodeobjects.com")
	assert.NoError(t, err)

	assert.Equal(t, 100000000, quotaUsage.QuotaLimit)
	assert.GreaterOrEqual(t, *quotaUsage.Usage, 0)
}
