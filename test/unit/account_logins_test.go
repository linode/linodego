package unit

import (
	"context"
	"fmt"
	"testing"

	"github.com/linode/linodego"
	"github.com/stretchr/testify/assert"
)

func TestAccountLogins_List(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("account_logins_list")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("account/logins", fixtureData)

	logins, err := base.Client.ListLogins(context.Background(), &linodego.ListOptions{})
	assert.NoError(t, err)

	assert.Len(t, logins, 1, "Expected one login record to be returned")

	login := logins[0]
	assert.Equal(t, 1234, login.ID, "Expected login ID to be 1234")
	assert.Equal(t, "2018-01-01 00:01:01 +0000 UTC", login.Datetime.String(), "Expected login datetime to be '2018-01-01T00:01:01'")
	assert.Equal(t, "192.0.2.0", login.IP, "Expected login IP to be '192.0.2.0'")
	assert.True(t, login.Restricted, "Expected login restricted to be true.")
	assert.Equal(t, "successful", login.Status, "Expected login status to be 'successful'")
	assert.Equal(t, "example_user", login.Username, "Expected login username to be 'example_user'")
}

func TestAccountLogin_Get(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("account_logins_get")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	loginID := 1234
	base.MockGet(fmt.Sprintf("account/logins/%d", loginID), fixtureData)

	login, err := base.Client.GetLogin(context.Background(), loginID)
	assert.NoError(t, err)

	assert.Equal(t, 1234, login.ID, "Expected login ID to be 1234")
	assert.Equal(t, "2018-01-01 00:01:01 +0000 UTC", login.Datetime.String(), "Expected login datetime to be '2018-01-01T00:01:01'")
	assert.Equal(t, "192.0.2.0", login.IP, "Expected login IP to be '192.0.2.0'")
	assert.True(t, login.Restricted, "Expected login restricted to be true")
	assert.Equal(t, "successful", login.Status, "Expected login status to be 'successful'")
	assert.Equal(t, "example_user", login.Username, "Expected login username to be 'example_user'")
}
