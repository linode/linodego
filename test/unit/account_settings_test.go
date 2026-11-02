package unit

import (
	"context"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/linode/linodego/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	assert.False(t, accountSettings.TFAEnforced)
	assert.Nil(t, accountSettings.LongviewSubscription, "Expected 'longview_subscription' to be nil")
	assert.True(t, accountSettings.BackupsEnabled, "Expected 'backups_enabled' to be true")
	assert.Equal(t, "active", *accountSettings.ObjectStorage, "Expected 'object_storage' to be 'active'")
	assert.Equal(
		t, linodego.LegacyConfigDefaultButLinodeAllowed, accountSettings.InterfacesForNewLinodes,
		"Expected 'interfaces_for_new_linodes' to be 'legacy_config_default_but_linode_allowed'",
	)
	assert.Equal(t, "linode/migrate", accountSettings.MaintenancePolicy, "Expected 'maintenance_policy' to be 'linode/migrate'")
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
		TFAEnforced:             Bool(true),
		InterfacesForNewLinodes: &i,
		MaintenancePolicy:       linodego.Pointer("linode/power_off_on"),
	}
	base.MockPut("account/settings", fixtureData)

	accountSettings, err := base.Client.UpdateAccountSettings(context.Background(), requestData)
	assert.NoError(t, err)
	assert.NotNil(t, accountSettings, "Account settings should not be nil")
	assert.False(t, accountSettings.Managed, "Expected 'managed' to be false")
	assert.True(t, accountSettings.NetworkHelper, "Expected 'network_helper' to be true")
	assert.True(t, accountSettings.TFAEnforced)
	assert.Nil(t, accountSettings.LongviewSubscription, "Expected 'longview_subscription' to be nil")
	assert.True(t, accountSettings.BackupsEnabled, "Expected 'backups_enabled' to be true")
	assert.Equal(t, "active", *accountSettings.ObjectStorage, "Expected 'object_storage' to be 'active'")
	assert.Equal(
		t, linodego.LegacyConfigDefaultButLinodeAllowed, accountSettings.InterfacesForNewLinodes,
		"Expected 'interfaces_for_new_linodes' to be 'legacy_config_default_but_linode_allowed'",
	)
	assert.Equal(t, "linode/power_off_on", accountSettings.MaintenancePolicy, "Expected 'maintenance_policy' to be 'linode/power_off_on'")
}

func TestAccountSettings_UpdateTFAEnforced(t *testing.T) {
	for _, testCase := range []struct {
		name     string
		enforced bool
		fixture  string
	}{
		{name: "enable", enforced: true, fixture: "account_settings_update"},
		{name: "disable", enforced: false, fixture: "account_settings_get"},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			fixtureData, err := fixtures.GetFixture(testCase.fixture)
			require.NoError(t, err)

			var base ClientBaseCase
			base.SetUp(t)
			t.Cleanup(func() {
				base.TearDown(t)
			})

			opts := linodego.AccountSettingsUpdateOptions{TFAEnforced: new(testCase.enforced)}
			httpmock.RegisterResponder("PUT", base.BaseURL+"account/settings", mockRequestBodyValidate(t, opts, fixtureData))

			settings, err := base.Client.UpdateAccountSettings(context.Background(), opts)
			require.NoError(t, err)
			assert.Equal(t, testCase.enforced, settings.TFAEnforced)
		})
	}
}
