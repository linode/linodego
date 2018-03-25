package golinode

import (
	"fmt"
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
	linode, err := client.GetInstance(6809519)
	if err != nil {
		t.Errorf("Error getting instance 6809519, expected *LinodeInstance, got error %v", err)
	}
	fmt.Printf("%#v \n", linode)
}
