package golinode

import (
	"testing"
)

func TestListInstanceBackups(t *testing.T) {
	client, err := createTestClient(debugAPI)
	if err != nil {
		t.Errorf("Error creating test client %v", err)
	}
	backups, err := client.GetInstanceBackups(6809519)
	if err != nil {
		t.Errorf("Error listing backups, expected struct, got error %v", err)
	}
	if len(backups.Automatic) > 0 {
		t.Errorf("Expected an empty list of automatic backups, but got %v", backups.Automatic)
	}
	if backups.Snapshot.Current != nil {
		t.Errorf("Expected empty current snapshot, but got %v", backups.Snapshot.Current)
	}
}
