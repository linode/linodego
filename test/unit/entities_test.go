package unit

import (
	"context"
	"testing"

	"github.com/linode/linodego"
	"github.com/stretchr/testify/assert"
)

func TestEntities_List(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("entities_list")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet(formatMockAPIPath("entities"), fixtureData)

	entities, err := base.Client.ListEntities(context.Background(), &linodego.ListOptions{})
	assert.NoError(t, err)

	assert.Equal(t, 7, entities[0].ID)
}

func TestEntitiesPermissions_Get(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("entities_permissions")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet(formatMockAPIPath("iam/users/Linode/permissions/linode/123"), fixtureData)

	perms, err := base.Client.GetEntityRoles(context.Background(), "Linode", "linode", 123)
	assert.NoError(t, err)

	assert.Equal(t, "rebuild_linode", perms[1])
}
