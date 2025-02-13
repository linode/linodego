package integration

import (
	"context"
	"reflect"
	"testing"

	. "github.com/linode/linodego"
)

func setupVPCWithSubnetWithInstance(
	t *testing.T,
	fixturesYaml string,
	modifiers ...instanceModifier,
) (
	*Client,
	*VPC,
	*VPCSubnet,
	*Instance,
	*InstanceConfig,
	func(),
	error,
) {
	t.Helper()
	client, fixtureTeardown := createTestClient(t, fixturesYaml)
	instance, instanceConfig, instanceTeardown, err := createInstanceWithoutDisks(
		t,
		client, true,
		modifiers...,
	)
	if err != nil {
		if instanceTeardown != nil {
			instanceTeardown()
		}
		t.Fatal(err)
	}

	vpc, vpcSubnet, vpcWithSubnetTeardown, err := createVPCWithSubnet(
		t,
		client,
		func(client *Client, options *VPCCreateOptions) {
			options.Region = instance.Region
		},
	)
	if err != nil {
		t.Error(err)
	}

	teardownAll := func() {
		instanceTeardown()
		vpcWithSubnetTeardown()
		fixtureTeardown()
	}
	return client, vpc, vpcSubnet, instance, instanceConfig, teardownAll, err
}

func setupInstanceWithVPCAndNATOneToOne(t *testing.T, fixturesYaml string) (
	*Client,
	*VPC,
	*VPCSubnet,
	*Instance,
	*InstanceConfig,
	func(),
) {
	t.Helper()
	client, vpc, vpcSubnet, instance, config, teardown, err := setupVPCWithSubnetWithInstance(
		t,
		fixturesYaml,
		func(client *Client, opts *InstanceCreateOptions) {
			opts.Region = getRegionsWithCaps(t, client, []string{"Linodes", "VPCs"})[0]
		},
	)
	if err != nil {
		if teardown != nil {
			teardown()
		}
		t.Fatal(err)
	}

	updateConfigOpts := config.GetUpdateOptions()
	updateConfigOpts.Interfaces = []InstanceConfigInterfaceCreateOptions{
		{
			Purpose:  Pointer(InterfacePurposeVPC),
			SubnetID: &vpcSubnet.ID,
			IPv4: &VPCIPv4{
				NAT1To1: Pointer("any"),
			},
		},
	}
	config, err = client.UpdateInstanceConfig(context.Background(), instance.ID, config.ID, updateConfigOpts)
	if err != nil {
		teardown()
		t.Fatal(err)
	}

	return client, vpc, vpcSubnet, instance, config, teardown
}

func setupInstanceWith3Interfaces(t *testing.T, fixturesYaml string) (
	*Client,
	*VPC,
	*VPCSubnet,
	*Instance,
	*InstanceConfig,
	func(),
) {
	t.Helper()
	client, vpc, vpcSubnet, instance, config, teardown, err := setupVPCWithSubnetWithInstance(
		t,
		fixturesYaml,
		func(client *Client, opts *InstanceCreateOptions) {
			opts.Region = getRegionsWithCaps(t, client, []string{"Linodes", "VPCs"})[0]
		},
	)
	if err != nil {
		if teardown != nil {
			teardown()
		}
		t.Fatal(err)
	}

	updateConfigOpts := config.GetUpdateOptions()
	updateConfigOpts.Interfaces = []InstanceConfigInterfaceCreateOptions{
		{
			Purpose: Pointer(InterfacePurposePublic),
		},
		{
			Purpose: Pointer(InterfacePurposeVLAN),
			Label:   Pointer("testvlan"),
		},
		{
			Purpose:  Pointer(InterfacePurposeVPC),
			SubnetID: &vpcSubnet.ID,
			IPv4: &VPCIPv4{
				NAT1To1: Pointer("any"),
			},
		},
	}
	config, err = client.UpdateInstanceConfig(context.Background(), instance.ID, config.ID, updateConfigOpts)
	if err != nil {
		teardown()
		t.Fatal(err)
	}

	return client, vpc, vpcSubnet, instance, config, teardown
}

func TestInstance_ConfigInterfaces_AppendDelete(t *testing.T) {
	client, _, subnet, instance, config, teardown, err := setupVPCWithSubnetWithInstance(
		t,
		"fixtures/TestInstance_ConfigInterfaces_AppendDelete",
		func(client *Client, opts *InstanceCreateOptions) {
			// Ensure we're in a region that supports VLANs
			opts.Region = getRegionsWithCaps(t, client, []string{"vlans", "VPCs"})[0]
		},
	)
	defer teardown()
	if err != nil {
		t.Fatal(err)
	}

	appednOpts := InstanceConfigInterfaceCreateOptions{
		Purpose:  Pointer(InterfacePurposeVPC),
		SubnetID: &subnet.ID,
	}

	intfc, err := client.AppendInstanceConfigInterface(
		context.Background(),
		instance.ID,
		config.ID,
		appednOpts,
	)
	if err != nil {
		t.Error(err)
	}

	if intfc.ID == 0 ||
		*appednOpts.Purpose != intfc.Purpose ||
		*appednOpts.SubnetID != *intfc.SubnetID {
		t.Errorf(
			"failed to append an interface to instance %v config %v",
			instance.ID,
			config.ID,
		)
	}

	interfaces, err := client.ListInstanceConfigInterfaces(
		context.Background(),
		instance.ID,
		config.ID,
	)
	if err != nil {
		t.Error(err)
	}

	interfacesLength := len(interfaces)

	err = client.DeleteInstanceConfigInterface(
		context.Background(),
		instance.ID,
		config.ID,
		intfc.ID,
	)
	if err != nil {
		t.Error(err)
	}

	updatedInterfaces, err := client.ListInstanceConfigInterfaces(
		context.Background(),
		instance.ID,
		config.ID,
	)
	if err != nil {
		t.Error(err)
	}
	if len(updatedInterfaces) > interfacesLength {
		t.Errorf(
			"failed to delete interface %v of config %v of instance %v",
			intfc.ID,
			config.ID,
			instance.ID,
		)
	}
}

func TestInstance_ConfigInterfaces_Reorder(t *testing.T) {
	client, _, _, instance, config, teardown := setupInstanceWith3Interfaces(
		t,
		"fixtures/TestInstance_ConfigInterfaces_Reorder",
	)
	defer teardown()

	desiredIDs := []int{
		config.Interfaces[1].ID,
		config.Interfaces[0].ID,
		config.Interfaces[2].ID,
	}
	err := client.ReorderInstanceConfigInterfaces(
		context.Background(),
		instance.ID,
		config.ID,
		InstanceConfigInterfacesReorderOptions{
			IDs: desiredIDs,
		},
	)
	if err != nil {
		t.Error(err)
	}

	reorderedInterfacesConfig, err := client.GetInstanceConfig(
		context.Background(),
		instance.ID,
		config.ID,
	)
	if err != nil {
		t.Error(err)
	}
	reorderedIDs := []int{
		reorderedInterfacesConfig.Interfaces[0].ID,
		reorderedInterfacesConfig.Interfaces[1].ID,
		reorderedInterfacesConfig.Interfaces[2].ID,
	}

	if !reflect.DeepEqual(reorderedIDs, desiredIDs) {
		t.Errorf(
			"interface IDs reordering failed, desired IDs: %v, appeared IDs: %v",
			desiredIDs,
			reorderedIDs,
		)
	}
}

func TestInstance_ConfigInterfaces_List_smoke(t *testing.T) {
	client, _, _, instance, config, teardown := setupInstanceWith3Interfaces(
		t,
		"fixtures/TestInstance_ConfigInterfaces_List",
	)
	defer teardown()

	interfaces, err := client.ListInstanceConfigInterfaces(
		context.Background(),
		instance.ID,
		config.ID,
	)
	if err != nil {
		t.Error(err)
	}

	if !(len(interfaces) == 3 &&
		interfaces[0].ID != 0 &&
		interfaces[1].ID != 0 &&
		interfaces[2].ID != 0) {
		t.Errorf("failed to list all interfaces of config %v", config.ID)
	}
}

// testing config interfaces update via config API
func TestInstance_ConfigInterfaces_Update(t *testing.T) {
	client, _, vpcSubnet, instance, config, teardown, err := setupVPCWithSubnetWithInstance(
		t,
		"fixtures/TestInstance_ConfigInterfaces_Update",
		func(client *Client, opts *InstanceCreateOptions) {
			// Ensure we're in a region that supports VLANs
			opts.Region = getRegionsWithCaps(t, client, []string{"vlans", "VPCs"})[0]
		},
	)
	defer teardown()
	if err != nil {
		t.Error(err)
	}

	updateConfigOpts := config.GetUpdateOptions()

	updateConfigOpts.Interfaces = []InstanceConfigInterfaceCreateOptions{
		{
			Purpose: Pointer(InterfacePurposePublic),
		},
		{
			Purpose: Pointer(InterfacePurposeVLAN),
			Label:   Pointer("testvlan"),
		},
		{
			Purpose:  Pointer(InterfacePurposeVPC),
			SubnetID: &vpcSubnet.ID,
			IPv4: &VPCIPv4{
				VPC: Pointer("192.168.0.87"),
			},
		},
	}

	_, err = client.UpdateInstanceConfig(context.Background(), instance.ID, config.ID, updateConfigOpts)
	if err != nil {
		t.Error(err)
	}

	result, err := client.GetInstanceConfig(context.Background(), instance.ID, config.ID)
	if err != nil {
		t.Error(err)
	}

	interfaceOptsList := make([]InstanceConfigInterfaceCreateOptions, len(result.Interfaces))
	for index, configInterface := range result.Interfaces {
		interfaceOptsList[index] = configInterface.GetCreateOptions()
	}

	// Compare each interface in the result with the expected values
	for i, instanceConfigInterface := range result.Interfaces {
		expectedInstanceConfigInterface := updateConfigOpts.Interfaces[i]
		checkInstanceConfigInterfaceMismatch(t, i, instanceConfigInterface, expectedInstanceConfigInterface)
	}

	// Ensure that a nil value will not update interfaces
	result, err = client.UpdateInstanceConfig(context.Background(), instance.ID, config.ID, InstanceConfigUpdateOptions{})
	if err != nil {
		t.Error(err)
	}
	interfaceOptsList = make([]InstanceConfigInterfaceCreateOptions, len(result.Interfaces))
	for index, configInterface := range result.Interfaces {
		interfaceOptsList[index] = configInterface.GetCreateOptions()
	}

	// Compare each interface in the result with the expected values
	for i, instanceConfigInterface := range result.Interfaces {
		expectedInstanceConfigInterface := updateConfigOpts.Interfaces[i]
		checkInstanceConfigInterfaceMismatch(t, i, instanceConfigInterface, expectedInstanceConfigInterface)
	}
}

// testing config interface update via interfaces API
func TestInstance_ConfigInterface_Update(t *testing.T) {
	client, _, vpcSubnet, instance, config, teardown, err := setupVPCWithSubnetWithInstance(
		t,
		"fixtures/TestInstance_ConfigInterface_Update",
		func(client *Client, opts *InstanceCreateOptions) {
			// Ensure we're in a region that supports VLANs
			opts.Region = getRegionsWithCaps(t, client, []string{"vlans", "VPCs"})[0]
		},
	)
	defer teardown()
	if err != nil {
		t.Error(err)
	}

	intfc, err := client.AppendInstanceConfigInterface(
		context.Background(),
		instance.ID,
		config.ID,
		InstanceConfigInterfaceCreateOptions{
			Purpose:  Pointer(InterfacePurposeVPC),
			SubnetID: &vpcSubnet.ID,
			IPRanges: &[]string{"192.168.0.5/32"},
		},
	)
	if err != nil {
		t.Errorf("an error occurs when appending an interface to config %v: %v", config.ID, err)
	}
	updateOpts := intfc.GetUpdateOptions()
	updateOpts.Primary = Pointer(true)

	updatedIntfc, err := client.UpdateInstanceConfigInterface(
		context.Background(),
		instance.ID,
		config.ID,
		intfc.ID,
		updateOpts,
	)
	if err != nil {
		t.Errorf("an error occurs when updating an interface in config %v", config.ID)
	}

	if updateOpts.Primary == nil || updatedIntfc.Primary != *updateOpts.Primary {
		t.Errorf("updating interface %v didn't succeed", intfc.ID)
	}

	if updatedIntfc.IPRanges[0] != "192.168.0.5/32" {
		t.Errorf("unexpected value for IPRanges: %s", updatedIntfc.IPRanges[0])
	}

	updateOpts.IPv4 = &VPCIPv4{
		VPC:     Pointer("192.168.0.10"),
		NAT1To1: Pointer("any"),
	}
	newIPRanges := make([]string, 0)
	updateOpts.IPRanges = &newIPRanges

	updatedIntfc, err = client.UpdateInstanceConfigInterface(
		context.Background(),
		instance.ID,
		config.ID,
		intfc.ID,
		updateOpts,
	)
	if err != nil {
		t.Errorf("an error occurs when updating an interface in config %v", config.ID)
	}

	if !(updatedIntfc.Primary == *updateOpts.Primary &&
		*updateOpts.IPv4.VPC == *updatedIntfc.IPv4.VPC) {
		t.Errorf("updating interface %v didn't succeed", intfc.ID)
	}

	if len(updatedIntfc.IPRanges) > 0 {
		t.Errorf("expected IPRanges to be empty, got %d entries", len(updatedIntfc.IPRanges))
	}
}

func TestInstance_Configs_List(t *testing.T) {
	client, instance, config, teardown, err := setupInstanceWithoutDisks(t, "fixtures/TestInstance_Configs_List", true)
	defer teardown()
	if err != nil {
		t.Error(err)
	}

	configs, err := client.ListInstanceConfigs(context.Background(), instance.ID, nil)
	if err != nil {
		t.Errorf("Error listing instance configs, expected struct, got error %v", err)
	}
	if len(configs) == 0 {
		t.Errorf("Expected a list of instance configs, but got %v", configs)
	}
	if configs[0].ID != config.ID {
		t.Errorf("Expected config id %d, got %d", configs[0].ID, config.ID)
	}
}

func TestInstance_Config_Update(t *testing.T) {
	client, instance, config, teardown, err := setupInstanceWithoutDisks(t, "fixtures/TestInstance_Config_Update", true)
	defer teardown()
	if err != nil {
		t.Error(err)
	}

	updateConfigOpts := InstanceConfigUpdateOptions{
		Label:      Pointer("go-conf-test-" + randLabel()),
		Devices:    &InstanceConfigDeviceMap{},
		RootDevice: Pointer("/dev/root"),
	}

	_, err = client.UpdateInstanceConfig(context.Background(), instance.ID, config.ID, updateConfigOpts)
	if err != nil {
		t.Error(err)
	}
}

func checkInstanceConfigInterfaceMismatch(t *testing.T, i int, icf InstanceConfigInterface, expectedIcf InstanceConfigInterfaceCreateOptions) {
	// Compare Purpose
	if icf.Purpose != *expectedIcf.Purpose {
		t.Errorf("Interface %d: Purpose mismatch: expected %v, got %v", i, *expectedIcf.Purpose, icf.Purpose)
	}

	// Compare Label
	if expectedIcf.Label != nil && icf.Label != *expectedIcf.Label {
		t.Errorf("Interface %d: Label mismatch: expected %v, got %v", i, *expectedIcf.Label, icf.Label)
	}

	// Compare SubnetID
	if icf.SubnetID != nil && expectedIcf.SubnetID != nil && *icf.SubnetID != *expectedIcf.SubnetID {
		t.Errorf("Interface %d: SubnetID mismatch: expected %v, got %v", i, *expectedIcf.SubnetID, *icf.SubnetID)
	}

	// Compare IPv4 VPC
	if icf.IPv4 != nil && expectedIcf.IPv4 != nil && *icf.IPv4.VPC != *expectedIcf.IPv4.VPC {
		t.Errorf("Interface %d: IPv4 VPC mismatch: expected %v, got %v", i, *expectedIcf.IPv4.VPC, *icf.IPv4.VPC)
	}

	// Compare IPRanges
	if icf.IPRanges != nil && expectedIcf.IPRanges != nil && !reflect.DeepEqual(icf.IPRanges, *expectedIcf.IPRanges) {
		t.Errorf("Interface %d: IPRanges mismatch: expected %v, got %v", i, *expectedIcf.IPRanges, icf.IPRanges)
	}
}
