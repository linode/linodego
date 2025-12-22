package integration

import (
	"context"
	"testing"
)

func TestIAM_GetAccountRolePermissions(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestIAM_GetAccountRolePermissions")
	defer teardown()

	rolePermissions, err := client.GetAccountRolePermissions(context.Background())
	if err != nil {
		t.Errorf("Error getting account role permissions: %s", err)
	}

	if rolePermissions == nil {
		t.Fatal("Expected account role permissions, got nil")
	}

	if len(rolePermissions.AccountAccess) == 0 && len(rolePermissions.EntityAccess) == 0 {
		t.Errorf("Expected account or entity access permissions, got none")
	}
}

func TestIAM_GetUserRolePermissions(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestIAM_GetUserRolePermissions")
	defer teardown()

	account, err := client.GetProfile(context.Background())
	if err != nil {
		t.Fatalf("Error getting account profile: %s", err)
	}

	rolePermissions, err := client.GetUserRolePermissions(context.Background(), account.Username)
	if err != nil {
		t.Errorf("Error getting user role permissions for %s: %s", account.Username, err)
	}

	if rolePermissions == nil {
		t.Fatal("Expected user role permissions, got nil")
	}

	if rolePermissions.AccountAccess == nil {
		t.Errorf("Expected AccountAccess field, got nil")
	}
	if rolePermissions.EntityAccess == nil {
		t.Errorf("Expected EntityAccess field, got nil")
	}
}

func TestIAM_UpdateUserRolePermissions(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestIAM_UpdateUserRolePermissions")
	defer teardown()

	account, err := client.GetProfile(context.Background())
	if err != nil {
		t.Fatalf("Error getting account profile: %s", err)
	}

	username := account.Username

	currentPermissions, err := client.GetUserRolePermissions(context.Background(), username)
	if err != nil {
		t.Fatalf("Error getting user role permissions: %s", err)
	}

	updateOpts := currentPermissions.GetUpdateOptions()

	updatedPermissions, err := client.UpdateUserRolePermissions(context.Background(), username, updateOpts)
	if err != nil {
		t.Errorf("Error updating user role permissions: %s", err)
	}

	if updatedPermissions == nil {
		t.Fatal("Expected updated permissions, got nil")
	}

	if len(updatedPermissions.AccountAccess) != len(updateOpts.AccountAccess) {
		t.Errorf("Expected %d AccountAccess entries, got %d",
			len(updateOpts.AccountAccess), len(updatedPermissions.AccountAccess))
	}
}

func TestIAM_GetUserAccountPermissions(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestIAM_GetUserAccountPermissions")
	defer teardown()

	// Get current user profile to obtain username
	account, err := client.GetProfile(context.Background())
	if err != nil {
		t.Fatalf("Error getting account profile: %s", err)
	}

	permissions, err := client.GetUserAccountPermissions(
		context.Background(),
		account.Username,
	)
	if err != nil {
		t.Fatalf("Error getting user account permissions for %s: %s",
			account.Username, err)
	}

	if permissions == nil {
		t.Fatal("Expected account permissions, got nil")
	}

	if len(permissions) == 0 {
		t.Errorf("Expected one or more account permissions, got none")
	}

	for _, perm := range permissions {
		if perm == "" {
			t.Errorf("Expected permission string to be non-empty")
		}
	}
}
