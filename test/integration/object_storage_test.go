package integration

import (
	"context"
	"testing"

	"github.com/jarcoal/httpmock"
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
	client := createMockClient(t)

	desiredResponse := make(map[string]interface{})

	httpmock.RegisterRegexpResponder("POST", mockRequestURL(t, "/object-storage/cancel"),
		httpmock.NewJsonResponderOrPanic(200, &desiredResponse))

	err := client.CancelObjectStorage(context.Background())
	if err != nil {
		t.Fatal(err)
	}
}
