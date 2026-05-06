package integration

import (
	"context"
	"testing"

	"github.com/linode/linodego"
	"github.com/stretchr/testify/assert"
)

func TestObjectStorageQuotas_Get(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestObjectStorageQuotas_Get")
	defer teardown()

	targetQuotaID := "obj-objects-us-ord-1.linodeobjects.com"
	quota, err := client.GetObjectStorageQuota(context.Background(), targetQuotaID)
	assert.NoError(t, err)

	assert.Equal(t, targetQuotaID, quota.QuotaID)
	assert.NotEmpty(t, quota.QuotaName)
	assert.NotEmpty(t, quota.EndpointType)
	assert.NotEmpty(t, quota.S3Endpoint)
	assert.NotEmpty(t, quota.Description)
	assert.Greater(t, quota.QuotaLimit, 0)
	assert.NotEmpty(t, quota.ResourceMetric)
	assert.NotEmpty(t, quota.QuotaType)
	assert.True(t, quota.HasUsage)
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
		assert.Equal(t, targetQuotaID, foundQuota.QuotaID)
		assert.NotEmpty(t, foundQuota.QuotaName)
		assert.NotEmpty(t, foundQuota.EndpointType)
		assert.NotEmpty(t, foundQuota.S3Endpoint)
		assert.NotEmpty(t, foundQuota.Description)
		assert.Greater(t, foundQuota.QuotaLimit, 0)
		assert.NotEmpty(t, foundQuota.ResourceMetric)
		assert.NotEmpty(t, foundQuota.QuotaType)
		assert.True(t, foundQuota.HasUsage)
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
