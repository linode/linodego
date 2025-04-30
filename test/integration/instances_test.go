package integration

import (
	"context"
	"encoding/base64"
	"slices"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/linode/linodego"
	"github.com/stretchr/testify/require"
)

type instanceModifier func(*linodego.Client, *linodego.InstanceCreateOptions)

func TestInstances_List_smoke(t *testing.T) {
	client, instance, _, teardown, err := setupInstanceWithoutDisks(
		t,
		"fixtures/TestInstances_List", true,
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
	client, instance, _, teardown, err := setupInstanceWithoutDisks(t, "fixtures/TestInstance_Get", true)
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

func TestInstance_GetTransfer(t *testing.T) {
	client, instance, _, teardown, err := setupInstanceWithoutDisks(t, "fixtures/TestInstance_GetTransfer", true)
	defer teardown()
	if err != nil {
		t.Error(err)
	}

	_, err = client.GetInstanceTransfer(context.Background(), instance.ID)
	if err != nil {
		t.Errorf("Error getting instance transfer, expected struct, got error %v", err)
	}
}

func TestInstance_GetMonthlyTransfer(t *testing.T) {
	client, instance, _, teardown, err := setupInstanceWithoutDisks(t, "fixtures/TestInstance_GetMonthlyTransfer", true)
	defer teardown()
	if err != nil {
		t.Error(err)
	}

	currentYear, currentMonth := time.Now().Year(), int(time.Now().Month())

	_, err = client.GetInstanceTransferMonthly(context.Background(), instance.ID, currentYear, currentMonth)
	if err != nil {
		t.Errorf("Error getting monthly instance transfer, expected struct, got error %v", err)
	}

	_, err = client.GetInstanceTransferMonthlyV2(context.Background(), instance.ID, currentYear, currentMonth)
	if err != nil {
		t.Errorf("Error getting monthly instance transfer, expected struct, got error %v", err)
	}
}

func TestInstance_ResetPassword(t *testing.T) {
	client, instance, teardown, err := setupInstance(
		t,
		"fixtures/TestInstance_ResetPassword", true,
		func(client *linodego.Client, options *linodego.InstanceCreateOptions) {
			boot := false
			options.Type = "g6-nanode-1"
			options.Booted = &boot
			options.RootPass = randPassword()
		},
	)

	defer teardown()
	if err != nil {
		t.Error(err)
	}

	instance, err = client.WaitForInstanceStatus(
		context.Background(),
		instance.ID,
		linodego.InstanceOffline,
		180,
	)
	if err != nil {
		t.Errorf("Error waiting for instance readiness for password reset: %s", err.Error())
	}

	err = client.ResetInstancePassword(
		context.Background(),
		instance.ID,
		linodego.InstancePasswordResetOptions{
			RootPass: randPassword(),
		},
	)
	if err != nil {
		t.Errorf("failed to reset instance password for instance with id %d: %v", instance.ID, err.Error())
	}
}

func TestInstance_Resize(t *testing.T) {
	client, instance, teardown, err := setupInstance(
		t,
		"fixtures/TestInstance_Resize", true,
		func(client *linodego.Client, options *linodego.InstanceCreateOptions) {
			boot := true
			options.Type = "g6-nanode-1"
			options.Booted = &boot
		},
	)

	defer teardown()
	if err != nil {
		t.Error(err)
	}

	instance, err = client.WaitForInstanceStatus(
		context.Background(),
		instance.ID,
		linodego.InstanceRunning,
		180,
	)
	if err != nil {
		t.Errorf("Error waiting for instance readiness for resize: %s", err.Error())
	}

	err = client.ResizeInstance(
		context.Background(),
		instance.ID,
		linodego.InstanceResizeOptions{
			Type:          "g6-standard-1",
			MigrationType: "warm",
		},
	)
	if err != nil {
		t.Errorf("failed to resize instance %d: %v", instance.ID, err.Error())
	}
}

func TestInstance_Migrate(t *testing.T) {
	client, instance, teardown, err := setupInstance(
		t,
		"fixtures/TestInstance_Migrate", true,
		func(client *linodego.Client, options *linodego.InstanceCreateOptions) {
			boot := true
			options.Type = "g6-nanode-1"
			options.Booted = &boot
		},
	)

	defer teardown()
	if err != nil {
		t.Error(err)
	}

	instance, err = client.WaitForInstanceStatus(
		context.Background(),
		instance.ID,
		linodego.InstanceRunning,
		180,
	)
	if err != nil {
		t.Errorf("Error waiting for instance readiness for migration: %s", err.Error())
	}

	upgrade := false

	err = client.MigrateInstance(
		context.Background(),
		instance.ID,
		linodego.InstanceMigrateOptions{
			Type:    "cold",
			Region:  "us-west",
			Upgrade: &upgrade,
		},
	)
	if err != nil {
		t.Errorf("failed to migrate instance %d: %v", instance.ID, err.Error())
	}
}

func TestInstance_MigrateToPG(t *testing.T) {
	client, clientTeardown := createTestClient(t, "fixtures/TestInstance_MigrateToPG")

	defer func() {
		clientTeardown()
	}()

	regions := getRegionsWithCaps(t, client, []string{"Placement Group"})

	pgOutboundCreateOpts := linodego.PlacementGroupCreateOptions{
		Label:                "linodego-test-" + getUniqueText(),
		Region:               regions[0],
		PlacementGroupType:   linodego.PlacementGroupTypeAntiAffinityLocal,
		PlacementGroupPolicy: linodego.PlacementGroupPolicyFlexible,
	}

	pgOutbound, err := client.CreatePlacementGroup(context.Background(), pgOutboundCreateOpts)
	if err != nil {
		t.Fatalf("failed to create placement group: %s", err)
	}

	instanceCreateOpts := linodego.InstanceCreateOptions{
		Label:    "go-test-ins-" + randLabel(),
		RootPass: randPassword(),
		Region:   regions[0],
		Type:     "g6-nanode-1",
		Image:    "linode/debian12",
		Booted:   linodego.Pointer(true),
		PlacementGroup: &linodego.InstanceCreatePlacementGroupOptions{
			ID: pgOutbound.ID,
		},
	}

	instance, err := client.CreateInstance(context.Background(), instanceCreateOpts)
	if err != nil {
		t.Fatalf("failed to create instance: %s", err)
	}

	instance, err = client.WaitForInstanceStatus(
		context.Background(),
		instance.ID,
		linodego.InstanceRunning,
		180,
	)
	if err != nil {
		t.Errorf("Error waiting for instance readiness for migration: %s", err.Error())
	}

	pgInboundCreateOpts := linodego.PlacementGroupCreateOptions{
		Label:                "linodego-test-" + getUniqueText(),
		Region:               regions[1],
		PlacementGroupType:   linodego.PlacementGroupTypeAntiAffinityLocal,
		PlacementGroupPolicy: linodego.PlacementGroupPolicyFlexible,
	}

	pgInbound, err := client.CreatePlacementGroup(context.Background(), pgInboundCreateOpts)
	if err != nil {
		t.Fatalf("failed to create placement group: %s", err)
	}

	upgrade := false

	err = client.MigrateInstance(
		context.Background(),
		instance.ID,
		linodego.InstanceMigrateOptions{
			Type:           "cold",
			Region:         regions[1],
			Upgrade:        &upgrade,
			PlacementGroup: &linodego.InstanceCreatePlacementGroupOptions{ID: pgInbound.ID},
		},
	)
	if err != nil {
		t.Errorf("failed to migrate instance %d: %v", instance.ID, err.Error())
	}

	pgInboundRefreshed, err := client.GetPlacementGroup(context.Background(), pgInbound.ID)
	if err != nil {
		t.Fatalf("failed to get placement group: %s", err)
	}

	pgOutboundRefreshed, err := client.GetPlacementGroup(context.Background(), pgOutbound.ID)
	if err != nil {
		t.Fatalf("failed to get placement group: %s", err)
	}

	require.Equal(t, pgInboundRefreshed.ID, pgInbound.ID)
	require.Equal(t, pgInboundRefreshed.Migrations.Inbound[0].LinodeID, instance.ID)
	require.Equal(t, pgOutboundRefreshed.Migrations.Outbound[0].LinodeID, instance.ID)

	if err := client.DeleteInstance(context.Background(), instance.ID); err != nil {
		t.Errorf("failed to delete instance: %s", err)
	}

	if err := client.DeletePlacementGroup(context.Background(), pgInboundRefreshed.ID); err != nil {
		t.Errorf("failed to delete placement group: %s", err)
	}

	if err := client.DeletePlacementGroup(context.Background(), pgOutboundRefreshed.ID); err != nil {
		t.Errorf("failed to delete placement group: %s", err)
	}
}

func TestInstance_Disks_List(t *testing.T) {
	client, instance, teardown, err := setupInstance(t, "fixtures/TestInstance_Disks_List", true)
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

func TestInstance_Disks_List_WithEncryption(t *testing.T) {
	client, instance, teardown, err := setupInstance(t, "fixtures/TestInstance_Disks_List_WithEncryption", true, func(c *linodego.Client, ico *linodego.InstanceCreateOptions) {
		ico.Region = getRegionsWithCaps(t, c, []string{"Disk Encryption"})[0]
	})
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

	// Disk Encryption should be enabled by default if not otherwise specified
	for _, disk := range disks {
		if disk.DiskEncryption != linodego.InstanceDiskEncryptionEnabled {
			t.Fatalf("expected disk encryption status: %s, got :%s", linodego.InstanceDiskEncryptionEnabled, disk.DiskEncryption)
		}
	}
}

func TestInstance_Disk_Resize(t *testing.T) {
	client, instance, _, teardown, err := setupInstanceWithoutDisks(t, "fixtures/TestInstance_Disk_Resize", true)
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
	client, instance1, teardown1, err := setupInstance(t, "fixtures/TestInstance_Disk_ListMultiple_Primary", true)
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

	client, instance2, _, teardown2, err := setupInstanceWithoutDisks(t, "fixtures/TestInstance_Disk_ListMultiple_Secondary", true)
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
		RootPass: randPassword(),
		Size:     2000,
	})
	if err != nil {
		t.Errorf("Error creating disk from private image: %s", err)
	}

	_, err = client.CreateInstanceDisk(context.Background(), instance2.ID, linodego.InstanceDiskCreateOptions{
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

func TestInstance_Disk_Clone(t *testing.T) {
	client, instance, _, teardown, err := setupInstanceWithoutDisks(t, "fixtures/TestInstance_Disk_Clone", true)
	defer teardown()
	if err != nil {
		t.Error(err)
	}

	instance, err = client.WaitForInstanceStatus(context.Background(), instance.ID, linodego.InstanceOffline, 180)
	if err != nil {
		t.Errorf("Error waiting for instance readiness for disk clone: %s", err)
	}

	disk, err := client.CreateInstanceDisk(context.Background(), instance.ID, linodego.InstanceDiskCreateOptions{
		Label:      "go-disk-test-" + randLabel(),
		Filesystem: "ext4",
		Image:      "linode/debian10",
		RootPass:   randPassword(),
		Size:       2000,
	})
	if err != nil {
		t.Errorf("Error creating disk for disk clone: %s", err)
	}

	instance, err = client.WaitForInstanceStatus(context.Background(), instance.ID, linodego.InstanceOffline, 180)
	if err != nil {
		t.Errorf("Error waiting for instance readiness after creating disk for disk clone: %s", err)
	}
	disk, err = client.WaitForInstanceDiskStatus(context.Background(), instance.ID, disk.ID, linodego.DiskReady, 180)
	if err != nil {
		t.Errorf("Error waiting for disk readiness for disk clone: %s", err)
	}

	opts := linodego.InstanceDiskCloneOptions{}

	_, err = client.CloneInstanceDisk(context.Background(), instance.ID, disk.ID, opts)
	if err != nil {
		t.Errorf("Error cloning instance disk: %s", err)
	}
}

func TestInstance_Disk_ResetPassword(t *testing.T) {
	client, instance, _, teardown, err := setupInstanceWithoutDisks(t, "fixtures/TestInstance_Disk_ResetPassword", true)
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
		Image:      "linode/debian10",
		RootPass:   randPassword(),
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

func TestInstance_NodeBalancers_List(t *testing.T) {
	client, _, _, node, teardown, err := setupNodeBalancerNode(t, "fixtures/TestInstance_NodeBalancers_List")
	defer teardown()
	if err != nil {
		t.Error(err)
	}

	privateIP := strings.Split(node.Address, ":")[0]

	instanceIPs, err := client.ListIPAddresses(context.Background(), nil)
	if err != nil {
		t.Error(err)
	}

	linodeID := 0

	for _, instanceIP := range instanceIPs {
		if instanceIP.Address == privateIP {
			linodeID = instanceIP.LinodeID
			break
		}
	}

	if linodeID == 0 {
		t.Errorf("Could not find instance with node's IP")
	}

	nodebalancers, err := client.ListInstanceNodeBalancers(context.Background(), linodeID, nil)
	if err != nil {
		t.Errorf("Error listing instance nodebalancers, expected struct, got error %v", err)
	}
	if len(nodebalancers) == 0 {
		t.Errorf("Expected an list of instance nodebalancers, but got %v", nodebalancers)
	}
}

func TestInstance_Volumes_List(t *testing.T) {
	client, instance, config, teardown, err := setupInstanceWithoutDisks(t, "fixtures/TestInstance_Volumes_List_Instance", true)
	defer teardown()
	if err != nil {
		t.Error(err)
	}

	volume, teardown, volErr := createVolume(t, client)

	_, err = client.WaitForVolumeStatus(context.Background(), volume.ID, linodego.VolumeActive, 500)
	if err != nil {
		t.Errorf("Error waiting for volume to be active, %s", err)
	}

	defer teardown()
	if volErr != nil {
		t.Error(err)
	}

	configOpts := linodego.InstanceConfigUpdateOptions{
		Label: "go-vol-test" + getUniqueText(),
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

	volumes, err := client.ListInstanceVolumes(context.Background(), instance.ID, nil)
	if err != nil {
		t.Errorf("Error listing instance volumes, expected struct, got error %v", err)
	}
	if len(volumes) == 0 {
		t.Errorf("Expected an list of instance volumes, but got %v", volumes)
	}
}

func TestInstance_CreateUnderFirewall(t *testing.T) {
	client, firewall, firewallTeardown, err := setupFirewall(
		t,
		[]firewallModifier{},
		"fixtures/TestInstance_CreateUnderFirewall",
	)
	defer firewallTeardown()

	if err != nil {
		t.Error(err)
	}
	_, _, teardownInstance, err := createInstanceWithoutDisks(
		t,
		client, true,
		func(_ *linodego.Client, options *linodego.InstanceCreateOptions) {
			options.FirewallID = firewall.ID
		},
	)
	defer teardownInstance()

	if err != nil {
		t.Error(err)
	}
}

func TestInstance_Rebuild(t *testing.T) {
	client, instance, _, teardown, err := setupInstanceWithoutDisks(
		t,
		"fixtures/TestInstance_Rebuild", true,
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
		Image: "linode/alpine3.19",
		Metadata: &linodego.InstanceMetadataOptions{
			UserData: base64.StdEncoding.EncodeToString([]byte("cool")),
		},
		RootPass: randPassword(),
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

func TestInstance_RebuildWithEncryption(t *testing.T) {
	client, instance, _, teardown, err := setupInstanceWithoutDisks(
		t,
		"fixtures/TestInstance_RebuildWithEncryption",
		true,
		func(client *linodego.Client, options *linodego.InstanceCreateOptions) {
			options.Region = getRegionsWithCaps(t, client, []string{"Disk Encryption"})[0]
			options.DiskEncryption = linodego.InstanceDiskEncryptionEnabled
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
		Image:          "linode/alpine3.19",
		RootPass:       randPassword(),
		Type:           "g6-standard-2",
		DiskEncryption: linodego.InstanceDiskEncryptionDisabled,
	}
	instance, err = client.RebuildInstance(context.Background(), instance.ID, rebuildOpts)
	if err != nil {
		t.Fatal(err)
	}

	if instance.DiskEncryption != linodego.InstanceDiskEncryptionDisabled {
		t.Fatalf("expected instance.DiskEncryption to be: %s, got: %s", linodego.InstanceDiskEncryptionDisabled, linodego.InstanceDiskEncryptionEnabled)
	}
}

func TestInstance_Clone(t *testing.T) {
	var targetRegion string

	client, instance, teardownOriginalLinode, err := setupInstance(
		t, "fixtures/TestInstance_Clone", true,
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
	_, inst, _, teardown, err := setupInstanceWithoutDisks(t, "fixtures/TestInstance_withMetadata", true,
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

func TestInstance_DiskEncryption(t *testing.T) {
	_, inst, teardown, err := setupInstance(t, "fixtures/TestInstance_DiskEncryption", true, func(c *linodego.Client, ico *linodego.InstanceCreateOptions) {
		ico.DiskEncryption = linodego.InstanceDiskEncryptionEnabled
		ico.Region = "us-east"
	})
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(teardown)

	if inst.DiskEncryption != linodego.InstanceDiskEncryptionEnabled {
		t.Fatalf("expected instance to have disk encryption enabled, got: %s, want: %s", inst.DiskEncryption, linodego.InstanceDiskEncryptionEnabled)
	}
}

func TestInstance_withBlockStorageEncryption(t *testing.T) {
	client, clientTeardown := createTestClient(t, "fixtures/TestInstance_withBlockStorageEncryption")

	inst, err := createInstance(t, client, true, func(client *linodego.Client, options *linodego.InstanceCreateOptions) {
		options.Region = getRegionsWithCaps(t, client, []string{"Linodes", "Block Storage Encryption"})[0]
		options.Label = "go-inst-test-create-bde"
	})
	require.NoError(t, err)

	defer func() {
		client.DeleteInstance(context.Background(), inst.ID)
		clientTeardown()
	}()

	// Filtering is not currently supported on capabilities
	require.True(t, slices.Contains(inst.Capabilities, "Block Storage Encryption"))
}

func TestInstance_withVPU(t *testing.T) {
	client, clientTeardown := createTestClient(t, "fixtures/TestInstance_withVPU")

	inst, err := createInstance(t, client, true, func(client *linodego.Client, options *linodego.InstanceCreateOptions) {
		options.Region = "us-lax"
		options.Type = "g1-accelerated-netint-vpu-t1u1-s"
		options.Label = "go-inst-test-create-vpu"
	})
	require.NoError(t, err)

	defer func() {
		client.DeleteInstance(context.Background(), inst.ID)
		clientTeardown()
	}()

	require.NotNil(t, inst.Specs.AcceleratedDevices)
}

func TestInstance_withPG(t *testing.T) {
	client, clientTeardown := createTestClient(t, "fixtures/TestInstance_withPG")

	pg, pgTeardown, err := createPlacementGroup(t, client)
	require.NoError(t, err)

	// Create an instance to assign to the PG
	inst, err := createInstance(t, client, true, func(client *linodego.Client, options *linodego.InstanceCreateOptions) {
		options.Region = pg.Region
		options.PlacementGroup = &linodego.InstanceCreatePlacementGroupOptions{
			ID: pg.ID,
		}
	})
	require.NoError(t, err)

	defer func() {
		client.DeleteInstance(context.Background(), inst.ID)
		pgTeardown()
		clientTeardown()
	}()

	require.NotNil(t, inst.PlacementGroup)
	require.Equal(t, inst.PlacementGroup.ID, pg.ID)
	require.Equal(t, inst.PlacementGroup.Label, pg.Label)
	require.Equal(t, inst.PlacementGroup.PlacementGroupType, pg.PlacementGroupType)
	require.Equal(t, inst.PlacementGroup.PlacementGroupPolicy, pg.PlacementGroupPolicy)
}

func createInstance(t *testing.T, client *linodego.Client, enableCloudFirewall bool, modifiers ...instanceModifier) (*linodego.Instance, error) {
	if t != nil {
		t.Helper()
	}

	createOpts := linodego.InstanceCreateOptions{
		Label:    "go-test-ins-" + randLabel(),
		RootPass: randPassword(),
		Region:   getRegionsWithCaps(t, client, []string{"linodes"})[0],
		Type:     "g6-nanode-1",
		Image:    "linode/debian12",
		Booted:   linodego.Pointer(false),
	}

	if enableCloudFirewall {
		createOpts.FirewallID = firewallID
	}

	for _, modifier := range modifiers {
		modifier(client, &createOpts)
	}
	return client.CreateInstance(context.Background(), createOpts)
}

func setupInstance(t *testing.T, fixturesYaml string, EnableCloudFirewall bool, modifiers ...instanceModifier) (*linodego.Client, *linodego.Instance, func(), error) {
	if t != nil {
		t.Helper()
	}
	client, fixtureTeardown := createTestClient(t, fixturesYaml)

	instance, err := createInstance(t, client, EnableCloudFirewall, modifiers...)
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

func createInstanceWithoutDisks(
	t *testing.T,
	client *linodego.Client,
	enableCloudFirewall bool,
	modifiers ...instanceModifier,
) (*linodego.Instance, *linodego.InstanceConfig, func(), error) {
	t.Helper()

	createOpts := linodego.InstanceCreateOptions{
		Label:  "go-test-ins-wo-disk-" + randLabel(),
		Region: getRegionsWithCaps(t, client, []string{"linodes"})[0],
		Type:   "g6-nanode-1",
		Booted: linodego.Pointer(false),
	}

	if enableCloudFirewall {
		createOpts.FirewallID = GetFirewallID()
	}

	for _, modifier := range modifiers {
		modifier(client, &createOpts)
	}

	instance, err := client.CreateInstance(context.Background(), createOpts)
	if err != nil {
		t.Errorf("Error creating test Instance: %s", err)
		return nil, nil, func() {}, err
	}
	configOpts := linodego.InstanceConfigCreateOptions{
		Label: "go-test-conf-" + randLabel(),
	}
	config, err := client.CreateInstanceConfig(context.Background(), instance.ID, configOpts)
	if err != nil {
		t.Errorf("Error creating config: %s", err)
		return nil, nil, func() {}, err
	}

	teardown := func() {
		if terr := client.DeleteInstance(context.Background(), instance.ID); terr != nil {
			t.Errorf("Error deleting test Instance: %s", terr)
		}
	}
	return instance, config, teardown, err
}

func setupInstanceWithoutDisks(t *testing.T, fixturesYaml string, enableCloudFirewall bool, modifiers ...instanceModifier) (*linodego.Client, *linodego.Instance, *linodego.InstanceConfig, func(), error) {
	t.Helper()
	client, fixtureTeardown := createTestClient(t, fixturesYaml)
	instance, config, instanceTeardown, err := createInstanceWithoutDisks(t, client, enableCloudFirewall, modifiers...)

	teardown := func() {
		instanceTeardown()
		fixtureTeardown()
	}
	return client, instance, config, teardown, err
}
