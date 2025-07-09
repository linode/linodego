package integration

import (
	"context"
	"testing"

	"github.com/linode/linodego"
)

func createInstanceWithLinodeInterfaces(
	t *testing.T,
	client *linodego.Client,
	enableCloudFirewall bool,
	interfaces []linodego.LinodeInterfaceCreateOptions,
	modifiers ...instanceModifier,
) (*linodego.Instance, func(), error) {
	if t != nil {
		t.Helper()
	}

	createOpts := linodego.InstanceCreateOptions{
		Label:               "go-test-intf-" + randLabel(),
		RootPass:            randPassword(),
		Region:              getRegionsWithCaps(t, client, []string{linodego.CapabilityLinodeInterfaces})[0],
		Type:                "g6-nanode-1",
		Image:               "linode/debian12",
		Booted:              linodego.Pointer(false),
		InterfaceGeneration: linodego.GenerationLinode,
		LinodeInterfaces:    interfaces,
	}

	if enableCloudFirewall {
		for i := range createOpts.Interfaces {
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

func setupInstanceWithLinodeInterfaces(
	t *testing.T,
	fixturesYaml string,
	EnableCloudFirewall bool,
	interfaces []linodego.LinodeInterfaceCreateOptions,
	modifiers ...instanceModifier,
) (*linodego.Client, *linodego.Instance, func(), error) {
	if t != nil {
		t.Helper()
	}
	client, fixtureTeardown := createTestClient(t, fixturesYaml)

	instance, teardownInstance, err := createInstanceWithLinodeInterfaces(t, client, EnableCloudFirewall, interfaces, modifiers...)
	if err != nil {
		t.Errorf("failed to create test instance: %s", err)
	}

	teardown := func() {
		teardownInstance()
		fixtureTeardown()
	}

	return client, instance, teardown, err
}

func TestInstance_CreateWithLinodeInterfaces(
	t *testing.T,
) {
	client, fixtureTeardown := createTestClient(t, "fixtures/TestInstance_CreateWithLinodeInterfaces")
	t.Cleanup(fixtureTeardown)

	testRegion := getRegionsWithCaps(t, client, []string{linodego.CapabilityVPCs, linodego.CapabilityLinodeInterfaces})[0]
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
		[]linodego.LinodeInterfaceCreateOptions{
			{
				FirewallID: linodego.Pointer(firewallID),
				Public: &linodego.PublicInterfaceCreateOptions{
					IPv4: &linodego.PublicInterfaceIPv4CreateOptions{
						Addresses: []linodego.PublicInterfaceIPv4AddressCreateOptions{
							{
								Address: "auto",
								Primary: linodego.Pointer(true),
							},
						},
					},
					IPv6: &linodego.PublicInterfaceIPv6CreateOptions{},
				},
			},
			{
				FirewallID: linodego.Pointer(firewallID),
				VPC: &linodego.VPCInterfaceCreateOptions{
					SubnetID: vpcSubnet.ID,
					IPv4: &linodego.VPCInterfaceIPv4CreateOptions{
						Addresses: []linodego.VPCInterfaceIPv4AddressCreateOptions{
							{
								Address:        "auto",
								Primary:        linodego.Pointer(true),
								NAT1To1Address: linodego.Pointer("auto"),
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
