package linodego_test

import (
	"context"
	"testing"

	"github.com/chiefy/linodego"
)

func TestGetType_missing(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestGetType_missing")
	defer teardown()

	i, err := client.GetType(context.TODO(), "does-not-exist")
	if err == nil {
		t.Errorf("should have received an error requesting a missing image, got %v", i)
	}
	e, ok := err.(*linodego.Error)
	if !ok {
		t.Errorf("should have received an Error requesting a missing image, got %v", e)
	}

	if e.Code != 404 {
		t.Errorf("should have received a 404 Code requesting a missing image, got %v", e.Code)
	}
}

func TestGetType_found(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestGetType_found")
	defer teardown()

	i, err := client.GetType(context.TODO(), "g6-standard-1")
	if err != nil {
		t.Errorf("Error getting image, expected struct, got %v and error %v", i, err)
	}
	if i.ID != "g6-standard-1" {
		t.Errorf("Expected a specific image, but got a different one %v", i)
	}
}
func TestListTypes(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestListTypes")
	defer teardown()

	i, err := client.ListTypes(context.TODO(), nil)
	if err != nil {
		t.Errorf("Error listing images, expected struct, got error %v", err)
	}
	if len(i) == 0 {
		t.Errorf("Expected a list of images, but got none %v", i)
	}
}
