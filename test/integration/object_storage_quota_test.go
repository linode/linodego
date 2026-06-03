package integration

import (
	"context"
	"strings"
	"testing"

	"github.com/linode/linodego/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// getExistingObjectStorageQuotaID is a helper function to retrieve an existing Object Storage quota ID
// It is used because accounts may differ in their QuotaID names
func getExistingObjectStorageQuotaID(t *testing.T, prefix string) string {
	targetQuotaID := "obj-objects-us-ord-10.linodeobjects.com"

	quotas, err := client.ListObjectStorageQuotas(context.Background(), &linodego.ListOptions{})
	assert.NoError(t, err, "Error listing ObjectStorageQuotas")

	for _, quota := range quotas {
		if strings.Contains(quota.QuotaID, prefix) {
			targetQuotaID = quota.QuotaID
			return targetQuotaID
		}
	}
	return targetQuotaID
}

func TestObjectStorageQuotas_Get(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestObjectStorageQuotas_Get")
	defer teardown()

	targetQuotaID := getExistingObjectStorageQuotaID(t, "obj-objects")
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
	require.NoError(t, err, "Error listing ObjectStorageQuotas")

	targetQuota := quotas[0]
	assert.NotEmpty(t, targetQuota.QuotaID)
	assert.NotEmpty(t, targetQuota.QuotaName)
	assert.NotEmpty(t, targetQuota.EndpointType)
	assert.NotEmpty(t, targetQuota.S3Endpoint)
	assert.NotEmpty(t, targetQuota.Description)
	assert.Greater(t, targetQuota.QuotaLimit, 0)
	assert.NotEmpty(t, targetQuota.ResourceMetric)
	assert.NotEmpty(t, targetQuota.QuotaType)
	assert.NotNil(t, targetQuota.HasUsage)
}

func TestObjectStorageQuotaUsage_Get(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestObjectStorageQuotaUsage_Get")
	defer teardown()

	objectQuotaID := getExistingObjectStorageQuotaID(t, "obj-objects")
	quotaUsage, err := client.GetObjectStorageQuotaUsage(context.Background(), objectQuotaID)
	assert.NoError(t, err)

	assert.GreaterOrEqual(t, quotaUsage.QuotaLimit, 100000000)
	assert.GreaterOrEqual(t, *quotaUsage.Usage, 0)
}
