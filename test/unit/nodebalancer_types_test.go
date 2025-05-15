package unit

import (
	"context"
	"testing"

	"github.com/linode/linodego"
	"github.com/stretchr/testify/assert"
)

func TestNodeBalancerTypes_List(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("nodebalancers_types_list")
	assert.NoError(t, err, "Expected no error when getting fixture")

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	// Mock the GET request for the node balancer types endpoint (mock as an array instead of paginatedResponse)
	base.MockGet("nodebalancers/types", fixtureData)

	// Call the ListNodeBalancerTypes method
	nodebalancerTypes, err := base.Client.ListNodeBalancerTypes(context.Background(), &linodego.ListOptions{})
	assert.NoError(t, err, "Expected no error when listing node balancer types")
	assert.NotEmpty(t, nodebalancerTypes, "Expected non-empty node balancer types list")

	// Validate the first node balancer type's details
	assert.Equal(t, "123", nodebalancerTypes[0].ID, "Expected node balancer type ID to match")
	assert.Equal(t, "NodeBalancer A", nodebalancerTypes[0].Label, "Expected node balancer type label to match")
	assert.Equal(t, 0.10, nodebalancerTypes[0].Price.Hourly, "Expected hourly price to match")
	assert.Equal(t, 10.00, nodebalancerTypes[0].Price.Monthly, "Expected monthly price to match")

	// Validate the second node balancer type's details
	assert.Equal(t, "456", nodebalancerTypes[1].ID, "Expected node balancer type ID to match")
	assert.Equal(t, "NodeBalancer B", nodebalancerTypes[1].Label, "Expected node balancer type label to match")
	assert.Equal(t, 0.15, nodebalancerTypes[1].Price.Hourly, "Expected hourly price to match")
	assert.Equal(t, 15.00, nodebalancerTypes[1].Price.Monthly, "Expected monthly price to match")

	// Access RegionPrice correctly from the embedded struct
	assert.NotEmpty(t, nodebalancerTypes[1].Price, "Expected price to be non-empty")
}
