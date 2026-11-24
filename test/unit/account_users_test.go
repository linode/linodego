package unit

import (
	"context"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/linode/linodego/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAccountUsers_List(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("account_users_list")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("account/users", fixtureData)

	users, err := base.Client.ListUsers(context.Background(), nil)
	if err != nil {
		t.Fatalf("Error listing users: %v", err)
	}

	assert.Equal(t, 1, len(users))
	user := users[0]
	assert.Equal(t, "jperez@linode.com", user.Email)
	assert.Equal(t, true, user.Restricted)
	assert.Equal(t, []string{"home-pc", "laptop"}, user.SSHKeys)
	assert.Equal(t, true, user.TFAEnabled)
	assert.True(t, user.TFAEnforced)
	assert.Equal(t, linodego.UserType("parent"), user.UserType)
	assert.Equal(t, "jsmith", user.Username)
	assert.Equal(t, "+5555555555", *user.VerifiedPhoneNumber)
}

func TestAccountUsers_Get(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("account_users_get")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("account/users/jsmith", fixtureData)

	user, err := base.Client.GetUser(context.Background(), "jsmith")
	if err != nil {
		t.Fatalf("Error getting user: %v", err)
	}

	assert.Equal(t, "jperez@linode.com", user.Email)
	assert.Equal(t, true, user.Restricted)
	assert.Equal(t, []string{"home-pc", "laptop"}, user.SSHKeys)
	assert.Equal(t, true, user.TFAEnabled)
	assert.False(t, user.TFAEnforced)
	assert.Equal(t, linodego.UserType("parent"), user.UserType)
	assert.Equal(t, "jsmith", user.Username)
	assert.Equal(t, "+5555555555", *user.VerifiedPhoneNumber)
}

func TestAccountUsers_Create(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("account_users_create")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	requestData := linodego.UserCreateOptions{
		Username:   "example_user",
		Email:      "example_user@linode.com",
		Restricted: true,
	}

	base.MockPost("account/users", fixtureData)

	user, err := base.Client.CreateUser(context.Background(), requestData)
	assert.NoError(t, err)

	assert.Equal(t, "example_user@linode.com", user.Email)
	assert.Equal(t, true, user.Restricted)
	assert.Equal(t, []string{"home-pc", "laptop"}, user.SSHKeys)
	assert.Equal(t, true, user.TFAEnabled)
	assert.True(t, user.TFAEnforced)
	assert.Equal(t, "example_user", user.Username)
	assert.Equal(t, "+5555555555", *user.VerifiedPhoneNumber)
}

func TestAccountUsers_Update(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("account_users_update")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	restricted := true

	requestData := linodego.UserUpdateOptions{
		Username:   "adevi",
		Email:      "jkowalski@linode.com",
		Restricted: &restricted,
	}

	base.MockPut("account/users/adevi", fixtureData)

	user, err := base.Client.UpdateUser(context.Background(), "adevi", requestData)
	assert.NoError(t, err)

	assert.Equal(t, "jkowalski@linode.com", user.Email)
	assert.Equal(t, true, user.Restricted)
	assert.Equal(t, []string{"home-pc", "laptop"}, user.SSHKeys)
	assert.Equal(t, true, user.TFAEnabled)
	assert.True(t, user.TFAEnforced)
	assert.Equal(t, linodego.UserType("parent"), user.UserType)
	assert.Equal(t, "adevi", user.Username)
	assert.Equal(t, "+5555555555", *user.VerifiedPhoneNumber)
}

func TestAccountUsers_Delete(t *testing.T) {
	client := createMockClient(t)

	httpmock.RegisterRegexpResponder("DELETE", mockRequestURL(t, "account/users/example-user"), httpmock.NewStringResponder(200, "{}"))

	if err := client.DeleteUser(context.Background(), "example-user"); err != nil {
		t.Fatal(err)
	}
}

func TestAccountUsers_ListTFAEnforced(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("account_users_list")
	require.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	t.Cleanup(func() {
		base.TearDown(t)
	})

	filter := `{"tfa_enforced":true}`
	httpmock.RegisterResponder("GET", base.BaseURL+"account/users", func(req *http.Request) (*http.Response, error) {
		assert.JSONEq(t, filter, req.Header.Get("X-Filter"))
		return httpmock.NewJsonResponse(http.StatusOK, fixtureData)
	})

	users, err := base.Client.ListUsers(context.Background(), linodego.NewListOptions(0, filter))
	require.NoError(t, err)
	require.Len(t, users, 1)
	assert.True(t, users[0].TFAEnforced)
}
