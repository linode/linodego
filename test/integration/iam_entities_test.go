package integration

import (
	"context"
	"testing"
)

func TestIAM_ListEntities(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestIAM_ListEntities")
	defer teardown()

	entities, err := client.ListEntities(context.Background(), nil)
	if err != nil {
		t.Fatalf("Error listing entities: %v", err)
	}

	if entities == nil {
		t.Fatal("Expected list of entities, got nil")
	}

	if len(entities) == 0 {
		t.Fatal("Expected one or more entities, got none")
	}

	for _, e := range entities {
		if e.ID == 0 {
			t.Errorf("Expected entity ID to be non-zero, got 0")
		}
		if e.Label == "" {
			t.Errorf("Expected entity label to be non-empty for entity ID %d", e.ID)
		}
		if e.Type == "" {
			t.Errorf("Expected entity type to be non-empty for entity ID %d", e.ID)
		}
	}
}
