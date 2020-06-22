package integration

import (
	"context"
	"testing"

	"github.com/linode/linodego"
	"github.com/linode/linodego/pkg/errors"
)

var (
	testNodePort                   = "8080"
	testNodeLabel                  = "test-label"
	testNodeWeight                 = 10
	testNodeBalancerNodeCreateOpts = linodego.NodeBalancerNodeCreateOptions{
		Label:  testNodeLabel,
		Weight: testNodeWeight,
		Mode:   linodego.ModeAccept,
	}
)

func TestCreateNodeBalancerNode(t *testing.T) {
	_, _, _, node, teardown, err := setupNodeBalancerNode(t, "fixtures/TestCreateNodeBalancerNode")
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

func TestUpdateNodeBalancerNode(t *testing.T) {
	client, nodebalancer, config, node, teardown, err := setupNodeBalancerNode(t, "fixtures/TestUpdateNodeBalancerNode")
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

func TestListNodeBalancerNodes(t *testing.T) {
	client, nodebalancer, config, _, teardown, err := setupNodeBalancerNode(t, "fixtures/TestListNodeBalancerNodes")
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

func TestListNodeBalancerNodesMultiplePages(t *testing.T) {
	//TODO: replace hand crafted fixtures with unit tests
	// This fixture was hand-crafted to render an empty page 1 result, with a single result on page 2
	// "results:1,data:[],page:1,pages:2"  .. "results:1,data[{...}],page:2,pages:2"
	client, nodebalancer, config, _, teardown, err := setupNodeBalancerNode(t, "fixtures/TestListNodeBalancerNodesMultiplePages")
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

func TestGetNodeBalancerNode(t *testing.T) {
	client, nodebalancer, config, node, teardown, err := setupNodeBalancerNode(t, "fixtures/TestGetNodeBalancerNode")
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

func TestRebuildNodeBalancer(t *testing.T) {
	client, nodebalancer, config, _, teardown, err := setupNodeBalancerNode(t, "fixtures/TestRebuildNodeBalancer")
	defer teardown()
	if err != nil {
		t.Error(err)
	}

	nbcRebuildOpts := config.GetRebuildOptions()

	nbcGot, err := client.RebuildNodeBalancerConfig(context.Background(), nodebalancer.ID, config.ID, nbcRebuildOpts)
	if err != nil {
		t.Errorf("Error rebuilding nodebalancer config %d: %v", config.ID, err)
	}
	if nbcGot.Port != config.Port {
		t.Errorf("RebuildNodeBalancerConfig did not return the expected port")
	}
}

func setupNodeBalancerNode(t *testing.T, fixturesYaml string) (*linodego.Client, *linodego.NodeBalancer, *linodego.NodeBalancerConfig, *linodego.NodeBalancerNode, func(), error) {
	t.Helper()
	var fixtureTeardown func()
	client, nodebalancer, config, fixtureTeardown, err := setupNodeBalancerConfig(t, fixturesYaml)
	if err != nil {
		t.Errorf("Error creating nodebalancer config, got error %v", err)
	}

	client, instance, instanceTeardown, err := setupInstance(t, fixturesYaml+"Instance")
	if err != nil {
		t.Error(err)
	}

	instanceIP, err := client.AddInstanceIPAddress(context.Background(), instance.ID, false)
	if err != nil {
		t.Error(err)
	}

	createOpts := testNodeBalancerNodeCreateOpts
	createOpts.Address = instanceIP.Address + ":" + testNodePort
	node, err := client.CreateNodeBalancerNode(context.Background(), nodebalancer.ID, config.ID, createOpts)
	if err != nil {
		t.Errorf("Error creating NodeBalancer Config Node, got error %v", err)
	}

	teardown := func() {
		// delete the NodeBalancerNode to exercise the code
		if err := client.DeleteNodeBalancerNode(context.Background(), nodebalancer.ID, config.ID, node.ID); err != nil {
			e, ok := err.(*errors.Error)
			// Tollerate 404 because Rebuild testing will delete all Nodes
			if !ok || e.Code != 404 {
				t.Errorf("Expected to delete a NodeBalancer Config Node, but got %v", err)
			}
		}
		fixtureTeardown()
		instanceTeardown()
	}
	return client, nodebalancer, config, node, teardown, err
}
