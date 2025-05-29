package unit

import (
	"context"
	"testing"

	"github.com/linode/linodego"
	"github.com/stretchr/testify/assert"
)

// Helper function to create *bool
func Bool(value bool) *bool {
	return &value
}

func TestAccountSettings_Get(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("account_settings_get")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("account/settings", fixtureData)

	accountSettings, err := base.Client.GetAccountSettings(context.Background())
	assert.NoError(t, err)
	assert.NotNil(t, accountSettings, "Account settings should not be nil")
	assert.False(t, accountSettings.Managed, "Expected 'managed' to be false")
	assert.True(t, accountSettings.NetworkHelper, "Expected 'network_helper' to be true")
	assert.Nil(t, accountSettings.LongviewSubscription, "Expected 'longview_subscription' to be nil")
	assert.True(t, accountSettings.BackupsEnabled, "Expected 'backups_enabled' to be true")
	assert.Equal(t, "active", *accountSettings.ObjectStorage, "Expected 'object_storage' to be 'active'")
	assert.Equal(
		t, linodego.LegacyConfigDefaultButLinodeAllowed, accountSettings.InterfacesForNewLinodes,
		"Expected 'object_storage' to be 'active'",
	)
}

func TestAccountSettings_Update(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("account_settings_update")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	i := linodego.LegacyConfigDefaultButLinodeAllowed
	requestData := linodego.AccountSettingsUpdateOptions{
		BackupsEnabled:          Bool(true),
		NetworkHelper:           Bool(true),
		InterfacesForNewLinodes: &i,
	}
	base.MockPut("account/settings", fixtureData)

	accountSettings, err := base.Client.UpdateAccountSettings(context.Background(), requestData)
	assert.NoError(t, err)
	assert.NotNil(t, accountSettings, "Account settings should not be nil")
	assert.False(t, accountSettings.Managed, "Expected 'managed' to be false")
	assert.True(t, accountSettings.NetworkHelper, "Expected 'network_helper' to be true")
	assert.Nil(t, accountSettings.LongviewSubscription, "Expected 'longview_subscription' to be nil")
	assert.True(t, accountSettings.BackupsEnabled, "Expected 'backups_enabled' to be true")
	assert.Equal(t, "active", *accountSettings.ObjectStorage, "Expected 'object_storage' to be 'active'")
	assert.Equal(
		t, linodego.LegacyConfigDefaultButLinodeAllowed, accountSettings.InterfacesForNewLinodes,
		"Expected 'object_storage' to be 'active'",
	)
}
