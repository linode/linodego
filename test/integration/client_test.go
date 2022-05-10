package integration

import (
	"testing"
)

func TestClient_Aliases(t *testing.T) {
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
