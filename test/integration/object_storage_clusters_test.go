package integration

import (
	"context"
	"testing"
)

func TestListObjectStorageClusters(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestListObjectStorageClusters")
	defer teardown()

	objectStorageClusters, err := client.ListObjectStorageClusters(context.Background(), nil)
	if err != nil {
		t.Errorf("Error listing objectStorageClusters, expected struct - error %v", err)
	}
	if len(objectStorageClusters) == 0 {
		t.Errorf("Expected a list of objectStorageClusters - %v", objectStorageClusters)
	}
}
