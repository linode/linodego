package unit

import (
	"context"
	"testing"

	"github.com/linode/linodego"
	"github.com/stretchr/testify/assert"
)

func TestNodeBalancer_Create(t *testing.T) {
	tests := []struct {
		name       string
		createOpts linodego.NodeBalancerCreateOptions
		fixture    string
		expectIPv4 string
	}{
		{
			name: "basic creation",
			createOpts: linodego.NodeBalancerCreateOptions{
				Label:  String("Test NodeBalancer"),
				Region: "us-east",
				Tags:   []string{"test", "example"},
			},
			fixture:    "nodebalancer_create",
			expectIPv4: "192.0.2.1", // whatever is in the fixture
		},
		{
			name: "creation with specific IPv4",
			createOpts: linodego.NodeBalancerCreateOptions{
				Label:  String("Test NodeBalancer IPv4"),
				Region: "us-east",
				Tags:   []string{"test", "example"},
				IPv4:   String("192.0.2.2"),
			},
			fixture:    "nodebalancer_create_with_ipv4",
			expectIPv4: "192.0.2.2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fixtureData, err := fixtures.GetFixture(tt.fixture)
			assert.NoError(t, err)

			var base ClientBaseCase
			base.SetUp(t)
			defer base.TearDown(t)

			base.MockPost("nodebalancers", fixtureData)

			nodebalancer, err := base.Client.CreateNodeBalancer(context.Background(), tt.createOpts)
			assert.NoError(t, err)

			assert.Equal(t, *tt.createOpts.Label, *nodebalancer.Label)
			assert.Equal(t, tt.createOpts.Region, nodebalancer.Region)
			assert.Equal(t, tt.createOpts.Tags, nodebalancer.Tags)
			assert.Equal(t, tt.expectIPv4, *nodebalancer.IPv4)
		})
	}
}

// Helper function if not already defined
func String(s string) *string {
	return &s
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
		Tags:  &[]string{"updated", "production"},
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
