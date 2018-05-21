package golinode

import (
	"testing"
)

const TestVolumeID = 5029

func TestListVolumes(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping test in short mode.")
	}
	client, err := createTestClient(debugAPI)
	if err != nil {
		t.Errorf("Error creating test client %v", err)
	}
	volumes, err := client.ListVolumes(nil)
	if err != nil {
		t.Errorf("Error listing instances, expected struct, got error %v", err)
	}
	if len(volumes) == 0 {
		t.Errorf("Expected a list of instances, but got %v", volumes)
	}
}

func TestGetVolume(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping test in short mode.")
	}
	client, err := createTestClient(debugAPI)
	if err != nil {
		t.Errorf("Error creating test client %v", err)
	}
	_, err = client.GetVolume(TestVolumeID)
	if err != nil {
		t.Errorf("Error getting volume %d, expected *LinodeVolume, got error %v", TestVolumeID, err)
	}
}
