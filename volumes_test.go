package golinode

import (
	"fmt"
	"testing"
)

func TestListVolumes(t *testing.T) {
	client, err := createTestClient(debugAPI)
	if err != nil {
		t.Errorf("Error creating test client %v", err)
	}
	volumes, err := client.ListVolumes()
	if err != nil {
		t.Errorf("Error listing instances, expected struct, got error %v", err)
	}
	if len(volumes) != 1 {
		t.Errorf("Expected a list of instances, but got %v", volumes)
	}
}

func TestGetVolume(t *testing.T) {
	client, err := createTestClient(debugAPI)
	if err != nil {
		t.Errorf("Error creating test client %v", err)
	}
	linode, err := client.GetVolume(4090913)
	if err != nil {
		t.Errorf("Error getting instance 1234, expected *LinodeVolume, got error %v", err)
	}
	fmt.Printf("%#v \n", linode)
}
