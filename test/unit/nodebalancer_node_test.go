package unit

import (
	"context"
	"testing"

	"github.com/linode/linodego"
	"github.com/stretchr/testify/assert"
)

func TestNodeBalancerNode_Create(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("nodebalancer_node_create")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	// Mock the POST request with the fixture response
	base.MockPost("nodebalancers/123/configs/456/nodes", fixtureData)

	createOpts := linodego.NodeBalancerNodeCreateOptions{
		Address: "192.168.1.1",
		Label:   "Test Node",
		Weight:  50,
		Mode:    linodego.ModeAccept,
	}
	node, err := base.Client.CreateNodeBalancerNode(context.Background(), 123, 456, createOpts)
	assert.NoError(t, err)

	assert.Equal(t, 789, node.ID)
	assert.Equal(t, "192.168.1.1", node.Address)
	assert.Equal(t, "Test Node", node.Label)
	assert.Equal(t, 50, node.Weight)
	assert.Equal(t, linodego.ModeAccept, node.Mode)
}

func TestNodeBalancerNode_Update(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("nodebalancer_node_update")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	// Mock the PUT request with the fixture response
	base.MockPut("nodebalancers/123/configs/456/nodes/789", fixtureData)

	updateOpts := linodego.NodeBalancerNodeUpdateOptions{
		Address: "192.168.1.2",
		Label:   "Updated Node",
		Weight:  60,
		Mode:    linodego.ModeDrain,
	}
	node, err := base.Client.UpdateNodeBalancerNode(context.Background(), 123, 456, 789, updateOpts)
	assert.NoError(t, err)

	assert.Equal(t, 789, node.ID)
	assert.Equal(t, "192.168.1.2", node.Address)
	assert.Equal(t, "Updated Node", node.Label)
	assert.Equal(t, 60, node.Weight)
	assert.Equal(t, linodego.ModeDrain, node.Mode)
}

func TestNodeBalancerNode_List(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("nodebalancer_node_list")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	// Mock the GET request with the fixture response
	base.MockGet("nodebalancers/123/configs/456/nodes", fixtureData)

	nodes, err := base.Client.ListNodeBalancerNodes(context.Background(), 123, 456, nil)
	assert.NoError(t, err)
	assert.Len(t, nodes, 2)

	assert.Equal(t, 789, nodes[0].ID)
	assert.Equal(t, "192.168.1.1", nodes[0].Address)
	assert.Equal(t, "Test Node", nodes[0].Label)
	assert.Equal(t, 50, nodes[0].Weight)
	assert.Equal(t, linodego.ModeAccept, nodes[0].Mode)
}

func TestNodeBalancerNode_Delete(t *testing.T) {
	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	// Mock the DELETE request with the fixture response
	base.MockDelete("nodebalancers/123/configs/456/nodes/789", nil)

	err := base.Client.DeleteNodeBalancerNode(context.Background(), 123, 456, 789)
	assert.NoError(t, err)
}
