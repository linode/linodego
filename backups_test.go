package linodego

import (
	"testing"
)

func TestListInstanceBackups(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping test in short mode.")
	}

	client, err := createTestClient(debugAPI)
	if err != nil {
		t.Errorf("Error creating test client %v", err)
	}
	backups, err := client.GetInstanceBackups(TestInstanceID)
	if err != nil {
		t.Errorf("Error listing backups, expected struct, got error %v", err)
	}
	if backups.Automatic != nil && len(backups.Automatic) != 1 {
		t.Errorf("Expected an empty list of automatic backups, but got %v", backups.Automatic)
	}
	if backups.Snapshot.Current != nil {
		t.Errorf("Expected empty current snapshot, but got %v", backups.Snapshot.Current)
	}
}
