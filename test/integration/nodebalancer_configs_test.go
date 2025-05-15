package integration

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/linode/linodego"
)

var TestNodeBalancerConfigCreateOpts = linodego.NodeBalancerConfigCreateOptions{
	Port:          80,
	Protocol:      linodego.ProtocolHTTP,
	Algorithm:     linodego.AlgorithmRoundRobin,
	CheckInterval: 60,
}

func TestNodeBalancerConfig_Create_smoke(t *testing.T) {
	_, _, config, teardown, err := setupNodeBalancerConfig(t, "fixtures/TestNodeBalancerConfig_Create")
	defer teardown()

	if err != nil {
		t.Errorf("Error creating NodeBalancer Config, got error %v", err)
	}

	expected := TestNodeBalancerConfigCreateOpts

	// cant compare Target, fixture IPs are sanitized
	if config.Port != expected.Port || config.Protocol != expected.Protocol {
		t.Errorf("NodeBalancerConfig did not match CreateOptions")
	}
}

func TestNodeBalancerConfig_Update(t *testing.T) {
	client, nodebalancer, config, teardown, err := setupNodeBalancerConfig(t, "fixtures/TestNodeBalancerConfig_Update")
	defer teardown()
	if err != nil {
		t.Error(err)
	}

	updateOpts := linodego.NodeBalancerConfigUpdateOptions{
		Port:          8080,
		Protocol:      linodego.ProtocolTCP,
		ProxyProtocol: linodego.ProxyProtocolV2,
		Algorithm:     linodego.AlgorithmLeastConn,
	}
	configUpdated, err := client.UpdateNodeBalancerConfig(context.Background(), nodebalancer.ID, config.ID, updateOpts)
	if err != nil {
		t.Errorf("Error updating NodeBalancer Config, %s", err)
	}
	if configUpdated.Port != updateOpts.Port ||
		string(updateOpts.Algorithm) != string(configUpdated.Algorithm) ||
		string(updateOpts.Protocol) != string(configUpdated.Protocol) ||
		string(updateOpts.ProxyProtocol) != string(configUpdated.ProxyProtocol) {
		t.Errorf("NodeBalancerConfig did not match UpdateOptions")
	}
}

func TestNodeBalancerConfigs_List(t *testing.T) {
	client, nodebalancer, _, teardown, err := setupNodeBalancerConfig(t, "fixtures/TestNodeBalancerConfigs_List")
	defer teardown()
	if err != nil {
		t.Error(err)
	}

	listOpts := linodego.NewListOptions(0, "")
	configs, err := client.ListNodeBalancerConfigs(context.Background(), nodebalancer.ID, listOpts)
	if err != nil {
		t.Errorf("Error listing nodebalancers configs, expected array, got error %v", err)
	}
	if len(configs) != listOpts.Results {
		t.Errorf("Expected ListNodeBalancerConfigs to match API result count")
	}
}

func TestNodeBalancerConfigs_ListMultiplePages(t *testing.T) {
	// This fixture was hand-crafted to render an empty page 1 result, with a single result on page 2
	// "results:1,data:[],page:1,pages:2"  .. "results:1,data[{...}],page:2,pages:2"
	client, nodebalancer, _, teardown, err := setupNodeBalancerConfig(t, "fixtures/TestNodeBalancerConfigs_ListMultiplePages")
	defer teardown()
	if err != nil {
		t.Error(err)
	}

	listOpts := linodego.NewListOptions(0, "")
	configs, err := client.ListNodeBalancerConfigs(context.Background(), nodebalancer.ID, listOpts)
	if err != nil {
		t.Errorf("Error listing nodebalancers configs, expected array, got error %v", err)
	}
	if len(configs) != listOpts.Results {
		t.Errorf("Expected ListNodeBalancerConfigs count to match API results count")
	}
}

func TestNodeBalancerConfig_Get(t *testing.T) {
	client, nodebalancer, config, teardown, err := setupNodeBalancerConfig(t, "fixtures/TestNodeBalancerConfig_Get")
	defer teardown()
	if err != nil {
		t.Error(err)
	}

	configGot, err := client.GetNodeBalancerConfig(context.Background(), nodebalancer.ID, config.ID)
	if configGot.Port != config.Port {
		t.Errorf("GetNodeBalancerConfig did not get the expected config")
	}
	if err != nil {
		t.Errorf("Error getting nodebalancer %d, got error %v", nodebalancer.ID, err)
	}
}

func TestNodeBalancerConfig_UDP(t *testing.T) {
	_, _, config, teardown, err := setupNodeBalancerConfig(
		t,
		"fixtures/TestNodeBalancerConfig_UDP",
		func(options *linodego.NodeBalancerConfigCreateOptions) {
			options.Protocol = linodego.ProtocolUDP
			options.UDPCheckPort = linodego.Pointer(1234)
		},
	)
	defer teardown()

	if err != nil {
		t.Errorf("Error creating NodeBalancer Config, got error %v", err)
	}

	require.Equal(t, linodego.ProtocolUDP, config.Protocol)
	require.Equal(t, 1234, config.UDPCheckPort)
	require.NotZero(t, config.UDPSessionTimeout)
}

func TestNodeBalancerConfig_Rebuild_InVPCWithInstance(t *testing.T) {
	client, nodebalancer, subnet, instanceVPCIP, teardown, err := setupNodeBalancerWithVPCAndInstance(t, "fixtures/TestNodeBalancerConfig_Rebuild_InVPCWithInstance")
	defer teardown()
	if err != nil {
		t.Error(err)
	}

	// Create a simple nodebalancer config
	config, err := client.CreateNodeBalancerConfig(context.Background(), nodebalancer.ID, TestNodeBalancerConfigCreateOpts)
	if err != nil {
		t.Fatalf("Error creating NodeBalancer Config, got error %v", err)
	}

	// Rebuild the nodebalancer config with the instance
	rebuildOpts := linodego.NodeBalancerConfigRebuildOptions{
		Port:          80,
		Protocol:      linodego.ProtocolHTTP,
		Algorithm:     linodego.AlgorithmRoundRobin,
		CheckInterval: 60,
		Nodes: []linodego.NodeBalancerConfigRebuildNodeOptions{
			{
				NodeBalancerNodeCreateOptions: linodego.NodeBalancerNodeCreateOptions{
					Address:  fmt.Sprintf("%s:80", instanceVPCIP),
					Mode:     linodego.ModeAccept,
					Weight:   1,
					Label:    "test",
					SubnetID: subnet.ID,
				},
			},
		},
	}

	config, err = client.RebuildNodeBalancerConfig(context.Background(), nodebalancer.ID, config.ID, rebuildOpts)
	if err != nil {
		t.Fatalf("Error creating NodeBalancer Config, got error %v", err)
	}

	// List nodebalancer nodes
	nodes, err := client.ListNodeBalancerNodes(context.Background(), nodebalancer.ID, config.ID, nil)
	if err != nil {
		t.Errorf("Error listing nodebalancer nodes: %s", err)
	}
	if len(nodes) != 1 {
		t.Errorf("Expected exactly one nodebalancer node, got %d", len(nodes))
	}
	if nodes[0].Address == "" {
		t.Errorf("Expected nodebalancer node address to be there, got %s", nodes[0].Address)
	}

	// get nodebalancer vpc config
	vpcConfigs, err := client.ListNodeBalancerVPCConfigs(context.Background(), nodebalancer.ID, nil)
	if err != nil {
		t.Errorf("Error listing nodebalancer VPC configs: %s", err)
	}
	if len(vpcConfigs) != 1 {
		t.Errorf("Expected exactly one nodebalancer VPC config, got %d", len(vpcConfigs))
	}
	if vpcConfigs[0].ID != nodes[0].VPCConfigID {
		t.Errorf("Expected nodebalancer VPC config ID to be the same as the nodebalancer node VPC config ID, got %d", vpcConfigs[0].ID)
	}
}

func createNodeBalancerConfig(
	t *testing.T,
	client *linodego.Client,
	nodeBalancerID int,
	modifiers ...func(options *linodego.NodeBalancerConfigCreateOptions),
) (*linodego.NodeBalancerConfig, func(), error) {
	t.Helper()

	createOpts := TestNodeBalancerConfigCreateOpts

	for _, modifier := range modifiers {
		modifier(&createOpts)
	}

	config, err := client.CreateNodeBalancerConfig(context.Background(), nodeBalancerID, createOpts)
	if err != nil {
		t.Fatalf("Error creating NodeBalancer Config, got error %v", err)
	}

	teardown := func() {
		// delete the NodeBalancerConfig to exercise the code
		if err := client.DeleteNodeBalancerConfig(context.Background(), nodeBalancerID, config.ID); err != nil {
			t.Fatalf("Expected to delete a NodeBalancer Config, but got %v", err)
		}
	}
	return config, teardown, err
}

func setupNodeBalancerConfig(t *testing.T, fixturesYaml string, modifiers ...func(options *linodego.NodeBalancerConfigCreateOptions)) (*linodego.Client, *linodego.NodeBalancer, *linodego.NodeBalancerConfig, func(), error) {
	t.Helper()
	var fixtureTeardown func()
	client, nodebalancer, fixtureTeardown, err := setupNodeBalancer(t, fixturesYaml, nil)
	if err != nil {
		t.Fatalf("Error creating nodebalancer, got error %v", err)
	}

	config, configTeardown, err := createNodeBalancerConfig(t, client, nodebalancer.ID, modifiers...)
	if err != nil {
		t.Fatalf("Error creating NodeBalancer Config, got error %v", err)
	}

	teardown := func() {
		configTeardown()
		fixtureTeardown()
	}
	return client, nodebalancer, config, teardown, err
}

func setupNodeBalancerWithVPCAndInstance(t *testing.T, fixturesYaml string) (*linodego.Client, *linodego.NodeBalancer, *linodego.VPCSubnet, string, func(), error) {
	t.Helper()
	var fixtureTeardown func()
	client, nodebalancer, _, subnet, fixtureTeardown, err := setupNodeBalancerWithVPC(t, fixturesYaml, func(client *linodego.Client, options *linodego.VPCCreateOptions) {
		options.Region = getRegionsWithCaps(t, client, []string{"Linodes", "VPCs"})[1]
	})
	if err != nil {
		t.Fatalf("Error creating nodebalancer, got error %v", err)
	}

	// Create an instance in the VPC subnet
	instance, _, instanceTeardown, err := createInstanceWithoutDisks(
		t,
		client,
		true,
		func(client *linodego.Client, opts *linodego.InstanceCreateOptions) {
			opts.Region = getRegionsWithCaps(t, client, []string{"Linodes", "VPCs"})[1]
			opts.Image = "linode/ubuntu22.04"
			opts.RootPass = "0o37Klm56P4ssw0rd"

			NAT1To1Any := "any"
			opts.Interfaces = []linodego.InstanceConfigInterfaceCreateOptions{
				{
					Purpose:  "vpc",
					SubnetID: &subnet.ID,
					IPv4: &linodego.VPCIPv4{
						NAT1To1: &NAT1To1Any,
					},
				},
			}
		},
	)
	if err != nil {
		if instanceTeardown != nil {
			instanceTeardown()
		}
		t.Fatal("Error creating instance: ", err)
	}

	instanceConfigs, err := client.ListInstanceConfigs(context.Background(), instance.ID, nil)
	if err != nil {
		t.Fatalf("Error listing instance configs: %s", err)
	}

	// Find the VPC interface and get its IP.
	var instanceVPCIP string
	for _, iface := range instanceConfigs[0].Interfaces {
		if iface.Purpose == "vpc" && iface.IPv4 != nil {
			instanceVPCIP = iface.IPv4.VPC
			break
		}
	}
	if instanceVPCIP == "" {
		t.Fatal("Failed to find VPC IP address for instance")
	}

	teardown := func() {
		instanceTeardown()
		fixtureTeardown()
	}

	return client, nodebalancer, subnet, instanceVPCIP, teardown, err
}
