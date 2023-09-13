package integration

import (
	"context"
	"encoding/base64"
	"reflect"
	"strconv"
	"testing"

	"github.com/linode/linodego"
)

type instanceModifier func(*linodego.Client, *linodego.InstanceCreateOptions)

func TestInstances_List(t *testing.T) {
	client, instance, _, teardown, err := setupInstanceWithoutDisks(
		t,
		"fixtures/TestInstances_List",
		func(client *linodego.Client, options *linodego.InstanceCreateOptions) {
			options.Region = "eu-west" // Override for metadata availability
		},
	)

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

	if linodes[0].HostUUID == "" {
		t.Errorf("failed to get instance HostUUID")
	}

	if linodes[0].HasUserData {
		t.Errorf("expected instance.HasUserData to be false, got true")
	}

	if linodes[0].Specs.GPUs < 0 {
		t.Errorf("failed to retrieve number of GPUs")
	}
}

func TestInstance_Get_smoke(t *testing.T) {
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

	if instance.HostUUID == "" {
		t.Errorf("failed to get instance HostUUID")
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
		Label:      "disk-test-" + randLabel(),
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

	imageLabel := "go-test-image-" + randLabel()
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
		Label:    "go-disk-test-" + randLabel(),
		Image:    image.ID,
		RootPass: "R34lBAdP455LONGLONGLONGLONG",
		Size:     2000,
	})
	if err != nil {
		t.Errorf("Error creating disk from private image: %s", err)
	}

	disk, err = client.CreateInstanceDisk(context.Background(), instance2.ID, linodego.InstanceDiskCreateOptions{
		Label: "go-disk-test-" + randLabel(),
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
		Label:      "go-disk-test-" + randLabel(),
		Filesystem: "ext4",
		Image:      "linode/debian9",
		RootPass:   "R34lBAdP455LONGLONGLONGLONG",
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
		Label:      "go-conf-test-" + randLabel(),
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
		func(client *linodego.Client, opts *linodego.InstanceCreateOptions) {
			// Ensure we're in a region that supports VLANs
			opts.Region = getRegionsWithCaps(t, client, []string{"vlans"})[0]
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
				Label:   instance.Label + "-r",
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
		Label: "go-vol-test" + randLabel(),
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
	client, instance, _, teardown, err := setupInstanceWithoutDisks(
		t,
		"fixtures/TestInstance_Rebuild",
		func(client *linodego.Client, options *linodego.InstanceCreateOptions) {
			options.Region = getRegionsWithCaps(t, client, []string{"Metadata"})[0]
		},
	)
	defer teardown()

	if err != nil {
		t.Error(err)
	}

	_, err = client.WaitForEventFinished(context.Background(), instance.ID, linodego.EntityLinode, linodego.ActionLinodeCreate, *instance.Created, 180)

	if err != nil {
		t.Errorf("Error waiting for instance created: %s", err)
	}

	rebuildOpts := linodego.InstanceRebuildOptions{
		Image: "linode/alpine3.15",
		Metadata: &linodego.InstanceMetadataOptions{
			UserData: base64.StdEncoding.EncodeToString([]byte("cool")),
		},
		RootPass: "R34lBAdP455LONGLONGLONGLONG",
		Type:     "g6-standard-2",
	}
	instance, err = client.RebuildInstance(context.Background(), instance.ID, rebuildOpts)

	if err != nil {
		t.Fatal(err)
	}

	if !instance.HasUserData {
		t.Fatal("expected instance.HasUserData to be true, got false")
	}
}

func TestInstance_Clone(t *testing.T) {
	var targetRegion string

	client, instance, teardownOriginalLinode, err := setupInstance(
		t, "fixtures/TestInstance_Clone",
		func(client *linodego.Client, options *linodego.InstanceCreateOptions) {
			targetRegion = getRegionsWithCaps(t, client, []string{"Metadata"})[0]

			options.Region = targetRegion
		})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(teardownOriginalLinode)

	_, err = client.WaitForEventFinished(
		context.Background(),
		instance.ID,
		linodego.EntityLinode,
		linodego.ActionLinodeCreate,
		*instance.Created,
		180,
	)

	if err != nil {
		t.Errorf("Error waiting for instance created: %s", err)
	}

	cloneOpts := linodego.InstanceCloneOptions{
		Region:    targetRegion,
		Type:      "g6-nanode-1",
		PrivateIP: true,
		Metadata: &linodego.InstanceMetadataOptions{
			UserData: base64.StdEncoding.EncodeToString([]byte("reallycooluserdata")),
		},
	}
	clonedInstance, err := client.CloneInstance(context.Background(), instance.ID, cloneOpts)

	t.Cleanup(func() {
		client.DeleteInstance(context.Background(), clonedInstance.ID)
	})

	if err != nil {
		t.Error(err)
	}

	_, err = client.WaitForEventFinished(
		context.Background(),
		instance.ID,
		linodego.EntityLinode,
		linodego.ActionLinodeClone,
		*clonedInstance.Created,
		240,
	)

	if err != nil {
		t.Fatal(err)
	}

	if clonedInstance.Image != instance.Image {
		t.Fatal("Clone instance image mismatched.")
	}

	clonedInstanceIPs, err := client.GetInstanceIPAddresses(
		context.Background(),
		clonedInstance.ID,
	)

	if err != nil {
		t.Fatal(err)
	}

	if len(clonedInstanceIPs.IPv4.Private) == 0 {
		t.Fatal("No private IPv4 assigned to the cloned instance.")
	}

	if !clonedInstance.HasUserData {
		t.Fatal("expected instance.HasUserData to be true, got false")
	}
}

func TestInstance_withMetadata(t *testing.T) {
	_, inst, _, teardown, err := setupInstanceWithoutDisks(t, "fixtures/TestInstance_withMetadata",
		func(client *linodego.Client, options *linodego.InstanceCreateOptions) {
			options.Metadata = &linodego.InstanceMetadataOptions{
				UserData: base64.StdEncoding.EncodeToString([]byte("reallycoolmetadata")),
			}
			options.Region = getRegionsWithCaps(t, client, []string{"Metadata"})[0]
		})
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(teardown)

	if !inst.HasUserData {
		t.Fatalf("expected instance.HasUserData to be true, got false")
	}
}

func createInstance(t *testing.T, client *linodego.Client, modifiers ...instanceModifier) (*linodego.Instance, error) {
	if t != nil {
		t.Helper()
	}

	booted := false
	createOpts := linodego.InstanceCreateOptions{
		Label:    "go-test-ins-" + randLabel(),
		RootPass: "R34lBAdP455LONGLONGLONGLONG",
		Region:   getRegionsWithCaps(t, client, []string{"linodes"})[0],
		Type:     "g6-nanode-1",
		Image:    "linode/debian9",
		Booted:   &booted,
	}

	for _, modifier := range modifiers {
		modifier(client, &createOpts)
	}
	return client.CreateInstance(context.Background(), createOpts)
}

func setupInstance(t *testing.T, fixturesYaml string, modifiers ...instanceModifier) (*linodego.Client, *linodego.Instance, func(), error) {
	if t != nil {
		t.Helper()
	}
	client, fixtureTeardown := createTestClient(t, fixturesYaml)

	instance, err := createInstance(t, client, modifiers...)
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
		Label:  "go-test-ins-wo-disk-" + randLabel(),
		Region: getRegionsWithCaps(t, client, []string{"linodes"})[0],
		Type:   "g6-nanode-1",
		Booted: &falseBool,
	}

	for _, modifier := range modifiers {
		modifier(client, &createOpts)
	}

	instance, err := client.CreateInstance(context.Background(), createOpts)
	if err != nil {
		t.Errorf("Error creating test Instance: %s", err)
		return nil, nil, nil, fixtureTeardown, err
	}
	configOpts := linodego.InstanceConfigCreateOptions{
		Label: "go-test-conf-" + randLabel(),
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
