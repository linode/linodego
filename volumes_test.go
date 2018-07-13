package linodego_test

import (
	"testing"

	"github.com/chiefy/linodego"
)

func TestCreateVolume(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping test in short mode.")
	}
	client, teardown := createTestClient(t, "fixtures/TestCreateVolume")
	defer teardown()

	createOpts := linodego.VolumeCreateOptions{
		Label:  "linodego-test-volume",
		Region: "us-west",
	}
	volume, err := client.CreateVolume(createOpts)
	if err != nil {
		t.Errorf("Error listing volumes, expected struct, got error %v", err)
	}
	if volume.ID == 0 {
		t.Errorf("Expected a volumes id, but got 0")
	}

	if err := client.DeleteVolume(volume.ID); err != nil {
		t.Errorf("Expected to delete a volume, but got %v", err)
	}
}

func TestListVolumes(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping test in short mode.")
	}
	client, teardown := createTestClient(t, "fixtures/TestListVolumes")
	defer teardown()

	volumes, err := client.ListVolumes(nil)
	if err != nil {
		t.Errorf("Error listing volumes, expected struct, got error %v", err)
	}
	if len(volumes) == 0 {
		t.Errorf("Expected a list of volumes, but got %v", volumes)
	}
}

func TestGetVolume(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping test in short mode.")
	}
	client, teardown := createTestClient(t, "fixtures/TestGetVolume")
	defer teardown()

	_, err := client.GetVolume(TestVolumeID)
	if err != nil {
		t.Errorf("Error getting volume %d, expected *LinodeVolume, got error %v", TestVolumeID, err)
	}
}
