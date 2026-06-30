package integration

import (
	"context"
	"testing"
	"time"

	"github.com/linode/linodego/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createInstanceWithLinodeInterfaces(
	t *testing.T,
	client *linodego.Client,
	enableCloudFirewall bool,
	interfaces []linodego.LinodeInstanceInterfaceCreateOptions,
	modifiers ...instanceModifier,
) (*linodego.Instance, func(), error) {
	if t != nil {
		t.Helper()
	}

	createOpts := linodego.InstanceCreateOptions{
		Label:                    "go-test-intf-" + randLabel(),
		RootPass:                 randPassword(),
		Region:                   getRegionsWithCaps(t, client, []linodego.RegionCapability{linodego.CapabilityLinodeInterfaces})[0],
		Type:                     "g6-nanode-1",
		Image:                    "linode/debian12",
		Booted:                   linodego.Pointer(false),
		InterfaceGeneration:      linodego.GenerationLinode,
		LinodeInstanceInterfaces: interfaces,
	}

	if enableCloudFirewall {
		for i := range createOpts.LinodeInterfaces {
			createOpts.LinodeInterfaces[i].FirewallID = linodego.Pointer(firewallID)
		}
	}

	for _, modifier := range modifiers {
		modifier(client, &createOpts)
	}
	instance, err := client.CreateInstance(context.Background(), createOpts)
	teardown := func() {
		if err := client.DeleteInstance(context.Background(), instance.ID); err != nil {
			if t != nil {
				t.Errorf("error deleting test Instance: %s", err)
			}
		}
	}
	return instance, teardown, err
}

func prepareMultipleRDMAInterfaces(amount int, subnet *linodego.VPCSubnet) []linodego.LinodeInstanceInterfaceCreateOptions {
	interfaces := make([]linodego.LinodeInstanceInterfaceCreateOptions, 0)

	for i := 1; i <= amount; i++ {
		interfaces = append(interfaces, linodego.LinodeInstanceInterfaceCreateOptions{
			LinodeInterfaceCreateOptions: linodego.LinodeInterfaceCreateOptions{
				FirewallID: linodego.Pointer(-1),
			},
			RDMAVPC: &linodego.RDMAVPCInterfaceCreateOptions{
				SubnetID: subnet.ID,
				IPv4: linodego.RDMAVPCInterfaceIPv4Options{
					Addresses: []linodego.RDMAVPCInterfaceIPv4AddressOptions{
						{Address: "auto", Primary: linodego.Pointer(true)},
					},
				},
			},
		})
	}
	return interfaces
}

func TestInstance_CreateWithLinodeInterfaces(
	t *testing.T,
) {
	client, fixtureTeardown := createTestClient(t, "fixtures/TestInstance_CreateWithLinodeInterfaces")
	t.Cleanup(fixtureTeardown)

	testRegion := getRegionsWithCaps(t, client, []linodego.RegionCapability{linodego.CapabilityVPCs, linodego.CapabilityLinodeInterfaces})[0]
	_, vpcSubnet, vpcTeardown, err := createVPCWithSubnet(
		t,
		client,
		func(c *linodego.Client, vo *linodego.VPCCreateOptions) {
			vo.Region = testRegion
		},
	)
	t.Cleanup(vpcTeardown)
	if err != nil {
		t.Fatalf("error creating a VPC with a subnet: %s", err)
	}

	instance, instanceTeardown, err := createInstanceWithLinodeInterfaces(
		t,
		client,
		true,
		[]linodego.LinodeInstanceInterfaceCreateOptions{
			{
				LinodeInterfaceCreateOptions: linodego.LinodeInterfaceCreateOptions{
					FirewallID: linodego.Pointer(firewallID),
					Public: &linodego.PublicInterfaceCreateOptions{
						IPv4: &linodego.PublicInterfaceIPv4CreateOptions{
							Addresses: []linodego.PublicInterfaceIPv4AddressCreateOptions{
								{
									Address: linodego.Pointer("auto"),
									Primary: linodego.Pointer(true),
								},
							},
						},
						IPv6: &linodego.PublicInterfaceIPv6CreateOptions{},
					},
				},
			},
			{
				LinodeInterfaceCreateOptions: linodego.LinodeInterfaceCreateOptions{
					FirewallID: linodego.Pointer(firewallID),
					VPC: &linodego.VPCInterfaceCreateOptions{
						SubnetID: vpcSubnet.ID,
						IPv4: &linodego.VPCInterfaceIPv4CreateOptions{
							Addresses: []linodego.VPCInterfaceIPv4AddressCreateOptions{
								{
									Address:        linodego.Pointer("auto"),
									Primary:        linodego.Pointer(true),
									NAT1To1Address: linodego.Pointer("auto"),
								},
							},
						},
					},
				},
			},
		},
		func(c *linodego.Client, opts *linodego.InstanceCreateOptions) {
			opts.Region = testRegion
		},
	)
	t.Cleanup(instanceTeardown)
	if err != nil {
		t.Fatalf("Error creating instance with interfaces: %s", err)
	}

	if instance.ID == 0 {
		t.Errorf("Expected a valid instance ID, got 0")
	}
}

func TestInstance_CreateWithRDMAVPCInterfaces(t *testing.T) {
	client, fixtureTeardown := createTestClient(t, "fixtures/TestInstance_CreateWithLinodeInterfaces")
	t.Cleanup(fixtureTeardown)

	//GPUDirect RDMA capability not available for now
	//region := getRegionsWithCaps(t, client, []string{linodego.CapabilityVPCs, linodego.CapabilityGPUDirectRDMA})
	testRegion := getRegionsWithCaps(t, client, []linodego.RegionCapability{linodego.CapabilityVPCs})[0]
	interfaceCreateOptions := make([]linodego.LinodeInstanceInterfaceCreateOptions, 0)

	// CREATE
	_, vpcSubnet, vpcTeardown, err := createVPCWithSubnet(
		t,
		client,
		func(c *linodego.Client, opts *linodego.VPCCreateOptions) {
			opts.Region = testRegion
			opts.VPCType = linodego.VPCTypeRegular
		},
	)
	require.NoErrorf(t, err, "Error creating RDMA VPC with subnet: %s", err)
	t.Cleanup(vpcTeardown)

	_, vpcSubnetRDMA, vpcRDMATeardown, err := createVPCWithSubnet(
		t,
		client,
		func(c *linodego.Client, opts *linodego.VPCCreateOptions) {
			opts.Region = testRegion
			opts.VPCType = linodego.VPCTypeRDMA
		},
	)
	require.NoErrorf(t, err, "Error creating RDMA VPC with subnet: %s", err)
	t.Cleanup(vpcRDMATeardown)

	// Add RDMA VPC interfaces
	multiRDMAInterfaces := prepareMultipleRDMAInterfaces(8, vpcSubnetRDMA)
	interfaceCreateOptions = append(interfaceCreateOptions, multiRDMAInterfaces...)

	//// Include at least one regular interface
	interfaceCreateOptions = append(interfaceCreateOptions, linodego.LinodeInstanceInterfaceCreateOptions{
		LinodeInterfaceCreateOptions: linodego.LinodeInterfaceCreateOptions{
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
	})

	instance, teardown, err := createInstanceWithLinodeInterfaces(
		t,
		client,
		false,
		interfaceCreateOptions,
		func(c *linodego.Client, opts *linodego.InstanceCreateOptions) {
			opts.Label = "go-test-rdma-" + randLabel()
			opts.RootPass = randPassword()
			opts.Image = "linode/ubuntu24.04"
			opts.Region = testRegion
			opts.Type = linodego.InstanceRDMAType
			opts.HostID = linodego.InstanceRDMAHostID
			opts.InterfaceGeneration = linodego.GenerationLinode
			opts.LinodeInstanceInterfaces = interfaceCreateOptions
		},
	)
	require.NoErrorf(t, err, "Error creating instance with RDMA interfaces: %s", err)
	t.Cleanup(teardown)

	instance, err = client.WaitForInstanceStatus(
		waitContext(t, 180*time.Second),
		instance.ID,
		linodego.InstanceOffline,
	)
	require.NoErrorf(t, err, "Error waiting for instance to be offline: %s", err)

	// READ
	// TODO: Defect for ListInterfaces needs to be resolved
	allInterfaces, err := client.ListInterfaces(context.Background(), instance.ID, nil)
	require.NoErrorf(t, err, "Error listing interfaces for RDMA instance: %s", err)
	assert.Equal(t, len(interfaceCreateOptions), len(allInterfaces), "Expected %d interfaces, got %d", len(interfaceCreateOptions), len(allInterfaces))

	basicRDMAInterface := allInterfaces[0]
	require.NotNil(t, basicRDMAInterface.RDMAVPC, "Expected interface to have RDMAVPC field populated")

	// UPDATE
	_, vpcSubnetRDMAUpdate, vpcRDMAUpdateTeardown, err := createVPCWithSubnet(
		t,
		client,
		func(c *linodego.Client, opts *linodego.VPCCreateOptions) {
			opts.Region = testRegion
			opts.VPCType = linodego.VPCTypeRDMA
		},
	)
	require.NoErrorf(t, err, "Error creating RDMA VPC with subnet: %s", err)
	t.Cleanup(vpcRDMAUpdateTeardown)

	updateOpts := linodego.LinodeInterfaceUpdateOptions{
		RDMAVPC: &linodego.RDMAVPCInterfaceUpdateOptions{
			SubnetID: vpcSubnetRDMAUpdate.ID,
		},
	}
	updatedRDMAInterface, err := client.UpdateInterface(context.Background(), instance.ID, basicRDMAInterface.ID, updateOpts)
	require.NoErrorf(t, err, "Error updating RDMA interface: %s", err)
	assert.Equal(t, basicRDMAInterface.ID, updatedRDMAInterface.ID, "Expected RDMA interface ID to remain the same after update")
	assert.Equal(t, basicRDMAInterface.RDMAVPC.SubnetID, vpcSubnetRDMAUpdate.ID, "Expected RDMA interface to be updated")

	// DELETE
	err = client.DeleteInterface(context.Background(), instance.ID, basicRDMAInterface.ID)
	require.Error(t, err, "Expected error deleting RDMA interface from RDMA instance")

	e, _ := err.(*linodego.Error)
	assert.Equal(t, 400, e.Code, "Expected error code 400, got: %d", e.Code)
	expectedErrorMessage := "RDMA VPC Interfaces cannot be deleted"
	assert.Contains(t, e.Message, expectedErrorMessage, "Expected error message to contain: %s, got: %s", expectedErrorMessage, e.Message)
}
