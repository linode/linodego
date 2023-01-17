package integration

import (
	"context"
	"testing"
)

func TestAccountLogins_List(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestAccountLogins_List")
	defer teardown()

	logins, err := client.ListLogins(context.Background(), nil)

	if err != nil {
		t.Errorf("Error getting Account Logins, expected struct, got error %v", err)
	}

	if len(logins) < 1 {
		t.Errorf("Expected to see at least one Account Login")
	}

	login := logins[0]

	response, err := client.GetLogin(context.Background(), login.ID)
	if err != nil {
		t.Errorf("Failed to get one Account Login: %v", err)
	}

	if response.Username != login.Username {
		t.Fatal("Recieved Account Login Username does not match source")
	}
}
