package unit

import (
	"context"
	"testing"

	"github.com/linode/linodego"
	"github.com/stretchr/testify/assert"
)

func TestNodeBalancer_Create(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("nodebalancer_create")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	// Mock the POST request with the fixture response
	base.MockPost("nodebalancers", fixtureData)

	label := "Test NodeBalancer"
	createOpts := linodego.NodeBalancerCreateOptions{
		Label:  &label,
		Region: linodego.Pointer("us-east"),
		Tags:   []string{"test", "example"},
	}
	nodebalancer, err := base.Client.CreateNodeBalancer(context.Background(), createOpts)
	assert.NoError(t, err)

	assert.Equal(t, 123, nodebalancer.ID)
	assert.Equal(t, "Test NodeBalancer", *nodebalancer.Label)
	assert.Equal(t, "us-east", nodebalancer.Region)
	assert.Equal(t, []string{"test", "example"}, nodebalancer.Tags)
}

func TestNodeBalancer_Get(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("nodebalancer_get")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	// Mock the GET request with the fixture response
	base.MockGet("nodebalancers/123", fixtureData)

	nodebalancer, err := base.Client.GetNodeBalancer(context.Background(), 123)
	assert.NoError(t, err)

	assert.Equal(t, 123, nodebalancer.ID, "Expected NodeBalancer ID to match")
	assert.Equal(t, "Existing NodeBalancer", *nodebalancer.Label, "Expected NodeBalancer label to match")
	assert.Equal(t, "us-west", nodebalancer.Region, "Expected NodeBalancer region to match")
}

func TestNodeBalancer_List(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("nodebalancers_list")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	// Mock the GET request with the fixture response
	base.MockGet("nodebalancers", fixtureData)

	nodebalancers, err := base.Client.ListNodeBalancers(context.Background(), &linodego.ListOptions{})
	assert.NoError(t, err)

	assert.Len(t, nodebalancers, 2, "Expected two NodeBalancers in the list")

	// Verify details of the first NodeBalancer
	assert.Equal(t, 123, nodebalancers[0].ID, "Expected first NodeBalancer ID to match")
	assert.Equal(t, "NodeBalancer A", *nodebalancers[0].Label, "Expected first NodeBalancer label to match")
	assert.Equal(t, "us-east", nodebalancers[0].Region, "Expected first NodeBalancer region to match")
	assert.Equal(t, []string{"tag1", "tag2"}, nodebalancers[0].Tags, "Expected first NodeBalancer tags to match")

	// Verify details of the second NodeBalancer
	assert.Equal(t, 456, nodebalancers[1].ID, "Expected second NodeBalancer ID to match")
	assert.Equal(t, "NodeBalancer B", *nodebalancers[1].Label, "Expected second NodeBalancer label to match")
	assert.Equal(t, "us-west", nodebalancers[1].Region, "Expected second NodeBalancer region to match")
	assert.Equal(t, []string{"tag3"}, nodebalancers[1].Tags, "Expected second NodeBalancer tags to match")
}

func TestNodeBalancer_Update(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("nodebalancer_update")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	// Mock the PUT request with the fixture response
	base.MockPut("nodebalancers/123", fixtureData)

	label := "Updated NodeBalancer"
	updateOpts := linodego.NodeBalancerUpdateOptions{
		Label: &label,
		Tags:  []string{"updated", "production"},
	}
	nodebalancer, err := base.Client.UpdateNodeBalancer(context.Background(), 123, updateOpts)
	assert.NoError(t, err)

	assert.Equal(t, 456, nodebalancer.ID)
	assert.Equal(t, "Updated NodeBalancer", *nodebalancer.Label)
	assert.Equal(t, []string{"updated", "production"}, nodebalancer.Tags)
}

func TestNodeBalancer_Delete(t *testing.T) {
	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	// Mock the DELETE request
	base.MockDelete("nodebalancers/123", nil)

	err := base.Client.DeleteNodeBalancer(context.Background(), 123)
	assert.NoError(t, err, "Expected no error when deleting NodeBalancer")
}
