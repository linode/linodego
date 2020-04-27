package integration

import (
	"context"
	"testing"
)

func TestClientAliases(t *testing.T) {
	client, _ := createTestClient(t, "")

	if client.Images == nil {
		t.Error("Expected alias for Images to return a *Resource")
	}
	if client.Instances == nil {
		t.Error("Expected alias for Instances to return a *Resource")
	}
	if client.InstanceSnapshots == nil {
		t.Error("Expected alias for Backups to return a *Resource")
	}
	if client.StackScripts == nil {
		t.Error("Expected alias for StackScripts to return a *Resource")
	}
	if client.Regions == nil {
		t.Error("Expected alias for Regions to return a *Resource")
	}
	if client.Volumes == nil {
		t.Error("Expected alias for Volumes to return a *Resource")
	}
}

func TestClient_APIResponseBadGateway(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestClient_APIResponseBadGateway")
	defer teardown()

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Expected Client to handle 502 from API Server")
		}
	}()

	_, err := client.ListImages(context.Background(), nil)

	if err == nil {
		t.Errorf("Error should be thrown on 502 Response from API")
	}

	expectedErr := "[502] Bad Gateway"
	if err.Error() != expectedErr {
		t.Errorf(`expected error to be "%s" but got "%s"`, expectedErr, err)
	}
}
