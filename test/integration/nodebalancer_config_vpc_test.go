package integration

import (
	"context"
	"testing"

	"github.com/linode/linodego"
	"github.com/stretchr/testify/require"
)

func TestNodeBalancerVpcConfig_List(t *testing.T) {
	client, nodebalancer, teardown, err := setupNodeBalancerVpcConfig(t, "fixtures/TestNodeBalancerVpcConfig_List")
	if err != nil {
		t.Errorf("Error setting up nodebalancer: %s", err)
	}
	defer teardown()

	configs, err := client.ListNodeBalancerVpcConfigs(context.Background(), nodebalancer.ID, nil)
	if err != nil {
		t.Errorf("Error listing nodebalancer VPC configs: %s", err)
	}

	// We expect the list to be not empty and have at least one VPC config
	require.NotEmpty(t, configs)
	require.Len(t, configs, 1)
}

func TestNodeBalancerVpcConfig_Get(t *testing.T) {
	client, nodebalancer, teardown, err := setupNodeBalancerVpcConfig(t, "fixtures/TestNodeBalancerVpcConfig_Get")
	if err != nil {
		t.Errorf("Error setting up nodebalancer: %s", err)
	}
	defer teardown()

	// Get the VPC config list for the nodebalancer (should only have one)
	configs, err := client.ListNodeBalancerVpcConfigs(context.Background(), nodebalancer.ID, nil)
	if err != nil {
		t.Errorf("Error listing nodebalancer VPC configs: %s", err)
	}
	require.NotEmpty(t, configs)
	require.Len(t, configs, 1)

	// Get the VPC config by ID
	config, err := client.GetNodeBalancerVpcConfig(context.Background(), nodebalancer.ID, configs[0].ID)
	if err != nil {
		t.Errorf("Error getting nodebalancer VPC config: %s", err)
	}
	require.NotNil(t, config)
	require.Equal(t, configs[0].ID, config.ID)
}

func setupNodeBalancerVpcConfig(t *testing.T, fixturesYaml string) (*linodego.Client, *linodego.NodeBalancer, func(), error) {
	t.Helper()
	client, fixtureTeardown := createTestClient(t, fixturesYaml)

	// 1. Create a VPC and Subnet.
	_, vpcSubnet, vpcTeardown, err := createVPCWithSubnet(t, client)
	if err != nil {
		t.Errorf("Error creating VPC and subnet: %s", err)
	}

	// 2. Create a NodeBalancer
	// We need a region that supports both NodeBalancers and VPC
	nbCreateOpts := linodego.NodeBalancerCreateOptions{
		Label:  &label,
		Region: getRegionsWithCaps(t, client, []string{"NodeBalancers", "VPCs"})[0],
		Vpcs: []*linodego.VPCConfig{ // We need to add this functionality to linodego
			{
				IPv4Range: vpcSubnet.IPv4,
				SubnetID:  vpcSubnet.ID,
			},
		},
	}
	nodebalancer, err := client.CreateNodeBalancer(context.Background(), nbCreateOpts)
	if err != nil {
		t.Errorf("Error creating nodebalancer: %s", err)
	}

	teardown := func() {
		// Delete resources in reverse order of creation.
		if nodebalancer != nil {
			err := client.DeleteNodeBalancer(context.Background(), nodebalancer.ID)
			if err != nil {
				t.Errorf("Error deleting nodebalancer: %s", err)
			}
		}
		// Only need to call vpcTeardown, which deletes the VPC and subnet
		vpcTeardown()

		fixtureTeardown()
	}
	return client, nodebalancer, teardown, err
}
