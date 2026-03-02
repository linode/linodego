package unit

import (
	"context"
	"testing"

	"github.com/linode/linodego"
	"github.com/stretchr/testify/assert"
)

func TestNodebalancerVPCConfig_List(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("nodebalancer_vpc_config_list")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("nodebalancers/123/vpcs", fixtureData)

	configs, err := base.Client.ListNodeBalancerVPCConfigs(context.Background(), 123, &linodego.ListOptions{})
	assert.NoError(t, err)

	assert.Len(t, configs, 2)

	assert.Equal(t, 6, configs[0].ID)
	assert.Equal(t, "10.0.0.12/30", configs[0].IPv4Range)
	assert.Equal(t, "", configs[0].IPv6Range)
	assert.Equal(t, 123, configs[0].NodeBalancerID)
	assert.Equal(t, 1, configs[0].SubnetID)
	assert.Equal(t, 1, configs[0].VPCID)
	assert.Equal(t, linodego.NodeBalancerVPCConfigPurposeBackend, configs[0].Purpose)

	assert.Equal(t, 7, configs[1].ID)
	assert.Equal(t, "10.0.1.8/30", configs[1].IPv4Range)
	assert.Equal(t, "", configs[1].IPv6Range)
	assert.Equal(t, 123, configs[1].NodeBalancerID)
	assert.Equal(t, 7, configs[1].SubnetID)
	assert.Equal(t, 1, configs[1].VPCID)
	assert.Equal(t, linodego.NodeBalancerVPCConfigPurposeFrontend, configs[1].Purpose)
}

func TestNodebalancerVPCConfig_BackendList(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("nodebalancer_vpc_config_backend_list")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("nodebalancers/123/backend_vpcs", fixtureData)

	configs, err := base.Client.ListNodeBalancerVPCBackendConfigs(context.Background(), 123, &linodego.ListOptions{})
	assert.NoError(t, err)

	assert.Len(t, configs, 2)

	assert.Equal(t, 9, configs[0].ID)
	assert.Equal(t, "10.0.0.12/30", configs[0].IPv4Range)
	assert.Equal(t, "", configs[0].IPv6Range)
	assert.Equal(t, 123, configs[0].NodeBalancerID)
	assert.Equal(t, 3, configs[0].SubnetID)
	assert.Equal(t, 1, configs[0].VPCID)
	assert.Equal(t, linodego.NodeBalancerVPCConfigPurposeBackend, configs[0].Purpose)

	assert.Equal(t, 10, configs[1].ID)
	assert.Equal(t, "10.0.1.16/30", configs[1].IPv4Range)
	assert.Equal(t, "", configs[1].IPv6Range)
	assert.Equal(t, 123, configs[1].NodeBalancerID)
	assert.Equal(t, 6, configs[1].SubnetID)
	assert.Equal(t, 1, configs[1].VPCID)
	assert.Equal(t, linodego.NodeBalancerVPCConfigPurposeBackend, configs[1].Purpose)
}

func TestNodebalancerVPCConfig_FrontendList(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("nodebalancer_vpc_config_frontend_list")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("nodebalancers/123/frontend_vpcs", fixtureData)

	configs, err := base.Client.ListNodeBalancerVPCFrontendConfigs(context.Background(), 123, &linodego.ListOptions{})
	assert.NoError(t, err)

	assert.Len(t, configs, 1)

	assert.Equal(t, 14, configs[0].ID)
	assert.Equal(t, "10.0.0.36/30", configs[0].IPv4Range)
	assert.Equal(t, "", configs[0].IPv6Range)
	assert.Equal(t, 123, configs[0].NodeBalancerID)
	assert.Equal(t, 3, configs[0].SubnetID)
	assert.Equal(t, 1, configs[0].VPCID)
	assert.Equal(t, linodego.NodeBalancerVPCConfigPurposeFrontend, configs[0].Purpose)
}

func TestNodebalancerVPCConfig_Get(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("nodebalancer_vpc_config_get")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("nodebalancers/123/vpcs/12345", fixtureData)

	config, err := base.Client.GetNodeBalancerVPCConfig(context.Background(), 123, 12345)
	assert.NoError(t, err)

	assert.Equal(t, 6, config.ID)
	assert.Equal(t, "10.0.0.12/30", config.IPv4Range)
	assert.Equal(t, "", config.IPv6Range)
	assert.Equal(t, 123, config.NodeBalancerID)
	assert.Equal(t, 1, config.SubnetID)
	assert.Equal(t, 1, config.VPCID)
	assert.Equal(t, linodego.NodeBalancerVPCConfigPurposeBackend, config.Purpose)
}
