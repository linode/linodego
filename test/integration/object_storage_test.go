package integration

import (
	"context"
	"testing"
)

func TestGetObjectStorageTransfer(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestObjectStorage_transfer")
	defer teardown()

	_, err := client.GetObjectStorageTransfer(context.Background())
	if err != nil {
		t.Errorf("unable to get object storage transfer data : %s", err)
	}
}

func TestCancelObjectStorage(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestObjectStorage_transfer")
	defer teardown()

	err := client.CancelObjectStorage(context.Background())
	if err != nil {
		t.Errorf("failed to cancel object storage : %s", err)
	}
}
