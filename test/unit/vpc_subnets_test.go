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
		IPv6: []linodego.VPCSubnetCreateOptionsIPv6{
			{
				linodego.Pointer("auto"),
			},
		},
	}

	subnet, err := base.Client.CreateVPCSubnet(context.Background(), subnetCreateOpts, 123)
	assert.NoError(t, err)

	assertJSONObjectsSimilar(t, subnet, subnet.GetCreateOptions())
	assertJSONObjectsSimilar(t, subnet, subnet.GetUpdateOptions())

	assert.Equal(t, 789, subnet.ID, "Expected subnet ID to match")
	assert.Equal(t, "Test Subnet", subnet.Label, "Expected subnet label to match")
	assert.Equal(t, "192.168.1.0/24", subnet.IPv4, "Expected subnet IPv4 to match")
	assert.Equal(t, "fd71:1140:a9d0::/52", subnet.IPv6[0].Range)
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

	assertJSONObjectsSimilar(t, subnet, subnet.GetCreateOptions())
	assertJSONObjectsSimilar(t, subnet, subnet.GetUpdateOptions())

	assert.Equal(t, 456, subnet.ID, "Expected subnet ID to match")
	assert.Equal(t, "Existing Subnet", subnet.Label, "Expected subnet label to match")
	assert.Equal(t, "192.168.2.0/24", subnet.IPv4, "Expected subnet IPv4 to match")
	assert.Equal(t, "fd71:1140:a9d0::/52", subnet.IPv6[0].Range)

	assert.Equal(t, 101, subnet.Linodes[0].ID, "Expected Linode ID to match")

	assert.Equal(t, 1, subnet.Linodes[0].Interfaces[0].ID)
	assert.True(t, subnet.Linodes[0].Interfaces[0].Active, "Expected interface to be active")
	assert.Equal(t, 4567, *subnet.Linodes[0].Interfaces[0].ConfigID)

	assert.Equal(t, 2, subnet.Linodes[0].Interfaces[1].ID)
	assert.False(t, subnet.Linodes[0].Interfaces[1].Active, "Expected interface to be inactive")
	assert.Equal(t, 4567, *subnet.Linodes[0].Interfaces[1].ConfigID)
}

func TestVPCSubnets_List(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("vpc_subnets_list")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("vpcs/123/subnets", fixtureData)

	subnets, err := base.Client.ListVPCSubnets(context.Background(), 123, &linodego.ListOptions{})
	assert.NoError(t, err, "Expected no error when listing subnets")
	assert.Len(t, subnets, 2, "Expected two subnets in the list")

	subnet := subnets[0]

	assertJSONObjectsSimilar(t, subnet, subnet.GetCreateOptions())
	assertJSONObjectsSimilar(t, subnet, subnet.GetUpdateOptions())

	assert.Equal(t, 123, subnet.ID, "Expected first subnet ID to match")
	assert.Equal(t, "Subnet A", subnet.Label, "Expected first subnet label to match")
	assert.Equal(t, "192.168.3.0/24", subnet.IPv4, "Expected first subnet IPv4 to match")
	assert.Equal(t, "fd71:1140:a9d0::/52", subnet.IPv6[0].Range, "Expected first subnet IPv6 to match")

	assert.Equal(t, 111, subnet.Linodes[0].ID)
	assert.Equal(t, true, subnet.Linodes[0].Interfaces[0].Active)
	assert.Nil(t, subnet.Linodes[0].Interfaces[0].ConfigID)
	assert.Equal(t, 421, subnet.Linodes[0].Interfaces[0].ID)

	subnet = subnets[1]

	assert.Equal(t, 124, subnet.ID, "Expected second subnet ID to match")
	assert.Equal(t, "Subnet B", subnet.Label, "Expected second subnet label to match")
	assert.Equal(t, "192.168.4.0/24", subnet.IPv4, "Expected second subnet IPv4 to match")
	assert.Empty(t, subnet.IPv6, 0, "Expected second subnet to not support IPv6")
	assert.Empty(t, subnet.Linodes, 0, "Expected second subnet to not have Linodes")
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

	assertJSONObjectsSimilar(t, subnet, subnet.GetCreateOptions())
	assertJSONObjectsSimilar(t, subnet, subnet.GetUpdateOptions())

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
