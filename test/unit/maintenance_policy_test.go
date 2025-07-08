package unit

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMaintenancePolicies_List(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("maintenance_policies_list")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("maintenance/policies", fixtureData)

	policies, err := base.Client.ListMaintenancePolicies(context.Background(), nil)
	assert.NoError(t, err)

	assert.Len(t, policies, 2, "Expected two maintenance policies to be returned")

	assert.Equal(t, "linode/migrate", policies[0].Slug)
	assert.Equal(t, "Migrate", policies[0].Label)
	assert.Equal(t, "Migrates the Linode to a new host while it remains fully operational. Recommended for maximizing availability.", policies[0].Description)
	assert.Equal(t, "migrate", policies[0].Type)
	assert.Equal(t, 300, policies[0].NotificationPeriodSec)
	assert.Equal(t, true, policies[0].IsDefault)

	assert.Equal(t, "linode/power_off_on", policies[1].Slug)
	assert.Equal(t, "Power-off/on", policies[1].Label)
	assert.Equal(t, "Powers off the Linode at the start of the maintenance event and reboots it once the maintenance finishes. Recommended for maximizing performance.", policies[1].Description)
	assert.Equal(t, "power_off_on", policies[1].Type)
	assert.Equal(t, 1800, policies[1].NotificationPeriodSec)
	assert.Equal(t, false, policies[1].IsDefault)
}
