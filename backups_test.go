package linodego_test

import (
	"context"
	"testing"
)

func TestListInstanceBackups(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping test in short mode.")
	}

	client, teardown := createTestClient(t, "fixtures/TestListInstanceBackups")
	defer teardown()

	backups, err := client.GetInstanceBackups(context.TODO(), TestInstanceID)
	if err != nil {
		t.Errorf("Error listing backups, expected struct, got error %v", err)
	}
	if backups.Automatic == nil || len(backups.Automatic) > 0 {
		t.Errorf("Expected to find no automatic backups, but got %v", backups.Automatic)
	}
	if backups.Snapshot.Current == nil {
		t.Errorf("Expected current snapshot, but got %v", backups.Snapshot.Current)
	}
}
