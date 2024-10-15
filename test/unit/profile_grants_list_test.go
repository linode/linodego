package unit

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGrantsList(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("profile_grants_list")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("profile/grants", fixtureData)

	grants, err := base.Client.GrantsList(context.Background())
	assert.NoError(t, err)

	assert.Equal(t, "read_only", string(*grants.Global.AccountAccess))
	assert.True(t, grants.Global.AddDatabases)
	assert.True(t, grants.Global.AddDomains)
	assert.True(t, grants.Global.AddFirewalls)
	assert.False(t, grants.Global.CancelAccount)

	assert.Len(t, grants.Database, 1)
	assert.Equal(t, 123, grants.Database[0].ID)
	assert.Equal(t, "example-entity", grants.Database[0].Label)
	assert.Equal(t, "read_only", string(grants.Database[0].Permissions))

	assert.Len(t, grants.Domain, 1)
	assert.Equal(t, 123, grants.Domain[0].ID)
	assert.Equal(t, "example-entity", grants.Domain[0].Label)
	assert.Equal(t, "read_only", string(grants.Domain[0].Permissions))

	assert.Len(t, grants.Firewall, 1)
	assert.Equal(t, 123, grants.Firewall[0].ID)
	assert.Equal(t, "example-entity", grants.Firewall[0].Label)
	assert.Equal(t, "read_only", string(grants.Firewall[0].Permissions))

	assert.Len(t, grants.Linode, 1)
	assert.Equal(t, 123, grants.Linode[0].ID)
	assert.Equal(t, "example-entity", grants.Linode[0].Label)
	assert.Equal(t, "read_only", string(grants.Linode[0].Permissions))
}
