package integration

import (
	"context"
	"testing"
)

func TestProfileLogins_List(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestProfileLogins_List")
	defer teardown()

	logins, err := client.ListProfileLogins(context.Background(), nil)

	if err != nil {
		t.Errorf("Error getting Profile Logins, expected struct, got error %v", err)
	}

	if len(logins) < 1 {
		t.Errorf("Expected to see at least one Profile Login")
	}

	login := logins[0]

	response, err := client.GetProfileLogin(context.Background(), login.ID)
	if err != nil {
		t.Errorf("Failed to get one Profile Login: %v", err)
	}

	if response.Username != login.Username {
		t.Fatal("Recieved Profile Login Username does not match source")
	}
}
