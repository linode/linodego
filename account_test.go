package linodego

import (
	"strings"
	"testing"
)

func TestGetAccount(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping test in short mode.")
	}

	client, err := createTestClient(debugAPI)
	if err != nil {
		t.Errorf("Error creating test client %v", err)
	}
	account, err := client.GetAccount()
	if err != nil {
		t.Errorf("Error getting Account, expected struct, got error %v", err)
	}

	if !strings.Contains(account.Email, "@") {
		t.Error("Error accessing Account, expected Email to contain '@'")
	}
}
