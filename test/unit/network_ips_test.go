package unit

import (
	"context"
	"testing"

	"github.com/linode/linodego"
	"github.com/stretchr/testify/assert"
)

func TestIPUpdateAddressV2(t *testing.T) {
	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	ip := "192.168.1.1"

	// Mock API response
	base.MockPut("networking/ips/"+ip, linodego.InstanceIP{
		Address:  ip,
		Reserved: true,
	})

	updatedIP, err := base.Client.UpdateIPAddressV2(context.Background(), ip, linodego.IPAddressUpdateOptionsV2{
		Reserved: linodego.Pointer(true),
	})
	assert.NoError(t, err, "Expected no error when updating IP address")
	assert.NotNil(t, updatedIP, "Expected non-nil updated IP address")
	assert.Equal(t, ip, updatedIP.Address, "Expected updated IP address to match")
	assert.True(t, updatedIP.Reserved, "Expected Reserved to be true")
}

func TestIPAllocateReserve(t *testing.T) {
	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	// Mock API response
	base.MockPost("networking/ips", linodego.InstanceIP{
		Address: "192.168.1.3",
		Region:  "us-east",
		Public:  true,
	})

	ip, err := base.Client.AllocateReserveIP(context.Background(), linodego.AllocateReserveIPOptions{
		Type:     "ipv4",
		Public:   true,
		Region:   "us-east",
		LinodeID: 12345,
	})
	assert.NoError(t, err, "Expected no error when allocating reserve IP")
	assert.NotNil(t, ip, "Expected non-nil allocated IP")
	assert.Equal(t, "192.168.1.3", ip.Address, "Expected allocated IP address to match")
	assert.Equal(t, "us-east", ip.Region, "Expected Region to match")
	assert.True(t, ip.Public, "Expected Public to be true")
	assert.Nil(t, ip.InterfaceID)
}

func TestIPAssignInstances(t *testing.T) {
	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	// Mock API response
	base.MockPost("networking/ips/assign", nil)

	err := base.Client.InstancesAssignIPs(context.Background(), linodego.LinodesAssignIPsOptions{
		Region: "us-east",
		Assignments: []linodego.LinodeIPAssignment{
			{Address: "192.168.1.10", LinodeID: 123},
		},
	})
	assert.NoError(t, err, "Expected no error when assigning IPs to instances")
}

func TestIPShareAddresses(t *testing.T) {
	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	// Mock API response
	base.MockPost("networking/ips/share", nil)

	err := base.Client.ShareIPAddresses(context.Background(), linodego.IPAddressesShareOptions{
		IPs:      []string{"192.168.1.20"},
		LinodeID: 456,
	})
	assert.NoError(t, err, "Expected no error when sharing IP addresses")
}

func TestIPAddresses_List(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("network_ip_addresses_list")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("networking/ips", fixtureData)

	ips, err := base.Client.ListIPAddresses(context.Background(), nil)
	assert.NoError(t, err)
	assert.NotEmpty(t, ips, "Expected ips to be returned.")

	ip := ips[0]
	assert.Equal(t, "197.1O7.143.141", ip.Address)
	assert.Equal(t, "197.1O7.143.1", ip.Gateway)
	assert.Equal(t, 456, *ip.InterfaceID)
	assert.Equal(t, 123, ip.LinodeID)
	assert.Equal(t, 24, ip.Prefix)
	assert.Equal(t, true, ip.Public)
	assert.Equal(t, "test.example.org", ip.RDNS)
	assert.Equal(t, "us-east", ip.Region)
	assert.Equal(t, "192.0.2.139", ip.SubnetMask)
	assert.Equal(t, linodego.InstanceIPType("ipv4"), ip.Type)
	assert.Equal(t, "192.0.2.1", ip.VPCNAT1To1.Address)
	assert.Equal(t, 101, ip.VPCNAT1To1.SubnetID)
	assert.Equal(t, 111, ip.VPCNAT1To1.VPCID)
}

func TestIPAddresses_Get(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("network_ip_address_get")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("networking/ips/97.107.143.141", fixtureData)

	ip, err := base.Client.GetIPAddress(context.Background(), "97.107.143.141")
	assert.NoError(t, err)

	assert.Equal(t, "97.107.143.141", ip.Address)
	assert.Equal(t, "97.107.143.1", ip.Gateway)
	assert.Equal(t, 456, *ip.InterfaceID)
	assert.Equal(t, 123, ip.LinodeID)
	assert.Equal(t, 24, ip.Prefix)
	assert.Equal(t, true, ip.Public)
	assert.Equal(t, "test.example.org", ip.RDNS)
	assert.Equal(t, "us-east", ip.Region)
	assert.Equal(t, "255.255.255.0", ip.SubnetMask)
	assert.Equal(t, linodego.InstanceIPType("ipv4"), ip.Type)
	assert.Equal(t, "192.168.0.42", ip.VPCNAT1To1.Address)
	assert.Equal(t, 101, ip.VPCNAT1To1.SubnetID)
	assert.Equal(t, 111, ip.VPCNAT1To1.VPCID)
}
