package unit

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIAMAccountRolePermissions_Get(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("iam_account_get")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet(formatMockAPIPath("iam/role-permissions"), fixtureData)

	perms, err := base.Client.GetAccountRolePermissions(context.Background())
	assert.NoError(t, err)

	assert.Equal(t, "linode", perms.EntityAccess[0].Type)
}

func TestIAMUserRolePermissions_Get(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("iam_user_get")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet(formatMockAPIPath("iam/users/Linode/role-permissions"), fixtureData)

	perms, err := base.Client.GetUserRolePermissions(context.Background(), "Linode")
	assert.NoError(t, err)

	assert.Equal(t, 1, perms.EntityAccess[0].ID)
}

func TestIAMUserRolePermissions_Update(t *testing.T) {
	updateFixtureData, err := fixtures.GetFixture("iam_user_update")
	assert.NoError(t, err)
	getFixtureData, err := fixtures.GetFixture("iam_user_get")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet(formatMockAPIPath("iam/users/Linode/role-permissions"), getFixtureData)
	base.MockPut(formatMockAPIPath("iam/users/Linode/role-permissions"), updateFixtureData)

	before, err := base.Client.GetUserRolePermissions(context.Background(), "Linode")
	opts := before.GetUpdateOptions()
	opts.AccountAccess = append(opts.AccountAccess, "test_admin")
	after, err := base.Client.UpdateUserRolePermissions(context.Background(), "Linode", opts)
	assert.NoError(t, err)

	assert.Equal(t, 1, after.EntityAccess[0].ID)
	assert.NotEqual(t, before.AccountAccess, after.AccountAccess)
}

func TestIAMUserAccountPermissions_Get(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("iam_user_account_permissions_get")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet(formatMockAPIPath("iam/users/Linode/permissions/account"), fixtureData)

	perms, err := base.Client.GetUserAccountPermissions(context.Background(), "Linode")
	assert.NoError(t, err)

	assert.Equal(t, "list_events", perms[0])
}
