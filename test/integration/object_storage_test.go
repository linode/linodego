package integration

import (
	"context"
	"testing"
)

func TestObjectStorage_Get_Transfer(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestObjectStorage_transfer")
	defer teardown()

	_, err := client.GetObjectStorageTransfer(context.Background())
	if err != nil {
		t.Errorf("unable to get object storage transfer data : %s", err)
	}
}

func TestObjectStorage_Cancel(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestObjectStorage_cancel")
	defer teardown()

	err := client.CancelObjectStorage(context.Background())
	if err != nil {
		t.Errorf("failed to cancel object storage : %s", err)
	}
}
