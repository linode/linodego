package golinode

import (
	"testing"
)

func TestListInstances(t *testing.T) {
	client, err := createTestClient(debugAPI)
	if err != nil {
		t.Errorf("Error creating test client %v", err)
	}
	linodes, err := client.ListInstances()
	if err != nil {
		t.Errorf("Error listing instances, expected struct, got error %v", err)
	}
	if len(linodes) != 1 {
		t.Errorf("Expected a list of instances, but got %v", linodes)
	}
}

func TestGetInstance(t *testing.T) {
	client, err := createTestClient(debugAPI)
	if err != nil {
		t.Errorf("Error creating test client %v", err)
	}
	instance, err := client.GetInstance(6809519)
	if err != nil {
		t.Errorf("Error getting instance 6809519, expected *LinodeInstance, got error %v", err)
	}
	if instance.Specs.Disk <= 0 {
		t.Errorf("Error in instance 6809519 spec for disk size, %v", instance.Specs)
	}
}

func TestListInstanceDisks(t *testing.T) {
	client, err := createTestClient(debugAPI)
	if err != nil {
		t.Errorf("Error creating test client %v", err)
	}
	disks, err := client.ListInstanceDisks(6809519)
	if err != nil {
		t.Errorf("Error listing instance disks, expected struct, got error %v", err)
	}
	if len(disks) != 1 {
		t.Errorf("Expected a list of instance disks, but got %v", disks)
	}
}

func TestListInstanceConfigs(t *testing.T) {
	client, err := createTestClient(debugAPI)
	if err != nil {
		t.Errorf("Error creating test client %v", err)
	}
	configs, err := client.ListInstanceConfigs(6809519)
	if err != nil {
		t.Errorf("Error listing instance configs, expected struct, got error %v", err)
	}
	if len(configs) != 1 {
		t.Errorf("Expected a list of instance configs, but got %v", configs)
	}
}

func TestListInstanceVolumes(t *testing.T) {
	client, err := createTestClient(debugAPI)
	if err != nil {
		t.Errorf("Error creating test client %v", err)
	}
	volumes, err := client.ListInstanceVolumes(6809519)
	if err != nil {
		t.Errorf("Error listing instance volumes, expected struct, got error %v", err)
	}
	if len(volumes) != 1 {
		t.Errorf("Expected a list of instance volumes, but got %v", volumes)
	}
}
