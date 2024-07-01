package integration

import (
    "context"
    "testing"

    "github.com/stretchr/testify/require"
    "github.com/linode/linodego"
)

func TestAccountSettings_Get(t *testing.T) {
    client, teardown := createTestClient(t, "fixtures/TestAccountSettings")
    defer teardown()

    settings, err := client.GetAccountSettings(context.Background())
    require.NoError(t, err, "Error getting Account Settings")

    require.True(t, settings.BackupsEnabled, "Expected BackupsEnabled to be true")
    require.True(t, settings.Managed, "Expected Managed to be true")
    require.True(t, settings.NetworkHelper, "Expected NetworkHelper to be true")
    require.NotNil(t, settings.LongviewSubscription, "Expected LongviewSubscription to be non-nil")
    require.Equal(t, "longview-3", *settings.LongviewSubscription, "Expected LongviewSubscription to be 'longview-3'")
    require.NotNil(t, settings.ObjectStorage, "Expected ObjectStorage to be non-nil")
    require.Equal(t, "active", *settings.ObjectStorage, "Expected ObjectStorage to be 'active'")
}

func TestAccountSettings_Update(t *testing.T) {
    client, teardown := createTestClient(t, "fixtures/TestAccountSettings")
    defer teardown()

    opts := linodego.AccountSettingsUpdateOptions{
        BackupsEnabled:      Bool(false),
        LongviewSubscription: String("longview-10"),
        NetworkHelper:       Bool(false),
    }

    settings, err := client.UpdateAccountSettings(context.Background(), opts)
    require.NoError(t, err, "Error updating Account Settings")

    require.False(t, settings.BackupsEnabled, "Expected BackupsEnabled to be false")
    require.False(t, settings.NetworkHelper, "Expected NetworkHelper to be false")
    require.NotNil(t, settings.LongviewSubscription, "Expected LongviewSubscription to be non-nil")
    require.Equal(t, "longview-10", *settings.LongviewSubscription, "Expected LongviewSubscription to be 'longview-10'")
}

func Bool(v bool) *bool { return &v }
func String(v string) *string { return &v }
