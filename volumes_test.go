package linodego_test

import (
	"testing"
)

const TestVolumeID = 7568

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
