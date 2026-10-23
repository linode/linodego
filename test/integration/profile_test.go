package integration

import (
	"context"
	"testing"

	"github.com/linode/linodego/v2"
	"github.com/stretchr/testify/require"
)

func TestProfile_Get_smoke(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestProfile_Get")
	defer teardown()

	i, err := client.GetProfile(context.Background())
	if err != nil {
		t.Errorf("Error getting profile: %s", err)
	}
	if len(i.Email) < 1 {
		t.Errorf("Expected profile email")
	}
}

func TestProfile_Update(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestProfile_Update")
	defer teardown()

	profile, err := client.GetProfile(context.Background())
	if err != nil {
		t.Errorf("Error getting profile: %s", err)
	}

	updateOpts := profile.GetUpdateOptions()
	if updateOpts.Email != profile.Email {
		t.Errorf("Expected matching Username from GetUpdateOptions, got: %v", updateOpts)
	}

	i, err := client.UpdateProfile(context.Background(), updateOpts)
	if err != nil {
		t.Errorf("Error updating profile: %s", err)
	}
	if i.Email != updateOpts.Email {
		t.Errorf("Expected profile email to be changed, but found %v", i)
	}
}

func TestProfile_TFAEnforced(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestProfile_TFAEnforced")
	t.Cleanup(teardown)
	warnSensitiveTest(t)

	ctx := context.Background()
	settings, err := client.GetAccountSettings(ctx)
	require.NoError(t, err)
	profile, err := client.GetProfile(ctx)
	require.NoError(t, err)
	exemptions, err := client.ListTFAOptionalUsers(ctx, nil)
	require.NoError(t, err)

	expectedEnforcement := settings.TFAEnforced
	for _, user := range exemptions {
		if user.Username == profile.Username {
			expectedEnforcement = false
			break
		}
	}
	require.Equal(t, expectedEnforcement, profile.TFAEnforced)

	updated, err := client.UpdateProfile(ctx, linodego.ProfileUpdateOptions{
		EmailNotifications: new(profile.EmailNotifications),
	})
	require.NoError(t, err)
	require.Equal(t, profile.TFAEnforced, updated.TFAEnforced)
	require.Equal(t, profile.EmailNotifications, updated.EmailNotifications)
}
