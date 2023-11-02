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
	client, instance, instanceConfig, instanceTeardown, err := setupInstanceWithoutDisks(
		t,
		fixturesYaml,
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
	}
	return client, vpc, vpcSubnet, instance, instanceConfig, teardownAll, err
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
			opts.Region = getRegionsWithCaps(t, client, []string{"VPCs"})[0]
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
			Purpose: InterfacePurposePublic,
		},
		{
			Purpose: InterfacePurposeVLAN,
			Label:   "testvlan",
		},
		{
			Purpose:  InterfacePurposeVPC,
			SubnetID: &vpcSubnet.ID,
			IPv4: &VPCIPv4{
				NAT1To1: "any",
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
		Purpose:  InterfacePurposeVPC,
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
		appednOpts.Purpose != intfc.Purpose ||
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

func TestInstance_ConfigInterfaces_List(t *testing.T) {
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
			Purpose: InterfacePurposePublic,
		},
		{
			Purpose: InterfacePurposeVLAN,
			Label:   "testvlan",
		},
		{
			Purpose:  InterfacePurposeVPC,
			SubnetID: &vpcSubnet.ID,
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

	if !reflect.DeepEqual(
		interfaceOptsList,
		updateConfigOpts.Interfaces,
	) {
		t.Error("failed to update linode interfaces: configs do not match")
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

	if !reflect.DeepEqual(
		interfaceOptsList,
		updateConfigOpts.Interfaces,
	) {
		t.Error("failed to update linode interfaces: configs do not match")
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
			Purpose:  InterfacePurposeVPC,
			SubnetID: &vpcSubnet.ID,
		},
	)
	if err != nil {
		t.Errorf("an error occurs when appending an interface to config %v: %v", config.ID, err)
	}
	updateOpts := intfc.GetUpdateOptions()
	updateOpts.Primary = true

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

	if !(updatedIntfc.Primary == updateOpts.Primary) {
		t.Errorf("updating interface %v didn't succeed", intfc.ID)
	}

	updateOpts.IPv4 = &VPCIPv4{
		VPC:     "192.168.0.10",
		NAT1To1: "any",
	}

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

	if !(updatedIntfc.Primary == updateOpts.Primary &&
		updateOpts.IPv4.VPC == updatedIntfc.IPv4.VPC) {
		t.Errorf("updating interface %v didn't succeed", intfc.ID)
	}
}

func TestInstance_Configs_List(t *testing.T) {
	client, instance, config, teardown, err := setupInstanceWithoutDisks(t, "fixtures/TestInstance_Configs_List")
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
	client, instance, config, teardown, err := setupInstanceWithoutDisks(t, "fixtures/TestInstance_Config_Update")
	defer teardown()
	if err != nil {
		t.Error(err)
	}

	updateConfigOpts := InstanceConfigUpdateOptions{
		Label:      "go-conf-test-" + randLabel(),
		Devices:    &InstanceConfigDeviceMap{},
		RootDevice: "/dev/root",
	}

	_, err = client.UpdateInstanceConfig(context.Background(), instance.ID, config.ID, updateConfigOpts)
	if err != nil {
		t.Error(err)
	}
}
