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
}
