package unit

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/linode/linodego"
)

func TestNodeBalancerFirewalls_List(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("nodebalancer_firewall_list")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	// Mock the API response for nodebalancer firewalls
	base.MockGet("nodebalancers/123/firewalls", fixtureData)

	firewalls, err := base.Client.ListNodeBalancerFirewalls(context.Background(), 123, nil)

	assert.NoError(t, err)
	assert.Len(t, firewalls, 2)

	// Check the details of the first firewall
	assert.Equal(t, 789, firewalls[0].ID)
	assert.Equal(t, "firewall-1", firewalls[0].Label)
	assert.Equal(t, linodego.FirewallStatus("enabled"), firewalls[0].Status)

	// Check the details of the second firewall
	assert.Equal(t, 790, firewalls[1].ID)
	assert.Equal(t, "firewall-2", firewalls[1].Label)
	assert.Equal(t, linodego.FirewallStatus("disabled"), firewalls[1].Status)
}
