package golinode

import (
	"testing"
)

func TestListStackscripts(t *testing.T) {
	client, err := createTestClient(debugAPI)
	if err != nil {
		t.Errorf("Error creating test client %v", err)
	}
	regions, err := client.ListRegions()
	if err != nil {
		t.Errorf("Error listing regions, expected struct - error %v", err)
	}
	if len(regions) == 0 {
		t.Errorf("Expected a list of regions - %v", regions)
	}
}
