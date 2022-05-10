package integration

import (
	"context"
	"testing"

	"github.com/linode/linodego"
)

func TestType_GetMissing(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestType_GetMissing")
	defer teardown()

	i, err := client.GetType(context.Background(), "does-not-exist")
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

func TestType_GetFound(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestType_GetFound")
	defer teardown()

	i, err := client.GetType(context.Background(), "g6-standard-1")
	if err != nil {
		t.Errorf("Error getting image, expected struct, got %v and error %v", i, err)
	}
	if i.ID != "g6-standard-1" {
		t.Errorf("Expected a specific image, but got a different one %v", i)
	}
}

func TestTypes_List(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestTypes_List")
	defer teardown()

	i, err := client.ListTypes(context.Background(), nil)
	if err != nil {
		t.Errorf("Error listing images, expected struct, got error %v", err)
	}
	if len(i) == 0 {
		t.Errorf("Expected a list of images, but got none %v", i)
	}
}
