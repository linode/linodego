package integration

import (
	"context"
	"testing"

	"github.com/linode/linodego"
)

func TestCreateVolume(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestCreateVolume")
	defer teardown()

	createOpts := linodego.VolumeCreateOptions{
		Label:  "linodego-test-volume",
		Region: "us-west",
	}
	volume, err := client.CreateVolume(context.Background(), createOpts)
	if err != nil {
		t.Errorf("Error listing volumes, expected struct, got error %v", err)
	}
	if volume.ID == 0 {
		t.Errorf("Expected a volumes id, but got 0")
	}

	assertDateSet(t, volume.Created)
	assertDateSet(t, volume.Updated)

	if err := client.DeleteVolume(context.Background(), volume.ID); err != nil {
		t.Errorf("Expected to delete a volume, but got %v", err)
	}
}

func TestRenameVolume(t *testing.T) {
	client, volume, teardown, err := setupVolume(t, "fixtures/TestRenameVolume")
	defer teardown()
	if err != nil {
		t.Error(err)
	}

	volume, err = client.RenameVolume(context.Background(), volume.ID, "test-volume-renamed")
	if err != nil || volume.Label != "test-volume-renamed" {
		t.Errorf("Error renaming volume, %s", err)
	}
}

func TestResizeVolume(t *testing.T) {
	client, volume, teardown, err := setupVolume(t, "fixtures/TestResizeVolume")
	defer teardown()

	if err != nil {
		t.Errorf("Error setting up volume test, %s", err)
	}

	if err := client.ResizeVolume(context.Background(), volume.ID, volume.Size+1); err != nil {
		t.Errorf("Error resizing volume, %s", err)
	}
}

func TestListVolumes(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestListVolumes")
	defer teardown()

	volumes, err := client.ListVolumes(context.Background(), nil)
	if err != nil {
		t.Errorf("Error listing volumes, expected struct, got error %v", err)
	}
	if len(volumes) == 0 {
		t.Errorf("Expected a list of volumes, but got %v", volumes)
	}
}

func TestGetVolume(t *testing.T) {
	client, volume, teardownVolume, errVolume := setupVolume(t, "fixtures/TestGetVolume")
	defer teardownVolume()
	if errVolume != nil {
		t.Error(errVolume)
	}

	_, err := client.GetVolume(context.Background(), volume.ID)
	if err != nil {
		t.Errorf("Error getting volume %d, expected *LinodeVolume, got error %v", volume.ID, err)
	}
	assertDateSet(t, volume.Created)
	assertDateSet(t, volume.Updated)
}

func TestWaitForVolumeLinodeID_nil(t *testing.T) {
	client, volume, teardown, err := setupVolume(t, "fixtures/TestWaitForVolumeLinodeID_nil")
	defer teardown()

	if err != nil {
		t.Errorf("Error setting up volume test, %s", err)
	}
	_, err = client.WaitForVolumeLinodeID(context.Background(), volume.ID, nil, 20)

	if err != nil {
		t.Errorf("Error getting volume %d, expected *LinodeVolume, got error %v", volume.ID, err)
	}
}

func TestWaitForVolumeLinodeIDInstance(t *testing.T) {
	client, instance, teardownInstance, errInstance := setupInstance(t, "fixtures/TestWaitForVolumeLinodeID_linode")
	if errInstance != nil {
		t.Errorf("Error setting up instance for volume test, %s", errInstance)
	}

	defer teardownInstance()

	createConfigOpts := linodego.InstanceConfigCreateOptions{
		Label:   "test-instance-volume",
		Devices: linodego.InstanceConfigDeviceMap{},
	}
	config, errConfig := client.CreateInstanceConfig(context.Background(), instance.ID, createConfigOpts)
	if errConfig != nil {
		t.Errorf("Error setting up instance config for volume test, %s", errConfig)
	}

	client, volume, teardownVolume, errVolume := setupVolume(t, "fixtures/TestWaitForVolumeLinodeID_volume")
	if errVolume != nil {
		t.Errorf("Error setting up volume test, %s", errVolume)
	}
	defer teardownVolume()

	attachOptions := linodego.VolumeAttachOptions{LinodeID: instance.ID, ConfigID: config.ID}
	if volumeAttached, err := client.AttachVolume(context.Background(), volume.ID, &attachOptions); err != nil {
		t.Errorf("Error attaching volume, %s", err)
	} else if volumeAttached.LinodeID == nil {
		t.Errorf("Could not attach test volume to test instance")
	}

	_, errWait := client.WaitForVolumeLinodeID(context.Background(), volume.ID, nil, 20)
	if errWait == nil {
		t.Errorf("Expected to timeout waiting for nil LinodeID on volume %d : %s", volume.ID, errWait)
	}

	client, teardownWait := createTestClient(t, "fixtures/TestWaitForVolumeLinodeID_waiting")
	defer teardownWait()

	_, errWait = client.WaitForVolumeLinodeID(context.Background(), volume.ID, &instance.ID, 20)
	if errWait != nil {
		t.Errorf("Error waiting for volume %d to attach to instance %d: %s", volume.ID, instance.ID, errWait)
	}
}

func setupVolume(t *testing.T, fixturesYaml string) (*linodego.Client, *linodego.Volume, func(), error) {
	t.Helper()
	var fixtureTeardown func()
	client, fixtureTeardown := createTestClient(t, fixturesYaml)

	createOpts := linodego.VolumeCreateOptions{
		Label:  "linodego-test-volume",
		Region: "us-west",
	}
	volume, err := client.CreateVolume(context.Background(), createOpts)
	if err != nil {
		t.Errorf("Error listing volumes, expected struct, got error %v", err)
	}

	teardown := func() {
		if terr := client.DeleteVolume(context.Background(), volume.ID); terr != nil {
			t.Errorf("Expected to delete a volume, but got %v", terr)
		}
		fixtureTeardown()
	}
	return client, volume, teardown, err
}
