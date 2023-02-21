package integration

import (
	"context"
	"testing"
)

func TestAccount_Get(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestAccount_Get")
	defer teardown()

	account, err := client.GetAccount(context.Background())
	if err != nil {
		t.Errorf("Error getting Account, expected struct, got error %v", err)
	}

	if len(account.Email) < 1 {
		t.Error("Error accessing Account, expected Email")
	}
}
