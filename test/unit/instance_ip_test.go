package unit

import (
	"context"
	"testing"

	"github.com/linode/linodego"
	"github.com/stretchr/testify/assert"
)

func TestInstanceIPAddresses_List(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("instance_ip_list")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("linode/instances/123/ips", fixtureData)

	ips, err := base.Client.GetInstanceIPAddresses(context.Background(), 123)
	assert.NoError(t, err)
	assert.NotNil(t, ips)
	assert.NotNil(t, ips.IPv4)
	assert.NotNil(t, ips.IPv6)

	// IPv4 Assertions
	assert.Len(t, ips.IPv4.Public, 1)
	assert.Equal(t, "192.0.2.1", ips.IPv4.Public[0].Address)

	// IPv6 Assertions
	assert.NotNil(t, ips.IPv6.SLAAC)
	assert.Equal(t, "2001:db8::1", ips.IPv6.SLAAC.Address)
}

func TestInstanceIPAddress_Get(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("instance_ip_get")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("linode/instances/123/ips/192.0.2.1", fixtureData)

	ip, err := base.Client.GetInstanceIPAddress(context.Background(), 123, "192.0.2.1")
	assert.NoError(t, err)
	assert.NotNil(t, ip)
	assert.Equal(t, "192.0.2.1", ip.Address)
	assert.Equal(t, "192.0.2.254", ip.Gateway)
	assert.Equal(t, 24, ip.Prefix)
	assert.True(t, ip.Public)
}

func TestInstanceIPAddress_Add(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("instance_ip_add")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockPost("linode/instances/123/ips", fixtureData)

	ip, err := base.Client.AddInstanceIPAddress(context.Background(), 123, true)
	assert.NoError(t, err)
	assert.NotNil(t, ip)
	assert.Equal(t, "198.51.100.1", ip.Address)
	assert.True(t, ip.Public)
}

func TestInstanceIPAddress_Update(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("instance_ip_update")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	rdns := "custom.reverse.dns"
	updateOpts := linodego.IPAddressUpdateOptions{
		RDNS: &rdns,
	}

	base.MockPut("linode/instances/123/ips/192.0.2.1", fixtureData)

	ip, err := base.Client.UpdateInstanceIPAddress(context.Background(), 123, "192.0.2.1", updateOpts)
	assert.NoError(t, err)
	assert.NotNil(t, ip)
	assert.Equal(t, "custom.reverse.dns", ip.RDNS)
}

func TestInstanceIPAddress_Delete(t *testing.T) {
	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockDelete("linode/instances/123/ips/192.0.2.1", nil)

	err := base.Client.DeleteInstanceIPAddress(context.Background(), 123, "192.0.2.1")
	assert.NoError(t, err)
}

func TestInstanceReservedIP_Assign(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("instance_ip_reserved")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	opts := linodego.InstanceReserveIPOptions{
		Type:    "ipv4",
		Public:  false,
		Address: "203.0.113.1",
	}

	base.MockPost("linode/instances/123/ips", fixtureData)

	ip, err := base.Client.AssignInstanceReservedIP(context.Background(), 123, opts)
	assert.NoError(t, err)
	assert.NotNil(t, ip)
	assert.Equal(t, "203.0.113.1", ip.Address)
	assert.False(t, ip.Public)
}
