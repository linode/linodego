package integration

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/linode/linodego/v2"
	"github.com/stretchr/testify/require"
)

func TestAccountTFAOptionalUsers_List(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestAccountTFAOptionalUsers_List")
	t.Cleanup(teardown)

	opts := linodego.NewListOptions(0, "")
	opts.PageSize = 25
	users, err := client.ListTFAOptionalUsers(context.Background(), opts)
	require.NoError(t, err)
	require.Len(t, users, opts.Results)

	for _, user := range users {
		require.NotEmpty(t, user.Username)
		require.Equal(t, "/v4/account/users/"+user.Username, user.URL)
	}
}

func TestAccountTFAOptionalUsers_Update(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestAccountTFAOptionalUsers_Update")
	t.Cleanup(teardown)

	ctx := context.Background()
	original, err := client.ListTFAOptionalUsers(ctx, nil)
	require.NoError(t, err)
	originalUsernames := make([]string, len(original))
	for index, user := range original {
		originalUsernames[index] = user.Username
	}

	user, deleteUser := createUser(t, client, func(opts *linodego.UserCreateOptions) {
		opts.Username = usernamePrefix + "tfaoptional"
		opts.Email = opts.Username + "@example.com"
		opts.Restricted = true
	})
	t.Cleanup(deleteUser)
	t.Cleanup(func() {
		_, err := client.UpdateTFAOptionalUsers(ctx, linodego.TFAOptionalUsersUpdateOptions{Usernames: originalUsernames})
		require.NoError(t, err, "failed to restore the original TFA exemption list")
	})

	usernames := append([]string{user.Username}, originalUsernames...)
	response, err := client.UpdateTFAOptionalUsers(ctx, linodego.TFAOptionalUsersUpdateOptions{Usernames: usernames})
	require.NoError(t, err)
	require.Equal(t, len(usernames), response.TFAOptionalUsersCount)
	require.Equal(t, "/v4/account/tfa-optional-users", response.URL)

	users, err := client.ListTFAOptionalUsers(ctx, nil)
	require.NoError(t, err)
	require.Len(t, users, len(usernames))
	require.Contains(t, users, linodego.TFAOptionalUser{Username: user.Username, URL: "/v4/account/users/" + user.Username})

	exempted, err := client.GetUser(ctx, user.Username)
	require.NoError(t, err)
	require.False(t, exempted.TFAEnforced)

	updated, err := client.UpdateUser(ctx, user.Username, linodego.UserUpdateOptions{Restricted: new(true)})
	require.NoError(t, err)
	require.False(t, updated.TFAEnforced)

	filter, err := json.Marshal(map[string]interface{}{"username": user.Username, "tfa_enforced": false})
	require.NoError(t, err)
	filtered, err := client.ListUsers(ctx, linodego.NewListOptions(0, string(filter)))
	require.NoError(t, err)
	require.Len(t, filtered, 1)
	require.Equal(t, user.Username, filtered[0].Username)
	require.False(t, filtered[0].TFAEnforced)

	response, err = client.UpdateTFAOptionalUsers(ctx, linodego.TFAOptionalUsersUpdateOptions{Usernames: originalUsernames})
	require.NoError(t, err)
	require.Equal(t, len(originalUsernames), response.TFAOptionalUsersCount)

	users, err = client.ListTFAOptionalUsers(ctx, nil)
	require.NoError(t, err)
	require.ElementsMatch(t, original, users)
}
