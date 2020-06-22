package integration

import (
	"context"

	"github.com/linode/linodego/pkg/errors"

	"testing"
)

func TestGetUser_missing(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestGetUser_missing")
	defer teardown()

	i, err := client.GetUser(context.Background(), "does-not-exist")
	if err == nil {
		t.Errorf("should have received an error requesting a missing user, got %v", i)
	}
	e, ok := err.(*errors.Error)
	if !ok {
		t.Errorf("should have received an Error requesting a missing user, got %v", e)
	}

	if e.Code != 404 {
		t.Errorf("should have received a 404 Code requesting a missing user, got %v", e.Code)
	}
}

/**
// -----------------------------------------------------------
// TODO(displague)
// User creation/list/updates require User email confirmation.
// Testing will be revisited.
// -----------------------------------------------------------
func TestGetUser_found(t *testing.T) {
	client, _, teardown, err := setupUser(t, "fixtures/TestGetUser_found")
	defer teardown()
	if err != nil {
		t.Fatal(err)
	}

	i, err := client.GetUser(context.Background(), "linodego-testuser")
	if err != nil {
		t.Errorf("Error getting user, expected struct, got %v and error %v", i, err)
	}
	if i.Username != "linodego-testuser" {
		t.Errorf("Expected a specific user, but got a different one %v", i)
	}
}

func TestUpdateUser(t *testing.T) {
	client, user, teardown, err := setupUser(t, "fixtures/TestUpdateUser")
	defer teardown()
	if err != nil {
		t.Fatal(err)
	}

	createOpts := user.GetCreateOptions()
	if createOpts.Restricted != user.Restricted {
		t.Errorf("Expected Restricted to match GetCreateOptions, got: %v", createOpts)
	}

	updateOpts := user.GetUpdateOptions()
	if updateOpts.Email != user.Email {
		t.Errorf("Expected matching Email from GetUpdateOptions, got: %v", updateOpts)
	}

	updateOpts.Email = "r_" + user.Email
	user, err = client.UpdateUser(context.Background(), user.Username, updateOpts)
	if err != nil {
		t.Errorf("Error listing users, expected struct, got error %v", err)
	}

	if user.Email != updateOpts.Email {
		t.Errorf("Expected a change in user Email, but got none %v", user)
	}
}

func TestListUsers(t *testing.T) {
	client, user, teardown, err := setupUser(t, "fixtures/TestListUsers")
	defer teardown()
	if err != nil {
		t.Fatal(err)
	}

	listOpts := NewListOptions(1, "{\"username\":\""+user.Username+"\"}")
	i, err := client.ListUsers(context.Background(), listOpts)
	if err != nil {
		t.Errorf("Error listing users, expected struct, got error %v", err)
	}
	if len(i) == 0 {
		t.Errorf("Expected a list of users, but got none %v", i)
	}
}

func setupUser(t *testing.T, fixturesYaml string) (*Client, *User, func(), error) {
	t.Helper()
	client, fixtureTeardown := createTestClient(t, fixturesYaml)

	// This scope must be <= the scope used for testing
	username := "linodego-testuser"

	createOpts := UserCreateOptions{
		Username:   username,
		Email:      username + "@example.com",
		Restricted: true,
	}
	user, err := client.CreateUser(context.Background(), createOpts)
	if err != nil {
		t.Fatalf("Error creating test User: %s", err)
	}

	teardown := func() {
		if user != nil {
			if err := client.DeleteUser(context.Background(), user.Username); err != nil {
				t.Errorf("Error deleting test User: %s", err)
			}
		}
		fixtureTeardown()
	}
	return client, user, teardown, err
}
**/
