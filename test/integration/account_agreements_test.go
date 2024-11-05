package integration

import (
	"context"
	"testing"
)

func TestAccountAgreements_Get(t *testing.T) {
	client, fixtureTeardown := createTestClient(t, "fixtures/TestAccountAgreements_List")
	defer fixtureTeardown()

	_, err := client.GetAccountAgreements(context.Background())
	if err != nil {
		t.Errorf("Error getting agreements, expected struct, got error %v", err)
	}
}
