package unit

import (
	"context"
	"testing"

	"github.com/linode/linodego"
	"github.com/stretchr/testify/assert"
)

func TestLinodeQuotas_Get(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("linode_quotas_get")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("linode/quotas/123", fixtureData)

	quota, err := base.Client.GetLinodeQuota(context.Background(), 123)
	if err != nil {
		t.Fatalf("Error getting linode quota: %v", err)
	}

	assert.Equal(t, 123, quota.QuotaID)
	assert.Equal(t, "Total vCPU for Dedicated Plans", quota.QuotaName)
	assert.Equal(t, "Maximum number of vCPUs assigned to Linodes with Dedicated Plans in this Region", quota.Description)
	assert.Equal(t, 20, quota.QuotaLimit)
	assert.Equal(t, "cpu", quota.ResourceMetric)
	assert.Equal(t, "us-lax", quota.RegionApplied)
}

func TestLinodeQuotas_List(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("linode_quotas_list")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("linode/quotas", fixtureData)

	quotas, err := base.Client.ListLinodeQuotas(context.Background(), &linodego.ListOptions{})
	if err != nil {
		t.Fatalf("Error listing linode quotas: %v", err)
	}

	if len(quotas) < 1 {
		t.Fatalf("Expected to get a list of linode quotas but failed.")
	}

	assert.Equal(t, 123, quotas[0].QuotaID)
	assert.Equal(t, "Total vCPU for Dedicated Plans", quotas[0].QuotaName)
	assert.Equal(t, "Maximum number of vCPUs assigned to Linodes with Dedicated Plans in this Region", quotas[0].Description)
	assert.Equal(t, 20, quotas[0].QuotaLimit)
	assert.Equal(t, "cpu", quotas[0].ResourceMetric)
	assert.Equal(t, "us-lax", quotas[0].RegionApplied)
}

func TestLinodeQuotaUsage_Get(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("linode_quotas_usage_get")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("linode/quotas/123/usage", fixtureData)

	quotaUsage, err := base.Client.GetLinodeQuotaUsage(context.Background(), 123)
	if err != nil {
		t.Fatalf("Error getting linode quota usage: %v", err)
	}

	assert.Equal(t, 20, quotaUsage.QuotaLimit)
	assert.Equal(t, 5, *quotaUsage.Usage)
}
