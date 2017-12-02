package golinode

import (
	"testing"
)

func TestListDistributions(t *testing.T) {
	client, err := createTestClient(debugAPI)
	if err != nil {
		t.Errorf("Error creating test client %v", err)
	}
	dists, err := client.ListDistributions()
	if err != nil {
		t.Errorf("Error listing distributions, expected struct, got error %v", err)
	}
	if len(dists) == 0 {
		t.Errorf("Expected a list of distributions, but got none %v", dists)
	}
}
