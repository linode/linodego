package integration

import (
	"context"
	"testing"
	"time"

	"github.com/linode/linodego"
)

type volumeModifier func(*linodego.Client, *linodego.VolumeCreateOptions)

func TestVolume_Create_smoke(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestVolume_Create")
	defer teardown()

	createOpts := linodego.VolumeCreateOptions{
		Label:  linodego.Pointer("go-vol-test-create"),
		Region: linodego.Pointer(getRegionsWithCaps(t, client, []string{"Linodes"})[0]),
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

	// volumes deleted too fast tend to stick, adding a few seconds to catch up
	time.Sleep(time.Second * 5)
	if err := client.DeleteVolume(context.Background(), volume.ID); err != nil {
		t.Errorf("Expected to delete a volume, but got %v", err)
	}
}

func TestVolume_Create_withEncryption(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestVolume_Create_withEncryption")
	defer teardown()

	createOpts := linodego.VolumeCreateOptions{
		Label:      linodego.Pointer("go-vol-test-create-encryption"),
		Region:     linodego.Pointer(getRegionsWithCaps(t, client, []string{"Linodes", "Block Storage Encryption"})[0]),
		Encryption: linodego.Pointer("enabled"),
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

	// volumes deleted too fast tend to stick, adding a few seconds to catch up
	time.Sleep(time.Second * 5)
	if err := client.DeleteVolume(context.Background(), volume.ID); err != nil {
		t.Errorf("Expected to delete a volume, but got %v", err)
	}
}

func TestVolume_Resize(t *testing.T) {
	client, volume, teardown, err := setupVolume(t, "fixtures/TestVolume_Resize")
	defer teardown()

	if err != nil {
		t.Errorf("Error setting up volume test, %s", err)
	}

	_, err = client.WaitForVolumeStatus(context.Background(), volume.ID, linodego.VolumeActive, 500)
	if err != nil {
		t.Errorf("Error waiting for volume to be active, %s", err)
	}

	opts := linodego.VolumeResizeOptions{
		Size: volume.Size + 1,
	}
	if err := client.ResizeVolume(context.Background(), volume.ID, opts); err != nil {
		t.Errorf("Error resizing volume, %s", err)
	}
}

func TestVolumes_List_smoke(t *testing.T) {
	client, volume, teardown, err := setupVolume(t, "fixtures/TestVolume_List")
	defer teardown()

	volumes, err := client.ListVolumes(context.Background(), nil)
	if err != nil {
		t.Errorf("Error listing volumes, expected struct, got error %v", err)
	}
	found := false
	for _, v := range volumes {
		if v.ID == volume.ID {
			found = true
		}
	}
	if !found {
		t.Errorf("%d volume not found in list", volume.ID)
	}
}

func TestVolume_Get(t *testing.T) {
	client, volume, teardownVolume, errVolume := setupVolume(t, "fixtures/TestVolume_Get")
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

func TestVolume_Get_withEncryption(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestVolume_Get_withEncryption")
	defer teardown()

	createOpts := linodego.VolumeCreateOptions{
		Label:      linodego.Pointer("go-vol-test-get-encryption"),
		Region:     linodego.Pointer(getRegionsWithCaps(t, client, []string{"Linodes", "Block Storage Encryption"})[0]),
		Encryption: linodego.Pointer("enabled"),
	}
	volume, err := client.CreateVolume(context.Background(), createOpts)
	if err != nil {
		t.Errorf("Error listing volumes, expected struct, got error %v", err)
	}
	if volume.ID == 0 {
		t.Errorf("Expected a volumes id, but got 0")
	}

	returnedVolume, err := client.GetVolume(context.Background(), volume.ID)
	if err != nil {
		t.Errorf("Error getting volume %d, expected *LinodeVolume, got error %v", volume.ID, err)
	}
	if returnedVolume.Encryption != "enabled" {
		t.Errorf("Expected volume encryption to be enabled, but got %v", returnedVolume.Encryption)
	}

	volumes, err := client.ListVolumes(context.Background(), nil)
	if err != nil {
		t.Errorf("Error listing volumes, expected struct, got error %v", err)
	}
	found := false
	for _, v := range volumes {
		if v.ID == volume.ID {
			found = true
			if v.Encryption != "enabled" {
				t.Errorf("Expected volume encryption to be enabled, but got %v", v.Encryption)
			}
		}
	}
	if !found {
		t.Errorf("%d volume not found in list", volume.ID)
	}

	// volumes deleted too fast tend to stick, adding a few seconds to catch up
	time.Sleep(time.Second * 5)
	if err := client.DeleteVolume(context.Background(), volume.ID); err != nil {
		t.Errorf("Expected to delete a volume, but got %v", err)
	}
}

func TestVolume_WaitForLinodeID_nil(t *testing.T) {
	client, volume, teardown, err := setupVolume(t, "fixtures/TestVolume_WaitForLinodeID_nil")
	defer teardown()

	if err != nil {
		t.Errorf("Error setting up volume test, %s", err)
	}
	_, err = client.WaitForVolumeLinodeID(context.Background(), volume.ID, nil, 20)
	if err != nil {
		t.Errorf("Error getting volume %d, expected *LinodeVolume, got error %v", volume.ID, err)
	}
}

func TestVolume_WaitForLinodeID(t *testing.T) {
	client, instance, teardownInstance, errInstance := setupInstance(t, "fixtures/TestVolume_WaitForLinodeID_linode", true)
	if errInstance != nil {
		t.Errorf("Error setting up instance for volume test, %s", errInstance)
	}

	defer teardownInstance()

	createConfigOpts := linodego.InstanceConfigCreateOptions{
		Label:   linodego.Pointer("go-config-test-wait"),
		Devices: linodego.InstanceConfigDeviceMap{},
	}
	config, errConfig := client.CreateInstanceConfig(context.Background(), instance.ID, createConfigOpts)
	if errConfig != nil {
		t.Errorf("Error setting up instance config for volume test, %s", errConfig)
	}

	client, volume, teardownVolume, errVolume := setupVolume(t, "fixtures/TestVolume_WaitForLinodeID_volume")
	if errVolume != nil {
		t.Errorf("Error setting up volume test, %s", errVolume)
	}
	defer teardownVolume()

	attachOptions := linodego.VolumeAttachOptions{LinodeID: instance.ID, ConfigID: &config.ID}
	if volumeAttached, err := client.AttachVolume(context.Background(), volume.ID, &attachOptions); err != nil {
		t.Errorf("Error attaching volume, %s", err)
	} else if volumeAttached.LinodeID == nil {
		t.Errorf("Could not attach test volume to test instance")
	}

	_, errWait := client.WaitForVolumeLinodeID(context.Background(), volume.ID, nil, 20)
	if errWait == nil {
		t.Errorf("Expected to timeout waiting for nil LinodeID on volume %d : %s", volume.ID, errWait)
	}

	client, teardownWait := createTestClient(t, "fixtures/TestVolume_WaitForLinodeID_waiting")
	defer teardownWait()

	v, errWait := client.WaitForVolumeLinodeID(context.Background(), volume.ID, &instance.ID, 20)
	if errWait != nil {
		t.Errorf("Error waiting for volume %d to attach to instance %d: %s", volume.ID, instance.ID, errWait)
	}
	if v.LinodeLabel != instance.Label {
		t.Error("Error linode label mismatched")
	}
}

func TestVolume_Update(t *testing.T) {
	client, volume, teardown, err := setupVolume(t, "fixtures/TestVolume_Update")
	if err != nil {
		t.Errorf("Error setting up volume test, %s", err)
	}
	defer teardown()
	updatedLabel := volume.Label + "-updated"
	updateOpts := linodego.VolumeUpdateOptions{
		Label: &updatedLabel,
	}
	volume, err = client.UpdateVolume(context.Background(), volume.ID, updateOpts)
	if err != nil {
		t.Errorf("Error updating volume, expected struct, got error %v", err)
	}
	if volume.ID == 0 {
		t.Errorf("Expected a volumes id, but got 0")
	}
	if volume.Label != updatedLabel {
		t.Errorf("Expected volume label to be equal to updated volume label")
	}
	assertDateSet(t, volume.Created)
	assertDateSet(t, volume.Updated)
}

func setupVolume(t *testing.T, fixturesYaml string) (*linodego.Client, *linodego.Volume, func(), error) {
	t.Helper()
	var fixtureTeardown func()
	client, fixtureTeardown := createTestClient(t, fixturesYaml)
	createOpts := linodego.VolumeCreateOptions{
		Label:  linodego.Pointer("go-vol-test-def"),
		Region: linodego.Pointer(getRegionsWithCaps(t, client, []string{"Linodes"})[0]),
	}
	volume, err := client.CreateVolume(context.Background(), createOpts)
	if err != nil {
		t.Errorf("Error creating volume, got error %v", err)
	}

	teardown := func() {
		// volumes deleted too fast tend to stick, adding a few seconds to catch up
		time.Sleep(time.Second * 5)
		if terr := client.DeleteVolume(context.Background(), volume.ID); terr != nil {
			t.Errorf("Expected to delete a volume, but got %v", terr)
		}
		fixtureTeardown()
	}
	return client, volume, teardown, err
}

func createVolume(
	t *testing.T,
	client *linodego.Client,
	vModifier ...volumeModifier,
) (*linodego.Volume, func(), error) {
	t.Helper()
	createOpts := linodego.VolumeCreateOptions{
		Label:  linodego.Pointer("go-vol-test" + randLabel()),
		Region: linodego.Pointer(getRegionsWithCaps(t, client, []string{"Linodes"})[0]),
	}

	for _, mod := range vModifier {
		mod(client, &createOpts)
	}

	v, err := client.CreateVolume(context.Background(), createOpts)
	if err != nil {
		t.Fatalf("failed to create volume: %s", err)
	}

	teardown := func() {
		if err := client.DeleteVolume(context.Background(), v.ID); err != nil {
			t.Errorf("failed to delete volume: %s", err)
		}
	}
	return v, teardown, err
}
