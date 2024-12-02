package unit

import (
	"context"
	"github.com/linode/linodego"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAccountUserGrants_Get(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("account_user_grants_get")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("account/users/example-user/grants", fixtureData)

	grants, err := base.Client.GetUserGrants(context.Background(), "example-user")
	if err != nil {
		t.Fatalf("Error getting grants: %v", err)
	}

	assert.Equal(t, 123, grants.Database[0].ID)
	assert.Equal(t, "example-entity", grants.Database[0].Label)
	assert.Equal(t, linodego.GrantPermissionLevel("read_only"), grants.Database[0].Permissions)
	assert.Equal(t, 123, grants.Domain[0].ID)
	assert.Equal(t, "example-entity", grants.Domain[0].Label)
	assert.Equal(t, linodego.GrantPermissionLevel("read_only"), grants.Domain[0].Permissions)
	assert.Equal(t, 123, grants.Firewall[0].ID)
	assert.Equal(t, "example-entity", grants.Firewall[0].Label)
	assert.Equal(t, linodego.GrantPermissionLevel("read_only"), grants.Firewall[0].Permissions)
	assert.Equal(t, 123, grants.Image[0].ID)
	assert.Equal(t, "example-entity", grants.Image[0].Label)
	assert.Equal(t, linodego.GrantPermissionLevel("read_only"), grants.Image[0].Permissions)
	assert.Equal(t, 123, grants.Linode[0].ID)
	assert.Equal(t, "example-entity", grants.Linode[0].Label)
	assert.Equal(t, linodego.GrantPermissionLevel("read_only"), grants.Linode[0].Permissions)
	assert.Equal(t, 123, grants.Longview[0].ID)
	assert.Equal(t, "example-entity", grants.Longview[0].Label)
	assert.Equal(t, linodego.GrantPermissionLevel("read_only"), grants.Longview[0].Permissions)
	assert.Equal(t, 123, grants.NodeBalancer[0].ID)
	assert.Equal(t, "example-entity", grants.NodeBalancer[0].Label)
	assert.Equal(t, linodego.GrantPermissionLevel("read_only"), grants.NodeBalancer[0].Permissions)
	assert.Equal(t, 123, grants.PlacementGroup[0].ID)
	assert.Equal(t, "example-entity", grants.PlacementGroup[0].Label)
	assert.Equal(t, linodego.GrantPermissionLevel("read_only"), grants.PlacementGroup[0].Permissions)
	assert.Equal(t, 123, grants.StackScript[0].ID)
	assert.Equal(t, "example-entity", grants.StackScript[0].Label)
	assert.Equal(t, linodego.GrantPermissionLevel("read_only"), grants.StackScript[0].Permissions)
	assert.Equal(t, 123, grants.Volume[0].ID)
	assert.Equal(t, "example-entity", grants.Volume[0].Label)
	assert.Equal(t, linodego.GrantPermissionLevel("read_only"), grants.Volume[0].Permissions)
	assert.Equal(t, 123, grants.VPC[0].ID)
	assert.Equal(t, "example-entity", grants.VPC[0].Label)
	assert.Equal(t, linodego.GrantPermissionLevel("read_only"), grants.VPC[0].Permissions)
	assert.Equal(t, linodego.GrantPermissionLevel("read_only"), *grants.Global.AccountAccess)
	assert.Equal(t, true, grants.Global.AddDatabases)
	assert.Equal(t, true, grants.Global.AddDomains)
	assert.Equal(t, true, grants.Global.AddFirewalls)
	assert.Equal(t, true, grants.Global.AddImages)
	assert.Equal(t, true, grants.Global.AddLinodes)
	assert.Equal(t, true, grants.Global.AddLongview)
	assert.Equal(t, true, grants.Global.AddNodeBalancers)
	assert.Equal(t, true, grants.Global.AddPlacementGroups)
	assert.Equal(t, true, grants.Global.AddStackScripts)
	assert.Equal(t, true, grants.Global.AddVolumes)
	assert.Equal(t, true, grants.Global.AddVPCs)
	assert.Equal(t, false, grants.Global.CancelAccount)
	assert.Equal(t, true, grants.Global.ChildAccountAccess)
	assert.Equal(t, true, grants.Global.LongviewSubscription)
}

func TestAccountGrants_Update(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("account_user_grants_update")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	globalGrants := linodego.GlobalUserGrants{
		AccountAccess: nil,
	}

	requestData := linodego.UserGrantsUpdateOptions{
		Global: globalGrants,
	}

	base.MockPut("account/users/example-user/grants", fixtureData)

	grants, err := base.Client.UpdateUserGrants(context.Background(), "example-user", requestData)
	assert.NoError(t, err)

	assert.Equal(t, 123, grants.Database[0].ID)
	assert.Equal(t, "example-entity", grants.Database[0].Label)
	assert.Equal(t, linodego.GrantPermissionLevel("read_only"), grants.Database[0].Permissions)
	assert.Equal(t, 123, grants.Domain[0].ID)
	assert.Equal(t, "example-entity", grants.Domain[0].Label)
	assert.Equal(t, linodego.GrantPermissionLevel("read_only"), grants.Domain[0].Permissions)
	assert.Equal(t, 123, grants.Firewall[0].ID)
	assert.Equal(t, "example-entity", grants.Firewall[0].Label)
	assert.Equal(t, linodego.GrantPermissionLevel("read_only"), grants.Firewall[0].Permissions)
	assert.Equal(t, 123, grants.Image[0].ID)
	assert.Equal(t, "example-entity", grants.Image[0].Label)
	assert.Equal(t, linodego.GrantPermissionLevel("read_only"), grants.Image[0].Permissions)
	assert.Equal(t, 123, grants.Linode[0].ID)
	assert.Equal(t, "example-entity", grants.Linode[0].Label)
	assert.Equal(t, linodego.GrantPermissionLevel("read_only"), grants.Linode[0].Permissions)
	assert.Equal(t, 123, grants.Longview[0].ID)
	assert.Equal(t, "example-entity", grants.Longview[0].Label)
	assert.Equal(t, linodego.GrantPermissionLevel("read_only"), grants.Longview[0].Permissions)
	assert.Equal(t, 123, grants.NodeBalancer[0].ID)
	assert.Equal(t, "example-entity", grants.NodeBalancer[0].Label)
	assert.Equal(t, linodego.GrantPermissionLevel("read_only"), grants.NodeBalancer[0].Permissions)
	assert.Equal(t, 123, grants.PlacementGroup[0].ID)
	assert.Equal(t, "example-entity", grants.PlacementGroup[0].Label)
	assert.Equal(t, linodego.GrantPermissionLevel("read_only"), grants.PlacementGroup[0].Permissions)
	assert.Equal(t, 123, grants.StackScript[0].ID)
	assert.Equal(t, "example-entity", grants.StackScript[0].Label)
	assert.Equal(t, linodego.GrantPermissionLevel("read_only"), grants.StackScript[0].Permissions)
	assert.Equal(t, 123, grants.Volume[0].ID)
	assert.Equal(t, "example-entity", grants.Volume[0].Label)
	assert.Equal(t, linodego.GrantPermissionLevel("read_only"), grants.Volume[0].Permissions)
	assert.Equal(t, 123, grants.VPC[0].ID)
	assert.Equal(t, "example-entity", grants.VPC[0].Label)
	assert.Equal(t, linodego.GrantPermissionLevel("read_only"), grants.VPC[0].Permissions)
	assert.Equal(t, linodego.GrantPermissionLevel("read_write"), *grants.Global.AccountAccess)
	assert.Equal(t, true, grants.Global.AddDatabases)
	assert.Equal(t, true, grants.Global.AddDomains)
	assert.Equal(t, true, grants.Global.AddFirewalls)
	assert.Equal(t, true, grants.Global.AddImages)
	assert.Equal(t, true, grants.Global.AddLinodes)
	assert.Equal(t, true, grants.Global.AddLongview)
	assert.Equal(t, true, grants.Global.AddNodeBalancers)
	assert.Equal(t, true, grants.Global.AddPlacementGroups)
	assert.Equal(t, true, grants.Global.AddStackScripts)
	assert.Equal(t, true, grants.Global.AddVolumes)
	assert.Equal(t, true, grants.Global.AddVPCs)
	assert.Equal(t, false, grants.Global.CancelAccount)
	assert.Equal(t, true, grants.Global.ChildAccountAccess)
	assert.Equal(t, true, grants.Global.LongviewSubscription)
}
