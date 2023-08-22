package integration

import (
	"context"
	"testing"

	"github.com/linode/linodego"
)

func TestAccountBetaPrograms_List(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestAccountBetaPrograms_List")
	defer teardown()

	betas, err := client.ListAccountBetaPrograms(context.Background(), &linodego.ListOptions{})
	if err != nil {
		t.Errorf("Error getting Account Beta programs, expected struct, got error %v", err)
	}

	if len(betas) == 0 {
		t.Errorf("Expected to see account beta program returned.")
	} else {
		assertDateSet(t, betas[0].Enrolled)
	}
}

func TestAccountBetaProgram_Get(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestAccountBetaProgram_Get")
	defer teardown()

	betaID := "cool-beta"

	// Enroll the account into a beta program.
	createOpts := linodego.AccountBetaProgramCreateOpts{ID: betaID}

	_, err := client.JoinBetaProgram(context.Background(), createOpts)
	if err != nil {
		t.Errorf("Error joining a Beta program, expected struct, got error %v", err)
	}

	beta, err := client.GetAccountBetaProgram(context.Background(), betaID)

	if err != nil {
		t.Errorf("Error getting an Account Beta program, expected struct, got error %v", err)
	}

	if beta.ID != betaID {
		t.Errorf("expected beta ID to be %s; got %s", betaID, beta.ID)
	}

}
