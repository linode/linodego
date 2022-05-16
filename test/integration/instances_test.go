package integration

import (
	"context"
	"fmt"
	"reflect"
	"strconv"
	"testing"
	"time"

	"github.com/linode/linodego"
)

type instanceModifier func(*linodego.InstanceCreateOptions)

func TestInstances_List(t *testing.T) {
	client, instance, _, teardown, err := setupInstanceWithoutDisks(t, "fixtures/TestInstances_List")
	defer teardown()

	if err != nil {
		t.Error(err)
	}

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

func TestInstance_Get(t *testing.T) {
	client, instance, _, teardown, err := setupInstanceWithoutDisks(t, "fixtures/TestInstance_Get")
	defer teardown()
	if err != nil {
		t.Error(err)
	}

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

	assertDateSet(t, instance.Created)
	assertDateSet(t, instance.Updated)
}

func TestInstance_Disks_List(t *testing.T) {
	client, instance, teardown, err := setupInstance(t, "fixtures/TestInstance_Disks_List")
	defer teardown()
	if err != nil {
		t.Error(err)
	}

	disks, err := client.ListInstanceDisks(context.Background(), instance.ID, nil)
	if err != nil {
		t.Errorf("Error listing instance disks, expected struct, got error %v", err)
	}
	if len(disks) == 0 {
		t.Errorf("Expected a list of instance disks, but got %v", disks)
	}
}

func TestInstance_Disk_Resize(t *testing.T) {
	client, instance, _, teardown, err := setupInstanceWithoutDisks(t, "fixtures/TestInstance_Disk_Resize")
	defer teardown()
	if err != nil {
		t.Error(err)
	}

	instance, err = client.WaitForInstanceStatus(context.Background(), instance.ID, linodego.InstanceOffline, 180)
	if err != nil {
		t.Errorf("Error waiting for instance readiness for resize: %s", err)
	}

	disk, err := client.CreateInstanceDisk(context.Background(), instance.ID, linodego.InstanceDiskCreateOptions{
		Label:      "test",
		Filesystem: "ext4",
		Size:       2000,
	})
	if err != nil {
		t.Errorf("Error creating disk for resize: %s", err)
	}

	disk, err = client.WaitForInstanceDiskStatus(context.Background(), instance.ID, disk.ID, linodego.DiskReady, 180)
	if err != nil {
		t.Errorf("Error waiting for disk readiness for resize: %s", err)
	}

	err = client.ResizeInstanceDisk(context.Background(), instance.ID, disk.ID, 4000)
	if err != nil {
		t.Errorf("Error resizing instance disk: %s", err)
	}
}

func TestInstance_Disk_ListMultiple(t *testing.T) {
	// This is a long running test
	client, instance1, teardown1, err := setupInstance(t, "fixtures/TestInstance_Disk_ListMultiple_Primary")
	defer teardown1()
	if err != nil {
		t.Error(err)
	}
	err = client.BootInstance(context.Background(), instance1.ID, 0)
	if err != nil {
		t.Error(err)
	}
	instance1, err = client.WaitForInstanceStatus(context.Background(), instance1.ID, linodego.InstanceRunning, 180)
	if err != nil {
		t.Errorf("Error waiting for instance readiness: %s", err)
	}

	disks, err := client.ListInstanceDisks(context.Background(), instance1.ID, nil)
	if err != nil {
		t.Error(err)
	}

	disk, err := client.WaitForInstanceDiskStatus(context.Background(), instance1.ID, disks[0].ID, linodego.DiskReady, 180)
	if err != nil {
		t.Errorf("Error waiting for disk readiness: %s", err)
	}

	imageLabel := fmt.Sprintf("linodego-test-image-%.d", time.Now().Second())
	imageCreateOptions := linodego.ImageCreateOptions{Label: imageLabel, DiskID: disk.ID}
	image, err := client.CreateImage(context.Background(), imageCreateOptions)

	defer client.DeleteImage(context.Background(), image.ID)
	if err != nil {
		t.Error(err)
	}

	client, instance2, _, teardown2, err := setupInstanceWithoutDisks(t, "fixtures/TestInstance_Disk_ListMultiple_Secondary")
	defer teardown2()
	if err != nil {
		t.Error(err)
	}
	instance2, err = client.WaitForInstanceStatus(context.Background(), instance2.ID, linodego.InstanceOffline, 180)
	if err != nil {
		t.Errorf("Error waiting for instance readiness: %s", err)
	}

	_, err = client.WaitForEventFinished(context.Background(), instance1.ID, linodego.EntityLinode, linodego.ActionDiskImagize, *disk.Created, 300)
	if err != nil {
		t.Errorf("Error waiting for imagize event: %s", err)
	}

	_, err = client.CreateInstanceDisk(context.Background(), instance2.ID, linodego.InstanceDiskCreateOptions{
		Label:    "linodego-test-instancedisk",
		Image:    image.ID,
		RootPass: "R34lBAdP455",
		Size:     2000,
	})
	if err != nil {
		t.Errorf("Error creating disk from private image: %s", err)
	}

	disk, err = client.CreateInstanceDisk(context.Background(), instance2.ID, linodego.InstanceDiskCreateOptions{
		Label: "linodego-test-2",
		Size:  2000,
	})

	if err != nil {
		t.Errorf("Error creating disk after a private image: %s", err)
	}

	disks, err = client.ListInstanceDisks(context.Background(), instance2.ID, nil)
	if err != nil {
		t.Errorf("Error listing instance disks, expected struct, got error %v", err)
	}
	if len(disks) != 2 {
		t.Errorf("Expected a list of instance disks, but got %v", disks)
	}
}

func TestInstance_Disk_ResetPassword(t *testing.T) {
	client, instance, _, teardown, err := setupInstanceWithoutDisks(t, "fixtures/TestInstance_Disk_ResetPassword")
	defer teardown()
	if err != nil {
		t.Error(err)
	}

	instance, err = client.WaitForInstanceStatus(context.Background(), instance.ID, linodego.InstanceOffline, 180)
	if err != nil {
		t.Errorf("Error waiting for instance readiness for password reset: %s", err)
	}

	disk, err := client.CreateInstanceDisk(context.Background(), instance.ID, linodego.InstanceDiskCreateOptions{
		Label:      "test",
		Filesystem: "ext4",
		Image:      "linode/debian9",
		RootPass:   "b4d_p455",
		Size:       2000,
	})
	if err != nil {
		t.Errorf("Error creating disk for password reset: %s", err)
	}

	instance, err = client.WaitForInstanceStatus(context.Background(), instance.ID, linodego.InstanceOffline, 180)
	if err != nil {
		t.Errorf("Error waiting for instance readiness after creating disk for password reset: %s", err)
	}
	disk, err = client.WaitForInstanceDiskStatus(context.Background(), instance.ID, disk.ID, linodego.DiskReady, 180)
	if err != nil {
		t.Errorf("Error waiting for disk readiness for password reset: %s", err)
	}

	err = client.PasswordResetInstanceDisk(context.Background(), instance.ID, disk.ID, "r34!_b4d_p455")
	if err != nil {
		t.Errorf("Error reseting password on instance disk: %s", err)
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

	updateConfigOpts := linodego.InstanceConfigUpdateOptions{
		Label:      "bar",
		Devices:    &linodego.InstanceConfigDeviceMap{},
		RootDevice: "/dev/root",
	}

	_, err = client.UpdateInstanceConfig(context.Background(), instance.ID, config.ID, updateConfigOpts)
	if err != nil {
		t.Error(err)
	}
}

func TestInstance_ConfigInterfaces_Update(t *testing.T) {
	client, instance, config, teardown, err := setupInstanceWithoutDisks(t,
		"fixtures/TestInstance_ConfigInterfaces_Update",
		func(opts *linodego.InstanceCreateOptions) {
			// Ensure we're in a region that supports VLANs
			opts.Region = "us-southeast"
		})
	defer teardown()
	if err != nil {
		t.Error(err)
	}

	updateConfigOpts := linodego.InstanceConfigUpdateOptions{
		Interfaces: []linodego.InstanceConfigInterface{
			{
				Purpose: linodego.InterfacePurposePublic,
			},
			{
				Purpose: linodego.InterfacePurposeVLAN,
				Label:   "linodego-cool-vlan",
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

	if !reflect.DeepEqual(result.Interfaces, updateConfigOpts.Interfaces) {
		t.Error("failed to update linode interfaces: configs do not match")
	}

	// Ensure that a nil value will not update interfaces
	_, err = client.UpdateInstanceConfig(context.Background(), instance.ID, config.ID, linodego.InstanceConfigUpdateOptions{})
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(result.Interfaces, updateConfigOpts.Interfaces) {
		t.Error("failed to update linode interfaces: configs do not match")
	}
}

func TestInstance_Volumes_List(t *testing.T) {
	client, instance, config, teardown, err := setupInstanceWithoutDisks(t, "fixtures/TestInstance_Volumes_List_Instance")
	defer teardown()
	if err != nil {
		t.Error(err)
	}

	clientVol, volume, teardown, err := setupVolume(t, "fixtures/TestInstance_Volumes_List")
	defer teardown()
	if err != nil {
		t.Error(err)
	}

	configOpts := linodego.InstanceConfigUpdateOptions{
		Label: "volume-test",
		Devices: &linodego.InstanceConfigDeviceMap{
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

func TestInstance_Rebuild(t *testing.T) {
	client, instance, _, teardown, err := setupInstanceWithoutDisks(t, "fixtures/TestInstance_Rebuild")
	defer teardown()

	if err != nil {
		t.Error(err)
	}

	_, err = client.WaitForEventFinished(context.Background(), instance.ID, linodego.EntityLinode, linodego.ActionLinodeCreate, *instance.Created, 180)

	if err != nil {
		t.Errorf("Error waiting for instance created: %s", err)
	}

	rebuildOpts := linodego.InstanceRebuildOptions{
		Image:    "linode/alpine3.11",
		RootPass: "R34lBAdP455",
	}
	instance, err = client.RebuildInstance(context.Background(), instance.ID, rebuildOpts)

	if err != nil {
		t.Error(err)
	}
}

func createInstance(t *testing.T, client *linodego.Client, modifiers ...instanceModifier) (*linodego.Instance, error) {
	if t != nil {
		t.Helper()
	}

	booted := false
	createOpts := linodego.InstanceCreateOptions{
		Label:    "linodego-test-instance",
		RootPass: "R34lBAdP455",
		Region:   "us-southeast",
		Type:     "g6-nanode-1",
		Image:    "linode/debian9",
		Booted:   &booted,
	}

	for _, modifier := range modifiers {
		modifier(&createOpts)
	}
	return client.CreateInstance(context.Background(), createOpts)
}

func setupInstance(t *testing.T, fixturesYaml string) (*linodego.Client, *linodego.Instance, func(), error) {
	if t != nil {
		t.Helper()
	}
	client, fixtureTeardown := createTestClient(t, fixturesYaml)

	instance, err := createInstance(t, client)
	if err != nil {
		t.Errorf("failed to create test instance: %s", err)
	}

	teardown := func() {
		if err := client.DeleteInstance(context.Background(), instance.ID); err != nil {
			if t != nil {
				t.Errorf("Error deleting test Instance: %s", err)
			}
		}
		fixtureTeardown()
	}
	return client, instance, teardown, err
}

func setupInstanceWithoutDisks(t *testing.T, fixturesYaml string, modifiers ...instanceModifier) (*linodego.Client, *linodego.Instance, *linodego.InstanceConfig, func(), error) {
	t.Helper()
	client, fixtureTeardown := createTestClient(t, fixturesYaml)
	falseBool := false
	createOpts := linodego.InstanceCreateOptions{
		Label:  "linodego-test-instance-wo-disk",
		Region: "us-southeast",
		Type:   "g6-nanode-1",
		Booted: &falseBool,
	}

	for _, modifier := range modifiers {
		modifier(&createOpts)
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
