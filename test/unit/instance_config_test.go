package unit

import (
	"context"
	"testing"

	"github.com/linode/linodego"
	"github.com/stretchr/testify/assert"
)

func TestInstanceConfig_List(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("instance_config_list")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("linode/instances/123/configs", fixtureData)

	configs, err := base.Client.ListInstanceConfigs(context.Background(), 123, nil)
	assert.NoError(t, err)
	assert.Len(t, configs, 2)

	assert.Equal(t, 1, configs[0].ID)
	assert.Equal(t, "config-1", configs[0].Label)

	iface := configs[0].Interfaces[0]
	assert.Equal(t, 1, iface.ID)
	assert.Equal(t, "vpc", string(iface.Purpose))
	assert.Equal(t, true, iface.Primary)
	assert.Equal(t, true, iface.Active)
	assert.Equal(t, 101, *iface.VPCID)
	assert.Equal(t, 202, *iface.SubnetID)

	assert.Equal(t, "vpc-1", iface.IPv4.VPC)
	assert.Equal(t, "203.0.113.1", *iface.IPv4.NAT1To1)

	assert.Len(t, iface.IPv6.SLAAC, 1)
	assert.Equal(t, "1234::5678/64", iface.IPv6.SLAAC[0].Range)
	assert.Equal(t, "1234::5678", iface.IPv6.SLAAC[0].Address)

	assert.Len(t, iface.IPv6.Ranges, 1)
	assert.Equal(t, "1234::5678/64", iface.IPv6.Ranges[0].Range)

	assert.Equal(t, true, iface.IPv6.IsPublic)

	assert.ElementsMatch(t, []string{"192.168.1.0/24"}, iface.IPRanges)
}

func TestInstanceConfig_Get(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("instance_config_get")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("linode/instances/123/configs/1", fixtureData)

	config, err := base.Client.GetInstanceConfig(context.Background(), 123, 1)
	assert.NoError(t, err)

	assertJSONObjectsSimilar(t, config, config.GetCreateOptions())
	assertJSONObjectsSimilar(t, config, config.GetUpdateOptions())

	assert.Equal(t, 1, config.ID)
	assert.Equal(t, "config-1", config.Label)

	iface := config.Interfaces[0]
	assert.Equal(t, 1, iface.ID)
	assert.Equal(t, "vpc", string(iface.Purpose))
	assert.Equal(t, true, iface.Primary)
	assert.Equal(t, true, iface.Active)
	assert.Equal(t, 101, *iface.VPCID)
	assert.Equal(t, 202, *iface.SubnetID)

	assert.Equal(t, "vpc-1", iface.IPv4.VPC)
	assert.Equal(t, "203.0.113.1", *iface.IPv4.NAT1To1)

	assert.Len(t, iface.IPv6.SLAAC, 1)
	assert.Equal(t, "1234::5678/64", iface.IPv6.SLAAC[0].Range)
	assert.Equal(t, "1234::5678", iface.IPv6.SLAAC[0].Address)

	assert.Len(t, iface.IPv6.Ranges, 1)
	assert.Equal(t, "1234::5678/64", iface.IPv6.Ranges[0].Range)

	assert.Equal(t, true, iface.IPv6.IsPublic)

	assert.ElementsMatch(t, []string{"192.168.1.0/24"}, iface.IPRanges)
}

func TestInstanceConfig_Create(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("instance_config_create")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	rootDevice := "/dev/sda"

	createOptions := linodego.InstanceConfigCreateOptions{
		Label:      "new-config",
		Kernel:     "linode/latest-64bit",
		RootDevice: &rootDevice,
	}

	base.MockPost("linode/instances/123/configs", fixtureData)

	config, err := base.Client.CreateInstanceConfig(context.Background(), 123, createOptions)
	assert.NoError(t, err)

	assertJSONObjectsSimilar(t, config, config.GetCreateOptions())
	assertJSONObjectsSimilar(t, config, config.GetUpdateOptions())

	assert.Equal(t, "new-config", config.Label)

	iface := config.Interfaces[0]
	assert.Equal(t, 1, iface.ID)
	assert.Equal(t, "vpc", string(iface.Purpose))
	assert.Equal(t, true, iface.Primary)
	assert.Equal(t, true, iface.Active)
	assert.Equal(t, 101, *iface.VPCID)
	assert.Equal(t, 202, *iface.SubnetID)

	assert.Equal(t, "vpc-1", iface.IPv4.VPC)
	assert.Equal(t, "203.0.113.1", *iface.IPv4.NAT1To1)

	assert.Len(t, iface.IPv6.SLAAC, 1)
	assert.Equal(t, "1234::5678/64", iface.IPv6.SLAAC[0].Range)
	assert.Equal(t, "1234::5678", iface.IPv6.SLAAC[0].Address)

	assert.Len(t, iface.IPv6.Ranges, 1)
	assert.Equal(t, "1234::5678/64", iface.IPv6.Ranges[0].Range)

	assert.Equal(t, true, iface.IPv6.IsPublic)

	assert.ElementsMatch(t, []string{"192.168.1.0/24"}, iface.IPRanges)
}

func TestInstanceConfig_Update(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("instance_config_update")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	updateOptions := linodego.InstanceConfigUpdateOptions{
		Label:      "updated-config",
		RootDevice: "/dev/sdb",
	}

	base.MockPut("linode/instances/123/configs/1", fixtureData)

	config, err := base.Client.UpdateInstanceConfig(context.Background(), 123, 1, updateOptions)
	assert.NoError(t, err)

	assertJSONObjectsSimilar(t, config, config.GetCreateOptions())
	assertJSONObjectsSimilar(t, config, config.GetUpdateOptions())

	assert.Equal(t, "updated-config", config.Label)

	iface := config.Interfaces[0]
	assert.Equal(t, 1, iface.ID)
	assert.Equal(t, "vpc", string(iface.Purpose))
	assert.Equal(t, true, iface.Primary)
	assert.Equal(t, true, iface.Active)
	assert.Equal(t, 101, *iface.VPCID)
	assert.Equal(t, 202, *iface.SubnetID)

	assert.Equal(t, "vpc-1", iface.IPv4.VPC)
	assert.Equal(t, "203.0.113.1", *iface.IPv4.NAT1To1)

	assert.Len(t, iface.IPv6.SLAAC, 1)
	assert.Equal(t, "1234::5678/64", iface.IPv6.SLAAC[0].Range)
	assert.Equal(t, "1234::5678", iface.IPv6.SLAAC[0].Address)

	assert.Len(t, iface.IPv6.Ranges, 1)
	assert.Equal(t, "1234::5678/64", iface.IPv6.Ranges[0].Range)

	assert.Equal(t, true, iface.IPv6.IsPublic)

	assert.ElementsMatch(t, []string{"192.168.1.0/24"}, iface.IPRanges)
}

func TestInstanceConfig_Delete(t *testing.T) {
	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockDelete("linode/instances/123/configs/1", nil)

	err := base.Client.DeleteInstanceConfig(context.Background(), 123, 1)
	assert.NoError(t, err)
}
