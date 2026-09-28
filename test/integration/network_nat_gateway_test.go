package integration

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/linode/linodego/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupNATGateway(
	t *testing.T,
	fixtureYaml string,
	modifiers ...func(*linodego.NATGatewayCreateOptions),
) (*linodego.Client, *linodego.NATGateway) {
	t.Helper()
	client, fixtureTeardown := createTestClient(t, fixtureYaml)
	createOpts := linodego.NATGatewayCreateOptions{
		Label: "go-test-nat-gateway-" + randLabel(),
		Region: getRegionsWithCaps(t, client, []linodego.RegionCapability{
			linodego.CapabilityCloudFirewall,
			linodego.CapabilityLinodes,
			linodego.CapabilityLinodeInterfaces,
			linodego.CapabilityNATGateway,
			linodego.CapabilityVPCs,
			linodego.CapabilityVPCDualStack,
			linodego.CapabilityVPCIPv6Stack,
		})[0],
	}

	for _, modifier := range modifiers {
		modifier(&createOpts)
	}

	gateway, err := client.CreateNATGateway(context.Background(), createOpts)
	require.NoErrorf(t, err, "Error creating NAT Gateway: %v", err)

	t.Cleanup(func() {
		if err := client.DeleteNATGateway(context.Background(), gateway.ID); err != nil {
			t.Errorf("Error deleting NAT Gateway: %v", err)
		}
		fixtureTeardown()
	})
	return client, gateway
}

func setupNATGatewayReservedIP(t *testing.T, client *linodego.Client, region string) *linodego.InstanceIP {
	t.Helper()
	reservedIP, err := client.ReserveIPAddress(context.Background(), linodego.ReserveIPOptions{
		Region: region,
	})
	require.NoErrorf(t, err, "Error creating NAT Gateway: %v", err)

	t.Cleanup(func() {
		if err := client.DeleteReservedIPAddress(context.Background(), reservedIP.Address); err != nil {
			t.Errorf("Failed to delete reserved IP %s: %v", reservedIP.Address, err)
		}
	})
	return reservedIP
}

func verifyNATGateway(t *testing.T, gatewayCreated *linodego.NATGateway, gatewayRead *linodego.NATGateway) {
	assert.Equal(t, gatewayCreated.ID, gatewayRead.ID)
	assert.Equal(t, gatewayCreated.Region, gatewayRead.Region)
	assert.Equal(t, gatewayCreated.Label, gatewayRead.Label)
	assert.Equal(t, 0, len(gatewayRead.Addresses))
	assert.Greater(t, gatewayRead.AddressAutoscaleMax, 0)
	assert.Equal(t, 4096, gatewayRead.DefaultPortsPerInterface)
	assert.Equal(t, 0, gatewayRead.PortsetAssignments)
	assert.Greater(t, gatewayRead.PortsetCapacity, 0)
}

func verifyNATGatewayAddress(t *testing.T, addressNAT *linodego.NATGatewayAddressObject, addressReserved string) {
	assert.Equal(t, addressReserved, addressNAT.Address)
	assert.Equal(t, false, addressNAT.InUse)
	assert.Equal(t, 0, addressNAT.InterfaceCount)
	assert.True(t, strings.HasSuffix(addressNAT.InterfaceURL, fmt.Sprintf("/%s/interfaces", addressReserved)))
	assert.Equal(t, 0, addressNAT.PortsetAssignments)
	assert.Greater(t, addressNAT.PortsetCapacity, 0)
}

func verifyNATGatewayInterface(t *testing.T, gatewayInterface linodego.NATGatewayInterface, instance linodego.Instance, reservedIP linodego.InstanceIP) {
	assert.Equal(t, gatewayInterface.Linode.ID, instance.ID)
	assert.Equal(t, gatewayInterface.Addresses[0], reservedIP.Address)
	assert.Equal(t, gatewayInterface.Portsets[0].Address, reservedIP.Address)
}

func TestNATGateway_Create_smoke(t *testing.T) {
	client, gatewayCreated := setupNATGateway(t, "fixtures/TestNATGateway_Create_smoke")

	gateway, err := client.GetNATGateway(context.Background(), gatewayCreated.ID)
	require.NoErrorf(t, err, "Error retrieving NAT Gateway: %v", err)
	verifyNATGateway(t, gatewayCreated, gateway)

	f := linodego.Filter{}
	f.AddField(linodego.Eq, "label", gatewayCreated.Label)
	filter, err := f.MarshalJSON()
	require.NoErrorf(t, err, "Error marshalling filter: %v", err)

	gateways, err := client.ListNATGateways(context.Background(), &linodego.ListOptions{Filter: string(filter)})
	require.NoErrorf(t, err, "Error listing NAT Gateways: %v", err)
	require.Equal(t, 1, len(gateways), "Expected exactly one NAT Gateway in the list")
	verifyNATGateway(t, gatewayCreated, &gateways[0])
}

func TestNATGateway_Update(t *testing.T) {
	client, gatewayCreated := setupNATGateway(t, "fixtures/TestNATGateway_Update")

	newLabel := gatewayCreated.Label + "-updated"
	updateOpts := linodego.NATGatewayUpdateOptions{Label: &newLabel}
	gatewayUpdated, err := client.UpdateNATGateway(context.Background(), gatewayCreated.ID, updateOpts)
	require.NoErrorf(t, err, "Error updating NAT Gateway: %v", err)

	gateway, err := client.GetNATGateway(context.Background(), gatewayUpdated.ID)
	require.NoErrorf(t, err, "Error retrieving NAT Gateway: %v", err)
	assert.Equal(t, newLabel, gateway.Label)
}

func TestNATGateway_AssignReservedIP(t *testing.T) {
	client, gatewayCreated := setupNATGateway(
		t,
		"fixtures/TestNATGateway_Assign_ReservedIP",
		func(opts *linodego.NATGatewayCreateOptions) {
			opts.UseAutoscaling = linodego.Pointer(false)
		},
	)
	reservedIP := setupNATGatewayReservedIP(t, client, gatewayCreated.Region)

	_, err := client.NATGatewayAddAddress(
		context.Background(),
		gatewayCreated.ID,
		linodego.NATGatewayAddAddressOptions{
			Address: reservedIP.Address,
		},
	)
	require.NoErrorf(t, err, "Error adding Reserved IP address to NAT Gateway: %v", err)

	addresses, err := client.NATGatewayListAddresses(context.Background(), gatewayCreated.ID, nil)
	require.NoErrorf(t, err, "Error retrieving list of NAT Gateway addresses: %v", err)
	require.Len(t, addresses, 1)
	verifyNATGatewayAddress(t, &addresses[0], reservedIP.Address)

	address, err := client.NATGatewayGetAddress(context.Background(), gatewayCreated.ID, reservedIP.Address)
	require.NoErrorf(t, err, "Error retrieving a NAT Gateway address: %v", err)
	verifyNATGatewayAddress(t, address, reservedIP.Address)

	err = client.NATGatewayDeleteAddress(context.Background(), gatewayCreated.ID, reservedIP.Address)
	require.NoErrorf(t, err, "Error deleting address form NAT Gateway: %v", err)

	addresses, err = client.NATGatewayListAddresses(context.Background(), gatewayCreated.ID, nil)
	require.NoErrorf(t, err, "Error retrieving list of NAT Gateway addresses: %v", err)
	require.Len(t, addresses, 0)
}

func TestNATGateway_GetLinodeInterfaces(t *testing.T) {
	client, gatewayCreated := setupNATGateway(
		t,
		"fixtures/TestNATGateway_GetLinodeInterfaces",
		func(opts *linodego.NATGatewayCreateOptions) {
			opts.UseAutoscaling = linodego.Pointer(false)
		},
	)
	reservedIP := setupNATGatewayReservedIP(t, client, gatewayCreated.Region)

	_, err := client.NATGatewayAddAddress(
		context.Background(),
		gatewayCreated.ID,
		linodego.NATGatewayAddAddressOptions{
			Address: reservedIP.Address,
		},
	)
	require.NoErrorf(t, err, "Error adding Reserved IP address to NAT Gateway: %v", err)

	vpc, vpcSubnet, vpcTeardown, err := createVPCWithSubnet(
		t,
		client,
		func(c *linodego.Client, vo *linodego.VPCCreateOptions) {
			vo.Region = gatewayCreated.Region
		},
	)
	t.Cleanup(vpcTeardown)
	require.NoErrorf(t, err, "Error creating VPC with subnet: %v", err)

	vpcSubnetUpdateOpts := linodego.VPCSubnetUpdateOptions{
		Label:      vpcSubnet.Label,
		NATGateway: linodego.Pointer(linodego.VPCSubnetUpdateOptionsNATGateway{ID: linodego.DoublePointer(gatewayCreated.ID)}),
	}
	vpcSubnet, err = client.UpdateVPCSubnet(context.Background(), vpc.ID, vpcSubnet.ID, vpcSubnetUpdateOpts)
	require.NoErrorf(t, err, "Error updating VPC subnet with NAT Gateway: %v", err)

	instance, instanceTeardown, err := createInstanceWithLinodeInterfaces(
		t,
		client,
		true,
		[]linodego.LinodeInterfaceCreateOptions{
			{
				FirewallID: linodego.Pointer(firewallID),
				VPC: &linodego.VPCInterfaceCreateOptions{
					SubnetID: vpcSubnet.ID,
					IPv4: &linodego.VPCInterfaceIPv4CreateOptions{
						Addresses: []linodego.VPCInterfaceIPv4AddressCreateOptions{
							{
								Address: linodego.Pointer("auto"),
								Primary: linodego.Pointer(true),
							},
						},
					},
				},
			},
		},
		func(c *linodego.Client, opts *linodego.InstanceCreateOptions) {
			opts.Type = "g6-standard-1"
			opts.FirewallID = firewallID
			opts.Region = gatewayCreated.Region
		},
	)
	t.Cleanup(instanceTeardown)
	require.NoErrorf(t, err, "Error creating instance with interfaces: %v", err)

	instanceInterfaces, err := client.ListInterfaces(context.Background(), instance.ID, nil)
	require.NoErrorf(t, err, "Error fetching Linode interfaces: %v", err)
	require.Len(t, instanceInterfaces, 1)

	gatewayInterfaces, err := client.NATGatewayListInterfaces(context.Background(), gatewayCreated.ID, nil)
	require.NoErrorf(t, err, "Error listing NAT Gateway interfaces: %v", err)
	require.Len(t, gatewayInterfaces, 1)
	assert.Equal(t, gatewayInterfaces[0].ID, instanceInterfaces[0].ID)
	verifyNATGatewayInterface(t, gatewayInterfaces[0], *instance, *reservedIP)

	gatewayInterfaces, err = client.NATGatewayListAddressInterfaces(context.Background(), gatewayCreated.ID, reservedIP.Address, nil)
	require.NoErrorf(t, err, "Error listing NAT Gateway interfaces: %v", err)
	assert.Equal(t, gatewayInterfaces[0].ID, instanceInterfaces[0].ID)
	verifyNATGatewayInterface(t, gatewayInterfaces[0], *instance, *reservedIP)

	address, err := client.NATGatewayGetAddress(context.Background(), gatewayCreated.ID, reservedIP.Address)
	require.NoErrorf(t, err, "Error retrieving a NAT Gateway address: %v", err)
	assert.Equal(t, true, address.InUse)
	assert.Equal(t, 1, address.InterfaceCount)
	assert.Equal(t, 1, address.PortsetAssignments)
}

func TestNATGateway_GetTypes(t *testing.T) {
	client, fixtureTeardown := createTestClient(t, "fixtures/TestNATGateway_GetTypes")
	defer fixtureTeardown()

	types, err := client.NATGatewayGetTypes(context.Background(), nil)
	require.NoErrorf(t, err, "Error retrieving a NAT Gateway types: %v", err)
	assert.Greater(t, len(types), 0)
	assert.Equal(t, "g1-natgateway", types[0].ID)
	assert.Equal(t, "NAT Gateway", types[0].Label)
	assert.True(t, types[0].Price.Hourly >= 0 || types[0].Price.Monthly >= 0)
}

func TestNATGateway_GetSettings(t *testing.T) {
	client, fixtureTeardown := createTestClient(t, "fixtures/TestNATGateway_GetSettings")
	defer fixtureTeardown()

	settings, err := client.NATGatewayGetSettings(context.Background())
	require.NoErrorf(t, err, "Error retrieving a NAT Gateway settings: %v", err)
	for _, port := range settings.AllowedPortsPerInterface {
		assert.Contains(t, []int{4096, 8192, 16384}, port)
	}
	assert.Greater(t, settings.MaximumAutoscalingAddressesPerNATGateway, 0)
	assert.Greater(t, settings.MaximumReservedAddressesPerNATGateway, 0)
}

// TODO: Move to VPC tests
func TestNATGateway_VPCSubnetWithNATGateway(t *testing.T) {}
