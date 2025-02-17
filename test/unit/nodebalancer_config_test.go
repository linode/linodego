package unit

import (
	"context"
	"testing"

	"github.com/linode/linodego"
	"github.com/stretchr/testify/assert"
)

func TestNodeBalancerConfig_List(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("nodebalancer_config_list")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("nodebalancers/123/configs", fixtureData)

	configs, err := base.Client.ListNodeBalancerConfigs(context.Background(), 123, nil)
	assert.NoError(t, err)
	assert.Len(t, configs, 2)

	assert.Equal(t, 456, configs[0].ID)
	assert.Equal(t, 80, configs[0].Port)
	assert.Equal(t, linodego.ProtocolHTTP, configs[0].Protocol)
}

func TestNodeBalancerConfig_Get(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("nodebalancer_config_get")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("nodebalancers/123/configs/456", fixtureData)

	config, err := base.Client.GetNodeBalancerConfig(context.Background(), 123, 456)
	assert.NoError(t, err)
	assert.NotNil(t, config)

	assert.Equal(t, 456, config.ID)
	assert.Equal(t, 80, config.Port)
	assert.Equal(t, linodego.ProtocolHTTP, config.Protocol)
}

func TestNodeBalancerConfig_Create(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("nodebalancer_config_create")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockPost("nodebalancers/123/configs", fixtureData)

	createOpts := linodego.NodeBalancerConfigCreateOptions{
		Port:       80,
		Protocol:   linodego.ProtocolHTTP,
		Algorithm:  linodego.AlgorithmRoundRobin,
		Stickiness: linodego.StickinessTable,
	}
	config, err := base.Client.CreateNodeBalancerConfig(context.Background(), 123, createOpts)
	assert.NoError(t, err)

	assert.Equal(t, 456, config.ID)
	assert.Equal(t, 80, config.Port)
	assert.Equal(t, linodego.ProtocolHTTP, config.Protocol)
}

func TestNodeBalancerConfig_Update(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("nodebalancer_config_update")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockPut("nodebalancers/123/configs/456", fixtureData)

	updateOpts := linodego.NodeBalancerConfigUpdateOptions{
		Port:       443,
		Protocol:   linodego.ProtocolHTTPS,
		Stickiness: linodego.StickinessNone,
	}
	config, err := base.Client.UpdateNodeBalancerConfig(context.Background(), 123, 456, updateOpts)
	assert.NoError(t, err)

	assert.Equal(t, 456, config.ID)
	assert.Equal(t, 443, config.Port)
	assert.Equal(t, linodego.ProtocolHTTPS, config.Protocol)
}

func TestNodeBalancerConfig_Delete(t *testing.T) {
	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockDelete("nodebalancers/123/configs/456", nil)

	err := base.Client.DeleteNodeBalancerConfig(context.Background(), 123, 456)
	assert.NoError(t, err)
}

func TestNodeBalancerConfig_Rebuild(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("nodebalancer_config_rebuild")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockPost("nodebalancers/123/configs/456/rebuild", fixtureData)

	rebuildOpts := linodego.NodeBalancerConfigRebuildOptions{
		Port:       443,
		Protocol:   linodego.ProtocolHTTPS,
		Stickiness: linodego.StickinessNone,
	}
	config, err := base.Client.RebuildNodeBalancerConfig(context.Background(), 123, 456, rebuildOpts)
	assert.NoError(t, err)

	assert.Equal(t, 456, config.ID)
	assert.Equal(t, 443, config.Port)
	assert.Equal(t, linodego.ProtocolHTTPS, config.Protocol)
}
