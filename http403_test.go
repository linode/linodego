package linodego_test

import (
	"context"
	"testing"
)

func TestListTypes_429(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestListTypes_429")
	defer teardown()

	_, err := client.ListTypes(context.Background(), nil)
	if err == nil {
		t.Errorf("Error listing images, expected error")
	}
}

func TestListTypes_429recovered(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestListTypes_429recovered")
	defer teardown()

	_, err := client.ListTypes(context.Background(), nil)
	if err != nil {
		t.Errorf("Error listing images, expected 429 to recover")
	}
}
