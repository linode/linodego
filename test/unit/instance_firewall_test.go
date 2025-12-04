package unit

import (
	"context"
	"testing"

	"github.com/linode/linodego"
	"github.com/stretchr/testify/assert"
)

func TestInstanceFirewalls_List(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("instance_firewall_list")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("linode/instances/123/firewalls", fixtureData)

	firewalls, err := base.Client.ListInstanceFirewalls(context.Background(), 123, nil)
	assert.NoError(t, err)
	assert.Len(t, firewalls, 2)

	assert.Equal(t, 456, firewalls[0].ID)
	assert.Equal(t, "firewall-1", firewalls[0].Label)

	assert.Equal(t, 789, firewalls[1].ID)
	assert.Equal(t, "firewall-2", firewalls[1].Label)
}

func TestInstanceFirewalls_Update(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("instance_firewall_update")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("linode/instances/123/firewalls", fixtureData)
	base.MockPut("linode/instances/123/firewalls", fixtureData)
	updateOpts := linodego.InstanceFirewallUpdateOptions{
		FirewallIDs: []int{789},
	}

	firewalls, err := base.Client.UpdateInstanceFirewalls(context.Background(), 123, updateOpts)
	assert.NoError(t, err)
	assert.NotNil(t, firewalls)
	assert.Len(t, firewalls, 1)
	assert.Equal(t, 789, firewalls[0].ID)
}
