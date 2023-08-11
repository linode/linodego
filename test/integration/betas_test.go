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

	// TODO: assert data after actual data can be retrieved from API.
	// No data is expected to be returned temporarily.
	if len(betas) != 0 {
		t.Errorf("Expected to see none beta program returned.")
	}
}
