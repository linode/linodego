package integration

import (
	"context"
	"testing"

	"github.com/linode/linodego"
)

func TestAccountNotifications_List(t *testing.T) {
	client, fixtureTeardown := createTestClient(t, "fixtures/TestAccountNotifications_List")
	defer fixtureTeardown()

	listOpts := linodego.NewListOptions(0, "")
	records, err := client.ListNotifications(context.Background(), listOpts)
	if err != nil {
		t.Errorf("Error listing notifications, expected array, got error %v", err)
	}
	if len(records) == 0 {
		t.Errorf("Expected ListNotifications to have some results")
	}
}
