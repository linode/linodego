package unit

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMaintenancePolicies_List(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("maintenance_policies_list")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("maintenance/policies", fixtureData)

	policies, err := base.Client.ListMaintenancePolicies(context.Background())
	assert.NoError(t, err)

	assert.Len(t, policies, 2, "Expected two maintenance policies to be returned")

	assert.Equal(t, "1", policies[0].ID)
	assert.Equal(t, "Default Migrate", policies[0].Name)
	assert.Equal(t, "predefined maintenance policy default for all linodes", policies[0].Description)
	assert.Equal(t, "migrate", policies[0].Type)
	assert.Equal(t, 3600, policies[0].NotificationPeriodSec)
	assert.Equal(t, true, policies[0].IsDefault)

	assert.Equal(t, "2", policies[1].ID)
	assert.Equal(t, "Default Power On/Off", policies[1].Name)
	assert.Equal(t, "predefined maintenance policy for general use cases", policies[1].Description)
	assert.Equal(t, "power on/off", policies[1].Type)
	assert.Equal(t, 1800, policies[1].NotificationPeriodSec)
	assert.Equal(t, false, policies[1].IsDefault)
}
