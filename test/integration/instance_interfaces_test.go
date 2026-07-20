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
	t.Helper()

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
		for i := range createOpts.LinodeInstanceInterfaces {
			createOpts.LinodeInstanceInterfaces[i].FirewallID = linodego.Pointer(firewallID)
		}
	}

	for _, modifier := range modifiers {
		modifier(client, &createOpts)
	}
	instance, err := client.CreateInstance(context.Background(), createOpts)
	require.NoErrorf(t, err, "Error creating test instance: %s", err)

	teardown := func() {
		// Use a fresh, independent context
		ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
		defer cancel()

		// Wait for the instance to be deleted
		p, err := client.NewEventPoller(
			context.Background(), instance.ID, linodego.EntityLinode, linodego.ActionLinodeDelete,
		)
		require.NoErrorf(t, err, "Error creating event poller for instance deletion: %s", err)

		err = client.DeleteInstance(context.Background(), instance.ID)
		require.NoErrorf(t, err, "Error deleting test instance: %s", err)

		// NOTE: The linode_delete event is sometimes reported as "failed" by the
		// API even though the instance is actually deleted successfully.
		// Rather than treating that as fatal, we log it and fall back to
		// GetInstance returning a 404 as the authoritative check that the instance
		// was actually deleted.
		event, err := p.WaitForFinished(ctx)
		if err != nil {
			t.Logf("Warning: Error waiting for instance deletion event (instance deletion will still be verified directly): %s", err)
		} else if event.Action != linodego.ActionLinodeDelete {
			t.Errorf("Expected event action %s, got %s", linodego.ActionLinodeDelete, event.Action)
		}

		_, err = client.GetInstance(context.Background(), instance.ID)
		require.Truef(t, linodego.IsNotFound(err), "Expected instance to be deleted (404), got: %v", err)
	}
	return instance, teardown, err
}

func createVPCWithSubnetAndType(t *testing.T, client *linodego.Client, region string, vpcType linodego.VPCType) (
	*linodego.VPC,
	*linodego.VPCSubnet,
	func(),
	error,
) {
	vpc, subnet, teardown, err := createVPCWithSubnet(
		t,
		client,
		func(c *linodego.Client, opts *linodego.VPCCreateOptions) {
			opts.Region = region
			opts.VPCType = vpcType
		},
	)
	return vpc, subnet, teardown, err
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
	t.Skip("Skipping test because Linode with RDMA interfaces requires manual infra changes at the moment")

	client, fixtureTeardown := createTestClient(t, "fixtures/TestInstance_CreateWithRDMAVPCInterfaces")
	t.Cleanup(fixtureTeardown)

	// GPUDirect RDMA capability not available for now
	// testRegion := getRegionsWithCaps(t, client, []linodego.RegionCapability{linodego.CapabilityVPCs, linodego.CapabilityGPUDirectRDMA})[0]
	testRegion := getRegionsWithCaps(t, client, []linodego.RegionCapability{linodego.CapabilityVPCs})[0]
	interfaceCreateOptions := make([]linodego.LinodeInstanceInterfaceCreateOptions, 0)

	// CREATE
	_, vpcSubnet, vpcTeardown, err := createVPCWithSubnetAndType(t, client, testRegion, linodego.VPCTypeRegular)
	require.NoErrorf(t, err, "Error creating Regular VPC with subnet: %s", err)
	t.Cleanup(vpcTeardown)

	_, vpcSubnetRDMA, vpcRDMATeardown, err := createVPCWithSubnetAndType(t, client, testRegion, linodego.VPCTypeRDMA)
	require.NoErrorf(t, err, "Error creating RDMA VPC with subnet: %s", err)
	t.Cleanup(vpcRDMATeardown)

	_, vpcSubnetRDMAUpdate, vpcRDMAUpdateTeardown, err := createVPCWithSubnetAndType(t, client, testRegion, linodego.VPCTypeRDMA)
	require.NoErrorf(t, err, "Error creating RDMA VPC with subnet: %s", err)
	t.Cleanup(vpcRDMAUpdateTeardown)

	// Include RDMA VPC interfaces
	multiRDMAInterfaces := prepareMultipleRDMAInterfaces(8, vpcSubnetRDMA)
	interfaceCreateOptions = append(interfaceCreateOptions, multiRDMAInterfaces...)

	// Include (at least one) regular interface
	interfaceCreateOptions = append(interfaceCreateOptions, linodego.LinodeInstanceInterfaceCreateOptions{
		LinodeInterfaceCreateOptions: linodego.LinodeInterfaceCreateOptions{
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

	instance, instanceTeardown, err := createInstanceWithLinodeInterfaces(
		t,
		client,
		false,
		interfaceCreateOptions,
		func(c *linodego.Client, opts *linodego.InstanceCreateOptions) {
			opts.Label = "go-test-rdma-" + randLabel()
			opts.RootPass = randPassword()
			opts.Image = "linode/ubuntu24.04"
			opts.Region = testRegion
			// opts.Type = linodego.InstanceRDMAType
			// opts.HostID = linodego.InstanceRDMAHostID
			opts.InterfaceGeneration = linodego.GenerationLinode
			opts.LinodeInstanceInterfaces = interfaceCreateOptions
		},
	)
	require.NoErrorf(t, err, "Error creating instance with RDMA interfaces: %s", err)
	t.Cleanup(instanceTeardown)

	instance, err = client.WaitForInstanceStatus(
		waitContext(t, 180*time.Second),
		instance.ID,
		linodego.InstanceOffline,
	)
	require.NoErrorf(t, err, "Error waiting for instance to be running: %s", err)

	// READ
	allInterfaces, err := client.ListInterfaces(context.Background(), instance.ID, nil)
	require.NoErrorf(t, err, "Error listing interfaces for RDMA instance: %s", err)
	assert.Equal(t, len(interfaceCreateOptions), len(allInterfaces), "Expected %d interfaces, got %d", len(interfaceCreateOptions), len(allInterfaces))

	basicRDMAInterface := allInterfaces[0]
	require.NotNil(t, basicRDMAInterface.RDMAVPC, "Expected interface to have RDMAVPC field populated")

	// UPDATE
	updateOpts := linodego.LinodeInterfaceUpdateOptions{
		RDMAVPC: &linodego.RDMAVPCInterfaceUpdateOptions{
			SubnetID: vpcSubnetRDMAUpdate.ID,
		},
	}
	updatedRDMAInterface, err := client.UpdateInterface(context.Background(), instance.ID, basicRDMAInterface.ID, updateOpts)
	require.NoErrorf(t, err, "Error updating RDMA interface: %s", err)
	require.NotNil(t, updatedRDMAInterface.RDMAVPC, "Expected updated interface to have RDMAVPC field populated")
	assert.Equal(t, basicRDMAInterface.ID, updatedRDMAInterface.ID, "Expected RDMA interface ID to remain the same after update")
	assert.Equal(t, vpcSubnetRDMAUpdate.ID, updatedRDMAInterface.RDMAVPC.SubnetID, "Expected RDMA interface to be updated")

	// DELETE
	err = client.DeleteInterface(context.Background(), instance.ID, basicRDMAInterface.ID)
	require.Error(t, err, "Expected error deleting RDMA interface from RDMA instance")

	var e *linodego.Error
	require.ErrorAsf(t, err, &e, "Expected error to be of type *linodego.Error, got: %T", err)
	assert.Equal(t, 400, e.Code, "Expected error code 400, got: %d", e.Code)
	expectedErrorMessage := "RDMA VPC Interfaces cannot be deleted"
	assert.Contains(t, e.Message, expectedErrorMessage, "Expected error message to contain: %s, got: %s", expectedErrorMessage, e.Message)
}
