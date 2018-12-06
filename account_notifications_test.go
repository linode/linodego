package linodego_test

import (
	"context"
	"testing"

	"github.com/linode/linodego"
)

func TestListNotifications(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping test in short mode.")
	}
	client, fixtureTeardown := createTestClient(t, "fixtures/TestListNotifications")
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
