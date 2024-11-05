package integration

import (
	"context"
	"testing"

	"github.com/linode/linodego"
)

func TestAccountMaintenances_List(t *testing.T) {
	client, fixtureTeardown := createTestClient(t, "fixtures/TestAccountMaintenances_List")
	defer fixtureTeardown()

	listOpts := linodego.NewListOptions(0, "")
	_, err := client.ListMaintenances(context.Background(), listOpts)
	if err != nil {
		t.Errorf("Error listing maintenances, expected array, got error %v", err)
	}
}
