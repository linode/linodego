package unit

import (
	"context"
	"testing"

	"github.com/linode/linodego/v2"
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
		Label:      "Test Subnet",
		IPv4:       "192.168.1.0/24",
		NATGateway: linodego.Pointer(linodego.VPCSubnetCreateOptionsNATGateway{ID: linodego.DoublePointer(42)}),
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
	assert.Equal(t, 42, subnet.NATGateway.ID, "Expected subnet NAT gateway ID to match")
	assert.Equal(t, "my-nat-gateway", subnet.NATGateway.Label, "Expected subnet NAT gateway label to match")
	assert.Equal(t, []string{"203.0.113.42"}, subnet.NATGateway.Addresses, "Expected subnet NAT gateway addresses to match")
	assert.Equal(t, 15, subnet.NATGateway.PortsetAssignments, "Expected subnet NAT gateway portset assignments to match")
	assert.Equal(t, 30, subnet.NATGateway.PortsetCapacity, "Expected subnet NAT gateway portset capacity to match")

	assert.Len(t, subnet.NATGateway.Portsets, 1, "Expected subnet NAT gateway to have one portset")
	assert.Equal(t, "203.0.113.42", subnet.NATGateway.Portsets[0].Address, "Expected portset address to match")
	assert.Len(t, subnet.NATGateway.Portsets[0].Ports, 1, "Expected portset to have one port range")
	assert.Equal(t, 2048, subnet.NATGateway.Portsets[0].Ports[0].Start, "Expected port range start to match")
	assert.Equal(t, 3071, subnet.NATGateway.Portsets[0].Ports[0].End, "Expected port range end to match")
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

	assert.Equal(t, 123, subnet.ID, "Expected subnet ID to match")
	assert.Equal(t, "Subnet A", subnet.Label, "Expected subnet label to match")
	assert.Equal(t, "192.0.2.13/24", subnet.IPv4, "Expected subnet IPv4 to match")
	assert.Equal(t, "fd71:1140:a9d0::/52", subnet.IPv6[0].Range)

	assert.Equal(t, 111, subnet.Linodes[0].ID, "Expected Linode ID to match")

	assert.Equal(t, 421, subnet.Linodes[0].Interfaces[0].ID)
	assert.True(t, subnet.Linodes[0].Interfaces[0].Active, "Expected interface to be active")
	assert.Equal(t, 4567, *subnet.Linodes[0].Interfaces[0].ConfigID)

	assert.Equal(t, 422, subnet.Linodes[0].Interfaces[1].ID)
	assert.False(t, subnet.Linodes[0].Interfaces[1].Active, "Expected interface to be inactive")
	assert.Nil(t, subnet.Linodes[0].Interfaces[1].ConfigID)

	assert.Equal(t, 123, subnet.Databases[0].ID)
	assert.Equal(t, "10.0.0.4/32", *subnet.Databases[0].IPv4Range)
	assert.Equal(t, "fda3:9c1b:5e2a:1::/64", subnet.Databases[0].IPv6Ranges[0])

	assert.Equal(t, 4444, subnet.Nodebalancers[0].ID)
	assert.Equal(t, "192.168.0.20/30", subnet.Nodebalancers[0].Ipv4Range)
	assert.Equal(t, "2830:0:0:4::/62", subnet.Nodebalancers[0].Ipv6Ranges[0].Range)

	assert.Equal(t, 42, subnet.NATGateway.ID, "Expected subnet NAT gateway ID to match")
	assert.Equal(t, "my-nat-gateway", subnet.NATGateway.Label, "Expected subnet NAT gateway label to match")
	assert.Equal(t, []string{"203.0.113.42"}, subnet.NATGateway.Addresses, "Expected subnet NAT gateway addresses to match")
	assert.Equal(t, 15, subnet.NATGateway.PortsetAssignments, "Expected subnet NAT gateway portset assignments to match")
	assert.Equal(t, 30, subnet.NATGateway.PortsetCapacity, "Expected subnet NAT gateway portset capacity to match")

	assert.Len(t, subnet.NATGateway.Portsets, 1, "Expected subnet NAT gateway to have one portset")
	assert.Equal(t, "203.0.113.42", subnet.NATGateway.Portsets[0].Address, "Expected portset address to match")
	assert.Len(t, subnet.NATGateway.Portsets[0].Ports, 1, "Expected portset to have one port range")
	assert.Equal(t, 2048, subnet.NATGateway.Portsets[0].Ports[0].Start, "Expected port range start to match")
	assert.Equal(t, 3071, subnet.NATGateway.Portsets[0].Ports[0].End, "Expected port range end to match")
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
	assert.Len(t, subnets, 2, "Expected two subnets in the list")

	assertJSONObjectsSimilar(t, subnet, subnet.GetCreateOptions())
	assertJSONObjectsSimilar(t, subnet, subnet.GetUpdateOptions())

	assert.Equal(t, 123, subnet.ID, "Expected first subnet ID to match")
	assert.Equal(t, "Subnet A", subnet.Label, "Expected first subnet label to match")
	assert.Equal(t, "192.0.2.13/24", subnet.IPv4, "Expected first subnet IPv4 to match")
	assert.Equal(t, "fd71:1140:a9d0::/52", subnet.IPv6[0].Range, "Expected first subnet IPv6 to match")

	assert.Equal(t, 111, subnet.Linodes[0].ID)

	assert.Equal(t, true, subnet.Linodes[0].Interfaces[0].Active)
	assert.Equal(t, 4567, *subnet.Linodes[0].Interfaces[0].ConfigID)
	assert.Equal(t, 421, subnet.Linodes[0].Interfaces[0].ID)

	assert.Equal(t, false, subnet.Linodes[0].Interfaces[1].Active)
	assert.Nil(t, subnet.Linodes[0].Interfaces[1].ConfigID)
	assert.Equal(t, 422, subnet.Linodes[0].Interfaces[1].ID)

	assert.Equal(t, 43, subnet.NATGateway.ID, "Expected subnet NAT gateway ID to match")
	assert.Equal(t, "my-nat-gateway", subnet.NATGateway.Label, "Expected subnet NAT gateway label to match")
	assert.Equal(t, []string{"203.0.113.42"}, subnet.NATGateway.Addresses, "Expected subnet NAT gateway addresses to match")
	assert.Equal(t, 15, subnet.NATGateway.PortsetAssignments, "Expected subnet NAT gateway portset assignments to match")
	assert.Equal(t, 30, subnet.NATGateway.PortsetCapacity, "Expected subnet NAT gateway portset capacity to match")

	assert.Len(t, subnet.NATGateway.Portsets, 1, "Expected subnet NAT gateway to have one portset")
	assert.Equal(t, "203.0.113.42", subnet.NATGateway.Portsets[0].Address, "Expected portset address to match")
	assert.Len(t, subnet.NATGateway.Portsets[0].Ports, 1, "Expected portset to have one port range")
	assert.Equal(t, 2048, subnet.NATGateway.Portsets[0].Ports[0].Start, "Expected port range start to match")
	assert.Equal(t, 3071, subnet.NATGateway.Portsets[0].Ports[0].End, "Expected port range end to match")

	subnet = subnets[1]

	assert.Equal(t, 124, subnet.ID, "Expected second subnet ID to match")
	assert.Equal(t, "Subnet B", subnet.Label, "Expected second subnet label to match")
	assert.Equal(t, "192.168.4.0/24", subnet.IPv4, "Expected second subnet IPv4 to match")
	assert.Empty(t, subnet.IPv6, 0, "Expected second subnet to not support IPv6")
	assert.Empty(t, subnet.Linodes, 0, "Expected second subnet to not have Linodes")

	assert.Equal(t, 123, subnet.Databases[0].ID)
	assert.Equal(t, "10.0.0.4/32", *subnet.Databases[0].IPv4Range)
	assert.Equal(t, "fda3:9c1b:5e2a:1::/64", subnet.Databases[0].IPv6Ranges[0])

	assert.Equal(t, 4444, subnet.Nodebalancers[0].ID)
	assert.Equal(t, "192.168.0.20/30", subnet.Nodebalancers[0].Ipv4Range)
	assert.Equal(t, "2830:0:0:4::/62", subnet.Nodebalancers[0].Ipv6Ranges[0].Range)
}

func TestVPCSubnet_Update(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("vpc_subnet_update")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockPut("vpcs/123/subnets/456", fixtureData)

	subnetUpdateOpts := linodego.VPCSubnetUpdateOptions{
		Label:      "Updated Subnet",
		NATGateway: linodego.Pointer(linodego.VPCSubnetUpdateOptionsNATGateway{ID: linodego.DoublePointer(42)}),
	}
	subnet, err := base.Client.UpdateVPCSubnet(context.Background(), 123, 456, subnetUpdateOpts)
	assert.NoError(t, err)

	assertJSONObjectsSimilar(t, subnet, subnet.GetCreateOptions())
	assertJSONObjectsSimilar(t, subnet, subnet.GetUpdateOptions())

	assert.Equal(t, 456, subnet.ID, "Expected subnet ID to match")
	assert.Equal(t, "Updated Subnet", subnet.Label, "Expected subnet label to match")
	assert.Equal(t, 42, subnet.NATGateway.ID, "Expected subnet NAT gateway ID to match")
	assert.Equal(t, "my-nat-gateway", subnet.NATGateway.Label, "Expected subnet NAT gateway label to match")
	assert.Equal(t, []string{"203.0.113.42"}, subnet.NATGateway.Addresses, "Expected subnet NAT gateway addresses to match")
	assert.Equal(t, 15, subnet.NATGateway.PortsetAssignments, "Expected subnet NAT gateway portset assignments to match")
	assert.Equal(t, 30, subnet.NATGateway.PortsetCapacity, "Expected subnet NAT gateway portset capacity to match")

	assert.Len(t, subnet.NATGateway.Portsets, 1, "Expected subnet NAT gateway to have one portset")
	assert.Equal(t, "203.0.113.42", subnet.NATGateway.Portsets[0].Address, "Expected portset address to match")
	assert.Len(t, subnet.NATGateway.Portsets[0].Ports, 1, "Expected portset to have one port range")
	assert.Equal(t, 2048, subnet.NATGateway.Portsets[0].Ports[0].Start, "Expected port range start to match")
	assert.Equal(t, 3071, subnet.NATGateway.Portsets[0].Ports[0].End, "Expected port range end to match")
}

func TestVPCSubnet_Delete(t *testing.T) {
	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockDelete("vpcs/123/subnets/456", nil)

	err := base.Client.DeleteVPCSubnet(context.Background(), 123, 456)
	assert.NoError(t, err, "Expected no error when deleting VPCSubnet")
}
