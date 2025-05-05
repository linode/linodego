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

	base.MockGet("vpcs/123/subnets/456", fixtureData)

	subnet, err := base.Client.GetVPCSubnet(context.Background(), 123, 456)
	assert.NoError(t, err)

	assert.Equal(t, 456, subnet.ID)
	assert.Equal(t, "10.0.1.0/24", subnet.IPv4)
	assert.Equal(t, "cool-vpc-subnet", subnet.Label)
	assert.Equal(t, 111, subnet.Linodes[0].ID)
	assert.Equal(t, true, subnet.Linodes[0].Interfaces[0].Active)
	assert.Equal(t, 4567, *subnet.Linodes[0].Interfaces[0].ConfigID)
	assert.Equal(t, 421, subnet.Linodes[0].Interfaces[0].ID)
}

func TestVPCSubnets_List(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("vpc_subnets_list")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("vpcs/123/subnets", fixtureData)

	subnets, err := base.Client.ListVPCSubnets(context.Background(), 123, &linodego.ListOptions{})
	subnet := subnets[0]
	assert.NoError(t, err, "Expected no error when listing subnets")
	assert.Len(t, subnets, 1, "Expected two subnets in the list")

	assert.Equal(t, 456, subnet.ID)
	assert.Equal(t, "192.0.2.13/24", subnet.IPv4)
	assert.Equal(t, "cool-vpc-subnet", subnet.Label)
	assert.Equal(t, 111, subnet.Linodes[0].ID)
	assert.Equal(t, true, subnet.Linodes[0].Interfaces[0].Active)
	assert.Nil(t, subnet.Linodes[0].Interfaces[0].ConfigID)
	assert.Equal(t, 421, subnet.Linodes[0].Interfaces[0].ID)
}

func TestVPCSubnet_Update(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("vpc_subnet_update")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

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

	base.MockDelete("vpcs/123/subnets/456", nil)

	err := base.Client.DeleteVPCSubnet(context.Background(), 123, 456)
	assert.NoError(t, err, "Expected no error when deleting VPCSubnet")
}
