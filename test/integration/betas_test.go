package integration

import (
	"context"
	"testing"

	"github.com/linode/linodego"
)

func TestBetaPrograms_List(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestBetaPrograms_List")
	defer teardown()

	betas, err := client.ListBetaPrograms(context.Background(), &linodego.ListOptions{})
	if err != nil {
		t.Errorf("Error getting Beta programs, expected struct, got error %v", err)
	}

	if len(betas) == 0 {
		t.Errorf("Expected to see beta program returned.")
	} else {
		assertDateSet(t, betas[0].Started)
	}
}

func TestBetaProgram_Get(t *testing.T) {
	client, teardown := createTestClient(t, "fixtures/TestBetaProgram_Get")
	defer teardown()

	betaID := "active_closed"
	beta, err := client.GetBetaProgram(context.Background(), betaID)
	if err != nil {
		t.Errorf("Error getting Beta program, expected struct, got error %v", err)
	}

	if beta.ID != betaID {
		t.Errorf("expected beta ID to be %s; got %s", betaID, beta.ID)
	}
}
