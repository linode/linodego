package unit

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInstance_NodeBalancers_List(t *testing.T) {
	fixtures := NewTestFixtures()

	fixtureData, err := fixtures.GetFixture("instance_nodebalancers_list")
	if err != nil {
		t.Fatalf("Failed to load fixture: %v", err)
	}

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("linode/instances/12345/nodebalancers", fixtureData)

	nodebalancers, err := base.Client.ListInstanceNodeBalancers(context.Background(), 12345, nil)
	if err != nil {
		t.Fatalf("Error listing instance nodebalancers: %v", err)
	}

	assert.Equal(t, 1, len(nodebalancers))
	nb := nodebalancers[0]

	assert.Equal(t, 0, nb.ClientConnThrottle)
	assert.Equal(t, "192.0.2.1.ip.linodeusercontent.com", *nb.Hostname)
	assert.Equal(t, 12345, nb.ID)
	assert.Equal(t, "203.0.113.1", *nb.IPv4)
	assert.Nil(t, nb.IPv6)
	assert.Equal(t, "balancer12345", *nb.Label)
	assert.Equal(t, "us-east", nb.Region)
	assert.Equal(t, []string{"example tag", "another example"}, nb.Tags)
	assert.Equal(t, 28.91200828552246, *nb.Transfer.In)
	assert.Equal(t, 3.5487728118896484, *nb.Transfer.Out)
	assert.Equal(t, 32.46078109741211, *nb.Transfer.Total)
}
