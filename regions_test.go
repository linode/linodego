package linodego

import (
	"testing"
)

func TestListRegions(t *testing.T) {
	client, err := createTestClient(debugAPI)
	if err != nil {
		t.Errorf("Error creating test client %v", err)
	}
	regions, err := client.ListRegions(nil)
	if err != nil {
		t.Errorf("Error listing regions, expected struct - error %v", err)
	}
	if len(regions) == 0 {
		t.Errorf("Expected a list of regions - %v", regions)
	}
}
