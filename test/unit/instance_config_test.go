package unit

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/linode/linodego"
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
	assert.Equal(t, 1, config.ID)
	assert.Equal(t, "config-1", config.Label)
}

func TestInstanceConfig_Create(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("instance_config_create")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	rootDevice := "/dev/sda"

	createOptions := linodego.InstanceConfigCreateOptions{
		Label:      linodego.Pointer("new-config"),
		Kernel:     linodego.Pointer("linode/latest-64bit"),
		RootDevice: &rootDevice,
	}

	base.MockPost("linode/instances/123/configs", fixtureData)

	config, err := base.Client.CreateInstanceConfig(context.Background(), 123, createOptions)
	assert.NoError(t, err)
	assert.Equal(t, "new-config", config.Label)
}

func TestInstanceConfig_Update(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("instance_config_update")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	updateOptions := linodego.InstanceConfigUpdateOptions{
		Label:      linodego.Pointer("updated-config"),
		RootDevice: linodego.Pointer("/dev/sdb"),
	}

	base.MockPut("linode/instances/123/configs/1", fixtureData)

	config, err := base.Client.UpdateInstanceConfig(context.Background(), 123, 1, updateOptions)
	assert.NoError(t, err)
	assert.Equal(t, "updated-config", config.Label)
}

func TestInstanceConfig_Delete(t *testing.T) {
	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockDelete("linode/instances/123/configs/1", nil)

	err := base.Client.DeleteInstanceConfig(context.Background(), 123, 1)
	assert.NoError(t, err)
}
