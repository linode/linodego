package integration

import (
	"context"
	"testing"

	"github.com/linode/linodego"
)

func TestCustomerBetaPrograms_List(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestCustomerBetaPrograms_List")
	defer teardown()

	betas, err := client.ListCustomerBetaPrograms(context.Background(), &linodego.ListOptions{})
	if err != nil {
		t.Errorf("Error getting Customer Beta programs, expected struct, got error %v", err)
	}

	if len(betas) == 0 {
		t.Errorf("Expected to see customer beta program returned.")
	} else {
		assertDateSet(t, betas[0].Enrolled)
	}
}

func TestCustomerBetaProgram_Get(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestCustomerBetaProgram_Get")
	defer teardown()

	betaID := "cool-beta"

	// Enroll the customer into a beta program.
	createOpts := linodego.CustomerBetaProgramCreateOpts{ID: betaID}

	_, err := client.CreateCustomerBetaProgram(context.Background(), createOpts)
	if err != nil {
		t.Errorf("Error creating Customer Beta program, expected struct, got error %v", err)
	}

	beta, err := client.GetCustomerBetaProgram(context.Background(), betaID)

	if err != nil {
		t.Errorf("Error getting Customer Beta program, expected struct, got error %v", err)
	}

	if beta.ID != betaID {
		t.Errorf("expected beta ID to be %s; got %s", betaID, beta.ID)
	}

}
