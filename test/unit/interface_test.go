package unit

import (
	"context"
	"testing"

	"github.com/linode/linodego"
	"github.com/stretchr/testify/assert"
)

func TestInterface_Get(t *testing.T) {
	fixtures := NewTestFixtures()

	fixtureData, err := fixtures.GetFixture("interface_get")
	if err != nil {
		t.Fatalf("Failed to load fixture: %v", err)
	}

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("linode/instances/123/interfaces/123", fixtureData)

	iface, err := base.Client.GetInterface(context.Background(), 123, 123)
	if err != nil {
		t.Fatalf("Error fetching interfaces: %v", err)
	}

	assert.Equal(t, 123, iface.ID)
	assert.Equal(t, 1, iface.Version)
	assert.Equal(t, false, *iface.DefaultRoute.IPv4)
	assert.Equal(t, "my_vlan", iface.VLAN.Label)
}

func TestInterface_List(t *testing.T) {
	fixtures := NewTestFixtures()

	fixtureData, err := fixtures.GetFixture("interface_list")
	if err != nil {
		t.Fatalf("Failed to load fixture: %v", err)
	}

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("linode/instances/123/interfaces", fixtureData)

	ifaces, err := base.Client.ListInterfaces(context.Background(), 123, nil)
	if err != nil {
		t.Fatalf("Error fetching interfaces: %v", err)
	}

	assert.Equal(t, 123, ifaces[0].ID)
	assert.Equal(t, 1, ifaces[0].Version)
}

func TestInterface_Delete(t *testing.T) {
	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockDelete("linode/instances/123/interfaces/123", nil)

	err := base.Client.DeleteInterface(context.Background(), 123, 123)
	assert.NoError(t, err)
}

func TestInterface_Create(t *testing.T) {
	fixtures := NewTestFixtures()

	fixtureData, err := fixtures.GetFixture("interface_create")
	if err != nil {
		t.Fatalf("Failed to load fixture: %v", err)
	}

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockPost("linode/instances/123/interfaces", fixtureData)

	opts := linodego.LinodeInterfaceCreateOptions{
		FirewallID: linodego.Pointer(123),
		Public:     nil,
	}

	iface, err := base.Client.CreateInterface(context.Background(), 123, opts)
	if err != nil {
		t.Fatalf("Error fetching interfaces: %v", err)
	}

	assert.Equal(t, 123, iface.ID)
	assert.Equal(t, "auto", iface.Public.IPv4.Addresses[0].Address)
}

func TestInterface_Update(t *testing.T) {
	fixtures := NewTestFixtures()

	fixtureData, err := fixtures.GetFixture("interface_update")
	if err != nil {
		t.Fatalf("Failed to load fixture: %v", err)
	}

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockPut("linode/instances/123/interfaces/123", fixtureData)

	opts := linodego.LinodeInterfaceUpdateOptions{
		DefaultRoute: &linodego.InterfaceDefaultRoute{
			IPv6: linodego.Pointer(true),
		},
	}

	iface, err := base.Client.UpdateInterface(context.Background(), 123, 123, opts)
	if err != nil {
		t.Fatalf("Error fetching interfaces: %v", err)
	}

	assert.Equal(t, 123, iface.ID)
	assert.Equal(t, false, *iface.DefaultRoute.IPv4)
	assert.Equal(t, true, *iface.DefaultRoute.IPv6)
}

func TestInterface_Upgrade(t *testing.T) {
	fixtures := NewTestFixtures()

	fixtureData, err := fixtures.GetFixture("interface_upgrade")
	if err != nil {
		t.Fatalf("Failed to load fixture: %v", err)
	}

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockPost("linode/instances/123/upgrade-interfaces", fixtureData)

	opts := linodego.LinodeInterfacesUpgradeOptions{
		ConfigID: linodego.Pointer(123),
		DryRun:   linodego.Pointer(false),
	}

	iface, err := base.Client.UpgradeInterfaces(context.Background(), 123, opts)
	if err != nil {
		t.Fatalf("Error fetching interfaces: %v", err)
	}

	assert.Equal(t, 123, iface.ConfigID)
	assert.Equal(t, false, iface.DryRun)
	assert.Equal(t, 123, iface.Interfaces[0].ID)
}

func TestInteface_ListFirewalls(t *testing.T) {
	fixtures := NewTestFixtures()

	fixtureData, err := fixtures.GetFixture("firewall_list")
	if err != nil {
		t.Fatalf("Failed to load fixture: %v", err)
	}

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("linode/instances/123/interfaces/123/firewalls", fixtureData)

	firewalls, err := base.Client.ListInterfaceFirewalls(context.Background(), 123, 123, nil)
	if err != nil {
		t.Fatalf("Error fetching firewalls: %v", err)
	}

	assert.Equal(t, 123, firewalls[0].ID)
}

func TestInteface_GetSettings(t *testing.T) {
	fixtures := NewTestFixtures()

	fixtureData, err := fixtures.GetFixture("interface_settings_get")
	if err != nil {
		t.Fatalf("Failed to load fixture: %v", err)
	}

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("linode/instances/123/interfaces/settings", fixtureData)

	settings, err := base.Client.GetInterfaceSettings(context.Background(), 123)
	if err != nil {
		t.Fatalf("Error fetching firewalls: %v", err)
	}

	assert.Equal(t, false, settings.NetworkHelper)
}

func TestInterface_UpdateSettings(t *testing.T) {
	fixtures := NewTestFixtures()

	fixtureData, err := fixtures.GetFixture("interface_settings_update")
	if err != nil {
		t.Fatalf("Failed to load fixture: %v", err)
	}

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockPut("linode/instances/123/interfaces/settings", fixtureData)

	opts := linodego.InterfaceSettingsUpdateOptions{
		NetworkHelper: linodego.Pointer(true),
	}

	settings, err := base.Client.UpdateInterfaceSettings(context.Background(), 123, opts)
	if err != nil {
		t.Fatalf("Error fetching interfaces: %v", err)
	}

	assert.Equal(t, true, settings.NetworkHelper)
}
