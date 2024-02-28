package integration

import (
	"context"
	"testing"

	"github.com/linode/linodego"
)

func TestAccountBetaPrograms(t *testing.T) {
	optInTest(t)

	client, teardown := createTestClient(t, "fixtures/TestAccountBetaPrograms")
	defer teardown()

	betas, err := client.ListBetaPrograms(context.Background(), linodego.NewListOptions(1, ""))

	if len(betas) == 0 {
		t.Log("No beta program is available during the test.")
		return
	}
	createOpts := linodego.AccountBetaProgramCreateOpts{ID: betas[0].ID}

	_, err = client.JoinBetaProgram(context.Background(), createOpts)
	if err != nil {
		t.Errorf("Error joining a Beta program, expected struct, got error %v", err)
	}

	accountBetas, err := client.ListAccountBetaPrograms(context.Background(), &linodego.ListOptions{})
	if err != nil {
		t.Errorf("Error getting Account Beta programs, expected struct, got error %v", err)
	}

	if len(accountBetas) == 0 {
		t.Errorf("Expected to see account beta program returned.")
	} else {
		assertDateSet(t, accountBetas[0].Enrolled)
		betaID := accountBetas[0].ID
		beta, err := client.GetAccountBetaProgram(context.Background(), betaID)
		if err != nil {
			t.Errorf("Error getting an Account Beta program, expected struct, got error %v", err)
		}
		if beta.ID != betaID {
			t.Errorf("expected beta ID to be %s; got %s", betaID, beta.ID)
		}
	}
}
