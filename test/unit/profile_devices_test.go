package unit

import (
	"context"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestProfileDevices_Get(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("profile_devices_get")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("profile/devices/123", fixtureData)

	device, err := base.Client.GetProfileDevice(context.Background(), 123)
	assert.NoError(t, err)

	assert.Equal(t, 123, device.ID)
	assert.Equal(t, "203.0.113.1", device.LastRemoteAddr)
	assert.Equal(t, "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3538.77 Safari/537.36 Vivaldi/2.1.1337.36", device.UserAgent)
}

func TestProfileDevices_List(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("profile_devices_list")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("profile/devices", fixtureData)

	devices, err := base.Client.ListProfileDevices(context.Background(), nil)
	assert.NoError(t, err)

	assert.Equal(t, 1, len(devices))
	device := devices[0]

	assert.Equal(t, 123, device.ID)
	assert.Equal(t, "203.0.113.1", device.LastRemoteAddr)
	assert.Equal(t, "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3538.77 Safari/537.36 Vivaldi/2.1.1337.36", device.UserAgent)
}

func TestProfileDevices_Delete(t *testing.T) {
	client := createMockClient(t)

	httpmock.RegisterRegexpResponder("DELETE", mockRequestURL(t, "profile/devices/123"), httpmock.NewStringResponder(200, "{}"))

	if err := client.DeleteProfileDevice(context.Background(), 123); err != nil {
		t.Fatal(err)
	}
}
