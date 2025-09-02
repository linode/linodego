package unit

import (
	"context"
	"testing"

	"github.com/linode/linodego"
	"github.com/stretchr/testify/assert"
)

func TestInstanceConfigInterface_List(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("instance_config_interface_list")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("linode/instances/123/configs/456/interfaces", fixtureData)

	interfaces, err := base.Client.ListInstanceConfigInterfaces(context.Background(), 123, 456)
	assert.NoError(t, err)
	assert.Len(t, interfaces, 2)

	assertJSONObjectsSimilar(t, interfaces[0], interfaces[0].GetCreateOptions())
	assertJSONObjectsSimilar(t, interfaces[0], interfaces[0].GetUpdateOptions())

	assert.Equal(t, 1, interfaces[0].ID)
	assert.Equal(t, "eth0", interfaces[0].Label)

	assert.Equal(t, "1234::5678/64", interfaces[0].IPv6.SLAAC[0].Range)
	assert.Equal(t, "1234::5678", interfaces[0].IPv6.SLAAC[0].Address)
	assert.Equal(t, "1234::5678/64", interfaces[0].IPv6.Ranges[0].Range)
	assert.True(t, interfaces[0].IPv6.IsPublic)
}

func TestInstanceConfigInterface_Get(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("instance_config_interface_get")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("linode/instances/123/configs/456/interfaces/1", fixtureData)

	iface, err := base.Client.GetInstanceConfigInterface(context.Background(), 123, 456, 1)
	assert.NoError(t, err)

	assertJSONObjectsSimilar(t, iface, iface.GetCreateOptions())
	assertJSONObjectsSimilar(t, iface, iface.GetUpdateOptions())

	assert.Equal(t, 1, iface.ID)
	assert.Equal(t, "eth0", iface.Label)

	assert.Equal(t, "1234::5678/64", iface.IPv6.SLAAC[0].Range)
	assert.Equal(t, "1234::5678", iface.IPv6.SLAAC[0].Address)
	assert.Equal(t, "1234::5678/64", iface.IPv6.Ranges[0].Range)
	assert.True(t, iface.IPv6.IsPublic)
}

func TestInstanceConfigInterface_Create(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("instance_config_interface_create")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	subnetID := 123
	nat1to1 := "192.168.1.1"

	createOptions := linodego.InstanceConfigInterfaceCreateOptions{
		Label:    "eth0",
		Purpose:  linodego.InterfacePurposeVPC,
		Primary:  true,
		SubnetID: &subnetID,
		IPv4: &linodego.VPCIPv4{
			NAT1To1: &nat1to1,
		},
		IPRanges: []string{"192.168.1.0/24"},
	}

	base.MockPost("linode/instances/123/configs/456/interfaces", fixtureData)

	iface, err := base.Client.AppendInstanceConfigInterface(context.Background(), 123, 456, createOptions)
	assert.NoError(t, err)

	assertJSONObjectsSimilar(t, iface, iface.GetCreateOptions())
	assertJSONObjectsSimilar(t, iface, iface.GetUpdateOptions())

	assert.Equal(t, "eth0", iface.Label)
	assert.True(t, iface.Primary)

	assert.Equal(t, "1234::5678/64", iface.IPv6.SLAAC[0].Range)
	assert.Equal(t, "1234::5678", iface.IPv6.SLAAC[0].Address)
	assert.Equal(t, "1234::5678/64", iface.IPv6.Ranges[0].Range)
	assert.True(t, iface.IPv6.IsPublic)
}

func TestInstanceConfigInterface_Update(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("instance_config_interface_update")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	nat1to1 := "192.168.1.1"
	ipRanges := []string{"192.168.1.0/24"}

	updateOptions := linodego.InstanceConfigInterfaceUpdateOptions{
		Primary: true,
		IPv4: &linodego.VPCIPv4{
			NAT1To1: &nat1to1,
		},
		IPRanges: &ipRanges,
	}

	base.MockPut("linode/instances/123/configs/456/interfaces/1", fixtureData)

	iface, err := base.Client.UpdateInstanceConfigInterface(context.Background(), 123, 456, 1, updateOptions)
	assert.NoError(t, err)

	assertJSONObjectsSimilar(t, iface, iface.GetCreateOptions())
	assertJSONObjectsSimilar(t, iface, iface.GetUpdateOptions())

	assert.True(t, iface.Primary)

	assert.Equal(t, "1234::5678/64", iface.IPv6.SLAAC[0].Range)
	assert.Equal(t, "1234::5678", iface.IPv6.SLAAC[0].Address)
	assert.Equal(t, "1234::5678/64", iface.IPv6.Ranges[0].Range)
	assert.True(t, iface.IPv6.IsPublic)
}

func TestInstanceConfigInterface_Delete(t *testing.T) {
	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockDelete("linode/instances/123/configs/456/interfaces/1", nil)

	err := base.Client.DeleteInstanceConfigInterface(context.Background(), 123, 456, 1)
	assert.NoError(t, err)
}

func TestInstanceConfigInterface_Reorder(t *testing.T) {
	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	reorderOptions := linodego.InstanceConfigInterfacesReorderOptions{
		IDs: []int{3, 1, 2},
	}

	base.MockPost("linode/instances/123/configs/456/interfaces/order", nil)

	err := base.Client.ReorderInstanceConfigInterfaces(context.Background(), 123, 456, reorderOptions)
	assert.NoError(t, err)
}
