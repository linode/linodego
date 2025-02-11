package unit

import (
	"context"
	"fmt"
	"github.com/jarcoal/httpmock"
	"github.com/linode/linodego"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFirewallDevice_List(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("firewall_device_list")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	firewallID := 123

	base.MockGet(formatMockAPIPath("networking/firewalls/%d/devices", firewallID), fixtureData)

	firewallDevices, err := base.Client.ListFirewallDevices(context.Background(), firewallID, &linodego.ListOptions{})
	assert.NoError(t, err)

	assert.Equal(t, 2, len(firewallDevices))

	for _, device := range firewallDevices {
		assert.NotNil(t, device.Entity)
		assert.NotEqual(t, 0, device.ID)

		assert.NotNil(t, device.Created)
		assert.NotNil(t, device.Updated)

		switch device.Entity.Type {
		case "linode":
			assert.Equal(t, 123, device.Entity.ID)
			assert.Equal(t, "my-linode", device.Entity.Label)
			assert.Equal(t, "/v4/linode/instances/123", device.Entity.URL)
		case "nodebalancer":
			assert.Equal(t, 321, device.Entity.ID)
			assert.Equal(t, "my-nodebalancer", device.Entity.Label)
			assert.Equal(t, "/v4/nodebalancers/123", device.Entity.URL)
		default:
			t.Fatalf("Unexpected device type: %s", device.Entity.Type)
		}
	}
}

func TestFirewallDevice_Get(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("firewall_device_get")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	firewallID := 123

	deviceID := 123

	base.MockGet(formatMockAPIPath("networking/firewalls/%d/devices/%d", firewallID, deviceID), fixtureData)

	firewallDevice, err := base.Client.GetFirewallDevice(context.Background(), firewallID, deviceID)
	assert.NoError(t, err)
	assert.NotNil(t, firewallDevice)

	assert.Equal(t, deviceID, firewallDevice.ID)
	assert.NotNil(t, firewallDevice.Entity)

	assert.Equal(t, 123, firewallDevice.Entity.ID)
	assert.Equal(t, "my-linode", firewallDevice.Entity.Label)
	assert.Equal(t, linodego.FirewallDeviceType("linode"), firewallDevice.Entity.Type)
	assert.Equal(t, "/v4/linode/instances/123", firewallDevice.Entity.URL)

	assert.NotNil(t, firewallDevice.Created)
	assert.NotNil(t, firewallDevice.Updated)
}

func TestFirewallDevice_Create(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("firewall_device_create")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	firewallID := 123

	requestData := linodego.FirewallDeviceCreateOptions{
		ID:   123,
		Type: "linode",
	}

	base.MockPost(formatMockAPIPath("networking/firewalls/%d/devices", firewallID), fixtureData)

	firewallDevice, err := base.Client.CreateFirewallDevice(context.Background(), firewallID, requestData)
	assert.NoError(t, err)
	assert.NotNil(t, firewallDevice)

	assert.NotNil(t, firewallDevice.Entity)

	assert.Equal(t, 123, firewallDevice.Entity.ID)
	assert.Equal(t, "my-linode", firewallDevice.Entity.Label)
	assert.Equal(t, linodego.FirewallDeviceType("linode"), firewallDevice.Entity.Type)
	assert.Equal(t, "/v4/linode/instances/123", firewallDevice.Entity.URL)

	assert.NotNil(t, firewallDevice.Created)
	assert.NotNil(t, firewallDevice.Updated)
}

func TestFirewallDevice_Delete(t *testing.T) {
	client := createMockClient(t)

	firewallID := 123

	deviceID := 123

	httpmock.RegisterRegexpResponder("DELETE", mockRequestURL(t, fmt.Sprintf("networking/firewalls/%d/devices/%d", firewallID, deviceID)),
		httpmock.NewStringResponder(200, "{}"))

	if err := client.DeleteFirewallDevice(context.Background(), firewallID, deviceID); err != nil {
		t.Fatal(err)
	}
}
