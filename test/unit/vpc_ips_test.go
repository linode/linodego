package unit

import (
	"context"
	"testing"

	"github.com/linode/linodego"
	"github.com/stretchr/testify/assert"
)

func TestVPCAllIPAddresses_List(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("vpc_ips_list")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("vpcs/ips", fixtureData)

	vpcIPs, err := base.Client.ListAllVPCIPAddresses(context.Background(), &linodego.ListOptions{})
	assert.NoError(t, err)

	assert.NotEmpty(t, vpcIPs, "Expected non-empty VPC IP addresses list")

	assert.NotNil(t, vpcIPs[0].Address, "Expected IP address to be non-nil")
	assert.Equal(t, "192.168.1.10", *vpcIPs[0].Address, "Expected IP address to match")
	assert.Equal(t, 123, vpcIPs[0].VPCID, "Expected VPC ID to match")

	assert.Equal(t, "fd71:1140:a9d0::/52", *vpcIPs[2].IPv6Range)
	assert.Equal(t, true, *vpcIPs[2].IPv6IsPublic)
	assert.Equal(t, "fd71:1140:a9d0::/52", vpcIPs[2].IPv6Addresses[0].SLAACAddress)
	assert.Equal(t, 125, vpcIPs[2].VPCID)
}

func TestVPCSpecificIPAddresses_List(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("vpc_specific_ips_list")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	vpcID := 123
	base.MockGet("vpcs/123/ips", fixtureData)

	vpcIPs, err := base.Client.ListVPCIPAddresses(context.Background(), vpcID, &linodego.ListOptions{})
	assert.NoError(t, err)

	assert.NotEmpty(t, vpcIPs, "Expected non-empty VPC IP addresses list for the specified VPC")

	assert.NotNil(t, vpcIPs[0].Address, "Expected IP address to be non-nil")
	assert.Equal(t, "192.168.1.20", *vpcIPs[0].Address, "Expected IP address to match")
	assert.Equal(t, vpcID, vpcIPs[0].VPCID, "Expected VPC ID to match")

	assert.Equal(t, "fd71:1140:a9d0::/52", *vpcIPs[2].IPv6Range)
	assert.Equal(t, true, *vpcIPs[2].IPv6IsPublic)
	assert.Equal(t, "fd71:1140:a9d0::/52", vpcIPs[2].IPv6Addresses[0].SLAACAddress)
	assert.Equal(t, 123, vpcIPs[2].VPCID)
}

func TestVPCSpecificIPv6Addresses_List(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("vpc_specific_ipv6s_list")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	vpcID := 123
	base.MockGet("vpcs/123/ipv6s", fixtureData)

	vpcIPs, err := base.Client.ListVPCIPv6Addresses(context.Background(), vpcID, &linodego.ListOptions{})
	assert.NoError(t, err)

	assert.NotEmpty(t, vpcIPs)

	assert.Equal(t, vpcID, vpcIPs[0].VPCID, "Expected VPC ID to match")
	assert.Nil(t, vpcIPs[0].Address)

	assert.True(t, *vpcIPs[0].IPv6IsPublic)
	assert.Equal(t, "1234::5678/64", *vpcIPs[0].IPv6Range)
	assert.Equal(t, "1234::5678", vpcIPs[0].IPv6Addresses[0].SLAACAddress)
}

func TestVPCAllIPv6Addresses_List(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("vpc_ipv6s_list")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	vpcID := 123
	base.MockGet("vpcs/ipv6s", fixtureData)

	vpcIPs, err := base.Client.ListAllVPCIPv6Addresses(context.Background(), &linodego.ListOptions{})
	assert.NoError(t, err)

	assert.NotEmpty(t, vpcIPs)

	assert.Equal(t, vpcID, vpcIPs[0].VPCID, "Expected VPC ID to match")
	assert.Nil(t, vpcIPs[0].Address)

	assert.True(t, *vpcIPs[0].IPv6IsPublic)
	assert.Equal(t, "1234::5678/64", *vpcIPs[0].IPv6Range)
	assert.Equal(t, "1234::5678", vpcIPs[0].IPv6Addresses[0].SLAACAddress)
}
