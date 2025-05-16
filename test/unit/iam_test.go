package unit

import (
	"context"
	"testing"

	"github.com/linode/linodego"
	"github.com/stretchr/testify/assert"
)

func TestUserRolePermissions_Get(t *testing.T) {
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

func TestUserRolePermissions_Update(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("iam_user_update")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	// TODO: get update role options
	opts := linodego.UserRolePermissionsUpdateOptions{}

	base.MockPut(formatMockAPIPath("iam/users/Linode/role-permissions"), fixtureData)

	perms, err := base.Client.UpdateUserRolePermissions(context.Background(), "Linode", opts)
	assert.NoError(t, err)

	assert.Equal(t, 1, perms.EntityAccess[0].ID)
}
