package integration

import (
	"context"

	"github.com/linode/linodego"
	. "github.com/linode/linodego"

	"testing"
)

const usernamePrefix = "linodegotest-"

type userModifier func(*linodego.UserCreateOptions)

func TestGetUser_missing(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestGetUser_missing")
	defer teardown()

	i, err := client.GetUser(context.Background(), "does-not-exist")
	if err == nil {
		t.Errorf("should have received an error requesting a missing user, got %v", i)
	}
	e, ok := err.(*Error)
	if !ok {
		t.Errorf("should have received an Error requesting a missing user, got %v", e)
	}

	if e.Code != 404 {
		t.Errorf("should have received a 404 Code requesting a missing user, got %v", e.Code)
	}
}

func TestGetUser(t *testing.T) {
	username := usernamePrefix + "getuser"
	email := usernamePrefix + "getuser@example.com"
	restricted := true

	client, _, teardown := setupUser(t, []userModifier{
		func(createOpts *linodego.UserCreateOptions) {
			createOpts.Username = username
			createOpts.Email = email
			createOpts.Restricted = restricted
		},
	}, "fixtures/TestGetUser")
	defer teardown()

	user, err := client.GetUser(context.TODO(), username)
	if err != nil {
		t.Fatalf("failed to get user (%s): %s", username, err)
	}

	if user.Email != email {
		t.Errorf("expected user email to be %s; got %s", email, user.Email)
	}
	if len(user.SSHKeys) != 0 {
		t.Error("expected user to have no SSH keys")
	}
	if !user.Restricted {
		t.Error("expected user to be restricted")
	}
}

func TestUpdateUser(t *testing.T) {
	username := usernamePrefix + "updateuser"
	email := usernamePrefix + "updateuser@example.com"
	restricted := false

	client, user, teardown := setupUser(t, []userModifier{
		func(createOpts *linodego.UserCreateOptions) {
			createOpts.Username = username
			createOpts.Email = email
			createOpts.Restricted = restricted
		},
	}, "fixtures/TestUpdateUser")
	defer teardown()

	updatedUsername := username + "-updated"
	restricted = true
	updateOpts := UserUpdateOptions{
		Username:   updatedUsername,
		Restricted: &restricted,
	}

	updated, err := client.UpdateUser(context.TODO(), username, updateOpts)
	if err != nil {
		t.Fatalf("failed to update user (%s): %s", username, err)
	}
	// update username to be deleted in teardown
	user.Username = updatedUsername

	if updated.Username != updatedUsername {
		t.Errorf("expected username to be %s; got %s", updatedUsername, updated.Username)
	}
	if !updated.Restricted {
		t.Error("expected user to be restricted")
	}
}

func TestListUsers(t *testing.T) {
	username := usernamePrefix + "listuser"
	email := usernamePrefix + "listuser@example.com"
	restricted := false

	client, _, teardown := setupUser(t, []userModifier{
		func(createOpts *linodego.UserCreateOptions) {
			createOpts.Username = username
			createOpts.Email = email
			createOpts.Restricted = restricted
		},
	}, "fixtures/TestListUsers")
	defer teardown()

	users, err := client.ListUsers(context.TODO(), nil)
	if err != nil {
		t.Fatalf("failed to get users: %s", err)
	}

	if len(users) == 0 {
		t.Fatalf("expected at least one user to be returned")
	}

	var newUser User
	for _, user := range users {
		if user.Username == username {
			newUser = user
		}
	}

	if newUser.Email != email {
		t.Errorf("expected user email to be %s; got %s", email, newUser.Email)
	}
	if len(newUser.SSHKeys) != 0 {
		t.Error("expected user to have no SSH keys")
	}
	if newUser.Restricted {
		t.Error("expected user to not be restricted")
	}
}

func createUser(t *testing.T, client *linodego.Client, userModifiers ...userModifier) (*User, func()) {
	t.Helper()

	var createOpts UserCreateOptions
	for _, modifier := range userModifiers {
		modifier(&createOpts)
	}

	user, err := client.CreateUser(context.TODO(), createOpts)
	if err != nil {
		t.Fatalf("failed to create test user: %s", err)
	}

	return user, func() {
		if err := client.DeleteUser(context.TODO(), user.Username); err != nil {
			t.Errorf("failed to delete test user (%s): %s", user.Username, err)
		}
	}
}

func setupUser(t *testing.T, userModifiers []userModifier, fixturesYaml string) (*Client, *User, func()) {
	t.Helper()
	client, fixtureTeardown := createTestClient(t, fixturesYaml)

	user, teardown := createUser(t, client, userModifiers...)
	return client, user, func() {
		teardown()
		fixtureTeardown()
	}
}
