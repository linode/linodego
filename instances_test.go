package linodego_test

import (
	"context"
	"strconv"
	"testing"

	"github.com/linode/linodego"
)

func TestListInstances(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping test in short mode.")
	}
	client, instance, _, teardown, err := setupInstanceWithoutDisks(t, "fixtures/TestListInstances")
	defer teardown()

	listOpts := linodego.NewListOptions(1, "{\"id\": "+strconv.Itoa(instance.ID)+"}")
	linodes, err := client.ListInstances(context.Background(), listOpts)
	if err != nil {
		t.Errorf("Error listing instances, expected struct, got error %v", err)
	}
	if len(linodes) != 1 {
		t.Errorf("Expected a list of instances, but got %v", linodes)
	}

	if linodes[0].ID != instance.ID {
		t.Errorf("Expected list of instances to include test instance, but got %v", linodes)
	}
}

func TestGetInstance(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping test in short mode.")
	}
	client, instance, _, teardown, err := setupInstanceWithoutDisks(t, "fixtures/TestGetInstance")
	defer teardown()

	instanceGot, err := client.GetInstance(context.Background(), instance.ID)
	if err != nil {
		t.Errorf("Error getting instance: %s", err)
	}
	if instanceGot.ID != instance.ID {
		t.Errorf("Expected instance ID %d to match %d", instanceGot.ID, instance.ID)
	}

	if instance.Specs.Disk <= 0 {
		t.Errorf("Error parsing instance spec for disk size: %v", instance.Specs)
	}
}

func TestListInstanceDisks(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping test in short mode.")
	}
	client, instance, teardown, err := setupInstance(t, "fixtures/TestListInstanceDisks")
	defer teardown()

	disks, err := client.ListInstanceDisks(context.Background(), instance.ID, nil)
	if err != nil {
		t.Errorf("Error listing instance disks, expected struct, got error %v", err)
	}
	if len(disks) == 0 {
		t.Errorf("Expected a list of instance disks, but got %v", disks)
	}
}

func TestListInstanceConfigs(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping test in short mode.")
	}
	client, instance, config, teardown, err := setupInstanceWithoutDisks(t, "fixtures/TestListInstanceConfigs")
	defer teardown()

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

func TestUpdateInstanceConfig(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping test in short mode.")
	}
	client, instance, config, teardown, err := setupInstanceWithoutDisks(t, "fixtures/TestUpdateInstanceConfig")
	defer teardown()

	updateConfigOpts := linodego.InstanceConfigUpdateOptions{
		Label:      "bar",
		Devices:    linodego.InstanceConfigDeviceMap{},
		RootDevice: "/dev/root",
	}

	_, err = client.UpdateInstanceConfig(context.Background(), instance.ID, config.ID, updateConfigOpts)
	if err != nil {
		t.Error(err)
	}
}

func TestListInstanceVolumes(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping test in short mode.")
	}

	client, instance, config, teardown, err := setupInstanceWithoutDisks(t, "fixtures/TestListInstanceVolumes_instance")
	defer teardown()
	if err != nil {
		t.Error(err)
	}

	clientVol, volume, teardown, err := setupVolume(t, "fixtures/TestListInstanceVolumes")
	defer teardown()
	if err != nil {
		t.Error(err)
	}

	configOpts := linodego.InstanceConfigUpdateOptions{
		Label: "volume-test",
		Devices: linodego.InstanceConfigDeviceMap{
			SDA: &linodego.InstanceConfigDevice{
				VolumeID: volume.ID,
			},
		},
	}
	_, err = client.UpdateInstanceConfig(context.Background(), instance.ID, config.ID, configOpts)
	if err != nil {
		t.Error(err)
	}

	volumes, err := clientVol.ListInstanceVolumes(context.Background(), instance.ID, nil)
	if err != nil {
		t.Errorf("Error listing instance volumes, expected struct, got error %v", err)
	}
	if len(volumes) == 0 {
		t.Errorf("Expected an list of instance volumes, but got %v", volumes)
	}
}

func setupInstance(t *testing.T, fixturesYaml string) (*linodego.Client, *linodego.Instance, func(), error) {
	t.Helper()
	client, fixtureTeardown := createTestClient(t, fixturesYaml)
	falseBool := false
	createOpts := linodego.InstanceCreateOptions{
		Label:    "linodego-test-instance",
		RootPass: "R34lBAdP455",
		Region:   "us-west",
		Type:     "g6-nanode-1",
		Image:    "linode/debian9",
		Booted:   &falseBool,
	}
	instance, err := client.CreateInstance(context.Background(), createOpts)
	if err != nil {
		t.Errorf("Error creating test Instance: %s", err)
	}

	teardown := func() {
		if err := client.DeleteInstance(context.Background(), instance.ID); err != nil {
			t.Errorf("Error deleting test Instance: %s", err)
		}
		fixtureTeardown()
	}
	return client, instance, teardown, err
}

func setupInstanceWithoutDisks(t *testing.T, fixturesYaml string) (*linodego.Client, *linodego.Instance, *linodego.InstanceConfig, func(), error) {
	t.Helper()
	client, fixtureTeardown := createTestClient(t, fixturesYaml)
	falseBool := false
	createOpts := linodego.InstanceCreateOptions{
		Label:  "linodego-test-instance",
		Region: "us-west",
		Type:   "g6-nanode-1",
		Booted: &falseBool,
	}
	instance, err := client.CreateInstance(context.Background(), createOpts)
	if err != nil {
		t.Errorf("Error creating test Instance: %s", err)
		return nil, nil, nil, fixtureTeardown, err
	}
	configOpts := linodego.InstanceConfigCreateOptions{
		Label: "linodego-test-config",
	}
	config, err := client.CreateInstanceConfig(context.Background(), instance.ID, configOpts)
	if err != nil {
		t.Errorf("Error creating config: %s", err)
		return nil, nil, nil, fixtureTeardown, err
	}

	teardown := func() {
		if terr := client.DeleteInstance(context.Background(), instance.ID); terr != nil {
			t.Errorf("Error deleting test Instance: %s", terr)
		}
		fixtureTeardown()
	}
	return client, instance, config, teardown, err
}
