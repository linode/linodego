package integration

import (
	"context"
	"strings"
	"testing"

	"github.com/linode/linodego"
)

var (
	clientConnThrottle         = 20
	label                      = "go-test-def"
	testNodeBalancerCreateOpts = linodego.NodeBalancerCreateOptions{
		Label:              &label,
		Region:             "us-west",
		ClientConnThrottle: &clientConnThrottle,
	}
)

func TestNodeBalancer_Create(t *testing.T) {
	_, nodebalancer, teardown, err := setupNodeBalancer(t, "fixtures/TestNodeBalancer_Create")
	defer teardown()

	if err != nil {
		t.Errorf("Error creating nodebalancer: %v", err)
	}

	// when comparing fixtures to random value Label will differ, compare the known suffix
	if !strings.Contains(*nodebalancer.Label, label) {
		t.Errorf("nodebalancer returned does not match nodebalancer create request")
	}

	assertDateSet(t, nodebalancer.Created)
	assertDateSet(t, nodebalancer.Updated)
}

func TestNodeBalancer_Update(t *testing.T) {
	client, nodebalancer, teardown, err := setupNodeBalancer(t, "fixtures/TestNodeBalancer_Update")
	defer teardown()
	if err != nil {
		t.Error(err)
	}

	renamedLabel := *nodebalancer.Label + "_r"
	updateOpts := linodego.NodeBalancerUpdateOptions{
		Label: &renamedLabel,
	}
	nodebalancer, err = client.UpdateNodeBalancer(context.Background(), nodebalancer.ID, updateOpts)

	if err != nil {
		t.Errorf("Error renaming nodebalancer, %s", err)
	}

	if !strings.Contains(*nodebalancer.Label, renamedLabel) {
		t.Errorf("nodebalancer returned does not match nodebalancer create request")
	}
}

func TestNodeBalancers_List(t *testing.T) {
	client, _, teardown, err := setupNodeBalancer(t, "fixtures/TestNodeBalancers_List")
	defer teardown()
	if err != nil {
		t.Error(err)
	}

	nodebalancers, err := client.ListNodeBalancers(context.Background(), nil)
	if err != nil {
		t.Errorf("Error listing nodebalancers, expected struct, got error %v", err)
	}
	if len(nodebalancers) == 0 {
		t.Errorf("Expected a list of nodebalancers, but got %v", nodebalancers)
	}
}

func TestNodeBalancer_Get(t *testing.T) {
	client, nodebalancer, teardown, err := setupNodeBalancer(t, "fixtures/TestNodeBalancer_Get")
	defer teardown()
	if err != nil {
		t.Error(err)
	}

	_, err = client.GetNodeBalancer(context.Background(), nodebalancer.ID)
	if err != nil {
		t.Errorf("Error getting nodebalancer %d, expected *NodeBalancer, got error %v", nodebalancer.ID, err)
	}
}

func setupNodeBalancer(t *testing.T, fixturesYaml string) (*linodego.Client, *linodego.NodeBalancer, func(), error) {
	t.Helper()
	var fixtureTeardown func()
	client, fixtureTeardown := createTestClient(t, fixturesYaml)
	createOpts := testNodeBalancerCreateOpts
	nodebalancer, err := client.CreateNodeBalancer(context.Background(), createOpts)
	if err != nil {
		t.Fatalf("Error listing nodebalancers, expected struct, got error %v", err)
	}

	teardown := func() {
		if err := client.DeleteNodeBalancer(context.Background(), nodebalancer.ID); err != nil {
			t.Errorf("Expected to delete a nodebalancer, but got %v", err)
		}
		fixtureTeardown()
	}
	return client, nodebalancer, teardown, err
}
