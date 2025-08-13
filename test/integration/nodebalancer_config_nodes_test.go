package integration

import (
	"context"
	"testing"

	"github.com/linode/linodego"
)

var (
	testNodePort                   = "8080"
	testNodeLabel                  = "go-node-test-def"
	testNodeWeight                 = 10
	testNodeBalancerNodeCreateOpts = linodego.NodeBalancerNodeCreateOptions{
		Label:  testNodeLabel,
		Weight: testNodeWeight,
		Mode:   linodego.ModeAccept,
	}
)

func TestNodeBalancerNode_Create_smoke(t *testing.T) {
	_, _, _, node, teardown, err := setupNodeBalancerNode(t, "fixtures/TestNodeBalancerNode_Create")
	defer teardown()

	if err != nil {
		t.Errorf("Error creating NodeBalancer Node, got error %v", err)
	}

	expected := testNodeBalancerNodeCreateOpts

	if node.Label != expected.Label ||
		node.Weight != expected.Weight ||
		node.Mode != expected.Mode {
		t.Errorf("NodeBalancerNode did not match CreateOptions - %v", node)
	}
}

func TestNodeBalancerNode_Update(t *testing.T) {
	client, nodebalancer, config, node, teardown, err := setupNodeBalancerNode(t, "fixtures/TestNodeBalancerNode_Update")
	defer teardown()
	if err != nil {
		t.Error(err)
	}

	updateOpts := linodego.NodeBalancerNodeUpdateOptions{
		Mode:   linodego.ModeDrain,
		Weight: testNodeWeight + 90,
		Label:  testNodeLabel + "_r",
	}
	nodeUpdated, err := client.UpdateNodeBalancerNode(context.Background(), nodebalancer.ID, config.ID, node.ID, updateOpts)
	if err != nil {
		t.Errorf("Error updating NodeBalancer Node, %s", err)
	}

	// fixture sanitization breaks predictability for this test, verify the prefix
	if string(updateOpts.Mode) != string(nodeUpdated.Mode) ||
		updateOpts.Label != nodeUpdated.Label ||
		updateOpts.Weight != nodeUpdated.Weight {
		t.Errorf("NodeBalancerNode did not match UpdateOptions")
	}
}

func TestNodeBalancerNodes_List(t *testing.T) {
	client, nodebalancer, config, _, teardown, err := setupNodeBalancerNode(t, "fixtures/TestNodeBalancerNodes_List")
	defer teardown()
	if err != nil {
		t.Error(err)
	}

	listOpts := linodego.NewListOptions(0, "")
	nodes, err := client.ListNodeBalancerNodes(context.Background(), nodebalancer.ID, config.ID, listOpts)
	if err != nil {
		t.Errorf("Error listing nodebalancers nodes, expected array, got error %v", err)
	}
	if len(nodes) != listOpts.Results {
		t.Errorf("Expected ListNodeBalancerNodes to match API result count")
	}
}

func TestNodeBalancerNodes_ListMultiplePages(t *testing.T) {
	// TODO: replace hand crafted fixtures with unit tests
	// This fixture was hand-crafted to render an empty page 1 result, with a single result on page 2
	// "results:1,data:[],page:1,pages:2"  .. "results:1,data[{...}],page:2,pages:2"
	client, nodebalancer, config, _, teardown, err := setupNodeBalancerNode(t, "fixtures/TestNodeBalancerNodes_ListMultiplePages")
	defer teardown()
	if err != nil {
		t.Error(err)
	}

	listOpts := linodego.NewListOptions(0, "")
	nodes, err := client.ListNodeBalancerNodes(context.Background(), nodebalancer.ID, config.ID, listOpts)
	if err != nil {
		t.Errorf("Error listing nodebalancers configs, expected array, got error %v", err)
	}
	if len(nodes) != listOpts.Results {
		t.Errorf("Expected ListNodeBalancerNodes count to match API results count")
	}
}

func TestNodeBalancerNode_Get(t *testing.T) {
	client, nodebalancer, config, node, teardown, err := setupNodeBalancerNode(t, "fixtures/TestNodeBalancerNode_Get")
	defer teardown()
	if err != nil {
		t.Error(err)
	}

	nodeGot, err := client.GetNodeBalancerNode(context.Background(), nodebalancer.ID, config.ID, node.ID)
	if nodeGot.Address != node.Address {
		t.Errorf("GetNodeBalancerNode did not get the expected node")
	}
	if err != nil {
		t.Errorf("Error getting nodebalancer %d, got error %v", nodebalancer.ID, err)
	}
}

func TestNodeBalancer_Rebuild(t *testing.T) {
	client, nodebalancer, config, node, teardown, err := setupNodeBalancerNode(t, "fixtures/TestNodeBalancer_Rebuild")
	defer teardown()
	if err != nil {
		t.Error(err)
	}

	nbcRebuildOpts := config.GetRebuildOptions()
	nbcRebuildOpts.Nodes = append(
		nbcRebuildOpts.Nodes,
		linodego.NodeBalancerConfigRebuildNodeOptions{
			NodeBalancerNodeCreateOptions: node.GetCreateOptions(),
			ID:                            node.ID,
		},
	)

	nbcGot, err := client.RebuildNodeBalancerConfig(
		context.Background(),
		nodebalancer.ID,
		config.ID,
		nbcRebuildOpts,
	)
	if err != nil {
		t.Errorf("Error rebuilding nodebalancer config %d: %v", config.ID, err)
	}
	if nbcGot.Port != config.Port {
		t.Errorf("RebuildNodeBalancerConfig did not return the expected port")
	}

	newNodes, err := client.ListNodeBalancerNodes(
		context.Background(),
		nodebalancer.ID,
		nbcGot.ID,
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}

	if newNodes[0].ID != node.ID {
		t.Fatalf("expected node ID to match; %d != %d", newNodes[0].ID, node.ID)
	}
}

func TestNodeBalancerNode_Create_InVPC(t *testing.T) {
	client, nodebalancer, subnet, instanceVPCIP, teardown, err := setupNodeBalancerWithVPCAndInstance(t, "fixtures/TestNodeBalancerNode_Create_InVPC")
	defer teardown()
	if err != nil {
		t.Error(err)
	}

	config, err := client.CreateNodeBalancerConfig(context.Background(), nodebalancer.ID, TestNodeBalancerConfigCreateOpts)
	if err != nil {
		t.Errorf("Error creating NodeBalancer Config, got error %v", err)
	}

	// Create a nodebalancer node in the VPC
	node, err := client.CreateNodeBalancerNode(context.Background(), nodebalancer.ID, config.ID, linodego.NodeBalancerNodeCreateOptions{
		Address:  instanceVPCIP + ":" + testNodePort,
		Mode:     linodego.ModeAccept,
		Weight:   10,
		Label:    "go-node-test-def",
		SubnetID: subnet.ID,
	})
	if err != nil {
		t.Fatalf("Error creating NodeBalancer Node, got error %v", err)
	}

	// get nodebalancer vpc config - cross check the nodebalancer node VPC config ID
	vpcConfigs, err := client.ListNodeBalancerVPCConfigs(context.Background(), nodebalancer.ID, nil)
	if err != nil {
		t.Errorf("Error listing nodebalancer VPC configs: %s", err)
	}
	if len(vpcConfigs) != 1 {
		t.Errorf("Expected exactly one nodebalancer VPC config, got %d", len(vpcConfigs))
	}
	if vpcConfigs[0].ID != node.VPCConfigID {
		t.Errorf("Expected nodebalancer VPC config ID to be the same as the nodebalancer node VPC config ID, got %d", vpcConfigs[0].ID)
	}
}

func TestNodeBalancerNode_List_InVPC(t *testing.T) {
	client, nodebalancer, subnet, instanceVPCIP, teardown, err := setupNodeBalancerWithVPCAndInstance(t, "fixtures/TestNodeBalancerNode_List_InVPC")
	defer teardown()
	if err != nil {
		t.Error(err)
	}

	config, err := client.CreateNodeBalancerConfig(context.Background(), nodebalancer.ID, TestNodeBalancerConfigCreateOpts)
	if err != nil {
		t.Errorf("Error creating NodeBalancer Config, got error %v", err)
	}

	node, err := client.CreateNodeBalancerNode(context.Background(), nodebalancer.ID, config.ID, linodego.NodeBalancerNodeCreateOptions{
		Address:  instanceVPCIP + ":" + testNodePort,
		Mode:     linodego.ModeAccept,
		Weight:   10,
		Label:    "go-node-test-def",
		SubnetID: subnet.ID,
	})
	if err != nil {
		t.Errorf("Error creating NodeBalancer Node, got error %v", err)
	}

	// Test listing nodebalancer nodes method
	nodes, err := client.ListNodeBalancerNodes(context.Background(), nodebalancer.ID, config.ID, nil)
	if err != nil {
		t.Fatalf("Error listing nodebalancer nodes: %s", err)
	}
	if len(nodes) != 1 {
		t.Errorf("Expected exactly one nodebalancer node, got %d", len(nodes))
	}
	if nodes[0].Address != instanceVPCIP+":"+testNodePort {
		t.Errorf("Expected nodebalancer node address to be the same as the instance VPC IP, got %s", nodes[0].Address)
	}
	if nodes[0].ID != node.ID {
		t.Errorf("Expected nodebalancer node ID to be the same as the nodebalancer node ID, got %d", nodes[0].ID)
	}
	if nodes[0].VPCConfigID != node.VPCConfigID {
		t.Errorf("Expected nodebalancer node VPC config ID to be the same as the nodebalancer node VPC config ID, got %d", nodes[0].VPCConfigID)
	}

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

func TestNodeBalancerNode_Update_InVPC(t *testing.T) {
	client, nodebalancer, subnet, instanceVPCIP, teardown, err := setupNodeBalancerWithVPCAndInstance(t, "fixtures/TestNodeBalancerNode_Update_InVPC")
	defer teardown()
	if err != nil {
		t.Error(err)
	}

	config, err := client.CreateNodeBalancerConfig(context.Background(), nodebalancer.ID, TestNodeBalancerConfigCreateOpts)
	if err != nil {
		t.Errorf("Error creating NodeBalancer Config, got error %v", err)
	}

	node, err := client.CreateNodeBalancerNode(context.Background(), nodebalancer.ID, config.ID, linodego.NodeBalancerNodeCreateOptions{
		Address:  instanceVPCIP + ":" + testNodePort,
		Mode:     linodego.ModeAccept,
		Weight:   10,
		Label:    "not-updated",
		SubnetID: subnet.ID,
	})
	if err != nil {
		t.Errorf("Error creating NodeBalancer Node, got error %v", err)
	}

	updateOpts := linodego.NodeBalancerNodeUpdateOptions{
		Address:  instanceVPCIP + ":" + testNodePort,
		Label:    "updated",
		SubnetID: subnet.ID,
	}

	node, err = client.UpdateNodeBalancerNode(context.Background(), nodebalancer.ID, config.ID, node.ID, updateOpts)
	if err != nil {
		t.Fatalf("Error updating NodeBalancer Node, got error %v", err)
	}
	if node.Label != "updated" {
		t.Errorf("Expected nodebalancer node label to be updated, got %s", node.Label)
	}

	vpcConfigs, err := client.ListNodeBalancerVPCConfigs(context.Background(), nodebalancer.ID, nil)
	if err != nil {
		t.Errorf("Error listing nodebalancer VPC configs: %s", err)
	}
	if len(vpcConfigs) != 1 {
		t.Errorf("Expected exactly one nodebalancer VPC config, got %d", len(vpcConfigs))
	}
	if vpcConfigs[0].ID != node.VPCConfigID {
		t.Errorf("Expected nodebalancer VPC config ID to be the same as the nodebalancer node VPC config ID, got %d", vpcConfigs[0].ID)
	}
}

func TestNodeBalancerNode_Get_InVPC(t *testing.T) {
	client, nodebalancer, subnet, instanceVPCIP, teardown, err := setupNodeBalancerWithVPCAndInstance(t, "fixtures/TestNodeBalancerNode_Get_InVPC")
	defer teardown()
	if err != nil {
		t.Error(err)
	}

	config, err := client.CreateNodeBalancerConfig(context.Background(), nodebalancer.ID, TestNodeBalancerConfigCreateOpts)
	if err != nil {
		t.Errorf("Error creating NodeBalancer Config, got error %v", err)
	}

	node, err := client.CreateNodeBalancerNode(context.Background(), nodebalancer.ID, config.ID, linodego.NodeBalancerNodeCreateOptions{
		Address:  instanceVPCIP + ":" + testNodePort,
		Mode:     linodego.ModeAccept,
		Weight:   10,
		Label:    "go-node-test-def",
		SubnetID: subnet.ID,
	})
	if err != nil {
		t.Errorf("Error creating NodeBalancer Node, got error %v", err)
	}

	nodeGot, err := client.GetNodeBalancerNode(context.Background(), nodebalancer.ID, config.ID, node.ID)
	if err != nil {
		t.Fatalf("Error getting NodeBalancer Node, got error %v", err)
	}
	if nodeGot.ID != node.ID {
		t.Errorf("Expected nodebalancer node ID to be the same as the nodebalancer node ID, got %d", nodeGot.ID)
	}
	if nodeGot.Address != node.Address {
		t.Errorf("Expected nodebalancer node address to be the same as the nodebalancer node address, got %s", nodeGot.Address)
	}
	if nodeGot.VPCConfigID != node.VPCConfigID {
		t.Errorf("Expected nodebalancer node VPC config ID to be the same as the nodebalancer node VPC config ID, got %d", nodeGot.VPCConfigID)
	}

	vpcConfigs, err := client.ListNodeBalancerVPCConfigs(context.Background(), nodebalancer.ID, nil)
	if err != nil {
		t.Errorf("Error listing nodebalancer VPC configs: %s", err)
	}
	if len(vpcConfigs) != 1 {
		t.Errorf("Expected exactly one nodebalancer VPC config, got %d", len(vpcConfigs))
	}
	if vpcConfigs[0].ID != nodeGot.VPCConfigID {
		t.Errorf("Expected nodebalancer VPC config ID to be the same as the nodebalancer node VPC config ID, got %d", vpcConfigs[0].ID)
	}
}

func setupNodeBalancerNode(
	t *testing.T,
	fixturesYaml string,
) (*linodego.Client, *linodego.NodeBalancer, *linodego.NodeBalancerConfig, *linodego.NodeBalancerNode, func(), error) {
	t.Helper()
	var fixtureTeardown func()
	client, nodebalancer, config, fixtureTeardown, err := setupNodeBalancerConfig(t, fixturesYaml)
	if err != nil {
		t.Fatalf("Error creating nodebalancer config, got error %v", err)
	}

	instance, err := createInstance(t, client, true)
	if err != nil {
		t.Errorf("failed to create test instance: %s", err)
	}

	instanceIP, err := client.AddInstanceIPAddress(context.Background(), instance.ID, false)
	if err != nil {
		t.Fatal(err)
	}

	createOpts := testNodeBalancerNodeCreateOpts
	createOpts.Address = instanceIP.Address + ":" + testNodePort
	node, err := client.CreateNodeBalancerNode(context.Background(), nodebalancer.ID, config.ID, createOpts)
	if err != nil {
		t.Fatalf("Error creating NodeBalancer Config Node, got error %v", err)
	}

	teardown := func() {
		// delete the NodeBalancerNode to exercise the code
		if err := client.DeleteNodeBalancerNode(context.Background(), nodebalancer.ID, config.ID, node.ID); err != nil {
			e, ok := err.(*linodego.Error)
			// Tollerate 404 because Rebuild testing will delete all Nodes
			if !ok || e.Code != 404 {
				t.Fatalf("Expected to delete a NodeBalancer Config Node, but got %v", err)
			}
		}
		// delete the instance
		if err := client.DeleteInstance(context.Background(), instance.ID); err != nil {
			if t != nil {
				t.Errorf("Error deleting test Instance: %s", err)
			}
		}
		fixtureTeardown()
	}
	return client, nodebalancer, config, node, teardown, err
}
