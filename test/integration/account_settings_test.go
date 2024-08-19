package integration

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/linode/linodego"
	"github.com/stretchr/testify/require"
)

func TestAccountSettings_Get(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestAccountSettings")
	defer teardown()

	// Mocking the API response
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mockSettings := linodego.AccountSettings{
		BackupsEnabled:       true,
		Managed:              true,
		NetworkHelper:        true,
		LongviewSubscription: String("longview-3"),
		ObjectStorage:        String("active"),
	}
	mockResponse, _ := json.Marshal(mockSettings)

	httpmock.RegisterResponder("GET", "https://api.linode.com/v4/account/settings",
		httpmock.NewStringResponder(200, string(mockResponse)))

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

	// Mocking the API response
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	opts := linodego.AccountSettingsUpdateOptions{
		BackupsEnabled:       Bool(false),
		LongviewSubscription: String("longview-10"),
		NetworkHelper:        Bool(false),
	}

	mockSettings := linodego.AccountSettings{
		BackupsEnabled:       false,
		NetworkHelper:        false,
		LongviewSubscription: String("longview-10"),
	}
	mockResponse, _ := json.Marshal(mockSettings)

	httpmock.RegisterResponder("PUT", "https://api.linode.com/v4/account/settings",
		httpmock.NewStringResponder(200, string(mockResponse)))

	settings, err := client.UpdateAccountSettings(context.Background(), opts)
	require.NoError(t, err, "Error updating Account Settings")

	require.False(t, settings.BackupsEnabled, "Expected BackupsEnabled to be false")
	require.False(t, settings.NetworkHelper, "Expected NetworkHelper to be false")
	require.NotNil(t, settings.LongviewSubscription, "Expected LongviewSubscription to be non-nil")
	require.Equal(t, "longview-10", *settings.LongviewSubscription, "Expected LongviewSubscription to be 'longview-10'")
}

func Bool(v bool) *bool       { return &v }
func String(v string) *string { return &v }
