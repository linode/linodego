package unit

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/linode/linodego"
)

func TestObjectStorageQuotas_Get(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("object_storage_quotas_get")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("object-storage/quotas/obj-objects-us-ord-1", fixtureData)

	quota, err := base.Client.GetObjectStorageQuota(context.Background(), "obj-objects-us-ord-1")
	if err != nil {
		t.Fatalf("Error getting object-storage quota: %v", err)
	}

	assert.Equal(t, "obj-objects-us-ord-1", quota.QuotaID)
	assert.Equal(t, "Object Storage Maximum Objects", quota.QuotaName)
	assert.Equal(t, "Maximum number of Objects this customer is allowed to have on this endpoint.", quota.Description)
	assert.Equal(t, "E1", quota.EndpointType)
	assert.Equal(t, "us-iad-1.linodeobjects.com", quota.S3Endpoint)
	assert.Equal(t, 50, quota.QuotaLimit)
	assert.Equal(t, "object", quota.ResourceMetric)
}

func TestObjectStorageQuotas_List(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("object_storage_quotas_list")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("object-storage/quotas", fixtureData)

	quotas, err := base.Client.ListObjectStorageQuotas(context.Background(), &linodego.ListOptions{})
	if err != nil {
		t.Fatalf("Error listing object-storage quotas: %v", err)
	}

	if len(quotas) < 1 {
		t.Fatalf("Expected to get a list of object-storage quotas but failed.")
	}

	assert.Equal(t, "obj-objects-us-ord-1", quotas[0].QuotaID)
	assert.Equal(t, "Object Storage Maximum Objects", quotas[0].QuotaName)
	assert.Equal(t, "Maximum number of Objects this customer is allowed to have on this endpoint.", quotas[0].Description)
	assert.Equal(t, "E1", quotas[0].EndpointType)
	assert.Equal(t, "us-iad-1.linodeobjects.com", quotas[0].S3Endpoint)
	assert.Equal(t, 50, quotas[0].QuotaLimit)
	assert.Equal(t, "object", quotas[0].ResourceMetric)
}

func TestObjectStorageQuotaUsage_Get(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("object_storage_quotas_usage_get")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("object-storage/quotas/obj-bucket-us-ord-1/usage", fixtureData)

	quotaUsage, err := base.Client.GetObjectStorageQuotaUsage(context.Background(), "obj-bucket-us-ord-1")
	if err != nil {
		t.Fatalf("Error getting object storage quota usage: %v", err)
	}

	assert.Equal(t, 100, quotaUsage.QuotaLimit)
	assert.Equal(t, 10, *quotaUsage.Usage)
}
