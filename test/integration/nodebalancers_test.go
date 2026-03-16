package integration

import (
	"context"
	"strings"
	"testing"

	"github.com/linode/linodego"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	clientConnThrottle = 20
	label              = "go-test-def"
	premium40gbRegions = []string{"us-iad"} // No DevCloud region for premium_40gb type
	premiumRegions     = []string{
		"nl-ams",
		"jp-tyo-3",
		"sg-sin-2",
		"de-fra-2",
		"in-bom-2",
		"gb-lon",
		"us-lax",
		"id-cgk",
		"us-mia",
		"it-mil",
		"jp-osa",
		"in-maa",
		"se-sto",
		"br-gru",
		"us-sea",
		"fr-par",
		"us-iad",
		"pl-labkrk-2", // DevCloud
	}
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

func TestNodeBalancer_Create_with_backend_vpc(t *testing.T) {
	client, nodebalancer, _, _, teardown, err := setupNodeBalancerWithVPC(t, "fixtures/TestNodeBalancer_With_VPC_Create")
	defer teardown()
	require.NoErrorf(t, err, "Error creating nodebalancer: %s", err)

	// when comparing fixtures to random value Label will differ, compare the known suffix
	if !strings.Contains(*nodebalancer.Label, label) {
		t.Errorf("nodebalancer returned does not match nodebalancer create request")
	}

	assertDateSet(t, nodebalancer.Created)
	assertDateSet(t, nodebalancer.Updated)
	assert.NotEmpty(t, nodebalancer.IPv4)
	assert.NotEmpty(t, nodebalancer.IPv6)
	assert.Equal(t, "public", string(nodebalancer.FrontendAddressType))
	assert.Empty(t, nodebalancer.FrontendVPCSubnetID)

	vpcConfigs, err := client.ListNodeBalancerVPCConfigs(context.Background(), nodebalancer.ID, nil)
	require.NoErrorf(t, err, "Error listing nodebalancer VPC configs: %s", err)
	require.Len(t, vpcConfigs, 1, "Expected exactly one nodebalancer VPC config, got %d", len(vpcConfigs))
	assert.Equal(t, "backend", string(vpcConfigs[0].Purpose))

	vpcConfig, err := client.GetNodeBalancerVPCConfig(context.Background(), nodebalancer.ID, vpcConfigs[0].ID)
	require.NoErrorf(t, err, "Error getting nodebalancer VPC config: %s", err)
	assert.Equal(t, "backend", string(vpcConfig.Purpose))

	// TODO: Uncomment when API implementation of /backend_vpcs and /frontend_vpcs endpoints is finished
	//backendVPCs, err := client.ListNodeBalancerVPCBackendConfigs(context.Background(), nodebalancer.ID, nil)
	//require.NoErrorf(t, err, "Error listing nodebalancer backend VPC configs: %s", err)
	//require.Len(t, backendVPCs, 1, "Expected exactly one backend VPC, got %d", len(backendVPCs))
	//assert.Equal(t, "backend", backendVPCs[0].Purpose)
	//
	//frontendVPCs, err := client.ListNodeBalancerVPCFrontendConfigs(context.Background(), nodebalancer.ID, nil)
	//require.NoErrorf(t, err, "Error listing nodebalancer frontend VPC configs: %s", err)
	//require.Len(t, frontendVPCs, 0, "Expected no frontend VPCs, got %d", len(frontendVPCs))
}

func TestNodeBalancer_Create_with_frontend_ipv4_only(t *testing.T) {
	_, nodebalancer, _, subnet, teardown, err := setupNodeBalancerWithFrontendVPC(
		t,
		"fixtures/TestNodeBalancer_With_Frontend_VPC_IPv4_Create",
		linodego.NBTypePremium,
		TestSubnetIPv4,
		"")
	defer teardown()
	require.NoError(t, err, "Error setting up NodeBalancer with frontend VPC")

	assert.Empty(t, nodebalancer.IPv6)
	assert.Equal(t, "192.168.0.2", *nodebalancer.IPv4)
	assert.Equal(t, "vpc", string(nodebalancer.FrontendAddressType))
	assert.Equal(t, subnet.ID, *nodebalancer.FrontendVPCSubnetID)

	// TODO: Uncomment when API implementation of /backend_vpcs and /frontend_vpcs endpoints is finished
	//backendVPCs, err := client.ListNodeBalancerVPCBackendConfigs(context.Background(), nodebalancer.ID, nil)
	//require.NoErrorf(t, err, "Error listing nodebalancer backend VPC configs: %s", err)
	//require.Len(t, backendVPCs, 0, "Expected no backend VPC, got %d", len(backendVPCs))
	//
	//frontendVPCs, err := client.ListNodeBalancerVPCFrontendConfigs(context.Background(), nodebalancer.ID, nil)
	//require.NoErrorf(t, err, "Error listing nodebalancer frontend VPC configs: %s", err)
	//require.Len(t, frontendVPCs, 1, "Expected exactly one frontend VPC, got %d", len(frontendVPCs))
	//assert.Equal(t, "frontend", frontendVPCs[0].Purpose)
}

func TestNodeBalancer_Create_with_frontend_ipv6_only_fail(t *testing.T) {
	client, _ := createTestClient(t, fixturesYaml)
	vpc, subnet, teardown, err := createVPCWithSubnet(t, client)
	defer teardown()
	require.NoError(t, err, "Error creating VPC with subnet for NodeBalancer frontend VPC test")

	createOpts := linodego.NodeBalancerCreateOptions{
		Label:  &label,
		Region: vpc.Region,
		Type:   linodego.NBTypePremium,
		FrontendVPCs: []linodego.NodeBalancerFrontendVPCOptions{{
			SubnetID:  subnet.ID,
			IPv6Range: "/62",
		}}}

	_, err = client.CreateNodeBalancer(context.Background(), createOpts)
	require.ErrorContainsf(
		t,
		err,
		"No IPv6 subnets available in VPC",
		"No expected error returned, actual error value: %s", err,
	)
}

func TestNodeBalancer_Create_with_frontend_and_default_type_fail(t *testing.T) {
	client, _ := createTestClient(t, fixturesYaml)
	vpc, subnet, teardown, err := createVPCWithSubnet(t, client)
	defer teardown()
	require.NoError(t, err, "Error creating VPC with subnet for NodeBalancer frontend VPC test")

	createOpts := linodego.NodeBalancerCreateOptions{
		Label:  &label,
		Region: vpc.Region,
		Type:   linodego.NBTypeCommon,
		FrontendVPCs: []linodego.NodeBalancerFrontendVPCOptions{{
			SubnetID: subnet.ID,
		}}}

	_, err = client.CreateNodeBalancer(context.Background(), createOpts)
	require.ErrorContainsf(
		t,
		err,
		"NodeBalancer with frontend VPC IP must be premium",
		"No expected error returned, actual error value: %s", err,
	)
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
		Region:             getRegionsWithCaps(t, client, []string{linodego.CapabilityNodeBalancers})[0],
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

func setupNodeBalancerWithFrontendVPC(
	t *testing.T,
	fixturesYaml string,
	nbType linodego.NodeBalancerPlanType,
	IPv4 string,
	IPv6 string,
) (*linodego.Client, *linodego.NodeBalancer, *linodego.VPC, *linodego.VPCSubnet, func(), error) {
	t.Helper()
	var fixtureTeardown func()
	client, fixtureTeardown := createTestClient(t, fixturesYaml)
	vpc, subnet, vpcTeardown, err := createVPCWithSubnet(t, client)
	require.NoError(t, err, "Error creating VPC with subnet for NodeBalancer frontend VPC test")

	nbModifiers := []nbModifier{
		func(createOpts *linodego.NodeBalancerCreateOptions) {
			createOpts.Type = nbType
			createOpts.FrontendVPCs = []linodego.NodeBalancerFrontendVPCOptions{{
				SubnetID:  subnet.ID,
				IPv4Range: IPv4,
				IPv6Range: IPv6,
			}}
		},
	}

	_, nodebalancer, nbTeardown, err := setupNodeBalancer(
		t,
		"fixtures/TestNodeBalancer_Frontend_VPC",
		nbModifiers,
	)
	require.NoError(t, err, "Error creating nodebalancer with frontend VPC")

	teardown := func() {
		nbTeardown()
		vpcTeardown()
		fixtureTeardown()
	}
	return client, nodebalancer, vpc, subnet, teardown, err
}
