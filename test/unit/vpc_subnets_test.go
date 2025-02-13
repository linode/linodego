package unit

import (
	"context"
	"testing"

	"github.com/linode/linodego"
	"github.com/stretchr/testify/assert"
)

func TestVPCSubnet_Create(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("vpc_subnet_create")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	// Mock the POST request with the fixture response
	base.MockPost("vpcs/123/subnets", fixtureData)

	subnetCreateOpts := linodego.VPCSubnetCreateOptions{
		Label: "Test Subnet",
		IPv4:  "192.168.1.0/24",
	}
	subnet, err := base.Client.CreateVPCSubnet(context.Background(), subnetCreateOpts, 123)
	assert.NoError(t, err)

	assert.Equal(t, 789, subnet.ID, "Expected subnet ID to match")
	assert.Equal(t, "Test Subnet", subnet.Label, "Expected subnet label to match")
	assert.Equal(t, "192.168.1.0/24", subnet.IPv4, "Expected subnet IPv4 to match")
}

func TestVPCSubnet_Get(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("vpc_subnet_get")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	// Mock the GET request with the fixture response
	base.MockGet("vpcs/123/subnets/456", fixtureData)

	subnet, err := base.Client.GetVPCSubnet(context.Background(), 123, 456)
	assert.NoError(t, err)

	assert.Equal(t, 456, subnet.ID, "Expected subnet ID to match")
	assert.Equal(t, "Existing Subnet", subnet.Label, "Expected subnet label to match")
	assert.Equal(t, "192.168.2.0/24", subnet.IPv4, "Expected subnet IPv4 to match")
	assert.Equal(t, 101, subnet.Linodes[0].ID, "Expected Linode ID to match")
	assert.True(t, subnet.Linodes[0].Interfaces[0].Active, "Expected interface to be active")
}

func TestVPCSubnets_List(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("vpc_subnets_list")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	// Mock the GET request with the fixture response
	base.MockGet("vpcs/123/subnets", fixtureData)

	subnets, err := base.Client.ListVPCSubnets(context.Background(), 123, &linodego.ListOptions{})
	assert.NoError(t, err, "Expected no error when listing subnets")
	assert.Len(t, subnets, 2, "Expected two subnets in the list")

	assert.Equal(t, 123, subnets[0].ID, "Expected first subnet ID to match")
	assert.Equal(t, "Subnet A", subnets[0].Label, "Expected first subnet label to match")
	assert.Equal(t, "192.168.3.0/24", subnets[0].IPv4, "Expected first subnet IPv4 to match")
}

func TestVPCSubnet_Update(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("vpc_subnet_update")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	// Mock the PUT request with the fixture response
	base.MockPut("vpcs/123/subnets/456", fixtureData)

	subnetUpdateOpts := linodego.VPCSubnetUpdateOptions{
		Label: "Updated Subnet",
	}
	subnet, err := base.Client.UpdateVPCSubnet(context.Background(), 123, 456, subnetUpdateOpts)
	assert.NoError(t, err)

	assert.Equal(t, 456, subnet.ID, "Expected subnet ID to match")
	assert.Equal(t, "Updated Subnet", subnet.Label, "Expected subnet label to match")
}

func TestVPCSubnet_Delete(t *testing.T) {
	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	// Mock the DELETE request
	base.MockDelete("vpcs/123/subnets/456",nil)

	err := base.Client.DeleteVPCSubnet(context.Background(), 123, 456)
	assert.NoError(t, err,"Expected no error when deleting VPCSubnet")
}
