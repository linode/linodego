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

func TestIAM_GetEntityRoles(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestIAM_GetEntityRoles")
	defer teardown()

	// Get current user
	profile, err := client.GetProfile(context.Background())
	if err != nil {
		t.Fatalf("Error getting account profile: %s", err)
	}

	entities, err := client.ListEntities(context.Background(), nil)
	if err != nil {
		t.Fatalf("Error listing entities: %s", err)
	}

	if len(entities) == 0 {
		t.Fatal("Expected at least one entity")
	}

	entity := entities[0]

	roles, err := client.GetEntityRoles(
		context.Background(),
		profile.Username,
		entity.Type,
		entity.ID,
	)
	if err != nil {
		t.Fatalf("Error getting entity roles for %s %d: %s",
			entity.Type, entity.ID, err)
	}

	if roles == nil {
		t.Fatal("Expected entity roles, got nil")
	}

	for _, role := range roles {
		if role == "" {
			t.Errorf("Expected role name to be non-empty")
		}
	}
}
