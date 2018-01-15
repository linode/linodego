package golinode

import (
	"testing"
)

func TestListImages(t *testing.T) {
	client, err := createTestClient(debugAPI)
	if err != nil {
		t.Errorf("Error creating test client %v", err)
	}
	i, err := client.ListImages()
	if err != nil {
		t.Errorf("Error listing images, expected struct, got error %v", err)
	}
	if len(i) == 0 {
		t.Errorf("Expected a list of images, but got none %v", i)
	}
}
