package integration

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNodeBalancerVPCConfig_List(t *testing.T) {
	client, nodebalancer, _, _, teardown, err := setupNodeBalancerWithVPC(t, "fixtures/TestNodeBalancerVpcConfig_List")
	if err != nil {
		t.Errorf("Error setting up nodebalancer: %s", err)
	}
	defer teardown()

	configs, err := client.ListNodeBalancerVPCConfigs(context.Background(), nodebalancer.ID, nil)
	if err != nil {
		t.Errorf("Error listing nodebalancer VPC configs: %s", err)
	}

	// We expect the list to be not empty and have at least one VPC config
	require.NotEmpty(t, configs)
	require.Len(t, configs, 1)
}

func TestNodeBalancerVPCConfig_Get(t *testing.T) {
	client, nodebalancer, _, _, teardown, err := setupNodeBalancerWithVPC(t, "fixtures/TestNodeBalancerVpcConfig_Get")
	if err != nil {
		t.Errorf("Error setting up nodebalancer: %s", err)
	}
	defer teardown()

	// Get the VPC config list for the nodebalancer (should only have one)
	configs, err := client.ListNodeBalancerVPCConfigs(context.Background(), nodebalancer.ID, nil)
	if err != nil {
		t.Errorf("Error listing nodebalancer VPC configs: %s", err)
	}
	require.NotEmpty(t, configs)
	require.Len(t, configs, 1)

	// Get the VPC config by ID
	config, err := client.GetNodeBalancerVPCConfig(context.Background(), nodebalancer.ID, configs[0].ID)
	if err != nil {
		t.Errorf("Error getting nodebalancer VPC config: %s", err)
	}
	require.NotNil(t, config)
	require.Equal(t, configs[0].ID, config.ID)
}
