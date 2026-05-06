package integration

import (
	"context"
	"testing"

	"github.com/linode/linodego"
	"github.com/stretchr/testify/assert"
)

func TestObjectStorageGlobalQuotas_List(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestObjectStorageGlobalQuotas_List")
	defer teardown()

	quotas, err := client.ListObjectStorageGlobalQuotas(context.Background(), &linodego.ListOptions{})
	assert.NoError(t, err, "Error listing ObjectStorageGlobalQuotas")

	// Just verify we got results and check structure of first quota if available
	if len(quotas) > 0 {
		firstQuota := quotas[0]
		assert.NotEmpty(t, firstQuota.QuotaID)
		assert.NotEmpty(t, firstQuota.QuotaName)
		assert.NotEmpty(t, firstQuota.QuotaType)
		assert.NotEmpty(t, firstQuota.Description)
		assert.NotEmpty(t, firstQuota.ResourceMetric)
		assert.Greater(t, firstQuota.QuotaLimit, 0)
		// HasUsage can be true or false, just check it's a valid field
	}
}

func TestObjectStorageGlobalQuotas_Get(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestObjectStorageGlobalQuotas_Get")
	defer teardown()

	// First list to get a valid quota ID
	quotas, err := client.ListObjectStorageGlobalQuotas(context.Background(), &linodego.ListOptions{})
	assert.NoError(t, err)

	if len(quotas) == 0 {
		t.Skip("No global quotas available to test Get")
	}

	// Use the first quota ID from the list
	quotaID := quotas[0].QuotaID

	quota, err := client.GetObjectStorageGlobalQuota(context.Background(), quotaID)
	assert.NoError(t, err)

	// Verify response structure
	assert.Equal(t, quotaID, quota.QuotaID)
	assert.NotEmpty(t, quota.QuotaName)
	assert.NotEmpty(t, quota.QuotaType)
	assert.NotEmpty(t, quota.Description)
	assert.NotEmpty(t, quota.ResourceMetric)
	assert.Greater(t, quota.QuotaLimit, 0)
}

func TestObjectStorageGlobalQuotaUsage_Get(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestObjectStorageGlobalQuotaUsage_Get")
	defer teardown()

	// First list to get a quota with usage support
	quotas, err := client.ListObjectStorageGlobalQuotas(context.Background(), &linodego.ListOptions{})
	assert.NoError(t, err)

	if len(quotas) == 0 {
		t.Skip("No global quotas available to test Usage")
	}

	// Find a quota that has usage support
	var quotaWithUsage *linodego.ObjectStorageGlobalQuota
	for _, q := range quotas {
		if q.HasUsage {
			quotaWithUsage = &q
			break
		}
	}

	if quotaWithUsage == nil {
		t.Skip("No global quotas with usage support available")
	}

	quotaUsage, err := client.GetObjectStorageGlobalQuotaUsage(context.Background(), quotaWithUsage.QuotaID)
	assert.NoError(t, err)

	// Verify response structure
	assert.Greater(t, quotaUsage.QuotaLimit, 0)
	if quotaUsage.Usage != nil {
		assert.GreaterOrEqual(t, *quotaUsage.Usage, 0)
	}
}
