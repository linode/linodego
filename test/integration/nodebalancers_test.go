package integration

import (
	"context"
	"strings"
	"testing"

	"github.com/linode/linodego"
	"github.com/stretchr/testify/require"
)

var (
	clientConnThrottle = 20
	label              = "go-test-def"
)

func TestNodeBalancer_Create_create_smoke(t *testing.T) {
	_, nodebalancer, teardown, err := setupNodeBalancer(t, "fixtures/TestNodeBalancer_Create", nil)
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

func TestNodeBalancer_Create_Type(t *testing.T) {
	_, nodebalancer, teardown, err := setupNodeBalancer(
		t,
		"fixtures/TestNodeBalancer_Create_Type",
		[]nbModifier{func(createOpts *linodego.NodeBalancerCreateOptions) {
			createOpts.Type = linodego.NBTypeCommon
		}},
	)
	defer teardown()

	if err != nil {
		t.Errorf("Error creating nodebalancer: %v", err)
	}

	// when comparing fixtures to random value Label will differ, compare the known suffix
	if !strings.Contains(*nodebalancer.Label, label) {
		t.Errorf("nodebalancer returned does not match nodebalancer create request")
	}
	// add this test case once the api supports returning it
	if nodebalancer.Type != linodego.NBTypeCommon {
		t.Errorf("nodebalancer returned type does not match the type of the nodebalancer create request")
	}

	assertDateSet(t, nodebalancer.Created)
	assertDateSet(t, nodebalancer.Updated)
}

func TestNodeBalancer_Create_with_ReservedIP(t *testing.T) {
	_, reserveIP, nodebalancer, teardown, err := setupNodeBalancerWithReservedIP(t, "fixtures/TestNodeBalancer_With_ReservedIP_Create")
	defer teardown()

	if err != nil {
		t.Errorf("Error creating nodebalancer: %v", err)
	}

	// when comparing fixtures to random value Label will differ, compare the known suffix
	if !strings.Contains(*nodebalancer.Label, label) {
		t.Errorf("nodebalancer returned does not match nodebalancer create request")
	}

	if reserveIP.Address != *nodebalancer.IPv4 {
		t.Errorf("nodebalancer address: %s does not matched requested reserved IP: %s", *nodebalancer.IPv4, reserveIP.Address)
	}

	assertDateSet(t, nodebalancer.Created)
	assertDateSet(t, nodebalancer.Updated)
}

func TestNodeBalancer_Create_with_vpc(t *testing.T) {
	_, nodebalancer, _, _, teardown, err := setupNodeBalancerWithVPC(t, "fixtures/TestNodeBalancer_With_VPC_Create")
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
	client, nodebalancer, teardown, err := setupNodeBalancer(t, "fixtures/TestNodeBalancer_Update", nil)
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

func TestNodeBalancers_List_smoke(t *testing.T) {
	client, _, teardown, err := setupNodeBalancer(t, "fixtures/TestNodeBalancers_List", nil)
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
	client, nodebalancer, teardown, err := setupNodeBalancer(t, "fixtures/TestNodeBalancer_Get", nil)
	defer teardown()
	if err != nil {
		t.Error(err)
	}

	_, err = client.GetNodeBalancer(context.Background(), nodebalancer.ID)
	if err != nil {
		t.Errorf("Error getting nodebalancer %d, expected *NodeBalancer, got error %v", nodebalancer.ID, err)
	}
}

func TestNodeBalancer_UDP(t *testing.T) {
	_, nodebalancer, teardown, err := setupNodeBalancer(
		t,
		"fixtures/TestNodeBalancer_UDP",
		[]nbModifier{
			func(options *linodego.NodeBalancerCreateOptions) {
				options.ClientUDPSessThrottle = linodego.Pointer(5)
			},
		},
	)
	defer teardown()
	if err != nil {
		t.Error(err)
	}

	require.Equal(t, 5, nodebalancer.ClientUDPSessThrottle)
}

type nbModifier func(options *linodego.NodeBalancerCreateOptions)

func setupNodeBalancer(t *testing.T, fixturesYaml string, nbModifiers []nbModifier) (*linodego.Client, *linodego.NodeBalancer, func(), error) {
	t.Helper()
	var fixtureTeardown func()
	client, fixtureTeardown := createTestClient(t, fixturesYaml)
	createOpts := linodego.NodeBalancerCreateOptions{
		Label:              &label,
		Region:             getRegionsWithCaps(t, client, []string{"NodeBalancers"})[0],
		ClientConnThrottle: &clientConnThrottle,
		FirewallID:         GetFirewallID(),
	}
	for _, modifier := range nbModifiers {
		modifier(&createOpts)
	}

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

func setupNodeBalancerWithReservedIP(t *testing.T, fixturesYaml string) (*linodego.Client, *linodego.InstanceIP, *linodego.NodeBalancer, func(), error) {
	t.Helper()
	var fixtureTeardown func()
	client, fixtureTeardown := createTestClient(t, fixturesYaml)
	reserveIP, err := client.ReserveIPAddress(context.Background(), linodego.ReserveIPOptions{
		Region: "us-east",
	})
	if err != nil {
		t.Fatalf("Failed to reserve IP %v", err)
	}
	t.Logf("Successfully reserved IP: %s", reserveIP.Address)

	createOpts := linodego.NodeBalancerCreateOptions{
		Label:              &label,
		Region:             "us-east",
		ClientConnThrottle: &clientConnThrottle,
		FirewallID:         GetFirewallID(),
		IPv4:               &reserveIP.Address,
	}

	nodebalancer, err := client.CreateNodeBalancer(context.Background(), createOpts)
	if err != nil {
		t.Fatalf("Error listing nodebalancers, expected struct, got error %v", err)
	}

	teardown := func() {
		if err := client.DeleteNodeBalancer(context.Background(), nodebalancer.ID); err != nil {
			t.Errorf("Expected to delete a nodebalancer, but got %v", err)
		}
		if err := client.DeleteReservedIPAddress(context.Background(), reserveIP.Address); err != nil {
			t.Errorf("Expected to delete a reserved IP, but got %v", err)
		}
		fixtureTeardown()
	}
	return client, reserveIP, nodebalancer, teardown, err
}

func setupNodeBalancerWithVPC(
	t *testing.T,
	fixturesYaml string,
	vpcModifier ...vpcModifier,
) (*linodego.Client, *linodego.NodeBalancer, *linodego.VPC, *linodego.VPCSubnet, func(), error) {
	t.Helper()
	var fixtureTeardown func()
	client, fixtureTeardown := createTestClient(t, fixturesYaml)
	vpc, subnet, vpcTeardown, err := createVPCWithSubnet(t, client, vpcModifier...)
	if err != nil {
		t.Errorf("Error creating vpc, got error %v", err)
	}
	createOpts := linodego.NodeBalancerCreateOptions{
		Label:              &label,
		Region:             vpc.Region,
		ClientConnThrottle: &clientConnThrottle,
		FirewallID:         GetFirewallID(),
		VPCs: []linodego.NodeBalancerVPCOptions{
			{
				IPv4Range: "192.168.0.64/30",
				IPv6Range: "",
				SubnetID:  subnet.ID,
			},
		},
	}

	nodebalancer, err := client.CreateNodeBalancer(context.Background(), createOpts)
	if err != nil {
		t.Fatalf("Error listing nodebalancers, expected struct, got error %v", err)
	}

	teardown := func() {
		if err := client.DeleteNodeBalancer(context.Background(), nodebalancer.ID); err != nil {
			t.Errorf("Expected to delete a nodebalancer, but got %v", err)
		}
		vpcTeardown()
		fixtureTeardown()
	}
	return client, nodebalancer, vpc, subnet, teardown, err
}
