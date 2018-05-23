package linodego

import (
	"testing"
)

func TestGetType_missing(t *testing.T) {
	client, err := createTestClient(debugAPI)
	if err != nil {
		t.Errorf("Error creating test client %v", err)
	}

	i, err := client.GetType("does-not-exist")
	if err == nil {
		t.Errorf("should have received an error requesting a missing image, got %v", i)
	}
	e, ok := err.(*Error)
	if !ok {
		t.Errorf("should have received an Error requesting a missing image, got %v", e)
	}

	if e.Code != 404 {
		t.Errorf("should have received a 404 Code requesting a missing image, got %v", e.Code)
	}
}

func TestGetType_found(t *testing.T) {
	client, err := createTestClient(debugAPI)
	if err != nil {
		t.Errorf("Error creating test client %v", err)
	}
	i, err := client.GetType("g6-standard-1")
	if err != nil {
		t.Errorf("Error getting image, expected struct, got %v and error %v", i, err)
	}
	if i.ID != "g6-standard-1" {
		t.Errorf("Expected a specific image, but got a different one %v", i)
	}
}
func TestListTypes(t *testing.T) {
	client, err := createTestClient(debugAPI)
	if err != nil {
		t.Errorf("Error creating test client %v", err)
	}
	i, err := client.ListTypes(nil)
	if err != nil {
		t.Errorf("Error listing images, expected struct, got error %v", err)
	}
	if len(i) == 0 {
		t.Errorf("Expected a list of images, but got none %v", i)
	}
}
